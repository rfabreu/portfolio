package chat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"unicode/utf8"
)

const maxMessages = 20

const (
	// MaxUserContentRunes matches the ChatWidget input's maxLength.
	MaxUserContentRunes = 2000
	// MaxAssistantContentRunes leaves headroom over Gemini's 1024 output-token cap.
	MaxAssistantContentRunes = 8000
	// MaxBodyBytes fits 20 messages at the assistant cap even when every
	// character is JSON-escaped (6 bytes) or 4 bytes of UTF-8.
	MaxBodyBytes = 1 << 20
)

type ChatRequest struct {
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

func NewHandler(provider Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				http.Error(w, `{"error":"request body too large"}`, http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if len(req.Messages) == 0 {
			http.Error(w, `{"error":"messages required"}`, http.StatusBadRequest)
			return
		}

		if len(req.Messages) > maxMessages {
			http.Error(w, `{"error":"conversation limit exceeded (max 20 messages)"}`, http.StatusBadRequest)
			return
		}

		for _, msg := range req.Messages {
			var limit int
			switch msg.Role {
			case "user":
				limit = MaxUserContentRunes
			case "assistant":
				limit = MaxAssistantContentRunes
			default:
				http.Error(w, `{"error":"invalid message role"}`, http.StatusBadRequest)
				return
			}
			if utf8.RuneCountInString(msg.Content) > limit {
				http.Error(w, `{"error":"message too long"}`, http.StatusBadRequest)
				return
			}
		}

		response, err := provider.GenerateResponse(r.Context(), SystemPrompt(), req.Messages)
		if err != nil {
			log.Printf("Chat error: %v", err)
			http.Error(w, `{"error":"failed to generate response"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Message: response})
	})
}
