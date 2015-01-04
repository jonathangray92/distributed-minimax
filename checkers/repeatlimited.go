package checkers

import "github.com/jonathangray92/distributed-minimax/game"

// Represents a checkers game with a draw end-state for when a given state has
// been repeated 3 times. It does this by maintaining both a fixed history up
// until the last move actually taken, as well as an temporary history of the
// move chain leading to this state from the most recently taken move.
type RepeatLimitedState struct {
	BasicState
	pastStates map[BasicState]int
	tempPast   *stateStack
}

// Constructs a new repeat-limited checkers game.
func NewRepeatLimitedGame() *RepeatLimitedState {
	return &RepeatLimitedState{
		BasicState: NewBasicGame(),
		pastStates: make(map[BasicState]int),
		tempPast:   nil,
	}
}

func (state RepeatLimitedState) MoveIterator() game.StateIterator {
	// Check for draw condition.
	if state.pastStates[state.BasicState] >= 3 {
		return func() game.State { return nil }
	}

	// Create an adaptor around the `BasicState.MoveIterator` to return
	// `*RepeatLimitedState`s.
	iterMoves := state.BasicState.MoveIterator()
	return func() game.State {
		move, ok := iterMoves().(BasicState)
		if !ok {
			return nil
		}
		return &RepeatLimitedState{
			BasicState: move,
			pastStates: state.pastStates,
			tempPast:   &stateStack{state.BasicState, state.tempPast},
		}
	}
}

func (state RepeatLimitedState) Value() game.Value {
	// Detect move loops in the fixed and temporary histories. If this move is
	// on such a loop, then it is of no value since it won't progress the game.

	// Check fixed history for this move.
	if state.pastStates[state.BasicState] > 0 {
		return 0
	}
	// Check temporary past for this move.
	for curr := state.tempPast; curr != nil; curr = curr.next {
		if curr.val == state.BasicState {
			return 0
		}
	}

	return state.BasicState.Value()
}

// Adds this move to the fixed history as a taken move, and clears the
// temporary history.
func (state *RepeatLimitedState) UseState() {
	state.pastStates[state.BasicState]++
	state.tempPast = nil
}

// Represents a stack of BasicStates.
type stateStack struct {
	val  BasicState
	next *stateStack
}
