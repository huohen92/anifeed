'use strict';

/* ================================================================
 * AniFeed Web UI
 * ================================================================ */

const ICONS = {
  dashboard: '<rect x="3" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="14" width="7" height="7" rx="1.5"/><rect x="3" y="14" width="7" height="7" rx="1.5"/>',
  rss: '<path d="M4 11a9 9 0 0 1 9 9"/><path d="M4 4a16 16 0 0 1 16 16"/><circle cx="5" cy="19" r="1"/>',
  tasks: '<path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>',
  history: '<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/>',
  logs: '<path d="M4 6h16"/><path d="M4 12h16"/><path d="M4 18h10"/>',
  settings: '<path d="M4 6h9"/><path d="M17 6h3"/><circledata/>',
  menu: '<path d="M4 6h16"/><path d="M4 12h16"/><path d="M4 18h16"/>',
  moon: '<path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8Z"/>',
  sun: '<circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4"/>',
  lock: '<rect x="4" y="10" width="16" height="11" rx="2"/><path d="M8 10V7a4 4 0 0 1 8 0v3"/>',
  plus: '<path d="M12 5v14"/><path d="M5 12h14"/>',
  play: '<path d="M7 4.5v15l13-7.5Z"/>',
  edit: '<path d="M11 4H5a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2h13a2 2 0 0 0 2-2v-6"/><path d="M18.5 2.5a2.1 2.1 0 0 1 3 3L12 15l-4 1 1-4Z"/>',
  trash: '<path d="M3 6h18"/><path d="M8 6V4a1 1 0 0 1 1-1h6a1 1 0 0 1 1 1v2"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6M14 11v6"/>',
  copy: '<rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>',
  refresh: '<path d="M21 12a9 9 0 1 1-3-6.7"/><path d="M21 3v6h-6"/>',
  login: '<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><path d="M10 17l5-5-5-5"/><path d="M15 12H3"/>',
  logout: '<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><path d="M16 17l5-5-5-5"/><path d="M21 12H9"/>',
  search: '<circle cx="11" cy="11" r="7"/><path d="M20 20l-3.5-3.5"/>',
  close: '<path d="M6 6l12 12M18 6L6 18"/>',
  qr: '<rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><path d="M14 14h3v3h-3z"/><path d="M20 14h1v3M14 20h3M20 20h1"/>',
  alert: '<path d="M12 3 2.5 20h19Z"/><path d="M12 9v5"/><circle cx="12" cy="17.2" r=".7" fill="currentColor"/>',
  check: '<circle cx="12" cy="12" r="9"/><path d="M8 12.5l2.6 2.6L16 9.7"/>',
  info: '<circle cx="12" cy="12" r="9"/><path d="M12 11v5"/><circle cx="12" cy="8" r=".7" fill="currentColor"/>',
  image: '<rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="M21 15l-5-5L5 21"/>',
  beaker: '<path d="M9 3h6"/><path d="M10 3v6l-5 9a2 2 0 0 0 1.7 3h10.6a2 2 0 0 0 1.7-3L14 9V3"/><path d="M7.5 15h9"/>',
  heart: '<path d="M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.6l-1-1a5.5 5.5 0 0 0-7.8 7.8l1 1L12 21l7.8-7.6 1-1a5.5 5.5 0 0 0 0-7.8Z"/>',
  external: '<path d="M14 3h7v7"/><path d="M21 3l-9 9"/><path d="M19 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2h5"/>',
  folder: '<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2Z"/>',
  clock: '<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/>',
  cloud: '<path d="M17.5 19a4.5 4.5 0 0 0 .5-8.97A6 6 0 0 0 6.1 9.5 4.5 4.5 0 0 0 6.5 19Z"/>',
  hash: '<path d="M5 9h14M5 15h14M9 4 7 20M17 4l-2 16"/>',
  database: '<ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.7 4 3 9 3s9-1.3 9-3V5"/><path d="M3 12c0 1.7 4 3 9 3s9-1.3 9-3"/>',
  server: '<rect x="2" y="3" width="20" height="7" rx="2"/><rect x="2" y="14" width="20" height="7" rx="2"/><path d="M6 6.5h.01M6 17.5h.01"/>',
  shield: '<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z"/>',
  arrow: '<path d="M5 12h14"/><path d="M13 6l6 6-6 6"/>',
  download: '<path d="M12 3v12"/><path d="M7 10l5 5 5-5"/><path d="M4 19h16"/>',
  flow: '<circle cx="6" cy="6" r="2"/><circle cx="18" cy="18" r="2"/><path d="M8 6h6a4 4 0 0 1 4 4v6"/>',
  filter: '<path d="M3 5h18l-7 8v6l-4 2v-8Z"/>',
  link: '<path d="M10 13a5 5 0 0 0 7 0l2-2a5 5 0 0 0-7-7l-1 1"/><path d="M14 11a5 5 0 0 0-7 0l-2 2a5 5 0 0 0 7 7l1-1"/>',
  key: '<circle cx="8" cy="15" r="4"/><path d="M11 12 20 3"/><path d="M17 3h4v4"/>',
  bell: '<path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9"/><path d="M10 21h4"/>',
  eye: '<path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z"/><circle cx="12" cy="12" r="3"/>',
  upload: '<path d="M12 16V4"/><path d="m7 9 5-5 5 5"/><path d="M5 20h14"/>',
  eyeoff: '<path d="M3 3l18 18"/><path d="M10.6 10.6a3 3 0 0 0 4.2 4.2"/><path d="M9.9 5.2A10.9 10.9 0 0 1 12 5c6.5 0 10 7 10 7a17.5 17.5 0 0 1-3.1 4.1M6.3 6.3A17.6 17.6 0 0 0 2 12s3.5 7 10 7c1.2 0 2.3-.2 3.3-.6"/>',
  terminal: '<path d="M4 5h16v14H4z"/><path d="M7 9l3 3-3 3M13 15h4"/>',
  save: '<path d="M5 3h11l5 5v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2Z"/><path d="M8 3v6h8V3"/><path d="M8 21v-7h8v7"/>',
  file: '<path d="M6 2h8l5 5v15H6z"/><path d="M14 2v5h5"/>',
};
ICONS.settings = '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.9l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.7 1.7 0 0 0-1.9-.3 1.7 1.7 0 0 0-1 1.5V21a2 2 0 1 1-4 0v-.1a1.7 1.7 0 0 0-1.1-1.6 1.7 1.7 0 0 0-1.9.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.7 1.7 0 0 0 .3-1.9 1.7 1.7 0 0 0-1.5-1H3a2 2 0 1 1 0-4h.1a1.7 1.7 0 0 0 1.6-1.1 1.7 1.7 0 0 0-.3-1.9l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.7 1.7 0 0 0 1.9.3H9a1.7 1.7 0 0 0 1-1.5V3a2 2 0 1 1 4 0v.1a1.7 1.7 0 0 0 1 1.5 1.7 1.7 0 0 0 1.9-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.7 1.7 0 0 0-.3 1.9V9a1.7 1.7 0 0 0 1.5 1H21a2 2 0 1 1 0 4h-.1a1.7 1.7 0 0 0-1.5 1Z"/>';



const SOURCES = {
  mikan: { label: 'Mikan', short: 'M', color: 'c0' },
  anibt: { label: 'AniBT', short: 'A', color: 'c1' },
  animegarden: { label: 'AnimeGarden', short: 'G', color: 'c2' },
  other: { label: '其他', short: 'R', color: 'c5' },
};

function detectSource(rawUrl) {
  try {
    const host = new URL(rawUrl).hostname.toLowerCase();
    if (host === 'mikanani.me' || host === 'mikanime.tv') return 'mikan';
    if (host === 'anibt.net') return 'anibt';
    if (host === 'animes.garden' || host === 'api.animes.garden') return 'animegarden';
  } catch { /* ignore */ }
  return 'other';
}

function sourceInfo(source, rawUrl) {
  const key = SOURCES[source] ? source : detectSource(rawUrl);
  return { key, ...SOURCES[key] };
}

function subgroupLabel(value, feedName = '') {
  let text = String(value || '').trim();
  if (feedName && text.startsWith(feedName + ' · ')) text = text.slice(feedName.length + 3);
  const index = text.lastIndexOf(' · ');
  if (index >= 0) text = text.slice(index + 3);
  return text.trim();
}

function bgmURLFromFeedURL(rawUrl) {
  try {
    const u = new URL(rawUrl);
    const id = u.searchParams.get('bgmId') || u.searchParams.get('bangumiId') || '';
    return id ? `https://bgm.tv/subject/${id}` : '';
  } catch {
    return '';
  }
}

function urlSubgroupName(rawUrl) {
  try {
    const slug = new URL(rawUrl).searchParams.get('groupSlug') || '';
    if (!slug) return '';
    const compound = slug.match(/^(.+?)(house|hub|sub|team|raws|anime)$/i);
    const word = value => value ? value.charAt(0).toUpperCase() + value.slice(1).toLowerCase() : '';
    return compound ? word(compound[1]) + word(compound[2]) : word(slug);
  } catch {
    return '';
  }
}

const DEFAULT_NOTIFY_EVENTS = ['run', 'added', 'completed', 'failed', 'idle', 'hourly', 'wash'];

const DEFAULT_NOTIFY_TEMPLATE = `通知
事件类型: \${action}
标题: \${title}
评分: \${score}
TMDB标题: \${themoviedbName}
BGM: \${bgmUrl}
季: \${season}
集: \${episode}
字幕组: \${subgroup}
进度: \${currentEpisodeNumber}/\${totalEpisodeNumber}
首播:  \${year}年\${month}月\${date}日
事件: \${text}
BGM集标题: \${bgmEpisodeTitle}`;

const PAGE_META = {
  dashboard: { title: '概览', subtitle: '', icon: 'dashboard' },
  feeds: { title: '订阅管理', subtitle: '', icon: 'rss' },
  tasks: { title: '离线任务', subtitle: '', icon: 'tasks' },
  history: { title: '历史记录', subtitle: '', icon: 'history' },
  logs: { title: '运行日志', subtitle: '', icon: 'logs' },
  settings: { title: '设置', subtitle: '', icon: 'settings' },
};

const state = {
  token: localStorage.getItem('anifeed_token') || '',
  status: null,
  config: null,
  meta: null,
  tasks: [],
  history: [],
  logs: [],
  route: '',
  feedQuery: '',
  feedSource: '',
  feedWeekdays: (() => {
    try { return JSON.parse(localStorage.getItem('anifeed_weekdays') || '{}'); } catch { return {}; }
  })(),
  covers: (() => {
    try {
      const raw = localStorage.getItem('anifeed_covers');
      return raw ? JSON.parse(raw) : {};
    } catch { return {}; }
  })(),
  rules: null,
  taskFilter: 'all',
  historyName: '',
  historyLimit: 50,
  settingsTab: 'form',
  settingsSites: null,
  logAuto: true,
  logMode: localStorage.getItem('anifeed_log_mode') || 'list',
  logLevel: 'all',
  logQuery: '',
  timers: [],
  jobTimer: null,
  logTimer: null,
  dirty: false,
  loginPollToken: 0,
};

/* ================================================================
 * Utilities
 * ================================================================ */
const $ = (sel, root = document) => root.querySelector(sel);
const $$ = (sel, root = document) => Array.from(root.querySelectorAll(sel));

function esc(value) {
  return String(value == null ? '' : value)
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
}

function icon(name) {
  return `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">${ICONS[name] || ''}</svg>`;
}

function fmtSize(bytes) {
  const n = Number(bytes) || 0;
  if (n < 1024) return `${n} B`;
  const units = ['KB', 'MB', 'GB', 'TB', 'PB'];
  let v = n / 1024, i = 0;
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++; }
  return `${v >= 100 ? v.toFixed(0) : v.toFixed(1)} ${units[i]}`;
}

function fmtPercent(value) {
  const n = Number(value) || 0;
  return `${n >= 100 ? 100 : n.toFixed(n < 10 ? 1 : 0)}%`;
}

function fmtDate(value) {
  if (!value) return '—';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);
  const pad = n => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

function timeAgo(value) {
  if (!value) return '';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return '';
  const diff = Math.max(0, Date.now() - d.getTime()) / 1000;
  if (diff < 60) return '刚刚';
  if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`;
  if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`;
  return `${Math.floor(diff / 86400)} 天前`;
}

function debounce(fn, ms = 200) {
  let timer = null;
  return function (...args) {
    clearTimeout(timer);
    timer = setTimeout(() => fn.apply(this, args), ms);
  };
}

function siteClass(site) {
  let hash = 0;
  for (const ch of String(site || '')) hash = (hash * 31 + ch.charCodeAt(0)) % 997;
  return `c${hash % 6}`;
}

function siteInitial(site) {
  const s = String(site || '?').replace(/^www\./, '');
  return s.slice(0, 1).toUpperCase() || '?';
}

function setTimer(fn, ms) {
  const id = setTimeout(fn, ms);
  state.timers.push(id);
  return id;
}

function clearTimers() {
  state.timers.forEach(clearTimeout);
  state.timers = [];
  if (state.jobTimer) { clearTimeout(state.jobTimer); state.jobTimer = null; }
  if (state.logTimer) { clearInterval(state.logTimer); state.logTimer = null; }
}

function setLoading(button, loading, label) {
  if (!button) return;
  if (loading) {
    button.dataset.label = button.innerHTML;
    button.disabled = true;
    button.innerHTML = `<span class="spinner"></span>${esc(label || '处理中')}`;
  } else {
    button.disabled = false;
    if (button.dataset.label) button.innerHTML = button.dataset.label;
  }
}

async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text);
    toast('已复制到剪贴板', 'success');
  } catch {
    const area = document.createElement('textarea');
    area.value = text;
    document.body.appendChild(area);
    area.select();
    try { document.execCommand('copy'); toast('已复制到剪贴板', 'success'); }
    catch { toast('复制失败，请手动选择文本', 'error'); }
    area.remove();
  }
}

/* ================================================================
 * API
 * ================================================================ */
async function api(path, options = {}) {
  const opts = { ...options };
  const headers = { Accept: 'application/json', ...(opts.headers || {}) };
  if (opts.body && typeof opts.body !== 'string') {
    headers['Content-Type'] = 'application/json';
    opts.body = JSON.stringify(opts.body);
  }
  if (state.token) headers['X-Access-Token'] = state.token;
  opts.headers = headers;

  let res;
  try {
    res = await fetch(path, opts);
  } catch (err) {
    throw new Error('无法连接 AniFeed 服务，请确认服务仍在运行');
  }
  const type = res.headers.get('content-type') || '';
  let data;
  if (type.includes('application/json')) {
    data = await res.json().catch(() => ({}));
  } else {
    data = await res.text().catch(() => '');
  }

  if (!res.ok) {
    if (res.status === 401 && data && data.code === 'token_required') {
      showTokenGate();
      const err = new Error(data.error || '需要访问令牌');
      err.code = 'token_required';
      throw err;
    }
    if (res.status === 401 && data && data.code === 'auth_required') {
      showAuthGate({ configured: true });
      const err = new Error(data.error || '需要登录');
      err.code = 'auth_required';
      throw err;
    }
    const message = (data && data.error) ? data.error : (typeof data === 'string' && data.trim()) || `请求失败 (${res.status})`;
    throw new Error(message);
  }
  return data;
}

/* ================================================================
 * Toast / Modal
 * ================================================================ */
function toast(message, type = 'success', timeout = 3600) {
  const root = $('#toast-root');
  const el = document.createElement('div');
  el.className = `toast ${type}`;
  const iconName = type === 'error' ? 'alert' : type === 'info' ? 'info' : 'check';
  el.innerHTML = `${icon(iconName)}<span>${esc(message)}</span>`;
  root.appendChild(el);
  setTimeout(() => {
    el.classList.add('out');
    setTimeout(() => el.remove(), 220);
  }, timeout);
}

function openModal({ title, subtitle = '', body = '', footer = '', size = '', onMount }) {
  const root = $('#modal-root');
  const backdrop = document.createElement('div');
  backdrop.className = 'modal-backdrop';
  backdrop.innerHTML = `
    <div class="modal ${size}" role="dialog" aria-modal="true">
      <div class="modal-head">
        <div><h3>${esc(title)}</h3>${subtitle ? `<p>${esc(subtitle)}</p>` : ''}</div>
        <button class="icon-btn close" data-action="modal-close" aria-label="关闭">${icon('close')}</button>
      </div>
      <div class="modal-body">${body}</div>
      ${footer ? `<div class="modal-foot">${footer}</div>` : ''}
    </div>`;
  root.appendChild(backdrop);
  requestAnimationFrame(() => backdrop.classList.add('show'));

  const handle = { el: backdrop, close: null, closed: false };
  const close = () => {
    handle.closed = true;
    backdrop.classList.remove('show');
    document.removeEventListener('keydown', onKey);
    setTimeout(() => backdrop.remove(), 180);
  };
  handle.close = close;
  const onKey = event => { if (event.key === 'Escape') close(); };
  document.addEventListener('keydown', onKey);
  backdrop.addEventListener('click', event => {
    if (event.target === backdrop) close();
    if (event.target.closest('[data-action="modal-close"]')) close();
  });

  if (onMount) onMount(backdrop, close, handle);
  return handle;
}

function confirmDialog({ title, message, confirmText = '确认', cancelText = '取消', danger = true }) {
  return new Promise(resolve => {
    let settled = false;
    const done = value => {
      if (settled) return;
      settled = true;
      resolve(value);
      modal.close();
    };
    const modal = openModal({
      title,
      size: 'narrow',
      body: `<p style="color:var(--muted);font-size:13px;line-height:1.7">${message}</p>`,
      footer: `<button class="btn" data-action="cancel">${esc(cancelText)}</button>
               <button class="btn ${danger ? 'btn-danger' : 'btn-primary'}" data-action="ok">${esc(confirmText)}</button>`,
      onMount(el) {
        el.querySelector('[data-action="cancel"]').addEventListener('click', () => done(false));
        el.querySelector('[data-action="ok"]').addEventListener('click', () => done(true));
        el.addEventListener('click', event => {
          if (event.target === el) done(false);
        });
        el.addEventListener('click', event => {
          if (event.target.closest('[data-action="modal-close"]')) done(false);
        });
      },
    });
  });
}

/* ================================================================
 * Chrome
 * ================================================================ */
function applyTheme(theme) {
  document.documentElement.dataset.theme = theme;
  localStorage.setItem('anifeed_theme', theme);
  const btn = $('#theme-btn');
  if (btn) btn.innerHTML = icon(theme === 'dark' ? 'sun' : 'moon');
}

function initTheme() {
  const saved = localStorage.getItem('anifeed_theme');
  const preferred = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  applyTheme(saved || preferred);
  $('#theme-btn').addEventListener('click', () => {
    applyTheme(document.documentElement.dataset.theme === 'dark' ? 'light' : 'dark');
  });
}

function initChrome() {
  $$('[data-icon]').forEach(el => { el.innerHTML = icon(el.dataset.icon); });
  initTheme();

  window.addEventListener('hashchange', () => { renderRoute(); });

}

function updateChrome() {
  const status = state.status;
  const pill = $('#login-pill');
  const pillText = $('#login-pill-text');
  const sideDot = $('#side-dot');
  const sideStatus = $('#side-status');

  pill.classList.remove('ok', 'err', 'busy');
  sideDot.classList.remove('ok', 'err', 'busy');

  if (!status) {
    pillText.textContent = '未连接';
    sideStatus.textContent = '未连接';
    sideDot.classList.add('err');
  } else if (status.login_state === 'ok') {
    const name = status.user && status.user.name ? status.user.name : '115 已登录';
    pillText.textContent = status.user && status.user.is_vip ? `${name} · VIP` : name;
    pill.classList.add('ok');
    sideDot.classList.add('ok');
    sideStatus.textContent = '服务正常';
  } else if (status.login_state === 'checking' || status.login_state === 'unknown') {
    pillText.textContent = '检测登录状态…';
    pill.classList.add('busy');
    sideDot.classList.add('busy');
    sideStatus.textContent = '检测中';
  } else if (status.login_state === 'nocookies') {
    pillText.textContent = '未登录 115';
    pill.classList.add('err');
    sideDot.classList.add('err');
    sideStatus.textContent = '未配置 Cookies';
  } else {
    pillText.textContent = '115 登录失效';
    pill.classList.add('err');
    sideDot.classList.add('err');
    sideStatus.textContent = '登录失效';
  }
  if (status) $('#version-text').textContent = `v${status.version}`;
}

async function loadStatus(force = false) {
  if (!force && state.status && Date.now() - (state.statusAt || 0) < 3000) {
    return state.status;
  }
  const status = await api('/api/status');
  state.status = status;
  state.statusAt = Date.now();
  updateChrome();
  return status;
}

/* ================================================================
 * Routing
 * ================================================================ */
function currentRoute() {
  const raw = (location.hash || '').replace(/^#\/?/, '').split('?')[0];
  return PAGE_META[raw] ? raw : 'dashboard';
}

function renderRoute() {
  clearTimers();
  const route = currentRoute();
  state.route = route;
  const meta = PAGE_META[route];
  document.title = 'AniFeed';
  $('#page-title').textContent = meta.title;
  $('#page-subtitle').textContent = meta.subtitle;
  $$('.nav-item').forEach(el => el.classList.toggle('active', el.dataset.route === route));
  $$('[data-mobile-route]').forEach(el => el.classList.toggle('active', el.dataset.mobileRoute === route));
  document.body.classList.toggle('hide-topbar-actions', route === 'logs');

  switch (route) {
    case 'dashboard': loadDashboard(); break;
    case 'feeds': loadFeedsPage(); break;
    case 'tasks': loadTasksPage(); break;
    case 'history': loadHistoryPage(); break;
    case 'logs': loadLogsPage(); break;
    case 'settings': loadSettingsPage(); break;
    default: loadDashboard();
  }
}

function renderFatal(message) {
  $('#view').innerHTML = `
    <div class="card">
      <div class="empty">
        <div class="empty-icon">${icon('alert')}</div>
        <strong>加载失败</strong>
        <p>${esc(message)}</p>
        <button class="btn btn-primary" id="retry-btn">${icon('refresh')}重试</button>
      </div>
    </div>`;
  $('#retry-btn').addEventListener('click', () => renderRoute());
}

/* ================================================================
 * Dashboard
 * ================================================================ */
async function loadDashboard() {
  const view = $('#view');
  const hasCache = !!state.status;
  if (hasCache) {
    renderDashboard();
  } else {
    view.innerHTML = `<div class="page-loading"><span class="spinner"></span><span>正在读取服务状态…</span></div>`;
  }

  let statusError = null;
  const jobs = [
    loadStatus().catch(err => { statusError = err; }),
    ensureConfig().catch(() => {}),
    api('/api/history?limit=200').then(data => { state.history = data.items || []; }).catch(() => {}),
    api('/api/logs?limit=12').then(data => { state.logs = data.lines || []; }).catch(() => {}),
  ];
  jobs.push(api('/api/tasks?limit=200').then(data => { state.tasks = data.tasks || []; }).catch(() => { state.tasks = []; }));
  await Promise.allSettled(jobs);
  if (!state.status) {
    renderFatal(statusError ? statusError.message : '无法读取服务状态');
    return;
  }
  renderDashboard();
  scheduleLoginPoll();
}

function scheduleLoginPoll() {
  if (state.route !== 'dashboard' || !state.status) return;
  if (!['checking', 'unknown'].includes(state.status.login_state)) return;
  const token = ++state.loginPollToken;
  setTimer(async () => {
    if (state.route !== 'dashboard' || state.loginPollToken !== token) return;
    try { await loadStatus(); } catch { /* keep polling */ }
    if (state.route !== 'dashboard' || state.loginPollToken !== token) return;
    renderDashboard();
    scheduleLoginPoll();
  }, 2200);
}

function dashboardStats() {
  const tasks = state.tasks || [];
  const running = tasks.filter(t => t.state === 'running').length;
  const done = tasks.filter(t => t.state === 'done').length;
  const failed = tasks.filter(t => t.state === 'failed').length;
  return { running, done, failed };
}

function renderDashboard() {
  const status = state.status;
  if (!status) return;
  const stats = dashboardStats();
  const loginOk = status.login_state === 'ok';
  const checking = ['checking', 'unknown'].includes(status.login_state);
  const user = status.user;
  const avatar = user && user.name ? esc(user.name.slice(0, 1).toUpperCase()) : icon('cloud');

  let heroTitle = '尚未登录 115 网盘';
  let heroDesc = '登录后即可把 RSS 条目提交到 115 离线下载。支持 Cookies 或扫码登录。';
  if (loginOk) {
    heroTitle = user && user.name ? `${user.name}${user.is_vip ? ' · VIP' : ''}` : '115 已登录';
    heroDesc = `UID ${user && user.id ? user.id : '—'} · 凭据保存在本地，可随时在设置中更换`;
  } else if (checking) {
    heroTitle = '正在检测 115 登录状态…';
    heroDesc = '首次加载需要向 115 验证 Cookies，请稍候。';
  } else if (status.login_state === 'error') {
    heroTitle = '115 登录已失效';
    heroDesc = status.login_error || 'Cookies 可能已过期，请重新登录。';
  }

  const banners = [];
  if (!status.config_exists) {
    banners.push(`<div class="banner">${icon('info')}<span>尚未创建 <code>config.toml</code>。在「设置」中保存后会自动生成，旧版 rss.json / .cookies 会自动迁移。</span><span class="banner-actions"><a class="btn btn-sm" href="#/settings">前往设置</a></span></div>`);
  }
  if (state.meta && state.meta.legacy_imported) {
    banners.push(`<div class="banner warn">${icon('alert')}<span>检测到旧版配置文件，当前显示的是合并后的配置。保存一次即可迁移到 <code>config.toml</code>。</span></div>`);
  }
  if (status.database_error) {
    banners.push(`<div class="banner warn">${icon('alert')}<span>数据库读取异常：${esc(status.database_error)}</span></div>`);
  }

  const recent = (state.tasks || []).slice(0, 6);
  let recentTitle = '最近任务';
  let recentSub = '来自 115 的最新离线下载';
  let recentHtml;
  if (recent.length) {
    recentHtml = `<div class="table-wrap"><table class="table recent-table">
      <thead><tr><th>任务</th><th>大小</th><th>状态</th><th>进度</th></tr></thead>
      <tbody>${recent.map(taskRowCompact).join('')}</tbody>
    </table></div>`;
  } else if ((state.history || []).length) {
    recentTitle = '最近提交';
    recentSub = '本地历史记录（登录 115 后可查看云端任务）';
    recentHtml = `<div class="table-wrap"><table class="table recent-table">
      <thead><tr><th>标题</th><th class="hide-sm">订阅</th><th class="hide-sm">时间</th></tr></thead>
      <tbody>${state.history.slice(0, 6).map(item => `<tr>
        <td class="name-cell" data-label="标题"><div class="title" title="${esc(item.title)}">${esc(item.title || '—')}</div><div class="meta">${esc(item.task_hash || '')}</div></td>
        <td class="hide-sm" data-label="订阅"><span class="badge badge-info">${esc(item.config_name || '未命名')}</span></td>
        <td class="numeric nowrap hide-sm" data-label="时间">${fmtDate(item.created_at)}</td>
      </tr>`).join('')}</tbody>
    </table></div>`;
  } else {
    recentHtml = emptyBlock('cloud', '还没有离线任务', '运行订阅或手动添加磁力链接后，任务会显示在这里。');
  }

  const allFeeds = state.config && state.config.rss ? state.config.rss : [];
  const today = new Date();
  const todayKey = `${today.getFullYear()}-${today.getMonth()}-${today.getDate()}`;
  const isToday = value => {
    const date = new Date(value);
    return !Number.isNaN(date.getTime()) && `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}` === todayKey;
  };
  const todayHistory = (state.history || []).filter(item => isToday(item.created_at));
  const animeFeeds = allFeeds.filter(feed => todayHistory.some(item =>
    (item.config_url && item.config_url === feed.url) || (item.config_name && item.config_name === feed.name)
  ));
  const animeGrid = animeFeeds.length
    ? `<div class="anime-grid">${animeFeeds.map(feed => dashboardAnimeCard(feed, allFeeds.indexOf(feed))).join('')}</div>`
    : `<div class="card">${emptyBlock('clock', '今天还没有订阅更新', '今天没有检测到新提交的番剧。', `<a class="btn btn-primary" href="#/feeds">查看全部订阅</a>`)}</div>`;

  const logs = (state.logs || []).slice(-9);
  const logHtml = logs.length
    ? `<pre class="log-pre" style="max-height:260px">${renderLogLines(logs)}</pre>`
    : emptyBlock('terminal', '暂无日志', '运行一次订阅后即可看到服务端输出。');

  $('#view').innerHTML = `
    ${banners.join('')}
    <section class="stats">
      <div class="card stat accent">
        <div class="stat-top"><span class="stat-icon">${icon('rss')}</span>订阅总数</div>
        <div class="stat-value">${status.feed_count}</div>
        <div class="stat-foot">${status.auto_enabled ? `自动运行每 ${status.auto_interval_minutes} 分钟` : '自动运行未开启'}</div>
      </div>
      <div class="card stat green">
        <div class="stat-top"><span class="stat-icon">${icon('download')}</span>今日提交</div>
        <div class="stat-value">${status.today_count}</div>
        <div class="stat-foot">历史累计 ${status.history_count} 条</div>
      </div>
      <div class="card stat amber">
        <div class="stat-top"><span class="stat-icon">${icon('cloud')}</span>云端任务</div>
        <div class="stat-value">${state.tasks.length}</div>
        <div class="stat-foot">${stats.running} 个运行中 · ${stats.done} 个已完成</div>
      </div>
      <div class="card stat cyan">
        <div class="stat-top"><span class="stat-icon">${icon('clock')}</span>提交成功率</div>
        <div class="stat-value">${state.tasks.length ? Math.round((stats.done / Math.max(1, state.tasks.length)) * 100) + '%' : '—'}</div>
        <div class="stat-foot">${stats.failed} 个失败任务</div>
      </div>
    </section>

    <section class="dashboard-anime">
      <div class="card-head">
        <div><h3>今天的订阅</h3><div class="sub">今天有更新的番剧 · 点击卡片可编辑</div></div>

      </div>
      ${animeGrid}
    </section>

    <section class="card">
      <div class="card-head">
        <div><h3>${recentTitle}</h3><div class="sub">${recentSub}</div></div>
        <div class="right"><a class="btn btn-sm btn-ghost" href="${recent.length ? '#/tasks' : '#/history'}">查看全部 ${icon('arrow')}</a></div>
      </div>
      ${recentHtml}
    </section>`;


  const runAll = $('#run-all-btn');
  if (runAll) runAll.addEventListener('click', runAllFeeds);
  hydrateFeedCovers($('#view'));
  refreshFeedWeekdays();
  $$('[data-dashboard-feed]').forEach(card => card.addEventListener('click', event => {
    const index = Number(card.dataset.dashboardFeed);
    if (event.target.closest('[data-dash-run]')) { runSingleFeed(state.config.rss[index]); return; }
    openFeedModal(index);
  }));
}

function dashboardAnimeCard(feed, index) {
  const source = sourceInfo(feed.source, feed.url);
  const cached = state.covers[feed.url];
  const enabled = feed.enabled !== false;
  return `<article class="anime-card ${enabled ? '' : 'disabled'}" data-dashboard-feed="${index}">
    <div class="anime-poster" data-feed-cover data-index="${index}">
      <span class="poster-fallback"><img src="/poster-placeholder.jpg?v=2" alt="AniFeed"></span>
      <img alt="${esc(cached ? cached.title || '' : '')}" ${cached && cached.cover ? `src="${esc(cached.cover)}"` : ''} loading="lazy" referrerpolicy="no-referrer">
      <span class="anime-badge">${completed ? '已完成' : disabled ? '未启用' : '追番中'}</span>
    </div>
    <div class="anime-body">
      <h4 title="${esc(feed.name)}">${esc(feed.name || '未命名订阅')}</h4>
      <div class="anime-meta">${esc(source.label)} · ${esc(feed.site || '自动识别')}</div>
      <div class="anime-foot">
        <span class="badge">${feed.rename_season ? `第 ${feed.rename_season} 季` : '自动季数'}</span>
        ${feed.tmdb_id ? `<span class="badge badge-success">TMDB</span>` : ''}
        <button class="btn btn-sm btn-ghost" type="button" data-dash-run title="运行订阅">${icon('play')}</button>
      </div>
    </div>
  </article>`;
}

function taskRowCompact(task) {
  const stateBadge = taskStateBadge(task.state);
  return `<tr>
    <td class="name-cell" data-label="任务"><div class="title" title="${esc(task.name || '未命名任务')}">${esc(task.name || '未命名任务')}</div><div class="meta">${esc((task.info_hash || '').slice(0, 18))}</div></td>
    <td class="numeric" data-label="大小">${fmtSize(task.size)}</td>
    <td data-label="状态">${stateBadge}</td>
    <td data-label="进度">${progressBar(task)}</td>
  </tr>`;
}

function taskStateBadge(stateName) {
  switch (stateName) {
    case 'running': return `<span class="badge badge-primary">运行中</span>`;
    case 'done': return `<span class="badge badge-success">已完成</span>`;
    case 'failed': return `<span class="badge badge-danger">失败</span>`;
    default: return `<span class="badge">等待中</span>`;
  }
}

function progressBar(task) {
  const pct = Math.max(0, Math.min(100, Number(task.percent) || 0));
  const cls = task.state === 'done' ? 'done' : task.state === 'failed' ? 'failed' : '';
  return `<div class="progress"><div class="progress-track"><div class="progress-fill ${cls}" style="width:${pct}%"></div></div><span class="progress-text">${fmtPercent(pct)}</span></div>`;
}

function emptyBlock(iconName, title, desc, action = '') {
  return `<div class="empty">
    <div class="empty-icon">${icon(iconName)}</div>
    <strong>${esc(title)}</strong>
    <p>${esc(desc)}</p>
    ${action}
  </div>`;
}

/* ================================================================
 * Feeds
 * ================================================================ */
async function ensureConfig() {
  if (state.config) return;
  const resp = await api('/api/config');
  applyConfigState(resp.config, resp.meta);
}

function applyConfigState(config, meta) {
  state.config = config;
  if (meta) state.meta = meta;
  document.body.classList.toggle('hide-hints', !!(config && config.lab && config.lab.hide_hints));
  state.settingsSites = Object.entries((config && config.sites) || {}).map(([host, site]) => ({
    host,
    https_agent: !!site.https_agent,
    headers: Object.entries(site.headers || {}).map(([k, v]) => ({ key: k, value: v })),
  }));
}

async function loadFeedsPage() {
  if (!state.config) {
    $('#view').innerHTML = `<div class="page-loading"><span class="spinner"></span><span>正在读取订阅配置…</span></div>`;
  }
  try {
    await ensureConfig();
    await loadStatus().catch(() => {});
    await api('/api/history?limit=500').then(data => { state.history = data.items || []; }).catch(() => {});
  } catch (err) {
    renderFatal(err.message);
    return;
  }
  renderFeedsPage();
}

function renderFeedsPage() {
  const feeds = (state.config && state.config.rss) || [];
  const legacyBanner = state.meta && state.meta.legacy_imported
    ? `<div class="banner warn">${icon('alert')}<span>当前订阅来自旧版配置，保存任意修改后会写入 <code>config.toml</code>。</span></div>` : '';

  $('#view').innerHTML = `
    ${legacyBanner}
    <div class="toolbar feed-toolbar">
      <div class="search-box">${icon('search')}<input class="input" id="feed-search" placeholder="搜索名称、站点或 URL" value="${esc(state.feedQuery)}"></div>
      <select class="select" id="feed-source-filter" style="width:auto;min-width:140px">
        <option value="">全部来源</option>
        ${Object.entries(SOURCES).map(([key, info]) => `<option value="${key}" ${state.feedSource === key ? 'selected' : ''}>${info.label}</option>`).join('')}
      </select>
      <select class="select" id="feed-status-filter" style="width:auto;min-width:120px">
        <option value="all" ${state.feedStatus === 'all' ? 'selected' : ''}>全部状态</option>
        <option value="enabled" ${state.feedStatus === 'enabled' ? 'selected' : ''}>启用</option>
        <option value="disabled" ${state.feedStatus === 'disabled' ? 'selected' : ''}>未启用</option>
      </select>
      <span class="badge badge-info">${feeds.length} 个订阅</span>
      <span class="spacer"></span>
      <button class="btn" id="manage-feeds">${icon('settings')}管理</button>
      <button class="btn" id="run-all-feeds">${icon('play')}运行全部</button>
      <button class="btn btn-icon-only" id="refresh-feeds" title="刷新">${icon('refresh')}</button>
      <button class="btn btn-primary" id="browse-feeds">${icon('plus')}添加</button>
    </div>
    <div id="feed-list"></div>`;

  const search = $('#feed-search');
  search.addEventListener('input', debounce(() => {
    state.feedQuery = search.value.trim().toLowerCase();
    renderFeedList();
  }, 160));
  $('#feed-source-filter').addEventListener('change', event => {
    state.feedSource = event.target.value;
    renderFeedList();
  });
  $('#feed-status-filter').addEventListener('change', event => {
    state.feedStatus = event.target.value;
    renderFeedList();
  });
  $('#browse-feeds').addEventListener('click', () => openBrowseModal('anibt', pick => {
    openFeedModal(-1, {
      source: pick.source,
      name: pick.title,
      url: pick.url,
    });
  }));
  $('#manage-feeds').addEventListener('click', openManageFeeds);
  $('#refresh-feeds').addEventListener('click', loadFeedsPage);
  $('#run-all-feeds').addEventListener('click', runAllFeeds);
  renderFeedList();
  refreshFeedWeekdays();
}

const WEEKDAY_NAMES = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];

function parseEpisodeNumber(title) {
  const value = String(title || '');
  const match = value.match(/S\d{1,2}\s*E(\d{1,3}(?:\.5)?)/i) ||
    value.match(/第\s*(\d{1,3})\s*[话集]/) ||
    value.match(/[\[【]\s*(\d{1,3})(?:v\d+)?\s*[\]】]/) ||
    value.match(/ - (\d{1,3})(?:\s|\[|\.|$)/) ||
    value.match(/\bE(?:P)?\s*(\d{1,3})\b/i);
  return match ? Number(match[1]) || 0 : 0;
}

function feedUpdateInfo(feed) {
  let latest = null;
  let latestTitle = '';
  let maxEpisode = 0;
  for (const item of state.history || []) {
    if ((feed.url && item.config_url === feed.url) || (feed.name && item.config_name === feed.name)) {
      const date = new Date(item.created_at);
      if (!Number.isNaN(date.getTime()) && (!latest || date > latest)) { latest = date; latestTitle = item.title || ''; }
      maxEpisode = Math.max(maxEpisode, parseEpisodeNumber(item.title));
    }
  }
  let episode = maxEpisode;
  if (!episode) episode = parseEpisodeNumber(latestTitle);
  const offset = Number(feed.rename_offset || (state.config.rename || {}).offset || 0);
  const adjustedEpisode = episode > 0 ? episode + offset : episode;
  const total = Number(feed.total_episodes) || 0;
  const completed = feed.completed === true || (total > 0 && adjustedEpisode >= total);
  const scheduled = state.feedWeekdays[feed.url];
  if (scheduled && Number(scheduled.weekday) > 0) {
    const weekday = Number(scheduled.weekday) % 7;
    const date = scheduled.date ? new Date(scheduled.date) : null;
    return { weekday, label: `每周${WEEKDAY_NAMES[weekday]}更新`, date, configured: true, episode: adjustedEpisode, total, completed };
  }
  if (!latest) return { weekday: -1, label: '未更新', date: null, configured: false, episode: adjustedEpisode, total, completed };
  return {
    weekday: latest.getDay(),
    label: `${WEEKDAY_NAMES[latest.getDay()]} · ${String(latest.getMonth() + 1).padStart(2, '0')}-${String(latest.getDate()).padStart(2, '0')}`,
    date: latest,
    configured: false,
    episode: adjustedEpisode,
    total,
    completed,
  };
}

function decoratedFeeds() {
  const feeds = (state.config && state.config.rss) || [];
  return feeds.map((feed, index) => ({ feed, index, update: feedUpdateInfo(feed) })).filter(({ feed, update }) => {
    const source = sourceInfo(feed.source, feed.url).key;
    if (state.feedSource && source !== state.feedSource) return false;
    if (!state.feedQuery) return true;
    const hay = `${feed.name || ''} ${feed.site || ''} ${feed.url || ''} ${feed.filter || ''} ${feed.exclude || ''} ${source}`.toLowerCase();
    return hay.includes(state.feedQuery);
  });
}

function renderFeedList() {
  const container = $('#feed-list');
  if (!container) return;
  const items = decoratedFeeds();
  if (!items.length) {
    container.innerHTML = `<div class="card">${emptyBlock('rss', state.feedQuery || state.feedSource ? '没有匹配的订阅' : '还没有订阅',
      state.feedQuery || state.feedSource ? '换个关键词或来源筛选试试。' : '从 Mikan、AniBT、AnimeGarden 或其他 RSS 添加订阅，开始把新番离线到 115。',
      `<button class="btn btn-primary" id="empty-add-feed">${icon('plus')}添加订阅</button>`)}</div>`;
    const btn = $('#empty-add-feed');
    if (btn) btn.addEventListener('click', () => openAddFlow());
    return;
  }

  const disabledItem = item => item.feed.enabled === false || item.update.completed;
  const enabledItems = items.filter(item => !disabledItem(item));
  const disabledItems = items.filter(disabledItem);
  const selectedItems = state.feedStatus === 'enabled' ? enabledItems : state.feedStatus === 'disabled' ? disabledItems : items;
  const active = selectedItems.filter(item => !disabledItem(item));
  const groups = new Map();
  for (const item of active) {
    const weekday = item.update.weekday;
    const key = weekday < 0 ? '未更新' : String(weekday);
    if (!groups.has(key)) groups.set(key, []);
    groups.get(key).push(item);
  }
  const order = key => key === '未更新' ? 99 : (Number(key) + 6) % 7;
  const sorted = [...groups.entries()].sort((a, b) => order(a[0]) - order(b[0]));
  let html = sorted.map(([key, list]) => {
    const name = key === '未更新' ? '未更新' : `星期${'日一二三四五六'[Number(key)]}`;
    const latest = list.map(item => item.update.date).filter(Boolean).sort((a, b) => b - a)[0];
    const configured = list.some(item => item.update.configured);
    const meta = configured ? `每周更新 · ${list.length} 个订阅` : (latest ? `${latest.getMonth() + 1}月${latest.getDate()}日 · ${list.length} 个订阅` : `${list.length} 个订阅`);
    return `<section class="weekday-group">
      <div class="weekday-title">${name}<span>${meta}</span></div>
      <div class="feed-grid">${list.map(({ feed, index }) => feedCard(feed, index)).join('')}</div>
    </section>`;
  }).join('');
  const showDisabled = state.feedStatus !== 'enabled' && disabledItems.length > 0;
  if (showDisabled) {
    html += `<section class="weekday-group completed-group">
      <div class="weekday-title">未启用<span>${disabledItems.length} 个订阅</span></div>
      <div class="feed-grid">${disabledItems.map(({ feed, index, update }) => feedCard(feed, index, update.completed)).join('')}</div>
    </section>`;
  }
  if (!selectedItems.length) {
    const label = state.feedStatus === 'disabled' ? '暂无未启用订阅' : state.feedStatus === 'enabled' ? '暂无启用订阅' : '没有匹配的订阅';
    html += `<div class="card">${emptyBlock('rss', label, '调整筛选条件，或添加新的订阅。', `<button class="btn btn-sm" id="show-all-feeds">查看全部订阅</button>`)}</div>`;
  }
  container.innerHTML = html;
  const showAll = $('#show-all-feeds');
  if (showAll) showAll.addEventListener('click', () => { state.feedStatus = 'all'; const select = $('#feed-status-filter'); if (select) select.value = 'all'; renderFeedList(); });

  hydrateFeedCovers(container);
  $$('[data-feed-poster]', container).forEach(poster => {
    poster.addEventListener('click', event => {
      const index = Number(poster.dataset.index);
      const feed = state.config && state.config.rss ? state.config.rss[index] : null;
      if (!feed) return;
      const action = event.target.closest('[data-poster-action]')?.dataset.posterAction;
      if (action === 'poster') openPosterModal(feed);
      else openFeedModal(index);
    });
  });
}

function saveCoverCache() {
  try {
    const entries = Object.entries(state.covers).slice(-200);
    localStorage.setItem('anifeed_covers', JSON.stringify(Object.fromEntries(entries)));
  } catch { /* quota */ }
}

function saveWeekdayCache() {
  try { localStorage.setItem('anifeed_weekdays', JSON.stringify(state.feedWeekdays)); } catch { /* quota */ }
}

async function refreshFeedWeekdays() {
  if (state.weekdayLoading) return;
  const feeds = (state.config && state.config.rss) || [];
  const now = Date.now();
  const pending = feeds.filter(feed => {
    if (!feed.tmdb_id) return false;
    const cached = state.feedWeekdays[feed.url];
    return !(cached && cached.expiresAt && cached.expiresAt > now);
  });
  if (!pending.length) return;
  state.weekdayLoading = true;
  const results = await Promise.all(pending.map(async feed => {
    try {
      const type = feed.tmdb_type || 'tv';
      const data = await api(`/api/feeds/weekday?tmdb_id=${encodeURIComponent(feed.tmdb_id)}&tmdb_type=${encodeURIComponent(type)}`);
      return data && data.weekday ? { feed, data } : null;
    } catch { return null; }
  }));
  let changed = false;
  for (const result of results) {
    if (!result) continue;
    state.feedWeekdays[result.feed.url] = { ...result.data, expiresAt: now + 12 * 60 * 60 * 1000 };
    changed = true;
  }
  state.weekdayLoading = false;
  if (changed) {
    saveWeekdayCache();
    if (state.route === 'feeds') renderFeedList();
    if (state.route === 'dashboard') renderDashboard();
  }
}

function applyFeedCover(el, data) {
  if (!el || !data || !data.cover) return;
  const img = $('img', el);
  if (!img) return;
  img.alt = data.title || '';
  img.addEventListener('load', () => el.classList.add('has-image'), { once: true });
  if (img.getAttribute('src') === data.cover) el.classList.add('has-image');
  img.src = data.cover;
  if (data.score > 0) {
    const score = $('.poster-score', el);
    if (score) {
      score.textContent = Number(data.score).toFixed(1);
      score.classList.remove('hidden');
    }
  }
}

async function hydrateFeedCovers(container) {
  const cards = $$('[data-feed-cover]', container);
  const pending = [];
  for (const el of cards) {
    const index = Number(el.dataset.index);
    const feed = state.config && state.config.rss ? state.config.rss[index] : null;
    if (!feed || !feed.url) continue;
    const key = feed.url;
    if (state.covers[key]) {
      applyFeedCover(el, state.covers[key]);
      continue;
    }
    pending.push({ el, feed, key });
  }
  let cursor = 0;
  const worker = async () => {
    while (cursor < pending.length) {
      const task = pending[cursor++];
      try {
        let coverURL = `/api/feeds/cover?url=${encodeURIComponent(task.feed.url)}&source=${encodeURIComponent(sourceInfo(task.feed.source, task.feed.url).key)}`;
        if (task.feed.tmdb_id) coverURL += `&tmdb_id=${encodeURIComponent(task.feed.tmdb_id)}&tmdb_type=${encodeURIComponent(task.feed.tmdb_type || 'tv')}`;
        if (task.feed.bgm_url) coverURL += `&bgm_url=${encodeURIComponent(task.feed.bgm_url)}`;
        const data = await api(coverURL);
        state.covers[task.key] = data;
        saveCoverCache();
        applyFeedCover(task.el, data);
      } catch {
        // Keep the placeholder; a later render can retry.
      }
    }
  };
  await Promise.all(Array.from({ length: 4 }, worker));
}

function feedCard(feed, index, completed = false) {
  const source = sourceInfo(feed.source, feed.url);
  const enabled = feed.enabled !== false;
  const disabled = !enabled || completed;
  const cachedCover = state.covers[feed.url];
  return `<article class="anime-card feed-manage-card ${disabled ? 'disabled' : ''}">
    <div class="anime-poster" data-feed-cover data-feed-poster data-index="${index}" data-source="${source.key}" title="点击编辑订阅">
      <span class="poster-fallback"><img src="/poster-placeholder.jpg?v=2" alt="AniFeed"></span>
      <img alt="${esc(cachedCover ? cachedCover.title || '' : '')}" ${cachedCover && cachedCover.cover ? `src="${esc(cachedCover.cover)}"` : ''} loading="lazy" referrerpolicy="no-referrer">
      <span class="anime-badge">${completed ? '已完成' : disabled ? '未启用' : '追番中'}</span>
      <button class="poster-action poster-action-left" type="button" data-poster-action="poster" title="修改海报">${icon('image')}</button>
      <button class="poster-action poster-action-right" type="button" data-poster-action="edit" title="编辑订阅">${icon('edit')}</button>
      <span class="poster-score hidden"></span>
    </div>
    <div class="anime-body">
      <h4 title="${esc(feed.name)}">${esc(feed.name || '未命名订阅')}</h4>
      <div class="anime-meta">${esc(source.label)} · ${esc(feed.site || '自动识别')}</div>
    </div>
  </article>`;
}

function syncFeedCover(feed, cover) {
  state.covers[feed.url] = { cover };
  const cards = $$('[data-feed-cover]');
  for (const el of cards) {
    const index = Number(el.dataset.index);
    const current = state.config && state.config.rss ? state.config.rss[index] : null;
    if (!current || current.url !== feed.url) continue;
    applyFeedCover(el, { cover });
  }
}

function openPosterModal(feed) {
  const source = sourceInfo(feed.source, feed.url).key;
  const effectiveBgm = feed.bgm_url || bgmURLFromFeedURL(feed.url);
  const coverQuery = `url=${encodeURIComponent(feed.url)}&source=${encodeURIComponent(source)}${feed.tmdb_id ? `&tmdb_id=${encodeURIComponent(feed.tmdb_id)}&tmdb_type=${encodeURIComponent(feed.tmdb_type || 'tv')}` : ''}${effectiveBgm ? `&bgm_url=${encodeURIComponent(effectiveBgm)}` : ''}`;
  openModal({
    title: '管理海报',
    subtitle: feed.name || feed.url,
    body: `
      <div class="poster-manage">
        <div class="poster-preview"><img id="poster-current" alt=""><span id="poster-empty">暂无海报</span></div>
        <div class="poster-actions">
          <p class="hint">自动顺序：本地图片 → BGM → TMDB → RSS 站点。也可以单独从某个来源更新。</p>
          <div class="poster-source-actions">
            <button class="btn btn-sm" id="poster-refresh-bgm" ${effectiveBgm ? '' : 'disabled title="未绑定 BGM 链接"'}>${icon('image')}从 BGM 获取</button>
            <button class="btn btn-sm" id="poster-refresh-tmdb" ${feed.tmdb_id ? '' : 'disabled title="未关联 TMDB"'}>${icon('image')}从 TMDB 获取</button>
            <button class="btn btn-sm" id="poster-refresh-rss">${icon('rss')}从 RSS 站点获取</button>
            <button class="btn btn-sm" id="poster-upload-btn">${icon('upload')}上传本地图片</button>
            <button class="btn btn-sm btn-ghost" id="poster-reset">${icon('refresh')}恢复自动</button>
          </div>
          <input type="file" id="poster-file" accept="image/jpeg,image/png,image/webp,image/gif" class="hidden">
        </div>
      </div>`,
    footer: `<span class="hint">海报保存在 config.toml 同级的 posters 目录中。</span><span class="spacer"></span><button class="btn" data-action="modal-close">关闭</button>`,
    onMount(el, close) {
      const img = el.querySelector('#poster-current');
      const empty = el.querySelector('#poster-empty');
      const file = el.querySelector('#poster-file');
      const load = async (bust = true) => {
        try {
          const data = await api(`/api/feeds/cover?${coverQuery}${bust ? `&_=${Date.now()}` : ''}`);
          if (data.cover) {
            img.src = data.cover;
            img.classList.add('show');
            empty.classList.add('hidden');
            syncFeedCover(feed, data.cover);
          }
        } catch { /* keep empty */ }
      };
      const apply = cover => {
        img.src = cover;
        img.classList.add('show');
        empty.classList.add('hidden');
        syncFeedCover(feed, cover);
      };
      el.querySelector('#poster-upload-btn').addEventListener('click', () => file.click());
      file.addEventListener('change', async () => {
        if (!file.files || !file.files[0]) return;
        const form = new FormData();
        form.append('url', feed.url);
        form.append('file', file.files[0]);
        try {
          const headers = state.token ? { 'X-Access-Token': state.token } : {};
          const res = await fetch('/api/feeds/poster/upload', { method: 'POST', body: form, headers });
          const data = await res.json().catch(() => ({}));
          if (!res.ok) throw new Error(data.error || `上传失败 (${res.status})`);
          apply(data.cover + `&_=${Date.now()}`);
          toast('海报已更新', 'success');
        } catch (err) {
          toast(err.message, 'error', 5200);
        }
      });
      const refreshFrom = async (mode, label) => {
        try {
          const data = await api('/api/feeds/poster/refresh', { method: 'POST', body: {
            url: feed.url, source, mode,
            tmdb_id: feed.tmdb_id || '', tmdb_type: feed.tmdb_type || 'tv',
            bgm_url: effectiveBgm || '',
          }});
          apply(data.cover + `&_=${Date.now()}`);
          toast(`海报已从 ${label} 更新`, 'success');
        } catch (err) {
          toast(err.message, 'error', 5200);
        }
      };
      el.querySelector('#poster-refresh-bgm')?.addEventListener('click', () => refreshFrom('bgm', 'BGM'));
      el.querySelector('#poster-refresh-tmdb')?.addEventListener('click', () => refreshFrom('tmdb', 'TMDB'));
      el.querySelector('#poster-refresh-rss')?.addEventListener('click', () => refreshFrom('rss', 'RSS 站点'));
      el.querySelector('#poster-reset').addEventListener('click', async () => {
        try {
          await api('/api/feeds/poster', { method: 'DELETE', body: { url: feed.url } });
          delete state.covers[feed.url];
          await load(true);
          toast('已恢复自动海报', 'success');
        } catch (err) {
          toast(err.message, 'error', 5200);
        }
      });
      load(false);
    },
  });
}

async function toggleFeed(index) {
  try {
    const body = structuredClone(state.config);
    const feed = body.rss[index];
    feed.enabled = feed.enabled === false;
    const resp = await api('/api/config', { method: 'PUT', body: { config: body } });
    state.config = resp.config;
    state.meta = resp.meta;
    toast(feed.enabled ? '订阅已启用' : '订阅已禁用', 'success');
    renderFeedsPage();
  } catch (err) {
    toast(err.message, 'error', 5200);
  }
}

function parseProviderFields(source, rawUrl) {
  const values = {};
  try {
    const u = new URL(rawUrl);
    const p = u.searchParams;
    values.site = u.hostname;
    if (source === 'mikan') {
      values.bangumiId = p.get('bangumiId') || '';
      values.subgroupId = p.get('subgroupid') || '';
    } else if (source === 'anibt') {
      values.bgmId = p.get('bgmId') || '';
      values.groupSlug = p.get('groupSlug') || '';
    } else if (source === 'animegarden') {
      values.subject = p.get('subject') || '';
      values.fansub = p.get('fansub') || '';
    }
  } catch { /* ignore invalid url */ }
  return values;
}

function buildProviderUrl(source, values, mikanHost) {
  const host = String(mikanHost || 'https://mikanani.me').replace(/\/+$/, '');
  const clean = key => String(values[key] || '').trim();
  if (source === 'mikan') {
    const id = clean('bangumiId');
    if (!id) return '';
    let url = `${host}/RSS/Bangumi?bangumiId=${encodeURIComponent(id)}`;
    if (clean('subgroupId')) url += `&subgroupid=${encodeURIComponent(clean('subgroupId'))}`;
    return url;
  }
  if (source === 'anibt') {
    const id = clean('bgmId');
    if (!id) return '';
    let url = `https://anibt.net/rss/anime.xml?bgmId=${encodeURIComponent(id)}`;
    if (clean('groupSlug')) url += `&groupSlug=${encodeURIComponent(clean('groupSlug'))}`;
    return url;
  }
  if (source === 'animegarden') {
    const subject = clean('subject');
    if (!subject) return '';
    let url = `https://api.animes.garden/feed.xml?subject=${encodeURIComponent(subject)}`;
    if (clean('fansub')) url += `&fansub=${encodeURIComponent(clean('fansub'))}`;
    return url;
  }
  return '';
}

function sourceFieldsHtml(source, values, mikanHost) {
  const input = (label, key, placeholder) => `<label class="field"><span>${label}</span><input class="input" data-provider="${key}" value="${esc(values[key] || '')}" placeholder="${esc(placeholder)}"></label>`;
  if (source === 'mikan') {
    return input('番剧 ID *', 'bangumiId', '例如 3548') + input('字幕组 ID', 'subgroupId', '可选，例如 12');
  }
  if (source === 'anibt') {
    return input('Bangumi 番剧 ID *', 'bgmId', '例如 3548') + input('字幕组 groupSlug', 'groupSlug', '可选，例如 vcb-studio');
  }
  if (source === 'animegarden') {
    return input('Bangumi 番剧 ID *', 'subject', '例如 3548') + input('字幕组名称', 'fansub', '可选，例如 VCB-Studio');
  }
  return `<div class="notice span-2">${icon('info')}<span>填写任意 RSS 地址。支持 enclosure、description 中的 magnet / ed2k，以及直接的 HTTP(S) 下载链接。</span></div>`;
}

function sourceHintHtml(source, mikanHost) {
  const host = String(mikanHost || 'https://mikanani.me').replace(/\/+$/, '');
  if (source === 'mikan') {
    return `在 <a href="${esc(host)}" target="_blank" rel="noopener">Mikan</a> 打开番剧页，RSS 地址中的 <code>bangumiId</code> 与 <code>subgroupid</code> 即所需参数。可直接填写完整 URL 覆盖自动生成结果。`;
  }
  if (source === 'anibt') {
    return `在 <a href="https://anibt.net/" target="_blank" rel="noopener">AniBT</a> 打开番剧页，字幕组链接携带 <code>bgmId</code> 与 <code>groupSlug</code>。`;
  }
  if (source === 'animegarden') {
    return `在 <a href="https://animes.garden/" target="_blank" rel="noopener">AnimeGarden</a> 选择番剧与字幕组，RSS 格式为 <code>feed.xml?subject=...&fansub=...</code>。`;
  }
  return '适用于其他站点提供的标准 RSS，程序会从条目中的 magnet、ed2k 或下载链接提取离线地址。';
}

function browseSiteURL(source) {
  if (source === 'mikan') return ((state.config && state.config.subscription) || {}).mikan_host || 'https://mikanani.me';
  if (source === 'anibt') return 'https://anibt.net/';
  if (source === 'animegarden') return 'https://animes.garden/';
  return '';
}

function openBrowseModal(source, onPick, initialQuery = '') {
  const allowed = ['mikan', 'anibt', 'animegarden', 'other'];
  const view = {
    source: allowed.includes(source) ? source : 'anibt',
    q: String(initialQuery || '').trim(),
    season: '',
    active: '',
    weeks: [],
    expanded: '',
    groups: {},
    data: null,
  };
  const sourceTabs = () => allowed.map(key => `<button type="button" class="source-tab ${view.source === key ? 'active' : ''}" data-browse-source="${key}">${SOURCES[key].label}</button>`).join('');

  openModal({
    title: '浏览订阅',
    subtitle: '搜索番剧、展开字幕组和资源预览，选择后自动填入添加订阅表单',
    size: 'browse',
    body: `
      <div class="source-tabs browse-source-tabs" id="browse-source-tabs">${sourceTabs()}</div>
      <div class="browse-toolbar">
        <div class="search-box">${icon('search')}<input class="input" id="browse-q" value="${esc(String(initialQuery || ''))}" placeholder="搜索番剧标题，或输入 id: 3548"></div>
        <button class="btn" id="browse-search">搜索</button>
        <select class="select hidden" id="browse-season" style="width:auto;min-width:130px"></select>
        <span class="spacer"></span>
        <span class="hint" id="browse-count"></span>
      </div>
      <div class="browse-weeks" id="browse-weeks"></div>
      <div class="browse-body" id="browse-body"><div class="page-loading" style="padding:46px 0"><span class="spinner"></span><span>正在加载…</span></div></div>`,
    footer: `<span class="hint">数据来自站点公开接口；添加后仍可修改名称、过滤规则和保存路径。</span><span class="spacer"></span><button class="btn" data-action="modal-close">关闭</button>`,
    onMount(el, close) {
      const body = $('#browse-body', el);
      const weeksBox = $('#browse-weeks', el);
      const seasonSelect = $('#browse-season', el);
      const count = $('#browse-count', el);
      const qInput = $('#browse-q', el);

      const setLoading = message => {
        body.innerHTML = `<div class="page-loading" style="padding:46px 0"><span class="spinner"></span><span>${esc(message)}</span></div>`;
        weeksBox.innerHTML = '';
        count.textContent = '';
      };

      const renderError = message => {
        weeksBox.innerHTML = '';
        count.textContent = '';
        const site = browseSiteURL(view.source);
        body.innerHTML = `<div class="empty" style="padding:36px 16px">
          <div class="empty-icon">${icon('alert')}</div>
          <strong>加载失败</strong>
          <p>${esc(message)}</p>
          ${site ? `<a class="btn btn-sm" href="${esc(site)}" target="_blank" rel="noopener">${icon('external')}打开站点</a>` : ''}
        </div>`;
      };

      const renderSeason = () => {
        const seasons = view.data.seasons || [];
        if (!seasons.length) {
          seasonSelect.classList.add('hidden');
          return;
        }
        seasonSelect.classList.remove('hidden');
        const current = view.season || view.data.current_season || seasons[0];
        seasonSelect.innerHTML = seasons.map(item => `<option value="${esc(item)}" ${item === current ? 'selected' : ''}>${esc(item)}</option>`).join('');
      };

      const renderWeeks = () => {
        if (!view.weeks.length) {
          weeksBox.innerHTML = '';
          return;
        }
        if (!view.active || !view.weeks.some(week => week.label === view.active)) {
          view.active = view.weeks[0].label;
        }
        weeksBox.innerHTML = view.weeks.map(week => `<button type="button" class="browse-week ${week.label === view.active ? 'active' : ''}" data-week="${esc(week.label)}">${esc(week.label)}<span>${week.items.length}</span></button>`).join('');
        $$('.browse-week', weeksBox).forEach(btn => btn.addEventListener('click', () => {
          view.active = btn.dataset.week;
          view.expanded = '';
          renderWeeks();
          renderItems();
        }));
      };

      const renderItems = () => {
        const week = view.weeks.find(item => item.label === view.active);
        if (!week || !week.items.length) {
          body.innerHTML = `<div class="empty" style="padding:36px 16px"><div class="empty-icon">${icon('search')}</div><strong>没有找到番剧</strong><p>换一个关键词或季度试试。</p></div>`;
          return;
        }
        body.innerHTML = week.items.map(item => {
          const cover = item.cover && /^https?:\/\//.test(item.cover)
            ? `<img class="browse-cover" loading="lazy" referrerpolicy="no-referrer" src="${esc(item.cover)}" alt="">`
            : `<span class="browse-cover fallback">${esc((item.title || '?').slice(0, 1))}</span>`;
          return `<div class="browse-anime" data-anime-id="${esc(item.id)}">
            ${cover}
            <div class="browse-anime-main">
              <div class="browse-anime-title">${esc(item.title)}${item.exists ? '<span class="badge badge-success">已订阅</span>' : ''}</div>
              <div class="browse-anime-meta">${esc(item.meta || '')}${item.rss_count ? ` · ${item.rss_count} 个资源` : ''}</div>
            </div>
            ${item.score > 0 ? `<span class="browse-score">${item.score.toFixed(1)}</span>` : ''}
            <span class="browse-chevron">${icon('arrow')}</span>
          </div>`;
        }).join('');

        $$('.browse-anime', body).forEach(row => {
          row.addEventListener('click', () => {
            const id = row.dataset.animeId;
            const item = week.items.find(entry => entry.id === id);
            if (!item) return;
            if (view.expanded === id) {
              view.expanded = '';
              const open = $('.browse-expand', body);
              if (open) open.remove();
              return;
            }
            view.expanded = id;
            $$('.browse-expand', body).forEach(node => node.remove());
            const expand = document.createElement('div');
            expand.className = 'browse-expand';
            expand.dataset.expandFor = id;
            expand.innerHTML = `<div class="page-loading" style="padding:22px"><span class="spinner"></span><span>正在读取字幕组…</span></div>`;
            row.insertAdjacentElement('afterend', expand);
            loadGroups(item, expand);
          });
        });
      };

      const renderGroups = (item, container, data) => {
        const groups = data.groups || [];
        if (!groups.length) {
          container.innerHTML = `<div class="empty" style="padding:22px"><p>没有找到可用字幕组。</p></div>`;
          return;
        }
        container.innerHTML = groups.map((group, index) => {
          const releases = (group.items || []).map(release => `<div class="browse-release">
            <div class="browse-release-title" title="${esc(release.title)}">${esc(release.title)}</div>
            <span class="browse-release-meta">${fmtSize(release.size)}${release.published_at ? ` · ${fmtDate(release.published_at)}` : ''}</span>
            <button class="btn btn-sm btn-ghost btn-icon-only" data-copy-magnet="${esc(release.magnet)}" title="复制磁力">${icon('copy')}</button>
          </div>`).join('');
          return `<div class="browse-group">
            <div class="browse-group-head">
              <div class="browse-group-name"><strong>${esc(group.name)}</strong><span>${esc(group.meta || '')}${group.latest ? ` · 最近 ${fmtDate(group.latest)}` : ''}</span></div>
              <button class="btn btn-sm btn-primary" data-add-group="${index}">${icon('plus')}添加订阅</button>
              ${releases ? `<button class="btn btn-sm btn-ghost" data-toggle-releases>资源 ${group.items.length}</button>` : ''}
            </div>
            ${releases ? `<div class="browse-releases hidden">${releases}</div>` : ''}
          </div>`;
        }).join('');

        $$('[data-add-group]', container).forEach(btn => btn.addEventListener('click', () => {
          const group = groups[Number(btn.dataset.addGroup)];
          close();
          onPick({ source: view.source, id: item.id, title: item.title, subgroup: group.name, url: group.rss });
        }));
        $$('[data-toggle-releases]', container).forEach(btn => btn.addEventListener('click', () => {
          const list = btn.closest('.browse-group').querySelector('.browse-releases');
          list.classList.toggle('hidden');
        }));
        $$('[data-copy-magnet]', container).forEach(btn => btn.addEventListener('click', () => copyText(btn.dataset.copyMagnet)));
      };

      const loadGroups = async (item, container) => {
        if (view.groups[item.id]) {
          renderGroups(item, container, view.groups[item.id]);
          return;
        }
        try {
          const data = await api(`/api/browse/groups?source=${encodeURIComponent(view.source)}&id=${encodeURIComponent(item.id)}`);
          view.groups[item.id] = data;
          renderGroups(item, container, data);
        } catch (err) {
          container.innerHTML = `<div class="form-error" style="margin:12px">${esc(err.message)}</div>`;
        }
      };

      const renderOther = () => {
        weeksBox.innerHTML = '';
        count.textContent = '';
        body.innerHTML = `<div class="browse-other">
          <div class="notice">${icon('info')}<span>粘贴任意标准 RSS 地址，提交后继续填写订阅名称和过滤规则。</span></div>
          <form class="field-row" id="browse-other-form">
            <input class="input mono" id="browse-other-url" placeholder="https://example.com/rss.xml" autocomplete="off">
            <button class="btn btn-primary" type="submit">${icon('plus')}添加这个 RSS</button>
          </form>
        </div>`;
        const form = $('#browse-other-form', el);
        form.addEventListener('submit', event => {
          event.preventDefault();
          const url = $('#browse-other-url', el).value.trim();
          if (!url) return;
          onPick({ source: 'other', title: '', url });
          close();
        });
      };

      const loadList = async () => {
        if (view.source === 'other') { renderOther(); return; }
        setLoading('正在从 ' + SOURCES[view.source].label + ' 加载…');
        try {
          const params = new URLSearchParams({ source: view.source });
          if (view.q) params.set('q', view.q);
          if (view.season) params.set('season', view.season);
          const data = await api(`/api/browse/list?${params}`);
          view.data = data;
          view.weeks = data.weeks || [];
          view.expanded = '';
          view.groups = {};
          if (data.current_season) view.season = data.current_season;
          count.textContent = data.total ? `${data.total} 部番剧` : '';
          renderSeason();
          renderWeeks();
          renderItems();
        } catch (err) {
          renderError(err.message);
        }
      };

      $('#browse-source-tabs', el).addEventListener('click', event => {
        const btn = event.target.closest('[data-browse-source]');
        if (!btn || btn.dataset.browseSource === view.source) return;
        view.source = btn.dataset.browseSource;
        view.q = '';
        view.season = '';
        view.data = null;
        qInput.value = '';
        $$('.source-tab', el).forEach(tab => tab.classList.toggle('active', tab.dataset.browseSource === view.source));
        loadList();
      });

      const doSearch = () => {
        view.q = qInput.value.trim();
        view.season = '';
        loadList();
      };
      $('#browse-search', el).addEventListener('click', doSearch);
      qInput.addEventListener('keydown', event => { if (event.key === 'Enter') { event.preventDefault(); doSearch(); } });
      seasonSelect.addEventListener('change', () => {
        view.season = seasonSelect.value;
        loadList();
      });

      loadList();
    },
  });
}

function openFolderPicker(onSelect, title = '选择 115 目录') {
  const stack = [{ cid: '0', name: '根目录' }];
  const current = () => stack[stack.length - 1];
  openModal({
    title,
    subtitle: '浏览 115 网盘目录，选择后会写入对应的 cid',
    body: `<div class="folder-crumb" id="folder-crumb"></div>
      <div class="folder-list" id="folder-list"><div class="page-loading" style="padding:32px 0"><span class="spinner"></span><span>正在读取目录…</span></div></div>
      <div class="folder-new">
        <input class="input" id="folder-new-name" placeholder="在当前目录新建文件夹">
        <button class="btn btn-sm" id="folder-new-btn">${icon('plus')}新建</button>
      </div>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="folder-select">${icon('check')}选择当前目录</button>`,
    onMount(el, close) {
      const crumb = $('#folder-crumb', el);
      const list = $('#folder-list', el);
      const newName = $('#folder-new-name', el);

      const renderCrumb = () => {
        crumb.innerHTML = stack.map((item, index) => `<button type="button" class="folder-crumb-item ${index === stack.length - 1 ? 'active' : ''}" data-crumb="${index}">${esc(item.name)}</button>`).join('<span class="folder-sep">/</span>');
        $$('[data-crumb]', crumb).forEach(btn => btn.addEventListener('click', () => {
          const index = Number(btn.dataset.crumb);
          stack.splice(index + 1);
          load(current().cid);
        }));
      };

      const load = async cid => {
        list.innerHTML = `<div class="page-loading" style="padding:32px 0"><span class="spinner"></span><span>正在读取目录…</span></div>`;
        try {
          const data = await api(`/api/fs/list?cid=${encodeURIComponent(cid)}`);
          renderCrumb();
          const entries = data.entries || [];
          if (!entries.length) {
            list.innerHTML = `<div class="empty" style="padding:28px"><div class="empty-icon">${icon('folder')}</div><p>此目录下没有子文件夹</p></div>`;
            return;
          }
          list.innerHTML = entries.map(entry => `<button type="button" class="folder-item" data-dir-cid="${esc(entry.id)}" data-dir-name="${esc(entry.name)}">${icon('folder')}<span>${esc(entry.name)}</span>${icon('arrow')}</button>`).join('');
          $$('[data-dir-cid]', list).forEach(btn => btn.addEventListener('click', () => {
            stack.push({ cid: btn.dataset.dirCid, name: btn.dataset.dirName });
            load(btn.dataset.dirCid);
          }));
        } catch (err) {
          list.innerHTML = `<div class="empty" style="padding:28px"><div class="empty-icon">${icon('lock')}</div><strong>无法读取目录</strong><p>${esc(err.message)}</p><a class="btn btn-sm" href="#/settings">${icon('login')}检查 115 登录</a></div>`;
        }
      };

      $('#folder-select', el).addEventListener('click', () => {
        close();
        onSelect({ cid: current().cid, path: stack.map(item => item.name).join(' / ') });
      });

      $('#folder-new-btn', el).addEventListener('click', async () => {
        const name = newName.value.trim();
        if (!name) { toast('请填写文件夹名称', 'error'); return; }
        try {
          const entry = await api('/api/fs/mkdir', { method: 'POST', body: { parent_cid: current().cid, name } });
          toast('文件夹已创建', 'success');
          stack.push({ cid: entry.id, name: entry.name });
          newName.value = '';
          load(entry.id);
        } catch (err) {
          toast(err.message, 'error', 5200);
        }
      });

      load('0');
    },
  });
}

function openRenameVarsHelp() {
  const vars = [
    ['${title}', '订阅名称（会自动去掉结尾的“ · 字幕组”）', '佐佐木与文鸟小哔 第二季'],
    ['${name}', '同 ${title}，兼容别名', '佐佐木与文鸟小哔 第二季'],
    ['${subgroup}', '从 BT 标题开头的方括号中提取字幕组', 'LoliHouse'],
    ['${season}', '订阅/全局设置的季数', '2'],
    ['${seasonFormat}', '两位季数', '02'],
    ['${episode}', '从 BT 标题解析的集数，支持 x.5', '1 / 1.5'],
    ['${episodeFormat}', '两位集数，x.5 保留小数', '01 / 01.5'],
    ['${itemTitle}', 'RSS 中的完整 BT 标题', '[LoliHouse] ... - 01 [WebRip 1080p]'],
    ['${resolution}', '从标题解析的分辨率', '1080p'],
    ['${tmdbid}', '订阅关联的 TMDB ID；未关联时为空', '127549'],
  ];
  openModal({
    title: '重命名模板变量',
    subtitle: '模板中的变量会替换为实际标题内容',
    size: 'browse',
    body: `<div class="table-wrap"><table class="table vars-table"><thead><tr><th>变量</th><th>说明</th><th>示例</th></tr></thead><tbody>${vars.map(([key, desc, example]) => `<tr><td><code>${esc(key)}</code></td><td>${esc(desc)}</td><td><code>${esc(example)}</code></td></tr>`).join('')}</tbody></table></div>
      <div class="notice" style="margin-top:14px">${icon('info')}<span>单文件重命名会自动保留原扩展名（如 .mkv/.mp4）。解析不到集数时不会重命名，也不会在离线后转移到保存文件夹下移动。季数可在订阅编辑中单独覆盖。</span></div>`,
    footer: `<span class="hint">修改模板后请保存配置。</span><span class="spacer"></span><button class="btn" data-action="modal-close">关闭</button>`,
  });
}

const MATCH_TAG_GROUPS = [
  { label: '分辨率', tags: [
    ['1080p', '(?i)(?:1080p|1920x1080)'],
    ['720p', '(?i)(?:720p|1280x720)'],
    ['4K', '(?i)(?:2160p|4K|3840x2160)'],
  ]},
  { label: '语言', tags: [
    ['简', '(?i)(?:简|CHS|GB|SC)'],
    ['繁', '(?i)(?:繁|CHT|BIG5|TC)'],
    ['日', '(?i)(?:日|JPN|JP)'],
    ['简繁', '(?i)(?:简繁|繁简|CHS&CHT|SC&TC)'],
    ['简日', '(?i)(?:简日|CHS&JPN|CHS_JP)'],
    ['繁日', '(?i)(?:繁日|CHT&JPN|CHT_JP)'],
  ]},
  { label: '字幕', tags: [
    ['内封', '(?i)(?:内封|内挂|软字幕)'],
    ['内嵌', '(?i)(?:内嵌|硬字幕)'],
    ['外挂', '(?i)(?:外挂|外置)'],
  ]},
  { label: '编码', tags: [
    ['HEVC', '(?i)(?:HEVC|H\\.265|x265)'],
    ['AVC', '(?i)(?:AVC|H\\.264|x264)'],
    ['10bit', '(?i)(?:10bit|10-bit)'],
    ['8bit', '(?i)(?:8bit|8-bit)'],
  ]},
  { label: '容器', tags: [
    ['MKV', '(?i)(?:MKV|\\.mkv)'],
    ['MP4', '(?i)(?:MP4|\\.mp4)'],
  ]},
];

function openAddFlow(prefill = null) {
  const selected = new Set();
  openModal({
    title: '匹配偏好',
    subtitle: '选择发布标题中必须包含的标签；不选表示不限制',
    size: 'match',
    body: `<div id="match-tags"></div><div class="rule-preview"><span>将生成的 filter</span><code id="match-filter-preview">（无，匹配全部）</code></div>`,
    footer: `<button class="btn btn-sm" id="match-all">全部（不限制）</button><span class="spacer"></span><button class="btn" id="match-skip">跳过</button><button class="btn btn-primary" id="match-confirm">${icon('check')}应用并继续</button>`,
    onMount(el, close) {
      const tagsBox = el.querySelector('#match-tags');
      const preview = el.querySelector('#match-filter-preview');
      const buildFilter = () => {
        const patterns = [];
        for (const group of MATCH_TAG_GROUPS) {
          for (const [label, pattern] of group.tags) {
            if (selected.has(label)) patterns.push(pattern);
          }
        }
        return patterns.length ? `/^${patterns.map(pattern => `(?=.*(?:${pattern}))`).join('')}.*/` : '';
      };
      const draw = () => {
        tagsBox.innerHTML = MATCH_TAG_GROUPS.map(group => `<div class="match-group"><span>${esc(group.label)}</span><div class="match-tags">${group.tags.map(([label]) => `<button type="button" class="match-chip ${selected.has(label) ? 'active' : ''}" data-match-tag="${esc(label)}">${esc(label)}</button>`).join('')}</div></div>`).join('');
        const filter = buildFilter();
        preview.textContent = filter || '（无，匹配全部）';
        $$('[data-match-tag]', tagsBox).forEach(btn => btn.addEventListener('click', () => {
          const tag = btn.dataset.matchTag;
          if (selected.has(tag)) selected.delete(tag); else selected.add(tag);
          draw();
        }));
      };
      el.querySelector('#match-all').addEventListener('click', () => { selected.clear(); draw(); });
      el.querySelector('#match-confirm').addEventListener('click', () => {
        const filter = buildFilter();
        close();
        openFeedModal(-1, { ...(prefill || {}), filter, exclude: '' });
      });
      el.querySelector('#match-skip').addEventListener('click', () => {
        close();
        openFeedModal(-1, { ...(prefill || {}), filter: '', exclude: '' });
      });
      draw();
    },
  });
}

function formatMissingEpisodes(missing) {
  const values = [...new Set((missing || []).map(Number).filter(n => n > 0))].sort((a, b) => a - b);
  const parts = [];
  let start = 0, prev = 0;
  for (const value of values) {
    if (!start) { start = prev = value; continue; }
    if (value === prev + 1) { prev = value; continue; }
    parts.push(start === prev ? String(start) : `${start}-${prev}`);
    start = prev = value;
  }
  if (start) parts.push(start === prev ? String(start) : `${start}-${prev}`);
  return parts.join('、');
}

function renderRuleResults(container, data) {
  const items = data.items || [];
  container._ruleData = data;
  container._ruleMode = 'matched';
  const draw = () => {
    const mode = container._ruleMode;
    const filtered = items.filter(item => {
      if (mode === 'auto') return item.suggested;
      if (mode === 'all') return true;
      return !item.filtered;
    });
    const visible = filtered.slice(0, 16);
    const summary = data.one_per_episode
      ? `RSS ${items.length} 条 · 规则匹配 ${data.matched} 条 · 将下载 ${Number(data.selected_count) || 0} 条`
      : `RSS ${items.length} 条 · 规则匹配 ${data.matched} 条 · 自动建议 ${Number(data.suggested_count) || 0} 条 · 实际提交 ${Number(data.selected_count) || 0} 条`;
    const rows = visible.map(item => {
      const status = item.status === 'ready'
        ? '<span class="badge badge-success">可重命名</span>'
        : '<span class="badge badge-danger">无法解析</span>';
      const rule = item.filtered
        ? `<span class="badge" title="${esc(item.filter_reason || '规则过滤')}">过滤</span>`
        : '<span class="badge badge-success">通过</span>';
      const suggestion = item.filtered
        ? '—'
        : item.suggested
          ? '<span class="badge badge-primary">建议</span>'
          : `<span class="badge" title="${esc(item.skip_reason || '自动建议忽略')}">忽略</span>`;
      const pick = item.filtered
        ? '—'
        : item.selected
          ? '<span class="badge badge-primary">将下载</span>'
          : `<span class="badge" title="${esc(item.skip_reason || '实际跳过')}">跳过</span>`;
      const renamed = item.renamed
        ? `<code class="preview-renamed">${esc(item.renamed)}</code>`
        : `<span class="form-error">${esc(item.error || '无法生成名称')}</span>`;
      return `<tr class="${item.filtered ? 'rule-filtered' : ''} ${item.suggested ? '' : 'preview-skipped'}">
        <td class="name-cell" data-label="标题"><div class="title">${esc(item.title)}</div></td>
        <td data-label="重命名后">${renamed}</td>
        <td class="numeric" data-label="集数">${item.episode ? esc(String(item.episode)) : '—'}</td>
        <td class="nowrap" data-label="规则">${rule}</td>
        <td class="nowrap" data-label="自动建议">${suggestion}</td>
        <td class="nowrap" data-label="下载">${pick}</td>
        <td class="nowrap" data-label="状态">${status}</td>
      </tr>`;
    }).join('');
    container.innerHTML = `
      <div class="rule-results-head">
        <span>${summary}</span>
        <span class="spacer"></span>
        <button type="button" class="btn btn-sm ${mode === 'auto' ? 'btn-primary' : ''}" id="rule-auto-view">自动</button>
        <button type="button" class="btn btn-sm ${mode === 'all' ? 'btn-primary' : ''}" id="rule-show-all">显示 RSS 全部内容</button>
      </div>
      ${(data.missing || []).length ? `<div class="notice" style="margin:6px 0 8px">${icon('alert')}<span>缺少第 <b>${esc(formatMissingEpisodes(data.missing))}</b> 集</span></div>` : ''}
      <div class="hint" style="margin:6px 0 8px">自动：只在规则匹配的条目中按集数分组，再按分辨率、语言、内封/内嵌、编码等评分选一个；被规则过滤的条目不参与自动。未限定语言时默认优先简体，选择繁体时自动在繁体候选里选。</div>
      ${!data.one_per_episode ? `<div class="notice" style="margin:6px 0 8px">${icon('info')}<span>当前仅识别并记录，不实际跳过；实际会提交全部规则匹配项。开启「每集只下载一个版本」后才按自动建议提交。</span></div>` : ''}
      ${visible.length ? `<div class="table-wrap rule-result-table-wrap">
        <table class="table rename-preview-table">
          <thead><tr><th>RSS / BT 标题</th><th>重命名后</th><th>集数</th><th>规则</th><th>自动建议</th><th>下载</th><th>状态</th></tr></thead>
          <tbody>${rows}</tbody>
        </table>
      </div>` : `<div class="empty" style="padding:22px"><span class="hint">当前视图没有条目。</span></div>`}
      ${filtered.length > visible.length ? `<div class="hint" style="margin-top:8px">还有 ${filtered.length - visible.length} 条未显示</div>` : ''}`;
    container.querySelector('#rule-auto-view')?.addEventListener('click', () => {
      container._ruleMode = container._ruleMode === 'auto' ? 'matched' : 'auto';
      draw();
    });
    container.querySelector('#rule-show-all')?.addEventListener('click', () => {
      container._ruleMode = container._ruleMode === 'all' ? 'matched' : 'all';
      draw();
    });
  };
  draw();
}

async function openItemPreview(feed) {
  if (!feed.url) {
    toast('请先填写 RSS 地址', 'error');
    return;
  }
  const backupCount = (feed.subgroups || []).filter(group => group.enabled !== false && group.url).length;
  openModal({
    title: '预览',
    subtitle: `${feed.name || feed.url} · 主 RSS${backupCount ? ` + ${backupCount} 个备用 RSS` : ''}`,
    size: 'rules',
    body: `<div class="page-loading" style="padding:40px 0"><span class="spinner"></span><span>正在读取规则匹配结果…</span></div>`,
    footer: `<button class="btn" data-action="modal-close">关闭</button><span class="spacer"></span><button class="btn btn-primary" id="feed-preview-save">${icon('save')}保存选择</button>`,
    onMount(el, close) {
      const body = el.querySelector('.modal-body');
      const backupSources = (feed.subgroups || []).filter(group => group.enabled !== false && group.url);
      const sourceBar = `<details class="preview-source-panel">
        <summary>RSS 来源 <span class="badge badge-info">${1 + backupSources.length}</span><span class="hint">主 RSS${backupCount ? ` + ${backupCount} 个备用 RSS` : ''}</span></summary>
        <div class="item-preview-sources">
          <div class="item-preview-source"><span class="rss-role-badge main">主 RSS</span>${(feed.subgroup_name || urlSubgroupName(feed.url)) ? `<strong>${esc(feed.subgroup_name || urlSubgroupName(feed.url))}</strong>` : ''}<code title="${esc(feed.url || '')}">${esc(feed.url || '')}</code></div>
          ${backupSources.map(group => `<div class="item-preview-source"><span class="rss-role-badge backup">备用 RSS</span><strong>${esc(subgroupLabel(group.name, feed.name) || '备用字幕组')}</strong><code title="${esc(group.url)}">${esc(group.url)}</code></div>`).join('')}
        </div>
      </details>`;
      const previewSources = [
        { role: 'main', subgroup: '', url: feed.url, source: feed.source, filter: feed.filter || '', exclude: feed.exclude || '', season: feed.rename_season || 0, offset: feed.rename_offset || 0 },
        ...backupSources.map(group => ({ role: 'backup', subgroup: group.name || '', url: group.url, source: feed.source, filter: group.filter || '', exclude: group.exclude || '', season: group.season || feed.rename_season || 0, offset: group.offset || 0 })),
      ];
      Promise.all(previewSources.map((source, order) => {
        const params = new URLSearchParams({
          name: feed.name || '', url: source.url || '', source: source.source || '',
          filter: source.filter, exclude: source.exclude,
          rename_season: String(source.season || 0), rename_offset: String(source.offset || 0),
          tmdb_id: feed.tmdb_id || '', tmdb_type: feed.tmdb_type || 'tv',
        });
        return api(`/api/rename/preview?${params}`).then(data => ({ source, order, data })).catch(error => ({ source, order, error: error.message }));
      })).then(results => {
        const failures = results.filter(result => result.error);
        const items = results.flatMap(result => ((result.data || {}).items || [])
          .map(item => ({ ...item, source_url: result.source.url, source_role: result.source.role, source_subgroup: result.source.subgroup, source_order: result.order })));
        items.sort((a, b) => {
          const ea = Number(a.episode) > 0 ? Number(a.episode) : 1e9;
          const eb = Number(b.episode) > 0 ? Number(b.episode) : 1e9;
          if (ea !== eb) return ea - eb;
          if (a.source_order !== b.source_order) return a.source_order - b.source_order;
          return String(a.title || '').localeCompare(String(b.title || ''));
        });
        if (!items.length && failures.length) throw new Error(failures.map(result => `${result.source.role === 'main' ? '主 RSS' : '备用 RSS'}：${result.error}`).join('；'));
        const allowed = new Map(items.map(item => [item.key, item.allowed]));
        const covered = new Map();
        for (const item of items) {
          const episode = Number(item.episode);
          if (!episode || item.filtered || !allowed.get(item.key)) continue;
          if (!covered.has(episode)) covered.set(episode, item);
        }
        if (!(state.config.lab || {}).multi_rss_coexist) {
          for (const item of items) {
            const episode = Number(item.episode);
            const first = covered.get(episode);
            if (!first || !episode || item.filtered || item.decided) continue;
            if (item.source_order > first.source_order) {
              item.coveredBy = first;
              allowed.set(item.key, false);
            }
          }
        }
        const selectedKeys = new Set();
        const draw = () => {
          const rows = items.map((item, index) => {
            const isAllowed = allowed.get(item.key);
            return `<tr class="${item.filtered ? 'preview-filtered-row' : ''}">
              <td data-label="选择"><input type="checkbox" class="preview-check" data-key="${esc(item.key)}" ${selectedKeys.has(item.key) ? 'checked' : ''}></td>
              <td class="nowrap" data-label="是否下载">${item.filtered ? '<span class="badge badge-warning">过滤</span>' : (item.coveredBy && !isAllowed ? `<span class="badge" title="已由 ${esc(item.coveredBy.source_subgroup || '主 RSS')} 覆盖">否</span>` : `<button type="button" class="badge ${isAllowed ? 'badge-primary' : ''}" data-preview-toggle="${index}">${isAllowed ? '是' : '否'}</button>`)}</td>
              <td class="nowrap" data-label="已下载">${item.downloaded ? '<span class="badge badge-success">是</span>' : '<span class="badge">否</span>'}</td>
              <td class="nowrap" data-label="主 RSS"><span class="badge ${item.source_role === 'main' ? 'badge-success' : 'badge'}">${item.source_role === 'main' ? '是' : '否'}</span></td>
              <td class="nowrap" data-label="字幕组">${esc(item.subgroup || '未知字幕组')}</td>
              <td class="name-cell" data-label="标题"><div class="title">${esc(item.title)}</div></td>
              <td data-label="重命名">${item.filtered ? `<span class="hint">${esc(item.filter_reason || '被当前规则过滤')}</span>` : (item.renamed ? `<code class="preview-renamed">${esc(item.renamed)}</code>` : '<span class="form-error">无法解析</span>')}</td>
            </tr>`;
          }).join('');
          const allowedCount = items.filter(item => !item.filtered && allowed.get(item.key)).length;
          const downloadedCount = items.filter(item => item.downloaded).length;
          const filteredCount = items.filter(item => item.filtered).length;
          const coveredCount = items.filter(item => item.coveredBy && !allowed.get(item.key)).length;
          body.innerHTML = `
            ${sourceBar}
            ${failures.length ? `<div class="notice warn" style="margin-bottom:8px">${icon('alert')}<span>部分 RSS 抓取失败：${failures.map(result => esc((result.source.role === 'main' ? '主 RSS' : (result.source.subgroup || '备用 RSS')) + '：' + result.error)).join('；')}</span></div>` : ''}
            <div class="preview-summary">
              <span><b>${items.length}</b> 项</span>
              <span>允许 <b>${allowedCount}</b></span>
              <span>禁止 <b>${items.length - allowedCount}</b></span>
              <span>已下载 <b>${downloadedCount}</b></span>
              ${filteredCount ? `<span>过滤 <b>${filteredCount}</b></span>` : ''}
              ${coveredCount ? `<span>已覆盖 <b>${coveredCount}</b></span>` : ''}
              <span class="spacer"></span>
              <span class="preview-selected">已选 <b>${selectedKeys.size}</b></span>
            </div>
            <div class="preview-toolbar preview-actions">
              <button class="btn btn-sm" id="preview-select-all">${selectedKeys.size === items.length && items.length ? '取消全选' : '全选'}</button>
              <button class="btn btn-sm" id="preview-allow">${icon('check')}允许</button>
              <button class="btn btn-sm" id="preview-deny">${icon('x')}禁止</button>
              <button class="btn btn-sm btn-danger" id="preview-delete-downloaded">${icon('trash')}删除已下载</button>
              <label class="check preview-delete-option"><input type="checkbox" id="preview-delete-files">同时删除云盘文件与离线记录</label>
            </div>
            ${items.length ? `<div class="table-wrap"><table class="table rename-preview-table">
              <thead><tr><th></th><th>是否下载</th><th>已下载</th><th>主 RSS</th><th>字幕组</th><th>标题</th><th>重命名</th></tr></thead>
              <tbody>${rows}</tbody>
            </table></div>` : `<div class="empty" style="padding:30px"><span class="hint">当前规则没有匹配到条目。</span></div>`}
            <div class="hint" style="margin-top:8px">这里调整的是规则匹配后的条目；保存后运行订阅会优先使用这些选择。已下载来自历史记录。</div>`;
          body.querySelectorAll('[data-preview-toggle]').forEach(button => button.addEventListener('click', () => {
            const index = Number(button.dataset.previewToggle);
            const key = items[index].key;
            const next = !allowed.get(key);
            allowed.set(key, next);
            if (next) items[index].coveredBy = null;
            draw();
          }));
          body.querySelectorAll('.preview-check').forEach(input => input.addEventListener('change', () => {
            if (input.checked) selectedKeys.add(input.dataset.key); else selectedKeys.delete(input.dataset.key);
            draw();
          }));
          body.querySelector('#preview-select-all')?.addEventListener('click', () => {
            if (selectedKeys.size === items.length) selectedKeys.clear();
            else items.forEach(item => selectedKeys.add(item.key));
            draw();
          });
          const setSelected = value => {
            if (!selectedKeys.size) { toast('请先勾选条目', 'error'); return; }
            for (const key of selectedKeys) {
              allowed.set(key, value);
              if (value) items.forEach(item => { if (item.key === key) item.coveredBy = null; });
            }
            draw();
          };
          body.querySelector('#preview-allow')?.addEventListener('click', () => setSelected(true));
          body.querySelector('#preview-deny')?.addEventListener('click', () => setSelected(false));
          body.querySelector('#preview-delete-downloaded')?.addEventListener('click', async () => {
            const targets = items.filter(item => selectedKeys.has(item.key));
            if (!targets.length) { toast('请先勾选条目', 'error'); return; }
            const deleteFiles = $('#preview-delete-files').checked;
            const ok = await confirmDialog({
              title: '删除下载记录',
              message: `将强制删除选中的 ${targets.length} 条下载历史，并删除对应的后台重命名/移动记录。<br><br>${deleteFiles ? '<strong>同时删除 115 云盘文件/文件夹和 115 云下载离线记录，不可恢复。</strong>' : '115 网盘中的文件和 115 云下载离线记录不会被删除。'}`,
              confirmText: '删除记录',
            });
            if (!ok) return;
            try {
              const result = await api('/api/feeds/preview/delete-history', { method: 'POST', body: {
                magnets: targets.map(item => item.magnet).filter(Boolean),
                delete_files: deleteFiles,
              } });
              targets.forEach(item => { item.downloaded = false; });
              toast(`已删除 ${Number(result.deleted) || 0} 条下载历史${Number(result.files_deleted) ? `，云盘文件 ${result.files_deleted} 个` : ''}${Number(result.files_skipped) ? `；${result.files_skipped} 个未找到 115 文件` : ''}`, 'success');
              draw();
            } catch (err) {
              toast(err.message, 'error', 5200);
            }
          });
        };
        draw();
        el.querySelector('#feed-preview-save').addEventListener('click', async () => {
          try {
            const bySource = new Map();
            for (const item of items) {
              const url = item.source_url || feed.url;
              if (!bySource.has(url)) bySource.set(url, []);
              bySource.get(url).push({ key: item.key, allowed: !!allowed.get(item.key) });
            }
            await Promise.all([...bySource.entries()].map(([url, decisions]) => api('/api/feeds/preview/decisions', { method: 'POST', body: { url, decisions } })));
            toast(`下载选择已保存（${bySource.size} 个 RSS）`, 'success');
            close();
          } catch (err) {
            toast(err.message, 'error', 5200);
          }
        });
      }).catch(err => {
        body.innerHTML = `<div class="form-error">${esc(err.message)}</div>`;
      });
    },
  });
}

async function openRulesModal(onApply, getPreviewContext = null) {
  let data = state.rules;
  if (!data) {
    try {
      data = await api('/api/rules');
      state.rules = data;
    } catch (err) {
      toast(err.message, 'error', 5200);
      return;
    }
  }
  const selected = {};
  const groups = new Set();
  for (const category of data.categories || []) selected[category.id] = '';
  const opposites = {
    language: { chs: ['cht', 'bilingual'], cht: ['chs', 'bilingual'], bilingual: [] },
    subtitle: { softsub: ['hardsub', 'external'], hardsub: ['softsub', 'external'], external: ['softsub', 'hardsub'] },
    container: { mkv: ['mp4'], mp4: ['mkv'] },
    resolution: { r1080: ['r720', 'r2160'], r720: ['r1080', 'r2160'], r2160: ['r1080', 'r720'] },
  };
  const optionById = {};
  for (const category of data.categories || []) {
    for (const option of category.options) optionById[`${category.id}:${option.id}`] = option;
  }
  const build = () => {
    const includes = [];
    const excludes = [];
    for (const category of data.categories || []) {
      const current = selected[category.id];
      if (!current) continue;
      const option = category.options.find(item => item.id === current);
      if (option) includes.push(option.pattern);
      for (const otherId of (opposites[category.id] || {})[current] || []) {
        const other = category.options.find(item => item.id === otherId);
        if (other) excludes.push(other.pattern);
      }
    }
    for (const group of data.groups || []) {
      if (groups.has(group.id)) excludes.push(group.pattern);
    }
    const filter = includes.length ? `/^${includes.map(pattern => `(?=.*(?:${pattern}))`).join('')}.*/` : '';
    const exclude = excludes.length ? `/(?:${excludes.map(pattern => `(?:${pattern})`).join('|')})/` : '';
    return { filter, exclude };
  };

  const initialContext = getPreviewContext ? getPreviewContext() : {};
  openModal({
    title: '规则与重命名',
    subtitle: '选择规则后实时显示匹配条目、重命名结果和同集优选；新规则会叠加到当前 filter/exclude',
    size: 'rules',
    body: `
      <div class="notice">${icon('info')}<span>规则库来源：<a href="${esc(data.source)}" target="_blank" rel="noopener">${esc(data.source)}</a>。选择偏好后会生成“必须包含”的 filter，并把同类的其他选项加入 exclude。</span></div>
      <div class="form-grid rules-manual">
        <label class="field"><span>白名单 filter</span><textarea class="textarea mono" id="rules-base-filter" rows="2" placeholder="标题包含即离线，或 /正则/；支持多行和 {{字幕组}}:正则">${esc(initialContext.filter || '')}</textarea><span class="hint">多行规则为 AND；规则库生成的内容会追加到这里。</span></label>
        <label class="field"><span>黑名单 exclude</span><textarea class="textarea mono" id="rules-base-exclude" rows="2" placeholder="标题包含即跳过，或 /正则/；支持多行">${esc(initialContext.exclude || '')}</textarea><span class="hint">多行规则为 OR；规则库生成的内容会追加到这里。</span></label>
      </div>
      <div id="rules-body"></div>
      <div class="rule-preview"><span>生成规则</span><code id="rules-filter-preview">—</code><code id="rules-exclude-preview">—</code></div>
      <div class="rule-results" id="rules-results"><span class="hint">选择规则后自动显示 RSS 匹配结果</span></div>`,
    footer: `<button class="btn" id="rules-best">${icon('filter')}自动应用最佳规则</button><button class="btn" id="rules-clear">清空</button><span class="spacer"></span><button class="btn" data-action="modal-close">取消</button><button class="btn btn-primary" id="rules-apply">${icon('check')}应用到订阅</button>`,
    onMount(el, close) {
      const body = el.querySelector('#rules-body');
      const filterPreview = el.querySelector('#rules-filter-preview');
      const excludePreview = el.querySelector('#rules-exclude-preview');
      const resultsBox = el.querySelector('#rules-results');
      let resultTimer = null;
      let resultToken = 0;
      let autoApplied = false;

      const loadResults = async () => {
        const context = getPreviewContext ? getPreviewContext() : null;
        if (!context || !context.url) {
          resultsBox.innerHTML = `<span class="hint">填写 RSS 地址后，这里会直接显示匹配到的条目。</span>`;
          return;
        }
        const token = ++resultToken;
        const { filter, exclude } = build();
        const baseFilter = (el.querySelector('#rules-base-filter')?.value || '').trim();
        const baseExclude = (el.querySelector('#rules-base-exclude')?.value || '').trim();
        const effectiveFilter = [baseFilter, filter].filter(value => value && value.trim()).join('\n');
        const effectiveExclude = [baseExclude, exclude].filter(value => value && value.trim()).join('\n');
        resultsBox.innerHTML = `<div class="page-loading" style="padding:18px"><span class="spinner"></span><span>正在匹配 RSS 条目…</span></div>`;
        try {
          const params = new URLSearchParams({
            name: context.name || '', url: context.url || '', source: context.source || '',
            filter: effectiveFilter, exclude: effectiveExclude,
            rename_season: String(context.rename_season || 0),
            rename_offset: String(context.rename_offset || 0),
            tmdb_id: context.tmdb_id || '', tmdb_type: context.tmdb_type || 'tv',
          });
          const data = await api(`/api/rename/preview?${params}`);
          if (token !== resultToken) return;
          renderRuleResults(resultsBox, data);
          if (!autoApplied && (state.config.offline || {}).auto_rules) {
            autoApplied = true;
            applyBestRules(true);
          }
        } catch (err) {
          if (token !== resultToken) return;
          resultsBox.innerHTML = `<div class="form-error" style="margin:10px 0">${esc(err.message)}</div>`;
        }
      };

      const scheduleResults = () => {
        clearTimeout(resultTimer);
        resultTimer = setTimeout(loadResults, 450);
      };
      ['rules-base-filter', 'rules-base-exclude'].forEach(id => {
        const input = el.querySelector('#' + id);
        if (input) input.addEventListener('input', scheduleResults);
      });

      const preferenceOrder = {
        language: ['chs', 'bilingual', 'cht'],
        subtitle: ['softsub', 'hardsub', 'external'],
        container: ['mkv', 'mp4'],
        resolution: ['r1080', 'r2160', 'r720'],
      };
      const applyBestRules = (silent = false) => {
        if (!resultsBox._ruleData) {
          toast('正在读取 RSS，请稍候再试', 'error');
          return;
        }
        const titles = (resultsBox._ruleData.items || []).map(item => item.title);
        const chosen = [];
        for (const category of data.categories || []) {
          const configured = (state.config && state.config.offline && state.config.offline[`auto_${category.id}`]) || [];
          const order = configured.length ? configured.filter(id => category.options.some(option => option.id === id)) : (preferenceOrder[category.id] || []);
          let picked = '';
          for (const optionId of order) {
            const option = category.options.find(item => item.id === optionId);
            if (!option) continue;
            let available = titles.length === 0;
            if (!available) {
              try {
                const re = new RegExp(option.pattern.replace(/^\(\?i\)/, ''), 'i');
                available = titles.some(title => re.test(title));
              } catch (err) {
                available = false;
              }
            }
            if (available) { picked = optionId; chosen.push(option.label); break; }
          }
          selected[category.id] = picked;
        }
        draw();
        if (!silent) toast(chosen.length ? `已自动选择：${chosen.join(' / ')}，点击「应用到订阅」保存` : '当前 RSS 没有可匹配的规则选项', chosen.length ? 'success' : 'error');
      };

      const draw = () => {
        const categories = (data.categories || []).map(category => {
          const chips = [`<button type="button" class="rule-chip ${!selected[category.id] ? 'active' : ''}" data-rule-cat="${category.id}" data-rule-id="">不限制</button>`];
          for (const option of category.options) {
            chips.push(`<button type="button" class="rule-chip ${selected[category.id] === option.id ? 'active' : ''}" data-rule-cat="${category.id}" data-rule-id="${option.id}" title="${esc(option.description || '')}">${esc(option.label)}</button>`);
          }
          return `<div class="rule-category"><div class="rule-category-label">${esc(category.label)}</div><div class="rule-chips">${chips.join('')}</div></div>`;
        }).join('');
        const groupChips = (data.groups || []).map(group => `<button type="button" class="rule-chip group ${groups.has(group.id) ? 'active' : ''}" data-group="${group.id}" title="${esc(group.description || '')}">${esc(group.label)}</button>`).join('');
        body.innerHTML = `<div class="rule-presets"><span>快速组合：</span>
          <button type="button" class="rule-chip" data-preset="chs-soft-mkv">简体 + 内封 + MKV</button>
          <button type="button" class="rule-chip" data-preset="chs-hard-mp4">简体 + 内嵌 + MP4</button>
          <button type="button" class="rule-chip" data-preset="cht-soft-mkv">繁体 + 内封 + MKV</button>
        </div>
        ${categories}
        <div class="rule-category"><div class="rule-category-label">字幕组排除预设</div><div class="rule-chips">${groupChips}</div></div>`;
        const { filter, exclude } = build();
        filterPreview.textContent = filter || '（无）';
        excludePreview.textContent = exclude || '（无）';
        scheduleResults();
        $$('[data-rule-cat]', body).forEach(btn => btn.addEventListener('click', () => {
          selected[btn.dataset.ruleCat] = btn.dataset.ruleId;
          draw();
        }));
        $$('[data-group]', body).forEach(btn => btn.addEventListener('click', () => {
          const id = btn.dataset.group;
          if (groups.has(id)) groups.delete(id); else groups.add(id);
          draw();
        }));
        $$('[data-preset]', body).forEach(btn => btn.addEventListener('click', () => {
          for (const key of Object.keys(selected)) selected[key] = '';
          groups.clear();
          const preset = btn.dataset.preset;
          if (preset === 'chs-soft-mkv') { selected.language = 'chs'; selected.subtitle = 'softsub'; selected.container = 'mkv'; }
          if (preset === 'chs-hard-mp4') { selected.language = 'chs'; selected.subtitle = 'hardsub'; selected.container = 'mp4'; }
          if (preset === 'cht-soft-mkv') { selected.language = 'cht'; selected.subtitle = 'softsub'; selected.container = 'mkv'; }
          draw();
        }));
      };

      el.querySelector('#rules-best').addEventListener('click', applyBestRules);
      el.querySelector('#rules-clear').addEventListener('click', () => {
        for (const key of Object.keys(selected)) selected[key] = '';
        groups.clear();
        draw();
      });
      el.querySelector('#rules-apply').addEventListener('click', () => {
        const { filter, exclude } = build();
        const baseFilter = (el.querySelector('#rules-base-filter')?.value || '').trim();
        const baseExclude = (el.querySelector('#rules-base-exclude')?.value || '').trim();
        const finalFilter = [baseFilter, filter].filter(value => value && value.trim()).join('\n');
        const finalExclude = [baseExclude, exclude].filter(value => value && value.trim()).join('\n');
        close();
        onApply(finalFilter, finalExclude);
        toast('已应用规则，保存订阅后生效', 'success');
      });
      draw();
    },
  });
}

function openTMDBSearch(options = {}) {
  const state = { type: options.type === 'movie' ? 'movie' : 'tv', query: options.query || '' };
  openModal({
    title: '关联 TMDB',
    subtitle: '搜索并选择对应的 TMDB 条目',
    size: 'browse',
    body: `
      <div class="browse-toolbar">
        <div class="search-box">${icon('search')}<input class="input" id="tmdb-q" placeholder="输入番剧名称" value="${esc(state.query)}"></div>
        <select class="select" id="tmdb-type" style="width:auto;min-width:100px"><option value="tv" ${state.type === 'tv' ? 'selected' : ''}>TV</option><option value="movie" ${state.type === 'movie' ? 'selected' : ''}>Movie</option></select>
        <button class="btn" id="tmdb-search">搜索</button>
      </div>
      <div id="tmdb-results"><div class="empty" style="padding:32px"><p>输入名称后搜索 TMDB。</p></div></div>`,
    footer: `<span class="hint">需要先在设置中填写 TmdbApiKey；选择后只关联 ID，不会修改订阅名称。</span><span class="spacer"></span><button class="btn" data-action="modal-close">关闭</button>`,
    onMount(el, close) {
      const results = el.querySelector('#tmdb-results');
      const input = el.querySelector('#tmdb-q');
      const type = el.querySelector('#tmdb-type');
      const search = async () => {
        const query = input.value.trim();
        if (!query) { toast('请输入搜索关键词', 'error'); return; }
        results.innerHTML = `<div class="page-loading" style="padding:36px 0"><span class="spinner"></span><span>正在搜索 TMDB…</span></div>`;
        try {
          const data = await api(`/api/tmdb/search?q=${encodeURIComponent(query)}&type=${encodeURIComponent(type.value)}`);
          renderTMDBResults(results, data.results || [], item => {
            close();
            options.onSelect?.(item);
          });
        } catch (err) {
          results.innerHTML = `<div class="empty" style="padding:30px"><div class="empty-icon">${icon('alert')}</div><strong>搜索失败</strong><p>${esc(err.message)}</p></div>`;
        }
      };
      el.querySelector('#tmdb-search').addEventListener('click', search);
      input.addEventListener('keydown', event => { if (event.key === 'Enter') { event.preventDefault(); search(); } });
      if (state.query) search();
    },
  });
}

function renderTMDBResults(container, results, onSelect) {
  if (!results.length) {
    container.innerHTML = `<div class="empty" style="padding:32px"><div class="empty-icon">${icon('search')}</div><strong>没有找到结果</strong><p>换一个关键词或切换 TV / Movie。</p></div>`;
    return;
  }
  container.innerHTML = `<div class="tmdb-results">${results.map((item, index) => {
    const date = item.first_air_date || item.release_date || '';
    const year = date ? date.slice(0, 4) : '';
    return `<div class="tmdb-result">
      ${item.poster ? `<img class="tmdb-poster" loading="lazy" referrerpolicy="no-referrer" src="${esc(item.poster)}" alt="">` : `<span class="tmdb-poster fallback">${esc((item.name || item.title || '?').slice(0, 1))}</span>`}
      <div class="tmdb-main"><div class="tmdb-name">${esc(item.name || item.title || '')} ${year ? `<span>${esc(year)}</span>` : ''}</div><div class="tmdb-original">${esc(item.original_name || '')}${item.vote_average ? ` · ${Number(item.vote_average).toFixed(1)}` : ''}</div><div class="tmdb-overview">${esc(item.overview || '')}</div></div>
      <button class="btn btn-sm btn-primary" data-tmdb-pick="${index}">${icon('check')}选择</button>
    </div>`;
  }).join('')}</div>`;
  $$('[data-tmdb-pick]', container).forEach(btn => btn.addEventListener('click', () => onSelect(results[Number(btn.dataset.tmdbPick)])));
}

function openFeedModal(index, prefill = null) {
  const editing = index >= 0;
  const preselected = !editing && !!prefill?.url;
  const feed = editing ? { ...state.config.rss[index] } : {
    site: '',
    name: prefill?.name || '',
    url: prefill?.url || '',
    source: prefill?.source || 'mikan',
    enabled: true,
    cid: prefill?.cid || '',
    savepath: prefill?.savepath || '',
    filter: prefill?.filter || '',
    exclude: prefill?.exclude || '',
    expiration: 0,
    tmdb_id: prefill?.tmdb_id || '',
    tmdb_type: prefill?.tmdb_type || 'tv',
  };
  const initialSource = sourceInfo(feed.source, feed.url).key;
  const mikanHost = ((state.config.subscription || {}).mikan_host || 'https://mikanani.me');
  const mainSubgroup = feed.subgroup_name || urlSubgroupName(feed.url);
  const initialBgm = feed.bgm_url || bgmURLFromFeedURL(feed.url);
  const local = {
    source: initialSource,
    subgroupName: mainSubgroup,
    bgmAuto: !feed.bgm_url,
    completed: !!feed.completed,
    provider: parseProviderFields(initialSource, feed.url),
    urlTouched: (editing || preselected) && !!feed.url,
    subgroups: Array.isArray(feed.subgroups) ? feed.subgroups : [],
    filter: feed.filter || '',
    exclude: feed.exclude || '',
    cid: feed.cid || '',
    savepath: feed.savepath || '',
  };
  const labMultiSubgroup = !!(state.config.lab && state.config.lab.multi_subgroup);

  const tabs = Object.entries(SOURCES).map(([key, info]) => `<button type="button" class="source-tab ${initialSource === key ? 'active' : ''}" data-source="${key}">${info.label}</button>`).join('');

  openModal({
    title: editing ? '编辑订阅' : '添加订阅',
    subtitle: editing ? '修改订阅设置' : (preselected ? '已选择来源，确认订阅设置后保存' : '选择来源后填写 ID，RSS 地址会自动生成；也可以直接粘贴完整 RSS 地址'),
    size: 'wide feed-edit',
    body: `
      <div class="source-tabs ${editing || preselected ? 'hidden' : ''}" id="source-tabs">${tabs}</div>
      <div class="source-help ${editing || preselected ? 'hidden' : ''}" id="source-help"></div>
      <div class="browse-cta ${editing || preselected ? 'hidden' : ''}" id="feed-browse-row"></div>
      <div class="form-grid provider-fields ${editing || preselected ? 'hidden' : ''}" id="provider-fields"></div>
      <div class="feed-edit-tabs" id="feed-edit-tabs">
        <button type="button" class="active" data-feed-tab="basic">基本</button>
        <button type="button" data-feed-tab="backup">备用 RSS</button>
        <button type="button" data-feed-tab="custom">自定义</button>
      </div>
      <form id="feed-form" class="feed-form" autocomplete="off">
        <section class="feed-tab-panel active" data-feed-panel="basic">
          <div class="feed-form-section">
            <div class="feed-form-section-head">基本信息</div>
            <div class="form-grid">
              <label class="field"><span>标题 *</span><input class="input" id="feed-name" value="${esc(feed.name)}" placeholder="例如：药屋少女的呢喃"></label>
              <label class="field"><span>启用订阅</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="feed-enabled" ${feed.enabled !== false ? 'checked' : ''}><span></span></span></span></label>
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">RSS 来源</div>
            <div class="form-grid">
              <div class="field span-2 feed-main-rss">
                <span class="field-label">主 RSS <span class="rss-role-badge main">主 RSS</span><button class="btn btn-sm btn-ghost" type="button" id="feed-main-rules">规则</button></span>
                <div class="feed-main-rss-row">
                  <input class="input" id="feed-subgroup-name" value="${esc(mainSubgroup)}" placeholder="字幕组名称，例如 LoliHouse">
                  <input class="input mono" id="feed-url" value="${esc(feed.url)}" placeholder="https://...">
                </div>
              </div>
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">TMDB 与 BGM</div>
            <div class="form-grid">
              <label class="field"><span>TMDB ID</span><div class="field-row"><input class="input mono" id="feed-tmdb-id" value="${esc(feed.tmdb_id || '')}" placeholder="自动匹配或手动填写"><button class="btn btn-sm" type="button" id="feed-tmdb-auto">${icon('refresh')}自动</button><button class="btn btn-sm" type="button" id="feed-tmdb-search">${icon('search')}搜索</button></div></label>
              <label class="field"><span>TMDB 类型</span><select class="select" id="feed-tmdb-type"><option value="tv" ${(feed.tmdb_type || 'tv') === 'tv' ? 'selected' : ''}>TV</option><option value="movie" ${feed.tmdb_type === 'movie' ? 'selected' : ''}>Movie</option></select></label>
              <label class="field span-2"><span>BGM 链接</span><input class="input mono" id="feed-bgm-url" value="${esc(initialBgm)}" placeholder="https://bgm.tv/subject/665390"></label>
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">季数与集数</div>
            <div class="form-grid">
              <label class="field"><span>季数</span><input class="input" id="feed-rename-season" type="number" min="0" max="1000" value="${Number(feed.rename_season || 0)}"></label>
              <label class="field"><span>集数偏移</span><input class="input" id="feed-rename-offset" type="number" min="-1000" max="1000" value="${Number(feed.rename_offset || 0)}"></label>
              <label class="field"><span>总集数</span><input class="input" id="feed-total-episodes" type="number" min="0" max="10000" value="${Number(feed.total_episodes || 0)}"></label>
            </div>
          </div>
        </section>

        <section class="feed-tab-panel" data-feed-panel="backup">
          <div class="feed-subgroups-head">
            <span class="rss-role-badge backup">备用 RSS</span>
            <span class="hint">与主 RSS 共用番剧标题，只区分字幕组名称、地址和规则；留空规则表示不做过滤。</span>
            ${labMultiSubgroup ? `<button class="btn btn-sm" type="button" id="feed-add-subgroup">${icon('plus')}添加备用 RSS</button>` : `<span class="badge">实验室未启用</span>`}
          </div>
          ${labMultiSubgroup ? `<div class="feed-subgroups-table">
            <div class="feed-subgroup-table-head"><span>字幕组</span><span>RSS</span><span>季 / 偏移</span><span>操作</span></div>
            <div id="feed-subgroups"></div>
          </div>` : `<div class="empty" style="padding:26px"><p>请先在「设置 → 实验室」开启多字幕组。</p></div>`}
        </section>

        <section class="feed-tab-panel" data-feed-panel="custom">
          <div class="form-grid">
            <label class="field"><span>保存文件夹来源</span><select class="select" id="feed-final-source">
              <option value="global1" ${feed.final_cid ? '' : (Number(feed.final_mode) === 2 ? '' : 'selected')}>全局保存文件夹 1</option>
              <option value="global2" ${feed.final_cid ? '' : (Number(feed.final_mode) === 2 ? 'selected' : '')}>全局保存文件夹 2（按番剧分文件夹）</option>
              <option value="custom" ${feed.final_cid ? 'selected' : ''}>自定义保存文件夹</option>
            </select><span class="hint">选择使用全局的哪个保存文件夹，或单独指定。</span></label>
            <label class="field" id="feed-final-custom-field" ${feed.final_cid ? '' : 'style="display:none"'}>
              <span>自定义保存文件夹</span>
              <div class="field-row"><input class="input" id="feed-final-cid-name" readonly placeholder="选择自定义保存文件夹"><input type="hidden" id="feed-final-cid" value="${esc(feed.final_cid || '')}"><button class="btn btn-sm" type="button" id="feed-final-cid-pick">${icon('folder')}选择</button></div>
              <span class="hint" id="feed-final-cid-path">${feed.final_cid ? '已单独设置保存文件夹' : '选择自定义保存文件夹'}</span>
            </label>
            <label class="field" id="feed-final-custom-mode-field" ${feed.final_cid ? '' : 'style="display:none"'}>
              <span>保存文件夹模式</span><select class="select" id="feed-final-custom-mode">
                <option value="1" ${Number(feed.final_custom_mode) !== 2 ? 'selected' : ''}>保存文件夹 1：全部放一起</option>
                <option value="2" ${Number(feed.final_custom_mode) === 2 ? 'selected' : ''}>保存文件夹 2 子路径：按子路径分文件夹</option>
              </select><span class="hint">保存文件夹 2 会使用全局子路径模板。</span>
            </label>
            <label class="field"><span>缓存过期（秒）</span><input class="input" id="feed-expiration" type="number" min="0" step="1" value="${feed.expiration || 0}"><span class="hint">0 使用全局缓存策略。</span></label>
            <div class="field span-2"><div class="field-row"><button class="btn btn-sm" type="button" id="feed-rules">${icon('filter')}规则与重命名</button><button class="btn btn-sm" type="button" id="feed-item-preview">${icon('eye')}预览</button><span class="hint">预览会合并主 RSS 和全部备用 RSS，并按集数排序。</span></div></div>
            <div class="field span-2">
              <div class="field-row">
                <button class="btn btn-sm" type="button" id="feed-preview">${icon('eye')}测试订阅</button>
                <span class="hint">抓取 RSS 并应用过滤规则，按来源显示匹配结果。</span>
              </div>
              <div id="feed-preview-box"></div>
            </div>
          </div>
        </section>
      </form>
      ${editing ? `<div class="feed-edit-actions">
        <button class="btn btn-soft" type="button" id="feed-action-run">${icon('play')}运行订阅</button>
        <button class="btn" type="button" id="feed-action-copy">${icon('link')}复制 RSS</button>
        <button class="btn" type="button" id="feed-action-toggle">${feed.enabled === false ? '启用订阅' : '禁用订阅'}</button>
        <span class="spacer"></span>
        <button class="btn btn-danger" type="button" id="feed-action-delete">${icon('trash')}删除订阅</button>
      </div>` : ''}`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="feed-save">${icon('save')}${editing ? '保存修改' : '添加订阅'}</button>`,
    onMount(el, close) {
      const form = $('#feed-form', el);
      if (form) form.addEventListener('submit', event => event.preventDefault());
      if (editing) {
        $('#feed-action-run', el).addEventListener('click', () => { close(); setTimeout(() => runSingleFeed(feed), 80); });
        $('#feed-action-copy', el).addEventListener('click', () => copyText(feed.url));
        $('#feed-action-toggle', el).addEventListener('click', () => { close(); toggleFeed(index); });
        $('#feed-action-delete', el).addEventListener('click', () => { close(); deleteFeed(index); });
      }
      let refreshBackupSources = null;
      const feedTabs = $$('[data-feed-tab]', el);
      const activateFeedTab = name => {
        feedTabs.forEach(tab => tab.classList.toggle('active', tab.dataset.feedTab === name));
        $$('[data-feed-panel]', el).forEach(panel => panel.classList.toggle('active', panel.dataset.feedPanel === name));
        if (name === 'backup' && refreshBackupSources) refreshBackupSources();
      };
      feedTabs.forEach(tab => tab.addEventListener('click', () => activateFeedTab(tab.dataset.feedTab)));
      const subgroupBox = $('#feed-subgroups', el);
      if (subgroupBox) {
        const subgroupRow = (group = {}, index = 0, total = 1) => {
          const enabled = group.enabled !== false;
          return `<div class="feed-subgroup-row" data-subgroup-row data-index="${index}">
            <input class="input" data-sg="name" value="${esc(subgroupLabel(group.name, feed.name))}" placeholder="字幕组名称">
            <input class="input mono" data-sg="url" value="${esc(group.url || '')}" placeholder="备用 RSS 地址">
            <div class="feed-subgroup-nums"><label>季 <input class="input" data-sg="season" type="number" min="0" max="1000" value="${Number(group.season || 0)}"></label><label>偏移 <input class="input" data-sg="offset" type="number" min="-1000" max="1000" value="${Number(group.offset || 0)}"></label></div>
            <div class="feed-subgroup-actions">
              <button class="btn btn-sm btn-ghost" type="button" data-set-main title="设为第一个 RSS（主 RSS）">设为主</button>
              <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-subgroup-up title="上移" ${index === 0 ? 'disabled' : ''}>↑</button>
              <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-subgroup-down title="下移" ${index === total - 1 ? 'disabled' : ''}>↓</button>
              <button class="btn btn-sm btn-ghost ${group.filter || group.exclude ? 'has-rules' : ''}" type="button" data-subgroup-rules title="${group.filter || group.exclude ? '已设置独立规则' : '未设置规则，不过滤'}">规则</button>
              <label class="check"><span class="switch"><input type="checkbox" data-sg="enabled" ${enabled ? 'checked' : ''}><span></span></span>启用</label>
              <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-remove-subgroup title="删除">${icon('trash')}</button>
              <input type="hidden" data-sg="filter" value="${esc(group.filter || '')}">
              <input type="hidden" data-sg="exclude" value="${esc(group.exclude || '')}">
            </div>
          </div>`;
        };
        const currentMain = () => ({
          name: $('#feed-subgroup-name', el).value.trim(),
          url: $('#feed-url', el).value.trim(),
          season: Number($('#feed-rename-season', el).value || 0),
          offset: Number($('#feed-rename-offset', el).value || 0),
          filter: local.filter || '',
          exclude: local.exclude || '',
          enabled: true,
        });
        const renderSubgroups = () => {
          const main = currentMain();
          const mainRow = `<div class="feed-subgroup-row feed-main-source-row" data-main-source-row>
            <div class="feed-subgroup-source"><span class="rss-role-badge main">主 RSS</span><strong>${esc(main.name || '未命名字幕组')}</strong></div>
            <div class="feed-subgroup-source"><code title="${esc(main.url)}">${esc(main.url)}</code></div>
            <div class="feed-subgroup-nums"><span>季 ${main.season || 0}</span><span>偏移 ${main.offset || 0}</span>${main.filter || main.exclude ? '<span class="badge badge-info">独立规则</span>' : ''}</div>
            <div class="feed-subgroup-actions"><button class="btn btn-sm btn-ghost" type="button" data-edit-main>编辑</button></div>
          </div>`;
          subgroupBox.innerHTML = mainRow + ((local.subgroups || []).map((group, index) => subgroupRow(group, index, (local.subgroups || []).length)).join('') || '<div class="hint">还没有备用字幕组</div>');
          subgroupBox.querySelectorAll('[data-edit-main]').forEach(button => button.addEventListener('click', () => activateFeedTab('basic')));
          subgroupBox.querySelectorAll('[data-set-main]').forEach(button => {
            button.addEventListener('click', () => {
              local.subgroups = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(readSubgroupRow);
              const index = Number(button.closest('[data-subgroup-row]').dataset.index);
              const selected = local.subgroups[index];
              if (!selected) return;
              const oldMain = currentMain();
              $('#feed-url', el).value = selected.url;
              $('#feed-subgroup-name', el).value = subgroupLabel(selected.name, feed.name);
              $('#feed-rename-season', el).value = selected.season || 0;
              $('#feed-rename-offset', el).value = selected.offset || 0;
              local.filter = selected.filter || '';
              local.exclude = selected.exclude || '';
              local.subgroups.splice(index, 1, { ...oldMain, name: subgroupLabel(oldMain.name, feed.name) });
              renderSubgroups();
            });
          });
          subgroupBox.querySelectorAll('[data-subgroup-up]').forEach(button => {
            button.addEventListener('click', () => {
              local.subgroups = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(readSubgroupRow);
              const index = Number(button.closest('[data-subgroup-row]').dataset.index);
              if (index > 0) [local.subgroups[index - 1], local.subgroups[index]] = [local.subgroups[index], local.subgroups[index - 1]];
              renderSubgroups();
            });
          });
          subgroupBox.querySelectorAll('[data-subgroup-down]').forEach(button => {
            button.addEventListener('click', () => {
              local.subgroups = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(readSubgroupRow);
              const index = Number(button.closest('[data-subgroup-row]').dataset.index);
              if (index < local.subgroups.length - 1) [local.subgroups[index + 1], local.subgroups[index]] = [local.subgroups[index], local.subgroups[index + 1]];
              renderSubgroups();
            });
          });
          subgroupBox.querySelectorAll('[data-subgroup-rules]').forEach(button => {
            button.addEventListener('click', () => {
              const rows = [...subgroupBox.querySelectorAll('[data-subgroup-row]')];
              local.subgroups = rows.map(readSubgroupRow);
              const index = rows.indexOf(button.closest('[data-subgroup-row]'));
              if (index < 0) return;
              const group = local.subgroups[index];
              const context = {
                ...feed,
                name: group.name || feed.name || '',
                url: group.url || feed.url || '',
                filter: group.filter || '',
                exclude: group.exclude || '',
                rename_season: group.season || feed.rename_season || 0,
                rename_offset: group.offset || feed.rename_offset || 0,
              };
              openRulesModal((filter, exclude) => {
                local.subgroups[index].filter = filter;
                local.subgroups[index].exclude = exclude;
                renderSubgroups();
              }, () => context);
            });
          });
          subgroupBox.querySelectorAll('[data-remove-subgroup]').forEach(button => {
            button.addEventListener('click', () => {
              const row = button.closest('[data-subgroup-row]');
              local.subgroups = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(readSubgroupRow).filter((_, i) => readSubgroupRow([...subgroupBox.querySelectorAll('[data-subgroup-row]')][i]) !== undefined);
              const index = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].indexOf(row);
              if (index >= 0) local.subgroups.splice(index, 1);
              renderSubgroups();
            });
          });
        };
        const readSubgroupRow = row => ({
          name: row.querySelector('[data-sg="name"]').value.trim(),
          url: row.querySelector('[data-sg="url"]').value.trim(),
          season: Number(row.querySelector('[data-sg="season"]').value || 0),
          offset: Number(row.querySelector('[data-sg="offset"]').value || 0),
          enabled: row.querySelector('[data-sg="enabled"]').checked,
          filter: row.querySelector('[data-sg="filter"]').value,
          exclude: row.querySelector('[data-sg="exclude"]').value,
        });
        refreshBackupSources = renderSubgroups;
        renderSubgroups();
        $('#feed-add-subgroup', el)?.addEventListener('click', () => {
          local.subgroups = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(readSubgroupRow);
          const source = sourceInfo(feed.source, feed.url).key;
          openBrowseModal(source, pick => {
            local.subgroups.push({
              name: pick.subgroup || pick.title || '备用字幕组',
              url: pick.url,
              season: Number(feed.rename_season || 0),
              offset: Number(feed.rename_offset || 0),
              enabled: true,
            });
            renderSubgroups();
          }, feed.name || '');
        });
      }
      const finalCidInput = $('#feed-final-cid', el);
      const finalCidName = $('#feed-final-cid-name', el);
      const hydrateFinalCid = async () => {
        const cid = finalCidInput ? finalCidInput.value.trim() : '';
        if (!cid) { if (finalCidName) finalCidName.value = ''; return; }
        if (finalCidName) finalCidName.value = '读取中…';
        try {
          const data = await api(`/api/fs/stat?cid=${encodeURIComponent(cid)}`);
          if (finalCidName) finalCidName.value = data.name || '已选择文件夹';
        } catch { if (finalCidName) finalCidName.value = '已选择文件夹'; }
      };
      hydrateFinalCid();
      const finalSource = $('#feed-final-source', el);
      const finalCustomField = $('#feed-final-custom-field', el);
      const finalCustomModeField = $('#feed-final-custom-mode-field', el);
      const updateFinalSource = () => {
        const custom = finalSource && finalSource.value === 'custom';
        if (finalCustomField) finalCustomField.style.display = custom ? '' : 'none';
        if (finalCustomModeField) finalCustomModeField.style.display = custom ? '' : 'none';
      };
      if (finalSource) finalSource.addEventListener('change', updateFinalSource);
      updateFinalSource();
      const finalPick = $('#feed-final-cid-pick', el);
      if (finalPick) finalPick.addEventListener('click', () => {
        openFolderPicker(({ cid, path }) => {
          finalCidInput.value = cid;
          finalCidName.value = path.split('/').filter(Boolean).pop() || path;
          const hint = $('#feed-final-cid-path', el);
          if (hint) hint.textContent = path;
        }, '选择该订阅的保存文件夹');
      });
      const providerBox = $('#provider-fields', el);
      const helpBox = $('#source-help', el);
      const browseRow = $('#feed-browse-row', el);
      const urlInput = $('#feed-url', el);

      const captureProviders = () => {
        $$('[data-provider]', providerBox).forEach(input => {
          local.provider[input.dataset.provider] = input.value.trim();
        });
      };

      const rebuildUrl = () => {
        if (local.source === 'other') return;
        const built = buildProviderUrl(local.source, local.provider, mikanHost);
        urlInput.value = built;
        urlInput.dispatchEvent(new Event('input', { bubbles: true }));
        local.urlTouched = false;
      };

      const renderSource = () => {
        $$('.source-tab', el).forEach(btn => btn.classList.toggle('active', btn.dataset.source === local.source));
        providerBox.innerHTML = sourceFieldsHtml(local.source, local.provider, mikanHost);
        helpBox.innerHTML = sourceHintHtml(local.source, mikanHost);
        if (local.source === 'other') {
          browseRow.innerHTML = '';
        } else {
          browseRow.innerHTML = `<button class="btn btn-sm" type="button" id="feed-browse">${icon('search')}在 ${SOURCES[local.source].label} 中浏览</button><span class="hint">搜索番剧、展开字幕组并一键填入</span>`;
          $('#feed-browse', el).addEventListener('click', () => {
            close();
            openBrowseModal(local.source, pick => {
              openFeedModal(-1, {
                source: pick.source,
                name: pick.title,
                url: pick.url,
              });
            });
          });
        }
        $$('[data-provider]', providerBox).forEach(input => {
          input.addEventListener('input', () => {
            local.provider[input.dataset.provider] = input.value.trim();
            rebuildUrl();
          });
        });
      };

      $$('.source-tab', el).forEach(btn => {
        btn.addEventListener('click', () => {
          if (btn.dataset.source === local.source) return;
          captureProviders();
          local.source = btn.dataset.source;
          local.provider = parseProviderFields(local.source, urlInput.value);
          local.urlTouched = false;
          if (local.source !== 'other') {
            urlInput.value = '';
            rebuildUrl();
          }
          renderSource();
        });
      });

      const bgmInput = $('#feed-bgm-url', el);
      const syncBgmFromURL = () => {
        if (!bgmInput || !local.bgmAuto) return;
        const detected = bgmURLFromFeedURL(urlInput.value);
        if (detected) bgmInput.value = detected;
      };
      urlInput.addEventListener('input', () => { local.urlTouched = true; syncBgmFromURL(); });
      bgmInput?.addEventListener('input', () => { local.bgmAuto = false; });

      $('#feed-rules', el).addEventListener('click', () => {
        openRulesModal((filter, exclude) => {
          local.filter = filter;
          local.exclude = exclude;
          markDirty();
        }, () => readFeedForm(el, local));
      });
      $('#feed-main-rules').addEventListener('click', () => {
        openRulesModal((filter, exclude) => {
          local.filter = filter;
          local.exclude = exclude;
          markDirty();
        }, () => readFeedForm(el, local));
      });
      $('#feed-item-preview', el).addEventListener('click', () => openItemPreview(readFeedForm(el, local)));
      const fillTotalFromTMDB = async (id, type) => {
        const totalInput = $('#feed-total-episodes', el);
        if (!totalInput || Number(totalInput.value) > 0 || !id) return;
        try {
          const data = await api(`/api/tmdb/detail?id=${encodeURIComponent(id)}&type=${encodeURIComponent(type || 'tv')}`);
          if (Number(data.number_of_episodes) > 0) {
            totalInput.value = String(Number(data.number_of_episodes));
            markDirty();
          }
        } catch { /* ignore */ }
      };
      const autoMatchTMDB = async (silent = false) => {
        const name = $('#feed-name', el).value.trim();
        if (!name || $('#feed-tmdb-id', el).value.trim()) return;
        if (!state.config || !state.config.tmdb || !state.config.tmdb.api_key) {
          if (!silent) $('#feed-tmdb-hint', el).textContent = '请先在设置 → TMDB 填写 TmdbApiKey';
          return;
        }
        const hint = $('#feed-tmdb-hint', el);
        if (!silent) hint.textContent = '正在自动匹配 TMDB…';
        try {
          const type = $('#feed-tmdb-type', el).value || 'tv';
          const data = await api(`/api/tmdb/auto?q=${encodeURIComponent(name)}&type=${encodeURIComponent(type)}`);
          if (data.matched && data.item) {
            $('#feed-tmdb-id', el).value = String(data.item.id);
            $('#feed-tmdb-type', el).value = data.item.first_air_date !== undefined ? 'tv' : 'movie';
            hint.textContent = `自动匹配：${data.item.name || data.item.title} (${data.item.id})`;
            markDirty();
          } else if (!silent) {
            hint.textContent = '未自动匹配，可点击「搜索」手动选择';
          }
        } catch (err) {
          if (!silent) hint.textContent = err.message;
        }
      };
      $('#feed-tmdb-auto', el).addEventListener('click', () => autoMatchTMDB(false));
      let tmdbTimer = null;
      $('#feed-name', el).addEventListener('input', () => {
        if ($('#feed-tmdb-id', el).value.trim()) return;
        clearTimeout(tmdbTimer);
        tmdbTimer = setTimeout(() => autoMatchTMDB(true), 800);
      });
      setTimeout(() => autoMatchTMDB(true), 500);
      $('#feed-tmdb-search', el).addEventListener('click', () => {
        const name = $('#feed-name', el).value.trim();
        openTMDBSearch({
          query: name.split(' · ')[0],
          type: $('#feed-tmdb-type', el).value,
          onSelect: item => {
            $('#feed-tmdb-id', el).value = String(item.id);
            $('#feed-tmdb-type', el).value = item.first_air_date !== undefined ? 'tv' : 'movie';
            $('#feed-tmdb-hint', el).textContent = `已关联：${item.name || item.title} (${item.id})`;
            markDirty();
          },
        });
      });
      renderSource();
      if (!editing && !preselected && local.source !== 'other') rebuildUrl();

      $('#feed-preview', el).addEventListener('click', async event => {
        const button = event.currentTarget;
        const box = $('#feed-preview-box', el);
        const payload = readFeedForm(el, local);
        if (!payload.url) { box.innerHTML = `<p class="form-error">请先填写 RSS 地址</p>`; return; }
        setLoading(button, true, '抓取中');
        box.innerHTML = '';
        try {
          const data = await api('/api/feeds/preview', { method: 'POST', body: payload });
          const sources = Array.isArray(data.sources) && data.sources.length ? data.sources : [{ role: 'main', name: payload.name || payload.url, url: payload.url, total: data.total, matched: data.matched, items: data.items || [] }];
          const sourceHtml = sources.map((source, sourceIndex) => {
            const items = source.items || [];
            const label = source.role === 'main' ? '主 RSS' : '备用 RSS';
            return `<details class="preview-source ${source.role}" open>
              <summary class="preview-source-head">
                <span class="rss-role-badge ${source.role === 'main' ? 'main' : 'backup'}">${label}</span>
                <strong>${esc(source.name || source.url)}</strong>
                ${source.subgroup ? `<span class="source-subgroup">${esc(subgroupLabel(source.subgroup, source.name))}</span>` : ''}
                <code>${esc(source.url || '')}</code>
                <span class="badge ${source.error ? 'badge-danger' : 'badge-primary'}">${source.error ? '抓取失败' : `匹配 ${source.matched} / ${source.total}`}</span>
                <span class="source-chevron">›</span>
              </summary>
              ${source.error ? `<p class="form-error" style="margin:8px 10px 10px">${esc(source.error)}</p>` : (items.length ? `<ul class="preview-list">${items.map((item, i) => `
                <li><span class="idx">${String(i + 1).padStart(2, '0')}</span><span class="t" title="${esc(item.title)}">${esc(item.title)}</span>${item.magnet ? '<span class="badge badge-success">有链接</span>' : '<span class="badge badge-warning">无链接</span>'}</li>`).join('')}</ul>` : `<div class="empty" style="padding:18px"><p>没有条目匹配当前规则。</p></div>`)}
            </details>`;
          }).join('');
          box.innerHTML = `<div class="preview-panel">
            <div class="preview-head"><span class="badge badge-primary">共匹配 ${data.matched} / ${data.total}</span><span>${sources.length > 1 ? `主 RSS + ${sources.length - 1} 个备用 RSS` : '仅主 RSS'}${data.truncated ? ' · 每组最多显示 30 条' : ''}</span></div>
            <div class="preview-sources">${sourceHtml}</div>
          </div>`;
        } catch (err) {
          box.innerHTML = `<p class="form-error">${esc(err.message)}</p>`;
        } finally {
          setLoading(button, false);
        }
      });

      $('#feed-save', el).addEventListener('click', async event => {
        const button = event.currentTarget;
        const payload = readFeedForm(el, local);
        if (!payload.name) { toast('请填写订阅名称', 'error'); return; }
        if (!payload.url) { toast('请填写 RSS 地址', 'error'); return; }
        setLoading(button, true, '保存中');
        try {
          const next = structuredClone(state.config.rss);
          if (editing) next[index] = payload; else next.push(payload);
          const body = structuredClone(state.config);
          body.rss = next;
          const resp = await api('/api/config', { method: 'PUT', body: { config: body } });
          state.config = resp.config;
          state.meta = resp.meta;
          state.dirty = false;
          close();
          toast(editing ? '订阅已更新' : '订阅已添加', 'success');
          if (state.route === 'feeds') renderFeedsPage();
          loadStatus().catch(() => {});
        } catch (err) {
          toast(err.message, 'error', 5200);
        } finally {
          setLoading(button, false);
        }
      });
    },
  });
}

function readFeedForm(root, local) {
  const subgroupBox = root.querySelector('#feed-subgroups');
  return {
    site: '',
    name: ($('#feed-name', root) || {}).value?.trim() || '',
    subgroup_name: ($('#feed-subgroup-name', root) || {}).value?.trim() || '',
    url: ($('#feed-url', root) || {}).value?.trim() || '',
    source: local.source,
    enabled: $('#feed-enabled', root).checked,
    cid: local.cid || '',
    savepath: local.savepath || '',
    filter: local.filter || '',
    exclude: local.exclude || '',
    expiration: Number(($('#feed-expiration', root) || {}).value || 0),
    rename_season: Number(($('#feed-rename-season', root) || {}).value || 0),
    rename_offset: Number(($('#feed-rename-offset', root) || {}).value || 0),
    total_episodes: Number(($('#feed-total-episodes', root) || {}).value || 0),
    completed: !!local.completed,
    final_cid: (($('#feed-final-source', root) || {}).value === 'custom') ? (($('#feed-final-cid', root) || {}).value?.trim() || '') : '',
    final_mode: (($('#feed-final-source', root) || {}).value === 'global2') ? 2 : 1,
    final_custom_mode: Number(($('#feed-final-custom-mode', root) || {}).value || 1),

    subgroups: subgroupBox ? [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(row => ({
      name: row.querySelector('[data-sg="name"]').value.trim(),
      url: row.querySelector('[data-sg="url"]').value.trim(),
      season: Number(row.querySelector('[data-sg="season"]').value || 0),
      offset: Number(row.querySelector('[data-sg="offset"]').value || 0),
      enabled: row.querySelector('[data-sg="enabled"]').checked,
      filter: row.querySelector('[data-sg="filter"]').value || '',
      exclude: row.querySelector('[data-sg="exclude"]').value || '',
    })).filter(group => group.url) : (local.subgroups || []),
    bgm_url: ($('#feed-bgm-url', root) || {}).value?.trim() || '',
    tmdb_id: ($('#feed-tmdb-id', root) || {}).value?.trim() || '',
    tmdb_type: ($('#feed-tmdb-type', root) || {}).value || 'tv',
  };
}


function openManageFeeds() {
  const selected = new Set();
  let query = '';
  let season = 'all';
  let status = 'all';
  const all = (state.config.rss || []).map((feed, index) => ({ feed, index, update: feedUpdateInfo(feed) }));
  const seasonValues = [...new Set(all.map(item => Number(item.feed.rename_season || 0)).filter(Boolean))].sort((a, b) => a - b);
  const statusOf = item => (item.feed.enabled !== false && !item.update.completed) ? 'enabled' : 'disabled';

  const list = () => all.filter(item => {
    const hay = `${item.feed.name || ''} ${item.feed.url || ''} ${item.feed.subgroup_name || ''}`.toLowerCase();
    if (query && !hay.includes(query)) return false;
    if (season !== 'all' && Number(item.feed.rename_season || 0) !== Number(season)) return false;
    if (status !== 'all' && statusOf(item) !== status) return false;
    return true;
  });

  let tableBox;
  let countBox;
  let selectAllBox;
  const draw = () => {
    const items = list();
    const visibleIndices = new Set(items.map(item => item.index));
    [...selected].forEach(index => { if (!visibleIndices.has(index)) selected.delete(index); });
    tableBox.innerHTML = `<div class="table-wrap"><table class="table manage-table">
      <thead><tr><th><input type="checkbox" id="manage-select-all" ${items.length && items.every(item => selected.has(item.index)) ? 'checked' : ''}></th><th>标题</th><th>季</th><th>字幕组</th><th>状态</th><th>进度</th><th>类型</th><th>URL</th></tr></thead>
      <tbody>${items.map(({ feed, index, update }) => {
        const total = Number(feed.total_episodes || update.total || 0);
        const progress = total ? `${Number(update.episode || 0)} / ${total}` : `${Number(update.episode || 0)} / —`;
        const enabled = feed.enabled !== false && !update.completed;
        return `<tr>
          <td data-label="选择"><input type="checkbox" data-manage-select="${index}" ${selected.has(index) ? 'checked' : ''}></td>
          <td class="name-cell" data-label="标题"><div class="title">${esc(feed.name || feed.url)}</div><div class="meta">${feed.tmdb_id ? `TMDB ${esc(feed.tmdb_id)}` : esc(feed.subgroup_name || '')}</div></td>
          <td data-label="季">${Number(feed.rename_season || 0) || '—'}</td>
          <td data-label="字幕组">${esc(feed.subgroup_name || urlSubgroupName(feed.url) || '—')}</td>
          <td data-label="状态"><span class="badge ${enabled ? 'badge-success' : ''}">${enabled ? '已启用' : '未启用'}</span></td>
          <td data-label="进度"><span class="badge ${total && Number(update.episode) >= total ? 'badge-success' : 'badge-warning'}">${progress}</span></td>
          <td data-label="类型"><span class="badge badge-info">${esc(feed.tmdb_type || 'tv')}</span></td>
          <td class="name-cell" data-label="URL"><div class="meta">${esc(feed.url || '')}</div></td>
        </tr>`;
      }).join('')}</tbody></table></div>`;
    tableBox.querySelectorAll('[data-manage-select]').forEach(input => input.addEventListener('change', () => {
      const index = Number(input.dataset.manageSelect);
      if (input.checked) selected.add(index); else selected.delete(index);
      countBox.textContent = `已选 ${selected.size} 项`;
      const selectAll = $('#manage-select-all', tableBox);
      const visible = list();
      if (selectAll) selectAll.checked = visible.length > 0 && visible.every(item => selected.has(item.index));
    }));
    const selectAll = $('#manage-select-all', tableBox);
    if (selectAll) selectAll.addEventListener('change', () => {
      if (selectAll.checked) list().forEach(item => selected.add(item.index));
      else list().forEach(item => selected.delete(item.index));
      draw();
    });
    countBox.textContent = `已选 ${selected.size} 项`;
  };

  const saveSelected = async (mutate, success) => {
    if (!selected.size) { toast('请先选择订阅', 'error'); return; }
    const body = structuredClone(state.config);
    mutate(body.rss, [...selected]);
    try {
      const resp = await api('/api/config', { method: 'PUT', body: { config: body } });
      applyConfigState(resp.config, resp.meta);
      toast(success, 'success');
      close();
      renderFeedsPage();
    } catch (err) {
      toast(err.message, 'error', 5200);
    }
  };

  let close;
  const modal = openModal({
    title: '订阅管理',
    subtitle: '批量启用、停用、更新总集数、导出或删除订阅',
    size: 'browse',
    body: `<div class="manage-toolbar">
      <div class="search-box">${icon('search')}<input class="input" id="manage-search" placeholder="搜索标题、字幕组或 URL"></div>
      <select class="select" id="manage-season"><option value="all">全部季数</option>${seasonValues.map(value => `<option value="${value}">第 ${value} 季</option>`).join('')}</select>
      <select class="select" id="manage-status"><option value="all">全部状态</option><option value="enabled">已启用</option><option value="disabled">未启用</option></select>
      <span class="spacer"></span><span class="badge" id="manage-count">已选 0 项</span>
    </div><div id="manage-table"></div>`,
    footer: `<button class="btn btn-sm" id="manage-export">${icon('download')}导出</button><button class="btn btn-sm" id="manage-total">${icon('refresh')}更新总集数</button><span class="spacer"></span><button class="btn btn-sm" id="manage-enable">启用</button><button class="btn btn-sm" id="manage-disable">禁用</button><button class="btn btn-sm btn-danger" id="manage-delete">删除</button><button class="btn" data-action="modal-close">关闭</button>`,
    onMount(el, closeFn) {
      close = closeFn;
      tableBox = el.querySelector('#manage-table');
      countBox = el.querySelector('#manage-count');
      draw();
      el.querySelector('#manage-search').addEventListener('input', debounce(event => { query = event.target.value.trim().toLowerCase(); draw(); }, 140));
      el.querySelector('#manage-season').addEventListener('change', event => { season = event.target.value; draw(); });
      el.querySelector('#manage-status').addEventListener('change', event => { status = event.target.value; draw(); });
      el.querySelector('#manage-enable').addEventListener('click', () => saveSelected((rss, indices) => {
        indices.forEach(index => { if (rss[index]) { rss[index].enabled = true; rss[index].completed = false; } });
      }, '已启用所选订阅'));
      el.querySelector('#manage-disable').addEventListener('click', () => saveSelected((rss, indices) => {
        indices.forEach(index => { if (rss[index]) rss[index].enabled = false; });
      }, '已停用所选订阅'));
      el.querySelector('#manage-delete').addEventListener('click', async () => {
        if (!selected.size) { toast('请先选择订阅', 'error'); return; }
        const ok = await confirmDialog({ title: '删除订阅', message: `确定删除选中的 ${selected.size} 个订阅吗？<br><br>只删除配置，不删除 115 文件。`, confirmText: '删除' });
        if (!ok) return;
        saveSelected((rss, indices) => { indices.sort((a, b) => b - a).forEach(index => rss.splice(index, 1)); }, '已删除所选订阅');
      });
      el.querySelector('#manage-total').addEventListener('click', async () => {
        if (!selected.size) { toast('请先选择订阅', 'error'); return; }
        const body = structuredClone(state.config);
        let updated = 0;
        for (const index of selected) {
          const feed = body.rss[index];
          if (!feed || !feed.tmdb_id) continue;
          try {
            const data = await api(`/api/tmdb/detail?id=${encodeURIComponent(feed.tmdb_id)}&type=${encodeURIComponent(feed.tmdb_type || 'tv')}`);
            if (Number(data.number_of_episodes) > 0) { feed.total_episodes = Number(data.number_of_episodes); updated++; }
          } catch { /* keep */ }
        }
        if (!updated) { toast('所选订阅没有可更新的 TMDB 总集数', 'error'); return; }
        try {
          const resp = await api('/api/config', { method: 'PUT', body: { config: body } });
          applyConfigState(resp.config, resp.meta);
          toast(`已更新 ${updated} 个订阅的总集数`, 'success');
          close(); renderFeedsPage();
        } catch (err) { toast(err.message, 'error', 5200); }
      });
      el.querySelector('#manage-export').addEventListener('click', () => {
        const items = selected.size ? [...selected].map(index => state.config.rss[index]).filter(Boolean) : state.config.rss;
        const blob = new Blob([JSON.stringify({ rss: items }, null, 2)], { type: 'application/json' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = 'anifeed-subscriptions.json';
        link.click();
        setTimeout(() => URL.revokeObjectURL(link.href), 1000);
        toast('订阅已导出', 'success');
      });
    },
  });
  return modal;
}

async function deleteFeed(index) {
  const feed = state.config.rss[index];
  const ok = await confirmDialog({
    title: '删除订阅',
    message: `确定删除订阅「${esc(feed.name || feed.url)}」吗？<br><br>只会从 <code>config.toml</code> 中移除该订阅，不会删除 115 网盘中的文件。`,
    confirmText: '删除',
  });
  if (!ok) return;
  try {
    const body = structuredClone(state.config);
    body.rss.splice(index, 1);
    const resp = await api('/api/config', { method: 'PUT', body: { config: body } });
    state.config = resp.config;
    state.meta = resp.meta;
    toast('订阅已删除', 'success');
    renderFeedsPage();
    loadStatus().catch(() => {});
  } catch (err) {
    toast(err.message, 'error', 5200);
  }
}

async function runSingleFeed(feed) {
  try {
    const job = await api('/api/feeds/run', { method: 'POST', body: feed });
    showJobModal(job);
  } catch (err) {
    toast(err.message, 'error', 5200);
  }
}

async function runAllFeeds() {
  try {
    const job = await api('/api/feeds/run-all', { method: 'POST', body: {} });
    showJobModal(job);
  } catch (err) {
    toast(err.message, 'error', 5200);
  }
}

/* ================================================================
 * Jobs
 * ================================================================ */
function showJobModal(job) {
  const modal = openModal({
    title: job.type === 'all' ? '正在运行全部订阅' : '正在运行订阅',
    subtitle: '任务在服务端执行，关闭窗口不会中断。',
    size: 'wide',
    body: `<div id="job-status" class="page-loading" style="padding:28px 0"><span class="spinner lg"></span><span>正在抓取 RSS 并提交 115 离线…</span></div><div id="job-results" class="job-results"></div>`,
    footer: `<span class="spacer"></span><button class="btn" data-action="modal-close" id="job-close">关闭</button>`,
    onMount(el, close, handle) {
      const poll = async () => {
        if (handle.closed) return;
        try {
          const current = await api(`/api/jobs/${job.id}`);
          renderJob(current, el);
          if (current.status === 'running') {
            state.jobTimer = setTimeout(poll, 1500);
          } else {
            state.jobTimer = null;
            loadStatus().catch(() => {});
            if (state.route === 'dashboard') loadDashboard();
            if (state.route === 'tasks') loadTasksPage();
            if (state.route === 'history') loadHistoryPage();
          }
        } catch (err) {
          $('#job-status', el).innerHTML = `<p class="form-error">${esc(err.message)}</p>`;
        }
      };
      poll();
    },
  });
  return modal;
}

function renderJob(job, root) {
  const statusBox = $('#job-status', root);
  const resultsBox = $('#job-results', root);
  if (!statusBox) return;
  if (job.status === 'running') {
    statusBox.innerHTML = `<span class="spinner lg"></span><span>正在执行：${esc(job.title)}…</span>`;
  } else if (job.status === 'error') {
    statusBox.innerHTML = `<div class="form-error">${esc(job.error || '任务执行失败')}</div>`;
  } else {
    statusBox.innerHTML = `<div class="form-success" style="width:100%">任务执行完成 · 开始于 ${fmtDate(job.startedAt)}${job.finishedAt ? ` · 结束于 ${fmtDate(job.finishedAt)}` : ''}</div>`;
  }
  if (!job.results || !job.results.length) {
    resultsBox.innerHTML = job.status === 'running' ? '' : `<p class="hint">没有产生结果。</p>`;
    return;
  }
  resultsBox.innerHTML = job.results.map(result => `
    <div class="job-result">
      <div class="info">
        <strong>${esc(result.name || result.url)}</strong>
        <span>匹配 ${result.matched} / ${result.total}${result.skipped ? ` · 同集跳过 ${result.skipped}` : ''} · 新增 ${result.added} · 缓存 ${result.cached}${result.empty ? ` · 空链接 ${result.empty}` : ''}${result.error ? ` · <span style="color:var(--danger)">${esc(result.error)}</span>` : ''}${result.warning ? ` · <span style="color:var(--warning)">${esc(result.warning)}</span>` : ''}</span>
      </div>
      <span class="num">+${result.added}</span>
    </div>`).join('');
}

/* ================================================================
 * Tasks
 * ================================================================ */
async function loadTasksPage() {
  $('#view').innerHTML = `<div class="page-loading"><span class="spinner"></span><span>正在读取 115 离线任务…</span></div>`;
  try {
    await loadStatus();
    if (!state.status || !state.status.has_cookies) {
      $('#view').innerHTML = `<div class="card">${emptyBlock('lock', '尚未登录 115',
        '离线任务需要 115 登录后才能读取。请先在设置中填写 Cookies 或扫码登录。',
        `<a class="btn btn-primary" href="#/settings">${icon('login')}去登录</a>`)}</div>`;
      return;
    }
    const data = await api('/api/tasks?limit=300');
    state.tasks = data.tasks || [];
  } catch (err) {
    $('#view').innerHTML = `<div class="card">${emptyBlock('cloud', '无法读取 115 任务', err.message,
      `<button class="btn btn-primary" id="tasks-retry">${icon('refresh')}重试</button>
       <a class="btn" href="#/settings">${icon('login')}检查登录</a>`)}</div>`;
    const btn = $('#tasks-retry');
    if (btn) btn.addEventListener('click', loadTasksPage);
    return;
  }
  renderTasksPage();
}

function renderTasksPage() {
  const tasks = state.tasks || [];
  const counts = {
    all: tasks.length,
    running: tasks.filter(t => t.state === 'running').length,
    done: tasks.filter(t => t.state === 'done').length,
    failed: tasks.filter(t => t.state === 'failed').length,
    waiting: tasks.filter(t => t.state === 'waiting').length,
  };
  const tabs = [
    ['all', '全部'], ['running', '运行中'], ['done', '已完成'], ['failed', '失败'], ['waiting', '等待中'],
  ];

  $('#view').innerHTML = `
    <div class="toolbar tasks-toolbar">
      <button class="btn" id="tasks-refresh">${icon('refresh')}刷新</button>
      <button class="btn" id="tasks-add-magnet">${icon('plus')}添加链接</button>
      <div class="segmented" id="task-tabs">
        ${tabs.map(([key, label]) => `<button data-task-filter="${key}" class="${state.taskFilter === key ? 'active' : ''}">${label}<span style="opacity:.55;margin-left:5px">${counts[key]}</span></button>`).join('')}
      </div>
      <span class="spacer"></span>
      <button class="btn btn-sm" id="clear-menu-btn">${icon('trash')}清理任务</button>
    </div>
    <div class="card" id="task-table-card"></div>
    <div class="sticky-save" style="margin-top:16px">
      <span class="hint" style="color:var(--muted)">提示：清理只作用于 115 离线任务列表，不会删除网盘里的文件。</span>
      <span class="spacer"></span>
      <button class="btn btn-sm btn-danger" data-clear="1">清除已完成</button>
      <button class="btn btn-sm btn-danger" data-clear="3">清除失败</button>
    </div>`;

  $('#tasks-refresh').addEventListener('click', loadTasksPage);
  $('#tasks-add-magnet').addEventListener('click', openMagnetModal);
  $('#task-tabs').addEventListener('click', event => {
    const btn = event.target.closest('[data-task-filter]');
    if (!btn) return;
    state.taskFilter = btn.dataset.taskFilter;
    $$('#task-tabs button').forEach(el => el.classList.toggle('active', el.dataset.taskFilter === state.taskFilter));
    renderTaskTable();
  });
  $('#clear-menu-btn').addEventListener('click', openClearMenu);
  $$('[data-clear]').forEach(btn => btn.addEventListener('click', () => clearTasks(Number(btn.dataset.clear))));
  renderTaskTable();
}

function visibleTasks() {
  const tasks = state.tasks || [];
  if (state.taskFilter === 'all') return tasks;
  return tasks.filter(t => t.state === state.taskFilter);
}

function renderTaskTable() {
  const container = $('#task-table-card');
  if (!container) return;
  const tasks = visibleTasks();
  if (!tasks.length) {
    container.innerHTML = emptyBlock('cloud', '没有任务', '当前筛选条件下没有 115 离线任务。',
      `<button class="btn btn-primary" id="empty-add-magnet">${icon('plus')}添加链接</button>`);
    const btn = $('#empty-add-magnet');
    if (btn) btn.addEventListener('click', openMagnetModal);
    return;
  }
  container.innerHTML = `<div class="table-wrap"><table class="table task-table">
    <thead><tr><th>任务名称</th><th>大小</th><th>状态</th><th>进度</th><th>Hash</th><th></th></tr></thead>
    <tbody>${tasks.map(task => `<tr>
      <td class="name-cell" data-label="任务"><div class="title" title="${esc(task.name)}">${esc(task.name || '未命名任务')}</div><div class="meta">${esc(task.url || '')}</div></td>
      <td class="numeric" data-label="大小">${fmtSize(task.size)}</td>
      <td data-label="状态">${taskStateBadge(task.state)}</td>
      <td data-label="进度">${progressBar(task)}</td>
      <td class="numeric" data-label="Hash">${esc((task.info_hash || '').slice(0, 16))}${task.info_hash && task.info_hash.length > 16 ? '…' : ''}</td>
      <td class="actions" data-label="操作"><button class="btn btn-sm btn-ghost task-copy-btn" data-copy-hash="${esc(task.info_hash)}" title="复制 Hash">${icon('copy')}<span class="task-copy-label">复制</span></button></td>
    </tr>`).join('')}</tbody>
  </table></div>`;
  $$('[data-copy-hash]', container).forEach(btn => {
    btn.addEventListener('click', () => copyText(btn.dataset.copyHash));
  });
}

function openClearMenu() {
  const items = [
    [1, '清除已完成', '只清除状态为已完成的任务'],
    [3, '清除失败', '只清除失败的任务'],
    [4, '清除运行中', '中断并清除正在运行的任务'],
    [2, '清除全部', '清除所有离线任务（不含文件）'],
    [5, '清除已完成并删除文件', '危险：同时删除网盘中的文件'],
    [6, '清除全部并删除文件', '危险：同时删除所有任务及文件'],
  ];
  const modal = openModal({
    title: '清理 115 任务',
    subtitle: '选择要清理的任务范围',
    size: 'narrow',
    body: `<div class="settings-stack" style="gap:8px">${items.map(([type, label, desc]) => `
      <button class="btn" style="height:auto;padding:11px 13px;justify-content:flex-start;text-align:left;flex-direction:column;align-items:flex-start;gap:2px" data-clear-type="${type}">
        <span style="font-weight:600">${label}</span><span class="hint" style="font-weight:400">${desc}</span>
      </button>`).join('')}</div>`,
  });
  $$('[data-clear-type]', modal.el).forEach(btn => {
    btn.addEventListener('click', () => {
      modal.close();
      clearTasks(Number(btn.dataset.clearType));
    });
  });
}

async function clearTasks(type) {
  const labels = {
    1: '清除已完成任务', 2: '清除全部任务', 3: '清除失败任务',
    4: '清除运行中任务', 5: '清除已完成任务并删除文件', 6: '清除全部任务并删除文件',
  };
  const danger = type >= 5;
  const ok = await confirmDialog({
    title: labels[type] || '清理任务',
    message: danger
      ? '此操作会同时删除 115 网盘中的文件，且不可恢复。确定继续吗？'
      : '确定要清理对应状态的 115 离线任务吗？',
    confirmText: '执行清理',
    danger: true,
  });
  if (!ok) return;
  try {
    await api('/api/tasks/clear', { method: 'POST', body: { type } });
    toast('清理请求已提交', 'success');
    await loadTasksPage();
  } catch (err) {
    toast(err.message, 'error', 5200);
  }
}

function openMagnetModal() {
  const modal = openModal({
    title: '添加离线链接',
    subtitle: '支持 magnet、ed2k、http(s) 与 ftp 链接，每行一个',
    size: 'wide',
    body: `<form class="form-grid" id="magnet-form">
      <label class="field span-2"><span>链接 *</span><textarea class="textarea" id="magnet-text" placeholder="magnet:?xt=urn:btih:...&#10;ed2k://|file|..."></textarea></label>
      <label class="field"><span>115 目录 ID (cid)</span><input class="input" id="magnet-cid" placeholder="可选，留空使用默认目录"></label>
      <label class="field"><span>保存路径</span><input class="input" id="magnet-savepath" placeholder="可选，例如：番剧/某番"></label>
    </form>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="magnet-save">${icon('plus')}提交任务</button>`,
    onMount(el, close) {
      $('#magnet-save', el).addEventListener('click', async event => {
        const btn = event.currentTarget;
        const text = $('#magnet-text', el).value.trim();
        if (!text) { toast('请填写至少一个链接', 'error'); return; }
        setLoading(btn, true, '提交中');
        try {
          const data = await api('/api/magnets', {
            method: 'POST',
            body: {
              text,
              cid: $('#magnet-cid', el).value.trim(),
              savepath: $('#magnet-savepath', el).value.trim(),
            },
          });
          close();
          toast(`已提交 ${data.added} / ${data.total} 个任务`, 'success');
          if (state.route === 'tasks') loadTasksPage();
          if (state.route === 'dashboard') loadDashboard();
          loadStatus().catch(() => {});
        } catch (err) {
          toast(err.message, 'error', 5200);
        } finally {
          setLoading(btn, false);
        }
      });
    },
  });
}

/* ================================================================
 * History
 * ================================================================ */
async function loadHistoryPage() {
  $('#view').innerHTML = `<div class="page-loading"><span class="spinner"></span><span>正在读取历史记录…</span></div>`;
  try {
    await ensureConfig().catch(() => {});
    const params = new URLSearchParams({ limit: String(state.historyLimit) });
    if (state.historyName) params.set('name', state.historyName);
    const data = await api(`/api/history?${params}`);
    state.history = data.items || [];
  } catch (err) {
    renderFatal(err.message);
    return;
  }
  renderHistoryPage();
}

function renderHistoryPage() {
  const feeds = (state.config && state.config.rss) || [];
  const names = Array.from(new Set(feeds.map(f => f.name).filter(Boolean))).sort();
  if (state.historyName && !names.includes(state.historyName)) names.unshift(state.historyName);

  const rows = state.history.length ? state.history.map(item => `<tr>
    <td class="numeric nowrap" data-label="时间">${fmtDate(item.created_at)}</td>
    <td class="name-cell" data-label="标题"><div class="title" title="${esc(item.title)}">${esc(item.title || '—')}</div><div class="meta">${esc(item.config_url || '')}</div></td>
    <td data-label="订阅"><span class="badge badge-info">${esc(item.config_name || '未命名')}</span></td>
    <td class="numeric" data-label="Hash">${esc((item.task_hash || '').slice(0, 18))}${item.task_hash && item.task_hash.length > 18 ? '…' : ''}</td>
    <td class="actions" data-label="操作"><button class="btn btn-sm btn-ghost btn-icon-only" data-copy="${esc(item.task_hash || '')}" title="复制任务 Hash">${icon('copy')}</button></td>
  </tr>`).join('') : '';

  $('#view').innerHTML = `
    <div class="toolbar history-toolbar">
      <button class="btn" id="history-refresh">${icon('refresh')}刷新</button>
      <select class="select" id="history-name" style="width:auto;min-width:180px">
        <option value="">全部订阅</option>
        ${names.map(n => `<option value="${esc(n)}" ${state.historyName === n ? 'selected' : ''}>${esc(n)}</option>`).join('')}
      </select>
      <select class="select" id="history-limit" style="width:auto">
        ${[20, 50, 100, 200].map(n => `<option value="${n}" ${state.historyLimit === n ? 'selected' : ''}>最近 ${n} 条</option>`).join('')}
      </select>
      <span class="spacer"></span>
      <span class="badge">${state.history.length} 条记录</span>
    </div>
    <div class="card">
      ${rows ? `<div class="table-wrap"><table class="table history-table">
        <thead><tr><th>时间</th><th>标题</th><th>订阅</th><th>任务 Hash</th><th></th></tr></thead>
        <tbody>${rows}</tbody></table></div>`
        : emptyBlock('history', '没有历史记录', '运行订阅后，提交成功的条目会记录在这里。', `<a class="btn btn-primary" href="#/feeds">${icon('rss')}去运行订阅</a>`)}
    </div>`;

  $('#history-refresh').addEventListener('click', loadHistoryPage);
  $('#history-name').addEventListener('change', event => { state.historyName = event.target.value; loadHistoryPage(); });
  $('#history-limit').addEventListener('change', event => { state.historyLimit = Number(event.target.value); loadHistoryPage(); });
  $$('[data-copy]').forEach(btn => btn.addEventListener('click', () => copyText(btn.dataset.copy)));
}

/* ================================================================
 * Logs
 * ================================================================ */
async function loadLogsPage() {
  if (state.logTimer) { clearInterval(state.logTimer); state.logTimer = null; }
  const view = $('#view');
  if (!view.querySelector('#log-pre') && !view.querySelector('#log-list')) {
    view.innerHTML = `<div class="page-loading"><span class="spinner"></span><span>正在读取日志…</span></div>`;
  }
  try {
    const data = await api('/api/logs?limit=500');
    state.logs = data.lines || [];
    renderLogsPage();
    if (state.logAuto) {
      state.logTimer = setInterval(refreshLogsOnly, 3000);
    }
  } catch (err) {
    renderFatal(err.message);
  }
}

function parseLogLine(line) {
  const match = String(line).match(/^(\d{4}[/-]\d{2}[/-]\d{2}[ T]\d{2}:\d{2}:\d{2}(?:\.\d+)?)\s+([\s\S]*)$/);
  const timestamp = match ? match[1] : '';
  const message = match ? match[2] : String(line);
  const lower = message.toLowerCase();
  let level = 'info';
  if (/\[error\]|\berror\b|失败|错误/.test(lower)) level = 'error';
  else if (/\[warning\]|\bwarn\b|警告/.test(lower)) level = 'warn';
  else if (/\[notify\]|\bsuccess\b|完成|成功/.test(lower)) level = 'ok';
  const time = timestamp ? timestamp.replace(/^\d{4}[/-]\d{2}[/-]\d{2}[ T]/, '').slice(0, 8) : '';
  return { time, message, level };
}

function logLevelText(level) {
  return ({ error: '错误', warn: '警告', ok: '完成', info: '信息' })[level] || '信息';
}

function visibleLogs() {
  return state.logs.map(parseLogLine).filter(item => {
    if (state.logLevel !== 'all' && item.level !== state.logLevel) return false;
    if (!state.logQuery) return true;
    return item.message.toLowerCase().includes(state.logQuery);
  });
}

function renderLogBody(keepBottom = true) {
  const list = $('#log-list');
  const pre = $('#log-pre');
  const target = list || pre;
  const atBottom = target ? target.scrollTop + target.clientHeight >= target.scrollHeight - 40 : true;
  const items = visibleLogs();
  if (list) {
    list.innerHTML = items.length ? items.map(item => `<div class="log-row ${item.level}">
      <span class="log-time">${esc(item.time || '—')}</span>
      <span class="log-level">${logLevelText(item.level)}</span>
      <span class="log-msg">${esc(item.message)}</span>
    </div>`).join('') : '<div class="log-empty">没有匹配的日志</div>';
  }
  if (pre) {
    const raw = state.logs.filter(line => {
      const item = parseLogLine(line);
      if (state.logLevel !== 'all' && item.level !== state.logLevel) return false;
      if (state.logQuery && !item.message.toLowerCase().includes(state.logQuery)) return false;
      return true;
    });
    pre.innerHTML = raw.length ? renderLogLines(raw) : '<span class="dim">暂无日志输出</span>';
  }
  const updated = list || pre;
  if (updated && keepBottom && atBottom) updated.scrollTop = updated.scrollHeight;
  const count = $('#log-count');
  if (count) count.textContent = `${items.length} 条`;
}

function renderLogsPage() {
  const auto = state.logAuto;
  const mode = state.logMode;
  $('#view').innerHTML = `
    <div class="log-box">
      <div class="log-head"><span class="dots"><i></i><i></i><i></i></span><span>AniFeed · server</span><span class="spacer"></span>
        <div class="segmented log-mode-switch" id="log-mode-switch">
          <button data-log-mode="list" class="${mode === 'list' ? 'active' : ''}">列表</button>
          <button data-log-mode="terminal" class="${mode === 'terminal' ? 'active' : ''}">终端</button>
        </div>
        <label class="check log-head-check"><span class="switch"><input type="checkbox" id="log-auto" ${auto ? 'checked' : ''}><span></span></span>自动刷新</label>
        <button class="btn btn-sm" id="log-refresh">${icon('refresh')}刷新</button>
        <span id="log-count">${state.logs.length} 条</span>
      </div>
      <div class="log-filter-bar">
        <select class="select" id="log-level">
          <option value="all" ${state.logLevel === 'all' ? 'selected' : ''}>全部级别</option>
          <option value="error" ${state.logLevel === 'error' ? 'selected' : ''}>错误</option>
          <option value="warn" ${state.logLevel === 'warn' ? 'selected' : ''}>警告</option>
          <option value="ok" ${state.logLevel === 'ok' ? 'selected' : ''}>完成</option>
          <option value="info" ${state.logLevel === 'info' ? 'selected' : ''}>信息</option>
        </select>
        <div class="search-box">${icon('search')}<input class="input" id="log-search" value="${esc(state.logQuery)}" placeholder="搜索日志内容"></div>
      </div>
      ${mode === 'list' ? '<div class="log-list" id="log-list"></div>' : '<pre class="log-pre" id="log-pre"></pre>'}
    </div>`;
  $('#log-refresh').addEventListener('click', loadLogsPage);
  $('#log-auto').addEventListener('change', event => {
    state.logAuto = event.target.checked;
    if (state.logTimer) { clearInterval(state.logTimer); state.logTimer = null; }
    if (state.logAuto) state.logTimer = setInterval(refreshLogsOnly, 3000);
  });
  $$('[data-log-mode]').forEach(button => button.addEventListener('click', () => {
    state.logMode = button.dataset.logMode;
    localStorage.setItem('anifeed_log_mode', state.logMode);
    renderLogsPage();
  }));
  $('#log-level').addEventListener('change', event => { state.logLevel = event.target.value; renderLogBody(false); });
  $('#log-search').addEventListener('input', event => { state.logQuery = event.target.value.trim().toLowerCase(); renderLogBody(false); });
  renderLogBody(false);
  const target = $('#log-list') || $('#log-pre');
  if (target) target.scrollTop = target.scrollHeight;
}

async function refreshLogsOnly() {
  if (state.route !== 'logs') return;
  try {
    const data = await api('/api/logs?limit=500');
    state.logs = data.lines || [];
    if (!$('#log-list') && !$('#log-pre')) return;
    renderLogBody(true);
  } catch { /* silent */ }
}

function renderLogLines(lines) {
  return lines.map(line => {
    const text = esc(line);
    if (/\[error\]|error|失败/i.test(line)) return `<span class="err">${text}</span>`;
    if (/\[warning\]|warn/i.test(line)) return `<span class="warn">${text}</span>`;
    if (/add \d+|finished|已启动|success/i.test(line)) return `<span class="ok">${text}</span>`;
    return text;
  }).join('\n');
}

async function loadSettingsPage() {
  if (!state.config) renderSettingsSkeleton();
  try {
    const resp = await api('/api/config');
    applyConfigState(resp.config, resp.meta);
    await loadStatus().catch(() => {});
  } catch (err) {
    renderFatal(err.message);
    return;
  }
  renderSettingsPage();
}

const PRIORITY_META = {
  language: { label: '语言', options: [['chs', '简体'], ['bilingual', '简繁'], ['cht', '繁体']] },
  subtitle: { label: '字幕', options: [['softsub', '内封'], ['hardsub', '内嵌'], ['external', '外挂']] },
  container: { label: '容器', options: [['mkv', 'MKV'], ['mp4', 'MP4']] },
  resolution: { label: '分辨率', options: [['r1080', '1080p'], ['r2160', '4K'], ['r720', '720p']] },
};
const PRIORITY_DEFAULTS = {
  language: ['chs', 'bilingual', 'cht'],
  subtitle: ['softsub', 'hardsub', 'external'],
  container: ['mkv', 'mp4'],
  resolution: ['r1080', 'r2160', 'r720'],
};

function initPriorityEditor(cfg) {
  const box = $('#priority-editor');
  if (!box) return;
  const orders = {};
  for (const key of Object.keys(PRIORITY_META)) {
    const configured = ((cfg.offline || {})[`auto_${key}`] || []).filter(id => PRIORITY_META[key].options.some(([optionId]) => optionId === id));
    orders[key] = configured.length ? configured : [...PRIORITY_DEFAULTS[key]];
  }
  state.autoPriority = orders;
  const draw = () => {
    box.innerHTML = Object.entries(PRIORITY_META).map(([key, meta]) => {
      const parts = orders[key].map((id, index) => {
        const label = (meta.options.find(([optionId]) => optionId === id) || [id, id])[1];
        return `<span class="priority-item">${esc(label)}<span class="priority-buttons"><button type="button" data-priority-move="up" data-priority-key="${key}" data-priority-index="${index}" ${index === 0 ? 'disabled' : ''}>↑</button><button type="button" data-priority-move="down" data-priority-key="${key}" data-priority-index="${index}" ${index === orders[key].length - 1 ? 'disabled' : ''}>↓</button></span></span>`;
      }).join('');
      return `<div class="priority-row"><span class="priority-label">${esc(meta.label)}</span><div class="priority-list">${parts}</div></div>`;
    }).join('');
    $$('[data-priority-move]', box).forEach(button => button.addEventListener('click', () => {
      const key = button.dataset.priorityKey;
      const index = Number(button.dataset.priorityIndex);
      const target = button.dataset.priorityMove === 'up' ? index - 1 : index + 1;
      if (target < 0 || target >= orders[key].length) return;
      [orders[key][index], orders[key][target]] = [orders[key][target], orders[key][index]];
      draw();
    }));
  };
  draw();
}

function initSettingsSections() {
  const cards = $$('.settings-stack > section.card');
  if (!cards.length) return;
  const groupByTitle = {
    '115 账号': 'connection',
    '115 登录': 'connection',
    '服务与访问': 'connection',
    '网络代理': 'connection',
    '数据库与缓存': 'system',
    '默认离线设置': 'offline',
    '自动运行': 'offline',
    '重命名': 'rename',
    'TMDB': 'rename',
    '通知': 'notify',
    '站点选项': 'notify',
    '配置文件': 'system',
    '当前状态': 'system',
    '多字幕组': 'lab',
    '喜好字幕组': 'lab',
    '洗板': 'lab',
    '自动备用字幕组': 'lab',
  };
  cards.forEach(card => {
    const title = (card.querySelector('.section-title h3') || {}).textContent || '';
    card.dataset.settingsGroup = groupByTitle[title.trim()] || 'connection';
  });
  const buttons = $$('[data-settings-section]');
  const apply = () => {
    buttons.forEach(button => button.classList.toggle('active', button.dataset.settingsSection === state.settingsSection));
    cards.forEach(card => {
      card.style.display = card.dataset.settingsGroup === state.settingsSection ? '' : 'none';
    });
  };
  buttons.forEach(button => button.addEventListener('click', () => {
    state.settingsSection = button.dataset.settingsSection;
    apply();
  }));
  apply();
}


function renderSettingsSkeleton() {
  const skeletonCard = (title, rows = 3) => `<section class="card">
    <div class="section-title"><span class="ico">${icon('settings')}</span><div><h3>${title}</h3><div class="skeleton-line short"></div></div></div>
    <div class="section-body">${Array.from({ length: rows }).map(() => `<div class="skeleton-line"></div>`).join('')}</div>
  </section>`;
  $('#view').innerHTML = `
    <div class="settings-grid"><div>
      <div class="subtabs"><button disabled>账号与连接</button><button disabled>离线与任务</button><button disabled>重命名与 TMDB</button><button disabled>通知与站点</button><button disabled>实验室</button><button disabled>系统信息</button></div>
      <div class="settings-loading-panel">
        <div class="settings-stack">${skeletonCard('115 账号', 3)}${skeletonCard('默认离线设置', 5)}${skeletonCard('通知与站点', 4)}</div>
        <div class="settings-loading-veil"><span class="spinner"></span><strong>正在读取配置…</strong></div>
      </div>
    </div></div>`;
}

function expandFinalPathPreview(template, sample) {
  const now = new Date();
  const seasonName = ['冬', '春', '夏', '秋'][Math.floor(now.getMonth() / 3)];
  const values = {
    title: sample.title || '冰之城墙',
    name: sample.title || '冰之城墙',
    site: sample.site || 'anibt.net',
    year: String(now.getFullYear()),
    month: String(now.getMonth() + 1).padStart(2, '0'),
    day: String(now.getDate()).padStart(2, '0'),
    season: `${now.getFullYear()}${seasonName}`,
    animeSeason: `Season ${sample.season || 1}`,
    animeSeasonCN: `第${sample.season || 1}季`,
    quarter: seasonName,
    seasonNumber: String(sample.season || 1),
    tmdbid: sample.tmdbID || '283428',
  };
  values.titleYearTmdb = `${values.title} (${values.year}) {tmdb-${values.tmdbid}}`;
  return String(template || '').replace(/\$\{(\w+)\}/g, (match, key) => values[key] !== undefined ? values[key] : match).replace(/^\/+|\/+$/g, '');
}

function openFinalPathHelp() {
  const rows = [
    ['${title}', '订阅标题', '冰之城墙'],
    ['${titleYearTmdb}', '标题 + 年份 + TMDB ID', '冰之城墙 (2026) {tmdb-283428}'],
    ['${tmdbid}', 'TMDB ID', '283428'],
    ['${animeSeason}', '番剧季数（英文）', 'Season 2'],
    ['${animeSeasonCN}', '番剧季数（中文）', '第2季'],
    ['${quarter}', '番剧首播季度', '秋'],
    ['${seasonNumber}', '番剧季数数字', '2'],
    ['${site}', 'RSS 站点域名', 'anibt.net'],
    ['${year}', '番剧首播年份', '2026'],
    ['${season}', '番剧首播年份 + 季度', '2026秋'],
    ['${month}', '番剧首播月份', '10'],
    ['${day}', '番剧首播日期', '23'],
  ];
  openModal({
    title: '保存文件夹 2 子路径变量',
    subtitle: '保存文件夹 2 子路径支持以下变量',
    size: 'narrow',
    body: `<div class="path-help-list">${rows.map(([key, desc, example]) => `<div class="path-help-row"><code>${esc(key)}</code><span>${esc(desc)}</span><em>${esc(example)}</em></div>`).join('')}</div>`,
    footer: `<div class="hint">路径会自动清理 115 不允许的字符；子文件夹会自动创建。</div><span class="spacer"></span><button class="btn" data-action="modal-close">关闭</button>`,
  });
}

function renderSettingsPage() {
  state.settingsTab = 'form';
  const cfg = state.config;
  const configuredEvents = ((cfg.notify || {}).events || []);
  const notifyEvents = configuredEvents.length ? configuredEvents : DEFAULT_NOTIFY_EVENTS;
  const renameHasSubgroup = String((cfg.rename || {}).template || '').includes('${subgroup}');
  if (!state.labWashFolders) {
    state.labWashFolders = (((cfg.lab || {}).wash_folders) || []).map(cid => ({ cid, name: '' }));
  }
  if (!state.labPreferred) {
    const presets = ['桜都字幕组', '北宇治字幕组', 'TSDM字幕组', 'LoliHouse', '拨雪寻春', 'ToonsHub'];
    const preferred = ((cfg.lab || {}).preferred_subgroups || []).map(item => String(item).trim()).filter(Boolean);
    const inactive = ((cfg.lab || {}).inactive_subgroups || []).map(item => String(item).trim()).filter(Boolean);
    state.labPreferred = [...preferred];
    state.labInactive = [...new Set([...inactive, ...presets.filter(item => !preferred.includes(item) && !inactive.includes(item))])];
  }
  const meta = state.meta || {};
  const loginState = state.status ? state.status.login_state : 'unknown';
  const loginStateText = ({ ok: '已登录', checking: '检测中', unknown: '检测中', nocookies: '未配置', error: '登录失效' })[loginState] || loginState;

  const accountUser = state.status && state.status.user;
  const accountName = accountUser && accountUser.name ? accountUser.name : (loginState === 'ok' ? '115 已登录' : '未登录 115');
  const accountAvatar = accountUser && accountUser.name ? esc(accountUser.name.slice(0, 1).toUpperCase()) : icon('cloud');
  const accountUid = accountUser && accountUser.id ? accountUser.id : '—';
  const loginHint = loginState === 'ok'
    ? `<span class="badge badge-success">已登录${state.status.user && state.status.user.name ? ' · ' + esc(state.status.user.name) : ''}</span>`
    : loginState === 'checking' || loginState === 'unknown'
      ? `<span class="badge badge-warning">检测中</span>`
      : loginState === 'nocookies'
        ? `<span class="badge">未配置</span>`
        : `<span class="badge badge-danger">登录失效</span>`;

  const legacyBanner = meta.legacy_imported
    ? `<div class="banner warn">${icon('alert')}<span>已合并旧版 rss.json / .cookies / node-site-config.json。保存配置后会写入 <code>${esc(meta.path)}</code>，旧文件保留不删除。</span></div>` : '';
  const restartHint = `<div class="notice">${icon('info')}<span>监听地址与端口在服务重启后生效；其余配置保存后立即热重载。</span></div>`;

  const formTab = `
    <div class="settings-stack">
      <section class="card">
        <div class="section-title"><span class="ico">${icon('cloud')}</span><div><h3>115 账号</h3><p>当前登录账号与连接状态</p></div><span style="margin-left:auto">${loginHint}</span></div>
        <div class="section-body">
          <div class="settings-account">
            <div class="settings-account-avatar">${accountAvatar}</div>
            <div class="settings-account-info">
              <span>115 CLOUD DRIVE</span>
              <h3>${esc(accountName)}${accountUser && accountUser.is_vip ? '<span class="hero-badge">VIP</span>' : ''}</h3>
              <p>UID ${esc(String(accountUid))}</p>
            </div>
          </div>
          <details class="advanced-settings">
            <summary>账号操作与 Cookies</summary>
            <div class="settings-account-actions">
              <button class="btn btn-soft" id="qr-login">${icon('qr')}扫码登录</button>
              <button class="btn btn-danger" id="logout-btn">${icon('logout')}退出 115 登录</button>
            </div>
            <label class="field">
              <span>115 Cookies</span>
              <input class="input mono" id="set-cookies" type="password" autocomplete="off" value="${esc(cfg.auth.cookies || '')}" placeholder="UID=...; CID=...; SEID=...; KID=...">
              <span class="hint">只显示一行，内容以 * 隐藏；留空并保存会清除登录状态。</span>
            </label>
          </details>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('server')}</span><div><h3>服务与访问</h3><p>Web UI 监听地址、登录保护与代理</p></div></div>
        <div class="section-body">
          <div class="settings-summary-row">
            <span><b>监听</b><code>${esc(cfg.server.host)}:${cfg.server.port}</code></span>
            <span><b>登录保护</b>${cfg.server.auth_enabled ? '<span class="badge badge-success">已开启</span>' : '<span class="badge">未开启</span>'}</span>
            <span><b>访问令牌</b>${cfg.server.access_token ? '<span class="badge badge-info">已设置</span>' : '<span class="badge">未设置</span>'}</span>
          </div>
          <details class="advanced-settings">
            <summary>高级设置</summary>
            <div class="form-grid">
              <label class="field"><span>监听地址</span><input class="input" id="set-host" value="${esc(cfg.server.host)}" placeholder="0.0.0.0"><span class="hint">0.0.0.0 表示监听所有网卡，127.0.0.1 仅本机。</span></label>
              <label class="field"><span>端口</span><input class="input" id="set-port" type="number" min="1" max="65535" value="${cfg.server.port}"></label>
              <label class="field span-2"><span>访问令牌（可选）</span><input class="input mono" id="set-token" value="${esc(cfg.server.access_token || '')}" placeholder="留空表示不启用"><span class="hint">设置后所有 API 请求都需要 <code>X-Access-Token</code>。UI 会自动记住当前令牌。</span></label>
              <label class="field span-2"><span>Web UI 登录保护</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auth-enabled" ${cfg.server.auth_enabled ? 'checked' : ''}><span></span></span><span>启用后打开 Web UI 需要账号和密码</span></span></label>
              <label class="field"><span>登录账号</span><input class="input" id="set-auth-username" value="${esc(cfg.server.auth_username || '')}" placeholder="例如 admin" autocomplete="username"></label>
              <label class="field"><span>登录密码</span><input class="input" id="set-auth-password" type="password" placeholder="${cfg.server.auth_enabled ? '留空表示不修改' : '默认 admin'}" autocomplete="new-password"><span class="hint">默认账号和密码都是 admin；修改后下次打开页面需要重新登录。</span></label>
              <label class="field span-2"><span>网络代理</span><input class="input mono" id="set-proxy" value="${esc(cfg.proxy.http)}" placeholder="http://127.0.0.1:10809"><span class="hint">启用站点代理时使用的 HTTP 代理地址。</span></label>
              <div class="field span-2"><button class="btn btn-sm" type="button" id="auth-logout">${icon('logout')}退出当前 Web UI 登录</button></div>
            </div>
            ${restartHint}
          </details>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('database')}</span><div><h3>数据库与缓存</h3><p>SQLite 历史库与 115 请求节流参数</p></div></div>
        <div class="section-body">
          <div class="form-grid">
            <label class="field span-2"><span>数据库路径</span><input class="input mono" id="set-db" value="${esc(cfg.database.path)}"><span class="hint">相对路径基于 config.toml 所在目录解析。</span></label>
            <label class="field span-2"><span>115 同步目录</span><div class="field-row"><input class="input mono" value="${esc(meta.database_path || '')}" readonly><span class="badge">解析后</span></div></label>
          </div>
          <label class="check"><span class="switch"><input type="checkbox" id="set-disable-cache" ${cfg.p115.disable_cache ? 'checked' : ''}><span></span></span>禁用历史去重（disable_cache），每次运行都会重新提交所有匹配条目</label>
          <div class="form-grid">
            <label class="field"><span>分块延迟（秒）</span><input class="input" id="set-chunk-delay" type="number" min="0" value="${cfg.p115.chunk_delay}"></label>
            <label class="field"><span>分块大小</span><input class="input" id="set-chunk-size" type="number" min="1" value="${cfg.p115.chunk_size}"></label>
            <label class="field"><span>冷却最小（毫秒）</span><input class="input" id="set-cooldown-min" type="number" min="0" value="${cfg.p115.cooldown_min_ms}"></label>
            <label class="field"><span>冷却最大（毫秒）</span><input class="input" id="set-cooldown-max" type="number" min="0" value="${cfg.p115.cooldown_max_ms}"></label>
          </div>
        </div>
      </section>


      <section class="card">
        <div class="section-title"><span class="ico">${icon('download')}</span><div><h3>默认离线设置</h3><p>所有订阅未单独指定时使用的默认值</p></div></div>
        <div class="section-body">
          <div class="form-grid">
            <label class="field"><span>离线文件夹</span>
              <div class="field-row"><input class="input" id="set-default-cid-name" readonly placeholder="使用 115 默认目录"><input type="hidden" id="set-default-cid" value="${esc((cfg.offline || {}).default_cid || '')}"><button class="btn btn-sm" type="button" id="set-default-cid-pick">${icon('folder')}选择</button></div>
              <span class="hint" id="set-default-cid-path">115 离线任务先下载到此文件夹</span>
            </label>
            <label class="field" id="set-final-cid-field" ${(cfg.offline || {}).dual_folder ? '' : 'style="display:none"'}>
              <span>保存文件夹 1</span>
              <div class="field-row"><input class="input" id="set-final-cid-name" readonly placeholder="未设置保存文件夹"><input type="hidden" id="set-final-cid" value="${esc((cfg.offline || {}).final_cid || '')}"><button class="btn btn-sm" type="button" id="set-final-cid-pick">${icon('folder')}选择</button></div>
              <span class="hint" id="set-final-cid-path">模式 1：所有视频移动到这个文件夹</span>
            </label>
            <label class="field" id="set-final-cid2-field" ${(cfg.offline || {}).dual_folder ? '' : 'style="display:none"'}>
              <span>保存文件夹 2</span>
              <div class="field-row"><input class="input" id="set-final-cid2-name" readonly placeholder="未设置保存文件夹 2"><input type="hidden" id="set-final-cid2" value="${esc((cfg.offline || {}).final_cid2 || '')}"><button class="btn btn-sm" type="button" id="set-final-cid2-pick">${icon('folder')}选择</button></div>
              <span class="hint" id="set-final-cid2-path">保存文件夹 2 的根文件夹，下面按模板为每个番剧建子文件夹</span>
            </label>
            <label class="field span-2" id="set-final-path-field" ${(cfg.offline || {}).dual_folder ? '' : 'style="display:none"'}>
              <span class="field-label">保存文件夹 2 子路径<button class="help-btn" type="button" id="final-path-help" title="查看所有变量">?</button></span>
              <div class="field-row"><input class="input mono" id="set-final-path-template" value="${esc(((cfg.offline || {}).final_path_template || '\${month}/${titleYearTmdb}/${animeSeason}'))}" placeholder="${'${month}/${titleYearTmdb}/${animeSeason}'}"><button class="btn btn-sm" type="button" id="final-path-preview">${icon('eye')}预览路径</button></div>
              <div class="path-preview" id="final-path-preview-output"></div>
            </label>
            <label class="field span-2"><span>全局排除规则</span><textarea class="textarea mono" id="set-global-exclude" rows="3" placeholder="720p&#10;/\\.5&#36;/">${esc(((cfg.offline || {}).global_exclude || []).join('\n'))}</textarea><span class="hint">每行一条，标题包含即跳过；支持 <code>/正则/</code>。对所有订阅生效。</span></label>
          </div>
          <div class="form-grid">
            <label class="field"><span>失败重试次数</span><input class="input" id="set-retry" type="number" min="0" max="100" value="${Number((cfg.offline || {}).retry || 0)}"><span class="hint">115 提交失败后重试的次数。</span></label>
            <label class="field"><span>跳过 x.5 集</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-skip-x5" ${(cfg.offline || {}).skip_x5 ? 'checked' : ''}><span></span></span><span>跳过 12.5 这类集数</span></span></label>
            <label class="field span-2"><span>运行后自动清理</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-clear" ${(cfg.offline || {}).auto_clear_completed ? 'checked' : ''}><span></span></span><span>全部/自动运行结束后，清理 115 已完成离线任务（不删除文件）</span></span></label>
            <label class="field span-2"><span>订阅完成后自动停用</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-disable-completed" ${(cfg.offline || {}).auto_disable_completed ? 'checked' : ''}><span></span></span><span>达到总集数后自动标记完成并改为未启用</span></span></label>
            <label class="field span-2"><span>离线后转移到保存文件夹</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-dual-folder" ${(cfg.offline || {}).dual_folder ? 'checked' : ''}><span></span></span><span>离线完成后（如启用重命名则先重命名）移动到默认保存文件夹</span></span></label>
            <label class="field span-2"><span>保存文件夹只保留视频和字幕</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-keep-only-video" ${(cfg.offline || {}).keep_only_video ? 'checked' : ''}><span></span></span><span>离线后转移到保存文件夹下，只把视频和字幕移动到保存文件夹；字体/其他文件留在离线目录，字幕自动追加 .sc/.tc 等语言后缀</span></span></label>
            <label class="field span-2"><span>同一 RSS 内每集只选一个版本</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-one-per-episode" ${(cfg.offline || {}).one_per_episode ? 'checked' : ''}><span></span></span><span>关闭时同一 RSS 内所有版本都会提交；开启后只提交评分最高的一个，不影响主 RSS 与备用 RSS 之间各自保留版本</span></span></label>
          </div>
          <label class="field span-2"><span>使用自动规则</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-rules" ${(cfg.offline || {}).auto_rules ? 'checked' : ''}><span></span></span><span>打开「规则与重命名」时，自动按下方优先级选择最佳规则（可再手动调整）</span></span></label>
          <div class="section-subtitle" style="margin-top:16px">自动应用最佳规则优先顺序</div>
          <div class="hint" style="margin-bottom:8px">从上到下依次尝试；当前 RSS 不存在的选项会自动跳过。</div>
          <div id="priority-editor"></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('edit')}</span><div><h3>重命名</h3><p>对 115 已完成的离线文件/文件夹按模板重命名</p></div><span style="margin-left:auto" class="badge ${(cfg.rename || {}).enabled ? 'badge-success' : ''}">${(cfg.rename || {}).enabled ? '已开启' : '未开启'}</span></div>
        <div class="section-body">
          <label class="field"><span>启用重命名</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-rename-enabled" ${(cfg.rename || {}).enabled ? 'checked' : ''}><span></span></span><span>离线完成后按模板重命名；解析不到集数时不会移动</span></span></label>
          <div class="form-grid">
            <label class="field span-2"><span class="label-help">重命名模板 <button class="help-dot" type="button" id="rename-vars-help" title="查看可用变量">?</button></span>
              <div class="field-row"><input class="input mono" id="set-rename-template" value="${esc((cfg.rename || {}).template || '[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}')}"><button class="btn btn-sm" type="button" id="rename-template-preview">${icon('eye')}预览</button></div>
              <div class="path-preview" id="rename-template-output"></div>
              <span class="hint">变量：<code>${'${title}'}</code> 订阅名、<code>${'${subgroup}'}</code> 字幕组、<code>${'${season}'}</code> / <code>${'${seasonFormat}'}</code>、<code>${'${episode}'}</code> / <code>${'${episodeFormat}'}</code>、<code>${'${resolution}'}</code>、<code>${'${tmdbid}'}</code>。</span>
            </label>
            <label class="field"><span>季数</span><input class="input" id="set-rename-season" type="number" min="0" max="1000" value="${Number((cfg.rename || {}).season || 1)}"></label>
            <label class="field"><span>集数偏移</span><input class="input" id="set-rename-offset" type="number" min="-1000" max="1000" value="${Number((cfg.rename || {}).offset || 0)}"><span class="hint">例如标题从 00 开始、实际是第 01 集时填 1。</span></label>
          </div>
          <div class="field-row"><button class="btn btn-sm" id="postprocess-now">${icon('refresh')}立即处理已完成任务</button><span class="hint">后台每 20 秒自动处理；当前待处理：${state.status ? Number(state.status.postprocess_pending || 0) : 0}</span></div>
          <div class="notice">${icon('info')}<span>重命名针对离线完成后的顶层文件或文件夹。离线后转移到保存文件夹下，启用重命名时只有重命名成功才会移动；未启用重命名则直接移动。解析不到集数时会保留在原目录。启用后「运行后自动清理」会暂停，避免任务被提前清理。</span></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('link')}</span><div><h3>TMDB</h3><p>用于关联番剧和 ${'${tmdbid}'} 重命名字段</p></div></div>
        <div class="section-body">
          <div class="form-grid">
            <label class="field"><span>TmdbApi</span><input class="input mono" id="set-tmdb-api" value="${esc((cfg.tmdb || {}).api || 'https://api.themoviedb.org')}" placeholder="https://api.themoviedb.org"></label>
            <label class="field"><span>TmdbApiKey</span><input class="input mono" id="set-tmdb-api-key" value="${esc((cfg.tmdb || {}).api_key || '')}" placeholder="TMDB API Key"></label>
            <label class="field"><span>TmdbImage</span><input class="input mono" id="set-tmdb-image" value="${esc((cfg.tmdb || {}).image || 'https://image.tmdb.org')}" placeholder="https://image.tmdb.org"></label>
          </div>
          <div class="notice">${icon('info')}<span>在订阅编辑中可搜索并关联 TMDB，重命名模板可使用 <code>${'${tmdbid}'}</code>。TMDB 需要有效的 API Key；图片地址用于搜索结果和 TMDB 海报。</span></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('clock')}</span><div><h3>自动运行</h3><p>定时检查启用中的订阅并提交新条目</p></div><span style="margin-left:auto" class="badge ${(cfg.auto || {}).enabled ? 'badge-success' : ''}">${(cfg.auto || {}).enabled ? '已开启' : '未开启'}</span></div>
        <div class="section-body">
          <label class="field"><span>启用自动运行</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-enabled" ${(cfg.auto || {}).enabled ? 'checked' : ''}><span></span></span><span>按间隔运行所有启用的订阅</span></span></label>
          <div class="form-grid">
            <label class="field"><span>自动订阅间隔（分钟）</span><input class="input" id="set-auto-interval" type="number" min="1" max="10080" value="${Number((cfg.auto || {}).interval_minutes || 30)}"><span class="hint">每隔多少分钟自动检查启用的订阅。</span></label>
            <label class="field"><span>启动时立即运行</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-run-on-start" ${(cfg.auto || {}).run_on_start ? 'checked' : ''}><span></span></span><span>服务启动后先跑一次</span></span></label>
            <label class="field"><span>RSS 超时（秒）</span><input class="input" id="set-rss-timeout" type="number" min="1" max="300" value="${Number((cfg.subscription || {}).rss_timeout_seconds || 20)}"><span class="hint">抓取 RSS 的网络超时时间。</span></label>
            <label class="field"><span>Mikan 站点地址</span><input class="input mono" id="set-mikan-host" value="${esc((cfg.subscription || {}).mikan_host || 'https://mikanani.me')}"><span class="hint">用于生成 Mikan RSS 地址，可改为镜像站。</span></label>
          </div>
          <div class="notice">${icon('info')}<span>自动运行只处理「启用」状态的订阅。115 未登录或已有运行中的作业时会自动跳过，并在日志中记录。</span></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('bell')}</span><div><h3>通知</h3><p>运行结束后推送纯文本摘要，可包含新增条目标题</p></div></div>
        <div class="section-body">
          <div class="form-grid">
            <label class="field"><span>通知方式</span><select class="select" id="set-notify-type">
              <option value="auto" ${((cfg.notify || {}).type || 'auto') === 'auto' ? 'selected' : ''}>自动识别</option>
              <option value="generic" ${((cfg.notify || {}).type) === 'generic' ? 'selected' : ''}>通用 Webhook</option>
              <option value="wecom" ${((cfg.notify || {}).type) === 'wecom' ? 'selected' : ''}>企业微信机器人</option>
            </select><span class="hint">自动识别企业微信机器人地址，否则按通用 JSON Webhook 发送。</span></label>
            <label class="field"><span>Webhook 地址</span><input class="input mono" id="set-webhook" value="${esc((cfg.notify || {}).webhook || '')}" placeholder="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=..."></label>
            <label class="field span-2"><span>推送内容</span><span class="check-row">
              <label class="check-inline"><input type="checkbox" id="notify-content-summary" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('summary') ? 'checked' : ''}>运行摘要</label>
              <label class="check-inline"><input type="checkbox" id="notify-content-details" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('details') ? 'checked' : ''}>订阅明细</label>
              <label class="check-inline"><input type="checkbox" id="notify-content-errors" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('errors') ? 'checked' : ''}>错误信息</label>
              <label class="check-inline"><input type="checkbox" id="notify-content-time" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('time') ? 'checked' : ''}>时间</label>
            </span><span class="hint">通知以纯文本发送；企业微信使用 text 消息，通用 Webhook 使用 text/plain。</span></label>
          </div>
          <label class="field span-2"><span>推送事件</span><span class="check-row">
            ${[['run','运行汇总'],['added','添加至云下载'],['completed','下载完成'],['failed','下载失败'],['idle','摸鱼提醒'],['hourly','整点更新'],['wash','洗板完成']].map(([key,label]) => `<label class="check-inline"><input type="checkbox" class="notify-event" data-event="${key}" ${notifyEvents.includes(key) ? 'checked' : ''}>${label}</label>`).join('')}
          </span><span class="hint">单独控制每种事件是否推送；关闭后该事件不会发送。</span></label>
          <div class="form-grid">
            <label class="field span-2"><span>摸鱼提醒</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="notify-idle-enabled" ${(cfg.notify || {}).idle_enabled ? 'checked' : ''}><span></span></span><span>订阅超过指定天数没有更新时发送提醒</span></span></label>
            <label class="field"><span>未更新天数</span><input class="input" id="notify-idle-days" type="number" min="1" max="3650" value="${Number((cfg.notify || {}).idle_days || 7)}"><span class="hint">每天最多提醒一次。</span></label>
          </div>
          <label class="field span-2"><span>通知模板</span><textarea class="textarea mono" id="set-notify-template" rows="13" spellcheck="false">${esc((cfg.notify || {}).template || DEFAULT_NOTIFY_TEMPLATE)}</textarea><span class="hint">变量：\${action} \${title} \${score} \${themoviedbName} \${bgmUrl} \${season} \${episode} \${subgroup} \${currentEpisodeNumber} \${totalEpisodeNumber} \${year} \${month} \${date} \${text} \${bgmEpisodeTitle}</span></label>
          <div class="form-grid">
            <label class="field span-2"><span>整点推送</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-notify-hourly" ${(cfg.notify || {}).hourly_enabled ? 'checked' : ''}><span></span></span><span>在指定整点发送整体更新情况；留空表示每小时</span></span></label>
            <label class="field span-2"><span>推送整点（0-23）</span><input class="input" id="set-notify-hours" value="${esc(((cfg.notify || {}).hourly_hours || []).join(', '))}" placeholder="例如 0, 6, 12, 18；留空表示每小时"><span class="hint">多个小时用逗号分隔，按服务器本地时间。</span></label>
            <label class="field span-2"><span>无更新时静默</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-notify-skip-no-change" ${(cfg.notify || {}).skip_no_change ? 'checked' : ''}><span></span></span><span>本次没有新增条目且没有错误时，不发送运行汇总通知</span></span></label>
            <label class="field"><span>下载超时</span><input class="input" id="set-notify-timeout" type="number" min="1" max="10080" value="${Number((cfg.notify || {}).download_timeout_minutes || 180)}"><span class="hint">超过该分钟数仍未离线完成时发送失败通知。</span></label>
            <label class="field"><span>成功时通知</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-notify-success" ${(cfg.notify || {}).on_success ? 'checked' : ''}><span></span></span><span>提交完成后推送</span></span></label>
            <label class="field"><span>失败时通知</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-notify-error" ${(cfg.notify || {}).on_error ? 'checked' : ''}><span></span></span><span>出现错误时推送</span></span></label>
          </div>
          <div class="field-row">
            <button class="btn btn-sm" id="notify-test">${icon('bell')}发送测试通知</button>
            <span class="hint">使用已保存的配置发送，修改后请先保存。</span>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('shield')}</span><div><h3>站点选项</h3><p>为特定域名启用代理或追加请求头</p></div><button class="btn btn-sm" id="add-site" style="margin-left:auto">${icon('plus')}添加站点</button></div>
        <div class="section-body" id="sites-editor"></div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('file')}</span><div><h3>配置文件</h3><p>当前实例使用的文件路径</p></div></div>
        <div class="section-body">
          <div class="kv">
            <div class="kv-row"><span>config.toml</span><code>${esc(meta.path || '—')}</code></div>
            <div class="kv-row"><span>Cookies 文件</span><code>${esc(meta.cookies_path || '—')}</code></div>
            <div class="kv-row"><span>数据库</span><code>${esc(meta.database_path || '—')}</code></div>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('info')}</span><div><h3>当前状态</h3><p>服务运行概览</p></div></div>
        <div class="section-body">
          <div class="kv">
            <div class="kv-row"><span>服务版本</span><code>v${esc(state.status ? state.status.version : '—')}</code></div>
            <div class="kv-row"><span>订阅数量</span><code>${state.status ? state.status.feed_count : '—'}</code></div>
            <div class="kv-row"><span>登录状态</span><code>${esc(loginStateText)}</code></div>
            <div class="kv-row"><span>自动运行</span><code>${(cfg.auto || {}).enabled ? `每 ${Number((cfg.auto || {}).interval_minutes || 30)} 分钟` : '未开启'}</code></div>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('beaker')}</span><div><h3>多字幕组</h3><p>实验功能：允许同一番剧保留多个字幕组</p></div></div>
        <div class="section-body">
          <label class="field span-2"><span>启用备用 RSS 字幕组</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-multi-subgroup" ${(cfg.lab || {}).multi_subgroup ? 'checked' : ''}><span></span></span><span>启用后可在订阅中维护主 RSS 和多个备用 RSS 字幕组</span></span></label>
          <label class="field span-2"><span>多 RSS 字幕组共存</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-multi-rss" ${(cfg.lab || {}).multi_rss_coexist ? 'checked' : ''}><span></span></span><span>开启后主 RSS 和备用 RSS 都下载，不做跨 RSS 集数覆盖</span></span></label>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('heart')}</span><div><h3>喜好字幕组</h3><p>用于洗板保留顺序和自动备用字幕组</p></div></div>
        <div class="section-body">
          <div class="field span-2"><span>字幕组优先级</span>
            <div id="lab-subgroup-priority" class="subgroup-priority"></div>
            <div class="field-row subgroup-add-row"><input class="input" id="lab-subgroup-add" placeholder="手动添加字幕组名称"><button class="btn btn-sm" type="button" id="lab-subgroup-add-btn">${icon('plus')}添加</button></div>
            <span class="hint">左侧从上到下优先；用于洗板保留和自动备用字幕组。可把字幕组移到右侧不启用。</span>
            <div class="notice" style="margin-top:8px">${icon('info')}<span>喜好字幕组依赖重命名后的标题包含字幕组名称。请确保重命名模板包含 <code>\${subgroup}</code>，否则洗板和自动备用无法识别字幕组。</span></div>
            ${renameHasSubgroup ? '' : `<div class="notice" style="margin-top:8px">${icon('alert')}<span>当前重命名模板不包含 <code>\${subgroup}</code>，喜好字幕组可能无法生效。</span></div>`}
          </div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('trash')}</span><div><h3>洗板</h3><p>按下载记录清理多个文件夹中的重复剧集</p></div></div>
        <div class="section-body">
          <label class="field span-2"><span>启用洗板</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-wash" ${(cfg.lab || {}).wash_enabled ? 'checked' : ''}><span></span></span><span>优先保留喜好字幕组；都不在喜好列表时保留简体且体积更大的版本</span></span></label>
          <label class="field"><span>洗板模式</span><select class="select" id="set-lab-wash-mode">
            <option value="safe" ${((cfg.lab || {}).wash_mode || 'safe') === 'safe' ? 'selected' : ''}>安全模式：只处理有下载记录的文件</option>
            <option value="full" ${((cfg.lab || {}).wash_mode || '') === 'full' ? 'selected' : ''}>完整模式：扫描文件夹内所有可识别剧集</option>
          </select><span class="hint">安全模式不会碰没有下载记录的文件。</span></label>
          <div class="form-grid">
            <label class="field"><span>定时洗板</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-wash-schedule" ${(cfg.lab || {}).wash_schedule_enabled ? 'checked' : ''}><span></span></span><span>按间隔自动执行洗板</span></span></label>
            <label class="field"><span>间隔（分钟）</span><input class="input" id="set-lab-wash-interval" type="number" min="1" max="100000" value="${Number((cfg.lab || {}).wash_interval_minutes || 1440)}"><span class="hint">默认 1440 分钟，即每天一次。</span></label>
          </div>
          <div class="field span-2"><span>洗板文件夹</span><div class="lab-folder-list" id="lab-wash-folders"></div><div class="field-row" style="margin-top:8px"><button class="btn btn-sm" type="button" id="lab-add-wash-folder">${icon('plus')}添加文件夹</button><button class="btn btn-sm btn-danger" type="button" id="lab-wash-now">${icon('trash')}立即洗板</button><span class="hint">先扫描预览，再确认删除；会发送逐文件删除通知。</span></div></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('refresh')}</span><div><h3>自动备用字幕组</h3><p>配合摸鱼提醒自动补充字幕组</p></div></div>
        <div class="section-body">
          <label class="field span-2"><span>自动添加备用字幕组</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-backup" ${(cfg.lab || {}).auto_backup_subgroup ? 'checked' : ''}><span></span></span><span>摸鱼时自动添加备用字幕组，且必须有缺失的那一集</span></span></label>
          <label class="field span-2"><span>自动应用规则</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-apply" ${(cfg.lab || {}).auto_apply_rules ? 'checked' : ''}><span></span></span><span>自动添加的备用字幕组使用自动应用最佳规则</span></span></label>
          <label class="field span-2"><span>显示表单提示</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-hide-hints" ${!(cfg.lab || {}).hide_hints ? 'checked' : ''}><span></span></span><span>关闭后隐藏字段下方说明文字；? / ! 帮助入口仍然保留</span></span></label>
          <div class="notice">${icon('info')}<span>以上为实验性选项，先保存配置；完整执行逻辑会在后续版本接入。</span></div>
        </div>
      </section>
    </div>`;

  const rawTab = `
    <section class="card">
      <div class="section-title"><span class="ico">${icon('file')}</span><div><h3>config.toml</h3><p>${esc(meta.path || '')}</p></div><button class="btn btn-sm" id="raw-reload" style="margin-left:auto">${icon('refresh')}重新读取</button></div>
      <div class="section-body">
        <div class="notice">${icon('info')}<span>原始编辑会保留注释，保存前会自动校验并备份为 <code>config.toml.bak</code>。</span></div>
        <textarea class="textarea raw-editor" id="raw-editor" spellcheck="false" placeholder="正在加载…"></textarea>
        <div class="field-row"><button class="btn btn-primary" id="raw-save">${icon('save')}保存 TOML</button><span class="hint" id="raw-path">${esc(meta.path || '')}</span></div>
      </div>
    </section>`;

  state.settingsSection = state.settingsSection || 'connection';
  $('#view').innerHTML = `
    ${legacyBanner}
    <div class="settings-grid">
      <div>
        ${state.settingsTab === 'form' ? `<div class="subtabs" id="settings-subtabs">
          <button data-settings-section="connection">账号与连接</button>
          <button data-settings-section="offline">离线与任务</button>
          <button data-settings-section="rename">重命名与 TMDB</button>
          <button data-settings-section="notify">通知与站点</button>
          <button data-settings-section="lab">实验室</button>
          <button data-settings-section="system">系统信息</button>
        </div>` : ''}
        <div id="settings-content">${state.settingsTab === 'form' ? formTab : rawTab}</div>
        <div class="sticky-save" id="settings-savebar">
          <span class="dirty" id="dirty-hint">所有修改在点击保存后写入 config.toml</span>
          <span class="spacer"></span>
          <button class="btn" id="settings-cancel">放弃修改</button>
          <button class="btn btn-primary" id="settings-save-2">${icon('save')}保存配置</button>
        </div>
      </div>
    </div>`;

  initPriorityEditor(cfg);
  initSettingsSections();
  hydrateFolderNames();
  renderLabSubgroupPriority();
  renderLabWashFolders();


  const logoutButton = $('#auth-logout');
  if (logoutButton) logoutButton.addEventListener('click', async () => {
    try { await api('/api/session/logout', { method: 'POST' }); } catch { /* ignore */ }
    location.reload();
  });
  const saveButtons = [$('#settings-save-2')].filter(Boolean);
  saveButtons.forEach(btn => btn.addEventListener('click', saveSettings));
  const reload = () => { state.dirty = false; loadSettingsPage(); };
  $('#settings-cancel').addEventListener('click', reload);

  if (state.settingsTab === 'form') {
    renderSitesEditor();
    $('#add-site').addEventListener('click', () => {
      state.settingsSites.push({ host: '', https_agent: false, headers: [] });
      renderSitesEditor(true);
      markDirty();
    });
    $('#qr-login').addEventListener('click', openQrModal);
    $('#logout-btn').addEventListener('click', logout);
    $('#notify-test').addEventListener('click', testNotification);
    $('#set-default-cid-pick').addEventListener('click', () => {
      openFolderPicker(({ cid, path }) => {
        $('#set-default-cid').value = cid;
        $('#set-default-cid-name').value = path.split('/').filter(Boolean).pop() || path;
        $('#set-default-cid-path').textContent = path;
        markDirty();
      }, '选择默认离线目录');
    });
    const washNow = $('#lab-wash-now');
    if (washNow) washNow.addEventListener('click', openWashPreview);
    const addWashFolder = $('#lab-add-wash-folder');
    if (addWashFolder) addWashFolder.addEventListener('click', () => {
      openFolderPicker(({ cid, path }) => {
        state.labWashFolders = state.labWashFolders || [];
        if (!state.labWashFolders.some(item => item.cid === cid)) {
          state.labWashFolders.push({ cid, name: path.split('/').filter(Boolean).pop() || path });
          renderLabWashFolders();
          markDirty();
        }
      }, '选择洗板文件夹');
    });
    $('#set-final-cid-pick').addEventListener('click', () => {
      openFolderPicker(({ cid, path }) => {
        $('#set-final-cid').value = cid;
        $('#set-final-cid-name').value = path.split('/').filter(Boolean).pop() || path;
        $('#set-final-cid-path').textContent = path;
        markDirty();
      }, '选择保存文件夹 1');
    });
    $('#set-final-cid2-pick').addEventListener('click', () => {
      openFolderPicker(({ cid, path }) => {
        $('#set-final-cid2').value = cid;
        $('#set-final-cid2-name').value = path.split('/').filter(Boolean).pop() || path;
        $('#set-final-cid2-path').textContent = path;
        markDirty();
      }, '选择保存文件夹 2');
    });
    const updateFinalCidField = () => {
      const dual = $('#set-dual-folder').checked;
      $('#set-final-cid-field').style.display = dual ? '' : 'none';
      $('#set-final-cid2-field').style.display = dual ? '' : 'none';
      $('#set-final-path-field').style.display = dual ? '' : 'none';
    };
    $('#set-dual-folder').addEventListener('change', () => {
      updateFinalCidField();
      markDirty();
    });
    updateFinalCidField();
    const finalPathPreview = $('#final-path-preview');
    if (finalPathPreview) finalPathPreview.addEventListener('click', () => {
      const output = $('#final-path-preview-output');
      const feed = (state.config.rss || [])[0] || {};
      const template = $('#set-final-path-template').value.trim();
      output.textContent = `示例：保存文件夹2/${expandFinalPathPreview(template, { title: feed.name, tmdbID: feed.tmdb_id, season: feed.rename_season })}/`;
    });
    const finalPathHelp = $('#final-path-help');
    if (finalPathHelp) finalPathHelp.addEventListener('click', openFinalPathHelp);
    const renamePreview = $('#rename-template-preview');
    if (renamePreview) renamePreview.addEventListener('click', () => {
      const feed = (state.config.rss || [])[0] || {};
      const template = $('#set-rename-template').value.trim() || '[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}';
      const season = Number($('#set-rename-season').value || 1);
      const episode = 3 + Number($('#set-rename-offset').value || 0);
      const values = {
        title: feed.name || '青之箱', name: feed.name || '青之箱',
        subgroup: feed.subgroup_name || urlSubgroupName(feed.url) || 'LoliHouse',
        season: String(season), seasonFormat: String(season).padStart(2, '0'),
        episode: String(episode), episodeFormat: String(episode).padStart(2, '0'),
        resolution: '1080p', tmdbid: feed.tmdb_id || '283428',
        itemTitle: `[${feed.subgroup_name || urlSubgroupName(feed.url) || 'LoliHouse'}] ${feed.name || '青之箱'} - 03 [WebRip 1080p]`,
      };
      const output = template.replace(/\$\{(\w+)\}/g, (match, key) => values[key] !== undefined ? values[key] : match);
      $('#rename-template-output').textContent = `示例：${output}`;
    });
    $('#rename-vars-help').addEventListener('click', openRenameVarsHelp);
    $('#postprocess-now').addEventListener('click', async () => {
      try {
        await api('/api/postprocess/run', { method: 'POST', body: {} });
        toast('已开始处理已完成任务', 'success');
        setTimeout(loadStatus, 2500);
      } catch (err) {
        toast(err.message, 'error', 5200);
      }
    });
    $('#settings-content').addEventListener('input', markDirty);
    $('#settings-content').addEventListener('change', markDirty);
  } else {
    loadRawConfig();
    $('#raw-reload').addEventListener('click', loadRawConfig);
    $('#raw-save').addEventListener('click', saveRawConfig);
    $('#raw-editor').addEventListener('input', markDirty);
  }
}

function markDirty() {
  state.dirty = true;
  const hint = $('#dirty-hint');
  if (hint) {
    hint.textContent = '有未保存的修改';
    hint.style.color = 'var(--warning)';
  }
}

function renderSitesEditor(focusLast = false) {
  const container = $('#sites-editor');
  if (!container) return;
  if (!state.settingsSites) state.settingsSites = [];
  if (!state.settingsSites.length) {
    container.innerHTML = `<div class="empty" style="padding:26px"><p>还没有站点选项。点击右上角「添加站点」为特定域名启用代理或追加请求头。</p></div>`;
    return;
  }
  container.innerHTML = state.settingsSites.map((site, index) => `
    <div class="site-editor" data-site-index="${index}">
      <div class="site-editor-head">
        <input class="input" data-site-host value="${esc(site.host)}" placeholder="域名，例如 share.dmhy.org" style="flex:1">
        <label class="check" style="white-space:nowrap"><span class="switch"><input type="checkbox" data-site-agent ${site.https_agent ? 'checked' : ''}><span></span></span>使用代理</label>
        <button class="btn btn-sm btn-danger" data-site-remove title="移除站点">${icon('trash')}</button>
      </div>
      <div class="field" style="gap:7px">
        <span class="hint">请求头</span>
        <div data-header-list>${site.headers.map((h, hi) => headerRow(h, hi)).join('')}</div>
        <button class="btn btn-sm btn-ghost" data-header-add style="align-self:flex-start">${icon('plus')}添加请求头</button>
      </div>
    </div>`).join('');

  $$('[data-site-remove]', container).forEach((btn, i) => btn.addEventListener('click', () => {
    state.settingsSites.splice(i, 1);
    renderSitesEditor();
    markDirty();
  }));
  $$('[data-header-add]', container).forEach((btn, i) => btn.addEventListener('click', () => {
    state.settingsSites[i].headers.push({ key: '', value: '' });
    cfg.lab = {
    multi_subgroup: $('#set-lab-multi-subgroup').checked,
    multi_rss_coexist: $('#set-lab-multi-rss').checked,
    wash_enabled: $('#set-lab-wash').checked,
    wash_mode: $('#set-lab-wash-mode').value,
    wash_schedule_enabled: $('#set-lab-wash-schedule').checked,
    wash_interval_minutes: Number($('#set-lab-wash-interval').value || 1440),
    wash_folders: (state.labWashFolders || []).map(item => item.cid).filter(Boolean),
    preferred_subgroups: [...(state.labPreferred || [])],
    inactive_subgroups: [...(state.labInactive || [])],
    auto_backup_subgroup: $('#set-lab-auto-backup').checked,
    auto_apply_rules: $('#set-lab-auto-apply').checked,
    hide_hints: !$('#set-lab-hide-hints').checked,
  };
  syncSitesFromDom();
    renderSitesEditor();
    markDirty();
  }));
  $$('[data-header-remove]', container).forEach(btn => btn.addEventListener('click', () => {
    const siteEl = btn.closest('[data-site-index]');
    const headerEl = btn.closest('[data-header-row]');
    cfg.lab = {
    multi_subgroup: $('#set-lab-multi-subgroup').checked,
    multi_rss_coexist: $('#set-lab-multi-rss').checked,
    wash_enabled: $('#set-lab-wash').checked,
    wash_mode: $('#set-lab-wash-mode').value,
    wash_schedule_enabled: $('#set-lab-wash-schedule').checked,
    wash_interval_minutes: Number($('#set-lab-wash-interval').value || 1440),
    wash_folders: (state.labWashFolders || []).map(item => item.cid).filter(Boolean),
    preferred_subgroups: [...(state.labPreferred || [])],
    inactive_subgroups: [...(state.labInactive || [])],
    auto_backup_subgroup: $('#set-lab-auto-backup').checked,
    auto_apply_rules: $('#set-lab-auto-apply').checked,
    hide_hints: !$('#set-lab-hide-hints').checked,
  };
  syncSitesFromDom();
    const index = Number(siteEl.dataset.siteIndex);
    const hi = Number(headerEl.dataset.headerIndex);
    state.settingsSites[index].headers.splice(hi, 1);
    renderSitesEditor();
    markDirty();
  }));
  if (focusLast) {
    const inputs = $$('[data-site-host]', container);
    if (inputs.length) inputs[inputs.length - 1].focus();
  }
}

function headerRow(header, index) {
  return `<div class="header-row" data-header-row data-header-index="${index}">
    <input class="input" data-header-key value="${esc(header.key)}" placeholder="Header 名称">
    <input class="input" data-header-value value="${esc(header.value)}" placeholder="值">
    <button class="btn btn-sm btn-ghost btn-icon-only" data-header-remove title="删除">${icon('close')}</button>
  </div>`;
}

function syncSitesFromDom() {
  const container = $('#sites-editor');
  if (!container) return;
  state.settingsSites = $$('.site-editor', container).map(el => ({
    host: el.querySelector('[data-site-host]').value.trim(),
    https_agent: el.querySelector('[data-site-agent]').checked,
    headers: $$('[data-header-row]', el).map(row => ({
      key: row.querySelector('[data-header-key]').value.trim(),
      value: row.querySelector('[data-header-value]').value,
    })).filter(h => h.key),
  }));
}

async function hydrateFolderNames() {
  const entries = [
    ['set-default-cid', 'set-default-cid-name', 'set-default-cid-path', '使用 115 默认目录'],
    ['set-final-cid', 'set-final-cid-name', 'set-final-cid-path', '未设置保存文件夹'],
    ['set-final-cid2', 'set-final-cid2-name', 'set-final-cid2-path', '未设置保存文件夹 2'],
  ];
  for (const [id, nameId, hintId, empty] of entries) {
    const hidden = $('#' + id);
    const name = $('#' + nameId);
    if (!hidden || !name) continue;
    const cid = hidden.value.trim();
    if (!cid) { name.value = ''; name.placeholder = empty; continue; }
    name.value = '读取中…';
    try {
      const data = await api(`/api/fs/stat?cid=${encodeURIComponent(cid)}`);
      name.value = data.name || '已选择文件夹';
      const hint = $('#' + hintId);
      if (hint) hint.textContent = data.name || '已选择文件夹';
    } catch {
      name.value = '已选择文件夹';
    }
  }
}


function renderLabSubgroupPriority() {
  const box = $('#lab-subgroup-priority');
  if (!box) return;
  const draw = () => {
    const item = (name, side, index, total) => `<div class="priority-tag">
      <span>${esc(name)}</span>
      ${side === 'enabled' ? `<button type="button" data-move="up" data-index="${index}" ${index === 0 ? 'disabled' : ''} title="上移">↑</button><button type="button" data-move="down" data-index="${index}" ${index === total - 1 ? 'disabled' : ''} title="下移">↓</button><button type="button" data-move="inactive" data-index="${index}" title="移到不启用">→</button>` : `<button type="button" data-move="enable" data-index="${index}" title="启用">←</button>`}
      <button type="button" data-move="remove" data-index="${index}" data-side="${side}" title="删除">×</button>
    </div>`;
    box.innerHTML = `<div class="subgroup-priority-columns">
      <div class="subgroup-priority-column"><h5>启用 <span class="badge">${state.labPreferred.length}</span></h5><div class="priority-tags">${state.labPreferred.map((name, i) => item(name, 'enabled', i, state.labPreferred.length)).join('') || '<span class="hint">没有启用的字幕组</span>'}</div></div>
      <div class="subgroup-priority-column"><h5>不启用 <span class="badge">${state.labInactive.length}</span></h5><div class="priority-tags">${state.labInactive.map((name, i) => item(name, 'inactive', i, state.labInactive.length)).join('') || '<span class="hint">没有不启用的字幕组</span>'}</div></div>
    </div>`;
    box.querySelectorAll('[data-move]').forEach(button => button.addEventListener('click', () => {
      const side = button.closest('.subgroup-priority-column').querySelector('h5').textContent.startsWith('启用') ? 'enabled' : 'inactive';
      const index = Number(button.dataset.index);
      const move = button.dataset.move;
      if (side === 'enabled') {
        if (move === 'up' && index > 0) [state.labPreferred[index - 1], state.labPreferred[index]] = [state.labPreferred[index], state.labPreferred[index - 1]];
        if (move === 'down' && index < state.labPreferred.length - 1) [state.labPreferred[index + 1], state.labPreferred[index]] = [state.labPreferred[index], state.labPreferred[index + 1]];
        if (move === 'inactive') {
          const [name] = state.labPreferred.splice(index, 1);
          if (!state.labInactive.includes(name)) state.labInactive.push(name);
        }
        if (move === 'remove') state.labPreferred.splice(index, 1);
      } else {
        if (move === 'enable') {
          const [name] = state.labInactive.splice(index, 1);
          if (!state.labPreferred.includes(name)) state.labPreferred.push(name);
        }
        if (move === 'remove') state.labInactive.splice(index, 1);
      }
      draw();
      markDirty();
    }));
  };
  draw();
  const add = () => {
    const input = $('#lab-subgroup-add');
    const name = (input.value || '').trim();
    if (!name) return;
    if (!state.labPreferred.includes(name)) state.labPreferred.push(name);
    state.labInactive = state.labInactive.filter(item => item !== name);
    input.value = '';
    draw();
    markDirty();
  };
  $('#lab-subgroup-add-btn')?.addEventListener('click', add);
  $('#lab-subgroup-add')?.addEventListener('keydown', event => { if (event.key === 'Enter') { event.preventDefault(); add(); } });
}

function renderLabWashFolders() {
  const box = $('#lab-wash-folders');
  if (!box) return;
  if (!state.labWashFolders) state.labWashFolders = [];
  box.innerHTML = state.labWashFolders.map((item, index) => `<span class="lab-folder-chip"><span>${esc(item.name || '已选择文件夹')}</span><button type="button" data-lab-folder-remove="${index}" title="移除">${icon('close')}</button></span>`).join('') || '<span class="hint">还没有添加洗板文件夹</span>';
  box.querySelectorAll('[data-lab-folder-remove]').forEach(button => button.addEventListener('click', () => {
    state.labWashFolders.splice(Number(button.dataset.labFolderRemove), 1);
    renderLabWashFolders(); markDirty();
  }));
  for (const item of state.labWashFolders) {
    if (!item.cid || item.name) continue;
    api(`/api/fs/stat?cid=${encodeURIComponent(item.cid)}`).then(data => {
      item.name = data.name || '';
      renderLabWashFolders();
    }).catch(() => { item.name = '文件夹名称读取失败'; renderLabWashFolders(); });
  }
}

async function openWashPreview() {
  if (state.dirty) {
    toast('请先保存实验室配置，再执行洗板', 'error');
    return;
  }
  try {
    const mode = $('#set-lab-wash-mode').value || 'safe';
    const data = await api('/api/lab/wash', { method: 'POST', body: { dry_run: true, mode } });
    const groups = data.groups || [];
    const rows = groups.map(group => `<div class="wash-group">
      <div class="wash-group-title">${esc(group.series)} · 第 ${Number(group.season) || 1} 季 · 第 ${esc(String(group.episode))} 集</div>
      <div class="wash-keep">保留：<code>${esc(group.keep.path)}</code></div>
      ${(group.delete || []).map(file => `<div class="wash-delete"><span>删除：<code>${esc(file.path)}</code></span><span class="hint">${esc(file.reason || '')}</span></div>`).join('')}
    </div>`).join('');
    openModal({
      title: '洗板预览',
      subtitle: `${data.mode === 'full' ? '完整模式' : '安全模式'} · 扫描 ${data.scanned} 个文件 · 发现 ${groups.length} 组重复`,
      size: 'rules',
      body: `<div class="wash-preview">${groups.length ? rows : '<div class="empty" style="padding:26px"><span class="hint">没有发现需要删除的重复剧集。</span></div>'}</div>`,
      footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-danger" id="wash-confirm" ${data.delete_count ? '' : 'disabled'}>${icon('trash')}确认删除 ${data.delete_count} 个文件</button>`,
      onMount(el, close) {
        el.querySelector('#wash-confirm')?.addEventListener('click', async event => {
          const button = event.currentTarget;
          setLoading(button, true, '删除中');
          try {
            const result = await api('/api/lab/wash', { method: 'POST', body: { dry_run: false, mode } });
            close();
            toast(`洗板完成：删除 ${result.deleted || 0} 个文件，失败 ${result.failed || 0} 个`, result.failed ? 'error' : 'success', 5200);
          } catch (err) {
            toast(err.message, 'error', 5200);
            setLoading(button, false);
          }
        });
      },
    });
  } catch (err) {
    toast(err.message, 'error', 5200);
  }
}

function collectSettings() {
  const cfg = structuredClone(state.config);
  cfg.auth.cookies = $('#set-cookies').value.trim();
  if (!cfg.auth.cookies_file) cfg.auth.cookies_file = '.cookies';
  cfg.server.host = $('#set-host').value.trim();
  cfg.server.port = Number($('#set-port').value || 0);
  cfg.server.access_token = $('#set-token').value.trim();
  cfg.server.auth_enabled = $('#set-auth-enabled').checked;
  cfg.server.auth_username = $('#set-auth-username').value.trim();
  cfg.server.auth_password = $('#set-auth-password').value;
  cfg.database.path = $('#set-db').value.trim();
  cfg.p115.disable_cache = $('#set-disable-cache').checked;
  cfg.p115.chunk_delay = Number($('#set-chunk-delay').value || 0);
  cfg.p115.chunk_size = Number($('#set-chunk-size').value || 0);
  cfg.p115.cooldown_min_ms = Number($('#set-cooldown-min').value || 0);
  cfg.p115.cooldown_max_ms = Number($('#set-cooldown-max').value || 0);
  cfg.proxy.http = $('#set-proxy').value.trim();
  cfg.subscription = {
    mikan_host: $('#set-mikan-host').value.trim(),
    rss_timeout_seconds: Number($('#set-rss-timeout').value || 20),
  };
  cfg.offline = {
    default_cid: $('#set-default-cid').value.trim(),
    default_savepath: '',
    global_exclude: $('#set-global-exclude').value.split(/\r?\n/).map(line => line.trim()).filter(Boolean),
    skip_x5: $('#set-skip-x5').checked,
    retry: Number($('#set-retry').value || 0),
    auto_clear_completed: $('#set-auto-clear').checked,
    auto_disable_completed: $('#set-auto-disable-completed').checked,
    one_per_episode: $('#set-one-per-episode').checked,
    keep_only_video: $('#set-keep-only-video').checked,
    auto_language: state.autoPriority ? state.autoPriority.language : PRIORITY_DEFAULTS.language,
    auto_subtitle: state.autoPriority ? state.autoPriority.subtitle : PRIORITY_DEFAULTS.subtitle,
    auto_container: state.autoPriority ? state.autoPriority.container : PRIORITY_DEFAULTS.container,
    auto_resolution: state.autoPriority ? state.autoPriority.resolution : PRIORITY_DEFAULTS.resolution,
    auto_rules: $('#set-auto-rules').checked,
    dual_folder: $('#set-dual-folder').checked,
    final_cid: $('#set-final-cid').value.trim(),
    final_mode: 1,
    final_cid2: $('#set-final-cid2').value.trim(),
    final_path_template: $('#set-final-path-template').value.trim() || '${month}/${titleYearTmdb}/${animeSeason}',
  };
  cfg.rename = {
    enabled: $('#set-rename-enabled').checked,
    template: $('#set-rename-template').value.trim(),
    season: Number($('#set-rename-season').value || 1),
    offset: Number($('#set-rename-offset').value || 0),
  };
  cfg.tmdb = {
    api: $('#set-tmdb-api').value.trim(),
    api_key: $('#set-tmdb-api-key').value.trim(),
    image: $('#set-tmdb-image').value.trim(),
  };
  cfg.auto = {
    enabled: $('#set-auto-enabled').checked,
    interval_minutes: Number($('#set-auto-interval').value || 30),
    run_on_start: $('#set-auto-run-on-start').checked,
  };
  cfg.notify = {
    type: $('#set-notify-type').value,
    webhook: $('#set-webhook').value.trim(),
    content: ['summary', 'details', 'errors', 'time'].filter(key => $(`#notify-content-${key}`).checked),
    on_success: $('#set-notify-success').checked,
    on_error: $('#set-notify-error').checked,
    idle_enabled: $('#notify-idle-enabled').checked,
    idle_days: Number($('#notify-idle-days').value || 7),
    events: [...document.querySelectorAll('.notify-event')].filter(el => el.checked).map(el => el.dataset.event),
    template: $('#set-notify-template').value,
    download_timeout_minutes: Number($('#set-notify-timeout').value || 180),
    skip_no_change: $('#set-notify-skip-no-change').checked,
    hourly_enabled: $('#set-notify-hourly').checked,
    hourly_hours: $('#set-notify-hours').value.split(/[,\s，]+/).map(Number).filter(n => Number.isInteger(n) && n >= 0 && n <= 23),
  };
  cfg.lab = {
    multi_subgroup: $('#set-lab-multi-subgroup').checked,
    multi_rss_coexist: $('#set-lab-multi-rss').checked,
    wash_enabled: $('#set-lab-wash').checked,
    wash_mode: $('#set-lab-wash-mode').value,
    wash_schedule_enabled: $('#set-lab-wash-schedule').checked,
    wash_interval_minutes: Number($('#set-lab-wash-interval').value || 1440),
    wash_folders: (state.labWashFolders || []).map(item => item.cid).filter(Boolean),
    preferred_subgroups: [...(state.labPreferred || [])],
    inactive_subgroups: [...(state.labInactive || [])],
    auto_backup_subgroup: $('#set-lab-auto-backup').checked,
    auto_apply_rules: $('#set-lab-auto-apply').checked,
    hide_hints: !$('#set-lab-hide-hints').checked,
  };
  syncSitesFromDom();
  cfg.sites = {};
  for (const site of state.settingsSites || []) {
    if (!site.host) continue;
    const headers = {};
    for (const h of site.headers) if (h.key) headers[h.key] = h.value;
    cfg.sites[site.host] = { https_agent: !!site.https_agent, headers };
  }
  return cfg;
}

async function testNotification(event) {
  const button = event.currentTarget;
  setLoading(button, true, '发送中');
  try {
    await api('/api/notify/test', { method: 'POST', body: {} });
    toast('测试通知已发送', 'success');
  } catch (err) {
    toast(err.message, 'error', 5200);
  } finally {
    setLoading(button, false);
  }
}

async function saveSettings(event) {
  const btn = event ? event.currentTarget : null;
  const cfg = collectSettings();
  if (cfg.offline.dual_folder && !cfg.offline.final_cid) {
    toast('启用离线后转移到保存文件夹时必须选择默认保存文件夹', 'error', 5200);
    return;
  }
  if (cfg.rename.enabled && !cfg.rename.template) {
    toast('启用重命名时必须填写重命名模板', 'error', 5200);
    return;
  }
  if (!cfg.notify.content.length) {
    toast('至少选择一项通知推送内容', 'error', 5200);
    return;
  }
  setLoading(btn, true, '保存中');
  try {
    const resp = await api('/api/config', { method: 'PUT', body: { config: cfg } });
    applyConfigState(resp.config, resp.meta);
    state.dirty = false;
    const newToken = (resp.config.server.access_token || '').trim();
    if (newToken !== state.token) {
      state.token = newToken;
      if (newToken) localStorage.setItem('anifeed_token', newToken);
      else localStorage.removeItem('anifeed_token');
    }
    toast('配置已保存并生效', 'success');
    await loadStatus().catch(() => {});
    renderSettingsPage();
  } catch (err) {
    toast(err.message, 'error', 6000);
  } finally {
    setLoading(btn, false);
  }
}

async function logout() {
  const ok = await confirmDialog({
    title: '退出 115 登录',
    message: '将清空本地 Cookies 文件并注销当前客户端。115 云端数据不受影响。',
    confirmText: '退出登录',
  });
  if (!ok) return;
  try {
    await api('/api/auth/logout', { method: 'POST', body: {} });
    toast('已退出登录', 'success');
    state.config = null;
    state.meta = null;
    await loadStatus().catch(() => {});
    loadSettingsPage();
  } catch (err) {
    toast(err.message, 'error');
  }
}

async function loadRawConfig() {
  const editor = $('#raw-editor');
  if (!editor) return;
  editor.value = '';
  editor.placeholder = '正在加载 config.toml…';
  try {
    const data = await api('/api/config/raw');
    editor.value = data.content || '';
    editor.placeholder = '';
    const path = $('#raw-path');
    if (path) path.textContent = data.path || '';
  } catch (err) {
    editor.placeholder = `读取失败：${err.message}`;
  }
}

async function saveRawConfig(event) {
  const btn = event.currentTarget;
  const content = $('#raw-editor').value;
  setLoading(btn, true, '保存中');
  try {
    const resp = await api('/api/config/raw', { method: 'PUT', body: { content } });
    applyConfigState(resp.config, resp.meta);
    const newToken = (resp.config.server.access_token || '').trim();
    if (newToken !== state.token) {
      state.token = newToken;
      if (newToken) localStorage.setItem('anifeed_token', newToken);
      else localStorage.removeItem('anifeed_token');
    }
    state.dirty = false;
    toast('config.toml 已保存并重载', 'success');
    await loadStatus().catch(() => {});
    renderSettingsPage();
  } catch (err) {
    toast(err.message, 'error', 6000);
  } finally {
    setLoading(btn, false);
  }
}

function openQrModal() {
  let pollTimer = null;
  let sessionId = null;
  const modal = openModal({
    title: '扫码登录 115',
    subtitle: '使用 115 手机 App 扫描二维码',
    size: 'narrow',
    body: `<div class="qr-wrap">
      <div id="qr-box"><div class="page-loading" style="padding:30px 0"><span class="spinner lg"></span><span>正在获取二维码…</span></div></div>
      <p id="qr-status">二维码有效期约 3 分钟</p>
      <button class="btn btn-sm hidden" id="qr-retry">${icon('refresh')}重新获取</button>
    </div>`,
    footer: `<button class="btn" data-action="modal-close">取消</button>`,
    onMount(el, close) {
      const stop = () => { if (pollTimer) { clearTimeout(pollTimer); pollTimer = null; } };
      const start = async () => {
        stop();
        const box = $('#qr-box', el);
        const status = $('#qr-status', el);
        const retry = $('#qr-retry', el);
        retry.classList.add('hidden');
        box.innerHTML = `<div class="page-loading" style="padding:30px 0"><span class="spinner lg"></span><span>正在获取二维码…</span></div>`;
        status.textContent = '二维码有效期约 3 分钟';
        try {
          const data = await api('/api/auth/qrcode', { method: 'POST', body: {} });
          sessionId = data.id;
          box.innerHTML = `<img class="qr-image" id="qr-image" alt="115 登录二维码" src="${data.image}">`;
          status.textContent = '请使用 115 App 扫码并确认登录';
          poll();
        } catch (err) {
          box.innerHTML = `<p class="form-error">${esc(err.message)}</p>`;
          retry.classList.remove('hidden');
        }
      };
      const poll = async () => {
        if (!sessionId) return;
        if (!document.body.contains(el)) return;
        try {
          const data = await api(`/api/auth/qrcode/${sessionId}`);
          if (data.status === 'waiting') {
            pollTimer = setTimeout(poll, 2000);
          } else if (data.status === 'success') {
            $('#qr-status', el).innerHTML = `<span class="badge badge-success">登录成功</span>`;
            $('#qr-box', el).innerHTML = `<div class="empty" style="padding:24px"><div class="empty-icon">${icon('check')}</div><strong>登录成功</strong><p>Cookies 已保存，正在重新加载配置。</p></div>`;
            setTimeout(() => { close(); state.config = null; state.meta = null; loadSettingsPage(); }, 900);
          } else {
            $('#qr-status', el).textContent = data.message || '二维码已失效';
            $('#qr-retry', el).classList.remove('hidden');
          }
        } catch (err) {
          $('#qr-status', el).textContent = err.message;
          $('#qr-retry', el).classList.remove('hidden');
        }
      };
      $('#qr-retry', el).addEventListener('click', start);
      el.addEventListener('click', event => {
        if (event.target.closest('[data-action="modal-close"]') || event.target === el) stop();
      });
      start();
    },
  });
  return modal;
}

/* ================================================================
 * Token gate & boot
 * ================================================================ */
function showTokenGate() {
  const gate = $('#token-gate');
  gate.classList.remove('hidden');
  setTimeout(() => $('#token-input').focus(), 60);
}

function hideTokenGate() {
  $('#token-gate').classList.add('hidden');
}

async function startApp() {
  await loadStatus();
  hideTokenGate();
  renderRoute();
}

function showAuthGate(auth) {
  hideTokenGate();
  const gate = $('#auth-gate');
  if (!gate) return;
  const configured = !!(auth && auth.configured);
  $('#auth-title').textContent = configured ? '' : '初始化登录';
  $('#auth-desc').textContent = configured ? '' : '首次使用，请设置 Web UI 登录账号和密码';
  $('#auth-title').classList.toggle('hidden', configured);
  $('#auth-desc').classList.toggle('hidden', configured);
  $('#auth-submit').textContent = configured ? '登录' : '创建账号并登录';
  $('#auth-confirm-field').classList.toggle('hidden', configured);
  $('#auth-error').classList.add('hidden');
  let remembered = null;
  try { remembered = JSON.parse(localStorage.getItem('anifeed_login') || 'null'); } catch { remembered = null; }
  const remember = !!(remembered && remembered.username && remembered.password);
  $('#auth-remember').checked = remember;
  $('#auth-user').value = remembered && remembered.username ? remembered.username : (configured && auth.username ? auth.username : 'admin');
  $('#auth-pass').value = remembered && remembered.password ? remembered.password : (configured ? '' : 'admin');
  $('#auth-confirm').value = configured ? '' : 'admin';
  gate.classList.remove('hidden');
  setTimeout(() => $('#auth-user').focus(), 80);
}

async function boot() {
  initChrome();
  try {
    const auth = await api('/api/auth/status');
    if (!auth.authenticated) {
      showAuthGate(auth);
      return;
    }
    await startApp();
  } catch (err) {
    if (err.code === 'token_required') {
      showTokenGate();
      return;
    }
    renderFatal(err.message);
  }
}

$('#auth-form').addEventListener('submit', async event => {
  event.preventDefault();
  const auth = await api('/api/auth/status').catch(() => ({}));
  const setup = !auth.configured;
  const username = $('#auth-user').value.trim();
  const password = $('#auth-pass').value;
  const error = $('#auth-error');
  error.classList.add('hidden');
  if (!username || !password) return;
  if (setup && password !== $('#auth-confirm').value) {
    error.textContent = '两次输入的密码不一致';
    error.classList.remove('hidden');
    return;
  }
  try {
    await api(setup ? '/api/auth/setup' : '/api/auth/login', { method: 'POST', body: { username, password } });
    if ($('#auth-remember').checked) {
      localStorage.setItem('anifeed_login', JSON.stringify({ username, password }));
    } else {
      localStorage.removeItem('anifeed_login');
    }
    $('#auth-gate').classList.add('hidden');
    await startApp();
    toast(setup ? '登录账号已创建' : '登录成功', 'success');
  } catch (err) {
    error.textContent = err.message || '登录失败';
    error.classList.remove('hidden');
  }
});

$('#token-form').addEventListener('submit', async event => {
  event.preventDefault();
  const value = $('#token-input').value.trim();
  const error = $('#token-error');
  if (!value) return;
  state.token = value;
  localStorage.setItem('anifeed_token', value);
  try {
    await loadStatus();
    error.classList.add('hidden');
    hideTokenGate();
    renderRoute();
    toast('令牌验证通过', 'success');
  } catch (err) {
    state.token = '';
    localStorage.removeItem('anifeed_token');
    error.textContent = err.message || '令牌错误';
    error.classList.remove('hidden');
  }
});

boot();
