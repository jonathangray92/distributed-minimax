// Package game defines the interfaces implemented by games used with the minimax package.
package game

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

	// An iterator over game states one move away from the current game state. The result
	// should not be nil, even if there are no possible moves.
	Moves() StateIterator

	// Converts the State to a serialized representation.
	EncodeState() ([]byte, error)

	// Populates the State based on the serialized representation provided by EncodeState.
	DecodeState([]byte) error

	// Returns an identifier unique to this game state. Two game states that originate
	// from different move sequences but are otherwise identical should have the same Id.
	// The Id can be any type of value provided it is comparable.
	Id() interface{}
}

// A StateIterator provides access to a sequence of States. Each call to a
// StateIterator returns the next State in the sequence, or nil if there are no
// more States in the sequence.
type StateIterator func() (next State)

// Value is the type of the heuristic values of game states.
type Value int
