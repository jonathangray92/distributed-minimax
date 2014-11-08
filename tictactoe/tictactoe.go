package tictactoe

import "github.com/jonathangray92/distributed-minimax/naive_minimax"

type player int

const (
	none player = iota
	x
	o
)

type Node struct {
	player player
	board  [3][3]player
}

func (n Node) Value() naive_minimax.Value {
	switch n.winner() {
	case none:
		return 0
	case x:
		return 1
	}
}

func (n Node) winner() player {
	for _, v := range [...][3]player{
		{n.board[0][0], n.board[0][1], n.board[0][2]},
		{n.board[1][0], n.board[1][1], n.board[1][2]},
		{n.board[2][0], n.board[2][1], n.board[2][2]},

		{n.board[0][0], n.board[1][0], n.board[2][0]},
		{n.board[0][1], n.board[1][1], n.board[2][1]},
		{n.board[0][2], n.board[1][2], n.board[2][2]},

		{n.board[0][0], n.board[1][1], n.board[2][2]},
		{n.board[0][2], n.board[1][1], n.board[2][0]},
	} {
		if v[0] != none && v[0] == v[1] && v[1] == v[2] {
			return v[0]
		}
	}

	return none
}
