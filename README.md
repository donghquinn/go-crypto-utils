# Go Crypto Utils

## Maintainer: donghquinn

## What is Crypto Util?
* It's crypto utils written by golang.
* It must be quite bothering creating secure keys, or encoding strings just for test.
    * In my case, I felt lots of bothersome every single moment I had to encode request bodies or fields for API Encryption/Decryption test with base64 and AES / SHA / Etc...
    * Plus, Generating Random String on MacOs has been removed. I loved this one.
* This is a simple tool for generating keys or encoding strings

## Providings
* Encrypt
![Encryption](screenshots/encrypt.png)
    * AES-CBC
    * AES-GCM
    * SHA-256 (Upcoming)
    * SHA-512 (Upcoming)

* Decrypt
![Decryption](screenshots/decrypt.png)
    * AES-CBC
    * AES-GCM
    * SHA-256 (Upcoming)
    * SHA-512 (Upcoming)
    
* Encode
![Encode](screenshots/encode.png)
    * Base64
    * Hex

* Decode
![Decode](screenshots/decode.png)
    * Base64
    * Hex

* Generate AES Key
![Key Generation](screenshots/keygen.png)

* Generate UUID
![UUID Generation](screenshots/uuid.png)

* Generate Random String
![Random String](screenshots/randomstring.png)