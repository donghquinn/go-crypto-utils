'use strict';

const path = require('node:path');
const { app, BrowserWindow, Menu, shell } = require('electron');
const { registerHandlers } = require('./ipc');

const isMac = process.platform === 'darwin';

function createWindow() {
  const win = new BrowserWindow({
    width: 900,
    height: 720,
    minWidth: 640,
    minHeight: 520,
    title: 'Secure Helper',
    show: false,
    backgroundColor: '#1b1b1f',
    webPreferences: {
      preload: path.join(__dirname, '..', 'preload', 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: true,
      spellcheck: false,
    },
  });

  win.once('ready-to-show', () => win.show());
  win.loadFile(path.join(__dirname, '..', 'renderer', 'index.html'));

  // Nothing in this app should navigate away or open a second window.
  win.webContents.setWindowOpenHandler(({ url }) => {
    if (url.startsWith('https://')) shell.openExternal(url);
    return { action: 'deny' };
  });
  win.webContents.on('will-navigate', (event) => event.preventDefault());

  return win;
}

function buildMenu() {
  Menu.setApplicationMenu(
    Menu.buildFromTemplate([
      ...(isMac ? [{ role: 'appMenu' }] : []),
      { role: 'fileMenu' },
      { role: 'editMenu' },
      { role: 'viewMenu' },
      { role: 'windowMenu' },
    ]),
  );
}

app.whenReady().then(() => {
  registerHandlers();
  buildMenu();
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', () => {
  if (!isMac) app.quit();
});
