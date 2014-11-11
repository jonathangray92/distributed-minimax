package tictactoe

import "github.com/jonathangray92/distributed-minimax/game"

type Player int

// Switches the Player between X and O. If the Player is None, this is a no-op.
func (p *Player) Toggle() { *p = -*p }

const (
	None Player = 0  // No player
	X    Player = 1  // Maximized player
	O    Player = -1 // Minimized player
)

// State implements game.State for tictactoe.
type State struct {
	// The next player to make a move from this State.
	Player Player

	// A 3x3 board storing which player (X or O) has claimed each square, or
	// None if the square is unclaimed.
	Board [3][3]Player
}

func (s State) Value() game.Value { return game.Value(s.winner()) }

// Returns the winner (X or O). If there is no winner in this game state,
// returns None.
func (s State) winner() Player {
	for _, v := range [...][3]Player{
		// Rows
		{s.Board[0][0], s.Board[0][1], s.Board[0][2]},
		{s.Board[1][0], s.Board[1][1], s.Board[1][2]},
		{s.Board[2][0], s.Board[2][1], s.Board[2][2]},

		// Columns
		{s.Board[0][0], s.Board[1][0], s.Board[2][0]},
		{s.Board[0][1], s.Board[1][1], s.Board[2][1]},
		{s.Board[0][2], s.Board[1][2], s.Board[2][2]},

		// Diagonals
		{s.Board[0][0], s.Board[1][1], s.Board[2][2]},
		{s.Board[0][2], s.Board[1][1], s.Board[2][0]},
	} {
		if v[0] != None && v[0] == v[1] && v[1] == v[2] {
			return v[0]
		}
	}

	return None
}

func (s State) MaximizingPlayer() bool { return s.Player == X }

func (s State) Moves() game.StateIterator {
	// If there is a winner, then there are no possible moves.
	if s.winner() != None {
		return func() game.State { return nil }
	}

	iNext, jNext := 0, 0

	iter := func() game.State {
		state := s
		for iNext < 3 {
			for jNext < 3 {
				i, j := iNext, jNext
				jNext++

				if state.Board[i][j] == None {
					// position is empty; perform the move and return
					state.Board[i][j] = state.Player
					state.Player.Toggle()
					return state
				}
			}
			iNext++
			jNext = 0
		}
		return nil
	}

	return iter
}

func (s State) EncodeState() ([]byte, error) { panic("Unimplemented.") }
func (s State) DecodeState(p []byte) error   { panic("Unimplemented.") }

func (s State) Id() interface{} { return s }
