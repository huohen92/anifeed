package rsssite

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/zhifengle/anifeed/config"
)

func TestDmhy(t *testing.T) {
	if os.Getenv("ANIFEED_TEST_LIVE") == "" {
		t.Skip("set ANIFEED_TEST_LIVE=1 to test the live Dmhy RSS")
	}
	dmhy := &Dmhy{}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://share.dmhy.org/topics/rss/rss.xml")
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range feed.Items {
		t.Log(dmhy.GetMagnetItem(item))
	}
}

func TestAcgnx(t *testing.T) {
	if os.Getenv("ANIFEED_TEST_LIVE") == "" {
		t.Skip("set ANIFEED_TEST_LIVE=1 to test the live Acgnx RSS")
	}
	acgnx := &Acgnx{}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://share.acgnx.net/rss.xml")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Items) == 0 {
		t.Fatal("expected at least one Acgnx item")
	}
	for _, item := range feed.Items[:1] {
		t.Log(acgnx.GetMagnetItem(item))
	}
}

func TestAnibt(t *testing.T) {
	anibt := &Anibt{}
	data, err := os.ReadFile("../test/anibt.rss")
	if err != nil {
		t.Fatal(err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Items) == 0 {
		t.Fatal("expected at least one anibt item")
	}

	got := anibt.GetMagnet(feed.Items[0])
	want := "magnet:?xt=urn:btih:8b8c2f0a461a212b0b7417289376ff243284edc6"
	if got != want {
		t.Fatalf("expected first magnet %q, got %q", want, got)
	}

	for _, item := range feed.Items {
		magnet := anibt.GetMagnet(item)
		if !strings.HasPrefix(magnet, "magnet:?xt=urn:btih:") {
			t.Fatalf("expected magnet URI for %q, got %q", item.Title, magnet)
		}
	}
}

func TestGetRssConfigByURL(t *testing.T) {
	rssConfig := GetRssConfigByURL("http://share.dmhy.org/topics/rss/rss.xml")
	t.Log(rssConfig)
}

func TestGetRssConfigByURLMatchesSchemeAndQueryOrder(t *testing.T) {
	SetRssJsonPath("../rss.json")
	t.Cleanup(func() {
		SetRssJsonPath("")
	})

	rssConfig := GetRssConfigByURL("https://share.dmhy.org/topics/rss/rss.xml?sort_id=2&team_id=0&order=date-desc&keyword=%E6%B0%B4%E6%98%9F%E7%9A%84%E9%AD%94%E5%A5%B3")
	if rssConfig == nil {
		t.Fatal("expected config, got nil")
	}
	if rssConfig.SavePath != "文件夹名称" {
		t.Fatalf("expected savepath 文件夹名称, got %q", rssConfig.SavePath)
	}
	if rssConfig.Cid != "" {
		t.Fatalf("expected empty cid from sample config, got %q", rssConfig.Cid)
	}
}

func TestReadRssConfigDictFromUserConfigDir(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")

	configDir := filepath.Join(homeDir, ".config", "anifeed")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	configFile := filepath.Join(configDir, "rss.json")
	configContent := `{"example.com":[{"name":"from-config-dir","url":"https://example.com/rss"}]}`
	if err := os.WriteFile(configFile, []byte(configContent), 0o600); err != nil {
		t.Fatalf("failed to create rss config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
		SetRssJsonPath("")
		RssConfigDict = nil
	})
	SetRssJsonPath("")
	RssConfigDict = nil

	configs := ReadRssConfigDict()
	if configs == nil {
		t.Fatalf("expected rss config to be read")
	}
	got := (*configs)["example.com"]
	if len(got) != 1 || got[0].Name != "from-config-dir" {
		t.Fatalf("unexpected rss config: %#v", got)
	}
	if !got[0].Enabled {
		t.Fatalf("feeds without explicit enabled should default to enabled: %#v", got[0])
	}
}

func TestCompileFilter(t *testing.T) {
	t.Parallel()
	// plain substring rules are not compiled into a regex
	if re, err := compileFilter("简体内嵌"); err != nil || re != nil {
		t.Errorf("compileFilter(plain) = (%v, %v), want (nil, nil)", re, err)
	}
	// /.../-wrapped rules compile into a regex
	if re, err := compileFilter("/简体|1080p/"); err != nil || re == nil {
		t.Errorf("compileFilter(/regex/) = (%v, %v), want non-nil", re, err)
	}
	// invalid regex returns an error instead of panicking
	if _, err := compileFilter("/[abc/"); err == nil {
		t.Error("compileFilter(invalid) = nil error, want error")
	}
}

func TestMatchRule(t *testing.T) {
	t.Parallel()
	// empty rule matches everything (no whitelist/blacklist restriction)
	if !matchRule("anything", "", nil) {
		t.Errorf("empty rule should match all titles")
	}
	// plain substring
	if !matchRule("简日双语 1080p", "简日双语", nil) {
		t.Errorf("substring present should match")
	}
	if matchRule("简体内嵌 720p", "简日双语", nil) {
		t.Errorf("substring absent should not match")
	}
	// regex (whitelist: match to include; blacklist: match to exclude)
	re, err := compileFilter("/简体|1080p/")
	if err != nil {
		t.Fatalf("compileFilter() error = %v", err)
	}
	if !matchRule("简体内嵌 1080p", "/简体|1080p/", re) {
		t.Errorf("regex match should return true")
	}
	if matchRule("繁體 720p", "/简体|1080p/", re) {
		t.Errorf("regex non-match should return false")
	}
}

func TestMagnetItemsFromFeedFilters(t *testing.T) {
	feed := &gofeed.Feed{Items: []*gofeed.Item{
		{Title: "WEB-DL keep", Link: "https://mikanani.me/Home/Episode/1111111111111111111111111111111111111111"},
		{Title: "WEB-DL exclude", Link: "https://mikanani.me/Home/Episode/2222222222222222222222222222222222222222"},
		{Title: "BluRay keep", Link: "https://mikanani.me/Home/Episode/3333333333333333333333333333333333333333"},
	}}

	tests := []struct {
		name    string
		filter  string
		exclude string
		want    []string
	}{
		{
			name:   "empty exclude keeps all matching filter items",
			filter: "WEB-DL",
			want:   []string{"WEB-DL keep", "WEB-DL exclude"},
		},
		{
			name:    "exclude removes matching items after filter",
			filter:  "WEB-DL",
			exclude: "exclude",
			want:    []string{"WEB-DL keep"},
		},
		{
			name:   "regular expression filter",
			filter: "/WEB-DL|BluRay/",
			want:   []string{"WEB-DL keep", "WEB-DL exclude", "BluRay keep"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := magnetItemsFromFeed(feed, &Mikanani{}, &RssConfig{Filter: tt.filter, Exclude: tt.exclude})
			if len(items) != len(tt.want) {
				t.Fatalf("got %d items, want %d: %#v", len(items), len(tt.want), items)
			}
			for i, want := range tt.want {
				if items[i].Title != want {
					t.Errorf("item %d title = %q, want %q", i, items[i].Title, want)
				}
				if items[i].Magnet == "" {
					t.Errorf("item %d has an empty magnet", i)
				}
			}
		})
	}
}

func TestGenericMagnetExtraction(t *testing.T) {
	t.Parallel()
	generic := &Generic{}
	feed := &gofeed.Feed{Items: []*gofeed.Item{
		{Title: "magnet enclosure", Enclosures: []*gofeed.Enclosure{{URL: "magnet:?xt=urn:btih:aaaa&dn=name"}}},
		{Title: "description", Description: "download magnet:?xt=urn:btih:bbbb&dn=x here"},
		{Title: "ed2k", Enclosures: []*gofeed.Enclosure{{URL: "ed2k://|file|a.mkv|1|ABCD|/"}}},
		{Title: "torrent url", Enclosures: []*gofeed.Enclosure{{URL: "https://example.com/a.torrent"}}},
		{Title: "direct", Enclosures: []*gofeed.Enclosure{{URL: "https://example.com/a.mkv"}}},
		{Title: "guid hash", GUID: "0123456789abcdef0123456789abcdef01234567"},
	}}
	got := make([]string, 0, len(feed.Items))
	for _, item := range feed.Items {
		got = append(got, generic.GetMagnet(item))
	}
	want := []string{
		"magnet:?xt=urn:btih:aaaa",
		"magnet:?xt=urn:btih:bbbb",
		"ed2k://|file|a.mkv|1|ABCD|/",
		"",
		"https://example.com/a.mkv",
		"magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567",
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item %d magnet = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGlobalExcludeAndSkipX5(t *testing.T) {
	t.Parallel()
	feed := &gofeed.Feed{Items: []*gofeed.Item{
		{Title: "第 01 话 1080p", Enclosures: []*gofeed.Enclosure{{URL: "magnet:?xt=urn:btih:1"}}},
		{Title: "第 01.5 话 1080p", Enclosures: []*gofeed.Enclosure{{URL: "magnet:?xt=urn:btih:2"}}},
		{Title: "第 02 话 720p", Enclosures: []*gofeed.Enclosure{{URL: "magnet:?xt=urn:btih:3"}}},
	}}
	cfg := &RssConfig{GlobalExcludes: []string{"720p"}, SkipX5: true}
	items, err := magnetItemsFromFeedStrict(feed, &Generic{}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Title != "第 01 话 1080p" {
		t.Fatalf("items = %#v", items)
	}
}

func TestApplyDefaults(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{Offline: config.OfflineConfig{
		DefaultCid:      "999",
		DefaultSavePath: "番剧/${title}",
		GlobalExclude:   []string{"720p"},
		SkipX5:          true,
	}}
	rss := &RssConfig{Name: "药屋少女", Url: "https://mikanani.me/RSS/Bangumi?bangumiId=1"}
	ApplyDefaults(cfg, rss)
	if rss.Source != config.SourceMikan {
		t.Errorf("Source = %q, want mikan", rss.Source)
	}
	if rss.Cid != "999" {
		t.Errorf("Cid = %q, want 999", rss.Cid)
	}
	if rss.SavePath != "" {
		t.Errorf("SavePath = %q, want empty (default save path removed)", rss.SavePath)
	}
	if len(rss.GlobalExcludes) != 1 || !rss.SkipX5 {
		t.Errorf("global rules = %#v, skip5 = %v", rss.GlobalExcludes, rss.SkipX5)
	}

	explicit := &RssConfig{
		Name:     "x",
		Url:      "https://example.com/feed.xml",
		Cid:      "123",
		SavePath: "自定义",
	}
	ApplyDefaults(cfg, explicit)
	if explicit.Cid != "123" || explicit.SavePath != "自定义" {
		t.Errorf("explicit values were overwritten: %#v", explicit)
	}
}

func TestGetSiteForConfigUsesSource(t *testing.T) {
	t.Parallel()
	if _, ok := getSiteForConfig(&RssConfig{Source: config.SourceMikan, Url: "https://mirror.example.com/rss"}).(*Mikanani); !ok {
		t.Error("mikan source should use Mikanani handler even on a custom host")
	}
	if _, ok := getSite("https://unknown.example.com/rss").(*Generic); !ok {
		t.Error("unknown host should fall back to Generic")
	}
	if _, ok := getSiteForConfig(&RssConfig{Source: config.SourceOther, Url: "https://nyaa.si/?page=rss"}).(*Nyaa); !ok {
		t.Error("other source should keep legacy Nyaa handler")
	}
}

func TestAniRSSGroupFilterSyntax(t *testing.T) {
	raw := `{{今晚月色真美}}:1080[Pp]
{{今晚月色真美}}:简
{{今晚月色真美}}:日
{{今晚月色真美}}:内嵌
{{今晚月色真美}}:avc|Avc|AVC`
	rules, err := compileFilterRules(raw)
	if err != nil {
		t.Fatal(err)
	}
	good := "【今晚月色真美】[青之箱 第二季 / アオのハコ Season 2][03][WebRip][1080P][AVC-8bit][AAC][简日内嵌]"
	if !matchFilterRules(good, rules, true) {
		t.Error("expected simplified AVC hardsub title to match")
	}
	bad := "【今晚月色真美】[青之箱 第二季 / アオのハコ Season 2][03][WebRip][1080P][AVC-8bit][AAC][繁日內嵌]"
	if matchFilterRules(bad, rules, true) {
		t.Error("traditional title should not match simplified filter")
	}
	other := "【其他字幕组】[青之箱][03][1080P][AVC][简日内嵌]"
	if matchFilterRules(other, rules, true) {
		t.Error("other subgroup should not match group-specific filter")
	}
}
