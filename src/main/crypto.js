'use strict';

// Port of the Go `biz` package (biz/aes.go, biz/pkcs.go, biz/key.go,
// biz/random.go, biz/sha.go, biz/uuid.go) onto node:crypto.
//
// The byte-level output is intentionally identical to the Go implementation so
// that data produced by the Fyne version of this app can still be read here.

const crypto = require('node:crypto');

const GCM_NONCE_SIZE = 12; // Go's cipher.NewGCM default nonce size
const GCM_TAG_SIZE = 16;

// --- PKCS#7 (biz/pkcs.go) -------------------------------------------------

function pkcs7Padding(data, blockSize) {
  const padding = blockSize - (data.length % blockSize);
  return Buffer.concat([data, Buffer.alloc(padding, padding)]);
}

function pkcs7Unpadding(data) {
  const length = data.length;
  if (length === 0) {
    throw new Error('data is empty');
  }
  const padding = data[length - 1];
  if (padding > length || padding === 0) {
    throw new Error('invalid padding');
  }
  for (const v of data.subarray(length - padding)) {
    if (v !== padding) {
      throw new Error('invalid padding');
    }
  }
  return data.subarray(0, length - padding);
}

// --- Key helpers ----------------------------------------------------------

// Go's aes.NewCipher selects AES-128/192/256 from the key length; node needs
// the variant spelled out in the algorithm name.
function cipherName(key, mode) {
  switch (key.length) {
    case 16:
    case 24:
    case 32:
      return `aes-${key.length * 8}-${mode}`;
    default:
      throw new Error(`crypto/aes: invalid key size ${key.length}`);
  }
}

// --- AES (biz/aes.go) -----------------------------------------------------

// EncryptAES256CBC: returns base64(iv || ciphertext).
function encryptAes256Cbc(plaintext, key) {
  const name = cipherName(key, 'cbc');
  const iv = crypto.randomBytes(16);
  const cipher = crypto.createCipheriv(name, key, iv);
  cipher.setAutoPadding(false); // padding applied explicitly, as in Go
  const padded = pkcs7Padding(plaintext, 16);
  const ciphertext = Buffer.concat([cipher.update(padded), cipher.final()]);
  return Buffer.concat([iv, ciphertext]).toString('base64');
}

// DecryptAES256CBC: takes base64(iv || ciphertext).
function decryptAes256Cbc(ciphertextBase64, key) {
  const withIv = decodeBase64Strict(ciphertextBase64);
  const name = cipherName(key, 'cbc');

  if (withIv.length < 16) {
    throw new Error('ciphertext is too short');
  }
  const iv = withIv.subarray(0, 16);
  const ciphertext = withIv.subarray(16);
  if (ciphertext.length % 16 !== 0) {
    throw new Error('ciphertext length is not a multiple of the block size');
  }

  const decipher = crypto.createDecipheriv(name, key, iv);
  decipher.setAutoPadding(false);
  const padded = Buffer.concat([decipher.update(ciphertext), decipher.final()]);
  return pkcs7Unpadding(padded);
}

// EncryptAES256GCM: returns nonce || ciphertext || tag (Go's gcm.Seal layout).
function encryptAes256Gcm(plaintext, key) {
  const name = cipherName(key, 'gcm');
  const nonce = crypto.randomBytes(GCM_NONCE_SIZE);
  const cipher = crypto.createCipheriv(name, key, nonce, {
    authTagLength: GCM_TAG_SIZE,
  });
  const ciphertext = Buffer.concat([cipher.update(plaintext), cipher.final()]);
  return Buffer.concat([nonce, ciphertext, cipher.getAuthTag()]);
}

function decryptAes256Gcm(sealed, key) {
  const name = cipherName(key, 'gcm');
  if (sealed.length < GCM_NONCE_SIZE + GCM_TAG_SIZE) {
    throw new Error('ciphertext too short');
  }
  const nonce = sealed.subarray(0, GCM_NONCE_SIZE);
  const tag = sealed.subarray(sealed.length - GCM_TAG_SIZE);
  const ciphertext = sealed.subarray(GCM_NONCE_SIZE, sealed.length - GCM_TAG_SIZE);

  const decipher = crypto.createDecipheriv(name, key, nonce, {
    authTagLength: GCM_TAG_SIZE,
  });
  decipher.setAuthTag(tag);
  try {
    return Buffer.concat([decipher.update(ciphertext), decipher.final()]);
  } catch {
    // Mirrors Go's opaque cipher: message authentication failed.
    throw new Error('cipher: message authentication failed');
  }
}

// --- Keys (biz/key.go) ----------------------------------------------------

function generateRandomAesKey(length) {
  if (length !== 16 && length !== 24 && length !== 32) {
    throw new Error('invalid key length: must be 16, 24, or 32 bytes');
  }
  return crypto.randomBytes(length);
}

// GenKey: hex and base64 representations of a key.
function genKey(key) {
  return { hex: key.toString('hex'), base64: key.toString('base64') };
}

// --- Random strings (biz/random.go) ---------------------------------------

const CHAR_SET =
  'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';

function generateCustomRandomString(length) {
  if (!Number.isInteger(length) || length <= 0) {
    throw new Error('length must be greater than 0');
  }
  let result = '';
  for (let i = 0; i < length; i += 1) {
    result += CHAR_SET[crypto.randomInt(CHAR_SET.length)];
  }
  return result;
}

// --- UUID (biz/uuid.go) ---------------------------------------------------

function createUuid() {
  return crypto.randomUUID(); // RFC 4122 v4, same shape as google/uuid
}

// --- Hashing (biz/sha.go) -------------------------------------------------

function hashData(data, method) {
  switch (method) {
    case 'SHA-256':
      return crypto.createHash('sha256').update(data).digest();
    case 'SHA-512':
      return crypto.createHash('sha512').update(data).digest();
    default:
      throw new Error('unsupported hash method');
  }
}

// --- Strict decoders ------------------------------------------------------
//
// Buffer.from() silently drops invalid characters; Go's encoding/hex and
// encoding/base64 reject them. Validate first so bad input surfaces as an
// error instead of as silently mangled bytes.

function decodeHexStrict(value) {
  if (value.length % 2 !== 0) {
    throw new Error('encoding/hex: odd length hex string');
  }
  if (!/^[0-9a-fA-F]*$/.test(value)) {
    throw new Error('encoding/hex: invalid byte in hex string');
  }
  return Buffer.from(value, 'hex');
}

function decodeBase64Strict(value) {
  if (value.length % 4 !== 0 || !/^[A-Za-z0-9+/]*={0,2}$/.test(value)) {
    throw new Error('illegal base64 data');
  }
  return Buffer.from(value, 'base64');
}

// DecodeKey (gui/build.gui.go): hex first, then base64.
function decodeKey(keyStr) {
  try {
    return decodeHexStrict(keyStr);
  } catch {
    /* fall through to base64 */
  }
  try {
    return decodeBase64Strict(keyStr);
  } catch {
    /* fall through to the shared error */
  }
  throw new Error('key must be in Hexadecimal or Base64 format');
}

module.exports = {
  GCM_NONCE_SIZE,
  GCM_TAG_SIZE,
  CHAR_SET,
  pkcs7Padding,
  pkcs7Unpadding,
  encryptAes256Cbc,
  decryptAes256Cbc,
  encryptAes256Gcm,
  decryptAes256Gcm,
  generateRandomAesKey,
  genKey,
  generateCustomRandomString,
  createUuid,
  hashData,
  decodeHexStrict,
  decodeBase64Strict,
  decodeKey,
};
