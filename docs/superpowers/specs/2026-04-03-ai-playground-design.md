# P2: AI Playground — Design Spec

> A dedicated `/playground` page where visitors play Tic-Tac-Toe against the Go API. The AI uses minimax with adaptive difficulty — it subtly adjusts its play to keep visitors engaged, winning just often enough to keep playing.

## Architecture

### Stateless API

`POST /api/game` — client sends the board state and session context, server returns the AI's move. No server-side sessions, no storage. Matches the existing chat endpoint pattern (client sends context, server responds).

### Request

```json
{
  "board": ["X", "", "O", "", "X", "", "", "", ""],
  "session": { "wins": 1, "losses": 2, "draws": 0 }
}
```

- `board`: 9-element array (flat 3x3 grid). `""` = empty, `"X"` = player (human), `"O"` = AI.
- `session`: running tally of game outcomes in this browser session. Tracked in React state, sent with each request.

### Response

```json
{
  "move": 5,
  "board": ["X", "", "O", "", "X", "O", "", "", ""],
  "status": "in_progress",
  "winner": ""
}
```

- `move`: index (0-8) where the AI placed `"O"`.
- `board`: updated board state after the AI's move.
- `status`: `"in_progress"` | `"player_win"` | `"ai_win"` | `"draw"`.
- `winner`: `"X"` | `"O"` | `""`.

### Go Packages

```
api/internal/game/
  ├── handler.go       — HTTP handler, request/response structs, validation
  ├── engine.go        — Minimax algorithm, move selection, adaptive difficulty
  ├── handler_test.go  — httptest-based handler tests
  └── engine_test.go   — Unit tests for minimax and difficulty logic
```

### Route Registration

In `api/cmd/server/main.go`:

```go
mux.Handle("POST /api/game", game.NewHandler())
```

No provider interface needed — the game engine is pure Go logic with no external dependencies. The existing CORS + rate limit middleware chain applies automatically.

## Game Engine

### Minimax with Alpha-Beta Pruning

The engine computes the optimal move for any board state via minimax with alpha-beta pruning. This is the core algorithm — it always knows the perfect play.

### Adaptive Engagement Model

The engine decides each turn whether to play optimally or suboptimally, based on session history:

```
engagementScore = (losses * 2 + draws) - (wins * 2)
```

- **Positive score** (visitor has been losing) — AI makes a deliberate mistake. Picks the 2nd or 3rd best move from minimax-ranked moves, not random. Leaves an opening without being obvious.
- **Negative score** (visitor has been winning) — AI plays optimally. Forces draws or wins.
- **Near zero** (score between -2 and 2 inclusive) — AI plays optimally ~60% of the time, slips ~40%.

### First Game Hook

If `wins + losses + draws == 0` (first game), the AI plays suboptimally on its early moves to give the visitor a realistic chance of winning. Hooks them immediately.

### Believable Mistakes

When the AI decides to slip:
1. Run minimax on all available empty cells, score each move.
2. Rank moves from best to worst.
3. Pick the move ranked 2nd or 3rd (not last — that would be suspiciously bad).
4. Never slip in a way that creates an immediate win for the player on their next move. The AI will let the player earn a 2-move setup, not hand them a win-in-one.

### Server-Side Validation

- Board must be a 9-element array containing only `""`, `"X"`, or `"O"`.
- It must be the AI's turn: count of `"X"` must equal count of `"O"` plus 1 (player moves first).
- The game must not already be over (no existing winner, board not full).
- Session fields must be non-negative integers.
- Return 400 with a JSON error message for any invalid state.

## Frontend

### New Page: `frontend/src/pages/playground.astro`

- Uses `BaseLayout` for consistent Nav/Footer.
- Section header with number label and title (matching existing section style).
- Brief intro text explaining the game.
- Contains `<GameBoard client:load />` React island.

### New Island: `frontend/src/islands/GameBoard.tsx`

Owns all game UI and state. Renders the 3x3 grid, status messages, session stats, and controls.

### State

```typescript
board: string[]           // 9-element array, "" | "X" | "O"
status: "waiting" | "playing" | "player_win" | "ai_win" | "draw"
isThinking: boolean       // true while waiting for API response
session: { wins: number, losses: number, draws: number }
```

### Game Flow

1. Page loads — empty board, status = `"waiting"`, prompt: "You're X. Make your move."
2. Player clicks an empty cell — board updates optimistically with `"X"`, `isThinking = true`.
3. `POST /api/game` fires with board + session.
4. API responds — board updates with AI's `"O"`, `isThinking = false`.
5. If `status` is terminal (`player_win` / `ai_win` / `draw`): show result message + "Play Again" button.
6. "Play Again" resets the board to empty but preserves session counts.
7. Session stats displayed subtly below the board: "W: 2 / L: 1 / D: 0".

### Cell Interaction

- Empty cells show a faint `"X"` preview on hover (opacity ~0.2).
- Clicking a cell during the AI's turn or on an occupied cell does nothing.
- AI's move appears after the API responds (no optimistic placement for the AI).

### Error Handling

- API unreachable — show "Can't reach the game server. Try again." with a retry button.
- Rate limited (429) — "Too many moves! Slow down."
- Invalid response — reset board to empty, show error message.

### Navigation

Add "Playground" to the nav items in:
- `frontend/src/components/Nav.astro` — `navItems` array
- Pass through to `MobileNav.tsx` via existing props

Link: `/playground`

### Design Style

Matches the existing site aesthetic:
- Dark background (`bg-base`), card surfaces (`bg-surface`, `border-surface-border`)
- Board cells: `bg-surface` with `border-surface-border`, sized ~80x80px on desktop
- X marks: `text-text-primary` (white)
- O marks: `text-accent-indigo`
- Winning cells: highlighted with `accent-indigo` background glow
- Buttons: same `bg-accent-indigo` / `rounded-btn` style as rest of site
- Typography: mono font for labels/stats, system font for messages
- Responsive: board scales down on mobile, remains playable at small sizes

## Testing

### Go Engine Tests (`engine_test.go`)

- Minimax returns optimal move on a board with one winning move available
- Minimax blocks opponent's winning move
- Adaptive difficulty: with high engagement score (visitor losing), AI picks suboptimal move
- Adaptive difficulty: with low engagement score (visitor winning), AI picks optimal move
- First game: AI plays suboptimally early
- Slip mechanic never creates an immediate loss for the AI

### Go Handler Tests (`handler_test.go`)

- Valid request returns 200 with move, updated board, and status
- Invalid board (wrong length) returns 400
- Invalid board (wrong characters) returns 400
- Wrong turn (not AI's turn) returns 400
- Game already over returns 400
- Negative session values return 400

### Frontend

No automated frontend tests (matches existing project pattern). Visual verification via dev server.
