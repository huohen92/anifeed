package webui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/zhifengle/anifeed/cloudfs"
	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/store"
)

type washFileView struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	Subgroup   string `json:"subgroup,omitempty"`
	Simplified bool   `json:"simplified"`
	HasRecord  bool   `json:"has_record"`
	Reason     string `json:"reason,omitempty"`
}

type washGroupView struct {
	Series  string         `json:"series"`
	Season  int            `json:"season"`
	Episode float64        `json:"episode"`
	Keep    washFileView   `json:"keep"`
	Delete  []washFileView `json:"delete"`
}

type washResponse struct {
	Mode        string          `json:"mode"`
	Folders     int             `json:"folders"`
	Scanned     int             `json:"scanned"`
	Groups      []washGroupView `json:"groups"`
	DeleteCount int             `json:"delete_count"`
	Deleted     int             `json:"deleted,omitempty"`
	Failed      int             `json:"failed,omitempty"`
	DryRun      bool            `json:"dry_run"`
}

func (s *Server) handleWash(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusInternalServerError, "配置尚未加载")
		return
	}
	var payload struct {
		DryRun bool   `json:"dry_run"`
		Mode   string `json:"mode"`
	}
	if err := decodeJSON(w, r, &payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}
	if !cfg.Lab.WashEnabled {
		writeAPIError(w, http.StatusBadRequest, "请先在实验室启用洗板")
		return
	}
	if len(cfg.Lab.WashFolders) == 0 {
		writeAPIError(w, http.StatusBadRequest, "请先添加洗板文件夹")
		return
	}
	mode := payload.Mode
	if mode != "full" {
		mode = "safe"
	}
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Minute)
	defer cancel()
	s.washRunMu.Lock()
	defer s.washRunMu.Unlock()
	resp, err := s.performWash(ctx, cfg, mode, payload.DryRun)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) performWash(ctx context.Context, cfg *config.Config, mode string, dryRun bool) (washResponse, error) {
	ag, err := s.ensureAgent()
	if err != nil {
		return washResponse{}, err
	}
	fs := ag.FileSystem()
	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		return washResponse{}, err
	}
	defer st.Close()
	finished, _ := st.FinishedEntries()
	records := make(map[string]store.FinishedEntry, len(finished))
	for _, entry := range finished {
		records[entry.EntryID] = entry
	}
	candidates := make([]washCandidate, 0, 128)
	scanned := 0
	for _, cid := range cfg.Lab.WashFolders {
		if strings.TrimSpace(cid) == "" {
			continue
		}
		rootName := cid
		if entry, statErr := fs.Stat(ctx, cid); statErr == nil && strings.TrimSpace(entry.Name) != "" {
			rootName = entry.Name
		}
		scanWashDir(ctx, fs, cid, rootName, mode, records, nil, &candidates, &scanned, 0)
	}
	groups := groupWashCandidates(candidates, cfg.Lab.PreferredSubgroups)
	resp := washResponse{Mode: mode, Folders: len(cfg.Lab.WashFolders), Scanned: scanned, Groups: groups, DryRun: dryRun}
	for _, group := range groups {
		resp.DeleteCount += len(group.Delete)
	}
	if dryRun {
		return resp, nil
	}
	deletedFiles := make([]washFileView, 0, resp.DeleteCount)
	failedFiles := make([]string, 0)
	for _, group := range groups {
		for _, file := range group.Delete {
			if err := fs.Delete(ctx, file.ID); err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("%s：%v", file.Path, err))
				continue
			}
			deletedFiles = append(deletedFiles, file)
		}
	}
	resp.Deleted = len(deletedFiles)
	resp.Failed = len(failedFiles)
	for _, group := range groups {
		for _, file := range group.Delete {
			reason := file.Reason
			if reason == "" {
				reason = "同集已有更优先版本"
			}
			log.Printf("[wash] 保留: %s | 删除: %s | 原因: %s", group.Keep.Path, file.Path, reason)
		}
	}
	for _, message := range failedFiles {
		log.Printf("[wash] 删除失败: %s", message)
	}
	s.sendWashNotification(cfg, mode, groups, deletedFiles, failedFiles)
	return resp, nil
}

type washCandidate struct {
	washFileView
	Series    string
	Season    int
	Episode   float64
	hasRecord bool
}

func scanWashDir(ctx context.Context, fs cloudfs.Driver, dirID, dirPath, mode string, records map[string]store.FinishedEntry, inherited *store.FinishedEntry, out *[]washCandidate, scanned *int, depth int) {
	if depth > 8 || len(*out) > 10000 {
		return
	}
	entries, err := fs.List(ctx, dirID)
	if err != nil {
		return
	}
	for _, entry := range entries {
		fullPath := path.Join(dirPath, entry.Name)
		if entry.IsDir() {
			nextRecord := inherited
			if record, ok := records[entry.ID]; ok {
				copy := record
				nextRecord = &copy
			}
			scanWashDir(ctx, fs, entry.ID, fullPath, mode, records, nextRecord, out, scanned, depth+1)
			continue
		}
		*scanned++
		if !p115.LooksLikeMediaFile(entry.Name) {
			continue
		}
		record, hasRecord := records[entry.ID]
		if !hasRecord && inherited != nil {
			record = *inherited
			hasRecord = true
		}
		if mode == "safe" && !hasRecord {
			continue
		}
		series, season, episode := washParseName(entry.Name)
		if hasRecord && record.ConfigName != "" {
			series = record.ConfigName
			if season == 0 {
				season = detectSeason(record.ConfigName)
			}
			if episode <= 0 {
				if parsed, err := p115.ExtractEpisode(record.Title); err == nil {
					episode = parsed
				}
			}
		}
		if episode <= 0 {
			continue
		}
		if season <= 0 {
			season = 1
		}
		subgroup := p115.ExtractSubgroup(entry.Name)
		if subgroup == "" || subgroup == "未知字幕组" {
			subgroup = washTrailingSubgroup(entry.Name)
		}
		*out = append(*out, washCandidate{
			washFileView: washFileView{ID: entry.ID, Name: entry.Name, Path: fullPath, Size: entry.Size, Subgroup: subgroup, Simplified: washIsSimplified(entry.Name), HasRecord: hasRecord},
			Series:       series, Season: season, Episode: episode, hasRecord: hasRecord,
		})
	}
}

func washParseName(name string) (string, int, float64) {
	base := strings.TrimSuffix(name, path.Ext(name))
	episode, err := p115.ExtractEpisode(base)
	if err != nil {
		return "", 0, 0
	}
	season := detectSeason(base)
	if season == 0 {
		season = 1
	}
	series := base
	// Trim at the most common season/episode markers.
	upper := strings.ToUpper(series)
	if index := strings.Index(upper, " S"); index > 0 {
		series = series[:index]
	} else if index := strings.Index(upper, ".S"); index > 0 {
		series = series[:index]
	} else if index := strings.Index(series, " - "); index > 0 {
		series = series[:index]
	} else if index := strings.Index(series, "["); index > 0 {
		series = series[:index]
	}
	series = strings.Trim(series, " ._-[]()（）")
	if series == "" {
		series = base
	}
	return series, season, episode
}

func washTrailingSubgroup(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	for _, sep := range []string{"-", "_", " "} {
		if index := strings.LastIndex(base, sep); index > 0 {
			value := strings.TrimSpace(base[index+1:])
			if value != "" && len([]rune(value)) <= 30 && !strings.ContainsAny(value, "[]()（）") {
				return value
			}
		}
	}
	return ""
}

func washIsSimplified(name string) bool {
	lower := strings.ToLower(name)
	simplified := strings.Contains(lower, "简") || strings.Contains(lower, "chs") || strings.Contains(lower, "gb") || strings.Contains(lower, "sc")
	traditional := strings.Contains(lower, "繁") || strings.Contains(lower, "cht") || strings.Contains(lower, "big5") || strings.Contains(lower, "tc")
	return simplified && !traditional
}

func groupWashCandidates(candidates []washCandidate, preferred []string) []washGroupView {
	grouped := make(map[string][]washCandidate)
	for _, candidate := range candidates {
		key := fmt.Sprintf("%s|%d|%g", strings.ToLower(candidate.Series), candidate.Season, candidate.Episode)
		grouped[key] = append(grouped[key], candidate)
	}
	var groups []washGroupView
	for _, list := range grouped {
		if len(list) < 2 {
			continue
		}
		sort.SliceStable(list, func(i, j int) bool {
			return washCandidateBetter(list[i], list[j], preferred)
		})
		keep := list[0]
		view := washGroupView{Series: keep.Series, Season: keep.Season, Episode: keep.Episode, Keep: keep.washFileView}
		for _, candidate := range list[1:] {
			candidate.Reason = washDeleteReason(keep, candidate, preferred)
			view.Delete = append(view.Delete, candidate.washFileView)
		}
		groups = append(groups, view)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].Series != groups[j].Series {
			return groups[i].Series < groups[j].Series
		}
		if groups[i].Season != groups[j].Season {
			return groups[i].Season < groups[j].Season
		}
		return groups[i].Episode < groups[j].Episode
	})
	return groups
}

func preferredRank(candidate washCandidate, preferred []string) int {
	name := strings.ToLower(candidate.Name + " " + candidate.Subgroup)
	for index, subgroup := range preferred {
		subgroup = strings.TrimSpace(strings.ToLower(subgroup))
		if subgroup != "" && strings.Contains(name, subgroup) {
			return index
		}
	}
	return 9999
}

func washCandidateBetter(left, right washCandidate, preferred []string) bool {
	leftRank, rightRank := preferredRank(left, preferred), preferredRank(right, preferred)
	if leftRank != rightRank {
		return leftRank < rightRank
	}
	if left.Simplified != right.Simplified {
		return left.Simplified
	}
	if left.Size != right.Size {
		return left.Size > right.Size
	}
	return left.Path < right.Path
}

func washDeleteReason(keep, deleted washCandidate, preferred []string) string {
	keepRank, deleteRank := preferredRank(keep, preferred), preferredRank(deleted, preferred)
	if keepRank != 9999 && deleteRank != 9999 {
		return "同集已有更优先的喜好字幕组 " + keep.Subgroup
	}
	if keep.Simplified != deleted.Simplified {
		return "同集已有简体版本"
	}
	if keep.Size > deleted.Size {
		return "同集已有体积更大的版本"
	}
	return "同集已有更优先的版本"
}

func (s *Server) sendWashNotification(cfg *config.Config, mode string, groups []washGroupView, deleted []washFileView, failed []string) {
	if cfg == nil || !notifyEventEnabled(cfg.Notify, "wash") || strings.TrimSpace(cfg.Notify.Webhook) == "" {
		return
	}
	modeText := "安全模式（只处理下载记录）"
	if mode == "full" {
		modeText = "完整模式（扫描文件夹）"
	}
	lines := []string{
		"AniFeed 洗板完成",
		"时间: " + time.Now().Format("2006-01-02 15:04"),
		"模式: " + modeText,
		fmt.Sprintf("发现重复组: %d，删除文件: %d，失败: %d", len(groups), len(deleted), len(failed)),
		"",
		"删除:",
	}
	for _, file := range deleted {
		line := "- " + file.Path
		if file.Subgroup != "" {
			line += "（" + file.Subgroup + "）"
		}
		if file.Reason != "" {
			line += "，" + file.Reason
		}
		lines = append(lines, line)
	}
	if len(deleted) == 0 {
		lines = append(lines, "- 无")
	}
	if len(failed) > 0 {
		lines = append(lines, "", "失败:")
		for _, message := range failed {
			lines = append(lines, "- "+message)
		}
	}
	if err := postRawNotification(cfg.Notify, strings.Join(lines, "\n")); err != nil {
		logNotifyError(err)
	}
}

func (s *Server) startWashScheduleLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	s.washScheduleMu.Lock()
	if s.lastScheduledWash.IsZero() {
		s.lastScheduledWash = time.Now()
	}
	s.washScheduleMu.Unlock()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			cfg := s.configSnapshot()
			if cfg == nil || !cfg.Lab.WashEnabled || !cfg.Lab.WashScheduleEnabled || len(cfg.Lab.WashFolders) == 0 {
				continue
			}
			interval := cfg.Lab.WashIntervalMinutes
			if interval <= 0 {
				interval = 1440
			}
			s.washScheduleMu.Lock()
			if s.washScheduledRunning || now.Sub(s.lastScheduledWash) < time.Duration(interval)*time.Minute {
				s.washScheduleMu.Unlock()
				continue
			}
			s.washScheduledRunning = true
			s.lastScheduledWash = now
			s.washScheduleMu.Unlock()
			go s.runScheduledWash(cfg, now)
		}
	}
}

func (s *Server) runScheduledWash(cfg *config.Config, started time.Time) {
	defer func() {
		s.washScheduleMu.Lock()
		s.washScheduledRunning = false
		s.washScheduleMu.Unlock()
	}()
	s.washRunMu.Lock()
	defer s.washRunMu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	mode := cfg.Lab.WashMode
	if mode != "full" {
		mode = "safe"
	}
	resp, err := s.performWash(ctx, cfg, mode, false)
	if err != nil {
		logNotifyError(fmt.Errorf("定时洗板失败: %w", err))
		return
	}
	log.Printf("[wash] 定时洗板完成，删除 %d 个文件，失败 %d 个", resp.Deleted, resp.Failed)
	_ = started
}
