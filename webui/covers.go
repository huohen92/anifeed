package webui

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/huohen92/anifeed/config"
)

var bangumiAPIBase = "https://api.bgm.tv"

type coverEntry struct {
	Cover   string
	Title   string
	Score   float64
	Expires time.Time
}

type coverResponse struct {
	Cover string  `json:"cover"`
	Title string  `json:"title,omitempty"`
	Score float64 `json:"score,omitempty"`
}

func (s *Server) cacheRemoteCover(rawURL string, info coverResponse) coverResponse {
	if strings.TrimSpace(info.Cover) == "" {
		return info
	}
	data, err := fetchBrowseBytes(info.Cover)
	if err != nil {
		return info
	}
	cover, err := s.savePoster(rawURL, data)
	if err != nil {
		return info
	}
	info.Cover = cover
	return info
}

func (s *Server) handleFeedCover(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, max-age=300")
	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	if rawURL == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少订阅 URL")
		return
	}
	if local := s.localPosterURL(rawURL); local != "" {
		writeJSON(w, http.StatusOK, coverResponse{Cover: local})
		return
	}
	bgmURL := strings.TrimSpace(r.URL.Query().Get("bgm_url"))
	if bgmID := bangumiIDFromValue(bgmURL); bgmID != "" {
		key := "bgm:" + bgmID
		s.coverMu.Lock()
		cached, ok := s.covers[key]
		valid := ok && time.Now().Before(cached.Expires)
		s.coverMu.Unlock()
		if valid && cached.Cover != "" {
			writeJSON(w, http.StatusOK, coverResponse{Cover: cached.Cover, Title: cached.Title, Score: cached.Score})
			return
		}
		if info, err := fetchBangumiCover(bgmID); err == nil && info.Cover != "" {
			info = s.cacheRemoteCover(rawURL, info)
			entry := coverEntry{Cover: info.Cover, Title: info.Title, Score: info.Score, Expires: time.Now().Add(6 * time.Hour)}
			s.coverMu.Lock()
			s.covers[key] = entry
			s.coverMu.Unlock()
			writeJSON(w, http.StatusOK, info)
			return
		}
	}
	source := config.NormalizeSource(r.URL.Query().Get("source"), rawURL)
	tmdbID := strings.TrimSpace(r.URL.Query().Get("tmdb_id"))
	if tmdbID != "" {
		mediaType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tmdb_type")))
		if mediaType != "movie" {
			mediaType = "tv"
		}
		key := "tmdb:" + mediaType + ":" + tmdbID
		s.coverMu.Lock()
		cached, ok := s.covers[key]
		valid := ok && time.Now().Before(cached.Expires)
		s.coverMu.Unlock()
		if valid && cached.Cover != "" {
			writeJSON(w, http.StatusOK, coverResponse{Cover: cached.Cover, Title: cached.Title, Score: cached.Score})
			return
		}
		if cfg := s.configSnapshot(); cfg != nil && cfg.TMDB.APIKey != "" {
			info, err := s.fetchTMDBPoster(tmdbID, mediaType)
			entry := coverEntry{Expires: time.Now().Add(6 * time.Hour)}
			if err == nil {
				entry.Cover, entry.Title, entry.Score = info.Cover, info.Title, info.Score
			} else {
				entry.Expires = time.Now().Add(2 * time.Minute)
			}
			s.coverMu.Lock()
			s.covers[key] = entry
			s.coverMu.Unlock()
			if err == nil && info.Cover != "" {
				info = s.cacheRemoteCover(rawURL, info)
				writeJSON(w, http.StatusOK, info)
				return
			}
		}
	}
	id := feedSourceID(source, rawURL)
	if id == "" {
		writeAPIError(w, http.StatusBadRequest, "该订阅没有可用于获取海报的站点 ID")
		return
	}
	key := source + ":" + id

	s.coverMu.Lock()
	if cached, ok := s.covers[key]; ok && time.Now().Before(cached.Expires) {
		s.coverMu.Unlock()
		if cached.Cover == "" {
			writeAPIError(w, http.StatusNotFound, "没有找到海报")
			return
		}
		writeJSON(w, http.StatusOK, coverResponse{Cover: cached.Cover, Title: cached.Title, Score: cached.Score})
		return
	}
	s.coverMu.Unlock()

	info, err := s.fetchCover(source, id, rawURL)
	entry := coverEntry{Expires: time.Now().Add(6 * time.Hour)}
	if err == nil && info.Cover != "" {
		info = s.cacheRemoteCover(rawURL, info)
		entry.Cover, entry.Title, entry.Score = info.Cover, info.Title, info.Score
	} else {
		entry.Expires = time.Now().Add(2 * time.Minute)
	}
	s.coverMu.Lock()
	s.covers[key] = entry
	s.coverMu.Unlock()
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	if info.Cover == "" {
		writeAPIError(w, http.StatusNotFound, "没有找到海报")
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func feedSourceID(source, rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	q := u.Query()
	switch config.NormalizeSource(source, rawURL) {
	case config.SourceMikan:
		return strings.TrimSpace(q.Get("bangumiId"))
	case config.SourceAniBT:
		return strings.TrimSpace(q.Get("bgmId"))
	case config.SourceAnimeGarden:
		return strings.TrimSpace(q.Get("subject"))
	default:
		return ""
	}
}

func (s *Server) fetchCover(source, id, rawURL string) (coverResponse, error) {
	switch source {
	case config.SourceMikan:
		return fetchMikanCover(id, rawURL)
	case config.SourceAniBT, config.SourceAnimeGarden:
		return fetchBangumiCover(id)
	default:
		return coverResponse{}, fmt.Errorf("该来源不支持海报")
	}
}

func fetchBangumiCover(id string) (coverResponse, error) {
	var payload struct {
		Name   string `json:"name"`
		NameCN string `json:"name_cn"`
		Images struct {
			Large  string `json:"large"`
			Common string `json:"common"`
			Medium string `json:"medium"`
			Small  string `json:"small"`
		} `json:"images"`
		Rating struct {
			Score float64 `json:"score"`
		} `json:"rating"`
	}
	var fetchErr error
	for attempt := 0; attempt < 3; attempt++ {
		fetchErr = fetchBrowseJSON(bangumiAPIBase+"/v0/subjects/"+url.PathEscape(id), &payload)
		if fetchErr == nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	if fetchErr != nil {
		return coverResponse{}, fmt.Errorf("获取 Bangumi 海报失败: %w", fetchErr)
	}
	cover := payload.Images.Large
	if cover == "" {
		cover = payload.Images.Common
	}
	if cover == "" {
		cover = payload.Images.Medium
	}
	if cover == "" {
		cover = payload.Images.Small
	}
	title := payload.NameCN
	if title == "" {
		title = payload.Name
	}
	return coverResponse{Cover: cover, Title: title, Score: payload.Rating.Score}, nil
}

func fetchMikanCover(id, rawURL string) (coverResponse, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return coverResponse{}, err
	}
	base := u.Scheme + "://" + u.Host
	pageURL := base + "/Home/Bangumi/" + url.PathEscape(id)
	body, err := fetchBrowseBytes(pageURL)
	if err != nil {
		return coverResponse{}, fmt.Errorf("获取 Mikan 海报失败: %w", err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return coverResponse{}, err
	}
	cover, _ := doc.Find(".content > img").First().Attr("src")
	return coverResponse{
		Cover: absoluteBrowseURL(base, cover),
		Title: strings.TrimSpace(doc.Find(".bangumi-title").First().Text()),
	}, nil
}

func (s *Server) fetchTMDBPoster(id, mediaType string) (coverResponse, error) {
	cfg := s.configSnapshot()
	if cfg == nil || cfg.TMDB.APIKey == "" {
		return coverResponse{}, fmt.Errorf("未配置 TMDB API Key")
	}
	values := url.Values{}
	values.Set("api_key", cfg.TMDB.APIKey)
	values.Set("language", "zh-CN")
	var item tmdbResult
	if err := s.fetchTMDB(cfg.TMDB, "/3/"+mediaType+"/"+url.PathEscape(id), values, &item); err != nil {
		return coverResponse{}, err
	}
	title := item.Name
	if title == "" {
		title = item.Title
	}
	return coverResponse{Cover: tmdbImageURL(cfg.TMDB.Image, item.PosterPath, "w500"), Title: title, Score: item.VoteAverage}, nil
}
