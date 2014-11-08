// This package defines a naive, sequential implmentation of the minimax algorithm.
package naive_minimax

// A Node represents a game state in a user-defined game.
type Node interface {
	Value() Value           // The heuristic value of the game state.
	Children() NodeIterator // An iterator over game states one move away from the current game state.
}

// A NodeIterator provides access to a sequence of Nodes.
type NodeIterator interface {
	Node() Node         // The current Node.
	Next() NodeIterator // A NodeIterator referencing the next Node in the sequence.
}

// Value is the type of the heuristic values of game states.
type Value int

// Minimax implements the minimax algorithm. It determines the best possible outcome of a
// game sequence starting at the provided node by evaluating the game tree rooted at the
// node up to the provided depth. Leaves of the evaluated game tree are valued based on
// their Node.Value. If maximizingPlayer is true, then the best outcome will be the one
// with the maximum value. Otherwise, it will be the one with the minimum value.
func Minimax(node Node, depth int, maximizingPlayer bool) Value {
	child := node.Children()

	if depth == 0 || child == nil {
		return node.Value()
	}

	var better func(a, b Value) Value
	if maximizingPlayer {
		better = max
	} else {
		better = min
	}

	best := Minimax(child.Node(), depth-1, !maximizingPlayer)

	for {
		child = child.Next()
		if child == nil {
			break
		}

		val := Minimax(child.Node(), depth-1, !maximizingPlayer)
		best = better(best, val)
	}

	return best
}

func min(a, b Value) Value {
	if a < b {
		return a
	}
	return b
}

func max(a, b Value) Value {
	if a > b {
		return a
	}
	return b
}
