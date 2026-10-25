package webui

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
)

func TestWeComNotificationPayload(t *testing.T) {
	var got struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errcode":0,"errmsg":"ok"}`))
	}))
	defer server.Close()

	payload := notifyPayload{Event: "test", Title: "AniFeed 测试通知", Status: "success", FeedCount: 1, Matched: 2, Added: 1, Timestamp: "2026-10-19T12:00:00+08:00"}
	err := postNotification(config.NotifyConfig{Type: "wecom", Webhook: server.URL}, payload)
	if err != nil {
		t.Fatal(err)
	}
	if got.MsgType != "text" || !strings.Contains(got.Text.Content, "AniFeed 测试通知") || !strings.Contains(got.Text.Content, "提交到 115：1") {
		t.Fatalf("payload = %+v", got)
	}
}

func TestWeComNotificationReportsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"errcode":93000,"errmsg":"invalid webhook url"}`))
	}))
	defer server.Close()
	err := postNotification(config.NotifyConfig{Type: "wecom", Webhook: server.URL}, notifyPayload{Title: "x", Status: "success"})
	if err == nil || !strings.Contains(err.Error(), "93000") {
		t.Fatalf("err = %v", err)
	}
}

func TestWeComAutoDetection(t *testing.T) {
	if !isWeComNotifier(config.NotifyConfig{Webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc"}) {
		t.Fatal("expected wecom URL auto-detection")
	}
	if isWeComNotifier(config.NotifyConfig{Webhook: "https://example.com/hook"}) {
		t.Fatal("generic URL should not be wecom")
	}
}

func TestNotificationContentSelection(t *testing.T) {
	payload := notifyPayload{Title: "运行完成", Status: "success", FeedCount: 1, Matched: 2, Added: 1, Timestamp: "2026-10-19T12:00:00+08:00", Results: []p115.FeedRunResult{{Name: "番剧", Matched: 2, Added: 1, AddedTitles: []string{"番剧 - 01"}}}}
	text := notificationText(config.NotifyConfig{Content: []string{"summary"}}, payload)
	if !strings.Contains(text, "状态：成功") || strings.Contains(text, "订阅更新") || strings.Contains(text, "时间：") {
		t.Fatalf("summary text = %q", text)
	}
	text = notificationText(config.NotifyConfig{Content: []string{"details", "time"}}, payload)
	if !strings.Contains(text, "订阅更新") || !strings.Contains(text, "番剧") || !strings.Contains(text, "番剧 - 01") || !strings.Contains(text, "时间：2026-10-19 12:00") {
		t.Fatalf("details text = %q", text)
	}
}

func TestNotificationIdleReminder(t *testing.T) {
	payload := notifyPayload{Title: "AniFeed 摸鱼提示", Status: "success", IdleDays: 7, IdleFeeds: []string{"青之箱（最后更新：2026-10-01）"}}
	text := notificationText(config.NotifyConfig{Content: []string{"summary"}}, payload)
	if !strings.Contains(text, "超过 7 天没有更新") || !strings.Contains(text, "青之箱") {
		t.Fatalf("idle text = %q", text)
	}
}

func TestNotifyTemplateRendering(t *testing.T) {
	fields := map[string]string{
		"action": "添加至云下载", "title": "冰之城墙 [tmdbid=283428]", "score": "7.4",
		"themoviedbName": "冰之城墙", "bgmUrl": "https://bgm.tv/subject/665390",
		"season": "2", "episode": "16", "subgroup": "LoliHouse",
		"currentEpisodeNumber": "3", "totalEpisodeNumber": "0",
		"year": "2026", "month": "10", "date": "1", "text": "已更新", "bgmEpisodeTitle": "S02E16",
	}
	text := renderNotifyTemplate(config.DefaultNotifyTemplate(), fields)
	for _, want := range []string{"添加至云下载", "冰之城墙 [tmdbid=283428]", "评分: 7.4", "LoliHouse", "进度: 3/0", "首播:  2026年10月1日", "BGM: https://bgm.tv/subject/665390"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in:\n%s", want, text)
		}
	}
	if got := bgmURL("https://anibt.net/rss/anime.xml?bgmId=665390&groupSlug=x"); got != "https://bgm.tv/subject/665390" {
		t.Fatalf("bgmURL = %q", got)
	}
}

func TestSendItemNotificationTemplate(t *testing.T) {
	var got string
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer webhook.Close()
	s, _ := newTestServerWith(t, nil)
	cfg := &config.Config{
		Notify: config.NotifyConfig{Webhook: webhook.URL, Type: "generic", OnSuccess: true, OnError: true, Template: config.DefaultNotifyTemplate()},
		Rename: config.RenameConfig{Enabled: true, Template: "${title} S${seasonFormat}E${episodeFormat}[${subgroup}]"},
		RSS:    map[string][]config.RssConfig{"anibt.net": {{Name: "冰之城墙", Url: "https://anibt.net/rss/anime.xml?bgmId=665390", RenameSeason: 2}}},
	}
	s.sendItemNotification(cfg, "冰之城墙", "https://anibt.net/rss/anime.xml?bgmId=665390", "[LoliHouse] 冰之城墙 - 16 [WebRip 1080p]", "添加至云下载", "", true)
	for _, want := range []string{"添加至云下载", "冰之城墙", "BGM: https://bgm.tv/subject/665390", "字幕组: LoliHouse", "集: 16", "事件: 冰之城墙 S02E16[LoliHouse]"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}

func TestEpisodeTitleFallback(t *testing.T) {
	s, _ := newTestServerWith(t, nil)
	cfg := &config.Config{}
	feed := &config.RssConfig{Name: "测试", Url: "https://anibt.net/rss/anime.xml?bgmId=665390"}
	if got := s.resolveEpisodeTitle(cfg, feed, 2, 16); got != "第16集" {
		t.Fatalf("title = %q", got)
	}
}

func TestItemNotificationAppliesEpisodeOffset(t *testing.T) {
	var got string
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer webhook.Close()
	s, _ := newTestServerWith(t, nil)
	cfg := &config.Config{
		Notify: config.NotifyConfig{Webhook: webhook.URL, Type: "generic", OnSuccess: true, Template: config.DefaultNotifyTemplate()},
		RSS:    map[string][]config.RssConfig{"anibt.net": {{Name: "青之箱", Url: "https://example.com/rss", RenameOffset: 25, TotalEpisodes: 37}}},
	}
	s.sendItemNotification(cfg, "青之箱", "https://example.com/rss", "【今晚月色真美】[青之箱 第二季][03][WebRip]", "添加至云下载", "", true)
	if !strings.Contains(got, "进度: 28/37") {
		t.Fatalf("progress missing in:\n%s", got)
	}
}

func TestNotificationErrorShowsURLAndReason(t *testing.T) {
	payload := notifyPayload{Title: "运行完成", Status: "error", FeedCount: 1, Results: []p115.FeedRunResult{{
		Name:  "青之箱",
		URL:   "https://anibt.net/rss/anime.xml?bgmId=548156",
		Error: `get rss from https://anibt.net/rss/anime.xml: Get "https://anibt.net/rss/anime.xml": dial tcp4: lookup anibt.net on 127.0.0.11:53: no such host`,
	}}}
	text := notificationText(config.NotifyConfig{Content: []string{"errors"}}, payload)
	for _, want := range []string{"RSS: https://anibt.net/rss/anime.xml?bgmId=548156", "原因：域名解析失败，无法解析 anibt.net"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "dial tcp4") {
		t.Fatalf("raw error leaked:\n%s", text)
	}
}
