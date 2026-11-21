package webui

import (
	"context"
	"log"
	"strings"
	"time"
)

// startAutoBackupMonitor watches failed 115 offline tasks. After the configured
// number of failed attempts, it runs the backup-subgroup discovery.
func (s *Server) startAutoBackupMonitor(ctx context.Context) {
	go func() {
		time.Sleep(25 * time.Second)
		ticker := time.NewTicker(45 * time.Second)
		defer ticker.Stop()
		failures := map[string]int{}
		handled := map[string]bool{}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cfg := s.configSnapshot()
				if cfg == nil || !cfg.Lab.AutoBackupSubgroup {
					continue
				}
				ag, err := s.ensureAgent()
				if err != nil {
					continue
				}
				tasks, err := ag.ListOfflineTasks(0, nil)
				if err != nil {
					continue
				}
				triggered := false
				for i := range tasks {
					task := &tasks[i]
					if !task.IsFailed() || strings.TrimSpace(task.Url) == "" || handled[task.InfoHash] {
						continue
					}
					failures[task.InfoHash]++
					threshold := cfg.Lab.AutoBackupFailures
					if threshold <= 0 {
						threshold = 2
					}
					if failures[task.InfoHash] < threshold {
						continue
					}
					handled[task.InfoHash] = true
					triggered = true
					log.Printf("[auto-backup] 115 失败 %d 次，开始检查备用字幕组: %s", failures[task.InfoHash], task.Name)
				}
				if triggered {
					go s.autoBackupSubgroups(cfg)
				}
			}
		}
	}()
}
