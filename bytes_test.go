package go_gb

import "testing"

func TestFromBytes(t *testing.T) { // GameBoy is little endian - if we read 12 34 in memory it's lsb msb
	v := []byte{0x12, 0x34}
	res := FromBytes(v)
	expected := uint16(0x3412)
	if res != expected {
		t.Errorf("expected %X, got %X\n", expected, res)
	}
}

func TestLsbMsbBytes(t *testing.T) {
	v := uint16(0xFE12)
	result := ToBytes(v, true)
	expected := []byte{0x12, 0xFE} // gameboy usually reads first lsb then msb
	for i := range result {
		if result[i] != expected[i] {
			t.Fatalf("expected %v, got %v\n", expected, result)
		}
	}
}

func TestReverse(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5}
	expected := []byte{5, 4, 3, 2, 1}
	Reverse(b)
	for i := range b {
		if expected[i] != b[i] {
			t.Errorf("expected %d, got %d\n", expected[i], b[i])
		}
	}
}
