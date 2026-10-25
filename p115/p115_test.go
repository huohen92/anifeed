package p115

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zhifengle/anifeed/rsssite"
)

func TestSuccessfulMagnetItemsPreservesHashAlignment(t *testing.T) {
	items := []rsssite.MagnetItem{{Title: "one"}, {Title: "two"}, {Title: "three"}}
	var rejected []string
	gotItems, gotHashes := successfulMagnetItems(items, []string{"hash-one", "", "hash-three"}, func(item rsssite.MagnetItem) {
		rejected = append(rejected, item.Title)
	})
	if len(gotItems) != 2 || gotItems[0].Title != "one" || gotItems[1].Title != "three" {
		t.Fatalf("unexpected successful items: %#v", gotItems)
	}
	if len(gotHashes) != 2 || gotHashes[0] != "hash-one" || gotHashes[1] != "hash-three" {
		t.Fatalf("unexpected successful hashes: %#v", gotHashes)
	}
	if len(rejected) != 1 || rejected[0] != "two" {
		t.Fatalf("unexpected rejected items: %#v", rejected)
	}
}

func TestLoadCookiesFromUserConfigDir(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")

	configDir := filepath.Join(homeDir, ".config", "anifeed")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, ".cookies"), []byte("UID=1; CID=2; SEID=3; KID=4"), 0o600); err != nil {
		t.Fatalf("failed to write cookies: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	cookies := LoadCookies()
	if cookies != "UID=1; CID=2; SEID=3; KID=4" {
		t.Fatalf("unexpected cookies: %q", cookies)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestResilientHTTPClientConvertsNetworkError(t *testing.T) {
	client := &resilientHTTPClient{inner: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	})}}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	defer resp.Body.Close()
	var payload struct {
		State   bool   `json:"state"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.State || !strings.Contains(payload.Message, "boom") {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestResilientHTTPClientFixesNilBody(t *testing.T) {
	client := &resilientHTTPClient{inner: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: nil, Header: http.Header{}}, nil
	})}}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Body == nil {
		t.Fatal("body is nil")
	}
	resp.Body.Close()
}
