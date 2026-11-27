package webui

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/rsssite"
)

// sourceSubgroupName resolves the subtitle-group label configured on the RSS
// source itself, falling back to the groupSlug in the URL and finally the
// first bracket in the torrent title.
func sourceSubgroupName(cfg *config.Config, feedName, feedURL, itemTitle string) string {
	if cfg != nil {
		// A backup source URL identifies the exact subtitle group. Prefer it
		// before matching by feed name, because all backup sources share the
		// same feed name.
		if feedURL != "" {
			for _, feeds := range cfg.RSS {
				for _, feed := range feeds {
					for _, group := range feed.Subgroups {
						if group.URL != feedURL {
							continue
						}
						name := cleanSubgroupLabel(group.Name, feed.Name)
						if name != "" && name != "未知字幕组" {
							return name
						}
						if name := subgroupFromURL(feedURL); name != "" {
							return name
						}
					}
				}
			}
		}
		for _, feeds := range cfg.RSS {
			for _, feed := range feeds {
				if feedName != "" && feed.Name != feedName && feed.Url != feedURL {
					continue
				}
				if feedURL == "" || feed.Url == feedURL {
					if name := cleanSubgroupLabel(feed.SubgroupName, feed.Name); name != "" {
						if subgroupASCII(name) {
							if extracted := p115.ExtractSubgroup(itemTitle); extracted != "" && extracted != "未知字幕组" && !subgroupASCII(extracted) {
								return extracted
							}
						}
						return name
					}
					if name := subgroupFromURL(feed.Url); name != "" {
						return name
					}
				}
				for _, group := range feed.Subgroups {
					if feedURL == "" || group.URL != feedURL {
						continue
					}
					if name := cleanSubgroupLabel(group.Name, feed.Name); name != "" {
						return name
					}
				}
			}
		}
	}
	extracted := p115.ExtractSubgroup(itemTitle)
	if extracted != "" && extracted != "未知字幕组" {
		return extracted
	}
	if name := subgroupFromURL(feedURL); name != "" {
		return name
	}
	return "未知字幕组"
}

func subgroupASCII(value string) bool {
	for _, r := range value {
		if r > 127 {
			return false
		}
	}
	return true
}

func cleanSubgroupLabel(value, feedName string) string {
	value = strings.TrimSpace(value)
	if feedName != "" && strings.HasPrefix(value, feedName+" · ") {
		value = strings.TrimSpace(strings.TrimPrefix(value, feedName+" · "))
	}
	if index := strings.LastIndex(value, " · "); index >= 0 {
		value = strings.TrimSpace(value[index+3:])
	}
	return value
}

func subgroupFromURL(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ""
	}
	slug := strings.TrimSpace(u.Query().Get("groupSlug"))
	if slug == "" {
		return ""
	}
	compound := regexpSubgroupSuffix.FindStringSubmatch(slug)
	if len(compound) == 3 {
		return titleWord(compound[1]) + titleWord(compound[2])
	}
	return titleWord(slug)
}

var regexpSubgroupSuffix = regexp.MustCompile(`^(.+?)(house|hub|sub|team|raws|anime)$`)

func titleWord(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(strings.ToLower(value))
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

// handleSubgroupResolve reads the first release titles and resolves the
// subtitle group as it is written by the RSS source. This lets a URL slug such
// as tsukigakirei align with the Chinese name used by rules and templates.
func (s *Server) handleSubgroupResolve(w http.ResponseWriter, r *http.Request) {
	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	if rawURL == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少 RSS 地址")
		return
	}
	if alias := s.subgroupAliasFromURL(rawURL); alias != "" {
		writeJSON(w, http.StatusOK, map[string]string{"subgroup": alias})
		return
	}
	feed, err := rsssite.FetchFeed(rawURL)
	if err != nil {
		if alias := s.subgroupAliasFromURL(rawURL); alias != "" {
			writeJSON(w, http.StatusOK, map[string]string{"subgroup": alias})
			return
		}
		writeAPIError(w, http.StatusBadGateway, "读取 RSS 失败: "+err.Error())
		return
	}
	for i, item := range feed.Items {
		if i >= 30 {
			break
		}
		if name := p115.ExtractSubgroup(item.Title); name != "" && name != "未知字幕组" {
			writeJSON(w, http.StatusOK, map[string]string{"subgroup": s.canonicalSubgroup(name)})
			return
		}
	}
	if alias := s.subgroupAliasFromURL(rawURL); alias != "" {
		writeJSON(w, http.StatusOK, map[string]string{"subgroup": alias})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"subgroup": ""})
}

// sourceFeedName finds the parent subscription that owns a main or backup RSS
// URL. Backup rows can share a generated label; the parent feed name is the
// stable anime title used by renaming.
func sourceFeedName(cfg *config.Config, fallback, feedURL string) string {
	if cfg != nil && feedURL != "" {
		for _, feeds := range cfg.RSS {
			for _, feed := range feeds {
				if feed.Url == feedURL {
					return feed.Name
				}
				for _, group := range feed.Subgroups {
					if group.URL == feedURL {
						return feed.Name
					}
				}
			}
		}
	}
	return fallback
}
