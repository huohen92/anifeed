/* AniFeed 主题/个性化运行时
 * 存储：localStorage.anifeed.theme.v1（明暗模式仍用 anifeed_theme，兼容原有切换按钮）
 * 页面加载时立即应用，避免闪烁；设置面板通过 window.AniFeedTheme 读写。
 */
(function () {
  var KEY = 'anifeed.theme.v1';
  var MODE_KEY = 'anifeed_theme';
  var DEFAULTS = {
    bgImage: '',
    bgBlur: 8,
    blur: 14,
    alpha: 0.72,
    accent: '#2563eb',
    nav: 'ink',
    mode: 'auto',
    radius: 16,
    elevation: 45,
  };
  var NAV_PRESETS = {
    ink:   { label: '纯白', light: '#ffffff', dark: '#121826', border: '#e3e9f4', borderDark: '#232c3e' },
    warm:  { label: '暖砂', light: '#fbf8f2', dark: '#1a1712', border: '#efe7d8', borderDark: '#2c2620' },
    mint:  { label: '松绿', light: '#f2f8f5', dark: '#111a16', border: '#e2ede7', borderDark: '#22302a' },
    slate: { label: '雾蓝', light: '#f2f6fb', dark: '#101722', border: '#e3eaf4', borderDark: '#232a31' },
  };
  var root = document.documentElement;
  var serverTheme = null;
  var serverPickName = '';

  function deviceKey() {
    return window.matchMedia && window.matchMedia('(max-width: 860px)').matches ? 'mobile' : 'desktop';
  }

  function serverImage() {
    if (!serverTheme) return null;
    var key = deviceKey();
    var settings = (serverTheme[key] || {});
    var images = ((serverTheme.images || {})[key] || []);
    if (!images.length) return null;
    if (settings.mode === 'none') return null;
    serverPickName = settings.selected || '';
    return images.find(function (item) { return item.name === serverPickName; }) || images[0];
  }

  function applyServerBackground() {
    if (serverTheme && (serverTheme[deviceKey()] || {}).mode === 'none') {
      root.style.setProperty('--af-bg-image', 'none');
      root.removeAttribute('data-af-bg');
      root.removeAttribute('data-af-glass');
      return true;
    }
    var image = serverImage();
    if (image) {
      root.style.setProperty('--af-bg-image', 'url("' + String(image.url).replace(/"/g, '\\"') + '")');
      root.setAttribute('data-af-bg', 'on');
      root.setAttribute('data-af-glass', 'on');
      return true;
    }
    return false;
  }

  function setServerTheme(data) {
    serverTheme = data || null;
    serverPickName = '';
    apply();
    return serverTheme;
  }

  function loadServerTheme() {
    return fetch('/api/theme', { credentials: 'same-origin', headers: { accept: 'application/json' } })
      .then(function (res) { return res.ok ? res.json() : null; })
      .then(function (data) {
        if (data) {
          serverTheme = data;
          apply();
          window.dispatchEvent(new CustomEvent('anifeed-theme-server', { detail: data }));
        }
        return data;
      })
      .catch(function () { return null; });
  }

  function read() {
    var saved = {};
    try { saved = JSON.parse(localStorage.getItem(KEY) || '{}') || {}; } catch (e) { saved = {}; }
    var t = {};
    Object.keys(DEFAULTS).forEach(function (k) { t[k] = saved[k] === undefined ? DEFAULTS[k] : saved[k]; });
    // 明暗模式以原有键为准
    var mode = localStorage.getItem(MODE_KEY);
    t.mode = mode ? mode : (saved.mode === 'light' || saved.mode === 'dark' ? saved.mode : 'auto');
    return t;
  }

  function write(t) {
    var payload = Object.assign({}, t);
    if (t.mode === 'auto') {
      delete payload.mode;
    } else {
      payload.mode = t.mode;
      localStorage.setItem(MODE_KEY, t.mode);
    }
    try {
      localStorage.setItem(KEY, JSON.stringify(payload));
    } catch (e) {
      // 超出配额（通常是背景图太大）
      try {
        localStorage.setItem(KEY, JSON.stringify(Object.assign({}, payload, { bgImage: '' })));
      } catch (e2) { /* 放弃持久化 */ }
      throw e;
    }
    apply(t);
    return read();
  }

  function resolveMode(t) {
    if (t.mode === 'light' || t.mode === 'dark') return t.mode;
    return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function apply(t) {
    t = t || read();
    var theme = resolveMode(t);
    if (root.dataset.theme !== theme) root.dataset.theme = theme;

    var nav = NAV_PRESETS[t.nav] || NAV_PRESETS.ink;
    root.style.setProperty('--af-nav-bg', theme === 'dark' ? nav.dark : nav.light);
    root.style.setProperty('--af-nav-border', theme === 'dark' ? nav.borderDark : nav.border);
    root.style.setProperty('--af-accent', t.accent || DEFAULTS.accent);
    root.style.setProperty('--af-blur', (Number(t.blur) || 0) + 'px');
    root.style.setProperty('--af-bg-blur', (Number(t.bgBlur) || 0) + 'px');
    root.style.setProperty('--af-alpha', String(Number(t.alpha) || 0.72));
    root.style.setProperty('--af-radius', (t.radius === undefined ? DEFAULTS.radius : t.radius) + 'px');
    root.style.setProperty('--af-elev', String((t.elevation === undefined ? DEFAULTS.elevation : t.elevation) / 100));
    var serverBg = applyServerBackground();
    if (!serverBg) {
      root.style.setProperty('--af-bg-image', t.bgImage ? 'url("' + t.bgImage.replace(/"/g, '\\"') + '")' : 'none');
      if (t.bgImage) {
        root.setAttribute('data-af-bg', 'on');
        root.setAttribute('data-af-glass', 'on');
      } else {
        root.removeAttribute('data-af-bg');
        root.removeAttribute('data-af-glass');
      }
    }
    root.setAttribute('data-af-nav', t.nav || 'ink');
  }

  function reset() {
    try { localStorage.removeItem(KEY); localStorage.removeItem(MODE_KEY); } catch (e) {}
    apply(read());
  }

  function exportDesign() {
    var t = read();
    return {
      app: 'anifeed',
      kind: 'theme',
      version: 1,
      exportedAt: new Date().toISOString(),
      theme: t,
    };
  }

  function importDesign(json) {
    var data = typeof json === 'string' ? JSON.parse(json) : json;
    var t = data && data.theme ? data.theme : data;
    if (!t || typeof t !== 'object') throw new Error('不是有效的主题文件');
    var next = read();
    Object.keys(DEFAULTS).forEach(function (k) {
      if (t[k] !== undefined) next[k] = t[k];
    });
    if (next.mode === 'auto') { try { localStorage.removeItem(MODE_KEY); } catch (e) {} }
    return write(next);
  }

  function onChange(cb) {
    if (typeof cb === 'function') listeners.push(cb);
  }
  var listeners = [];

  function save(t) {
    var next = write(t);
    listeners.forEach(function (fn) { try { fn(next); } catch (e) {} });
    return next;
  }


  /* 背景图来源解析：直接图片 / 返回 JSON 的 API 都可以 */
  function probeImage(src) {
    return new Promise(function (resolve, reject) {
      if (!src) return resolve('');
      if (String(src).indexOf('data:') === 0) return resolve(src);
      var img = new Image();
      var done = false;
      img.onload = function () { if (!done) { done = true; resolve(src); } };
      img.onerror = function () { if (!done) { done = true; reject(new Error('不是可用的图片地址')); } };
      img.src = src;
      setTimeout(function () { if (!done) { done = true; reject(new Error('图片加载超时')); } }, 12000);
    });
  }

  function findImageUrl(node, depth) {
    if (depth > 6 || node == null) return '';
    if (typeof node === 'string') {
      if (/^https?:\/\//i.test(node) && /(\.(png|jpe?g|webp|gif|avif|svg)([?#]|$))|image|img|cover|pic/i.test(node)) return node;
      return '';
    }
    if (Array.isArray(node)) {
      for (var i = 0; i < node.length; i++) { var a = findImageUrl(node[i], depth + 1); if (a) return a; }
      return '';
    }
    if (typeof node === 'object') {
      var keys = ['url', 'image', 'img', 'imageUrl', 'image_url', 'cover', 'pic', 'src', 'data', 'result', 'images'];
      for (var k = 0; k < keys.length; k++) {
        if (node[keys[k]] !== undefined) { var b = findImageUrl(node[keys[k]], depth + 1); if (b) return b; }
      }
      for (var key in node) { if (Object.prototype.hasOwnProperty.call(node, key)) { var c = findImageUrl(node[key], depth + 1); if (c) return c; } }
    }
    return '';
  }

  function extractImageUrl(text) {
    try {
      var j = JSON.parse(text);
      var fromJson = findImageUrl(j, 0);
      if (fromJson) return fromJson;
    } catch (e) { /* 不是 JSON */ }
    var m = String(text).match(/https?:\/\/[^\s"'<>]+\.(?:png|jpe?g|webp|gif|avif|svg)(?:\?[^\s"'<>]*)?/i);
    return m ? m[0] : '';
  }

  function resolveSource(raw) {
    var value = String(raw || '').trim();
    return probeImage(value).then(function (ok) {
      return ok;
    }, function () {
      // 当作 API：取回内容再找图片地址
      return fetch(value, { headers: { accept: 'application/json,image/*,*/*' } }).catch(function(){ throw new Error('地址打不开（网络或跨域限制）'); })
        .then(function (res) { return res.text().then(function (t) { return { res: res, text: t }; }); })
        .then(function (r) {
          if (r.res.ok === false) throw new Error('接口返回 ' + r.res.status);
          var url = extractImageUrl(r.text);
          if (!url) throw new Error('返回内容里没有找到图片地址');
          try { url = new URL(url, value).href; } catch (e) {}
          return probeImage(url);
        });
    });
  }

  window.AniFeedTheme = {
    resolveSource: resolveSource,
    DEFAULTS: DEFAULTS,
    NAV_PRESETS: NAV_PRESETS,
    read: read,
    save: save,
    apply: apply,
    reset: reset,
    exportDesign: exportDesign,
    importDesign: importDesign,
    onChange: onChange,
    resolveMode: resolveMode,
    setServerTheme: setServerTheme,
    loadServerTheme: loadServerTheme,
    serverTheme: function () { return serverTheme; },
  };

  apply(read());
  loadServerTheme();
  if (window.addEventListener) {
    var lastDevice = deviceKey();
    window.addEventListener('resize', function () {
      if (deviceKey() !== lastDevice) {
        lastDevice = deviceKey();
        serverPickName = '';
        apply(read());
      }
    });
  }
  if (window.matchMedia) {
    var mq = window.matchMedia('(prefers-color-scheme: dark)');
    var onScheme = function () { if (read().mode === 'auto') apply(read()); };
    if (mq.addEventListener) mq.addEventListener('change', onScheme);
    else if (mq.addListener) mq.addListener(onScheme);
  }
})();
