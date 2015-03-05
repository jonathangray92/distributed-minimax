package main

import (
	"log"
	"time"
	gameImpl "github.com/jonathangray92/distributed-minimax/connect4"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
)

func makeRequest(stub *rpc.UserServiceClient, state *gameImpl.State, timeLimitMillis uint64) *gameImpl.State {
	var request rpc.DoWorkRequest
	var response rpc.DoWorkResponse

	// construct request
	request.State, _ = state.EncodeState()
	request.TimeLimitMillis = &timeLimitMillis

	// make blocking request and calculate end-to-end latency
	log.Printf("making request, time limit %v millis\n", timeLimitMillis)
	rpcStartTime := time.Now()
	err := stub.DoWork(&request, &response)
	rpcDuration := time.Since(rpcStartTime)
	log.Printf("rpc finished in %v\n", rpcDuration)
	if err != nil {
		log.Fatalf("DoWork rpc failed: %v\n", err)
	}

	// decode response
	nextMove := new(gameImpl.State)
	nextMove.DecodeState(response.Move)
	return nextMove
}

func main() {

	// connect to master node and get rpc stub
    stub, client, err := rpc.DialUserService("tcp", "localhost:14783")
    if err != nil {
        log.Fatalf("dialing master's user service failed: %v\n", err)
    }
    defer client.Close()

	// send a request
	nextMove := makeRequest(stub, gameImpl.NewInitialState(), 1000)
	log.Printf("best move: %+v\n", nextMove)
}
