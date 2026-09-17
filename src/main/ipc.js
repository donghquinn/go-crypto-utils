'use strict';

const { ipcMain, clipboard } = require('electron');
const ops = require('./operations');

// Every renderer request runs through here: the handler returns a plain
// {ok, value} / {ok:false, error} envelope so a thrown error becomes a message
// in the UI rather than an unhandled rejection.
function handle(channel, fn) {
  ipcMain.handle(channel, async (_event, payload) => {
    try {
      return { ok: true, value: await fn(payload) };
    } catch (err) {
      return { ok: false, error: err instanceof Error ? err.message : String(err) };
    }
  });
}

function registerHandlers() {
  handle('crypto:encrypt', (payload) => ops.encrypt(payload));
  handle('crypto:decrypt', (payload) => ops.decrypt(payload));
  handle('encode:base64', (payload) => ops.encodeBase64(payload.text));
  handle('encode:hex', (payload) => ops.encodeHex(payload.text));
  handle('decode:base64', (payload) => ops.decodeBase64(payload.text));
  handle('decode:hex', (payload) => ops.decodeHex(payload.text));
  handle('key:generate', (payload) => ops.generateKey(payload.keyType));
  handle('key:generateAs', (payload) => ops.generateKeyAs(payload.format));
  handle('random:string', (payload) => ops.randomString(payload.length));
  handle('uuid:create', () => ops.uuid());
  handle('clipboard:write', (payload) => {
    if (!payload || payload.text === '') throw new Error('Nothing to copy.');
    clipboard.writeText(payload.text);
    return true;
  });
}

module.exports = { registerHandlers };
