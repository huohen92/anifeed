package webui

import (
	"net/http"
	"strconv"
)

func (s *Server) handleQBFallbackList(w http.ResponseWriter, r *http.Request) {
	ag, err := s.ensureAgent()
	if err != nil || ag.StoreInstance == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "数据库尚未就绪")
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	records, err := ag.StoreInstance.QBFallbackRecords(limit)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": records})
}
