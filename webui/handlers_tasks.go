package webui

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/deadblue/elevengo"
)

type taskView struct {
	InfoHash string  `json:"info_hash"`
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Status   int     `json:"status"`
	State    string  `json:"state"`
	Percent  float64 `json:"percent"`
	URL      string  `json:"url"`
	FileID   string  `json:"file_id"`
}

type tasksResponse struct {
	Tasks []taskView `json:"tasks"`
	Total int        `json:"total"`
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	tasks, err := ag.ListOfflineTasks(limit, nil)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "获取 115 任务失败: "+err.Error())
		return
	}
	resp := tasksResponse{Tasks: make([]taskView, 0, len(tasks)), Total: len(tasks)}
	for i := range tasks {
		t := &tasks[i]
		resp.Tasks = append(resp.Tasks, taskView{
			InfoHash: t.InfoHash,
			Name:     t.Name,
			Size:     t.Size,
			Status:   t.Status,
			State:    taskState(t),
			Percent:  t.Percent,
			URL:      t.Url,
			FileID:   t.FileId,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

func taskState(t *elevengo.OfflineTask) string {
	switch {
	case t.IsRunning():
		return "running"
	case t.IsDone():
		return "done"
	case t.IsFailed():
		return "failed"
	default:
		return "waiting"
	}
}

type clearTasksPayload struct {
	Type int `json:"type"`
}

func (s *Server) handleClearTasks(w http.ResponseWriter, r *http.Request) {
	var payload clearTasksPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if payload.Type < 1 || payload.Type > 6 {
		writeAPIError(w, http.StatusBadRequest, "清除类型必须是 1-6")
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	if err := ag.OfflineClear(payload.Type - 1); err != nil {
		writeAPIError(w, http.StatusBadGateway, "清除任务失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "type": payload.Type})
}

type magnetsPayload struct {
	Links    []string `json:"links"`
	Text     string   `json:"text"`
	Cid      string   `json:"cid"`
	SavePath string   `json:"savepath"`
}

func (s *Server) handleAddMagnets(w http.ResponseWriter, r *http.Request) {
	var payload magnetsPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	links := make([]string, 0, len(payload.Links))
	for _, link := range payload.Links {
		links = append(links, splitLinks(link)...)
	}
	links = append(links, splitLinks(payload.Text)...)

	seen := make(map[string]bool, len(links))
	unique := make([]string, 0, len(links))
	for _, link := range links {
		link = strings.TrimSpace(link)
		if link == "" || seen[link] {
			continue
		}
		if !validLink(link) {
			writeAPIError(w, http.StatusBadRequest, "不支持的链接类型: "+truncate(link, 80))
			return
		}
		seen[link] = true
		unique = append(unique, link)
	}
	if len(unique) == 0 {
		writeAPIError(w, http.StatusBadRequest, "没有可添加的链接")
		return
	}
	if len(unique) > 500 {
		writeAPIError(w, http.StatusBadRequest, "一次最多添加 500 个链接")
		return
	}

	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	added, err := ag.AddMagnetTasks(unique, strings.TrimSpace(payload.Cid), strings.TrimSpace(payload.SavePath))
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "添加离线任务失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"added": added, "total": len(unique)})
}

func splitLinks(text string) []string {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == '\r' || r == '\t'
	})
	return fields
}

func validLink(link string) bool {
	lower := strings.ToLower(link)
	switch {
	case strings.HasPrefix(lower, "magnet:"),
		strings.HasPrefix(lower, "ed2k://"),
		strings.HasPrefix(lower, "http://"),
		strings.HasPrefix(lower, "https://"),
		strings.HasPrefix(lower, "ftp://"):
		return true
	default:
		return false
	}
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

// handleLegacyAdd keeps the original POST /add API working.
func (s *Server) handleLegacyAdd(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Tasks    []string `json:"tasks"`
		Cid      string   `json:"cid"`
		SavePath string   `json:"savepath,omitempty"`
	}
	if err := decodeJSONLenient(r, &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(task.Tasks) == 0 {
		http.Error(w, "tasks is empty", http.StatusBadRequest)
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	ag.AddMagnetTask(task.Tasks, task.Cid, task.SavePath)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("message success"))
}
