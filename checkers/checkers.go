package checkers

import "bytes"

import "github.com/jonathangray92/distributed-minimax/game"

// Represents a game state for checkers.
type BasicState struct {
	Turn                    Turn
	CurrentPlayer, Opponent Pieces
}

// Returns the initial game state for checkers.
func NewBasicGame() BasicState {
	return BasicState{
		Turn: Player1,
		CurrentPlayer: Pieces{
			Pawns: 0x0000000000FFFFFF & checkerMask,
		},
		Opponent: Pieces{
			Pawns: 0xFFFFFF0000000000 & checkerMask,
		},
	}
}

func (s BasicState) Value() game.Value {
	player1, player2 := s.CurrentPlayer, s.Opponent
	if !s.Turn.isPlayer1 {
		player1, player2 = player2, player1
	}

	var accum game.Value

	for pos := startPos; pos.valid(); pos = pos.succ() {
		if piece := player1.pieceAt(pos); piece != NoPieces {
			if piece.hasKing() {
				accum += 5
			} else {
				accum += 2
			}
		} else if piece := player2.pieceAt(pos); piece != NoPieces {
			if piece.hasKing() {
				accum -= 5
			} else {
				accum -= 2
			}
		}
	}

	return accum
}

func (s BasicState) MaximizingPlayer() bool       { return s.Turn.isPlayer1 }
func (s BasicState) EncodeState() ([]byte, error) { panic("Unimplemented") }
func (s BasicState) DecodeState([]byte) error     { panic("Unimplemented") }
func (s BasicState) Id() interface{}              { return s }

func (s BasicState) MoveIterator() game.StateIterator {
	if s.CurrentPlayer == NoPieces || s.Opponent == NoPieces {
		return func() game.State { return nil }
	}

	nextPos := startPos

	var jumps Jumps
	var steps Steps
	var moves Moves = &jumps

	var foundMove, pushingSteps bool

	iter := func() game.State {
	outer:
		for {
			for hasMove := moves.HasMove(); hasMove; {
				foundMove = true
				if move, ok := moves.AdvanceMoves(); ok {
					return move
				}
			}

			for nextPos.valid() {
				pos := nextPos
				nextPos = nextPos.succ()
				if s.CurrentPlayer.contains(pos) {
					moves.InitiateMove(s, s.CurrentPlayer.pieceAt(pos))
					continue outer
				}
			}

			if !foundMove && !pushingSteps {
				nextPos = startPos
				moves = &steps
				pushingSteps = true
			} else {
				return nil
			}

		}
	}

	return iter
}

func (s BasicState) String() string {
	var player1, player2 Pieces
	if s.Turn.isPlayer1 {
		player1, player2 = s.CurrentPlayer, s.Opponent
	} else {
		player1, player2 = s.Opponent, s.CurrentPlayer
	}

	var buf stateStringBuilder
	for pos := startPos; pos.valid(); pos = pos.succ() {
		piece := ' '
		switch {
		case player1.contains(pos):
			if player1.pieceAt(pos).hasKing() {
				piece = 'ⓧ' // black king
			} else {
				piece = 'x' // black pawn
			}
		case player2.contains(pos):
			if player2.pieceAt(pos).hasKing() {
				piece = 'ⓞ' // white king
			} else {
				piece = 'o' // white pawn
			}
		}

		buf.writePiece(piece)
	}

	return buf.String()
}

// Builds the string representation of a board from a sequence of pieces.
type stateStringBuilder struct {
	buf      bytes.Buffer
	row, col int
}

// Writes the next piece to the string being built. Each piece is represented
// as a unicode rune.
func (sb *stateStringBuilder) writePiece(piece rune) {
	if sb.col == 0 {
		sb.buf.WriteRune('|')
	}

	if sb.row%2 == 0 {
		sb.buf.WriteString(" |")
		sb.buf.WriteRune(piece)
		sb.buf.WriteRune('|')
	} else {
		sb.buf.WriteRune(piece)
		sb.buf.WriteString("| |")
	}

	sb.col++
	if sb.col >= size/2 {
		sb.col = 0
		sb.row++
		sb.buf.WriteRune('\n')
	}
}

// Returns the constructed string.
func (sb *stateStringBuilder) String() string {
	return sb.buf.String()
}

// The number of rows and columns in the board.
const size = 8

// A bitvector representation of a board, where a 1 indicates that the
// indictated square is occupied, and a 0 indicates it is not. Does not
// distinguish between different types of pieces or their owners. The
// representation is row major with the highest bit representing the top
// leftmost square on the board from player 1's perspective.
type positions uint64

const (
	startPos    positions = 1 << (size*size - 2) // top leftmost position
	checkerMask positions = 0x55AA55AA55AA55AA   // ((0101 0101)(1010 1010)){4}
	homeRow     positions = 0x00000000000000FF   // (0000 0000){7}(1111 1111)
	endRow      positions = 0xFF00000000000000   // (1111 1111)(0000 0000){7}
)

// Returns the next valid checkers position. Progresses from left to right
// and top to bottom from player 1's perspective.
func (pos positions) succ() positions {
	pos = pos >> 1
	for pos&^checkerMask != 0 {
		pos = pos >> 1
	}
	return pos
}

// Returns whether pos shows any positions as occupied.
func (pos positions) valid() bool { return pos != 0 }

// Represents which player's turn it is.
type Turn struct {
	isPlayer1  bool      // Whether the current player is player 1.
	kingingRow positions // The row where the current player's pawns would be promoted to kings.
}

// Returns the opposite Turn from turn.
func (turn Turn) toggle() Turn {
	return Turn{isPlayer1: !turn.isPlayer1, kingingRow: turn.kingingRow ^ (homeRow | endRow)}
}

// Value representing player 1's turn.
var Player1 = Turn{isPlayer1: true, kingingRow: endRow}

// Value representing player 2's turn.
var Player2 = Turn{isPlayer1: false, kingingRow: homeRow}
