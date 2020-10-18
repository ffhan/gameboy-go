package memory

import (
	"fmt"
	go_gb "go-gb"
)

const (
	MiB = 1 << 20
	KiB = 1 << 10
)

type noMBC struct {
	rom [ROMBankNEnd + 1]byte
	ram []byte
}

func (m *noMBC) ReadBytes(pointer, n uint16) []byte {
	if pointer <= ROMBankNEnd {
		return m.rom[pointer : pointer+n]
	} else if ExternalRAMStart <= pointer && pointer <= ExternalRAMEnd {
		if m.ram == nil {
			return make([]byte, n)
		}
		address := pointer - ExternalRAMStart
		return m.ram[address : address+n]
	}
	panic(fmt.Errorf("invalid memory access at address %X", pointer))
}

func (m *noMBC) Read(pointer uint16) byte {
	if pointer <= ROMBankNEnd {
		return m.rom[pointer]
	} else if ExternalRAMStart <= pointer && pointer <= ExternalRAMEnd {
		if m.ram == nil {
			return 0
		}
		address := pointer - ExternalRAMStart
		return m.ram[address]
	}
	panic(fmt.Errorf("invalid memory access at address %X", pointer))
}

func (m *noMBC) StoreBytes(pointer uint16, bytes []byte) {
	if m.ram == nil || pointer < ExternalRAMStart {
		return
	}
	address := int(pointer) - int(ExternalRAMStart)
	if address < 0 {
		return
	}
	copy(m.ram[address:address+len(bytes)], bytes)
}

func (m *noMBC) Store(pointer uint16, val byte) {
	if m.ram == nil || pointer < ExternalRAMStart {
		return
	}
	address := int(pointer) - int(ExternalRAMStart)
	if address < 0 {
		return
	}
	m.ram[address] = val
}

func (m *noMBC) LoadRom(bytes []byte) int {
	n := len(m.rom)
	if n > len(bytes) {
		n = len(bytes)
	}
	copy(m.rom[:], bytes[:n])
	return n
}

type mbc1 struct {
	romBank   *bank
	ramBank   *bank
	ramEnable bool

	selectedRomBank byte
	selectedRamBank byte
	ramBankingMode  bool
}

func NewMbc1(romBank *bank, ramBank *bank) *mbc1 {
	selectedRomBank := byte(1)
	if romBank.numOfParts == 1 {
		selectedRomBank = 0
	}
	return &mbc1{romBank: romBank, ramBank: ramBank, selectedRomBank: selectedRomBank}
}

func (m *mbc1) ReadBytes(pointer, n uint16) []byte {
	return go_gb.ReadBytes(m, pointer, n)
}

func (m *mbc1) Read(pointer uint16) byte {
	if pointer <= ROMBank0End {
		return m.romBank.Read(0, pointer)
	} else if pointer <= ROMBankNEnd {
		return m.romBank.Read(uint16(m.selectedRomBank), pointer-ROMBankNStart)
	} else if ExternalRAMStart <= pointer && pointer <= ExternalRAMEnd {
		if !m.ramEnable {
			return 0xFF
		}
		if m.ramBankingMode {
			return m.ramBank.Read(uint16(m.selectedRamBank), pointer-ExternalRAMStart)
		}
		return m.ramBank.Read(0, pointer-ExternalRAMStart)
	}
	panic(fmt.Errorf("invalid address %X", pointer))
}

func (m *mbc1) StoreBytes(pointer uint16, bytes []byte) {
	val := go_gb.FromBytes(bytes)
	if pointer <= 0x1FFF {
		m.ramEnable = val&0xFF == 0x0A
	} else if pointer <= ROMBank0End {
		val := byte(val & 0x1F)
		m.selectedRomBank = (m.selectedRomBank & 0xE0) | val
		switch m.selectedRomBank {
		case 0x00, 0x20, 0x40, 0x60:
			m.selectedRomBank += 1
		}
	} else if pointer <= 0x5FFF {
		val &= 0x03
		if m.ramBankingMode {
			m.selectedRamBank = byte(val)
		} else {
			m.selectedRomBank |= byte(val << 5)
			switch m.selectedRomBank {
			case 0x00, 0x20, 0x40, 0x60:
				m.selectedRomBank += 1
			}
		}
	} else if pointer <= ROMBankNEnd {
		m.ramBankingMode = val&0x01 == 0x01
	} else {
		if m.ramBank != nil && m.ramEnable {
			m.ramBank.StoreBytes(uint16(m.selectedRamBank), pointer-ExternalRAMStart, bytes)
		}
	}
}

func (m *mbc1) Store(pointer uint16, val byte) {
	m.StoreBytes(pointer, []byte{val})
}

func (m *mbc1) LoadRom(bytes []byte) int {
	return m.romBank.LoadRom(bytes)
}
