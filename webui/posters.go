package webui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
)

var posterKeyPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

func posterKeyForURL(rawURL string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(rawURL)))
	return hex.EncodeToString(sum[:16])
}

func (s *Server) posterDir() string {
	return filepath.Join(filepath.Dir(s.ConfigPath()), "posters")
}

func (s *Server) localPosterURL(rawURL string) string {
	key := posterKeyForURL(rawURL)
	if !posterKeyPattern.MatchString(key) {
		return ""
	}
	for _, ext := range []string{".jpg", ".png", ".webp", ".gif"} {
		path := filepath.Join(s.posterDir(), key+ext)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return "/api/feeds/poster?key=" + key + "&v=" + fmt.Sprint(info.ModTime().Unix())
		}
	}
	return ""
}

func (s *Server) savePoster(rawURL string, data []byte) (string, error) {
	contentType := http.DetectContentType(data)
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	default:
		if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
			ext = ".webp"
		} else {
			return "", fmt.Errorf("不支持的图片格式: %s", contentType)
		}
	}
	key := posterKeyForURL(rawURL)
	dir := s.posterDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	for _, oldExt := range []string{".jpg", ".png", ".webp", ".gif"} {
		_ = os.Remove(filepath.Join(dir, key+oldExt))
	}
	path := filepath.Join(dir, key+ext)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	_ = os.Chmod(path, 0o644)
	return "/api/feeds/poster?key=" + key + "&v=" + fmt.Sprint(time.Now().Unix()), nil
}

func (s *Server) deletePoster(rawURL string) error {
	key := posterKeyForURL(rawURL)
	for _, ext := range []string{".jpg", ".png", ".webp", ".gif"} {
		_ = os.Remove(filepath.Join(s.posterDir(), key+ext))
	}
	return nil
}

func (s *Server) handlePosterGet(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if !posterKeyPattern.MatchString(key) {
		http.NotFound(w, r)
		return
	}
	for _, ext := range []string{".jpg", ".png", ".webp", ".gif"} {
		path := filepath.Join(s.posterDir(), key+ext)
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			http.Error(w, "读取海报失败", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
		return
	}
	http.NotFound(w, r)
}

func (s *Server) handlePosterUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8<<20)
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		writeAPIError(w, http.StatusBadRequest, "上传格式错误: "+err.Error())
		return
	}
	rawURL := strings.TrimSpace(r.FormValue("url"))
	if rawURL == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少订阅 URL")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "缺少图片文件")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 5<<20+1))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "读取图片失败: "+err.Error())
		return
	}
	if len(data) > 5<<20 {
		writeAPIError(w, http.StatusBadRequest, "图片不能超过 5MB")
		return
	}
	cover, err := s.savePoster(rawURL, data)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"cover": cover})
}

type posterRefreshPayload struct {
	URL      string `json:"url"`
	Source   string `json:"source"`
	BgmURL   string `json:"bgm_url"`
	TmdbID   string `json:"tmdb_id"`
	TmdbType string `json:"tmdb_type"`
	Mode     string `json:"mode"`
}

func (s *Server) handlePosterRefresh(w http.ResponseWriter, r *http.Request) {
	var payload posterRefreshPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	payload.URL = strings.TrimSpace(payload.URL)
	if payload.URL == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少订阅 URL")
		return
	}
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusInternalServerError, "配置尚未加载")
		return
	}
	var (
		info coverResponse
		err  error
	)
	mode := strings.ToLower(strings.TrimSpace(payload.Mode))
	bgmID := bangumiIDFromValue(payload.BgmURL)
	tmdbID := strings.TrimSpace(payload.TmdbID)
	source := config.NormalizeSource(payload.Source, payload.URL)
	sourceID := feedSourceID(source, payload.URL)
	switch mode {
	case "bgm":
		if bgmID == "" {
			writeAPIError(w, http.StatusBadRequest, "该订阅未绑定 BGM 链接")
			return
		}
		info, err = fetchBangumiCover(bgmID)
	case "tmdb":
		if tmdbID == "" || cfg.TMDB.APIKey == "" {
			writeAPIError(w, http.StatusBadRequest, "该订阅未关联 TMDB，或未配置 TMDB API Key")
			return
		}
		info, err = s.fetchTMDBPoster(tmdbID, payload.TmdbType)
	case "rss":
		if sourceID == "" {
			writeAPIError(w, http.StatusBadRequest, "该 RSS 来源没有可用于获取海报的站点 ID")
			return
		}
		info, err = s.fetchCover(source, sourceID, payload.URL)
	default:
		if bgmID != "" {
			info, err = fetchBangumiCover(bgmID)
		} else if tmdbID != "" && cfg.TMDB.APIKey != "" {
			info, err = s.fetchTMDBPoster(tmdbID, payload.TmdbType)
		} else if sourceID != "" {
			info, err = s.fetchCover(source, sourceID, payload.URL)
		} else {
			writeAPIError(w, http.StatusBadRequest, "该订阅没有可用的海报来源")
			return
		}
	}
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	if info.Cover == "" {
		writeAPIError(w, http.StatusNotFound, "没有找到海报")
		return
	}
	data, err := fetchBrowseBytes(info.Cover)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "下载海报失败: "+err.Error())
		return
	}
	cover, err := s.savePoster(payload.URL, data)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"cover": cover, "title": info.Title})
}

func (s *Server) handlePosterReset(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		URL string `json:"url"`
	}
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if strings.TrimSpace(payload.URL) == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少订阅 URL")
		return
	}
	if err := s.deletePoster(payload.URL); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.coverMu.Lock()
	s.covers = make(map[string]coverEntry)
	s.coverMu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func bangumiIDFromValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if u, err := url.Parse(value); err == nil && u.Hostname() != "" {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		for i, part := range parts {
			if part == "subject" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
		return ""
	}
	if _, err := strconv.Atoi(value); err == nil {
		return value
	}
	return ""
}
