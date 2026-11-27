package ruleengine

import "testing"

func TestMatchAllApplicableAndIgnoreOtherGroups(t *testing.T) {
	rules, err := Parse(`{{今晚月色真美}}:1080[Pp]
{{今晚月色真美}}:简
{{今晚月色真美}}:日
{{今晚月色真美}}:内嵌
{{今晚月色真美}}:avc|Avc|AVC`)
	if err != nil {
		t.Fatal(err)
	}
	good := "【今晚月色真美】[青之箱 第二季][03][WebRip][1080P][AVC-8bit][AAC][简日内嵌]"
	if miss := rules.FirstUnmatched(good, "今晚月色真美"); miss != nil {
		t.Errorf("good title should pass, missed rule: %s", miss.Raw)
	}
	bad := "【今晚月色真美】[青之箱][03][1080P][AVC][繁日内嵌]"
	if miss := rules.FirstUnmatched(bad, "今晚月色真美"); miss == nil {
		t.Error("traditional title should fail the simplified rule")
	}
	other := "【其他字幕组】[青之箱][03][1080P][AVC][简日内嵌]"
	if miss := rules.FirstUnmatched(other, "其他字幕组"); miss != nil {
		t.Errorf("other subgroup should ignore scoped rules, got: %s", miss.Raw)
	}
}

func TestExcludeAnyAndGroupScope(t *testing.T) {
	rules, err := Parse(`{{桜都字幕组}}:繁体内嵌
{{MingY}}:繁日内嵌`)
	if err != nil {
		t.Fatal(err)
	}
	if hit := rules.FirstMatch("[桜都字幕组] 迷宫饭 [24][1080p][繁体内嵌]", "桜都字幕组"); hit == nil {
		t.Error("sakurato traditional hardsub should be excluded")
	}
	if hit := rules.FirstMatch("[MingY] 前辈是伪娘 [03][1080p][繁日内嵌]", "桜都字幕组"); hit != nil {
		t.Errorf("rule for another subgroup must not apply: %s", hit.Raw)
	}
	if hit := rules.FirstMatch("[MingY] 前辈是伪娘 [03][1080p][简日内嵌]", "MingY"); hit != nil {
		t.Errorf("simplified hardsub should not be excluded: %s", hit.Raw)
	}
}

func TestRawRegexSlashAndInvalid(t *testing.T) {
	if _, err := Parse(`^(?=.*\[LoliHouse\])(?=.*1080p).*`); err != nil {
		t.Fatalf("raw regex should parse: %v", err)
	}
	if _, err := Parse(`/简体|1080p/`); err != nil {
		t.Fatalf("slash regex should parse: %v", err)
	}
	if _, err := Parse(`/[abc/`); err == nil {
		t.Error("invalid regex should fail")
	}
	if _, err := Parse(`{{LoliHouse}}1080p`); err == nil {
		t.Error("malformed subgroup prefix should fail")
	}
}

func TestEvaluateOrderAndGlobal(t *testing.T) {
	match, _ := Parse(`1080[Pp]`)
	exclude, _ := Parse(`{{LoliHouse}}:繁体内嵌`)
	global, _ := Parse(`720[Pp]`)
	good := Evaluate("[LoliHouse] 某番 [01][1080p][简体内嵌]", "LoliHouse", match, exclude, global)
	if !good.Matched || good.Excluded {
		t.Errorf("good item should pass: %#v", good)
	}
	bad := Evaluate("[LoliHouse] 某番 [01][1080p][繁体内嵌]", "LoliHouse", match, exclude, global)
	if bad.Matched || !bad.Excluded {
		t.Errorf("exclude should win: %#v", bad)
	}
	globalHit := Evaluate("[LoliHouse] 某番 [01][720p][简体内嵌]", "LoliHouse", match, exclude, global)
	if !globalHit.Excluded || globalHit.Reason != "命中全局排除" {
		t.Errorf("global exclude should apply: %#v", globalHit)
	}
}

func TestLiteralFallbackForLegacyRule(t *testing.T) {
	// A legacy plain-text rule containing "[" is not a valid regexp, but must
	// remain usable instead of making the whole configuration fail to load.
	rules, err := Parse(`[简体内嵌]`)
	if err != nil {
		t.Fatal(err)
	}
	if hit := rules.FirstMatch("某番 [简体内嵌] 1080p", ""); hit == nil {
		t.Error("legacy literal rule should still match")
	}
}

func TestGroupBlockSyntax(t *testing.T) {
	rules, err := Parse(`{{今晚月色真美}}:
1080[Pp]
简日内嵌`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules.Rules))
	}
	if rules.Rules[0].Group != "今晚月色真美" || rules.Rules[1].Group != "今晚月色真美" {
		t.Fatalf("block group was not applied: %#v", rules.Rules)
	}
	if miss := rules.FirstUnmatched("[今晚月色真美][01][1080P][简日内嵌]", "今晚月色真美"); miss != nil {
		t.Errorf("unexpected miss: %s", miss.Raw)
	}
}
