package webui

import (
	"net/http"
	"strings"

	"github.com/huohen92/anifeed/store"
)

type previewDecision struct {
	Key     string `json:"key"`
	Allowed bool   `json:"allowed"`
}

type previewDecisionPayload struct {
	URL       string            `json:"url"`
	Decisions []previewDecision `json:"decisions"`
}

func (s *Server) handleSavePreviewDecisions(w http.ResponseWriter, r *http.Request) {
	var payload previewDecisionPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	payload.URL = strings.TrimSpace(payload.URL)
	if payload.URL == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少 RSS 地址")
		return
	}
	decisions := make(map[string]bool, len(payload.Decisions))
	for _, decision := range payload.Decisions {
		key := strings.TrimSpace(decision.Key)
		if key == "" {
			continue
		}
		decisions[key] = decision.Allowed
	}
	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer st.Close()
	if err := st.SavePreviewDecisions(payload.URL, decisions); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "saved": len(decisions)})
}
