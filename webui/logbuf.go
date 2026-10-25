package webui

import (
	"strings"
	"sync"
)

// LogBuffer keeps the most recent log lines for the Web UI log viewer.
type LogBuffer struct {
	mu    sync.Mutex
	lines []string
	max   int
}

// NewLogBuffer creates a ring buffer holding at most max lines.
func NewLogBuffer(max int) *LogBuffer {
	if max <= 0 {
		max = 500
	}
	return &LogBuffer{max: max}
}

// Write implements io.Writer. It splits input on newlines and timestamps each
// line's arrival by relying on log's own prefix.
func (b *LogBuffer) Write(p []byte) (int, error) {
	text := strings.TrimRight(string(p), "\n")
	if text == "" {
		return len(p), nil
	}
	lines := strings.Split(text, "\n")

	b.mu.Lock()
	b.lines = append(b.lines, lines...)
	if len(b.lines) > b.max {
		b.lines = append([]string(nil), b.lines[len(b.lines)-b.max:]...)
	}
	b.mu.Unlock()
	return len(p), nil
}

// Tail returns up to n most recent lines.
func (b *LogBuffer) Tail(n int) []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n <= 0 || n > len(b.lines) {
		n = len(b.lines)
	}
	out := make([]string, n)
	copy(out, b.lines[len(b.lines)-n:])
	return out
}
