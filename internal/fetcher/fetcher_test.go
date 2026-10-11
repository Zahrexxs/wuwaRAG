package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchPage(t *testing.T) {
	mockResponse := `{
		"parse": {
			"title": "Simulacrum of the Heart",
			"pageid": 77596,
			"wikitext": {
				"*": "==Dialogue==\n:'''Rover:''' Test dialogue"
			},
			"sections": [
				{
					"toclevel": 1,
					"level": "2",
					"line": "Dialogue",
					"number": "1",
					"index": "1"
				}
			]
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	page, err := client.FetchPage(context.Background(), "Simulacrum of the Heart")
	if err != nil {
		t.Fatalf("FetchPage failed: %v", err)
	}

	if page.Title != "Simulacrum of the Heart" {
		t.Errorf("Expected title 'Simulacrum of the Heart', got: %s", page.Title)
	}
	if len(page.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(page.Sections))
	}
	if page.Wikitext != "==Dialogue==\n:'''Rover:''' Test dialogue" {
		t.Errorf("Unexpected wikitext: %s", page.Wikitext)
	}
}

func TestSearchPages(t *testing.T) {
	mockResponse := `{
		"query": {
			"prefixsearch": [
				{"title": "Simulacrum of the Heart", "pageid": 77596},
				{"title": "Moonlike Heart, Mortal Longing", "pageid": 77594}
			]
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	titles, err := client.SearchPages(context.Background(), "Simulacrum")
	if err != nil {
		t.Fatalf("SearchPages failed: %v", err)
	}

	if len(titles) != 2 {
		t.Fatalf("Expected 2 titles, got: %d", len(titles))
	}
	if titles[0] != "Simulacrum of the Heart" {
		t.Errorf("Expected first title to be 'Simulacrum of the Heart', got: %s", titles[0])
	}
}
