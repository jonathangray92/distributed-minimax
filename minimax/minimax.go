package minimax

import (
	"github.com/jonathangray92/distributed-minimax/game"
	"time"
	"log"
)

// Implements the minimax algorithm with alpha-beta pruning
// Should always return the same result as naiveminimax.Minimax, hopefully faster
// Follows the pseudo-code at http://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
//
// N.B. this function allows the caller to specify a value map which overrides
// the heuristic values of certain nodes. If a value map is specified and this
// function is called on a state that is present as a key in that map, the
// function will return immediately with the value provided in the map. This
// functionality is used by the mastser's result aggregator.
func AlphaBetaImpl(state game.State, maxDepth int, alpha game.Value, beta game.Value, valueMap map[interface{}]game.Value) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {

	// if we are using a value map, check for its value there and return early
	// if found
	if valueMap != nil {
		value, found := valueMap[state.Id()]
		if found {
			return value, nil, 1
		}
	}

	// get an iterator over the children of this state
	iterMoves := state.MoveIterator()
	move := iterMoves()

	// if we are at a leaf node, return the heuristic value of the state
	if maxDepth == 0 || move == nil {
		stateVal := state.Value()
		return stateVal, nil, 1
	}

	numStatesAnalyzed = 0;

	if state.MaximizingPlayer() {
		bestVal = game.MinValue
		for ; move != nil && alpha < beta ; move = iterMoves() {
			val, _, childStatesAnalyzed := AlphaBetaImpl(move, maxDepth - 1, alpha, beta, valueMap)
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
			val, _, childStatesAnalyzed := AlphaBetaImpl(move, maxDepth - 1, alpha, beta, valueMap)
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
// to their proper values and uses no value map
func Minimax(state game.State, maxDepth int) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {
	return AlphaBeta(state, maxDepth, game.MinValue, game.MaxValue)
}

// calls the alpha-beta minimax algorithm with no value map (described below)
func AlphaBeta(state game.State, maxDepth int, alpha game.Value, beta game.Value) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {
	return AlphaBetaImpl(state, maxDepth, alpha, beta, nil)
}

// calls the alpha-beta minimax algorithm with proper initial alpha/beta, large max depth, and value map
func AlphaBetaWithValueMap(state game.State, valueMap map[interface{}]game.Value) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {
	return AlphaBetaImpl(state, 100, game.MinValue, game.MaxValue, valueMap)
}

// Implements iterative deepening alpha-beta minimax in a separate goroutine,
// updating the return values every time the AlphaBeta function returns.
// When timeLimit has elapsed, returns the value of the most recent (deepest)
// search, and quit the worker goroutine
func TimeLimitedAlphaBeta(state game.State, timeLimit time.Duration) (bestVal game.Value, bestMove game.State, numStatesAnalyzed int) {

	// this channel will fire after the time limit has been reached
	timeLimitReached := make(chan bool, 1)

	// this function will exit when
	go func() {
		for depth := 1; ; depth += 1 {
			select {
				case <-timeLimitReached:
					log.Println("returning from goroutine")
					return
				default:
					log.Println("analyzing maxDepth %v", depth)
					bestVal, bestMove, numStatesAnalyzed = Minimax(state, depth)
			}
		}
	}()

	// when the time limit is up, signal that the worker goroutine should return, and exit immediately
	time.Sleep(timeLimit)
	timeLimitReached <- true
	return
}
