package webui

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/huohen92/anifeed/ruleengine"
)

// SubgroupMatchTemplate is a reusable matching preset for one subtitle group.
// Templates live next to config.toml so they survive restarts but stay out of
// the hand-written TOML configuration.
type SubgroupMatchTemplate struct {
	ID       string   `json:"id"`
	Subgroup string   `json:"subgroup"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	Aliases  []string `json:"aliases,omitempty"`
	Filter   string   `json:"filter"`
	Exclude  string   `json:"exclude"`
}

type matchTemplateStore struct {
	Templates []SubgroupMatchTemplate `json:"templates"`
}

func (s *Server) matchTemplatesPath() string {
	return filepath.Join(filepath.Dir(s.ConfigPath()), "match_templates.json")
}

func (s *Server) loadMatchTemplates() matchTemplateStore {
	store := matchTemplateStore{Templates: []SubgroupMatchTemplate{}}
	data, err := os.ReadFile(s.matchTemplatesPath())
	if err != nil {
		if os.IsNotExist(err) {
			store.Templates = presetMatchTemplates()
		}
		return store
	}
	_ = json.Unmarshal(data, &store)
	if store.Templates == nil {
		store.Templates = []SubgroupMatchTemplate{}
	}
	return store
}

func (s *Server) saveMatchTemplates(store matchTemplateStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(s.matchTemplatesPath())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp := s.matchTemplatesPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.matchTemplatesPath())
}

func validateMatchTemplate(t SubgroupMatchTemplate) (SubgroupMatchTemplate, error) {
	t.Subgroup = strings.TrimSpace(t.Subgroup)
	t.Name = strings.TrimSpace(t.Name)
	if t.Subgroup == "" {
		return t, fmt.Errorf("字幕组不能为空")
	}
	if t.Name == "" {
		return t, fmt.Errorf("模板名称不能为空")
	}
	if _, err := ruleengine.Parse(t.Filter); err != nil {
		return t, fmt.Errorf("匹配规则：%w", err)
	}
	if _, err := ruleengine.Parse(t.Exclude); err != nil {
		return t, fmt.Errorf("排除规则：%w", err)
	}
	t.Aliases = cleanTemplateTags(t.Aliases)
	t.Tags = cleanTemplateTags(t.Tags)
	return t, nil
}

func cleanTemplateTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	seen := map[string]bool{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		out = append(out, tag)
	}
	return out
}

func (s *Server) handleMatchTemplates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.loadMatchTemplates())
	case http.MethodPut:
		var payload matchTemplateStore
		if err := decodeJSON(w, r, &payload); err != nil {
			writeAPIError(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		seen := map[string]bool{}
		clean := make([]SubgroupMatchTemplate, 0, len(payload.Templates))
		for _, item := range payload.Templates {
			item, err := validateMatchTemplate(item)
			if err != nil {
				writeAPIError(w, http.StatusBadRequest, err.Error())
				return
			}
			if strings.TrimSpace(item.ID) == "" || seen[item.ID] {
				random := make([]byte, 6)
				_, _ = rand.Read(random)
				item.ID = "tpl-" + hex.EncodeToString(random)
			}
			seen[item.ID] = true
			clean = append(clean, item)
		}
		if err := s.saveMatchTemplates(matchTemplateStore{Templates: clean}); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, matchTemplateStore{Templates: clean})
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}
