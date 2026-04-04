package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rfabreu/portfolio-api/internal/chat"
	"github.com/rfabreu/portfolio-api/internal/game"
	"github.com/rfabreu/portfolio-api/internal/middleware"
)

func main() {
	port := getEnv("PORT", "8080")
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:4321")
	rpm, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "30"))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY not set — chat endpoint will fail")
	}
	provider := chat.NewGeminiProvider(geminiKey, "")
	mux.Handle("POST /api/chat", chat.NewHandler(provider))
	mux.Handle("POST /api/game", game.NewHandler())

	handler := middleware.CORS(allowedOrigins)(middleware.RateLimit(rpm)(mux))

	log.Printf("API server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
