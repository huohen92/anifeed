package webui

import (
	"testing"

	"github.com/dlclark/regexp2"
)

func TestRulePatternsCompileAndMatch(t *testing.T) {
	resp := rulesResponse{
		Categories: []ruleCategory{
			{ID: "language", Options: []ruleOption{
				{ID: "chs", Pattern: `(?:简体|简中|简日|CHS|GB|SC)`},
				{ID: "cht", Pattern: `(?:繁体|繁中|繁日|CHT|BIG5|TC)`},
			}},
			{ID: "subtitle", Options: []ruleOption{{ID: "hardsub", Pattern: `(?:内嵌|硬字幕|硬字)`}}},
		},
		Groups: []groupRule{
			{ID: "sakurato", Pattern: `^(?=.*\[桜都字幕组\])(?=.*(?:繁体内嵌|繁日内嵌)).*`},
			{ID: "up21", Pattern: `^(?=.*\[Up to 21\°C\])(?!.*\b(?:B-Global)\b).*`},
		},
	}
	for _, cat := range resp.Categories {
		for _, opt := range cat.Options {
			if _, err := regexp2.Compile(opt.Pattern, 0); err != nil {
				t.Errorf("%s pattern: %v", opt.ID, err)
			}
		}
	}
	for _, g := range resp.Groups {
		if _, err := regexp2.Compile(g.Pattern, 0); err != nil {
			t.Errorf("%s pattern: %v", g.ID, err)
		}
	}

	cases := []struct {
		pattern, title string
		want           bool
	}{
		{`(?:简体|简中|简日|CHS|GB|SC)`, "[LoliHouse] 药屋少女的呢喃 - 01 [CHS].mkv", true},
		{`(?:繁体|繁中|繁日|CHT|BIG5|TC)`, "[LoliHouse] 药屋少女的呢喃 - 01 [CHT].mkv", true},
		{`(?:内嵌|硬字幕|硬字)`, "[桜都字幕组] 迷宫饭 [24][1080p][繁体内嵌]", true},
		{`^(?=.*\[桜都字幕组\])(?=.*(?:繁体内嵌|繁日内嵌)).*`, "[桜都字幕组] 迷宫饭 [24][1080p][繁体内嵌]", true},
		{`^(?=.*\[Up to 21\°C\])(?!.*\b(?:B-Global)\b).*`, "[Up to 21°C] 败北女角太多了！ - 07 (Baha 1920x1080 AVC AAC MP4)", true},
		{`^(?=.*\[Up to 21\°C\])(?!.*\b(?:B-Global)\b).*`, "[Up to 21°C] 败犬女主太多了！ - 07 (B-Global 3840x2160 HEVC AAC MKV)", false},
	}
	for _, c := range cases {
		re, err := regexp2.Compile(c.pattern, 0)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := re.MatchString(c.title)
		if got != c.want {
			t.Errorf("pattern %q title %q = %v, want %v", c.pattern, c.title, got, c.want)
		}
	}
}

func TestRuleCompositionSupportsLookaheads(t *testing.T) {
	filter := `^(?=.*(?:(?i)(?:简体|简中|简日|CHS|GB|SC)))(?=.*(?:(?i)(?:内封|内挂|软字幕|字幕内封|CHS&CHT|SC&TC)))(?=.*(?:(?i)(?:MKV|\.mkv))).*`
	re, err := regexp2.Compile(filter, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, title := range []string{
		"[LoliHouse] 药屋少女的呢喃 - 01 [简体内封][1080p].mkv",
		"[Sub] Anime [CHS][内封][MKV]",
	} {
		if ok, _ := re.MatchString(title); !ok {
			t.Errorf("filter did not match %q", title)
		}
	}
	if ok, _ := re.MatchString("[Sub] Anime [CHT][内嵌].mp4"); ok {
		t.Error("filter matched unwanted title")
	}
}
