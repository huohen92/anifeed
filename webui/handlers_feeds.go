package webui

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/rsssite"
)

type feedPayload struct {
	Site            string                `json:"site"`
	Name            string                `json:"name"`
	SubgroupName    string                `json:"subgroup_name"`
	BgmURL          string                `json:"bgm_url"`
	URL             string                `json:"url"`
	Source          string                `json:"source"`
	Enabled         bool                  `json:"enabled"`
	Cid             string                `json:"cid"`
	SavePath        string                `json:"savepath"`
	Filter          string                `json:"filter"`
	Exclude         string                `json:"exclude"`
	Expiration      uint                  `json:"expiration"`
	RenameSeason    int                   `json:"rename_season"`
	RenameOffset    int                   `json:"rename_offset"`
	UpdateWeekday   int                   `json:"update_weekday"`
	TotalEpisodes   int                   `json:"total_episodes"`
	Completed       bool                  `json:"completed"`
	FinalCid        string                `json:"final_cid"`
	FinalMode       int                   `json:"final_mode"`
	FinalCustomMode int                   `json:"final_custom_mode"`
	Subgroups       []rsssite.RssSubgroup `json:"subgroups"`
	TmdbID          string                `json:"tmdb_id"`
	TmdbType        string                `json:"tmdb_type"`
}

func (p feedPayload) validate() error {
	rawURL := strings.TrimSpace(p.URL)
	if rawURL == "" {
		return fmt.Errorf("RSS 地址不能为空")
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("RSS 地址无效: %s", rawURL)
	}
	return nil
}

func (p feedPayload) rssConfig() *rsssite.RssConfig {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		name = p.URL
	}
	return &rsssite.RssConfig{
		Name:         name,
		SubgroupName: strings.TrimSpace(p.SubgroupName),
		BgmURL:       strings.TrimSpace(p.BgmURL),
		Url:          strings.TrimSpace(p.URL),
		Source:       config.NormalizeSource(p.Source, p.URL),
		Enabled:      p.Enabled,
		Cid:          strings.TrimSpace(p.Cid),
		SavePath:     strings.TrimSpace(p.SavePath),
		Filter:       strings.TrimSpace(p.Filter),
		Exclude:      strings.TrimSpace(p.Exclude),
		Expiration:   p.Expiration,
		RenameSeason: p.RenameSeason,
		RenameOffset: p.RenameOffset,
		Subgroups:    p.Subgroups,
		TmdbID:       strings.TrimSpace(p.TmdbID),
		TmdbType:     strings.TrimSpace(p.TmdbType),
	}
}

type previewItem struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Magnet string `json:"magnet"`
}

type previewSource struct {
	Role     string        `json:"role"`
	Name     string        `json:"name"`
	Subgroup string        `json:"subgroup,omitempty"`
	URL      string        `json:"url"`
	Total    int           `json:"total"`
	Matched  int           `json:"matched"`
	Items    []previewItem `json:"items"`
	Error    string        `json:"error,omitempty"`
}

type previewResponse struct {
	Total     int             `json:"total"`
	Matched   int             `json:"matched"`
	Truncated bool            `json:"truncated"`
	Items     []previewItem   `json:"items"`
	Sources   []previewSource `json:"sources"`
}

func (s *Server) handlePreviewFeed(w http.ResponseWriter, r *http.Request) {
	var payload feedPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if err := payload.validate(); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	cfg := payload.rssConfig()
	rsssite.ApplyDefaults(s.configSnapshot(), cfg)

	sources := make([]previewSource, 0, 1+len(payload.Subgroups))
	mainSubgroup := strings.TrimSpace(cfg.SubgroupName)
	if mainSubgroup == "" {
		mainSubgroup = subgroupFromURL(cfg.Url)
	}
	mainSource := previewRSSSource(cfg, "main", cfg.Name, mainSubgroup, cfg.Url)
	sources = append(sources, mainSource)
	for _, subgroup := range payload.Subgroups {
		if !subgroup.Enabled || strings.TrimSpace(subgroup.URL) == "" {
			continue
		}
		copy := *cfg
		copy.Url = strings.TrimSpace(subgroup.URL)
		copy.Filter = subgroup.Filter
		copy.Exclude = subgroup.Exclude
		if subgroup.Season > 0 {
			copy.RenameSeason = subgroup.Season
		}
		copy.RenameOffset = subgroup.Offset
		name := strings.TrimSpace(subgroup.Name)
		if name == "" {
			name = copy.Url
		}
		sources = append(sources, previewRSSSource(&copy, "backup", cfg.Name, name, copy.Url))
	}

	resp := previewResponse{Sources: sources, Items: make([]previewItem, 0)}
	for _, source := range sources {
		resp.Total += source.Total
		resp.Matched += source.Matched
		if source.Role == "main" {
			resp.Items = source.Items
		}
	}
	for _, source := range sources {
		if len(source.Items) >= 30 {
			resp.Truncated = true
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

func previewRSSSource(cfg *rsssite.RssConfig, role, name, subgroup, rawURL string) previewSource {
	source := previewSource{Role: role, Name: name, Subgroup: subgroup, URL: rawURL, Items: make([]previewItem, 0)}
	total, items, err := rsssite.GetMagnetItemListWithStats(cfg)
	if err != nil {
		source.Error = err.Error()
		return source
	}
	source.Total = total
	source.Matched = len(items)
	for i, item := range items {
		if i >= 30 {
			break
		}
		source.Items = append(source.Items, previewItem{Title: item.Title, Link: item.Link, Magnet: item.Magnet})
	}
	return source
}

func (s *Server) handleRunFeed(w http.ResponseWriter, r *http.Request) {
	var payload feedPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if err := payload.validate(); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}

	cfg := payload.rssConfig()
	if current, ok := lookupConfiguredFeed(s.configSnapshot(), cfg.Url); ok {
		cfg = current
	}
	rsssite.ApplyDefaults(s.configSnapshot(), cfg)

	job, err := s.startRunFeedJob(ag, cfg.Name, func() []p115.FeedRunResult {
		return ag.RunRssConfigGroup(cfg)
	})
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}

func (s *Server) handleRunAllFeeds(w http.ResponseWriter, r *http.Request) {
	var ignored struct{}
	if r.ContentLength > 0 {
		_ = decodeJSON(w, r, &ignored)
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	job, err := s.startRunAllJob(ag, "all", "运行全部订阅")
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}

func lookupConfiguredFeed(cfg *config.Config, rawURL string) (*rsssite.RssConfig, bool) {
	if cfg == nil {
		return nil, false
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if feed.Url != rawURL {
				continue
			}
			enabled := true
			if feed.Enabled != nil {
				enabled = *feed.Enabled
			}
			return &rsssite.RssConfig{
				Name:       feed.Name,
				Url:        feed.Url,
				Source:     config.NormalizeSource(feed.Source, feed.Url),
				Enabled:    enabled,
				Cid:        feed.Cid,
				SavePath:   feed.SavePath,
				Filter:     feed.Filter,
				Exclude:    feed.Exclude,
				Expiration: feed.Expiration,
			}, true
		}
	}
	return nil, false
}
