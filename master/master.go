package main

import (
	"fmt"
	"log"
	"sync/atomic"
	gameImpl "github.com/jonathangray92/distributed-minimax/connect4"
	"github.com/jonathangray92/distributed-minimax/game"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
//	"code.google.com/p/goprotobuf/proto"
)

type Result *rpc.GetWorkRequest_Result

// global variables
var workQueue = make(chan game.State, 100)
var resultAggregator ResultAggregator
var workInProgress uint32  // 0 iff no ongoing work in progress; use atomic methods

// magic protobuf declaration
type SlaveService int
type UserService int

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
 * User-facing RPC code
 */
func (t *UserService) DoWork(request *rpc.DoWorkRequest, response *rpc.DoWorkResponse) error {
	log.Printf("DoWork called")

	// check and set the workInProgress global flag so that only one user's
	// request is ever being processed
	if !atomic.CompareAndSwapUint32(&workInProgress, workInProgress, 1) {
		log.Printf("work already in progress, returning nil")
		return nil
	}

	// decode state from request
	rootState := new(gameImpl.State)
	rootState.DecodeState(request.State)

	// the result should be sent on this channel
	bestMoveChan := make(chan []byte)

	// populate work queue
	numExpectedResults := populateWorkQueueFromRootState(workQueue, rootState)

	// initialize result aggregator
	resultAggregator = NewResultAggregator(numExpectedResults,
		func(minState Result, maxState Result) {
			log.Printf("minState: %+v worth %v\n", minState.State, *minState.Value)
			log.Printf("maxState: %+v worth %v\n", maxState.State, *maxState.Value)
			if rootState.MaximizingPlayer() {
				bestMoveChan <- maxState.State
			} else {
				bestMoveChan <- minState.State
			}
			// set the flag to allow new user requests
			atomic.StoreUint32(&workInProgress, 0)
		})

	// wait for best move and respond to user
	response.Move = <-bestMoveChan
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

/**
 * Main function. Boot up User and Slave services.
 */
func main() {

	// listen for slaves on known port
	slavePort := 14782
	log.Printf("Slave service listening on port %v\n", slavePort)
	go rpc.ListenAndServeSlaveService("tcp", fmt.Sprint("localhost:", slavePort), new(SlaveService))

	// listen for user on known port
	// N.B. call in this goroutine so we don't exit immediately
	userPort := 14783
	log.Printf("User service listening on port %v\n", userPort)
	rpc.ListenAndServeUserService("tcp", fmt.Sprint("localhost:", userPort), new(UserService))
}
