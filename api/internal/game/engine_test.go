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

func TestMinimax_TakesWinningMove(t *testing.T) {
	// O can win at index 5
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
				empty := 0
				for _, c := range test {
					if c == "" {
						empty++
					}
				}
				if empty > 1 {
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
