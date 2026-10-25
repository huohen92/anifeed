package rsssite

import (
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

// Generic handles arbitrary RSS feeds. It looks for magnet/ed2k links in
// enclosures, descriptions and item GUIDs, plus direct HTTP(S) enclosures.
type Generic struct{}

var (
	genericMagnetPattern = regexp.MustCompile(`magnet:\?[^\s"'<>]+`)
	genericEd2kPattern   = regexp.MustCompile(`ed2k://[^\s"'<>]+`)
	genericInfoHash      = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
)

func (g *Generic) GetMagnet(item *gofeed.Item) string {
	for _, enclosure := range item.Enclosures {
		link := strings.TrimSpace(enclosure.URL)
		lower := strings.ToLower(link)
		switch {
		case strings.HasPrefix(lower, "magnet:?"):
			return trimMagnetDN(link)
		case strings.HasPrefix(lower, "ed2k://"):
			return link
		case strings.HasPrefix(lower, "http://"), strings.HasPrefix(lower, "https://"), strings.HasPrefix(lower, "ftp://"):
			path := lower
			if idx := strings.IndexAny(path, "?#"); idx >= 0 {
				path = path[:idx]
			}
			if strings.HasSuffix(path, ".torrent") {
				continue
			}
			return link
		}
	}

	for _, text := range []string{item.Link, item.Description, item.Content} {
		if magnet := genericMagnetPattern.FindString(text); magnet != "" {
			return trimMagnetDN(magnet)
		}
		if ed2k := genericEd2kPattern.FindString(text); ed2k != "" {
			return ed2k
		}
	}

	guid := strings.TrimSpace(item.GUID)
	if genericInfoHash.MatchString(guid) {
		return "magnet:?xt=urn:btih:" + strings.ToLower(guid)
	}
	return ""
}

func (g *Generic) GetMagnetItem(item *gofeed.Item) MagnetItem {
	return MagnetItem{
		Title:       item.Title,
		Link:        item.Link,
		Magnet:      g.GetMagnet(item),
		Description: item.Description,
		Content:     item.Content,
	}
}
