package go_gb

import "testing"

func TestHelper(t *testing.T) {
	o := Unprefixed[0xC4]
	expected1 := "CALL NZ, a16"
	if o.String() != expected1 {
		t.Errorf("expected '%s', got '%s'\n", expected1, o.String())
	}
	o = Unprefixed[0x02]
	expected2 := "LD (BC), A"
	if o.String() != expected2 {
		t.Errorf("expected '%s', got '%s'\n", expected2, o.String())
	}
}
