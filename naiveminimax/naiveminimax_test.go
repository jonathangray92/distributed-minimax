package naiveminimax_test

import (
	"encoding/gob"
	"os"
	"testing"

	"github.com/jonathangray92/distributed-minimax/tictactoe"
)

func TestMinimax(t *testing.T) {
	const FileName = "tictactoe_table.gob"
	file, err := os.Open(FileName)
	if err != nil {
		t.Fatalf("error opening test file: %v\n", err)
	}
	defer file.Close()

	var table1 tictactoe.Table
	gob.NewDecoder(file).Decode(&table1)

	table2 := tictactoe.MakeTable()

	if len(table1) != len(table2) {
		t.Fatalf("Incorrect number of beginning-of-move states for tictactoe. Expecting %v, found %v.\n", len(table1), len(table2))
	}

	for k, row1 := range table1 {
		row2 := table2[tictactoe.Normalize(k)]

		if len(row1) != len(row2) {
			t.Fatalf("Incorrect number of possible moves for following start state. Expecting %v, found %v.\n%v", len(row1), len(row2), k)
		}

	outer:
		for _, state2 := range row2 {
			for _, state1 := range row1 {
				if tictactoe.Normalize(state1) == tictactoe.Normalize(state2) {
					continue outer
				}
			}

			t.Fatalf("Incorrect move possibility. Start state:\n%v\nEnd state:\n%v", k, state2)
		}
	}
}
