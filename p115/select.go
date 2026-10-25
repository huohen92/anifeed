package p115

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zhifengle/anifeed/rsssite"
)

var versionPattern = regexp.MustCompile(`(?i)v(\d+)`)

// Selection describes whether one RSS item should be submitted when multiple
// versions of the same episode exist.
type Selection struct {
	Selected bool
	Score    int
	Episode  float64
	Version  int
	Reason   string
}

// SelectEpisodes keeps one item per episode. Items without a parseable episode
// are kept as-is. Ties keep the earliest RSS item.
func SelectEpisodes(items []rsssite.MagnetItem, enabled bool) []Selection {
	result := make([]Selection, len(items))
	if !enabled {
		for i := range items {
			result[i] = Selection{Selected: true}
		}
		return result
	}
	type group struct {
		best    int
		score   int
		version int
	}
	groups := map[int]*group{}
	for i := range items {
		episode, err := ExtractEpisode(items[i].Title)
		score := scoreRelease(items[i].Title)
		version := releaseVersion(items[i].Title)
		result[i] = Selection{Episode: episode, Score: score, Version: version, Selected: true}
		if err != nil {
			continue
		}
		key := int(episode * 2)
		current, ok := groups[key]
		if !ok {
			groups[key] = &group{best: i, score: score, version: version}
			continue
		}
		if score > current.score || (score == current.score && version > current.version) {
			result[current.best].Selected = false
			result[current.best].Reason = fmt.Sprintf("同集已有更优版本（%d 分，v%d）", score, version)
			current.best = i
			current.score = score
			current.version = version
		} else {
			result[i].Selected = false
			result[i].Reason = fmt.Sprintf("同集已有更优版本（%d 分，v%d）", current.score, current.version)
		}
	}
	return result
}

func scoreRelease(title string) int {
	lower := strings.ToLower(title)
	score := 0
	switch {
	case strings.Contains(lower, "2160p") || strings.Contains(lower, "4k") || strings.Contains(lower, "3840x2160"):
		score += 80
	case strings.Contains(lower, "1080p") || strings.Contains(lower, "1920x1080"):
		score += 60
	case strings.Contains(lower, "720p") || strings.Contains(lower, "1280x720"):
		score += 20
	}
	switch {
	case strings.Contains(lower, "简体") || strings.Contains(lower, "简中") || strings.Contains(lower, "简日") || strings.Contains(lower, "chs") || strings.Contains(lower, "gb") || strings.Contains(lower, "sc"):
		score += 30
	case strings.Contains(lower, "简繁") || strings.Contains(lower, "chs&cht") || strings.Contains(lower, "sc&tc"):
		score += 20
	case strings.Contains(lower, "繁体") || strings.Contains(lower, "繁中") || strings.Contains(lower, "cht") || strings.Contains(lower, "big5") || strings.Contains(lower, "tc"):
		score += 10
	}
	if strings.Contains(lower, "内封") || strings.Contains(lower, "软字幕") {
		score += 15
	} else if strings.Contains(lower, "内嵌") || strings.Contains(lower, "硬字幕") {
		score += 8
	}
	if strings.Contains(lower, "hevc") || strings.Contains(lower, "h.265") || strings.Contains(lower, "x265") {
		score += 6
	} else if strings.Contains(lower, "avc") || strings.Contains(lower, "h.264") || strings.Contains(lower, "x264") {
		score += 3
	}
	if strings.Contains(lower, "10bit") || strings.Contains(lower, "10-bit") {
		score += 4
	} else if strings.Contains(lower, "8bit") || strings.Contains(lower, "8-bit") {
		score += 1
	}
	if strings.Contains(lower, "mkv") {
		score += 3
	} else if strings.Contains(lower, "mp4") {
		score += 1
	}
	for index, subgroup := range optsSnapshot().preferredSubgroups {
		subgroup = strings.TrimSpace(subgroup)
		if subgroup == "" {
			continue
		}
		if extract := strings.ToLower(parseSubgroup(title)); extract == strings.ToLower(subgroup) || strings.Contains(lower, strings.ToLower(subgroup)) {
			score += 1000 - index*10
			break
		}
	}
	return score
}

func releaseVersion(title string) int {
	matches := versionPattern.FindAllStringSubmatch(title, -1)
	if len(matches) == 0 {
		return 1
	}
	last := matches[len(matches)-1]
	value := 0
	for _, r := range last[1] {
		if r < '0' || r > '9' {
			break
		}
		value = value*10 + int(r-'0')
	}
	if value == 0 {
		return 1
	}
	return value
}
