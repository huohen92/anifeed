package webui

import (
	"context"
	"log"
	"time"

	"github.com/huohen92/anifeed/config"
)

// autoLoop checks the [auto] configuration periodically and triggers scheduled
// RSS runs. A short ticker lets config changes take effect without a restart.
func (s *Server) autoLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var next time.Time
	started := false
	for {
		cfg := s.configSnapshot()
		if cfg != nil && cfg.Auto.Enabled {
			now := time.Now()
			if !started && cfg.Auto.RunOnStart {
				next = now
				started = true
			}
			if next.IsZero() {
				next = now.Add(autoInterval(cfg))
			}
			if !next.After(now) {
				s.runAutoJob()
				next = now.Add(autoInterval(cfg))
			}
		} else {
			next = time.Time{}
			started = false
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func autoInterval(cfg *config.Config) time.Duration {
	minutes := cfg.Auto.IntervalMinutes
	if minutes <= 0 {
		minutes = 30
	}
	return time.Duration(minutes) * time.Minute
}

func (s *Server) runAutoJob() {
	if s.jobs.isRunning() {
		log.Printf("[auto] 跳过定时任务：已有运行中的作业")
		return
	}
	ag, err := s.ensureAgent()
	if err != nil {
		log.Printf("[auto] 跳过定时任务: %v", err)
		return
	}
	if _, err := s.startRunAllJob(ag, "auto", "自动运行全部订阅"); err != nil {
		log.Printf("[auto] 启动定时任务失败: %v", err)
	}
}
