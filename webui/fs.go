package webui

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/huohen92/anifeed/cloudfs"
)

type fsEntryView struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size,omitempty"`
}

type fsListResponse struct {
	Cid       string        `json:"cid"`
	Entries   []fsEntryView `json:"entries"`
	Total     int           `json:"total"`
	Truncated bool          `json:"truncated"`
}

func (s *Server) handleFSList(w http.ResponseWriter, r *http.Request) {
	cid := strings.TrimSpace(r.URL.Query().Get("cid"))
	if cid == "" {
		cid = "0"
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	entries, err := ag.ListDirectories(cid, 500)
	if err != nil {
		// 115 偶发返回空响应，elevengo 会 panic；等待后重试一次。
		time.Sleep(700 * time.Millisecond)
		entries, err = ag.ListDirectories(cid, 500)
	}
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "读取 115 目录失败: "+err.Error())
		return
	}
	dirs := make([]fsEntryView, 0, len(entries))
	for _, entry := range entries {
		dirs = append(dirs, fsEntryView{ID: entry.ID, Name: entry.Name, IsDir: true})
	}
	sort.SliceStable(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	resp := fsListResponse{Cid: cid, Entries: dirs, Total: len(dirs)}
	if len(dirs) >= 500 {
		resp.Truncated = true
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleFSStat(w http.ResponseWriter, r *http.Request) {
	cid := strings.TrimSpace(r.URL.Query().Get("cid"))
	if cid == "" || cid == "0" {
		writeJSON(w, http.StatusOK, fsEntryView{ID: "0", Name: "根目录", IsDir: true})
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	entry, err := ag.FileSystem().Stat(ctx, cid)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "读取 115 文件夹失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, fsEntryView{ID: entry.ID, Name: entry.Name, IsDir: entry.IsDir()})
}

func (s *Server) handleFSMkdir(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ParentCid string `json:"parent_cid"`
		Name      string `json:"name"`
	}
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" || payload.Name == "." || payload.Name == ".." ||
		strings.ContainsAny(payload.Name, `/\\`) {
		writeAPIError(w, http.StatusBadRequest, "文件夹名称不合法")
		return
	}
	if payload.ParentCid == "" {
		payload.ParentCid = "0"
	}
	ag, err := s.ensureAgent()
	if err != nil {
		writeAPIError(w, http.StatusConflict, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	entry, err := mkdirWithRecover(ctx, ag.FileSystem(), payload.ParentCid, payload.Name)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "创建 115 文件夹失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, fsEntryView{ID: entry.ID, Name: entry.Name, IsDir: true})
}

func mkdirWithRecover(ctx context.Context, fs cloudfs.Driver, parentCid, name string) (entry cloudfs.Entry, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			entry = cloudfs.Entry{}
			err = fmt.Errorf("115 创建目录接口异常: %v", recovered)
		}
	}()
	return fs.Mkdir(ctx, parentCid, name)
}
