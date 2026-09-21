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

/**
 * Ordered: the first pattern that matches a filename wins.
 *
 * electron-builder.yml builds exactly three target types, x64 and arm64 each.
 * Anything not matched here is not offered on the download site.
 */
const RULES = [
  { re: /-arm64\.dmg$/i,       os: "mac",   arch: "arm64",  kind: "Apple silicon" },
  { re: /-x64\.dmg$/i,         os: "mac",   arch: "x64",    kind: "Intel" },

  { re: /-x64-setup\.exe$/i,   os: "win",   arch: "x64",    kind: "Installer" },
  { re: /-arm64-setup\.exe$/i, os: "win",   arch: "arm64",  kind: "Installer (ARM)" },

  { re: /-arm64\.AppImage$/i,  os: "linux", arch: "arm64",  kind: "AppImage (ARM)" },
  { re: /\.AppImage$/i,        os: "linux", arch: "x86_64", kind: "AppImage" },
];

/** Sort key: platform, then the build people are most likely to want. */
const OS_ORDER = { mac: 0, win: 1, linux: 2 };
const KIND_ORDER = [
  "Apple silicon", "Intel",
  "Installer", "Installer (ARM)",
  "AppImage", "AppImage (ARM)",
];

/** Anything else in the directory (checksums, the manifest itself) is not a download. */
const DISTRIBUTABLE = /\.(dmg|exe|AppImage)$/i;

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
