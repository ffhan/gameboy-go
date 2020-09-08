package cpu

import "testing"

func TestOptable(t *testing.T) {
	if len(optable) != 256 {
		t.Fatal("invalid table size")
	}
	if len(cbOptable) != 256 {
		t.Fatal("invalid CB table size")
	}
}
