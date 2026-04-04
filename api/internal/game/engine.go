package game

import (
	"math"
	"math/rand"
)

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
