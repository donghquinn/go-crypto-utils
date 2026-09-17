'use strict';

// Tab-level operations, ported from the Fyne `gui` package.
//
// Wire-format note: the Fyne Encryption tab base64-encoded the *already*
// base64-encoded output of EncryptAES256CBC, and the Decryption tab undid both
// layers. That double encoding is preserved here so ciphertext written by the
// Go build still decrypts in this one. AES-GCM was encoded exactly once.

const {
  encryptAes256Cbc,
  decryptAes256Cbc,
  encryptAes256Gcm,
  decryptAes256Gcm,
  generateRandomAesKey,
  genKey,
  generateCustomRandomString,
  createUuid,
  decodeHexStrict,
  decodeBase64Strict,
  decodeKey,
} = require('./crypto');

const METHODS = ['AES-CBC', 'AES-GCM'];

const KEY_SIZES = {
  'AES-128 (16bytes)': 16,
  'AES-192 (24bytes)': 24,
  'AES-256 (32bytes)': 32,
};

function required(value, fieldName) {
  if (value === '' || value == null) {
    throw new Error(`Please enter ${fieldName}.`);
  }
  return value;
}

function encrypt({ method, text, key: keyInput }) {
  required(text, 'text');
  required(keyInput, 'an AES key');
  if (!METHODS.includes(method)) {
    throw new Error(`Unsupported method: ${method}`);
  }

  const key = decodeKey(keyInput);
  const plaintext = Buffer.from(text, 'utf8');

  if (method === 'AES-CBC') {
    // legacy double encoding, see note above
    return Buffer.from(encryptAes256Cbc(plaintext, key), 'utf8').toString('base64');
  }
  return encryptAes256Gcm(plaintext, key).toString('base64');
}

function decrypt({ method, data, key: keyInput }) {
  required(data, 'encrypted data');
  required(keyInput, 'a valid AES key');
  if (!METHODS.includes(method)) {
    throw new Error(`Unsupported method: ${method}`);
  }

  const key = decodeKey(keyInput);

  let encrypted;
  try {
    encrypted = decodeBase64Strict(data);
  } catch {
    throw new Error('Invalid Base64 data.');
  }

  if (method === 'AES-CBC') {
    return decryptAes256Cbc(encrypted.toString('utf8'), key).toString('utf8');
  }
  return decryptAes256Gcm(encrypted, key).toString('utf8');
}

function encodeBase64(text) {
  required(text, 'text to encode');
  return Buffer.from(text, 'utf8').toString('base64');
}

function encodeHex(text) {
  required(text, 'text to encode');
  return Buffer.from(text, 'utf8').toString('hex');
}

function decodeBase64(text) {
  required(text, 'a Base64 string to decode');
  try {
    return decodeBase64Strict(text).toString('utf8');
  } catch (err) {
    throw new Error(`Error decoding Base64 string: ${err.message}`);
  }
}

function decodeHex(text) {
  required(text, 'a Hex string to decode');
  try {
    return decodeHexStrict(text).toString('utf8');
  } catch (err) {
    throw new Error(`Error decoding Hex string: ${err.message}`);
  }
}

function generateKey(keyType) {
  const length = KEY_SIZES[keyType] ?? 32;
  return genKey(generateRandomAesKey(length));
}

// The "Generate New Key" dialog on the Encryption tab: always AES-256, emitted
// in the requested representation.
function generateKeyAs(format) {
  const { hex, base64 } = genKey(generateRandomAesKey(32));
  if (format === 'Hexadecimal') return hex;
  if (format === 'Base64') return base64;
  throw new Error('Please enter a key format.');
}

function randomString(length) {
  return generateCustomRandomString(length);
}

function uuid() {
  return createUuid();
}

module.exports = {
  METHODS,
  KEY_SIZES,
  encrypt,
  decrypt,
  encodeBase64,
  encodeHex,
  decodeBase64,
  decodeHex,
  generateKey,
  generateKeyAs,
  randomString,
  uuid,
};
