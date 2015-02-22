package minimax

import "github.com/jonathangray92/distributed-minimax/game"

// Implements the minimax algorithm with alpha-beta pruning
// Should always return the same result as naiveminimax.Minimax, hopefully faster
// Follows the pseudo-code at http://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
func AlphaBeta(state game.State, maxDepth int, alpha game.Value, beta game.Value) (bestVal game.Value, bestMove game.State) {

	var (
		iterMoves = state.MoveIterator()
		move      = iterMoves()
	)

	if maxDepth == 0 || move == nil {
		// leaf node: return heuristic value of state
		stateVal := state.Value()
		return stateVal, nil
	}

	if state.MaximizingPlayer() {
		bestVal = game.MinValue
		for ; move != nil && alpha < beta ; move = iterMoves() {
			val, _ := AlphaBeta(move, maxDepth - 1, alpha, beta)
			if val > bestVal {
				bestVal, bestMove = val, move
			}
			if val > alpha {
				alpha = val
			}
		}
	} else {
		bestVal = game.MaxValue
		for ; move != nil && alpha < beta ; move = iterMoves() {
			val, _ := AlphaBeta(move, maxDepth - 1, alpha, beta)
			if val < bestVal {
				bestVal, bestMove = val, move
			}
			if val < beta {
				beta = val
			}
		}
	}

	return bestVal, bestMove
}


// convenience wrapper over the AlphaBeta function which initializes alpha and beta
// to their proper values
func Minimax(state game.State, maxDepth int) (bestVal game.Value, bestMove game.State) {
	bestVal, bestMove = AlphaBeta(state, maxDepth, game.MinValue, game.MaxValue)
	return
}
