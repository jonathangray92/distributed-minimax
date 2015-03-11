package main

import (
	"testing"
	"github.com/jonathangray92/distributed-minimax/connect4"
)

func TestGetChanContentsAsSlice(t *testing.T) {
	cc := make(chan chan slaveWork, 3)
	cc <- make(chan slaveWork, 1)
	cc <- make(chan slaveWork, 1)
	cc <- make(chan slaveWork, 1)
	slice := getChanContentsAsSlice(cc)
	if len(slice) != 3 {
		t.Fatalf("len(slice) == %v\n", len(slice))
	}
}

// given 1 slave, the root state should be not be used and must be expanded
func TestCreateJobsFromRootStateWith1Slave(t *testing.T) {
	rootState := connect4.NewInitialState()
	jobs := createJobsFromRootState(rootState, 1)
	if len(jobs) == 1 {
		t.Fatal("should have expanded root state")
	}
}

// given 2 slaves, the root state should be expanded 1 ply
func TestCreateJobsFromRootStateWith2Slaves(t *testing.T) {
	rootState := connect4.NewInitialState()
	jobs := createJobsFromRootState(rootState, 2)
	// the root connect 4 state has 7 children
	if len(jobs) != 7 {
		t.Fatalf("expected 7 jobs, got %v\n", len(jobs))
	}
}

// given lots of slaves, should expand some states 2 plies deep
func TestCreateJobsFromRootStateWithManySlaves(t *testing.T) {
	rootState := connect4.NewInitialState()
	jobs := createJobsFromRootState(rootState, 20)
	if len(jobs) < 20 {
		t.Fatalf("expected 20+ jobs, got %v\n", len(jobs))
	}
}
