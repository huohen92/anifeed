package webui

import (
	"net/http"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/rsssite"
)

type seasonSuggestResponse struct {
	Season int    `json:"season"`
	Reason string `json:"reason"`
}

func (s *Server) handleSeasonSuggest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}
	var req struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		TmdbID string `json:"tmdb_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	cfg := s.configSnapshot()
	result := seasonSuggestResponse{Season: 1, Reason: "默认第一季"}
	if req.TmdbID != "" && cfg != nil && cfg.TMDB.APIKey != "" {
		if seasons, err := s.tmdbSeasonNumbers(cfg, req.TmdbID); err == nil {
			if len(seasons) == 1 {
				result = seasonSuggestResponse{Season: seasons[0], Reason: "TMDB 只有一季"}
			} else if len(seasons) > 1 {
				hint := detectSeason(req.Name)
				if hint == 0 {
					hint = detectSeason(req.URL)
				}
				if hint > 0 && containsInt(seasons, hint) {
					result = seasonSuggestResponse{Season: hint, Reason: "从订阅名称识别季数"}
				}
			}
		}
	}
	if result.Season <= 1 && req.URL != "" {
		items := rsssite.GetMagnetItemList(&rsssite.RssConfig{Url: req.URL, Source: config.NormalizeSource("", req.URL)})
		for i := range items {
			hint := detectSeason(items[i].Title)
			if hint > 1 {
				// Only trust the title's season label when TMDB also has it.
				// If TMDB has one season, keep TMDB's season 1.
				if req.TmdbID == "" {
					result = seasonSuggestResponse{Season: hint, Reason: "从最近 RSS 标题识别季数"}
				}
				break
			}
		}
		// If TMDB has multiple seasons and RSS has no season label, use the
		// closest TMDB air date to the newest RSS item.
		if result.Season <= 1 && req.TmdbID != "" && cfg != nil {
			if episodes, err := s.tmdbEpisodes(cfg, req.TmdbID); err == nil && len(episodes) > 0 {
				var latest time.Time
				for i := range items {
					if items[i].Published.After(latest) {
						latest = items[i].Published
					}
				}
				if !latest.IsZero() {
					best := tmdbEpisode{}
					bestDiff := 10 * 24 * time.Hour
					for _, ep := range episodes {
						d := latest.Sub(ep.AirDate)
						if d < 0 {
							d = -d
						}
						if d < bestDiff {
							bestDiff, best = d, ep
						}
					}
					if best.Season > 0 {
						result = seasonSuggestResponse{Season: best.Season, Reason: "按最新 RSS 发布日期匹配 TMDB 季"}
					}
				}
			}
		}
	}
	if result.Season <= 0 {
		result.Season = 1
	}
	writeJSON(w, http.StatusOK, result)
}

func containsInt(values []int, target int) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func (s *Server) tmdbSeasonNumbers(cfg *config.Config, id string) ([]int, error) {
	var detail struct {
		Seasons []struct {
			SeasonNumber int `json:"season_number"`
		} `json:"seasons"`
	}
	values := map[string]string{"api_key": cfg.TMDB.APIKey, "language": "zh-CN"}
	if err := s.fetchTMDBQuery(cfg, "/3/tv/"+id, values, &detail); err != nil {
		return nil, err
	}
	var out []int
	for _, season := range detail.Seasons {
		if season.SeasonNumber > 0 {
			out = append(out, season.SeasonNumber)
		}
	}
	return out, nil
}

var _ = p115.ExtractEpisode
var _ = strings.TrimSpace
