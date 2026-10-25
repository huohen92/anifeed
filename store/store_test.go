package store

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/huohen92/anifeed/rsssite"
	_ "github.com/mattn/go-sqlite3"
)

func TestStore(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	s := New(db)
	err := s.SaveMagnetItems([]rsssite.MagnetItem{
		{
			Title:       "test",
			Link:        "test",
			Magnet:      "magnet:?xt=urn:btih:aa",
			Description: "test",
			Content:     "test",
		},
		{
			Title:       "test",
			Link:        "test",
			Magnet:      "magnet:?xt=urn:btih:uniquehash&dn=test",
			Description: "test",
			Content:     "test",
		},
	}, "feed1", "https://example.com/rss", []string{"hash1", "hash2"})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	exists := s.HasItem("magnet:?xt=urn:btih:aa")
	if !exists {
		t.Errorf("Error: %v", err)
	}
	exists = s.HasMagnetByXt("magnet:?xt=urn:btih:uniquehash&dn=test2")
	if !exists {
		t.Errorf("Error: %v", err)
	}
	notExists := s.HasItem("magnet:?xt=urn:btih:bb")
	if notExists {
		t.Errorf("Error: %v", err)
	}
}

func TestNewWithPathCreatesParentDir(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "nested", "db.sqlite")
	s, err := NewWithPath(dbPath)
	if err != nil {
		t.Fatalf("NewWithPath failed: %v", err)
	}
	defer s.Close()

	if s.DBInstance == nil {
		t.Fatal("expected DBInstance to be set")
	}
	if err := s.SaveMagnetItems([]rsssite.MagnetItem{{
		Title:  "test",
		Link:   "test",
		Magnet: "magnet:?xt=urn:btih:path",
	}}, "feed2", "https://example.com/rss2", []string{"h3"}); err != nil {
		t.Fatalf("SaveMagnetItems failed: %v", err)
	}
}

func TestMigrateSchemaAddsHistoryColumns(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	New(db) // runs initSchema + migrateSchema
	cols, err := tableColumns(db, "rss_items")
	if err != nil {
		t.Fatalf("tableColumns: %v", err)
	}
	for _, col := range []string{"configName", "configUrl", "taskHash"} {
		if !cols[col] {
			t.Errorf("column %s missing after migration", col)
		}
	}
	// idempotent: running again must not error
	if err := migrateSchema(db); err != nil {
		t.Errorf("migrateSchema second run: %v", err)
	}
}

func TestMigrateSchemaOnOldTable(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	// Simulate a table created by a previous program version (no history columns)
	if _, err := db.Exec("CREATE TABLE rss_items (id INTEGER PRIMARY KEY AUTOINCREMENT, link VARCHAR(255), title VARCHAR(255), guid VARCHAR(255), pubDate DATETIME, creator VARCHAR(255), summary TEXT, content VARCHAR(255), isoDate DATETIME, categories VARCHAR(255), contentSnippet VARCHAR(255), done TINYINT(1) DEFAULT 0, magnet VARCHAR(255) NOT NULL, createdAt DATETIME NOT NULL, updatedAt DATETIME NOT NULL)"); err != nil {
		t.Fatalf("create old table: %v", err)
	}
	if err := migrateSchema(db); err != nil {
		t.Fatalf("migrateSchema: %v", err)
	}
	s := &Store{DBInstance: db}
	if err := s.SaveMagnetItems([]rsssite.MagnetItem{{Title: "x", Magnet: "magnet:?xt=urn:btih:01"}}, "feedA", "https://a.example/rss", []string{"h1"}); err != nil {
		t.Fatalf("SaveMagnetItems after migration: %v", err)
	}
}

func TestSaveMagnetItemsRollsBackWholeChunk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	s := New(db)
	defer s.Close()
	if _, err := db.Exec(`CREATE TRIGGER reject_bad_magnet
		BEFORE INSERT ON rss_items
		WHEN NEW.magnet = 'bad'
		BEGIN
			SELECT RAISE(ABORT, 'rejected by test');
		END`); err != nil {
		t.Fatalf("create trigger: %v", err)
	}

	err = s.SaveMagnetItems([]rsssite.MagnetItem{
		{Title: "good", Magnet: "good"},
		{Title: "bad", Magnet: "bad"},
	}, "feed", "https://example.com/rss", []string{"hash-good", "hash-bad"})
	if err == nil {
		t.Fatal("expected SaveMagnetItems to fail")
	}
	var count int
	if err := db.QueryRow("SELECT count(*) FROM rss_items").Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("transaction left %d rows, want 0", count)
	}
}

func TestSaveMagnetItemsRejectsEmptyHash(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	s := New(db)
	defer s.Close()

	err = s.SaveMagnetItems(
		[]rsssite.MagnetItem{{Title: "failed", Magnet: "magnet:?xt=urn:btih:failed"}},
		"feed", "https://example.com/rss", []string{""},
	)
	if err == nil {
		t.Fatal("expected empty task hash to be rejected")
	}
	var count int
	if err := db.QueryRow("SELECT count(*) FROM rss_items").Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("empty hash inserted %d rows, want 0", count)
	}
}

func TestListHistory(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	s := New(db)
	items := []rsssite.MagnetItem{
		{Title: "first", Magnet: "magnet:?xt=urn:btih:01"},
		{Title: "second", Magnet: "magnet:?xt=urn:btih:02"},
		{Title: "third", Magnet: "magnet:?xt=urn:btih:03"},
	}
	if err := s.SaveMagnetItems(items[:2], "feedA", "https://a.example/rss", []string{"hA1", "hA2"}); err != nil {
		t.Fatalf("SaveMagnetItems A: %v", err)
	}
	if err := s.SaveMagnetItems(items[2:], "feedB", "https://b.example/rss", []string{"hB1"}); err != nil {
		t.Fatalf("SaveMagnetItems B: %v", err)
	}
	// all, newest first
	all, err := s.ListHistory(0, "")
	if err != nil {
		t.Fatalf("ListHistory: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len(all) = %d, want 3", len(all))
	}
	if all[0].Title != "third" || all[0].ConfigName != "feedB" {
		t.Errorf("newest = %+v, want third/feedB", all[0])
	}
	if all[0].TaskHash != "hB1" {
		t.Errorf("TaskHash = %q, want hB1", all[0].TaskHash)
	}
	// filter by feed name
	a, err := s.ListHistory(0, "feedA")
	if err != nil {
		t.Fatalf("ListHistory feedA: %v", err)
	}
	if len(a) != 2 {
		t.Fatalf("len(feedA) = %d, want 2", len(a))
	}
	// limit
	lim, err := s.ListHistory(1, "")
	if err != nil {
		t.Fatalf("ListHistory limit: %v", err)
	}
	if len(lim) != 1 {
		t.Fatalf("len(limit1) = %d, want 1", len(lim))
	}
}

func TestListHistoryNullColumns(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	New(db) // initSchema + migrateSchema
	// Insert a row the way the old code did: no configName/configUrl/taskHash
	// (these default to NULL), simulating pre-migration records.
	if _, err := db.Exec("INSERT INTO rss_items (link, title, magnet, done, createdAt, updatedAt) VALUES (?,?,?,?,?,?)",
		"link", "legacy item", "magnet:?xt=urn:btih:legacy", 0, time.Now(), time.Now()); err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}
	items, err := (&Store{DBInstance: db}).ListHistory(0, "")
	if err != nil {
		t.Fatalf("ListHistory with NULL columns: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].ConfigName != "" || items[0].TaskHash != "" {
		t.Errorf("NULL columns = %q/%q, want empty", items[0].ConfigName, items[0].TaskHash)
	}
	if items[0].Title != "legacy item" {
		t.Errorf("Title = %q, want %q", items[0].Title, "legacy item")
	}
}

func TestPostProcessLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	s := New(db)
	defer s.Close()
	items := []rsssite.MagnetItem{{Title: "[Sub] Anime - 01", Magnet: "magnet:?xt=urn:btih:abc"}}
	if err := s.SaveMagnetItemsWithPostprocess(items, "番剧", "https://example.com/rss", []string{"hash-post"}, "cid-1", 0, 0); err != nil {
		t.Fatal(err)
	}
	pending, err := s.ListPendingPostTasks(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 1 {
		t.Fatalf("pending = %+v", pending)
	}
	if pending[0].TaskHash != "hash-post" || pending[0].ConfigName != "番剧" || pending[0].OfflineCid != "cid-1" {
		t.Fatalf("pending row = %+v", pending[0])
	}
	count, err := s.CountPendingPostTasks()
	if err != nil || count != 1 {
		t.Fatalf("count = %d, err = %v", count, err)
	}
	if err := s.CompletePostTask("hash-post"); err != nil {
		t.Fatal(err)
	}
	if count, _ := s.CountPendingPostTasks(); count != 0 {
		t.Fatalf("count after complete = %d", count)
	}
}

func TestPostProcessFailureRetries(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	s := New(db)
	defer s.Close()
	if err := s.SaveMagnetItemsWithPostprocess([]rsssite.MagnetItem{{Title: "x", Magnet: "magnet:?xt=urn:btih:x"}}, "a", "https://a", []string{"h"}, "", 0, 0); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := s.FailPostTask("h", "temporary", 3); err != nil {
			t.Fatal(err)
		}
	}
	if count, _ := s.CountPendingPostTasks(); count != 1 {
		t.Fatalf("should still be pending after 2 attempts, count=%d", count)
	}
	if err := s.FailPostTask("h", "final", 3); err != nil {
		t.Fatal(err)
	}
	if count, _ := s.CountPendingPostTasks(); count != 0 {
		t.Fatalf("should be failed after 3 attempts, count=%d", count)
	}
}

func TestPreviewDecisionsRoundTrip(t *testing.T) {
	st, err := NewWithPath(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	decisions := map[string]bool{"a": true, "b": false}
	if err := st.SavePreviewDecisions("https://example.com/rss", decisions); err != nil {
		t.Fatal(err)
	}
	got, err := st.PreviewDecisions("https://example.com/rss")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || !got["a"] || got["b"] {
		t.Fatalf("decisions = %+v", got)
	}
}

func TestDeleteDownloadedByMagnets(t *testing.T) {
	st, err := NewWithPath(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	_, err = st.DBInstance.Exec("INSERT INTO rss_items (title, magnet, taskHash, createdAt, updatedAt) VALUES (?, ?, ?, datetime('now'), datetime('now'))", "t", "magnet:?xt=urn:btih:a", "hash-a")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.DBInstance.Exec("INSERT INTO offline_postprocess (task_hash, title, config_name, offline_cid, state, created_at, updated_at) VALUES (?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))", "hash-a", "t", "feed", "cid"); err != nil {
		t.Fatal(err)
	}
	deleted, err := st.DeleteDownloadedByMagnets([]string{"magnet:?xt=urn:btih:a"})
	if err != nil || deleted != 1 {
		t.Fatalf("deleted=%d err=%v", deleted, err)
	}
	if st.HasItem("magnet:?xt=urn:btih:a") {
		t.Fatal("history still present")
	}
	var count int
	st.DBInstance.QueryRow("SELECT count(*) FROM offline_postprocess WHERE task_hash = 'hash-a'").Scan(&count)
	if count != 0 {
		t.Fatal("postprocess still present")
	}
}

func TestFinalEntriesRoundTrip(t *testing.T) {
	st, err := NewWithPath(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	_, err = st.DBInstance.Exec("INSERT INTO offline_postprocess (task_hash,title,config_name,offline_cid,state,created_at,updated_at) VALUES ('h1','t','feed','cid','pending',datetime('now'),datetime('now'))")
	if err != nil {
		t.Fatal(err)
	}
	entries := []PostEntry{{ID: "v1", Name: "video.mkv"}, {ID: "s1", Name: "video.sc.ass"}}
	if err := st.CompletePostTaskWithEntries("h1", entries); err != nil {
		t.Fatal(err)
	}
	got, err := st.FinalEntriesForTaskHashes([]string{"h1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got["h1"]) != 2 || got["h1"][0].ID != "v1" || got["h1"][1].Name != "video.sc.ass" {
		t.Fatalf("entries = %+v", got)
	}
}

func TestIdleReminderDedup(t *testing.T) {
	st, err := NewWithPath(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	first, err := st.MarkIdleReminder("https://example.com/rss", "2026-10-19")
	if err != nil || !first {
		t.Fatalf("first = %v err=%v", first, err)
	}
	second, err := st.MarkIdleReminder("https://example.com/rss", "2026-10-19")
	if err != nil || second {
		t.Fatalf("second = %v err=%v", second, err)
	}
}

func TestFailedTaskHashesAndDelete(t *testing.T) {
	st, err := NewWithPath(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	items := []rsssite.MagnetItem{{Link: "l", Title: "t", Magnet: "magnet:?xt=urn:btih:a"}}
	if err := st.SaveMagnetItemsWithPostprocess(items, "番剧", "https://example.com/rss", []string{"hash-a"}, "cid", 0, 0); err != nil {
		t.Fatal(err)
	}
	if err := st.FailPostTask("hash-a", "failed", 1); err != nil {
		t.Fatal(err)
	}
	hashes, err := st.FailedTaskHashesForFeed("https://example.com/rss", "")
	if err != nil || len(hashes) != 1 || hashes[0] != "hash-a" {
		t.Fatalf("hashes=%v err=%v", hashes, err)
	}
	if err := st.DeleteTaskRecords(hashes); err != nil {
		t.Fatal(err)
	}
	if st.HasItem("magnet:?xt=urn:btih:a") {
		t.Fatal("history row was not deleted")
	}
}
