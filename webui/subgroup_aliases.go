package webui

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// SubgroupAlias maps one canonical subtitle-group name to RSS/URL aliases.
type SubgroupAlias struct {
	Main    string   `json:"main"`
	Aliases []string `json:"aliases"`
}

type subgroupAliasStore struct {
	Aliases []SubgroupAlias `json:"aliases"`
}

func (s *Server) subgroupAliasesPath() string {
	return filepath.Join(filepath.Dir(s.ConfigPath()), "subgroup_aliases.json")
}

func defaultSubgroupAliases() []SubgroupAlias {
	return []SubgroupAlias{
		{Main: "今晚月色真美", Aliases: []string{"Tsukigakirei", "tsukigakirei"}},
	}
}

func (s *Server) loadSubgroupAliases() subgroupAliasStore {
	store := subgroupAliasStore{Aliases: []SubgroupAlias{}}
	data, err := os.ReadFile(s.subgroupAliasesPath())
	if err != nil {
		if os.IsNotExist(err) {
			store.Aliases = defaultSubgroupAliases()
		}
		return store
	}
	_ = json.Unmarshal(data, &store)
	if store.Aliases == nil {
		store.Aliases = []SubgroupAlias{}
	}
	return store
}

func (s *Server) saveSubgroupAliases(store subgroupAliasStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.subgroupAliasesPath()), 0o755); err != nil {
		return err
	}
	tmp := s.subgroupAliasesPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.subgroupAliasesPath())
}

func normalizeSubgroupAlias(value string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r > 127 {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (s *Server) canonicalSubgroup(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	target := normalizeSubgroupAlias(name)
	for _, item := range s.loadSubgroupAliases().Aliases {
		for _, candidate := range append([]string{item.Main}, item.Aliases...) {
			if normalizeSubgroupAlias(candidate) == target {
				return item.Main
			}
		}
	}
	return name
}

func cleanAliases(main string, aliases []string) []string {
	out := make([]string, 0, len(aliases))
	seen := map[string]bool{normalizeSubgroupAlias(main): true}
	for _, alias := range aliases {
		alias = strings.TrimSpace(alias)
		key := normalizeSubgroupAlias(alias)
		if alias == "" || key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, alias)
	}
	return out
}

func (s *Server) handleSubgroupAliases(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.loadSubgroupAliases())
	case http.MethodPut:
		var payload subgroupAliasStore
		if err := decodeJSON(w, r, &payload); err != nil {
			writeAPIError(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		clean := make([]SubgroupAlias, 0, len(payload.Aliases))
		seen := map[string]bool{}
		for _, item := range payload.Aliases {
			item.Main = strings.TrimSpace(item.Main)
			if item.Main == "" {
				writeAPIError(w, http.StatusBadRequest, "主名称不能为空")
				return
			}
			key := normalizeSubgroupAlias(item.Main)
			if seen[key] {
				writeAPIError(w, http.StatusBadRequest, "主名称重复："+item.Main)
				return
			}
			seen[key] = true
			item.Aliases = cleanAliases(item.Main, item.Aliases)
			clean = append(clean, item)
		}
		if err := s.saveSubgroupAliases(subgroupAliasStore{Aliases: clean}); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, subgroupAliasStore{Aliases: clean})
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// subgroupAliasFromURL resolves a groupSlug through the alias registry, then
// through template aliases for installations that have not created an alias.
func (s *Server) subgroupAliasFromURL(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ""
	}
	slug := normalizeSubgroupAlias(u.Query().Get("groupSlug"))
	if slug == "" {
		return ""
	}
	for _, item := range s.loadSubgroupAliases().Aliases {
		for _, candidate := range append([]string{item.Main}, item.Aliases...) {
			if normalizeSubgroupAlias(candidate) == slug {
				return item.Main
			}
		}
	}
	for _, template := range s.loadMatchTemplates().Templates {
		for _, candidate := range append([]string{template.Subgroup}, template.Aliases...) {
			if normalizeSubgroupAlias(candidate) == slug {
				return template.Subgroup
			}
		}
	}
	return ""
}
