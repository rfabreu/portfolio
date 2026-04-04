package game_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rfabreu/portfolio-api/internal/game"
)

func TestHandler_ValidMove(t *testing.T) {
	handler := game.NewHandler()

	reqBody := game.GameRequest{
		Board:   [9]string{"X", "", "", "", "", "", "", "", ""},
		Session: game.Session{Wins: 0, Losses: 0, Draws: 0},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp game.GameResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Move < 0 || resp.Move > 8 {
		t.Errorf("move out of range: %d", resp.Move)
	}
	if resp.Board[resp.Move] != "O" {
		t.Errorf("expected O at move index %d, got %q", resp.Move, resp.Board[resp.Move])
	}
	if resp.Status == "" {
		t.Error("expected non-empty status")
	}
}

func TestHandler_InvalidBoardLength(t *testing.T) {
	handler := game.NewHandler()

	body := []byte(`{"board":["X","",""],"session":{"wins":0,"losses":0,"draws":0}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_InvalidBoardCharacters(t *testing.T) {
	handler := game.NewHandler()

	reqBody := game.GameRequest{
		Board:   [9]string{"X", "Z", "", "", "", "", "", "", ""},
		Session: game.Session{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_WrongTurn(t *testing.T) {
	handler := game.NewHandler()

	reqBody := game.GameRequest{
		Board:   [9]string{"X", "O", "", "", "", "", "", "", ""},
		Session: game.Session{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_GameAlreadyOver(t *testing.T) {
	handler := game.NewHandler()

	reqBody := game.GameRequest{
		Board:   [9]string{"X", "X", "X", "O", "O", "", "", "", ""},
		Session: game.Session{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_NegativeSessionValues(t *testing.T) {
	handler := game.NewHandler()

	reqBody := game.GameRequest{
		Board:   [9]string{"X", "", "", "", "", "", "", "", ""},
		Session: game.Session{Wins: -1, Losses: 0, Draws: 0},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
