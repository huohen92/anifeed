package webui

import "testing"

func TestWashGroupPrefersListedSubgroup(t *testing.T) {
	candidates := []washCandidate{
		{washFileView: washFileView{ID: "a", Name: "侦探已死.2021.S02E01.第1集-ANi.mkv", Path: "/a/ANi.mkv", Size: 1400, Subgroup: "ANi", Simplified: true}, Series: "侦探已死", Season: 2, Episode: 1},
		{washFileView: washFileView{ID: "b", Name: "侦探已死.2021.S02E01.第1集-LoliHouse.mkv", Path: "/b/LoliHouse.mkv", Size: 1200, Subgroup: "LoliHouse", Simplified: true}, Series: "侦探已死", Season: 2, Episode: 1},
	}
	groups := groupWashCandidates(candidates, []string{"LoliHouse", "ANi"})
	if len(groups) != 1 {
		t.Fatalf("groups = %d", len(groups))
	}
	if groups[0].Keep.Subgroup != "LoliHouse" {
		t.Fatalf("keep = %+v", groups[0].Keep)
	}
	if len(groups[0].Delete) != 1 || groups[0].Delete[0].Subgroup != "ANi" {
		t.Fatalf("delete = %+v", groups[0].Delete)
	}
}

func TestWashParseName(t *testing.T) {
	series, season, episode := washParseName("侦探已死.2021.S02E01.第1集-LoliHouse.mkv")
	if series != "侦探已死.2021" || season != 2 || episode != 1 {
		t.Fatalf("parse = %q %d %v", series, season, episode)
	}
	if !washIsSimplified("侦探已死 简日") || washIsSimplified("侦探已死 繁中") {
		t.Fatal("simplified detection failed")
	}
}
