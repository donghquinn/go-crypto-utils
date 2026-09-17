# Secure Helper

A cross-platform desktop app for everyday cryptographic chores — AES
encryption/decryption, Base64 and Hex encoding, AES key generation, UUIDs and
random strings. No command line, no throwaway scripts.

**Maintainer:** donghquinn

Built with Electron. The original Go/Fyne implementation is archived under
[`legacy/go/`](legacy/go/README.md).

## Features

- **Encryption / Decryption** — AES-CBC and AES-GCM (128/192/256-bit keys)
- **Encoding / Decoding** — Base64 and Hexadecimal
- **Key Generation** — AES-128, AES-192, AES-256, shown in Hex and Base64
- **UUID Generation** — RFC 4122 v4
- **Random Strings** — 1–128 characters, `[a-zA-Z0-9]`
- **Cross-platform** — macOS, Windows and Linux installers from one build host

All cryptography runs in the main process on Node's `node:crypto`; nothing
leaves the machine and the app makes no network requests.

## Requirements

- Node.js 20+ and npm (development only — end users just run the installer)

## Development

```bash
npm install
npm start     # launch the app
npm test      # unit tests for the crypto and operations layers
```

## Building installers

```bash
npm run dist          # macOS + Windows + Linux
npm run dist:mac      # .dmg and .zip        (x64, arm64)
npm run dist:win      # NSIS setup, portable (x64, arm64)
npm run dist:linux    # AppImage, .deb, .tar.gz (x64, arm64)
```

Artifacts land in `dist/`. All three platforms build from macOS — no MinGW,
Wine or Docker needed. CI (`.github/workflows/build.yml`) also builds each
platform on its native runner and attaches the artifacts to a `v*` tag release.

### Code signing

Builds are unsigned. macOS users get a Gatekeeper warning on first launch
(right-click → Open, or `xattr -dr com.apple.quarantine "/Applications/Secure Helper.app"`);
Windows shows a SmartScreen prompt. To sign, set `mac.identity` in
`electron-builder.yml` and provide `CSC_LINK` / `CSC_KEY_PASSWORD` (plus
`APPLE_ID`, `APPLE_APP_SPECIFIC_PASSWORD`, `APPLE_TEAM_ID` for notarization).

## Usage

| Tab | What it does |
| --- | --- |
| Encryption | Encrypt text with AES-CBC or AES-GCM; generate a key inline |
| Decryption | Decrypt Base64 ciphertext with the matching key |
| Encode | Text → Base64 or Hex |
| Decode | Base64 or Hex → text |
| Key Generation | Random AES-128/192/256 key, in Hex and Base64 |
| UUID | Generate a v4 UUID |
| Random String | Alphanumeric string, length set by a slider |

Keys are accepted as Hexadecimal or Base64. The input is tried as Hex first and
falls back to Base64, matching the previous build.

## Compatibility with the Go build

Ciphertext written by the Go/Fyne version decrypts here, and vice versa. This
was verified both directions against the original `biz` package across AES-128,
AES-192 and AES-256 with Hex and Base64 keys.

Two details are carried over deliberately:

- **AES-CBC output is Base64-encoded twice.** `EncryptAES256CBC` returned a
  Base64 string, and the old Encryption tab Base64-encoded it again; the
  Decryption tab undid both layers. Changing it would break existing
  ciphertext, so the format is preserved.
- **AES-GCM output is Base64-encoded once**, over `nonce || ciphertext || tag`.

One behaviour was **not** carried over: the old Encryption tab discarded the
error from `EncryptAES256CBC`, so an invalid key silently produced an empty
result. Errors now surface in the UI.

## Project structure

```
├── src/
│   ├── main/
│   │   ├── main.js          # app lifecycle, window, menu
│   │   ├── ipc.js           # IPC handler registration
│   │   ├── operations.js    # tab-level operations (port of gui/)
│   │   └── crypto.js        # AES, PKCS#7, keys, UUID, random (port of biz/)
│   ├── preload/preload.js   # contextBridge API over a fixed channel list
│   └── renderer/            # index.html, styles.css, renderer.js
├── test/crypto.test.js
├── build/icon.png           # 1024×1024 source for every platform icon
├── electron-builder.yml
├── legacy/go/               # archived Go/Fyne implementation
└── screenshots/
```

### Security posture

`contextIsolation` and `sandbox` are on, `nodeIntegration` is off, the renderer
has a restrictive CSP, navigation and new windows are blocked, and the preload
exposes only a fixed list of IPC channels. The renderer never touches key
material beyond passing it through to the main process.

## Screenshots

![Encryption](screenshots/encrypt.png)
![Decryption](screenshots/decrypt.png)
![Encoding](screenshots/encode.png)
![Decoding](screenshots/decode.png)
![Key Generation](screenshots/keygen.png)
![UUID Generation](screenshots/uuid.png)
![Random String](screenshots/randomstring.png)

## License

See [LICENSE](LICENSE).
