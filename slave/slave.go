package main

import (
	"log"
	"time"
	proto "code.google.com/p/goprotobuf/proto"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
	"github.com/jonathangray92/distributed-minimax/game"
	gameImpl "github.com/jonathangray92/distributed-minimax/connect4"
	minimax "github.com/jonathangray92/distributed-minimax/minimax"
)

type Result struct {
	State []byte
	Value game.Value
}

func getWork(stub *rpc.SlaveServiceClient, result *Result) (*gameImpl.State, time.Duration, *Result, error) {
	// create request
	var request rpc.GetWorkRequest
	var response rpc.GetWorkResponse
	if result != nil {
		request.Result = &rpc.GetWorkRequest_Result {
			State: result.State,
			Value: proto.Int32(int32(result.Value)),
		}
	}
	// make request
	log.Printf("GetWork\n")
	err := stub.GetWork(&request, &response)
	log.Printf("GotWork\n")
	if err != nil {
		log.Fatalf("master service responded with error: %v\n", err)
		return nil, 0, nil, err
	}
	// deserialize game state and time limit from response
	state := new(gameImpl.State)
	err = state.DecodeState(response.GetState())
	if err != nil {
		log.Fatalf("error in decoding state: %v\n", err)
		return nil, 0, nil, err
	}
	timeLimitDuration := time.Duration(response.GetTimeLimitMillis()) * time.Millisecond
	log.Printf("timeLimitDuration %v\n", timeLimitDuration)

	return state, timeLimitDuration, &Result{State: response.GetState()}, nil
}

func main() {
	stub, client, err := rpc.DialSlaveService("tcp", "localhost:14782")
	if err != nil {
		log.Fatalf("dialing master service failed: %v\n", err)
	}
	defer client.Close()

	var lastResult *Result = nil
	for {
		state, timeLimit, result, err := getWork(stub, lastResult)
		if err != nil {
			log.Fatalf("getWork returned err %v\n", err)
		}
		log.Printf("work: %+v\n", state)

		// analyze state and save the value to the lastResult
		log.Printf("starting minimax with time limit %v\n", timeLimit)
		value, _, numStatesAnalyzed := minimax.TimeLimitedAlphaBeta(state, timeLimit)
		result.Value = value
		lastResult = result
		log.Printf("value: %v\n", value)
		log.Printf("numStatesAnalyzed: %v\n", numStatesAnalyzed)
	}
}
