package connect4

import (
	"testing"
	"github.com/jonathangray92/distributed-minimax/game"
)

// Test the makeMove function on an initial state
func TestMakeMove(t *testing.T) {
	// start with empty state
	start := NewInitialState()

	// Move in column 3
	next := start.makeMove(3)

	// Verify that the proper move was made
	if next.Y != 0 || next.IsXMove || next.X != (1 << (5*7+3)) {
		t.Error("Invalid next state %+v", next)
	}
}

// Test that a win is recognized
func TestWinCondition(t *testing.T) {
	// start with empty state
	state := NewInitialState()

	// X plays in column 3 and Y in column 4 until X wins
	// Just before the end, check that nobody has won yet
	state = state.makeMove(3)
	state = state.makeMove(4)
	state = state.makeMove(3)
	state = state.makeMove(4)
	state = state.makeMove(3)
	state = state.makeMove(4)

	// check that nobody has won
	if state.MoveIterator()() == nil {
		t.Error("state %+v should have moves, but doesn't", state)
	}
	if state.Value() == game.MaxValue || state.Value() == game.MinValue {
		t.Error("state %+v should not have win value %v", state, state.Value())
	}

	// now we make X win...
	state = state.makeMove(3)

	// ...and check for the proper results
	if state.MoveIterator()() != nil {
		t.Error("state %+v should have no moves, but does", state)
	}
	if state.Value() != game.MaxValue {
		t.Error("win for X should not have value %v", state.Value())
	}
}

func TestStateValue(t *testing.T) {
	// start with empty state and make a single move for X
	state := NewInitialState()
	state = state.makeMove(3)

	// X is closer to a win than Y (who has not played), so the value should be positive
	xValue := state.Value()
	if xValue <= 0 {
		t.Error("should have positive value, not %v", xValue)
	}

	// Y makes a move in the corner, which is less valuable
	state = state.makeMove(0)
	yValue := state.Value()
	if yValue <= 0 {
		t.Error("should have positive value, not %v", yValue)
	}
	if yValue >= xValue {
		t.Error("new value %v should be less than old value %v", yValue, xValue)
	}
}

func TestBlah(t *testing.T) {
	s := State{X:68719476736, Y:0, IsXMove:false}
	s.PrintState()
}