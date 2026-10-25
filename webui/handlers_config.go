package webui

import (
	"net/http"
	"strings"
)

type configPayload struct {
	Config EditableConfig `json:"config"`
}

type rawConfigPayload struct {
	Content string `json:"content"`
}

func (s *Server) handleGetConfig(w http.ResponseWriter, _ *http.Request) {
	resp, err := s.loadEditable()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handlePutConfig(w http.ResponseWriter, r *http.Request) {
	var payload configPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "配置格式错误: "+err.Error())
		return
	}
	resp, err := s.saveEditable(&payload.Config)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetRawConfig(w http.ResponseWriter, _ *http.Request) {
	content, err := s.rawConfig()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"path":    s.ConfigPath(),
		"content": content,
	})
}

func (s *Server) handlePutRawConfig(w http.ResponseWriter, r *http.Request) {
	var payload rawConfigPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if strings.TrimSpace(payload.Content) == "" {
		writeAPIError(w, http.StatusBadRequest, "配置内容不能为空")
		return
	}
	if err := s.saveRawConfig(payload.Content); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := s.loadEditable()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetFeeds(w http.ResponseWriter, _ *http.Request) {
	resp, err := s.loadEditable()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"feeds": resp.Config.RSS,
		"meta":  resp.Meta,
	})
}
