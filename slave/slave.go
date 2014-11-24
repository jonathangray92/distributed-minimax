package main

import (
	"log"
	proto "code.google.com/p/goprotobuf/proto"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
	"github.com/jonathangray92/distributed-minimax/game"
	"github.com/jonathangray92/distributed-minimax/bvttt"
	minimax "github.com/jonathangray92/distributed-minimax/naiveminimax"
)

type Result struct {
	State []byte
	Value game.Value
}

func getWork(stub *rpc.SlaveServiceClient, result *Result) (*bvttt.State, *Result, error) {
	// create request
	var request rpc.GetWorkRequest
	var response rpc.GetWorkResponse
	if result != nil {
		request.Result = &rpc.GetWorkRequest_Result{
			State: result.State,
			Value: proto.Int32(int32(result.Value)),
		}
	}
	// make request
	err := stub.GetWork(&request, &response)
	if err != nil {
		log.Fatalf("master service responded with error: %v\n", err)
		return nil, nil, err
	}
	// deserialize game state in response
	state := new(bvttt.State)
	err = state.DecodeState(response.GetState())
	if err != nil {
		log.Fatalf("error in decoding state: %v\n", err)
		return nil, nil, err
	}
	return state, &Result{State: response.GetState()}, nil
}

func main() {
	stub, client, err := rpc.DialSlaveService("tcp", "localhost:14782")
	if err != nil {
		log.Fatalf("dialing master service failed: %v\n", err)
	}
	defer client.Close()

	var lastResult *Result = nil
	for {
		state, result, _ := getWork(stub, lastResult)
		log.Printf("work: %+v\n", state)

		// analyze state and save the value to the lastResult
		maxDepth := 4
		value, _ := minimax.Minimax(state, maxDepth)
		result.Value = value
		lastResult = result
		log.Printf("value: %v\n", value)
	}
}
