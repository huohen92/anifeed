package webui

import (
	"net/http"
	"strconv"
)

func (s *Server) handleListJobs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"jobs": s.jobs.list()})
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.jobs.get(id)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "任务不存在")
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	limit := 200
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 800 {
			limit = n
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"lines": s.logBuf.Tail(limit)})
}
