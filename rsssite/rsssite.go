package rsssite

import (
	"fmt"
	"log"
	urlPkg "net/url"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/mmcdole/gofeed"
	"github.com/zhifengle/anifeed/request"
)

var x5EpisodePattern = regexp.MustCompile(`(?:^|[^\d])\d{1,4}\.5(?:[^\d]|$)`)

type MagnetSite interface {
	GetMagnet(item *gofeed.Item) string
	GetMagnetItem(item *gofeed.Item) MagnetItem
}

type MagnetItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Magnet      string `json:"magnet"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func getSiteForConfig(config *RssConfig) MagnetSite {
	switch config.Source {
	case "mikan":
		return &Mikanani{}
	case "anibt":
		return &Anibt{}
	case "animegarden":
		return &Generic{}
	case "other":
		// Keep the legacy site-specific handlers for Nyaa/Dmhy/ACGNX/RSSHub;
		// unknown hosts fall back to Generic inside getSite.
		return getSite(config.Url)
	default:
		return getSite(config.Url)
	}
}

func getSite(url string) MagnetSite {
	name := url
	if strings.HasPrefix(url, "http") {
		urlObj, _ := urlPkg.Parse(url)
		name = urlObj.Host
	}
	switch name {
	case "mikanani.me", "mikanime.tv":
		return &Mikanani{}
	case "nyaa.si", "sukebei.nyaa.si":
		return &Nyaa{}
	case "share.dmhy.org":
		return &Dmhy{}
	case "share.acgnx.se", "share.acgnx.net", "www.acgnx.se":
		return &Acgnx{}
	case "rsshub.app":
		return &Rsshub{}
	case "anibt.net":
		return &Anibt{}
	case "animes.garden", "api.animes.garden":
		return &Generic{}
	default:
		return &Generic{}
	}
}

func GetFeed(url string) *gofeed.Feed {
	feed, err := FetchFeed(url)
	if err != nil {
		log.Printf("[error] %s\n", err)
		return nil
	}
	return feed
}

// FetchFeed downloads and parses a feed, returning the underlying error.
func FetchFeed(url string) (*gofeed.Feed, error) {
	res, err := request.Get(url, nil)
	if err != nil {
		return nil, fmt.Errorf("get rss from %s: %w", url, err)
	}
	feed, err := gofeed.NewParser().ParseString(res)
	if err != nil {
		return nil, fmt.Errorf("parse rss from %s: %w", url, err)
	}
	return feed, nil
}

func GetMagnetItemList(config *RssConfig) []MagnetItem {
	_, items, err := GetMagnetItemListWithStats(config)
	if err != nil {
		log.Printf("[error] %s\n", err)
		return nil
	}
	return items
}

// GetMagnetItemListWithStats fetches a feed and returns the total item count
// together with the items matching the configured filter/exclude rules.
func GetMagnetItemListWithStats(config *RssConfig) (total int, items []MagnetItem, err error) {
	site := getSiteForConfig(config)
	feed, err := FetchFeed(config.Url)
	if err != nil {
		return 0, nil, err
	}
	items, err = magnetItemsFromFeedStrict(feed, site, config)
	if err != nil {
		return len(feed.Items), nil, err
	}
	return len(feed.Items), items, nil
}

func magnetItemsFromFeed(feed *gofeed.Feed, site MagnetSite, config *RssConfig) []MagnetItem {
	items, _ := magnetItemsFromFeedStrict(feed, site, config)
	return items
}

// PreviewMagnetItem carries one RSS item and whether it survived the rules.
type PreviewMagnetItem struct {
	Item     MagnetItem
	Filtered bool
	Reason   string
}

func magnetItemsFromFeedStrict(feed *gofeed.Feed, site MagnetSite, config *RssConfig) ([]MagnetItem, error) {
	detailed, err := magnetItemsFromFeedDetailed(feed, site, config)
	if err != nil {
		return nil, err
	}
	items := make([]MagnetItem, 0, len(detailed))
	for _, entry := range detailed {
		if !entry.Filtered {
			items = append(items, entry.Item)
		}
	}
	return items, nil
}

func magnetItemsFromFeedDetailed(feed *gofeed.Feed, site MagnetSite, config *RssConfig) ([]PreviewMagnetItem, error) {
	itemList := make([]PreviewMagnetItem, 0, len(feed.Items))
	whiteRules, err := compileFilterRules(config.Filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter %q: %w", config.Filter, err)
	}
	blackRules, err := compileFilterRules(config.Exclude)
	if err != nil {
		return nil, fmt.Errorf("invalid exclude %q: %w", config.Exclude, err)
	}

	type globalRule struct {
		rule string
		re   *regexp2.Regexp
	}
	globalRules := make([]globalRule, 0, len(config.GlobalExcludes))
	for _, rule := range config.GlobalExcludes {
		if strings.TrimSpace(rule) == "" {
			continue
		}
		re, err := compileFilter(rule)
		if err != nil {
			return nil, fmt.Errorf("invalid global exclude %q: %w", rule, err)
		}
		globalRules = append(globalRules, globalRule{rule: rule, re: re})
	}

	for _, item := range feed.Items {
		entry := PreviewMagnetItem{Item: site.GetMagnetItem(item)}
		if config.SkipX5 && x5EpisodePattern.MatchString(item.Title) {
			entry.Filtered = true
			entry.Reason = "x5 集已跳过"
			itemList = append(itemList, entry)
			continue
		}
		globallyExcluded := false
		for _, rule := range globalRules {
			if matchRule(item.Title, rule.rule, rule.re) {
				globallyExcluded = true
				break
			}
		}
		if globallyExcluded {
			entry.Filtered = true
			entry.Reason = "命中全局排除"
			itemList = append(itemList, entry)
			continue
		}
		if !matchFilterRules(item.Title, whiteRules, true) {
			entry.Filtered = true
			entry.Reason = "未通过 filter"
			itemList = append(itemList, entry)
			continue
		}
		if matchFilterRules(item.Title, blackRules, false) {
			entry.Filtered = true
			entry.Reason = "命中 exclude"
			itemList = append(itemList, entry)
			continue
		}
		itemList = append(itemList, entry)
	}
	return itemList, nil
}

// GetMagnetItemsPreview returns every RSS item, including items filtered out by
// rules, so the UI can display both matched and ignored entries.
func GetMagnetItemsPreview(config *RssConfig) (total int, items []PreviewMagnetItem, err error) {
	site := getSiteForConfig(config)
	feed, err := FetchFeed(config.Url)
	if err != nil {
		return 0, nil, err
	}
	items, err = magnetItemsFromFeedDetailed(feed, site, config)
	if err != nil {
		return len(feed.Items), nil, err
	}
	return len(feed.Items), items, nil
}

type filterRule struct {
	group   string
	literal string
	re      *regexp2.Regexp
}

// compileFilterRules parses newline-separated rules. A line may be a plain
// substring, /regex/, or ani-rss style {{subgroup}}:regex.
func compileFilterRules(raw string) ([]filterRule, error) {
	var rules []filterRule
	for lineNo, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "{{") {
			end := strings.Index(line, "}}")
			if end < 0 || end+2 >= len(line) || line[end+2] != ':' {
				return nil, fmt.Errorf("line %d: expected {{subgroup}}:regex", lineNo+1)
			}
			group := strings.TrimSpace(line[2:end])
			pattern := strings.TrimSpace(line[end+3:])
			if group == "" || pattern == "" {
				return nil, fmt.Errorf("line %d: empty subgroup or pattern", lineNo+1)
			}
			re, err := regexp2.Compile(pattern, 0)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo+1, err)
			}
			rules = append(rules, filterRule{group: group, re: re})
			continue
		}
		re, err := compileFilter(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo+1, err)
		}
		rules = append(rules, filterRule{literal: line, re: re})
	}
	return rules, nil
}

func matchFilterRules(title string, rules []filterRule, all bool) bool {
	if len(rules) == 0 {
		return all
	}
	for _, rule := range rules {
		matched := false
		if rule.group != "" {
			if strings.Contains(title, "["+rule.group+"]") || strings.Contains(title, "【"+rule.group+"】") {
				matched = rule.re != nil && matchRegex(rule.re, title)
			}
		} else if rule.re != nil {
			matched = matchRegex(rule.re, title)
		} else {
			matched = strings.Contains(title, rule.literal)
		}
		if all && !matched {
			return false
		}
		if !all && matched {
			return true
		}
	}
	return all
}

func matchRegex(re *regexp2.Regexp, title string) bool {
	matched, _ := re.MatchString(title)
	return matched
}

// compileFilter compiles a /regex/ rule into a regexp2 pattern.
// A plain substring rule returns nil and is matched via strings.Contains.
func compileFilter(rule string) (*regexp2.Regexp, error) {
	if strings.HasPrefix(rule, "/") && strings.HasSuffix(rule, "/") && len(rule) > 2 {
		return regexp2.Compile(rule[1:len(rule)-1], 0)
	}
	return nil, nil
}

// matchRule reports whether title matches the filter rule.
// An empty rule matches everything.
func matchRule(title, rule string, re *regexp2.Regexp) bool {
	if rule == "" {
		return true
	}
	if re != nil {
		m, _ := re.MatchString(title)
		return m
	}
	return strings.Contains(title, rule)
}

func GetMagnetByEnclosure(item *gofeed.Item) string {
	if len(item.Enclosures) == 0 {
		return ""
	}
	// find enclosure by type == "application/x-bittorrent" or url has prefix magnet:?
	for _, enclosure := range item.Enclosures {
		if enclosure.Type == "application/x-bittorrent" || strings.HasPrefix(enclosure.URL, "magnet:?") {
			lst := strings.Split(enclosure.URL, "&dn=")
			if len(lst) != 2 {
				return enclosure.URL
			}
			return lst[0]
		}
	}
	return ""
}
