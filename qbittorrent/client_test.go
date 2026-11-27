package qbittorrent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientAddAndVersion(t *testing.T) {
	var gotAdd url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			_, _ = w.Write([]byte("Ok."))
		case "/api/v2/app/version":
			_, _ = w.Write([]byte("v5.0.0"))
		case "/api/v2/torrents/add":
			_ = r.ParseForm()
			gotAdd = r.PostForm
			_, _ = w.Write([]byte("Ok."))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client, err := New(srv.URL, "admin", "pass")
	if err != nil {
		t.Fatal(err)
	}
	version, err := client.Version(context.Background())
	if err != nil || version != "v5.0.0" {
		t.Fatalf("version = %q, err = %v", version, err)
	}
	if err := client.Add(context.Background(), "magnet:?xt=urn:btih:abc", "/downloads", "anifeed"); err != nil {
		t.Fatal(err)
	}
	if gotAdd.Get("urls") != "magnet:?xt=urn:btih:abc" || gotAdd.Get("savepath") != "/downloads" || gotAdd.Get("category") != "anifeed" {
		t.Fatalf("add form = %v", gotAdd)
	}
}

func TestClientAddRejectsFailsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/auth/login" {
			_, _ = w.Write([]byte("Ok."))
			return
		}
		_, _ = w.Write([]byte("Fails."))
	}))
	defer srv.Close()
	client, err := New(srv.URL, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Add(context.Background(), "magnet:?xt=urn:btih:abc", "", ""); err == nil {
		t.Fatal("expected add error for Fails. response")
	}
}

func TestClientTorrents(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			_, _ = w.Write([]byte("Ok."))
		case "/api/v2/torrents/info":
			if r.URL.Query().Get("category") != "anifeed" {
				t.Errorf("category = %q", r.URL.Query().Get("category"))
			}
			_, _ = w.Write([]byte(`[{"hash":"abc","name":"episode 1","progress":1,"state":"uploading"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	client, err := New(srv.URL, "", "")
	if err != nil {
		t.Fatal(err)
	}
	torrents, err := client.Torrents(context.Background(), "anifeed")
	if err != nil {
		t.Fatal(err)
	}
	if len(torrents) != 1 || torrents[0].Hash != "abc" || torrents[0].Progress != 1 {
		t.Fatalf("torrents = %#v", torrents)
	}
}

func TestAPIKeyAndAddOptions(t *testing.T) {
	var got url.Values
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/torrents/add" {
			auth = r.Header.Get("Authorization")
			_ = r.ParseForm()
			got = r.PostForm
			_, _ = w.Write([]byte("Ok."))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()
	client, err := NewWithAPIKey(srv.URL, "qbt_test")
	if err != nil {
		t.Fatal(err)
	}
	err = client.AddWithOptions(context.Background(), "magnet:?xt=urn:btih:abc", AddOptions{
		SavePath: "/downloads", Category: "anifeed", Tags: "a,b", ContentLayout: "Subfolder",
		DownloadLimitKB: 100, UploadLimitKB: 50, RatioLimit: 2, SeedingTimeLimit: 60,
		InactiveSeedingTimeLimit: 30, UseDownloadPath: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if auth != "Bearer qbt_test" {
		t.Fatalf("authorization = %q", auth)
	}
	if got.Get("tags") != "a,b" || got.Get("contentLayout") != "Subfolder" ||
		got.Get("dlLimit") != "102400" || got.Get("upLimit") != "51200" ||
		got.Get("ratioLimit") != "2" || got.Get("seedingTimeLimit") != "60" ||
		got.Get("inactiveSeedingTimeLimit") != "30" || got.Get("useDownloadPath") != "true" {
		t.Fatalf("add form = %v", got)
	}
}
