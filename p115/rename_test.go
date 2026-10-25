package p115

import "testing"

func TestBuildRenameName(t *testing.T) {
	opts := RenameOptions{Enabled: true, Template: "[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}", Season: 1, TmdbID: "486456"}
	cases := []struct {
		feed, item, want string
	}{
		{"佐佐木与文鸟小哔 第二季", "[LoliHouse] Sasaki to P-chan S2 - 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv", "[LoliHouse] 佐佐木与文鸟小哔 第二季 S01E01"},
		{"佐佐木与文鸟小哔 第二季 · LoliHouse", "[LoliHouse] Sasaki to P-chan S2 - 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv", "[LoliHouse] 佐佐木与文鸟小哔 第二季 S01E01"},
		{"测试番剧", "[Nix-Raws] Anime S01E12 [CR WEB-DL 1080p AVC AAC]", "[Nix-Raws] 测试番剧 S01E12"},
		{"测试番剧", "[ANi] Anime - 11 [1080P][Baha][WEB-DL]", "[ANi] 测试番剧 S01E11"},
		{"测试番剧", "The.Show.S02E02.1080p.WEB-DL.mkv", "[未知字幕组] 测试番剧 S01E02"},
		{"测试番剧", "[Sub] Anime - 10 [1080p]", "[Sub] 测试番剧 S01E10"},
		{"测试番剧", "[Sub] Anime - 11 [1080p]", "[Sub] 测试番剧 S01E11"},
		{"测试番剧", "[Sub] Anime - 01.5 [1080p]", "[Sub] 测试番剧 S01E01.5"},
		{"测试番剧", "[Sub] Anime 第05话 [1080p]", "[Sub] 测试番剧 S01E05"},
	}
	for _, c := range cases {
		got, err := BuildRenameName(c.feed, c.item, opts)
		if err != nil {
			t.Fatalf("%q: %v", c.item, err)
		}
		if got != c.want {
			t.Errorf("%q => %q, want %q", c.item, got, c.want)
		}
	}
}

func TestBuildRenameNameOffsetAndNoEpisode(t *testing.T) {
	opts := RenameOptions{Enabled: true, Template: "${title} - ${episode}", Season: 2, Offset: 1}
	got, err := BuildRenameName("番剧", "[Sub] Anime - 10 [1080p]", opts)
	if err != nil {
		t.Fatal(err)
	}
	if got != "番剧 - 11" {
		t.Errorf("got %q", got)
	}
	if _, err := BuildRenameName("番剧", "[Sub] Anime [1080p]", opts); err == nil {
		t.Fatal("expected episode parse error")
	}
}

func TestPreserveFileExtension(t *testing.T) {
	if got := PreserveFileExtension("剧集 S01E01", "[Sub] 剧集 - 01.mkv", false); got != "剧集 S01E01.mkv" {
		t.Errorf("got %q", got)
	}
	if got := PreserveFileExtension("剧集 S01E01.MKV", "old.mkv", false); got != "剧集 S01E01.MKV" {
		t.Errorf("got %q", got)
	}
	if got := PreserveFileExtension("剧集 S01E01", "old folder", true); got != "剧集 S01E01" {
		t.Errorf("got %q", got)
	}
}

func TestBuildRenameNameTmdbID(t *testing.T) {
	opts := RenameOptions{Enabled: true, Template: "${title} [tmdbid=${tmdbid}] S${seasonFormat}E${episodeFormat}", Season: 2, TmdbID: "486456"}
	got, err := BuildRenameName("测试番剧", "[Sub] Anime - 01 [1080p]", opts)
	if err != nil {
		t.Fatal(err)
	}
	if got != "测试番剧 [tmdbid=486456] S02E01" {
		t.Fatalf("got %q", got)
	}
}

func TestExtractEpisodeIgnoresSeasonNumber(t *testing.T) {
	cases := map[string]float64{
		"【今晚月色真美】[青之箱 第二季 / アオのハコ Season 2 / Ao no Hako Season 2][03][WebRip][1080P][AVC-8bit][AAC][繁日內嵌]": 3,
		"【今晚月色真美】[青之箱 第二季 / アオのハコ Season 2 / Ao no Hako Season 2][02][WebRip][1080P][AVC-8bit][AAC][简日内嵌]": 2,
		"【今晚月色真美】[青之箱 第二季 / アオのハコ Season 2 / Ao no Hako Season 2][01][WebRip][1080P][AVC-8bit][AAC][简日内嵌]": 1,
	}
	for title, want := range cases {
		got, err := ExtractEpisode(title)
		if err != nil {
			t.Fatalf("%q: %v", title, err)
		}
		if got != want {
			t.Errorf("%q => %v, want %v", title, got, want)
		}
	}
	if _, err := ExtractEpisode("[Sub] Anime Season 2 [1080p]"); err == nil {
		t.Fatal("season number must not be treated as an episode")
	}
}
