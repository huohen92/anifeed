package webui

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type tmdbCacheDisk struct {
	Episodes []tmdbEpisode `json:"episodes"`
	At       time.Time     `json:"at"`
}

func (s *Server) tmdbCachePath() string {
	if s.configPath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(s.configPath), ".tmdb-cache.json")
}

func tmdbCacheKey(apiKey, id string) string {
	sum := sha256.Sum256([]byte(apiKey + "|" + id))
	return hex.EncodeToString(sum[:])
}

func (s *Server) loadTMDBDiskCache() {
	path := s.tmdbCachePath()
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var saved map[string]tmdbCacheDisk
	if json.Unmarshal(data, &saved) != nil {
		return
	}
	now := time.Now()
	s.tmdbCacheMu.Lock()
	defer s.tmdbCacheMu.Unlock()
	if s.tmdbCache == nil {
		s.tmdbCache = make(map[string]tmdbEpisodesCacheEntry)
	}
	for key, entry := range saved {
		if now.Sub(entry.At) < 6*time.Hour {
			s.tmdbCache[key] = tmdbEpisodesCacheEntry{episodes: entry.Episodes, at: entry.At}
		}
	}
}

func (s *Server) saveTMDBDiskCacheLocked() {
	path := s.tmdbCachePath()
	if path == "" {
		return
	}
	now := time.Now()
	saved := make(map[string]tmdbCacheDisk, len(s.tmdbCache))
	for key, entry := range s.tmdbCache {
		if now.Sub(entry.at) < 6*time.Hour {
			saved[key] = tmdbCacheDisk{Episodes: entry.episodes, At: entry.at}
		}
	}
	data, err := json.Marshal(saved)
	if err != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmdb-cache-*.tmp")
	if err != nil {
		return
	}
	name := tmp.Name()
	defer os.Remove(name)
	_ = tmp.Chmod(0o600)
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	_ = os.Rename(name, path)
}
