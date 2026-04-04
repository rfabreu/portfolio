package game

import (
	"encoding/json"
	"net/http"
)

// rawGameRequest is used for decoding so we can validate board length strictly.
type rawGameRequest struct {
	Board   []string `json:"board"`
	Session Session  `json:"session"`
}

type GameRequest struct {
	Board   Board   `json:"board"`
	Session Session `json:"session"`
}

type GameResponse struct {
	Move   int    `json:"move"`
	Board  Board  `json:"board"`
	Status string `json:"status"`
	Winner string `json:"winner"`
}

func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw rawGameRequest
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		if len(raw.Board) != 9 {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		var req GameRequest
		req.Session = raw.Session
		copy(req.Board[:], raw.Board)

		if err := validateBoard(req.Board); err != "" {
			http.Error(w, `{"error":"`+err+`"}`, http.StatusBadRequest)
			return
		}

		if req.Session.Wins < 0 || req.Session.Losses < 0 || req.Session.Draws < 0 {
			http.Error(w, `{"error":"session values must be non-negative"}`, http.StatusBadRequest)
			return
		}

		// Check if game is already over
		if checkWinner(req.Board) != "" || isBoardFull(req.Board) {
			http.Error(w, `{"error":"game is already over"}`, http.StatusBadRequest)
			return
		}

		// Validate turn: X count must equal O count + 1
		xCount, oCount := 0, 0
		for _, cell := range req.Board {
			if cell == "X" {
				xCount++
			} else if cell == "O" {
				oCount++
			}
		}
		if xCount != oCount+1 {
			http.Error(w, `{"error":"not AI turn"}`, http.StatusBadRequest)
			return
		}

		// AI makes a move
		move := SelectMove(req.Board, req.Session)
		board := req.Board
		board[move] = "O"

		// Determine game status after AI move
		status := "in_progress"
		winner := ""
		if chk := checkWinner(board); chk != "" {
			winner = chk
			if chk == "O" {
				status = "ai_win"
			} else {
				status = "player_win"
			}
		} else if isBoardFull(board) {
			status = "draw"
		}

		resp := GameResponse{
			Move:   move,
			Board:  board,
			Status: status,
			Winner: winner,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

func validateBoard(b Board) string {
	for _, cell := range b {
		if cell != "" && cell != "X" && cell != "O" {
			return "invalid board value"
		}
	}
	return ""
}
