package webui

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deadblue/elevengo"
	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/request"
	"github.com/huohen92/anifeed/version"
)

//go:embed assets
var assetsFS embed.FS

// ErrNoCookies is returned when 115 credentials have not been configured yet.
var ErrNoCookies = errors.New("尚未配置 115 Cookies，请先登录")

// ErrJobRunning is returned when another feed run is already in progress.
var ErrJobRunning = errors.New("已有运行中的任务，请稍后再试")

// ErrTokenRequired is returned when the optional access token is missing.
var ErrTokenRequired = errors.New("需要访问令牌")

// Server serves the anifeed Web UI and its JSON API.
type Server struct {
	mu           sync.RWMutex
	cfg          *config.Config
	source       *config.ConfigSource
	configPath   string
	configExists bool

	agentMu  sync.Mutex
	agent    *p115.Agent
	agentErr error
	agentAt  time.Time

	loginMu       sync.Mutex
	loginChecking bool
	loginChecked  bool
	loginOK       bool
	loginErr      string
	user          *elevengo.UserInfo
	loginAt       time.Time

	qrMu       sync.Mutex
	qrSessions map[string]*qrSession

	coverMu              sync.Mutex
	covers               map[string]coverEntry
	weekdays             map[string]weekdayCacheEntry
	notifyMetaCache      map[string]notifyMetaEntry
	episodeTitles        map[string]string
	hourlyMu             sync.Mutex
	lastHourly           string
	washRunMu            sync.Mutex
	washScheduleMu       sync.Mutex
	lastScheduledWash    time.Time
	washScheduledRunning bool

	logBuf        *LogBuffer
	prevLog       io.Writer
	authMu        sync.Mutex
	authSessions  map[string]authSession
	logMu         sync.Mutex
	jobs          *jobManager
	postProcessMu sync.Mutex
	startedAt     time.Time
	httpSrv       *http.Server
	versionText   string
}

// New creates a Web UI server from a loaded configuration.
func New(cfg *config.Config, source *config.ConfigSource) *Server {
	if source == nil {
		source = &config.ConfigSource{}
	}
	path := source.TOMLPath
	if path == "" {
		path = config.DefaultConfigPath()
	}
	_, err := os.Stat(path)
	return &Server{
		cfg:             cfg,
		source:          source,
		configPath:      path,
		configExists:    err == nil,
		qrSessions:      make(map[string]*qrSession),
		covers:          make(map[string]coverEntry),
		weekdays:        make(map[string]weekdayCacheEntry),
		notifyMetaCache: make(map[string]notifyMetaEntry),
		episodeTitles:   make(map[string]string),
		logBuf:          NewLogBuffer(800),
		jobs:            newJobManager(),
		authSessions:    make(map[string]authSession),
		startedAt:       time.Now(),
		versionText:     version.Version,
	}
}

// ConfigPath returns the config.toml path managed by the UI.
func (s *Server) ConfigPath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.configPath
}

// Handler builds the HTTP routing table.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/status", s.handleStatus)
	mux.HandleFunc("GET /api/config", s.handleGetConfig)
	mux.HandleFunc("PUT /api/config", s.handlePutConfig)
	mux.HandleFunc("GET /api/config/raw", s.handleGetRawConfig)
	mux.HandleFunc("PUT /api/config/raw", s.handlePutRawConfig)
	mux.HandleFunc("GET /api/feeds", s.handleGetFeeds)
	mux.HandleFunc("GET /api/feeds/cover", s.handleFeedCover)
	mux.HandleFunc("GET /api/feeds/weekday", s.handleFeedWeekday)
	mux.HandleFunc("GET /api/feeds/poster", s.handlePosterGet)
	mux.HandleFunc("POST /api/feeds/poster/upload", s.handlePosterUpload)
	mux.HandleFunc("POST /api/feeds/poster/refresh", s.handlePosterRefresh)
	mux.HandleFunc("DELETE /api/feeds/poster", s.handlePosterReset)
	mux.HandleFunc("GET /api/rename/preview", s.handleRenamePreview)
	mux.HandleFunc("POST /api/feeds/preview/decisions", s.handleSavePreviewDecisions)
	mux.HandleFunc("POST /api/feeds/preview/delete-history", s.handleDeletePreviewHistory)
	mux.HandleFunc("GET /api/tmdb/search", s.handleTMDBSearch)
	mux.HandleFunc("GET /api/tmdb/detail", s.handleTMDBDetail)
	mux.HandleFunc("GET /api/tmdb/auto", s.handleTMDBAuto)
	mux.HandleFunc("GET /api/rules", s.handleRules)
	mux.HandleFunc("GET /api/browse/list", s.handleBrowseList)
	mux.HandleFunc("GET /api/browse/groups", s.handleBrowseGroups)
	mux.HandleFunc("GET /api/fs/list", s.handleFSList)
	mux.HandleFunc("GET /api/fs/stat", s.handleFSStat)
	mux.HandleFunc("POST /api/lab/wash", s.handleWash)
	mux.HandleFunc("POST /api/fs/mkdir", s.handleFSMkdir)
	mux.HandleFunc("POST /api/feeds/preview", s.handlePreviewFeed)
	mux.HandleFunc("POST /api/feeds/run", s.handleRunFeed)
	mux.HandleFunc("POST /api/feeds/run-all", s.handleRunAllFeeds)
	mux.HandleFunc("GET /api/jobs", s.handleListJobs)
	mux.HandleFunc("GET /api/jobs/{id}", s.handleGetJob)
	mux.HandleFunc("GET /api/tasks", s.handleListTasks)
	mux.HandleFunc("POST /api/tasks/clear", s.handleClearTasks)
	mux.HandleFunc("POST /api/magnets", s.handleAddMagnets)
	mux.HandleFunc("GET /api/history", s.handleHistory)
	mux.HandleFunc("GET /api/logs", s.handleLogs)
	mux.HandleFunc("POST /api/notify/test", s.handleTestNotify)
	mux.HandleFunc("POST /api/postprocess/run", s.handleRunPostProcess)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("POST /api/auth/qrcode", s.handleStartQrcode)
	mux.HandleFunc("GET /api/auth/qrcode/{id}", s.handlePollQrcode)
	mux.HandleFunc("GET /api/auth/status", s.handleAuthStatus)
	mux.HandleFunc("POST /api/auth/login", s.handleAuthLogin)
	mux.HandleFunc("POST /api/auth/setup", s.handleAuthSetup)
	mux.HandleFunc("POST /api/session/logout", s.handleAuthLogout)
	mux.HandleFunc("POST /add", s.handleLegacyAdd)

	mux.Handle("/", staticHandler())
	return s.withSecurity(s.withAuth(mux))
}

// Start listens and serves until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	s.installLogCapture()

	s.mu.RLock()
	host := s.cfg.Server.Host
	port := s.cfg.Server.Port
	s.mu.RUnlock()
	if host == "" {
		host = "0.0.0.0"
	}
	if port == 0 {
		port = 8115
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.close()
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	s.httpSrv = &http.Server{
		Handler:           s.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      5 * time.Minute,
		IdleTimeout:       2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	log.Printf("AniFeed Web UI 已启动: http://%s", displayURL(host, port))
	if s.hasCookies() {
		s.refreshLogin(false)
	}
	go s.autoLoop(ctx)
	go s.startPostProcessLoop(ctx)
	go s.startHourlyNotifyLoop(ctx)
	go s.startWashScheduleLoop(ctx)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.httpSrv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Web UI 关闭失败: %v", err)
		}
	}()

	err = s.httpSrv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		s.close()
		return nil
	}
	s.close()
	return err
}

func (s *Server) close() {
	s.agentMu.Lock()
	if s.agent != nil {
		if err := s.agent.StoreInstance.Close(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
		s.agent = nil
	}
	s.agentMu.Unlock()

	s.logMu.Lock()
	if s.prevLog != nil {
		log.SetOutput(s.prevLog)
		s.prevLog = nil
	}
	s.logMu.Unlock()
}

func (s *Server) installLogCapture() {
	s.logMu.Lock()
	defer s.logMu.Unlock()
	if s.prevLog != nil {
		return
	}
	s.prevLog = log.Writer()
	log.SetOutput(io.MultiWriter(s.prevLog, s.logBuf))
}

func (s *Server) hasCookies() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg != nil && s.cfg.Auth.Cookies != ""
}

func (s *Server) configSnapshot() *config.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *Server) metaSnapshot() (string, *config.ConfigSource, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.configPath, s.source, s.configExists
}

// ensureAgent lazily creates the 115 client. Creation also validates cookies
// against 115, so failures are cached briefly to keep the UI responsive.
func (s *Server) ensureAgent() (*p115.Agent, error) {
	s.agentMu.Lock()
	defer s.agentMu.Unlock()

	if s.agent != nil {
		return s.agent, nil
	}
	if s.agentErr != nil && time.Since(s.agentAt) < 5*time.Second {
		return nil, s.agentErr
	}
	cfg := s.configSnapshot()
	if cfg == nil || cfg.Auth.Cookies == "" {
		err := ErrNoCookies
		s.agentErr = err
		s.agentAt = time.Now()
		return nil, err
	}
	ag, err := p115.NewAgent(cfg.Auth.Cookies)
	if err != nil {
		wrapped := fmt.Errorf("创建 115 客户端失败（Cookies 可能已失效）: %w", err)
		s.agentErr = wrapped
		s.agentAt = time.Now()
		return nil, wrapped
	}
	s.agent = ag
	s.agentErr = nil
	return ag, nil
}

func (s *Server) closeAgent() {
	s.agentMu.Lock()
	defer s.agentMu.Unlock()
	if s.agent != nil {
		if err := s.agent.StoreInstance.Close(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
		s.agent = nil
	}
	s.agentErr = nil
	s.agentAt = time.Time{}
}

// refreshLogin checks 115 credentials in the background. force ignores the
// 60-second cache. The UI polls /api/status for the result.
func (s *Server) refreshLogin(force bool) {
	s.loginMu.Lock()
	if s.loginChecking {
		s.loginMu.Unlock()
		return
	}
	if !force && s.loginChecked && time.Since(s.loginAt) < 60*time.Second {
		s.loginMu.Unlock()
		return
	}
	s.loginChecking = true
	s.loginMu.Unlock()

	go func() {
		var (
			info *elevengo.UserInfo
			err  error
		)
		defer func() {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("115 登录检测异常: %v", recovered)
				info = nil
				log.Printf("[error] %v", err)
			}
			s.loginMu.Lock()
			s.loginChecking = false
			s.loginChecked = true
			s.loginAt = time.Now()
			if err != nil {
				s.loginOK = false
				s.loginErr = err.Error()
				s.user = nil
			} else {
				s.loginOK = true
				s.loginErr = ""
				s.user = info
			}
			s.loginMu.Unlock()
		}()

		ag, agentErr := s.ensureAgent()
		if agentErr != nil {
			err = agentErr
			return
		}
		info = &elevengo.UserInfo{}
		if uerr := ag.Agent.UserGet(info); uerr != nil {
			info = nil
			err = fmt.Errorf("115 用户信息获取失败: %w", uerr)
		}
	}()
}

func (s *Server) loginSnapshot() (checking, checked, ok bool, errMsg string, user *elevengo.UserInfo) {
	s.loginMu.Lock()
	defer s.loginMu.Unlock()
	return s.loginChecking, s.loginChecked, s.loginOK, s.loginErr, s.user
}

// reloadFromDisk re-reads config.toml, refreshes HTTP/proxy caches, replaces
// the 115 client and re-checks the login state.
func (s *Server) reloadFromDisk() error {
	cfg, source, err := config.Load(config.CLIParams{})
	if err != nil {
		return err
	}
	p115.SetOption(p115.Option{
		DisableCache:       cfg.P115.DisableCache,
		ChunkDelay:         cfg.P115.ChunkDelay,
		ChunkSize:          cfg.P115.ChunkSize,
		CooldownMinMs:      cfg.P115.CooldownMinMs,
		CooldownMaxMs:      cfg.P115.CooldownMaxMs,
		Retry:              cfg.Offline.Retry,
		PostProcess:        cfg.Rename.Enabled || cfg.Offline.DualFolder,
		OnePerEpisode:      cfg.Offline.OnePerEpisode,
		MultiRssCoexist:    cfg.Lab.MultiRssCoexist,
		DatabasePath:       cfg.Database.Path,
		PreferredSubgroups: cfg.Lab.PreferredSubgroups,
	})
	if err := request.ReloadConfig(); err != nil {
		log.Printf("[warning] 刷新站点配置失败: %v", err)
	}

	s.closeAgent()

	s.mu.Lock()
	s.cfg = cfg
	s.source = source
	if source.TOMLPath != "" {
		s.configPath = source.TOMLPath
		_, statErr := os.Stat(source.TOMLPath)
		s.configExists = statErr == nil
	} else {
		s.configExists = false
	}
	s.mu.Unlock()

	s.loginMu.Lock()
	s.loginChecked = false
	s.loginErr = ""
	s.user = nil
	s.loginMu.Unlock()
	s.refreshLogin(true)
	return nil
}

func (s *Server) accessToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.cfg == nil {
		return ""
	}
	return s.cfg.Server.AccessToken
}

func (s *Server) withSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; img-src 'self' data: https: http:; style-src 'self' 'unsafe-inline'; script-src 'self'; connect-src 'self'")

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			if origin := r.Header.Get("Origin"); origin != "" && r.Header.Get("Sec-Fetch-Site") != "same-origin" {
				if !sameHost(origin, r.Host) {
					writeAPIError(w, http.StatusForbidden, "跨域请求被拒绝")
					return
				}
			}
		}

		path := r.URL.Path
		protected := (path == "/add" || strings.HasPrefix(path, "/api/")) && !strings.HasPrefix(path, "/api/auth/") && path != "/api/session/logout"
		if protected && path != "/api/health" && !(r.Method == http.MethodGet && path == "/api/feeds/poster") {
			if token := s.accessToken(); token != "" && !tokenMatches(token, r) {
				w.Header().Set("WWW-Authenticate", "X-Access-Token")
				writeAPIErrorCode(w, http.StatusUnauthorized, "需要访问令牌", "token_required")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func tokenMatches(token string, r *http.Request) bool {
	return subtleCompare(token, r.Header.Get("X-Access-Token"))
}

func sameHost(origin, host string) bool {
	u, err := urlParse(origin)
	if err != nil {
		return false
	}
	originHost := strings.Trim(u.Hostname(), "[]")
	requestHost := host
	if parsedHost, _, splitErr := net.SplitHostPort(host); splitErr == nil {
		requestHost = parsedHost
	}
	requestHost = strings.Trim(requestHost, "[]")
	return strings.EqualFold(originHost, requestHost)
}

func displayURL(host string, port int) string {
	shown := host
	if shown == "" || shown == "0.0.0.0" || shown == "::" || shown == "*" {
		shown = "127.0.0.1"
	}
	return net.JoinHostPort(shown, strconv.Itoa(port))
}

func staticHandler() http.Handler {
	sub, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			data, err := fs.ReadFile(sub, "index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(data)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

func (s *Server) databasePath() string {
	cfg := s.configSnapshot()
	if cfg == nil {
		return "db.sqlite"
	}
	return cfg.Database.Path
}

func (s *Server) cookiesPath() string {
	path, source, _ := s.metaSnapshot()
	dir := filepath.Dir(path)
	if source.TOMLPath != "" {
		if tc, err := config.LoadTOML(source.TOMLPath); err == nil && tc.Auth.CookiesFile != "" {
			return config.ResolveCookiesFile(tc.Auth.CookiesFile, source.TOMLPath)
		}
	}
	if source.CookiesPath != "" {
		return source.CookiesPath
	}
	return filepath.Join(dir, ".cookies")
}
