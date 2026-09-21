#!/usr/bin/env node
/**
 * Build releases/latest.json for the download site.
 *
 * Usage: node scripts/build-manifest.mjs <dir> <version> <publicBaseUrl>
 *
 * Walks <dir> for installer artifacts, classifies each by filename (the names
 * come from electron-builder.yml), and prints the manifest to stdout. The site
 * fetches this file so a new release shows up without redeploying Pages.
 */

import { readdirSync, statSync } from "node:fs";
import { join, relative } from "node:path";

const [dir, version, baseUrlRaw] = process.argv.slice(2);

if (!dir || !version || !baseUrlRaw) {
  console.error("usage: build-manifest.mjs <dir> <version> <publicBaseUrl>");
  process.exit(1);
}

const baseUrl = baseUrlRaw.replace(/\/+$/, "");
const prefix = `v${version}`;

/** Ordered: the first pattern that matches a filename wins. */
const RULES = [
  { re: /-arm64\.dmg$/i,            os: "mac",   arch: "arm64",  kind: "Apple silicon · .dmg" },
  { re: /-x64\.dmg$/i,              os: "mac",   arch: "x64",    kind: "Intel · .dmg" },
  { re: /\.dmg$/i,                  os: "mac",   arch: "x64",    kind: ".dmg" },
  { re: /-arm64-mac\.zip$/i,        os: "mac",   arch: "arm64",  kind: "Apple silicon · .zip" },
  { re: /-mac\.zip$/i,              os: "mac",   arch: "x64",    kind: "Intel · .zip" },

  { re: /-arm64-setup\.exe$/i,      os: "win",   arch: "arm64",  kind: "Installer" },
  { re: /-x64-setup\.exe$/i,        os: "win",   arch: "x64",    kind: "Installer" },
  { re: /-portable\.exe$/i,         os: "win",   arch: "x64",    kind: "Portable" },
  { re: /setup\.exe$/i,             os: "win",   arch: "x64",    kind: "Installer" },

  { re: /-arm64\.AppImage$/i,       os: "linux", arch: "arm64",  kind: "AppImage" },
  { re: /\.AppImage$/i,             os: "linux", arch: "x86_64", kind: "AppImage" },
  { re: /_arm64\.deb$/i,            os: "linux", arch: "arm64",  kind: "Debian package" },
  { re: /_amd64\.deb$/i,            os: "linux", arch: "amd64",  kind: "Debian package" },
  { re: /-arm64\.tar\.gz$/i,        os: "linux", arch: "arm64",  kind: "Tarball" },
  { re: /\.tar\.gz$/i,              os: "linux", arch: "x64",    kind: "Tarball" },
];

/** Sort key: platform, then the arch people are most likely to want. */
const OS_ORDER = { mac: 0, win: 1, linux: 2 };
const KIND_ORDER = [
  "Apple silicon · .dmg", "Intel · .dmg", ".dmg",
  "Apple silicon · .zip", "Intel · .zip",
  "Installer", "Portable",
  "AppImage", "Debian package", "Tarball",
];

/** Anything else in the directory (unpacked trees, debug logs) is not a download. */
const DISTRIBUTABLE = /\.(dmg|zip|exe|AppImage|deb|rpm|pkg|snap|tar\.gz|tar\.xz)$/i;

function* walk(root) {
  for (const entry of readdirSync(root, { withFileTypes: true })) {
    if (entry.name.endsWith("-unpacked") || entry.name.startsWith(".")) continue;
    const full = join(root, entry.name);
    if (entry.isDirectory()) yield* walk(full);
    else if (entry.isFile() && DISTRIBUTABLE.test(entry.name)) yield full;
  }
}

const files = [];
const skipped = [];

for (const path of [...walk(dir)].sort()) {
  const name = path.split("/").pop();
  if (/\.blockmap$/i.test(name)) continue;

  const rule = RULES.find((r) => r.re.test(name));
  if (!rule) {
    skipped.push(relative(dir, path));
    continue;
  }

  // A matrix build can produce the same artifact on more than one runner.
  if (files.some((f) => f.name === name)) continue;

  files.push({
    name,
    os: rule.os,
    arch: rule.arch,
    kind: rule.kind,
    size: statSync(path).size,
    url: `${baseUrl}/${prefix}/${encodeURIComponent(name)}`,
  });
}

files.sort((a, b) => {
  const os = OS_ORDER[a.os] - OS_ORDER[b.os];
  if (os !== 0) return os;
  const kind = KIND_ORDER.indexOf(a.kind) - KIND_ORDER.indexOf(b.kind);
  if (kind !== 0) return kind;
  return a.name.localeCompare(b.name);
});

if (files.length === 0) {
  console.error(`no installer artifacts found under ${dir}`);
  process.exit(1);
}

if (skipped.length) {
  console.error(`unclassified (not listed on the site): ${skipped.join(", ")}`);
}

process.stdout.write(
  JSON.stringify(
    {
      version,
      tag: `v${version}`,
      baseUrl,
      prefix,
      count: files.length,
      files,
    },
    null,
    2,
  ) + "\n",
);
