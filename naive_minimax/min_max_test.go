package naive_minimax

import "testing"

func TestMin(t *testing.T) {
	if min(1, 2) != 1 {
		t.Fail()
	}
}

func TestMax(t *testing.T) {
	if max(1, 2) != 2 {
		t.Fail()
	}
}
