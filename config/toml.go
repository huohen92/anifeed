package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/pelletier/go-toml/v2"
)

// TOMLConfig represents the structure of config.toml
type TOMLConfig struct {
	Auth         TOMLAuthConfig            `toml:"auth"`
	Server       TOMLServerConfig          `toml:"server"`
	Database     TOMLDatabaseConfig        `toml:"database"`
	P115         TOMLP115Config            `toml:"p115"`
	Proxy        TOMLProxyConfig           `toml:"proxy"`
	Subscription TOMLSubscriptionConfig    `toml:"subscription"`
	Offline      TOMLOfflineConfig         `toml:"offline"`
	Rename       TOMLRenameConfig          `toml:"rename"`
	TMDB         TOMLTMDBConfig            `toml:"tmdb"`
	Auto         TOMLAutoConfig            `toml:"auto"`
	Notify       TOMLNotifyConfig          `toml:"notify"`
	Lab          TOMLLabConfig             `toml:"lab"`
	RSS          []TOMLRSSConfig           `toml:"rss"`
	Sites        map[string]TOMLSiteConfig `toml:"sites"`
	Compat       TOMLCompatConfig          `toml:"compat"`
}

// TOMLSubscriptionConfig holds source and network defaults.
type TOMLSubscriptionConfig struct {
	MikanHost         string `toml:"mikan_host"`
	RSSTimeoutSeconds int    `toml:"rss_timeout_seconds"`
}

// TOMLOfflineConfig holds defaults applied to every offline task.
type TOMLOfflineConfig struct {
	DefaultCid           string   `toml:"default_cid"`
	DefaultSavePath      string   `toml:"default_savepath"`
	GlobalExclude        []string `toml:"global_exclude"`
	SkipX5               bool     `toml:"skip_x5"`
	Retry                int      `toml:"retry"`
	AutoClearCompleted   bool     `toml:"auto_clear_completed"`
	AutoDisableCompleted bool     `toml:"auto_disable_completed"`
	DualFolder           bool     `toml:"dual_folder"`
	FinalCid             string   `toml:"final_cid"`
	FinalMode            int      `toml:"final_mode"`
	FinalCid2            string   `toml:"final_cid2"`
	FinalPathTemplate    string   `toml:"final_path_template"`
	OnePerEpisode        *bool    `toml:"one_per_episode"`
	KeepOnlyVideo        bool     `toml:"keep_only_video"`
	AutoLanguage         []string `toml:"auto_language"`
	AutoSubtitle         []string `toml:"auto_subtitle"`
	AutoContainer        []string `toml:"auto_container"`
	AutoResolution       []string `toml:"auto_resolution"`
	AutoRules            bool     `toml:"auto_rules"`
}

// TOMLTMDBConfig holds TMDB API settings.
type TOMLTMDBConfig struct {
	API    string `toml:"api"`
	APIKey string `toml:"api_key"`
	Image  string `toml:"image"`
}

// TOMLRenameConfig controls renaming completed offline entries.
type TOMLRenameConfig struct {
	Enabled  bool   `toml:"enabled"`
	Template string `toml:"template"`
	Season   int    `toml:"season"`
	Offset   int    `toml:"offset"`
}

// TOMLAutoConfig controls scheduled RSS runs.
type TOMLAutoConfig struct {
	Enabled         bool `toml:"enabled"`
	IntervalMinutes int  `toml:"interval_minutes"`
	RunOnStart      bool `toml:"run_on_start"`
}

// TOMLNotifyConfig controls webhook notifications.
type TOMLNotifyConfig struct {
	Type                   string   `toml:"type"`
	Webhook                string   `toml:"webhook"`
	Content                []string `toml:"content"`
	Events                 []string `toml:"events"`
	OnSuccess              *bool    `toml:"on_success"`
	OnError                *bool    `toml:"on_error"`
	IdleEnabled            bool     `toml:"idle_enabled"`
	IdleDays               int      `toml:"idle_days"`
	Template               string   `toml:"template"`
	DownloadTimeoutMinutes int      `toml:"download_timeout_minutes"`
	SkipNoChange           bool     `toml:"skip_no_change"`
	HourlyEnabled          bool     `toml:"hourly_enabled"`
	HourlyHours            []int    `toml:"hourly_hours"`
}

// TOMLCompatConfig controls migration behaviour for legacy JSON files.
type TOMLCompatConfig struct {
	// DisableLegacy stops reading rss.json / node-site-config.json / .cookies
	// once config.toml is the source of truth.
	DisableLegacy bool `toml:"disable_legacy"`
}

// TOMLAuthConfig represents authentication configuration
type TOMLAuthConfig struct {
	CookiesFile string `toml:"cookies_file"` // Relative to config.toml directory
	Cookies     string `toml:"cookies"`      // Direct cookie string
}

// TOMLServerConfig represents server configuration
type TOMLServerConfig struct {
	Host             string `toml:"host"`         // Listen address, default "0.0.0.0"
	Port             int    `toml:"port"`         // Valid range: 1-65535
	AccessToken      string `toml:"access_token"` // Optional token required by the Web UI API
	AuthEnabled      bool   `toml:"auth_enabled"`
	AuthUsername     string `toml:"auth_username"`
	AuthPasswordHash string `toml:"auth_password_hash"`
	AuthPasswordSalt string `toml:"auth_password_salt"`
}

// TOMLDatabaseConfig represents database configuration
type TOMLDatabaseConfig struct {
	Path string `toml:"path"` // SQLite database file path
}

// TOMLP115Config represents 115 cloud storage configuration
type TOMLP115Config struct {
	DisableCache  bool `toml:"disable_cache"`
	ChunkDelay    int  `toml:"chunk_delay"`
	ChunkSize     int  `toml:"chunk_size"`
	CooldownMinMs int  `toml:"cooldown_min_ms"`
	CooldownMaxMs int  `toml:"cooldown_max_ms"`
}

// TOMLProxyConfig represents proxy configuration
type TOMLProxyConfig struct {
	HTTP string `toml:"http"` // HTTP proxy URL
}

// TOMLRSSConfig represents a single RSS feed configuration
type TOMLRSSSubgroup struct {
	Name    string `toml:"name"`
	URL     string `toml:"url"`
	Season  int    `toml:"season"`
	Offset  int    `toml:"offset"`
	Filter  string `toml:"filter"`
	Exclude string `toml:"exclude"`
	Enabled *bool  `toml:"enabled"`
}

type TOMLRSSConfig struct {
	Site            string            `toml:"site"`           // Required: host for grouping
	Name            string            `toml:"name"`           // Required: feed name
	SubgroupName    string            `toml:"subgroup_name"`  // Optional: main RSS subgroup label
	BgmURL          string            `toml:"bgm_url"`        // Optional: Bangumi subject URL
	URL             string            `toml:"url"`            // Required: RSS feed URL
	Source          string            `toml:"source"`         // Optional: mikan|anibt|animegarden|other
	Enabled         *bool             `toml:"enabled"`        // Optional: defaults to true
	Cid             string            `toml:"cid"`            // Optional: 115 directory ID
	SavePath        string            `toml:"savepath"`       // Optional: save path
	Filter          string            `toml:"filter"`         // Optional: content filter (whitelist)
	Exclude         string            `toml:"exclude"`        // Optional: content filter (blacklist)
	Expiration      uint              `toml:"expiration"`     // Optional: cache expiration
	RenameSeason    int               `toml:"rename_season"`  // Optional: per-feed rename season
	RenameOffset    int               `toml:"rename_offset"`  // Optional: per-feed episode offset
	UpdateWeekday   int               `toml:"update_weekday"` // Optional: 1=Monday ... 7=Sunday
	TotalEpisodes   int               `toml:"total_episodes"` // Optional: total episodes for completion
	Completed       bool              `toml:"completed"`
	FinalCid        string            `toml:"final_cid"`
	FinalMode       int               `toml:"final_mode"`
	FinalCustomMode int               `toml:"final_custom_mode"`
	Subgroups       []TOMLRSSSubgroup `toml:"subgroups"`
	TmdbID          string            `toml:"tmdb_id"`   // Optional: TMDB subject id
	TmdbType        string            `toml:"tmdb_type"` // Optional: tv|movie
}

// TOMLSiteConfig represents site-specific HTTP configuration
type TOMLSiteConfig struct {
	HTTPSAgent bool              `toml:"https_agent"` // Enable proxy for this site
	Headers    map[string]string `toml:"headers"`     // Custom HTTP headers
}

// FindConfigFile discovers config.toml in standard locations
// Returns path and true if found, empty string and false otherwise
func FindConfigFile() (string, bool) {
	return findFile("config.toml", false)
}

// LoadTOML parses config.toml from the given path
func LoadTOML(path string) (*TOMLConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.toml at %s: %w", path, err)
	}

	cfg, err := ParseTOML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config.toml at %s: %w", path, err)
	}
	return cfg, nil
}

// ParseTOML parses and validates config.toml content.
func ParseTOML(data []byte) (*TOMLConfig, error) {
	var cfg TOMLConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid TOML: %w", err)
	}
	if err := validateTOMLConfig(&cfg, "config.toml"); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validateTOMLConfig validates the parsed TOML configuration
func validateTOMLConfig(cfg *TOMLConfig, path string) error {
	// Validate server host
	if strings.ContainsAny(cfg.Server.Host, " /\\") || strings.Contains(cfg.Server.Host, "://") {
		return fmt.Errorf("invalid server host %q in %s", cfg.Server.Host, path)
	}

	// Validate server port range
	if cfg.Server.Port != 0 && (cfg.Server.Port < 1 || cfg.Server.Port > 65535) {
		return fmt.Errorf("invalid server port %d in %s: must be between 1 and 65535", cfg.Server.Port, path)
	}

	// Validate proxy URL format
	if cfg.Proxy.HTTP != "" {
		u, err := url.Parse(cfg.Proxy.HTTP)
		if err != nil {
			return fmt.Errorf("invalid proxy URL %q in %s: %w", cfg.Proxy.HTTP, path, err)
		}
		// Ensure it has a scheme (http:// or https://)
		if u.Scheme == "" {
			return fmt.Errorf("invalid proxy URL %q in %s: missing scheme (http:// or https://)", cfg.Proxy.HTTP, path)
		}
	}

	// Validate subscription defaults
	if cfg.Subscription.MikanHost != "" {
		u, err := url.Parse(cfg.Subscription.MikanHost)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid subscription.mikan_host %q in %s", cfg.Subscription.MikanHost, path)
		}
	}
	if cfg.Subscription.RSSTimeoutSeconds < 0 || cfg.Subscription.RSSTimeoutSeconds > 300 {
		return fmt.Errorf("subscription.rss_timeout_seconds in %s must be between 1 and 300", path)
	}
	if cfg.Offline.Retry < 0 || cfg.Offline.Retry > 100 {
		return fmt.Errorf("offline.retry in %s must be between 0 and 100", path)
	}
	if cfg.Offline.DualFolder && strings.TrimSpace(cfg.Offline.FinalCid) == "" {
		return fmt.Errorf("offline.final_cid is required when offline.dual_folder is enabled in %s", path)
	}
	if cfg.Rename.Enabled && strings.TrimSpace(cfg.Rename.Template) == "" {
		return fmt.Errorf("rename.template is required when rename.enabled is enabled in %s", path)
	}
	if cfg.Rename.Season < 0 || cfg.Rename.Season > 1000 {
		return fmt.Errorf("rename.season in %s must be between 0 and 1000", path)
	}
	if cfg.Rename.Offset < -1000 || cfg.Rename.Offset > 1000 {
		return fmt.Errorf("rename.offset in %s must be between -1000 and 1000", path)
	}
	if cfg.TMDB.API != "" {
		u, err := url.Parse(cfg.TMDB.API)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid tmdb.api %q in %s", cfg.TMDB.API, path)
		}
	}
	if cfg.TMDB.Image != "" {
		u, err := url.Parse(cfg.TMDB.Image)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid tmdb.image %q in %s", cfg.TMDB.Image, path)
		}
	}
	for i, rss := range cfg.RSS {
		if rss.TmdbType != "" && rss.TmdbType != "tv" && rss.TmdbType != "movie" {
			return fmt.Errorf("RSS entry %d tmdb_type in %s must be tv or movie", i, path)
		}
	}
	for i, rss := range cfg.RSS {
		for _, subgroup := range rss.Subgroups {
			if strings.TrimSpace(subgroup.URL) == "" {
				return fmt.Errorf("rss %q 的备用字幕组缺少 URL", rss.Name)
			}
			if subgroup.Season < 0 || subgroup.Season > 1000 {
				return fmt.Errorf("rss %q 的备用字幕组季数不合法", rss.Name)
			}
			if subgroup.Offset < -1000 || subgroup.Offset > 1000 {
				return fmt.Errorf("rss %q 的备用字幕组集数偏移不合法", rss.Name)
			}
		}
		if rss.TotalEpisodes < 0 || rss.TotalEpisodes > 10000 {
			return fmt.Errorf("rss %q total_episodes 超出范围 [0, 10000]", rss.Name)
		}
		if rss.UpdateWeekday < 0 || rss.UpdateWeekday > 7 {
			return fmt.Errorf("rss %q update_weekday 超出范围 [0, 7]", rss.Name)
		}
		if rss.RenameOffset < -1000 || rss.RenameOffset > 1000 {
			return fmt.Errorf("rss %q rename_offset 超出范围 [-1000, 1000]", rss.Name)
		}
		if rss.RenameSeason < 0 || rss.RenameSeason > 1000 {
			return fmt.Errorf("RSS entry %d rename_season in %s must be between 0 and 1000", i, path)
		}
	}
	for i, rule := range cfg.Offline.GlobalExclude {
		if err := validateFilterRule(rule, fmt.Sprintf("offline.global_exclude[%d]", i)); err != nil {
			return err
		}
	}
	if cfg.Auto.IntervalMinutes < 0 || cfg.Auto.IntervalMinutes > 10080 {
		return fmt.Errorf("auto.interval_minutes in %s must be between 1 and 10080", path)
	}
	for _, section := range cfg.Notify.Content {
		switch section {
		case "summary", "details", "errors", "time":
		default:
			return fmt.Errorf("invalid notify.content %q in %s", section, path)
		}
	}
	switch strings.ToLower(strings.TrimSpace(cfg.Notify.Type)) {
	case "", "auto", "generic", "wecom":
	default:
		return fmt.Errorf("invalid notify.type %q in %s", cfg.Notify.Type, path)
	}
	if cfg.Notify.Webhook != "" {
		u, err := url.Parse(cfg.Notify.Webhook)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return fmt.Errorf("invalid notify.webhook %q in %s", cfg.Notify.Webhook, path)
		}
	}

	// Validate RSS required fields
	for i, rss := range cfg.RSS {
		if rss.Site == "" {
			return fmt.Errorf("RSS entry %d in %s missing required field 'site'", i, path)
		}
		if rss.Source != "" && NormalizeSource(rss.Source, rss.URL) != strings.ToLower(strings.TrimSpace(rss.Source)) {
			switch strings.ToLower(strings.TrimSpace(rss.Source)) {
			case SourceMikan, "mikanani", "mikanime", SourceAniBT, "ani-bt", "anibt.net",
				SourceAnimeGarden, "anime-garden", "animes.garden", SourceOther, "generic", "rss":
			default:
				return fmt.Errorf("RSS entry %d in %s has invalid source %q", i, path, rss.Source)
			}
		}
		if rss.Name == "" {
			return fmt.Errorf("RSS entry %d in %s missing required field 'name'", i, path)
		}
		if rss.URL == "" {
			return fmt.Errorf("RSS entry %d in %s missing required field 'url'", i, path)
		}
		// Validate RSS URL format
		u, err := url.Parse(rss.URL)
		if err != nil {
			return fmt.Errorf("invalid RSS URL %q in entry %d of %s: %w", rss.URL, i, path, err)
		}
		// Ensure it has a scheme (http:// or https://)
		if u.Scheme == "" {
			return fmt.Errorf("invalid RSS URL %q in entry %d of %s: missing scheme (http:// or https://)", rss.URL, i, path)
		}
		if err := validateFilterRule(rss.Filter, fmt.Sprintf("RSS entry %d filter", i)); err != nil {
			return err
		}
		if err := validateFilterRule(rss.Exclude, fmt.Sprintf("RSS entry %d exclude", i)); err != nil {
			return err
		}
	}

	return nil
}

// validateFilterRule ensures a /regex/ filter rule is compilable.
func validateFilterRule(rule, field string) error {
	if !strings.HasPrefix(rule, "/") || !strings.HasSuffix(rule, "/") || len(rule) <= 2 {
		return nil
	}
	if _, err := regexp2.Compile(rule[1:len(rule)-1], 0); err != nil {
		return fmt.Errorf("%s has invalid regular expression %q: %w", field, rule, err)
	}
	return nil
}

// TransformTOMLRSS converts []TOMLRSSConfig to map[string][]RssConfig grouped by site
func TransformTOMLRSS(tomlRSS []TOMLRSSConfig) map[string][]RssConfig {
	result := make(map[string][]RssConfig)

	for _, tr := range tomlRSS {
		enabled := true
		if tr.Enabled != nil {
			enabled = *tr.Enabled
		}
		rssConfig := RssConfig{
			Name:            tr.Name,
			SubgroupName:    tr.SubgroupName,
			BgmURL:          tr.BgmURL,
			Url:             tr.URL,
			Source:          NormalizeSource(tr.Source, tr.URL),
			Enabled:         &enabled,
			Cid:             tr.Cid,
			SavePath:        tr.SavePath,
			Filter:          tr.Filter,
			Exclude:         tr.Exclude,
			Expiration:      tr.Expiration,
			RenameSeason:    tr.RenameSeason,
			RenameOffset:    tr.RenameOffset,
			UpdateWeekday:   tr.UpdateWeekday,
			TotalEpisodes:   tr.TotalEpisodes,
			Completed:       tr.Completed,
			FinalCid:        tr.FinalCid,
			FinalMode:       tr.FinalMode,
			FinalCustomMode: tr.FinalCustomMode,
			Subgroups: func() []RssSubgroup {
				result := make([]RssSubgroup, 0, len(tr.Subgroups))
				for _, group := range tr.Subgroups {
					result = append(result, RssSubgroup{Name: group.Name, URL: group.URL, Season: group.Season, Offset: group.Offset, Filter: group.Filter, Exclude: group.Exclude, Enabled: group.Enabled})
				}
				return result
			}(),
			TmdbID:   tr.TmdbID,
			TmdbType: tr.TmdbType,
		}

		result[tr.Site] = append(result[tr.Site], rssConfig)
	}

	return result
}

// TransformTOMLSites converts map[string]TOMLSiteConfig to map[string]SiteConfig
func TransformTOMLSites(tomlSites map[string]TOMLSiteConfig) map[string]SiteConfig {
	result := make(map[string]SiteConfig)

	for host, ts := range tomlSites {
		siteConfig := SiteConfig{
			Headers: ts.Headers,
		}

		// Convert boolean HTTPSAgent to string for backward compatibility
		if ts.HTTPSAgent {
			siteConfig.HttpsAgent = "true"
		}

		result[host] = siteConfig
	}

	return result
}

// ResolveCookiesFile resolves cookies_file path relative to config.toml directory
// If the path is absolute, it is returned as-is
// If the path is relative, it is resolved relative to the directory containing config.toml
func ResolveCookiesFile(cookiesFile string, tomlPath string) string {
	if cookiesFile == "" {
		return ""
	}

	if filepath.IsAbs(cookiesFile) || strings.HasPrefix(cookiesFile, "/") || isWindowsAbsPath(cookiesFile) {
		return cookiesFile
	}

	if tomlPath != "" {
		if strings.Contains(tomlPath, "/") && !strings.Contains(tomlPath, "\\") {
			return path.Join(path.Dir(tomlPath), cookiesFile)
		}
		return filepath.Join(filepath.Dir(tomlPath), cookiesFile)
	}

	return cookiesFile
}

// isWindowsAbsPath reports whether p looks like a Windows absolute path even
// when the process runs on a non-Windows platform.
func isWindowsAbsPath(p string) bool {
	if len(p) < 3 {
		return false
	}
	c := p[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
		return false
	}
	return p[1] == ':' && (p[2] == '/' || p[2] == '\\')
}

// ResolveDatabasePath resolves database path relative to config.toml directory
// If the path is absolute, it is returned as-is
// If the path is relative, it is resolved relative to the directory containing config.toml
func ResolveDatabasePath(dbPath string, tomlPath string) string {
	if dbPath == "" {
		return ""
	}

	if filepath.IsAbs(dbPath) || isWindowsAbsPath(dbPath) {
		return dbPath
	}

	if tomlPath != "" {
		return filepath.Join(filepath.Dir(tomlPath), dbPath)
	}

	return dbPath
}

// TOMLLabConfig holds experimental options.
type TOMLLabConfig struct {
	MultiSubgroup       bool     `toml:"multi_subgroup"`
	MultiRssCoexist     bool     `toml:"multi_rss_coexist"`
	WashEnabled         bool     `toml:"wash_enabled"`
	WashMode            string   `toml:"wash_mode"`
	WashScheduleEnabled bool     `toml:"wash_schedule_enabled"`
	WashIntervalMinutes int      `toml:"wash_interval_minutes"`
	WashFolders         []string `toml:"wash_folders"`
	PreferredSubgroups  []string `toml:"preferred_subgroups"`
	InactiveSubgroups   []string `toml:"inactive_subgroups"`
	AutoBackupSubgroup  bool     `toml:"auto_backup_subgroup"`
	AutoApplyRules      bool     `toml:"auto_apply_rules"`
	HideHints           bool     `toml:"hide_hints"`
}
