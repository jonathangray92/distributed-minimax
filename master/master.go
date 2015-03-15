package main

import (
	"fmt"
	"log"
	"github.com/jonathangray92/distributed-minimax/game"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
//	"code.google.com/p/goprotobuf/proto"
)

// this is not really enforced, just used for slice allocation
const MAX_SLAVES = 100

// global ResultAggregator used by slave threads
var resultAggregator ResultAggregator

// global flag (0 iff no ongoing work in progress)
// use atomic methods!!!
var workInProgress uint32

// A number of game states to analyze, and a time limit.
// Each call to GetWork() will result in one instance of slaveWork being sent
// to the slave node.
type slaveWork struct {
	states []game.State;
	timeLimitMillis uint64;
}

// each slave calling GetWork() adds a channel to slaveChanChan
// when the user calls DoWork(), work for each slave will be sent over the
// channel
var slaveChanChan = make(chan chan slaveWork, MAX_SLAVES)

// Main function. Boot up User and Slave services.
func main() {

	// listen for slaves on known port
	slavePort := 14782
	log.Printf("Slave service listening on port %v\n", slavePort)
	go rpc.ListenAndServeSlaveService("tcp", fmt.Sprint("0.0.0.0:", slavePort), new(SlaveService))

	// listen for user on known port
	// N.B. call in this goroutine so we don't exit immediately
	userPort := 14783
	log.Printf("User service listening on port %v\n", userPort)
	rpc.ListenAndServeUserService("tcp", fmt.Sprint("0.0.0.0:", userPort), new(UserService))
}
