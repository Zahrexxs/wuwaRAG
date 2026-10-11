# Project: Wuthering Waves Lore RAG System

## 1. Project Overview
A Retrieval-Augmented Generation (RAG) RESTful API web application designed to answer questions based on the Wuthering Waves game lore. The system specifically targets the "Simulacrum of the Heart" main quest. It features an automated data ingestion pipeline from the Fandom MediaWiki API, semantic search capabilities, and a robust LLM fallback mechanism to handle free-tier API rate limits gracefully.

## 2. Architecture & Tech Stack
*   **Backend Framework:** Golang (Gin or Fiber)
*   **Data Source:** MediaWiki API (`wutheringwaves.fandom.com/api.php`)
*   **Vector Database:** Qdrant (deployed via Docker)
*   **LLM Gateway & Fallback:** LiteLLM (deployed as a Docker proxy container)
*   **LLM Providers:** Google Gemini API (Primary), Groq API / HuggingFace (Fallback)
*   **Embedding Model:** `all-MiniLM-L6-v2` (via HuggingFace API or local)
*   **Deployment:** Docker & Docker Compose

## 3. Project Directory Structure
```text
.
├── cmd/
│   ├── api/                 # Main entry point for REST API
│   │   └── main.go
│   └── ingest/              # CLI tool for fetching and embedding wiki data
│       └── main.go
├── internal/
│   ├── config/              # Environment variables and configuration
│   ├── fetcher/             # MediaWiki API HTTP client
│   ├── parser/              # Regex and text cleaning logic for wikitext
│   ├── qdrant/              # Qdrant vector database client setup
│   ├── rag/                 # Retrieval and prompt generation logic
│   └── server/              # Gin/Fiber router and handlers
├── deploy/
│   ├── docker-compose.yml
│   └── litellm_config.yaml  # LiteLLM routing and fallback rules
├── go.mod
└── go.sum
