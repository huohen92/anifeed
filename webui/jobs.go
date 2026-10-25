package webui

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/huohen92/anifeed/p115"
)

// Job is an asynchronous feed run shown in the UI.
type Job struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Title      string               `json:"title"`
	Status     string               `json:"status"` // running | done | error
	StartedAt  time.Time            `json:"startedAt"`
	FinishedAt *time.Time           `json:"finishedAt,omitempty"`
	Results    []p115.FeedRunResult `json:"results,omitempty"`
	Error      string               `json:"error,omitempty"`
}

type jobManager struct {
	mu      sync.Mutex
	jobs    map[string]*Job
	order   []string
	running bool
}

func newJobManager() *jobManager {
	return &jobManager{jobs: make(map[string]*Job)}
}

func (m *jobManager) start(kind, title string, fn func() ([]p115.FeedRunResult, error)) (Job, error) {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return Job{}, ErrJobRunning
	}
	id := strconv.FormatInt(time.Now().UnixNano(), 36)
	job := &Job{
		ID:        id,
		Type:      kind,
		Title:     title,
		Status:    "running",
		StartedAt: time.Now(),
	}
	m.jobs[id] = job
	m.order = append(m.order, id)
	if len(m.order) > 30 {
		delete(m.jobs, m.order[0])
		m.order = m.order[1:]
	}
	m.running = true
	m.mu.Unlock()

	snapshot := *job
	go func() {
		var (
			results []p115.FeedRunResult
			err     error
		)
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					err = fmt.Errorf("任务异常退出: %v", recovered)
				}
			}()
			results, err = fn()
		}()

		now := time.Now()
		m.mu.Lock()
		defer m.mu.Unlock()
		job.FinishedAt = &now
		job.Results = results
		if err != nil {
			job.Status = "error"
			job.Error = err.Error()
		} else {
			job.Status = "done"
		}
		m.running = false
	}()
	return snapshot, nil
}

func (m *jobManager) get(id string) (Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	job, ok := m.jobs[id]
	if !ok {
		return Job{}, false
	}
	return copyJob(job), true
}

func (m *jobManager) list() []Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Job, 0, len(m.order))
	for i := len(m.order) - 1; i >= 0; i-- {
		if job, ok := m.jobs[m.order[i]]; ok {
			out = append(out, copyJob(job))
		}
	}
	return out
}

func (m *jobManager) isRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func copyJob(job *Job) Job {
	out := *job
	out.Results = append([]p115.FeedRunResult(nil), job.Results...)
	if job.FinishedAt != nil {
		t := *job.FinishedAt
		out.FinishedAt = &t
	}
	return out
}
