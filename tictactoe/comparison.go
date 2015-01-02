package tictactoe

func Normalize(s State) State {
	norm, norm_hash := s, hash(s)

	foldIntoNorm := func(next State) {
		next_hash := hash(next)
		if next_hash < norm_hash {
			norm, norm_hash = next, next_hash
		}
	}

	foldIntoNorm(transpose(s))

	s_flip_v := flip_v(s)
	foldIntoNorm(s_flip_v)
	foldIntoNorm(transpose(s_flip_v))

	s_flip_h := flip_h(s)
	foldIntoNorm(s_flip_h)
	foldIntoNorm(transpose(s_flip_h))

	s_rot_180 := flip_h(s_flip_v)
	foldIntoNorm(s_rot_180)
	foldIntoNorm(transpose(s_rot_180))

	return norm
}

func hash(s State) int {
	accum := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v := s.Board[i][j] + 1
			accum = accum*3 + int(v)
		}
	}
	return accum
}

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
