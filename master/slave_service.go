package main

import (
	"log"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
	"github.com/jonathangray92/distributed-minimax/game"

    // change this line and re-build to make this work with different games
    // this must be a package which includes a type called State that implements
    // the State interface from the "game" package in this repo
    gameImpl "github.com/jonathangray92/distributed-minimax/connect4"
)

// magic protobuf declaration
type SlaveService int

// RPC request code. Executed in a goroutine when a slave calls GetWork
//// Pops a state from the work queue and sends it to the slave for processing.
//  
// If the slave has included a Result from past work, submit this work to the
// aggregator.
func (t *SlaveService) GetWork(request *rpc.GetWorkRequest, response *rpc.GetWorkResponse) error {
    log.Printf("slave connected with %v past results\n", len(request.GetResult()))

    // decode and handle results of previous work, if provided
	for _, encodedResult := range request.GetResult() {
		decodedResult := Result{
			State: new(gameImpl.State),
			Value: game.Value(encodedResult.GetValue()),
			NumStatesAnalyzed: uint64(encodedResult.GetNumStatesAnalyzed()),
		}
		decodedResult.State.DecodeState(encodedResult.GetState())
		resultAggregator.AddResult(decodedResult)
	}

    // make a channel on which work will be provided, and wait for it
    c := make(chan slaveWork, 1)
    slaveChanChan <- c
    work := <-c

    // encode all states
    stateEncodings := make([][]byte, len(work.states))
    for i, state := range work.states {
        bs, err := state.EncodeState()
        if err != nil {
            log.Printf("error: %v\n", err)
            return err
        }
        stateEncodings[i] = bs
    }

    // build response
    response.State = stateEncodings
    response.TimeLimitMillis = &work.timeLimitMillis
    return nil
}

