package webui

import (
	"net/http"
	"time"

	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/store"
)

type userView struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	IsVip bool   `json:"is_vip"`
}

type statusResponse struct {
	Version             string    `json:"version"`
	StartedAt           time.Time `json:"started_at"`
	ConfigPath          string    `json:"config_path"`
	ConfigExists        bool      `json:"config_exists"`
	DatabasePath        string    `json:"database_path"`
	CookiesPath         string    `json:"cookies_path"`
	HasCookies          bool      `json:"has_cookies"`
	LoginState          string    `json:"login_state"` // nocookies | checking | unknown | ok | error
	LoggedIn            bool      `json:"logged_in"`
	LoginError          string    `json:"login_error,omitempty"`
	User                *userView `json:"user,omitempty"`
	FeedCount           int       `json:"feed_count"`
	HistoryCount        int       `json:"history_count"`
	TodayCount          int       `json:"today_count"`
	DatabaseError       string    `json:"database_error,omitempty"`
	Host                string    `json:"host"`
	Port                int       `json:"port"`
	AccessTokenRequired bool      `json:"access_token_required"`
	RunningJob          bool      `json:"running_job"`
	AutoEnabled         bool      `json:"auto_enabled"`
	AutoIntervalMinutes int       `json:"auto_interval_minutes"`
	RenameEnabled       bool      `json:"rename_enabled"`
	DualFolder          bool      `json:"dual_folder"`
	PostProcessPending  int       `json:"postprocess_pending"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"version": s.versionText,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusInternalServerError, "配置尚未加载")
		return
	}

	// Refresh the 115 login state in the background; the response below is
	// always fast and never blocks on the 115 API.
	s.refreshLogin(false)
	checking, checked, loginOK, loginErr, user := s.loginSnapshot()

	_, _, configExists := s.metaSnapshot()
	hasCookies := cfg.Auth.Cookies != ""
	loginState := "unknown"
	switch {
	case !hasCookies:
		loginState = "nocookies"
	case checking:
		loginState = "checking"
	case !checked:
		loginState = "unknown"
	case loginOK:
		loginState = "ok"
	default:
		loginState = "error"
	}
	if loginState == "nocookies" {
		loginErr = ""
	}

	resp := statusResponse{
		Version:             s.versionText,
		StartedAt:           s.startedAt,
		ConfigPath:          s.ConfigPath(),
		ConfigExists:        configExists,
		DatabasePath:        s.databasePath(),
		CookiesPath:         s.cookiesPath(),
		HasCookies:          hasCookies,
		LoginState:          loginState,
		LoggedIn:            loginState == "ok",
		LoginError:          loginErr,
		FeedCount:           feedCount(cfg),
		Host:                cfg.Server.Host,
		Port:                cfg.Server.Port,
		AccessTokenRequired: cfg.Server.AccessToken != "",
		RunningJob:          s.jobs.isRunning(),
		AutoEnabled:         cfg.Auto.Enabled,
		AutoIntervalMinutes: cfg.Auto.IntervalMinutes,
		RenameEnabled:       cfg.Rename.Enabled,
		DualFolder:          cfg.Offline.DualFolder,
	}
	if user != nil && loginState == "ok" {
		resp.User = &userView{ID: user.Id, Name: user.Name, IsVip: user.IsVip}
	}

	// Local history counters are cheap; a failure must not break the page.
	if db, err := store.NewWithPath(s.databasePath()); err != nil {
		resp.DatabaseError = err.Error()
	} else {
		defer db.Close()
		if n, err := db.CountHistory(); err == nil {
			resp.HistoryCount = n
		} else {
			resp.DatabaseError = err.Error()
		}
		if n, err := db.CountPendingPostTasks(); err == nil {
			resp.PostProcessPending = n
		}
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if n, err := db.CountHistorySince(start); err == nil {
			resp.TodayCount = n
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func feedCount(cfg *config.Config) int {
	if cfg == nil {
		return 0
	}
	total := 0
	for _, feeds := range cfg.RSS {
		total += len(feeds)
	}
	return total
}
