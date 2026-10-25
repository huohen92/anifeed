package webui

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
)

// sourceSubgroupName resolves the subtitle-group label configured on the RSS
// source itself, falling back to the groupSlug in the URL and finally the
// first bracket in the torrent title.
func sourceSubgroupName(cfg *config.Config, feedName, feedURL, itemTitle string) string {
	if cfg != nil {
		for _, feeds := range cfg.RSS {
			for _, feed := range feeds {
				if feedName != "" && feed.Name != feedName && feed.Url != feedURL {
					continue
				}
				if feedURL == "" || feed.Url == feedURL {
					if name := cleanSubgroupLabel(feed.SubgroupName, feed.Name); name != "" {
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
	return p115.ExtractSubgroup(itemTitle)
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
