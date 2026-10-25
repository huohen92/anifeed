package request

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	urlPkg "net/url"
	"sync"
	"time"

	"github.com/zhifengle/anifeed/config"
)

var (
	ReqSiteConfig  = ReadNodeSiteConfig()
	proxyURL       = getProxyURL()
	ua             = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	clientMap      = make(map[string]*http.Client)
	clientMapMu    sync.RWMutex
	siteConfigMu   sync.RWMutex
	timeoutMu      sync.RWMutex
	requestTimeout = 20 * time.Second
)

// SetTimeout changes the HTTP timeout used for RSS/site requests. Existing
// clients keep their old timeout until ReloadConfig clears the client cache.
func SetTimeout(seconds int) {
	if seconds <= 0 {
		seconds = 20
	}
	timeoutMu.Lock()
	requestTimeout = time.Duration(seconds) * time.Second
	timeoutMu.Unlock()
}

func snapshotTimeout() time.Duration {
	timeoutMu.RLock()
	defer timeoutMu.RUnlock()
	return requestTimeout
}

// ReloadConfig re-reads site headers and proxy settings and drops cached HTTP
// clients. The Web UI calls it after config.toml is saved.
func ReloadConfig() error {
	cfg, _, err := config.LoadWithOptions(config.CLIParams{}, config.LoadOptions{Sites: true})
	if err != nil {
		return err
	}
	sites := make(NodeSiteConfig)
	for host, site := range cfg.Sites {
		sites[host] = SiteConfig{
			HttpsAgent: site.HttpsAgent,
			Headers:    site.Headers,
		}
	}

	// Close idle connections before swapping the configuration. The lock order
	// is clientMap -> siteConfig, matching getClientByReq.
	clientMapMu.Lock()
	for host, client := range clientMap {
		client.CloseIdleConnections()
		delete(clientMap, host)
	}
	clientMapMu.Unlock()

	siteConfigMu.Lock()
	ReqSiteConfig = sites
	proxyURL = cfg.Proxy.HTTP
	siteConfigMu.Unlock()
	SetTimeout(cfg.Subscription.RSSTimeoutSeconds)
	return nil
}

func snapshotSiteConfig(host string) (SiteConfig, bool) {
	siteConfigMu.RLock()
	defer siteConfigMu.RUnlock()
	cfg, ok := ReqSiteConfig[host]
	return cfg, ok
}

func snapshotProxyURL() string {
	siteConfigMu.RLock()
	defer siteConfigMu.RUnlock()
	return proxyURL
}

// getProxyURL loads the proxy URL from configuration
func getProxyURL() string {
	cfg, _, err := config.LoadWithOptions(config.CLIParams{}, config.LoadOptions{})
	if err != nil {
		// Fallback to default on error
		return "http://127.0.0.1:10809"
	}
	return cfg.Proxy.HTTP
}

type SiteConfig struct {
	HttpsAgent string            `json:"httpsAgent,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type NodeSiteConfig = map[string]SiteConfig

func ReadNodeSiteConfig() NodeSiteConfig {
	// Use config.Load() to get unified configuration
	cfg, _, err := config.LoadWithOptions(config.CLIParams{}, config.LoadOptions{Sites: true})
	if err != nil {
		// Fallback to legacy behavior on error
		sites, err := config.LoadLegacySites()
		if err != nil {
			return make(NodeSiteConfig)
		}
		legacyConfig := make(NodeSiteConfig)
		for host, site := range sites {
			legacyConfig[host] = SiteConfig{
				HttpsAgent: site.HttpsAgent,
				Headers:    site.Headers,
			}
		}
		return legacyConfig
	}

	// Convert config.SiteConfig to request.SiteConfig
	result := make(NodeSiteConfig)
	for host, site := range cfg.Sites {
		result[host] = SiteConfig{
			HttpsAgent: site.HttpsAgent,
			Headers:    site.Headers,
		}
	}

	return result
}

func getClientByReq(req *http.Request) *http.Client {
	host := req.URL.Host
	clientMapMu.RLock()
	client := clientMap[host]
	clientMapMu.RUnlock()
	if client != nil {
		return client
	}

	// Snapshot site settings before taking the client-map lock so a config
	// reload can never deadlock with client creation.
	curConfig, hasConfig := snapshotSiteConfig(host)
	proxy := snapshotProxyURL()

	clientMapMu.Lock()
	defer clientMapMu.Unlock()

	if clientMap[host] != nil {
		return clientMap[host]
	}

	var p func(*http.Request) (*urlPkg.URL, error)
	if hasConfig && curConfig.HttpsAgent != "" {
		u, _ := http.ProxyFromEnvironment(req)
		if u == nil {
			proxyURLParsed, _ := urlPkg.Parse(proxy)
			p = http.ProxyURL(proxyURLParsed)
		} else {
			p = http.ProxyFromEnvironment
		}
	}
	transport := NewIPv4Transport()
	transport.Proxy = p
	transport.DisableCompression = true
	transport.TLSHandshakeTimeout = 10 * time.Second
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Transport: transport,
		Timeout:   snapshotTimeout(),
		Jar:       jar,
	}
	clientMap[host] = client
	return client
}

// NewIPv4Transport returns an HTTP transport that only dials IPv4. This
// avoids containers that resolve AAAA records without an IPv6 route.
func NewIPv4Transport() *http.Transport {
	dialer := &net.Dialer{Timeout: 15 * time.Second, KeepAlive: 30 * time.Second}
	return &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp4", addr)
		},
	}
}

// NewIPv4Client builds a client that is restricted to IPv4.
func NewIPv4Client(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout, Transport: NewIPv4Transport()}
}

func Request(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)

	if curConfig, ok := snapshotSiteConfig(req.URL.Host); ok {
		if curConfig.Headers != nil {
			for k, v := range curConfig.Headers {
				req.Header.Set(k, v)
			}
		}
	}
	client := getClientByReq(req)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}

func GetByte(url string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	res, err := Request("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

func Get(url string, headers map[string]string) (string, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	body, err := GetByte(url, headers)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PostJson(url string, body []byte, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json; charset=UTF-8"
	res, err := Request("POST", url, bytes.NewBuffer(body), headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

func PostForm(url string, data urlPkg.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	res, err := Request("POST", url, bytes.NewBufferString(data.Encode()), headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}
