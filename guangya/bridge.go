package guangya

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deadblue/elevengo"
	"github.com/huohen92/anifeed/p115"
	"github.com/skip2/go-qrcode"
)

//go:embed web/*
var webFS embed.FS

const gyClientID = "aMe-8VSlkrbQXpUR"

type Config struct {
	Cookies   string `json:"cookies"`
	DirID     string `json:"dir_id"`
	UpdatedAt string `json:"updated_at"`
}

type FileItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Index int    `json:"index"`
}

type Task struct {
	ID           string     `json:"id"`
	URL          string     `json:"url"`
	Name         string     `json:"name"`
	Status       string     `json:"status"`
	StatusCode   int        `json:"status_code"`
	Progress     float64    `json:"progress"`
	Files        []FileItem `json:"files"`
	CopyStage    string     `json:"copy_stage"`
	CopyProgress float64    `json:"copy_progress"`
	Copied       int        `json:"copied"`
	CopyError    string     `json:"copy_error"`
	Rapid        bool       `json:"rapid"`
	ParentID     string     `json:"parent_id"`
	Raw          any        `json:"raw,omitempty"`
}

type QRState struct {
	Session   string    `json:"session"`
	Image     string    `json:"image"`
	UserCode  string    `json:"user_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Bridge struct {
	mu        sync.Mutex
	dataDir   string
	cfg       Config
	token     string
	refresh   string
	tasks     map[string]*Task
	qr        map[string]*gyDevice
	doneFiles map[string]bool
	pollOnce  sync.Once
	startedAt time.Time
}

type gyDevice struct {
	DeviceCode string
	UserCode   string
	Created    time.Time
}

type gyClient struct {
	token     string
	refresh   string
	onRefresh func(access, refresh string)
	client    *http.Client
}

func New(dataDir, cookies, dirID string) *Bridge {
	b := &Bridge{dataDir: dataDir, tasks: map[string]*Task{}, qr: map[string]*gyDevice{}, doneFiles: map[string]bool{}, startedAt: time.Now()}
	b.loadConfig()
	if cookies != "" {
		b.cfg.Cookies = cookies
	}
	if dirID != "" {
		b.cfg.DirID = dirID
	}
	b.loadTokens()
	b.loadTasks()
	_ = os.MkdirAll(filepath.Join(dataDir, "tmp"), 0o755)
	return b
}

func (b *Bridge) loadConfig() {
	b.cfg.DirID = "0"
	raw, err := os.ReadFile(filepath.Join(b.dataDir, "config.json"))
	if err == nil {
		_ = json.Unmarshal(raw, &b.cfg)
	}
	if strings.TrimSpace(b.cfg.Cookies) == "" {
		for _, path := range []string{filepath.Join("..", ".cookies"), ".cookies"} {
			if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) != "" {
				b.cfg.Cookies = strings.TrimSpace(string(data))
				break
			}
		}
	}
	if b.cfg.DirID == "" {
		b.cfg.DirID = "0"
	}
}

func (b *Bridge) saveConfig() error {
	b.cfg.UpdatedAt = time.Now().Format(time.RFC3339)
	data, _ := json.MarshalIndent(b.cfg, "", "  ")
	if err := os.MkdirAll(b.dataDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(b.dataDir, "config.json"), data, 0o600)
}

func (b *Bridge) loadTasks() {
	raw, err := os.ReadFile(filepath.Join(b.dataDir, "tasks.json"))
	if err != nil {
		return
	}
	var list []*Task
	if json.Unmarshal(raw, &list) != nil {
		return
	}
	for _, t := range list {
		if t != nil && t.ID != "" {
			switch t.CopyStage {
			case "下载中", "上传中", "115 秒传检查", "准备复制":
				t.CopyStage = "已中断，点击重试"
			}
			b.tasks[t.ID] = t
		}
	}
}

func (b *Bridge) saveTasks() {
	b.mu.Lock()
	list := make([]*Task, 0, len(b.tasks))
	for _, t := range b.tasks {
		list = append(list, t)
	}
	b.mu.Unlock()
	raw, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(b.dataDir, "tasks.json"), raw, 0o600)
}

func (b *Bridge) loadTokens() {
	raw, err := os.ReadFile(filepath.Join(b.dataDir, "guangya.json"))
	if err != nil {
		return
	}
	var v struct {
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}
	if json.Unmarshal(raw, &v) == nil {
		b.token, b.refresh = v.Access, v.Refresh
	}
}

func (b *Bridge) saveTokens() error {
	data, _ := json.MarshalIndent(map[string]string{"access_token": b.token, "refresh_token": b.refresh}, "", "  ")
	if err := os.MkdirAll(b.dataDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(b.dataDir, "guangya.json"), data, 0o600)
}

func (b *Bridge) gy() *gyClient {
	b.mu.Lock()
	defer b.mu.Unlock()
	return &gyClient{token: b.token, refresh: b.refresh, onRefresh: func(access, refresh string) {
		b.mu.Lock()
		b.token, b.refresh = access, refresh
		b.mu.Unlock()
		_ = b.saveTokens()
	}, client: &http.Client{Timeout: 60 * time.Second}}
}

func (c *gyClient) post(api string, body any) (map[string]any, error) {
	out, err := c.postOnce(api, body)
	if err != nil && (strings.Contains(err.Error(), "无效token") || strings.Contains(err.Error(), "401")) && c.refresh != "" {
		if rerr := c.refreshAccess(); rerr == nil {
			return c.postOnce(api, body)
		}
	}
	return out, err
}

func (c *gyClient) refreshAccess() error {
	body, _ := json.Marshal(map[string]any{"client_id": gyClientID, "grant_type": "refresh_token", "refresh_token": c.refresh})
	req, err := http.NewRequest(http.MethodPost, "https://account.guangyapan.com/v1/auth/token", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 AniFeed-Guanya")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	access, _ := data["access_token"].(string)
	refresh, _ := data["refresh_token"].(string)
	if access == "" {
		return fmt.Errorf("刷新光鸭登录失败: %s", string(raw))
	}
	c.token = access
	if refresh != "" {
		c.refresh = refresh
	}
	if c.onRefresh != nil {
		c.onRefresh(access, refresh)
	}
	return nil
}

func (c *gyClient) postOnce(api string, body any) (map[string]any, error) {
	var payload io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		payload = strings.NewReader(string(data))
	}
	req, err := http.NewRequest(http.MethodPost, api, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 AniFeed-Guanya")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("光鸭返回异常 (%d): %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	if resp.StatusCode == 401 || intValue(out, "code") == 117 {
		return out, fmt.Errorf("无效token")
	}
	if msg, _ := out["error"].(string); msg != "" && msg != "authorization_pending" {
		desc, _ := out["error_description"].(string)
		return out, fmt.Errorf("%s %s", msg, desc)
	}
	if strings.Contains(api, "api.guangyapan.com") {
		if code := intValue(out, "code"); code > 0 && code != 200 {
			msg := stringValue(out, "msg", "message", "error_message")
			if msg == "" {
				msg = fmt.Sprintf("光鸭接口错误码 %d", code)
			}
			return out, fmt.Errorf("%s", msg)
		}
		if msg := stringValue(out, "msg"); msg != "" && msg != "success" {
			return out, fmt.Errorf("%s", msg)
		}
	}
	return out, nil
}

func (b *Bridge) startQR() (QRState, error) {
	c := b.gy()
	c.token = ""
	data, err := c.post("https://account.guangyapan.com/v1/auth/device/code", map[string]any{
		"client_id": gyClientID, "scope": "user",
	})
	if err != nil {
		return QRState{}, err
	}
	code, _ := data["device_code"].(string)
	userCode, _ := data["user_code"].(string)
	uri, _ := data["verification_uri_complete"].(string)
	if uri == "" {
		uri, _ = data["verification_uri"].(string)
	}
	if code == "" || uri == "" {
		return QRState{}, fmt.Errorf("光鸭没有返回登录二维码")
	}
	png, err := qrcode.Encode(uri, qrcode.Medium, 320)
	if err != nil {
		return QRState{}, err
	}
	session := randomID()
	b.mu.Lock()
	b.qr[session] = &gyDevice{DeviceCode: code, UserCode: userCode, Created: time.Now()}
	b.mu.Unlock()
	return QRState{Session: session, Image: "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), UserCode: userCode, ExpiresAt: time.Now().Add(120 * time.Second)}, nil
}

func (b *Bridge) pollQR(session string) (map[string]any, error) {
	b.mu.Lock()
	dev := b.qr[session]
	b.mu.Unlock()
	if dev == nil {
		return nil, errors.New("二维码会话已失效")
	}
	c := b.gy()
	c.token = ""
	data, err := c.post("https://account.guangyapan.com/v1/auth/token", map[string]any{
		"client_id": gyClientID, "grant_type": "urn:ietf:params:oauth:grant-type:device_code", "device_code": dev.DeviceCode,
	})
	if err != nil {
		if strings.Contains(err.Error(), "authorization_pending") || strings.Contains(err.Error(), "Precondition") {
			return map[string]any{"status": "pending"}, nil
		}
		return nil, err
	}
	access, _ := data["access_token"].(string)
	refresh, _ := data["refresh_token"].(string)
	if access == "" {
		return map[string]any{"status": "pending"}, nil
	}
	b.mu.Lock()
	b.token, b.refresh = access, refresh
	delete(b.qr, session)
	b.mu.Unlock()
	_ = b.saveTokens()
	b.StartPoll()
	return map[string]any{"status": "ok"}, nil
}

func randomID() string {
	var buf [12]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}

func (b *Bridge) StartPoll() {
	b.pollOnce.Do(func() {
		go func() {
			for {
				time.Sleep(3 * time.Second)
				b.refreshTasks()
			}
		}()
	})
}

func (b *Bridge) AddTask(rawURL, parentID string, fileIndexes []int) (*Task, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, errors.New("请填写 magnet 或 ed2k 链接")
	}
	c := b.gy()
	payload := map[string]any{"url": rawURL}
	if strings.TrimSpace(parentID) != "" {
		payload["parentId"] = strings.TrimSpace(parentID)
	}
	if len(fileIndexes) > 0 {
		payload["fileIndexes"] = fileIndexes
	}
	data, err := c.post("https://api.guangyapan.com/cloudcollection/v1/create_task", payload)
	if err != nil {
		return nil, err
	}
	id := deepFindID(data, "taskId", "taskld", "taskID", "taskid", "task_id")
	if id == "" {
		return nil, fmt.Errorf("光鸭没有返回任务 ID: %v", data)
	}
	task := &Task{ID: id, URL: rawURL, Name: basenameURL(rawURL), Status: "等待中", ParentID: parentID}
	b.mu.Lock()
	b.tasks[id] = task
	b.mu.Unlock()
	b.StartPoll()
	b.saveTasks()
	go b.refreshTasks()
	return task, nil
}

func (b *Bridge) previewFiles(rawURL string) ([]FileItem, error) {
	c := b.gy()
	data, err := c.post("https://api.guangyapan.com/cloudcollection/v1/resolve_res", map[string]any{"url": strings.TrimSpace(rawURL)})
	if err != nil {
		return nil, err
	}
	files := collectFiles(data)
	if len(files) == 0 {
		for _, m := range collectMaps(data, func(m map[string]any) bool {
			return stringValue(m, "fileName", "name") != "" && floatValue(m, "fileSize", "size") > 0
		}) {
			files = append(files, FileItem{Name: stringValue(m, "fileName", "name"), Size: int64(floatValue(m, "fileSize", "size"))})
			break
		}
	}
	return files, nil
}

func (b *Bridge) listFolders(parentID string) ([]map[string]any, error) {
	c := b.gy()
	payload := map[string]any{"page": 0, "pageSize": 1000}
	if strings.TrimSpace(parentID) != "" && parentID != "0" {
		payload["parentId"] = parentID
	}
	data, err := c.post("https://api.guangyapan.com/userres/v1/file/get_file_list", payload)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []map[string]any
	for _, m := range collectMaps(data, func(m map[string]any) bool {
		id := findID(m, "fileId", "file_id")
		if id == "" {
			return false
		}
		if v, ok := m["isDir"].(bool); ok {
			return v
		}
		if v, ok := m["fileType"]; ok {
			if n, ok := v.(float64); ok {
				return n == 4
			}
		}
		if v, ok := m["dirType"]; ok {
			if n, ok := v.(float64); ok {
				return n == 1
			}
		}
		if v, ok := m["resType"]; ok {
			if n, ok := v.(float64); ok {
				return n == 2
			}
		}
		return false
	}) {
		id := findID(m, "fileId", "file_id")
		if seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, map[string]any{"id": id, "name": stringValue(m, "fileName", "name", "title")})
	}
	sort.Slice(out, func(i, j int) bool { return out[i]["name"].(string) < out[j]["name"].(string) })
	return out, nil
}

func (b *Bridge) gyAssets() (map[string]any, error) {
	c := b.gy()
	return c.post("https://api.guangyapan.com/assets/v1/get_assets", nil)
}

func (b *Bridge) list115Folders(parentID string) ([]map[string]any, error) {
	b.mu.Lock()
	cookies := b.cfg.Cookies
	b.mu.Unlock()
	if strings.TrimSpace(cookies) == "" {
		return nil, errors.New("未配置 115 Cookies")
	}
	agent, err := p115.NewAgent(cookies)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(parentID) == "" {
		parentID = "0"
	}
	it, err := agent.Agent.FileIterate(parentID)
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for _, f := range it.Items() {
		if !f.IsDirectory {
			continue
		}
		out = append(out, map[string]any{"id": f.FileId, "name": f.Name})
	}
	sort.Slice(out, func(i, j int) bool { return out[i]["name"].(string) < out[j]["name"].(string) })
	return out, nil
}

func (b *Bridge) get115FolderName(id string) (string, error) {
	b.mu.Lock()
	cookies := b.cfg.Cookies
	b.mu.Unlock()
	if strings.TrimSpace(cookies) == "" {
		return "", errors.New("未配置 115 Cookies")
	}
	if strings.TrimSpace(id) == "" || id == "0" {
		return "根目录", nil
	}
	agent, err := p115.NewAgent(cookies)
	if err != nil {
		return "", err
	}
	var file elevengo.File
	if err := agent.Agent.FileGet(id, &file); err != nil {
		return "", err
	}
	if !file.IsDirectory {
		return "", errors.New("目标不是文件夹")
	}
	return file.Name, nil
}

func (b *Bridge) importTask(id string) (*Task, error) {
	c := b.gy()
	data, err := c.post("https://api.guangyapan.com/cloudcollection/v1/list_task", map[string]any{"taskIds": []string{id}, "page": 0, "pageSize": 100})
	if err != nil {
		return nil, err
	}
	maps := collectMaps(data, func(m map[string]any) bool { return findID(m, "taskId", "taskld", "taskID", "taskid", "task_id") == id })
	if len(maps) == 0 {
		return nil, errors.New("没有找到该光鸭任务")
	}
	raw := maps[0]
	b.mu.Lock()
	t := b.tasks[id]
	if t == nil {
		t = &Task{ID: id, URL: stringValue(raw, "url"), Name: stringValue(raw, "taskName", "fileName", "name", "title")}
		if t.Name == "" {
			t.Name = basenameURL(t.URL)
		}
		b.tasks[id] = t
	}
	t.Raw = raw
	t.StatusCode = intValue(raw, "status", "state")
	t.Status = statusText(t.StatusCode)
	t.Progress = floatValue(raw, "progress", "percent", "rate")
	t.Files = []FileItem{}
	if fileID := findID(raw, "fileId", "file_id"); fileID != "" {
		t.Files = []FileItem{{ID: fileID, Name: stringValue(raw, "fileName", "name", "title"), Size: int64(floatValue(raw, "totalSize", "fileSize", "size"))}}
	}
	b.mu.Unlock()
	b.saveTasks()
	b.maybeAutoCopy(id)
	return t, nil
}

func basenameURL(raw string) string {
	if i := strings.Index(raw, "&dn="); i >= 0 {
		s := raw[i+4:]
		if j := strings.IndexByte(s, '&'); j >= 0 {
			s = s[:j]
		}
		if v, err := url.QueryUnescape(s); err == nil {
			s = v
		}
		return strings.TrimSpace(s)
	}
	return "未命名任务"
}

func (b *Bridge) refreshTasks() {
	b.mu.Lock()
	ids := make([]string, 0, len(b.tasks))
	for id := range b.tasks {
		ids = append(ids, id)
	}
	token := b.token
	b.mu.Unlock()
	if token == "" {
		return
	}
	c := &gyClient{token: token, client: &http.Client{Timeout: 45 * time.Second}}
	data, err := c.post("https://api.guangyapan.com/cloudcollection/v1/list_task", map[string]any{"taskIds": ids, "page": 0, "pageSize": 100})
	if err != nil {
		log.Printf("光鸭任务查询失败: %v", err)
		return
	}
	for _, raw := range collectMaps(data, func(m map[string]any) bool {
		return findID(m, "taskId", "taskld", "taskID", "taskid", "task_id") != "" && (m["status"] != nil || m["state"] != nil)
	}) {
		id := findID(raw, "taskId", "taskld", "taskID", "taskid", "task_id")
		b.mu.Lock()
		t := b.tasks[id]
		if t == nil {
			created := time.Unix(int64(floatValue(raw, "createTime", "createdAt")), 0)
			t = &Task{ID: id, URL: stringValue(raw, "url"), Name: stringValue(raw, "taskName", "fileName", "name", "title")}
			if t.Name == "" {
				t.Name = basenameURL(t.URL)
			}
			if created.IsZero() || created.Before(b.startedAt.Add(-30*time.Minute)) {
				t.CopyStage = "历史任务，未自动复制"
			}
			b.tasks[id] = t
		}
		if t != nil {
			t.Raw = raw
			t.StatusCode = intValue(raw, "status", "state")
			t.Status = statusText(t.StatusCode)
			t.Progress = floatValue(raw, "progress", "percent", "rate")
			if t.Progress <= 0 && t.StatusCode == 2 {
				t.Progress = 100
			}
			name := stringValue(raw, "taskName", "fileName", "name", "title")
			if name != "" {
				t.Name = name
			}
			if fileID := findID(raw, "fileId", "file_id"); fileID != "" && len(t.Files) == 0 {
				t.Files = []FileItem{{ID: fileID, Name: name, Size: int64(floatValue(raw, "totalSize", "fileSize", "size"))}}
			}
		}
		b.mu.Unlock()
		b.maybeAutoCopy(id)
	}
	b.saveTasks()
}

func (b *Bridge) maybeAutoCopy(id string) {
	b.mu.Lock()
	t := b.tasks[id]
	hasCookies := strings.TrimSpace(b.cfg.Cookies) != ""
	if t == nil || !hasCookies || t.StatusCode != 2 || t.CopyStage != "" {
		b.mu.Unlock()
		return
	}
	t.CopyStage = "准备复制"
	b.mu.Unlock()
	go b.CopyTask(id)
}

func (b *Bridge) taskFiles(id string) ([]FileItem, error) {
	b.mu.Lock()
	t := b.tasks[id]
	token := b.token
	b.mu.Unlock()
	if t == nil {
		return nil, errors.New("任务不存在")
	}
	c := &gyClient{token: token, client: &http.Client{Timeout: 60 * time.Second}}
	data, err := c.post("https://api.guangyapan.com/userres/v1/file/get_info_by_task_id", map[string]any{"taskId": id})
	if err != nil {
		return nil, err
	}
	files := collectFiles(data)
	if len(files) == 0 && t.Name != "" {
		search, err := c.post("https://api.guangyapan.com/userres/v1/file/search_files", map[string]any{"name": t.Name, "page": 0, "pageSize": 100})
		if err == nil {
			files = collectFiles(search)
		}
	}
	b.mu.Lock()
	if t := b.tasks[id]; t != nil {
		t.Files = files
	}
	b.mu.Unlock()
	return files, nil
}

func collectFiles(v any) []FileItem {
	seen := map[string]bool{}
	var out []FileItem
	for _, m := range collectMaps(v, func(m map[string]any) bool {
		return findID(m, "fileId", "file_id") != "" && fileTypeIsFile(m)
	}) {
		f := FileItem{ID: findID(m, "fileId", "file_id"), Name: stringValue(m, "fileName", "name", "title"), Size: int64(floatValue(m, "fileSize", "size")), Index: intValue(m, "fileIndex", "file_index", "index")}
		if f.ID != "" && !seen[f.ID] {
			seen[f.ID] = true
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func fileTypeIsFile(m map[string]any) bool {
	if v, ok := m["isDir"]; ok {
		if b, ok := v.(bool); ok {
			return !b
		}
	}
	if v, ok := m["fileType"]; ok {
		if n, ok := v.(float64); ok {
			return n != 4
		}
	}
	if v, ok := m["dirType"]; ok {
		if n, ok := v.(float64); ok {
			return n != 1
		}
	}
	if v, ok := m["resType"]; ok {
		if n, ok := v.(float64); ok {
			return n != 2
		}
	}
	return true
}

func collectMaps(v any, pred func(map[string]any) bool) []map[string]any {
	var out []map[string]any
	var walk func(any)
	walk = func(x any) {
		switch y := x.(type) {
		case map[string]any:
			if pred(y) {
				out = append(out, y)
			}
			for _, v := range y {
				walk(v)
			}
		case []any:
			for _, v := range y {
				walk(v)
			}
		}
	}
	walk(v)
	return out
}

func findID(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch s := v.(type) {
			case string:
				if strings.TrimSpace(s) != "" {
					return strings.TrimSpace(s)
				}
			case float64:
				return strconv.FormatInt(int64(s), 10)
			}
		}
	}
	return ""
}
func deepFindID(v any, keys ...string) string {
	var found string
	var walk func(any)
	walk = func(x any) {
		if found != "" {
			return
		}
		switch y := x.(type) {
		case map[string]any:
			if id := findID(y, keys...); id != "" {
				found = id
				return
			}
			for _, child := range y {
				walk(child)
			}
		case []any:
			for _, child := range y {
				walk(child)
			}
		}
	}
	walk(v)
	return found
}

func stringValue(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
func intValue(m map[string]any, keys ...string) int {
	for _, k := range keys {
		switch v := m[k].(type) {
		case float64:
			return int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	return -1
}
func floatValue(m map[string]any, keys ...string) float64 {
	for _, k := range keys {
		switch v := m[k].(type) {
		case float64:
			return v
		case string:
			if n, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64); err == nil {
				return n
			}
		}
	}
	return 0
}
func statusText(code int) string {
	switch code {
	case 0:
		return "等待中"
	case 1, 4:
		return "离线中"
	case 2:
		return "已完成"
	case 3:
		return "重试中"
	case 5:
		return "失败/部分完成"
	}
	return "未知"
}

func (b *Bridge) CopyTask(id string) {
	b.mu.Lock()
	t := b.tasks[id]
	b.mu.Unlock()
	if t == nil {
		return
	}
	files := t.Files
	if len(files) == 0 {
		files, _ = b.taskFiles(id)
	}
	if len(files) == 0 {
		b.setCopy(id, "失败", 0, "没有找到可复制的光鸭文件")
		return
	}
	b.setCopy(id, "准备复制", 0, "")
	b.mu.Lock()
	cfg := b.cfg
	b.mu.Unlock()
	if strings.TrimSpace(cfg.Cookies) == "" {
		b.setCopy(id, "失败", 0, "未配置 115 Cookies")
		return
	}
	agent, err := p115.NewAgent(cfg.Cookies)
	if err != nil {
		b.setCopy(id, "失败", 0, "115 登录失败: "+err.Error())
		return
	}
	dirID := strings.TrimSpace(cfg.DirID)
	if dirID == "" {
		dirID = "0"
	}
	total := len(files)
	copied := 0
	for _, f := range files {
		key := copyKey(f)
		b.mu.Lock()
		if b.doneFiles[key] {
			b.mu.Unlock()
			copied++
			continue
		}
		b.doneFiles[key] = true
		b.mu.Unlock()
		if err := b.copyOne(agent, id, f, dirID); err != nil {
			b.mu.Lock()
			delete(b.doneFiles, key)
			b.mu.Unlock()
			b.setCopy(id, "失败", 0, err.Error())
			return
		}
		copied++
		b.setCopy(id, "已完成", float64(copied)/float64(total)*100, "")
		b.mu.Lock()
		if t := b.tasks[id]; t != nil {
			t.Copied = copied
		}
		b.mu.Unlock()
	}
}

func (b *Bridge) copyOne(agent *p115.Agent, taskID string, f FileItem, dirID string) error {
	c := b.gy()
	info, err := c.post("https://api.guangyapan.com/userres/v1/get_res_download_url", map[string]any{"fileId": f.ID})
	if err != nil {
		return fmt.Errorf("获取光鸭下载链接失败: %w", err)
	}
	downloadURL := deepString(info, "signedURL")
	if downloadURL == "" {
		return fmt.Errorf("光鸭没有返回下载链接")
	}
	b.setCopy(taskID, "下载中", 5, "")
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("光鸭下载失败: HTTP %d", resp.StatusCode)
	}
	tmp, err := os.CreateTemp(filepath.Join(b.dataDir, "tmp"), "gy-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	total := resp.ContentLength
	var done int64
	buf := make([]byte, 1<<20)
	for {
		n, er := resp.Body.Read(buf)
		if n > 0 {
			if _, ew := tmp.Write(buf[:n]); ew != nil {
				tmp.Close()
				return ew
			}
			done += int64(n)
			if total > 0 {
				b.setCopy(taskID, "下载中", 5+float64(done)/float64(total)*45, "")
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			tmp.Close()
			return er
		}
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	fh, err := os.Open(tmpPath)
	if err != nil {
		return err
	}
	defer fh.Close()
	b.setCopy(taskID, "115 秒传检查", 55, "")
	ticket := &elevengo.UploadTicket{}
	// use elevengo ticket via wrapper exposes Agent field
	if err := agent.Agent.UploadCreateTicket(dirID, f.Name, fh, ticket); err != nil {
		return fmt.Errorf("115 上传初始化失败: %w", err)
	}
	if ticket.Exist {
		b.setCopy(taskID, "秒传成功", 100, "")
		return nil
	}
	b.setCopy(taskID, "上传中", 70, "")
	if _, err := fh.Seek(0, io.SeekStart); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, ticket.Url, fh)
	if err != nil {
		return err
	}
	for k, v := range ticket.Header {
		req.Header.Set(k, v)
	}
	up, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer up.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(up.Body, 1<<20))
	if up.StatusCode < 200 || up.StatusCode >= 300 {
		return fmt.Errorf("115 上传失败: HTTP %d %s", up.StatusCode, string(body))
	}
	b.setCopy(taskID, "上传完成", 100, "")
	return nil
}

func copyKey(f FileItem) string {
	name := strings.ToLower(strings.TrimSpace(f.Name))
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	// Ignore 光鸭's duplicate suffix like " (1)".
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] != ')' {
			continue
		}
		if j := strings.LastIndex(base[:i], " ("); j >= 0 {
			if n, err := strconv.Atoi(strings.TrimSuffix(base[j+2:i], " ")); err == nil {
				_ = n
				base = base[:j]
			}
		}
		break
	}
	base = strings.ToLower(base)
	base = strings.NewReplacer(" ", "", "\t", "", "-", "", "_", "").Replace(base)
	return fmt.Sprintf("%d|%s", f.Size, base)
}

func (b *Bridge) setCopy(id, stage string, progress float64, errText string) {
	b.mu.Lock()
	if t := b.tasks[id]; t != nil {
		t.CopyStage = stage
		t.CopyProgress = progress
		t.CopyError = errText
		if stage == "秒传成功" {
			t.Rapid = true
		}
	}
	b.mu.Unlock()
	b.saveTasks()
}

func deepString(v any, key string) string {
	var out string
	var walk func(any)
	walk = func(x any) {
		switch y := x.(type) {
		case map[string]any:
			if s, ok := y[key].(string); ok && s != "" && out == "" {
				out = s
			}
			for _, v := range y {
				walk(v)
			}
		case []any:
			for _, v := range y {
				walk(v)
			}
		}
	}
	walk(v)
	return out
}

func (b *Bridge) snapshot() map[string]any {
	b.mu.Lock()
	defer b.mu.Unlock()
	tasks := make([]*Task, 0, len(b.tasks))
	for _, t := range b.tasks {
		tasks = append(tasks, t)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID > tasks[j].ID })
	return map[string]any{"logged_in": b.token != "", "tasks": tasks, "config": map[string]any{"dir_id": b.cfg.DirID, "has_cookies": b.cfg.Cookies != ""}}
}

func (b *Bridge) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", b.handle)
	mux.HandleFunc("/", b.web)
	return mux
}

func (b *Bridge) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	fail := func(code int, msg string) { w.WriteHeader(code); _ = enc.Encode(map[string]any{"error": msg}) }
	switch {
	case r.Method == "GET" && r.URL.Path == "/api/status":
		_ = enc.Encode(b.snapshot())
	case r.Method == "POST" && r.URL.Path == "/api/gy/qr":
		state, err := b.startQR()
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(state)
	case r.Method == "GET" && r.URL.Path == "/api/gy/qr/status":
		res, err := b.pollQR(r.URL.Query().Get("session"))
		if err != nil {
			fail(400, err.Error())
			return
		}
		_ = enc.Encode(res)
	case r.Method == "POST" && r.URL.Path == "/api/tasks":
		var p struct {
			URL         string `json:"url"`
			ParentID    string `json:"parent_id"`
			FileIndexes []int  `json:"file_indexes"`
		}
		_ = json.NewDecoder(r.Body).Decode(&p)
		t, err := b.AddTask(p.URL, p.ParentID, p.FileIndexes)
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(t)
	case r.Method == "GET" && r.URL.Path == "/api/gy/assets":
		data, err := b.gyAssets()
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(data)
	case r.Method == "GET" && r.URL.Path == "/api/115/folders":
		folders, err := b.list115Folders(r.URL.Query().Get("parent_id"))
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(map[string]any{"folders": folders, "parent_id": r.URL.Query().Get("parent_id")})
	case r.Method == "GET" && r.URL.Path == "/api/115/folder":
		name, err := b.get115FolderName(r.URL.Query().Get("id"))
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(map[string]any{"id": r.URL.Query().Get("id"), "name": name})
	case r.Method == "POST" && r.URL.Path == "/api/gy/preview":
		var p struct {
			URL string `json:"url"`
		}
		_ = json.NewDecoder(r.Body).Decode(&p)
		files, err := b.previewFiles(p.URL)
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(map[string]any{"files": files})
	case r.Method == "GET" && r.URL.Path == "/api/gy/folders":
		folders, err := b.listFolders(r.URL.Query().Get("parent_id"))
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(map[string]any{"folders": folders, "parent_id": r.URL.Query().Get("parent_id")})
	case r.Method == "POST" && r.URL.Path == "/api/tasks/import":
		var p struct {
			ID string `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&p)
		if strings.TrimSpace(p.ID) == "" {
			fail(400, "缺少任务 ID")
			return
		}
		task, err := b.importTask(strings.TrimSpace(p.ID))
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(task)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/tasks/") && strings.HasSuffix(r.URL.Path, "/files"):
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/tasks/"), "/files")
		files, err := b.taskFiles(id)
		if err != nil {
			fail(502, err.Error())
			return
		}
		_ = enc.Encode(files)
	case r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/api/tasks/") && strings.HasSuffix(r.URL.Path, "/copy"):
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/tasks/"), "/copy")
		b.mu.Lock()
		t := b.tasks[id]
		b.mu.Unlock()
		if t == nil {
			fail(404, "任务不存在")
			return
		}
		go b.CopyTask(id)
		_ = enc.Encode(map[string]any{"ok": true})
	case r.Method == "POST" && r.URL.Path == "/api/test115":
		b.mu.Lock()
		cookies := b.cfg.Cookies
		b.mu.Unlock()
		if strings.TrimSpace(cookies) == "" {
			fail(400, "未配置 115 Cookies")
			return
		}
		agent, err := p115.NewAgent(cookies)
		if err != nil {
			fail(502, "115 初始化失败: "+err.Error())
			return
		}
		var info elevengo.UserInfo
		if err := agent.Agent.UserGet(&info); err != nil {
			fail(502, "115 登录失败: "+err.Error())
			return
		}
		_ = enc.Encode(map[string]any{"ok": true, "name": info.Name, "vip": info.IsVip})
	case r.Method == "GET" && r.URL.Path == "/api/config":
		b.mu.Lock()
		cfg := b.cfg
		b.mu.Unlock()
		_ = enc.Encode(cfg)
	case r.Method == "POST" && r.URL.Path == "/api/config":
		var p struct {
			Cookies string `json:"cookies"`
			DirID   string `json:"dir_id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&p)
		b.mu.Lock()
		if strings.TrimSpace(p.Cookies) != "" {
			b.cfg.Cookies = strings.TrimSpace(p.Cookies)
		}
		b.cfg.DirID = strings.TrimSpace(p.DirID)
		if b.cfg.DirID == "" {
			b.cfg.DirID = "0"
		}
		b.mu.Unlock()
		if err := b.saveConfig(); err != nil {
			fail(500, err.Error())
			return
		}
		_ = enc.Encode(map[string]any{"ok": true})
	default:
		fail(404, "not found")
	}
}

func (b *Bridge) web(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	data, err := webFS.ReadFile("web/" + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	switch filepath.Ext(path) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	_, _ = w.Write(data)
}

func cleanupTemp(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		_ = os.Remove(filepath.Join(dir, entry.Name()))
	}
}
