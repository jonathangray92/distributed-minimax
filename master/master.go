package main

import (
	"fmt"
	"log"
	"github.com/jonathangray92/distributed-minimax/bvttt"
	rpc "github.com/jonathangray92/distributed-minimax/proto"
//	"code.google.com/p/goprotobuf/proto"
)

var STATE1 = bvttt.State{IsXMove: true, X: 0050, Y: 0600}  // X should win in 1
var STATE2 = bvttt.State{IsXMove: true, X: 0401, Y: 0020}  // X should win in 2
var STATE3 = bvttt.State{IsXMove: false, X: 0150, Y: 0600} // X should win in 2

var workQueue = make(chan bvttt.State, 100)

type SlaveService int

func (t *SlaveService) GetWork(request *rpc.GetWorkRequest, response *rpc.GetWorkResponse) error {
	log.Printf("request: %v\n", request.GetResult())
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
	workQueue <- STATE1
	workQueue <- STATE2
	workQueue <- STATE3

	// listen on known port
	port := 14782
	log.Printf("Slave service listening on port %v\n", port)
	rpc.ListenAndServeSlaveService("tcp", fmt.Sprint("localhost:", port), new(SlaveService))
}
