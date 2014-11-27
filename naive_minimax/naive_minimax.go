// This package defines a naive, sequential implementation of the minimax algorithm.
package naive_minimax

import "github.com/jonathangray92/distributed-minimax/game"

// Minimax implements the minimax algorithm. It determines the best possible outcome of a
// game sequence starting at the provided state by evaluating the game tree rooted at
// that state up to the provided maxDepth. Leaves of the evaluated game tree are valued
// based on their State.Value. If state.MaximizingPlayer() is true, then the best outcome
// will be the one with the maximum value. Otherwise, it will be the one with the minimum
// value. The function returns the value of the best outcome and the game state from
// after the first move towards that outcome. If the provided state is a leaf node or
// maxDepth is zero, bestMove will be nil and bestVal will be state.Value().
func Minimax(state game.State, maxDepth int) (bestVal game.Value, bestMove game.State) {

	var (
		iterMoves = state.MoveIterator()
		move      = iterMoves()
	)

	if maxDepth == 0 || move == nil {
		// we are at a leaf of the (potentially truncated) tree
		return state.Value(), nil
	}

	// function used to determine which outcomes are preferred for the current
	// player
	var better func(a, b game.Value) bool
	if state.MaximizingPlayer() {
		better = greater
	} else {
		better = less
	}

	// priming the loop
	bestMove = move
	bestVal, _ = Minimax(bestMove, maxDepth-1)

	// compare all possible moves to determine which one is best
	for {
		move = iterMoves()
		if move == nil {
			break
		}

		val, _ := Minimax(move, maxDepth-1)

		if better(val, bestVal) {
			bestVal, bestMove = val, move
		}
	}

	return bestVal, bestMove
}

func less(a, b game.Value) bool    { return a < b }
func greater(a, b game.Value) bool { return a > b }
