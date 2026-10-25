package webui

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/rsssite"
	"github.com/zhifengle/anifeed/store"
)

func newTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	return newTestServerWith(t, nil)
}

func newTestServerWith(t *testing.T, setup func(dir string)) (*Server, string) {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", dir)
	if setup != nil {
		setup(dir)
	}

	cfg, source, err := config.Load(config.CLIParams{})
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}
	return New(cfg, source), dir
}

func doRequest(t *testing.T, s *Server, method, path string, body any, headers map[string]string) (*httptest.ResponseRecorder, []byte) {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	return rec, rec.Body.Bytes()
}

func decodeConfig(t *testing.T, data []byte) ConfigResponse {
	t.Helper()
	var resp ConfigResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("decode config response: %v\n%s", err, data)
	}
	return resp
}

func TestConfigAPIWritesAndReloads(t *testing.T) {
	legacy := `{"mikanani.me":[{"name":"旧订阅","url":"https://mikanani.me/rss"}]}`
	s, dir := newTestServerWith(t, func(dir string) {
		if err := os.WriteFile(filepath.Join(dir, "rss.json"), []byte(legacy), 0o644); err != nil {
			t.Fatal(err)
		}
	})

	rec, body := doRequest(t, s, http.MethodGet, "/api/config", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/config status = %d, body = %s", rec.Code, body)
	}
	resp := decodeConfig(t, body)
	if len(resp.Config.RSS) != 1 || resp.Config.RSS[0].Name != "旧订阅" {
		t.Fatalf("legacy RSS not exposed: %+v", resp.Config.RSS)
	}
	if !resp.Meta.LegacyImported {
		t.Error("Meta.LegacyImported = false, want true")
	}

	// Save a normal config through the UI API.
	resp.Config.Auth.Cookies = "UID=1; CID=2; SEID=3; KID=4"
	resp.Config.RSS = []EditableRSS{{
		Site:     "mikanani.me",
		Name:     "新订阅",
		URL:      "https://mikanani.me/RSS/Bangumi?bangumiId=2739",
		Source:   "mikan",
		Enabled:  true,
		Cid:      "123",
		SavePath: "番剧/测试",
		Filter:   "简体",
		Exclude:  "720p",
	}}
	resp.Config.Subscription = EditableSubscription{
		MikanHost:         "https://mikanani.me",
		RSSTimeoutSeconds: 15,
	}
	resp.Config.Offline = EditableOffline{
		DefaultCid:         "999",
		DefaultSavePath:    "番剧/${title}",
		GlobalExclude:      []string{"720p", "/x5/"},
		SkipX5:             true,
		Retry:              1,
		AutoClearCompleted: true,
	}
	resp.Config.Auto = EditableAuto{Enabled: true, IntervalMinutes: 15, RunOnStart: true}
	resp.Config.Notify = EditableNotify{Webhook: "https://example.com/hook", OnSuccess: true, OnError: false}
	resp.Config.Sites = map[string]EditableSite{
		"share.dmhy.org": {HTTPSAgent: true, Headers: map[string]string{"Referer": "https://share.dmhy.org/"}},
	}
	rec, body = doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": resp.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT /api/config status = %d, body = %s", rec.Code, body)
	}
	updated := decodeConfig(t, body)
	if len(updated.Config.RSS) != 1 || updated.Config.RSS[0].Name != "新订阅" {
		t.Fatalf("updated RSS = %+v", updated.Config.RSS)
	}
	if !updated.Meta.ConfigExists {
		t.Error("Meta.ConfigExists = false after save")
	}
	if updated.Config.Offline.DefaultCid != "999" || !updated.Config.Offline.SkipX5 || updated.Config.Offline.Retry != 1 {
		t.Errorf("offline settings not persisted: %+v", updated.Config.Offline)
	}
	if !updated.Config.Auto.Enabled || updated.Config.Auto.IntervalMinutes != 15 {
		t.Errorf("auto settings not persisted: %+v", updated.Config.Auto)
	}
	if updated.Config.Notify.Webhook != "https://example.com/hook" {
		t.Errorf("notify settings not persisted: %+v", updated.Config.Notify)
	}
	if updated.Config.RSS[0].Source != "mikan" || !updated.Config.RSS[0].Enabled {
		t.Errorf("feed source/enabled not persisted: %+v", updated.Config.RSS[0])
	}

	// Files on disk.
	configData, err := os.ReadFile(filepath.Join(dir, "config.toml"))
	if err != nil {
		t.Fatalf("config.toml not written: %v", err)
	}
	text := string(configData)
	if !strings.Contains(text, "disable_legacy = true") {
		t.Errorf("config.toml missing disable_legacy:\n%s", text)
	}
	if !strings.Contains(text, `name = "新订阅"`) {
		t.Errorf("config.toml missing feed:\n%s", text)
	}
	cookies, err := os.ReadFile(filepath.Join(dir, ".cookies"))
	if err != nil {
		t.Fatalf(".cookies not written: %v", err)
	}
	if strings.TrimSpace(string(cookies)) != "UID=1; CID=2; SEID=3; KID=4" {
		t.Errorf(".cookies = %q", cookies)
	}

	// Emptying the feed list must not fall back to rss.json.
	updated.Config.RSS = []EditableRSS{}
	rec, body = doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": updated.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT empty RSS status = %d, body = %s", rec.Code, body)
	}
	rec, body = doRequest(t, s, http.MethodGet, "/api/config", nil, nil)
	final := decodeConfig(t, body)
	if len(final.Config.RSS) != 0 {
		t.Fatalf("RSS = %+v after deleting all, legacy fallback happened", final.Config.RSS)
	}
}

func TestRawConfigAPI(t *testing.T) {
	s, dir := newTestServer(t)

	rec, body := doRequest(t, s, http.MethodGet, "/api/config/raw", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET raw status = %d, body = %s", rec.Code, body)
	}
	var raw struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(raw.Content, "[server]") {
		t.Errorf("raw config = %s", raw.Content)
	}

	content := `# handwritten comment
[server]
host = "127.0.0.1"
port = 8123

[[rss]]
site = "nyaa.si"
name = "原始配置订阅"
url = "https://nyaa.si/?page=rss"
`
	rec, body = doRequest(t, s, http.MethodPut, "/api/config/raw", map[string]string{"content": content}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT raw status = %d, body = %s", rec.Code, body)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.toml.bak")); err == nil {
		// The default config did not exist, so no backup is expected. This is
		// only informational; the next save must create one.
		t.Log("backup existed unexpectedly")
	}
	rec, body = doRequest(t, s, http.MethodGet, "/api/config/raw", nil, nil)
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(raw.Content, "# handwritten comment") {
		t.Errorf("raw comment lost:\n%s", raw.Content)
	}

	rec, body = doRequest(t, s, http.MethodPut, "/api/config/raw", map[string]string{"content": "[server\n"}, nil)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("invalid raw PUT status = %d, want 400, body = %s", rec.Code, body)
	}
}

func TestAccessTokenMiddleware(t *testing.T) {
	s, _ := newTestServer(t)

	resp, _ := s.loadEditable()
	resp.Config.Server.AccessToken = "s3cret"
	rec, body := doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": resp.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT config status = %d, body = %s", rec.Code, body)
	}

	rec, _ = doRequest(t, s, http.MethodGet, "/api/status", nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GET status without token = %d, want 401", rec.Code)
	}
	rec, _ = doRequest(t, s, http.MethodGet, "/api/status", nil, map[string]string{"X-Access-Token": "s3cret"})
	if rec.Code != http.StatusOK {
		t.Errorf("GET status with token = %d, want 200", rec.Code)
	}
	rec, _ = doRequest(t, s, http.MethodGet, "/api/status?token=s3cret", nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("query-string token = %d, want 401", rec.Code)
	}
	rec, _ = doRequest(t, s, http.MethodGet, "/api/health", nil, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("GET health without token = %d, want 200", rec.Code)
	}
	rec, _ = doRequest(t, s, http.MethodGet, "/", nil, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("GET / without token = %d, want 200", rec.Code)
	}
	rec, _ = doRequest(t, s, http.MethodPost, "/api/feeds/run-all",
		map[string]any{}, map[string]string{"X-Access-Token": "s3cret", "Origin": "https://evil.example"})
	if rec.Code != http.StatusForbidden {
		t.Errorf("cross-origin POST = %d, want 403", rec.Code)
	}
}

func TestHistoryAndStatus(t *testing.T) {
	s, dir := newTestServer(t)
	dbPath := filepath.Join(dir, "db.sqlite")
	st, err := store.NewWithPath(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	items := []rsssite.MagnetItem{
		{Title: "第一条", Link: "https://example.com/1", Magnet: "magnet:?xt=urn:btih:aaa"},
		{Title: "第二条", Link: "https://example.com/2", Magnet: "magnet:?xt=urn:btih:bbb"},
	}
	if err := st.SaveMagnetItems(items, "测试订阅", "https://example.com/rss", []string{"hash1", "hash2"}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	rec, body := doRequest(t, s, http.MethodGet, "/api/history?limit=10", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("history status = %d, body = %s", rec.Code, body)
	}
	var history historyResponse
	if err := json.Unmarshal(body, &history); err != nil {
		t.Fatal(err)
	}
	if len(history.Items) != 2 {
		t.Fatalf("history items = %d, want 2", len(history.Items))
	}
	if history.Items[0].Title != "第二条" {
		t.Errorf("newest first: got %q", history.Items[0].Title)
	}

	rec, body = doRequest(t, s, http.MethodGet, "/api/history?name=测试订阅", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("filtered history status = %d", rec.Code)
	}

	rec, body = doRequest(t, s, http.MethodGet, "/api/status", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, body)
	}
	var status statusResponse
	if err := json.Unmarshal(body, &status); err != nil {
		t.Fatal(err)
	}
	if status.HistoryCount != 2 {
		t.Errorf("HistoryCount = %d, want 2", status.HistoryCount)
	}
	if status.TodayCount != 2 {
		t.Errorf("TodayCount = %d, want 2", status.TodayCount)
	}
	if status.LoginState != "nocookies" {
		t.Errorf("LoginState = %q, want nocookies", status.LoginState)
	}
	if status.StartedAt.IsZero() || time.Since(status.StartedAt) > time.Hour {
		t.Errorf("StartedAt = %v", status.StartedAt)
	}
}

func TestNoCookiesConflictsAndValidation(t *testing.T) {
	s, _ := newTestServer(t)

	cases := []struct {
		method string
		path   string
		body   any
		status int
	}{
		{http.MethodGet, "/api/tasks", nil, http.StatusConflict},
		{http.MethodPost, "/api/feeds/run", feedPayload{Name: "x", URL: "https://mikanani.me/rss"}, http.StatusConflict},
		{http.MethodPost, "/api/feeds/run-all", map[string]any{}, http.StatusConflict},
		{http.MethodPost, "/api/magnets", magnetsPayload{Links: []string{"magnet:?xt=urn:btih:123"}}, http.StatusConflict},
		{http.MethodPost, "/api/feeds/preview", feedPayload{URL: "not-a-url"}, http.StatusBadRequest},
		{http.MethodPost, "/api/magnets", magnetsPayload{Links: []string{"not-a-link"}}, http.StatusBadRequest},
		{http.MethodPost, "/api/tasks/clear", map[string]any{"type": 9}, http.StatusBadRequest},
	}
	for _, tc := range cases {
		rec, body := doRequest(t, s, tc.method, tc.path, tc.body, nil)
		if rec.Code != tc.status {
			t.Errorf("%s %s = %d, want %d, body = %s", tc.method, tc.path, rec.Code, tc.status, body)
		}
	}
}

func TestStaticIndexAndSecurityHeaders(t *testing.T) {
	s, _ := newTestServer(t)

	rec, body := doRequest(t, s, http.MethodGet, "/", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / = %d", rec.Code)
	}
	if !strings.Contains(string(body), "AniFeed") {
		t.Errorf("index.html does not mention anifeed")
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q", got)
	}
	if got := rec.Header().Get("Content-Security-Policy"); got == "" {
		t.Error("missing Content-Security-Policy")
	}

	rec, _ = doRequest(t, s, http.MethodGet, "/app.js", nil, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("GET /app.js = %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Cache-Control"), "no-cache") {
		t.Errorf("Cache-Control = %q", rec.Header().Get("Cache-Control"))
	}
}

func TestRawConfigRollback(t *testing.T) {
	s, _ := newTestServer(t)

	good := []byte("[server]\nport = 8115\n")
	rec, body := doRequest(t, s, http.MethodPut, "/api/config/raw", map[string]string{"content": string(good)}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("initial raw save = %d, body = %s", rec.Code, body)
	}

	bad := "[auth]\ncookies_file = \"missing.cookies\"\n\n[server]\nport = 8115\n"
	rec, body = doRequest(t, s, http.MethodPut, "/api/config/raw", map[string]string{"content": bad}, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("bad raw save = %d, want 400, body = %s", rec.Code, body)
	}

	rec, body = doRequest(t, s, http.MethodGet, "/api/config", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("config did not recover after rollback: %d, body = %s", rec.Code, body)
	}
	resp := decodeConfig(t, body)
	if !resp.Meta.ConfigExists || resp.Config.Server.Port != 8115 {
		t.Fatalf("rolled back config = %+v", resp.Config)
	}
}

func TestLogoutClearsCookies(t *testing.T) {
	s, dir := newTestServer(t)
	resp, err := s.loadEditable()
	if err != nil {
		t.Fatal(err)
	}
	resp.Config.Auth.Cookies = "UID=1; CID=2; SEID=3; KID=4"
	rec, body := doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": resp.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("save cookies = %d, body = %s", rec.Code, body)
	}
	if data, err := os.ReadFile(filepath.Join(dir, ".cookies")); err != nil || len(data) == 0 {
		t.Fatalf("cookies file not written: %v %q", err, data)
	}

	rec, body = doRequest(t, s, http.MethodPost, "/api/auth/logout", map[string]any{}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("logout = %d, body = %s", rec.Code, body)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".cookies"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(data)) != "" {
		t.Errorf("cookies after logout = %q, want empty", data)
	}
	rec, body = doRequest(t, s, http.MethodGet, "/api/status", nil, nil)
	var status statusResponse
	if err := json.Unmarshal(body, &status); err != nil {
		t.Fatal(err)
	}
	if status.HasCookies || status.LoginState != "nocookies" {
		t.Errorf("status after logout = %+v", status)
	}
}

func TestNotifyWebhook(t *testing.T) {
	var gotBody string
	var gotType string
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("webhook method = %s", r.Method)
		}
		gotType = r.Header.Get("Content-Type")
		data, _ := io.ReadAll(r.Body)
		gotBody = string(data)
		w.WriteHeader(http.StatusOK)
	}))
	defer webhook.Close()

	s, _ := newTestServer(t)
	resp, err := s.loadEditable()
	if err != nil {
		t.Fatal(err)
	}
	resp.Config.Notify = EditableNotify{Type: "generic", Webhook: webhook.URL, Content: []string{"summary"}, OnSuccess: true, OnError: true}
	rec, body := doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": resp.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("save notify config = %d, body = %s", rec.Code, body)
	}

	rec, body = doRequest(t, s, http.MethodPost, "/api/notify/test", map[string]any{}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("test notify = %d, body = %s", rec.Code, body)
	}
	if !strings.Contains(gotType, "text/plain") {
		t.Errorf("content type = %q", gotType)
	}
	if !strings.Contains(gotBody, "AniFeed 测试通知") || !strings.Contains(gotBody, "状态：成功") {
		t.Errorf("notification body = %q", gotBody)
	}
}

func TestFeedSourceNormalization(t *testing.T) {
	s, _ := newTestServer(t)
	resp, err := s.loadEditable()
	if err != nil {
		t.Fatal(err)
	}
	resp.Config.RSS = []EditableRSS{{
		Name:    "AnimeGarden 测试",
		URL:     "https://api.animes.garden/feed.xml?subject=123&fansub=VCB",
		Source:  "anime-garden",
		Enabled: false,
	}}
	rec, body := doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": resp.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("save source = %d, body = %s", rec.Code, body)
	}
	updated := decodeConfig(t, body)
	if updated.Config.RSS[0].Source != "animegarden" {
		t.Errorf("source = %q, want animegarden", updated.Config.RSS[0].Source)
	}
	if updated.Config.RSS[0].Enabled {
		t.Error("disabled feed became enabled")
	}
}

func TestAutoInterval(t *testing.T) {
	cfg := &config.Config{Auto: config.AutoConfig{IntervalMinutes: 5}}
	if got := autoInterval(cfg); got != 5*time.Minute {
		t.Errorf("autoInterval() = %v, want 5m", got)
	}
	cfg.Auto.IntervalMinutes = 0
	if got := autoInterval(cfg); got != 30*time.Minute {
		t.Errorf("default autoInterval() = %v, want 30m", got)
	}
}

func TestBrowseAniBTAndGroups(t *testing.T) {
	old := anibtBrowseBase
	defer func() { anibtBrowseBase = old }()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/seasons/anime", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"currentSeason":"2026秋","requestedSeason":"2026秋","availableSeasons":["2026秋"],"byWeekday":[{"weekdayLabel":"周一","animes":[{"bgmId":624923,"cover":"https://example.com/a.jpg","title":{"primary":"Test Anime","chinese":"测试番剧"},"rating":6.4,"format":"TV","episodes":12,"rssReleaseCount":14}]}]}}`)
	})
	mux.HandleFunc("/api/anime/groups", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"title":{"primary":"Test Anime","chinese":"测试番剧"},"groups":[{"slug":"sub","name":"字幕组 A","lastUpdatedAt":1792246630105,"items":[{"title":"第 01 话","magnet":"magnet:?xt=urn:btih:111","size":123,"publishedAt":1792246630105}]}]}}`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	anibtBrowseBase = ts.URL

	s, _ := newTestServer(t)
	list, err := s.browseAniBT("", "2026秋")
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 || len(list.Weeks) != 1 || list.Weeks[0].Items[0].Title != "测试番剧" {
		t.Fatalf("list = %+v", list)
	}
	groups, err := s.browseAniBTGroups("624923")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups.Groups) != 1 || !strings.Contains(groups.Groups[0].RSS, "bgmId=624923") {
		t.Fatalf("groups = %+v", groups)
	}
	if len(groups.Groups[0].Items) != 1 || groups.Groups[0].Items[0].Magnet == "" {
		t.Fatalf("group items = %+v", groups.Groups[0].Items)
	}
}

func TestBrowseAnimeGardenAndGroups(t *testing.T) {
	old := animeGardenBrowseBase
	defer func() { animeGardenBrowseBase = old }()

	mux := http.NewServeMux()
	mux.HandleFunc("/subjects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"status":"OK","subjects":[{"id":621843,"name":"测试番剧","search":{"include":["Test Anime"]},"activedAt":"2026-10-05T00:00:00.000Z","isArchived":false}]}`)
	})
	mux.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"status":"OK","resources":[{"title":"第 01 话","magnet":"magnet:?xt=urn:btih:222","size":456,"createdAt":"2026-10-06T00:00:00.000Z","fansub":{"id":1,"name":"字幕组 B"}}]}`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	animeGardenBrowseBase = ts.URL

	s, _ := newTestServer(t)
	list, err := s.browseAnimeGarden("测试")
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 || list.Weeks[0].Label != "周一" {
		t.Fatalf("list = %+v", list)
	}
	groups, err := s.browseAnimeGardenGroups("621843")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups.Groups) != 1 || !strings.Contains(groups.Groups[0].RSS, "subject=621843") {
		t.Fatalf("groups = %+v", groups)
	}
}

func TestBrowseMikanParsing(t *testing.T) {
	mux := http.NewServeMux()
	listHTML := `<html><body>
<div class="date-select"><span class="date-text">2026秋</span><div class="dropdown-menu"><ul><li><a data-year="2026" data-season="秋">2026 秋</a></li></ul></div></div>
<div class="sk-bangumi"><div>周一</div><ul class="an-ul"><li><a href="/Home/Bangumi/3548">药屋少女的呢喃</a><span data-src="/img/cover.jpg"></span></li></ul></div>
</body></html>`
	mux.HandleFunc("/Home/Search", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, listHTML)
	})
	mux.HandleFunc("/Home/BangumiCoverFlowByDayOfWeek", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, listHTML)
	})
	mux.HandleFunc("/Home/SearchUnused", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `<html><body>
<div class="date-select"><span class="date-text">2026秋</span><div class="dropdown-menu"><ul><li><a data-year="2026" data-season="秋">2026 秋</a></li></ul></div></div>
<div class="sk-bangumi"><div>周一</div><ul class="an-ul"><li><a href="/Home/Bangumi/3548">药屋少女的呢喃</a><span data-src="/img/cover.jpg"></span></li></ul></div>
</body></html>`)
	})
	mux.HandleFunc("/Home/Bangumi/3548", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `<html><body>
<div class="bangumi-title">药屋少女的呢喃</div>
<div class="leftbar-item"><a class="subgroup-name" data-anchor="#213">简日字幕组</a><span class="date">周一</span></div>
<div id="213"><a class="mikan-rss" href="/RSS/Bangumi?bangumiId=3548&subgroupid=213">RSS</a></div>
<table><tbody><tr><td>x</td><td><a>第 01 话</a><a data-clipboard-text="magnet:?xt=urn:btih:333">magnet</a><a href="/a.torrent">torrent</a></td><td>100 MB</td><td>2026-10-01</td></tr></tbody></table>
</body></html>`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	s, _ := newTestServer(t)
	s.cfg.Subscription.MikanHost = ts.URL

	list, err := s.browseMikan("", "2026秋")
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 || list.Weeks[0].Items[0].ID != "3548" {
		t.Fatalf("list = %+v", list)
	}
	if !strings.HasPrefix(list.Weeks[0].Items[0].Cover, ts.URL) {
		t.Errorf("cover = %q", list.Weeks[0].Items[0].Cover)
	}

	groups, err := s.browseMikanGroups("3548")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups.Groups) != 1 || !strings.Contains(groups.Groups[0].RSS, "bangumiId=3548") {
		t.Fatalf("groups = %+v", groups)
	}
	if groups.Groups[0].Items[0].Size != 100*1024*1024 {
		t.Errorf("size = %d", groups.Groups[0].Items[0].Size)
	}
}

func TestFSListRequiresLogin(t *testing.T) {
	s, _ := newTestServer(t)
	rec, body := doRequest(t, s, http.MethodGet, "/api/fs/list?cid=0", nil, nil)
	if rec.Code != http.StatusConflict {
		t.Errorf("fs list without login = %d, body = %s", rec.Code, body)
	}
}

func TestParseSizeText(t *testing.T) {
	cases := map[string]int64{
		"100 MB":  100 * 1024 * 1024,
		"4.5GB":   4831838208,
		"1.5 GiB": 1610612736,
		"512KB":   524288,
		"":        0,
		"bad":     0,
	}
	for text, want := range cases {
		if got := parseSizeText(text); got != want {
			t.Errorf("parseSizeText(%q) = %d, want %d", text, got, want)
		}
	}
}

func TestMikanBrowseBases(t *testing.T) {
	cfg := &config.Config{Subscription: config.SubscriptionConfig{MikanHost: "https://mikanani.me"}}
	bases := mikanBrowseBases(cfg)
	if len(bases) != 2 || bases[1] != "https://mikanime.tv" {
		t.Errorf("bases = %#v, want primary + mirror", bases)
	}
	cfg.Subscription.MikanHost = "https://example.com"
	if got := mikanBrowseBases(cfg); len(got) != 1 || got[0] != "https://example.com" {
		t.Errorf("custom bases = %#v", got)
	}
}

func TestWebUIAuthSetupLoginAndProtectAPI(t *testing.T) {
	s, _ := newTestServerWith(t, nil)
	rec, body := doRequest(t, s, http.MethodPost, "/api/auth/setup", map[string]any{"username": "admin", "password": "secret123"}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("setup status=%d body=%s", rec.Code, body)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("setup did not set session cookie")
	}
	rec, body = doRequest(t, s, http.MethodGet, "/api/config", nil, nil)
	if rec.Code != http.StatusUnauthorized || !strings.Contains(string(body), "auth_required") {
		t.Fatalf("unauthenticated config status=%d body=%s", rec.Code, body)
	}
	rec, _ = doRequest(t, s, http.MethodGet, "/api/config", nil, map[string]string{"Cookie": cookies[0].Name + "=" + cookies[0].Value})
	if rec.Code != http.StatusOK {
		t.Fatalf("authenticated config status=%d", rec.Code)
	}
	rec, body = doRequest(t, s, http.MethodPost, "/api/auth/login", map[string]any{"username": "admin", "password": "wrong"}, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password status=%d body=%s", rec.Code, body)
	}
	rec, _ = doRequest(t, s, http.MethodPost, "/api/auth/login", map[string]any{"username": "admin", "password": "secret123"}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login status=%d", rec.Code)
	}
}

func TestSameHostIgnoresDefaultPort(t *testing.T) {
	if !sameHost("http://127.0.0.1:80", "127.0.0.1") {
		t.Fatal("default port should be ignored")
	}
	if !sameHost("http://localhost:8115", "localhost:8115") {
		t.Fatal("same host/port should match")
	}
	if sameHost("http://evil.example", "127.0.0.1") {
		t.Fatal("different hosts must not match")
	}
}

func TestNotifyAndAutoRuleFlagsRoundTrip(t *testing.T) {
	s, _ := newTestServerWith(t, nil)
	rec, body := doRequest(t, s, http.MethodGet, "/api/config", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET config status=%d body=%s", rec.Code, body)
	}
	resp := decodeConfig(t, body)
	resp.Config.Notify.SkipNoChange = true
	resp.Config.Notify.HourlyEnabled = true
	resp.Config.Notify.HourlyHours = []int{6, 18}
	resp.Config.Notify.Events = []string{"completed", "failed"}
	resp.Config.Offline.AutoRules = true
	resp.Config.Lab.MultiSubgroup = true
	resp.Config.Lab.PreferredSubgroups = []string{"LoliHouse"}
	resp.Config.Lab.WashEnabled = true
	resp.Config.Lab.WashFolders = []string{"123456"}
	rec, body = doRequest(t, s, http.MethodPut, "/api/config", map[string]any{"config": resp.Config}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT config status=%d body=%s", rec.Code, body)
	}
	resp = decodeConfig(t, body)
	if !resp.Config.Notify.SkipNoChange || !resp.Config.Notify.HourlyEnabled || len(resp.Config.Notify.HourlyHours) != 2 || len(resp.Config.Notify.Events) != 2 || !resp.Config.Offline.AutoRules || !resp.Config.Lab.MultiSubgroup || !resp.Config.Lab.WashEnabled || len(resp.Config.Lab.PreferredSubgroups) != 1 {
		t.Fatalf("flags not persisted: notify=%+v offline=%+v", resp.Config.Notify, resp.Config.Offline)
	}
}
