package naive_minimax_test

import "testing"
import . "github.com/jonathangray92/distributed-minimax/tictactoe"
import "github.com/jonathangray92/distributed-minimax/naive_minimax"

func TestMinimax(t *testing.T) {
	for i := range tests {
		_, output := naive_minimax.Minimax(tests[i].input, 10)
		if output, ok := output.(State); !ok || !Equivalent(output, tests[i].output) {
			t.Errorf("test %v failed.", i)
		}
	}
}

type test struct {
	input, output State
}

var tests = [...]test{
	0: {
		input: State{
			Player: X,
			Board: [3][3]Player{
				{X, None, X},
				{None, O, O},
				{O, None, None},
			},
		},
		output: State{
			Player: O,
			Board: [3][3]Player{
				{X, X, X},
				{None, O, O},
				{O, None, None},
			},
		},
	},

	1: {
		input: State{
			Player: X,
			Board: [3][3]Player{
				{X, O, None},
				{None, O, X},
				{None, None, None},
			},
		},
		output: State{
			Player: O,
			Board: [3][3]Player{
				{X, O, None},
				{None, O, X},
				{None, X, None},
			},
		},
	},

	2: {
		input: State{
			Player: O,
			Board: [3][3]Player{
				{None, None, None},
				{X, O, X},
				{None, None, None},
			},
		},
		output: State{
			Player: X,
			Board: [3][3]Player{
				{O, None, None},
				{X, O, X},
				{None, None, None},
			},
		},
	},
}
