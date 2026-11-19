const $ = s => document.querySelector(s);
let qrTimer = null;
const state = { parentID: '', preview: [], selected: new Set(), folderPath: ['根目录'], folderParents: [], dirID: '', dirName: '', dir115Path: ['根目录'], dir115Parents: [], dirDirty: false, cfgDirty: false };
function toast(msg, type='ok'){const el=document.createElement('div');el.className='toast '+type;el.textContent=msg;$('#toast-root').appendChild(el);setTimeout(()=>el.remove(),4200);}
async function api(path, opts={}){const res=await fetch(path,{headers:{'Content-Type':'application/json'},...opts});const data=await res.json().catch(()=>({}));if(!res.ok||data.error)throw new Error(data.error||('HTTP '+res.status));return data;}
function setLogin(logged){$('#login-pill').textContent=logged?'光鸭已登录':'光鸭未登录';$('#login-pill').classList.toggle('ok',logged);$('#add-card').style.display=logged?'':'none';
 $('#login-card').style.display=logged?'none':'';}
function fmtTraffic(n){return n?((n/1024**3).toFixed(2)+' GB'):'—';}
function fmtSize(n){if(!n)return '';const u=['B','KB','MB','GB','TB'];let i=0;while(n>=1024&&i<u.length-1){n/=1024;i++}return n.toFixed(i?1:0)+u[i];}
function statusClass(s){if(s==='已完成')return 'ok';if(s==='失败/部分完成')return 'err';if(s==='离线中'||s==='等待中')return 'warn';return '';}
function renderTasks(tasks){
 const box=$('#tasks');
 if(!tasks.length){box.innerHTML='<div class="empty">暂无任务</div>';return;}
 box.innerHTML=tasks.map(t=>{
  const files=(t.files||[]).map(f=>`<div class="file-row"><span title="${escapeHTML(f.name)}">${escapeHTML(f.name)}</span><b>${fmtSize(f.size)}</b></div>`).join('');
  const copyPct=Number(t.copy_progress||0).toFixed(0);
  const stage=t.copy_stage||'未开始';
  const canCopy=t.status==='已完成';
  const busy=['下载中','上传中','115 秒传检查','准备复制'].includes(stage);
  return `<div class="task">
    <div class="task-top"><div class="task-name">${escapeHTML(t.name||t.url)}</div><span class="badge ${statusClass(t.status)}">${t.status}</span></div>
    <div class="progress-line"><span>光鸭</span><div class="bar"><i style="width:${Math.max(0,Math.min(100,Number(t.progress)||0))}%"></i></div><b>${Number(t.progress||0).toFixed(0)}%</b></div>
    <div class="progress-line"><span>115</span><div class="bar"><i style="width:${Math.max(0,Math.min(100,Number(copyPct)))}%"></i></div><b>${stage}${stage==='已完成'?'':' '+copyPct+'%'}</b></div>
    <div class="task-foot">
      ${t.files&&t.files.length?`<details class="task-files"><summary>${t.files.length} 个文件</summary><div class="file-list">${files}</div></details>`:''}
      <span class="spacer"></span>
      ${t.rapid?'<span class="badge ok">秒传</span>':''}
      ${canCopy?`<button class="btn small ${t.copy_stage==='失败'||t.copy_stage==='已中断，点击重试'?'':'primary'}" data-copy="${t.id}" ${busy?'disabled':''}>${t.copy_stage==='失败'||t.copy_stage==='已中断，点击重试'?'重试复制':'复制到 115'}</button>`:''}
    </div>
    ${t.copy_error?`<div class="hint" style="color:var(--err)">${escapeHTML(t.copy_error)}</div>`:''}
  </div>`;
 }).join('');
 box.querySelectorAll('[data-copy]').forEach(btn=>btn.addEventListener('click',async()=>{
  btn.disabled=true;try{await api(`api/tasks/${btn.dataset.copy}/copy`,{method:'POST',body:'{}'});toast('已开始复制');refresh();}catch(e){toast(e.message,'err');btn.disabled=false;}
 }));
}
function escapeHTML(v){return String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));}
function markCfgDirty(){
 state.cfgDirty=true;
 const el=$('#save-state'); if(el){el.textContent='有未保存修改';el.className='save-state dirty';}
}
function markCfgSaved(){
 state.cfgDirty=false;
 const el=$('#save-state'); if(el){el.textContent='已保存';el.className='save-state';}
}
async function refresh(){
 try{
  const s=await api('api/status');
  setLogin(s.logged_in); renderTasks(s.tasks||[]);
  $('#gy-status').textContent=s.logged_in?'已登录':'未登录';$('#gy-dot').className='dot '+(s.logged_in?'ok':'');
  if(s.config){
    if(!state.dirDirty) $('#cfg-dir').value=s.config.dir_id||'0';
    $('#cfg-hint').textContent=s.config.has_cookies?'115 Cookies 已配置':'未配置 115 Cookies';
    $('#cookie-state').textContent=s.config.has_cookies?'已配置，留空表示不修改':'未配置，请填写 Cookies';
    $('#p115-status').textContent=s.config.has_cookies?'已配置 Cookies':'未配置';
    $('#p115-dot').className='dot '+(s.config.has_cookies?'ok':'');
    if(!state.dirDirty && s.config.dir_id && s.config.dir_id!==state.dirID){state.dirID=s.config.dir_id;load115Name(s.config.dir_id);}
    if(!state.dirDirty && !s.config.dir_id){state.dirID='';state.dirName='根目录';$('#target-folder').textContent='根目录';$('#cfg-dir-name').value='根目录';}
  }
 }catch(e){toast(e.message,'err');}
}
async function startQR(){
 $('#qr-box').innerHTML='<div class="hint">正在生成二维码…</div>';
 try{
  const q=await api('api/gy/qr',{method:'POST',body:'{}'});
  $('#qr-box').innerHTML=`<img src="${q.image}" alt="光鸭登录二维码">`;
  $('#qr-hint').textContent=`请用光鸭 App 扫码登录${q.user_code?'，设备码 '+q.user_code:''}`;
  clearInterval(qrTimer);
  qrTimer=setInterval(async()=>{
    try{const r=await api('api/gy/qr/status?session='+encodeURIComponent(q.session));if(r.status==='ok'){clearInterval(qrTimer);toast('光鸭登录成功');refresh();}}catch(e){clearInterval(qrTimer);$('#qr-hint').textContent=e.message;}
  },2000);
 }catch(e){$('#qr-box').innerHTML='<button class="btn primary" id="qr-start">重新生成二维码</button>';$('#qr-start').onclick=startQR;toast(e.message,'err');}
}
$('#qr-start').onclick=startQR;
async function loadPreview(){
 const url=$('#task-url').value.trim();if(!url)return toast('请先填写链接','err');
 const box=$('#preview-files');box.innerHTML='<span class="hint">正在解析文件列表…</span>';
 try{
  const r=await api('api/gy/preview',{method:'POST',body:JSON.stringify({url})});
  state.preview=r.files||[];state.selected=new Set(state.preview.map(f=>f.index));
  if(!state.preview.length){box.innerHTML='<span class="hint">没有解析到文件，可直接添加任务。</span>';return;}
  box.innerHTML=state.preview.map(f=>`<label class="preview-file"><input type="checkbox" data-idx="${f.index}" checked> ${escapeHTML(f.name)} <span>${fmtSize(f.size)}</span></label>`).join('');
  box.querySelectorAll('input[data-idx]').forEach(cb=>cb.addEventListener('change',()=>{const i=Number(cb.dataset.idx);if(cb.checked)state.selected.add(i);else state.selected.delete(i);}));
 }catch(e){box.innerHTML=`<span class="hint" style="color:var(--err)">${escapeHTML(e.message)}</span>`;}
}
function closeFolderModal(){document.querySelector('#folder-modal')?.remove();}
async function loadFolder(parent){
 try{
  const r=await api('api/gy/folders?parent_id='+encodeURIComponent(parent));
  const folders=r.folders||[];
  let modal=document.querySelector('#folder-modal');
  if(!modal){modal=document.createElement('div');modal.id='folder-modal';modal.className='modal-backdrop';document.body.appendChild(modal);}
  modal.innerHTML=`<div class="modal folder-modal"><div class="card-head"><div><h2>选择光鸭文件夹</h2><p>${escapeHTML(state.folderPath.join(' / '))}</p></div><button class="btn small" id="folder-close">关闭</button></div>
   <div class="folder-list">${parent?`<button class="folder-item" data-folder-parent="1">← 上一级</button>`:''}
   ${folders.map(f=>`<button class="folder-item" data-folder-id="${f.id}">📁 ${escapeHTML(f.name)}</button>`).join('')||'<div class="empty">没有子文件夹</div>'}</div>
   <div class="modal-foot"><button class="btn" id="folder-root">选择根目录</button><span class="spacer"></span><button class="btn primary" id="folder-use">选择当前文件夹</button></div></div>`;
  modal.querySelector('#folder-close').onclick=closeFolderModal;
  modal.querySelector('#folder-root').onclick=()=>{state.parentID='';state.folderPath=['根目录'];closeFolderModal();$('#folder-label').textContent='光鸭目录：根目录';};
  modal.querySelector('#folder-use').onclick=()=>{closeFolderModal();$('#folder-label').textContent='光鸭目录：'+state.folderPath.join(' / ');};
  modal.querySelectorAll('[data-folder-id]').forEach(btn=>btn.onclick=()=>{state.folderParents.push(state.parentID);state.parentID=btn.dataset.folderId;state.folderPath.push(btn.textContent.replace('📁 ',''));loadFolder(state.parentID);});
  modal.querySelector('[data-folder-parent]')?.addEventListener('click',()=>{state.folderPath.pop();state.parentID=state.folderParents.pop()||'';loadFolder(state.parentID);});
 }catch(e){toast(e.message,'err');}
}
$('#task-preview').onclick=loadPreview;
$('#folder-pick').onclick=()=>loadFolder(state.parentID);
$('#cfg-cookies').addEventListener('input',markCfgDirty);
$('#task-url').addEventListener('input',()=>{state.preview=[];state.selected.clear();$('#preview-files').innerHTML='';});
async function load115Name(id){
 try{const r=await api('api/115/folder?id='+encodeURIComponent(id||'0'));state.dirName=r.name||'已选择';}catch{state.dirName='已选择文件夹';}
 $('#cfg-dir-name').value=state.dirName;$('#target-folder').textContent=state.dirName;
}
async function load115Folder(parent){
 try{
  const r=await api('api/115/folders?parent_id='+encodeURIComponent(parent));
  let modal=document.querySelector('#folder-modal-115');
  if(!modal){modal=document.createElement('div');modal.id='folder-modal-115';modal.className='modal-backdrop';document.body.appendChild(modal);}
  modal.innerHTML=`<div class="modal folder-modal"><div class="card-head"><div><h2>选择 115 文件夹</h2><p>${escapeHTML(state.dir115Path.join(' / '))}</p></div><button class="btn small" id="f115-close">关闭</button></div>
   <div class="folder-list">${parent&&parent!=='0'?`<button class="folder-item" data-f115-parent="1">← 上一级</button>`:''}
   ${(r.folders||[]).map(f=>`<button class="folder-item" data-f115-id="${f.id}">📁 ${escapeHTML(f.name)}</button>`).join('')||'<div class="empty">没有子文件夹</div>'}</div>
   <div class="modal-foot"><button class="btn" id="f115-root">选择根目录</button><span class="spacer"></span><button class="btn primary" id="f115-use">选择当前文件夹</button></div></div>`;
  modal.querySelector('#f115-close').onclick=()=>modal.remove();
  modal.querySelector('#f115-root').onclick=()=>{state.dirID='0';state.dirName='根目录';state.dirDirty=true;markCfgDirty();state.dir115Path=['根目录'];state.dir115Parents=[];modal.remove();$('#cfg-dir').value='0';$('#cfg-dir-name').value='根目录';$('#target-folder').textContent='根目录';};
  modal.querySelector('#f115-use').onclick=()=>{
   const id=parent||'0';
   state.dirID=id; state.dirName=state.dir115Path[state.dir115Path.length-1]||'根目录'; state.dirDirty=true;markCfgDirty();
   $('#cfg-dir').value=id; $('#cfg-dir-name').value=state.dirName; $('#target-folder').textContent=state.dirName;
   modal.remove();
  };
  modal.querySelectorAll('[data-f115-id]').forEach(btn=>btn.onclick=()=>{state.dir115Parents.push(parent);state.dir115Path.push(btn.textContent.replace('📁 ',''));load115Folder(btn.dataset.f115Id);});
  modal.querySelector('[data-f115-parent]')?.addEventListener('click',()=>{const p=state.dir115Parents.pop()||'0';state.dir115Path.pop();load115Folder(p);});
 }catch(e){toast(e.message,'err');}
}
$('#cfg-dir-pick').onclick=()=>load115Folder(state.dirID||'0');
$('#task-add').onclick=async()=>{
 const url=$('#task-url').value.trim();if(!url)return toast('请填写链接','err');
 $('#task-add').disabled=true;try{await api('api/tasks',{method:'POST',body:JSON.stringify({url,parent_id:state.parentID,file_indexes:Array.from(state.selected)})});$('#task-url').value='';$('#preview-files').innerHTML='';state.preview=[];state.selected.clear();toast('已添加到光鸭');refresh();}catch(e){toast(e.message,'err');}finally{$('#task-add').disabled=false;}
};
$('#task-import').onclick=async()=>{
 const id=$('#task-import-id').value.trim();if(!id)return toast('请填写任务 ID','err');
 try{await api('api/tasks/import',{method:'POST',body:JSON.stringify({id})});toast('任务已导入');$('#task-import-id').value='';refresh();}catch(e){toast(e.message,'err');}
};
$('#refresh').onclick=refresh;
$('#test115').onclick=async()=>{
 try{const r=await api('api/test115',{method:'POST',body:'{}'});toast('115 已登录：'+(r.name||'')+(r.vip?' VIP':''));}catch(e){toast(e.message,'err');}
};
$('#cfg-save').onclick=async()=>{
 const btn=$('#cfg-save'), stateEl=$('#save-state');
 btn.disabled=true;btn.textContent='保存中…';stateEl.textContent='保存中…';stateEl.className='save-state saving';
 try{
  await api('api/config',{method:'POST',body:JSON.stringify({cookies:$('#cfg-cookies').value,dir_id:$('#cfg-dir').value.trim()||'0'})});
  state.dirDirty=false; state.dirID=''; $('#cfg-cookies').value=''; markCfgSaved(); toast('设置已保存');
  await refresh();
 }catch(e){stateEl.textContent='保存失败';stateEl.className='save-state dirty';toast(e.message,'err');}
 finally{btn.disabled=false;btn.textContent='保存设置';}
};
async function refreshAssets(){
 try{const d=await api('api/gy/assets');const t=d.data?.highSpeedTraffic;if(t){$('#traffic').textContent=fmtTraffic(t.remained)+' / '+fmtTraffic(t.total);}}catch{}
}
setInterval(refresh,2500);setInterval(refreshAssets,60000);refresh();refreshAssets();
