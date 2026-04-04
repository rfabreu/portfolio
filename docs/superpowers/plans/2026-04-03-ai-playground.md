# P2: AI Playground — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a `/playground` page where visitors play Tic-Tac-Toe against the Go API, with an adaptive engagement model that subtly controls difficulty to keep visitors playing.

**Architecture:** A stateless `POST /api/game` endpoint receives the board and session stats, computes the AI move via minimax with alpha-beta pruning and adaptive difficulty, returns the updated board. A React island (`GameBoard.tsx`) on a new Astro page (`/playground`) renders the game UI.

**Tech Stack:** Go 1.26 (stdlib net/http, no external deps), React 19, TypeScript, Astro 6

---

## File Structure

```
api/internal/game/
  ├── engine.go        — Board types, minimax, adaptive move selection
  ├── engine_test.go   — Unit tests for minimax and difficulty logic
  ├── handler.go       — HTTP handler, request/response structs, validation
  └── handler_test.go  — httptest-based handler tests

api/cmd/server/main.go — Add game route registration

frontend/src/
  ├── pages/playground.astro   — New page with GameBoard island
  ├── islands/GameBoard.tsx    — React island: game UI, state, API communication
  └── components/Nav.astro     — Add Playground nav item
```

---

### Task 1: Implement the game engine core — board types and win detection

**Files:**
- Create: `api/internal/game/engine.go`
- Create: `api/internal/game/engine_test.go`

- [ ] **Step 1: Write failing tests for win detection**

Create `api/internal/game/engine_test.go`:

```go
package game

import "testing"

func TestCheckWinner_XWinsRow(t *testing.T) {
	board := Board{"X", "X", "X", "", "O", "O", "", "", ""}
	winner := checkWinner(board)
	if winner != "X" {
		t.Errorf("expected X, got %q", winner)
	}
}

func TestCheckWinner_OWinsColumn(t *testing.T) {
	board := Board{"O", "X", "", "O", "X", "", "O", "", "X"}
	winner := checkWinner(board)
	if winner != "O" {
		t.Errorf("expected O, got %q", winner)
	}
}

func TestCheckWinner_OWinsDiagonal(t *testing.T) {
	board := Board{"O", "X", "", "X", "O", "", "", "", "O"}
	winner := checkWinner(board)
	if winner != "O" {
		t.Errorf("expected O, got %q", winner)
	}
}

func TestCheckWinner_NoWinner(t *testing.T) {
	board := Board{"X", "O", "X", "X", "O", "O", "O", "X", "X"}
	winner := checkWinner(board)
	if winner != "" {
		t.Errorf("expected empty, got %q", winner)
	}
}

func TestCheckWinner_EmptyBoard(t *testing.T) {
	board := Board{"", "", "", "", "", "", "", "", ""}
	winner := checkWinner(board)
	if winner != "" {
		t.Errorf("expected empty, got %q", winner)
	}
}

func TestIsBoardFull_Full(t *testing.T) {
	board := Board{"X", "O", "X", "X", "O", "O", "O", "X", "X"}
	if !isBoardFull(board) {
		t.Error("expected board to be full")
	}
}

func TestIsBoardFull_NotFull(t *testing.T) {
	board := Board{"X", "", "X", "", "O", "", "", "", ""}
	if isBoardFull(board) {
		t.Error("expected board to not be full")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/game/... -v`
Expected: FAIL — types and functions not defined.

- [ ] **Step 3: Implement board types and win detection**

Create `api/internal/game/engine.go`:

```go
package game

// Board is a flat 3x3 grid. Index 0-8, left-to-right, top-to-bottom.
// Values: "" (empty), "X" (player/human), "O" (AI).
type Board [9]string

// Session tracks game outcomes for adaptive difficulty.
type Session struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
	Draws  int `json:"draws"`
}

var winLines = [][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // columns
	{0, 4, 8}, {2, 4, 6},             // diagonals
}

func checkWinner(b Board) string {
	for _, line := range winLines {
		a, bc, c := b[line[0]], b[line[1]], b[line[2]]
		if a != "" && a == bc && bc == c {
			return a
		}
	}
	return ""
}

func isBoardFull(b Board) bool {
	for _, cell := range b {
		if cell == "" {
			return false
		}
	}
	return true
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/game/... -v`
Expected: All 7 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add api/internal/game/engine.go api/internal/game/engine_test.go
git commit -m "feat(game): add board types and win detection"
```

---

### Task 2: Implement minimax with alpha-beta pruning

**Files:**
- Modify: `api/internal/game/engine.go`
- Modify: `api/internal/game/engine_test.go`

- [ ] **Step 1: Write failing tests for minimax**

Append to `api/internal/game/engine_test.go`:

```go
func TestMinimax_TakesWinningMove(t *testing.T) {
	// O can win at index 2
	board := Board{"X", "X", "", "O", "O", "", "", "", "X"}
	move := bestMove(board)
	if move != 5 {
		t.Errorf("expected AI to win at 5, got %d", move)
	}
}

func TestMinimax_BlocksOpponentWin(t *testing.T) {
	// X is about to win at index 2 — AI must block
	board := Board{"X", "X", "", "O", "O", "X", "", "", ""}
	move := bestMove(board)
	if move != 2 {
		t.Errorf("expected AI to block at 2, got %d", move)
	}
}

func TestMinimax_PrefersCenter(t *testing.T) {
	// On a near-empty board, center is optimal
	board := Board{"X", "", "", "", "", "", "", "", ""}
	move := bestMove(board)
	if move != 4 {
		t.Errorf("expected center (4), got %d", move)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/game/... -v -run TestMinimax`
Expected: FAIL — `bestMove` not defined.

- [ ] **Step 3: Implement minimax with alpha-beta pruning**

Append to `api/internal/game/engine.go`:

```go
import "math"

// minimax returns the score of the board for the maximizing player (O/AI).
// +10 = AI wins, -10 = player wins, 0 = draw. Depth penalty for faster wins.
func minimax(b Board, depth int, alpha, beta int, isMaximizing bool) int {
	winner := checkWinner(b)
	if winner == "O" {
		return 10 - depth
	}
	if winner == "X" {
		return depth - 10
	}
	if isBoardFull(b) {
		return 0
	}

	if isMaximizing {
		best := math.MinInt32
		for i := 0; i < 9; i++ {
			if b[i] != "" {
				continue
			}
			b[i] = "O"
			score := minimax(b, depth+1, alpha, beta, false)
			b[i] = ""
			if score > best {
				best = score
			}
			if best > alpha {
				alpha = best
			}
			if beta <= alpha {
				break
			}
		}
		return best
	}

	best := math.MaxInt32
	for i := 0; i < 9; i++ {
		if b[i] != "" {
			continue
		}
		b[i] = "X"
		score := minimax(b, depth+1, alpha, beta, true)
		b[i] = ""
		if score < best {
			best = score
		}
		if best < beta {
			beta = best
		}
		if beta <= alpha {
			break
		}
	}
	return best
}

// bestMove returns the optimal move index for the AI (O) using minimax.
func bestMove(b Board) int {
	best := math.MinInt32
	move := -1
	for i := 0; i < 9; i++ {
		if b[i] != "" {
			continue
		}
		b[i] = "O"
		score := minimax(b, 0, math.MinInt32, math.MaxInt32, false)
		b[i] = ""
		if score > best {
			best = score
			move = i
		}
	}
	return move
}
```

Note: update the import block at the top of `engine.go` to include `"math"`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/game/... -v`
Expected: All 10 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add api/internal/game/engine.go api/internal/game/engine_test.go
git commit -m "feat(game): add minimax with alpha-beta pruning"
```

---

### Task 3: Implement adaptive difficulty and move selection

**Files:**
- Modify: `api/internal/game/engine.go`
- Modify: `api/internal/game/engine_test.go`

- [ ] **Step 1: Write failing tests for adaptive difficulty**

Append to `api/internal/game/engine_test.go`:

```go
func TestRankedMoves_ReturnsAllEmpty(t *testing.T) {
	board := Board{"X", "", "", "", "", "", "", "", ""}
	moves := rankedMoves(board)
	if len(moves) != 8 {
		t.Errorf("expected 8 moves, got %d", len(moves))
	}
	// Best move should be first
	if moves[0].index != bestMove(board) {
		t.Errorf("expected best move %d first, got %d", bestMove(board), moves[0].index)
	}
}

func TestSelectMove_OptimalWhenWinning(t *testing.T) {
	// Visitor has been winning (negative engagement) — AI should play optimally
	board := Board{"X", "", "", "", "", "", "", "", ""}
	session := Session{Wins: 3, Losses: 0, Draws: 0}
	move := SelectMove(board, session)
	optimal := bestMove(board)
	if move != optimal {
		t.Errorf("expected optimal move %d, got %d", optimal, move)
	}
}

func TestSelectMove_SuboptimalWhenLosing(t *testing.T) {
	// Visitor has been losing a lot (high positive engagement) — AI should slip
	board := Board{"X", "", "", "", "", "", "", "", ""}
	session := Session{Wins: 0, Losses: 5, Draws: 0}
	optimal := bestMove(board)
	// Run multiple times — at least one should differ from optimal
	slipped := false
	for i := 0; i < 50; i++ {
		move := SelectMove(board, session)
		if move != optimal {
			slipped = true
			break
		}
	}
	if !slipped {
		t.Error("expected AI to slip at least once in 50 tries when visitor is losing heavily")
	}
}

func TestSelectMove_FirstGameSlips(t *testing.T) {
	// First game — AI should play suboptimally
	board := Board{"X", "", "", "", "", "", "", "", ""}
	session := Session{Wins: 0, Losses: 0, Draws: 0}
	optimal := bestMove(board)
	slipped := false
	for i := 0; i < 50; i++ {
		move := SelectMove(board, session)
		if move != optimal {
			slipped = true
			break
		}
	}
	if !slipped {
		t.Error("expected AI to slip at least once in 50 tries during first game")
	}
}

func TestSelectMove_NeverCreatesImmediateLoss(t *testing.T) {
	// Board where player is one move from winning at index 6
	// AI must not leave index 6 open when slipping
	board := Board{"X", "O", "", "X", "O", "", "", "", "X"}
	session := Session{Wins: 0, Losses: 10, Draws: 0} // Heavy losing — force slip
	for i := 0; i < 100; i++ {
		move := SelectMove(board, session)
		// After AI places O at 'move', check if player can win next turn
		test := board
		test[move] = "O"
		if checkWinner(test) != "" {
			continue // AI won, that's fine
		}
		// Check all remaining empty cells for player win
		for j := 0; j < 9; j++ {
			if test[j] != "" {
				continue
			}
			check := test
			check[j] = "X"
			if checkWinner(check) == "X" {
				// Player can win in one move — is it the ONLY option?
				// Count how many empty cells are left
				empty := 0
				for _, c := range test {
					if c == "" {
						empty++
					}
				}
				if empty > 1 {
					// AI had other options but chose one that lets player win immediately
					// This is acceptable only if ALL moves lead to player winning next turn
					allLose := true
					for k := 0; k < 9; k++ {
						if board[k] != "" {
							continue
						}
						alt := board
						alt[k] = "O"
						if checkWinner(alt) == "O" {
							allLose = false
							break
						}
						canWin := false
						for m := 0; m < 9; m++ {
							if alt[m] != "" {
								continue
							}
							altCheck := alt
							altCheck[m] = "X"
							if checkWinner(altCheck) == "X" {
								canWin = true
								break
							}
						}
						if !canWin {
							allLose = false
							break
						}
					}
					if !allLose {
						t.Errorf("AI move %d allows player to win immediately at %d (run %d)", move, j, i)
						return
					}
				}
			}
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/game/... -v -run "TestRankedMoves|TestSelectMove"`
Expected: FAIL — `rankedMoves`, `SelectMove` not defined.

- [ ] **Step 3: Implement adaptive difficulty**

Append to `api/internal/game/engine.go`:

```go
import "math/rand"

type scoredMove struct {
	index int
	score int
}

// rankedMoves returns all available moves ranked by minimax score (best first).
func rankedMoves(b Board) []scoredMove {
	var moves []scoredMove
	for i := 0; i < 9; i++ {
		if b[i] != "" {
			continue
		}
		b[i] = "O"
		score := minimax(b, 0, math.MinInt32, math.MaxInt32, false)
		b[i] = ""
		moves = append(moves, scoredMove{index: i, score: score})
	}
	// Sort descending by score (best move first)
	for i := 0; i < len(moves); i++ {
		for j := i + 1; j < len(moves); j++ {
			if moves[j].score > moves[i].score {
				moves[i], moves[j] = moves[j], moves[i]
			}
		}
	}
	return moves
}

// wouldAllowImmediateWin checks if placing O at index lets X win on the next turn.
func wouldAllowImmediateWin(b Board, index int) bool {
	test := b
	test[index] = "O"
	if checkWinner(test) != "" {
		return false // AI wins here — not a problem
	}
	for i := 0; i < 9; i++ {
		if test[i] != "" {
			continue
		}
		check := test
		check[i] = "X"
		if checkWinner(check) == "X" {
			return true
		}
	}
	return false
}

// SelectMove picks the AI's move based on the board and session context.
// It uses minimax for optimal play and the engagement model to decide when to slip.
func SelectMove(b Board, session Session) int {
	moves := rankedMoves(b)
	if len(moves) == 0 {
		return -1
	}
	if len(moves) == 1 {
		return moves[0].index
	}

	engagementScore := (session.Losses*2 + session.Draws) - (session.Wins * 2)
	isFirstGame := session.Wins+session.Losses+session.Draws == 0
	shouldSlip := false

	if isFirstGame {
		shouldSlip = true
	} else if engagementScore > 2 {
		// Visitor losing — slip to give them a chance
		shouldSlip = true
	} else if engagementScore >= -2 {
		// Balanced — slip 40% of the time
		shouldSlip = rand.Intn(100) < 40
	}
	// engagementScore < -2: visitor winning — play optimally (shouldSlip stays false)

	if !shouldSlip {
		return moves[0].index
	}

	// Pick 2nd or 3rd best move that doesn't allow an immediate win
	for _, m := range moves[1:] {
		if len(moves) > 2 && m.index == moves[len(moves)-1].index {
			continue // Skip the worst move
		}
		if !wouldAllowImmediateWin(b, m.index) {
			return m.index
		}
	}

	// All suboptimal moves allow immediate wins — play optimally
	return moves[0].index
}
```

Note: update the import block at the top of `engine.go` to include `"math/rand"`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/game/... -v`
Expected: All 15 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add api/internal/game/engine.go api/internal/game/engine_test.go
git commit -m "feat(game): add adaptive difficulty and move selection"
```

---

### Task 4: Implement the HTTP handler with validation

**Files:**
- Create: `api/internal/game/handler.go`
- Create: `api/internal/game/handler_test.go`

- [ ] **Step 1: Write failing tests for the handler**

Create `api/internal/game/handler_test.go`:

```go
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

	// O has same count as X — it's not AI's turn
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

	// X already won
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

func TestHandler_PlayerWinsBeforeAIMove(t *testing.T) {
	handler := game.NewHandler()

	// Player just won — send the winning board
	// X wins on top row, but it's "AI's turn" by count (3X, 2O)
	reqBody := game.GameRequest{
		Board:   [9]string{"X", "X", "X", "O", "O", "", "", "", ""},
		Session: game.Session{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Game is already over — should reject
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/game/... -v -run "TestHandler"`
Expected: FAIL — `GameRequest`, `GameResponse`, `NewHandler` not defined.

- [ ] **Step 3: Implement the handler**

Create `api/internal/game/handler.go`:

```go
package game

import (
	"encoding/json"
	"net/http"
)

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
		var req GameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

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
		if w := checkWinner(board); w != "" {
			winner = w
			if w == "O" {
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
	for i, cell := range b {
		if cell != "" && cell != "X" && cell != "O" {
			_ = i
			return "invalid board value"
		}
	}
	return ""
}
```

Note: The handler uses `w` as both the http.ResponseWriter variable name and the winner variable name. Rename the winner check variable to avoid shadowing:

Replace the status determination block with:

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/game/... -v`
Expected: All tests PASS (15 engine + 7 handler = 22 total).

- [ ] **Step 5: Commit**

```bash
git add api/internal/game/handler.go api/internal/game/handler_test.go
git commit -m "feat(game): add HTTP handler with validation"
```

---

### Task 5: Register the game route in main.go

**Files:**
- Modify: `api/cmd/server/main.go`

- [ ] **Step 1: Add the game route**

In `api/cmd/server/main.go`, add the import:

```go
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
```

Then add the game route after the chat route registration (after line 31):

```go
	mux.Handle("POST /api/game", game.NewHandler())
```

- [ ] **Step 2: Run all API tests**

Run: `cd api && go vet ./... && go test ./... -v`
Expected: All tests pass, no vet issues.

- [ ] **Step 3: Test manually with curl**

Start the server: `cd api && go run ./cmd/server &`

```bash
curl -X POST http://localhost:8080/api/game \
  -H "Content-Type: application/json" \
  -d '{"board":["X","","","","","","","",""],"session":{"wins":0,"losses":0,"draws":0}}'
```

Expected: 200 response with `move`, `board`, `status`, `winner` fields. Kill the server after.

- [ ] **Step 4: Commit**

```bash
git add api/cmd/server/main.go
git commit -m "feat(game): register POST /api/game route"
```

---

### Task 6: Create the playground Astro page

**Files:**
- Create: `frontend/src/pages/playground.astro`

- [ ] **Step 1: Create the page**

Create `frontend/src/pages/playground.astro`:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
import Footer from '../components/Footer.astro';
import GameBoard from '../islands/GameBoard.tsx';
---

<BaseLayout title="Playground | Rafael Abreu" description="Play Tic-Tac-Toe against Rafael's AI — built with Go and minimax.">
  <Nav />
  <main class="min-h-[80vh] py-20 px-6">
    <div class="container mx-auto max-w-6xl">
      <div class="mb-2">
        <span class="font-mono text-accent-indigo text-xs tracking-[3px]">05 // PLAYGROUND</span>
      </div>
      <h1 class="text-4xl md:text-5xl font-black mb-4" style="letter-spacing: -0.06em;">
        Beat My AI
      </h1>
      <p class="text-text-secondary text-sm md:text-base max-w-lg leading-relaxed mb-10">
        A game of Tic-Tac-Toe powered by a Go backend running minimax with alpha-beta pruning. You're X. Make your move.
      </p>
      <GameBoard client:load />
    </div>
  </main>
  <Footer />
</BaseLayout>
```

- [ ] **Step 2: Build and verify**

Run: `cd frontend && npm run build`
Expected: Build succeeds, 6 pages built (was 5, now includes `/playground`).

- [ ] **Step 3: Commit**

```bash
git add frontend/src/pages/playground.astro
git commit -m "feat(playground): add Astro page for AI playground"
```

---

### Task 7: Create the GameBoard React island

**Files:**
- Create: `frontend/src/islands/GameBoard.tsx`

- [ ] **Step 1: Create the GameBoard component**

Create `frontend/src/islands/GameBoard.tsx`:

```tsx
import { useState } from 'react';

interface Session {
  wins: number;
  losses: number;
  draws: number;
}

type Status = 'waiting' | 'playing' | 'player_win' | 'ai_win' | 'draw';

const API_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080';

export default function GameBoard() {
  const [board, setBoard] = useState<string[]>(Array(9).fill(''));
  const [status, setStatus] = useState<Status>('waiting');
  const [isThinking, setIsThinking] = useState(false);
  const [session, setSession] = useState<Session>({ wins: 0, losses: 0, draws: 0 });
  const [error, setError] = useState<string | null>(null);
  const [winLine, setWinLine] = useState<number[] | null>(null);

  const WIN_LINES = [
    [0, 1, 2], [3, 4, 5], [6, 7, 8],
    [0, 3, 6], [1, 4, 7], [2, 5, 8],
    [0, 4, 8], [2, 4, 6],
  ];

  function findWinLine(b: string[]): number[] | null {
    for (const line of WIN_LINES) {
      const [a, bc, c] = line;
      if (b[a] && b[a] === b[bc] && b[bc] === b[c]) {
        return line;
      }
    }
    return null;
  }

  async function handleCellClick(index: number) {
    if (board[index] !== '' || isThinking || status === 'player_win' || status === 'ai_win' || status === 'draw') {
      return;
    }

    setError(null);
    const newBoard = [...board];
    newBoard[index] = 'X';
    setBoard(newBoard);
    setStatus('playing');

    // Check if player just won (shouldn't happen in normal play since AI goes second, but handle it)
    const playerWin = findWinLine(newBoard);
    if (playerWin) {
      setStatus('player_win');
      setWinLine(playerWin);
      setSession((s) => ({ ...s, wins: s.wins + 1 }));
      return;
    }

    // Check if board is full after player move (draw)
    if (newBoard.every((c) => c !== '')) {
      setStatus('draw');
      setSession((s) => ({ ...s, draws: s.draws + 1 }));
      return;
    }

    // Ask AI for its move
    setIsThinking(true);
    try {
      const resp = await fetch(`${API_URL}/api/game`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ board: newBoard, session }),
      });

      if (!resp.ok) {
        if (resp.status === 429) {
          throw new Error('Too many moves! Slow down.');
        }
        throw new Error("Can't reach the game server. Try again.");
      }

      const data = await resp.json();
      setBoard(data.board);

      if (data.status === 'ai_win') {
        setStatus('ai_win');
        setWinLine(findWinLine(data.board));
        setSession((s) => ({ ...s, losses: s.losses + 1 }));
      } else if (data.status === 'draw') {
        setStatus('draw');
        setSession((s) => ({ ...s, draws: s.draws + 1 }));
      } else if (data.status === 'player_win') {
        setStatus('player_win');
        setWinLine(findWinLine(data.board));
        setSession((s) => ({ ...s, wins: s.wins + 1 }));
      } else {
        setStatus('playing');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong.');
      // Revert the player's move on error
      setBoard(board);
      setStatus(board.every((c) => c === '') ? 'waiting' : 'playing');
    } finally {
      setIsThinking(false);
    }
  }

  function resetGame() {
    setBoard(Array(9).fill(''));
    setStatus('waiting');
    setIsThinking(false);
    setError(null);
    setWinLine(null);
  }

  function statusMessage(): string {
    switch (status) {
      case 'waiting':
        return "You're X. Make your move.";
      case 'playing':
        return isThinking ? 'AI is thinking...' : 'Your turn.';
      case 'player_win':
        return 'You win!';
      case 'ai_win':
        return 'AI wins. Try again?';
      case 'draw':
        return "It's a draw.";
    }
  }

  const isGameOver = status === 'player_win' || status === 'ai_win' || status === 'draw';

  return (
    <div className="flex flex-col items-center">
      <div className="font-mono text-xs text-text-muted tracking-wider mb-6">
        {statusMessage()}
      </div>

      <div className="grid grid-cols-3 gap-2 mb-6">
        {board.map((cell, i) => {
          const isWinCell = winLine?.includes(i);
          return (
            <button
              key={i}
              onClick={() => handleCellClick(i)}
              disabled={isThinking || isGameOver}
              className={`w-20 h-20 sm:w-24 sm:h-24 rounded-card border text-2xl sm:text-3xl font-black transition-all duration-200 ${
                isWinCell
                  ? 'bg-accent-indigo/20 border-accent-indigo/60'
                  : 'bg-surface border-surface-border hover:border-accent-indigo/30'
              } ${
                cell === '' && !isThinking && !isGameOver
                  ? 'cursor-pointer'
                  : 'cursor-default'
              }`}
              aria-label={cell || `Empty cell ${i}`}
            >
              {cell === 'X' && <span className="text-text-primary">X</span>}
              {cell === 'O' && <span className="text-accent-indigo">O</span>}
              {cell === '' && !isThinking && !isGameOver && (
                <span className="text-text-primary opacity-0 hover:opacity-20 transition-opacity">X</span>
              )}
            </button>
          );
        })}
      </div>

      {isGameOver && (
        <button
          onClick={resetGame}
          className="px-8 py-3 bg-accent-indigo rounded-btn text-white text-sm font-semibold hover:bg-accent-indigo/90 transition-colors duration-200 mb-4"
        >
          Play Again
        </button>
      )}

      {error && (
        <div className="text-red-400 text-xs text-center mb-4">
          {error}{' '}
          <button onClick={resetGame} className="underline">
            Reset
          </button>
        </div>
      )}

      <div className="font-mono text-xs text-text-muted tracking-wider">
        W: {session.wins} / L: {session.losses} / D: {session.draws}
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Verify no TypeScript errors**

Run: `cd frontend && npx astro check`
Expected: 0 errors.

- [ ] **Step 3: Build the project**

Run: `cd frontend && npm run build`
Expected: Build succeeds, 6 pages built.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/islands/GameBoard.tsx
git commit -m "feat(playground): add GameBoard React island"
```

---

### Task 8: Add Playground to navigation

**Files:**
- Modify: `frontend/src/components/Nav.astro`

- [ ] **Step 1: Add Playground nav item**

In `frontend/src/components/Nav.astro`, update the `navItems` array:

```javascript
const navItems = [
  { label: 'About', href: '#about' },
  { label: 'Projects', href: '#projects' },
  { label: 'Skills', href: '#skills' },
  { label: 'Playground', href: '/playground' },
  { label: 'Contact', href: '#contact' },
];
```

- [ ] **Step 2: Build and verify**

Run: `cd frontend && npm run build`
Expected: Build succeeds, 6 pages built.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/Nav.astro
git commit -m "feat(playground): add Playground to navigation"
```

---

### Task 9: Visual verification

**Files:**
- Possibly modify: `frontend/src/islands/GameBoard.tsx`

- [ ] **Step 1: Start both servers**

Terminal 1: `cd api && go run ./cmd/server`
Terminal 2: `cd frontend && npm run dev`

Open `http://localhost:4321/playground` in a browser.

- [ ] **Step 2: Verify game flow**

- Click a cell — X appears, AI responds with O after a brief pause
- Play a full game to completion — result message appears, "Play Again" button shows
- Click "Play Again" — board resets, session stats preserved
- Session stats update correctly (W/L/D counter)
- Hover on empty cells shows faint X preview

- [ ] **Step 3: Verify navigation**

- Playground link appears in desktop nav and mobile hamburger menu
- Clicking "Playground" from the homepage navigates to `/playground`
- Nav on the playground page still shows all links

- [ ] **Step 4: Verify adaptive difficulty**

- Play 3-4 games in a row losing — AI should start making weaker moves
- After winning a couple games — AI should tighten up

- [ ] **Step 5: Verify error handling**

- Stop the API server, try to make a move — error message appears with reset option
- Restart API, reset game — should work again

- [ ] **Step 6: Commit any tuning changes**

```bash
git add -A
git commit -m "fix: tune playground visual parameters"
```

(Skip if no changes needed.)

---

### Task 10: Final build verification, push, and PR

**Files:**
- No new changes expected

- [ ] **Step 1: Run all checks**

```bash
cd frontend && npx astro check
cd frontend && npm run build
cd api && go vet ./... && go test ./... -v
```

Expected: All pass.

- [ ] **Step 2: Verify git status**

Run: `git status`
Expected: Clean working tree on `feature/ai-playground`.

- [ ] **Step 3: Push and create PR**

```bash
git push -u origin feature/ai-playground
gh pr create --base main --head feature/ai-playground \
  --title "feat: AI Playground — Tic-Tac-Toe with adaptive difficulty" \
  --body "## Summary
- New \`POST /api/game\` endpoint with minimax + alpha-beta pruning
- Adaptive engagement model — AI adjusts difficulty based on session history
- Casino-style difficulty: visitor wins ~40%, draws ~30%, loses ~30%
- First game is deliberately easier to hook visitors
- New \`/playground\` page with React island game board
- Playground added to site navigation
- Full server-side validation, 22+ Go tests

## Test plan
- [ ] Play multiple games, verify adaptive difficulty feels natural
- [ ] Verify session stats (W/L/D) track correctly across games
- [ ] Mobile: board scales, cells are tappable
- [ ] API down: error message displays, reset works
- [ ] Build passes (astro check + npm run build + go vet + go test)"
```
