package webui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/deadblue/elevengo"
	"github.com/zhifengle/anifeed/cloudfs"
	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/p115"
	"github.com/zhifengle/anifeed/store"
)

var (
	chineseSeasonPattern = regexp.MustCompile(`第\s*([一二三四五六七八九十两0-9]+)\s*[季期]`)
	englishSeasonPattern = regexp.MustCompile(`(?i)(?:Season|S)\s*(\d{1,2})`)
	shortSeasonPattern   = regexp.MustCompile(`(?i)(\d{1,2})(?:st|nd|rd|th)\s+Season`)
	chineseSeasons       = map[string]int{"一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9, "十": 10}
)

func postProcessEnabled(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	if cfg.Rename.Enabled || cfg.Offline.DualFolder {
		return true
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if strings.TrimSpace(feed.FinalCid) != "" || feed.FinalMode > 0 || feed.FinalCustomMode > 0 {
				return true
			}
		}
	}
	return false
}

// startPostProcessLoop polls completed 115 offline tasks and renames/moves
// them. It is intentionally independent from the RSS scheduler.
func (s *Server) startPostProcessLoop(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(8 * time.Second):
		}
		s.runPostProcess()
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.runPostProcess()
			}
		}
	}()
}

func renameSeasonForFeed(cfg *config.Config, name, feedURL string) int {
	return renameSeasonForFeedTitle(cfg, name, feedURL, "")
}

func renameSeasonForFeedTitle(cfg *config.Config, name, feedURL, title string) int {
	season := 1
	if cfg != nil && cfg.Rename.Season > 0 {
		season = cfg.Rename.Season
	}
	if cfg == nil {
		return season
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if (feedURL != "" && feed.Url == feedURL) || (name != "" && feed.Name == name) {
				if feed.RenameSeason > 0 {
					return feed.RenameSeason
				}
				if detected := detectSeason(feed.Name); detected > 0 {
					return detected
				}
				if detected := detectSeason(title); detected > 0 {
					return detected
				}
				return season
			}
		}
	}
	if detected := detectSeason(name); detected > 0 {
		return detected
	}
	if detected := detectSeason(title); detected > 0 {
		return detected
	}
	return season
}

func detectSeason(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	if match := chineseSeasonPattern.FindStringSubmatch(text); len(match) == 2 {
		if value, ok := chineseSeasons[match[1]]; ok {
			return value
		}
		if value, err := strconv.Atoi(match[1]); err == nil && value > 0 {
			return value
		}
	}
	if match := englishSeasonPattern.FindStringSubmatch(text); len(match) == 2 {
		if value, err := strconv.Atoi(match[1]); err == nil && value > 0 {
			return value
		}
	}
	if match := shortSeasonPattern.FindStringSubmatch(text); len(match) == 2 {
		if value, err := strconv.Atoi(match[1]); err == nil && value > 0 {
			return value
		}
	}
	return 0
}

func renameOffsetForFeed(cfg *config.Config, feed *config.RssConfig) int {
	if feed != nil && feed.RenameOffset != 0 {
		return feed.RenameOffset
	}
	if cfg != nil {
		return cfg.Rename.Offset
	}
	return 0
}

func renameOffsetForFeedConfig(cfg *config.Config, name string) int {
	if cfg != nil {
		if feed, ok := findFeedConfig(cfg, name, ""); ok {
			return renameOffsetForFeed(cfg, &feed)
		}
	}
	if cfg != nil {
		return cfg.Rename.Offset
	}
	return 0
}

func (s *Server) runPostProcess() {
	s.postProcessMu.Lock()
	defer s.postProcessMu.Unlock()
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Printf("[warning] 离线后处理异常: %v", recovered)
		}
	}()

	cfg := s.configSnapshot()
	if !postProcessEnabled(cfg) {
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		return
	}
	pending, err := ag.StoreInstance.ListPendingPostTasks(50)
	if err != nil {
		log.Printf("[warning] 读取待处理离线任务失败: %v", err)
		return
	}
	if len(pending) == 0 {
		return
	}
	tasks, err := ag.ListOfflineTasks(0, nil)
	if err != nil {
		log.Printf("[warning] 获取 115 离线任务失败: %v", err)
		return
	}
	byHash := make(map[string]*elevengo.OfflineTask, len(tasks))
	for i := range tasks {
		byHash[tasks[i].InfoHash] = &tasks[i]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for _, item := range pending {
		task, ok := byHash[item.TaskHash]
		if !ok {
			if time.Since(item.CreatedAt) > 24*time.Hour {
				_ = ag.StoreInstance.FailPostTask(item.TaskHash, "115 中找不到该离线任务，可能已被清理", 1)
				go s.sendItemNotification(cfg, item.ConfigName, item.SourceURL, item.Title, "下载失败", "", false)
			}
			continue
		}
		if task.IsFailed() {
			_ = ag.StoreInstance.FailPostTask(item.TaskHash, "115 离线任务失败", 1)
			go s.sendItemNotification(cfg, item.ConfigName, item.SourceURL, item.Title, "下载失败", "", false)
			continue
		}
		if !task.IsDone() {
			timeout := cfg.Notify.DownloadTimeoutMinutes
			if timeout <= 0 {
				timeout = 180
			}
			if time.Since(item.CreatedAt) > time.Duration(timeout)*time.Minute {
				reason := fmt.Sprintf("超过 %d 分钟未完成", timeout)
				_ = ag.StoreInstance.FailPostTask(item.TaskHash, reason, 1)
				go s.sendItemNotification(cfg, item.ConfigName, item.SourceURL, item.Title, "下载失败", "", false)
			}
			continue
		}
		firstAir := ""
		if feed := feedForName(cfg, item.ConfigName); feed != nil {
			firstAir = s.loadNotifyMeta(cfg, feed).FirstAir
		}
		entries, err := processPostTask(ctx, ag.FileSystem(), task, item, cfg, firstAir)
		if err != nil {
			log.Printf("[postprocess] %s: %v", item.Title, err)
			_ = ag.StoreInstance.FailPostTask(item.TaskHash, err.Error(), 5)
			continue
		}
		if err := ag.StoreInstance.CompletePostTaskWithEntries(item.TaskHash, entries); err != nil {
			log.Printf("[warning] 标记离线后处理完成失败: %v", err)
			continue
		}
		log.Printf("[postprocess] 完成: %s", item.Title)
		go s.sendItemNotification(cfg, item.ConfigName, item.SourceURL, item.Title, "下载完成", postEventName(entries, predictedRenameName(cfg, feedForName(cfg, item.ConfigName), item.Title)), true)
	}
}

type postProcessFS interface {
	Stat(ctx context.Context, entryID string) (cloudfs.Entry, error)
	Lookup(ctx context.Context, parentID, name string) (cloudfs.Entry, error)
	List(ctx context.Context, dirID string) ([]cloudfs.Entry, error)
	Rename(ctx context.Context, entryID, newName string) (cloudfs.Entry, error)
	Move(ctx context.Context, targetDirID, entryID string) (cloudfs.Entry, error)
	Delete(ctx context.Context, entryID string) error
}

func processPostTask(ctx context.Context, fs postProcessFS, task *elevengo.OfflineTask, pending store.PostTask, cfg *config.Config, firstAir string) (finalEntries []store.PostEntry, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("115 文件操作异常: %v", recovered)
		}
	}()
	if fs == nil {
		return nil, fmt.Errorf("115 文件系统不可用")
	}

	entry, locateErr := locateTaskEntry(ctx, fs, task, pending, cfg)
	if locateErr != nil {
		return nil, locateErr
	}
	feedFinalCid := ""
	feedFinalMode := 0
	feedCustomMode := 0
	if feed, ok := findFeedConfig(cfg, pending.ConfigName, ""); ok {
		feedFinalCid = strings.TrimSpace(feed.FinalCid)
		if feed.FinalMode == 2 {
			feedFinalMode = 2
		} else if feed.FinalMode == 1 {
			feedFinalMode = 1
		}
		if feed.FinalCustomMode == 2 {
			feedCustomMode = 2
		}
	}

	var media mediaFiles
	if cfg.Rename.Enabled || ((cfg.Offline.DualFolder || feedFinalCid != "" || feedFinalMode > 0) && cfg.Offline.KeepOnlyVideo) {
		media, err = findMediaFiles(ctx, fs, entry)
		if err != nil {
			return nil, err
		}
	}
	if cfg.Rename.Enabled {
		tmdbID := ""
		if feed, ok := findFeedConfig(cfg, pending.ConfigName, ""); ok {
			tmdbID = strings.TrimSpace(feed.TmdbID)
		}
		season := renameSeasonForFeedTitle(cfg, pending.ConfigName, "", pending.Title)
		offset := renameOffsetForFeedConfig(cfg, pending.ConfigName)
		if pending.RenameSeason > 0 {
			season = pending.RenameSeason
		}
		if pending.RenameSeason > 0 || pending.RenameOffset != 0 {
			offset = pending.RenameOffset
		}
		newName, buildErr := p115.BuildRenameName(pending.ConfigName, pending.Title, p115.RenameOptions{
			Enabled:  true,
			Template: cfg.Rename.Template,
			Season:   season,
			Offset:   offset,
			TmdbID:   tmdbID,
			Subgroup: sourceSubgroupName(cfg, pending.ConfigName, pending.SourceURL, pending.Title),
		})
		if buildErr != nil {
			return nil, fmt.Errorf("生成重命名失败: %w", buildErr)
		}
		if entry.IsDir() {
			if len(media.videos) == 0 {
				return nil, fmt.Errorf("离线目录中没有找到视频文件")
			}
			for index := range media.videos {
				target := p115.PreserveFileExtension(newName, media.videos[index].Name, false)
				if index > 0 {
					target = appendNameIndex(target, index+1)
				}
				if target == "" {
					return nil, fmt.Errorf("重命名结果为空")
				}
				if target == media.videos[index].Name {
					continue
				}
				if _, err := fs.Rename(ctx, media.videos[index].ID, target); err != nil {
					return nil, fmt.Errorf("重命名视频 %q -> %q 失败: %v", media.videos[index].Name, target, err)
				}
				log.Printf("[postprocess] 重命名视频: %s -> %s", media.videos[index].Name, target)
			}
			languageCounts := map[string]int{}
			for _, subtitle := range media.subtitles {
				language := detectSubtitleLanguage(subtitle.Name, pending.Title)
				languageCounts[language]++
				count := languageCounts[language]
				base := newName
				if count > 1 {
					base = fmt.Sprintf("%s - %d", newName, count)
				}
				targetBase := fmt.Sprintf("%s.%s", base, language)
				target := p115.PreserveFileExtension(targetBase, subtitle.Name, false)
				if target == subtitle.Name {
					continue
				}
				if _, err := fs.Rename(ctx, subtitle.ID, target); err != nil {
					return nil, fmt.Errorf("重命名字幕 %q -> %q 失败: %v", subtitle.Name, target, err)
				}
				log.Printf("[postprocess] 重命名字幕: %s -> %s", subtitle.Name, target)
			}
		} else {
			newName = p115.PreserveFileExtension(newName, entry.Name, false)
			if newName == "" {
				return nil, fmt.Errorf("重命名结果为空")
			}
			if newName != entry.Name {
				entry, err = fs.Rename(ctx, entry.ID, newName)
				if err != nil {
					return nil, fmt.Errorf("重命名 %q -> %q 失败: %v", entry.Name, newName, err)
				}
				log.Printf("[postprocess] 重命名: %s -> %s", pending.Title, newName)
			}
		}
	}

	if cfg.Offline.DualFolder || feedFinalCid != "" || feedFinalMode > 0 {
		finalCid := strings.TrimSpace(cfg.Offline.FinalCid)
		useSubfolder := false
		if feedFinalCid != "" {
			finalCid = feedFinalCid
			useSubfolder = feedCustomMode == 2
		} else if feedFinalMode == 2 {
			finalCid = strings.TrimSpace(cfg.Offline.FinalCid2)
			useSubfolder = true
		}
		if finalCid == "" {
			return nil, fmt.Errorf("保存文件夹为空")
		}
		if useSubfolder {
			target, targetErr := ensureFinalPath(ctx, fs, finalCid, cfg, pending, firstAir)
			if targetErr != nil {
				return nil, targetErr
			}
			finalCid = target
		}
		if cfg.Offline.KeepOnlyVideo && entry.IsDir() {
			if len(media.videos) == 0 {
				return nil, fmt.Errorf("离线目录中没有找到视频文件")
			}
			for _, item := range append(append([]cloudfs.Entry{}, media.videos...), media.subtitles...) {
				if item.ParentID == finalCid {
					finalEntries = append(finalEntries, store.PostEntry{ID: item.ID, Name: item.Name})
					continue
				}
				moved, moveErr := fs.Move(ctx, finalCid, item.ID)
				if moveErr != nil {
					return nil, fmt.Errorf("移动媒体文件 %q 到保存文件夹失败: %v", item.Name, moveErr)
				}
				finalEntries = append(finalEntries, store.PostEntry{ID: moved.ID, Name: moved.Name})
				log.Printf("[postprocess] 移动媒体文件: %s -> cid=%s", item.Name, finalCid)
			}
			if empty, cleanupErr := cleanupEmptyDir(ctx, fs, entry); cleanupErr != nil {
				log.Printf("[postprocess] 清理空目录失败: %v", cleanupErr)
			} else if empty {
				log.Printf("[postprocess] 已删除空目录: %s", entry.Name)
			}
		} else if entry.ParentID != finalCid {
			moved, moveErr := fs.Move(ctx, finalCid, entry.ID)
			if moveErr != nil {
				return nil, fmt.Errorf("移动 %q 到保存文件夹失败: %v", entry.Name, moveErr)
			}
			finalEntries = append(finalEntries, store.PostEntry{ID: moved.ID, Name: moved.Name})
			log.Printf("[postprocess] 移动: %s -> cid=%s", entry.Name, finalCid)
		}
	}
	return finalEntries, nil
}

func ensureFinalPath(ctx context.Context, fs postProcessFS, baseCid string, cfg *config.Config, pending store.PostTask, firstAir string) (string, error) {
	template := strings.TrimSpace(cfg.Offline.FinalPathTemplate)
	if template == "" {
		return baseCid, nil
	}
	feedName := pending.ConfigName
	site := ""
	if feed, ok := findFeedConfig(cfg, pending.ConfigName, ""); ok {
		if strings.TrimSpace(feed.Name) != "" {
			feedName = feed.Name
		}
		if u, err := url.Parse(feed.Url); err == nil {
			site = u.Hostname()
		}
	}
	seasonNumber := renameSeasonForFeedTitle(cfg, pending.ConfigName, "", pending.Title)
	if pending.RenameSeason > 0 {
		seasonNumber = pending.RenameSeason
	}
	tmdbID := ""
	if feed, ok := findFeedConfig(cfg, pending.ConfigName, ""); ok {
		tmdbID = strings.TrimSpace(feed.TmdbID)
	}
	airTime := time.Now()
	if parsed, err := time.Parse("2006-01-02", strings.TrimSpace(firstAir)); err == nil {
		airTime = parsed
	}
	expanded := config.ExpandFinalPathWithMeta(template, feedName, site, seasonNumber, airTime.Year(), tmdbID, airTime)
	if strings.TrimSpace(expanded) == "" {
		return baseCid, nil
	}
	mkdirFS, ok := fs.(interface {
		Mkdir(context.Context, string, string) (cloudfs.Entry, error)
	})
	if !ok {
		return "", fmt.Errorf("当前文件系统不支持创建最终子文件夹")
	}
	current := baseCid
	for _, segment := range strings.Split(expanded, "/") {
		segment = strings.TrimSpace(segment)
		if segment == "" || segment == "." || segment == ".." {
			continue
		}
		entry, err := fs.Lookup(ctx, current, segment)
		if err == nil && entry.ID != "" {
			current = entry.ID
			continue
		}
		created, createErr := mkdirFS.Mkdir(ctx, current, segment)
		if createErr != nil {
			return "", fmt.Errorf("创建最终子文件夹 %q 失败: %v", segment, createErr)
		}
		current = created.ID
	}
	return current, nil
}

func locateTaskEntry(ctx context.Context, fs postProcessFS, task *elevengo.OfflineTask, pending store.PostTask, cfg *config.Config) (cloudfs.Entry, error) {
	if task.FileId != "" {
		if entry, err := fs.Stat(ctx, task.FileId); err == nil && entry.ID != "" && entry.ID != "0" {
			return entry, nil
		}
	}
	cidSet := make([]string, 0, 3)
	seen := map[string]bool{}
	for _, cid := range []string{strings.TrimSpace(pending.OfflineCid), strings.TrimSpace(cfg.Offline.DefaultCid), "0"} {
		if cid == "" || seen[cid] {
			continue
		}
		seen[cid] = true
		cidSet = append(cidSet, cid)
	}
	for _, cid := range cidSet {
		if entry, err := fs.Lookup(ctx, cid, task.Name); err == nil && entry.ID != "" {
			return entry, nil
		}
	}
	searcher, canSearch := fs.(interface {
		Search(context.Context, string, string, cloudfs.SearchOptions) ([]cloudfs.Entry, error)
	})
	if canSearch {
		for _, cid := range cidSet {
			matches, err := searcher.Search(ctx, cid, task.Name, cloudfs.SearchOptions{})
			if err != nil {
				continue
			}
			for _, entry := range matches {
				if entry.Name == task.Name || strings.EqualFold(entry.Name, task.Name) {
					return entry, nil
				}
			}
		}
	}
	return cloudfs.Entry{}, fmt.Errorf("无法定位已完成文件 %q: 115 任务文件 ID 无效，且目录中未找到", task.Name)
}

func cleanupEmptyDir(ctx context.Context, fs postProcessFS, entry cloudfs.Entry) (bool, error) {
	if !entry.IsDir() {
		return false, nil
	}
	children, err := fs.List(ctx, entry.ID)
	if err != nil {
		return false, err
	}
	remaining := 0
	for _, child := range children {
		if child.IsDir() {
			empty, childErr := cleanupEmptyDir(ctx, fs, child)
			if childErr != nil {
				return false, childErr
			}
			if !empty {
				remaining++
			}
			continue
		}
		remaining++
	}
	if remaining > 0 {
		return false, nil
	}
	if err := fs.Delete(ctx, entry.ID); err != nil {
		return false, err
	}
	return true, nil
}

var videoExtensions = map[string]bool{
	".mkv": true, ".mp4": true, ".avi": true, ".mov": true, ".wmv": true,
	".flv": true, ".webm": true, ".m4v": true, ".ts": true, ".m2ts": true,
	".rmvb": true, ".rm": true, ".mpg": true, ".mpeg": true, ".iso": true,
}

var subtitleExtensions = map[string]bool{
	".ass": true, ".ssa": true, ".srt": true, ".vtt": true,
	".sub": true, ".idx": true, ".sup": true, ".ttml": true, ".dfxp": true,
}

func isVideoFile(name string) bool {
	index := strings.LastIndex(name, ".")
	if index < 0 {
		return false
	}
	return videoExtensions[strings.ToLower(name[index:])]
}

func isSubtitleFile(name string) bool {
	index := strings.LastIndex(name, ".")
	if index < 0 {
		return false
	}
	return subtitleExtensions[strings.ToLower(name[index:])]
}

type mediaFiles struct {
	videos    []cloudfs.Entry
	subtitles []cloudfs.Entry
}

func findMediaFiles(ctx context.Context, fs postProcessFS, entry cloudfs.Entry) (mediaFiles, error) {
	if !entry.IsDir() {
		return mediaFiles{videos: []cloudfs.Entry{entry}}, nil
	}
	files, err := listMediaFiles(ctx, fs, entry.ID, 0)
	if err != nil {
		return mediaFiles{}, err
	}
	return files, nil
}

func listMediaFiles(ctx context.Context, fs postProcessFS, dirID string, depth int) (mediaFiles, error) {
	result := mediaFiles{}
	if depth > 2 {
		return result, nil
	}
	entries, err := fs.List(ctx, dirID)
	if err != nil {
		return result, fmt.Errorf("读取离线目录失败: %v", err)
	}
	for _, child := range entries {
		if child.IsDir() {
			nested, nestedErr := listMediaFiles(ctx, fs, child.ID, depth+1)
			if nestedErr != nil {
				return result, nestedErr
			}
			result.videos = append(result.videos, nested.videos...)
			result.subtitles = append(result.subtitles, nested.subtitles...)
			continue
		}
		if isVideoFile(child.Name) {
			result.videos = append(result.videos, child)
			continue
		}
		if isSubtitleFile(child.Name) {
			result.subtitles = append(result.subtitles, child)
		}
	}
	return result, nil
}

func detectSubtitleLanguage(name, title string) string {
	if language := subtitleLanguageFromText(name); language != "" {
		return language
	}
	if language := subtitleLanguageFromText(title); language != "" {
		return language
	}
	return "und"
}

func subtitleLanguageFromText(value string) string {
	value = strings.ToLower(value)
	switch {
	case strings.Contains(value, "简") || strings.Contains(value, "chs") || strings.Contains(value, "sc") || strings.Contains(value, "gb"):
		return "sc"
	case strings.Contains(value, "繁") || strings.Contains(value, "cht") || strings.Contains(value, "tc") || strings.Contains(value, "big5"):
		return "tc"
	case strings.Contains(value, "日") || strings.Contains(value, "jp"):
		return "jp"
	case strings.Contains(value, "英") || strings.Contains(value, "eng"):
		return "en"
	default:
		return ""
	}
}

func appendNameIndex(name string, index int) string {
	ext := ""
	if idx := strings.LastIndex(name, "."); idx > 0 {
		ext = name[idx:]
		name = name[:idx]
	}
	return fmt.Sprintf("%s - %d%s", name, index, ext)
}

func (s *Server) handleRunPostProcess(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if !postProcessEnabled(cfg) {
		writeAPIError(w, http.StatusBadRequest, "请先启用重命名或离线后转移到保存文件夹")
		return
	}
	go s.runPostProcess()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
