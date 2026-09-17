'use strict';

const $ = (id) => document.getElementById(id);

// --- toast ----------------------------------------------------------------

let toastTimer = null;

function toast(message, kind = 'info') {
  const el = $('toast');
  el.textContent = message;
  el.className = `toast show ${kind}`;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => {
    el.className = 'toast';
  }, kind === 'error' ? 4200 : 2200);
}

// Unwraps the {ok, value} envelope from the main process; a failure surfaces as
// a toast and a null result, so callers just check for null.
async function call(fn, payload) {
  const res = await fn(payload);
  if (!res.ok) {
    toast(res.error, 'error');
    return null;
  }
  return res.value;
}

// --- tabs -----------------------------------------------------------------

// Wires one tab strip. Each strip is scoped to its own <nav> so the Encode and
// Decode sub-strips do not deselect each other.
function wireTabStrip(nav, tabSelector, panelPrefix, attr) {
  const tabs = [...nav.querySelectorAll(tabSelector)];
  const select = (name) => {
    for (const tab of tabs) {
      const active = tab.dataset[attr] === name;
      tab.setAttribute('aria-selected', String(active));
      $(`${panelPrefix}${tab.dataset[attr]}`).classList.toggle('active', active);
    }
  };
  for (const tab of tabs) {
    tab.addEventListener('click', () => select(tab.dataset[attr]));
  }
  if (tabs.length) select(tabs[0].dataset[attr]);
}

function wireTabs(navSelector, tabSelector, panelPrefix, attr) {
  for (const nav of document.querySelectorAll(navSelector)) {
    wireTabStrip(nav, tabSelector, panelPrefix, attr);
  }
}

// --- copy buttons ---------------------------------------------------------

function wireCopyButtons() {
  for (const btn of document.querySelectorAll('[data-copy]')) {
    btn.addEventListener('click', async () => {
      const source = $(btn.dataset.copy);
      if (!source.value) {
        toast('Nothing to copy.');
        return;
      }
      const ok = await call(window.api.clipboardWrite, { text: source.value });
      if (ok) toast(btn.dataset.copyMessage ?? 'Copied.', 'success');
    });
  }
}

const selectedRadio = (name) =>
  document.querySelector(`input[name="${name}"]:checked`).value;

// --- generate-key modal (gui/build.gui.go: GenerateKeyDialog) -------------

function openGenerateKeyDialog(targetEntry) {
  const dlg = document.createElement('dialog');
  dlg.innerHTML = `
    <h2>Generate AES Key</h2>
    <label class="field-label">Select Key Format:</label>
    <div class="radio-group">
      <label><input type="radio" name="genkey-format" value="Hexadecimal" checked /> Hexadecimal</label>
      <label><input type="radio" name="genkey-format" value="Base64" /> Base64</label>
    </div>
    <div class="row">
      <button class="primary" value="generate">Generate</button>
      <button class="secondary" value="close">Close</button>
    </div>`;

  dlg.querySelector('button[value="close"]').addEventListener('click', () => dlg.close());
  dlg.querySelector('button[value="generate"]').addEventListener('click', async () => {
    const format = selectedRadio('genkey-format');
    const key = await call(window.api.keyGenerateAs, { format });
    if (key === null) return;
    targetEntry.value = key;
    dlg.close();
  });
  dlg.addEventListener('close', () => dlg.remove());

  document.body.appendChild(dlg);
  dlg.showModal();
}

// --- panels ---------------------------------------------------------------

function wireEncryption() {
  $('encrypt-genkey').addEventListener('click', () =>
    openGenerateKeyDialog($('encrypt-key')),
  );

  $('encrypt-run').addEventListener('click', async () => {
    const result = await call(window.api.cryptoEncrypt, {
      method: selectedRadio('encrypt-method'),
      text: $('encrypt-input').value,
      key: $('encrypt-key').value,
    });
    if (result !== null) $('encrypt-result').value = result;
  });
}

function wireDecryption() {
  $('decrypt-run').addEventListener('click', async () => {
    const result = await call(window.api.cryptoDecrypt, {
      method: selectedRadio('decrypt-method'),
      data: $('decrypt-input').value,
      key: $('decrypt-key').value,
    });
    if (result !== null) $('decrypt-result').value = result;
  });
}

function wireTransform(buttonId, inputId, resultId, apiFn) {
  $(buttonId).addEventListener('click', async () => {
    const result = await call(apiFn, { text: $(inputId).value });
    if (result !== null) $(resultId).value = result;
  });
}

function wireKeygen() {
  $('keygen-run').addEventListener('click', async () => {
    const key = await call(window.api.keyGenerate, {
      keyType: selectedRadio('keygen-type'),
    });
    if (key === null) return;
    $('keygen-hex').value = key.hex;
    $('keygen-base64').value = key.base64;
  });
}

function wireUuid() {
  $('uuid-run').addEventListener('click', async () => {
    const value = await call(window.api.uuidCreate);
    if (value !== null) $('uuid-result').value = value;
  });
}

function wireRandom() {
  const slider = $('random-slider');
  const label = $('random-length-label');
  const result = $('random-result');

  // Matches the Fyne slider: regenerate as the length changes.
  const regenerate = async () => {
    const length = Number(slider.value);
    label.textContent = `Length: ${length} bytes`;
    const value = await call(window.api.randomString, { length });
    if (value !== null) result.value = value;
  };

  slider.addEventListener('input', regenerate);
}

// --- boot -----------------------------------------------------------------

wireTabs('nav.tabs', '.tab', 'panel-', 'panel');
wireTabs('nav.subtabs', '.subtab', 'subpanel-', 'subpanel');
wireCopyButtons();
wireEncryption();
wireDecryption();
wireTransform('encode-base64-run', 'encode-base64-input', 'encode-base64-result', window.api.encodeBase64);
wireTransform('encode-hex-run', 'encode-hex-input', 'encode-hex-result', window.api.encodeHex);
wireTransform('decode-base64-run', 'decode-base64-input', 'decode-base64-result', window.api.decodeBase64);
wireTransform('decode-hex-run', 'decode-hex-input', 'decode-hex-result', window.api.decodeHex);
wireKeygen();
wireUuid();
wireRandom();
