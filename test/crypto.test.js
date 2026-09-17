'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');

const c = require('../src/main/crypto');
const ops = require('../src/main/operations');

const KEY32 = Buffer.from('0'.repeat(64), 'hex');

test('pkcs7 pads to the block size and round-trips', () => {
  for (const len of [0, 1, 15, 16, 17, 32]) {
    const data = Buffer.alloc(len, 0xab);
    const padded = c.pkcs7Padding(data, 16);
    assert.equal(padded.length % 16, 0);
    assert.ok(padded.length > data.length);
    assert.deepEqual(Buffer.from(c.pkcs7Unpadding(padded)), data);
  }
});

test('pkcs7 unpadding rejects bad padding', () => {
  assert.throws(() => c.pkcs7Unpadding(Buffer.alloc(0)), /empty/);
  assert.throws(() => c.pkcs7Unpadding(Buffer.from([1, 2, 0])), /invalid padding/);
  assert.throws(() => c.pkcs7Unpadding(Buffer.from([1, 2, 9])), /invalid padding/);
  assert.throws(() => c.pkcs7Unpadding(Buffer.from([1, 3, 3])), /invalid padding/);
});

test('AES-CBC round-trips for every valid key size', () => {
  for (const size of [16, 24, 32]) {
    const key = c.generateRandomAesKey(size);
    const sealed = c.encryptAes256Cbc(Buffer.from('안녕하세요 secure'), key);
    assert.equal(c.decryptAes256Cbc(sealed, key).toString('utf8'), '안녕하세요 secure');
  }
});

test('AES-CBC emits base64(iv || ciphertext)', () => {
  const sealed = Buffer.from(c.encryptAes256Cbc(Buffer.from('x'), KEY32), 'base64');
  assert.equal(sealed.length, 16 + 16); // one IV + one padded block
});

test('AES-CBC rejects a short or misaligned ciphertext', () => {
  assert.throws(() => c.decryptAes256Cbc(Buffer.alloc(8).toString('base64'), KEY32), /too short/);
  assert.throws(() => c.decryptAes256Cbc(Buffer.alloc(24).toString('base64'), KEY32), /multiple of the block size/);
});

test('AES-GCM round-trips and emits nonce || ciphertext || tag', () => {
  const sealed = c.encryptAes256Gcm(Buffer.from('hello'), KEY32);
  assert.equal(sealed.length, c.GCM_NONCE_SIZE + 5 + c.GCM_TAG_SIZE);
  assert.equal(c.decryptAes256Gcm(sealed, KEY32).toString('utf8'), 'hello');
});

test('AES-GCM rejects tampered ciphertext', () => {
  const sealed = c.encryptAes256Gcm(Buffer.from('hello'), KEY32);
  sealed[c.GCM_NONCE_SIZE] ^= 0xff;
  assert.throws(() => c.decryptAes256Gcm(sealed, KEY32), /authentication failed/);
  assert.throws(() => c.decryptAes256Gcm(Buffer.alloc(4), KEY32), /too short/);
});

test('key generation rejects invalid lengths', () => {
  for (const size of [16, 24, 32]) {
    assert.equal(c.generateRandomAesKey(size).length, size);
  }
  assert.throws(() => c.generateRandomAesKey(20), /invalid key length/);
});

test('decodeKey prefers hex, falls back to base64, rejects garbage', () => {
  assert.deepEqual(c.decodeKey('00ff'), Buffer.from([0x00, 0xff]));
  assert.deepEqual(c.decodeKey('aGVsbG8='), Buffer.from('hello'));
  assert.throws(() => c.decodeKey('not a key!'), /Hexadecimal or Base64/);
});

test('strict decoders reject what Buffer.from would silently drop', () => {
  assert.throws(() => c.decodeHexStrict('0g'), /invalid byte/);
  assert.throws(() => c.decodeHexStrict('abc'), /odd length/);
  assert.throws(() => c.decodeBase64Strict('a!b='), /illegal base64/);
  assert.throws(() => c.decodeBase64Strict('abcde'), /illegal base64/);
});

test('random strings use the documented charset and length', () => {
  const s = c.generateCustomRandomString(128);
  assert.equal(s.length, 128);
  assert.match(s, /^[a-zA-Z0-9]+$/);
  assert.throws(() => c.generateCustomRandomString(0), /greater than 0/);
});

test('uuid is RFC 4122 v4', () => {
  assert.match(c.createUuid(), /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/);
});

test('hashData matches known digests', () => {
  assert.equal(
    c.hashData(Buffer.from('abc'), 'SHA-256').toString('hex'),
    'ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad',
  );
  assert.equal(c.hashData(Buffer.from('abc'), 'SHA-512').length, 64);
  assert.throws(() => c.hashData(Buffer.alloc(0), 'MD5'), /unsupported/);
});

// --- tab-level operations -------------------------------------------------

test('encrypt/decrypt round-trip through the tab operations', () => {
  const { hex, base64 } = ops.generateKey('AES-256 (32bytes)');
  for (const key of [hex, base64]) {
    for (const method of ['AES-CBC', 'AES-GCM']) {
      const sealed = ops.encrypt({ method, text: 'round trip', key });
      assert.equal(ops.decrypt({ method, data: sealed, key }), 'round trip');
    }
  }
});

test('AES-CBC keeps the legacy double base64 wire format', () => {
  const key = KEY32.toString('hex');
  const outer = ops.encrypt({ method: 'AES-CBC', text: 'x', key });
  // The outer layer decodes to an inner base64 string, not to raw bytes.
  const inner = Buffer.from(outer, 'base64').toString('utf8');
  assert.match(inner, /^[A-Za-z0-9+/]+={0,2}$/);
  assert.equal(Buffer.from(inner, 'base64').length, 32);
});

test('AES-GCM is encoded exactly once', () => {
  const key = KEY32.toString('hex');
  const outer = ops.encrypt({ method: 'AES-GCM', text: 'x', key });
  assert.equal(Buffer.from(outer, 'base64').length, c.GCM_NONCE_SIZE + 1 + c.GCM_TAG_SIZE);
});

test('operations validate empty input like the Fyne tabs did', () => {
  assert.throws(() => ops.encrypt({ method: 'AES-CBC', text: '', key: 'ff' }), /Please enter text\./);
  assert.throws(() => ops.encrypt({ method: 'AES-CBC', text: 'x', key: '' }), /Please enter an AES key\./);
  assert.throws(() => ops.decrypt({ method: 'AES-CBC', data: '', key: 'ff' }), /Please enter encrypted data\./);
  assert.throws(() => ops.encodeBase64(''), /Please enter text to encode\./);
  assert.throws(() => ops.decodeHex(''), /Please enter a Hex string to decode\./);
});

test('an invalid key size surfaces instead of producing empty output', () => {
  // The Fyne build discarded this error and showed a blank result.
  assert.throws(
    () => ops.encrypt({ method: 'AES-CBC', text: 'x', key: 'aabb' }),
    /invalid key size 2/,
  );
});

test('encode/decode operations round-trip', () => {
  assert.equal(ops.encodeBase64('hello'), 'aGVsbG8=');
  assert.equal(ops.decodeBase64('aGVsbG8='), 'hello');
  assert.equal(ops.encodeHex('hello'), '68656c6c6f');
  assert.equal(ops.decodeHex('68656c6c6f'), 'hello');
  assert.throws(() => ops.decodeBase64('!!!!'), /Error decoding Base64 string/);
  assert.throws(() => ops.decodeHex('zz'), /Error decoding Hex string/);
});

test('generateKey honours each key type and defaults to 256-bit', () => {
  assert.equal(Buffer.from(ops.generateKey('AES-128 (16bytes)').hex, 'hex').length, 16);
  assert.equal(Buffer.from(ops.generateKey('AES-192 (24bytes)').hex, 'hex').length, 24);
  assert.equal(Buffer.from(ops.generateKey('AES-256 (32bytes)').hex, 'hex').length, 32);
  assert.equal(Buffer.from(ops.generateKey('nonsense').hex, 'hex').length, 32);
});

test('generateKeyAs emits the requested representation', () => {
  assert.match(ops.generateKeyAs('Hexadecimal'), /^[0-9a-f]{64}$/);
  assert.match(ops.generateKeyAs('Base64'), /^[A-Za-z0-9+/]{43}=$/);
  assert.throws(() => ops.generateKeyAs(''), /key format/);
});
