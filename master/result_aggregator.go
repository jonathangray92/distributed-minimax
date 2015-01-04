package main

type ResultAggregator interface {

	/**
	 * Call this once to set a callback that will execute after AddResult()
	 * has been called the given number of times.
	 *
	 * The callback will be called with two Result args, holding the moves that
	 * yield the min and max value.
	 *
	 * The first arg is the number of AddResult() calls before the callback
	 * is triggered.
	 */
	SetCallback(int, func(Result, Result))

	/**
	 * Call this with a new result delivered by a slave
	 *
	 * If a callback is set and this method has been called the correct number of
	 * times, call the callback with the best move given the results aggregated.
	 */
	AddResult(Result)
}

type resultAggregatorImpl struct {
	remainingCalls int
	minMove, maxMove Result
	callback func(Result, Result)
}

func (r *resultAggregatorImpl) SetCallback(numCalls int, callback func(Result, Result)) {
	r.remainingCalls = numCalls
	r.callback = callback
}

func (r *resultAggregatorImpl) AddResult(result Result) {
	if r.minMove == nil || *result.Value < *r.minMove.Value {
		r.minMove = result
	}
	if r.maxMove == nil || *result.Value > *r.maxMove.Value {
		r.maxMove = result
	}
	r.remainingCalls--
	if r.remainingCalls == 0 {
		go r.callback(r.minMove, r.maxMove)
	}
}

func NewResultAggregator(numCalls int, callback func(Result, Result)) ResultAggregator {
	aggregator := new(resultAggregatorImpl)
	aggregator.SetCallback(numCalls, callback)
	return aggregator
}

