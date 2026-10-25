package rsssite

import (
	urlPkg "net/url"
	"sort"
	"strings"

	"github.com/zhifengle/anifeed/config"
)

var (
	RssConfigDict map[string][]RssConfig
	rssJsonPath   string
)

type RssSubgroup struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Season  int    `json:"season,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Filter  string `json:"filter,omitempty"`
	Exclude string `json:"exclude,omitempty"`
	Enabled bool   `json:"enabled"`
}

type RssConfig struct {
	Name           string        `json:"name"`
	SubgroupName   string        `json:"subgroup_name,omitempty"`
	BgmURL         string        `json:"bgm_url,omitempty"`
	Url            string        `json:"url"`
	Source         string        `json:"source,omitempty"`
	Enabled        bool          `json:"enabled"`
	Completed      bool          `json:"completed,omitempty"`
	Cid            string        `json:"cid,omitempty"`
	SavePath       string        `json:"savepath,omitempty"`
	Filter         string        `json:"filter,omitempty"`
	Exclude        string        `json:"exclude,omitempty"`
	Expiration     uint          `json:"expiration,omitempty"`
	GlobalExcludes []string      `json:"-"`
	SkipX5         bool          `json:"-"`
	RenameSeason   int           `json:"rename_season,omitempty"`
	RenameOffset   int           `json:"rename_offset,omitempty"`
	TmdbID         string        `json:"tmdb_id,omitempty"`
	TmdbType       string        `json:"tmdb_type,omitempty"`
	Subgroups      []RssSubgroup `json:"subgroups,omitempty"`
}

func SetRssJsonPath(p string) {
	rssJsonPath = p
}

func ReadRssConfigDict() *map[string][]RssConfig {
	// Use config.Load() with RSSPath from SetRssJsonPath()
	cfg, _, err := config.LoadWithOptions(config.CLIParams{RSSPath: rssJsonPath}, config.LoadOptions{RSS: true})
	if err != nil {
		return nil
	}

	// Convert config.RssConfig to rsssite.RssConfig
	result := make(map[string][]RssConfig)
	for site, configs := range cfg.RSS {
		siteConfigs := make([]RssConfig, len(configs))
		for i, c := range configs {
			siteConfigs[i] = RssConfig{
				Name:         c.Name,
				SubgroupName: c.SubgroupName,
				BgmURL:       c.BgmURL,
				Url:          c.Url,
				Source:       c.Source,
				Enabled:      c.IsEnabled(),
				Completed:    c.Completed,
				Cid:          c.Cid,
				SavePath:     c.SavePath,
				Filter:       c.Filter,
				Exclude:      c.Exclude,
				Expiration:   c.Expiration,
				RenameSeason: c.RenameSeason,
				RenameOffset: c.RenameOffset,
				TmdbID:       c.TmdbID,
				TmdbType:     c.TmdbType,
				Subgroups:    transformSubgroups(c.Subgroups),
			}
			ApplyDefaults(cfg, &siteConfigs[i])
		}
		result[site] = siteConfigs
	}

	RssConfigDict = result
	return &result
}

func GetRssConfigByURL(url string) *RssConfig {
	urlObj, err := urlPkg.Parse(url)
	if err != nil {
		return nil
	}
	ReadRssConfigDict()
	configs, ok := RssConfigDict[urlObj.Host]
	if !ok {
		return fallbackRssConfig(url)
	}
	for _, config := range configs {
		if isSameRSSURL(config.Url, url) {
			return &config
		}
	}
	return fallbackRssConfig(url)
}

func fallbackRssConfig(rawURL string) *RssConfig {
	fallback := &RssConfig{
		Url:     rawURL,
		Source:  config.DetectSource(rawURL),
		Enabled: true,
	}
	if cfg, _, err := config.LoadWithOptions(config.CLIParams{}, config.LoadOptions{RSS: true}); err == nil {
		ApplyDefaults(cfg, fallback)
	}
	return fallback
}

// ApplyDefaults fills the global offline defaults on a feed config.
func ApplyDefaults(cfg *config.Config, rss *RssConfig) {
	if cfg == nil || rss == nil {
		return
	}
	if rss.Source == "" {
		rss.Source = config.DetectSource(rss.Url)
	} else {
		rss.Source = config.NormalizeSource(rss.Source, rss.Url)
	}
	if rss.Cid == "" {
		rss.Cid = cfg.Offline.DefaultCid
	}
	rss.GlobalExcludes = append([]string(nil), cfg.Offline.GlobalExclude...)
	rss.SkipX5 = cfg.Offline.SkipX5
}

func urlHost(raw string) string {
	u, err := urlPkg.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func isSameRSSURL(left string, right string) bool {
	if left == right {
		return true
	}
	return normalizeRSSURL(left) == normalizeRSSURL(right)
}

func normalizeRSSURL(raw string) string {
	u, err := urlPkg.Parse(raw)
	if err != nil {
		return raw
	}

	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if port != "" {
		host = host + ":" + port
	}

	path := strings.TrimSuffix(u.EscapedPath(), "/")
	query := u.Query()
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	queryParts := make([]string, 0)
	for _, key := range keys {
		values := query[key]
		sort.Strings(values)
		escapedKey := urlPkg.QueryEscape(key)
		for _, value := range values {
			queryParts = append(queryParts, escapedKey+"="+urlPkg.QueryEscape(value))
		}
	}

	normalized := host + path
	if len(queryParts) > 0 {
		normalized += "?" + strings.Join(queryParts, "&")
	}
	return normalized
}

func transformSubgroups(values []config.RssSubgroup) []RssSubgroup {
	result := make([]RssSubgroup, 0, len(values))
	for _, value := range values {
		result = append(result, RssSubgroup{Name: value.Name, URL: value.URL, Season: value.Season, Offset: value.Offset, Filter: value.Filter, Exclude: value.Exclude, Enabled: value.IsEnabled()})
	}
	return result
}
