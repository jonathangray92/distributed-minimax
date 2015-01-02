package naiveminimax

import "testing"

func TestLess(t *testing.T) {
	if less(1, 2) != true {
		t.Fail()
	}
}

func TestGreater(t *testing.T) {
	if greater(1, 2) != false {
		t.Fail()
	}
}
