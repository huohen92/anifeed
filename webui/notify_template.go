package webui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/request"
	"github.com/zhifengle/anifeed/store"
)

type notifyMetaEntry struct {
	Name     string
	Score    float64
	FirstAir string
	Total    int
	Expires  time.Time
}

type notifyMeta struct {
	Name     string
	Score    float64
	FirstAir string
	Total    int
}

func (s *Server) sendAddedNotifications(cfg *config.Config, results []p115.FeedRunResult) {
	if cfg == nil || cfg.Notify.Webhook == "" {
		return
	}
	for _, result := range results {
		var feed *config.RssConfig
		if found, ok := findFeedConfig(cfg, result.Name, result.URL); ok {
			feed = &found
		}
		for _, title := range result.AddedTitles {
			feedName := result.Name
			feedURL := result.URL
			if feed != nil {
				feedName = feed.Name
				feedURL = feed.Url
			}
			s.sendItemNotification(cfg, feedName, feedURL, title, "添加至云下载", "", true)
		}
	}
}

func (s *Server) sendItemNotification(cfg *config.Config, feedName, feedURL, itemTitle, action, eventName string, success bool) {
	if cfg == nil || cfg.Notify.Webhook == "" {
		return
	}
	if success && !cfg.Notify.OnSuccess {
		return
	}
	if !success && !cfg.Notify.OnError {
		return
	}
	event := "added"
	switch action {
	case "下载完成":
		event = "completed"
	case "下载失败":
		event = "failed"
	}
	if !notifyEventEnabled(cfg.Notify, event) {
		return
	}
	template := cfg.Notify.Template
	if strings.TrimSpace(template) == "" {
		template = config.DefaultNotifyTemplate()
	}
	feed := config.RssConfig{Name: feedName, Url: feedURL}
	if found, ok := findFeedConfig(cfg, feedName, feedURL); ok {
		feed = found
	}
	meta := s.loadNotifyMeta(cfg, &feed)

	title := feed.Name
	if title == "" {
		title = itemTitle
	}
	tmdbID := strings.TrimSpace(feed.TmdbID)
	if tmdbID != "" {
		title = fmt.Sprintf("%s [tmdbid=%s]", title, tmdbID)
	}
	season := renameSeasonForFeed(cfg, feed.Name, feed.Url)
	episode, _ := p115.ExtractEpisode(itemTitle)
	episodeTitle := s.resolveEpisodeTitle(cfg, &feed, season, int(episode))
	subgroup := sourceSubgroupName(cfg, feedName, feedURL, itemTitle)
	total := meta.Total
	if feed.TotalEpisodes > 0 {
		total = feed.TotalEpisodes
	}
	currentEpisode := int(episode)
	if currentEpisode > 0 {
		currentEpisode += renameOffsetForFeed(cfg, &feed)
		if currentEpisode < 0 {
			currentEpisode = 0
		}
	}
	firstAir := meta.FirstAir
	year, month, date := "", "", ""
	if parsed, err := time.Parse("2006-01-02", firstAir); err == nil {
		year = strconv.Itoa(parsed.Year())
		month = strconv.Itoa(int(parsed.Month()))
		date = strconv.Itoa(parsed.Day())
	}

	if strings.TrimSpace(eventName) == "" {
		eventName = predictedRenameName(cfg, &feed, itemTitle)
	}
	fields := map[string]string{
		"action":         action,
		"title":          title,
		"score":          "",
		"themoviedbName": meta.Name,
		"bgmUrl": func() string {
			if strings.TrimSpace(feed.BgmURL) != "" {
				return strings.TrimSpace(feed.BgmURL)
			}
			return bgmURL(feed.Url)
		}(),
		"season":               strconv.Itoa(season),
		"episode":              formatEpisodeText(episode),
		"subgroup":             subgroup,
		"currentEpisodeNumber": strconv.Itoa(currentEpisode),
		"totalEpisodeNumber":   strconv.Itoa(total),
		"year":                 year,
		"month":                month,
		"date":                 date,
		"text":                 eventName,
		"bgmEpisodeTitle":      episodeTitle,
	}
	if meta.Score > 0 {
		fields["score"] = strconv.FormatFloat(meta.Score, 'f', 1, 64)
	}
	if fields["themoviedbName"] == "" {
		fields["themoviedbName"] = feed.Name
	}
	rendered := renderNotifyTemplate(template, fields)
	if err := postRawNotification(cfg.Notify, rendered); err != nil {
		logNotifyError(err)
	}
}

func (s *Server) loadNotifyMeta(cfg *config.Config, feed *config.RssConfig) notifyMeta {
	if cfg == nil || feed == nil || strings.TrimSpace(feed.TmdbID) == "" || cfg.TMDB.APIKey == "" {
		return notifyMeta{Name: feed.Name}
	}
	key := "tv:" + feed.TmdbID
	s.coverMu.Lock()
	cached, ok := s.notifyMetaCache[key]
	s.coverMu.Unlock()
	if ok && time.Now().Before(cached.Expires) {
		return notifyMeta{Name: cached.Name, Score: cached.Score, FirstAir: cached.FirstAir, Total: cached.Total}
	}
	values := url.Values{}
	values.Set("api_key", cfg.TMDB.APIKey)
	values.Set("language", "zh-CN")
	var detail struct {
		Name             string  `json:"name"`
		VoteAverage      float64 `json:"vote_average"`
		FirstAirDate     string  `json:"first_air_date"`
		NumberOfEpisodes int     `json:"number_of_episodes"`
		Seasons          []struct {
			EpisodeCount int `json:"episode_count"`
		} `json:"seasons"`
	}
	meta := notifyMeta{Name: feed.Name}
	if err := s.fetchTMDB(cfg.TMDB, "/3/tv/"+url.PathEscape(feed.TmdbID), values, &detail); err == nil {
		meta.Name = detail.Name
		meta.Score = detail.VoteAverage
		meta.FirstAir = detail.FirstAirDate
		meta.Total = detail.NumberOfEpisodes
		if meta.Total <= 0 {
			for _, season := range detail.Seasons {
				meta.Total += season.EpisodeCount
			}
		}
		if meta.Name == "" {
			meta.Name = feed.Name
		}
		entry := notifyMetaEntry{Name: meta.Name, Score: meta.Score, FirstAir: meta.FirstAir, Total: meta.Total, Expires: time.Now().Add(6 * time.Hour)}
		s.coverMu.Lock()
		if s.notifyMetaCache == nil {
			s.notifyMetaCache = make(map[string]notifyMetaEntry)
		}
		s.notifyMetaCache[key] = entry
		s.coverMu.Unlock()
	}
	return meta
}

func predictedRenameName(cfg *config.Config, feed *config.RssConfig, itemTitle string) string {
	if cfg == nil || feed == nil {
		return itemTitle
	}
	opts := p115.RenameOptions{
		Enabled:  cfg.Rename.Enabled,
		Template: cfg.Rename.Template,
		Season:   renameSeasonForFeed(cfg, feed.Name, feed.Url),
		Offset:   renameOffsetForFeed(cfg, feed),
		TmdbID:   feed.TmdbID,
	}
	if name, err := p115.BuildRenameName(feed.Name, itemTitle, opts); err == nil && strings.TrimSpace(name) != "" {
		return name
	}
	return itemTitle
}

func feedForName(cfg *config.Config, name string) *config.RssConfig {
	if cfg == nil {
		return nil
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if feed.Name == name {
				copy := feed
				return &copy
			}
		}
	}
	return nil
}

func postEventName(entries []store.PostEntry, fallback string) string {
	for _, entry := range entries {
		if p115.LooksLikeMediaFile(entry.Name) {
			return entry.Name
		}
	}
	if len(entries) > 0 && strings.TrimSpace(entries[0].Name) != "" {
		return entries[0].Name
	}
	return fallback
}

func (s *Server) resolveEpisodeTitle(cfg *config.Config, feed *config.RssConfig, season, episode int) string {
	if episode <= 0 {
		return "第" + formatEpisodeText(float64(episode)) + "集"
	}
	fallback := fmt.Sprintf("第%d集", episode)
	if cfg == nil || feed == nil || cfg.TMDB.APIKey == "" {
		return fallback
	}
	key := fmt.Sprintf("%s:%d:%d", feed.Url, season, episode)
	s.coverMu.Lock()
	if title, ok := s.episodeTitles[key]; ok {
		s.coverMu.Unlock()
		return title
	}
	s.coverMu.Unlock()
	title := ""
	if id := bgmIDFromURL(feed.Url); id != "" {
		title = fetchBgmEpisodeTitle(id, episode)
	}
	if strings.TrimSpace(title) == "" && strings.TrimSpace(feed.TmdbID) != "" {
		values := url.Values{}
		values.Set("api_key", cfg.TMDB.APIKey)
		values.Set("language", "zh-CN")
		var data struct {
			Name string `json:"name"`
		}
		path := fmt.Sprintf("/3/tv/%s/season/%d/episode/%d", url.PathEscape(feed.TmdbID), maxInt(season, 1), episode)
		if err := s.fetchTMDB(cfg.TMDB, path, values, &data); err == nil {
			title = strings.TrimSpace(data.Name)
		}
	}
	if strings.TrimSpace(title) == "" {
		title = fallback
	}
	s.coverMu.Lock()
	if s.episodeTitles == nil {
		s.episodeTitles = make(map[string]string)
	}
	s.episodeTitles[key] = title
	s.coverMu.Unlock()
	return title
}

func fetchBgmEpisodeTitle(subjectID string, episode int) string {
	client := request.NewIPv4Client(6 * time.Second)
	raw := fmt.Sprintf("https://api.bgm.tv/v0/episodes?subject_id=%s&type=0&limit=100&offset=0", url.QueryEscape(subjectID))
	req, err := http.NewRequest(http.MethodGet, raw, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "AniFeed/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	var payload struct {
		Data []struct {
			Ep     int    `json:"ep"`
			Sort   int    `json:"sort"`
			Name   string `json:"name"`
			NameCN string `json:"name_cn"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ""
	}
	for _, item := range payload.Data {
		if item.Ep != episode && item.Sort != episode {
			continue
		}
		if strings.TrimSpace(item.NameCN) != "" {
			return item.NameCN
		}
		return strings.TrimSpace(item.Name)
	}
	return ""
}

func bgmIDFromURL(feedURL string) string {
	u, err := url.Parse(feedURL)
	if err != nil {
		return ""
	}
	id := strings.TrimSpace(u.Query().Get("bgmId"))
	if id == "" {
		id = strings.TrimSpace(u.Query().Get("bangumiId"))
	}
	return id
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func bgmURL(feedURL string) string {
	u, err := url.Parse(feedURL)
	if err != nil {
		return ""
	}
	id := strings.TrimSpace(u.Query().Get("bgmId"))
	if id == "" {
		id = strings.TrimSpace(u.Query().Get("bangumiId"))
	}
	if id == "" {
		return ""
	}
	return "https://bgm.tv/subject/" + id
}

func formatEpisodeText(episode float64) string {
	if episode <= 0 {
		return "0"
	}
	if episode == float64(int(episode)) {
		return strconv.Itoa(int(episode))
	}
	return strconv.FormatFloat(episode, 'f', -1, 64)
}

func renderNotifyTemplate(template string, fields map[string]string) string {
	for key, value := range fields {
		template = strings.ReplaceAll(template, "${"+key+"}", value)
	}
	return strings.TrimSpace(template)
}

func postRawNotification(cfg config.NotifyConfig, text string) error {
	if len(text) > 2000 {
		text = truncateUTF8(text, 2000) + "…"
	}
	if isWeComNotifier(cfg) {
		return postWeCom(cfg.Webhook, text)
	}
	return postPlainWebhook(cfg.Webhook, text)
}

func logNotifyError(err error) {
	log.Printf("[warning] 条目通知发送失败: %v", err)
}
