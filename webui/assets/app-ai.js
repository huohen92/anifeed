/* Lazily loaded AI page module. */
(function () {
'use strict';
function loadAIPage() {
  const view = $('#view');
  const welcome = '你好，我可以查询订阅、115 任务、下载历史和运行状态。修改订阅需要二次确认。';
  if (!state.aiMessages) state.aiMessages = [{ role: 'assistant', content: welcome }];
  const suggestions = ['现在有哪些订阅？', '最近失败的 115 任务', '最近下载了什么？', '备用 RSS 有哪些？'];

  const bubble = m => {
    const isUser = m.role === 'user';
    return `<div class="chat-row ${isUser ? 'is-user' : 'is-ai'}">
      ${isUser ? '' : '<span class="chat-avatar">AI</span>'}
      <div class="chat-bubble">${esc(m.content).replace(/\n/g, '<br>')}</div>
    </div>`;
  };

  const paintDisabled = () => {
    view.innerHTML = `<div class="chat-shell chat-disabled">
      <div class="chat-empty">
        <span class="chat-avatar lg">AI</span>
        <strong>AI 助手未启用</strong>
        <p>在「设置 → 账号与连接 → AI 设置」里填写 API 地址、模型与 Key 后开启。</p>
        <a class="btn btn-primary" href="#/settings">去启用</a>
      </div>
    </div>`;
  };

  const paint = () => {
    const empty = state.aiMessages.length <= 1;
    view.innerHTML = `<div class="chat-shell">
      <div class="chat-top">
        <span class="chat-mode" id="ai-mode">${esc(state.aiMode || '')}</span>
        <button class="chat-ghost" id="ai-clear" type="button">清空</button>
        <a class="chat-ghost" href="#/settings">设置</a>
      </div>
      <div class="chat-body" id="ai-messages">
        ${empty ? `<div class="chat-empty">
          <span class="chat-avatar lg">AI</span>
          <strong>有什么可以帮你？</strong>
          <p>${esc(welcome)}</p>
          <div class="chat-chips">${suggestions.map(q => `<button type="button" class="chat-chip" data-ai-suggest="${esc(q)}">${esc(q)}</button>`).join('')}</div>
        </div>` : state.aiMessages.map(bubble).join('') + (state.aiBusy ? `<div class="chat-row is-ai"><span class="chat-avatar">AI</span><div class="chat-bubble thinking"><i></i><i></i><i></i></div></div>` : '')}
      </div>
      <div class="chat-compose">
        <textarea class="chat-input" id="ai-input" rows="1" placeholder="问点什么…（Enter 发送，Shift+Enter 换行）" ${state.aiBusy ? 'disabled' : ''}></textarea>
        <button class="chat-send" id="ai-send" type="button" ${state.aiBusy ? 'disabled' : ''} aria-label="发送">${icon('arrow')}</button>
      </div>
      <div class="chat-foot">AI 仅使用配置摘要和查询结果；不会读取 Cookie、API Key 或密码。</div>
    </div>`;

    const input = $('#ai-input');
    const grow = () => { input.style.height = 'auto'; input.style.height = Math.min(140, input.scrollHeight) + 'px'; };
    const scrollDown = () => { const box = $('#ai-messages'); if (box) box.scrollTop = box.scrollHeight; };

    const send = async () => {
      const text = input.value.trim();
      if (!text || state.aiBusy) return;
      state.aiMessages.push({ role: 'user', content: text });
      input.value = '';
      state.aiBusy = true;
      paint();
      const msgs = state.aiMessages.filter(m => m.role === 'user' || m.role === 'assistant').slice(-20);
      try {
        const res = await api('/api/ai/chat', { method: 'POST', body: { messages: msgs } });
        state.aiMessages.push({ role: 'assistant', content: res.reply || '（AI 没有返回内容）' });
        state.aiBusy = false;
        paint();
        if (res.pending) openAIApproval(res.pending);
      } catch (e) {
        state.aiMessages.push({ role: 'assistant', content: '请求失败：' + e.message });
        state.aiBusy = false;
        paint();
      }
    };

    $('#ai-send').addEventListener('click', send);
    input.addEventListener('keydown', e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send(); } });
    input.addEventListener('input', grow);
    $$('[data-ai-suggest]').forEach(btn => btn.addEventListener('click', () => { input.value = btn.dataset.aiSuggest; input.focus(); grow(); }));
    $('#ai-clear').addEventListener('click', () => { state.aiMessages = [{ role: 'assistant', content: welcome }]; state.aiBusy = false; paint(); });
    grow();
    scrollDown();
  };

  paint();
  api('/api/config').then(r => {
    const ai = (r.config || {}).ai || {};
    state.aiEnabled = !!ai.enabled;
    if (!ai.enabled) { paintDisabled(); return; }
    state.aiMode = ai.permission === 'full' ? '完整权限' : '只读';
    const el = $('#ai-mode');
    if (el) el.textContent = state.aiMode;
  }).catch(err => console.debug('[AniFeed]', err));
}
function openAIApproval(pending) {
  const modal = openModal({
    title: '确认 AI 操作', subtitle: '该操作会修改配置，请确认',
    body: `<p>${esc(pending.summary)}</p><div class="notice">${icon('alert')}<span>确认后立即执行，无法自动撤销。</span></div>`,
    footer: `<button class="btn" data-action="modal-close">取消</button><span class="spacer"></span><button class="btn btn-danger" id="ai-confirm">确认执行</button>`,
  });
  modal.el.querySelector('#ai-confirm').addEventListener('click', async () => {
    modal.close();
    try { const r = await api('/api/ai/action', { method: 'POST', body: { token: pending.token } }); toast(r.message || '操作已执行', 'success'); await loadSettingsPage(); }
    catch (e) { toast(e.message, 'error', 5000); }
  });
}


window.AniFeedAI = { loadPage: loadAIPage };
})();
