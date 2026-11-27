/* Lazily loaded Settings page module. */
(function () {
'use strict';
async function loadSettingsPage() {
  const epoch = state.routeEpoch;
  if (state.config) {
    renderSettingsPage();
  } else {
    renderSettingsSkeleton();
  }
  try {
    const resp = await api('/api/config');
    applyConfigState(resp.config, resp.meta);
    await loadStatus().catch(err => console.debug('[AniFeed]', err));
  } catch (err) {
    if (state.routeEpoch === epoch) renderFatal(err.message);
    return;
  }
  if (state.routeEpoch !== epoch) return;
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
  orders.custom = {};
  for (const key of Object.keys(PRIORITY_META)) orders.custom[key] = [];
  for (const entry of ((cfg.offline || {}).auto_custom || [])) {
    const raw = String(entry || '');
    const at = raw.indexOf('|');
    if (at <= 0) continue;
    const key = raw.slice(0, at);
    const value = raw.slice(at + 1).trim();
    if (value && orders.custom[key]) orders.custom[key].push(value);
  }
  state.autoPriority = orders;
  const draw = () => {
    box.innerHTML = Object.entries(PRIORITY_META).map(([key, meta]) => {
      const parts = orders[key].map((id, index) => {
        const label = (meta.options.find(([optionId]) => optionId === id) || [id, id])[1];
        return `<span class="priority-item">${esc(label)}<span class="priority-buttons"><button type="button" data-priority-move="up" data-priority-key="${key}" data-priority-index="${index}" ${index === 0 ? 'disabled' : ''}>↑</button><button type="button" data-priority-move="down" data-priority-key="${key}" data-priority-index="${index}" ${index === orders[key].length - 1 ? 'disabled' : ''}>↓</button><button type="button" data-priority-remove="${key}" data-priority-index="${index}" title="移除">×</button></span></span>`;
      }).join('');
      const available = meta.options.filter(([id]) => !orders[key].includes(id));
      const customParts = (orders.custom[key] || []).map((value, index) => `<span class="priority-item custom">${esc(value)}<span class="priority-buttons"><button type="button" data-priority-custom-remove="${key}" data-priority-index="${index}" title="移除">×</button></span></span>`).join('');
      const add = `<span class="priority-add"><select data-priority-select="${key}" class="select select-compact">${available.length ? available.map(([id, label]) => `<option value="${id}">${esc(label)}</option>`).join('') : '<option value="">全部已添加</option>'}</select><button type="button" class="btn btn-sm btn-icon-only" data-priority-add="${key}" title="添加" ${available.length ? '' : 'disabled'}>+</button></span><span class="priority-add custom-add"><input class="input input-compact" data-priority-custom-input="${key}" placeholder="自定义正则"><button type="button" class="btn btn-sm btn-icon-only" data-priority-custom-add="${key}" title="添加自定义">+</button></span>`;
      return `<div class="priority-row"><span class="priority-label">${esc(meta.label)}</span><div class="priority-list">${parts || '<span class="hint">未设置</span>'}${customParts}${add}</div></div>`;
    }).join('') + `<div class="priority-footer"><button type="button" class="btn btn-sm btn-ghost" data-priority-reset>恢复默认顺序</button><span class="hint">每个分类都可以单独移除、添加和排序</span></div>`;
    box.querySelector('[data-priority-reset]')?.addEventListener('click', () => {
      for (const key of Object.keys(PRIORITY_DEFAULTS)) orders[key] = [...PRIORITY_DEFAULTS[key]];
      draw();
    });
    $$('[data-priority-move]', box).forEach(button => button.addEventListener('click', () => {
      const key = button.dataset.priorityKey;
      const index = Number(button.dataset.priorityIndex);
      const target = button.dataset.priorityMove === 'up' ? index - 1 : index + 1;
      if (target < 0 || target >= orders[key].length) return;
      [orders[key][index], orders[key][target]] = [orders[key][target], orders[key][index]];
      draw();
    }));
    $$('[data-priority-remove]', box).forEach(button => button.addEventListener('click', () => {
      const key = button.dataset.priorityRemove;
      const index = Number(button.dataset.priorityIndex);
      if (index >= 0) orders[key].splice(index, 1);
      draw();
    }));
    $$('[data-priority-add]', box).forEach(button => button.addEventListener('click', () => {
      const key = button.dataset.priorityAdd;
      const select = button.parentElement.querySelector('[data-priority-select]');
      const id = select?.value;
      if (!id) return;
      if (!orders[key].includes(id)) orders[key].push(id);
      draw();
    }));
    $$('[data-priority-custom-remove]', box).forEach(button => button.addEventListener('click', () => {
      const key = button.dataset.priorityCustomRemove;
      orders.custom[key].splice(Number(button.dataset.priorityIndex), 1);
      draw();
    }));
    $$('[data-priority-custom-add]', box).forEach(button => button.addEventListener('click', () => {
      const key = button.dataset.priorityCustomAdd;
      const input = button.parentElement.querySelector('[data-priority-custom-input]');
      const value = String(input?.value || '').trim();
      if (!value) return;
      if (!orders.custom[key].includes(value)) orders.custom[key].push(value);
      draw();
    }));
    $$('[data-priority-custom-input]', box).forEach(input => input.addEventListener('keydown', event => {
      if (event.key === 'Enter') { event.preventDefault(); input.parentElement.querySelector('[data-priority-custom-add]')?.click(); }
    }));
  };
  draw();
}

function initSettingsSections() {
  const cards = $$('.settings-stack > section.card');
  if (!cards.length) return;
  const groupByTitle = {
    '115 账号': 'connection',
    '个性化': 'personal',
    'AI 助手': 'connection',
    '服务与访问': 'connection',
    '默认离线设置': 'download',
    'qBittorrent': 'download',
    'qBittorrent': 'download',
    '自动运行': 'download',
    '重命名': 'organize',
    'TMDB': 'organize',
    '洗版': 'organize',
    '多字幕组': 'subgroups',
    '喜好字幕组': 'subgroups',
    '自动备用字幕组': 'subgroups',
    '字幕组匹配模板': 'subgroups',
    '自动应用最佳规则': 'subgroups',
    '字幕组别名': 'subgroups',
    '通知': 'download',
    '数据库与缓存': 'about',
    '关于': 'about',
  };
  const sectionMeta = {
    connection: ['账号与连接', '115 账号、AI 助手、Web 访问和代理设置'],
    download: ['下载与通知', '115、qBittorrent、自动运行和通知设置'],
    organize: ['整理与命名', '重命名、TMDB 关联和洗版'],
    subgroups: ['字幕组与规则', '喜好顺序、匹配模板、别名和自动规则'],
    about: ['系统与关于', '数据库、版本和项目信息'],
    personal: ['外观与数据', '背景图、玻璃效果、强调色与导入导出'],
  };

  const foldMeta = {
    '115 账号': ['账号与登录', '登录状态、Cookie 和默认保存位置'],
    '默认离线设置': ['保存与任务行为', '保存目录、任务行为、自动规则和优先级'],
    'qBittorrent': ['qBittorrent 兜底', '本地下载器地址、账号和兜底阈值'],
    '自动运行': ['高级运行设置', '自动运行间隔、并发和通知行为'],
    '重命名': ['命名规则', '命名模板、季数偏移和立即处理'],
    'TMDB': ['TMDB 接口', '接口地址、API Key 和图片代理'],
    '洗版': ['高级整理设置', '洗版、清理和迁移规则'],
    '通知': ['通知与推送', '推送渠道、事件开关和消息模板'],
    '多字幕组': ['实验功能', '多字幕组与多 RSS 共存'],
    '自动备用字幕组': ['高级备用 RSS 设置', '失败触发、自动偏移和跳过策略'],
    '喜好字幕组': ['字幕组偏好', '自动备用和优先选择使用的字幕组顺序'],
    '字幕组匹配模板': ['模板维护', '维护字幕组预设和临时规则'],
    '自动应用最佳规则': ['自动规则优先级', '语言、字幕、容器和分辨率的自动选择顺序'],
    '字幕组别名': ['别名映射', '主名称与别名的自动映射'],
    'AI 助手': ['AI 设置', 'DeepSeek 接口、权限与密钥'],
  };
  cards.forEach(card => {
    const title = (card.querySelector('.section-title h3') || {}).textContent || '';
    card.dataset.settingsGroup = groupByTitle[title.trim()] || 'connection';
    const fold = foldMeta[title.trim()];
    if (fold && !card.closest('.settings-fold')) {
      const details = document.createElement('details');
      details.className = 'settings-fold';
      details.dataset.settingsGroup = card.dataset.settingsGroup;
      details.innerHTML = `<summary><span class="fold-title">${fold[0]}</span><small>${fold[1]}</small><span class="fold-state"></span></summary><div class="settings-fold-body"></div>`;
      card.parentNode.insertBefore(details, card);
      details.querySelector('.settings-fold-body').appendChild(card);
    }
  });
  // 连接与账号：115 账号 -> AI 助手 -> 服务与访问（放最后）
  const CONNECTION_ORDER = { '115 账号': 0, 'AI 助手': 1, '服务与访问': 2 };
  const stack = document.querySelector('.settings-stack');
  if (stack) {
    const conn = [...stack.children].filter(node => (node.dataset && node.dataset.settingsGroup) === 'connection');
    conn.sort((a, b) => {
      const titleOf = node => {
        const h = node.querySelector('.section-title h3');
        if (h) return h.textContent.trim();
        const t = node.querySelector('.fold-title');
        return t ? t.textContent.trim() : '';
      };
      const ra = CONNECTION_ORDER[titleOf(a)] ?? 99;
      const rb = CONNECTION_ORDER[titleOf(b)] ?? 99;
      return ra - rb;
    });
    conn.forEach(node => stack.appendChild(node));
  }
  // 关于页：关于 -> 数据库与缓存
  const ABOUT_ORDER = { '关于': 0, '数据库与缓存': 1 };
  if (stack) {
    const about = [...stack.children].filter(node => (node.dataset && node.dataset.settingsGroup) === 'about');
    about.sort((a, b) => {
      const titleOf = node => { const h = node.querySelector('.section-title h3'); return h ? h.textContent.trim() : ''; };
      return (ABOUT_ORDER[titleOf(a)] ?? 99) - (ABOUT_ORDER[titleOf(b)] ?? 99);
    });
    about.forEach(node => stack.appendChild(node));
  }
  const buttons = $$('[data-settings-section]');
  if (!buttons.some(button => button.dataset.settingsSection === state.settingsSection)) state.settingsSection = 'connection';
  const intro = $('#settings-section-intro');
  const apply = () => {
    const meta = sectionMeta[state.settingsSection] || sectionMeta.connection;
    buttons.forEach(button => button.classList.toggle('active', button.dataset.settingsSection === state.settingsSection));
    cards.forEach(card => {
      card.style.display = card.dataset.settingsGroup === state.settingsSection ? '' : 'none';
    });
    $$('.settings-fold').forEach(fold => {
      fold.style.display = fold.dataset.settingsGroup === state.settingsSection ? '' : 'none';
    });
    if (intro) intro.innerHTML = `<strong>${meta[0]}</strong><span>${meta[1]}</span>`;
    const savebar = document.getElementById('settings-savebar');
    if (savebar) savebar.style.display = state.settingsSection === 'personal' ? 'none' : '';
  };
  buttons.forEach(button => button.addEventListener('click', () => {
    state.settingsSection = button.dataset.settingsSection;
    apply();
  }));
  apply();
}



/* ===== 下拉多选 ===== */
function initMultiSelects() {
  const defs = [
    { trigger: '#ms-content-trigger', menu: '#ms-content-menu', label: '#ms-content-label', items: '#notify-content-summary, #notify-content-details, #notify-content-errors, #notify-content-time' },
    { trigger: '#ms-events-trigger', menu: '#ms-events-menu', label: '#ms-events-label', items: '.notify-event' },
  ];
  defs.forEach(function (def) {
    const trigger = $(def.trigger);
    const menu = $(def.menu);
    const label = $(def.label);
    if (!trigger || !menu || !label) return;
    const sync = function () {
      const names = $$(def.items).filter(function (el) { return el.checked; }).map(function (el) {
        const row = el.closest('.ms-item, .check-inline');
        if (!row) return el.dataset.event || '';
        const copy = row.cloneNode(true);
        const input = copy.querySelector('input');
        if (input) input.remove();
        return copy.textContent.trim();
      }).filter(Boolean);
      label.textContent = names.length ? names.join('、') : '未选择';
      label.classList.toggle('is-empty', names.length === 0);
    };
    trigger.addEventListener('click', function (event) {
      event.stopPropagation();
      const open = menu.classList.toggle('open');
      trigger.setAttribute('aria-expanded', open ? 'true' : 'false');
    });
    $$(def.items).forEach(function (el) { el.addEventListener('change', sync); });
    document.addEventListener('click', function (event) {
      if (!menu.contains(event.target) && !trigger.contains(event.target)) menu.classList.remove('open');
    });
    sync();
  });
}


/* ===== 开关行：开关在前，说明折叠 ===== */
function enhanceSwitchRows() {
  const content = document.getElementById('settings-content');
  if (!content) return;
  [...content.querySelectorAll('.check')].forEach(check => {
    const field = check.closest('.field');
    if (!field || field.dataset.switchRow) return;
    const switchEl = check.querySelector('.switch');
    if (!switchEl) return;
    const titleEl = [...field.children].find(el => el !== check && el.tagName === 'SPAN' && !el.classList.contains('switch'));
    if (!titleEl) return;
    const descEl = [...check.children].find(el => el !== switchEl);
    field.dataset.switchRow = '1';
    field.classList.add('switch-row');

    const row = document.createElement('div');
    row.className = 'switch-row-inner';
    row.appendChild(switchEl);

    const title = document.createElement('span');
    title.className = 'switch-title';
    title.textContent = titleEl.textContent.trim();
    row.appendChild(title);

    const descText = descEl ? descEl.textContent.trim() : '';
    if (descText) {
      const more = document.createElement('button');
      more.type = 'button';
      more.className = 'switch-more';
      more.textContent = '?';
      more.setAttribute('aria-label', '查看说明');
      const open = () => field.classList.add('desc-open');
      const close = () => { if (!field.dataset.descPinned) field.classList.remove('desc-open'); };
      more.addEventListener('mouseenter', open);
      more.addEventListener('mouseleave', close);
      more.addEventListener('click', event => {
        event.preventDefault();
        event.stopPropagation();
        const nowPinned = !field.dataset.descPinned;
        if (nowPinned) field.dataset.descPinned = '1'; else delete field.dataset.descPinned;
        field.classList.toggle('desc-open', nowPinned);
      });
      row.appendChild(more);
    }

    const spacer = document.createElement('span');
    spacer.className = 'switch-spacer';
    row.appendChild(spacer);

    field.appendChild(row);
    titleEl.remove();
    if (descText) {
      const desc = document.createElement('div');
      desc.className = 'switch-desc';
      desc.textContent = descText;
      field.appendChild(desc);
    }
    check.remove();
  });
}

function renderGlobalExcludeRules() {
  const box = document.getElementById('global-exclude-chips');
  if (!box) return;
  const cfg = state.config || {};
  if (!Array.isArray(state.globalExcludes)) {
    state.globalExcludes = [...(((cfg.offline || {}).global_exclude) || [])];
  }
  const rules = state.globalExcludes;
  if (!rules.length) {
    box.innerHTML = '<span class="hint">还没有全局排除规则</span>';
  } else {
    box.innerHTML = rules.map((rule, index) => {
      const match = String(rule).match(/^\{\{(.+?)\}\}:(.*)$/);
      const group = match ? match[1] : '';
      const pattern = match ? match[2] : rule;
      return `<span class="global-exclude-chip" title="${esc(rule)}"><i class="global-exclude-dot"></i>${group ? `<b>${esc(group)}</b>` : ''}<code>${esc(pattern)}</code><button type="button" data-global-exclude-remove="${index}" aria-label="删除规则">×</button></span>`;
    }).join('');
    box.querySelectorAll('[data-global-exclude-remove]').forEach(button => button.addEventListener('click', () => {
      state.globalExcludes.splice(Number(button.dataset.globalExcludeRemove), 1);
      renderGlobalExcludeRules();
    }));
  }
  const add = document.getElementById('global-exclude-add');
  if (add && !add.dataset.bound) {
    add.dataset.bound = '1';
    add.addEventListener('click', () => openGlobalExcludeModal());
  }
  const clear = document.getElementById('global-exclude-clear');
  if (clear && !clear.dataset.bound) {
    clear.dataset.bound = '1';
    clear.addEventListener('click', async () => {
      if (!state.globalExcludes.length) return;
      const ok = await confirmDialog({ title: '清空全局排除', message: '确定清空所有全局排除规则吗？', confirmText: '清空' });
      if (!ok) return;
      state.globalExcludes = [];
      renderGlobalExcludeRules();
    });
  }
}

function openGlobalExcludeModal() {
  openModal({
    title: '添加排除规则',
    subtitle: '留空字幕组则对所有字幕组生效',
    size: 'narrow',
    body: `<div class="form-grid">
      <label class="field span-2"><span>字幕组</span><input class="input" id="ge-subgroup" placeholder="留空匹配所有字幕组"></label>
      <label class="field span-2"><span>正则</span><input class="input mono" id="ge-pattern" placeholder="如 720、简、\\d-\\d"><span class="hint">支持正则表达式；任意一条命中即排除。</span></label>
    </div>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="ge-add">${icon('plus')}添加</button>`,
    onMount(el, close) {
      const submit = () => {
        const group = el.querySelector('#ge-subgroup').value.trim();
        const pattern = el.querySelector('#ge-pattern').value.trim();
        if (!pattern) { toast('请填写正则', 'error'); return; }
        const rule = group ? `{{${group}}}:${pattern}` : pattern;
        if (!Array.isArray(state.globalExcludes)) state.globalExcludes = [];
        if (state.globalExcludes.includes(rule)) { toast('这条规则已存在', 'error'); return; }
        state.globalExcludes.push(rule);
        renderGlobalExcludeRules();
        close();
        toast('已添加排除规则，保存配置后生效', 'success');
      };
      el.querySelector('#ge-add').addEventListener('click', submit);
      el.querySelector('#ge-pattern').addEventListener('keydown', event => { if (event.key === 'Enter') submit(); });
      setTimeout(() => el.querySelector('#ge-subgroup').focus(), 60);
    },
  });
}

async function renderSubgroupAliasSettings() {
  const box = document.getElementById('subgroup-alias-settings');
  if (!box) return;
  let aliases = [];
  try {
    aliases = await loadSubgroupAliases(true);
  } catch (err) {
    box.innerHTML = `<div class="form-error">读取字幕组别名失败：${esc(err.message)}</div>`;
    return;
  }
  const render = () => {
    if (!aliases.length) {
      box.innerHTML = '<div class="empty" style="padding:22px"><span class="hint">还没有字幕组别名。</span></div>';
      return;
    }
    box.innerHTML = `<div class="match-template-list">${aliases.map(item => `
      <div class="match-template-item">
        <div class="match-template-main">
          <div class="match-template-title"><strong>${esc(item.main || '未命名')}</strong><span class="badge">主名称</span></div>
          <div class="match-tags">${(item.aliases || []).map(alias => `<span>${esc(alias)}</span>`).join('') || '<span class="hint">暂无别名</span>'}</div>
        </div>
        <div class="match-template-actions">
          <button class="btn btn-sm" type="button" data-alias-edit="${esc(item.main)}">编辑</button>
          <button class="btn btn-sm btn-danger" type="button" data-alias-delete="${esc(item.main)}">删除</button>
        </div>
      </div>`).join('')}</div>`;
    box.querySelectorAll('[data-alias-edit]').forEach(button => button.addEventListener('click', () => {
      const item = aliases.find(value => value.main === button.dataset.aliasEdit);
      if (item) openSubgroupAliasEditor(item);
    }));
    box.querySelectorAll('[data-alias-delete]').forEach(button => button.addEventListener('click', async () => {
      const item = aliases.find(value => value.main === button.dataset.aliasDelete);
      if (!item) return;
      const ok = await confirmDialog({ title: '删除字幕组别名', message: `确定删除「${esc(item.main)}」的别名映射吗？`, confirmText: '删除' });
      if (!ok) return;
      try {
        const next = aliases.filter(value => value.main !== item.main);
        const data = await api('/api/subgroup-aliases', { method: 'PUT', body: { aliases: next } });
        state.subgroupAliases = data.aliases || [];
        aliases = state.subgroupAliases;
        render();
        toast('别名已删除', 'success');
      } catch (err) { toast(err.message, 'error', 5200); }
    }));
  };
  render();
  const add = document.getElementById('subgroup-alias-add');
  if (add && !add.dataset.bound) {
    add.dataset.bound = '1';
    add.addEventListener('click', () => openSubgroupAliasEditor(null));
  }
}

function openSubgroupAliasEditor(item) {
  const editing = !!item;
  const current = item || { main: '', aliases: [] };
  const aliases = current.aliases || [];
  openModal({
    title: editing ? '编辑字幕组别名' : '添加字幕组别名',
    subtitle: 'RSS 中的英文名会自动转换为这里的主名称',
    size: 'narrow',
    body: `<div class="form-grid">
      <label class="field span-2"><span>主名称 *</span><input class="input" id="sa-main" value="${esc(current.main || '')}" placeholder="例如 今晚月色真美"></label>
      <label class="field"><span>别名 1</span><input class="input mono" id="sa-alias1" value="${esc(aliases[0] || '')}" placeholder="例如 Tsukigakirei"></label>
      <label class="field"><span>别名 2</span><input class="input mono" id="sa-alias2" value="${esc(aliases[1] || '')}" placeholder="例如 tsukigakirei"></label>
      <label class="field span-2"><span>更多别名</span><textarea class="textarea mono" id="sa-extra" rows="2" placeholder="每行一个；可选">${esc(aliases.slice(2).join('\n'))}</textarea></label>
    </div>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="sa-save">${icon('save')}保存别名</button>`,
    onMount(el, close) {
      const save = el.querySelector('#sa-save');
      save.addEventListener('click', async () => {
        const main = el.querySelector('#sa-main').value.trim();
        const nextAliases = [el.querySelector('#sa-alias1').value, el.querySelector('#sa-alias2').value, ...el.querySelector('#sa-extra').value.split(/\r?\n/)]
          .map(value => value.trim()).filter(Boolean);
        if (!main) { toast('请填写主名称', 'error'); return; }
        setLoading(save, true, '保存中');
        try {
          const list = [...(state.subgroupAliases || [])];
          const index = list.findIndex(value => value.main === main);
          const next = { main, aliases: nextAliases };
          if (index >= 0) list[index] = next; else list.push(next);
          const data = await api('/api/subgroup-aliases', { method: 'PUT', body: { aliases: list } });
          state.subgroupAliases = data.aliases || [];
          close();
          await renderSubgroupAliasSettings();
          toast(editing ? '别名已更新' : '别名已添加', 'success');
        } catch (err) {
          toast(err.message, 'error', 5200);
          setLoading(save, false);
        }
      });
    },
  });
}

async function renderMatchTemplateSettings() {
  const box = document.getElementById('match-template-settings');
  if (!box) return;
  let templates = [];
  try {
    templates = await loadMatchTemplates(true);
  } catch (err) {
    box.innerHTML = `<div class="form-error">读取模板失败：${esc(err.message)}</div>`;
    return;
  }
  const render = () => {
    if (!templates.length) {
      box.innerHTML = `<div class="empty" style="padding:22px"><span class="hint">还没有匹配模板。添加模板后，添加订阅时可在「匹配」页选择。</span></div>`;
      return;
    }
    box.innerHTML = `<div class="match-template-list">${templates.map(template => `
      <div class="match-template-item">
        <div class="match-template-main">
          <div class="match-template-title"><strong>${esc(template.name || '未命名模板')}</strong><span class="badge">${esc(template.subgroup || '未设置字幕组')}</span><span class="match-template-tags">${(((template.tags || []).length ? template.tags : String(template.name || '').split(/[\s/]+/)).filter(Boolean)).map(tag => `<span>${esc(tag)}</span>`).join('')}</span></div>
          <div class="match-template-rules"><code>匹配：${esc(template.filter || '全部')}</code><code>排除：${esc(template.exclude || '无')}</code></div>
        </div>
        <div class="match-template-actions">
          <button class="btn btn-sm" type="button" data-template-edit="${esc(template.id)}">编辑</button>
          <button class="btn btn-sm btn-danger" type="button" data-template-delete="${esc(template.id)}">删除</button>
        </div>
      </div>`).join('')}</div>`;
    box.querySelectorAll('[data-template-edit]').forEach(button => button.addEventListener('click', () => {
      const template = templates.find(item => item.id === button.dataset.templateEdit);
      if (template) openMatchTemplateEditor(template);
    }));
    box.querySelectorAll('[data-template-delete]').forEach(button => button.addEventListener('click', async () => {
      const template = templates.find(item => item.id === button.dataset.templateDelete);
      if (!template) return;
      const ok = await confirmDialog({ title: '删除匹配模板', message: `确定删除「${esc(template.name || '未命名模板')}」吗？`, confirmText: '删除' });
      if (!ok) return;
      try {
        const next = templates.filter(item => item.id !== template.id);
        const data = await api('/api/match-templates', { method: 'PUT', body: { templates: next } });
        state.matchTemplates = data.templates || [];
        templates = state.matchTemplates;
        render();
        toast('匹配模板已删除', 'success');
      } catch (err) { toast(err.message, 'error', 5200); }
    }));
  };
  render();
  const add = document.getElementById('match-template-add');
  if (add && !add.dataset.bound) {
    add.dataset.bound = '1';
    add.addEventListener('click', () => openMatchTemplateEditor(null));
  }
}

function openMatchTemplateEditor(template) {
  const editing = !!template;
  const item = template || { id: '', subgroup: '', name: '', tags: [], filter: '', exclude: '' };
  openModal({
    title: editing ? '编辑匹配模板' : '添加匹配模板',
    subtitle: '模板按字幕组主名称匹配；RSS 英文名请使用「字幕组别名」映射',
    size: 'narrow',
    body: `<div class="form-grid">
      <label class="field"><span>字幕组名称 *</span><input class="input" id="mt-subgroup" value="${esc(item.subgroup || '')}" placeholder="例如 喵萌奶茶屋"></label>
      <label class="field"><span>模板名称 *</span><input class="input" id="mt-name" value="${esc(item.name || '')}" placeholder="例如 简繁日内封"></label>
      <label class="field span-2"><span>匹配规则</span><textarea class="textarea mono" id="mt-filter" rows="3" placeholder="每行一条正则；可留空">${esc(item.filter || '')}</textarea></label>
      <label class="field span-2"><span>排除规则</span><textarea class="textarea mono" id="mt-exclude" rows="3" placeholder="每行一条正则；可留空">${esc(item.exclude || '')}</textarea></label>
    </div>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-primary" id="mt-save">${icon('save')}保存模板</button>`,
    onMount(el, close) {
      const save = el.querySelector('#mt-save');
      save.addEventListener('click', async () => {
        const next = {
          id: item.id || '',
          subgroup: el.querySelector('#mt-subgroup').value.trim(),
          name: el.querySelector('#mt-name').value.trim(),
          aliases: Array.isArray(item.aliases) ? item.aliases : [],
          tags: Array.isArray(item.tags) ? item.tags : [],
          filter: el.querySelector('#mt-filter').value.trim(),
          exclude: el.querySelector('#mt-exclude').value.trim(),
        };
        if (!next.subgroup || !next.name) { toast('请填写字幕组和模板名称', 'error'); return; }
        setLoading(save, true, '保存中');
        try {
          const templates = [...(state.matchTemplates || [])];
          const index = templates.findIndex(value => value.id === next.id);
          if (index >= 0) templates[index] = next; else templates.push(next);
          const data = await api('/api/match-templates', { method: 'PUT', body: { templates } });
          state.matchTemplates = data.templates || [];
          close();
          await renderMatchTemplateSettings();
          toast(editing ? '匹配模板已更新' : '匹配模板已添加', 'success');
        } catch (err) {
          toast(err.message, 'error', 5200);
          setLoading(save, false);
        }
      });
    },
  });
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
          <div class="settings-account-block">
            <div class="settings-account-actions">
              <button class="btn btn-soft" id="qr-login">${icon('qr')}扫码登录</button>
              <button class="btn btn-danger" id="logout-btn">${icon('logout')}退出 115 登录</button>
            </div>
            <label class="field">
              <span>115 Cookies</span>
              <input class="input mono" id="set-cookies" type="password" autocomplete="off" value="${esc(cfg.auth.cookies || '')}" placeholder="UID=...; CID=...; SEID=...; KID=...">
              <span class="hint">只显示一行，内容以 * 隐藏；留空并保存会清除登录状态。</span>
            </label>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('server')}</span><div><h3>服务与访问</h3><p>Web UI 监听地址、登录保护与代理</p></div></div>
        <div class="section-body">
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
        <div class="section-title"><span class="ico">${icon('info')}</span><div><h3>AI 助手</h3><p>DeepSeek 兼容 API，用于查询和受控操作</p></div><span style="margin-left:auto" class="badge ${(cfg.ai || {}).enabled ? 'badge-success' : ''}">${(cfg.ai || {}).enabled ? '已开启' : '未开启'}</span></div>
        <div class="section-body">
          <label class="field"><span>启用 AI 助手</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-ai-enabled" ${(cfg.ai || {}).enabled ? 'checked' : ''}><span></span></span><span>启用后可在「AI 助手」页面查询；修改/删除操作需要二次确认</span></span></label>
          <div class="form-grid">
            <label class="field"><span>模板</span><select class="select" id="set-ai-provider"><option value="deepseek" selected>DeepSeek</option></select><span class="hint">当前支持 DeepSeek 兼容 Chat Completions API。</span></label>
            <label class="field"><span>权限</span><select class="select" id="set-ai-permission"><option value="readonly" ${(cfg.ai || {}).permission !== 'full' ? 'selected' : ''}>部分权限（只读）</option><option value="full" ${(cfg.ai || {}).permission === 'full' ? 'selected' : ''}>完整权限（操作需二次确认）</option></select></label>
            <label class="field span-2"><span>API 地址</span><input class="input mono" id="set-ai-base" value="${esc((cfg.ai || {}).base_url || 'https://api.deepseek.com')}" placeholder="https://api.deepseek.com"></label>
            <label class="field"><span>模型</span><input class="input mono" id="set-ai-model" value="${esc((cfg.ai || {}).model || 'deepseek-chat')}" placeholder="deepseek-chat"></label>
            <label class="field"><span>API Key</span><input class="input" id="set-ai-key" type="password" value="${esc((cfg.ai || {}).api_key || '')}" placeholder="sk-..."></label>
          </div>
          <div class="notice">${icon('info')}<span>AI 不会直接删除或修改订阅；完整权限下也必须经过页面二次确认。</span></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('database')}</span><div><h3>数据库与缓存</h3><p>SQLite 历史库与 115 请求节流参数</p></div></div>
        <div class="section-body">
          <div class="form-grid">
            <label class="field span-2"><span>数据库路径</span><input class="input mono" id="set-db" value="${esc(cfg.database.path)}"><span class="hint">相对路径基于 config.toml 所在目录解析。</span></label>
            <label class="field span-2"><span>115 同步目录</span><div class="field-row"><input class="input mono" value="${esc(meta.database_path || '')}" readonly><span class="badge">解析后</span></div></label>
            <label class="field"><span>历史保留天数</span><input class="input" id="set-history-retention" type="number" min="30" max="3650" value="${Number((cfg.offline || {}).history_retention_days || 365)}"><span class="hint">仅清理旧 RSS 历史，不删除下载文件；保留待后处理任务。</span></label>
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
          <div class="section-subtitle">存储位置</div>
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
            <div class="field span-2"><span>全局排除规则</span><div class="global-exclude-box"><div class="global-exclude-chips" id="global-exclude-chips"></div><div class="global-exclude-actions"><button class="btn btn-sm" type="button" id="global-exclude-add">${icon('plus')}添加排除规则</button><button class="btn btn-sm btn-ghost" type="button" id="global-exclude-clear">清空</button></div></div><span class="hint">一个标签一条规则，任一命中即跳过；支持正则表达式和 <code>{{字幕组}}:正则</code>，对所有订阅生效。</span></div>
          </div>
          <div class="section-subtitle" style="margin-top:16px">任务行为</div>
          <div class="form-grid">
            <label class="field"><span>失败重试次数</span><input class="input" id="set-retry" type="number" min="0" max="100" value="${Number((cfg.offline || {}).retry || 0)}"><span class="hint">115 提交失败后重试的次数。</span></label>
            <label class="field"><span>跳过 x.5 集</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-skip-x5" ${(cfg.offline || {}).skip_x5 ? 'checked' : ''}><span></span></span><span>跳过 12.5 这类集数</span></span></label>
            <label class="field span-2"><span>运行后自动清理</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-clear" ${(cfg.offline || {}).auto_clear_completed ? 'checked' : ''}><span></span></span><span>全部/自动运行结束后，清理 115 已完成离线任务（不删除文件）</span></span></label>
            <label class="field span-2"><span>订阅完成后自动停用</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-disable-completed" ${(cfg.offline || {}).auto_disable_completed ? 'checked' : ''}><span></span></span><span>达到总集数后自动标记完成并改为未启用</span></span></label>
            <label class="field span-2"><span>离线后转移到保存文件夹</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-dual-folder" ${(cfg.offline || {}).dual_folder ? 'checked' : ''}><span></span></span><span>离线完成后（如启用重命名则先重命名）移动到默认保存文件夹</span></span></label>
            <label class="field span-2"><span>保存文件夹只保留视频和字幕</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-keep-only-video" ${(cfg.offline || {}).keep_only_video ? 'checked' : ''}><span></span></span><span>离线后转移到保存文件夹下，只把视频和字幕移动到保存文件夹；字体/其他文件留在离线目录，字幕自动追加 .sc/.tc 等语言后缀</span></span></label>
            <label class="field span-2"><span>同一 RSS 内每集只选一个版本</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-one-per-episode" ${(cfg.offline || {}).one_per_episode ? 'checked' : ''}><span></span></span><span>关闭时同一 RSS 内所有版本都会提交；开启后只提交评分最高的一个，不影响主 RSS 与备用 RSS 之间各自保留版本</span></span></label>
            <label class="field span-2"><span>优先选择 v 数更大的版本</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-prefer-higher-version" ${(cfg.offline || {}).prefer_higher_version ? 'checked' : ''}><span></span></span><span>开启后同集 01v2 会优先于 01，即使 01 的画质或字幕评分更高；仅在「每集只选一个版本」开启时生效</span></span></label>
          </div>
        </div>
      </section>



      <section class="card">
        <div class="section-title"><span class="ico">${icon('download')}</span><div><h3>qBittorrent</h3><p>本地下载器：115 失败后的最后兜底</p></div><button class="help-dot" type="button" id="qb-settings-help" title="查看 qBittorrent 设置说明" style="margin-left:auto">?</button><span class="badge ${(cfg.downloader || {}).qbittorrent && cfg.downloader.qbittorrent.enabled ? 'badge-success' : ''}">${(cfg.downloader || {}).qbittorrent && cfg.downloader.qbittorrent.enabled ? '已开启' : '未开启'}</span></div>
        <div class="section-body">
          <label class="field"><span>启用 qBittorrent</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-qb-enabled" ${(cfg.downloader || {}).qbittorrent && cfg.downloader.qbittorrent.enabled ? 'checked' : ''}><span></span></span><span>115 失败达到阈值后提交本地 qBittorrent</span></span></label>
          <div class="form-grid">
            <label class="field span-2"><span>WebUI 地址</span><input class="input mono" id="set-qb-url" value="${esc(((cfg.downloader || {}).qbittorrent || {}).url || '')}" placeholder="http://127.0.0.1:8080"></label>
            <label class="field"><span>用户名</span><input class="input" id="set-qb-user" value="${esc(((cfg.downloader || {}).qbittorrent || {}).username || '')}" placeholder="admin"></label>
            <label class="field"><span>密码</span><input class="input" id="set-qb-pass" type="password" value="${esc(((cfg.downloader || {}).qbittorrent || {}).password || '')}"></label>
            <label class="field span-2"><span>API Key（qB 5.2+，优先于用户名密码）</span><input class="input mono" id="set-qb-api-key" type="password" value="${esc(((cfg.downloader || {}).qbittorrent || {}).api_key || '')}" placeholder="qbt_..."></label>
            <label class="field"><span>保存路径</span><input class="input mono" id="set-qb-save" value="${esc(((cfg.downloader || {}).qbittorrent || {}).save_path || '')}" placeholder="留空使用 qB 默认路径"></label>
            <label class="field"><span>分类</span><input class="input" id="set-qb-category" value="${esc(((cfg.downloader || {}).qbittorrent || {}).category || '')}" placeholder="anifeed"></label>
            <label class="field"><span>标签</span><input class="input" id="set-qb-tags" value="${esc(((cfg.downloader || {}).qbittorrent || {}).tags || '')}" placeholder="anifeed,115-fallback"></label>
            <label class="field"><span>内容布局</span><select class="select" id="set-qb-layout">${["", "Original", "Subfolder", "NoSubfolder"].map(v => `<option value="${v}" ${(((cfg.downloader || {}).qbittorrent || {}).content_layout || "") === v ? "selected" : ""}>${v || "默认"}</option>`).join("")}</select></label>
            <label class="field"><span>下载限速 KB/s</span><input class="input" id="set-qb-dl-limit" type="number" min="0" value="${Number(((cfg.downloader || {}).qbittorrent || {}).download_limit_kb || 0)}"></label>
            <label class="field"><span>上传限速 KB/s</span><input class="input" id="set-qb-up-limit" type="number" min="0" value="${Number(((cfg.downloader || {}).qbittorrent || {}).upload_limit_kb || 0)}"></label>
            <label class="field"><span>分享率</span><input class="input" id="set-qb-ratio" type="number" value="${Number(((cfg.downloader || {}).qbittorrent || {}).ratio_limit || 0)}"><span class="hint">0 使用 qB 全局设置。</span></label>
            <label class="field"><span>做种时长/分钟</span><input class="input" id="set-qb-seed-time" type="number" value="${Number(((cfg.downloader || {}).qbittorrent || {}).seeding_time_limit || 0)}"></label>
            <label class="field"><span>非活跃时长/分钟</span><input class="input" id="set-qb-inactive" type="number" value="${Number(((cfg.downloader || {}).qbittorrent || {}).inactive_seeding_time_limit || 0)}"></label>
            <label class="field"><span>使用 qB 临时下载路径</span><span class="check" style="height:32px"><span class="switch"><input type="checkbox" id="set-qb-use-download-path" ${(((cfg.downloader || {}).qbittorrent || {}).use_download_path) ? "checked" : ""}><span></span></span><span>最终路径仍使用保存路径</span></span></label>
            <label class="field"><span>115 失败阈值</span><input class="input" id="set-qb-threshold" type="number" min="1" max="20" value="${Number(((cfg.downloader || {}).qbittorrent || {}).fallback_threshold || 3)}"><span class="hint">115 连续失败达到次数后触发 qBittorrent。</span></label>
          </div>
          <div class="field-row"><button class="btn btn-sm" id="qb-test" type="button">${icon('link')}测试连接</button><span class="hint">测试会登录 qBittorrent WebUI 并读取版本。</span></div>
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
          <div class="field-row"><button class="btn btn-sm" id="postprocess-now">${icon('refresh')}立即处理已完成任务</button><span class="hint">处理 115 离线目录中已完成任务的文件/文件夹；后台每 20 秒自动处理；当前待处理：${state.status ? Number(state.status.postprocess_pending || 0) : 0}</span></div>
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
            <label class="field span-2"><span>推送内容</span><div class="ms" data-ms="content"><button type="button" class="ms-trigger" id="ms-content-trigger"><span class="ms-label" id="ms-content-label">选择内容</span><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M7 10l5 5 5-5"/></svg></button><div class="ms-menu" id="ms-content-menu">
              <label class="check-inline"><input type="checkbox" id="notify-content-summary" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('summary') ? 'checked' : ''}>运行摘要</label>
              <label class="check-inline"><input type="checkbox" id="notify-content-details" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('details') ? 'checked' : ''}>订阅明细</label>
              <label class="check-inline"><input type="checkbox" id="notify-content-errors" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('errors') ? 'checked' : ''}>错误信息</label>
              <label class="check-inline"><input type="checkbox" id="notify-content-time" ${((cfg.notify || {}).content || ['summary','details','errors','time']).includes('time') ? 'checked' : ''}>时间</label>
            </div></div><span class="hint">通知以纯文本发送；企业微信使用 text 消息，通用 Webhook 使用 text/plain。</span></label>
          </div>
          <label class="field span-2"><span>推送事件</span><div class="ms" data-ms="events"><button type="button" class="ms-trigger" id="ms-events-trigger"><span class="ms-label" id="ms-events-label">选择事件</span><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M7 10l5 5 5-5"/></svg></button><div class="ms-menu" id="ms-events-menu">
            ${[['run','运行汇总'],['added','添加至云下载'],['completed','下载完成'],['failed','下载失败'],['idle','摸鱼提醒'],['hourly','整点更新'],['wash','洗版完成']].map(([key,label]) => `<label class="check-inline"><input type="checkbox" class="notify-event" data-event="${key}" ${notifyEvents.includes(key) ? 'checked' : ''}>${label}</label>`).join('')}
          </div></div><span class="hint">单独控制每种事件是否推送；关闭后该事件不会发送。</span></label>
          <div class="form-grid">
            <label class="field span-2"><span>摸鱼提醒</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="notify-idle-enabled" ${(cfg.notify || {}).idle_enabled ? 'checked' : ''}><span></span></span><span>订阅超过指定天数没有更新时发送提醒</span></span></label>
            <label class="field"><span>未更新天数</span><input class="input" id="notify-idle-days" type="number" min="1" max="3650" value="${Number((cfg.notify || {}).idle_days || 7)}"><span class="hint">每天最多提醒一次。</span></label>
          </div>
          <label class="field span-2"><span class="label-help">通知模板 <button class="help-dot" type="button" id="notify-template-preview" title="预览通知">?</button></span><textarea class="textarea mono" id="set-notify-template" rows="13" spellcheck="false">${esc((cfg.notify || {}).template || DEFAULT_NOTIFY_TEMPLATE)}</textarea><div class="field-row"><button class="btn btn-sm" type="button" id="notify-template-reset">恢复默认模板</button></div><span class="hint">变量：\${action} \${title} \${score} \${themoviedbName} \${bgmUrl} \${season} \${episode} \${subgroup} \${currentEpisodeNumber} \${totalEpisodeNumber} \${year} \${month} \${date} \${text} \${bgmEpisodeTitle}</span></label>
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
        <div class="section-title"><span class="ico">${icon('info')}</span><div><h3>关于</h3><p>AniFeed 项目信息</p></div></div>
        <div class="section-body">
          <div class="about-head"><span class="brand-mark"><img src="/logo.png" alt="AniFeed"></span><div><h2>AniFeed</h2><p>RSS → 115 离线 · 自托管动漫下载管理</p></div></div>
          <div class="kv">
            <div class="kv-row"><span>版本</span><code>v${esc(state.status ? state.status.version : '0.02')}</code></div>
            <div class="kv-row"><span>GitHub</span><a href="https://github.com/huohen92/anifeed" target="_blank" rel="noopener noreferrer">github.com/huohen92/anifeed</a></div>
            <div class="kv-row"><span>许可</span><code>MIT</code></div>
          </div>
          <div class="notice" style="margin-top:14px">${icon('info')}<span>115 为主下载器；qBittorrent 为本地兜底下载器。</span></div>
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
        <div class="section-title"><span class="ico">${icon('heart')}</span><div><h3>喜好字幕组</h3><p>用于洗版保留顺序和自动备用字幕组</p></div></div>
        <div class="section-body">
          <div class="field span-2"><span>字幕组优先级</span>
            <div id="lab-subgroup-priority" class="subgroup-priority"></div>
            <div class="field-row subgroup-add-row"><input class="input" id="lab-subgroup-add" placeholder="手动添加字幕组名称"><button class="btn btn-sm" type="button" id="lab-subgroup-add-btn">${icon('plus')}添加</button></div>
            <span class="hint">左侧从上到下优先；用于洗版保留和自动备用字幕组。可把字幕组移到右侧不启用。</span>
            <div class="notice" style="margin-top:8px">${icon('info')}<span>喜好字幕组依赖重命名后的标题包含字幕组名称。请确保重命名模板包含 <code>\${subgroup}</code>，否则洗版和自动备用无法识别字幕组。</span></div>
            ${renameHasSubgroup ? '' : `<div class="notice" style="margin-top:8px">${icon('alert')}<span>当前重命名模板不包含 <code>\${subgroup}</code>，喜好字幕组可能无法生效。</span></div>`}
          </div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('trash')}</span><div><h3>洗版</h3><p>按下载记录清理多个文件夹中的重复剧集</p></div></div>
        <div class="section-body">
          <label class="field span-2"><span>启用洗版</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-wash" ${(cfg.lab || {}).wash_enabled ? 'checked' : ''}><span></span></span><span>优先保留喜好字幕组；都不在喜好列表时保留简体且体积更大的版本</span></span></label>
          <label class="field"><span>洗版模式</span><select class="select" id="set-lab-wash-mode">
            <option value="safe" ${((cfg.lab || {}).wash_mode || 'safe') === 'safe' ? 'selected' : ''}>安全模式：只处理有下载记录的文件</option>
            <option value="full" ${((cfg.lab || {}).wash_mode || '') === 'full' ? 'selected' : ''}>完整模式：扫描文件夹内所有可识别剧集</option>
          </select><span class="hint">安全模式不会碰没有下载记录的文件。</span></label>
          <div class="form-grid">
            <label class="field"><span>定时洗版</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-wash-schedule" ${(cfg.lab || {}).wash_schedule_enabled ? 'checked' : ''}><span></span></span><span>按间隔自动执行洗版</span></span></label>
            <label class="field"><span>间隔（分钟）</span><input class="input" id="set-lab-wash-interval" type="number" min="1" max="100000" value="${Number((cfg.lab || {}).wash_interval_minutes || 1440)}"><span class="hint">默认 1440 分钟，即每天一次。</span></label>
          </div>
          <div class="field span-2"><span>洗版文件夹</span><div class="lab-folder-list" id="lab-wash-folders"></div><div class="field-row" style="margin-top:8px"><button class="btn btn-sm" type="button" id="lab-add-wash-folder">${icon('plus')}添加文件夹</button><button class="btn btn-sm btn-danger" type="button" id="lab-wash-now">${icon('trash')}立即洗版</button><span class="hint">先扫描预览，再确认删除；会发送逐文件删除通知。</span></div></div>
        </div>
      </section>

      <section class="card">
        <div class="section-title"><span class="ico">${icon('refresh')}</span><div><h3>自动备用字幕组</h3><p>主 RSS 连续添加失败，或摸鱼超时后自动补充字幕组</p></div><button class="help-dot" type="button" id="auto-backup-logic-help" title="查看完整逻辑" style="margin-left:auto">?</button></div>
        <div class="section-body">
          <label class="field span-2"><span>自动添加备用字幕组</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-backup" ${(cfg.lab || {}).auto_backup_subgroup ? 'checked' : ''}><span></span></span><span>主 RSS 连续多次添加失败后，自动补充有缺失集数的备用字幕组 RSS</span></span></label>
          <label class="field"><span>主 RSS 连续添加失败次数</span><input class="input" id="set-lab-auto-backup-failures" type="number" min="1" max="20" value="${Number((cfg.lab || {}).auto_backup_failures || 2)}"><span class="hint">主 RSS 连续添加失败达到该次数后触发检查。</span></label>
          <label class="field span-2"><span>摸鱼后自动补充备用 RSS</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-backup-idle" ${(cfg.lab || {}).auto_backup_idle ? 'checked' : ''}><span></span></span><span>订阅超过「通知 → 摸鱼提醒天数」没有更新时，自动检查缺失集数并补充备用 RSS，不要求配置 Webhook</span></span></label>
          <label class="field span-2"><span>跳过使用集数偏移的番剧</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-backup-skip-offset" ${(cfg.lab || {}).auto_backup_skip_offset ? 'checked' : ''}><span></span></span><span>主 RSS 已设置集数偏移时不自动添加备用 RSS；备用 RSS 本身按 TMDB 顺序，不需要继承偏移</span></span></label>
          <label class="field span-2"><span>自动校正备用 RSS 偏移</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-offset" ${(cfg.lab || {}).auto_offset ? 'checked' : ''}><span></span></span><span>用 TMDB 分集 + 发布日期匹配，自动计算备用 RSS 的季数和偏移</span></span></label>
          <div class="field span-2"><div class="notice">${icon('info')}<span>优先级：先执行「跳过使用集数偏移的番剧」；跳过命中后不会自动添加备用 RSS，也不会计算偏移。未跳过时才由「自动校正备用 RSS 偏移」计算季数和偏移。</span></div></div>
          <label class="field"><span>偏移匹配样本数</span><div class="field-row"><input class="input" id="set-lab-auto-offset-samples" type="number" min="1" max="20" value="${Number((cfg.lab || {}).auto_offset_samples || 3)}"><button class="help-dot" type="button" id="auto-offset-samples-help" title="查看说明">?</button></div></label>
          <label class="field span-2"><span>自动应用规则</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-auto-apply" ${(cfg.lab || {}).auto_apply_rules ? 'checked' : ''}><span></span></span><span>自动添加的备用字幕组使用自动应用最佳规则</span></span></label>
          <label class="field span-2"><span>显示表单提示</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-lab-hide-hints" ${!(cfg.lab || {}).hide_hints ? 'checked' : ''}><span></span></span><span>关闭后隐藏字段说明文字与 ? / ! 帮助入口（给熟练用户的清爽模式）</span></span></label>
          <div class="notice">${icon('info')}<span>自动备用相关选项已接入运行流程；摸鱼触发使用「通知 → 摸鱼提醒天数」，不需要配置 Webhook。</span></div>
        </div>
      </section>
      <section class="card">
        <div class="section-title"><span class="ico">${icon('filter')}</span><div><h3>字幕组匹配模板</h3><p>添加订阅时按字幕组快速选择匹配和排除规则</p></div><button class="btn btn-sm" type="button" id="match-template-add">${icon('plus')}添加模板</button></div>
        <div class="section-body">
          <div class="notice">${icon('info')}<span>模板按字幕组名称匹配；添加订阅时选择「匹配」标签页即可套用。选择「全部」不添加规则。</span></div>
          <div id="match-template-settings" class="match-template-settings"><span class="hint">正在读取模板…</span></div>
        </div>
      </section>
      <section class="card">
        <div class="section-title"><span class="ico">${icon('filter')}</span><div><h3>自动应用最佳规则</h3><p>语言、字幕、容器和分辨率的自动选择顺序</p></div><span style="margin-left:auto" class="badge ${(cfg.offline || {}).auto_rules ? 'badge-success' : ''}">${(cfg.offline || {}).auto_rules ? '已开启' : '未开启'}</span></div>
        <div class="section-body">
          <label class="field span-2"><span>使用自动规则</span><span class="check" style="height:36px"><span class="switch"><input type="checkbox" id="set-auto-rules" ${(cfg.offline || {}).auto_rules ? 'checked' : ''}><span></span></span><span>打开「规则与匹配」时，自动按下方优先级选择最佳规则（可再手动调整）</span></span></label>
          <div class="hint" style="margin:8px 0">每个分类都可以单独添加、移除和排序；当前 RSS 不存在的选项会自动跳过。</div>
          <div id="priority-editor"></div>
        </div>
      </section>
      <section class="card">
        <div class="section-title"><span class="ico">${icon('users')}</span><div><h3>字幕组别名</h3><p>把 RSS 中的英文名映射到中文字幕组主名称</p></div><button class="btn btn-sm" type="button" id="subgroup-alias-add">${icon('plus')}添加别名</button></div>
        <div class="section-body">
          <div class="notice">${icon('info')}<span>添加订阅时先读取 RSS 标题；标题无法识别时按 URL 的 groupSlug 查询别名。匹配到别名后，字幕组字段和规则都会使用主名称。</span></div>
          <div id="subgroup-alias-settings" class="match-template-settings"><span class="hint">正在读取别名…</span></div>
        </div>
      </section>
      <section class="card">
        <div class="section-title"><span class="ico">${icon('image')}</span><div><h3>个性化</h3><p>背景、玻璃效果、强调色与导航配色</p></div></div>
        <div class="section-body">
          <div class="notice">${icon('info')}<span>视觉效果保存在本机浏览器；背景图片上传到服务器，可分别管理电脑端和移动端，修改即时生效。</span></div>
          <div class="af-blocks">
          <div class="af-block">
            <div class="af-block-title">背景图片 <small>上传到服务器；电脑端和移动端分别管理</small></div>
            <div class="af-device-tabs" id="af-bg-device">
              <button type="button" data-device="desktop">电脑端</button>
              <button type="button" data-device="mobile">移动端</button>
            </div>
            <div class="af-row">
              <input class="input" id="af-bg-url" placeholder="https://example.com/bg.jpg 或 API 地址">
              <button class="btn btn-sm" id="af-bg-apply" type="button">应用</button>
              <button class="btn btn-sm" id="af-bg-upload" type="button">${icon('upload')}批量上传</button>
              <button class="btn btn-sm btn-danger" id="af-bg-clear" type="button">清除选择</button>
              <input type="file" id="af-bg-file" accept="image/*" multiple hidden>
            </div>
            <div class="af-row af-bg-mode-row">
              <select class="select" id="af-bg-mode"><option value="single">显示指定图片</option><option value="none">不显示背景</option></select>
              <span class="af-hint" id="af-bg-count">0 张图片</span>
            </div>
            <div class="af-bg-gallery" id="af-bg-gallery"></div>
            <div class="af-bg-preview" id="af-bg-preview">未设置背景</div>
            <div class="af-hint">图片保存在服务器；电脑端和移动端可使用不同背景。建议单张不超过 10MB。</div>
            <label class="field" style="margin-top:10px"><span>图片模糊 <span class="af-slider-value" id="af-bg-blur-val">8px</span></span><input class="af-slider" id="af-bg-blur" type="range" min="0" max="30" step="1"></label>
          </div>
          <div class="af-block">
            <div class="af-block-title">玻璃效果 <small>设置背景图后自动启用</small></div>
            <div class="af-grid">
              <label class="field"><span>面板模糊 <span class="af-slider-value" id="af-blur-val">14px</span></span><input class="af-slider" id="af-blur" type="range" min="0" max="30" step="1"></label>
              <label class="field"><span>面板不透明度 <span class="af-slider-value" id="af-alpha-val">72%</span></span><input class="af-slider" id="af-alpha" type="range" min="35" max="100" step="1"></label>
            </div>
          </div>
          <div class="af-block">
            <div class="af-block-title">圆角与悬浮 <small>卡片、按钮、输入框统一生效</small></div>
            <div class="af-grid">
              <label class="field"><span>圆角 <span class="af-slider-value" id="af-radius-val">16px</span></span><input class="af-slider" id="af-radius" type="range" min="4" max="28" step="1"></label>
              <label class="field"><span>悬浮阴影 <span class="af-slider-value" id="af-elev-val">45%</span></span><input class="af-slider" id="af-elev" type="range" min="0" max="100" step="1"></label>
            </div>
          </div>
          <div class="af-block">
            <div class="af-block-title">导航配色 <small>侧边栏底色，和内容区区分开</small></div>
            <div class="af-swatches" id="af-nav"></div>
          </div>
          <div class="af-block">
            <div class="af-block-title">强调色 <small>链接、激活图标与重点标记</small></div>
            <div class="af-swatches" id="af-accent"></div>
          </div>
          <div class="af-block">
            <div class="af-block-title">主题模式</div>
            <div class="af-seg" id="af-mode">
              <button type="button" data-mode="light">浅色</button>
              <button type="button" data-mode="dark">深色</button>
              <button type="button" data-mode="auto">跟随系统</button>
            </div>
          </div>
          <div class="af-block">
            <div class="af-block-title">导入 / 导出 <small>主题配置可带到其他设备</small></div>
            <div class="af-row">
              <button class="btn btn-sm" id="af-export-2" type="button">${icon('download')}导出 JSON</button>
              <button class="btn btn-sm" id="af-import" type="button">${icon('upload')}导入 JSON</button>
              <input type="file" id="af-import-file" accept="application/json,.json" hidden>
              <button class="btn btn-sm btn-danger" id="af-reset" type="button">恢复默认</button>
            </div>
            <div class="af-hint" id="af-io-hint">「导出 JSON」导出的是主题 JSON（背景图、模糊、配色、导航色、明暗模式），与下方「导入 JSON」对应。</div>
            <div class="af-row" style="margin-top:8px"><button class="btn btn-primary btn-sm" id="af-save" type="button">${icon('save')}保存个性化</button><span class="af-hint">修改会即时生效，并保存在本机浏览器</span></div>
          </div>
          </div>
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
          <button data-settings-section="download">下载与通知</button>
          <button data-settings-section="organize">整理与命名</button>
          <button data-settings-section="subgroups">字幕组与规则</button>
          <button data-settings-section="personal">外观与数据</button>
          <button data-settings-section="about">系统与关于</button>
        </div>
        <div class="settings-section-intro" id="settings-section-intro"></div>` : ''}
        <div id="settings-content">${state.settingsTab === 'form' ? formTab : rawTab}</div>
        <div class="sticky-save" id="settings-savebar">
          <span class="spacer"></span>
          <button class="btn" id="settings-cancel">放弃修改</button>
          <button class="btn btn-primary" id="settings-save-2">${icon('save')}保存配置</button>
        </div>
      </div>
    </div>`;

  initPriorityEditor(cfg);
  initSettingsSections();
  initMultiSelects();
  enhanceSwitchRows();
  initPersonalization();
  hydrateFolderNames();
  renderLabSubgroupPriority();
  renderLabWashFolders();
  renderMatchTemplateSettings();
  renderSubgroupAliasSettings();
  renderGlobalExcludeRules();


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
    $('#qr-login').addEventListener('click', openQrModal);
    $('#logout-btn').addEventListener('click', logout);
    $('#notify-test').addEventListener('click', testNotification);

    $('#auto-backup-logic-help')?.addEventListener('click', () => openModal({
      title: '自动备用 RSS 完整逻辑',
        size: 'narrow',
      subtitle: '从失败检测到 Season + Offset',
      body: `<div class="kv">
        <div class="kv-row"><span>1. 总开关</span><code>设置中的自动添加备用字幕组</code></div>
        <div class="kv-row"><span>2. 单订阅</span><code>启用自动加入备用 RSS</code></div>
        <div class="kv-row"><span>3. 触发</span><code>主 RSS 连续添加 115 失败，或摸鱼达到指定天数</code></div>
        <div class="kv-row"><span>4. 缺失集</span><code>对比历史记录寻找缺失集</code></div>
        <div class="kv-row"><span>5. 字幕组</span><code>匹配喜好字幕组，且该组近期有缺失集</code></div>
        <div class="kv-row"><span>6. 编号判断</span><code>按季编号 / 绝对编号；超出季长度时按绝对编号</code></div>
        <div class="kv-row"><span>7. 季结构</span><code>AniList 季集数 + TMDB 季结构</code></div>
        <div class="kv-row"><span>8. TMDB 对应季</span><code>直接确定季，Offset 0</code></div>
        <div class="kv-row"><span>9. TMDB 单季</span><code>AniList 多季时按累计集数计算 Offset</code></div>
        <div class="kv-row"><span>10. 后备</span><code>无法直接映射时按发布日期匹配 TMDB，优先 ANi 搬运组</code></div>
        <div class="kv-row"><span>11. 校验</span><code>样本数达到设置值且一致性 ≥60%</code></div>
        <div class="kv-row"><span>12. 写入</span><code>只写入新添加的备用 RSS，不修改主 RSS/已有备用 RSS</code></div>
      </div><div class="hint" style="margin-top:10px">「跳过使用集数偏移的番剧」优先级最高；命中后不会添加备用 RSS，也不会执行偏移计算。摸鱼触发使用「通知 → 摸鱼提醒天数」，并且不需要配置 Webhook。</div>`,
      footer: '<button class="btn btn-primary" data-action="modal-close">知道了</button>',
    }));
    $('#auto-offset-samples-help')?.addEventListener('click', () => openModal({
      title: '偏移匹配样本数说明',
        size: 'narrow',
      subtitle: '不是下载次数，也不是重试次数',
      body: `<div class="notice">${icon('info')}<span>样本数指：有多少条 RSS 条目成功匹配到 TMDB，并且计算出的「TMDB 季 + 偏移」完全一致。</span></div>
        <div class="kv" style="margin-top:12px">
          <div class="kv-row"><span>1 条条目</span><code>解析集数 + 发布日期 → 匹配一集 TMDB</code></div>
          <div class="kv-row"><span>候选偏移</span><code>TMDB 集数 - RSS 集数</code></div>
          <div class="kv-row"><span>样本</span><code>同一个「季 + 偏移」出现的条目数</code></div>
        </div>
        <div class="hint" style="margin-top:10px">示例：4 条 RSS 都得到 Season 1 / Offset +14，则样本数为 4；达到设置值且一致性 ≥60% 后，才会把 Season 和 Offset 写入新添加的备用 RSS。样本不足时仍会添加备用 RSS，但 Offset 保持 0。</div>`,
      footer: '<button class="btn btn-primary" data-action="modal-close">知道了</button>',
    }));
    $('#qb-settings-help')?.addEventListener('click', openQBHelp);
    $('#notify-template-preview')?.addEventListener('click', openNotifyPreview);
    $('#notify-template-reset')?.addEventListener('click', () => { const el = $('#set-notify-template'); if (el) el.value = DEFAULT_NOTIFY_TEMPLATE; });
    $('#qb-test')?.addEventListener('click', async () => {
      const btn = $('#qb-test'); btn.disabled = true; btn.textContent = '连接中…';
      try { const r = await api('/api/downloaders/qb/test', { method: 'POST', body: {} }); toast('qBittorrent 已连接：' + (r.version || 'OK'), 'success'); }
      catch (e) { toast(e.message, 'error', 5000); }
      finally { btn.disabled = false; btn.innerHTML = `${icon('link')}测试连接`; }
    });
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
      }, '选择洗版文件夹');
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
    auto_backup_failures: Number($('#set-lab-auto-backup-failures').value || 2),
    auto_backup_skip_offset: $('#set-lab-auto-backup-skip-offset').checked,
    auto_backup_idle: $('#set-lab-auto-backup-idle').checked,
    auto_offset: $('#set-lab-auto-offset').checked,
    auto_offset_samples: Number($('#set-lab-auto-offset-samples').value || 3),
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
    auto_backup_failures: Number($('#set-lab-auto-backup-failures').value || 2),
    auto_backup_skip_offset: $('#set-lab-auto-backup-skip-offset').checked,
    auto_backup_idle: $('#set-lab-auto-backup-idle').checked,
    auto_offset: $('#set-lab-auto-offset').checked,
    auto_offset_samples: Number($('#set-lab-auto-offset-samples').value || 3),
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
  box.innerHTML = state.labWashFolders.map((item, index) => `<span class="lab-folder-chip"><span>${esc(item.name || '已选择文件夹')}</span><button type="button" data-lab-folder-remove="${index}" title="移除">${icon('close')}</button></span>`).join('') || '<span class="hint">还没有添加洗版文件夹</span>';
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
    toast('请先保存实验室配置，再执行洗版', 'error');
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
      title: '洗版预览',
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
            toast(`洗版完成：删除 ${result.deleted || 0} 个文件，失败 ${result.failed || 0} 个`, result.failed ? 'error' : 'success', 5200);
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
  cfg.downloader = {
    primary_115: {
      enabled: true,
      offline_cid: $('#set-default-cid').value.trim(),
      save_path: '',
    },

    qbittorrent: {
      enabled: $('#set-qb-enabled').checked,
      url: $('#set-qb-url').value.trim(),
      username: $('#set-qb-user').value.trim(),
      password: $('#set-qb-pass').value,
      api_key: $('#set-qb-api-key').value.trim(),
      tags: $('#set-qb-tags').value.trim(),
      content_layout: $('#set-qb-layout').value,
      download_limit_kb: Number($('#set-qb-dl-limit').value || 0),
      upload_limit_kb: Number($('#set-qb-up-limit').value || 0),
      ratio_limit: Number($('#set-qb-ratio').value || 0),
      seeding_time_limit: Number($('#set-qb-seed-time').value || 0),
      inactive_seeding_time_limit: Number($('#set-qb-inactive').value || 0),
      use_download_path: $('#set-qb-use-download-path').checked,
      save_path: $('#set-qb-save').value.trim(),
      category: $('#set-qb-category').value.trim(),
      fallback_threshold: Number($('#set-qb-threshold').value || 3),
    },
  };
  cfg.offline = {
    default_cid: $('#set-default-cid').value.trim(),
    default_savepath: '',
    global_exclude: [...(state.globalExcludes || [])],
    skip_x5: $('#set-skip-x5').checked,
    retry: Number($('#set-retry').value || 0),
    auto_clear_completed: $('#set-auto-clear').checked,
    auto_disable_completed: $('#set-auto-disable-completed').checked,
    one_per_episode: $('#set-one-per-episode').checked,
    prefer_higher_version: $('#set-prefer-higher-version').checked,
    keep_only_video: $('#set-keep-only-video').checked,
    auto_language: state.autoPriority ? state.autoPriority.language : PRIORITY_DEFAULTS.language,
    auto_subtitle: state.autoPriority ? state.autoPriority.subtitle : PRIORITY_DEFAULTS.subtitle,
    auto_container: state.autoPriority ? state.autoPriority.container : PRIORITY_DEFAULTS.container,
    auto_resolution: state.autoPriority ? state.autoPriority.resolution : PRIORITY_DEFAULTS.resolution,
    auto_custom: state.autoPriority && state.autoPriority.custom ? Object.entries(state.autoPriority.custom).flatMap(([key, values]) => (values || []).map(value => key + '|' + value)) : [],
    auto_rules: $('#set-auto-rules').checked,
    dual_folder: $('#set-dual-folder').checked,
    final_cid: $('#set-final-cid').value.trim(),
    final_mode: 1,
    final_cid2: $('#set-final-cid2').value.trim(),
    final_path_template: $('#set-final-path-template').value.trim() || '${month}/${titleYearTmdb}/${animeSeason}',
    history_retention_days: Number($('#set-history-retention').value || 365),
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
  cfg.ai = {
    enabled: $('#set-ai-enabled').checked,
    provider: $('#set-ai-provider').value,
    base_url: $('#set-ai-base').value.trim(),
    model: $('#set-ai-model').value.trim(),
    api_key: $('#set-ai-key').value.trim(),
    permission: $('#set-ai-permission').value,
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
    auto_backup_failures: Number($('#set-lab-auto-backup-failures').value || 2),
    auto_backup_skip_offset: $('#set-lab-auto-backup-skip-offset').checked,
    auto_backup_idle: $('#set-lab-auto-backup-idle').checked,
    auto_offset: $('#set-lab-auto-offset').checked,
    auto_offset_samples: Number($('#set-lab-auto-offset-samples').value || 3),
    auto_apply_rules: $('#set-lab-auto-apply').checked,
    hide_hints: !$('#set-lab-hide-hints').checked,
  };
  cfg.sites = (state.config && state.config.sites) || {};
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
      if (newToken) sessionStorage.setItem('anifeed_token', newToken);
      else sessionStorage.removeItem('anifeed_token');
    }
    toast('配置已保存并生效', 'success');
    await loadStatus().catch(err => console.debug('[AniFeed]', err));
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
    await loadStatus().catch(err => console.debug('[AniFeed]', err));
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
      if (newToken) sessionStorage.setItem('anifeed_token', newToken);
      else sessionStorage.removeItem('anifeed_token');
    }
    state.dirty = false;
    toast('config.toml 已保存并重载', 'success');
    await loadStatus().catch(err => console.debug('[AniFeed]', err));
    renderSettingsPage();
  } catch (err) {
    toast(err.message, 'error', 6000);
  } finally {
    setLoading(btn, false);
  }
}


window.AniFeedSettings = { loadPage: loadSettingsPage };
})();
