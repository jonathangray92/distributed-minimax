package main

import (
	"log"
	"time"
	proto "code.google.com/p/goprotobuf/proto"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
	"github.com/jonathangray92/distributed-minimax/game"
	minimax "github.com/jonathangray92/distributed-minimax/minimax"

	// change this line and re-build to make this work with different games
	// this must be a package which includes a type called State that implements
	// the State interface from the "game" package in this repo
	gameImpl "github.com/jonathangray92/distributed-minimax/connect4"
)

type Result struct {
	State game.State
	Value game.Value
	NumStatesAnalyzed int
}

// helper to handle making a GetWork RPC call
// returns a list of game states to work and a time limit
func getWork(stub *rpc.SlaveServiceClient, results []Result) ([]game.State, time.Duration, error) {

	var request rpc.GetWorkRequest
	var response rpc.GetWorkResponse
	log.Printf("calling GetWork\n")

	// encode results of previous round
	request.Result = encodeResults(results)

	// make request
	err := stub.GetWork(&request, &response)
	if err != nil {
		log.Fatalf("master service responded with error: %v\n", err)
		return nil, 0, err
	}

	// deserialize response states
	return decodeGetWorkResponse(response)
}

// Main function: connect to master node and do work endlessly
func main() {

	// connect to master node
	stub, client, err := rpc.DialSlaveService("tcp", "localhost:14782")
	if err != nil {
		log.Fatalf("dialing master service failed: %v\n", err)
	}
	defer client.Close()

	// endlessly get work, do work, and return results
	var lastResults []Result = nil
	for {
		// get job states and time limit
		states, timeLimit, err := getWork(stub, lastResults)
		if err != nil {
			log.Fatalf("getWork returned err %v\n", err)
		}
		log.Printf("analyzing %v states in %v\n", len(states), timeLimit)

		// spawn a goroutine for each job
		// each goroutine is passed the same channel over which to send results
		resultsChan := make(chan Result, len(states))
		for _, state := range states {
			go func(state game.State) {
				value, _, numStatesAnalyzed := minimax.TimeLimitedAlphaBeta(state, timeLimit)
				resultsChan <- Result{
					State: state,
					Value: value,
					NumStatesAnalyzed: numStatesAnalyzed,
				}
			}(state)
		}

		// wait for the correct number of results to come over the channel
		lastResults = make([]Result, len(states))
		for i := range states {
			lastResults[i] = <-resultsChan
		}
	}
}

// convert a slice of Result structs to protobuf format
func encodeResults(results []Result) []*rpc.GetWorkRequest_Result {

	// allocate space
	encodedResults := make([]*rpc.GetWorkRequest_Result, len(results))

	// encode each result
	for i, result := range results {
		encodedState, err := result.State.EncodeState()
		if err != nil {
			log.Fatalf("error encoding result state %+v\n", result.State)
		}
		encodedResults[i] = &rpc.GetWorkRequest_Result{
			State: encodedState,
			Value: proto.Int32(int32(result.Value)),
			NumStatesAnalyzed: proto.Int64(int64(result.NumStatesAnalyzed)),
		}
	}

	return encodedResults
}

// convert a protobuf GetWorkResponse to a slice of states to work and a time
// limit
func decodeGetWorkResponse(response rpc.GetWorkResponse) ([]game.State, time.Duration, error) {

	// decode each work state
	encodedStates := response.GetState()
	decodedStates := make([]game.State, len(encodedStates))
	for i, s := range encodedStates {
		decodedStates[i] = new(gameImpl.State)
		err := decodedStates[i].DecodeState(s)
		if err != nil {
			log.Fatalf("error in decoding state: %v\n", err)
			return nil, 0, err
		}
	}

	// convert time limit to time.Duration
	timeLimit := time.Duration(response.GetTimeLimitMillis()) * time.Millisecond

	return decodedStates, timeLimit, nil
}

