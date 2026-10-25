package webui

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/rsssite"
	"github.com/zhifengle/anifeed/store"
)

type renamePreviewItem struct {
	Title        string  `json:"title"`
	Renamed      string  `json:"renamed"`
	Status       string  `json:"status"`
	Error        string  `json:"error,omitempty"`
	Episode      float64 `json:"episode,omitempty"`
	TmdbID       string  `json:"tmdbid,omitempty"`
	Selected     bool    `json:"selected"`
	Suggested    bool    `json:"suggested"`
	Filtered     bool    `json:"filtered"`
	FilterReason string  `json:"filter_reason,omitempty"`
	Score        int     `json:"score,omitempty"`
	SkipReason   string  `json:"skip_reason,omitempty"`
	Key          string  `json:"key"`
	Magnet       string  `json:"magnet"`
	Subgroup     string  `json:"subgroup"`
	Downloaded   bool    `json:"downloaded"`
	Allowed      bool    `json:"allowed"`
	Decided      bool    `json:"decided"`
}

type renamePreviewResponse struct {
	Feed           string              `json:"feed"`
	Template       string              `json:"template"`
	Enabled        bool                `json:"enabled"`
	Total          int                 `json:"total"`
	Matched        int                 `json:"matched"`
	SelectedCount  int                 `json:"selected_count"`
	SuggestedCount int                 `json:"suggested_count"`
	OnePerEpisode  bool                `json:"one_per_episode"`
	Items          []renamePreviewItem `json:"items"`
	Missing        []int               `json:"missing,omitempty"`
}

func (s *Server) handleRenamePreview(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusInternalServerError, "配置尚未加载")
		return
	}
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	feedURL := strings.TrimSpace(r.URL.Query().Get("url"))
	source := strings.TrimSpace(r.URL.Query().Get("source"))
	if name == "" && feedURL == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少订阅名称或 URL")
		return
	}
	feed, found := findFeedConfig(cfg, name, feedURL)
	if found {
		if name == "" {
			name = feed.Name
		}
		if feedURL == "" {
			feedURL = feed.Url
		}
	}
	if feedURL == "" {
		writeAPIError(w, http.StatusBadRequest, "订阅缺少 RSS 地址")
		return
	}
	if feed.Name == "" {
		feed.Name = name
	}
	if feedURL != "" {
		feed.Url = feedURL
	}
	if feed.Source == "" {
		feed.Source = source
	}
	query := r.URL.Query()
	if query.Has("source") && strings.TrimSpace(query.Get("source")) != "" {
		feed.Source = strings.TrimSpace(query.Get("source"))
	}
	if query.Has("filter") {
		feed.Filter = query.Get("filter")
	}
	if query.Has("exclude") {
		feed.Exclude = query.Get("exclude")
	}
	if query.Has("rename_season") {
		if season, err := strconv.Atoi(query.Get("rename_season")); err == nil && season >= 0 {
			feed.RenameSeason = season
		}
	}
	if query.Has("rename_offset") {
		if offset, err := strconv.Atoi(query.Get("rename_offset")); err == nil && offset >= -1000 && offset <= 1000 {
			feed.RenameOffset = offset
		}
	}
	if query.Has("tmdb_id") {
		feed.TmdbID = strings.TrimSpace(query.Get("tmdb_id"))
	}
	if query.Has("tmdb_type") {
		feed.TmdbType = strings.TrimSpace(query.Get("tmdb_type"))
	}
	rssConfig := rsssite.RssConfig{
		Name:         feed.Name,
		Url:          feed.Url,
		Source:       feed.Source,
		Cid:          feed.Cid,
		SavePath:     feed.SavePath,
		Filter:       feed.Filter,
		Exclude:      feed.Exclude,
		Expiration:   feed.Expiration,
		RenameSeason: feed.RenameSeason,
		RenameOffset: feed.RenameOffset,
		TmdbID:       feed.TmdbID,
		TmdbType:     feed.TmdbType,
	}
	rsssite.ApplyDefaults(cfg, &rssConfig)

	decisions := map[string]bool{}
	previewStore, storeOK := func() (*store.Store, bool) {
		db, dbErr := store.NewWithPath(s.databasePath())
		if dbErr != nil {
			return nil, false
		}
		return db, true
	}()
	if storeOK {
		defer previewStore.Close()
		if values, decErr := previewStore.PreviewDecisions(feed.Url); decErr == nil {
			decisions = values
		}
	}

	total, detailed, err := rsssite.GetMagnetItemsPreview(&rssConfig)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "读取 RSS 失败: "+err.Error())
		return
	}
	items := make([]rsssite.MagnetItem, 0, len(detailed))
	passedIndex := make([]int, len(detailed))
	for i, entry := range detailed {
		if entry.Filtered {
			passedIndex[i] = -1
			continue
		}
		passedIndex[i] = len(items)
		items = append(items, entry.Item)
	}
	seasonTitle := ""
	if len(items) > 0 {
		seasonTitle = items[0].Title
	}
	season := renameSeasonForFeedConfig(cfg, &feed, seasonTitle)
	tmdbID := strings.TrimSpace(feed.TmdbID)
	sourceSubgroup := sourceSubgroupName(cfg, feed.Name, feedURL, "")
	opts := p115.RenameOptions{
		Enabled:  true,
		Template: cfg.Rename.Template,
		Season:   season,
		Offset:   renameOffsetForFeed(cfg, &feed),
		TmdbID:   tmdbID,
		Subgroup: sourceSubgroup,
	}
	suggested := p115.SelectEpisodes(items, true)
	actual := p115.SelectEpisodes(items, cfg.Offline.OnePerEpisode)
	resp := renamePreviewResponse{
		Feed:          feed.Name,
		Template:      cfg.Rename.Template,
		Enabled:       cfg.Rename.Enabled,
		Total:         total,
		Matched:       len(items),
		OnePerEpisode: cfg.Offline.OnePerEpisode,
		Items:         make([]renamePreviewItem, 0, len(detailed)),
	}
	for detailedIndex, entry := range detailed {
		item := entry.Item
		preview := renamePreviewItem{
			Title:        item.Title,
			TmdbID:       tmdbID,
			Filtered:     entry.Filtered,
			FilterReason: entry.Reason,
			Key:          p115.ItemKey(item),
			Magnet:       item.Magnet,
			Subgroup:     sourceSubgroupName(cfg, feed.Name, feedURL, item.Title),
		}
		if entry.Filtered {
			preview.SkipReason = entry.Reason
		} else {
			index := passedIndex[detailedIndex]
			preview.Selected = actual[index].Selected
			preview.Suggested = suggested[index].Selected
			preview.Score = suggested[index].Score
			preview.SkipReason = suggested[index].Reason
			preview.Allowed = actual[index].Selected
			if preview.Selected {
				resp.SelectedCount++
			}
			if preview.Suggested {
				resp.SuggestedCount++
			}
		}
		if storeOK && item.Magnet != "" {
			preview.Downloaded = previewStore.HasItem(item.Magnet)
		}
		if allowed, ok := decisions[preview.Key]; ok {
			preview.Allowed = allowed
			preview.Decided = true
		}
		renamed, buildErr := p115.BuildRenameName(feed.Name, item.Title, opts)
		if buildErr != nil {
			preview.Status = "error"
			preview.Error = buildErr.Error()
		} else {
			// A BT title is not a filename; appending its last "." segment
			// would produce things like ". S2 - 01 [WebRip...]". The real
			// extension is preserved later against the actual file name.
			preview.Renamed = renamed
			preview.Status = "ready"
			if episode, err := p115.ExtractEpisode(item.Title); err == nil {
				preview.Episode = episode
			}
		}
		resp.Items = append(resp.Items, preview)
	}
	sort.SliceStable(resp.Items, func(i, j int) bool {
		left, right := resp.Items[i], resp.Items[j]
		leftEp, rightEp := left.Episode, right.Episode
		if leftEp <= 0 {
			leftEp = 1e9
		}
		if rightEp <= 0 {
			rightEp = 1e9
		}
		if leftEp != rightEp {
			return leftEp < rightEp
		}
		if left.Filtered != right.Filtered {
			return !left.Filtered
		}
		if left.Suggested != right.Suggested {
			return left.Suggested
		}
		return false
	})
	missingItems := make([]renamePreviewItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		if !item.Filtered {
			missingItems = append(missingItems, item)
		}
	}
	resp.Missing = missingEpisodes(missingItems)
	writeJSON(w, http.StatusOK, resp)
}

func findFeedConfig(cfg *config.Config, name, feedURL string) (config.RssConfig, bool) {
	if cfg == nil {
		return config.RssConfig{}, false
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if (feedURL != "" && feed.Url == feedURL) || (name != "" && feed.Name == name) {
				return feed, true
			}
		}
	}
	return config.RssConfig{}, false
}

func renameSeasonForFeedConfig(cfg *config.Config, feed *config.RssConfig, title string) int {
	if feed == nil {
		return renameSeasonForFeedTitle(cfg, "", "", title)
	}
	return renameSeasonForFeedTitle(cfg, feed.Name, feed.Url, title)
}

func missingEpisodes(items []renamePreviewItem) []int {
	seen := map[int]bool{}
	min, max := 0, 0
	for _, item := range items {
		if item.Status != "ready" || item.Episode <= 0 {
			continue
		}
		episode := int(item.Episode)
		seen[episode] = true
		if min == 0 || episode < min {
			min = episode
		}
		if episode > max {
			max = episode
		}
	}
	if max <= 0 || len(seen) < 2 || max > 500 {
		return nil
	}
	var missing []int
	for episode := min + 1; episode < max; episode++ {
		if !seen[episode] {
			missing = append(missing, episode)
		}
	}
	return missing
}
