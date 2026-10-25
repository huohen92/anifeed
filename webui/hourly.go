package webui

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/store"
)

func (s *Server) startHourlyNotifyLoop(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			if now.Minute() != 0 {
				continue
			}
			cfg := s.configSnapshot()
			if cfg == nil || !cfg.Notify.HourlyEnabled || !notifyEventEnabled(cfg.Notify, "hourly") || strings.TrimSpace(cfg.Notify.Webhook) == "" {
				continue
			}
			if len(cfg.Notify.HourlyHours) > 0 && !containsHour(cfg.Notify.HourlyHours, now.Hour()) {
				continue
			}
			key := now.Format("2006-01-02 15")
			s.hourlyMu.Lock()
			if s.lastHourly == key {
				s.hourlyMu.Unlock()
				continue
			}
			s.lastHourly = key
			s.hourlyMu.Unlock()
			go s.sendHourlySummary()
		}
	}
}

func (s *Server) sendHourlySummary() {
	cfg := s.configSnapshot()
	if cfg == nil || !cfg.Notify.HourlyEnabled || !notifyEventEnabled(cfg.Notify, "hourly") || strings.TrimSpace(cfg.Notify.Webhook) == "" {
		return
	}
	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		log.Printf("[warning] 整点推送读取数据库失败: %v", err)
		return
	}
	defer st.Close()
	history, err := st.ListHistory(500, "")
	if err != nil {
		log.Printf("[warning] 整点推送读取历史失败: %v", err)
		return
	}
	now := time.Now()
	today := 0
	latest := map[string]store.HistoryItem{}
	maxEpisode := map[string]float64{}
	for _, item := range history {
		if sameLocalDate(item.CreatedAt, now) {
			today++
		}
		if _, ok := latest[item.ConfigName]; !ok {
			latest[item.ConfigName] = item
		}
		if episode, parseErr := p115.ExtractEpisode(item.Title); parseErr == nil && episode > maxEpisode[item.ConfigName] {
			maxEpisode[item.ConfigName] = episode
		}
	}
	enabled := 0
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if feed.IsEnabled() {
				enabled++
			}
		}
	}
	pending, _ := st.CountPendingPostTasks()
	lines := []string{
		"AniFeed 整点更新",
		"时间: " + now.Format("2006-01-02 15:04"),
		fmt.Sprintf("订阅: %d 个，启用 %d 个", countFeeds(cfg), enabled),
		fmt.Sprintf("今天提交: %d 个任务", today),
		fmt.Sprintf("等待下载完成: %d 个", pending),
		"",
		"订阅更新:",
	}
	added := 0
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if !feed.IsEnabled() {
				continue
			}
			item, ok := latest[feed.Name]
			progress := "—"
			lastUpdate := "—"
			if ok {
				lastUpdate = item.CreatedAt.Format("01-02 15:04")
				episode := maxEpisode[feed.Name]
				if episode <= 0 {
					if parsed, epErr := p115.ExtractEpisode(item.Title); epErr == nil {
						episode = parsed
					}
				}
				if episode > 0 {
					episode += float64(renameOffsetForFeed(cfg, &feed))
					if episode < 0 {
						episode = 0
					}
					progress = formatEpisodeText(episode)
					total := feed.TotalEpisodes
					if total <= 0 {
						total = s.loadNotifyMeta(cfg, &feed).Total
					}
					if total > 0 {
						progress += fmt.Sprintf("/%d", total)
					}
				}
			}
			lines = append(lines, fmt.Sprintf("- %s：进度 %s，最近更新 %s", feed.Name, progress, lastUpdate))
			added++
			if added >= 20 {
				break
			}
		}
		if added >= 20 {
			break
		}
	}
	if added == 0 {
		lines = append(lines, "- 暂无更新记录")
	}
	if err := postRawNotification(cfg.Notify, strings.Join(lines, "\n")); err != nil {
		log.Printf("[warning] 整点推送发送失败: %v", err)
		return
	}
	log.Printf("[notify] 整点更新推送已发送")
}

func sameLocalDate(value, now time.Time) bool {
	y1, m1, d1 := value.Local().Date()
	y2, m2, d2 := now.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func countFeeds(cfg *config.Config) int {
	n := 0
	for _, feeds := range cfg.RSS {
		n += len(feeds)
	}
	return n
}

func containsHour(hours []int, target int) bool {
	for _, hour := range hours {
		if hour == target {
			return true
		}
	}
	return false
}
