package webui

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/deadblue/elevengo"
	"github.com/zhifengle/anifeed/p115"
)

type qrSession struct {
	login   *p115.WebQrcodeLogin
	created time.Time
}

func (s *Server) handleStartQrcode(w http.ResponseWriter, r *http.Request) {
	var ignored struct{}
	if err := decodeJSON(w, r, &ignored); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	login, err := p115.StartQrcodeLogin()
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "获取登录二维码失败: "+err.Error())
		return
	}
	id, err := randomID()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.qrMu.Lock()
	for key, session := range s.qrSessions {
		if time.Since(session.created) > 5*time.Minute {
			delete(s.qrSessions, key)
		}
	}
	s.qrSessions[id] = &qrSession{login: login, created: time.Now()}
	s.qrMu.Unlock()

	contentType := http.DetectContentType(login.Image())
	if contentType != "image/png" && contentType != "image/jpeg" {
		contentType = "image/png"
	}
	dataURI := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(login.Image())
	writeJSON(w, http.StatusOK, map[string]any{
		"id":         id,
		"image":      dataURI,
		"expires_in": 180,
	})
}

func (s *Server) handlePollQrcode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.qrMu.Lock()
	session, ok := s.qrSessions[id]
	s.qrMu.Unlock()
	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{"status": "expired", "message": "二维码已过期，请重新获取"})
		return
	}
	if session.login.Expired() {
		s.dropQRSession(id)
		writeJSON(w, http.StatusOK, map[string]any{"status": "expired", "message": "二维码已过期，请重新获取"})
		return
	}

	done, err := session.login.Poll()
	if err != nil {
		if errors.Is(err, elevengo.ErrQrcodeCancelled) {
			s.dropQRSession(id)
			writeJSON(w, http.StatusOK, map[string]any{"status": "cancelled", "message": "登录已取消"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "waiting", "message": "等待扫码确认"})
		return
	}
	if !done {
		writeJSON(w, http.StatusOK, map[string]any{"status": "waiting", "message": "等待扫码确认"})
		return
	}

	cookies := session.login.Cookies()
	s.dropQRSession(id)

	resp, err := s.loadEditable()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "登录成功，但读取配置失败: "+err.Error())
		return
	}
	resp.Config.Auth.Cookies = cookies
	if resp.Config.Auth.CookiesFile == "" {
		resp.Config.Auth.CookiesFile = ".cookies"
	}
	if _, err := s.saveEditable(&resp.Config); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "登录成功，但保存 Cookies 失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "success", "message": "登录成功"})
}

func (s *Server) dropQRSession(id string) {
	s.qrMu.Lock()
	delete(s.qrSessions, id)
	s.qrMu.Unlock()
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	var ignored struct{}
	if err := decodeJSON(w, r, &ignored); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	resp, err := s.loadEditable()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Config.Auth.Cookies = ""
	if resp.Config.Auth.CookiesFile == "" {
		resp.Config.Auth.CookiesFile = ".cookies"
	}
	if _, err := s.saveEditable(&resp.Config); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "退出登录失败: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func randomID() (string, error) {
	var buf [12]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}
