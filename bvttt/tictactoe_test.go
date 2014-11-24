package bvttt

import (
	"reflect"
	"testing"
	"github.com/jonathangray92/distributed-minimax/game"
)

var valueTests = []struct {
	state State
	expected game.Value
}{
	{State{IsXMove: false, X: 0000, Y: 0000}, TIE},
	{State{IsXMove: false, X: 0070, Y: 0402}, XWINS},
	{State{IsXMove: true,  X: 0106, Y: 0631}, YWINS},
	{State{IsXMove: false, X: 0523, Y: 0254}, XWINS},
	{State{IsXMove: false, X: 0520, Y: 0250}, TIE},
}

func TestValue(t *testing.T) {
	for _, test := range valueTests {
		if val := test.state.Value(); val != test.expected {
			t.Errorf("State {X: 0%o, Y: 0%o} has value %v, expected %v",
					test.state.X,
					test.state.Y,
					val,
					test.expected)
		}
	}
}

func TestEncodeState(t *testing.T) {
	for _, test := range valueTests {
		bytes, _ := test.state.EncodeState()
		decoded := State{}
		decoded.DecodeState(bytes)
		if decoded != test.state {
			t.Errorf("State %v encoded to %v, decoded to %v", test.state, bytes, decoded)
		}
	}
}

var moveTests = []struct {
	state State
	expected map[State]struct{}
}{
	{state: State{IsXMove: true, X: 0050, Y: 0600}, expected: map[State]struct{}{
		// all possible moves
		State{IsXMove: false, X: 0051, Y: 0600}: struct{}{},
		State{IsXMove: false, X: 0052, Y: 0600}: struct{}{},
		State{IsXMove: false, X: 0054, Y: 0600}: struct{}{},
		State{IsXMove: false, X: 0070, Y: 0600}: struct{}{},
		State{IsXMove: false, X: 0150, Y: 0600}: struct{}{},
	}},
	{state: State{IsXMove: false, X: 0111, Y: 0022}, expected: map[State]struct{}{
		// no moves for game that has already been won
	}},
}

func TestMoves(t *testing.T) {
	for _, test := range moveTests {
		nextStateIter := test.state.Moves()
		stateSet := make(map[State]struct{})
		for {
			nextState := nextStateIter()
			if nextState == nil {
				break
			}
			stateSet[*nextState.(*State)] = struct{}{}
		}
		if !reflect.DeepEqual(stateSet, test.expected) {
			t.Errorf("Got %v, expected %v", stateSet, test.expected)
		}
	}
}
