package qbittorrent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Client is a small qBittorrent Web API v2 client.
type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

func New(baseURL, username, password string) (*Client, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("请填写 qBittorrent WebUI 地址")
	}
	jar, _ := cookiejar.New(nil)
	return &Client{
		baseURL: baseURL, username: username, password: password,
		http: &http.Client{Timeout: 30 * time.Second, Jar: jar},
	}, nil
}

func (c *Client) do(ctx context.Context, method, path string, form url.Values) ([]byte, int, error) {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Referer", c.baseURL)
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Accept", "application/json")
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return data, resp.StatusCode, nil
}

// Login authenticates and stores the qBittorrent session cookie.
func (c *Client) Login(ctx context.Context) error {
	form := url.Values{"username": {c.username}, "password": {c.password}}
	data, code, err := c.do(ctx, http.MethodPost, "/api/v2/auth/login", form)
	if err != nil {
		return fmt.Errorf("连接 qBittorrent 失败: %w", err)
	}
	text := strings.TrimSpace(string(data))
	if (code != http.StatusOK && code != http.StatusNoContent) || strings.HasPrefix(strings.ToLower(text), "<!doctype") || strings.HasPrefix(strings.ToLower(text), "<html") {
		if strings.Contains(strings.ToLower(text), "<html") || strings.Contains(strings.ToLower(text), "<!doctype") {
			return fmt.Errorf("qBittorrent 返回了网页而不是 API（HTTP %d）：请关闭 WebUI 的 Host header validation，或检查反向代理和 IP 封禁设置", code)
		}
		return fmt.Errorf("qBittorrent 登录失败（HTTP %d）: %s", code, text)
	}
	if text != "Ok." && text != "" {
		return fmt.Errorf("qBittorrent 登录失败: %s", text)
	}
	return nil
}

// Version returns the qBittorrent version.
func (c *Client) Version(ctx context.Context) (string, error) {
	if err := c.Login(ctx); err != nil {
		return "", err
	}
	data, code, err := c.do(ctx, http.MethodGet, "/api/v2/app/version", nil)
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(string(data))
	if (code != http.StatusOK && code != http.StatusNoContent) || strings.HasPrefix(strings.ToLower(text), "<!doctype") || strings.HasPrefix(strings.ToLower(text), "<html") {
		return "", fmt.Errorf("qBittorrent 返回了网页而不是 API（HTTP %d）：请检查 WebUI 访问权限或反向代理", code)
	}
	return text, nil
}

// Add submits a magnet/torrent URL.
func (c *Client) Add(ctx context.Context, torrentURL, savePath, category string) error {
	if strings.TrimSpace(torrentURL) == "" {
		return fmt.Errorf("链接为空")
	}
	if err := c.Login(ctx); err != nil {
		return err
	}
	form := url.Values{
		"urls":    {torrentURL},
		"autoTMM": {"false"},
		"paused":  {"false"},
	}
	if strings.TrimSpace(savePath) != "" {
		form.Set("savepath", strings.TrimSpace(savePath))
	}
	if strings.TrimSpace(category) != "" {
		form.Set("category", strings.TrimSpace(category))
	}
	data, code, err := c.do(ctx, http.MethodPost, "/api/v2/torrents/add", form)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("qBittorrent 返回 HTTP %d: %s", code, strings.TrimSpace(string(data)))
	}
	return nil
}
