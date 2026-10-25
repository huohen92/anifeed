package webui

import (
	"context"
	"strings"
	"testing"

	"github.com/deadblue/elevengo"
	"github.com/zhifengle/anifeed/cloudfs"
	"github.com/zhifengle/anifeed/config"
	"github.com/zhifengle/anifeed/store"
)

type fakePostFS struct {
	entry       cloudfs.Entry
	children    []cloudfs.Entry
	renameTo    string
	renamed     map[string]string
	moveTo      string
	movedIDs    []string
	deletedIDs  []string
	lookupCalls int
}

func (f *fakePostFS) Stat(context.Context, string) (cloudfs.Entry, error) { return f.entry, nil }
func (f *fakePostFS) Lookup(_ context.Context, _ string, name string) (cloudfs.Entry, error) {
	f.lookupCalls++
	return cloudfs.Entry{ID: "found", Name: name, ParentID: "offline", Type: cloudfs.EntryTypeFile}, nil
}
func (f *fakePostFS) Rename(_ context.Context, id, name string) (cloudfs.Entry, error) {
	f.renameTo = name
	if f.renamed == nil {
		f.renamed = map[string]string{}
	}
	f.renamed[id] = name
	f.entry.Name = name
	for i := range f.children {
		if f.children[i].ID == id {
			f.children[i].Name = name
		}
	}
	return f.entry, nil
}
func (f *fakePostFS) List(_ context.Context, _ string) ([]cloudfs.Entry, error) {
	return f.children, nil
}
func (f *fakePostFS) Search(_ context.Context, _ string, _ string, _ cloudfs.SearchOptions) ([]cloudfs.Entry, error) {
	return f.children, nil
}
func (f *fakePostFS) Delete(_ context.Context, entryID string) error {
	f.deletedIDs = append(f.deletedIDs, entryID)
	remaining := f.children[:0]
	for _, child := range f.children {
		if child.ID != entryID {
			remaining = append(remaining, child)
		}
	}
	f.children = remaining
	return nil
}
func (f *fakePostFS) Move(_ context.Context, target, entryID string) (cloudfs.Entry, error) {
	f.moveTo = target
	f.movedIDs = append(f.movedIDs, entryID)
	remaining := f.children[:0]
	for _, child := range f.children {
		if child.ID != entryID {
			remaining = append(remaining, child)
		}
	}
	f.children = remaining
	f.entry.ParentID = target
	return f.entry, nil
}

func postConfig(rename, dual bool) *config.Config {
	return &config.Config{
		Offline: config.OfflineConfig{DualFolder: dual, FinalCid: "final-cid"},
		Rename:  config.RenameConfig{Enabled: rename, Template: "[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}", Season: 1},
	}
}

func TestProcessPostTaskRenamesAndMoves(t *testing.T) {
	fs := &fakePostFS{entry: cloudfs.Entry{ID: "f1", Name: "[Sub] Anime - 01.mkv", ParentID: "offline", Type: cloudfs.EntryTypeFile}}
	task := &elevengo.OfflineTask{FileId: "f1", Name: "[Sub] Anime - 01.mkv"}
	pending := store.PostTask{Title: "[Sub] Anime - 01.mkv", ConfigName: "测试番剧", OfflineCid: "offline"}
	if _, err := processPostTask(context.Background(), fs, task, pending, postConfig(true, true), ""); err != nil {
		t.Fatal(err)
	}
	if fs.renameTo != "[Sub] 测试番剧 S01E01.mkv" {
		t.Errorf("renameTo = %q", fs.renameTo)
	}
	if fs.moveTo != "final-cid" {
		t.Errorf("moveTo = %q", fs.moveTo)
	}
}

func TestProcessPostTaskMovesWithoutRename(t *testing.T) {
	fs := &fakePostFS{entry: cloudfs.Entry{ID: "f1", Name: "folder", ParentID: "offline", Type: cloudfs.EntryTypeDirectory}}
	task := &elevengo.OfflineTask{FileId: "f1", Name: "folder"}
	pending := store.PostTask{Title: "[Sub] Anime [1080p]", ConfigName: "番剧"}
	if _, err := processPostTask(context.Background(), fs, task, pending, postConfig(false, true), ""); err != nil {
		t.Fatal(err)
	}
	if fs.moveTo != "final-cid" || fs.renameTo != "" {
		t.Errorf("moveTo=%q renameTo=%q", fs.moveTo, fs.renameTo)
	}
}

func TestProcessPostTaskDoesNotMoveAfterRenameFailure(t *testing.T) {
	fs := &fakePostFS{entry: cloudfs.Entry{ID: "f1", Name: "folder", ParentID: "offline", Type: cloudfs.EntryTypeDirectory}}
	task := &elevengo.OfflineTask{FileId: "f1", Name: "folder"}
	pending := store.PostTask{Title: "[Sub] Anime [1080p]", ConfigName: "番剧"}
	_, err := processPostTask(context.Background(), fs, task, pending, postConfig(true, true), "")
	if err == nil || !strings.Contains(err.Error(), "集数") {
		t.Fatalf("err = %v", err)
	}
	if fs.moveTo != "" {
		t.Errorf("moved after rename failure: %q", fs.moveTo)
	}
}

func TestProcessPostTaskFallsBackToLookup(t *testing.T) {
	fs := &fakePostFS{entry: cloudfs.Entry{}}
	task := &elevengo.OfflineTask{Name: "folder", FileId: ""}
	pending := store.PostTask{Title: "folder", ConfigName: "番剧", OfflineCid: "offline"}
	// A zero-value entry makes Stat fail the ID check; Lookup returns the entry.
	fake := &lookupOnlyFS{fakePostFS: fs}
	if _, err := processPostTask(context.Background(), fake, task, pending, postConfig(false, true), ""); err != nil {
		t.Fatal(err)
	}
	if fake.lookupCalls != 1 || fake.moveTo != "final-cid" {
		t.Errorf("lookups=%d moveTo=%q", fake.lookupCalls, fake.moveTo)
	}
}

type lookupOnlyFS struct {
	*fakePostFS
}

func (f *lookupOnlyFS) Stat(context.Context, string) (cloudfs.Entry, error) {
	return cloudfs.Entry{}, context.DeadlineExceeded
}

func TestProcessPostTaskRenamesVideoInsideFolder(t *testing.T) {
	folder := cloudfs.Entry{ID: "dir1", Name: "Anime S01E01", ParentID: "offline", Type: cloudfs.EntryTypeDirectory}
	video := cloudfs.Entry{ID: "v1", Name: "[Sub] Anime - 01.mkv", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	subtitle := cloudfs.Entry{ID: "s1", Name: "Anime.ass", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	fs := &fakePostFS{entry: folder, children: []cloudfs.Entry{video, subtitle}}
	task := &elevengo.OfflineTask{FileId: "dir1", Name: folder.Name}
	pending := store.PostTask{Title: "[Sub] Anime - 01.mkv", ConfigName: "测试番剧", OfflineCid: "offline"}
	if _, err := processPostTask(context.Background(), fs, task, pending, postConfig(true, true), ""); err != nil {
		t.Fatal(err)
	}
	if fs.renamed["v1"] != "[Sub] 测试番剧 S01E01.mkv" {
		t.Errorf("video rename = %q", fs.renamed["v1"])
	}
	if fs.renamed["s1"] != "[Sub] 测试番剧 S01E01.und.ass" {
		t.Errorf("subtitle rename = %q", fs.renamed["s1"])
	}
	if fs.moveTo != "final-cid" || len(fs.movedIDs) != 1 || fs.movedIDs[0] != "dir1" {
		t.Errorf("moveTo=%q moved=%v, folder should be moved", fs.moveTo, fs.movedIDs)
	}
}

func TestProcessPostTaskKeepOnlyVideo(t *testing.T) {
	folder := cloudfs.Entry{ID: "dir1", Name: "Anime S01E01", ParentID: "offline", Type: cloudfs.EntryTypeDirectory}
	video := cloudfs.Entry{ID: "v1", Name: "[Sub] Anime - 01.mkv", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	subtitle := cloudfs.Entry{ID: "s1", Name: "Anime.ass", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	fs := &fakePostFS{entry: folder, children: []cloudfs.Entry{video, subtitle}}
	task := &elevengo.OfflineTask{FileId: "dir1", Name: folder.Name}
	pending := store.PostTask{Title: "[Sub] Anime - 01.mkv", ConfigName: "测试番剧", OfflineCid: "offline"}
	cfg := postConfig(true, true)
	cfg.Offline.KeepOnlyVideo = true
	if _, err := processPostTask(context.Background(), fs, task, pending, cfg, ""); err != nil {
		t.Fatal(err)
	}
	if len(fs.movedIDs) != 2 || fs.movedIDs[0] != "v1" || fs.movedIDs[1] != "s1" {
		t.Errorf("movedIDs = %v, want video and subtitle", fs.movedIDs)
	}
	if len(fs.deletedIDs) != 1 || fs.deletedIDs[0] != "dir1" {
		t.Errorf("deletedIDs = %v, want empty source dir removed", fs.deletedIDs)
	}
}

func TestProcessPostTaskRenamesSubtitleLanguages(t *testing.T) {
	folder := cloudfs.Entry{ID: "dir1", Name: "Anime S01E01", ParentID: "offline", Type: cloudfs.EntryTypeDirectory}
	video := cloudfs.Entry{ID: "v1", Name: "Anime - 01.mkv", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	sc := cloudfs.Entry{ID: "s1", Name: "Anime.sc.ass", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	tc := cloudfs.Entry{ID: "s2", Name: "Anime.tc.ass", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	fs := &fakePostFS{entry: folder, children: []cloudfs.Entry{video, sc, tc}}
	task := &elevengo.OfflineTask{FileId: "dir1", Name: folder.Name}
	pending := store.PostTask{Title: "[简繁内封] Anime - 01.mkv", ConfigName: "测试番剧", OfflineCid: "offline"}
	cfg := postConfig(true, true)
	cfg.Offline.KeepOnlyVideo = true
	if _, err := processPostTask(context.Background(), fs, task, pending, cfg, ""); err != nil {
		t.Fatal(err)
	}
	if fs.renamed["s1"] != "[简繁内封] 测试番剧 S01E01.sc.ass" {
		t.Errorf("sc subtitle = %q", fs.renamed["s1"])
	}
	if fs.renamed["s2"] != "[简繁内封] 测试番剧 S01E01.tc.ass" {
		t.Errorf("tc subtitle = %q", fs.renamed["s2"])
	}
	if len(fs.movedIDs) != 3 {
		t.Errorf("movedIDs = %v, want video + 2 subtitles", fs.movedIDs)
	}
}

func TestProcessPostTaskKeepsDirWhenFilesRemain(t *testing.T) {
	folder := cloudfs.Entry{ID: "dir1", Name: "Anime S01E01", ParentID: "offline", Type: cloudfs.EntryTypeDirectory}
	video := cloudfs.Entry{ID: "v1", Name: "Anime - 01.mkv", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	font := cloudfs.Entry{ID: "f1", Name: "font.ttf", ParentID: "dir1", Type: cloudfs.EntryTypeFile}
	fs := &fakePostFS{entry: folder, children: []cloudfs.Entry{video, font}}
	task := &elevengo.OfflineTask{FileId: "dir1", Name: folder.Name}
	pending := store.PostTask{Title: "Anime - 01.mkv", ConfigName: "测试番剧", OfflineCid: "offline"}
	cfg := postConfig(true, true)
	cfg.Offline.KeepOnlyVideo = true
	if _, err := processPostTask(context.Background(), fs, task, pending, cfg, ""); err != nil {
		t.Fatal(err)
	}
	if len(fs.deletedIDs) != 0 {
		t.Errorf("deletedIDs = %v, want source dir kept", fs.deletedIDs)
	}
}

func TestMissingEpisodes(t *testing.T) {
	items := []renamePreviewItem{
		{Episode: 1, Status: "ready"},
		{Episode: 2, Status: "ready"},
		{Episode: 4, Status: "ready"},
		{Episode: 5, Status: "error"},
	}
	got := missingEpisodes(items)
	if len(got) != 1 || got[0] != 3 {
		t.Fatalf("missing = %v, want [3]", got)
	}
}

func TestRenameSeasonForFeed(t *testing.T) {
	cfg := &config.Config{Rename: config.RenameConfig{Season: 1}}
	cfg.RSS = map[string][]config.RssConfig{
		"anibt.net": {
			{Name: "佐佐木第二季", Url: "https://a", RenameSeason: 2},
			{Name: "女友的朋友", Url: "https://b"},
		},
	}
	if got := renameSeasonForFeed(cfg, "佐佐木第二季", "https://a"); got != 2 {
		t.Errorf("season = %d, want 2", got)
	}
	if got := renameSeasonForFeed(cfg, "女友的朋友", "https://b"); got != 1 {
		t.Errorf("season = %d, want 1", got)
	}
}

func TestMagnetInfoHash(t *testing.T) {
	hash := "64cbf15a2dab04f7547fd2da336b920e26e19d08"
	if got := magnetInfoHash("magnet:?xt=urn:btih:" + hash + "&dn=test"); got != hash {
		t.Fatalf("hash = %q", got)
	}
	if got := magnetInfoHash("not-a-magnet"); got != "" {
		t.Fatalf("invalid hash = %q", got)
	}
}

func TestRenameOffsetForFeed(t *testing.T) {
	cfg := &config.Config{Rename: config.RenameConfig{Offset: 1}}
	cfg.RSS = map[string][]config.RssConfig{
		"site": {
			{Name: "A", Url: "https://a", RenameOffset: 2},
			{Name: "B", Url: "https://b"},
		},
	}
	if got := renameOffsetForFeedConfig(cfg, "A"); got != 2 {
		t.Errorf("offset A = %d, want 2", got)
	}
	if got := renameOffsetForFeedConfig(cfg, "B"); got != 1 {
		t.Errorf("offset B = %d, want global 1", got)
	}
}

func TestDetectSeason(t *testing.T) {
	cases := map[string]int{
		"青之箱 第二季":          2,
		"Anime Season 3":   3,
		"Anime S2 - 01":    2,
		"Anime 2nd Season": 2,
		"青之箱":              0,
	}
	for input, want := range cases {
		if got := detectSeason(input); got != want {
			t.Errorf("detectSeason(%q) = %d, want %d", input, got, want)
		}
	}
}
