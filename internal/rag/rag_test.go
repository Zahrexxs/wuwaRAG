package rag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ajieblr-dev/wuwaRAG/internal/qdrant"
)

func TestDeterministicEmbedder(t *testing.T) {
	embedder := NewHuggingFaceEmbedder("", "all-MiniLM-L6-v2", "", 384)
	vec, err := embedder.Embed(context.Background(), "Sentinel Hsin and the Waking Moon Festival")
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vec) != 384 {
		t.Fatalf("Expected 384 dimensions, got %d", len(vec))
	}

	// Verify normalization
	var sumSq float64
	for _, val := range vec {
		sumSq += float64(val * val)
	}
	if sumSq < 0.99 || sumSq > 1.01 {
		t.Errorf("Expected normalized vector magnitude ~1.0, got %f", sumSq)
	}
}

func TestRAGService(t *testing.T) {
	// Mock Qdrant server
	qdrantMux := http.NewServeMux()
	qdrantMux.HandleFunc("/collections/wuwa_lore/points/search", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"result": []map[string]interface{}{
				{
					"id":      "chunk-1",
					"version": 0,
					"score":   0.88,
					"payload": map[string]interface{}{
						"title":      "Simulacrum of the Heart",
						"section":    "Dialogue",
						"subsection": "Talk to tea stall owner",
						"speakers":   []string{"Rover", "Tea Stall Owner"},
						"content":    "Tea Stall Owner: That'd be our new Moon Fox Rice Cakes, no question. They're Lady Hsin's favorite treat!",
					},
				},
			},
			"status": "ok",
		}
		json.NewEncoder(w).Encode(res)
	})
	qdrantServer := httptest.NewServer(qdrantMux)
	defer qdrantServer.Close()

	// Mock LiteLLM proxy server
	llmMux := http.NewServeMux()
	llmMux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"model": "wuwa-lore-model",
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Lady Hsin's favorite treat is Moon Fox Rice Cakes, as stated by the Tea Stall Owner in Xianling Village.",
					},
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens": 120,
				"total_tokens":  150,
			},
		}
		json.NewEncoder(w).Encode(res)
	})
	llmServer := httptest.NewServer(llmMux)
	defer llmServer.Close()

	qClient := qdrant.NewClient(qdrantServer.URL)
	embedder := NewHuggingFaceEmbedder("", "", "", 384)
	llmClient := NewLLMClient(llmServer.URL, "dummy-key", "wuwa-lore-model", "", "")

	svc := NewService(qClient, embedder, llmClient, "wuwa_lore", 3, 0.1)

	res, err := svc.AnswerQuestion(context.Background(), "What is Lady Hsin's favorite treat?")
	if err != nil {
		t.Fatalf("AnswerQuestion failed: %v", err)
	}

	if len(res.Sources) != 1 {
		t.Errorf("Expected 1 source snippet, got %d", len(res.Sources))
	}
	if res.Model != "wuwa-lore-model" {
		t.Errorf("Expected model 'wuwa-lore-model', got %s", res.Model)
	}
	if res.Answer == "" {
		t.Errorf("Expected non-empty answer")
	}
}
