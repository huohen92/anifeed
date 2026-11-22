package webui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/rsssite"
)

type aniSeason struct {
	Season   int
	Episodes int
	Start    string
}

type aniMedia struct {
	ID       int `json:"id"`
	Episodes int `json:"episodes"`
	Title    struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	StartDate struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"startDate"`
}

func (s *Server) anilistSeasons(cfg *config.Config, feed *config.RssConfig) ([]aniSeason, error) {
	titles := []string{feed.Name}
	if cfg.TMDB.APIKey != "" && feed.TmdbID != "" {
		if names, err := s.tmdbTVTitles(cfg, feed.TmdbID); err == nil {
			titles = append(titles, names...)
		}
	}
	seen := map[int]bool{}
	var media []aniMedia
	for _, title := range uniqueStrings(titles) {
		if strings.TrimSpace(title) == "" {
			continue
		}
		list, err := searchAniList(title)
		if err != nil {
			continue
		}
		base := normalizeAnimeBase(title)
		for _, m := range list {
			if seen[m.ID] {
				continue
			}
			names := []string{m.Title.Romaji, m.Title.English, m.Title.Native}
			match := false
			for _, n := range names {
				if normalizeAnimeBase(n) == base || strings.Contains(normalizeAnimeBase(n), base) || strings.Contains(base, normalizeAnimeBase(n)) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
			seen[m.ID] = true
			media = append(media, m)
		}
	}
	if len(media) == 0 {
		return nil, fmt.Errorf("AniList 未找到番剧")
	}
	sort.Slice(media, func(i, j int) bool {
		return media[i].StartDate.Year*10000+media[i].StartDate.Month*100+media[i].StartDate.Day <
			media[j].StartDate.Year*10000+media[j].StartDate.Month*100+media[j].StartDate.Day
	})
	out := make([]aniSeason, 0, len(media))
	for i, m := range media {
		out = append(out, aniSeason{Season: i + 1, Episodes: m.Episodes, Start: fmt.Sprintf("%04d-%02d-%02d", m.StartDate.Year, m.StartDate.Month, m.StartDate.Day)})
	}
	return out, nil
}

func (s *Server) suggestByAniList(cfg *config.Config, feed *config.RssConfig, rssURL string) (*offsetSuggestion, error) {
	ani, err := s.anilistSeasons(cfg, feed)
	if err != nil {
		return nil, err
	}
	tmdbSeasons, err := s.tmdbSeasonNumbers(cfg, feed.TmdbID)
	if err != nil {
		return nil, err
	}
	items := rsssite.GetMagnetItemList(&rsssite.RssConfig{Url: rssURL, Source: config.NormalizeSource(feed.Source, rssURL)})
	if len(items) == 0 {
		return nil, fmt.Errorf("RSS 没有条目")
	}
	hint := 0
	for i := range items {
		if h := detectSeason(items[i].Title); h > 0 {
			hint = h
			break
		}
	}
	// TMDB has the same season: no offset is needed.
	firstEpisode, _ := p115.ExtractEpisode(items[0].Title)
	if hint > 0 && hint <= len(ani) && firstEpisode > 0 && int(firstEpisode) > ani[hint-1].Episodes {
		// The title says "Season 2", but episode numbers exceed that season's
		// length. This is an absolute-numbered release (e.g. LoliHouse 31).
		hint = 0
	}
	if hint > 0 && containsInt(tmdbSeasons, hint) {
		return &offsetSuggestion{Season: hint, Offset: 0, Samples: len(items), Total: len(items), Confidence: 1}, nil
	}
	// TMDB collapsed all seasons into one: use AniList season lengths.
	if len(tmdbSeasons) == 1 && hint > 0 && hint <= len(ani) {
		offset := 0
		for i := 0; i < hint-1; i++ {
			offset += ani[i].Episodes
		}
		return &offsetSuggestion{Season: tmdbSeasons[0], Offset: offset, Samples: len(items), Total: len(items), Confidence: 1}, nil
	}
	// Absolute numbering with one TMDB season: already in TMDB order.
	if hint == 0 && len(tmdbSeasons) == 1 {
		return &offsetSuggestion{Season: tmdbSeasons[0], Offset: 0, Samples: len(items), Total: len(items), Confidence: 1}, nil
	}
	// Absolute numbering with multiple TMDB seasons: map by AniList cumulative counts.
	if hint == 0 && len(tmdbSeasons) > 1 && len(ani) == len(tmdbSeasons) {
		ep, err := p115.ExtractEpisode(items[0].Title)
		if err == nil && ep > 0 {
			cum := 0
			for i, season := range ani {
				if int(ep) <= cum+season.Episodes {
					return &offsetSuggestion{Season: tmdbSeasons[i], Offset: -cum, Samples: len(items), Total: len(items), Confidence: 1}, nil
				}
				cum += season.Episodes
			}
		}
	}
	return nil, fmt.Errorf("AniList/TMDB 季结构无法直接映射")
}

func (s *Server) tmdbTVTitles(cfg *config.Config, id string) ([]string, error) {
	var detail struct {
		Name         string `json:"name"`
		OriginalName string `json:"original_name"`
	}
	values := map[string]string{"api_key": cfg.TMDB.APIKey, "language": "zh-CN"}
	if err := s.fetchTMDBQuery(cfg, "/3/tv/"+id, values, &detail); err != nil {
		return nil, err
	}
	return []string{detail.Name, detail.OriginalName}, nil
}

func searchAniList(search string) ([]aniMedia, error) {
	query := `query($search:String){ Page(perPage:10){ media(search:$search,type:ANIME){ id episodes title{romaji english native} startDate{year month day} } } }`
	body, _ := json.Marshal(map[string]any{"query": query, "variables": map[string]string{"search": search}})
	req, err := http.NewRequest(http.MethodPost, "https://graphql.anilist.co", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AniList HTTP %d", resp.StatusCode)
	}
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var parsed struct {
		Data struct {
			Page struct {
				Media []aniMedia `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	return parsed.Data.Page.Media, nil
}

func normalizeAnimeBase(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacements := []string{"第二季", "第三季", "第四季", "season 2", "season 3", "season 4", "2nd season", "3rd season", "4th season", "s2", "s3", "s4"}
	for _, r := range replacements {
		value = strings.ReplaceAll(value, r, "")
	}
	value = strings.NewReplacer("/", " ", "：", " ", ":", " ", "-", " ", "_", " ").Replace(value)
	return strings.Join(strings.Fields(value), " ")
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
