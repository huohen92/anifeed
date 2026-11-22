package webui

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/store"
)

type aiPendingAction struct {
	Action  string         `json:"action"`
	Args    map[string]any `json:"args"`
	Summary string         `json:"summary"`
	Expires time.Time      `json:"-"`
}

type aiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aiChatRequest struct {
	Messages []aiMessage `json:"messages"`
}

type aiChatResponse struct {
	Reply   string         `json:"reply"`
	Pending *aiPendingView `json:"pending,omitempty"`
}

type aiPendingView struct {
	Token   string `json:"token"`
	Summary string `json:"summary"`
}

type aiToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type aiAPIMessage struct {
	Role       string       `json:"role"`
	Content    string       `json:"content,omitempty"`
	ToolCalls  []aiToolCall `json:"tool_calls,omitempty"`
	ToolCallID string       `json:"tool_call_id,omitempty"`
}

type aiAPIResponse struct {
	Choices []struct {
		Message aiAPIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (s *Server) aiClientConfig() (*config.AIConfig, error) {
	cfg := s.configSnapshot()
	if cfg == nil {
		return nil, fmt.Errorf("配置未加载")
	}
	if !cfg.AI.Enabled {
		return nil, fmt.Errorf("AI 助手未启用")
	}
	if strings.TrimSpace(cfg.AI.APIKey) == "" {
		return nil, fmt.Errorf("请先填写 AI API Key")
	}
	if strings.TrimSpace(cfg.AI.BaseURL) == "" {
		return nil, fmt.Errorf("请填写 AI Base URL")
	}
	return &cfg.AI, nil
}

func (s *Server) handleAIChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}
	aiCfg, err := s.aiClientConfig()
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req aiChatRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if len(req.Messages) == 0 {
		writeAPIError(w, http.StatusBadRequest, "消息不能为空")
		return
	}
	if len(req.Messages) > 20 {
		req.Messages = req.Messages[len(req.Messages)-20:]
	}
	messages := []aiAPIMessage{{
		Role:    "system",
		Content: s.aiSystemPrompt(aiCfg),
	}}
	for _, msg := range req.Messages {
		role := strings.TrimSpace(msg.Role)
		if role != "user" && role != "assistant" {
			continue
		}
		content := strings.TrimSpace(msg.Content)
		if len(content) > 4000 {
			content = content[:4000]
		}
		if content == "" {
			continue
		}
		messages = append(messages, aiAPIMessage{Role: role, Content: content})
	}
	if len(messages) == 1 {
		writeAPIError(w, http.StatusBadRequest, "请输入问题")
		return
	}

	for round := 0; round < 3; round++ {
		reply, err := s.callAI(r.Context(), aiCfg, messages)
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, err.Error())
			return
		}
		if len(reply.ToolCalls) == 0 {
			writeJSON(w, http.StatusOK, aiChatResponse{Reply: strings.TrimSpace(reply.Content)})
			return
		}
		messages = append(messages, reply)
		for _, call := range reply.ToolCalls {
			if call.Type != "" && call.Type != "function" {
				continue
			}
			if isAIWriteTool(call.Function.Name) {
				view, err := s.proposeAIAction(aiCfg, call)
				if err != nil {
					writeAPIError(w, http.StatusForbidden, err.Error())
					return
				}
				writeJSON(w, http.StatusOK, aiChatResponse{
					Reply:   "该操作需要你二次确认后才会执行。",
					Pending: view,
				})
				return
			}
			result, err := s.runAIReadTool(call.Function.Name, call.Function.Arguments)
			if err != nil {
				result = map[string]any{"error": err.Error()}
			}
			data, _ := json.Marshal(result)
			messages = append(messages, aiAPIMessage{Role: "tool", ToolCallID: call.ID, Content: string(data)})
		}
	}
	writeJSON(w, http.StatusOK, aiChatResponse{Reply: "查询步骤过多，请缩小问题范围后重试。"})
}

func (s *Server) handleAIAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}
	cfg := s.configSnapshot()
	if cfg == nil || cfg.AI.Permission != "full" {
		writeAPIError(w, http.StatusForbidden, "当前 AI 权限不允许执行修改操作")
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	s.aiMu.Lock()
	pending, ok := s.aiPending[req.Token]
	if ok {
		delete(s.aiPending, req.Token)
	}
	s.aiMu.Unlock()
	if !ok || time.Now().After(pending.Expires) {
		writeAPIError(w, http.StatusBadRequest, "确认已过期，请重新发起")
		return
	}
	msg, err := s.executeAIAction(cfg, pending)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "message": msg})
}

func (s *Server) aiSystemPrompt(aiCfg *config.AIConfig) string {
	context := s.aiContextJSON()
	return `你是 AniFeed 的运维助手。AniFeed 是 RSS 自动离线到 115 网盘的动漫管理工具。
你主要回答查询问题。只能根据下面的数据和工具结果回答，不要编造。
禁止输出 Cookie、API Key、密码、Webhook、115 会话等敏感信息。
需要修改、删除、启用、停用订阅时，必须调用提供的工具；不要声称已经执行，系统会要求用户二次确认。
当前权限：` + aiCfg.Permission + `（readonly=只读，full=完整权限但修改仍需二次确认）。

当前数据摘要：
` + context
}

func (s *Server) aiContextJSON() string {
	cfg := s.configSnapshot()
	type feedView struct {
		Name      string `json:"name"`
		Enabled   bool   `json:"enabled"`
		Source    string `json:"source"`
		Total     int    `json:"total_episodes"`
		Offset    int    `json:"rename_offset"`
		BackupRSS int    `json:"backup_rss"`
	}
	type view struct {
		Version string     `json:"version"`
		Feeds   []feedView `json:"feeds"`
		Recent  []string   `json:"recent_history"`
	}
	v := view{Version: versionText()}
	if cfg != nil {
		for _, feeds := range cfg.RSS {
			for _, f := range feeds {
				v.Feeds = append(v.Feeds, feedView{Name: f.Name, Enabled: f.IsEnabled(), Source: f.Source, Total: f.TotalEpisodes, Offset: f.RenameOffset, BackupRSS: len(f.Subgroups)})
			}
		}
	}
	if st, err := store.NewWithPath(s.databasePath()); err == nil {
		defer st.Close()
		if items, err := st.ListHistory(10, ""); err == nil {
			for _, item := range items {
				v.Recent = append(v.Recent, item.Title)
			}
		}
	}
	data, _ := json.Marshal(v)
	return string(data)
}

func (s *Server) callAI(ctx context.Context, cfg *config.AIConfig, messages []aiAPIMessage) (aiAPIMessage, error) {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	endpoint := base
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint += "/chat/completions"
	}
	body := map[string]any{
		"model": strings.TrimSpace(cfg.Model), "messages": messages,
		"temperature": 0.2, "tools": aiTools(), "tool_choice": "auto",
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return aiAPIMessage{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(cfg.APIKey))
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return aiAPIMessage{}, fmt.Errorf("AI 接口连接失败: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return aiAPIMessage{}, fmt.Errorf("AI 接口返回 HTTP %d：%s", resp.StatusCode, truncateUTF8(string(raw), 300))
	}
	var parsed aiAPIResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return aiAPIMessage{}, fmt.Errorf("AI 返回格式错误")
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return aiAPIMessage{}, fmt.Errorf("AI 接口错误：%s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return aiAPIMessage{}, fmt.Errorf("AI 没有返回内容")
	}
	return parsed.Choices[0].Message, nil
}

func aiTools() []map[string]any {
	return []map[string]any{
		aiTool("list_feeds", "列出所有订阅及其启用状态、总集数、偏移和备用 RSS 数量", map[string]any{}),
		aiTool("get_status", "获取 AniFeed 当前运行状态摘要", map[string]any{}),
		aiTool("list_tasks", "列出最近 20 个 115 离线任务", map[string]any{}),
		aiTool("list_history", "列出最近 20 条下载历史", map[string]any{}),
		aiTool("set_feed_enabled", "启用或停用订阅；需要用户二次确认", map[string]any{"name": map[string]any{"type": "string"}, "enabled": map[string]any{"type": "boolean"}}, "name", "enabled"),
		aiTool("delete_feed", "删除订阅；需要用户二次确认", map[string]any{"name": map[string]any{"type": "string"}}, "name"),
	}
}

func aiTool(name, description string, props map[string]any, required ...string) map[string]any {
	if required == nil {
		required = []string{}
	}
	return map[string]any{
		"type": "function",
		"function": map[string]any{
			"name": name, "description": description,
			"parameters": map[string]any{"type": "object", "properties": props, "required": required},
		},
	}
}

func isAIWriteTool(name string) bool { return name == "set_feed_enabled" || name == "delete_feed" }

func (s *Server) proposeAIAction(aiCfg *config.AIConfig, call aiToolCall) (*aiPendingView, error) {
	if aiCfg.Permission != "full" {
		return nil, fmt.Errorf("当前为只读权限，不能执行修改或删除操作")
	}
	args := map[string]any{}
	if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
		return nil, fmt.Errorf("操作参数无效")
	}
	name, _ := args["name"].(string)
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("缺少订阅名称")
	}
	summary := ""
	switch call.Function.Name {
	case "set_feed_enabled":
		enabled, _ := args["enabled"].(bool)
		summary = fmt.Sprintf("将订阅「%s」设置为%s", name, map[bool]string{true: "启用", false: "停用"}[enabled])
	case "delete_feed":
		summary = fmt.Sprintf("删除订阅「%s」", name)
	default:
		return nil, fmt.Errorf("不支持的操作")
	}
	tokenBytes := make([]byte, 18)
	_, _ = rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)
	s.aiMu.Lock()
	if s.aiPending == nil {
		s.aiPending = map[string]aiPendingAction{}
	}
	s.aiPending[token] = aiPendingAction{Action: call.Function.Name, Args: args, Summary: summary, Expires: time.Now().Add(5 * time.Minute)}
	s.aiMu.Unlock()
	return &aiPendingView{Token: token, Summary: summary}, nil
}

func (s *Server) runAIReadTool(name, rawArgs string) (any, error) {
	cfg := s.configSnapshot()
	switch name {
	case "list_feeds":
		var out []map[string]any
		if cfg != nil {
			for _, feeds := range cfg.RSS {
				for _, f := range feeds {
					out = append(out, map[string]any{"name": f.Name, "enabled": f.IsEnabled(), "total": f.TotalEpisodes, "offset": f.RenameOffset, "backups": len(f.Subgroups)})
				}
			}
		}
		return out, nil
	case "get_status":
		return map[string]any{"version": versionText(), "feeds": len(allFeeds(cfg))}, nil
	case "list_tasks":
		ag, err := s.ensureAgent()
		if err != nil {
			return nil, err
		}
		tasks, err := ag.ListOfflineTasks(20, nil)
		if err != nil {
			return nil, err
		}
		var out []map[string]any
		for i := range tasks {
			out = append(out, map[string]any{"name": tasks[i].Name, "state": taskState(&tasks[i]), "percent": tasks[i].Percent})
		}
		return out, nil
	case "list_history":
		st, err := store.NewWithPath(s.databasePath())
		if err != nil {
			return nil, err
		}
		defer st.Close()
		items, err := st.ListHistory(20, "")
		if err != nil {
			return nil, err
		}
		var out []map[string]any
		for _, item := range items {
			out = append(out, map[string]any{"title": item.Title, "feed": item.ConfigName, "time": item.CreatedAt})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("未知查询工具")
	}
}

func allFeeds(cfg *config.Config) []config.RssConfig {
	var out []config.RssConfig
	if cfg == nil {
		return out
	}
	for _, feeds := range cfg.RSS {
		out = append(out, feeds...)
	}
	return out
}

func (s *Server) executeAIAction(cfg *config.Config, pending aiPendingAction) (string, error) {
	name, _ := pending.Args["name"].(string)
	resp, err := s.loadEditable()
	if err != nil {
		return "", err
	}
	wanted := strings.TrimSpace(name)
	exact := -1
	partial := -1
	for i := range resp.Config.RSS {
		feedName := strings.TrimSpace(resp.Config.RSS[i].Name)
		if strings.EqualFold(feedName, wanted) || strings.EqualFold(strings.TrimSpace(resp.Config.RSS[i].URL), wanted) {
			exact = i
			break
		}
		if strings.Contains(feedName, wanted) {
			if partial >= 0 {
				return "", fmt.Errorf("订阅「%s」匹配到多个，请提供完整名称", name)
			}
			partial = i
		}
	}
	if exact < 0 {
		exact = partial
	}
	if exact < 0 {
		return "", fmt.Errorf("未找到订阅「%s」", name)
	}
	switch pending.Action {
	case "set_feed_enabled":
		enabled, _ := pending.Args["enabled"].(bool)
		resp.Config.RSS[exact].Enabled = enabled
	case "delete_feed":
		resp.Config.RSS = append(resp.Config.RSS[:exact], resp.Config.RSS[exact+1:]...)
	default:
		return "", fmt.Errorf("不支持的操作")
	}
	log.Printf("[ai] %s", pending.Summary)
	if _, err := s.saveEditable(&resp.Config); err != nil {
		return "", err
	}
	return pending.Summary + "，已执行", nil
}

func versionText() string {
	// Kept local to avoid importing version in every AI helper.
	return "0.03"
}
