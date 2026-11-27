package webui

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/rsssite"
)

type tmdbEpisode struct {
	Season  int
	Episode int
	AirDate time.Time
	Name    string
}

type offsetSuggestion struct {
	Season     int     `json:"season"`
	Offset     int     `json:"offset"`
	Samples    int     `json:"samples"`
	Total      int     `json:"total"`
	Confidence float64 `json:"confidence"`
}

// suggestOffset aligns RSS releases to TMDB episodes by air date. TMDB remains
// the source of truth; dates only select the episode to compare.
func (s *Server) suggestOffset(cfg *config.Config, feed *config.RssConfig, rssURL string) (*offsetSuggestion, error) {
	if cfg == nil || feed == nil || strings.TrimSpace(feed.TmdbID) == "" || cfg.TMDB.APIKey == "" {
		return nil, fmt.Errorf("缺少 TMDB ID 或 API Key")
	}
	episodes, err := s.tmdbEpisodes(cfg, feed.TmdbID)
	if err != nil {
		return nil, err
	}
	if len(episodes) == 0 {
		return nil, fmt.Errorf("TMDB 没有分集数据")
	}
	items := rsssite.GetMagnetItemList(&rsssite.RssConfig{
		Url: rssURL, Source: config.NormalizeSource(feed.Source, rssURL),
	})
	if len(items) == 0 {
		return nil, fmt.Errorf("RSS 没有条目")
	}
	// ANi is a reposting group with reliable timestamps. When its items are
	// present, prefer them for date matching instead of mixing groups.
	if ani := preferAniItems(items); len(ani) >= 2 {
		items = ani
	}
	hasSeasonHint := func(season int) bool {
		for _, ep := range episodes {
			if ep.Season == season {
				return true
			}
		}
		return false
	}
	counts := map[string]int{}
	total := 0
	for i := range items {
		item := &items[i]
		episode, epErr := p115.ExtractEpisode(item.Title)
		if epErr != nil || episode <= 0 || item.Published.IsZero() {
			continue
		}
		seasonHint := detectSeason(item.Title)
		restrictSeason := seasonHint > 0 && hasSeasonHint(seasonHint)
		best := tmdbEpisode{}
		bestDiff := 1000 * time.Hour
		for _, candidate := range episodes {
			if restrictSeason && candidate.Season != seasonHint {
				continue
			}
			diff := item.Published.Sub(candidate.AirDate)
			if diff < 0 {
				diff = -diff
			}
			if diff < bestDiff && diff <= 7*24*time.Hour {
				best = candidate
				bestDiff = diff
			}
		}
		if best.Episode == 0 {
			continue
		}
		key := fmt.Sprintf("%d|%d", best.Season, best.Episode-int(episode))
		counts[key]++
		total++
	}
	if total == 0 {
		return nil, fmt.Errorf("无法按播出日期匹配 TMDB 分集")
	}
	type candidate struct {
		key   string
		count int
	}
	var list []candidate
	for key, count := range counts {
		list = append(list, candidate{key: key, count: count})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].count > list[j].count })
	var season, offset int
	if _, err := fmt.Sscanf(list[0].key, "%d|%d", &season, &offset); err != nil {
		return nil, fmt.Errorf("偏移结果无效")
	}
	minSamples := cfg.Lab.AutoOffsetSamples
	if minSamples <= 0 {
		minSamples = 3
	}
	conf := float64(list[0].count) / float64(total)
	if list[0].count < minSamples {
		return nil, fmt.Errorf("有效匹配不足：%d/%d", list[0].count, minSamples)
	}
	if conf < 0.6 {
		return nil, fmt.Errorf("匹配一致性不足：%.0f%%", conf*100)
	}
	return &offsetSuggestion{Season: season, Offset: offset, Samples: list[0].count, Total: total, Confidence: conf}, nil
}

type tmdbEpisodesCacheEntry struct {
	episodes []tmdbEpisode
	at       time.Time
}

func (s *Server) tmdbEpisodes(cfg *config.Config, tmdbID string) ([]tmdbEpisode, error) {
	cacheKey := tmdbCacheKey(cfg.TMDB.APIKey, tmdbID)
	s.tmdbCacheMu.Lock()
	if entry, ok := s.tmdbCache[cacheKey]; ok && time.Since(entry.at) < 6*time.Hour {
		out := append([]tmdbEpisode(nil), entry.episodes...)
		s.tmdbCacheMu.Unlock()
		return out, nil
	}
	s.tmdbCacheMu.Unlock()
	values := map[string]string{"api_key": cfg.TMDB.APIKey, "language": "zh-CN"}
	var detail struct {
		Seasons []struct {
			SeasonNumber int `json:"season_number"`
		} `json:"seasons"`
	}
	if err := s.fetchTMDBQuery(cfg, "/3/tv/"+tmdbID, values, &detail); err != nil {
		return nil, err
	}
	var out []tmdbEpisode
	for _, season := range detail.Seasons {
		if season.SeasonNumber <= 0 {
			continue
		}
		var payload struct {
			Episodes []struct {
				EpisodeNumber int    `json:"episode_number"`
				AirDate       string `json:"air_date"`
				Name          string `json:"name"`
			} `json:"episodes"`
		}
		if err := s.fetchTMDBQuery(cfg, fmt.Sprintf("/3/tv/%s/season/%d", tmdbID, season.SeasonNumber), values, &payload); err != nil {
			continue
		}
		for _, ep := range payload.Episodes {
			air, _ := time.Parse("2006-01-02", ep.AirDate)
			if ep.EpisodeNumber > 0 && !air.IsZero() {
				out = append(out, tmdbEpisode{Season: season.SeasonNumber, Episode: ep.EpisodeNumber, AirDate: air, Name: ep.Name})
			}
		}
	}
	s.tmdbCacheMu.Lock()
	if s.tmdbCache == nil {
		s.tmdbCache = make(map[string]tmdbEpisodesCacheEntry)
	}
	s.tmdbCache[cacheKey] = tmdbEpisodesCacheEntry{episodes: append([]tmdbEpisode(nil), out...), at: time.Now()}
	s.saveTMDBDiskCacheLocked()
	s.tmdbCacheMu.Unlock()
	return out, nil
}

func (s *Server) handleOffsetSuggest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "配置未加载")
		return
	}
	var feed *config.RssConfig
	for _, feeds := range cfg.RSS {
		for i := range feeds {
			if strings.EqualFold(strings.TrimSpace(feeds[i].Name), strings.TrimSpace(req.Name)) {
				feed = &feeds[i]
				break
			}
		}
	}
	if feed == nil {
		writeAPIError(w, http.StatusNotFound, "未找到订阅")
		return
	}
	targetURL := strings.TrimSpace(req.URL)
	if targetURL == "" {
		targetURL = feed.Url
	}
	suggestion, err := s.suggestByAniList(cfg, feed, targetURL)
	if err != nil {
		suggestion, err = s.suggestOffset(cfg, feed, targetURL)
	}
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, suggestion)
}

func preferAniItems(items []rsssite.MagnetItem) []rsssite.MagnetItem {
	var out []rsssite.MagnetItem
	for _, item := range items {
		subgroup := strings.ToLower(strings.TrimSpace(p115.ExtractSubgroup(item.Title)))
		if subgroup == "ani" {
			out = append(out, item)
		}
	}
	return out
}
