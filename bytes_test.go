package go_gb

import "testing"

func TestUnifyBytes(t *testing.T) {
	v := []byte{0x12, 0x34}
	res := UnifyBytes(v)
	expected := uint16(0x1234)
	if res != expected {
		t.Errorf("expected %X, got %X\n", expected, res)
	}
}

func TestSeparateUint16(t *testing.T) {
	v := uint16(0xFE12)
	result := SeparateUint16(v)
	expected := []byte{0xFE, 0x12}
	for i := range result {
		if result[i] != expected[i] {
			t.Fatalf("expected %v, got %v\n", expected, result)
		}
	}
}
