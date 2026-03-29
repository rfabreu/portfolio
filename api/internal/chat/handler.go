package chat

import (
	"encoding/json"
	"net/http"
)

const maxMessages = 20

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

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

		response, err := provider.GenerateResponse(r.Context(), SystemPrompt(), req.Messages)
		if err != nil {
			http.Error(w, `{"error":"failed to generate response"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Message: response})
	})
}
