package webui

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
)

type weekdayCacheEntry struct {
	Weekday int
	Date    string
	Episode int
	Expires time.Time
}

func (s *Server) handleFeedWeekday(w http.ResponseWriter, r *http.Request) {
	tmdbID := strings.TrimSpace(r.URL.Query().Get("tmdb_id"))
	if tmdbID == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少 TMDB ID")
		return
	}
	cfg := s.configSnapshot()
	if cfg == nil || cfg.TMDB.APIKey == "" {
		writeAPIError(w, http.StatusBadRequest, "请先配置 TmdbApiKey")
		return
	}
	key := "tv:" + tmdbID
	s.coverMu.Lock()
	cached, ok := s.weekdays[key]
	s.coverMu.Unlock()
	if ok && time.Now().Before(cached.Expires) {
		writeJSON(w, http.StatusOK, map[string]any{
			"weekday": cached.Weekday, "date": cached.Date, "episode": cached.Episode,
		})
		return
	}
	weekday, date, episode, err := s.fetchLatestEpisodeWeekday(cfg.TMDB, tmdbID)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	entry := weekdayCacheEntry{Weekday: weekday, Date: date, Episode: episode, Expires: time.Now().Add(6 * time.Hour)}
	s.coverMu.Lock()
	if s.weekdays == nil {
		s.weekdays = make(map[string]weekdayCacheEntry)
	}
	s.weekdays[key] = entry
	s.coverMu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"weekday": weekday, "date": date, "episode": episode})
}

func (s *Server) fetchLatestEpisodeWeekday(cfg config.TMDBConfig, tmdbID string) (int, string, int, error) {
	values := url.Values{}
	values.Set("api_key", cfg.APIKey)
	values.Set("language", "zh-CN")
	var detail struct {
		Seasons []struct {
			SeasonNumber int    `json:"season_number"`
			AirDate      string `json:"air_date"`
		} `json:"seasons"`
	}
	if err := s.fetchTMDB(cfg, "/3/tv/"+url.PathEscape(tmdbID), values, &detail); err != nil {
		return 0, "", 0, err
	}
	now := time.Now()
	var latest time.Time
	latestEpisode := 0
	for _, season := range detail.Seasons {
		if season.SeasonNumber <= 0 {
			continue
		}
		var data struct {
			Episodes []struct {
				EpisodeNumber int    `json:"episode_number"`
				AirDate       string `json:"air_date"`
			} `json:"episodes"`
		}
		seasonValues := url.Values{}
		seasonValues.Set("api_key", cfg.APIKey)
		seasonValues.Set("language", "zh-CN")
		path := "/3/tv/" + url.PathEscape(tmdbID) + "/season/" + strconv.Itoa(season.SeasonNumber)
		if err := s.fetchTMDB(cfg, path, seasonValues, &data); err != nil {
			continue
		}
		for _, episode := range data.Episodes {
			if episode.AirDate == "" {
				continue
			}
			date, err := time.ParseInLocation("2006-01-02", episode.AirDate, time.Local)
			if err != nil || date.After(now) {
				continue
			}
			if date.After(latest) {
				latest = date
				latestEpisode = episode.EpisodeNumber
			}
		}
	}
	if latest.IsZero() {
		return 0, "", 0, errNoEpisodeDate
	}
	weekday := int(latest.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return weekday, latest.Format("2006-01-02"), latestEpisode, nil
}

var errNoEpisodeDate = errors.New("TMDB 没有已播出的集数日期")
