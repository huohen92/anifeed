package webui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/request"
)

var (
	anibtBrowseBase        = "https://anibt.net"
	animeGardenBrowseBase  = "https://api.animes.garden"
	mikanTrailingIDPattern = regexp.MustCompile(`(\d+)/?$`)
)

type browseItem struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Cover    string  `json:"cover,omitempty"`
	Score    float64 `json:"score,omitempty"`
	Meta     string  `json:"meta,omitempty"`
	RSSCount int     `json:"rss_count,omitempty"`
	Exists   bool    `json:"exists"`
}

type browseWeek struct {
	Label string       `json:"label"`
	Items []browseItem `json:"items"`
}

type browseListResponse struct {
	Source        string       `json:"source"`
	Title         string       `json:"title"`
	Seasons       []string     `json:"seasons"`
	CurrentSeason string       `json:"current_season,omitempty"`
	Weeks         []browseWeek `json:"weeks"`
	Total         int          `json:"total"`
}

type browseRelease struct {
	Title       string `json:"title"`
	Magnet      string `json:"magnet"`
	Size        int64  `json:"size"`
	PublishedAt string `json:"published_at,omitempty"`
}

type browseGroup struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	RSS    string          `json:"rss"`
	Latest string          `json:"latest,omitempty"`
	Items  []browseRelease `json:"items,omitempty"`
	Meta   string          `json:"meta,omitempty"`
}

type browseGroupsResponse struct {
	Source string        `json:"source"`
	Title  string        `json:"title"`
	Groups []browseGroup `json:"groups"`
}

func (s *Server) handleBrowseList(w http.ResponseWriter, r *http.Request) {
	source := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("source")))
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	season := strings.TrimSpace(r.URL.Query().Get("season"))

	var (
		resp *browseListResponse
		err  error
	)
	switch source {
	case "anibt":
		resp, err = s.browseAniBT(q, season)
	case "animegarden":
		resp, err = s.browseAnimeGarden(q)
	case "mikan":
		resp, err = s.browseMikan(q, season)
	default:
		writeAPIError(w, http.StatusBadRequest, "不支持的浏览来源: "+source)
		return
	}
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleBrowseGroups(w http.ResponseWriter, r *http.Request) {
	source := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("source")))
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeAPIError(w, http.StatusBadRequest, "缺少番剧 ID")
		return
	}

	var (
		resp *browseGroupsResponse
		err  error
	)
	switch source {
	case "anibt":
		resp, err = s.browseAniBTGroups(id)
	case "animegarden":
		resp, err = s.browseAnimeGardenGroups(id)
	case "mikan":
		resp, err = s.browseMikanGroups(id)
	default:
		writeAPIError(w, http.StatusBadRequest, "不支持的浏览来源: "+source)
		return
	}
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) feedExists(source, id string) bool {
	cfg := s.configSnapshot()
	if cfg == nil {
		return false
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if feed.Url == "" {
				continue
			}
			u, err := url.Parse(feed.Url)
			if err != nil {
				continue
			}
			values := u.Query()
			switch source {
			case "anibt":
				if values.Get("bgmId") == id {
					return true
				}
			case "animegarden":
				if values.Get("subject") == id {
					return true
				}
			case "mikan":
				if values.Get("bangumiId") == id || strings.Contains(u.Path, "/Bangumi/"+id) {
					return true
				}
			}
		}
	}
	return false
}

func fetchBrowseBytes(rawURL string) ([]byte, error) {
	resp, err := request.Request(http.MethodGet, rawURL, nil, map[string]string{"Accept": "*/*"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求 %s 失败: HTTP %d", rawURL, resp.StatusCode)
	}
	return body, nil
}

func fetchBrowseJSON(rawURL string, dst any) error {
	body, err := fetchBrowseBytes(rawURL)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("解析 %s 返回数据失败: %w", rawURL, err)
	}
	return nil
}

func absoluteBrowseURL(base, ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return ref
	}
	b, err := url.Parse(base)
	if err != nil {
		return ref
	}
	r, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return b.ResolveReference(r).String()
}

/* ---------------- AniBT ---------------- */

type anibtListPayload struct {
	Data struct {
		CurrentSeason    string   `json:"currentSeason"`
		RequestedSeason  string   `json:"requestedSeason"`
		AvailableSeasons []string `json:"availableSeasons"`
		ByWeekday        []struct {
			WeekdayLabel string `json:"weekdayLabel"`
			Animes       []struct {
				BgmID int    `json:"bgmId"`
				Cover string `json:"cover"`
				Title struct {
					Primary string `json:"primary"`
					Chinese string `json:"chinese"`
				} `json:"title"`
				Rating          float64 `json:"rating"`
				Format          string  `json:"format"`
				Episodes        int     `json:"episodes"`
				RSSReleaseCount int     `json:"rssReleaseCount"`
			} `json:"animes"`
		} `json:"byWeekday"`
	} `json:"data"`
}

func (s *Server) browseAniBT(q, season string) (*browseListResponse, error) {
	u, _ := url.Parse(anibtBrowseBase + "/api/seasons/anime")
	values := u.Query()
	if season != "" {
		values.Set("season", season)
	}
	if q != "" {
		values.Set("query", q)
	}
	u.RawQuery = values.Encode()

	var payload anibtListPayload
	if err := fetchBrowseJSON(u.String(), &payload); err != nil {
		return nil, err
	}

	resp := &browseListResponse{
		Source:        "anibt",
		Title:         "AniBT",
		Seasons:       payload.Data.AvailableSeasons,
		CurrentSeason: payload.Data.CurrentSeason,
		Weeks:         make([]browseWeek, 0, len(payload.Data.ByWeekday)),
	}
	for _, week := range payload.Data.ByWeekday {
		items := make([]browseItem, 0, len(week.Animes))
		for _, anime := range week.Animes {
			if q == "" && anime.RSSReleaseCount == 0 {
				continue
			}
			title := anime.Title.Chinese
			if title == "" {
				title = anime.Title.Primary
			}
			id := strconv.Itoa(anime.BgmID)
			meta := anime.Format
			if anime.Episodes > 0 {
				meta = fmt.Sprintf("%s · %d 话", meta, anime.Episodes)
			}
			items = append(items, browseItem{
				ID:       id,
				Title:    title,
				Cover:    anime.Cover,
				Score:    anime.Rating,
				Meta:     meta,
				RSSCount: anime.RSSReleaseCount,
				Exists:   s.feedExists("anibt", id),
			})
		}
		if len(items) == 0 {
			continue
		}
		sort.SliceStable(items, func(i, j int) bool { return items[i].Score > items[j].Score })
		resp.Weeks = append(resp.Weeks, browseWeek{Label: week.WeekdayLabel, Items: items})
		resp.Total += len(items)
	}
	if resp.CurrentSeason == "" {
		resp.CurrentSeason = payload.Data.RequestedSeason
	}
	return resp, nil
}

func (s *Server) browseAniBTGroups(bgmID string) (*browseGroupsResponse, error) {
	u := anibtBrowseBase + "/api/anime/groups?bgmId=" + url.QueryEscape(bgmID)
	var payload struct {
		Data struct {
			Title struct {
				Primary string `json:"primary"`
				Chinese string `json:"chinese"`
			} `json:"title"`
			Groups []struct {
				Slug          string `json:"slug"`
				Name          string `json:"name"`
				LastUpdatedAt int64  `json:"lastUpdatedAt"`
				Items         []struct {
					Title       string `json:"title"`
					Magnet      string `json:"magnet"`
					Size        int64  `json:"size"`
					PublishedAt int64  `json:"publishedAt"`
				} `json:"items"`
			} `json:"groups"`
		} `json:"data"`
	}
	if err := fetchBrowseJSON(u, &payload); err != nil {
		return nil, err
	}

	resp := &browseGroupsResponse{Source: "anibt", Title: payload.Data.Title.Chinese}
	if resp.Title == "" {
		resp.Title = payload.Data.Title.Primary
	}
	for _, group := range payload.Data.Groups {
		rss := fmt.Sprintf("%s/rss/anime.xml?bgmId=%s&groupSlug=%s", anibtBrowseBase, url.QueryEscape(bgmID), url.QueryEscape(group.Slug))
		items := make([]browseRelease, 0, len(group.Items))
		for i, item := range group.Items {
			if i >= 20 {
				break
			}
			items = append(items, browseRelease{
				Title:       item.Title,
				Magnet:      item.Magnet,
				Size:        item.Size,
				PublishedAt: millisToRFC3339(item.PublishedAt),
			})
		}
		resp.Groups = append(resp.Groups, browseGroup{
			ID:     group.Slug,
			Name:   group.Name,
			RSS:    rss,
			Latest: millisToRFC3339(group.LastUpdatedAt),
			Items:  items,
			Meta:   fmt.Sprintf("%d 个资源", len(group.Items)),
		})
	}
	sort.SliceStable(resp.Groups, func(i, j int) bool { return resp.Groups[i].Latest > resp.Groups[j].Latest })
	return resp, nil
}

func millisToRFC3339(millis int64) string {
	if millis <= 0 {
		return ""
	}
	return time.UnixMilli(millis).Format(time.RFC3339)
}

/* ---------------- AnimeGarden ---------------- */

func (s *Server) browseAnimeGarden(q string) (*browseListResponse, error) {
	var payload struct {
		Subjects []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Actived  string `json:"activedAt"`
			Archived bool   `json:"isArchived"`
			Search   struct {
				Include []string `json:"include"`
			} `json:"search"`
		} `json:"subjects"`
	}
	if err := fetchBrowseJSON(animeGardenBrowseBase+"/subjects", &payload); err != nil {
		return nil, err
	}

	query := strings.ToLower(q)
	weeks := map[string][]browseItem{}
	for _, subject := range payload.Subjects {
		if subject.Archived {
			continue
		}
		if query != "" {
			haystack := strings.ToLower(subject.Name + " " + strings.Join(subject.Search.Include, " "))
			if !strings.Contains(haystack, query) {
				continue
			}
		}
		actived, err := time.Parse(time.RFC3339, subject.Actived)
		if err != nil {
			actived = time.Now()
		}
		label := weekdayLabel(actived.Weekday())
		id := strconv.Itoa(subject.ID)
		weeks[label] = append(weeks[label], browseItem{
			ID:     id,
			Title:  subject.Name,
			Meta:   "AnimeGarden",
			Exists: s.feedExists("animegarden", id),
		})
	}

	resp := &browseListResponse{Source: "animegarden", Title: "AnimeGarden"}
	for _, label := range weekdayOrder {
		items := weeks[label]
		if len(items) == 0 {
			continue
		}
		sort.SliceStable(items, func(i, j int) bool { return items[i].Title < items[j].Title })
		resp.Weeks = append(resp.Weeks, browseWeek{Label: label, Items: items})
		resp.Total += len(items)
	}
	return resp, nil
}

func (s *Server) browseAnimeGardenGroups(subjectID string) (*browseGroupsResponse, error) {
	u, _ := url.Parse(animeGardenBrowseBase + "/resources")
	values := u.Query()
	values.Set("subject", subjectID)
	values.Set("pageSize", "200")
	values.Set("duplicate", "false")
	u.RawQuery = values.Encode()

	var payload struct {
		Resources []struct {
			Title     string `json:"title"`
			Magnet    string `json:"magnet"`
			Size      int64  `json:"size"`
			CreatedAt string `json:"createdAt"`
			Fansub    *struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"fansub"`
		} `json:"resources"`
	}
	if err := fetchBrowseJSON(u.String(), &payload); err != nil {
		return nil, err
	}

	type groupAcc struct {
		group browseGroup
		count int
	}
	groups := map[string]*groupAcc{}
	for _, resource := range payload.Resources {
		if resource.Fansub == nil || resource.Magnet == "" {
			continue
		}
		key := strconv.Itoa(resource.Fansub.ID)
		acc, ok := groups[key]
		if !ok {
			rss := fmt.Sprintf("%s/feed.xml?subject=%s&fansub=%s", animeGardenBrowseBase, url.QueryEscape(subjectID), url.QueryEscape(resource.Fansub.Name))
			acc = &groupAcc{group: browseGroup{
				ID:   key,
				Name: resource.Fansub.Name,
				RSS:  rss,
			}}
			groups[key] = acc
		}
		acc.count++
		if resource.CreatedAt > acc.group.Latest {
			acc.group.Latest = resource.CreatedAt
		}
		if len(acc.group.Items) < 20 {
			acc.group.Items = append(acc.group.Items, browseRelease{
				Title:       resource.Title,
				Magnet:      resource.Magnet,
				Size:        resource.Size,
				PublishedAt: resource.CreatedAt,
			})
		}
	}

	resp := &browseGroupsResponse{Source: "animegarden"}
	for _, acc := range groups {
		acc.group.Meta = fmt.Sprintf("%d 个资源", acc.count)
		resp.Groups = append(resp.Groups, acc.group)
	}
	sort.SliceStable(resp.Groups, func(i, j int) bool { return resp.Groups[i].Latest > resp.Groups[j].Latest })
	return resp, nil
}

var weekdayOrder = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}

func weekdayLabel(day time.Weekday) string {
	switch day {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	default:
		return "周日"
	}
}

/* ---------------- Mikan ---------------- */

func (s *Server) browseMikan(q, season string) (*browseListResponse, error) {
	bases := mikanBrowseBases(s.configSnapshot())
	var lastErr error
	for _, base := range bases {
		resp, err := s.browseMikanWithBase(base, q, season)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (s *Server) browseMikanWithBase(base, q, season string) (*browseListResponse, error) {
	resp := &browseListResponse{Source: "mikan", Title: "Mikan", Weeks: []browseWeek{}}

	if strings.HasPrefix(strings.ToLower(q), "id:") {
		id := strings.TrimSpace(q[3:])
		if id == "" {
			return nil, fmt.Errorf("请输入番剧 ID")
		}
		item, err := s.mikanAnimeInfo(base, id)
		if err != nil {
			return nil, err
		}
		resp.Weeks = append(resp.Weeks, browseWeek{Label: "搜索结果", Items: []browseItem{item}})
		resp.Total = 1
		return resp, nil
	}

	rawURL := base + "/Home/Search?searchstr=" + url.QueryEscape(q)
	if q == "" {
		year, seasonName := seasonParts(season)
		rawURL = fmt.Sprintf("%s/Home/BangumiCoverFlowByDayOfWeek?year=%d&seasonStr=%s", base, year, url.QueryEscape(seasonName))
	}
	body, err := fetchBrowseBytes(rawURL)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("解析 Mikan 页面失败: %w", err)
	}

	resp.Seasons, resp.CurrentSeason = parseMikanSeasons(doc, base)
	if resp.CurrentSeason == "" {
		year, seasonName := seasonParts(season)
		resp.CurrentSeason = fmt.Sprintf("%d%s", year, seasonName)
	}

	weekBlocks := doc.Find(".sk-bangumi")
	if weekBlocks.Length() > 0 {
		weekBlocks.Each(func(_ int, block *goquery.Selection) {
			label := strings.TrimSpace(block.Children().First().Text())
			if label == "" {
				label = "番剧"
			}
			items := parseMikanItems(block, base, s)
			if len(items) > 0 {
				resp.Weeks = append(resp.Weeks, browseWeek{Label: label, Items: items})
				resp.Total += len(items)
			}
		})
	} else {
		items := parseMikanItems(doc.Selection, base, s)
		if len(items) > 0 {
			resp.Weeks = append(resp.Weeks, browseWeek{Label: "搜索结果", Items: items})
			resp.Total = len(items)
		}
	}
	return resp, nil
}

func (s *Server) mikanAnimeInfo(base, id string) (browseItem, error) {
	pageURL := base + "/Home/Bangumi/" + url.PathEscape(id)
	body, err := fetchBrowseBytes(pageURL)
	if err != nil {
		return browseItem{}, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return browseItem{}, err
	}
	title := strings.TrimSpace(doc.Find(".bangumi-title").First().Text())
	cover, _ := doc.Find(".content > img").First().Attr("src")
	return browseItem{
		ID:     id,
		Title:  title,
		Cover:  absoluteBrowseURL(base, cover),
		Exists: s.feedExists("mikan", id),
	}, nil
}

func parseMikanItems(root *goquery.Selection, base string, s *Server) []browseItem {
	items := make([]browseItem, 0)
	root.Find(".an-ul li").Each(func(_ int, li *goquery.Selection) {
		anchor := li.Find("a").First()
		href, _ := anchor.Attr("href")
		title := strings.TrimSpace(anchor.Text())
		cover, _ := li.Find("span[data-src]").First().Attr("data-src")
		id := ""
		if match := mikanTrailingIDPattern.FindStringSubmatch(href); len(match) == 2 {
			id = match[1]
		}
		if title == "" || href == "" {
			return
		}
		items = append(items, browseItem{
			ID:     id,
			Title:  title,
			Cover:  absoluteBrowseURL(base, cover),
			Exists: s.feedExists("mikan", id),
		})
	})
	return items
}

func parseMikanSeasons(doc *goquery.Document, base string) ([]string, string) {
	seasons := make([]string, 0)
	current := ""
	dateText := strings.TrimSpace(doc.Find(".date-select .date-text").First().Text())
	doc.Find(".date-select .dropdown-menu a[data-year][data-season]").Each(func(_ int, a *goquery.Selection) {
		year, _ := a.Attr("data-year")
		season, _ := a.Attr("data-season")
		label := year + season
		if year == "" || season == "" {
			return
		}
		seasons = append(seasons, label)
		if dateText != "" && strings.Contains(dateText, label) {
			current = label
		}
	})
	return seasons, current
}

func (s *Server) browseMikanGroups(idOrURL string) (*browseGroupsResponse, error) {
	bases := mikanBrowseBases(s.configSnapshot())
	var lastErr error
	for _, base := range bases {
		resp, err := s.browseMikanGroupsWithBase(base, idOrURL)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (s *Server) browseMikanGroupsWithBase(base, idOrURL string) (*browseGroupsResponse, error) {
	pageURL := idOrURL
	if !strings.HasPrefix(pageURL, "http://") && !strings.HasPrefix(pageURL, "https://") {
		pageURL = base + "/Home/Bangumi/" + url.PathEscape(idOrURL)
	}
	body, err := fetchBrowseBytes(pageURL)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	resp := &browseGroupsResponse{Source: "mikan", Title: strings.TrimSpace(doc.Find(".bangumi-title").First().Text())}

	doc.Find(".leftbar-item").Each(func(_ int, block *goquery.Selection) {
		link := block.Find("a.subgroup-name").First()
		name := strings.TrimSpace(link.Text())
		anchor, _ := link.Attr("data-anchor")
		if name == "" || anchor == "" {
			return
		}
		rss := ""
		if rssLink := doc.Find(anchor + " .mikan-rss").First(); rssLink.Length() > 0 {
			href, _ := rssLink.Attr("href")
			rss = absoluteBrowseURL(base, href)
		}
		group := browseGroup{
			ID:     strings.TrimPrefix(anchor, "#"),
			Name:   name,
			RSS:    rss,
			Latest: strings.TrimSpace(block.Find(".date").First().Text()),
		}
		table := doc.Find(anchor).First().Next()
		table.Find("tbody tr").Each(func(index int, tr *goquery.Selection) {
			if index >= 20 {
				return
			}
			links := tr.Find("a")
			title := strings.TrimSpace(links.First().Text())
			magnet, _ := links.Eq(1).Attr("data-clipboard-text")
			if title == "" || magnet == "" {
				return
			}
			size := int64(0)
			tDs := tr.Find("td")
			if tDs.Length() > 2 {
				size = parseSizeText(tDs.Eq(2).Text())
			}
			published := ""
			if tDs.Length() > 3 {
				published = strings.TrimSpace(tDs.Eq(3).Text())
			}
			group.Items = append(group.Items, browseRelease{
				Title:       title,
				Magnet:      magnet,
				Size:        size,
				PublishedAt: published,
			})
		})
		group.Meta = fmt.Sprintf("%d 个资源", len(group.Items))
		resp.Groups = append(resp.Groups, group)
	})
	return resp, nil
}

func mikanBrowseBases(cfg *config.Config) []string {
	base := configMikanHost(cfg)
	bases := []string{base}
	if strings.Contains(base, "mikanani.me") {
		bases = append(bases, "https://mikanime.tv")
	}
	return bases
}

func configMikanHost(cfg *config.Config) string {
	host := "https://mikanani.me"
	if cfg != nil && strings.TrimSpace(cfg.Subscription.MikanHost) != "" {
		host = strings.TrimSpace(cfg.Subscription.MikanHost)
	}
	return strings.TrimRight(host, "/")
}

func seasonParts(season string) (int, string) {
	now := time.Now()
	year, name := now.Year(), mikanSeasonName(now.Month())
	season = strings.TrimSpace(season)
	if season == "" {
		return year, name
	}
	digits := ""
	for _, r := range season {
		if r >= '0' && r <= '9' {
			digits += string(r)
			continue
		}
		if strings.ContainsRune("冬春夏秋", r) {
			name = string(r)
		}
	}
	if len(digits) >= 4 {
		if parsed, err := strconv.Atoi(digits[:4]); err == nil {
			year = parsed
		}
	}
	return year, name
}

func mikanSeasonName(month time.Month) string {
	switch {
	case month >= 1 && month <= 3:
		return "冬"
	case month >= 4 && month <= 6:
		return "春"
	case month >= 7 && month <= 9:
		return "夏"
	default:
		return "秋"
	}
}

var sizeTextPattern = regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)\s*(B|KB|MB|GB|TB|KIB|MIB|GIB|TIB)?$`)

func parseSizeText(text string) int64 {
	match := sizeTextPattern.FindStringSubmatch(strings.TrimSpace(text))
	if len(match) == 0 {
		return 0
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0
	}
	unit := strings.ToUpper(match[2])
	multipliers := map[string]float64{
		"": 1, "B": 1, "KB": 1024, "MB": 1024 * 1024, "GB": 1024 * 1024 * 1024,
		"TB":  1024 * 1024 * 1024 * 1024,
		"KIB": 1024, "MIB": 1024 * 1024, "GIB": 1024 * 1024 * 1024,
		"TIB": 1024 * 1024 * 1024 * 1024,
	}
	if multiplier, ok := multipliers[unit]; ok {
		return int64(value * multiplier)
	}
	return 0
}
