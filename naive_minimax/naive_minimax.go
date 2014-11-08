// This package defines a naive, sequential implementation of the minimax algorithm.
package naive_minimax

// A State represents a game state in a user-defined game.
type State interface {
	// The heuristic value of the game state.
	// The result should be positive if the maximized player has won,
	// negative if the minimized player has won, and 0 if the outcome is a
	// draw.
	Value() Value

	// MaximizingPlayer is true if the next player to play for the game state
	// is the maximized player, false if it is the minimized player.
	MaximizingPlayer() bool

	// An iterator over game states one move away from the current game state.
	// The result should be nil if the game is in an end state.
	Children() StateIterator

	// Converts the State to a serialized representation.
	EncodeState() ([]byte, err)

	// Populates the State based on the serialized representation provided by EncodeState.
	DecodeState([]byte) err
}

// A StateIterator provides access to a sequence of States.
type StateIterator interface {
	// Povides the next State in the sequence. Returns nil if there are no more states.
	Next() State
}

// Value is the type of the heuristic values of game states.
type Value int

// Minimax implements the minimax algorithm. It determines the best possible outcome of a
// game sequence starting at the provided state by evaluating the game tree rooted at the
// state up to the provided maxDepth. Leaves of the evaluated game tree are valued based
// on their State.Value. If the State.MaximizingPlayer is true, then the best outcome will
// be the one with the maximum value. Otherwise, it will be the one with the minimum value.
func Minimax(state State, maxDepth int) Value {
	child := state.Children()
	maximizingPlayer := state.MaximizingPlayer()

	if maxDepth == 0 || child == nil {
		return state.Value()
	}

	var better func(a, b Value) Value
	if maximizingPlayer {
		better = max
	} else {
		better = min
	}

	best := Minimax(child.State(), maxDepth-1, !maximizingPlayer)

	for {
		child = child.Next()
		if child == nil {
			break
		}

		val := Minimax(child.State(), maxDepth-1, !maximizingPlayer)
		best = better(best, val)
	}

	return best
}

func min(a, b Value) Value {
	if a < b {
		return a
	}
	return b
}

func max(a, b Value) Value {
	if a > b {
		return a
	}
	return b
}
