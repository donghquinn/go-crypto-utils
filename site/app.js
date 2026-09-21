/* Secure Helper site — download links, OS detection, screenshot switcher. */
(function () {
  "use strict";

  var root = document.documentElement;
  var version = root.dataset.version;
  var base = (root.dataset.downloads || "").replace(/\/+$/, "");

  // The Pages workflow substitutes these. If it didn't run, fall back to the
  // GitHub release page rather than linking at a literal placeholder.
  var unresolved = /^__[A-Z_]+__$/;
  if (unresolved.test(version)) version = "";
  if (unresolved.test(base)) base = "";

  var RELEASES = "https://github.com/donghquinn/go-crypto-utils/releases";

  /* ---------- download table ---------- */

  // Artifact names follow electron-builder.yml. `arch` is the label shown to
  // the reader; `file` is templated with the version.
  var ARTIFACTS = [
    { os: "mac", kind: "Apple silicon · .dmg", arch: "arm64", file: "Secure Helper-{v}-arm64.dmg" },
    { os: "mac", kind: "Intel · .dmg", arch: "x64", file: "Secure Helper-{v}-x64.dmg" },
    { os: "mac", kind: "Apple silicon · .zip", arch: "arm64", file: "Secure Helper-{v}-arm64-mac.zip" },
    { os: "mac", kind: "Intel · .zip", arch: "x64", file: "Secure Helper-{v}-mac.zip" },

    { os: "win", kind: "Installer", arch: "x64", file: "Secure Helper-{v}-x64-setup.exe" },
    { os: "win", kind: "Installer", arch: "arm64", file: "Secure Helper-{v}-arm64-setup.exe" },
    { os: "win", kind: "Portable", arch: "x64", file: "Secure Helper-{v}-portable.exe" },

    { os: "linux", kind: "AppImage", arch: "x86_64", file: "Secure Helper-{v}-x86_64.AppImage" },
    { os: "linux", kind: "AppImage", arch: "arm64", file: "Secure Helper-{v}-arm64.AppImage" },
    { os: "linux", kind: "Debian package", arch: "amd64", file: "go-crypto-utils_{v}_amd64.deb" },
    { os: "linux", kind: "Debian package", arch: "arm64", file: "go-crypto-utils_{v}_arm64.deb" },
    { os: "linux", kind: "Tarball", arch: "x64", file: "go-crypto-utils-{v}.tar.gz" },
    { os: "linux", kind: "Tarball", arch: "arm64", file: "go-crypto-utils-{v}-arm64.tar.gz" }
  ];

  function urlFor(file) {
    if (!base || !version) return RELEASES;
    return base + "/v" + version + "/" + encodeURIComponent(file);
  }

  function render(items) {
    ["mac", "win", "linux"].forEach(function (os) {
      var list = document.querySelector('[data-list="' + os + '"]');
      if (!list) return;
      list.textContent = "";

      items.filter(function (a) { return a.os === os; }).forEach(function (a) {
        var li = document.createElement("li");
        var link = document.createElement("a");
        link.href = a.url || urlFor(a.file.replace("{v}", version));
        link.setAttribute("download", "");

        var label = document.createElement("span");
        label.textContent = a.kind;
        var arch = document.createElement("span");
        arch.className = "arch";
        arch.textContent = a.arch;

        link.append(label, arch);
        li.append(link);
        list.append(li);
      });
    });
  }

  /* ---------- which platform is this ---------- */

  function detect() {
    var ua = navigator.userAgent;
    var plat = (navigator.userAgentData && navigator.userAgentData.platform) ||
      navigator.platform || "";
    var hay = (plat + " " + ua).toLowerCase();

    if (/android/.test(hay)) return null;
    if (/iphone|ipad|ipod/.test(hay)) return null;
    if (/win/.test(hay)) return "win";
    if (/mac|darwin/.test(hay)) return "mac";
    if (/linux|x11|cros/.test(hay)) return "linux";
    return null;
  }

  // Apple silicon can't be read from the UA, so offer the arm64 build (the
  // common case since 2020) and leave the Intel build one click away.
  var PRIMARY = {
    mac: { label: "Download for macOS", note: "Apple silicon · .dmg", file: "Secure Helper-{v}-arm64.dmg" },
    win: { label: "Download for Windows", note: "64-bit installer", file: "Secure Helper-{v}-x64-setup.exe" },
    linux: { label: "Download for Linux", note: "x86_64 · AppImage", file: "Secure Helper-{v}-x86_64.AppImage" }
  };

  function setPrimary(items) {
    var os = detect();
    var btn = document.getElementById("primary-dl");
    if (!btn) return;

    var labelEl = btn.querySelector("[data-primary-label]");
    var noteEl = btn.querySelector("[data-primary-note]");

    if (!os) {
      labelEl.textContent = "Download";
      noteEl.textContent = "macOS, Windows or Linux";
      btn.href = "#download";
      return;
    }

    var want = PRIMARY[os];
    var name = want.file.replace("{v}", version);
    var match = items.filter(function (a) {
      return a.file && a.file.replace("{v}", version) === name;
    })[0];

    labelEl.textContent = want.label;
    noteEl.textContent = want.note;
    btn.href = (match && match.url) || urlFor(name);
    if (btn.href !== RELEASES) btn.setAttribute("download", "");

    var card = document.querySelector('.dl-card[data-os="' + os + '"]');
    if (card) card.classList.add("is-yours");
  }

  /* ---------- live manifest ---------- */

  // CI writes latest.json to the bucket root alongside each v<version>/
  // directory. If it is readable (bucket CORS permitting), it wins over the
  // built-in list — a new release shows up without redeploying the site.
  function refresh() {
    if (!base) return Promise.resolve(null);
    return fetch(base + "/latest.json", { cache: "no-cache" })
      .then(function (r) { return r.ok ? r.json() : null; })
      .catch(function () { return null; });
  }

  function fromManifest(manifest) {
    if (!manifest || !Array.isArray(manifest.files) || !manifest.files.length) return null;
    var items = manifest.files.filter(function (f) { return f.os && f.url; });
    return items.length ? items : null;
  }

  /* ---------- screenshots ---------- */

  function screenshots() {
    var tabs = document.querySelectorAll(".shot-tabs button");
    var img = document.getElementById("shot-img");
    var cap = document.getElementById("shot-cap");
    if (!tabs.length || !img) return;

    tabs.forEach(function (tab) {
      tab.addEventListener("click", function () {
        tabs.forEach(function (t) { t.setAttribute("aria-selected", "false"); });
        tab.setAttribute("aria-selected", "true");
        img.src = "assets/screenshots/" + tab.dataset.shot + ".png";
        img.alt = tab.textContent + " tab";
        cap.textContent = tab.dataset.cap;
      });
    });
  }

  /* ---------- boot ---------- */

  document.querySelectorAll("[data-version-text]").forEach(function (el) {
    el.textContent = version || "latest";
  });

  render(ARTIFACTS);
  setPrimary(ARTIFACTS);
  screenshots();

  refresh().then(function (manifest) {
    var items = fromManifest(manifest);
    if (!items) return;
    if (manifest.version) {
      version = manifest.version;
      document.querySelectorAll("[data-version-text]").forEach(function (el) {
        el.textContent = version;
      });
    }
    render(items);
    setPrimary(items);
  });
})();
