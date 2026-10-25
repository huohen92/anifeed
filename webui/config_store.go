package webui

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zhifengle/anifeed/config"
)

// EditableAuth is the [auth] section exposed to the UI.
type EditableAuth struct {
	CookiesFile string `json:"cookies_file"`
	Cookies     string `json:"cookies"`
}

// EditableServer is the [server] section exposed to the UI.
type EditableServer struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	AccessToken  string `json:"access_token"`
	AuthEnabled  bool   `json:"auth_enabled"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	passwordHash string
	passwordSalt string
}

// EditableDatabase is the [database] section exposed to the UI.
type EditableDatabase struct {
	Path string `json:"path"`
}

// EditableP115 is the [p115] section exposed to the UI.
type EditableP115 struct {
	DisableCache  bool `json:"disable_cache"`
	ChunkDelay    int  `json:"chunk_delay"`
	ChunkSize     int  `json:"chunk_size"`
	CooldownMinMs int  `json:"cooldown_min_ms"`
	CooldownMaxMs int  `json:"cooldown_max_ms"`
}

// EditableProxy is the [proxy] section exposed to the UI.
type EditableProxy struct {
	HTTP string `json:"http"`
}

// EditableSubscription holds source and network defaults.
type EditableSubscription struct {
	MikanHost         string `json:"mikan_host"`
	RSSTimeoutSeconds int    `json:"rss_timeout_seconds"`
}

// EditableOffline holds defaults applied to every offline task.
type EditableOffline struct {
	DefaultCid           string   `json:"default_cid"`
	DefaultSavePath      string   `json:"default_savepath"`
	GlobalExclude        []string `json:"global_exclude"`
	SkipX5               bool     `json:"skip_x5"`
	Retry                int      `json:"retry"`
	AutoClearCompleted   bool     `json:"auto_clear_completed"`
	AutoDisableCompleted bool     `json:"auto_disable_completed"`
	OnePerEpisode        bool     `json:"one_per_episode"`
	KeepOnlyVideo        bool     `json:"keep_only_video"`
	AutoLanguage         []string `json:"auto_language"`
	AutoSubtitle         []string `json:"auto_subtitle"`
	AutoContainer        []string `json:"auto_container"`
	AutoResolution       []string `json:"auto_resolution"`
	AutoRules            bool     `json:"auto_rules"`
	DualFolder           bool     `json:"dual_folder"`
	FinalCid             string   `json:"final_cid"`
	FinalMode            int      `json:"final_mode"`
	FinalCid2            string   `json:"final_cid2"`
	FinalPathTemplate    string   `json:"final_path_template"`
}

// EditableTMDB holds TMDB API settings.
type EditableTMDB struct {
	API    string `json:"api"`
	APIKey string `json:"api_key"`
	Image  string `json:"image"`
}

// EditableRename holds completed-download rename settings.
type EditableRename struct {
	Enabled  bool   `json:"enabled"`
	Template string `json:"template"`
	Season   int    `json:"season"`
	Offset   int    `json:"offset"`
}

// EditableAuto holds scheduled-run settings.
type EditableAuto struct {
	Enabled         bool `json:"enabled"`
	IntervalMinutes int  `json:"interval_minutes"`
	RunOnStart      bool `json:"run_on_start"`
}

// EditableNotify holds webhook notification settings.
type EditableNotify struct {
	Type                   string   `json:"type"`
	Webhook                string   `json:"webhook"`
	Content                []string `json:"content"`
	Events                 []string `json:"events"`
	OnSuccess              bool     `json:"on_success"`
	OnError                bool     `json:"on_error"`
	IdleEnabled            bool     `json:"idle_enabled"`
	IdleDays               int      `json:"idle_days"`
	Template               string   `json:"template"`
	DownloadTimeoutMinutes int      `json:"download_timeout_minutes"`
	SkipNoChange           bool     `json:"skip_no_change"`
	HourlyEnabled          bool     `json:"hourly_enabled"`
	HourlyHours            []int    `json:"hourly_hours"`
}

// EditableLab holds experimental settings.
type EditableLab struct {
	MultiSubgroup       bool     `json:"multi_subgroup"`
	MultiRssCoexist     bool     `json:"multi_rss_coexist"`
	WashEnabled         bool     `json:"wash_enabled"`
	WashMode            string   `json:"wash_mode"`
	WashScheduleEnabled bool     `json:"wash_schedule_enabled"`
	WashIntervalMinutes int      `json:"wash_interval_minutes"`
	WashFolders         []string `json:"wash_folders"`
	PreferredSubgroups  []string `json:"preferred_subgroups"`
	InactiveSubgroups   []string `json:"inactive_subgroups"`
	AutoBackupSubgroup  bool     `json:"auto_backup_subgroup"`
	AutoApplyRules      bool     `json:"auto_apply_rules"`
	HideHints           bool     `json:"hide_hints"`
}

// EditableSubgroup is one extra RSS source for a subscription.
type EditableSubgroup struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Season  int    `json:"season"`
	Offset  int    `json:"offset"`
	Filter  string `json:"filter"`
	Exclude string `json:"exclude"`
	Enabled bool   `json:"enabled"`
}

// EditableRSS is one [[rss]] entry.
type EditableRSS struct {
	Site            string             `json:"site"`
	Name            string             `json:"name"`
	SubgroupName    string             `json:"subgroup_name"`
	BgmURL          string             `json:"bgm_url"`
	URL             string             `json:"url"`
	Source          string             `json:"source"`
	Enabled         bool               `json:"enabled"`
	Cid             string             `json:"cid"`
	SavePath        string             `json:"savepath"`
	Filter          string             `json:"filter"`
	Exclude         string             `json:"exclude"`
	Expiration      uint               `json:"expiration"`
	RenameSeason    int                `json:"rename_season"`
	RenameOffset    int                `json:"rename_offset"`
	UpdateWeekday   int                `json:"update_weekday"`
	TotalEpisodes   int                `json:"total_episodes"`
	Completed       bool               `json:"completed"`
	FinalCid        string             `json:"final_cid"`
	FinalMode       int                `json:"final_mode"`
	FinalCustomMode int                `json:"final_custom_mode"`
	Subgroups       []EditableSubgroup `json:"subgroups"`
	TmdbID          string             `json:"tmdb_id"`
	TmdbType        string             `json:"tmdb_type"`
}

// EditableSite is one [sites."host"] entry.
type EditableSite struct {
	HTTPSAgent bool              `json:"https_agent"`
	Headers    map[string]string `json:"headers"`
}

// EditableConfig is the full configuration shape used by the UI forms.
type EditableConfig struct {
	Auth         EditableAuth            `json:"auth"`
	Server       EditableServer          `json:"server"`
	Database     EditableDatabase        `json:"database"`
	P115         EditableP115            `json:"p115"`
	Proxy        EditableProxy           `json:"proxy"`
	Subscription EditableSubscription    `json:"subscription"`
	Offline      EditableOffline         `json:"offline"`
	Rename       EditableRename          `json:"rename"`
	TMDB         EditableTMDB            `json:"tmdb"`
	Auto         EditableAuto            `json:"auto"`
	Notify       EditableNotify          `json:"notify"`
	Lab          EditableLab             `json:"lab"`
	RSS          []EditableRSS           `json:"rss"`
	Sites        map[string]EditableSite `json:"sites"`
}

// ConfigMeta describes where the configuration lives.
type ConfigMeta struct {
	Path           string `json:"path"`
	ConfigExists   bool   `json:"config_exists"`
	CookiesPath    string `json:"cookies_path"`
	DatabasePath   string `json:"database_path"`
	LegacyImported bool   `json:"legacy_imported"`
	TokenRequired  bool   `json:"token_required"`
}

// ConfigResponse is returned by GET/PUT /api/config.
type ConfigResponse struct {
	Config EditableConfig `json:"config"`
	Meta   ConfigMeta     `json:"meta"`
}

func (s *Server) loadEditable() (*ConfigResponse, error) {
	cfg := s.configSnapshot()
	if cfg == nil {
		return nil, fmt.Errorf("配置尚未加载")
	}
	path, source, exists := s.metaSnapshot()
	if source == nil {
		source = &config.ConfigSource{}
	}
	dir := filepath.Dir(path)

	var tcfg *config.TOMLConfig
	if source.TOMLPath != "" {
		parsed, err := config.LoadTOML(source.TOMLPath)
		if err != nil {
			return nil, err
		}
		tcfg = parsed
	}

	ec := &EditableConfig{
		Auth: EditableAuth{Cookies: cfg.Auth.Cookies},
		Server: EditableServer{
			Host:         cfg.Server.Host,
			Port:         cfg.Server.Port,
			AccessToken:  cfg.Server.AccessToken,
			AuthEnabled:  cfg.Server.AuthEnabled,
			AuthUsername: cfg.Server.AuthUsername,
			passwordHash: cfg.Server.AuthPasswordHash,
			passwordSalt: cfg.Server.AuthPasswordSalt,
		},
		Database: EditableDatabase{Path: "db.sqlite"},
		P115: EditableP115{
			DisableCache:  cfg.P115.DisableCache,
			ChunkDelay:    cfg.P115.ChunkDelay,
			ChunkSize:     cfg.P115.ChunkSize,
			CooldownMinMs: cfg.P115.CooldownMinMs,
			CooldownMaxMs: cfg.P115.CooldownMaxMs,
		},
		Proxy: EditableProxy{HTTP: cfg.Proxy.HTTP},
		Subscription: EditableSubscription{
			MikanHost:         cfg.Subscription.MikanHost,
			RSSTimeoutSeconds: cfg.Subscription.RSSTimeoutSeconds,
		},
		Offline: EditableOffline{
			DefaultCid:           cfg.Offline.DefaultCid,
			DefaultSavePath:      cfg.Offline.DefaultSavePath,
			GlobalExclude:        append([]string{}, cfg.Offline.GlobalExclude...),
			SkipX5:               cfg.Offline.SkipX5,
			Retry:                cfg.Offline.Retry,
			AutoClearCompleted:   cfg.Offline.AutoClearCompleted,
			AutoDisableCompleted: cfg.Offline.AutoDisableCompleted,
			OnePerEpisode:        cfg.Offline.OnePerEpisode,
			KeepOnlyVideo:        cfg.Offline.KeepOnlyVideo,
			AutoLanguage:         append([]string(nil), cfg.Offline.AutoLanguage...),
			AutoSubtitle:         append([]string(nil), cfg.Offline.AutoSubtitle...),
			AutoContainer:        append([]string(nil), cfg.Offline.AutoContainer...),
			AutoResolution:       append([]string(nil), cfg.Offline.AutoResolution...),
			AutoRules:            cfg.Offline.AutoRules,
			DualFolder:           cfg.Offline.DualFolder,
			FinalCid:             cfg.Offline.FinalCid,
			FinalMode:            cfg.Offline.FinalMode,
			FinalCid2:            cfg.Offline.FinalCid2,
			FinalPathTemplate:    cfg.Offline.FinalPathTemplate,
		},
		Rename: EditableRename{
			Enabled:  cfg.Rename.Enabled,
			Template: cfg.Rename.Template,
			Season:   cfg.Rename.Season,
			Offset:   cfg.Rename.Offset,
		},
		TMDB: EditableTMDB{
			API:    cfg.TMDB.API,
			APIKey: cfg.TMDB.APIKey,
			Image:  cfg.TMDB.Image,
		},
		Auto: EditableAuto{
			Enabled:         cfg.Auto.Enabled,
			IntervalMinutes: cfg.Auto.IntervalMinutes,
			RunOnStart:      cfg.Auto.RunOnStart,
		},
		Notify: EditableNotify{
			Type:    cfg.Notify.Type,
			Webhook: cfg.Notify.Webhook,
			Content: append([]string{}, cfg.Notify.Content...),
			Events: func() []string {
				if len(cfg.Notify.Events) == 0 {
					return config.DefaultNotifyEvents()
				}
				return append([]string{}, cfg.Notify.Events...)
			}(),
			OnSuccess:              cfg.Notify.OnSuccess,
			OnError:                cfg.Notify.OnError,
			IdleEnabled:            cfg.Notify.IdleEnabled,
			IdleDays:               cfg.Notify.IdleDays,
			Template:               cfg.Notify.Template,
			DownloadTimeoutMinutes: cfg.Notify.DownloadTimeoutMinutes,
			SkipNoChange:           cfg.Notify.SkipNoChange,
			HourlyEnabled:          cfg.Notify.HourlyEnabled,
			HourlyHours:            append([]int(nil), cfg.Notify.HourlyHours...),
		},
		Lab: EditableLab{
			MultiSubgroup:       cfg.Lab.MultiSubgroup,
			MultiRssCoexist:     cfg.Lab.MultiRssCoexist,
			WashEnabled:         cfg.Lab.WashEnabled,
			WashMode:            cfg.Lab.WashMode,
			WashScheduleEnabled: cfg.Lab.WashScheduleEnabled,
			WashIntervalMinutes: cfg.Lab.WashIntervalMinutes,
			WashFolders:         append([]string(nil), cfg.Lab.WashFolders...),
			PreferredSubgroups:  append([]string(nil), cfg.Lab.PreferredSubgroups...),
			InactiveSubgroups:   append([]string(nil), cfg.Lab.InactiveSubgroups...),
			AutoBackupSubgroup:  cfg.Lab.AutoBackupSubgroup,
			AutoApplyRules:      cfg.Lab.AutoApplyRules,
			HideHints:           cfg.Lab.HideHints,
		},
		RSS:   make([]EditableRSS, 0),
		Sites: make(map[string]EditableSite),
	}

	// cookies_file: keep the configured value when possible.
	cookiesFile := ".cookies"
	if tcfg != nil && tcfg.Auth.CookiesFile != "" {
		cookiesFile = tcfg.Auth.CookiesFile
	} else if source.CookiesPath != "" {
		cookiesFile = relPathOrAbs(source.CookiesPath, dir)
	}
	ec.Auth.CookiesFile = cookiesFile

	// database.path: prefer the literal configured value.
	if tcfg != nil && tcfg.Database.Path != "" {
		ec.Database.Path = tcfg.Database.Path
	} else if cfg.Database.Path != "" {
		ec.Database.Path = relPathOrAbs(cfg.Database.Path, dir)
	}

	// Prefer the literal config.toml order so adding/editing a feed does not
	// reorder the UI. Legacy-only configurations fall back to site sorting.
	if tcfg != nil && len(tcfg.RSS) > 0 {
		for _, rss := range tcfg.RSS {
			enabled := true
			if rss.Enabled != nil {
				enabled = *rss.Enabled
			}
			ec.RSS = append(ec.RSS, EditableRSS{
				Site:            rss.Site,
				Name:            rss.Name,
				SubgroupName:    rss.SubgroupName,
				BgmURL:          rss.BgmURL,
				URL:             rss.URL,
				Source:          config.NormalizeSource(rss.Source, rss.URL),
				Enabled:         enabled,
				Cid:             rss.Cid,
				SavePath:        rss.SavePath,
				Filter:          rss.Filter,
				Exclude:         rss.Exclude,
				Expiration:      rss.Expiration,
				RenameSeason:    rss.RenameSeason,
				RenameOffset:    rss.RenameOffset,
				UpdateWeekday:   rss.UpdateWeekday,
				TotalEpisodes:   rss.TotalEpisodes,
				Completed:       rss.Completed,
				FinalCid:        rss.FinalCid,
				FinalMode:       rss.FinalMode,
				FinalCustomMode: rss.FinalCustomMode,
				Subgroups:       editableSubgroupsFromTOML(rss.Subgroups),
				TmdbID:          rss.TmdbID,
				TmdbType:        rss.TmdbType,
			})
		}
	} else {
		sites := make([]string, 0, len(cfg.RSS))
		for site := range cfg.RSS {
			sites = append(sites, site)
		}
		sort.Strings(sites)
		for _, site := range sites {
			for _, rss := range cfg.RSS[site] {
				enabled := true
				if rss.Enabled != nil {
					enabled = *rss.Enabled
				}
				ec.RSS = append(ec.RSS, EditableRSS{
					Site:            site,
					Name:            rss.Name,
					SubgroupName:    rss.SubgroupName,
					BgmURL:          rss.BgmURL,
					URL:             rss.Url,
					Source:          config.NormalizeSource(rss.Source, rss.Url),
					Enabled:         enabled,
					Cid:             rss.Cid,
					SavePath:        rss.SavePath,
					Filter:          rss.Filter,
					Exclude:         rss.Exclude,
					Expiration:      rss.Expiration,
					RenameSeason:    rss.RenameSeason,
					RenameOffset:    rss.RenameOffset,
					UpdateWeekday:   rss.UpdateWeekday,
					TotalEpisodes:   rss.TotalEpisodes,
					Completed:       rss.Completed,
					FinalCid:        rss.FinalCid,
					FinalMode:       rss.FinalMode,
					FinalCustomMode: rss.FinalCustomMode,
					Subgroups:       editableSubgroups(rss.Subgroups),
					TmdbID:          rss.TmdbID,
					TmdbType:        rss.TmdbType,
				})
			}
		}
	}

	for host, site := range cfg.Sites {
		headers := make(map[string]string, len(site.Headers))
		for k, v := range site.Headers {
			headers[k] = v
		}
		ec.Sites[host] = EditableSite{HTTPSAgent: site.HttpsAgent != "", Headers: headers}
	}

	legacy := false
	if source.TOMLPath == "" {
		legacy = len(ec.RSS) > 0 || ec.Auth.Cookies != "" || len(ec.Sites) > 0
	} else if tcfg != nil && len(tcfg.RSS) == 0 && len(cfg.RSS) > 0 {
		legacy = true
	}

	cookiesPath := config.ResolveCookiesFile(cookiesFile, path)
	if cookiesFile == "" {
		cookiesPath = filepath.Join(dir, ".cookies")
	}
	dbPath := ec.Database.Path
	if !filepath.IsAbs(dbPath) {
		dbPath = filepath.Join(dir, dbPath)
	}

	meta := ConfigMeta{
		Path:           path,
		ConfigExists:   exists,
		CookiesPath:    cookiesPath,
		DatabasePath:   dbPath,
		LegacyImported: legacy,
		TokenRequired:  cfg.Server.AccessToken != "",
	}
	return &ConfigResponse{Config: *ec, Meta: meta}, nil
}

func (s *Server) saveEditable(ec *EditableConfig) (*ConfigResponse, error) {
	if ec == nil {
		return nil, fmt.Errorf("缺少配置内容")
	}
	if ec.Server.AuthEnabled && strings.TrimSpace(ec.Server.AuthPassword) == "" {
		s.mu.RLock()
		if s.cfg != nil {
			ec.Server.passwordHash = s.cfg.Server.AuthPasswordHash
			ec.Server.passwordSalt = s.cfg.Server.AuthPasswordSalt
		}
		s.mu.RUnlock()
	}
	path, _, _ := s.metaSnapshot()

	tc, cookiesAbs, err := buildTOMLConfig(ec, path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(cookiesAbsDir(cookiesAbs), 0o755); err != nil {
		return nil, fmt.Errorf("创建 Cookies 目录失败: %w", err)
	}
	if err := os.WriteFile(cookiesAbs, []byte(strings.TrimSpace(ec.Auth.Cookies)), 0o600); err != nil {
		return nil, fmt.Errorf("写入 Cookies 文件失败: %w", err)
	}

	if err := config.SaveTOML(path, tc); err != nil {
		return nil, err
	}
	if err := s.reloadFromDisk(); err != nil {
		return nil, fmt.Errorf("配置已保存，但重新加载失败: %w", err)
	}

	resp, err := s.loadEditable()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// buildTOMLConfig validates and converts the editable form into TOML config.
func buildTOMLConfig(ec *EditableConfig, configPath string) (*config.TOMLConfig, string, error) {
	dir := filepath.Dir(configPath)

	host := strings.TrimSpace(ec.Server.Host)
	if host == "" {
		host = "0.0.0.0"
	}
	if strings.Contains(host, "://") || strings.ContainsAny(host, " /\\") {
		return nil, "", fmt.Errorf("监听地址格式不正确: %q", ec.Server.Host)
	}
	port := ec.Server.Port
	if port == 0 {
		port = 8115
	}
	if port < 1 || port > 65535 {
		return nil, "", fmt.Errorf("端口必须在 1-65535 之间")
	}

	if ec.P115.ChunkDelay < 0 || ec.P115.ChunkSize < 1 || ec.P115.CooldownMinMs < 0 || ec.P115.CooldownMaxMs < 0 {
		return nil, "", fmt.Errorf("115 参数不合法：分块大小至少为 1，其余参数不能为负")
	}
	if ec.P115.CooldownMaxMs != 0 && ec.P115.CooldownMaxMs < ec.P115.CooldownMinMs {
		return nil, "", fmt.Errorf("冷却最大值不能小于最小值")
	}
	for _, section := range ec.Notify.Content {
		switch section {
		case "summary", "details", "errors", "time":
		default:
			return nil, "", fmt.Errorf("通知内容选项不合法: %s", section)
		}
	}
	if ec.Notify.Events == nil {
		ec.Notify.Events = config.DefaultNotifyEvents()
	}
	if len(ec.Notify.Events) == 0 {
		return nil, "", fmt.Errorf("请至少选择一个推送事件")
	}
	if ec.Notify.IdleEnabled && (ec.Notify.IdleDays < 1 || ec.Notify.IdleDays > 3650) {
		return nil, "", fmt.Errorf("摸鱼提醒天数必须在 1-3650 之间")
	}
	ec.Notify.HourlyHours = normalizeEditableHours(ec.Notify.HourlyHours)
	if ec.Offline.FinalMode != 2 {
		ec.Offline.FinalMode = 1
	}
	if ec.Offline.FinalMode == 2 && strings.TrimSpace(ec.Offline.FinalCid2) == "" {
		return nil, "", fmt.Errorf("保存文件夹 2 已启用，请选择文件夹")
	}
	if strings.TrimSpace(ec.Offline.FinalPathTemplate) == "" {
		ec.Offline.FinalPathTemplate = "${month}/${titleYearTmdb}/${animeSeason}"
	}
	if ec.Lab.WashIntervalMinutes <= 0 {
		ec.Lab.WashIntervalMinutes = 1440
	}
	if ec.Lab.WashIntervalMinutes < 1 || ec.Lab.WashIntervalMinutes > 100000 {
		return nil, "", fmt.Errorf("定时洗板间隔必须在 1-100000 分钟之间")
	}
	if ec.Notify.DownloadTimeoutMinutes == 0 {
		ec.Notify.DownloadTimeoutMinutes = 180
	}
	if ec.Notify.DownloadTimeoutMinutes < 1 || ec.Notify.DownloadTimeoutMinutes > 10080 {
		return nil, "", fmt.Errorf("下载超时时间必须在 1-10080 分钟之间")
	}
	switch strings.ToLower(strings.TrimSpace(ec.Notify.Type)) {
	case "", "auto", "generic", "wecom":
	default:
		return nil, "", fmt.Errorf("通知方式必须是 auto、generic 或 wecom")
	}
	if ec.Offline.DualFolder && strings.TrimSpace(ec.Offline.FinalCid) == "" {
		return nil, "", fmt.Errorf("启用「离线后转移到保存文件夹」时必须选择默认保存文件夹")
	}
	if ec.Rename.Enabled && strings.TrimSpace(ec.Rename.Template) == "" {
		return nil, "", fmt.Errorf("启用重命名时必须填写重命名模板")
	}
	authUsername := strings.TrimSpace(ec.Server.AuthUsername)
	authHash := strings.TrimSpace(ec.Server.passwordHash)
	authSalt := strings.TrimSpace(ec.Server.passwordSalt)
	if ec.Server.AuthEnabled {
		if authUsername == "" {
			return nil, "", fmt.Errorf("启用 Web UI 登录时必须填写账号")
		}
		if password := strings.TrimSpace(ec.Server.AuthPassword); password != "" {
			if len([]rune(password)) < 8 {
				return nil, "", fmt.Errorf("Web UI 登录密码至少需要 8 位")
			}
			authSalt = newPasswordSalt()
			authHash = hashPassword(password, authSalt)
		}
		if authHash == "" || authSalt == "" {
			return nil, "", fmt.Errorf("启用 Web UI 登录时必须设置密码")
		}
	} else {
		authUsername, authHash, authSalt = "", "", ""
	}
	if ec.Rename.Season < 0 || ec.Rename.Season > 1000 {
		return nil, "", fmt.Errorf("季数必须在 0-1000 之间")
	}
	if ec.Rename.Offset < -1000 || ec.Rename.Offset > 1000 {
		return nil, "", fmt.Errorf("集数偏移必须在 -1000 到 1000 之间")
	}
	for _, raw := range []struct {
		name  string
		value string
	}{{"TmdbApi", ec.TMDB.API}, {"TmdbImage", ec.TMDB.Image}} {
		value := strings.TrimSpace(raw.value)
		if value == "" {
			continue
		}
		u, err := url.Parse(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return nil, "", fmt.Errorf("%s 格式不正确", raw.name)
		}
	}

	proxy := strings.TrimSpace(ec.Proxy.HTTP)
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return nil, "", fmt.Errorf("代理地址格式不正确，应形如 http://127.0.0.1:10809")
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return nil, "", fmt.Errorf("代理地址仅支持 http/https")
		}
	}

	dbPath := strings.TrimSpace(ec.Database.Path)
	if dbPath == "" {
		dbPath = "db.sqlite"
	}

	cookiesFile := strings.TrimSpace(ec.Auth.CookiesFile)
	if cookiesFile == "" {
		cookiesFile = ".cookies"
	}
	cookiesAbs := config.ResolveCookiesFile(cookiesFile, configPath)
	if cookiesAbs == "" {
		cookiesAbs = filepath.Join(dir, ".cookies")
	}
	if info, err := os.Stat(cookiesAbs); err == nil && info.IsDir() {
		return nil, "", fmt.Errorf("Cookies 路径是目录: %s", cookiesAbs)
	}
	if samePath(cookiesAbs, configPath) {
		return nil, "", fmt.Errorf("Cookies 文件不能与 config.toml 相同")
	}
	if strings.TrimSpace(dbPath) != "" && samePath(config.ResolveDatabasePath(dbPath, configPath), configPath) {
		return nil, "", fmt.Errorf("数据库路径不能与 config.toml 相同")
	}

	tc := &config.TOMLConfig{
		Auth: config.TOMLAuthConfig{
			CookiesFile: relPathOrAbs(cookiesAbs, dir),
		},
		Server: config.TOMLServerConfig{
			Host:             host,
			Port:             port,
			AccessToken:      strings.TrimSpace(ec.Server.AccessToken),
			AuthEnabled:      ec.Server.AuthEnabled,
			AuthUsername:     authUsername,
			AuthPasswordHash: authHash,
			AuthPasswordSalt: authSalt,
		},
		Database: config.TOMLDatabaseConfig{Path: dbPath},
		P115: config.TOMLP115Config{
			DisableCache:  ec.P115.DisableCache,
			ChunkDelay:    ec.P115.ChunkDelay,
			ChunkSize:     ec.P115.ChunkSize,
			CooldownMinMs: ec.P115.CooldownMinMs,
			CooldownMaxMs: ec.P115.CooldownMaxMs,
		},
		Proxy: config.TOMLProxyConfig{HTTP: proxy},
		Subscription: config.TOMLSubscriptionConfig{
			MikanHost:         strings.TrimSpace(ec.Subscription.MikanHost),
			RSSTimeoutSeconds: ec.Subscription.RSSTimeoutSeconds,
		},
		Offline: config.TOMLOfflineConfig{
			DefaultCid:           strings.TrimSpace(ec.Offline.DefaultCid),
			DefaultSavePath:      strings.TrimSpace(ec.Offline.DefaultSavePath),
			GlobalExclude:        nonEmptyRules(ec.Offline.GlobalExclude),
			SkipX5:               ec.Offline.SkipX5,
			Retry:                ec.Offline.Retry,
			AutoClearCompleted:   ec.Offline.AutoClearCompleted,
			AutoDisableCompleted: ec.Offline.AutoDisableCompleted,
			OnePerEpisode:        &ec.Offline.OnePerEpisode,
			KeepOnlyVideo:        ec.Offline.KeepOnlyVideo,
			AutoLanguage:         ec.Offline.AutoLanguage,
			AutoSubtitle:         ec.Offline.AutoSubtitle,
			AutoContainer:        ec.Offline.AutoContainer,
			AutoResolution:       ec.Offline.AutoResolution,
			AutoRules:            ec.Offline.AutoRules,
			DualFolder:           ec.Offline.DualFolder,
			FinalCid:             strings.TrimSpace(ec.Offline.FinalCid),
			FinalMode: func() int {
				if ec.Offline.FinalMode == 2 {
					return 2
				}
				return 1
			}(),
			FinalCid2:         strings.TrimSpace(ec.Offline.FinalCid2),
			FinalPathTemplate: strings.TrimSpace(ec.Offline.FinalPathTemplate),
		},
		Rename: config.TOMLRenameConfig{
			Enabled:  ec.Rename.Enabled,
			Template: strings.TrimSpace(ec.Rename.Template),
			Season:   ec.Rename.Season,
			Offset:   ec.Rename.Offset,
		},
		TMDB: config.TOMLTMDBConfig{
			API:    strings.TrimSpace(ec.TMDB.API),
			APIKey: strings.TrimSpace(ec.TMDB.APIKey),
			Image:  strings.TrimSpace(ec.TMDB.Image),
		},
		Auto: config.TOMLAutoConfig{
			Enabled:         ec.Auto.Enabled,
			IntervalMinutes: ec.Auto.IntervalMinutes,
			RunOnStart:      ec.Auto.RunOnStart,
		},
		Notify: config.TOMLNotifyConfig{
			Type:                   strings.TrimSpace(ec.Notify.Type),
			Webhook:                strings.TrimSpace(ec.Notify.Webhook),
			Content:                notifyContentOrDefault(ec.Notify.Content),
			Events:                 normalizeEditableEvents(ec.Notify.Events),
			OnSuccess:              &ec.Notify.OnSuccess,
			OnError:                &ec.Notify.OnError,
			IdleEnabled:            ec.Notify.IdleEnabled,
			IdleDays:               ec.Notify.IdleDays,
			Template:               ec.Notify.Template,
			DownloadTimeoutMinutes: ec.Notify.DownloadTimeoutMinutes,
			SkipNoChange:           ec.Notify.SkipNoChange,
			HourlyEnabled:          ec.Notify.HourlyEnabled,
			HourlyHours:            append([]int(nil), ec.Notify.HourlyHours...),
		},
		Lab: config.TOMLLabConfig{
			MultiSubgroup:       ec.Lab.MultiSubgroup,
			MultiRssCoexist:     ec.Lab.MultiRssCoexist,
			WashEnabled:         ec.Lab.WashEnabled,
			WashMode:            ec.Lab.WashMode,
			WashScheduleEnabled: ec.Lab.WashScheduleEnabled,
			WashIntervalMinutes: ec.Lab.WashIntervalMinutes,
			WashFolders:         nonEmptyRules(ec.Lab.WashFolders),
			PreferredSubgroups:  nonEmptyRules(ec.Lab.PreferredSubgroups),
			InactiveSubgroups:   nonEmptyRules(ec.Lab.InactiveSubgroups),
			AutoBackupSubgroup:  ec.Lab.AutoBackupSubgroup,
			AutoApplyRules:      ec.Lab.AutoApplyRules,
			HideHints:           ec.Lab.HideHints,
		},
		RSS:    make([]config.TOMLRSSConfig, 0, len(ec.RSS)),
		Sites:  make(map[string]config.TOMLSiteConfig, len(ec.Sites)),
		Compat: config.TOMLCompatConfig{DisableLegacy: true},
	}

	for i, rss := range ec.RSS {
		name := strings.TrimSpace(rss.Name)
		rawURL := strings.TrimSpace(rss.URL)
		if name == "" {
			return nil, "", fmt.Errorf("第 %d 个订阅缺少名称", i+1)
		}
		if rawURL == "" {
			return nil, "", fmt.Errorf("订阅「%s」缺少 RSS 地址", name)
		}
		u, err := url.Parse(rawURL)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return nil, "", fmt.Errorf("订阅「%s」的 RSS 地址无效: %s", name, rawURL)
		}
		if rss.TotalEpisodes < 0 || rss.TotalEpisodes > 10000 {
			return nil, "", fmt.Errorf("订阅 %q 的总集数必须在 0-10000 之间", rss.Name)
		}
		if rss.UpdateWeekday < 0 || rss.UpdateWeekday > 7 {
			return nil, "", fmt.Errorf("订阅 %q 的更新星期必须在 0-7 之间", rss.Name)
		}
		if rss.RenameOffset < -1000 || rss.RenameOffset > 1000 {
			return nil, "", fmt.Errorf("订阅 %q 的集数偏移超出范围 [-1000, 1000]", rss.Name)
		}
		if rss.RenameSeason < 0 || rss.RenameSeason > 1000 {
			return nil, "", fmt.Errorf("订阅「%s」的重命名季数必须在 0-1000 之间", name)
		}
		site := strings.TrimSpace(rss.Site)
		if site == "" {
			site = u.Host
		}
		enabled := rss.Enabled
		tc.RSS = append(tc.RSS, config.TOMLRSSConfig{
			Site:          site,
			Name:          name,
			SubgroupName:  strings.TrimSpace(rss.SubgroupName),
			BgmURL:        strings.TrimSpace(rss.BgmURL),
			URL:           rawURL,
			Source:        config.NormalizeSource(rss.Source, rawURL),
			Enabled:       &enabled,
			Cid:           strings.TrimSpace(rss.Cid),
			SavePath:      strings.TrimSpace(rss.SavePath),
			Filter:        strings.TrimSpace(rss.Filter),
			Exclude:       strings.TrimSpace(rss.Exclude),
			Expiration:    rss.Expiration,
			RenameSeason:  rss.RenameSeason,
			RenameOffset:  rss.RenameOffset,
			UpdateWeekday: rss.UpdateWeekday,
			TotalEpisodes: rss.TotalEpisodes,
			Completed:     rss.Completed,
			FinalCid:      strings.TrimSpace(rss.FinalCid),
			FinalMode: func() int {
				if rss.FinalMode == 2 {
					return 2
				}
				return 1
			}(),
			FinalCustomMode: func() int {
				if rss.FinalCustomMode == 2 {
					return 2
				}
				return 1
			}(),
			Subgroups: tomlSubgroups(rss.Subgroups),
			TmdbID:    strings.TrimSpace(rss.TmdbID),
			TmdbType:  strings.TrimSpace(rss.TmdbType),
		})
	}

	for host, site := range ec.Sites {
		host = strings.TrimSpace(host)
		if host == "" {
			return nil, "", fmt.Errorf("站点域名不能为空")
		}
		headers := make(map[string]string, len(site.Headers))
		for k, v := range site.Headers {
			k = strings.TrimSpace(k)
			if k == "" {
				continue
			}
			headers[k] = v
		}
		tc.Sites[host] = config.TOMLSiteConfig{HTTPSAgent: site.HTTPSAgent, Headers: headers}
	}

	return tc, cookiesAbs, nil
}

func notifyContentOrDefault(content []string) []string {
	if len(content) == 0 {
		return config.DefaultNotifyContent()
	}
	return content
}

func cookiesAbsDir(path string) string {
	dir := filepath.Dir(path)
	if dir == "" {
		return "."
	}
	return dir
}

func relPathOrAbs(path, dir string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	base, err := filepath.Abs(dir)
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(base, abs)
	if err != nil {
		return path
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return abs
	}
	return rel
}

func (s *Server) rawConfig() (string, error) {
	path, _, _ := s.metaSnapshot()
	if data, err := os.ReadFile(path); err == nil {
		return string(data), nil
	}
	resp, err := s.loadEditable()
	if err != nil {
		return "", err
	}
	tc, _, err := buildTOMLConfig(&resp.Config, path)
	if err != nil {
		return "", err
	}
	// A brand-new config gets self-contained inline cookies so that saving the
	// raw text works even without the separate .cookies file.
	if resp.Config.Auth.Cookies != "" {
		tc.Auth.Cookies = resp.Config.Auth.Cookies
		tc.Auth.CookiesFile = ""
	}
	data, err := config.MarshalTOML(tc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) saveRawConfig(content string) error {
	path, _, _ := s.metaSnapshot()
	previous, previousErr := os.ReadFile(path)
	if err := config.SaveTOMLBytes(path, []byte(content)); err != nil {
		return err
	}
	if err := s.reloadFromDisk(); err != nil {
		// Roll back so a bad cookies_file / path cannot leave the service with
		// a config that no longer loads.
		if previousErr == nil {
			_ = os.WriteFile(path, previous, 0o600)
		} else {
			_ = os.Remove(path)
		}
		_ = s.reloadFromDisk()
		return fmt.Errorf("配置已回滚：%w", err)
	}
	return nil
}

func nonEmptyRules(rules []string) []string {
	out := make([]string, 0, len(rules))
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule != "" {
			out = append(out, rule)
		}
	}
	return out
}

func samePath(left, right string) bool {
	leftAbs, err := filepath.Abs(left)
	if err != nil {
		return false
	}
	rightAbs, err := filepath.Abs(right)
	if err != nil {
		return false
	}
	return leftAbs == rightAbs
}

func normalizeEditableHours(values []int) []int {
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

func normalizeEditableEvents(values []string) []string {
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

func editableSubgroups(values []config.RssSubgroup) []EditableSubgroup {
	result := make([]EditableSubgroup, 0, len(values))
	for _, value := range values {
		result = append(result, EditableSubgroup{Name: value.Name, URL: value.URL, Season: value.Season, Offset: value.Offset, Filter: value.Filter, Exclude: value.Exclude, Enabled: value.IsEnabled()})
	}
	return result
}

func tomlSubgroups(values []EditableSubgroup) []config.TOMLRSSSubgroup {
	result := make([]config.TOMLRSSSubgroup, 0, len(values))
	for _, value := range values {
		enabled := value.Enabled
		result = append(result, config.TOMLRSSSubgroup{Name: strings.TrimSpace(value.Name), URL: strings.TrimSpace(value.URL), Season: value.Season, Offset: value.Offset, Filter: strings.TrimSpace(value.Filter), Exclude: strings.TrimSpace(value.Exclude), Enabled: &enabled})
	}
	return result
}

func editableSubgroupsFromTOML(values []config.TOMLRSSSubgroup) []EditableSubgroup {
	result := make([]EditableSubgroup, 0, len(values))
	for _, value := range values {
		enabled := true
		if value.Enabled != nil {
			enabled = *value.Enabled
		}
		result = append(result, EditableSubgroup{Name: value.Name, URL: value.URL, Season: value.Season, Offset: value.Offset, Filter: value.Filter, Exclude: value.Exclude, Enabled: enabled})
	}
	return result
}
