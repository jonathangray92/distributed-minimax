package checkers

// A collection of incomplete moves.
type Moves interface {
	// Inserts a move into the collection. Takes a game state and a piece to
	// move. If more than one move is possible for the provided piece, this
	// will result in each possible move eventually being returned by
	// AdvanceMoves.
	InitiateMove(state State, piece Pieces)

	// Causes moves in the collection to progress. If this results in any of
	// the moves being complete, a completed move--represented as the game
	// state after the move is made--is returned and ok is true.  Otherwise, ok
	// is false. The returned move, if any, is removed from the collection.
	// Additional completed moves which aren't returned on this call to
	// AdvanceMoves are returned from subsequent calls.
	AdvanceMoves() (completeMove State, ok bool)

	// Returns whether there are any active moves in the collection.
	HasMove() bool
}

// A collection of step moves.
type Steps struct{ moveStack }

func (steps *Steps) InitiateMove(state State, piece Pieces) {
	pushRelativeMoves(steps, move{State: state, start: piece, end: piece, captures: NoPieces})
}

func (steps *Steps) AdvanceMoves() (completeMove State, ok bool) {
	return steps.pop().apply(), true
}

// A collection of jump moves.
type Jumps struct{ moveStack }

func (jumps *Jumps) InitiateMove(state State, piece Pieces) {
	pushRelativeMoves(jumps, move{State: state, start: piece, end: piece, captures: NoPieces})
}

func (jumps *Jumps) AdvanceMoves() (completeMove State, ok bool) {
	jump := jumps.pop()

	if !pushRelativeMoves(jumps, jump) {
		return jump.apply(), true
	}

	return State{}, false
}

// A stack of in-progress moves.
type moveStack []move

func (ms *moveStack) HasMove() bool {
	return len(*ms) > 0
}

// Push a move.
func (ms *moveStack) push(moveVal move) {
	*ms = append(*ms, moveVal)
}

// Pop the top move.
func (ms *moveStack) pop() (moveVal move) {
	moveVal, *ms = (*ms)[len(*ms)-1], (*ms)[:len(*ms)-1]
	return moveVal
}

// Represents an in-progress move.
type move struct {
	State             // The initial state for the move.
	start, end Pieces // The start and end positions of the piece to be moved.
	captures   Pieces // The pieces to be captured during the move.
}

// Returns the end state of the move.
func (m move) apply() State {
	return State{
		Turn:          m.Turn.toggle(),
		CurrentPlayer: m.Opponent.capture(m.captures),
		Opponent:      m.CurrentPlayer.doMove(m.start, m.end).kingPositions(m.Turn.kingingRow),
	}
}

// Represents a value that can push moves with a provided direction.
type movesWithDirPusher interface {
	// Pushes moves relative to an in-progress move (from) in the direction
	// indicated by the provided advance function. Returns whether any moves
	// were pushed.
	pushMovesWithDir(from move, advance advanceFunc) (didPush bool)
}

// A function which returns the two possible outcomes of advancing a piece
// either forward or backward. If either possible move in that direction
type advanceFunc func(Pieces) [2]Pieces

// Pushes moves in the appropriate direction: forward for player 1, backward
// for player2, or both if the piece is a pawn. Returns whether any moves were
// pushed.
func pushRelativeMoves(moves movesWithDirPusher, from move) (didPush bool) {
	var didPushForward, didPushBackward bool

	if from.Turn.isPlayer1 || from.start.hasKing() {
		didPushForward = moves.pushMovesWithDir(from, Pieces.forward)
	}

	if !from.Turn.isPlayer1 || from.start.hasKing() {
		didPushBackward = moves.pushMovesWithDir(from, Pieces.backward)
	}

	return didPushForward || didPushBackward
}

// Implement movesWithDirPusher for Steps.

func (steps *Steps) pushMovesWithDir(from move, advance advanceFunc) (didPush bool) {
	newEnds := advance(from.start)
	for i := range newEnds {
		if newEnds[i] != NoPieces {
			if !from.CurrentPlayer.combinedWith(from.Opponent).contains(newEnds[i].positions()) {
				steps.push(move{State: from.State, start: from.start, end: newEnds[i], captures: NoPieces})
				didPush = true
			}
		}
	}
	return didPush
}

// Implement movesWithDirPusher for Jumps.

func (jumps *Jumps) pushMovesWithDir(from move, advance advanceFunc) (didPush bool) {
	skips := advance(from.end)
	for i := range skips {
		newCapture := from.Opponent.capture(from.captures).pieceAt(skips[i].positions())
		newEnd := advance(skips[i])[i]
		if newCapture != NoPieces && newEnd != NoPieces {
			if !from.CurrentPlayer.combinedWith(from.Opponent).contains(newEnd.positions()) {
				jumps.push(move{State: from.State, start: from.start, end: newEnd, captures: from.captures.combinedWith(newCapture)})
				didPush = true
			}
		}
	}
	return didPush
}
