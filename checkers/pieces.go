package checkers

// Represents the positions of kings and pawns. Does not distinguish the owner
// of the pieces.
type Pieces struct {
	Pawns, Kings positions
}

// A Pieces value representing a lack of any pieces.
var NoPieces Pieces = Pieces{}

// Returns the union of pieces in p and other.
func (p Pieces) combinedWith(other Pieces) Pieces {
	return Pieces{
		Pawns: p.Pawns | other.Pawns,
		Kings: p.Kings | other.Kings,
	}
}

// Returns the pieces in p at the provided positions.
func (p Pieces) pieceAt(pos positions) Pieces {
	return Pieces{
		Pawns: p.Pawns & pos,
		Kings: p.Kings & pos,
	}
}

// Returns whether p contains a king.
func (p Pieces) hasKing() bool {
	return p.Kings != 0
}

// Returns p with the from Pieces removed and the to Pieces inserted.
func (p Pieces) doMove(from, to Pieces) Pieces {
	return Pieces{
		Pawns: p.Pawns&^from.Pawns | to.Pawns,
		Kings: p.Kings&^from.Kings | to.Kings,
	}
}

// Returns p with the captured pieces removed.
func (p Pieces) capture(captured Pieces) Pieces {
	return Pieces{
		Pawns: p.Pawns &^ captured.Pawns,
		Kings: p.Kings &^ captured.Kings,
	}
}

// Returns whether p contains a piece at any of the provided positions.
func (p Pieces) contains(pos positions) bool {
	return (p.Pawns|p.Kings)&pos != 0
}

// Returns the positions occupied by the pieces in p.
func (p Pieces) positions() positions {
	return p.Pawns | p.Kings
}

// Returns 2 sets of pieces, where all the pieces in p are move forward to the
// left and forward to the right. Pieces where this move would put them off the
// end of the board are discarded.
func (p Pieces) forward() [2]Pieces {
	return [2]Pieces{
		{
			Pawns: (p.Pawns << (size - 1)) & checkerMask,
			Kings: (p.Kings << (size - 1)) & checkerMask,
		},
		{
			Pawns: (p.Pawns << (size + 1)) & checkerMask,
			Kings: (p.Kings << (size + 1)) & checkerMask,
		},
	}
}

// Returns 2 sets of pieces, where all the pieces in p are move backward to the
// left and backward to the right. Pieces where this move would put them off the
// end of the board are discarded.
func (p Pieces) backward() [2]Pieces {
	return [2]Pieces{
		{
			Pawns: (p.Pawns >> (size - 1)) & checkerMask,
			Kings: (p.Kings >> (size - 1)) & checkerMask,
		},
		{
			Pawns: (p.Pawns >> (size + 1)) & checkerMask,
			Kings: (p.Kings >> (size + 1)) & checkerMask,
		},
	}
}

// Returns p with any pawns at the positions in row turned into kings.
func (p Pieces) kingPositions(row positions) Pieces {
	return Pieces{
		Pawns: p.Pawns &^ row,
		Kings: p.Kings | p.Pawns&row,
	}
}
