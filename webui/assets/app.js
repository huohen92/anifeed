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
  ai: '<path d="M12 3.2l1.85 4.5 4.5 1.85-4.5 1.85L12 15.9l-1.85-4.5L5.65 9.55l4.5-1.85L12 3.2z"/><path d="M18 15l.8 1.95 1.95.8-1.95.8L18 20.5l-.8-1.95-1.95-.8 1.95-.8L18 15z"/>',
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
  ai: { title: 'AI 助手', subtitle: '查询订阅、任务和运行状态', icon: 'info' },
};

const state = {
  token: sessionStorage.getItem('anifeed_token') || '',
  status: null,
  config: null,
  meta: null,
  tasks: [],
  tasksLoaded: false,
  taskVisible: 20,
  history: [],
  historyLoaded: false,
  logs: [],
  logsLoaded: false,
  route: '',
  routeEpoch: 0,
  feedQuery: '',
  feedSource: '',
  feedStatus: 'enabled',
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
  matchTemplates: null,
  subgroupAliases: null,
  taskFilter: 'all',
  historyName: '',
  historyLimit: 50,
  settingsTab: 'form',
  settingsSites: null,
  globalExcludes: null,
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


/* ===== 字幕组别名 ===== */
async function loadSubgroupAliases(force) {
  if (!force && Array.isArray(state.subgroupAliases)) return state.subgroupAliases;
  const data = await api('/api/subgroup-aliases');
  state.subgroupAliases = Array.isArray(data.aliases) ? data.aliases : [];
  return state.subgroupAliases;
}

/* ===== 字幕组匹配模板 ===== */
async function loadMatchTemplates(force) {
  if (!force && Array.isArray(state.matchTemplates)) return state.matchTemplates;
  const data = await api('/api/match-templates');
  state.matchTemplates = Array.isArray(data.templates) ? data.templates : [];
  return state.matchTemplates;
}
function templateMatchesSubgroup(template, subgroup) {
  const target = String(subgroup || '').trim().toLowerCase();
  if (!target) return false;
  const candidates = [template.subgroup, ...(template.aliases || [])].map(value => String(value || '').trim().toLowerCase()).filter(Boolean);
  return candidates.includes(target);
}
function scopeTemplateRules(raw, subgroup) {
  const group = String(subgroup || '').trim();
  if (!group) return String(raw || '').trim();
  return String(raw || '').split(/\r?\n/).map(line => {
    line = line.trim();
    if (!line) return '';
    if (line.startsWith('{{')) return line;
    return '{{' + group + '}}:' + line;
  }).filter(Boolean).join('\n');
}

/* ===== 季度标注（仅前端显示，不参与后端重命名） ===== */
const QUARTER_KEY = 'anifeed.quarters.v3';
function quarterMap() {
  try { return JSON.parse(localStorage.getItem(QUARTER_KEY) || '{}') || {}; } catch (e) { return {}; }
}
function saveQuarterMap(map) {
  try { localStorage.setItem(QUARTER_KEY, JSON.stringify(map)); } catch (e) {}
}
function quarterOf(feed) {
  return (feed && feed.url && quarterMap()[feed.url]) || '';
}
function quarterLabel(q) {
  if (!q) return '';
  const parts = String(q).split('-');
  if (parts.length !== 2) return String(q);
  return parts[0] + '-' + parts[1];
}


/* 按实际放送日期推算季度：优先取与当前进度对应的那一集，其次取已播出的最后一集 */
function quarterFromEpisodeList(list, feed) {
  if (!Array.isArray(list) || !list.length) return '';
  const progress = (typeof feedProgress === 'function') ? feedProgress(feed) : {};
  const target = Number(progress.mainDone || 0) || Number(progress.allDone || 0) || 0;
  const season = Number((feed && feed.rename_season) || 0);
  let pool = season > 0 ? list.filter(function (e) { return Number(e.season) === season; }) : [];
  if (!pool.length) pool = list.slice();
  let pick = null;
  if (target > 0) pick = pool.find(function (e) { return Number(e.episode) === target; }) || null;
  if (!pick) {
    const today = new Date().toISOString().slice(0, 10);
    const aired = pool.filter(function (e) { return e.air_date && e.air_date <= today; });
    pick = aired.length ? aired[aired.length - 1] : null;
  }
  return pick ? quarterFromDate(pick.air_date) : '';
}
function quarterFromDate(dateStr) {
  if (!dateStr) return '';
  const m = String(dateStr).match(/^(\d{4})-(\d{2})/);
  if (!m) return '';
  const month = Number(m[2]);
  const q = month <= 3 ? '01' : month <= 6 ? '04' : month <= 9 ? '07' : '10';
  return m[1] + '-' + q;
}

function setQuarter(feed, q) {
  const map = quarterMap();
  if (!feed || !feed.url) return;
  if (q) map[feed.url] = q; else delete map[feed.url];
  saveQuarterMap(map);
}

function quarterFilterValues() {
  const map = quarterMap();
  return [...new Set(Object.keys(map).map(function (k) { return map[k]; }))].filter(Boolean).sort().reverse();
}

function quarterOptions(selected) {
  const labels = { 1: '01', 2: '04', 3: '07', 4: '10' };
  const now = new Date();
  let y = now.getFullYear();
  let q = Math.floor(now.getMonth() / 3) + 1;
  const out = [];
  for (let i = 0; i < 20; i++) {
    const month = labels[q];
    out.push({ value: y + '-' + month, label: y + '-' + month });
    q -= 1;
    if (q < 1) { q = 4; y -= 1; }
  }
  if (selected && !out.some(function (item) { return item.value === selected; })) {
    out.unshift({ value: selected, label: quarterLabel(selected) });
  }
  return out;
}

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

function magnetFromHash(hash) {
  const raw = String(hash || '').trim();
  if (!raw) return '';
  if (/^(magnet:|ed2k:|https?:)/i.test(raw)) return raw;
  return 'magnet:?xt=urn:btih:' + raw;
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

  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 60000);
  opts.signal = controller.signal;
  let res;
  try {
    res = await fetch(path, opts);
  } catch (err) {
    if (err && err.name === 'AbortError') throw new Error('请求超时，请稍后重试');
    throw new Error('无法连接 AniFeed 服务，请确认服务仍在运行');
  } finally {
    clearTimeout(timeout);
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

const modalStack = [];
let modalTitleSeq = 0;

function openModal({ title, subtitle = '', body = '', footer = '', size = '', onMount }) {
  const root = $('#modal-root');
  const backdrop = document.createElement('div');
  backdrop.className = 'modal-backdrop';
  const titleId = 'modal-title-' + (++modalTitleSeq);
  backdrop.innerHTML = `
    <div class="modal ${size}" role="dialog" aria-modal="true" aria-labelledby="${titleId}" tabindex="-1">
      <div class="modal-head">
        <div><h3 id="${titleId}">${esc(title)}</h3>${subtitle ? `<p>${esc(subtitle)}</p>` : ''}</div>
        <button class="icon-btn close" data-action="modal-close" aria-label="关闭">${icon('close')}</button>
      </div>
      <div class="modal-body">${body}</div>
      ${footer ? `<div class="modal-foot">${footer}</div>` : ''}
    </div>`;
  root.appendChild(backdrop);
  requestAnimationFrame(() => backdrop.classList.add('show'));

  const previousFocus = document.activeElement;
  const handle = { el: backdrop, close: null, closed: false };
  modalStack.push(handle);
  const focusable = () => Array.from(backdrop.querySelectorAll('button:not([disabled]), a[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])')).filter(el => el.offsetParent !== null);
  const close = () => {
    if (handle.closed) return;
    handle.closed = true;
    const index = modalStack.indexOf(handle);
    if (index >= 0) modalStack.splice(index, 1);
    backdrop.classList.remove('show');
    document.removeEventListener('keydown', onKey);
    setTimeout(() => backdrop.remove(), 180);
    if (previousFocus && previousFocus.isConnected && modalStack.length === 0) previousFocus.focus({ preventScroll: true });
  };
  handle.close = close;
  const onKey = event => {
    if (modalStack[modalStack.length - 1] !== handle) return;
    if (event.key === 'Escape') { event.preventDefault(); close(); return; }
    if (event.key !== 'Tab') return;
    const items = focusable();
    if (!items.length) { event.preventDefault(); backdrop.querySelector('.modal')?.focus(); return; }
    const first = items[0], last = items[items.length - 1];
    if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
    else if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
  };
  document.addEventListener('keydown', onKey);
  backdrop.addEventListener('click', event => {
    if (event.target === backdrop) close();
    if (event.target.closest('[data-action="modal-close"]')) close();
  });

  if (onMount) onMount(backdrop, close, handle);
  requestAnimationFrame(() => { const items = focusable(); (items[0] || backdrop.querySelector('.modal'))?.focus({ preventScroll: true }); });
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

function initMobileKeyboardGuard() {
  const lockBackgroundViewport = () => {
    const root = document.documentElement;
    root.style.setProperty('--af-bg-width', window.innerWidth + 'px');
    root.style.setProperty('--af-bg-height', window.innerHeight + 'px');
  };
  const update = () => {
    const vv = window.visualViewport;
    const keyboard = vv ? Math.max(0, window.innerHeight - vv.height - (vv.offsetTop || 0)) : 0;
    const focused = document.activeElement && /^(INPUT|TEXTAREA|SELECT)$/.test(document.activeElement.tagName);
    const keyboardOpen = !!focused && keyboard > 120;
    document.body.classList.toggle('keyboard-open', keyboardOpen);
    if (!keyboardOpen) lockBackgroundViewport();
  };
  if (window.visualViewport) {
    window.visualViewport.addEventListener('resize', update);
    window.visualViewport.addEventListener('scroll', update);
  }
  document.addEventListener('focusin', () => setTimeout(update, 50));
  document.addEventListener('focusout', () => setTimeout(update, 120));
  window.addEventListener('resize', () => { if (!document.body.classList.contains('keyboard-open')) lockBackgroundViewport(); });
  lockBackgroundViewport();
  update();
}

function initChrome() {
  $$('[data-icon]').forEach(el => { el.innerHTML = icon(el.dataset.icon); });
  initTheme();
  initMobileKeyboardGuard();

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
  document.body.dataset.route = route;
  window.scrollTo(0, 0);
  state.routeEpoch = (state.routeEpoch || 0) + 1;
  const meta = PAGE_META[route];
  document.title = 'AniFeed';
  $('#page-title').textContent = meta.title;
  $('#page-subtitle').textContent = meta.subtitle;
  $$('.nav-item').forEach(el => el.classList.toggle('active', el.dataset.route === route));
  $$('[data-mobile-route]').forEach(el => el.classList.toggle('active', el.dataset.mobileRoute === route));
  document.querySelectorAll('#ai-fab, #ai-fab-mobile').forEach(el => el.classList.toggle('active', route === 'ai'));
  document.body.classList.toggle('hide-topbar-actions', route === 'logs');

  switch (route) {
    case 'dashboard': loadDashboard(); break;
    case 'feeds': loadFeedsPage(); break;
    case 'tasks': loadTasksPage(); break;
    case 'history': loadHistoryPage(); break;
    case 'logs': loadLogsPage(); break;
    case 'settings': loadSettingsPage(); break;
    case 'ai': loadAIPage(); break;
    default: loadDashboard();
  }
}

function loadAssetScript(src) {
  return new Promise((resolve, reject) => {
    const existing = document.querySelector('script[data-lazy-module="' + src + '"]');
    if (existing) { existing.addEventListener('load', resolve, { once: true }); existing.addEventListener('error', reject, { once: true }); return; }
    const script = document.createElement('script');
    script.src = src;
    script.dataset.lazyModule = src;
    script.onload = resolve;
    script.onerror = () => reject(new Error('加载模块失败: ' + src));
    document.body.appendChild(script);
  });
}

async function ensureAIModule() {
  if (window.AniFeedAI) return;
  const app = document.querySelector('script[src*="app.min.js"]');
  const query = app ? (new URL(app.src).search || '') : '';
  await loadAssetScript('/app-ai.min.js' + query);
}

async function loadAIPage() {
  await ensureAIModule();
  return window.AniFeedAI.loadPage();
}

function syncAINav() {
  const apply = ai => {
    state.aiEnabled = !!((ai || {}).enabled);
    document.querySelectorAll('#ai-fab, #ai-fab-mobile').forEach(el => el.classList.toggle('is-disabled', !state.aiEnabled));
  };
  if (state.config) { apply((state.config.ai || {})); return; }
  api('/api/config').then(r => apply((r.config || {}).ai)).catch(err => console.debug('[AniFeed]', err));
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
function skeletonAnimeCard(kind = 'dashboard') {
  const badges = kind === 'feeds'
    ? '<div class="skeleton-line short" style="height:12px;margin:10px 9px"></div>'
    : '<div class="skeleton-line short" style="height:11px;margin:9px 9px 0"></div><div class="skeleton-line" style="height:11px;margin:5px 9px"></div><div class="skeleton-line short" style="height:11px;margin:0 9px 9px"></div>';
  return '<article class="anime-card skeleton-card"><div class="skeleton-poster"></div>' + badges + '</article>';
}

function skeletonStats(count = 4) {
  return '<div class="stats">' + Array.from({ length: count }).map(() => '<div class="stat skeleton-stat"><div class="skeleton-line short"></div><div class="skeleton-line"></div><div class="skeleton-line short"></div></div>').join('') + '</div>';
}

function skeletonToolbar() {
  return '<div class="toolbar skeleton-toolbar"><div class="skeleton-block skeleton-toolbar-search"></div><div class="skeleton-block skeleton-toolbar-pill"></div><div class="skeleton-block skeleton-toolbar-pill"></div><span class="spacer"></span><div class="skeleton-block skeleton-toolbar-btn"></div><div class="skeleton-block skeleton-toolbar-btn"></div></div>';
}

function pageSkeleton(kind = 'list') {
  const grid = '<div class="feed-grid">' + Array.from({ length: 6 }).map(() => skeletonAnimeCard('feeds')).join('') + '</div>';
  if (kind === 'dashboard') {
    return skeletonStats(4)
      + '<section class="dashboard-anime"><div class="card-head"><div><div class="skeleton-line short"></div><div class="skeleton-line short"></div></div></div><div class="anime-grid">' + Array.from({ length: 6 }).map(() => skeletonAnimeCard('dashboard')).join('') + '</div></section>'
      + '<section class="card"><div class="section-body"><div class="skeleton-line"></div><div class="skeleton-line short"></div><div class="skeleton-line"></div></div></section>';
  }
  if (kind === 'feeds') {
    return skeletonToolbar() + '<div class="skeleton-line short" style="margin:2px 2px 10px;height:16px"></div>' + grid;
  }
  if (kind === 'tasks') {
    const rows = Array.from({ length: 5 }).map(() => '<div class="skeleton-table-row"><div class="skeleton-line short"></div><div class="skeleton-line"></div><div class="skeleton-line short"></div></div>').join('');
    return skeletonToolbar() + '<div class="card"><div class="section-body">' + rows + '</div></div><div class="card"><div class="section-body"><div class="skeleton-line short"></div><div class="skeleton-line"></div></div></div>';
  }
  if (kind === 'history') {
    const rows = Array.from({ length: 6 }).map(() => '<div class="skeleton-table-row"><div class="skeleton-line short"></div><div class="skeleton-line"></div><div class="skeleton-line short"></div></div>').join('');
    return skeletonToolbar() + '<div class="card"><div class="section-body">' + rows + '</div></div>';
  }
  if (kind === 'logs') {
    const bars = Array.from({ length: 12 }).map(() => '<div class="skeleton-log-line"></div>').join('');
    return '<div class="log-box skeleton-log-box"><div class="log-head"><span class="dots"><i></i><i></i><i></i></span><span>AniFeed · server</span><span class="spacer"></span><div class="skeleton-log-control"></div><div class="skeleton-log-control short"></div></div><div class="log-filter-bar"><div class="skeleton-log-control"></div><div class="skeleton-log-control wide"></div></div><div class="log-skeleton-body">' + bars + '</div></div>';
  }
  return skeletonToolbar() + grid;
}

function renderBootSkeleton() {
  const route = currentRoute();
  if (route === 'dashboard') { renderDashboardSkeleton(); return; }
  if (route === 'settings') { renderSettingsSkeleton(); return; }
  $('#view').innerHTML = pageSkeleton(route);
}

function renderDashboardSkeleton() {
  $('#view').innerHTML = pageSkeleton('dashboard');
}

async function loadDashboard() {
  const epoch = state.routeEpoch;
  const hasCache = !!state.status;
  if (hasCache) {
    renderDashboard();
  } else {
    renderDashboardSkeleton();
  }

  let statusError = null;
  await Promise.allSettled([
    loadStatus().catch(err => { statusError = err; }),
    ensureConfig().catch(err => console.debug('[AniFeed]', err)),
  ]);
  if (!state.config) {
    try { await ensureConfig(); } catch { /* status-only */ }
  }
  if (state.routeEpoch !== epoch) return;
  if (!state.status) {
    renderFatal(statusError ? statusError.message : '无法读取服务状态');
    return;
  }
  renderDashboard();
  scheduleLoginPoll();

  // Tasks/history/logs are not needed for the first paint. Load them after
  // the shell is visible so a slow 115 request cannot block the dashboard.
  const dataJobs = [
    api('/api/history?limit=200').then(data => { state.history = data.items || []; if (state.history.length) state.historyLoaded = true; }).catch(err => console.debug('[AniFeed]', err)),
    api('/api/logs?limit=12').then(data => { state.logs = data.lines || []; }).catch(err => console.debug('[AniFeed]', err)),
    api('/api/tasks?limit=200').then(data => { state.tasks = data.tasks || []; state.tasksLoaded = true; }).catch(() => { state.tasks = []; }),
  ];
  Promise.allSettled(dataJobs).then(() => {
    if (state.routeEpoch !== epoch || state.route !== 'dashboard') return;
    renderDashboard();
  });
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
  const todayFeeds = allFeeds.filter(feed => todayHistory.some(item =>
    (item.config_url && item.config_url === feed.url) || (item.config_name && item.config_name === feed.name)
  ));
  const animeFeeds = todayFeeds.length ? todayFeeds : allFeeds.slice(0, 12);
  const animeHeading = todayFeeds.length ? '今天的订阅' : '订阅一览';
  const animeSub = todayFeeds.length ? '今天有更新的番剧 · 点击卡片可编辑' : '当前全部订阅 · 今天暂无更新';
  const animeGrid = animeFeeds.length
    ? `<div class="anime-grid">${animeFeeds.map(feed => dashboardAnimeCard(feed, allFeeds.indexOf(feed))).join('')}</div>`
    : `<div class="card">${emptyBlock('clock', '还没有订阅', '添加 RSS 订阅后，番剧会显示在这里。', `<a class="btn btn-primary" href="#/feeds">添加订阅</a>`)}</div>`;

  const logs = (state.logs || []).slice(-9);
  const logHtml = logs.length
    ? `<pre class="log-pre" style="max-height:260px">${renderLogLines(logs)}</pre>`
    : emptyBlock('terminal', '暂无日志', '运行一次订阅后即可看到服务端输出。');

  $('#view').innerHTML = `
    ${banners.join('')}
    <section class="stats">
      <div class="card stat accent">
        <span class="stat-icon">${icon('rss')}</span>
        <div class="stat-copy">
          <div class="stat-top">订阅总数</div>
          <div class="stat-value">${status.feed_count}</div>
          <div class="stat-foot">${status.auto_enabled ? `自动运行每 ${status.auto_interval_minutes} 分钟` : '自动运行未开启'}</div>
        </div>
      </div>
      <div class="card stat green">
        <span class="stat-icon">${icon('download')}</span>
        <div class="stat-copy">
          <div class="stat-top">今日提交</div>
          <div class="stat-value">${status.today_count}</div>
          <div class="stat-foot">历史累计 ${status.history_count} 条</div>
        </div>
      </div>
      <div class="card stat amber">
        <span class="stat-icon">${icon('cloud')}</span>
        <div class="stat-copy">
          <div class="stat-top">云端任务</div>
          <div class="stat-value">${state.tasks.length}</div>
          <div class="stat-foot">${stats.running} 个运行中 · ${stats.done} 个已完成</div>
        </div>
      </div>
      <div class="card stat cyan">
        <span class="stat-icon">${icon('clock')}</span>
        <div class="stat-copy">
          <div class="stat-top">提交成功率</div>
          <div class="stat-value">${state.tasks.length ? Math.round((stats.done / Math.max(1, state.tasks.length)) * 100) + '%' : '—'}</div>
          <div class="stat-foot">${stats.failed} 个失败任务</div>
        </div>
      </div>
    </section>

    <section class="dashboard-anime">
      <div class="card-head">
        <div><h3>${animeHeading}</h3><div class="sub">${animeSub}</div></div>
        <div class="right"><a class="btn btn-sm btn-ghost" href="#/feeds">查看全部 ${icon('arrow')}</a></div>
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
  const update = feedUpdateInfo(feed);
  const completed = !!update.completed;
  const disabled = !enabled || completed;
  return `<article class="anime-card ${disabled ? 'disabled' : ''}" data-dashboard-feed="${index}">
    <div class="anime-poster ${cached && cached.cover ? 'has-image' : ''}" data-feed-cover data-index="${index}">
      <span class="poster-fallback"><img src="/logo.webp" alt="AniFeed"></span>
      <img class="poster-image" alt="${esc(cached ? cached.title || '' : '')}" ${cached && cached.cover ? `src="${esc(cached.cover)}"` : ''} loading="lazy" referrerpolicy="no-referrer">
      <span class="anime-badge">${completed ? '已完成' : disabled ? '未启用' : '追番中'}</span>
    </div>
    <div class="anime-body">
      <h4 title="${esc(feed.name)}">${esc(feed.name || '未命名订阅')}</h4>
      <div class="anime-meta">${esc(source.label)} · ${esc(feed.site || '自动识别')}</div>
      <div class="anime-foot">
        <span class="badge">${feed.rename_season ? `第 ${feed.rename_season} 季` : '自动季数'}</span>
        ${(() => { const p = feedProgress(feed); const t = p.total || '*'; return `<span class="badge badge-main" title="主 RSS">主 ${p.mainDone}/${t}</span><span class="badge badge-all" title="含备用 RSS">合计 ${p.allDone}/${t}</span>`; })()}
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


/* ===== 播出季度补全：按 TMDB 首播日期批量生成 YYYY-MM ===== */
let quarterHydrating = false;
async function hydrateQuarters() {
  if (quarterHydrating || !state.config || !Array.isArray(state.config.rss)) return;
  const pending = state.config.rss.filter(function (feed) {
    return feed && feed.tmdb_id && feed.url && !quarterOf(feed);
  });
  if (!pending.length) return;
  quarterHydrating = true;
  for (const feed of pending.slice(0, 40)) {
    try {
      const type = feed.tmdb_type || 'tv';
      const season = Number(feed.rename_season || 0);
      const seasonQuery = type === 'tv' && season > 1 ? '&season=' + season : '';
      let q = '';
      try {
        const eps = await api('/api/tmdb/episodes?id=' + encodeURIComponent(feed.tmdb_id));
        q = quarterFromEpisodeList(eps.episodes, feed);
      } catch (e) { /* 集数表不可用时退回首播日期 */ }
      if (!q) {
        const data = await api('/api/tmdb/detail?id=' + encodeURIComponent(feed.tmdb_id) + '&type=' + encodeURIComponent(type) + seasonQuery);
        q = quarterFromDate(data.season_air_date || data.first_air_date || data.release_date || '');
      }
      if (q) setQuarter(feed, q);
    } catch (e) { /* 单个失败忽略 */ }
  }
  quarterHydrating = false;
  if (location.hash.indexOf('#/feeds') >= 0) renderFeedsPage();
}

function applyConfigState(config, meta) {
  state.config = config;
  if (meta) state.meta = meta;
  document.body.classList.toggle('hide-hints', !!(config && config.lab && config.lab.hide_hints));
  setTimeout(hydrateQuarters, 600);
  state.settingsSites = Object.entries((config && config.sites) || {}).map(([host, site]) => ({
    host,
    https_agent: !!site.https_agent,
    headers: Object.entries(site.headers || {}).map(([k, v]) => ({ key: k, value: v })),
  }));
  // A slow first config response must repaint the overview once it arrives.
  if (state.route === 'dashboard' && state.status && document.querySelector('.stat')) {
    try { renderDashboard(); } catch { /* ignore transient repaint errors */ }
  }
}

async function loadFeedsPage() {
  const epoch = state.routeEpoch;
  if (state.config) {
    renderFeedsPage();
  } else {
    $('#view').innerHTML = pageSkeleton('feeds');
  }
  try {
    await ensureConfig();
    await loadStatus().catch(err => console.debug('[AniFeed]', err));
    await api('/api/history?limit=500').then(data => { state.history = data.items || []; state.historyLoaded = true; }).catch(err => console.debug('[AniFeed]', err));
  } catch (err) {
    if (state.routeEpoch === epoch) renderFatal(err.message);
    return;
  }
  if (state.routeEpoch !== epoch) return;
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
    openNewFeedWithRules({
      source: pick.source,
      name: pick.title,
      subgroup_name: pick.subgroup || '',
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

function feedProgress(feed) {
  const mainURL = String(feed.url || '').trim();
  const groups = (feed.subgroups || []).filter(group => group.enabled !== false && String(group.url || '').trim());
  const mainOffset = Number(feed.rename_offset || (state.config && state.config.rename ? state.config.rename.offset : 0) || 0);
  const groupByURL = new Map(groups.map(group => [String(group.url || '').trim(), group]));
  const mainEpisodes = new Set();
  const allEpisodes = new Set();
  for (const item of state.history || []) {
    const itemURL = String(item.config_url || '').trim();
    const group = groupByURL.get(itemURL);
    const nameMatch = feed.name && item.config_name === feed.name;
    if (!nameMatch && itemURL !== mainURL && !group) continue;
    const episode = parseEpisodeNumber(item.title);
    if (!episode) continue;
    const isMain = itemURL === mainURL || (!itemURL && !group);
    const offset = isMain ? mainOffset : Number(group && group.offset != null ? group.offset : mainOffset);
    const adjusted = episode + offset;
    if (adjusted <= 0) continue;
    allEpisodes.add(adjusted);
    if (isMain) mainEpisodes.add(adjusted);
  }
  return {
    mainDone: mainEpisodes.size,
    allDone: allEpisodes.size,
    total: Number(feed.total_episodes) || 0,
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
    if (btn) btn.addEventListener('click', () => openBrowseModal('anibt', pick => openNewFeedWithRules({
      source: pick.source, name: pick.title, subgroup_name: pick.subgroup || '', url: pick.url,
    })));
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
  const img = el.querySelector('.poster-image');
  if (!img) return;
  img.alt = data.title || '';
  const sameSource = img.getAttribute('src') === data.cover;
  const reveal = () => { if (img.naturalWidth > 0) el.classList.add('has-image'); };
  if (!sameSource) {
    el.classList.remove('has-image');
    img.addEventListener('load', reveal, { once: true });
    img.addEventListener('error', () => el.classList.remove('has-image'), { once: true });
    img.src = data.cover;
  } else if (el.classList.contains('has-image')) {
    reveal();
  } else {
    img.addEventListener('load', reveal, { once: true });
    img.addEventListener('error', () => el.classList.remove('has-image'), { once: true });
    if (img.complete) reveal();
  }
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
  const update = feedUpdateInfo(feed);
  const progress = feedProgress(feed);
  const totalText = progress.total || '*';
  const subgroup = subgroupLabel(feed.subgroup_name || urlSubgroupName(feed.url), feed.name);
  return `<article class="anime-card feed-manage-card ${disabled ? 'disabled' : ''}" title="${esc(feed.name)}">
    <div class="anime-poster ${cachedCover && cachedCover.cover ? 'has-image' : ''}" data-feed-cover data-feed-poster data-index="${index}" data-source="${source.key}" title="点击编辑订阅">
      <span class="poster-fallback"><img src="/logo.webp" alt="AniFeed"></span>
      <img class="poster-image" alt="${esc(cachedCover ? cachedCover.title || '' : '')}" ${cachedCover && cachedCover.cover ? `src="${esc(cachedCover.cover)}"` : ''} loading="lazy" referrerpolicy="no-referrer">
      <span class="anime-badge">${completed ? '已完成' : disabled ? '未启用' : '追番中'}</span>
      <div class="feed-cover-overlay">
        <h4>${esc(feed.name || '未命名订阅')}</h4>
        <div>第 ${Number(feed.rename_season) || 1} 季${subgroup ? ` · ${esc(subgroup)}` : ''}</div>
        <div class="feed-cover-progress">主 ${progress.mainDone} · 合计 ${progress.allDone}/${totalText}</div>
        <div class="feed-card-meta" title="季 ${Number(feed.rename_season) || 1} · 播出 ${esc(quarterOf(feed) || '—')} · ${esc(subgroup || '—')} · ${completed ? '已完成' : disabled ? '未启用' : '已启用'} · ${progress.allDone}/${totalText} · ${esc(feed.tmdb_type || 'tv')}">
          <span>S${Number(feed.rename_season) || 1}</span><span>${esc(quarterOf(feed) || '—')}</span><span>${esc(subgroup || '—')}</span><span>${completed ? '已完成' : disabled ? '未启用' : '已启用'}</span><span>${progress.allDone}/${totalText}</span><span>${esc(feed.tmdb_type || 'tv')}</span>
        </div>
      </div>
      <button class="poster-action poster-action-left" type="button" data-poster-action="poster" title="修改海报">${icon('image')}</button>
      <button class="poster-action poster-action-right" type="button" data-poster-action="edit" title="编辑订阅">${icon('edit')}</button>
      <span class="poster-score hidden"></span>
    </div>
    <div class="feed-card-tags">
      <span class="badge">${esc(source.label)}</span>
      <span class="badge">主 ${progress.mainDone} · 合计 ${progress.allDone}/${totalText}</span>
      ${update.completed ? '<span class="badge badge-warning">已完成</span>' : ''}
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
              <div class="browse-group-name ${releases ? 'browse-group-toggle' : ''}" ${releases ? 'data-toggle-group role="button" tabindex="0" aria-expanded="false"' : ''}>
                <strong>${esc(group.name)}</strong>
                <span>${esc(group.meta || '')}${group.latest ? ` · 最近 ${fmtDate(group.latest)}` : ''}</span>
              </div>
              <button class="btn btn-sm btn-primary" data-add-group="${index}">${icon('plus')}添加订阅</button>
            </div>
            ${releases ? `<div class="browse-releases hidden">${releases}</div>` : ''}
          </div>`;
        }).join('');

        $$('[data-add-group]', container).forEach(btn => btn.addEventListener('click', () => {
          const group = groups[Number(btn.dataset.addGroup)];
          close();
          onPick({ source: view.source, id: item.id, title: item.title, subgroup: group.name, url: group.rss });
        }));
        container.querySelectorAll('[data-toggle-group]').forEach(toggle => {
          const group = toggle.closest('.browse-group');
          const list = group && group.querySelector('.browse-releases');
          if (!list) return;
          const setOpen = open => {
            list.classList.toggle('hidden', !open);
            toggle.classList.toggle('open', open);
            toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
          };
          toggle.addEventListener('click', () => setOpen(list.classList.contains('hidden')));
          toggle.addEventListener('keydown', event => {
            if (event.key === 'Enter' || event.key === ' ') {
              event.preventDefault();
              setOpen(list.classList.contains('hidden'));
            }
          });
        });
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

function openQBHelp() {
  openModal({
    title: 'qBittorrent 设置说明',
    subtitle: '115 失败后的本地兜底下载器',
    size: 'browse',
    body: `<div class="notice">${icon('info')}<span>运行方式：115 离线任务连续失败达到阈值后，AniFeed 通过 qB WebUI API 提交同一个磁力/种子，并记录提交状态。qB 完成下载后，AniFeed 只读取任务状态，不会直接移动 qB 文件。</span></div>
      <div class="table-wrap" style="margin-top:14px"><table class="table"><thead><tr><th>设置</th><th>如何填写</th></tr></thead><tbody>
        <tr><td>WebUI 地址</td><td>填 qB 容器内可访问的地址，如 <code>http://127.0.0.1:8080</code>；不要填宿主机路径或带 <code>/api</code>。</td></tr>
        <tr><td>API Key</td><td>qB 5.2+ 在 WebUI 设置中生成，通常以 <code>qbt_</code> 开头。填写后优先使用 Bearer API Key，不再需要用户名密码。</td></tr>
        <tr><td>用户名 / 密码</td><td>仅在未填写 API Key 时使用；旧版 qB 可用。</td></tr>
        <tr><td>保存路径</td><td>填 qB 容器内路径，如 <code>/downloads</code>。留空使用 qB 默认保存路径。不是宿主机路径。</td></tr>
        <tr><td>分类 / 标签</td><td>可选，用于在 qB 中筛选 AniFeed 任务；标签多个用逗号分隔。</td></tr>
        <tr><td>内容布局</td><td>Original 保留种子目录；Subfolder 创建任务子文件夹；NoSubfolder 不创建子文件夹；默认交给 qB。</td></tr>
        <tr><td>限速</td><td>单位 KB/s；填 0 表示不覆盖 qB 全局限速。</td></tr>
        <tr><td>分享率 / 做种时长 / 非活跃时长</td><td>填 0 使用 qB 全局设置；做种时长单位分钟。</td></tr>
        <tr><td>使用 qB 临时下载路径</td><td>开启后使用 qB 的临时下载目录；最终目录仍由保存路径决定。需先在 qB 内配置临时目录。</td></tr>
        <tr><td>115 失败阈值</td><td>同一个 115 失败任务连续检查达到该次数后提交 qB；默认 3 次。</td></tr>
      </tbody></table></div>
      <div class="notice" style="margin-top:14px">${icon('alert')}<span>如果以后要做 qB 下载完成后的文件重命名/移动，需要让 qB 和 AniFeed 访问同一个目录，或配置路径映射。当前版本仅提交任务、查询完成状态和展示兜底记录，不会移动 qB 文件。</span></div>`,
    footer: '<button class="btn btn-primary" data-action="modal-close">知道了</button>',
  });
}

function openNotifyPreview() {
  const raw = ($('#set-notify-template') || {}).value || DEFAULT_NOTIFY_TEMPLATE;
  const fields = {
    action: '添加至云下载', title: '青之箱', score: '8.2', themoviedbName: '青之箱',
    bgmUrl: 'https://bgm.tv/subject/548156', season: '1', episode: '6',
    subgroup: 'LoliHouse', currentEpisodeNumber: '6', totalEpisodeNumber: '37',
    year: '2026', month: '11', date: '27', text: '[LoliHouse] 青之箱 - 06',
    bgmEpisodeTitle: '第6集', failureCount: '0',
  };
  const rendered = raw.replace(/\$\{(\w+)\}/g, (m, key) => fields[key] === undefined ? m : fields[key])
    .split('\n').filter(line => line.trim() && !/^[^：:]+[：:]\s*$/.test(line.trim())).join('\n');
  const summary = `AniFeed 运行完成
状态：成功
订阅 4 · 发现 3 · 提交 2
时间：2026-11-27 07:40

更新：
- 青之箱：提交 1
  [LoliHouse] 青之箱 - 06`;
  openModal({
    title: '通知预览',
    subtitle: '条目模板与运行汇总的实际排版',
    size: 'browse',
    body: `<div class="section-subtitle">条目通知（模板）</div><pre class="log-pre" style="max-height:240px;border-radius:10px">${esc(rendered)}</pre>
      <div class="section-subtitle" style="margin-top:14px">运行汇总（非模板）</div><pre class="log-pre" style="max-height:240px;border-radius:10px">${esc(summary)}</pre>`,
    footer: '<button class="btn btn-primary" data-action="modal-close">关闭</button>',
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


function openNewFeedWithRules(prefill = {}) {
  const context = {
    ...prefill,
    subgroup_name: prefill.subgroup_name || prefill.subgroup || '',
    filter: prefill.filter || '',
    exclude: prefill.exclude || '',
  };
  const loading = openModal({
    title: '规则与匹配',
    subtitle: '先配置规则，再填写订阅基本信息',
    size: 'rules',
    body: '<div class="stats"><div class="stat"><div class="skeleton-line short"></div><div class="skeleton-line"></div></div><div class="stat"><div class="skeleton-line short"></div><div class="skeleton-line"></div></div></div><div class="card"><div class="section-body"><div class="skeleton-line"></div><div class="skeleton-line short"></div><div class="skeleton-line"></div></div></div>',
  });
  openRulesModal((filter, exclude) => {
    openFeedModal(-1, { ...prefill, filter, exclude });
  }, () => context).finally(() => loading.close());
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
    const summary = `RSS ${items.length} 条 · 规则匹配 ${data.matched} 条 · 实际提交 ${Number(data.selected_count) || 0} 条`;
    const rows = visible.map(item => {
      const renamed = item.renamed
        ? `<code class="preview-renamed">${esc(item.renamed)}</code>`
        : `<span class="form-error">${esc(item.error || '无法生成名称')}</span>`;
      return `<tr class="${item.filtered ? 'rule-filtered' : ''} ${item.suggested ? '' : 'preview-skipped'}">
        <td class="numeric" data-label="集数">${item.episode ? esc(String(item.episode)) : '—'}</td>
        <td class="name-cell" data-label="标题"><div class="title">${esc(item.title)}</div></td>
        <td data-label="重命名后">${renamed}</td>
      </tr>`;
    }).join('');
    container.innerHTML = `
      <div class="rule-results-head">
        <span>${summary}</span>
        <span class="spacer"></span>
        <button type="button" class="btn btn-sm ${mode === 'all' ? 'btn-primary' : ''}" id="rule-show-all">显示 RSS 全部内容</button>
      </div>
      ${(data.missing || []).length ? `<div class="notice" style="margin:6px 0 8px">${icon('alert')}<span>缺少第 <b>${esc(formatMissingEpisodes(data.missing))}</b> 集</span></div>` : ''}
      ${!data.one_per_episode ? `<div class="notice" style="margin:6px 0 8px">${icon('info')}<span>当前仅识别并记录，不实际跳过；实际会提交全部规则匹配项。开启「每集只下载一个版本」后才按自动建议提交。</span></div>` : ''}
      ${visible.length ? `<div class="table-wrap rule-result-table-wrap">
        <table class="table rename-preview-table rule-result-compact">
          <thead><tr><th>集数</th><th>RSS / BT 标题</th><th>重命名后</th></tr></thead>
          <tbody>${rows}</tbody>
        </table>
      </div>` : `<div class="empty" style="padding:22px"><span class="hint">当前视图没有条目。</span></div>`}
      ${filtered.length > visible.length ? `<div class="hint" style="margin-top:8px">还有 ${filtered.length - visible.length} 条未显示</div>` : ''}`;
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
          <div class="item-preview-source"><span class="badge badge-primary role-badge">主 RSS</span>${(feed.subgroup_name || urlSubgroupName(feed.url)) ? `<strong>${esc(feed.subgroup_name || urlSubgroupName(feed.url))}</strong>` : ''}<code title="${esc(feed.url || '')}">${esc(feed.url || '')}</code></div>
          ${backupSources.map(group => `<div class="item-preview-source"><span class="badge role-badge">备用 RSS</span><strong>${esc(subgroupLabel(group.name, feed.name) || '备用字幕组')}</strong><code title="${esc(group.url)}">${esc(group.url)}</code></div>`).join('')}
        </div>
      </details>`;
      const previewSources = [
        { role: 'main', subgroup: feed.subgroup_name || '', url: feed.url, source: feed.source, filter: feed.filter || '', exclude: feed.exclude || '', season: feed.rename_season || 0, offset: feed.rename_offset || 0 },
        ...backupSources.map(group => ({ role: 'backup', subgroup: group.name || '', url: group.url, source: feed.source, filter: group.filter || '', exclude: group.exclude || '', season: group.season || feed.rename_season || 0, offset: group.offset || 0 })),
      ];
      Promise.all(previewSources.map((source, order) => {
        const params = new URLSearchParams({
          name: feed.name || '', url: source.url || '', source: source.source || '',
          subgroup: source.subgroup || '',
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
              <td class="nowrap preview-col-status" data-label="是否下载">${item.filtered ? '<span class="badge badge-warning">过滤</span>' : (item.coveredBy && !isAllowed ? `<span class="badge" title="已由 ${esc(item.coveredBy.source_subgroup || '主 RSS')} 覆盖">否</span>` : `<button type="button" class="badge ${isAllowed ? 'badge-primary' : ''}" data-preview-toggle="${index}">${isAllowed ? '是' : '否'}</button>`)}</td>
              <td class="nowrap preview-col-downloaded" data-label="已下载">${item.downloaded ? '<span class="badge badge-success">是</span>' : '<span class="badge">否</span>'}</td>
              <td class="nowrap preview-col-subgroup" data-label="字幕组"><span class="rss-role-badge ${item.source_role === 'main' ? 'main' : 'backup'}">${item.source_role === 'main' ? '主' : '备'}</span> ${esc(item.subgroup || '未知字幕组')}</td>
              <td class="name-cell" data-label="标题"><div class="preview-mobile-meta"><span class="${item.filtered ? 'meta-warn' : (isAllowed ? 'meta-ok' : 'meta-muted')}">下载 ${item.filtered ? '过滤' : (isAllowed ? '是' : '否')}</span><span class="${item.downloaded ? 'meta-ok' : 'meta-muted'}">已下载 ${item.downloaded ? '是' : '否'}</span><span class="${item.source_role === 'main' ? 'meta-primary' : 'meta-info'}">${item.source_role === 'main' ? '主' : '备'} ${esc(item.subgroup || '未知字幕组')}</span></div><div class="title">${esc(item.title)}</div><div class="preview-mobile-rename">${item.filtered ? esc(item.filter_reason || '被当前规则过滤') : (item.renamed ? esc(item.renamed) : '无法解析')}</div></td>
              <td class="preview-col-rename" data-label="重命名">${item.filtered ? `<span class="hint">${esc(item.filter_reason || '被当前规则过滤')}</span>` : (item.renamed ? `<code class="preview-renamed">${esc(item.renamed)}</code>` : '<span class="form-error">无法解析</span>')}</td>
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
            ${items.length ? `<div class="table-wrap"><table class="table rename-preview-table item-preview-table">
              <colgroup><col style="width:34px"><col style="width:64px"><col style="width:62px"><col style="width:150px"><col><col style="width:240px"></colgroup>
              <thead><tr><th>选</th><th>下载</th><th>已下载</th><th>字幕组</th><th>标题</th><th>重命名</th></tr></thead>
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
      if (groups.has(group.id)) excludes.push(group.group ? `{{${group.group}}}:${group.pattern}` : group.pattern);
    }
    return { filter: includes.join('\n'), exclude: excludes.join('\n') };
  };

  const initialContext = getPreviewContext ? getPreviewContext() : {};
  let matchTemplates = [];
  try { matchTemplates = await loadMatchTemplates(); } catch (err) { matchTemplates = []; }
  openModal({
    title: '规则与匹配',
    subtitle: (String(initialContext.name || '').trim() ? String(initialContext.name).trim() + ' · ' : '') + '选择字幕组模板，或使用常用规则临时组合；实时显示匹配结果',
    size: 'rules',
    body: `
      <div class="rules-template-panel"><div class="rules-template-head"><strong>字幕组匹配模板</strong><span id="rules-template-subgroup"></span><span class="hint">选择「全部」表示不限制</span></div><div id="rules-templates" class="rules-templates"></div></div>
      <div class="rules-active-summary" id="rules-active-summary"></div>
      <div class="rules-quick-panel">
        <details class="rules-quick-fold">
          <summary><strong>常用规则</strong><span class="hint">点击展开选择</span></summary>
          <div class="rules-quick-head">
            <span class="hint">点击追加到</span>
            <div class="rules-quick-target" id="rules-quick-target">
              <button type="button" class="active" data-quick-target="filter">匹配</button>
              <button type="button" data-quick-target="exclude">排除</button>
            </div>
          </div>
          <div class="rules-quick-chips" id="rules-quick-chips"></div>
        </details>
        <div class="rule-bubble-columns">
          <div class="rule-bubble-column">
            <div class="rule-bubble-head"><span>匹配</span><small>全部命中</small></div>
            <div class="rule-bubble-list" id="rules-filter-bubbles"></div>
            <div class="rule-bubble-add"><input class="input mono" id="rules-filter-input" placeholder="输入正则，回车添加"><button class="btn btn-sm" type="button" data-bubble-add="filter">添加</button></div>
          </div>
          <div class="rule-bubble-column">
            <div class="rule-bubble-head"><span>排除</span><small>任一命中即过滤</small></div>
            <div class="rule-bubble-list" id="rules-exclude-bubbles"></div>
            <div class="rule-bubble-add"><input class="input mono" id="rules-exclude-input" placeholder="输入正则，回车添加"><button class="btn btn-sm" type="button" data-bubble-add="exclude">添加</button></div>
          </div>
        </div>
        <textarea class="hidden" id="rules-base-filter">${esc(initialContext.filter || '')}</textarea>
        <textarea class="hidden" id="rules-base-exclude">${esc(initialContext.exclude || '')}</textarea>
      </div>
      <div class="hidden" aria-hidden="true">
        <div id="rules-body"></div>
        <div class="rule-preview"><span>生成规则</span><code id="rules-filter-preview">—</code><code id="rules-exclude-preview">—</code></div>
        <button class="btn btn-sm" type="button" id="rules-clear">清空快捷预设</button>
      </div>
      <div class="rule-results" id="rules-results"><span class="hint">选择模板后自动显示 RSS 匹配结果</span></div>`,
    footer: `<span class="spacer"></span><button class="btn" data-action="modal-close">取消</button><button class="btn btn-primary" id="rules-apply">${icon('check')}应用到订阅</button>`,
    onMount(el, close) {
      const body = el.querySelector('#rules-body');
      const filterPreview = el.querySelector('#rules-filter-preview');
      const excludePreview = el.querySelector('#rules-exclude-preview');
      const resultsBox = el.querySelector('#rules-results');
      let resultTimer = null;
      let resultToken = 0;
      let autoApplied = false;
      const templateSubgroup = String(initialContext.subgroup_name || '').trim();
      const baseFilterEl = el.querySelector('#rules-base-filter');
      const baseExcludeEl = el.querySelector('#rules-base-exclude');
      const templatesBox = el.querySelector('#rules-templates');
      const activeSummary = el.querySelector('#rules-active-summary');
      const usableTemplates = templateSubgroup
        ? matchTemplates.filter(template => templateMatchesSubgroup(template, templateSubgroup))
        : [];
      let activeTemplate = '';
      const templateFilter = template => scopeTemplateRules(template.filter, template.subgroup || templateSubgroup);
      const templateExclude = template => scopeTemplateRules(template.exclude, template.subgroup || templateSubgroup);
      const bubbleEls = {
        filter: el.querySelector('#rules-filter-bubbles'),
        exclude: el.querySelector('#rules-exclude-bubbles'),
      };
      const bubbleInputs = {
        filter: el.querySelector('#rules-filter-input'),
        exclude: el.querySelector('#rules-exclude-input'),
      };
      const bubbleRules = target => String((target === 'filter' ? baseFilterEl : baseExcludeEl).value || '').split(/\r?\n/).map(line => line.trim()).filter(Boolean);
      const setBubbleRules = (target, lines) => {
        (target === 'filter' ? baseFilterEl : baseExcludeEl).value = Array.from(new Set(lines.filter(Boolean))).join('\n');
      };
      const refreshBubblePreview = () => {
        autoApplied = true;
        activeTemplate = '';
        for (const key of Object.keys(selected)) selected[key] = '';
        groups.clear();
        renderTemplates();
        draw();
        scheduleResults();
      };
      function renderBubbles() {
        for (const target of ['filter', 'exclude']) {
          const box = bubbleEls[target];
          if (!box) continue;
          const lines = bubbleRules(target);
          if (!lines.length) {
            box.innerHTML = `<span class="rule-bubble-empty">${target === 'filter' ? '全部' : '无'}</span>`;
            continue;
          }
          box.innerHTML = lines.map((line, index) => `<span class="rule-bubble" title="${esc(line)}"><code>${esc(line)}</code><button type="button" data-bubble-remove="${target}" data-bubble-index="${index}" aria-label="删除">×</button></span>`).join('');
          box.querySelectorAll('[data-bubble-remove]').forEach(button => button.addEventListener('click', () => {
            const list = bubbleRules(button.dataset.bubbleRemove);
            list.splice(Number(button.dataset.bubbleIndex), 1);
            setBubbleRules(button.dataset.bubbleRemove, list);
            renderBubbles();
            refreshBubblePreview();
          }));
        }
      }
      const addBubble = (target, raw) => {
        const value = String(raw || '').trim();
        if (!value) return;
        const lines = bubbleRules(target);
        lines.push(value);
        setBubbleRules(target, lines);
        renderBubbles();
        refreshBubblePreview();
      };
      for (const target of ['filter', 'exclude']) {
        const input = bubbleInputs[target];
        const button = el.querySelector(`[data-bubble-add="${target}"]`);
        const submit = () => { addBubble(target, input.value); input.value = ''; };
        button?.addEventListener('click', submit);
        input?.addEventListener('keydown', event => {
          if (event.key === 'Enter') { event.preventDefault(); submit(); }
        });
      }

      const renderTemplates = () => {
        if (!templatesBox) return;
        const subgroupLabel = el.querySelector('#rules-template-subgroup');
        const feedLabel = String(initialContext.name || '').trim();
        if (subgroupLabel) subgroupLabel.textContent = feedLabel ? '当前订阅：' + feedLabel : '当前订阅未设置';
        const option = (id, label, template) => {
          const tags = template ? ((template.tags || []).length ? template.tags : String(template.name || '').split(/[\s/]+/).filter(Boolean)) : [];
          return `<button type="button" class="match-option ${activeTemplate === id ? 'active' : ''}" data-rule-template="${esc(id)}">
            <span class="match-radio"></span>
            <span class="match-option-main">${template ? `<span class="match-option-name">${esc(template.name || '未命名模板')}</span><span class="match-tags">${tags.map(tag => `<span>${esc(tag)}</span>`).join('')}</span>` : `<span class="match-option-name">${esc(label)}</span>`}</span>
          </button>`;
        };
        const cards = usableTemplates.map(template => option(template.id, template.name, template));
        cards.push(option('__all', '全部', null));
        cards.push(option('__auto', '自动应用最佳规则', null));
        if (!usableTemplates.length) cards.push('<p class="hint" style="margin:6px 0 0">当前字幕组没有匹配模板，可直接在下方手动填写规则。</p>');
        templatesBox.innerHTML = cards.join('');
        if (activeSummary) {
          const template = usableTemplates.find(item => item.id === activeTemplate);
          const label = template ? template.name : (activeTemplate === '__all' ? '全部' : (activeTemplate === '__auto' ? '自动应用最佳规则' : '自定义规则'));
          const filterLines = String(baseFilterEl?.value || '').split(/\r?\n/).filter(line => line.trim()).length;
          const excludeLines = String(baseExcludeEl?.value || '').split(/\r?\n/).filter(line => line.trim()).length;
          activeSummary.innerHTML = `<span>当前</span><strong>${esc(label)}</strong><span>匹配 ${filterLines} 条 · 排除 ${excludeLines} 条</span>`;
        }
        templatesBox.querySelectorAll('[data-rule-template]').forEach(button => button.addEventListener('click', () => {
          const id = button.dataset.ruleTemplate;
          if (id === '__auto') {
            activeTemplate = '__auto';
            applyBestRules(false, true);
            return;
          }
          activeTemplate = id;
          autoApplied = true;
          for (const key of Object.keys(selected)) selected[key] = '';
          groups.clear();
          if (id === '__all') {
            baseFilterEl.value = '';
            baseExcludeEl.value = '';
          } else {
            const template = usableTemplates.find(item => item.id === id);
            if (template) {
              baseFilterEl.value = templateFilter(template);
              baseExcludeEl.value = templateExclude(template);
            }
          }
          renderTemplates();
          renderBubbles();
          draw();
          scheduleResults();
        }));
      };
      const currentFilter = (baseFilterEl?.value || '').trim();
      const currentExclude = (baseExcludeEl?.value || '').trim();
      const matchedTemplate = usableTemplates.find(template => currentFilter === templateFilter(template).trim() && currentExclude === templateExclude(template).trim());
      activeTemplate = matchedTemplate ? matchedTemplate.id : (currentFilter || currentExclude ? '__custom' : '__all');
      renderTemplates();
      renderBubbles();
      [baseFilterEl, baseExcludeEl].forEach(input => input?.addEventListener('input', () => {
        activeTemplate = '';
        autoApplied = true;
        for (const key of Object.keys(selected)) selected[key] = '';
        groups.clear();
        renderTemplates();
      }));




      const QUICK_RULE_CHIPS = [
        ['1080P', '1080[Pp]'], ['720P', '720[Pp]'], ['4K', '2160[Pp]|4K'],
        ['简', '简'], ['繁', '繁'], ['简日', '简日'], ['繁日', '繁日'], ['简繁日', '简繁日'],
        ['CHS', '(?i)CHS'], ['CHT', '(?i)CHT'], ['GB', '(?i)GB'], ['BIG5', '(?i)BIG5'], ['简中', '简中'], ['繁中', '繁中'],
        ['内封', '内封'], ['内嵌', '内嵌'], ['外挂', '外挂'],
        ['MKV', '(?i)MKV'], ['MP4', '(?i)MP4'],
        ['x.5', '\\d{1,4}\\.5'], ['合集', '合集'], ['特别篇', '特别篇'],
      ];
      let quickTarget = 'filter';
      const quickBox = el.querySelector('#rules-quick-chips');
      const targetBox = el.querySelector('#rules-quick-target');
      const renderQuick = () => {
        if (quickBox) {
          quickBox.innerHTML = QUICK_RULE_CHIPS.map(([label, value]) => `<button type="button" class="rule-chip" data-quick-rule="${esc(value)}">${esc(label)}</button>`).join('');
          quickBox.querySelectorAll('[data-quick-rule]').forEach(button => button.addEventListener('click', () => {
            const target = quickTarget === 'exclude' ? baseExcludeEl : baseFilterEl;
            const value = button.dataset.quickRule;
            const lines = target.value.split(/\r?\n/).map(line => line.trim()).filter(Boolean);
            if (!lines.includes(value)) lines.push(value);
            target.value = lines.join('\n');
            autoApplied = true;
            activeTemplate = '';
            for (const key of Object.keys(selected)) selected[key] = '';
            groups.clear();
            renderTemplates();
            renderBubbles();
            draw();
            scheduleResults();
          }));
        }
        if (targetBox) {
          targetBox.querySelectorAll('[data-quick-target]').forEach(button => button.classList.toggle('active', button.dataset.quickTarget === quickTarget));
        }
      };
      targetBox?.querySelectorAll('[data-quick-target]').forEach(button => button.addEventListener('click', () => {
        quickTarget = button.dataset.quickTarget;
        renderQuick();
      }));
      renderQuick();


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
            subgroup: context.subgroup_name || '',
            filter: effectiveFilter, exclude: effectiveExclude,
            rename_season: String(context.rename_season || 0),
            rename_offset: String(context.rename_offset || 0),
            tmdb_id: context.tmdb_id || '', tmdb_type: context.tmdb_type || 'tv',
          });
          const data = await api(`/api/rename/preview?${params}`);
          if (token !== resultToken) return;
          renderRuleResults(resultsBox, data);
          if (data.total > 0 && Number(data.matched) === 0) {
            resultsBox.insertAdjacentHTML('afterbegin', `<div class="notice warn" style="margin:6px 0 8px">${icon('alert')}<span>当前规则没有匹配任何条目。请检查语言/字幕/内封等标签是否存在于实际标题，或检查排除规则是否命中了全部条目。</span></div>`);
          }
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
      const templateFitsTitles = (template, titles) => {
        const lines = String(template.filter || '').split(/\r?\n/).map(line => line.trim()).filter(Boolean);
        if (!lines.length) return true;
        const matchesLine = (line, title) => {
          let pattern = line.replace(/^\{\{.+?\}\}:/, '');
          let flags = '';
          if (pattern.startsWith('(?i)')) { flags = 'i'; pattern = pattern.slice(4); }
          try { return new RegExp(pattern, flags).test(title); } catch (_) { return false; }
        };
        return lines.every(line => titles.some(title => matchesLine(line, title)));
      };
      const applyBestRules = (silent = false) => {
        if (!resultsBox._ruleData) {
          toast('正在读取 RSS，请稍候再试', 'error');
          return;
        }
        autoApplied = true;
        const titles = (resultsBox._ruleData.items || []).map(item => item.title);
        const autoCustomLines = () => Object.values((state.autoPriority && state.autoPriority.custom) || {}).flat().map(line => String(line || '').trim()).filter(Boolean).filter(line => {
          let pattern = line.replace(/^\{\{.+?\}\}:/, '');
          if (pattern.startsWith('(?i)')) pattern = pattern.slice(4);
          try { return titles.some(title => new RegExp(pattern, 'i').test(title)); } catch (_) { return false; }
        });
        const existingFilter = String(baseFilterEl?.value || '').trim();
        const existingExclude = String(baseExcludeEl?.value || '').trim();
        if (existingFilter || existingExclude) {
          activeTemplate = '__auto';
          renderTemplates();
          renderBubbles();
          if (!silent) toast('已保留字幕组已有规则', 'success');
          return;
        }
        const autoTemplate = usableTemplates.find(template => templateFitsTitles(template, titles));
        if (autoTemplate) {
          activeTemplate = '__auto';
          baseFilterEl.value = [templateFilter(autoTemplate), ...autoCustomLines()].filter(Boolean).join('\n');
          baseExcludeEl.value = templateExclude(autoTemplate);
          for (const key of Object.keys(selected)) selected[key] = '';
          groups.clear();
          renderTemplates();
          renderBubbles();
          draw();
          if (!silent) toast('已自动应用规则，匹配/排除气泡已更新', 'success');
          return;
        }
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
        const autoRules = build();
        const customLines = autoCustomLines();
        baseFilterEl.value = [autoRules.filter, ...customLines].filter(Boolean).join('\n');
        baseExcludeEl.value = autoRules.exclude;
        for (const key of Object.keys(selected)) selected[key] = '';
        groups.clear();
        activeTemplate = '__auto';
        renderTemplates();
        renderBubbles();
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
    subgroup_name: prefill?.subgroup_name || prefill?.subgroup || '',
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
    auto_backup_idle: !!((state.config.lab || {}).auto_backup_idle),
    auto_backup_failures: !!((state.config.lab || {}).auto_backup_subgroup),
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
    title: editing ? '修改订阅' : '添加订阅',
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
        <button type="button" data-feed-tab="custom">高级</button>
      </div>
      <form id="feed-form" class="feed-form" autocomplete="off">
        <section class="feed-tab-panel active" data-feed-panel="basic">
          ${editing ? (() => { const p = feedProgress(feed); const t = p.total || '*'; const pct = v => p.total > 0 ? Math.min(100, Math.round(v / p.total * 100)) : 0; return `<div class="feed-progress-summary">
            <div class="feed-progress-card"><div><span>主 RSS 进度</span><strong>${p.mainDone} / ${t}</strong></div><div class="feed-progress-bar"><i style="width:${pct(p.mainDone)}%"></i></div></div>
            <div class="feed-progress-card"><div><span>主 + 备用 RSS</span><strong>${p.allDone} / ${t}</strong></div><div class="feed-progress-bar info"><i style="width:${pct(p.allDone)}%"></i></div></div>
          </div>`; })() : ''}
          <div class="feed-form-section">
            <div class="feed-form-section-head">基本信息</div>
            <div class="form-grid">
              <label class="field span-2"><span>标题 *</span><input class="input" id="feed-name" value="${esc(feed.name)}" placeholder="例如：药屋少女的呢喃"></label>
              <label class="field span-2 feed-tmdb-line"><span>TMDB ID</span>
                <div class="feed-tmdb-row">
                  <input class="input mono" id="feed-tmdb-id" value="${esc(feed.tmdb_id || '')}" placeholder="输入 ID">
                  <button class="btn btn-sm btn-icon-only" type="button" id="feed-tmdb-auto" title="自动匹配 TMDB ID" aria-label="自动匹配 TMDB ID">${icon('refresh')}</button>
                  <button class="btn btn-sm btn-icon-only" type="button" id="feed-tmdb-search" title="搜索 TMDB" aria-label="搜索 TMDB">${icon('search')}</button>
                  <select class="select select-compact" id="feed-tmdb-type"><option value="tv" ${(feed.tmdb_type || 'tv') === 'tv' ? 'selected' : ''}>TV</option><option value="movie" ${feed.tmdb_type === 'movie' ? 'selected' : ''}>Movie</option></select>
                </div>
              </label>
              <label class="field span-2"><span>BGM 链接</span><input class="input mono" id="feed-bgm-url" value="${esc(initialBgm)}" placeholder="https://bgm.tv/subject/665390"></label>
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">RSS 来源</div>
            <div class="feed-rss-source-row">
              <input class="input feed-rss-subgroup" id="feed-subgroup-name" value="${esc(mainSubgroup)}" placeholder="字幕组名称，例如 LoliHouse">
              <button class="btn btn-sm" type="button" id="feed-subgroup-detect">${icon('refresh')}获取</button>
              <input class="input mono feed-rss-url" id="feed-url" value="${esc(feed.url)}" placeholder="https://...">
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">季数与集数</div>
            <div class="form-grid feed-season-grid">
              <label class="field"><span>季数</span>
                <div class="field-row"><input class="input input-compact" id="feed-rename-season" type="number" min="0" max="1000" placeholder="季数" value="${Number(feed.rename_season || 0)}"><button class="btn btn-sm btn-icon-only" type="button" id="feed-season-suggest" title="自动识别季数">${icon('refresh')}</button></div>
                <span class="hint">用于重命名，如 S01</span>
              </label>
              <label class="field"><span>集数偏移</span><input class="input input-compact" id="feed-rename-offset" type="number" min="-1000" max="1000" value="${Number(feed.rename_offset || 0)}"><span class="hint">RSS 集数与实际相差几集</span></label>
              <label class="field"><span>总集数</span><div class="field-row"><input class="input input-compact" id="feed-total-episodes" type="number" min="0" max="10000" placeholder="总集数" value="${Number(feed.total_episodes || 0)}"><button class="btn btn-sm btn-icon-only" type="button" id="feed-total-fetch" title="从 TMDB 自动获取总集数">${icon('refresh')}</button></div><span class="hint">0 表示未设置；已关联 TMDB 时可自动获取</span></label>
              <label class="field"><span>播出季度</span><div class="feed-quarter-editor"><input class="input input-compact" id="feed-quarter-year" type="number" min="2000" max="2100" placeholder="年份" value="${quarterOf(feed) ? quarterOf(feed).slice(0,4) : new Date().getFullYear()}"><select class="select" id="feed-quarter-month">${["01","04","07","10"].map(m => `<option value="${m}" ${((quarterOf(feed) ? quarterOf(feed).slice(5) : String(Math.floor(new Date().getMonth()/3) === 0 ? "01" : Math.floor(new Date().getMonth()/3) === 1 ? "04" : Math.floor(new Date().getMonth()/3) === 2 ? "07" : "10")) === m) ? "selected" : ""}>${m}</option>`).join("")}</select></div><span class="hint">仅用于显示和筛选，不参与重命名</span></label>
            </div>
          </div>
          <div class="feed-form-section feed-enabled-section">
            <div class="feed-form-section-head">订阅状态</div>
            <label class="field span-2"><span>启用订阅</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="feed-enabled" ${feed.enabled !== false ? 'checked' : ''}><span></span></span><span>关闭后不会参与自动运行</span></span></label>
            <label class="field span-2"><span>主 RSS 完成后自动停用</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="feed-auto-disable" ${feed.auto_disable ? 'checked' : ''}><span></span></span><span>主 RSS 下载进度达到总集数后自动关闭启用</span></span></label>
          </div>
          ${editing ? `<div class="feed-edit-actions">
            <button class="btn" type="button" id="feed-action-toggle">${feed.enabled === false ? '启用订阅' : '禁用订阅'}</button>
            <button class="btn" type="button" id="feed-action-reset">重置为未完成</button>
            <span class="spacer"></span>
            <button class="btn btn-danger" type="button" id="feed-action-delete">${icon('trash')}删除订阅</button>
          </div>` : ''}
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
          <div class="feed-form-section">
            <div class="feed-form-section-head">保存位置</div>
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
            <label class="field"><span>缓存过期（秒）</span><input class="input input-compact" id="feed-expiration" type="number" min="0" step="1" value="${feed.expiration || 0}"><span class="hint">0 使用全局缓存策略。</span></label>
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">备用 RSS 策略</div>
            <div class="form-grid">
              <label class="field span-2"><span>启用自动加入备用 RSS</span><span class="check" style="height:32px"><span class="switch"><input type="checkbox" id="feed-auto-backup" ${feed.auto_backup !== false ? 'checked' : ''}><span></span></span><span>仅对该订阅生效</span><button class="help-dot" type="button" id="feed-auto-backup-help" title="查看说明">?</button></span></label>
              <label class="field span-2"><span>连续失败后自动加入备用 RSS</span><span class="check" style="height:32px"><span class="switch"><input type="checkbox" id="feed-auto-backup-failures" ${feed.auto_backup_failures !== false ? 'checked' : ''}><span></span></span><span>主 RSS 连续添加失败达到全局次数后触发</span></span></label>
              <label class="field span-2"><span>摸鱼后自动加入备用 RSS</span><span class="check" style="height:32px"><span class="switch"><input type="checkbox" id="feed-auto-backup-idle" ${feed.auto_backup_idle ? 'checked' : ''}><span></span></span><span>超过摸鱼提醒天数后触发；全局「摸鱼后自动补充」也需开启</span></span></label>

              <label class="field span-2"><span>自动备用 RSS 使用 TMDB 集数偏移</span><span class="check" style="height:32px"><span class="switch"><input type="checkbox" id="feed-auto-offset" ${feed.auto_offset ? 'checked' : ''}><span></span></span><span>TMDB 自动校正</span><button class="help-dot" type="button" id="feed-auto-offset-help" title="查看说明">?</button></span></label>
            </div>
          </div>
          <div class="feed-form-section">
            <div class="feed-form-section-head">规则与匹配</div>
            <div class="feed-rule-summary" id="feed-rule-summary"></div>
            <div class="feed-action-row">
              <button class="btn btn-sm" type="button" id="feed-rules">${icon('filter')}配置规则与匹配</button>
              <button class="btn btn-sm" type="button" id="feed-item-preview">${icon('eye')}预览条目</button>
              <button class="btn btn-sm" type="button" id="feed-preview">${icon('zap')}测试订阅</button>
            </div>
            <p class="feed-action-hint">配置规则与匹配：选择字幕组模板，或手动编辑匹配/排除规则并实时预览；预览条目：合并主 RSS 与全部备用 RSS 按集数排序；测试订阅：抓取并应用规则显示匹配结果。</p>
            <div id="feed-preview-box"></div>
          </div>
        </section>
      </form>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="feed-save">${icon('save')}${editing ? '保存修改' : '添加订阅'}</button>`,
    onMount(el, close) {
      const form = $('#feed-form', el);
      if (form) form.addEventListener('submit', event => event.preventDefault());
      const updateRuleSummary = () => {
        const box = $('#feed-rule-summary', el);
        if (!box) return;
        const filter = String(local.filter || '').trim();
        const exclude = String(local.exclude || '').trim();
        box.innerHTML = `<div><span>匹配</span><code>${esc(filter || '全部')}</code></div><div><span>排除</span><code>${esc(exclude || '无')}</code></div><button class="btn btn-sm btn-ghost" type="button" id="feed-rules-summary">调整</button>`;
        box.querySelector('#feed-rules-summary')?.addEventListener('click', () => $('#feed-rules', el).click());
      };
      updateRuleSummary();
      if (editing) {
        $('#feed-action-toggle', el).addEventListener('click', () => { close(); toggleFeed(index); });
        $('#feed-action-delete', el).addEventListener('click', () => { close(); deleteFeed(index); });
        $('#feed-action-reset', el).addEventListener('click', () => {
          close();
          const modal = openModal({
            title: '重置订阅完成状态',
            subtitle: '重新启用该订阅',
            body: `<div class="notice">${icon('info')}<span>将订阅「${esc(feed.name || feed.url)}」恢复为未完成并启用。</span></div>
              <label class="check" style="margin-top:12px"><input type="checkbox" id="reset-clear-history" checked>同时清除该订阅的下载历史</label>
              <div class="hint" style="margin-top:6px">只有清除历史，才能避免达到总集数后再次被自动停用。</div>`,
            footer: '<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="reset-confirm">重置并启用</button>',
          });
          modal.el.querySelector('#reset-confirm').addEventListener('click', async () => {
            const clearHistory = !!modal.el.querySelector('#reset-clear-history')?.checked;
            modal.close();
            try {
              const r = await api('/api/feeds/reset', { method: 'POST', body: { name: feed.name, url: feed.url, clear_history: clearHistory } });
              toast(r.deleted_history ? `已重置并清除 ${r.deleted_history} 条历史` : '已重置为未完成并启用', 'success');
              await loadSettingsPage();
              if (state.route === 'feeds') loadFeedsPage();
            } catch (e) { toast(e.message, 'error', 5000); }
          });
        });
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
          return `<div class="feed-subgroup-row feed-subgroup-card" data-subgroup-row data-index="${index}">
            <div class="feed-subgroup-top">
              <span class="rss-role-badge backup">备用</span>
              <input class="input feed-subgroup-name" data-sg="name" value="${esc(subgroupLabel(group.name, feed.name))}" placeholder="字幕组名称">
              <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-subgroup-detect title="根据 RSS 标题识别字幕组">${icon('refresh')}</button>
              <label class="check feed-subgroup-enabled"><span class="switch"><input type="checkbox" data-sg="enabled" ${enabled ? 'checked' : ''}><span></span></span>启用</label>
              <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-remove-subgroup title="删除">${icon('trash')}</button>
            </div>
            <input class="input mono feed-subgroup-url" data-sg="url" value="${esc(group.url || '')}" placeholder="备用 RSS 地址">
            <div class="feed-subgroup-bottom">
              <div class="feed-subgroup-meta">
                <label>季 <input class="feed-subgroup-num" data-sg="season" type="number" min="0" max="1000" value="${Number(group.season || 0)}"></label>
                <label>偏移 <input class="feed-subgroup-num" data-sg="offset" type="number" min="-1000" max="1000" value="${Number(group.offset || 0)}"></label>
              </div>
              <div class="feed-subgroup-actions">
                <button class="btn btn-sm btn-ghost" type="button" data-set-main title="设为第一个 RSS（主 RSS）">设主</button>
                <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-subgroup-up title="上移" ${index === 0 ? 'disabled' : ''}>↑</button>
                <button class="btn btn-sm btn-ghost btn-icon-only" type="button" data-subgroup-down title="下移" ${index === total - 1 ? 'disabled' : ''}>↓</button>
                <button class="btn btn-sm btn-ghost ${group.filter || group.exclude ? 'has-rules' : ''}" type="button" data-subgroup-rules title="${group.filter || group.exclude ? '已设置独立规则' : '未设置规则，不过滤'}">规则</button>
              </div>
            </div>
            <input type="hidden" data-sg="filter" value="${esc(group.filter || '')}">
            <input type="hidden" data-sg="exclude" value="${esc(group.exclude || '')}">
            <input type="hidden" data-sg="start_episode" value="${Number(group.start_episode || 0)}">
          </div>`;
        };

        const currentMain = () => ({
          name: $('#feed-subgroup-name', el).value.trim(),
          url: $('#feed-url', el).value.trim(),
          season: Number($('#feed-rename-season', el).value || 0),
          offset: Number($('#feed-rename-offset', el).value || 0),
          filter: local.filter || '',
          exclude: local.exclude || '',
          start_episode: 0,
          enabled: true,
        });
        const renderSubgroups = () => {
          const main = currentMain();
          const mainRow = `<div class="feed-main-source-card" data-main-source-row>
            <div class="feed-main-source-top"><span class="badge badge-primary role-badge">主 RSS</span><strong>${esc(main.name || '未命名字幕组')}</strong><button class="btn btn-sm btn-ghost" type="button" data-edit-main>编辑</button></div>
            <code class="feed-main-source-url" title="${esc(main.url)}">${esc(main.url)}</code>
            <div class="feed-main-source-meta"><span>季 ${main.season || 0}</span><span>偏移 ${main.offset || 0}</span>${main.filter || main.exclude ? '<span class="badge badge-info">独立规则</span>' : ''}</div>
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
          subgroupBox.querySelectorAll('[data-subgroup-detect]').forEach(button => {
            button.addEventListener('click', () => {
              const row = button.closest('[data-subgroup-row]');
              const url = row.querySelector('[data-sg="url"]').value.trim();
              detectSubgroup(button, url, name => { row.querySelector('[data-sg="name"]').value = name; });
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
                name: feed.name || '',
                subgroup_name: group.name || '',
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
          start_episode: Number(row.querySelector('[data-sg="start_episode"]').value || 0),
        });
        refreshBackupSources = renderSubgroups;
        renderSubgroups();
        $('#feed-add-subgroup', el)?.addEventListener('click', () => {
          local.subgroups = [...subgroupBox.querySelectorAll('[data-subgroup-row]')].map(readSubgroupRow);
          const source = sourceInfo(feed.source, feed.url).key;
          openBrowseModal(source, pick => {
            const group = {
              name: pick.subgroup || pick.title || '备用字幕组',
              url: pick.url,
              season: Number(feed.rename_season || 0),
              offset: Number(feed.rename_offset || 0),
              start_episode: 0,
              enabled: true,
              filter: '',
              exclude: '',
            };
            // Only attach the backup RSS after the rules page is confirmed.
            openRulesModal((filter, exclude) => {
              group.filter = filter;
              group.exclude = exclude;
              local.subgroups.push(group);
              renderSubgroups();
            }, () => ({
              name: feed.name || '',
              subgroup_name: group.name,
              url: group.url,
              season: group.season,
              offset: group.offset,
              filter: '',
              exclude: '',
            }));
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
      const subgroupInput = $('#feed-subgroup-name', el);
      let subgroupResolveTimer = null;
      let subgroupResolveToken = 0;
      let subgroupAutoResolved = false;
      let subgroupTouched = false;
      const scheduleSubgroupResolve = () => {
        if (editing) return;
        clearTimeout(subgroupResolveTimer);
        subgroupResolveTimer = setTimeout(async () => {
          const rawURL = String(urlInput.value || '').trim();
          if (!rawURL || subgroupTouched) return;
          const current = String(subgroupInput.value || '').trim();
          const slug = urlSubgroupName(rawURL);
          if (current && current.toLowerCase() !== String(slug || '').toLowerCase() && !subgroupAutoResolved) return;
          const token = ++subgroupResolveToken;
          try {
            const data = await api('/api/feeds/subgroup-resolve?url=' + encodeURIComponent(rawURL));
            if (token !== subgroupResolveToken || !data.subgroup) return;
            subgroupInput.value = data.subgroup;
            local.subgroupName = data.subgroup;
            subgroupAutoResolved = true;
          } catch (_) { /* keep the URL-derived label when the RSS cannot be read */ }
        }, 550);
      };
      subgroupInput?.addEventListener('input', () => { subgroupTouched = true; subgroupAutoResolved = false; });
      urlInput.addEventListener('input', () => {
        if (subgroupAutoResolved) { subgroupInput.value = ''; subgroupAutoResolved = false; }
        subgroupTouched = false;
        scheduleSubgroupResolve();
      });
      setTimeout(scheduleSubgroupResolve, 300);

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
              openNewFeedWithRules({
                source: pick.source,
                name: pick.title,
                subgroup_name: pick.subgroup || '',
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
          updateRuleSummary();
          markDirty();
        }, () => readFeedForm(el, local));
      });
      const detectSubgroup = async (button, rawURL, apply) => {
        const url = String(rawURL || '').trim();
        if (!url) { toast('请先填写 RSS 地址', 'error'); return; }
        setLoading(button, true, '识别中');
        try {
          const data = await api('/api/feeds/subgroup-resolve?url=' + encodeURIComponent(url));
          if (!data.subgroup) { toast('没有从 RSS 中识别到字幕组', 'error'); return; }
          apply(data.subgroup);
          toast('已识别字幕组：' + data.subgroup, 'success');
        } catch (err) {
          toast('识别失败：' + err.message, 'error', 5200);
        } finally {
          setLoading(button, false);
        }
      };
      $('#feed-subgroup-detect')?.addEventListener('click', () => {
        detectSubgroup($('#feed-subgroup-detect', el), urlInput.value, name => {
          $('#feed-subgroup-name', el).value = name;
          local.subgroupName = name;
        });
      });
      $('#feed-auto-backup-help', el)?.addEventListener('click', () => openModal({
        title: '自动加入备用 RSS 说明',
        size: 'narrow',
        subtitle: '仅对该订阅生效',
        body: `<div class="notice">${icon('info')}<span>满足任一条件后会检查该订阅的缺失集和喜好字幕组，并尝试添加备用 RSS：主 RSS 连续添加失败达到设定次数；或开启「摸鱼后自动补充」且超过摸鱼提醒天数。全局「自动添加备用字幕组」也必须开启；关闭本开关可单独跳过该订阅。</span></div>`,
        footer: '<button class="btn btn-primary" data-action="modal-close">知道了</button>',
      }));
      $('#feed-auto-offset-help', el)?.addEventListener('click', () => openModal({
        title: '自动备用 RSS 偏移说明',
        size: 'narrow',
        subtitle: '默认关闭，按需开启',
        body: `<div class="notice">${icon('info')}<span>若番剧存在分季放送，例如 BGM 第一季 1-13、第二季 14-26，而 TMDB 统一为第一季 1-26，开启该功能后，系统会用 TMDB 分集和发布日期计算备用 RSS 的季数与偏移，提高准确度。</span></div><div class="hint" style="margin-top:10px">需要填写 TMDB ID 和 TMDB API Key；匹配样本不足时会保留 Offset 0。</div>`,
        footer: '<button class="btn btn-primary" data-action="modal-close">知道了</button>',
      }));
      $('#feed-item-preview', el).addEventListener('click', () => openItemPreview(readFeedForm(el, local)));
      const suggestSeason = async () => {
        const seasonInput = $('#feed-rename-season', el);
        if (!seasonInput) return;
        if (Number(seasonInput.value) > 0 && seasonInput.dataset.autofilled !== '1') return;
        try {
          const r = await api('/api/feeds/season-suggest', { method: 'POST', body: {
            name: ($('#feed-name', el) || {}).value || '',
            url: ($('#feed-url', el) || {}).value || '',
            tmdb_id: ($('#feed-tmdb-id', el) || {}).value || '',
          }});
          if (r.season > 0 && (Number(seasonInput.value) === 0 || seasonInput.dataset.autofilled === '1')) {
            seasonInput.value = r.season; seasonInput.dataset.autofilled = '1'; markDirty();
            toast('已识别季数：第 ' + r.season + ' 季', 'success');
          }
        } catch { /* 自动识别失败不阻断添加 */ }
      };
      let seasonTimer = null;
      const autoSeason = !editing;
      const scheduleSeason = () => {
        if (!autoSeason) return;
        clearTimeout(seasonTimer);
        seasonTimer = setTimeout(suggestSeason, 500);
      };
      $('#feed-url', el)?.addEventListener('change', scheduleSeason);
      $('#feed-name', el)?.addEventListener('change', scheduleSeason);
      $('#feed-tmdb-id', el)?.addEventListener('change', scheduleSeason);
      $('#feed-season-suggest', el)?.addEventListener('click', suggestSeason);
      $('#feed-total-fetch', el)?.addEventListener('click', async () => {
        const button = $('#feed-total-fetch', el);
        const tmdbID = ($('#feed-tmdb-id', el) || {}).value?.trim() || '';
        const type = ($('#feed-tmdb-type', el) || {}).value || 'tv';
        if (!tmdbID) { toast('请先关联 TMDB ID', 'error'); return; }
        setLoading(button, true);
        try {
          const data = await api('/api/tmdb/detail?id=' + encodeURIComponent(tmdbID) + '&type=' + encodeURIComponent(type));
          const total = Number(data.number_of_episodes || 0);
          if (!(total > 0)) { toast('TMDB 没有返回总集数', 'error'); return; }
          $('#feed-total-episodes', el).value = total;
          markDirty();
          toast('已获取总集数：' + total, 'success');
        } catch (err) {
          toast('获取总集数失败：' + err.message, 'error', 5200);
        } finally {
          setLoading(button, false);
        }
      });
      $('#feed-rename-season', el)?.addEventListener('input', event => { event.target.dataset.autofilled = ''; });
      let quarterTouched = false;
      const readQuarter = () => {
        const yearInput = $('#feed-quarter-year', el);
        const monthInput = $('#feed-quarter-month', el);
        const year = Number(yearInput && yearInput.value);
        const month = monthInput && monthInput.value;
        if (!year || year < 2000 || year > 2100 || !month) return '';
        return year + '-' + month;
      };
      const applyQuarter = () => {
        const value = readQuarter();
        if (!value) return;
        quarterTouched = true;
        setQuarter(feed, value);
        toast('已标注为 ' + quarterLabel(value), 'success');
      };
      $('#feed-quarter-year', el)?.addEventListener('change', applyQuarter);
      $('#feed-quarter-month', el)?.addEventListener('change', applyQuarter);
      if (autoSeason) setTimeout(suggestSeason, 300);
      const refreshQuarterHint = (q) => {
        if (!q || q.length !== 7) return;
        const year = $('#feed-quarter-year', el);
        const month = $('#feed-quarter-month', el);
        if (year) year.value = q.slice(0, 4);
        if (month) month.value = q.slice(5);
      };
      const loadQuarter = async (id, type) => {
        if (!id) return;
        if (quarterOf(feed)) { refreshQuarterHint(quarterOf(feed)); return; }
        try {
          const season = Number(feed.rename_season || 0);
          const seasonQuery = (type || 'tv') === 'tv' && season > 1 ? '&season=' + season : '';
          let q = '';
          try {
            const eps = await api('/api/tmdb/episodes?id=' + encodeURIComponent(id));
            q = quarterFromEpisodeList(eps.episodes, feed);
          } catch (e) { /* ignore */ }
          if (!q) {
            const data = await api('/api/tmdb/detail?id=' + encodeURIComponent(id) + '&type=' + encodeURIComponent(type || 'tv') + seasonQuery);
            q = quarterFromDate(data.season_air_date || data.first_air_date || data.release_date || '');
          }
          if (q && !quarterTouched) { setQuarter(feed, q); refreshQuarterHint(q); }
        } catch (e) { /* ignore */ }
      };
      if (feed.tmdb_id) loadQuarter(feed.tmdb_id, feed.tmdb_type);
      $('#feed-tmdb-id', el)?.addEventListener('change', ev => { loadQuarter(ev.target.value.trim(), $('#feed-tmdb-type', el).value); });
      const fillTotalFromTMDB = async (id, type) => {
        const totalInput = $('#feed-total-episodes', el);
        if (!id) return;
        try {
          const data = await api(`/api/tmdb/detail?id=${encodeURIComponent(id)}&type=${encodeURIComponent(type || 'tv')}`);
          const q = quarterFromDate(data.first_air_date || data.release_date || '');
          if (q) {
            setQuarter(feed, q);
            const year = $('#feed-quarter-year', el);
            const month = $('#feed-quarter-month', el);
            if (year) year.value = q.slice(0, 4);
            if (month) month.value = q.slice(5);
          }
          if (totalInput && Number(totalInput.value) <= 0 && Number(data.number_of_episodes) > 0) {
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
          const selectedQuarter = readQuarter();
          if (editing && feed.url && feed.url !== payload.url) setQuarter(feed, '');
          if (selectedQuarter) setQuarter(payload, selectedQuarter);
          close();
          toast(editing ? '订阅已更新' : '订阅已添加', 'success');
          if (state.route === 'feeds') renderFeedsPage();
          loadStatus().catch(err => console.debug('[AniFeed]', err));
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
    auto_disable: !!($('#feed-auto-disable', root) || {}).checked,
    cid: local.cid || '',
    savepath: local.savepath || '',
    filter: local.filter || '',
    exclude: local.exclude || '',
    expiration: Number(($('#feed-expiration', root) || {}).value || 0),
    rename_season: Number(($('#feed-rename-season', root) || {}).value || 0),
    rename_offset: Number(($('#feed-rename-offset', root) || {}).value || 0),
    total_episodes: Number(($('#feed-total-episodes', root) || {}).value || 0),
    auto_offset: !!($('#feed-auto-offset', root) || {}).checked,
    auto_backup: !!($('#feed-auto-backup', root) || {}).checked,
    auto_backup_idle: !!($('#feed-auto-backup-idle', root) || {}).checked,
    auto_backup_failures: !!($('#feed-auto-backup-failures', root) || {}).checked,
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
      start_episode: Number(row.querySelector('[data-sg="start_episode"]').value || 0),
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
    if (season === 'none') { if (quarterOf(item.feed)) return false; }
    else if (season !== 'all' && quarterOf(item.feed) !== season) return false;
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
      <thead><tr><th><input type="checkbox" id="manage-select-all" ${items.length && items.every(item => selected.has(item.index)) ? 'checked' : ''}></th><th>标题</th><th>季</th><th>播出季度</th><th>字幕组</th><th>状态</th><th>进度</th><th>类型</th><th>URL</th></tr></thead>
      <tbody>${items.map(({ feed, index, update }) => {
        const total = Number(feed.total_episodes || update.total || 0);
        const progress = total ? `${Number(update.episode || 0)} / ${total}` : `${Number(update.episode || 0)} / —`;
        const enabled = feed.enabled !== false && !update.completed;
        return `<tr>
          <td data-label="选择"><input type="checkbox" data-manage-select="${index}" ${selected.has(index) ? 'checked' : ''}></td>
          <td class="name-cell" data-label="标题"><div class="title">${esc(feed.name || feed.url)}</div><div class="meta">${feed.tmdb_id ? `TMDB ${esc(feed.tmdb_id)}` : esc(feed.subgroup_name || '')}</div></td>
          <td data-label="季">${Number(feed.rename_season || 0) || '—'}</td>
          <td data-label="播出季度">${quarterOf(feed) || '—'}</td>
          <td data-label="字幕组">${esc(feed.subgroup_name || urlSubgroupName(feed.url) || '—')}</td>
          <td data-label="状态"><span class="badge ${enabled ? 'badge-success' : ''}">${enabled ? '已启用' : '未启用'}</span></td>
          <td data-label="进度"><span class="badge ${total && Number(update.episode) >= total ? 'badge-success' : 'badge-warning'}">${progress}</span></td>
          <td data-label="类型"><span class="badge badge-info">${esc(feed.tmdb_type || 'tv')}</span></td>
          <td data-label="URL"><button class="btn btn-sm url-reveal" type="button" data-url="${esc(feed.url || '')}">查看</button></td>
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
      <select class="select" id="manage-season"><option value="all">全部播出季度</option><option value="none">未标注</option>${quarterFilterValues().map(value => `<option value="${value}">${quarterLabel(value)}</option>`).join('')}</select>
      <select class="select" id="manage-status"><option value="all">全部状态</option><option value="enabled">已启用</option><option value="disabled">未启用</option></select>
      <span class="spacer"></span><span class="badge" id="manage-count">已选 0 项</span>
    </div><div id="manage-table"></div>`,
    footer: `<button class="btn btn-sm" id="manage-export">${icon('download')}导出</button><button class="btn btn-sm" id="manage-total">${icon('refresh')}更新总集数</button><span class="spacer"></span><button class="btn btn-sm" id="manage-enable">启用</button><button class="btn btn-sm" id="manage-disable">禁用</button><button class="btn btn-sm btn-danger" id="manage-delete">删除</button><button class="btn" data-action="modal-close">关闭</button>`,
    onMount(el, closeFn) {
      close = closeFn;
      tableBox = el.querySelector('#manage-table');
      countBox = el.querySelector('#manage-count');
      draw();
      tableBox.addEventListener('click', event => {
        const btn = event.target.closest('[data-url]');
        if (!btn) return;
        event.stopPropagation();
        const url = btn.dataset.url || '';
        if (!url) { toast('该订阅没有地址', 'error'); return; }
        const modal = openModal({
          title: '订阅地址',
          subtitle: btn.closest('tr')?.querySelector('.title')?.textContent?.trim() || '',
          body: '<div class="url-box">' + esc(url) + '</div>',
          footer: '<button class="btn btn-sm" id="url-copy">' + icon('copy') + '复制</button><span class="spacer"></span><button class="btn" data-action="modal-close">关闭</button>',
        });
        modal.el.querySelector('#url-copy').addEventListener('click', async () => {
          try { await navigator.clipboard.writeText(url); toast('已复制到剪贴板', 'success'); }
          catch (e) { toast('复制失败，请手动选择', 'error'); }
        });
      });
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
    loadStatus().catch(err => console.debug('[AniFeed]', err));
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
            loadStatus().catch(err => console.debug('[AniFeed]', err));
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
async function loadTasksPage(force = false) {
  const epoch = state.routeEpoch;
  if (state.tasksLoaded) {
    renderTasksPage();
    if (!force) return;
  } else {
    $('#view').innerHTML = pageSkeleton('tasks');
  }
  try {
    if (!state.status) await loadStatus();
    if (state.routeEpoch !== epoch) return;
    if (!state.status || !state.status.has_cookies) {
      $('#view').innerHTML = `${subPageTabs('tasks')}<div class="card">${emptyBlock('lock', '尚未登录 115',
        '离线任务需要 115 登录后才能读取。请先在设置中填写 Cookies 或扫码登录。',
        `<a class="btn btn-primary" href="#/settings">${icon('login')}去登录</a>`)}</div>`;
      return;
    }
    const data = await api('/api/tasks?limit=100');
    state.tasks = data.tasks || [];
    state.tasksLoaded = true;
  } catch (err) {
    if (state.routeEpoch !== epoch) return;
    $('#view').innerHTML = `${subPageTabs('tasks')}<div class="card">${emptyBlock('cloud', '无法读取 115 任务', err.message,
      `<button class="btn btn-primary" id="tasks-retry">${icon('refresh')}重试</button>
       <a class="btn" href="#/settings">${icon('login')}检查登录</a>`)}</div>`;
    const btn = $('#tasks-retry');
    if (btn) btn.addEventListener('click', loadTasksPage);
    return;
  }
  if (state.routeEpoch !== epoch) return;
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

  $('#view').innerHTML = `${subPageTabs('tasks')}
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
    <div class="card" id="qb-fallback-card"></div>
    <div class="sticky-save" style="margin-top:16px">
      <span class="hint" style="color:var(--muted)">提示：清理只作用于 115 离线任务列表，不会删除网盘里的文件。</span>
      <span class="spacer"></span>
      <button class="btn btn-sm btn-danger" data-clear="1">清除已完成</button>
      <button class="btn btn-sm btn-danger" data-clear="3">清除失败</button>
    </div>`;

  $('#tasks-refresh').addEventListener('click', () => loadTasksPage(true));
  $('#tasks-add-magnet').addEventListener('click', openMagnetModal);
  $('#task-tabs').addEventListener('click', event => {
    const btn = event.target.closest('[data-task-filter]');
    if (!btn) return;
    state.taskFilter = btn.dataset.taskFilter;
    state.taskVisible = 20;
    $$('#task-tabs button').forEach(el => el.classList.toggle('active', el.dataset.taskFilter === state.taskFilter));
    renderTaskTable();
  });
  $('#clear-menu-btn').addEventListener('click', openClearMenu);
  $$('[data-clear]').forEach(btn => btn.addEventListener('click', () => clearTasks(Number(btn.dataset.clear))));
  renderTaskTable();
}

async function loadQBFallbackCard() {
  const box = $('#qb-fallback-card');
  if (!box) return;
  box.innerHTML = '<div class="section-title"><span class="ico">' + icon('download') + '</span><div><h3>qB 兜底记录</h3><p>115 失败后已提交到 qBittorrent 的任务</p></div></div><div class="section-body"><div class="skeleton-line"></div><div class="skeleton-line short"></div></div>';
  try {
    const data = await api('/api/qb/fallback?limit=20');
    if (!box.isConnected) return;
    const items = data.items || [];
    if (!items.length) {
      box.innerHTML = '<div class="section-title"><span class="ico">' + icon('download') + '</span><div><h3>qB 兜底记录</h3><p>115 失败后已提交到 qBittorrent 的任务</p></div></div>' + emptyBlock('cloud', '暂无 qB 兜底任务', '115 连续失败达到阈值后，任务会出现在这里。');
      return;
    }
    box.innerHTML = '<div class="section-title"><span class="ico">' + icon('download') + '</span><div><h3>qB 兜底记录</h3><p>115 失败后已提交到 qBittorrent 的任务</p></div></div>' +
      '<div class="table-wrap"><table class="table task-table"><thead><tr><th>任务名称</th><th>状态</th><th>提交时间</th><th>Hash</th></tr></thead><tbody>' +
      items.map(item => '<tr><td class="name-cell" data-label="任务"><div class="title">' + esc(item.title || '未命名任务') + '</div></td><td data-label="状态"><span class="badge ' + (item.state === 'completed' ? 'badge-success' : 'badge-info') + '">' + (item.state === 'completed' ? '已完成' : '已提交') + '</span></td><td data-label="提交时间">' + esc(new Date(item.created_at).toLocaleString('zh-CN')) + '</td><td class="numeric" data-label="Hash">' + esc((item.info_hash || '').slice(0, 16)) + '</td></tr>').join('') +
      '</tbody></table></div>';
  } catch (err) {
    if (!box.isConnected) return;
    box.innerHTML = '<div class="section-title"><span class="ico">' + icon('download') + '</span><div><h3>qB 兜底记录</h3><p>115 失败后已提交到 qBittorrent 的任务</p></div></div>' + emptyBlock('alert', '读取失败', err.message || '暂时无法读取 qB 兜底记录。');
  }
}

function visibleTasks() {
  const tasks = state.tasks || [];
  if (state.taskFilter === 'all') return tasks;
  return tasks.filter(t => t.state === state.taskFilter);
}

function renderTaskTable() {
  const container = $('#task-table-card');
  if (!container) return;
  const allTasks = visibleTasks();
  const tasks = allTasks.slice(0, state.taskVisible);
  if (!allTasks.length) {
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
      <td class="actions" data-label="操作"><button class="btn btn-sm btn-ghost task-copy-btn" data-copy-hash="${esc(magnetFromHash(task.url || task.info_hash))}" title="复制磁力链接">${icon('copy')}<span class="task-copy-label">复制磁力</span></button></td>
    </tr>`).join('')}</tbody>
  </table></div>${allTasks.length > tasks.length ? `<div style="text-align:center;padding:12px"><button class="btn btn-sm" id="task-load-more">加载更多（${allTasks.length - tasks.length}）</button></div>` : ''}`;
  const loadMore = $('#task-load-more');
  if (loadMore) loadMore.addEventListener('click', () => { state.taskVisible += 20; renderTaskTable(); });
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
    await loadTasksPage(true);
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
          loadStatus().catch(err => console.debug('[AniFeed]', err));
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
  const epoch = state.routeEpoch;
  if (state.historyLoaded) {
    renderHistoryPage();
  } else {
    $('#view').innerHTML = pageSkeleton('history');
  }
  try {
    await ensureConfig().catch(err => console.debug('[AniFeed]', err));
    const params = new URLSearchParams({ limit: String(state.historyLimit) });
    if (state.historyName) params.set('name', state.historyName);
    const data = await api(`/api/history?${params}`);
    state.history = data.items || [];
    state.historyLoaded = true;
  } catch (err) {
    if (state.routeEpoch === epoch) renderFatal(err.message);
    return;
  }
  if (state.routeEpoch !== epoch) return;
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
    <td class="actions" data-label="操作"><button class="btn btn-sm btn-ghost btn-icon-only" data-copy="${esc(magnetFromHash(item.task_hash))}" title="复制磁力链接">${icon('copy')}</button></td>
  </tr>`).join('') : '';

  $('#view').innerHTML = `${subPageTabs('history')}
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
  const epoch = state.routeEpoch;
  if (state.logTimer) { clearInterval(state.logTimer); state.logTimer = null; }
  const view = $('#view');
  if (state.logsLoaded) {
    renderLogsPage();
  } else if (!view.querySelector('#log-pre') && !view.querySelector('#log-list')) {
    view.innerHTML = pageSkeleton('logs');
  }
  try {
    const data = await api('/api/logs?limit=500');
    state.logs = data.lines || [];
    state.logsLoaded = true;
    if (state.routeEpoch !== epoch) return;
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


/* ===== details 展开/收起高度动画 ===== */
document.addEventListener('click', function (event) {
  const summary = event.target.closest('details > summary');
  if (!summary) return;
  const details = summary.parentElement;
  if (!details || !details.classList.contains('settings-fold') && !details.classList.contains('advanced-settings')) return;
  const body = details.classList.contains('settings-fold')
    ? details.querySelector(':scope > .settings-fold-body')
    : [...details.children].find(function (el) { return el.tagName !== 'SUMMARY'; });
  if (!body) return;
  event.preventDefault();
  if (details.open) {
    const h = body.scrollHeight;
    body.style.overflow = 'hidden';
    body.style.height = h + 'px';
    body.style.transition = 'height .22s ease, opacity .18s ease';
    requestAnimationFrame(function () { body.style.height = '0px'; body.style.opacity = '0'; });
    setTimeout(function () { details.open = false; body.style = ''; }, 230);
  } else {
    details.open = true;
    const h = body.scrollHeight;
    body.style.overflow = 'hidden';
    body.style.height = '0px';
    body.style.opacity = '0';
    body.style.transition = 'height .22s ease, opacity .18s ease';
    requestAnimationFrame(function () { body.style.height = h + 'px'; body.style.opacity = '1'; });
    setTimeout(function () { body.style = ''; }, 240);
  }
});
/* ===== 历史页子标签：提交历史 / 离线任务 ===== */
function subPageTabs(active) {
  return '<div class="segmented subpage-tabs">' +
    '<button type="button" data-subpage="history" class="' + (active === 'history' ? 'active' : '') + '">提交历史</button>' +
    '<button type="button" data-subpage="tasks" class="' + (active === 'tasks' ? 'active' : '') + '">离线任务</button>' +
    '</div>';
}

document.addEventListener('click', function (event) {
  const btn = event.target.closest('[data-subpage]');
  if (!btn) return;
  const tab = btn.dataset.subpage;
  if (tab === 'tasks') { state.routeEpoch++; loadTasksPage(true); }
  else { state.routeEpoch++; loadHistoryPage(); }
});
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
  $('#log-help')?.addEventListener('click', () => openModal({
    title: '运行日志说明',
    subtitle: '用于排查 RSS、115、通知和后处理问题',
    body: `<div class="kv"><div class="kv-row"><span>信息</span><code>正常运行过程</code></div><div class="kv-row"><span>完成</span><code>通知、任务完成等成功事件</code></div><div class="kv-row"><span>警告</span><code>可继续运行但需要注意的问题</code></div><div class="kv-row"><span>错误</span><code>请求失败、添加失败或后处理异常</code></div></div><div class="hint" style="margin-top:10px">日志按行显示；搜索只过滤当前已加载的日志，不读取历史文件。敏感信息不会显示。</div>`,
    footer: '<button class="btn btn-primary" data-action="modal-close">知道了</button>',
  }));
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

function renderSettingsSkeleton() {
  const skeletonCard = (iconName, title, subtitle, rows = 3) => `<section class="card">
    <div class="section-title"><span class="ico">${icon(iconName)}</span><div><h3>${title}</h3><p>${subtitle}</p></div></div>
    <div class="section-body">${Array.from({ length: rows }).map(() => '<div class="skeleton-line"></div>').join('')}</div>
  </section>`;
  $('#view').innerHTML = `
    <div class="settings-grid"><div>
      <div class="subtabs"><button disabled>账号与连接</button><button disabled>下载与通知</button><button disabled>整理与命名</button><button disabled>字幕组与规则</button><button disabled>外观与数据</button><button disabled>系统与关于</button></div>
      <div class="settings-loading-panel">
        <div class="settings-stack">
          ${skeletonCard('cloud', '115 账号', '当前登录账号与连接状态', 3)}
          ${skeletonCard('server', '服务与访问', 'Web UI 监听地址、登录保护与代理', 4)}
          ${skeletonCard('info', 'AI 助手', 'DeepSeek 兼容 API，用于查询和受控操作', 3)}
          ${skeletonCard('database', '数据库与缓存', 'SQLite 历史库与 115 请求节流参数', 3)}
        </div>
      </div>
    </div></div>`;
}


function loadAssetScript(src) {
  return new Promise((resolve, reject) => {
    const existing = document.querySelector('script[data-lazy-module="' + src + '"]');
    if (existing) { existing.addEventListener('load', resolve, { once: true }); existing.addEventListener('error', reject, { once: true }); return; }
    const script = document.createElement('script');
    script.src = src;
    script.dataset.lazyModule = src;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error('加载模块失败: ' + src));
    document.body.appendChild(script);
  });
}

async function ensureSettingsModule() {
  if (window.AniFeedSettings) return;
  const app = document.querySelector('script[src*="app.min.js"]');
  const query = app ? (new URL(app.src).search || '') : '';
  await loadAssetScript('/app-settings.min.js' + query);
}

async function loadSettingsPage() {
  await ensureSettingsModule();
  return window.AniFeedSettings.loadPage();
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
  await Promise.allSettled([loadStatus(), ensureConfig()]);
  hideTokenGate();
  renderRoute();
  syncAINav();
  if (window.AniFeedTheme && window.AniFeedTheme.loadServerTheme) window.AniFeedTheme.loadServerTheme();
  // Warm the rules/templates cache so the first "add subscription" can open
  // the rules editor immediately instead of waiting on another API round-trip.
  setTimeout(() => {
    if (!state.rules) api('/api/rules').then(data => { state.rules = data; }).catch(err => console.debug('[AniFeed]', err));
    if (!Array.isArray(state.matchTemplates)) loadMatchTemplates().catch(err => console.debug('[AniFeed]', err));
  }, 350);
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

function initTaskRowNavigation() {
  document.addEventListener('click', event => {
    const row = event.target.closest('.recent-table tbody tr, .task-table tbody tr');
    if (!row || event.target.closest('button, a, input, select, textarea')) return;
    if (state.route === 'history') return;
    location.hash = '#/history';
  });
}

async function initInputGuards() {
  document.addEventListener('click', event => {
    if (event.target.closest('input, button, select, textarea, a, .switch')) return;
    const label = event.target.closest('label');
    if (!label || !label.querySelector('.switch')) return;
    event.preventDefault();
  }, true);
}

async function boot() {
  renderBootSkeleton();
  window.addEventListener('unhandledrejection', event => {
    console.error('[AniFeed] unhandled promise rejection', event.reason);
  });
  initInputGuards();
  initTaskRowNavigation();
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
  if (state.authPending) return;
  state.authPending = true;
  const submitBtn = $('#auth-submit');
  const submitLabel = submitBtn ? submitBtn.textContent : '';
  if (submitBtn) { submitBtn.disabled = true; submitBtn.textContent = '登录中…'; }
  try {
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
    toast(error.textContent, 'error', 4200);
  }
  } finally {
    state.authPending = false;
    if (submitBtn) { submitBtn.disabled = false; submitBtn.textContent = submitLabel; }
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
    sessionStorage.removeItem('anifeed_token');
    error.textContent = err.message || '令牌错误';
    error.classList.remove('hidden');
  }
});

boot();

/* ===== 个性化设置：绑定与持久化 ===== */
function initPersonalization() {
  const T = window.AniFeedTheme;
  if (!T) return;
  const $id = id => document.getElementById(id);
  const ACCENTS = [['#2563eb', '云蓝'], ['#a2431b', '赭石'], ['#1a1916', '墨黑'], ['#0f6b5c', '松绿'], ['#245c8a', '藏蓝'], ['#8f3a5c', '绛紫'], ['#8a6a2f', '鎏金'], ['#b3342a', '朱红'], ['#4a6b3d', '竹青'], ['#33557f', '靛蓝'], ['#6a5a78', '灰紫'], ['#6d4b33', '茶褐'], ['#4a6068', '青灰']];
  const hint = (msg, isErr) => {
    const el = $id('af-io-hint');
    if (!el) return;
    el.textContent = msg || '';
    el.style.color = isErr ? 'var(--danger)' : 'var(--muted)';
  };
  const saveAll = () => {
    const t = T.read();
    try {
      T.save(t);
      hint('个性化设置已保存');
    } catch (e) {
      hint('保存失败：' + (e && e.message ? e.message : '浏览器拒绝写入'), true);
    }
  };
  const saveButton = $id('af-save');
  if (saveButton) saveButton.onclick = saveAll;

  const save = patch => {
    try {
      T.save(Object.assign({}, T.read(), patch));
    } catch (e) {
      hint('保存失败：' + (e && e.message ? e.message : '图片可能过大'), true);
      return false;
    }
    draw();
    return true;
  };

  function draw() {
    const t = T.read();

    const navBox = $id('af-nav');
    if (navBox) {
      navBox.innerHTML = Object.entries(T.NAV_PRESETS).map(([key, v]) =>
        `<button type="button" class="af-swatch af-swatch-nav ${t.nav === key ? 'is-active' : ''}" data-nav="${key}" title="${v.label}" style="background:${v.light}"></button>`).join('');
      navBox.onclick = ev => {
        const b = ev.target.closest('[data-nav]');
        if (!b) return;
        save({ nav: b.dataset.nav });
      };
    }

    const accBox = $id('af-accent');
    if (accBox) {
      accBox.innerHTML = ACCENTS.map(([c, name]) =>
        `<button type="button" class="af-swatch ${t.accent.toLowerCase() === c ? 'is-active' : ''}" data-accent="${c}" title="${name}" style="background:${c}"></button>`).join('')
        + `<input type="color" class="af-swatch-custom" id="af-accent-custom" value="${t.accent}" title="自定义颜色">`;
      accBox.onclick = ev => {
        const b = ev.target.closest('[data-accent]');
        if (!b) return;
        save({ accent: b.dataset.accent });
      };
      const custom = $id('af-accent-custom');
      if (custom) custom.oninput = () => save({ accent: custom.value });
    }

    const modeBox = $id('af-mode');
    if (modeBox) {
      const mode = T.read().mode;
      modeBox.querySelectorAll('[data-mode]').forEach(b => b.classList.toggle('is-active', b.dataset.mode === mode));
      modeBox.onclick = ev => {
        const b = ev.target.closest('[data-mode]');
        if (!b) return;
        save({ mode: b.dataset.mode });
      };
    }

    const radius = $id('af-radius');
    if (radius) {
      const r = t.radius === undefined ? T.DEFAULTS.radius : t.radius;
      radius.value = r;
      $id('af-radius-val').textContent = r + 'px';
      radius.oninput = () => { $id('af-radius-val').textContent = radius.value + 'px'; save({ radius: Number(radius.value) }); };
    }
    const elev = $id('af-elev');
    if (elev) {
      const e = t.elevation === undefined ? T.DEFAULTS.elevation : t.elevation;
      elev.value = e;
      $id('af-elev-val').textContent = e + '%';
      elev.oninput = () => { $id('af-elev-val').textContent = elev.value + '%'; save({ elevation: Number(elev.value) }); };
    }

    const blur = $id('af-blur');
    if (blur) {
      blur.value = t.blur;
      $id('af-blur-val').textContent = t.blur + 'px';
      blur.oninput = () => { $id('af-blur-val').textContent = blur.value + 'px'; save({ blur: Number(blur.value) }); };
    }
    const bgBlur = $id('af-bg-blur');
    if (bgBlur) {
      const value = t.bgBlur === undefined ? 8 : Number(t.bgBlur);
      bgBlur.value = value;
      $id('af-bg-blur-val').textContent = value + 'px';
      bgBlur.oninput = () => { $id('af-bg-blur-val').textContent = bgBlur.value + 'px'; save({ bgBlur: Number(bgBlur.value) }); };
    }
    const alpha = $id('af-alpha');
    if (alpha) {
      const pct = Math.round(t.alpha * 100);
      alpha.value = pct;
      $id('af-alpha-val').textContent = pct + '%';
      alpha.oninput = () => { $id('af-alpha-val').textContent = alpha.value + '%'; save({ alpha: Number(alpha.value) / 100 }); };
    }

    const prev = $id('af-bg-preview');
    if (prev) {
      if (t.bgImage) {
        prev.style.backgroundImage = `url("${t.bgImage}")`;
        prev.classList.add('has-image');
        prev.textContent = '';
      } else {
        prev.style.backgroundImage = '';
        prev.classList.remove('has-image');
        prev.textContent = '未设置背景';
      }
    }
    const url = $id('af-bg-url');
    if (url && !String(t.bgImage).startsWith('data:')) url.value = t.bgImage || '';
    if (serverTheme) drawServer();
  }

  const applyUrl = async () => {
    const url = $id('af-bg-url');
    if (!url) return;
    const value = url.value.trim();
    if (!value) { hint('地址是空的：填写图片/API 地址再点应用；要清空背景请点「清除」', true); return; }
    if (!/^https?:\/\//i.test(value) && !value.startsWith('data:')) {
      hint('请填写 http(s) 图片地址，或使用「上传」', true);
      return;
    }
    const btn = $id('af-bg-apply');
    const label = btn ? btn.textContent : '';
    if (btn) { btn.disabled = true; btn.textContent = '加载中…'; }
    hint('正在检查图片…');
    try {
      const src = T.resolveSource ? await T.resolveSource(value) : value;
      if (save({ bgImage: src })) {
        hint(src === value ? '背景已应用' : '接口返回的是图片地址，已自动解析并应用');
      }
    } catch (e) {
      hint('图片无法加载：' + (e && e.message ? e.message : '地址不可用') + '（已保留原背景）', true);
    } finally {
      if (btn) { btn.disabled = false; btn.textContent = label; }
    }
  };
  const bgApply = $id('af-bg-apply');
  if (bgApply) bgApply.onclick = applyUrl;
  const bgUrl = $id('af-bg-url');
  if (bgUrl) bgUrl.addEventListener('keydown', ev => { if (ev.key === 'Enter') { ev.preventDefault(); applyUrl(); } });
  const bgClear = $id('af-bg-clear');
  if (bgClear) bgClear.onclick = async () => {
    const u = $id('af-bg-url'); if (u) u.value = '';
    save({ bgImage: '' });
    if (serverTheme) {
      serverTheme[bgDevice] = { mode: 'none', selected: '' };
      await saveServer();
    } else {
      hint('背景已清除');
    }
  };
  const bgUpload = $id('af-bg-upload');
  const bgFile = $id('af-bg-file');
  if (bgUpload && bgFile) {
    bgUpload.onclick = () => bgFile.click();
    bgFile.onchange = async () => {
      const files = Array.from(bgFile.files || []);
      if (!files.length) return;
      hint('正在上传到服务器…');
      try {
        for (const file of files) {
          if (file.size > 10 * 1024 * 1024) throw new Error(file.name + ' 超过 10MB');
          const form = new FormData();
          form.append('file', file);
          const res = await fetch('/api/theme/upload?device=' + encodeURIComponent(bgDevice), { method: 'POST', body: form, headers: { accept: 'application/json' } });
          const data = await res.json().catch(() => ({}));
          if (!res.ok) throw new Error(data.error || '上传失败');
          serverTheme = data.theme;
        }
        bgFile.value = '';
        applyServer(); drawServer();
        hint('已上传 ' + files.length + ' 张图片');
      } catch (e) { hint(e.message, true); bgFile.value = ''; }
    };
  }

  let serverTheme = null;
  let bgDevice = localStorage.getItem('anifeed_bg_device') || ((window.matchMedia && window.matchMedia('(max-width: 860px)').matches) ? 'mobile' : 'desktop');
  const serverImages = () => (serverTheme && serverTheme.images && serverTheme.images[bgDevice]) || [];
  const serverDevice = () => (serverTheme && serverTheme[bgDevice]) || { mode: 'single', selected: '' };
  const pickServerImage = () => {
    const images = serverImages();
    if (!images.length) return null;
    const cfg = serverDevice();
    return images.find(item => item.name === cfg.selected) || images[0];
  };
  const applyServer = () => {
    if (!serverTheme) return;
    window.AniFeedTheme.setServerTheme(serverTheme);
    const picked = pickServerImage();
    const prev = $id('af-bg-preview');
    if (prev) {
      if (picked) { prev.style.backgroundImage = `url("${picked.url}")`; prev.classList.add('has-image'); prev.textContent = ''; }
      else { prev.style.backgroundImage = ''; prev.classList.remove('has-image'); prev.textContent = '未设置背景'; }
    }
  };
  const saveServer = async () => {
    if (!serverTheme) return;
    try {
      const resp = await api('/api/theme', { method: 'PUT', body: { desktop: serverTheme.desktop, mobile: serverTheme.mobile } });
      serverTheme = resp;
      applyServer();
      drawServer();
      hint('背景设置已保存到服务器');
    } catch (e) { hint('保存背景失败：' + e.message, true); }
  };
  function drawServer() {
    const tabs = $id('af-bg-device');
    if (tabs) {
      tabs.querySelectorAll('[data-device]').forEach(btn => btn.classList.toggle('active', btn.dataset.device === bgDevice));
    }
    if (!serverTheme) return;
    const images = serverImages();
    const cfg = serverDevice();
    const mode = $id('af-bg-mode');
    if (mode) { mode.value = cfg.mode || 'single'; mode.onchange = async () => { serverTheme[bgDevice] = { mode: mode.value, selected: cfg.selected || '' }; await saveServer(); }; }
    const count = $id('af-bg-count');
    if (count) count.textContent = `${images.length} 张图片`;
    const gallery = $id('af-bg-gallery');
    if (!gallery) return;
    if (!images.length) {
      gallery.innerHTML = '<div class="af-bg-empty">还没有上传背景图片</div>';
    } else {
      const selected = cfg.mode === 'single' ? (cfg.selected || images[0].name) : '';
      gallery.innerHTML = images.map(img => `<div class="af-bg-item ${img.name === selected ? 'active' : ''}"><img src="${img.url}" alt=""><div class="af-bg-item-actions"><button type="button" class="btn btn-sm ${img.name === selected ? 'btn-primary' : ''}" data-bg-pick="${img.name}">${img.name === selected ? '使用中' : '使用'}</button>${img.builtin ? '<span class="af-bg-builtin">内置</span>' : `<button type="button" class="btn btn-sm btn-ghost" data-bg-remove="${img.name}" title="删除">×</button>`}</div></div>`).join('');
      gallery.onclick = async event => {
        const pick = event.target.closest('[data-bg-pick]');
        const del = event.target.closest('[data-bg-remove]');
        if (pick) {
          serverTheme[bgDevice] = { mode: 'single', selected: pick.dataset.bgPick };
          await saveServer();
        }
        if (del) {
          try {
            serverTheme = await api('/api/theme/image/delete?name=' + encodeURIComponent(del.dataset.bgRemove), { method: 'DELETE' });
            applyServer(); drawServer(); hint('图片已删除');
          } catch (e) { hint('删除失败：' + e.message, true); }
        }
      };
    }
    applyServer();
  }
  async function loadServer() {
    try {
      serverTheme = await api('/api/theme');
      window.AniFeedTheme.setServerTheme(serverTheme);
      drawServer();
    } catch (e) { hint('读取服务器背景失败：' + e.message, true); }
  }
  const doExport = () => {
    const data = T.exportDesign();
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'anifeed-theme-' + new Date().toISOString().slice(0, 10) + '.json';
    document.body.appendChild(a);
    a.click();
    setTimeout(() => { URL.revokeObjectURL(a.href); a.remove(); }, 0);
    hint('已导出主题 JSON');
  };
  ['af-export', 'af-export-2'].forEach(id => { const b = $id(id); if (b) b.onclick = doExport; });

  const importBtn = $id('af-import');
  const importFile = $id('af-import-file');
  if (importBtn && importFile) {
    importBtn.onclick = () => importFile.click();
    importFile.onchange = () => {
      const file = importFile.files && importFile.files[0];
      if (!file) return;
      const reader = new FileReader();
      reader.onload = () => {
        try {
          T.importDesign(String(reader.result));
          draw();
          hint('已导入主题：' + file.name);
        } catch (e) {
          hint('导入失败：' + (e && e.message ? e.message : '文件格式不正确'), true);
        }
        importFile.value = '';
      };
      reader.onerror = () => hint('读取文件失败', true);
      reader.readAsText(file);
    };
  }

  const reset = $id('af-reset');
  if (reset) reset.onclick = () => { T.reset(); draw(); hint('已恢复默认主题'); };

  const deviceTabs = $id('af-bg-device');
  if (deviceTabs) {
    deviceTabs.onclick = event => {
      const btn = event.target.closest('[data-device]');
      if (!btn || btn.dataset.device === bgDevice) return;
      bgDevice = btn.dataset.device;
      localStorage.setItem('anifeed_bg_device', bgDevice);
      drawServer();
    };
  }

  draw();
  loadServer();
}
