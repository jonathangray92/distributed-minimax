package tictactoe

func Equivalent(a, b State) bool {
	if a == b || a == transpose(b) {
		return true
	}

	b_flip_v := flip_v(b)
	if a == b_flip_v || a == transpose(b_flip_v) {
		return true
	}

	b_flip_h := flip_h(b)
	if a == b_flip_h || a == transpose(b_flip_h) {
		return true
	}

	b_rot_180 := flip_h(b_flip_v)
	if a == b_rot_180 || a == transpose(b_rot_180) {
		return true
	}

	return false
}

func flip_v(s State) State {
	out := s
	for i := range s.Board {
		for j := range s.Board[i] {
			out.Board[2-i][j] = s.Board[i][j]
		}
	}
	return out
}

func flip_h(s State) State {
	out := s
	for i := range s.Board {
		for j := range s.Board[i] {
			out.Board[i][2-j] = s.Board[i][j]
		}
	}
	return out
}

func transpose(s State) State {
	out := s
	for i := range s.Board {
		for j := range s.Board[i] {
			out.Board[j][i] = s.Board[i][j]
		}
	}
	return out
}
