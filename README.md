# Go Crypto Utils

A cross-platform GUI application for cryptographic operations, key generation, and encoding/decoding utilities built with Go and Fyne.

## Maintainer: donghquinn

## Overview

Go Crypto Utils is a desktop application that provides a user-friendly interface for common cryptographic operations. It eliminates the need for command-line tools or writing custom scripts for encryption, decryption, encoding, and key generation tasks.

### Key Features

- **Encryption/Decryption**: AES-256-CBC and AES-256-GCM algorithms
- **Hashing**: SHA-256 and SHA-512 support
- **Encoding/Decoding**: Base64 and Hexadecimal
- **Key Generation**: Secure AES key generation (128, 192, 256-bit)
- **UUID Generation**: RFC 4122 compliant UUIDs
- **Random String Generation**: Customizable random string generation
- **Cross-platform**: Supports macOS and Windows

## Requirements

- Go 1.23.3 or later
- Fyne v2.5.2+ framework

## Installation

### Prerequisites

1. Install Go from [golang.org](https://golang.org/download/)
2. Install Fyne command-line tools:
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```
3. Add Go binary path to your environment:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

### Building from Source

#### macOS
```bash
# Using build script
./build-macos.sh

# Or using Makefile
make macos
```

#### Windows (Cross-compilation from macOS)
1. Install MinGW-w64:
   ```bash
   brew install mingw-w64
   ```

2. Set environment variables:
   ```bash
   export GOOS=windows
   export GOARCH=amd64
   export CGO_ENABLED=1
   export CC=x86_64-w64-mingw32-gcc
   ```

3. Build:
   ```bash
   # Using build script
   ./build-windows64.sh
   
   # Or using Makefile
   make windows
   ```

#### Using Makefile (Build All)
```bash
make all  # Builds both Windows and macOS versions
```

## Usage

Launch the application and use the tabbed interface to access different features:

### Available Tabs

1. **Encrypt** - Encrypt plaintext using AES algorithms
   - AES-256-CBC
   - AES-256-GCM

2. **Decrypt** - Decrypt ciphertext using AES algorithms
   - AES-256-CBC
   - AES-256-GCM

3. **Encode** - Encode data to different formats
   - Base64 encoding
   - Hexadecimal encoding

4. **Decode** - Decode data from different formats
   - Base64 decoding
   - Hexadecimal decoding

5. **Key Generation** - Generate secure AES keys
   - AES-128 (16 bytes)
   - AES-192 (24 bytes)
   - AES-256 (32 bytes)
   - Output in both hexadecimal and Base64 formats

6. **UUID Generation** - Generate RFC 4122 compliant UUIDs

7. **Random String** - Generate random strings with customizable parameters

## Project Structure

```
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── Makefile             # Build automation
├── biz/                 # Business logic package
│   ├── aes.go          # AES encryption/decryption
│   ├── key.go          # Key generation utilities
│   ├── pkcs.go         # PKCS padding functions
│   ├── random.go       # Random string generation
│   ├── sha.go          # SHA hashing functions
│   └── uuid.go         # UUID generation
├── gui/                 # GUI components
│   ├── common.go       # Shared GUI utilities
│   ├── encrypt.gui.go  # Encryption tab
│   ├── decrypt.gui.go  # Decryption tab
│   ├── encode.gui.go   # Encoding tab
│   ├── decode.gui.go   # Decoding tab (includes base64 and hex)
│   ├── key.gui.go      # Key generation tab
│   ├── uuid.gui.go     # UUID generation tab
│   └── random.gui.go   # Random string tab
├── screenshots/         # Application screenshots
└── icon.png            # Application icon
```

## Dependencies

- **fyne.io/fyne/v2**: Cross-platform GUI framework
- **github.com/google/uuid**: UUID generation library

## Screenshots

![Encryption](screenshots/encrypt.png)
![Decryption](screenshots/decrypt.png)
![Encoding](screenshots/encode.png)
![Decoding](screenshots/decode.png)
![Key Generation](screenshots/keygen.png)
![UUID Generation](screenshots/uuid.png)
![Random String](screenshots/randomstring.png)

## License

See [LICENSE](LICENSE) file for details.