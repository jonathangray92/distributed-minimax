package main

import (
	"log"
	"sync/atomic"
	"github.com/jonathangray92/distributed-minimax/game"
	rpc "github.com/jonathangray92/distributed-minimax/proto"

	// change this line and re-build to make this work with different games
	// this must be a package which includes a type called State that implements
	// the State interface from the "game" package in this repo
	gameImpl "github.com/jonathangray92/distributed-minimax/connect4"
)

// magic protobuf declarations
type UserService int

// User-facing RPC code
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

    // get a slice of all slave channels
    // we need this to know how many slave channels we have
    slaveChannels := getChanContentsAsSlice(slaveChanChan)

    // create at least as many jobs as there are slaves to do work
    jobs := createJobsFromRootState(rootState, len(slaveChannels))

	// create a channel on which the best move will be sent, and initialize
    // the result aggregator to expect one result for every job we have just
	// created
    bestMoveChan := make(chan []byte)
    resultAggregator = NewResultAggregator(rootState, len(jobs),
        func(bestMove game.State) {
			if bestMove == nil {
				log.Fatalf("resultAggregator returned nil bestMove\n")
			}
			// encode best move
			encodedState, err := bestMove.EncodeState()
			if err != nil {
				log.Fatalf("error encoding best move: %+v\n", bestMove)
			}
			// send encoded best move along channel to the user
			bestMoveChan <- encodedState
            // set the flag to allow new user requests
            atomic.StoreUint32(&workInProgress, 0)
        })

    // distribute jobs to slaves
    distributeJobsToSlaves(jobs, slaveChannels, request.GetTimeLimitMillis())

    // wait for best move and respond to user
    response.Move = <-bestMoveChan
	log.Printf("returning bestMove to user\n")
    return nil
}

// get a slice of all slave channels
// slaveChannels := getChanContentsAsSlice(slaveChanChan)
// `go get generics`
func getChanContentsAsSlice(c chan chan slaveWork) []chan slaveWork {
    slice := make([]chan slaveWork, 0, MAX_SLAVES)
    for {
        select {
            case e := <-c:
                slice = append(slice, e)
            default:
                return slice
        }
    }
}

// create at least as many jobs as there are slaves to do work
func createJobsFromRootState(rootState game.State, numSlaves int) []game.State {

    log.Printf("expanding rootState for %v slaves\n", numSlaves)

    // allocate a jobs slice and initialize with root state
    jobs := make([]game.State, 0, 2*numSlaves)
    jobs = append(jobs, rootState)

	// pop the head node and push its children until we have enough jobs
	//
	// N.B. we require that the root state is expanded at least one ply, even
	// if there is only one slave. Currently, the slaves send only the
	// resulting value of each job state, not the corresponding best move. If a
	// single slave works the root state and sends its value, the master does
	// not know which move to make. The simplest workaround is to ensure that
	// the root state is never worked directly by always expanding at least one
	// ply
    for ; len(jobs) == 1 || len(jobs) < numSlaves; {
        head := jobs[0]
        jobs = jobs[1:]
        stateIter := head.MoveIterator()
        for s := stateIter(); s != nil; s = stateIter() {
            jobs = append(jobs, s)
        }
    }

    log.Printf("created %v jobs from root state\n", len(jobs))
    return jobs
}

// Given a list of jobs and a list of slave channels, distribute jobs to slaves.
// Constructs a slaveWork instance for each slave, populates it, and sends it on
// the slave channel.
func distributeJobsToSlaves(jobs []game.State, slaveChannels []chan slaveWork, timeLimitMillis uint64) {

    numJobs := len(jobs)
    numSlaves := len(slaveChannels)
    log.Printf("distributing %v jobs to %v slaves\n", numJobs, numSlaves)

    // check that we have 1+ jobs per slave
    if numSlaves > numJobs {
        panic("more slaves than jobs")
    }

    // assign jobs to slaves round-robin
    // with 5 jobs and 3 slaves:
    // - slave 0 gets jobs 0, 3
    // - slave 1 gets jobs 1, 4
    // - slave 2 gets job 2
    maxJobsPerSlave := numJobs / numSlaves + 1
    for i, slaveChannel := range slaveChannels {
        work := slaveWork{states: make([]game.State, 0, maxJobsPerSlave), timeLimitMillis: timeLimitMillis}
        for j := i; j < numJobs; j += numSlaves {
            work.states = append(work.states, jobs[j])
        }
        slaveChannel <- work
    }
}
