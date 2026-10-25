package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// MarshalTOML renders a TOMLConfig as a deterministic, human-readable
// config.toml. Comments from a hand-written file are not preserved; use the
// raw editor when that matters.
func MarshalTOML(cfg *TOMLConfig) ([]byte, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	var b strings.Builder
	b.WriteString("# anifeed configuration\n")
	b.WriteString("# 由 anifeed Web UI 保存；手写注释请使用「原始配置」编辑器。\n\n")

	// [auth]
	b.WriteString("[auth]\n")
	if cfg.Auth.CookiesFile != "" {
		writeKV(&b, "cookies_file", quote(cfg.Auth.CookiesFile))
	}
	if cfg.Auth.Cookies != "" {
		writeKV(&b, "cookies", quote(cfg.Auth.Cookies))
	}
	b.WriteString("\n")

	// [server]
	b.WriteString("[server]\n")
	host := cfg.Server.Host
	if host == "" {
		host = "0.0.0.0"
	}
	writeKV(&b, "host", quote(host))
	port := cfg.Server.Port
	if port == 0 {
		port = 8115
	}
	writeKV(&b, "port", strconv.Itoa(port))
	if cfg.Server.AccessToken != "" {
		writeKV(&b, "access_token", quote(cfg.Server.AccessToken))
	}
	writeKV(&b, "auth_enabled", strconv.FormatBool(cfg.Server.AuthEnabled))
	if cfg.Server.AuthUsername != "" {
		writeKV(&b, "auth_username", quote(cfg.Server.AuthUsername))
	}
	if cfg.Server.AuthPasswordHash != "" {
		writeKV(&b, "auth_password_hash", quote(cfg.Server.AuthPasswordHash))
	}
	if cfg.Server.AuthPasswordSalt != "" {
		writeKV(&b, "auth_password_salt", quote(cfg.Server.AuthPasswordSalt))
	}
	b.WriteString("\n")

	// [database]
	b.WriteString("[database]\n")
	dbPath := cfg.Database.Path
	if dbPath == "" {
		dbPath = "db.sqlite"
	}
	writeKV(&b, "path", quote(dbPath))
	b.WriteString("\n")

	// [p115]
	b.WriteString("[p115]\n")
	writeKV(&b, "disable_cache", strconv.FormatBool(cfg.P115.DisableCache))
	writeKV(&b, "chunk_delay", strconv.Itoa(defaultInt(cfg.P115.ChunkDelay, 2)))
	writeKV(&b, "chunk_size", strconv.Itoa(defaultInt(cfg.P115.ChunkSize, 200)))
	writeKV(&b, "cooldown_min_ms", strconv.Itoa(defaultInt(cfg.P115.CooldownMinMs, 1000)))
	writeKV(&b, "cooldown_max_ms", strconv.Itoa(defaultInt(cfg.P115.CooldownMaxMs, 1100)))
	b.WriteString("\n")

	// [proxy]
	b.WriteString("[proxy]\n")
	proxy := cfg.Proxy.HTTP
	if proxy == "" {
		proxy = "http://127.0.0.1:10809"
	}
	writeKV(&b, "http", quote(proxy))
	b.WriteString("\n")

	// [subscription]
	b.WriteString("[subscription]\n")
	mikanHost := cfg.Subscription.MikanHost
	if mikanHost == "" {
		mikanHost = "https://mikanani.me"
	}
	writeKV(&b, "mikan_host", quote(mikanHost))
	timeout := cfg.Subscription.RSSTimeoutSeconds
	if timeout == 0 {
		timeout = 20
	}
	writeKV(&b, "rss_timeout_seconds", strconv.Itoa(timeout))
	b.WriteString("\n")

	// [offline]
	b.WriteString("[offline]\n")
	writeKV(&b, "default_cid", quote(cfg.Offline.DefaultCid))
	writeKV(&b, "default_savepath", quote(cfg.Offline.DefaultSavePath))
	b.WriteString("global_exclude = ")
	b.WriteString(stringArray(cfg.Offline.GlobalExclude))
	b.WriteString("\n")
	writeKV(&b, "skip_x5", strconv.FormatBool(cfg.Offline.SkipX5))
	writeKV(&b, "retry", strconv.Itoa(cfg.Offline.Retry))
	writeKV(&b, "auto_clear_completed", strconv.FormatBool(cfg.Offline.AutoClearCompleted))
	writeKV(&b, "auto_disable_completed", strconv.FormatBool(cfg.Offline.AutoDisableCompleted))
	writeKV(&b, "dual_folder", strconv.FormatBool(cfg.Offline.DualFolder))
	writeKV(&b, "final_cid", quote(cfg.Offline.FinalCid))
	finalMode := cfg.Offline.FinalMode
	if finalMode != 2 {
		finalMode = 1
	}
	writeKV(&b, "final_mode", strconv.Itoa(finalMode))
	writeKV(&b, "final_cid2", quote(cfg.Offline.FinalCid2))
	finalPath := cfg.Offline.FinalPathTemplate
	if strings.TrimSpace(finalPath) == "" {
		finalPath = "${month}/${titleYearTmdb}/${animeSeason}"
	}
	writeKV(&b, "final_path_template", quote(finalPath))
	writeKV(&b, "one_per_episode", strconv.FormatBool(boolOrDefault(cfg.Offline.OnePerEpisode, false)))
	writeKV(&b, "keep_only_video", strconv.FormatBool(cfg.Offline.KeepOnlyVideo))
	language, subtitle, container, resolution := DefaultAutoPriority()
	writeKV(&b, "auto_language", stringArray(orDefault(cfg.Offline.AutoLanguage, language)))
	writeKV(&b, "auto_subtitle", stringArray(orDefault(cfg.Offline.AutoSubtitle, subtitle)))
	writeKV(&b, "auto_container", stringArray(orDefault(cfg.Offline.AutoContainer, container)))
	writeKV(&b, "auto_resolution", stringArray(orDefault(cfg.Offline.AutoResolution, resolution)))
	writeKV(&b, "auto_rules", strconv.FormatBool(cfg.Offline.AutoRules))
	b.WriteString("\n")

	// [rename]
	b.WriteString("[rename]\n")
	writeKV(&b, "enabled", strconv.FormatBool(cfg.Rename.Enabled))
	template := cfg.Rename.Template
	if template == "" {
		template = DefaultRenameConfig().Template
	}
	writeKV(&b, "template", quote(template))
	season := cfg.Rename.Season
	if season == 0 {
		season = 1
	}
	writeKV(&b, "season", strconv.Itoa(season))
	writeKV(&b, "offset", strconv.Itoa(cfg.Rename.Offset))
	b.WriteString("\n")

	// [tmdb]
	b.WriteString("[tmdb]\n")
	writeKV(&b, "api", quote(cfg.TMDB.API))
	writeKV(&b, "api_key", quote(cfg.TMDB.APIKey))
	writeKV(&b, "image", quote(cfg.TMDB.Image))
	b.WriteString("\n")

	// [auto]
	b.WriteString("[auto]\n")
	writeKV(&b, "enabled", strconv.FormatBool(cfg.Auto.Enabled))
	interval := cfg.Auto.IntervalMinutes
	if interval == 0 {
		interval = 30
	}
	writeKV(&b, "interval_minutes", strconv.Itoa(interval))
	writeKV(&b, "run_on_start", strconv.FormatBool(cfg.Auto.RunOnStart))
	b.WriteString("\n")

	// [notify]
	b.WriteString("[notify]\n")
	notifyType := cfg.Notify.Type
	if notifyType == "" {
		notifyType = "auto"
	}
	writeKV(&b, "type", quote(notifyType))
	writeKV(&b, "webhook", quote(cfg.Notify.Webhook))
	events := cfg.Notify.Events
	if len(events) == 0 {
		events = DefaultNotifyEvents()
	}
	writeKV(&b, "events", stringArray(events))
	content := cfg.Notify.Content
	if len(content) == 0 {
		content = DefaultNotifyContent()
	}
	b.WriteString("content = ")
	b.WriteString(stringArray(content))
	b.WriteString("\n")
	writeKV(&b, "on_success", strconv.FormatBool(boolOrDefault(cfg.Notify.OnSuccess, true)))
	writeKV(&b, "on_error", strconv.FormatBool(boolOrDefault(cfg.Notify.OnError, true)))
	writeKV(&b, "idle_enabled", strconv.FormatBool(cfg.Notify.IdleEnabled))
	idleDays := cfg.Notify.IdleDays
	if idleDays <= 0 {
		idleDays = 7
	}
	writeKV(&b, "idle_days", strconv.Itoa(idleDays))
	template = cfg.Notify.Template
	if strings.TrimSpace(template) == "" {
		template = DefaultNotifyTemplate()
	}
	writeKV(&b, "template", quote(template))
	timeout = cfg.Notify.DownloadTimeoutMinutes
	if timeout <= 0 {
		timeout = 180
	}
	writeKV(&b, "download_timeout_minutes", strconv.Itoa(timeout))
	writeKV(&b, "skip_no_change", strconv.FormatBool(cfg.Notify.SkipNoChange))
	writeKV(&b, "hourly_enabled", strconv.FormatBool(cfg.Notify.HourlyEnabled))
	hours := cfg.Notify.HourlyHours
	b.WriteString("hourly_hours = [")
	for i, hour := range hours {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Itoa(hour))
	}
	b.WriteString("]\n")
	b.WriteString("\n")

	// [lab]
	b.WriteString("[lab]\n")
	writeKV(&b, "multi_subgroup", strconv.FormatBool(cfg.Lab.MultiSubgroup))
	writeKV(&b, "multi_rss_coexist", strconv.FormatBool(cfg.Lab.MultiRssCoexist))
	writeKV(&b, "wash_enabled", strconv.FormatBool(cfg.Lab.WashEnabled))
	writeKV(&b, "wash_mode", quote(normalizeWashMode(cfg.Lab.WashMode)))
	writeKV(&b, "hide_hints", strconv.FormatBool(cfg.Lab.HideHints))
	writeKV(&b, "wash_schedule_enabled", strconv.FormatBool(cfg.Lab.WashScheduleEnabled))
	washInterval := cfg.Lab.WashIntervalMinutes
	if washInterval <= 0 {
		washInterval = 1440
	}
	writeKV(&b, "wash_interval_minutes", strconv.Itoa(washInterval))
	b.WriteString("wash_folders = ")
	b.WriteString(stringArray(cfg.Lab.WashFolders))
	b.WriteString("\n")
	b.WriteString("preferred_subgroups = ")
	b.WriteString(stringArray(cfg.Lab.PreferredSubgroups))
	b.WriteString("\n")
	b.WriteString("inactive_subgroups = ")
	b.WriteString(stringArray(cfg.Lab.InactiveSubgroups))
	b.WriteString("\n")
	writeKV(&b, "auto_backup_subgroup", strconv.FormatBool(cfg.Lab.AutoBackupSubgroup))
	writeKV(&b, "auto_apply_rules", strconv.FormatBool(cfg.Lab.AutoApplyRules))
	b.WriteString("\n")

	// [[rss]]
	for _, rss := range cfg.RSS {
		b.WriteString("[[rss]]\n")
		writeKV(&b, "site", quote(rss.Site))
		writeKV(&b, "name", quote(rss.Name))
		if strings.TrimSpace(rss.SubgroupName) != "" {
			writeKV(&b, "subgroup_name", quote(rss.SubgroupName))
		}
		if strings.TrimSpace(rss.BgmURL) != "" {
			writeKV(&b, "bgm_url", quote(rss.BgmURL))
		}
		writeKV(&b, "url", quote(rss.URL))
		if rss.Source != "" {
			writeKV(&b, "source", quote(rss.Source))
		}
		if rss.Enabled != nil {
			writeKV(&b, "enabled", strconv.FormatBool(*rss.Enabled))
		}
		if rss.Cid != "" {
			writeKV(&b, "cid", quote(rss.Cid))
		}
		if rss.SavePath != "" {
			writeKV(&b, "savepath", quote(rss.SavePath))
		}
		if rss.Filter != "" {
			writeKV(&b, "filter", quote(rss.Filter))
		}
		if rss.Exclude != "" {
			writeKV(&b, "exclude", quote(rss.Exclude))
		}
		if rss.Expiration != 0 {
			writeKV(&b, "expiration", strconv.FormatUint(uint64(rss.Expiration), 10))
		}
		if rss.RenameSeason != 0 {
			writeKV(&b, "rename_season", strconv.Itoa(rss.RenameSeason))
		}
		if rss.RenameOffset != 0 {
			writeKV(&b, "rename_offset", strconv.Itoa(rss.RenameOffset))
		}
		if rss.UpdateWeekday != 0 {
			writeKV(&b, "update_weekday", strconv.Itoa(rss.UpdateWeekday))
		}
		if rss.TotalEpisodes != 0 {
			writeKV(&b, "total_episodes", strconv.Itoa(rss.TotalEpisodes))
		}
		if rss.Completed {
			writeKV(&b, "completed", "true")
		}
		if strings.TrimSpace(rss.FinalCid) != "" {
			writeKV(&b, "final_cid", quote(rss.FinalCid))
		}
		if rss.FinalMode == 2 {
			writeKV(&b, "final_mode", "2")
		}
		if rss.FinalCustomMode == 2 {
			writeKV(&b, "final_custom_mode", "2")
		}
		if len(rss.Subgroups) > 0 {
			b.WriteString("subgroups = [")
			for i, group := range rss.Subgroups {
				if i > 0 {
					b.WriteString(", ")
				}
				enabled := true
				if group.Enabled != nil {
					enabled = *group.Enabled
				}
				b.WriteString(fmt.Sprintf("{name=%s,url=%s,season=%d,offset=%d,filter=%s,exclude=%s,enabled=%t}", quote(group.Name), quote(group.URL), group.Season, group.Offset, quote(group.Filter), quote(group.Exclude), enabled))
			}
			b.WriteString("]\n")
		}
		if rss.TmdbID != "" {
			writeKV(&b, "tmdb_id", quote(rss.TmdbID))
		}
		if rss.TmdbType != "" {
			writeKV(&b, "tmdb_type", quote(rss.TmdbType))
		}
		b.WriteString("\n")
	}

	if cfg.Compat.DisableLegacy {
		b.WriteString("[compat]\n")
		writeKV(&b, "disable_legacy", "true")
		b.WriteString("\n")
	}

	// [sites."host"] sorted for stable output
	siteKeys := make([]string, 0, len(cfg.Sites))
	for host := range cfg.Sites {
		siteKeys = append(siteKeys, host)
	}
	sort.Strings(siteKeys)
	for _, host := range siteKeys {
		site := cfg.Sites[host]
		b.WriteString("[sites.")
		b.WriteString(quote(host))
		b.WriteString("]\n")
		writeKV(&b, "https_agent", strconv.FormatBool(site.HTTPSAgent))
		if len(site.Headers) > 0 {
			b.WriteString("\n[sites.")
			b.WriteString(quote(host))
			b.WriteString(".headers]\n")
			headerKeys := make([]string, 0, len(site.Headers))
			for key := range site.Headers {
				headerKeys = append(headerKeys, key)
			}
			sort.Strings(headerKeys)
			for _, key := range headerKeys {
				writeKV(&b, quote(key), quote(site.Headers[key]))
			}
		}
		b.WriteString("\n")
	}

	return []byte(strings.TrimRight(b.String(), "\n") + "\n"), nil
}

// SaveTOML validates, backs up and atomically writes cfg to path.
func SaveTOML(path string, cfg *TOMLConfig) error {
	data, err := MarshalTOML(cfg)
	if err != nil {
		return err
	}
	return SaveTOMLBytes(path, data)
}

// SaveTOMLBytes validates and atomically writes raw TOML content to path.
// The previous file is copied to <path>.bak before replacement.
func SaveTOMLBytes(path string, data []byte) error {
	if _, err := ParseTOML(data); err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("empty config path")
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	if old, err := os.ReadFile(path); err == nil {
		if bytes.Equal(old, data) {
			return nil
		}
		if err := os.WriteFile(path+".bak", old, 0o600); err != nil {
			return fmt.Errorf("failed to back up config: %w", err)
		}
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func writeKV(b *strings.Builder, key, value string) {
	b.WriteString(key)
	b.WriteString(" = ")
	b.WriteString(value)
	b.WriteString("\n")
}

func boolOrDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func orDefault(values, fallback []string) []string {
	if len(values) == 0 {
		return fallback
	}
	return values
}

func stringArray(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, quote(value))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func defaultInt(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}

// quote renders s as a TOML basic string. Unlike strconv.Quote it never emits
// \xNN escapes, which TOML does not support.
func quote(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		default:
			if r < 0x20 || r == 0x7f {
				fmt.Fprintf(&b, `\u%04X`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}
