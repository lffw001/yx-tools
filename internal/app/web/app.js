'use strict';

const $ = s => document.querySelector(s);
const state = {
  colos: [],
  picked: [],
  results: [],
  sortKey: '',
  sortDir: 'desc',
  running: false,
  hasToken: false,
  system: {},
  pool: 1000,
  httping: false,
  noDL: false,
};

function toast(msg, kind) {
  const el = document.createElement('div');
  el.className = 'toast' + (kind ? ' ' + kind : '');
  el.textContent = msg;
  $('#toasts').appendChild(el);
  setTimeout(() => {
    el.style.opacity = '0';
    el.style.transition = 'opacity .3s';
    setTimeout(() => el.remove(), 300);
  }, 3600);
}

async function api(path, opts) {
  const r = await fetch(path, opts);
  const text = await r.text();
  let data = null;
  try { data = text ? JSON.parse(text) : null; } catch (_) {}
  if (!r.ok) throw new Error((data && data.error) || text || ('HTTP ' + r.status));
  return data;
}

function renderChips() {
  const box = $('#coloChips');
  box.innerHTML = '';
  state.picked.forEach(code => {
    const c = state.colos.find(x => x.code === code);
    const el = document.createElement('div');
    el.className = 'chip';
    el.innerHTML = `<b>${code}</b>${c ? c.name : ''}<span>&times;</span>`;
    el.querySelector('span').onclick = () => {
      state.picked = state.picked.filter(x => x !== code);
      renderChips();
    };
    box.appendChild(el);
  });
  if (typeof setPing === 'function') setPing(state.picked.length > 0 || state.httping);
}

function renderColoList(q) {
  const list = $('#coloList');
  q = (q || '').trim().toLowerCase();
  if (!q) { list.classList.remove('show'); return; }
  const hit = state.colos.filter(c =>
    c.code.toLowerCase().includes(q) || c.name.includes(q) ||
    c.country.includes(q) || c.region.includes(q)
  ).slice(0, 40);
  list.innerHTML = '';
  if (!hit.length) { list.classList.remove('show'); return; }
  hit.forEach(c => {
    const el = document.createElement('div');
    el.className = 'colo-item';
    el.innerHTML = `<span>${c.name} <code>${c.code}</code></span><code>${c.country}</code>`;
    el.onclick = () => {
      if (!state.picked.includes(c.code)) state.picked.push(c.code);
      $('#coloSearch').value = '';
      list.classList.remove('show');
      renderChips();
    };
    list.appendChild(el);
  });
  list.classList.add('show');
}

function fmtSpeed(v) {
  const cls = v >= 5 ? 'g' : v >= 1 ? 'y' : 'r';
  return `<span class="${cls}">${v.toFixed(2)}</span>`;
}
function fmtDelay(v) {
  const cls = v <= 100 ? 'g' : v <= 250 ? 'y' : 'r';
  const txt = v > 0 && v < 10 ? v.toFixed(2) : v.toFixed(0);
  return `<span class="${cls}">${txt}</span>`;
}

function visibleRows() {
  const q = $('#filterText').value.trim().toLowerCase();
  let rows = state.results;
  if (q) {
    rows = rows.filter(r =>
      r.ip.toLowerCase().includes(q) ||
      (r.colo || '').toLowerCase().includes(q) ||
      (r.colo_name || '').includes(q)
    );
  }
  if (state.sortKey) {
    const k = state.sortKey, dir = state.sortDir === 'asc' ? 1 : -1;
    rows = rows.slice().sort((a, b) => {
      let x = k === 'loss' ? a.loss_rate : a[k];
      let y = k === 'loss' ? b.loss_rate : b[k];
      if (typeof x === 'string') return x.localeCompare(y) * dir;
      return (x - y) * dir;
    });
  }
  return rows;
}

function renderTable() {
  const rows = visibleRows();
  const tb = $('#tbody');
  tb.innerHTML = '';
  $('#emptyBox').classList.toggle('hidden', rows.length > 0);
  rows.forEach((r, i) => {
    const tr = document.createElement('tr');
    tr.innerHTML =
      `<td class="c-idx">${i + 1}</td>` +
      `<td class="mono">${r.ip}</td>` +
      `<td class="c-num mono">${r.port}</td>` +
      `<td class="c-num mono">${fmtDelay(r.delay)}</td>` +
      `<td class="c-num mono">${fmtSpeed(r.speed)}</td>` +
      `<td class="c-num mono">${(r.loss_rate * 100).toFixed(0)}%</td>` +
      `<td>${r.colo_name || '-'}${r.colo ? ' <code style="opacity:.6">' + r.colo + '</code>' : ''}</td>` +
      `<td class="c-act"><button class="copy" title="复制 IP:端口">⧉</button></td>`;
    tr.querySelector('.copy').onclick = () => {
      navigator.clipboard.writeText(`${r.ip}:${r.port}`).then(
        () => toast('已复制 ' + r.ip + ':' + r.port, 'ok'),
        () => toast('复制失败', 'err')
      );
    };
    tb.appendChild(tr);
  });
  $('#statResult').textContent = '结果 ' + state.results.length;
}

function setRunning(on, keepProgress) {
  state.running = on;
  $('#btnStart').classList.toggle('hidden', on);
  $('#btnStop').classList.toggle('hidden', !on);
  $('#statusDot').className = 'dot' + (on ? ' run' : '');
  if (keepProgress) return;
  const fill = $('#progressFill');
  fill.style.width = '';
  fill.className = on ? 'indet' : 'idle';
}

function setProgress(cur, total) {
  const fill = $('#progressFill');
  if (total > 0) {
    const pct = Math.min(100, Math.round(cur / total * 100));
    fill.className = '';
    fill.style.width = pct + '%';
    return pct;
  }
  fill.className = 'indet';
  fill.style.width = '';
  return null;
}

function connectEvents() {
  const es = new EventSource('/api/events');
  es.onmessage = ev => {
    let e;
    try { e = JSON.parse(ev.data); } catch (_) { return; }
    if (e.message) {
      $('#statusText').textContent = e.total > 0
        ? `${e.message}  ${e.current}/${e.total}`
        : e
