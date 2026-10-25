package config

import (
	"encoding/json"
	"os"
	"strings"
)

// RssConfig represents RSS feed configuration (local copy to avoid circular dependency)
// RssSubgroup is an optional extra RSS source for one subscription.
type RssSubgroup struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Season  int    `json:"season,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Filter  string `json:"filter,omitempty"`
	Exclude string `json:"exclude,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

func (s RssSubgroup) IsEnabled() bool {
	return s.Enabled == nil || *s.Enabled
}

type RssConfig struct {
	Name            string        `json:"name"`
	SubgroupName    string        `json:"subgroup_name,omitempty"`
	BgmURL          string        `json:"bgm_url,omitempty"`
	Url             string        `json:"url"`
	Source          string        `json:"source,omitempty"`
	Enabled         *bool         `json:"enabled,omitempty"`
	Cid             string        `json:"cid,omitempty"`
	SavePath        string        `json:"savepath,omitempty"`
	Filter          string        `json:"filter,omitempty"`
	Exclude         string        `json:"exclude,omitempty"`
	Expiration      uint          `json:"expiration,omitempty"`
	RenameSeason    int           `json:"rename_season,omitempty"`
	RenameOffset    int           `json:"rename_offset,omitempty"`
	UpdateWeekday   int           `json:"update_weekday,omitempty"`
	TotalEpisodes   int           `json:"total_episodes,omitempty"`
	Completed       bool          `json:"completed,omitempty"`
	FinalCid        string        `json:"final_cid,omitempty"`
	FinalMode       int           `json:"final_mode,omitempty"`
	FinalCustomMode int           `json:"final_custom_mode,omitempty"`
	Subgroups       []RssSubgroup `json:"subgroups,omitempty"`
	TmdbID          string        `json:"tmdb_id,omitempty"`
	TmdbType        string        `json:"tmdb_type,omitempty"`
}

// IsEnabled reports whether a feed participates in scheduled/all-feed runs.
// Missing enabled fields default to true for backward compatibility.
func (r RssConfig) IsEnabled() bool {
	return r.Enabled == nil || *r.Enabled
}

// SiteConfig represents site-specific HTTP configuration (local copy to avoid circular dependency)
type SiteConfig struct {
	HttpsAgent string            `json:"httpsAgent,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// P115Option represents 115 cloud storage options (local copy to avoid circular dependency)
type P115Option struct {
	DisableCache  bool
	ChunkDelay    int
	ChunkSize     int
	CooldownMinMs int
	CooldownMaxMs int
}

// Config represents the unified configuration
type Config struct {
	Auth         AuthConfig
	Server       ServerConfig
	Database     DatabaseConfig
	P115         P115Config
	Proxy        ProxyConfig
	Subscription SubscriptionConfig
	Offline      OfflineConfig
	Rename       RenameConfig
	TMDB         TMDBConfig
	Auto         AutoConfig
	Notify       NotifyConfig
	Lab          LabConfig
	RSS          map[string][]RssConfig // Keyed by site host
	Sites        map[string]SiteConfig  // Keyed by site host
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Cookies string
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host             string
	Port             int
	AccessToken      string
	AuthEnabled      bool
	AuthUsername     string
	AuthPasswordHash string
	AuthPasswordSalt string
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string
}

// P115Config represents 115 cloud storage configuration
type P115Config struct {
	DisableCache  bool
	ChunkDelay    int
	ChunkSize     int
	CooldownMinMs int
	CooldownMaxMs int
}

// ProxyConfig represents proxy configuration
type ProxyConfig struct {
	HTTP string
}

// ConfigSource tracks the source of loaded configuration for cookie save operations
type ConfigSource struct {
	TOMLPath    string // Path to config.toml if loaded
	CookiesPath string // Path to cookies file that was read
}

// CLIParams captures command-line parameters
type CLIParams struct {
	Cookies          string
	RSSPath          string
	Host             string
	HostSet          bool
	Port             int
	PortSet          bool
	DisableCache     bool
	DisableCacheSet  bool
	ChunkDelay       int
	ChunkDelaySet    bool
	ChunkSize        int
	ChunkSizeSet     bool
	CooldownMinMs    int
	CooldownMinMsSet bool
	CooldownMaxMs    int
	CooldownMaxMsSet bool
}

// LegacyConfig holds loaded legacy configuration
type LegacyConfig struct {
	RSS         map[string][]RssConfig
	Sites       map[string]SiteConfig
	Cookies     string
	CookiesPath string
}

// LoadOptions controls which file-backed sections should be loaded.
type LoadOptions struct {
	RSS   bool
	Sites bool
	Auth  bool
}

func LoadAllOptions() LoadOptions {
	return LoadOptions{RSS: true, Sites: true, Auth: true}
}

// ToP115Option converts P115Config to P115Option
func (c *P115Config) ToP115Option() P115Option {
	return P115Option{
		DisableCache:  c.DisableCache,
		ChunkDelay:    c.ChunkDelay,
		ChunkSize:     c.ChunkSize,
		CooldownMinMs: c.CooldownMinMs,
		CooldownMaxMs: c.CooldownMaxMs,
	}
}

// Resolve merges configuration from multiple sources according to priority rules
// Priority: CLI > TOML > Legacy > Defaults
func Resolve(cli CLIParams, toml *TOMLConfig, tomlPath string, legacy *LegacyConfig) *Config {
	cfg := &Config{
		RSS:          make(map[string][]RssConfig),
		Sites:        make(map[string]SiteConfig),
		Subscription: DefaultSubscriptionConfig(),
		Offline:      DefaultOfflineConfig(),
		Rename:       DefaultRenameConfig(),
		TMDB:         DefaultTMDBConfig(),
		Auto:         DefaultAutoConfig(),
		Notify:       DefaultNotifyConfig(),
	}

	// Auth priority: CLI cookies > TOML cookies > TOML cookies_file > legacy .cookies
	if cli.Cookies != "" {
		cfg.Auth.Cookies = cli.Cookies
	} else if toml != nil && toml.Auth.Cookies != "" {
		cfg.Auth.Cookies = toml.Auth.Cookies
	} else if toml != nil && toml.Auth.CookiesFile != "" {
		// Read cookies from file specified in TOML
		cookiesPath := ResolveCookiesFile(toml.Auth.CookiesFile, tomlPath)
		if data, err := os.ReadFile(cookiesPath); err == nil {
			cfg.Auth.Cookies = string(data)
		}
	} else if legacy != nil {
		cfg.Auth.Cookies = legacy.Cookies
	}

	// Server priority: CLI host/port > TOML host/port > defaults
	if cli.HostSet {
		cfg.Server.Host = cli.Host
	} else if toml != nil && toml.Server.Host != "" {
		cfg.Server.Host = toml.Server.Host
	} else {
		cfg.Server.Host = "0.0.0.0" // Default: listen on all interfaces
	}
	if cli.PortSet || cli.Port != 0 {
		cfg.Server.Port = cli.Port
	} else if toml != nil && toml.Server.Port != 0 {
		cfg.Server.Port = toml.Server.Port
	} else {
		cfg.Server.Port = 8115 // Default
	}
	if toml != nil {
		cfg.Server.AccessToken = toml.Server.AccessToken
		cfg.Server.AuthEnabled = toml.Server.AuthEnabled
		cfg.Server.AuthUsername = toml.Server.AuthUsername
		cfg.Server.AuthPasswordHash = toml.Server.AuthPasswordHash
		cfg.Server.AuthPasswordSalt = toml.Server.AuthPasswordSalt
	}

	// Database priority: TOML database.path > existing db.sqlite > default "db.sqlite"
	if toml != nil && toml.Database.Path != "" {
		// Resolve relative path based on config.toml directory
		cfg.Database.Path = ResolveDatabasePath(toml.Database.Path, tomlPath)
	} else if path, ok := findFile("db.sqlite", false); ok {
		cfg.Database.Path = path
	} else {
		cfg.Database.Path = "db.sqlite" // Default
	}

	// P115 priority: CLI params > TOML > defaults
	// DisableCache: CLI or TOML (boolean OR logic)
	cfg.P115.DisableCache = (cli.DisableCacheSet && cli.DisableCache) || (!cli.DisableCacheSet && cli.DisableCache) || (toml != nil && toml.P115.DisableCache)

	// ChunkDelay: CLI > TOML > default 2
	if cli.ChunkDelaySet || cli.ChunkDelay != 0 {
		cfg.P115.ChunkDelay = cli.ChunkDelay
	} else if toml != nil && toml.P115.ChunkDelay != 0 {
		cfg.P115.ChunkDelay = toml.P115.ChunkDelay
	} else {
		cfg.P115.ChunkDelay = 2 // Default
	}

	// ChunkSize: CLI > TOML > default 200
	if cli.ChunkSizeSet || cli.ChunkSize != 0 {
		cfg.P115.ChunkSize = cli.ChunkSize
	} else if toml != nil && toml.P115.ChunkSize != 0 {
		cfg.P115.ChunkSize = toml.P115.ChunkSize
	} else {
		cfg.P115.ChunkSize = 200 // Default
	}

	// CooldownMinMs: CLI > TOML > default 1000
	if cli.CooldownMinMsSet || cli.CooldownMinMs != 0 {
		cfg.P115.CooldownMinMs = cli.CooldownMinMs
	} else if toml != nil && toml.P115.CooldownMinMs != 0 {
		cfg.P115.CooldownMinMs = toml.P115.CooldownMinMs
	} else {
		cfg.P115.CooldownMinMs = 1000 // Default
	}

	// CooldownMaxMs: CLI > TOML > default 1100
	if cli.CooldownMaxMsSet || cli.CooldownMaxMs != 0 {
		cfg.P115.CooldownMaxMs = cli.CooldownMaxMs
	} else if toml != nil && toml.P115.CooldownMaxMs != 0 {
		cfg.P115.CooldownMaxMs = toml.P115.CooldownMaxMs
	} else {
		cfg.P115.CooldownMaxMs = 1100 // Default
	}

	// RSS priority: CLI --rss path > TOML [[rss]] > legacy rss.json
	if cli.RSSPath != "" {
		// Load RSS from CLI-specified path
		// This will be handled by the caller (Load function)
		// For now, we just mark that CLI takes priority
		cfg.RSS = make(map[string][]RssConfig)
	} else if toml != nil && len(toml.RSS) > 0 {
		cfg.RSS = TransformTOMLRSS(toml.RSS)
	} else if legacy != nil && legacy.RSS != nil {
		cfg.RSS = legacy.RSS
	}

	// Sites priority: TOML [sites] > legacy node-site-config.json
	if toml != nil && len(toml.Sites) > 0 {
		cfg.Sites = TransformTOMLSites(toml.Sites)
	} else if legacy != nil && legacy.Sites != nil {
		cfg.Sites = legacy.Sites
	}

	// Proxy priority: TOML [proxy].http > default "http://127.0.0.1:10809"
	if toml != nil && toml.Proxy.HTTP != "" {
		cfg.Proxy.HTTP = toml.Proxy.HTTP
	} else {
		cfg.Proxy.HTTP = "http://127.0.0.1:10809" // Default
	}

	if toml != nil {
		// [subscription]
		if toml.Subscription.MikanHost != "" {
			cfg.Subscription.MikanHost = toml.Subscription.MikanHost
		}
		if toml.Subscription.RSSTimeoutSeconds > 0 {
			cfg.Subscription.RSSTimeoutSeconds = toml.Subscription.RSSTimeoutSeconds
		}
		// [offline]
		cfg.Offline.DefaultCid = toml.Offline.DefaultCid
		cfg.Offline.DefaultSavePath = toml.Offline.DefaultSavePath
		cfg.Offline.GlobalExclude = append([]string(nil), toml.Offline.GlobalExclude...)
		cfg.Offline.SkipX5 = toml.Offline.SkipX5
		cfg.Offline.Retry = toml.Offline.Retry
		cfg.Offline.AutoClearCompleted = toml.Offline.AutoClearCompleted
		cfg.Offline.AutoDisableCompleted = toml.Offline.AutoDisableCompleted
		cfg.Offline.DualFolder = toml.Offline.DualFolder
		cfg.Offline.FinalCid = toml.Offline.FinalCid
		cfg.Offline.FinalMode = toml.Offline.FinalMode
		if cfg.Offline.FinalMode != 2 {
			cfg.Offline.FinalMode = 1
		}
		cfg.Offline.FinalCid2 = toml.Offline.FinalCid2
		if strings.TrimSpace(toml.Offline.FinalPathTemplate) != "" {
			cfg.Offline.FinalPathTemplate = toml.Offline.FinalPathTemplate
		} else {
			cfg.Offline.FinalPathTemplate = "${month}/${titleYearTmdb}/${animeSeason}"
		}
		if toml.Offline.OnePerEpisode != nil {
			cfg.Offline.OnePerEpisode = *toml.Offline.OnePerEpisode
		}
		cfg.Offline.KeepOnlyVideo = toml.Offline.KeepOnlyVideo
		language, subtitle, container, resolution := DefaultAutoPriority()
		cfg.Offline.AutoLanguage = normalizePriority(toml.Offline.AutoLanguage, language, map[string]bool{"chs": true, "cht": true, "bilingual": true})
		cfg.Offline.AutoSubtitle = normalizePriority(toml.Offline.AutoSubtitle, subtitle, map[string]bool{"softsub": true, "hardsub": true, "external": true})
		cfg.Offline.AutoContainer = normalizePriority(toml.Offline.AutoContainer, container, map[string]bool{"mkv": true, "mp4": true})
		cfg.Offline.AutoResolution = normalizePriority(toml.Offline.AutoResolution, resolution, map[string]bool{"r1080": true, "r2160": true, "r720": true})
		cfg.Offline.AutoRules = toml.Offline.AutoRules
		// [rename]
		cfg.Rename.Enabled = toml.Rename.Enabled
		if strings.TrimSpace(toml.Rename.Template) != "" {
			cfg.Rename.Template = toml.Rename.Template
		}
		if toml.Rename.Season > 0 {
			cfg.Rename.Season = toml.Rename.Season
		}
		cfg.Rename.Offset = toml.Rename.Offset
		// [tmdb]
		cfg.TMDB.API = toml.TMDB.API
		if cfg.TMDB.API == "" {
			cfg.TMDB.API = "https://api.themoviedb.org"
		}
		cfg.TMDB.APIKey = toml.TMDB.APIKey
		cfg.TMDB.Image = toml.TMDB.Image
		if cfg.TMDB.Image == "" {
			cfg.TMDB.Image = "https://image.tmdb.org"
		}
		// [auto]
		cfg.Auto.Enabled = toml.Auto.Enabled
		if toml.Auto.IntervalMinutes > 0 {
			cfg.Auto.IntervalMinutes = toml.Auto.IntervalMinutes
		}
		cfg.Auto.RunOnStart = toml.Auto.RunOnStart
		// [notify]
		cfg.Notify.Type = toml.Notify.Type
		if strings.TrimSpace(cfg.Notify.Type) == "" {
			cfg.Notify.Type = "auto"
		}
		cfg.Notify.Webhook = toml.Notify.Webhook
		if len(toml.Notify.Content) > 0 {
			cfg.Notify.Content = append([]string(nil), toml.Notify.Content...)
		}
		cfg.Notify.Events = normalizeNotifyEvents(toml.Notify.Events)
		if toml.Notify.OnSuccess != nil {
			cfg.Notify.OnSuccess = *toml.Notify.OnSuccess
		}
		if toml.Notify.OnError != nil {
			cfg.Notify.OnError = *toml.Notify.OnError
		}
		cfg.Notify.IdleEnabled = toml.Notify.IdleEnabled
		if toml.Notify.IdleDays > 0 {
			cfg.Notify.IdleDays = toml.Notify.IdleDays
		}
		if strings.TrimSpace(toml.Notify.Template) != "" {
			cfg.Notify.Template = toml.Notify.Template
		}
		if toml.Notify.DownloadTimeoutMinutes > 0 {
			cfg.Notify.DownloadTimeoutMinutes = toml.Notify.DownloadTimeoutMinutes
		}
		cfg.Notify.SkipNoChange = toml.Notify.SkipNoChange
		cfg.Notify.HourlyEnabled = toml.Notify.HourlyEnabled
		cfg.Notify.HourlyHours = normalizeHours(toml.Notify.HourlyHours)
		cfg.Lab = LabConfig{
			MultiSubgroup:       toml.Lab.MultiSubgroup,
			MultiRssCoexist:     toml.Lab.MultiRssCoexist,
			WashEnabled:         toml.Lab.WashEnabled,
			WashMode:            normalizeWashMode(toml.Lab.WashMode),
			WashScheduleEnabled: toml.Lab.WashScheduleEnabled,
			WashIntervalMinutes: func() int {
				if toml.Lab.WashIntervalMinutes > 0 {
					return toml.Lab.WashIntervalMinutes
				}
				return 1440
			}(),
			WashFolders:        append([]string(nil), toml.Lab.WashFolders...),
			PreferredSubgroups: append([]string(nil), toml.Lab.PreferredSubgroups...),
			InactiveSubgroups:  append([]string(nil), toml.Lab.InactiveSubgroups...),
			AutoBackupSubgroup: toml.Lab.AutoBackupSubgroup,
			AutoApplyRules:     toml.Lab.AutoApplyRules,
			HideHints:          toml.Lab.HideHints,
		}
	}

	// Normalise feed metadata (legacy JSON has no source/enabled fields).
	for site, feeds := range cfg.RSS {
		for i := range feeds {
			if feeds[i].Source == "" {
				feeds[i].Source = DetectSource(feeds[i].Url)
			} else {
				feeds[i].Source = NormalizeSource(feeds[i].Source, feeds[i].Url)
			}
			if feeds[i].Enabled == nil {
				enabled := true
				feeds[i].Enabled = &enabled
			}
		}
		cfg.RSS[site] = feeds
	}

	return cfg
}

// Load discovers and loads configuration from all sources
// Returns merged configuration and source tracking information
func Load(cliParams CLIParams) (*Config, *ConfigSource, error) {
	return LoadWithOptions(cliParams, LoadAllOptions())
}

// LoadWithOptions discovers and loads only the requested file-backed sections.
func LoadWithOptions(cliParams CLIParams, options LoadOptions) (*Config, *ConfigSource, error) {
	source := &ConfigSource{}

	// Step 1: Try to find and load config.toml
	var tomlConfig *TOMLConfig
	var tomlPath string
	if path, found := FindConfigFile(); found {
		cfg, err := LoadTOML(path)
		if err != nil {
			return nil, nil, err
		}
		tomlConfig = cfg
		tomlPath = path
		source.TOMLPath = path
	}

	if tomlConfig != nil {
		tomlCopy := *tomlConfig
		if !options.RSS {
			tomlCopy.RSS = nil
		}
		if !options.Sites {
			tomlCopy.Sites = nil
		}
		if !options.Auth {
			tomlCopy.Auth = TOMLAuthConfig{}
		}
		tomlConfig = &tomlCopy
	}

	// Step 2: Load legacy configuration for sections missing in TOML.
	// The Web UI sets [compat].disable_legacy after migrating so that an empty
	// [[rss]] list in config.toml no longer falls back to rss.json.
	legacyConfig := &LegacyConfig{}
	disableLegacy := tomlConfig != nil && tomlConfig.Compat.DisableLegacy

	// Load legacy RSS if not in TOML and no CLI override
	if !disableLegacy && options.RSS && cliParams.RSSPath == "" && (tomlConfig == nil || len(tomlConfig.RSS) == 0) {
		rss, err := LoadLegacyRSS()
		if err != nil {
			return nil, nil, err
		}
		legacyConfig.RSS = rss
	}

	// Load legacy sites if not in TOML
	if !disableLegacy && options.Sites && (tomlConfig == nil || len(tomlConfig.Sites) == 0) {
		sites, err := LoadLegacySites()
		if err != nil {
			return nil, nil, err
		}
		legacyConfig.Sites = sites
	}

	// Load legacy cookies if not in TOML and no CLI override
	if !disableLegacy && options.Auth && cliParams.Cookies == "" && (tomlConfig == nil || (tomlConfig.Auth.Cookies == "" && tomlConfig.Auth.CookiesFile == "")) {
		cookies, cookiesPath, err := LoadLegacyCookies()
		if err != nil {
			return nil, nil, err
		}
		legacyConfig.Cookies = cookies
		legacyConfig.CookiesPath = cookiesPath
		source.CookiesPath = cookiesPath
	}

	// Track cookies path from TOML if configured
	if options.Auth && tomlConfig != nil && tomlConfig.Auth.CookiesFile != "" {
		cookiesPath := ResolveCookiesFile(tomlConfig.Auth.CookiesFile, tomlPath)
		data, err := os.ReadFile(cookiesPath)
		if err != nil {
			return nil, nil, err
		}
		tomlConfig.Auth.Cookies = string(data)
		source.CookiesPath = cookiesPath
	}

	// Step 3: Resolve configuration by priority
	cfg := Resolve(cliParams, tomlConfig, tomlPath, legacyConfig)

	// Step 4: Handle CLI RSS path override
	if options.RSS && cliParams.RSSPath != "" {
		// Load RSS from CLI-specified path
		rss, err := LoadRSSFromPath(cliParams.RSSPath)
		if err != nil {
			return nil, nil, err
		}
		cfg.RSS = rss
	}

	return cfg, source, nil
}

// LoadRSSFromPath loads RSS configuration from a specific file path
func LoadRSSFromPath(path string) (map[string][]RssConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string][]RssConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func normalizeHours(values []int) []int {
	seen := make(map[int]bool, len(values))
	result := make([]int, 0, len(values))
	for _, value := range values {
		if value < 0 || value > 23 || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func normalizeWashMode(value string) string {
	if value == "full" {
		return "full"
	}
	return "safe"
}
