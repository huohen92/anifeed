package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Subscription sources supported by the Web UI.
const (
	SourceMikan       = "mikan"
	SourceAniBT       = "anibt"
	SourceAnimeGarden = "animegarden"
	SourceOther       = "other"
)

// OfflineConfig holds defaults applied to every RSS subscription.
type OfflineConfig struct {
	DefaultCid           string
	DefaultSavePath      string
	GlobalExclude        []string
	SkipX5               bool
	Retry                int
	AutoClearCompleted   bool
	AutoDisableCompleted bool
	// DualFolder enables moving a completed task from the offline folder to
	// FinalCid after renaming.
	DualFolder        bool
	FinalCid          string
	FinalMode         int
	FinalCid2         string
	FinalPathTemplate string
	OnePerEpisode     bool
	KeepOnlyVideo     bool
	AutoRules         bool

	// AutoPriority orders the options used by "apply best rules".
	AutoLanguage   []string
	AutoSubtitle   []string
	AutoContainer  []string
	AutoResolution []string
}

// RenameConfig controls renaming completed offline entries.
type RenameConfig struct {
	Enabled  bool
	Template string
	Season   int
	Offset   int
}

// AutoConfig controls scheduled RSS runs.
type AutoConfig struct {
	Enabled         bool
	IntervalMinutes int
	RunOnStart      bool
}

// NotifyConfig controls the optional webhook notification.
type NotifyConfig struct {
	Type                   string // auto | generic | wecom
	Webhook                string
	Content                []string // summary | details | errors | time
	Events                 []string // run | added | completed | failed | idle | hourly | wash
	OnSuccess              bool
	OnError                bool
	IdleEnabled            bool
	IdleDays               int
	Template               string
	DownloadTimeoutMinutes int
	SkipNoChange           bool
	HourlyEnabled          bool
	HourlyHours            []int
}

// TMDBConfig holds TMDB API settings.
type TMDBConfig struct {
	API    string
	APIKey string
	Image  string
}

// SubscriptionConfig holds source and network defaults.
type SubscriptionConfig struct {
	MikanHost         string
	RSSTimeoutSeconds int
}

// DetectSource infers a subscription source from its RSS URL.
func DetectSource(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return SourceOther
	}
	switch strings.ToLower(u.Hostname()) {
	case "mikanani.me", "mikanime.tv":
		return SourceMikan
	case "anibt.net":
		return SourceAniBT
	case "animes.garden", "api.animes.garden":
		return SourceAnimeGarden
	default:
		return SourceOther
	}
}

// NormalizeSource validates a source name and falls back to URL detection.
func NormalizeSource(source, rawURL string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case SourceMikan, "mikanani", "mikanime":
		return SourceMikan
	case SourceAniBT, "ani-bt", "anibt.net":
		return SourceAniBT
	case SourceAnimeGarden, "anime-garden", "animes.garden":
		return SourceAnimeGarden
	case SourceOther, "generic", "rss":
		return SourceOther
	default:
		return DetectSource(rawURL)
	}
}

// SourceLabel returns a human-readable source name.
func SourceLabel(source string) string {
	switch source {
	case SourceMikan:
		return "Mikan"
	case SourceAniBT:
		return "AniBT"
	case SourceAnimeGarden:
		return "AnimeGarden"
	default:
		return "其他"
	}
}

// ExpandSavePath expands a save path template such as "番剧/${title}".
// Supported variables: ${title}, ${name}, ${site}, ${year}, ${month},
// ${day}, ${season}. Substituted values are sanitised so they remain valid
// 115 folder names.
func ExpandSavePath(template, name, site string, now time.Time) string {
	template = strings.TrimSpace(template)
	if template == "" {
		return ""
	}
	season := fmt.Sprintf("%d%s", now.Year(), seasonName(now.Month()))
	replacements := map[string]string{
		"${title}":  sanitizePathSegment(name),
		"${name}":   sanitizePathSegment(name),
		"${site}":   sanitizePathSegment(site),
		"${year}":   fmt.Sprintf("%d", now.Year()),
		"${month}":  fmt.Sprintf("%02d", int(now.Month())),
		"${day}":    fmt.Sprintf("%02d", now.Day()),
		"${season}": sanitizePathSegment(season),
	}
	for key, value := range replacements {
		template = strings.ReplaceAll(template, key, value)
	}
	template = strings.Trim(strings.TrimSpace(template), "/")
	return template
}

// ExpandFinalPath expands the final-folder template and additionally supports
// the anime season number, e.g. ${animeSeason} -> "Season 2".
func ExpandFinalPath(template, name, site string, seasonNumber int, now time.Time) string {
	return ExpandFinalPathWithMeta(template, name, site, seasonNumber, now.Year(), "", now)
}

// ExpandFinalPathWithMeta also supports ${titleYearTmdb} and ${tmdbid}.
func ExpandFinalPathWithMeta(template, name, site string, seasonNumber, year int, tmdbID string, now time.Time) string {
	expanded := ExpandSavePath(template, name, site, now)
	if seasonNumber <= 0 {
		seasonNumber = 1
	}
	if year <= 0 {
		year = now.Year()
	}
	expanded = strings.ReplaceAll(expanded, "${animeSeason}", sanitizePathSegment(fmt.Sprintf("Season %d", seasonNumber)))
	expanded = strings.ReplaceAll(expanded, "${animeSeasonCN}", sanitizePathSegment(fmt.Sprintf("第%d季", seasonNumber)))
	expanded = strings.ReplaceAll(expanded, "${quarter}", sanitizePathSegment(seasonName(now.Month())))
	expanded = strings.ReplaceAll(expanded, "${seasonNumber}", strconv.Itoa(seasonNumber))
	tmdbID = strings.TrimSpace(tmdbID)
	titleWithYear := fmt.Sprintf("%s (%d)", strings.TrimSpace(name), year)
	if tmdbID != "" {
		titleWithYear += fmt.Sprintf(" {tmdb-%s}", tmdbID)
	}
	expanded = strings.ReplaceAll(expanded, "${titleYearTmdb}", sanitizePathSegment(titleWithYear))
	expanded = strings.ReplaceAll(expanded, "${tmdbid}", sanitizePathSegment(tmdbID))
	return strings.Trim(strings.TrimSpace(expanded), "/")
}

func seasonName(month time.Month) string {
	switch {
	case month >= 1 && month <= 3:
		return "冬"
	case month >= 4 && month <= 6:
		return "春"
	case month >= 7 && month <= 9:
		return "夏"
	default:
		return "秋"
	}
}

func sanitizePathSegment(value string) string {
	value = strings.TrimSpace(value)
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_", "?", "_",
		"\"", "_", "<", "_", ">", "_", "|", "_",
	)
	value = replacer.Replace(value)
	value = strings.Trim(value, ". ")
	runes := []rune(value)
	if len(runes) > 80 {
		value = string(runes[:80])
	}
	if value == "" {
		value = "未命名"
	}
	return value
}

// DefaultOfflineConfig returns the defaults used when config.toml does not
// contain an [offline] section.
func DefaultOfflineConfig() OfflineConfig {
	language, subtitle, container, resolution := DefaultAutoPriority()
	return OfflineConfig{
		Retry:          0,
		OnePerEpisode:  false,
		AutoLanguage:   language,
		AutoSubtitle:   subtitle,
		AutoContainer:  container,
		AutoResolution: resolution,
	}
}

// DefaultRenameConfig returns the default rename template.
func DefaultRenameConfig() RenameConfig {
	return RenameConfig{
		Template: "[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}",
		Season:   1,
	}
}

// DefaultAutoConfig returns the default scheduled-run configuration.
func DefaultAutoConfig() AutoConfig {
	return AutoConfig{IntervalMinutes: 30}
}

// DefaultNotifyConfig returns the default notification configuration.
func DefaultNotifyConfig() NotifyConfig {
	return NotifyConfig{
		Type:      "auto",
		Content:   DefaultNotifyContent(),
		Events:    DefaultNotifyEvents(),
		OnSuccess: true,
		OnError:   true,
	}
}

// DefaultNotifyContent returns the sections included in text notifications.
func DefaultNotifyContent() []string {
	return []string{"summary", "details", "errors", "time"}
}

// DefaultTMDBConfig returns the default TMDB endpoints.
// DefaultSubscriptionConfig returns source and network defaults.
func DefaultTMDBConfig() TMDBConfig {
	return TMDBConfig{API: "https://api.themoviedb.org", Image: "https://image.tmdb.org"}
}

func DefaultSubscriptionConfig() SubscriptionConfig {
	return SubscriptionConfig{
		MikanHost:         "https://mikanani.me",
		RSSTimeoutSeconds: 20,
	}
}

// DefaultNotifyTemplate is used for per-item update notifications.
func DefaultNotifyTemplate() string {
	return `通知
事件类型: ${action}
标题: ${title}
评分: ${score}
TMDB标题: ${themoviedbName}
BGM: ${bgmUrl}
季: ${season}
集: ${episode}
字幕组: ${subgroup}
进度: ${currentEpisodeNumber}/${totalEpisodeNumber}
首播:  ${year}年${month}月${date}日
事件: ${text}
BGM集标题: ${bgmEpisodeTitle}`
}

// DefaultNotifyEvents returns all supported notification event keys.
func DefaultNotifyEvents() []string {
	return []string{"run", "added", "completed", "failed", "idle", "hourly", "wash"}
}

func normalizeNotifyEvents(values []string) []string {
	allowed := map[string]bool{"run": true, "added": true, "completed": true, "failed": true, "idle": true, "hourly": true, "wash": true}
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if allowed[value] && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}
