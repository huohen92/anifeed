package p115

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// RenameOptions controls how a completed offline entry is renamed.
type RenameOptions struct {
	Enabled  bool
	Template string
	Season   int
	Offset   int
	TmdbID   string
	Subgroup string
}

var (
	renameSubgroupPattern = regexp.MustCompile(`^\s*[\[【]([^\]】]{1,60})[\]】]`)
	renameEpisodePatterns = []*regexp.Regexp{
		// SxxExx (the most reliable form).
		regexp.MustCompile(`(?i)S\d{1,2}\s*E(\d{1,3}(?:\.5)?)`),
		// [03], [03v2], 【03】. Four-digit years/resolutions are not matched.
		regexp.MustCompile(`[\[【]\s*(\d{1,3}(?:\.5)?)\s*(?:v\d+)?\s*[\]】]`),
		// 第05话 / 第05集.
		regexp.MustCompile(`(?i)第\s*(\d{1,3}(?:\.5)?)\s*[话話集]`),
		// " - 01" / " - 01.5" (common Mikan / AniBT form).
		regexp.MustCompile(`(?i)(?:^|\s)-\s*(\d{1,3}(?:\.5)?)\s*(?:v\d+)?(?:\s|\[|\.|$)`),
		// Explicit E01 / EP01.
		regexp.MustCompile(`(?i)(?:^|[\s\[【])(?:E|EP)\s*(\d{1,3}(?:\.5)?)\s*(?:v\d+)?(?:\s|\[|\.|$)`),
	}
	renameResolutionPattern = regexp.MustCompile(`(?i)(\d{3,4}p|4k|8k)`)
	renameTechPattern       = regexp.MustCompile(`(?i)^(?:1080p?|720p?|2160p?|4k|8k|hevc|avc|web-?dl|webrip|bdrip|bluray|flac|aac|mp4|mkv|简|繁|chs|cht|gb|big5|字幕|生肉|熟肉|合集|v\d+)$`)
	renameIllegalChars      = regexp.MustCompile(`[\\/:*?"<>|\x00-\x1f]`)
)

// BuildRenameName builds a filename / folder name from the feed name and the
// original RSS item title.
func BuildRenameName(feedName, itemTitle string, opt RenameOptions) (string, error) {
	if !opt.Enabled {
		return "", nil
	}
	template := strings.TrimSpace(opt.Template)
	if template == "" {
		template = "[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}"
	}
	episode, err := parseEpisode(itemTitle)
	if err != nil {
		return "", err
	}
	episode += float64(opt.Offset)
	if episode < 0 {
		return "", fmt.Errorf("集数偏移后小于 0")
	}
	season := opt.Season
	if season <= 0 {
		season = 1
	}

	episodeText := formatEpisode(episode, false)
	episodeFormat := formatEpisode(episode, true)
	seasonText := strconv.Itoa(season)
	seasonFormat := fmt.Sprintf("%02d", season)

	itemTitle = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(itemTitle, "\n", " "), "\t", " "))
	itemTitle = strings.TrimSuffix(itemTitle, "+NCOPED")
	itemTitle = strings.TrimSpace(itemTitle)
	subgroup := strings.TrimSpace(opt.Subgroup)
	if subgroup == "" {
		subgroup = parseSubgroup(itemTitle)
	}
	feedName = cleanFeedTitle(feedName, subgroup)
	resolution := ""
	if match := renameResolutionPattern.FindStringSubmatch(itemTitle); len(match) == 2 {
		resolution = strings.ToUpper(match[1])
	}

	replacements := map[string]string{
		"${title}":         feedName,
		"${name}":          feedName,
		"${subgroup}":      subgroup,
		"${season}":        seasonText,
		"${seasonFormat}":  seasonFormat,
		"${episode}":       episodeText,
		"${episodeFormat}": episodeFormat,
		"${itemTitle}":     itemTitle,
		"${resolution}":    resolution,
		"${tmdbid}":        opt.TmdbID,
		"${tmdb_id}":       opt.TmdbID,
		"${tmdbId}":        opt.TmdbID,
	}
	for key, value := range replacements {
		template = strings.ReplaceAll(template, key, value)
	}
	template = strings.TrimSpace(template)
	if template == "" {
		return "", fmt.Errorf("重命名模板结果为空")
	}
	return SanitizeFilename(template), nil
}

// ExtractEpisode parses an episode number from an item title.
func ExtractEpisode(title string) (float64, error) {
	return parseEpisode(title)
}

// LooksLikeMediaFile reports whether a name has a common video/file extension.
func LooksLikeMediaFile(name string) bool {
	lower := strings.ToLower(name)
	for _, ext := range []string{".mkv", ".mp4", ".avi", ".ts", ".m2ts", ".webm", ".mov", ".flv", ".wmv", ".m4v"} {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func parseEpisode(title string) (float64, error) {
	for _, pattern := range renameEpisodePatterns {
		match := pattern.FindStringSubmatch(title)
		if len(match) != 2 {
			continue
		}
		value, err := strconv.ParseFloat(match[1], 64)
		if err == nil {
			return value, nil
		}
	}
	return 0, fmt.Errorf("无法从标题解析集数: %s", title)
}

func formatEpisode(episode float64, pad bool) string {
	if pad {
		if episode == float64(int(episode)) {
			return fmt.Sprintf("%02d", int(episode))
		}
		return fmt.Sprintf("%02d.5", int(episode))
	}
	if episode == float64(int(episode)) {
		return fmt.Sprintf("%d", int(episode))
	}
	return fmt.Sprintf("%d.5", int(episode))
}

func cleanFeedTitle(title, subgroup string) string {
	title = strings.TrimSpace(title)
	if subgroup == "" || subgroup == "未知字幕组" {
		return title
	}
	for _, sep := range []string{" · ", " - ", " | "} {
		suffix := sep + subgroup
		if strings.HasSuffix(strings.ToLower(title), strings.ToLower(suffix)) {
			return strings.TrimSpace(title[:len(title)-len(suffix)])
		}
	}
	return title
}

func parseSubgroup(title string) string {
	match := renameSubgroupPattern.FindStringSubmatch(title)
	if len(match) != 2 {
		return "未知字幕组"
	}
	subgroup := strings.TrimSpace(match[1])
	lower := strings.ToLower(subgroup)
	if subgroup == "" || renameTechPattern.MatchString(subgroup) ||
		strings.Contains(lower, "1080p") || strings.Contains(lower, "720p") ||
		strings.Contains(lower, "2160p") || strings.Contains(lower, "web-") ||
		strings.Contains(lower, "webrip") || strings.Contains(lower, "bdrip") {
		return "未知字幕组"
	}
	return subgroup
}

// SanitizeFilename makes a string safe for 115 file/folder names.
func SanitizeFilename(name string) string {
	name = renameIllegalChars.ReplaceAllString(name, " ")
	name = strings.Join(strings.Fields(name), " ")
	name = strings.Trim(name, " .")
	if utf8.RuneCountInString(name) > 200 {
		runes := []rune(name)
		name = string(runes[:200])
	}
	return strings.TrimSpace(name)
}

// PreserveFileExtension appends the original file extension when the template
// did not include one.
func PreserveFileExtension(newName, originalName string, isDir bool) string {
	if isDir {
		return newName
	}
	ext := originalName
	if idx := strings.LastIndex(originalName, "."); idx >= 0 && idx < len(originalName)-1 {
		ext = originalName[idx:]
	} else {
		return newName
	}
	if strings.HasSuffix(strings.ToLower(newName), strings.ToLower(ext)) {
		return newName
	}
	return newName + ext
}
