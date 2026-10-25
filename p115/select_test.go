package p115

import (
	"testing"

	"github.com/zhifengle/anifeed/rsssite"
)

func TestSelectEpisodesPrefersBetterVersion(t *testing.T) {
	items := []rsssite.MagnetItem{
		{Title: "[Sub] Anime - 01 [720p][繁][内嵌][AVC][8bit].mp4"},
		{Title: "[Sub] Anime - 01 [1080p][简][内封][HEVC][10bit].mkv"},
		{Title: "[Sub] Anime - 02 [1080p][简][内封][HEVC][10bit].mkv"},
	}
	got := SelectEpisodes(items, true)
	if !got[1].Selected || got[0].Selected {
		t.Fatalf("episode 1 selection = %+v / %+v", got[0], got[1])
	}
	if !got[2].Selected {
		t.Fatalf("episode 2 should be selected: %+v", got[2])
	}
	if got[0].Reason == "" {
		t.Fatal("skipped version should have a reason")
	}
}

func TestSelectEpisodesKeepsUnparseableAndHalfEpisodes(t *testing.T) {
	items := []rsssite.MagnetItem{
		{Title: "[Sub] Anime - 01 [1080p]"},
		{Title: "[Sub] Anime - 01.5 [1080p]"},
		{Title: "[Sub] Anime 合集 [BDRip]"},
	}
	got := SelectEpisodes(items, true)
	for i, selection := range got {
		if !selection.Selected {
			t.Errorf("item %d should be selected: %+v", i, selection)
		}
	}
}

func TestSelectEpisodesDisabledSelectsAll(t *testing.T) {
	items := []rsssite.MagnetItem{{Title: "[Sub] Anime - 01 [720p]"}, {Title: "[Sub] Anime - 01 [1080p]"}}
	got := SelectEpisodes(items, false)
	if !got[0].Selected || !got[1].Selected {
		t.Fatalf("disabled selection = %+v", got)
	}
}

func TestSelectEpisodesUserReportedTitles(t *testing.T) {
	prefix := "【今晚月色真美】[青之箱 第二季 / アオのハコ Season 2 / Ao no Hako Season 2]"
	items := []rsssite.MagnetItem{
		{Title: prefix + "[03][WebRip][1080P][AVC-8bit][AAC][繁日內嵌]"},
		{Title: prefix + "[03][WebRip][1080P][AVC-8bit][AAC][简日内嵌]"},
		{Title: prefix + "[02][WebRip][1080P][AVC-8bit][AAC][繁日內嵌]"},
		{Title: prefix + "[02][WebRip][1080P][AVC-8bit][AAC][简日内嵌]"},
		{Title: prefix + "[01][WebRip][1080P][AVC-8bit][AAC][繁日內嵌]"},
		{Title: prefix + "[01][WebRip][1080P][AVC-8bit][AAC][简日内嵌]"},
	}
	got := SelectEpisodes(items, true)
	for i, selection := range got {
		want := i%2 == 1
		if selection.Selected != want {
			t.Errorf("item %d selected = %v, want %v (%+v)", i, selection.Selected, want, selection)
		}
	}
}

func TestSelectEpisodesPrefersSimplifiedOverBilingual(t *testing.T) {
	items := []rsssite.MagnetItem{
		{Title: "[Sub] Anime - 01 [1080p][简繁][内封][HEVC].mkv"},
		{Title: "[Sub] Anime - 01 [1080p][简体][内封][HEVC].mkv"},
	}
	got := SelectEpisodes(items, true)
	if !got[1].Selected || got[0].Selected {
		t.Fatalf("simplified should win: %+v / %+v", got[0], got[1])
	}
}

func TestSelectWithPreviewDecisions(t *testing.T) {
	items := []rsssite.MagnetItem{
		{Title: "[Sub] Anime - 01 [720p]", Magnet: "magnet:?xt=urn:btih:a"},
		{Title: "[Sub] Anime - 01 [1080p]", Magnet: "magnet:?xt=urn:btih:b"},
		{Title: "[Sub] Anime - 02 [1080p]", Magnet: "magnet:?xt=urn:btih:c"},
	}
	auto := SelectEpisodes(items, true)
	decisions := map[string]bool{
		ItemKey(items[0]): true,
		ItemKey(items[1]): false,
	}
	selected, skipped := selectWithDecisions(items, auto, decisions)
	if len(selected) != 2 || skipped != 1 {
		t.Fatalf("selected=%d skipped=%d", len(selected), skipped)
	}
	if selected[0].Magnet != items[0].Magnet || selected[1].Magnet != items[2].Magnet {
		t.Fatalf("unexpected selection: %+v", selected)
	}
}

func TestSelectEpisodesPrefersHigherVersion(t *testing.T) {
	items := []rsssite.MagnetItem{
		{Title: "[Sub] Anime - 01v1 [1080p]"},
		{Title: "[Sub] Anime - 01v3 [1080p]"},
		{Title: "[Sub] Anime - 01v2 [1080p]"},
	}
	got := SelectEpisodes(items, true)
	if !got[1].Selected || got[0].Selected || got[2].Selected {
		t.Fatalf("v3 should win: %+v", got)
	}
}
