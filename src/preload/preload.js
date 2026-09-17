'use strict';

const { contextBridge, ipcRenderer } = require('electron');

// Fixed channel list: the renderer can only reach these, and only with a plain
// serialisable payload.
const CHANNELS = [
  'crypto:encrypt',
  'crypto:decrypt',
  'encode:base64',
  'encode:hex',
  'decode:base64',
  'decode:hex',
  'key:generate',
  'key:generateAs',
  'random:string',
  'uuid:create',
  'clipboard:write',
];

const api = {};
for (const channel of CHANNELS) {
  const name = channel
    .replace(/:(.)/, (_m, c) => c.toUpperCase())
    .replace(/[:-]/g, '');
  api[name] = (payload) => ipcRenderer.invoke(channel, payload ?? {});
}

contextBridge.exposeInMainWorld('api', Object.freeze(api));
