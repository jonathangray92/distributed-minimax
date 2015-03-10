package main

import (
	"github.com/jonathangray92/distributed-minimax/game"
	"github.com/jonathangray92/distributed-minimax/minimax"
	"log"
)

type Result struct {
	State game.State
	Value game.Value
}

type ResultAggregator interface {

	/**
	 * The arg is the number of AddResult() calls that will be executed before
	 * the callback is triggered. It should be equal to the number of jobs
	 * created and distributed to slaves.
	 */
	SetExpectedResultsCount(int)

	/**
	 * Call this once to set a callback that will execute after AddResult()
	 * has been called the given number of times.
	 *
	 * The results of the slaves will be aggregated to determine the best move
	 * for the user. The callback will be called with this best move.
	 */
	SetCallback(func(game.State))

	/**
	 * Call this with a new result delivered by a slave
	 *
	 * If a callback is set and this method has been called the correct number of
	 * times, call the callback with the best move given the results aggregated.
	 */
	AddResult(Result)
}

type resultAggregatorImpl struct {
	rootState game.State
	resultMap map[interface{}]game.Value  // interface{} is the return type of game.State.Id()
	resultChan chan Result
	remainingCalls int
	callback func(game.State)
}

func (r *resultAggregatorImpl) SetExpectedResultsCount(numCalls int) {
	r.remainingCalls = numCalls
}

func (r *resultAggregatorImpl) SetCallback(callback func(game.State)) {
	r.callback = callback
}

func (r *resultAggregatorImpl) AddResult(result Result) {
	log.Printf("AddResult %+v\n", result.State.Id())
	r.resultChan <- result
}

// this should be run exactly once in a separate goroutine so that map accesses
// are sync'd
func (r *resultAggregatorImpl) aggregate() {
	for {
		result := <-r.resultChan
		r.resultMap[result.State.Id()] = result.Value
		r.remainingCalls--
		if r.remainingCalls == 0 {
			_, bestMove, _ := minimax.AlphaBetaWithValueMap(r.rootState, r.resultMap)
			r.callback(bestMove)
		}
	}
}

func NewResultAggregator(rootState game.State, numCalls int, callback func(game.State)) ResultAggregator {
	aggregator := new(resultAggregatorImpl)
	aggregator.rootState = rootState
	aggregator.resultMap = make(map[interface{}]game.Value)
	aggregator.resultChan = make(chan Result, numCalls)
	aggregator.SetExpectedResultsCount(numCalls)
	aggregator.SetCallback(callback)
	go aggregator.aggregate()
	return aggregator
}

