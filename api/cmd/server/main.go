package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rfabreu/portfolio-api/internal/chat"
	"github.com/rfabreu/portfolio-api/internal/game"
	"github.com/rfabreu/portfolio-api/internal/middleware"
)

func main() {
	port := getEnv("PORT", "8080")
	cfg, warnings := loadConfig(os.Getenv)
	cfg.allowedOrigins = getEnv("ALLOWED_ORIGINS", "http://localhost:4321")
	for _, w := range warnings {
		log.Printf("WARNING: %s", w)
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY not set — chat endpoint will fail")
	}
	provider := chat.NewGeminiProvider(geminiKey, "")

	log.Printf("Limits: chat %d/min and %d/%s per client, game %d/min per client, chat %d calls/%s in total; trust Fly-Client-IP: %t",
		cfg.chatRPM, cfg.chatClientDaily, chat.DailyCapWindow, cfg.gameRPM, cfg.chatDailyCap, chat.DailyCapWindow, cfg.trustFlyClientIP)
	log.Printf("API server starting on :%s", port)
	if err := newServer(":"+port, newHandler(cfg, provider)).ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// newHandler builds the routes. /health is unlimited so Fly health checks never
// fail; chat and game each get their own per-client limiter, and chat also
// draws from the global Gemini call budget. Limiters wrap the route handlers,
// so CORS preflights and 404/405 responses never reach them.
func newHandler(cfg config, provider chat.Provider) http.Handler {
	limit := func(n int, window time.Duration) func(http.Handler) http.Handler {
		return middleware.RateLimit(middleware.RateLimitConfig{
			Limit:            n,
			Window:           window,
			TrustFlyClientIP: cfg.trustFlyClientIP,
			Now:              cfg.now,
		})
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})
	capped := chat.NewCappedProvider(provider, cfg.chatDailyCap, chat.DailyCapWindow, cfg.now)
	// The per-client daily limit stops one client spending the whole global budget.
	chatLimits := func(h http.Handler) http.Handler {
		return limit(cfg.chatRPM, time.Minute)(limit(cfg.chatClientDaily, chat.DailyCapWindow)(h))
	}
	mux.Handle("POST /api/chat", chatLimits(chat.NewHandler(capped)))
	mux.Handle("POST /api/game", limit(cfg.gameRPM, time.Minute)(game.NewHandler()))

	return middleware.CORS(cfg.allowedOrigins)(mux)
}

func newServer(addr string, handler http.Handler) *http.Server {
	const readTimeout = 15 * time.Second
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       readTimeout,
		// The write deadline starts once headers are read, so it must cover
		// the body read plus the Gemini call for chat errors to reach the client.
		WriteTimeout: readTimeout + chat.RequestTimeout + 15*time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
