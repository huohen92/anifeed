package webui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/request"
	"github.com/huohen92/anifeed/store"
)

type notifyPayload struct {
	Event     string               `json:"event"`
	Title     string               `json:"title"`
	Status    string               `json:"status"`
	FeedCount int                  `json:"feed_count"`
	Matched   int                  `json:"matched"`
	Added     int                  `json:"added"`
	Error     string               `json:"error,omitempty"`
	Timestamp string               `json:"timestamp"`
	Results   []p115.FeedRunResult `json:"results,omitempty"`
	IdleFeeds []string             `json:"idle_feeds,omitempty"`
	IdleDays  int                  `json:"idle_days,omitempty"`
}

// startRunAllJob runs every enabled feed and appends notifications and
// optional task cleanup.
func (s *Server) startRunAllJob(ag *p115.Agent, kind, title string) (Job, error) {
	return s.jobs.start(kind, title, func() ([]p115.FeedRunResult, error) {
		results, err := ag.RunAllRssTasks()
		s.afterRun(kind, ag, results, err)
		return results, err
	})
}

// startRunFeedJob runs a single feed with the same post-processing as a full run.
func (s *Server) startRunFeedJob(ag *p115.Agent, name string, run func() []p115.FeedRunResult) (Job, error) {
	return s.jobs.start("feed", "运行订阅："+name, func() ([]p115.FeedRunResult, error) {
		results := run()
		var err error
		for _, result := range results {
			if result.Error != "" {
				err = fmt.Errorf("%s", result.Error)
				break
			}
		}
		s.afterRun("feed", ag, results, err)
		return results, nil
	})
}

func (s *Server) afterRun(kind string, ag *p115.Agent, results []p115.FeedRunResult, runErr error) {
	cfg := s.configSnapshot()
	if cfg == nil {
		return
	}

	if ag != nil && (kind == "all" || kind == "auto") && cfg.Offline.AutoClearCompleted && runErr == nil {
		if postProcessEnabled(cfg) {
			log.Printf("[auto] 重命名/离线后转移启用，跳过自动清理，等待后处理完成")
		} else if err := ag.OfflineClear(0); err != nil {
			log.Printf("[warning] 自动清理已完成任务失败: %v", err)
		} else {
			log.Printf("[auto] 已清理 115 已完成离线任务")
		}
	}

	payload := buildNotifyPayload(kind, results, runErr)
	skipNoChange := cfg.Notify.SkipNoChange && payload.Added == 0 && payload.Status == "success"
	if !skipNoChange && notifyEventEnabled(cfg.Notify, "run") {
		s.sendNotification(cfg.Notify, payload)
	}
	if payload.Added > 0 {
		go s.sendAddedNotifications(cfg, results)
	}
	if kind == "auto" || kind == "all" || kind == "feed" {
		if cfg.Offline.AutoDisableCompleted && ag != nil {
			go func() {
				time.Sleep(3 * time.Second)
				s.autoDisableCompletedFeeds(ag)
			}()
		}
	}
	if kind == "auto" || kind == "all" {
		go s.checkIdleFeeds(cfg)
	}
}

func (s *Server) checkIdleFeeds(cfg *config.Config) {
	if cfg == nil || !cfg.Notify.IdleEnabled || !notifyEventEnabled(cfg.Notify, "idle") || strings.TrimSpace(cfg.Notify.Webhook) == "" {
		return
	}
	days := cfg.Notify.IdleDays
	if days <= 0 {
		days = 7
	}
	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		return
	}
	defer st.Close()
	now := time.Now()
	dayKey := now.Format("2006-01-02")
	var reminders []string
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if !feed.IsEnabled() {
				continue
			}
			last, ok, lookupErr := st.LatestHistoryForFeed(feed.Url)
			if lookupErr != nil || !ok || now.Sub(last) < time.Duration(days)*24*time.Hour {
				continue
			}
			first, markErr := st.MarkIdleReminder(feed.Url, dayKey)
			if markErr != nil || !first {
				continue
			}
			reminders = append(reminders, fmt.Sprintf("%s（最后更新：%s）", feed.Name, last.Format("2006-01-02")))
		}
	}
	if len(reminders) == 0 {
		return
	}
	payload := notifyPayload{
		Event:     "idle",
		Title:     "AniFeed 摸鱼提示",
		Status:    "success",
		FeedCount: len(reminders),
		Timestamp: now.Format(time.RFC3339),
		IdleFeeds: reminders,
		IdleDays:  days,
	}
	if err := postNotification(cfg.Notify, payload); err != nil {
		log.Printf("[warning] 摸鱼提示发送失败: %v", err)
		return
	}
	log.Printf("[notify] 摸鱼提示已发送")
}

func buildNotifyPayload(kind string, results []p115.FeedRunResult, runErr error) notifyPayload {
	payload := notifyPayload{
		Event:     "run",
		Status:    "success",
		FeedCount: len(results),
		Timestamp: time.Now().Format(time.RFC3339),
		Results:   results,
	}
	var failures []string
	for _, result := range results {
		payload.Matched += result.Matched
		payload.Added += result.Added
		if result.Error != "" {
			failures = append(failures, result.Name+": "+result.Error)
		}
	}
	switch kind {
	case "all", "auto":
		payload.Title = "AniFeed 全部订阅运行完成"
	default:
		payload.Title = "AniFeed 订阅运行完成"
		if len(results) == 1 {
			payload.Title = "AniFeed 订阅运行完成：" + results[0].Name
		}
	}
	if runErr != nil {
		payload.Status = "error"
		payload.Error = runErr.Error()
	} else if len(failures) > 0 {
		payload.Status = "error"
		payload.Error = strings.Join(failures, "; ")
	}
	return payload
}

func notifyEventEnabled(cfg config.NotifyConfig, event string) bool {
	events := cfg.Events
	if len(events) == 0 {
		events = config.DefaultNotifyEvents()
	}
	for _, value := range events {
		if value == event {
			return true
		}
	}
	return false
}

func (s *Server) sendNotification(cfg config.NotifyConfig, payload notifyPayload) {
	if cfg.Webhook == "" {
		return
	}
	if payload.Status == "success" && !cfg.OnSuccess {
		return
	}
	if payload.Status == "error" && !cfg.OnError {
		return
	}
	if err := postNotification(cfg, payload); err != nil {
		log.Printf("[warning] 通知发送失败: %v", err)
		return
	}
	log.Printf("[notify] 通知已发送")
}

func (s *Server) sendTestNotification() error {
	cfg := s.configSnapshot()
	if cfg == nil || cfg.Notify.Webhook == "" {
		return fmt.Errorf("尚未配置 Webhook 地址")
	}
	payload := notifyPayload{
		Event:     "test",
		Title:     "AniFeed 测试通知",
		Status:    "success",
		FeedCount: 0,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return postNotification(cfg.Notify, payload)
}

func postNotification(cfg config.NotifyConfig, payload notifyPayload) error {
	content := notificationText(cfg, payload)
	if isWeComNotifier(cfg) {
		return postWeCom(cfg.Webhook, content)
	}
	return postPlainWebhook(cfg.Webhook, content)
}

func isWeComNotifier(cfg config.NotifyConfig) bool {
	switch strings.ToLower(strings.TrimSpace(cfg.Type)) {
	case "wecom":
		return true
	case "generic":
		return false
	}
	u, err := url.Parse(cfg.Webhook)
	if err != nil {
		return false
	}
	return strings.EqualFold(u.Hostname(), "qyapi.weixin.qq.com") && strings.Contains(u.Path, "/cgi-bin/webhook/send")
}

func notificationText(cfg config.NotifyConfig, payload notifyPayload) string {
	sections := make(map[string]bool, len(cfg.Content))
	if len(cfg.Content) == 0 {
		for _, key := range config.DefaultNotifyContent() {
			sections[key] = true
		}
	} else {
		for _, key := range cfg.Content {
			sections[key] = true
		}
	}
	status := "成功"
	if payload.Status != "success" {
		status = "失败"
	}
	lines := []string{payload.Title}
	if len(payload.IdleFeeds) > 0 {
		lines = append(lines, "", fmt.Sprintf("以下订阅超过 %d 天没有更新：", payload.IdleDays))
		for _, feed := range payload.IdleFeeds {
			lines = append(lines, "- "+feed)
		}
	}
	if sections["summary"] {
		lines = append(lines,
			"状态："+status,
			fmt.Sprintf("检查订阅：%d 个", payload.FeedCount),
			fmt.Sprintf("发现新内容：%d 条", payload.Matched),
			fmt.Sprintf("提交到 115：%d 个", payload.Added),
		)
	}
	if sections["time"] && payload.Timestamp != "" {
		lines = append(lines, "时间："+formatNotificationTime(payload.Timestamp))
	}
	if sections["details"] && len(payload.Results) > 0 {
		lines = append(lines, "", "订阅更新：")
		for _, result := range payload.Results {
			name := result.Name
			if name == "" {
				name = result.URL
			}
			line := fmt.Sprintf("- %s：发现 %d 条，提交 %d 个", name, result.Matched, result.Added)
			if result.Cached > 0 {
				line += fmt.Sprintf("，以前下载过 %d 条", result.Cached)
			}
			if result.Skipped > 0 {
				line += fmt.Sprintf("，重复版本跳过 %d 条", result.Skipped)
			}
			if result.Rejected > 0 {
				line += fmt.Sprintf("，115 未接受 %d 条", result.Rejected)
			}
			lines = append(lines, line)
			for _, title := range result.AddedTitles {
				title = strings.TrimSpace(title)
				if len([]rune(title)) > 80 {
					title = string([]rune(title)[:80]) + "…"
				}
				lines = append(lines, "  · "+title)
			}
		}
	}
	if sections["errors"] {
		errorResults := make([]p115.FeedRunResult, 0)
		for _, result := range payload.Results {
			if result.Error != "" {
				errorResults = append(errorResults, result)
			}
		}
		if len(errorResults) > 0 {
			lines = append(lines, "", "错误：")
			for _, result := range errorResults {
				name := result.Name
				if name == "" {
					name = result.URL
				}
				lines = append(lines, "- "+name)
				if result.URL != "" {
					lines = append(lines, "  RSS: "+result.URL)
				}
				lines = append(lines, "  原因："+friendlyNotifyError(result.Error, result.URL))
			}
		} else if payload.Error != "" {
			lines = append(lines, "", "错误：", "- 原因："+friendlyNotifyError(payload.Error, ""))
		}
	}
	content := strings.Join(lines, "\n")
	if len(content) > 2000 {
		content = truncateUTF8(content, 2000) + "…"
	}
	return content
}

func friendlyNotifyError(message, rawURL string) string {
	lower := strings.ToLower(message)
	host := ""
	if u, err := url.Parse(rawURL); err == nil {
		host = u.Hostname()
	}
	if strings.Contains(lower, "no such host") || strings.Contains(lower, "server misbehaving") || strings.Contains(lower, "name resolution") {
		if host == "" {
			host = "RSS 域名"
		}
		return "域名解析失败，无法解析 " + host
	}
	if strings.Contains(lower, "network is unreachable") || strings.Contains(lower, "no route to host") {
		return "网络不可达，请检查服务器网络或代理设置"
	}
	if strings.Contains(lower, "connection refused") {
		return "连接被拒绝，目标服务器或代理没有响应"
	}
	if strings.Contains(lower, "i/o timeout") || strings.Contains(lower, "deadline exceeded") || strings.Contains(lower, "timeout") {
		return "连接超时，请检查网络或稍后重试"
	}
	if strings.Contains(lower, "certificate") || strings.Contains(lower, "tls") || strings.Contains(lower, "x509") {
		return "HTTPS 证书校验失败"
	}
	if strings.Contains(lower, "invalid filter") {
		return message
	}
	// Strip the repetitive Go/RSS wrapper and keep the meaningful tail.
	if index := strings.LastIndex(message, ": "); index >= 0 && index+2 < len(message) {
		tail := strings.TrimSpace(message[index+2:])
		if tail != "" && len([]rune(tail)) < 180 {
			return tail
		}
	}
	if host != "" {
		return "请求 RSS 失败（" + host + "）"
	}
	return "请求 RSS 失败"
}

func formatNotificationTime(value string) string {
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.Local().Format("2006-01-02 15:04")
	}
	return value
}

func truncateUTF8(value string, maxBytes int) string {
	if len(value) <= maxBytes {
		return value
	}
	runes := []rune(value)
	for len(string(runes)) > maxBytes {
		runes = runes[:len(runes)-1]
	}
	return string(runes)
}

func postWeCom(webhook, content string) error {
	body, err := json.Marshal(map[string]any{
		"msgtype": "text",
		"text":    map[string]string{"content": content},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "AniFeed-wecom/1.0")
	client := request.NewIPv4Client(10 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("企业微信 HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("企业微信返回格式错误: %w", err)
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("企业微信返回错误 %d: %s", result.ErrCode, result.ErrMsg)
	}
	return nil
}

func postPlainWebhook(webhook, content string) error {
	req, err := http.NewRequest(http.MethodPost, webhook, strings.NewReader(content))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("User-Agent", "AniFeed-webhook/1.0")
	client := request.NewIPv4Client(10 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

func (s *Server) autoDisableCompletedFeeds(ag *p115.Agent) {
	cfg := s.configSnapshot()
	if cfg == nil || ag == nil || ag.StoreInstance == nil {
		return
	}
	history, err := ag.StoreInstance.ListHistory(1000, "")
	if err != nil {
		return
	}
	latest := make(map[string]float64)
	for _, item := range history {
		episode, epErr := p115.ExtractEpisode(item.Title)
		if epErr == nil && episode > latest[item.ConfigName] {
			latest[item.ConfigName] = episode
		}
	}
	targets := make(map[string]bool)
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if feed.TotalEpisodes <= 0 || feed.Completed || (feed.Enabled != nil && !*feed.Enabled) {
				continue
			}
			offset := feed.RenameOffset
			if offset == 0 {
				offset = cfg.Rename.Offset
			}
			if latest[feed.Name]+float64(offset) >= float64(feed.TotalEpisodes) {
				targets[feed.Name] = true
			}
		}
	}
	if len(targets) == 0 {
		return
	}
	resp, err := s.loadEditable()
	if err != nil {
		return
	}
	changed := false
	for i := range resp.Config.RSS {
		if targets[resp.Config.RSS[i].Name] {
			resp.Config.RSS[i].Enabled = false
			resp.Config.RSS[i].Completed = true
			changed = true
		}
	}
	if !changed {
		return
	}
	if _, err := s.saveEditable(&resp.Config); err != nil {
		log.Printf("[auto] 自动停用已完成订阅失败: %v", err)
		return
	}
	log.Printf("[auto] 已将 %d 个达到总集数的订阅标记完成并停用", len(targets))
}
