package webui

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/huohen92/anifeed/qbittorrent"
)

func (s *Server) handleQBTest(w http.ResponseWriter, r *http.Request) {
	cfg := s.configSnapshot()
	if cfg == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "配置未加载")
		return
	}
	q := cfg.Downloaders.QBittorrent
	client, err := qbittorrent.New(q.URL, q.Username, q.Password)
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
				client, err := qbittorrent.New(cfg.Downloaders.QBittorrent.URL, cfg.Downloaders.QBittorrent.Username, cfg.Downloaders.QBittorrent.Password)
				if err != nil {
					continue
				}
				tasks, err := ag.ListOfflineTasks(0, nil)
				if err != nil {
					continue
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
					if err := client.Add(ctx, task.Url, cfg.Downloaders.QBittorrent.SavePath, cfg.Downloaders.QBittorrent.Category); err != nil {
						log.Printf("[qbittorrent] 提交失败 %s: %v", task.Name, err)
						continue
					}
					started[task.InfoHash] = true
					log.Printf("[qbittorrent] 115 失败 %d 次，已转 qBittorrent: %s", failures[task.InfoHash], task.Name)
				}
			}
		}
	}()
}
