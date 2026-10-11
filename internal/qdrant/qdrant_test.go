package qdrant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQdrantClient(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("all systems operational"))
	})

	mux.HandleFunc("/collections/wuwa_lore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Collection info
			res := map[string]interface{}{
				"result": map[string]interface{}{
					"status":        "green",
					"points_count":  42,
					"vectors_count": 42,
				},
			}
			json.NewEncoder(w).Encode(res)
		} else if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result": true, "status": "ok"}`))
		}
	})

	mux.HandleFunc("/collections/wuwa_lore/points", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": {"operation_id": 1, "status": "completed"}}`))
	})

	mux.HandleFunc("/collections/wuwa_lore/points/search", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"result": []map[string]interface{}{
				{
					"id":      "1234-5678",
					"version": 0,
					"score":   0.92,
					"payload": map[string]interface{}{
						"title":   "Simulacrum of the Heart",
						"content": "Rover speaks with Abby",
					},
				},
			},
			"status": "ok",
		}
		json.NewEncoder(w).Encode(res)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := NewClient(ts.URL)

	// Test Ping
	if err := client.Ping(context.Background()); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	// Test EnsureCollection
	if err := client.EnsureCollection(context.Background(), "wuwa_lore", 384); err != nil {
		t.Fatalf("EnsureCollection failed: %v", err)
	}

	// Test UpsertPoints
	points := []Point{
		{
			ID:     "1234-5678",
			Vector: make([]float32, 384),
			Payload: map[string]interface{}{
				"title":   "Simulacrum of the Heart",
				"content": "Rover speaks with Abby",
			},
		},
	}
	if err := client.UpsertPoints(context.Background(), "wuwa_lore", points, 10); err != nil {
		t.Fatalf("UpsertPoints failed: %v", err)
	}

	// Test Search
	results, err := client.Search(context.Background(), "wuwa_lore", make([]float32, 384), 5, 0.5)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	if results[0].Score < 0.9 {
		t.Errorf("Expected score ~0.92, got %f", results[0].Score)
	}

	// Test GetCollectionInfo
	info, err := client.GetCollectionInfo(context.Background(), "wuwa_lore")
	if err != nil {
		t.Fatalf("GetCollectionInfo failed: %v", err)
	}
	if info.PointsCount != 42 {
		t.Errorf("Expected 42 points, got %d", info.PointsCount)
	}
}
