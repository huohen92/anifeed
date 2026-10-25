package webui

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBangumiCover(t *testing.T) {
	old := bangumiAPIBase
	defer func() { bangumiAPIBase = old }()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v0/subjects/486456" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"name":"Japanese","name_cn":"中文名","images":{"large":"https://img.example/cover.jpg"},"rating":{"score":7.4}}`)
	}))
	defer server.Close()
	bangumiAPIBase = server.URL
	info, err := fetchBangumiCover("486456")
	if err != nil {
		t.Fatal(err)
	}
	if info.Cover != "https://img.example/cover.jpg" || info.Title != "中文名" || info.Score != 7.4 {
		t.Fatalf("info = %+v", info)
	}
}
