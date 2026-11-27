package webui

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPasswordHashingAndLegacyCompatibility(t *testing.T) {
	const password = "correct horse battery staple"
	salt := newPasswordSalt()
	hashed := hashPassword(password, salt)
	if !verifyPassword(password, salt, hashed) {
		t.Fatal("new password hash did not verify")
	}
	if verifyPassword("wrong", salt, hashed) {
		t.Fatal("wrong password verified")
	}

	legacy := sha256.Sum256([]byte(salt + ":" + password))
	if !verifyPassword(password, salt, hex.EncodeToString(legacy[:])) {
		t.Fatal("legacy sha256 password hash did not verify")
	}
}

func TestRequestHTTPSHonorsForwardedProto(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "http://example.com/api/auth/login", nil)
	if requestIsHTTPS(r) {
		t.Fatal("plain request should not be HTTPS")
	}
	r.Header.Set("X-Forwarded-Proto", "https")
	if !requestIsHTTPS(r) {
		t.Fatal("X-Forwarded-Proto=https should mark request as HTTPS")
	}
	r.Header.Set("X-Forwarded-Proto", "http, https")
	if requestIsHTTPS(r) {
		t.Fatal("only the first forwarded proto should be trusted")
	}
}

func TestAuthSessionSurvivesServerRestart(t *testing.T) {
	s, _ := newTestServerWith(t, nil)
	rec, body := doRequest(t, s, http.MethodPost, "/api/auth/setup", map[string]any{"username": "admin", "password": "secret123"}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("setup = %d %s", rec.Code, body)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("missing session cookie")
	}
	s2 := New(s.cfg, s.source)
	rec, body = doRequest(t, s2, http.MethodGet, "/api/config", nil, map[string]string{"Cookie": cookies[0].Name + "=" + cookies[0].Value})
	if rec.Code != http.StatusOK {
		t.Fatalf("session did not survive restart: %d %s", rec.Code, body)
	}
}

func TestPasswordChangeInvalidatesSessions(t *testing.T) {
	s, _ := newTestServerWith(t, nil)
	rec, body := doRequest(t, s, http.MethodPost, "/api/auth/setup", map[string]any{"username": "admin", "password": "secret123"}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("setup = %d %s", rec.Code, body)
	}
	cookies := rec.Result().Cookies()
	resp, err := s.loadEditable()
	if err != nil {
		t.Fatal(err)
	}
	resp.Config.Server.AuthPassword = "newsecret123"
	if _, err := s.saveEditable(&resp.Config); err != nil {
		t.Fatal(err)
	}
	rec, _ = doRequest(t, s, http.MethodGet, "/api/config", nil, map[string]string{"Cookie": cookies[0].Name + "=" + cookies[0].Value})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("old session after password change = %d, want 401", rec.Code)
	}
}
