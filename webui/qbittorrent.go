package webui

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

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

func (s *Server) startQBFallbackLoop(ctx context.Context) {
	go func() {
		time.Sleep(18 * time.Second)
		ticker := time.NewTicker(45 * time.Second)
		defer ticker.Stop()
		failures := map[string]int{}
		started := map[string]bool{}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cfg := s.configSnapshot()
				if cfg == nil || !cfg.Downloaders.QBittorrent.Enabled {
					continue
				}
				ag, err := s.ensureAgent()
				if err != nil {
					continue
				}
				if ag.StoreInstance != nil {
					if saved, err := ag.StoreInstance.QBFallbackHashes(); err == nil {
						for hash := range saved {
							started[hash] = true
						}
					}
				}
				client, err := newQBClient(cfg.Downloaders.QBittorrent)
				if err != nil {
					continue
				}
				tasks, err := ag.ListOfflineTasks(0, nil)
				if err != nil {
					continue
				}
				if torrents, err := client.Torrents(ctx, cfg.Downloaders.QBittorrent.Category); err == nil {
					for _, torrent := range torrents {
						if torrent.Progress >= 1 && ag.StoreInstance != nil {
							_ = ag.StoreInstance.MarkQBFallbackCompleted(torrent.Hash)
						}
					}
				}
				stillFailed := make(map[string]bool, len(tasks))
				for i := range tasks {
					task := &tasks[i]
					if !task.IsFailed() || strings.TrimSpace(task.Url) == "" {
						continue
					}
					stillFailed[task.InfoHash] = true
				}
				for hash := range failures {
					if !stillFailed[hash] {
						delete(failures, hash)
					}
				}
				for i := range tasks {
					task := &tasks[i]
					if !task.IsFailed() || strings.TrimSpace(task.Url) == "" || started[task.InfoHash] {
						continue
					}
					failures[task.InfoHash]++
					threshold := cfg.Downloaders.QBittorrent.FallbackThreshold
					if threshold <= 0 {
						threshold = 3
					}
					if failures[task.InfoHash] < threshold {
						continue
					}
					q := cfg.Downloaders.QBittorrent
					if err := client.AddWithOptions(ctx, task.Url, qbittorrent.AddOptions{
						SavePath: q.SavePath, Category: q.Category, Tags: q.Tags, ContentLayout: q.ContentLayout,
						DownloadLimitKB: q.DownloadLimitKB, UploadLimitKB: q.UploadLimitKB,
						RatioLimit: q.RatioLimit, SeedingTimeLimit: q.SeedingTimeLimit,
						InactiveSeedingTimeLimit: q.InactiveSeedingTimeLimit, UseDownloadPath: q.UseDownloadPath,
					}); err != nil {
						log.Printf("[qbittorrent] 提交失败 %s: %v", task.Name, err)
						continue
					}
					started[task.InfoHash] = true
					if ag.StoreInstance != nil {
						_ = ag.StoreInstance.MarkQBFallbackSubmitted(task.InfoHash, task.Name)
					}
					log.Printf("[qbittorrent] 115 失败 %d 次，已转 qBittorrent: %s", failures[task.InfoHash], task.Name)
				}
			}
		}
	}()
}
