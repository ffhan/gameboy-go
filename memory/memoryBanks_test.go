package memory

import "testing"

func TestBank_Read(t *testing.T) {
	b := newBank(3, 16)
	for i := 0; i < 3*16; i++ {
		b.memory[i] = byte(i)
	}
	for bankId := uint16(0); bankId < 3; bankId++ {
		for i := uint16(0); i < 16; i++ {
			expected := byte(bankId*16 + i)
			result := b.Read(bankId, i)
			if result != expected {
				t.Errorf("expected %d, got %d\n", expected, result)
			}
		}
	}
}
