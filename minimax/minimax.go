package minimax

import "github.com/jonathangray92/distributed-minimax/game"

// Implements the minimax algorithm with alpha-beta pruning
// Should always return the same result as naiveminimax.Minimax, hopefully faster
// Follows the pseudo-code at http://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
func AlphaBeta(state game.State, maxDepth int, alpha game.Value, beta game.Value) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {

	var (
		iterMoves = state.MoveIterator()
		move      = iterMoves()
	)

	if maxDepth == 0 || move == nil {
		// leaf node: return heuristic value of state
		stateVal := state.Value()
		return stateVal, nil, 1
	}

	numStatesAnalyzed = 0;

	if state.MaximizingPlayer() {
		bestVal = game.MinValue
		for ; move != nil && alpha < beta ; move = iterMoves() {
			val, _, childStatesAnalyzed := AlphaBeta(move, maxDepth - 1, alpha, beta)
			numStatesAnalyzed += childStatesAnalyzed
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
			val, _, childStatesAnalyzed := AlphaBeta(move, maxDepth - 1, alpha, beta)
			numStatesAnalyzed += childStatesAnalyzed
			if val < bestVal {
				bestVal, bestMove = val, move
			}
			if val < beta {
				beta = val
			}
		}
	}

	return bestVal, bestMove, numStatesAnalyzed
}


// convenience wrapper over the AlphaBeta function which initializes alpha and beta
// to their proper values
func Minimax(state game.State, maxDepth int) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {
	bestVal, bestMove, numStatesAnalyzed = AlphaBeta(state, maxDepth, game.MinValue, game.MaxValue)
	return
}
