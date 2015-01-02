package tictactoe

import (
	"fmt"

	"github.com/jonathangray92/distributed-minimax/game"
	"github.com/jonathangray92/distributed-minimax/naiveminimax"
)

// Represents one-to-many mapping between States
type Table map[State][]State

// Returns a Table that maps from each possible normalized non-winning State to a list of possible
// moves that might be recommended by the minimax algorithm.
func MakeTable() Table {
	var (
		state       = NewGame()
		stateVal, _ = naiveminimax.Minimax(state, 10)
	)

	return make(Table).initRec(state, stateVal)
}

func (table Table) initRec(state State, bestVal game.Value) Table {
	iterMoves := state.MoveIterator()

	row, exists := table[state]
	if exists {
		return table
	}
	defer func() {
		if len(row) > 0 {
			table[state] = row
		}
	}()

outer:
	for {
		switch move := iterMoves().(type) {
		case State:
			var (
				moveNorm   = Normalize(move)
				moveVal, _ = naiveminimax.Minimax(move, 10)
			)

			for i := range row {
				if moveNorm == Normalize(row[i]) {
					table.initRec(moveNorm, moveVal)
					continue outer
				}
			}

			if moveVal == bestVal {
				row = append(row, move)
			}

			table.initRec(moveNorm, moveVal)
		case nil:
			return table
		default:
			panic(fmt.Errorf("main.Walk: wrong type of move from a State: %T\n", move))
		}
	}
}
