// Package ruleengine implements ANI-RSS compatible title matching and
// exclusion rules.
//
// A rule is one line of text. Every rule is a regular expression. A rule may
// be scoped to one subtitle group with the {{group}}: prefix:
//
//	{{LoliHouse}}:1080[Pp]
//	^(?=.*\[桜都字幕组\])(?=.*繁体内嵌).*
//	/简体|CHS/i
//
// Match rules are AND: every rule applicable to the current subgroup must
// match. Exclude rules are OR: the first applicable matching rule wins.
// Scoped rules that belong to a different subgroup are ignored completely.
package ruleengine

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

// Rule is one parsed line.
type Rule struct {
	Raw     string
	Group   string
	Pattern string
	re      *regexp2.Regexp
}

// Applies reports whether the rule is active for the given subgroup label.
func (r Rule) Applies(subgroup string) bool {
	if r.Group == "" {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(r.Group), strings.TrimSpace(subgroup))
}

// Match reports whether the title matches the rule.
func (r Rule) Match(title string) bool {
	if r.re == nil {
		return false
	}
	ok, err := r.re.MatchString(title)
	return err == nil && ok
}

// RuleSet is an ordered collection of rules.
type RuleSet struct {
	Rules []Rule
}

// Parse parses newline-separated rules.
//
// Blank lines and lines beginning with # are ignored. Inline /.../ delimiters
// are accepted for compatibility with older AniFeed files. Plain text is also
// a valid regular expression and matches by substring unless it contains
// regexp syntax.
func Parse(raw string) (RuleSet, error) {
	set := RuleSet{}
	blockGroup := ""
	for lineNo, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// AI prompts sometimes emit a group block header followed by one rule
		// per line:
		//   {{今晚月色真美}}:
		//   1080[Pp]
		//   简日内嵌
		// Treat that as a compact block form and apply the group to each line.
		if strings.HasPrefix(line, "{{") {
			if end := strings.Index(line, "}}"); end > 0 && strings.TrimSpace(line[end+2:]) == ":" {
				group := strings.TrimSpace(line[2:end])
				if group == "" {
					return RuleSet{}, fmt.Errorf("第 %d 行：字幕组名称不能为空", lineNo+1)
				}
				blockGroup = group
				continue
			}
		}
		if blockGroup != "" && !strings.HasPrefix(line, "{{") {
			line = "{{" + blockGroup + "}}:" + line
		}
		rule, err := parseLine(line)
		if err != nil {
			return RuleSet{}, fmt.Errorf("第 %d 行：%w", lineNo+1, err)
		}
		set.Rules = append(set.Rules, rule)
	}
	return set, nil
}

func parseLine(line string) (Rule, error) {
	rule := Rule{Raw: line}
	pattern := line
	if strings.HasPrefix(pattern, "{{") {
		end := strings.Index(pattern, "}}")
		if end < 0 || end+2 >= len(pattern) || pattern[end+2] != ':' {
			return Rule{}, fmt.Errorf("字幕组规则格式应为 {{字幕组}}:正则")
		}
		rule.Group = strings.TrimSpace(pattern[2:end])
		pattern = strings.TrimSpace(pattern[end+3:])
		if rule.Group == "" {
			return Rule{}, fmt.Errorf("字幕组名称不能为空")
		}
	}
	if pattern == "" {
		return Rule{}, fmt.Errorf("规则不能为空")
	}
	explicitRegex := strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") && len(pattern) > 2
	if explicitRegex {
		pattern = pattern[1 : len(pattern)-1]
	}
	if strings.TrimSpace(pattern) == "" {
		return Rule{}, fmt.Errorf("规则不能为空")
	}
	re, err := regexp2.Compile(pattern, 0)
	if err != nil {
		// Preserve legacy plain-text rules that are not valid regexes, for
		// example a title fragment containing a lone "[" character.
		if !explicitRegex && !strings.ContainsAny(pattern, "(?)^$\\|()") {
			escaped := regexp.QuoteMeta(pattern)
			if fallback, ferr := regexp2.Compile(escaped, 0); ferr == nil {
				re = fallback
				pattern = escaped
				err = nil
			}
		}
		if err != nil {
			return Rule{}, fmt.Errorf("正则表达式无效：%w", err)
		}
	}
	// Keep user supplied expressions from blocking the scheduler forever.
	re.MatchTimeout = 250 * time.Millisecond
	rule.Pattern = pattern
	rule.re = re
	return rule, nil
}

// Empty reports whether the set has no rules.
func (s RuleSet) Empty() bool { return len(s.Rules) == 0 }

// Applicable returns rules that apply to the subgroup.
func (s RuleSet) Applicable(subgroup string) []Rule {
	out := make([]Rule, 0, len(s.Rules))
	for _, rule := range s.Rules {
		if rule.Applies(subgroup) {
			out = append(out, rule)
		}
	}
	return out
}

// FirstUnmatched returns the first applicable match rule that does not match.
// It returns nil when all applicable rules match, including the empty set.
func (s RuleSet) FirstUnmatched(title, subgroup string) *Rule {
	for i := range s.Rules {
		rule := s.Rules[i]
		if !rule.Applies(subgroup) {
			continue
		}
		if !rule.Match(title) {
			return &rule
		}
	}
	return nil
}

// FirstMatch returns the first applicable rule that matches the title.
func (s RuleSet) FirstMatch(title, subgroup string) *Rule {
	for i := range s.Rules {
		rule := s.Rules[i]
		if !rule.Applies(subgroup) {
			continue
		}
		if rule.Match(title) {
			return &rule
		}
	}
	return nil
}

// Decision describes the result of applying match/exclude/global rules.
type Decision struct {
	Matched    bool
	Excluded   bool
	Reason     string
	Rule       string
	Subgroup   string
	Applicable int
}

// Evaluate applies match, exclude and global exclude sets in ANI-RSS order.
func Evaluate(title, subgroup string, match, exclude, global RuleSet) Decision {
	decision := Decision{Matched: true, Subgroup: subgroup}
	if hit := global.FirstMatch(title, subgroup); hit != nil {
		decision.Matched = false
		decision.Excluded = true
		decision.Reason = "命中全局排除"
		decision.Rule = hit.Raw
		return decision
	}
	if miss := match.FirstUnmatched(title, subgroup); miss != nil {
		decision.Matched = false
		decision.Reason = "未通过匹配"
		decision.Rule = miss.Raw
		return decision
	}
	if hit := exclude.FirstMatch(title, subgroup); hit != nil {
		decision.Matched = false
		decision.Excluded = true
		decision.Reason = "命中排除"
		decision.Rule = hit.Raw
		return decision
	}
	return decision
}
