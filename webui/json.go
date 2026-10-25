package webui

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const maxBodyBytes = 4 << 20

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type apiError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiError{Error: message})
}

func writeAPIErrorCode(w http.ResponseWriter, status int, message, code string) {
	writeJSON(w, status, apiError{Error: message, Code: code})
}

// decodeJSON decodes a JSON request body. Mutating API endpoints require the
// application/json content type, which also blocks simple form-based CSRF.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(ctype)), "application/json") {
		return errors.New("请求 Content-Type 必须是 application/json")
	}
	body := http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("请求体只能包含一个 JSON 对象")
	}
	return nil
}

func decodeJSONLenient(r *http.Request, dst any) error {
	return json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes)).Decode(dst)
}

func urlParse(raw string) (*url.URL, error) {
	return url.Parse(raw)
}

func subtleCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
