package rsssite

import (
	"fmt"
	"log"
	urlPkg "net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/huohen92/anifeed/request"
	"github.com/huohen92/anifeed/ruleengine"
	"github.com/mmcdole/gofeed"
)

var x5EpisodePattern = regexp.MustCompile(`(?:^|[^\d])\d{1,4}\.5(?:[^\d]|$)`)

type MagnetSite interface {
	GetMagnet(item *gofeed.Item) string
	GetMagnetItem(item *gofeed.Item) MagnetItem
}

type MagnetItem struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Magnet      string    `json:"magnet"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Published   time.Time `json:"published,omitempty"`
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

var feedCache sync.Map // url -> feedCacheEntry
type feedCacheEntry struct {
	feed *gofeed.Feed
	at   time.Time
}

// FetchFeed downloads and parses a feed, returning the underlying error.
// A short cache prevents preview, rule tests, and a run from downloading the
// same RSS several times in one burst.
func FetchFeed(url string) (*gofeed.Feed, error) {
	if cached, ok := feedCache.Load(url); ok {
		entry := cached.(feedCacheEntry)
		if time.Since(entry.at) < 15*time.Second {
			return entry.feed, nil
		}
	}
	feed, err := fetchFeedUncached(url)
	if err != nil {
		return nil, err
	}
	feedCache.Store(url, feedCacheEntry{feed: feed, at: time.Now()})
	return feed, nil
}

func fetchFeedUncached(url string) (*gofeed.Feed, error) {
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

func subgroupNameForFeed(configured string, feed *gofeed.Feed) string {
	configured = strings.TrimSpace(configured)
	if configured != "" && !isASCII(configured) {
		return configured
	}
	for i, item := range feed.Items {
		if i >= 30 {
			break
		}
		match := subgroupTitlePattern.FindStringSubmatch(item.Title)
		if len(match) != 2 {
			continue
		}
		name := strings.TrimSpace(match[1])
		if name != "" && !isASCII(name) {
			return name
		}
	}
	return configured
}

func isASCII(value string) bool {
	for _, r := range value {
		if r > 127 {
			return false
		}
	}
	return true
}

var subgroupTitlePattern = regexp.MustCompile(`^\s*[\[【]([^\]】]{1,60})[\]】]`)

func magnetItemsFromFeedDetailed(feed *gofeed.Feed, site MagnetSite, config *RssConfig) ([]PreviewMagnetItem, error) {
	matchRules, err := ruleengine.Parse(config.Filter)
	if err != nil {
		return nil, fmt.Errorf("匹配规则格式错误: %w", err)
	}
	excludeRules, err := ruleengine.Parse(config.Exclude)
	if err != nil {
		return nil, fmt.Errorf("排除规则格式错误: %w", err)
	}
	globalRules, err := ruleengine.Parse(strings.Join(config.GlobalExcludes, "\n"))
	if err != nil {
		return nil, fmt.Errorf("全局排除规则格式错误: %w", err)
	}
	subgroup := subgroupNameForFeed(config.SubgroupName, feed)
	itemList := make([]PreviewMagnetItem, 0, len(feed.Items))
	for _, item := range feed.Items {
		entry := PreviewMagnetItem{Item: site.GetMagnetItem(item)}
		if config.SkipX5 && x5EpisodePattern.MatchString(item.Title) {
			entry.Filtered = true
			entry.Reason = "x5 集已跳过"
			itemList = append(itemList, entry)
			continue
		}
		decision := ruleengine.Evaluate(item.Title, subgroup, matchRules, excludeRules, globalRules)
		if !decision.Matched {
			entry.Filtered = true
			entry.Reason = decision.Reason
			if decision.Rule != "" {
				entry.Reason += "：" + decision.Rule
			}
		}
		itemList = append(itemList, entry)
	}
	return itemList, nil
}

type cachedFeed struct {
	feed *gofeed.Feed
	at   time.Time
}

var previewFeedCache = struct {
	sync.Mutex
	entries map[string]cachedFeed
}{entries: map[string]cachedFeed{}}

func fetchFeedForPreview(rawURL string, ttl time.Duration) (*gofeed.Feed, error) {
	now := time.Now()
	previewFeedCache.Lock()
	entry, ok := previewFeedCache.entries[rawURL]
	previewFeedCache.Unlock()
	if ok && now.Sub(entry.at) < ttl {
		return entry.feed, nil
	}
	feed, err := FetchFeed(rawURL)
	if err != nil {
		return nil, err
	}
	previewFeedCache.Lock()
	for key, value := range previewFeedCache.entries {
		if now.Sub(value.at) > 30*time.Minute {
			delete(previewFeedCache.entries, key)
		}
	}
	previewFeedCache.entries[rawURL] = cachedFeed{feed: feed, at: now}
	previewFeedCache.Unlock()
	return feed, nil
}

// GetMagnetItemsPreviewCached is used by the interactive rules preview. It
// reuses a parsed RSS feed for a few minutes so editing rules does not hammer
// the source. Normal scheduled/manual runs still call the non-cached path.
func GetMagnetItemsPreviewCached(config *RssConfig, ttl time.Duration) (int, []PreviewMagnetItem, error) {
	site := getSiteForConfig(config)
	feed, err := fetchFeedForPreview(config.Url, ttl)
	if err != nil {
		return 0, nil, err
	}
	items, err := magnetItemsFromFeedDetailed(feed, site, config)
	if err != nil {
		return len(feed.Items), nil, err
	}
	return len(feed.Items), items, nil
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

func itemPublished(item *gofeed.Item) time.Time {
	if item == nil {
		return time.Time{}
	}
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, strings.TrimSpace(item.Published)); err == nil {
			return t
		}
	}
	return time.Time{}
}
