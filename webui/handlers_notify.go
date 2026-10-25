package webui

import (
	"net/http"
)

func (s *Server) handleTestNotify(w http.ResponseWriter, r *http.Request) {
	var ignored struct{}
	if err := decodeJSON(w, r, &ignored); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if err := s.sendTestNotification(); err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
