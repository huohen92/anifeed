package webui

import "testing"

func TestCleanTMDBQuery(t *testing.T) {
	cases := map[string]string{
		"青之箱 第二季 · 今晚月色真美":             "青之箱",
		"佐佐木与文鸟小哔 第二季 · LoliHouse":     "佐佐木与文鸟小哔",
		"魔法少女育成计划 restart · LoliHouse": "魔法少女育成计划 restart",
	}
	for input, want := range cases {
		if got := cleanTMDBQuery(input); got != want {
			t.Errorf("cleanTMDBQuery(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestBestTMDBMatch(t *testing.T) {
	results := []tmdbResult{
		{ID: 1, Name: "青之箱", Original: "Ao no Hako"},
		{ID: 2, Name: "青之箱 第二季", Original: "Ao no Hako Season 2"},
	}
	item, score := bestTMDBMatch("青之箱", results)
	if item.ID != 1 || score < 500 {
		t.Fatalf("match = %+v score=%d", item, score)
	}
	if item, score := bestTMDBMatch("完全无关的名字", results); score >= 500 || item.ID != 0 {
		t.Fatalf("unexpected match = %+v score=%d", item, score)
	}
}
