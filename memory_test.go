package go_gb

import "testing"

func TestMemoryBus_Read(t *testing.T) {
	m := NewMemoryBus()
	var b [0xFFFF + 1]byte
	for i := range b {
		b[i] = byte(i)
	}
	copy(m.completeMem[:], b[:])
	for i := 0; i <= 0xFFFF; i++ {
		var mc MC
		val := m.Read(uint16(i), &mc)
		if mc != 1 {
			t.Error("MC should be 1")
		}
		if ECHORAMStart <= uint16(i) && uint16(i) <= ECHORAMEnd {
			j := byte((uint16(i) - ECHORAMStart) + WRAMBank0Start)
			if val != j {
				t.Errorf("memlocation %d: expected %X, got %X\n", i, byte(i), val)
			}
			continue
		}
		if val != byte(i) {
			t.Errorf("memlocation %d: expected %X, got %X\n", i, byte(i), val)
		}
	}
}

func TestMemoryBus_Store(t *testing.T) {
	m := NewMemoryBus()
	for i := 0; i <= 0xFFFF; i++ {
		if ECHORAMStart <= uint16(i) && uint16(i) <= ECHORAMEnd {
			continue
		}
		var mc MC
		m.Store(uint16(i), byte(i), &mc)
		if mc != 1 {
			t.Error("MC should be 1")
		}
	}
	for i := 0; i <= 0xFFFF; i++ {
		val := m.completeMem[i]
		if ECHORAMStart <= uint16(i) && uint16(i) <= ECHORAMEnd {
			if val != 0 { // real echo location should be empty, not set
				t.Errorf("memlocation %d: expected %X, got %X\n", i, 0, val)
			}
			continue
		}
		if val != byte(i) {
			t.Errorf("memlocation %d: expected %X, got %X\n", i, byte(i), val)
		}
	}
}
