package webui

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/huohen92/anifeed/config"
)

type tmdbResult struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Title            string  `json:"title"`
	Original         string  `json:"original_name"`
	Date             string  `json:"first_air_date"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	Overview         string  `json:"overview"`
	VoteAverage      float64 `json:"vote_average"`
	NumberOfEpisodes int     `json:"number_of_episodes,omitempty"`
	Poster           string  `json:"poster,omitempty"`
}

func (s *Server) handleTMDBSearch(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusInternalServerError, "配置尚未加载")
		return
	}
	if strings.TrimSpace(cfg.TMDB.APIKey) == "" {
		writeAPIError(w, http.StatusBadRequest, "请先在设置中填写 TmdbApiKey")
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少搜索关键词")
		return
	}
	mediaType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("type")))
	if mediaType != "movie" {
		mediaType = "tv"
	}
	results, err := s.searchTMDB(cfg.TMDB, query, mediaType)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

func (s *Server) searchTMDB(cfg config.TMDBConfig, query, mediaType string) ([]tmdbResult, error) {
	values := url.Values{}
	values.Set("api_key", cfg.APIKey)
	values.Set("query", query)
	values.Set("language", "zh-CN")
	values.Set("include_adult", "false")
	var payload struct {
		Results []tmdbResult `json:"results"`
	}
	if err := s.fetchTMDB(cfg, "/3/search/"+mediaType, values, &payload); err != nil {
		return nil, err
	}
	for i := range payload.Results {
		payload.Results[i].Poster = tmdbImageURL(cfg.Image, payload.Results[i].PosterPath, "w342")
		if payload.Results[i].Name == "" {
			payload.Results[i].Name = payload.Results[i].Title
		}
	}
	return payload.Results, nil
}

func (s *Server) handleTMDBAuto(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusInternalServerError, "配置尚未加载")
		return
	}
	if strings.TrimSpace(cfg.TMDB.APIKey) == "" {
		writeAPIError(w, http.StatusBadRequest, "请先在设置中填写 TmdbApiKey")
		return
	}
	name := cleanTMDBQuery(r.URL.Query().Get("q"))
	if name == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少订阅名称")
		return
	}
	mediaType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("type")))
	if mediaType != "movie" {
		mediaType = "tv"
	}
	results, err := s.searchTMDB(cfg.TMDB, name, mediaType)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	best, score := bestTMDBMatch(name, results)
	if score < 500 {
		writeJSON(w, http.StatusOK, map[string]any{"matched": false, "query": name, "candidates": results})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"matched": true, "query": name, "item": best, "candidates": results})
}

func cleanTMDBQuery(raw string) string {
	value := strings.TrimSpace(raw)
	if idx := strings.Index(value, " · "); idx > 0 {
		value = value[:idx]
	}
	for _, marker := range []string{"[tmdbid=", "{tmdbid=", "{tmdbib="} {
		if idx := strings.Index(value, marker); idx > 0 {
			value = value[:idx]
		}
	}
	replacements := []struct{ old, new string }{
		{"第二季", ""}, {"第三季", ""}, {"第四季", ""}, {"第一季", ""},
		{"第二期", ""}, {"第三期", ""}, {"第四期", ""}, {"第一期", ""},
	}
	for _, re := range replacements {
		value = strings.ReplaceAll(value, re.old, re.new)
	}
	value = strings.ReplaceAll(value, "Season 2", "")
	value = strings.ReplaceAll(value, "Season 3", "")
	value = strings.ReplaceAll(value, "Season 4", "")
	value = strings.ReplaceAll(value, "Season 1", "")
	value = strings.TrimSpace(value)
	return value
}

func normalizeTMDBName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r > 127:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func bestTMDBMatch(query string, results []tmdbResult) (tmdbResult, int) {
	q := normalizeTMDBName(query)
	var best tmdbResult
	bestScore := 0
	for _, item := range results {
		for _, name := range []string{item.Name, item.Title, item.Original} {
			n := normalizeTMDBName(name)
			if n == "" {
				continue
			}
			score := 0
			switch {
			case n == q:
				score = 1000
			case strings.Contains(n, q) || strings.Contains(q, n):
				score = 700
				if len(n) > 0 && len(q) > 0 {
					if len(n) > len(q) {
						score += len(q) * 20 / len(n)
					} else {
						score += len(n) * 20 / len(q)
					}
				}
			}
			if score > bestScore {
				bestScore = score
				best = item
			}
		}
	}
	return best, bestScore
}

func (s *Server) handleTMDBDetail(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil || strings.TrimSpace(cfg.TMDB.APIKey) == "" {
		writeAPIError(w, http.StatusBadRequest, "请先在设置中填写 TmdbApiKey")
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少 TMDB ID")
		return
	}
	mediaType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("type")))
	if mediaType != "movie" {
		mediaType = "tv"
	}
	values := url.Values{}
	values.Set("api_key", cfg.TMDB.APIKey)
	values.Set("language", "zh-CN")
	var item tmdbResult
	if err := s.fetchTMDB(cfg.TMDB, "/3/"+mediaType+"/"+url.PathEscape(id), values, &item); err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	item.Poster = tmdbImageURL(cfg.TMDB.Image, item.PosterPath, "w500")
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) fetchTMDB(cfg config.TMDBConfig, path string, values url.Values, dst any) error {
	api := strings.TrimRight(strings.TrimSpace(cfg.API), "/")
	if api == "" {
		api = "https://api.themoviedb.org"
	}
	if !strings.HasSuffix(api, "/3") {
		api += "/3"
	}
	// path is /3/...
	if strings.HasPrefix(path, "/3/") {
		path = strings.TrimPrefix(path, "/3")
	}
	rawURL := api + path
	if len(values) > 0 {
		rawURL += "?" + values.Encode()
	}
	if err := fetchBrowseJSON(rawURL, dst); err != nil {
		return fmt.Errorf("TMDB 请求失败: %w", err)
	}
	return nil
}

func tmdbImageURL(base, path, size string) string {
	if path == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		base = "https://image.tmdb.org"
	}
	return base + "/t/p/" + size + path
}

func parseTMDBID(raw string) int {
	value, _ := strconv.Atoi(strings.TrimSpace(raw))
	return value
}
