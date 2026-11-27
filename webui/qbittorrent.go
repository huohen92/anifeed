package webui

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/deadblue/elevengo"
	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/qbittorrent"
)

func (s *Server) handleQBTest(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "配置未加载")
		return
	}
	q := cfg.Downloaders.QBittorrent
	client, err := newQBClient(q)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	version, err := client.Version(r.Context())
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "version": version})
}

func newQBClient(q config.DownloaderConfig) (*qbittorrent.Client, error) {
	if strings.TrimSpace(q.APIKey) != "" {
		return qbittorrent.NewWithAPIKey(q.URL, q.APIKey)
	}
	return qbittorrent.New(q.URL, q.Username, q.Password)
}

// observeQBFallback inspects an already-fetched 115 task list. It is called
// only from normal user/API flows, so AniFeed never polls 115 just to look for
// failed tasks.
func (s *Server) observeQBFallback(ctx context.Context, cfg *config.Config, tasks []elevengo.OfflineTask) {
	if cfg == nil || !cfg.Downloaders.QBittorrent.Enabled {
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		return
	}
	q := cfg.Downloaders.QBittorrent
	client, err := newQBClient(q)
	if err != nil {
		return
	}
	if torrents, terr := client.Torrents(ctx, q.Category); terr == nil {
		for _, torrent := range torrents {
			if torrent.Progress >= 1 && ag.StoreInstance != nil {
				_ = ag.StoreInstance.MarkQBFallbackCompleted(torrent.Hash)
			}
		}
	}

	s.qbFallbackMu.Lock()
	if s.qbFallbackFailures == nil {
		s.qbFallbackFailures = map[string]int{}
	}
	if s.qbFallbackStarted == nil {
		s.qbFallbackStarted = map[string]bool{}
		if ag.StoreInstance != nil {
			if saved, err := ag.StoreInstance.QBFallbackHashes(); err == nil {
				for hash := range saved {
					s.qbFallbackStarted[hash] = true
				}
			}
		}
	}
	threshold := q.FallbackThreshold
	if threshold <= 0 {
		threshold = 3
	}
	observed := map[string]bool{}
	var candidates []elevengo.OfflineTask
	for i := range tasks {
		task := tasks[i]
		if !task.IsFailed() || strings.TrimSpace(task.Url) == "" || s.qbFallbackStarted[task.InfoHash] {
			continue
		}
		observed[task.InfoHash] = true
		s.qbFallbackFailures[task.InfoHash]++
		if s.qbFallbackFailures[task.InfoHash] < threshold {
			continue
		}
		s.qbFallbackStarted[task.InfoHash] = true
		candidates = append(candidates, task)
	}
	for hash := range s.qbFallbackFailures {
		if !observed[hash] {
			delete(s.qbFallbackFailures, hash)
		}
	}
	s.qbFallbackMu.Unlock()

	for _, task := range candidates {
		if err := client.AddWithOptions(ctx, task.Url, qbittorrent.AddOptions{
			SavePath: q.SavePath, Category: q.Category, Tags: q.Tags, ContentLayout: q.ContentLayout,
			DownloadLimitKB: q.DownloadLimitKB, UploadLimitKB: q.UploadLimitKB,
			RatioLimit: q.RatioLimit, SeedingTimeLimit: q.SeedingTimeLimit,
			InactiveSeedingTimeLimit: q.InactiveSeedingTimeLimit, UseDownloadPath: q.UseDownloadPath,
		}); err != nil {
			s.qbFallbackMu.Lock()
			delete(s.qbFallbackStarted, task.InfoHash)
			s.qbFallbackMu.Unlock()
			log.Printf("[qbittorrent] 提交失败 %s: %v", task.Name, err)
			continue
		}
		if ag.StoreInstance != nil {
			_ = ag.StoreInstance.MarkQBFallbackSubmitted(task.InfoHash, task.Name)
		}
		log.Printf("[qbittorrent] 115 失败任务已转为 qBittorrent: %s", task.Name)
	}
}
