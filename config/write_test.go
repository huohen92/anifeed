package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func sampleTOMLConfig() *TOMLConfig {
	return &TOMLConfig{
		Auth: TOMLAuthConfig{CookiesFile: ".cookies"},
		Server: TOMLServerConfig{
			Host:        "0.0.0.0",
			Port:        8115,
			AccessToken: "secret-token",
		},
		Database: TOMLDatabaseConfig{Path: "data/app.db"},
		P115: TOMLP115Config{
			DisableCache:  true,
			ChunkDelay:    3,
			ChunkSize:     120,
			CooldownMinMs: 500,
			CooldownMaxMs: 800,
		},
		Proxy: TOMLProxyConfig{HTTP: "http://127.0.0.1:7890"},
		Subscription: TOMLSubscriptionConfig{
			MikanHost:         "https://mikanani.me",
			RSSTimeoutSeconds: 25,
		},
		Offline: TOMLOfflineConfig{
			DefaultCid:         "999888",
			DefaultSavePath:    "番剧/${title}",
			GlobalExclude:      []string{"720p", "/\\d+\\.5/"},
			SkipX5:             true,
			Retry:              1,
			AutoClearCompleted: true,
			DualFolder:         true,
			FinalCid:           "777888",
			FinalMode:          1,
			FinalPathTemplate:  "${month}/${titleYearTmdb}/${animeSeason}",
			OnePerEpisode:      boolPtr(true),
			AutoLanguage:       []string{"chs", "bilingual", "cht"},
			AutoSubtitle:       []string{"softsub", "hardsub", "external"},
			AutoContainer:      []string{"mkv", "mp4"},
			AutoResolution:     []string{"r1080", "r2160", "r720"},
		},
		Rename: TOMLRenameConfig{
			Enabled:  true,
			Template: "[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}",
			Season:   2,
			Offset:   1,
		},
		Auto: TOMLAutoConfig{Enabled: true, IntervalMinutes: 15, RunOnStart: true},
		Notify: TOMLNotifyConfig{
			Type:                   "auto",
			Webhook:                "https://example.com/hook",
			Content:                DefaultNotifyContent(),
			Events:                 DefaultNotifyEvents(),
			OnSuccess:              boolPtr(true),
			OnError:                boolPtr(false),
			IdleDays:               7,
			Template:               DefaultNotifyTemplate(),
			DownloadTimeoutMinutes: 180,
			HourlyHours:            []int{},
		},
		RSS: []TOMLRSSConfig{
			{
				Subgroups: []TOMLRSSSubgroup{{Name: "备份组", URL: "https://example.com/sub.xml", Season: 2, Offset: 1, Enabled: boolPtr(true)}},
				Site:      "mikanani.me",
				Name:      "测试 \"订阅\" \\ 名称",
				URL:       "https://mikanani.me/RSS/Bangumi?bangumiId=2739&subgroupid=12",
				Source:    "mikan",
				Enabled:   boolPtr(true),
				Cid:       "123456",
				SavePath:  "番剧/测试",
				Filter:    "/简体|1080p/",
				Exclude:   "720p",
				TmdbID:    "486456",
				TmdbType:  "tv",
			},
			{
				Site:       "nyaa.si",
				Name:       "VCB-Studio",
				URL:        "https://nyaa.si/?page=rss&u=VCB-Studio",
				Source:     "other",
				Enabled:    boolPtr(false),
				Expiration: 3600,
			},
		},
		Sites: map[string]TOMLSiteConfig{
			"share.dmhy.org": {
				HTTPSAgent: true,
				Headers: map[string]string{
					"Referer":    "https://share.dmhy.org/",
					"User-Agent": "anifeed \"test\"",
				},
			},
			"115.com": {
				Headers: map[string]string{"Cookie": "UID=1; CID=2; SEID=3; KID=4"},
			},
		},
		Compat: TOMLCompatConfig{DisableLegacy: true},
	}
}

func TestMarshalTOMLRoundTrip(t *testing.T) {
	want := sampleTOMLConfig()
	data, err := MarshalTOML(want)
	if err != nil {
		t.Fatalf("MarshalTOML() error = %v", err)
	}
	if !strings.Contains(string(data), "[compat]") {
		t.Fatalf("marshaled config missing [compat] section:\n%s", data)
	}

	got, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v\n%s", err, data)
	}
	if got.Server != want.Server {
		t.Errorf("Server = %+v, want %+v", got.Server, want.Server)
	}
	if got.Database != want.Database {
		t.Errorf("Database = %+v, want %+v", got.Database, want.Database)
	}
	if got.P115 != want.P115 {
		t.Errorf("P115 = %+v, want %+v", got.P115, want.P115)
	}
	if got.Proxy != want.Proxy {
		t.Errorf("Proxy = %+v, want %+v", got.Proxy, want.Proxy)
	}
	if got.Auth != want.Auth {
		t.Errorf("Auth = %+v, want %+v", got.Auth, want.Auth)
	}
	if got.Subscription != want.Subscription {
		t.Errorf("Subscription = %+v, want %+v", got.Subscription, want.Subscription)
	}
	if !reflect.DeepEqual(got.Offline, want.Offline) {
		t.Errorf("Offline = %+v, want %+v", got.Offline, want.Offline)
	}
	if !reflect.DeepEqual(got.Rename, want.Rename) {
		t.Errorf("Rename = %+v, want %+v", got.Rename, want.Rename)
	}
	if got.Auto != want.Auto {
		t.Errorf("Auto = %+v, want %+v", got.Auto, want.Auto)
	}
	if !reflect.DeepEqual(got.Notify, want.Notify) {
		t.Errorf("Notify = %+v, want %+v", got.Notify, want.Notify)
	}
	if !got.Compat.DisableLegacy {
		t.Error("Compat.DisableLegacy = false, want true")
	}
	if len(got.RSS) != len(want.RSS) {
		t.Fatalf("len(RSS) = %d, want %d", len(got.RSS), len(want.RSS))
	}
	for i := range want.RSS {
		if !reflect.DeepEqual(got.RSS[i], want.RSS[i]) {
			t.Errorf("RSS[%d] = %+v, want %+v", i, got.RSS[i], want.RSS[i])
		}
	}
	for host, wantSite := range want.Sites {
		gotSite, ok := got.Sites[host]
		if !ok {
			t.Fatalf("missing site %q", host)
		}
		if gotSite.HTTPSAgent != wantSite.HTTPSAgent {
			t.Errorf("site %q HTTPSAgent = %v, want %v", host, gotSite.HTTPSAgent, wantSite.HTTPSAgent)
		}
		for k, v := range wantSite.Headers {
			if gotSite.Headers[k] != v {
				t.Errorf("site %q header %q = %q, want %q", host, k, gotSite.Headers[k], v)
			}
		}
	}
}

func TestMarshalTOMLDefaults(t *testing.T) {
	data, err := MarshalTOML(&TOMLConfig{})
	if err != nil {
		t.Fatalf("MarshalTOML() error = %v", err)
	}
	cfg, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}
	if cfg.Server.Host != "0.0.0.0" || cfg.Server.Port != 8115 {
		t.Errorf("server defaults = %+v", cfg.Server)
	}
	if cfg.Database.Path != "db.sqlite" {
		t.Errorf("database default = %q", cfg.Database.Path)
	}
	if cfg.Proxy.HTTP == "" {
		t.Error("proxy default is empty")
	}
	if cfg.P115.ChunkSize != 200 {
		t.Errorf("chunk_size = %d, want 200", cfg.P115.ChunkSize)
	}
	if cfg.Subscription.MikanHost != "https://mikanani.me" || cfg.Subscription.RSSTimeoutSeconds != 20 {
		t.Errorf("subscription defaults = %+v", cfg.Subscription)
	}
	if cfg.Auto.IntervalMinutes != 30 {
		t.Errorf("auto interval = %d, want 30", cfg.Auto.IntervalMinutes)
	}
}

func boolPtr(value bool) *bool { return &value }

func TestExpandSavePath(t *testing.T) {
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.Local)
	got := ExpandSavePath("番剧/${title}/${season}/${year}-${month}", "药屋/少女", "mikanani.me", now)
	want := "番剧/药屋_少女/2026秋/2026-10"
	if got != want {
		t.Errorf("ExpandSavePath() = %q, want %q", got, want)
	}
	if got := ExpandSavePath("", "x", "y", now); got != "" {
		t.Errorf("empty template = %q, want empty", got)
	}
}

func TestDetectSource(t *testing.T) {
	cases := map[string]string{
		"https://mikanani.me/RSS/Bangumi?bangumiId=1":  "mikan",
		"https://mikanime.tv/rss":                      "mikan",
		"https://anibt.net/rss/anime.xml?bgmId=1":      "anibt",
		"https://api.animes.garden/feed.xml?subject=1": "animegarden",
		"https://example.com/feed.xml":                 "other",
	}
	for raw, want := range cases {
		if got := DetectSource(raw); got != want {
			t.Errorf("DetectSource(%q) = %q, want %q", raw, got, want)
		}
	}
	if got := NormalizeSource("anime-garden", ""); got != SourceAnimeGarden {
		t.Errorf("NormalizeSource(anime-garden) = %q", got)
	}
}

func TestSaveTOMLBytesBackupAndValidation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	old := []byte("[server]\nport = 8115\n")
	if err := os.WriteFile(path, old, 0o644); err != nil {
		t.Fatal(err)
	}

	next := []byte("[server]\nport = 9000\n")
	if err := SaveTOMLBytes(path, next); err != nil {
		t.Fatalf("SaveTOMLBytes() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(next) {
		t.Errorf("file content = %q, want %q", got, next)
	}
	backup, err := os.ReadFile(path + ".bak")
	if err != nil {
		t.Fatalf("backup missing: %v", err)
	}
	if string(backup) != string(old) {
		t.Errorf("backup = %q, want %q", backup, old)
	}

	if err := SaveTOMLBytes(path, []byte("[server]\nport = 99999\n")); err == nil {
		t.Error("SaveTOMLBytes() accepted invalid port")
	}
	if err := SaveTOMLBytes(path, []byte("[server\n")); err == nil {
		t.Error("SaveTOMLBytes() accepted invalid TOML")
	}
}

func TestValidateTOMLFilterRegex(t *testing.T) {
	data := []byte(`
[[rss]]
site = "mikanani.me"
name = "bad"
url = "https://mikanani.me/rss"
filter = "/[abc/"
`)
	if _, err := ParseTOML(data); err == nil {
		t.Fatal("ParseTOML() accepted invalid regular expression")
	}
}

func TestDisableLegacyPreventsFallback(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", dir)

	if err := os.WriteFile("rss.json", []byte(`{"mikanani.me":[{"name":"legacy","url":"https://mikanani.me/rss"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// Without the compat flag, an empty config.toml falls back to rss.json.
	if err := os.WriteFile("config.toml", []byte("[server]\nport = 8115\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, _, err := LoadWithOptions(CLIParams{}, LoadAllOptions())
	if err != nil {
		t.Fatalf("LoadWithOptions() error = %v", err)
	}
	if len(cfg.RSS["mikanani.me"]) != 1 {
		t.Fatalf("legacy fallback RSS = %v, want 1 feed", cfg.RSS)
	}

	// With disable_legacy, rss.json must be ignored even when [[rss]] is empty.
	content := "[compat]\ndisable_legacy = true\n\n[server]\nport = 8115\n"
	if err := os.WriteFile("config.toml", []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, _, err = LoadWithOptions(CLIParams{}, LoadAllOptions())
	if err != nil {
		t.Fatalf("LoadWithOptions() error = %v", err)
	}
	if len(cfg.RSS) != 0 {
		t.Fatalf("RSS = %v, want empty after disable_legacy", cfg.RSS)
	}
}

func TestRenameAndDualFolderValidation(t *testing.T) {
	base := `[server]
host = "0.0.0.0"
port = 8115
[compat]
disable_legacy = true
`
	dup := base + "[offline]\ndual_folder = true\n"
	cfg, err := ParseTOML([]byte(dup))
	if err == nil {
		t.Fatal("expected dual_folder without final_cid to fail")
	}
	_ = cfg

	final := base + "[offline]\ndual_folder = true\nfinal_cid = \"123\"\n[rename]\nenabled = true\ntemplate = \"${title} S${seasonFormat}E${episodeFormat}\"\n"
	if _, err := ParseTOML([]byte(final)); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
}

func TestRenameDefaults(t *testing.T) {
	if DefaultRenameConfig().Template == "" || DefaultRenameConfig().Season != 1 {
		t.Fatalf("unexpected rename defaults: %+v", DefaultRenameConfig())
	}
}

func TestExpandFinalPathAnimeSeason(t *testing.T) {
	now := time.Date(2026, 10, 23, 12, 0, 0, 0, time.Local)
	got := ExpandFinalPath("番剧/${title}/${animeSeason}", "青之箱", "anibt.net", 2, now)
	if got != "番剧/青之箱/Season 2" {
		t.Fatalf("ExpandFinalPath = %q", got)
	}
	if got := ExpandFinalPath("${seasonNumber}/${title}", "青之箱", "anibt.net", 1, now); got != "1/青之箱" {
		t.Fatalf("seasonNumber = %q", got)
	}
}

func TestExpandFinalPathTitleYearTmdb(t *testing.T) {
	now := time.Date(2026, 10, 23, 12, 0, 0, 0, time.Local)
	got := ExpandFinalPathWithMeta("${titleYearTmdb}", "冰之城墙", "anibt.net", 2, 2026, "283428", now)
	if got != "冰之城墙 (2026) {tmdb-283428}" {
		t.Fatalf("titleYearTmdb = %q", got)
	}
}
