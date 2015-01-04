package checkers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jonathangray92/distributed-minimax/game"
	minimax "github.com/jonathangray92/distributed-minimax/naiveminimax"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestCheckers(t *testing.T) {
	var start game.State = NewRepeatLimitedGame()
	var chosenMove game.State
	for state := start; state != nil; state = chosenMove {
		state.(*RepeatLimitedState).UseState()
		fmt.Println(state)
		_, chosenMove = minimax.Minimax(state, 4)
	}
}
