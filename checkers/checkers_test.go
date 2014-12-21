package checkers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jonathangray92/distributed-minimax/game"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestCheckers(t *testing.T) {
	start := game.State(NewGame())

	var chosenMove game.State
	for state := start; state != nil; state = chosenMove {
		fmt.Println(state)
		iterMoves := state.MoveIterator()
		chosenMove = iterMoves()
		for move := chosenMove; move != nil; move = iterMoves() {
			if rand.Float32() < 0.2 {
				chosenMove = move
			}
		}
	}
}
