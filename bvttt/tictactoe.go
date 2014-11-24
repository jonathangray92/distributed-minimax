package bvttt

import (
	"encoding/binary"
	"github.com/jonathangray92/distributed-minimax/game"
)

const (
	XWINS = game.Value(1)
	YWINS = game.Value(-1)
	TIE   = game.Value(0)
)

// defines a set of bitmasks that represent win conditions in tic tac toe
// each state is represented in octal, where the 0th digit is the top-left
// square and the 8th digit is the bottom-right square.
//
// e.g. the octal bitmask 0521 = 0b101010001 represents the following mask:
//
//	 1 0 0
//	 0 1 0
//	 1 0 1
//
// N.B. winMasks cannot be a const, but should be treated like one
//
var winMasks = [...]uint16 {
	// horizontal lines
	0700, 0070, 0007,
	// vertical lines
	0444, 0222, 0111,
	// diagonal lines
	0421, 0124,
}

func boolToUint32(b bool) uint32 {
	if b {
		return 1
	} else {
		return 0
	}
}

// State implements game.State for tictactoe.
type State struct {
	X, Y    uint16
	IsXMove bool
}

// Value is positive if X wins,
func (s *State) Value() game.Value {
	for _, mask := range winMasks {
		if s.X & mask == mask { return XWINS }
		if s.Y & mask == mask { return YWINS }
	}
	return TIE
}

func (s *State) MaximizingPlayer() bool {
	return s.IsXMove
}

// return a new state with the position at oneHotMask played
// DOES NOT check that the move is valid. Assumes the caller checked.
func (s *State) makeMove(oneHotMask uint16) *State {
	if s.IsXMove {
		return &State{IsXMove: false, X: s.X | oneHotMask, Y: s.Y}
	} else {
		return &State{IsXMove: true, X: s.X, Y: s.Y | oneHotMask}
	}
}

func (s *State) MoveIterator() game.StateIterator {
	// check for win condition
	if s.Value() != 0 {
		return func() game.State { return nil }
	}

	usedMask := s.X | s.Y
	oneHotMask := uint16(1 << 9)  // shift right before using

	return func() game.State {
		for oneHotMask > 1 {
			oneHotMask >>= 1
			if usedMask & oneHotMask == 0 {
				return s.makeMove(oneHotMask)
			}
		}
		return nil
	}
}

// Use the following encoding:
// IsXMove: bit 31
// X moves: bits 16-24
// Y moves: bits 0-9
func (s *State) encodeUint32() uint32{
	return (boolToUint32(s.IsXMove) << 31) | (uint32(s.X) << 16) | uint32(s.Y)
}

func (s *State) EncodeState() ([]byte, error) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, s.encodeUint32())
	return bs, nil
}

func (s *State) DecodeState(p []byte) error {
	u := binary.BigEndian.Uint32(p)
	s.IsXMove = (u & (1<<31)) != 0			// bit 31
	s.X = uint16((u & 0x1ff0000) >> 16)		// bits 16-24
	s.Y = uint16((u & 0x1ff))				// bits 0-8
	return nil
}

func (s *State) Id() interface{}              { return s.encodeUint32() }
