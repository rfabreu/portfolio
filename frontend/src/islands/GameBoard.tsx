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
    const prevBoard = [...board];
    const newBoard = [...board];
    newBoard[index] = 'X';
    setBoard(newBoard);
    setStatus('playing');

    // Check if player just won
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
      if (!Array.isArray(data?.board) || data.board.length !== 9) {
        throw new Error('Invalid response from game server.');
      }
      setBoard(data.board);

      if (data.status === 'ai_win') {
        setStatus('ai_win');
        setWinLine(findWinLine(data.board));
        setSession((s) => ({ ...s, losses: s.losses + 1 }));
      } else if (data.status === 'draw') {
        setStatus('draw');
        setSession((s) => ({ ...s, draws: s.draws + 1 }));
      } else {
        setStatus('playing');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong.');
      setBoard(prevBoard);
      setStatus(prevBoard.every((c) => c === '') ? 'waiting' : 'playing');
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
