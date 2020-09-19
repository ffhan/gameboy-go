package memory

import (
	go_gb "go-gb"
	"testing"
)

func TestMemoryBus_Read_Store(t *testing.T) {
	m := NewMMU()
	var b [0xFFFF + 1]byte
	for i := range b {
		b[i] = byte(i)
	}
	b[CartridgeTypeAddr] = 0x08    // ROM+RAM
	b[CartridgeROMSizeAddr] = 0x05 // 1MByte in 64 banks
	b[CartridgeRAMSizeAddr] = 0x03 // 32 KByte in 4 banks
	m.Init(b[:], go_gb.GB)

	for i := VRAMStart; i <= VRAMEnd; i++ {
		m.Store(i, byte(i))
	}

	for i := uint(ExternalRAMStart); i < 0xFFFF+1; i++ {
		m.Store(uint16(i), byte(i))
	}

	for i := 0; i <= 0xFFFF; i++ {
		val := m.Read(uint16(i))
		if uint16(i) == CartridgeTypeAddr {
			if val != 0x08 {
				t.Fatalf("memlocation %X: expected %X, got %X\n", i, 0x02, val)
			}
			continue
		} else if uint16(i) == CartridgeROMSizeAddr {
			if val != 0x05 {
				t.Fatalf("memlocation %X: expected %X, got %X\n", i, 0x05, val)
			}
			continue
		} else if uint16(i) == CartridgeRAMSizeAddr {
			if val != 0x03 {
				t.Fatalf("memlocation %X: expected %X, got %X\n", i, 0x03, val)
			}
			continue
		} else if ECHORAMStart <= uint16(i) && uint16(i) <= ECHORAMEnd {
			j := byte((uint16(i) - ECHORAMStart) + WRAMBank0Start)
			if val != j {
				t.Fatalf("memlocation %X: expected %X, got %X\n", i, byte(i), val)
			}
			continue
		}
		if val != byte(i) {
			t.Fatalf("memlocation %X: expected %X, got %X\n", i, byte(i), val)
		}
	}
}
