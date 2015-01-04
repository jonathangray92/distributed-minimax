package main

import (
	"fmt"
	"log"
	"github.com/jonathangray92/distributed-minimax/bvttt"
	"github.com/jonathangray92/distributed-minimax/game"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
//	"code.google.com/p/goprotobuf/proto"
)

type Result *rpc.GetWorkRequest_Result

// mock user-submitted state
// expected result: X should win in 2
var ROOT_STATE = &bvttt.State{IsXMove: true, X: 0401, Y: 0120}

// global variables
var workQueue = make(chan game.State, 100)
var resultAggregator ResultAggregator

// magic protobuf declaration
type SlaveService int

/**
 * RPC request code. Executed in a goroutine when a slave calls GetWork
 *
 * Pops a state from the work queue and sends it to the slave for processing.
 *
 * If the slave has included a Result from past work, submit this work to the
 * aggregator.
 */
func (t *SlaveService) GetWork(request *rpc.GetWorkRequest, response *rpc.GetWorkResponse) error {
	log.Printf("request: %v\n", request.GetResult())
	if (request.GetResult() != nil) {
		resultAggregator.AddResult(request.GetResult())
	}
	state := <-workQueue
	bs, err := state.EncodeState()
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}
	response.State = bs
	return nil
}

/**
 * Given a "root" game state, find all the states one ply deep and add them to
 * "queue".
 *
 * Return the number of states appended to the queue
 */
func populateWorkQueueFromRootState(queue chan game.State, root game.State) int {
	// get states one ply deep
	stateIter := root.MoveIterator()
	count := 0
	for nextState := stateIter(); nextState != nil; nextState = stateIter() {
		queue <- nextState
		count++
	}
	return count
}

func main() {

	// initialize work queue
	log.Printf("root state: %+v\n", ROOT_STATE)
	numExpectedResults := populateWorkQueueFromRootState(workQueue, ROOT_STATE)

	// initialize result aggregator
	resultAggregator = NewResultAggregator(numExpectedResults,
		func(minState Result, maxState Result) {
			log.Printf("minState: %+v worth %v\n", minState.State, *minState.Value)
			log.Printf("maxState: %+v worth %v\n", maxState.State, *maxState.Value)
			bestMove := new(bvttt.State)
			if ROOT_STATE.MaximizingPlayer() {
				bestMove.DecodeState(maxState.State)
			} else {
				bestMove.DecodeState(minState.State)
			}
			log.Printf("best move: %+v\n", bestMove)
		})

	// listen on known port
	port := 14782
	log.Printf("Slave service listening on port %v\n", port)
	rpc.ListenAndServeSlaveService("tcp", fmt.Sprint("localhost:", port), new(SlaveService))
}
