package p115

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/huohen92/anifeed/rsssite"
)

// ItemKey returns a stable key for preview decisions.
func ItemKey(item rsssite.MagnetItem) string {
	value := item.Magnet
	if value == "" {
		value = item.Title + "|" + item.Link
	}
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

// ExtractSubgroup returns the bracket subgroup at the start of a title.
func ExtractSubgroup(title string) string {
	return parseSubgroup(title)
}

// selectWithDecisions applies saved per-item preview decisions. Denied items
// are skipped; explicitly allowed items bypass the automatic episode choice.
func selectWithDecisions(items []rsssite.MagnetItem, auto []Selection, decisions map[string]bool) ([]rsssite.MagnetItem, int) {
	selected := make([]rsssite.MagnetItem, 0, len(items))
	skipped := 0
	for i, item := range items {
		if allowed, ok := decisions[ItemKey(item)]; ok {
			if !allowed {
				skipped++
				continue
			}
			selected = append(selected, item)
			continue
		}
		if !auto[i].Selected {
			skipped++
			continue
		}
		selected = append(selected, item)
	}
	return selected, skipped
}
