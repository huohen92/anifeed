package webui

import (
	"testing"

	"github.com/huohen92/anifeed/config"
)

func TestSourceSubgroupNamePrefersBackupURL(t *testing.T) {
	cfg := &config.Config{RSS: map[string][]config.RssConfig{
		"anibt": {{
			Name: "青之箱",
			Url:  "https://anibt.net/rss/anime.xml?bgmId=1&groupSlug=main",
			Subgroups: []config.RssSubgroup{{
				Name: "ANi",
				URL:  "https://anibt.net/rss/anime.xml?bgmId=1&groupSlug=ani",
			}},
		}},
	}}
	got := sourceSubgroupName(cfg, "青之箱", "https://anibt.net/rss/anime.xml?bgmId=1&groupSlug=ani", "[ANi] 青之箱 - 01")
	if got != "ANi" {
		t.Fatalf("sourceSubgroupName = %q, want ANi", got)
	}
}

func TestSourceFeedNameUsesBackupParent(t *testing.T) {
	cfg := &config.Config{RSS: map[string][]config.RssConfig{
		"anibt": {{
			Name: "侦探已经死了。",
			Url:  "https://anibt.net/rss/anime.xml?bgmId=393010&groupSlug=lolihouse",
			Subgroups: []config.RssSubgroup{{
				Name: "ANi",
				URL:  "https://anibt.net/rss/anime.xml?bgmId=393010&groupSlug=ani",
			}},
		}},
	}}
	got := sourceFeedName(cfg, "ANi", "https://anibt.net/rss/anime.xml?bgmId=393010&groupSlug=ani")
	if got != "侦探已经死了。" {
		t.Fatalf("sourceFeedName = %q, want parent title", got)
	}
}
