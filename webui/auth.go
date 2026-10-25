package webui

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"
)

const authCookieName = "anifeed_session"
const authSessionTTL = 7 * 24 * time.Hour
const passwordIterations = 210000

type authSession struct {
	Username  string
	ExpiresAt time.Time
}

func newPasswordSalt() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(buf)
}

// hashPassword derives a PBKDF2-HMAC-SHA256 key. The encoded form is
// self-describing so older sha256 hashes can still be verified and migrated.
func hashPassword(password, salt string) string {
	sum := pbkdf2SHA256([]byte(password), []byte(salt), passwordIterations, 32)
	return fmt.Sprintf("pbkdf2-sha256$%d$%s$%s", passwordIterations, salt, hex.EncodeToString(sum))
}

func verifyPassword(password, salt, stored string) bool {
	if strings.HasPrefix(stored, "pbkdf2-sha256$") {
		parts := strings.Split(stored, "$")
		if len(parts) != 4 {
			return false
		}
		iterations, err := strconv.Atoi(parts[1])
		if err != nil || iterations < 10000 || iterations > 2000000 {
			return false
		}
		sum := pbkdf2SHA256([]byte(password), []byte(parts[2]), iterations, 32)
		return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum)), []byte(parts[3])) == 1
	}
	// Backward compatibility with the early sha256(salt + ":" + password) format.
	sum := sha256.Sum256([]byte(salt + ":" + password))
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(stored)) == 1
}

func pbkdf2SHA256(password, salt []byte, iterations, keyLen int) []byte {
	hLen := sha256.Size
	numBlocks := (keyLen + hLen - 1) / hLen
	out := make([]byte, 0, numBlocks*hLen)
	var block [4]byte
	for i := 1; i <= numBlocks; i++ {
		binary.BigEndian.PutUint32(block[:], uint32(i))
		mac := hmac.New(sha256.New, password)
		mac.Write(salt)
		mac.Write(block[:])
		u := mac.Sum(nil)
		t := append([]byte(nil), u...)
		for j := 1; j < iterations; j++ {
			mac = hmac.New(sha256.New, password)
			mac.Write(u)
			u = mac.Sum(nil)
			for k := range t {
				t[k] ^= u[k]
			}
		}
		out = append(out, t...)
	}
	return out[:keyLen]
}

var loginLimiter = struct {
	sync.Mutex
	attempts map[string]loginAttempt
}{attempts: make(map[string]loginAttempt)}

type loginAttempt struct {
	Count   int
	ResetAt time.Time
}

func loginKey(r *http.Request, username string) string {
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i > 0 {
		host = host[:i]
	}
	return host + "|" + strings.ToLower(strings.TrimSpace(username))
}

func loginAllowed(key string) bool {
	loginLimiter.Lock()
	defer loginLimiter.Unlock()
	attempt, ok := loginLimiter.attempts[key]
	if !ok || time.Now().After(attempt.ResetAt) {
		delete(loginLimiter.attempts, key)
		return true
	}
	return attempt.Count < 5
}

func loginFailed(key string) {
	loginLimiter.Lock()
	defer loginLimiter.Unlock()
	now := time.Now()
	for existingKey, existing := range loginLimiter.attempts {
		if now.After(existing.ResetAt) {
			delete(loginLimiter.attempts, existingKey)
		}
	}
	attempt := loginLimiter.attempts[key]
	if attempt.ResetAt.IsZero() || time.Now().After(attempt.ResetAt) {
		attempt = loginAttempt{ResetAt: time.Now().Add(5 * time.Minute)}
	}
	attempt.Count++
	loginLimiter.attempts[key] = attempt
}

func loginSucceeded(key string) {
	loginLimiter.Lock()
	defer loginLimiter.Unlock()
	delete(loginLimiter.attempts, key)
}

func (s *Server) authConfig() (enabled bool, username, hash, salt string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.cfg == nil {
		return false, "", "", ""
	}
	return s.cfg.Server.AuthEnabled, s.cfg.Server.AuthUsername, s.cfg.Server.AuthPasswordHash, s.cfg.Server.AuthPasswordSalt
}

func (s *Server) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enabled, username, hash, _ := s.authConfig()
		path := r.URL.Path
		if strings.HasPrefix(path, "/api/auth/") || path == "/api/health" || path == "/api/session/logout" {
			next.ServeHTTP(w, r)
			return
		}
		if !enabled || username == "" || hash == "" {
			next.ServeHTTP(w, r)
			return
		}
		if s.authSessionValid(r) {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(path, "/api/") || path == "/add" {
			writeAPIErrorCode(w, http.StatusUnauthorized, "需要登录", "auth_required")
			return
		}
		// The SPA shell and static assets are public so the login screen can render.
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authSessionValid(r *http.Request) bool {
	cookie, err := r.Cookie(authCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}
	s.authMu.Lock()
	defer s.authMu.Unlock()
	session, ok := s.authSessions[cookie.Value]
	if !ok || time.Now().After(session.ExpiresAt) {
		delete(s.authSessions, cookie.Value)
		return false
	}
	return true
}

func (s *Server) createAuthSession(w http.ResponseWriter, r *http.Request, username string) {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	token := hex.EncodeToString(buf)
	s.authMu.Lock()
	for key, session := range s.authSessions {
		if time.Now().After(session.ExpiresAt) {
			delete(s.authSessions, key)
		}
	}
	s.authSessions[token] = authSession{Username: username, ExpiresAt: time.Now().Add(authSessionTTL)}
	s.authMu.Unlock()
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(authSessionTTL.Seconds()),
	})
}

func (s *Server) clearAuthSession(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(authCookieName); err == nil {
		s.authMu.Lock()
		delete(s.authSessions, cookie.Value)
		s.authMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: authCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
}

func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	enabled, username, hash, _ := s.authConfig()
	configured := enabled && username != "" && hash != ""
	writeJSON(w, http.StatusOK, map[string]any{
		"configured":    configured,
		"authenticated": configured && s.authSessionValid(r),
		"username":      username,
	})
}

func (s *Server) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	enabled, username, hash, salt := s.authConfig()
	if !enabled || username == "" || hash == "" {
		writeAPIError(w, http.StatusBadRequest, "尚未启用 Web UI 登录")
		return
	}
	key := loginKey(r, payload.Username)
	if !loginAllowed(key) {
		writeAPIError(w, http.StatusTooManyRequests, "登录失败次数过多，请 5 分钟后再试")
		return
	}
	if strings.TrimSpace(payload.Username) != username || !verifyPassword(payload.Password, salt, hash) {
		loginFailed(key)
		writeAPIError(w, http.StatusUnauthorized, "账号或密码错误")
		return
	}
	loginSucceeded(key)
	s.createAuthSession(w, r, username)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleAuthSetup(w http.ResponseWriter, r *http.Request) {
	enabled, username, hash, _ := s.authConfig()
	if enabled && username != "" && hash != "" {
		writeAPIError(w, http.StatusBadRequest, "Web UI 登录已经初始化")
		return
	}
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	payload.Username = strings.TrimSpace(payload.Username)
	if payload.Username == "" {
		writeAPIError(w, http.StatusBadRequest, "账号不能为空")
		return
	}
	if len([]rune(strings.TrimSpace(payload.Password))) < 8 {
		writeAPIError(w, http.StatusBadRequest, "密码至少需要 8 位")
		return
	}
	resp, err := s.loadEditable()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Config.Server.AuthEnabled = true
	resp.Config.Server.AuthUsername = payload.Username
	resp.Config.Server.AuthPassword = payload.Password
	if _, err := s.saveEditable(&resp.Config); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.createAuthSession(w, r, payload.Username)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	s.clearAuthSession(w, r)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
