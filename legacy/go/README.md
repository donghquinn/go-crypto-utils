# Legacy Go / Fyne build (archived)

This is the original implementation, kept for reference. The app is now an
Electron project at the repository root — see the top-level `README.md`.

Nothing here is built by the current toolchain. It still compiles on its own:

```bash
cd legacy/go
go build -o bin/app-macos main.go
```

## Why it was replaced

Fyne requires CGO, which made cross-compilation painful: the Windows build
needed a MinGW-w64 toolchain on the build host, and there was no Linux target at
all. `electron-builder` produces macOS, Windows and Linux installers from a
single machine with no cross-toolchain.

## Wire compatibility

The Electron build reads ciphertext produced here. `src/main/crypto.js` is a
byte-for-byte port of `biz/`, and `src/main/operations.js` reproduces the
tab-level encoding from `gui/` — including the double Base64 wrapping that
`gui/encrypt.gui.go` applied to AES-CBC output.
