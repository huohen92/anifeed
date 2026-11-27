package webui

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type themeDeviceSettings struct {
	Mode     string `json:"mode"`
	Selected string `json:"selected"`
}

type themeSettings struct {
	Desktop themeDeviceSettings `json:"desktop"`
	Mobile  themeDeviceSettings `json:"mobile"`
}

type themeImage struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	Modified int64  `json:"modified"`
	Builtin  bool   `json:"builtin,omitempty"`
}

type themeResponse struct {
	Desktop themeDeviceSettings `json:"desktop"`
	Mobile  themeDeviceSettings `json:"mobile"`
	Images  struct {
		Desktop []themeImage `json:"desktop"`
		Mobile  []themeImage `json:"mobile"`
	} `json:"images"`
}

var themeImageName = regexp.MustCompile(`^(desktop|mobile)-[0-9]+-[a-f0-9]{8}\.(jpg|png|webp|gif)$`)

func (s *Server) themeDir() string {
	return filepath.Join(filepath.Dir(s.ConfigPath()), "theme")
}

func (s *Server) themeSettingsPath() string {
	return filepath.Join(s.themeDir(), "theme.json")
}

func (s *Server) loadThemeSettings() themeSettings {
	out := themeSettings{
		Desktop: themeDeviceSettings{Mode: "single"},
		Mobile:  themeDeviceSettings{Mode: "single"},
	}
	data, err := os.ReadFile(s.themeSettingsPath())
	if err != nil {
		return out
	}
	_ = json.Unmarshal(data, &out)
	if out.Desktop.Mode != "none" {
		out.Desktop.Mode = "single"
	}
	if out.Mobile.Mode != "none" {
		out.Mobile.Mode = "single"
	}
	return out
}

func (s *Server) saveThemeSettings(settings themeSettings) error {
	if err := os.MkdirAll(s.themeDir(), 0o755); err != nil {
		return err
	}
	settings.Desktop.Mode = normalizeThemeMode(settings.Desktop.Mode)
	settings.Mobile.Mode = normalizeThemeMode(settings.Mobile.Mode)
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.themeSettingsPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.themeSettingsPath())
}

func normalizeThemeMode(mode string) string {
	if strings.ToLower(strings.TrimSpace(mode)) == "none" {
		return "none"
	}
	return "single"
}

func (s *Server) listThemeImages(device string) []themeImage {
	entries, err := os.ReadDir(s.themeDir())
	if err != nil {
		return []themeImage{}
	}
	out := make([]themeImage, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !themeImageName.MatchString(name) || !strings.HasPrefix(name, device+"-") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		out = append(out, themeImage{
			Name:     name,
			URL:      "/api/theme/image?name=" + name,
			Size:     info.Size(),
			Modified: info.ModTime().Unix(),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Modified > out[j].Modified })
	return out
}

var builtinThemeFiles = map[string][]string{
	"desktop": {"desktop-1.webp", "desktop-2.webp", "desktop-3.webp", "desktop-4.webp"},
	"mobile":  {"mobile-1.webp", "mobile-2.webp", "mobile-3.webp", "mobile-4.webp"},
}

func builtinThemeImages(device string) []themeImage {
	files := builtinThemeFiles[device]
	out := make([]themeImage, 0, len(files))
	for i, file := range files {
		out = append(out, themeImage{
			Name:    fmt.Sprintf("builtin-%s-%d", device, i+1),
			URL:     "/backgrounds/" + file,
			Builtin: true,
		})
	}
	return out
}

func (s *Server) themeResponse() themeResponse {
	settings := s.loadThemeSettings()
	out := themeResponse{Desktop: settings.Desktop, Mobile: settings.Mobile}
	out.Images.Desktop = append(builtinThemeImages("desktop"), s.listThemeImages("desktop")...)
	out.Images.Mobile = append(builtinThemeImages("mobile"), s.listThemeImages("mobile")...)
	return out
}

func (s *Server) handleTheme(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.themeResponse())
	case http.MethodPut:
		var payload themeSettings
		if err := decodeJSON(w, r, &payload); err != nil {
			writeAPIError(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		if err := s.saveThemeSettings(payload); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, s.themeResponse())
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

func (s *Server) handleThemeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}
	device := strings.TrimSpace(r.URL.Query().Get("device"))
	if device != "desktop" && device != "mobile" {
		writeAPIError(w, http.StatusBadRequest, "设备参数无效")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)
	if err := r.ParseMultipartForm(12 << 20); err != nil {
		writeAPIError(w, http.StatusBadRequest, "上传格式错误: "+err.Error())
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "缺少图片文件")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 10<<20+1))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "读取图片失败: "+err.Error())
		return
	}
	if len(data) > 10<<20 {
		writeAPIError(w, http.StatusBadRequest, "图片不能超过 10MB")
		return
	}
	ext := themeImageExt(data)
	if ext == "" {
		writeAPIError(w, http.StatusBadRequest, "只支持 JPG、PNG、WebP、GIF")
		return
	}
	if err := os.MkdirAll(s.themeDir(), 0o755); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	random := make([]byte, 4)
	_, _ = rand.Read(random)
	name := fmt.Sprintf("%s-%d-%s%s", device, time.Now().UnixNano(), hex.EncodeToString(random), ext)
	if err := os.WriteFile(filepath.Join(s.themeDir(), name), data, 0o644); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	settings := s.loadThemeSettings()
	deviceSettings := &settings.Desktop
	if device == "mobile" {
		deviceSettings = &settings.Mobile
	}
	if deviceSettings.Selected == "" {
		deviceSettings.Selected = name
		deviceSettings.Mode = "single"
		_ = s.saveThemeSettings(settings)
	}
	writeJSON(w, http.StatusOK, map[string]any{"name": name, "url": "/api/theme/image?name=" + name, "theme": s.themeResponse()})
}

func (s *Server) handleThemeImage(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if !themeImageName.MatchString(name) || filepath.Base(name) != name {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join(s.themeDir(), name)
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "读取图片失败", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

func (s *Server) handleThemeDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeAPIError(w, http.StatusMethodNotAllowed, "仅支持 DELETE")
		return
	}
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if !themeImageName.MatchString(name) || filepath.Base(name) != name {
		writeAPIError(w, http.StatusBadRequest, "图片名称无效")
		return
	}
	_ = os.Remove(filepath.Join(s.themeDir(), name))
	settings := s.loadThemeSettings()
	if settings.Desktop.Selected == name {
		settings.Desktop.Selected = ""
	}
	if settings.Mobile.Selected == name {
		settings.Mobile.Selected = ""
	}
	_ = s.saveThemeSettings(settings)
	writeJSON(w, http.StatusOK, s.themeResponse())
}

func themeImageExt(data []byte) string {
	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	}
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return ".webp"
	}
	return ""
}
