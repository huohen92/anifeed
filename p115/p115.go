package p115

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/deadblue/elevengo"
	"github.com/deadblue/elevengo/option"
	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/request"
	"github.com/huohen92/anifeed/rsssite"
	"github.com/huohen92/anifeed/store"
)

type agentOptions struct {
	disableCache       bool
	chunkDelay         int
	chunkSize          int
	cooldownMinMs      uint
	cooldownMaxMs      uint
	databasePath       string
	retry              int
	postProcess        bool
	onePerEpisode      bool
	multiRssCoexist    bool
	preferredSubgroups []string
}

var (
	agentOptionsMu sync.RWMutex
	currentOptions = agentOptions{
		chunkDelay:    2,
		chunkSize:     200,
		databasePath:  "db.sqlite",
		onePerEpisode: false,
	}
)

func optsSnapshot() agentOptions {
	agentOptionsMu.RLock()
	defer agentOptionsMu.RUnlock()
	return currentOptions
}

type Option struct {
	DisableCache       bool
	ChunkDelay         int
	ChunkSize          int
	CooldownMinMs      int
	CooldownMaxMs      int
	Retry              int
	PostProcess        bool
	OnePerEpisode      bool
	MultiRssCoexist    bool
	DatabasePath       string
	PreferredSubgroups []string
}

func SetOption(opt Option) {
	next := agentOptions{
		disableCache:       opt.DisableCache,
		chunkDelay:         2,
		chunkSize:          200,
		databasePath:       "db.sqlite",
		retry:              opt.Retry,
		postProcess:        opt.PostProcess,
		onePerEpisode:      opt.OnePerEpisode,
		multiRssCoexist:    opt.MultiRssCoexist,
		preferredSubgroups: append([]string(nil), opt.PreferredSubgroups...),
	}
	if opt.ChunkDelay > 0 {
		next.chunkDelay = opt.ChunkDelay
	}
	if opt.ChunkSize > 0 {
		next.chunkSize = opt.ChunkSize
	}
	if opt.CooldownMinMs > 0 {
		next.cooldownMinMs = uint(opt.CooldownMinMs)
	}
	if opt.CooldownMaxMs > 0 {
		next.cooldownMaxMs = uint(opt.CooldownMaxMs)
	}
	if next.cooldownMinMs > 0 && next.cooldownMaxMs == 0 {
		next.cooldownMaxMs = next.cooldownMinMs
	}
	if next.cooldownMaxMs > 0 && next.cooldownMaxMs < next.cooldownMinMs {
		next.cooldownMaxMs = next.cooldownMinMs
	}
	if opt.DatabasePath != "" {
		next.databasePath = opt.DatabasePath
	}
	if next.retry < 0 {
		next.retry = 0
	}
	agentOptionsMu.Lock()
	currentOptions = next
	agentOptionsMu.Unlock()
}

// offlineAdd submits URLs with the configured number of retries.
func (ag *Agent) offlineAdd(urls []string, cid, savepath string) ([]string, error) {
	opts := optsSnapshot()
	attempts := opts.retry + 1
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		hashes, err := ag.Agent.OfflineAddUrl(urls, &option.OfflineAddOptions{SaveDirId: cid, SavePath: savepath})
		if err == nil {
			return hashes, nil
		}
		lastErr = err
		if strings.Contains(strings.ToLower(err.Error()), "invalid ec data") {
			lastErr = fmt.Errorf("115 拒绝添加任务（invalid ec data）：任务或文件可能已经存在，请先删除云盘文件/下载记录后重试")
		}
		log.Printf("[warning] add offline task failed (attempt %d/%d): %s\n", attempt+1, attempts, lastErr)
		if attempt < attempts-1 {
			time.Sleep(time.Second * time.Duration(opts.chunkDelay))
		}
	}
	return nil, lastErr
}

type Agent struct {
	Agent         *elevengo.Agent
	StoreInstance *store.Store
}

func parseCookies(cookiesString string) map[string]string {
	cookies := make(map[string]string)

	// Split the cookies string into individual cookies
	cookiePairs := strings.Split(cookiesString, ";")

	// Parse each cookie into key-value pair
	for _, cookiePair := range cookiePairs {
		cookie := strings.TrimSpace(cookiePair)
		cookieParts := strings.SplitN(cookie, "=", 2)
		if len(cookieParts) == 2 {
			key := cookieParts[0]
			value := cookieParts[1]
			cookies[key] = value
		}
	}

	return cookies
}

func New() (*Agent, error) {
	cookies, err := LoadCookiesWithError()
	if err != nil {
		return nil, err
	}
	if cookies != "" {
		agent, err := NewAgent(cookies)
		// cookies is invalid
		if err != nil {
			return nil, err
		}
		return agent, nil
	}
	return nil, errors.New(".cookies is empty or not exist")
}

func NewAgentByQrcode() (*Agent, error) {
	cookies, err := LoadCookiesWithError()
	if err != nil {
		return nil, err
	}
	if cookies != "" {
		agent, err := NewAgent(cookies)
		// cookies is invalid
		if err != nil {
			return QrcodeLogin()
		}
		return agent, nil
	}
	return QrcodeLogin()
}
func NewAgentByConfig() (*Agent, error) {
	config := request.ReadNodeSiteConfig()
	if p115Config, ok := config["115.com"]; ok {
		cookies, ok := p115Config.Headers["cookie"]
		if !ok {
			cookies = p115Config.Headers["Cookie"]
		}
		if cookies == "" {
			return nil, errors.New("115 cookie is empty")
		}
		return NewAgent(cookies)
	}
	return nil, errors.New("no 115.com config in node-site-config.json")
}

func NewAgent(cookies string) (result *Agent, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = nil
			err = fmt.Errorf("初始化 115 客户端异常: %v", recovered)
			log.Printf("[error] %v", err)
		}
	}()
	agent := newElevengoAgent()
	if agent == nil {
		return nil, errors.New("初始化 115 客户端失败")
	}
	cookiesMap := parseCookies(cookies)
	err = agent.CredentialImport(&elevengo.Credential{
		UID: cookiesMap["UID"], CID: cookiesMap["CID"], SEID: cookiesMap["SEID"],
		KID: cookiesMap["KID"],
	})
	if err != nil {
		return nil, err
	}
	storeInstance, err := store.NewWithPath(optsSnapshot().databasePath)
	if err != nil {
		return nil, err
	}
	return &Agent{
		Agent:         agent,
		StoreInstance: storeInstance,
	}, nil
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

// FeedRunResult summarises one RSS feed run for the Web UI.
type FeedRunResult struct {
	Site            string   `json:"site"`
	Name            string   `json:"name"`
	URL             string   `json:"url"`
	Total           int      `json:"total"`
	Matched         int      `json:"matched"`
	Cached          int      `json:"cached"`
	Empty           int      `json:"empty"`
	Added           int      `json:"added"`
	AddedTitles     []string `json:"added_titles,omitempty"`
	Skipped         int      `json:"skipped"`
	Rejected        int      `json:"rejected"`
	Covered         int      `json:"covered"`
	CoveredEpisodes []int    `json:"covered_episodes,omitempty"`
	Error           string   `json:"error,omitempty"`
	Warning         string   `json:"warning,omitempty"`
}

// cleanupFailedTasks removes failed 115 offline tasks and their local rows
// before retrying a feed, so a failed item can be submitted again.
func (ag *Agent) cleanupFailedTasks(config *rsssite.RssConfig) {
	if ag == nil || ag.StoreInstance == nil || config == nil {
		return
	}
	hashes, err := ag.StoreInstance.FailedTaskHashesForFeed(config.Url, config.Name)
	if err != nil || len(hashes) == 0 {
		if err != nil {
			log.Printf("[warning] [%s] 读取失败离线任务失败: %v\n", config.Name, err)
		}
		return
	}
	for _, chunk := range chunkBy(hashes, 200) {
		if err := ag.Agent.OfflineDelete(chunk, option.OfflineDelete()); err != nil {
			log.Printf("[warning] [%s] 删除 115 失败任务失败: %v\n", config.Name, err)
		}
	}
	if err := ag.StoreInstance.DeleteTaskRecords(hashes); err != nil {
		log.Printf("[warning] [%s] 清理失败任务记录失败: %v\n", config.Name, err)
		return
	}
	log.Printf("[%s] 已清理 %d 个失败的离线任务，准备重新添加\n", config.Name, len(hashes))
}

func mergeCovered(covered map[int]bool, episodes []int) {
	for _, episode := range episodes {
		covered[episode] = true
	}
}

// RunRssConfigGroup runs the main feed and all enabled subgroup feeds.
func (ag *Agent) RunRssConfigGroup(config *rsssite.RssConfig) []FeedRunResult {
	if config == nil {
		return nil
	}
	mainResult := ag.runRssConfig(config, nil)
	results := []FeedRunResult{mainResult}
	covered := make(map[int]bool)
	if !optsSnapshot().multiRssCoexist {
		mergeCovered(covered, mainResult.CoveredEpisodes)
	}
	for _, subgroup := range config.Subgroups {
		if !subgroup.Enabled || strings.TrimSpace(subgroup.URL) == "" {
			continue
		}
		copy := *config
		copy.Url = subgroup.URL
		copy.Subgroups = nil
		copy.Filter = subgroup.Filter
		copy.Exclude = subgroup.Exclude
		if subgroup.Season > 0 {
			copy.RenameSeason = subgroup.Season
		}
		copy.RenameOffset = subgroup.Offset
		result := ag.runRssConfig(&copy, covered)
		if !optsSnapshot().multiRssCoexist {
			mergeCovered(covered, result.CoveredEpisodes)
		}
		results = append(results, result)
	}
	return results
}

// RunRssConfig fetches one feed and adds all new matching items to 115.
func (ag *Agent) RunRssConfig(config *rsssite.RssConfig) FeedRunResult {
	return ag.runRssConfig(config, nil)
}

func (ag *Agent) runRssConfig(config *rsssite.RssConfig, skipEpisodes map[int]bool) FeedRunResult {
	opts := optsSnapshot()
	result := FeedRunResult{
		Name: config.Name,
		URL:  config.Url,
	}
	if config.Name == "" {
		result.Name = config.Url
	}
	if site := siteHost(config.Url); site != "" {
		result.Site = site
	}

	ag.cleanupFailedTasks(config)

	total, magnetItems, err := rsssite.GetMagnetItemListWithStats(config)
	if err != nil {
		result.Error = err.Error()
		log.Printf("[error] [%s] %s\n", result.Name, err)
		return result
	}
	result.Total = total
	result.Matched = len(magnetItems)
	decisions, _ := ag.StoreInstance.PreviewDecisions(config.Url)

	if len(skipEpisodes) > 0 {
		kept := make([]rsssite.MagnetItem, 0, len(magnetItems))
		for _, item := range magnetItems {
			episode, epErr := ExtractEpisode(item.Title)
			if epErr == nil && episode > 0 && skipEpisodes[int(episode)] {
				if allowed, decided := decisions[ItemKey(item)]; !decided || !allowed {
					result.Covered++
					log.Printf("[%s] episode %g already covered by earlier RSS, skip: %s", config.Name, episode, item.Title)
					continue
				}
			}
			kept = append(kept, item)
		}
		magnetItems = kept
	}

	selected := SelectEpisodes(magnetItems, opts.onePerEpisode)
	selectedItems, skipped := selectWithDecisions(magnetItems, selected, decisions)
	result.Skipped += skipped
	if result.Skipped > 0 {
		log.Printf("[%s] skipped %d duplicate episode version(s)\n", config.Name, result.Skipped)
	}
	filteredItems := make([]rsssite.MagnetItem, 0, len(selectedItems))
	coveredEpisodes := make([]int, 0)
	addCovered := func(item rsssite.MagnetItem) {
		episode, err := ExtractEpisode(item.Title)
		if err != nil || episode <= 0 {
			return
		}
		key := int(episode)
		for _, existing := range coveredEpisodes {
			if existing == key {
				return
			}
		}
		coveredEpisodes = append(coveredEpisodes, key)
	}
	for _, item := range selectedItems {
		if item.Magnet == "" {
			result.Empty++
			continue
		}
		if !opts.disableCache && ag.StoreInstance.HasItem(item.Magnet) {
			result.Cached++
			addCovered(item)
			continue
		}
		filteredItems = append(filteredItems, item)
	}
	if result.Empty != 0 {
		log.Printf("[warning] [%s] has %d empty task\n", config.Name, result.Empty)
	}
	result.CoveredEpisodes = coveredEpisodes
	if len(filteredItems) == 0 {
		log.Printf("[%s] has 0 new task\n", result.Name)
		return result
	}

	chunks := chunkBy(filteredItems, opts.chunkSize)
	for i, items := range chunks {
		urls := make([]string, 0, len(items))
		for _, item := range items {
			urls = append(urls, item.Magnet)
		}
		hashes, err := ag.offlineAdd(urls, config.Cid, config.SavePath)
		if err != nil {
			result.Error = fmt.Sprintf("add offline task: %v", err)
			log.Printf("[error] [%s] add offline task: %s\n", result.Name, err)
			return result
		}
		successfulItems, successfulHashes := successfulMagnetItems(items, hashes, func(item rsssite.MagnetItem) {
			log.Printf("[warning] [%s] offline task was not accepted: %s\n", config.Name, item.Title)
		})
		result.Added += len(successfulItems)
		for _, item := range successfulItems {
			result.AddedTitles = append(result.AddedTitles, item.Title)
			addCovered(item)
		}
		result.Rejected += len(items) - len(successfulItems)
		log.Printf("[%s] [%s] add %d/%d tasks\n", result.Name, config.Url, len(successfulItems), len(urls))
		var saveErr error
		if opts.postProcess {
			saveErr = ag.StoreInstance.SaveMagnetItemsWithPostprocess(successfulItems, config.Name, config.Url, successfulHashes, config.Cid, config.RenameSeason, config.RenameOffset)
		} else {
			saveErr = ag.StoreInstance.SaveMagnetItems(successfulItems, config.Name, config.Url, successfulHashes)
		}
		if saveErr != nil {
			log.Printf("Save magnet items error: %s\n", saveErr)
			result.Warning = fmt.Sprintf("save history: %v", saveErr)
		}
		if i < len(chunks)-1 {
			time.Sleep(time.Second * time.Duration(opts.chunkDelay))
		}
	}
	return result
}

func siteHost(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Host
}

func successfulMagnetItems(items []rsssite.MagnetItem, hashes []string, rejected func(rsssite.MagnetItem)) ([]rsssite.MagnetItem, []string) {
	successfulItems := make([]rsssite.MagnetItem, 0, len(items))
	successfulHashes := make([]string, 0, len(items))
	for i, item := range items {
		if i >= len(hashes) || hashes[i] == "" {
			if rejected != nil {
				rejected(item)
			}
			continue
		}
		successfulItems = append(successfulItems, item)
		successfulHashes = append(successfulHashes, hashes[i])
	}
	return successfulItems, successfulHashes
}

func (ag *Agent) AddRssUrlTask(url string) {
	result, err := ag.RunRssURL(url)
	if err != nil {
		log.Printf("[error] %s\n", err)
		return
	}
	if result.Error != "" {
		log.Printf("[error] [%s] %s\n", result.Name, result.Error)
	}
}

// RunRssURL runs a single feed URL using the matching config.toml entry (or
// default options when no entry exists).
func (ag *Agent) RunRssURL(url string) (FeedRunResult, error) {
	config := rsssite.GetRssConfigByURL(url)
	if config == nil {
		return FeedRunResult{}, fmt.Errorf("invalid rss url: %s", url)
	}
	result := ag.RunRssConfig(config)
	return result, nil
}

// RunAllRssTasks runs every configured feed in a deterministic order.
func (ag *Agent) RunAllRssTasks() ([]FeedRunResult, error) {
	rssDict := rsssite.ReadRssConfigDict()
	if rssDict == nil {
		pwd := os.Getenv("PWD")
		return nil, fmt.Errorf("rss config not found: %s", pwd)
	}

	sites := make([]string, 0, len(*rssDict))
	for site := range *rssDict {
		sites = append(sites, site)
	}
	sort.Strings(sites)

	opts := optsSnapshot()
	total := 0
	for _, site := range sites {
		for _, config := range (*rssDict)[site] {
			if config.Enabled {
				total++
			}
		}
	}

	var results []FeedRunResult
	for _, site := range sites {
		configs := (*rssDict)[site]
		for i := range configs {
			config := configs[i]
			if !config.Enabled || config.Completed {
				log.Printf("[%s] disabled or completed, skipped\n", config.Name)
				continue
			}
			groupResults := ag.RunRssConfigGroup(&config)
			for i := range groupResults {
				if groupResults[i].Site == "" {
					groupResults[i].Site = site
				}
			}
			results = append(results, groupResults...)
			if len(results) < total {
				time.Sleep(time.Second * time.Duration(opts.chunkDelay))
			}
		}
	}
	return results, nil
}

func (ag *Agent) ExecuteAllRssTask() {
	results, err := ag.RunAllRssTasks()
	if err != nil {
		log.Printf("[error] %s\n", err)
		return
	}
	added := 0
	for _, result := range results {
		added += result.Added
	}
	log.Printf("[all] finished %d feed(s), added %d task(s)\n", len(results), added)
}

// AddMagnetTasks adds magnet/ed2k/http links and returns the number of tasks
// accepted by 115.
func (ag *Agent) AddMagnetTasks(magnets []string, cid string, savepath string) (int, error) {
	opts := optsSnapshot()
	added := 0
	chunks := chunkBy(magnets, opts.chunkSize)
	for i, urls := range chunks {
		hashes, err := ag.offlineAdd(urls, cid, savepath)
		if err != nil {
			log.Printf("Add offline error: %s\n", err)
			return added, err
		}
		for _, hash := range hashes {
			if hash != "" {
				added++
			}
		}
		log.Printf("[magnet] add %d tasks\n", len(urls))
		if i < len(chunks)-1 {
			time.Sleep(time.Second * time.Duration(opts.chunkDelay))
		}
	}
	return added, nil
}

func (ag *Agent) AddMagnetTask(magnets []string, cid string, savepath string) {
	if _, err := ag.AddMagnetTasks(magnets, cid, savepath); err != nil {
		log.Printf("Add offline error: %s\n", err)
	}
}
func (ag *Agent) OfflineClear(num int) (err error) {
	flag := elevengo.OfflineClearFlag(num)
	return ag.Agent.OfflineClear(flag)
}

// ListOfflineTasks returns current offline downloading tasks from 115.
// A positive limit stops remote pagination as soon as enough matching tasks
// have been collected. A nil include function matches every task.
func (ag *Agent) ListOfflineTasks(limit int, include func(*elevengo.OfflineTask) bool) (tasks []elevengo.OfflineTask, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			tasks = nil
			err = fmt.Errorf("115 离线任务接口异常: %v", recovered)
		}
	}()
	it, err := ag.Agent.OfflineIterate()
	if err != nil {
		return nil, err
	}
	seen := 0
	for _, t := range it.Items() {
		seen++
		if include != nil && !include(t) {
			continue
		}
		tasks = append(tasks, *t)
		if limit > 0 && len(tasks) >= limit {
			return tasks, nil
		}
	}
	if seen < it.Count() {
		return nil, fmt.Errorf("offline task listing stopped early: received %d of %d tasks", seen, it.Count())
	}
	return tasks, nil
}

func SaveCookies(agent *elevengo.Agent) {
	cr := &elevengo.Credential{}
	agent.CredentialExport(cr)
	cookies := fmt.Sprintf("UID=%s; CID=%s; SEID=%s; KID=%s", cr.UID, cr.CID, cr.SEID, cr.KID)

	// Determine save location using config package
	savePath := determineCookieSavePath()
	if err := os.WriteFile(savePath, []byte(cookies), 0600); err != nil {
		log.Printf("failed to save cookies to %s: %v\n", savePath, err)
	}
}

func LoadCookies() string {
	cookies, err := LoadCookiesWithError()
	if err != nil {
		return ""
	}
	return cookies
}

func LoadCookiesWithError() (string, error) {
	cfg, _, err := config.LoadWithOptions(config.CLIParams{}, config.LoadOptions{Auth: true})
	if err != nil {
		return "", err
	}
	return cfg.Auth.Cookies, nil
}

// determineCookieSavePath determines where to save cookies based on configuration source
// Priority: original file > TOML cookies_file > current directory .cookies
func determineCookieSavePath() string {
	_, source, err := config.LoadWithOptions(config.CLIParams{}, config.LoadOptions{Auth: true})
	if err != nil {
		// Fallback to current directory if config loading fails
		return config.ExistingCookiePathOrDefault()
	}

	// Priority 1: Use the file that was originally read
	if source.CookiesPath != "" {
		return source.CookiesPath
	}

	// Priority 2: Use TOML-configured cookies_file
	if source.TOMLPath != "" {
		cfg, err := config.LoadTOML(source.TOMLPath)
		if err == nil && cfg.Auth.CookiesFile != "" {
			return config.ResolveCookiesFile(cfg.Auth.CookiesFile, source.TOMLPath)
		}
	}

	// Priority 3: Fallback to current directory
	return config.ExistingCookiePathOrDefault()
}

func QrcodeLogin() (*Agent, error) {
	agent := newElevengoAgent()
	if agent == nil {
		return nil, errors.New("初始化 115 客户端失败")
	}
	session := &elevengo.QrcodeSession{}
	// @TODO: add option; default is tv
	err := agent.QrcodeStart(session, option.Qrcode().LoginTv())
	if err != nil {
		return nil, err
	}
	err = DisplayQrcode(session.Image)
	if err != nil {
		return nil, err
	}
	after := time.Now().Add(2 * time.Minute)
	for {
		time.Sleep(200 * time.Millisecond)
		success, err := agent.QrcodePoll(session)
		if success {
			SaveCookies(agent)
			DisposeQrcode()

			storeInstance, err := store.NewWithPath(optsSnapshot().databasePath)
			if err != nil {
				return nil, err
			}
			return &Agent{
				Agent:         agent,
				StoreInstance: storeInstance,
			}, nil
		}
		if err != nil && err == elevengo.ErrQrcodeCancelled {
			return nil, errors.New("login cancelled")
		}
		if time.Now().After(after) {
			return nil, errors.New("login timed out")
		}
	}
}

const fallbackAppVersion = "36.0.1"

func newElevengoAgent() *elevengo.Agent {
	return safeElevengoNew(buildAgentOptions(fallbackAppVersion))
}

func buildAgentOptions(version string) *option.AgentOptions {
	agentOpts := optsSnapshot()
	opts := option.Agent().WithHttpClient(newResilientHTTPClient())
	if agentOpts.cooldownMinMs > 0 || agentOpts.cooldownMaxMs > 0 {
		minMs := agentOpts.cooldownMinMs
		maxMs := agentOpts.cooldownMaxMs
		if maxMs == 0 {
			maxMs = minMs
		}
		if maxMs < minMs {
			maxMs = minMs
		}
		opts.WithCooldown(minMs, maxMs)
	}
	if version != "" {
		opts.WithVersion(version)
	}
	return opts
}

type resilientHTTPClient struct {
	inner *http.Client
}

func newResilientHTTPClient() *resilientHTTPClient {
	return &resilientHTTPClient{inner: request.NewIPv4Client(30 * time.Second)}
}

// Do prevents elevengo from receiving (nil body, nil error), which triggers a
// panic in its JSON parser when a network request fails. Network errors are
// converted into a normal JSON error response so callers get a useful error.
func (c *resilientHTTPClient) Do(req *http.Request) (*http.Response, error) {
	attempts := 1
	if req.Method == http.MethodGet {
		attempts = 2
	}
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		resp, err := c.inner.Do(req)
		if err == nil {
			if resp.Body == nil {
				resp.Body = http.NoBody
			}
			return resp, nil
		}
		lastErr = err
		if req.Context().Err() != nil {
			break
		}
	}
	return syntheticErrorResponse(req, lastErr), nil
}

func syntheticErrorResponse(req *http.Request, cause error) *http.Response {
	message := "115 网络请求失败"
	if cause != nil {
		message += ": " + cause.Error()
	}
	payload, _ := json.Marshal(map[string]any{"state": false, "code": -1, "message": message})
	return &http.Response{
		Status:        "502 Bad Gateway",
		StatusCode:    http.StatusBadGateway,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          io.NopCloser(bytes.NewReader(payload)),
		ContentLength: int64(len(payload)),
		Request:       req,
	}
}

func safeElevengoNew(opts *option.AgentOptions) (agent *elevengo.Agent) {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Printf("[warning] 115 客户端初始化异常: %v", recovered)
			agent = nil
		}
	}()
	return elevengo.New(opts)
}
