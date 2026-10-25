package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zhifengle/anifeed/rsssite"
)

type Store struct {
	DBInstance *sql.DB
}

// New creates a new Store instance with the given database connection.
// If db is nil, it opens a SQLite database at the default path "db.sqlite".
// Deprecated: Use NewWithPath to specify a custom database path.
func New(db *sql.DB) *Store {
	if db == nil {
		store, err := NewWithPath("db.sqlite")
		if err != nil {
			panic(err)
		}
		return store
	}
	if err := initSchema(db); err != nil {
		panic(err)
	}
	if err := migrateSchema(db); err != nil {
		panic(err)
	}
	return &Store{
		DBInstance: db,
	}
}

// NewWithPath creates a new Store instance with a database at the specified path.
// If path is empty, it defaults to "db.sqlite".
func NewWithPath(path string) (*Store, error) {
	if path == "" {
		path = "db.sqlite"
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrateSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{
		DBInstance: db,
	}, nil
}

func initSchema(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE if not exists `rss_items` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `link` VARCHAR(255), `title` VARCHAR(255), `guid` VARCHAR(255), `pubDate` DATETIME, `creator` VARCHAR(255), `summary` TEXT, `content` VARCHAR(255), `isoDate` DATETIME, `categories` VARCHAR(255), `contentSnippet` VARCHAR(255), `done` TINYINT(1) DEFAULT 0, `magnet` VARCHAR(255) NOT NULL, `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL)"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE if not exists `sites_status` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `name` VARCHAR(255), `needLogin` TINYINT(1), `abnormalOp` TINYINT(1), `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL)"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE if not exists `offline_postprocess` (`task_hash` VARCHAR(255) PRIMARY KEY, `title` TEXT, `config_name` VARCHAR(255), `offline_cid` VARCHAR(255), `state` VARCHAR(32) NOT NULL DEFAULT 'pending', `attempts` INTEGER NOT NULL DEFAULT 0, `error` TEXT, `final_entries` TEXT, `created_at` DATETIME NOT NULL, `updated_at` DATETIME NOT NULL)"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE if not exists `preview_decisions` (`feed_url` TEXT NOT NULL, `item_key` TEXT NOT NULL, `allowed` TINYINT(1) NOT NULL, `updated_at` DATETIME NOT NULL, PRIMARY KEY (`feed_url`, `item_key`))"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE if not exists `idle_reminders` (`feed_url` TEXT NOT NULL, `day` TEXT NOT NULL, `sent_at` DATETIME NOT NULL, PRIMARY KEY (`feed_url`, `day`))"); err != nil {
		return err
	}
	// Older databases do not have final_entries yet. ALTER TABLE has no
	// IF NOT EXISTS in SQLite, so ignore "duplicate column" errors.
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `final_entries` TEXT")
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `source_url` TEXT")
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `rename_season` INTEGER NOT NULL DEFAULT 0")
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `rename_offset` INTEGER NOT NULL DEFAULT 0")
	return nil
}

// migrateSchema adds columns introduced after the initial schema to existing
// rss_items tables. Idempotent: skips columns that already exist.
func migrateSchema(db *sql.DB) error {
	cols, err := tableColumns(db, "rss_items")
	if err != nil {
		return err
	}
	additions := []struct {
		col string
		sql string
	}{
		{"configName", "ALTER TABLE `rss_items` ADD COLUMN `configName` VARCHAR(255)"},
		{"configUrl", "ALTER TABLE `rss_items` ADD COLUMN `configUrl` VARCHAR(255)"},
		{"taskHash", "ALTER TABLE `rss_items` ADD COLUMN `taskHash` VARCHAR(255)"},
	}
	for _, a := range additions {
		if cols[a.col] {
			continue
		}
		if _, err := db.Exec(a.sql); err != nil {
			return err
		}
	}
	if _, err := db.Exec("CREATE TABLE if not exists `offline_postprocess` (`task_hash` VARCHAR(255) PRIMARY KEY, `title` TEXT, `config_name` VARCHAR(255), `offline_cid` VARCHAR(255), `state` VARCHAR(32) NOT NULL DEFAULT 'pending', `attempts` INTEGER NOT NULL DEFAULT 0, `error` TEXT, `final_entries` TEXT, `created_at` DATETIME NOT NULL, `updated_at` DATETIME NOT NULL)"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE if not exists `preview_decisions` (`feed_url` TEXT NOT NULL, `item_key` TEXT NOT NULL, `allowed` TINYINT(1) NOT NULL, `updated_at` DATETIME NOT NULL, PRIMARY KEY (`feed_url`, `item_key`))"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE if not exists `idle_reminders` (`feed_url` TEXT NOT NULL, `day` TEXT NOT NULL, `sent_at` DATETIME NOT NULL, PRIMARY KEY (`feed_url`, `day`))"); err != nil {
		return err
	}
	// Older databases do not have final_entries yet. ALTER TABLE has no
	// IF NOT EXISTS in SQLite, so ignore "duplicate column" errors.
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `final_entries` TEXT")
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `source_url` TEXT")
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `rename_season` INTEGER NOT NULL DEFAULT 0")
	_, _ = db.Exec("ALTER TABLE `offline_postprocess` ADD COLUMN `rename_offset` INTEGER NOT NULL DEFAULT 0")
	return nil
}

// tableColumns returns the set of column names present in a table.
func tableColumns(db *sql.DB, table string) (map[string]bool, error) {
	rows, err := db.Query("PRAGMA table_info(`" + table + "`)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols := make(map[string]bool)
	for rows.Next() {
		var cid, notNull, pk int
		var name, typ string
		var dflt any
		if err := rows.Scan(&cid, &name, &typ, &notNull, &dflt, &pk); err != nil {
			return nil, err
		}
		cols[name] = true
	}
	return cols, rows.Err()
}

// SaveMagnetItems records accepted magnet items with their RSS feed context
// and the non-empty 115 offline task hashes returned by OfflineAddUrl.
// Each item must have a corresponding hash.
func (s *Store) SaveMagnetItems(items []rsssite.MagnetItem, configName, configUrl string, hashes []string) error {
	return s.SaveMagnetItemsWithPostprocess(items, configName, configUrl, hashes, "", 0, 0)
}

// SaveMagnetItemsWithPostprocess records accepted items and creates a pending
// row for the post-download worker (rename / final folder move).
func (s *Store) SaveMagnetItemsWithPostprocess(items []rsssite.MagnetItem, configName, configUrl string, hashes []string, offlineCid string, renameSeason, renameOffset int) error {
	if len(items) == 0 {
		return nil
	}
	if len(hashes) != len(items) {
		return fmt.Errorf("magnet item/hash count mismatch: %d items, %d hashes", len(items), len(hashes))
	}
	for i, hash := range hashes {
		if hash == "" {
			return fmt.Errorf("magnet item %d has an empty task hash", i)
		}
	}
	tx, err := s.DBInstance.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	for i, item := range items {
		stmt := "INSERT INTO rss_items (`link`,`title`,`content`,`magnet`,`configName`,`configUrl`,`taskHash`,`done`,`createdAt`,`updatedAt`) VALUES (?,?,?,?,?,?,?,?,?,?)"
		_, err := tx.Exec(stmt, item.Link, item.Title, item.Content, item.Magnet, configName, configUrl, hashes[i], 0, now, now)
		if err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT OR IGNORE INTO `offline_postprocess` (`task_hash`,`title`,`config_name`,`offline_cid`,`source_url`,`rename_season`,`rename_offset`,`state`,`attempts`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?, 'pending', 0, ?, ?)", hashes[i], item.Title, configName, offlineCid, configUrl, renameSeason, renameOffset, now, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// PostTask is one completed-download entry waiting for rename/move.
type PostTask struct {
	TaskHash     string
	Title        string
	ConfigName   string
	OfflineCid   string
	SourceURL    string
	RenameSeason int
	RenameOffset int
	State        string
	Attempts     int
	Error        string
	CreatedAt    time.Time
}

// PostEntry identifies one final file/folder in 115 after post-processing.
type PostEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CompletePostTaskWithEntries marks a task done and records the final entries
// that were renamed/moved, so history deletion can remove the actual files.
func (s *Store) CompletePostTaskWithEntries(taskHash string, entries []PostEntry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	_, err = s.DBInstance.Exec("UPDATE offline_postprocess SET state = 'done', error = '', final_entries = ?, updated_at = ? WHERE task_hash = ?", string(data), time.Now(), taskHash)
	return err
}

// FinalEntriesForTaskHashes returns recorded final entries grouped by task hash.
func (s *Store) FinalEntriesForTaskHashes(hashes []string) (map[string][]PostEntry, error) {
	result := make(map[string][]PostEntry)
	if len(hashes) == 0 {
		return result, nil
	}
	for start := 0; start < len(hashes); start += 500 {
		end := start + 500
		if end > len(hashes) {
			end = len(hashes)
		}
		chunk := hashes[start:end]
		placeholders := strings.TrimRight(strings.Repeat("?,", len(chunk)), ",")
		args := make([]any, 0, len(chunk))
		for _, hash := range chunk {
			args = append(args, hash)
		}
		rows, err := s.DBInstance.Query("SELECT task_hash, COALESCE(final_entries,'') FROM offline_postprocess WHERE task_hash IN ("+placeholders+")", args...)
		if err != nil {
			return result, err
		}
		for rows.Next() {
			var hash, raw string
			if err := rows.Scan(&hash, &raw); err != nil {
				rows.Close()
				return result, err
			}
			if raw == "" {
				continue
			}
			var entries []PostEntry
			if err := json.Unmarshal([]byte(raw), &entries); err != nil {
				continue
			}
			result[hash] = entries
		}
		rows.Close()
	}
	return result, nil
}

// ListPendingPostTasks returns pending post-processing rows, oldest first.
func (s *Store) ListPendingPostTasks(limit int) ([]PostTask, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.DBInstance.Query("SELECT task_hash, COALESCE(title,''), COALESCE(config_name,''), COALESCE(offline_cid,''), COALESCE(source_url,''), rename_season, rename_offset, state, attempts, COALESCE(error,''), created_at FROM offline_postprocess WHERE state = 'pending' ORDER BY created_at ASC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []PostTask
	for rows.Next() {
		var t PostTask
		if err := rows.Scan(&t.TaskHash, &t.Title, &t.ConfigName, &t.OfflineCid, &t.SourceURL, &t.RenameSeason, &t.RenameOffset, &t.State, &t.Attempts, &t.Error, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// CountPendingPostTasks counts pending rows.
func (s *Store) CountPendingPostTasks() (int, error) {
	var count int
	err := s.DBInstance.QueryRow("SELECT count(*) FROM offline_postprocess WHERE state = 'pending'").Scan(&count)
	return count, err
}

// CompletePostTask marks a post-processing row as done.
func (s *Store) CompletePostTask(taskHash string) error {
	_, err := s.DBInstance.Exec("UPDATE offline_postprocess SET state = 'done', error = '', updated_at = ? WHERE task_hash = ?", time.Now(), taskHash)
	return err
}

// FailPostTask increments attempts and records the last error. It keeps the
// row pending until maxAttempts is reached, then marks it failed.
func (s *Store) FailPostTask(taskHash, message string, maxAttempts int) error {
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	var attempts int
	if err := s.DBInstance.QueryRow("SELECT attempts FROM offline_postprocess WHERE task_hash = ?", taskHash).Scan(&attempts); err != nil {
		return err
	}
	state := "pending"
	if attempts+1 >= maxAttempts {
		state = "failed"
	}
	_, err := s.DBInstance.Exec("UPDATE offline_postprocess SET state = ?, attempts = attempts + 1, error = ?, updated_at = ? WHERE task_hash = ?", state, message, time.Now(), taskHash)
	return err
}

func (s *Store) HasItem(magnet string) bool {
	var count int
	s.DBInstance.QueryRow("SELECT count(*) AS num FROM rss_items WHERE magnet = ?", magnet).Scan(&count)
	return count > 0
}

// @TODO 替换 HasItem. 注意目前 magnet 存的长度是 VARCHAR(255)。有tracker的长URI会存不了.
func (s *Store) HasMagnetByXt(magnet string) bool {
	var count int
	u, err := url.Parse(magnet)
	if err != nil {
		return false
	}
	params := u.Query()
	xt := params.Get("xt")
	s.DBInstance.QueryRow("SELECT count(*) AS num FROM rss_items WHERE magnet LIKE ?", "%"+xt+"%").Scan(&count)
	return count > 0
}

// PreviewDecisions returns per-item allow/deny decisions for a feed.
func (s *Store) PreviewDecisions(feedURL string) (map[string]bool, error) {
	decisions := make(map[string]bool)
	rows, err := s.DBInstance.Query("SELECT item_key, allowed FROM preview_decisions WHERE feed_url = ?", feedURL)
	if err != nil {
		return decisions, err
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var allowed bool
		if err := rows.Scan(&key, &allowed); err != nil {
			return decisions, err
		}
		decisions[key] = allowed
	}
	return decisions, rows.Err()
}

// SavePreviewDecisions replaces the decisions for one feed.
func (s *Store) SavePreviewDecisions(feedURL string, decisions map[string]bool) error {
	tx, err := s.DBInstance.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("DELETE FROM preview_decisions WHERE feed_url = ?", feedURL); err != nil {
		return err
	}
	now := time.Now()
	for key, allowed := range decisions {
		if _, err := tx.Exec("INSERT INTO preview_decisions (feed_url, item_key, allowed, updated_at) VALUES (?, ?, ?, ?)", feedURL, key, allowed, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// HistoryItem is a single submission record from the rss_items table.
type HistoryItem struct {
	ID         int64
	Title      string
	ConfigName string
	ConfigUrl  string
	CreatedAt  time.Time
	TaskHash   string
}

// ListHistory returns the most recent submission records, newest first.
// A non-empty name filters by RSS feed config name. limit <= 0 defaults to 20.
func (s *Store) ListHistory(limit int, name string) ([]HistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	query := "SELECT id, COALESCE(title,''), COALESCE(configName,''), COALESCE(configUrl,''), createdAt, COALESCE(taskHash,'') FROM rss_items"
	var args []any
	if name != "" {
		query += " WHERE configName = ?"
		args = append(args, name)
	}
	query += " ORDER BY id DESC LIMIT ?"
	args = append(args, limit)
	rows, err := s.DBInstance.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []HistoryItem
	for rows.Next() {
		var h HistoryItem
		if err := rows.Scan(&h.ID, &h.Title, &h.ConfigName, &h.ConfigUrl, &h.CreatedAt, &h.TaskHash); err != nil {
			return nil, err
		}
		items = append(items, h)
	}
	return items, rows.Err()
}

// CountHistory returns the total number of submissions.
func (s *Store) CountHistory() (int, error) {
	var count int
	err := s.DBInstance.QueryRow("SELECT count(*) FROM rss_items").Scan(&count)
	return count, err
}

// CountHistorySince returns the number of submissions created at or after t.
func (s *Store) CountHistorySince(t time.Time) (int, error) {
	var count int
	err := s.DBInstance.QueryRow("SELECT count(*) FROM rss_items WHERE createdAt >= ?", t).Scan(&count)
	return count, err
}

func (s *Store) Close() error {
	return s.DBInstance.Close()
}

// SubmittedItem is a recorded RSS submission used by the rename preview.
type SubmittedItem struct {
	ID       int64
	Title    string
	TaskHash string
}

// ListSubmittedItems returns submitted items for one feed, oldest first.
func (s *Store) ListSubmittedItems(configName, configURL string, limit int) ([]SubmittedItem, error) {
	if limit <= 0 {
		limit = 200
	}
	query := "SELECT id, COALESCE(title,''), COALESCE(taskHash,'') FROM rss_items WHERE taskHash != ''"
	var args []any
	if strings.TrimSpace(configURL) != "" {
		query += " AND configUrl = ?"
		args = append(args, configURL)
	} else if strings.TrimSpace(configName) != "" {
		query += " AND configName = ?"
		args = append(args, configName)
	} else {
		return nil, nil
	}
	query += " ORDER BY id ASC LIMIT ?"
	args = append(args, limit)
	rows, err := s.DBInstance.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SubmittedItem
	for rows.Next() {
		var item SubmittedItem
		if err := rows.Scan(&item.ID, &item.Title, &item.TaskHash); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// DeleteDownloadedByMagnets removes submission history and pending
// postprocess records for the given magnets. It returns the number of
// history rows deleted. Files in 115 are not touched.
func (s *Store) DeleteDownloadedByMagnets(magnets []string) (int64, error) {
	unique := make([]string, 0, len(magnets))
	seen := map[string]bool{}
	for _, magnet := range magnets {
		magnet = strings.TrimSpace(magnet)
		if magnet == "" || seen[magnet] {
			continue
		}
		seen[magnet] = true
		unique = append(unique, magnet)
	}
	if len(unique) == 0 {
		return 0, nil
	}
	tx, err := s.DBInstance.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var hashes []string
	for start := 0; start < len(unique); start += 500 {
		end := start + 500
		if end > len(unique) {
			end = len(unique)
		}
		chunk := unique[start:end]
		placeholders := strings.TrimRight(strings.Repeat("?,", len(chunk)), ",")
		args := make([]any, 0, len(chunk))
		for _, magnet := range chunk {
			args = append(args, magnet)
		}
		rows, err := tx.Query("SELECT taskHash FROM rss_items WHERE magnet IN ("+placeholders+") AND taskHash != ''", args...)
		if err != nil {
			return 0, err
		}
		for rows.Next() {
			var hash string
			if err := rows.Scan(&hash); err != nil {
				rows.Close()
				return 0, err
			}
			if hash != "" {
				hashes = append(hashes, hash)
			}
		}
		rows.Close()
	}
	for start := 0; start < len(hashes); start += 500 {
		end := start + 500
		if end > len(hashes) {
			end = len(hashes)
		}
		chunk := hashes[start:end]
		placeholders := strings.TrimRight(strings.Repeat("?,", len(chunk)), ",")
		args := make([]any, 0, len(chunk))
		for _, hash := range chunk {
			args = append(args, hash)
		}
		if _, err := tx.Exec("DELETE FROM offline_postprocess WHERE task_hash IN ("+placeholders+")", args...); err != nil {
			return 0, err
		}
	}
	var deleted int64
	for start := 0; start < len(unique); start += 500 {
		end := start + 500
		if end > len(unique) {
			end = len(unique)
		}
		chunk := unique[start:end]
		placeholders := strings.TrimRight(strings.Repeat("?,", len(chunk)), ",")
		args := make([]any, 0, len(chunk))
		for _, magnet := range chunk {
			args = append(args, magnet)
		}
		res, err := tx.Exec("DELETE FROM rss_items WHERE magnet IN ("+placeholders+")", args...)
		if err != nil {
			return 0, err
		}
		n, _ := res.RowsAffected()
		deleted += n
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return deleted, nil
}

// TaskHashesForMagnets returns taskHash values for matching history rows.
func (s *Store) TaskHashesForMagnets(magnets []string) (map[string]string, error) {
	result := make(map[string]string)
	if len(magnets) == 0 {
		return result, nil
	}
	for start := 0; start < len(magnets); start += 500 {
		end := start + 500
		if end > len(magnets) {
			end = len(magnets)
		}
		chunk := magnets[start:end]
		placeholders := strings.TrimRight(strings.Repeat("?,", len(chunk)), ",")
		args := make([]any, 0, len(chunk))
		for _, magnet := range chunk {
			args = append(args, magnet)
		}
		rows, err := s.DBInstance.Query("SELECT magnet, taskHash FROM rss_items WHERE magnet IN ("+placeholders+") AND taskHash != ''", args...)
		if err != nil {
			return result, err
		}
		for rows.Next() {
			var magnet, hash string
			if err := rows.Scan(&magnet, &hash); err != nil {
				rows.Close()
				return result, err
			}
			result[magnet] = hash
		}
		rows.Close()
	}
	return result, nil
}

// HistoryRow is the minimal history data needed to reconstruct final names.
type HistoryRow struct {
	Magnet     string
	Title      string
	ConfigName string
	ConfigURL  string
	TaskHash   string
}

// HistoryRowsForMagnets returns local history rows for the given magnets.
func (s *Store) HistoryRowsForMagnets(magnets []string) ([]HistoryRow, error) {
	if len(magnets) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(magnets)), ",")
	args := make([]any, 0, len(magnets))
	for _, magnet := range magnets {
		args = append(args, magnet)
	}
	rows, err := s.DBInstance.Query("SELECT COALESCE(magnet,''), COALESCE(title,''), COALESCE(configName,''), COALESCE(configUrl,''), COALESCE(taskHash,'') FROM rss_items WHERE magnet IN ("+placeholders+")", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []HistoryRow
	for rows.Next() {
		var row HistoryRow
		if err := rows.Scan(&row.Magnet, &row.Title, &row.ConfigName, &row.ConfigURL, &row.TaskHash); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

// LatestHistoryForFeed returns the newest local history time for a feed.
func (s *Store) LatestHistoryForFeed(feedURL string) (time.Time, bool, error) {
	var value time.Time
	err := s.DBInstance.QueryRow("SELECT createdAt FROM rss_items WHERE configUrl = ? ORDER BY id DESC LIMIT 1", feedURL).Scan(&value)
	if err == sql.ErrNoRows {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}
	return value, true, nil
}

// MarkIdleReminder records that a feed was reminded on a day. It returns true
// only for the first call of that feed/day.
func (s *Store) MarkIdleReminder(feedURL, day string) (bool, error) {
	res, err := s.DBInstance.Exec("INSERT OR IGNORE INTO idle_reminders (feed_url, day, sent_at) VALUES (?, ?, ?)", feedURL, day, time.Now())
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// FailedTaskHashesForFeed returns 115 hashes of failed offline/postprocess
// rows for a feed. These tasks are removed before retrying the feed.
func (s *Store) FailedTaskHashesForFeed(configURL, configName string) ([]string, error) {
	query := "SELECT DISTINCT o.task_hash FROM offline_postprocess o JOIN rss_items r ON r.taskHash = o.task_hash WHERE o.state = 'failed'"
	var args []any
	if strings.TrimSpace(configURL) != "" {
		query += " AND r.configUrl = ?"
		args = append(args, configURL)
	} else if strings.TrimSpace(configName) != "" {
		query += " AND r.configName = ?"
		args = append(args, configName)
	} else {
		return nil, nil
	}
	rows, err := s.DBInstance.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var hashes []string
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, err
		}
		if strings.TrimSpace(hash) != "" {
			hashes = append(hashes, hash)
		}
	}
	return hashes, rows.Err()
}

// DeleteTaskRecords removes local history and postprocess rows for hashes.
func (s *Store) DeleteTaskRecords(hashes []string) error {
	if len(hashes) == 0 {
		return nil
	}
	tx, err := s.DBInstance.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, hash := range hashes {
		if strings.TrimSpace(hash) == "" {
			continue
		}
		if _, err := tx.Exec("DELETE FROM rss_items WHERE taskHash = ?", hash); err != nil {
			return err
		}
		if _, err := tx.Exec("DELETE FROM offline_postprocess WHERE task_hash = ?", hash); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// FinishedEntry associates one final 115 file/folder with a completed task.
type FinishedEntry struct {
	TaskHash   string
	ConfigName string
	Title      string
	EntryID    string
	EntryName  string
}

// FinishedEntries returns final entries for completed post-process tasks.
func (s *Store) FinishedEntries() ([]FinishedEntry, error) {
	rows, err := s.DBInstance.Query("SELECT o.task_hash, COALESCE(r.configName,''), COALESCE(r.title,''), COALESCE(o.final_entries,'') FROM offline_postprocess o LEFT JOIN rss_items r ON r.taskHash = o.task_hash WHERE o.state = 'done' AND o.final_entries IS NOT NULL AND o.final_entries != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []FinishedEntry
	for rows.Next() {
		var taskHash, configName, title, raw string
		if err := rows.Scan(&taskHash, &configName, &title, &raw); err != nil {
			return nil, err
		}
		var entries []PostEntry
		if err := json.Unmarshal([]byte(raw), &entries); err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.ID == "" {
				continue
			}
			result = append(result, FinishedEntry{TaskHash: taskHash, ConfigName: configName, Title: title, EntryID: entry.ID, EntryName: entry.Name})
		}
	}
	return result, rows.Err()
}
