package webui

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/deadblue/elevengo/option"

	"github.com/zhifengle/anifeed/cloudfs"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/store"
)

type deleteHistoryPayload struct {
	Magnets     []string `json:"magnets"`
	DeleteFiles bool     `json:"delete_files"`
}

func (s *Server) handleDeletePreviewHistory(w http.ResponseWriter, r *http.Request) {
	var payload deleteHistoryPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	magnets := make([]string, 0, len(payload.Magnets))
	for _, magnet := range payload.Magnets {
		if value := strings.TrimSpace(magnet); value != "" {
			magnets = append(magnets, value)
		}
	}
	if len(magnets) == 0 {
		writeAPIError(w, http.StatusBadRequest, "没有选择条目")
		return
	}
	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer st.Close()
	filesDeleted := 0
	filesSkipped := 0
	if payload.DeleteFiles {
		storedHashes, hashErr := st.TaskHashesForMagnets(magnets)
		if hashErr != nil {
			writeAPIError(w, http.StatusInternalServerError, hashErr.Error())
			return
		}
		hashSet := make(map[string]bool)
		for _, hash := range storedHashes {
			if hash = strings.ToLower(strings.TrimSpace(hash)); hash != "" {
				hashSet[hash] = true
			}
		}
		for _, magnet := range magnets {
			if hash := magnetInfoHash(magnet); hash != "" {
				hashSet[hash] = true
			}
		}
		ag, agentErr := s.ensureAgent()
		if agentErr != nil {
			writeAPIError(w, http.StatusConflict, agentErr.Error())
			return
		}
		hashList := make([]string, 0, len(hashSet))
		for hash := range hashSet {
			hashList = append(hashList, hash)
		}
		finalEntries, _ := st.FinalEntriesForTaskHashes(hashList)
		fs := ag.FileSystem()
		for _, entries := range finalEntries {
			for _, entry := range entries {
				if entry.ID == "" {
					continue
				}
				if delErr := fs.Delete(context.Background(), entry.ID); delErr != nil {
					continue
				}
				filesDeleted++
			}
		}
		// Older postprocess rows predate final_entries. Reconstruct the
		// expected renamed name and search the final folder as a fallback.
		if cfg := s.configSnapshot(); cfg != nil {
			if history, historyErr := st.HistoryRowsForMagnets(magnets); historyErr == nil {
				searcher, canSearch := fs.(interface {
					Search(context.Context, string, string, cloudfs.SearchOptions) ([]cloudfs.Entry, error)
				})
				if canSearch {
					for _, row := range history {
						if len(finalEntries[row.TaskHash]) > 0 {
							continue
						}
						feed, ok := findFeedConfig(cfg, row.ConfigName, row.ConfigURL)
						if !ok {
							continue
						}
						base, buildErr := p115.BuildRenameName(row.ConfigName, row.Title, p115.RenameOptions{
							Enabled:  true,
							Template: cfg.Rename.Template,
							Season:   renameSeasonForFeedTitle(cfg, row.ConfigName, row.ConfigURL, row.Title),
							Offset:   cfg.Rename.Offset,
							TmdbID:   strings.TrimSpace(feed.TmdbID),
						})
						if buildErr != nil || base == "" {
							continue
						}
						if strings.TrimSpace(cfg.Offline.FinalCid) == "" {
							continue
						}
						matches, searchErr := searcher.Search(context.Background(), cfg.Offline.FinalCid, base, cloudfs.SearchOptions{})
						if searchErr != nil {
							continue
						}
						for _, entry := range matches {
							if entry.Name == base || strings.HasPrefix(entry.Name, base+".") {
								if delErr := fs.Delete(context.Background(), entry.ID); delErr == nil {
									filesDeleted++
								}
							}
						}
					}
				}
			}
		}
		for hash := range hashSet {
			// 115's offline delete API removes the task and, with this flag,
			// the completed file/folder created by that task.
			if delErr := ag.Agent.OfflineDelete([]string{hash}, option.OfflineDelete().DeleteDownloadedFiles()); delErr != nil {
				filesSkipped++
				continue
			}
			filesDeleted++
		}
	}
	deleted, err := st.DeleteDownloadedByMagnets(magnets)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "deleted": deleted, "files_deleted": filesDeleted, "files_skipped": filesSkipped})
}

func magnetInfoHash(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	xt := strings.ToLower(strings.TrimSpace(parsed.Query().Get("xt")))
	xt = strings.TrimPrefix(xt, "urn:btih:")
	xt = strings.TrimSpace(xt)
	if len(xt) == 40 || len(xt) == 32 {
		return xt
	}
	return ""
}
