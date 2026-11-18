package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/rfabreu/portfolio-api/internal/middleware"
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
	Messages messageList `json:"messages"`
}

var errTooManyMessages = errors.New("too many messages")

// messageList decodes one element at a time so an oversized array is
// rejected before every element is allocated.
type messageList []Message

func (m *messageList) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if tok == nil {
		return nil
	}
	if tok != json.Delim('[') {
		return errors.New("messages must be an array")
	}
	var msgs messageList
	for dec.More() {
		if len(msgs) == maxMessages {
			return errTooManyMessages
		}
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			return err
		}
		msgs = append(msgs, msg)
	}
	*m = msgs
	return nil
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
			if errors.Is(err, errTooManyMessages) {
				http.Error(w, `{"error":"conversation limit exceeded (max 20 messages)"}`, http.StatusBadRequest)
				return
			}
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if len(req.Messages) == 0 {
			http.Error(w, `{"error":"messages required"}`, http.StatusBadRequest)
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
			var capErr *CapExceededError
			if errors.As(err, &capErr) {
				middleware.WriteJSONError(w, http.StatusServiceUnavailable, "chat is unavailable right now (daily limit reached)", capErr.RetryAfter)
				return
			}
			log.Printf("Chat error: %v", err)
			http.Error(w, `{"error":"failed to generate response"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Message: response})
	})
}
