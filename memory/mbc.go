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
	if pointer <= ROMBank0End {
		return m.romBank.ReadBytes(0, pointer, n)
	} else if pointer <= ROMBankNEnd {
		return m.romBank.ReadBytes(uint16(m.selectedRomBank), pointer, n)
	} else if ExternalRAMStart <= pointer && pointer <= ExternalRAMEnd {
		if !m.ramEnable {
			return make([]byte, n)
		}
		if m.ramBankingMode {
			return m.ramBank.ReadBytes(uint16(m.selectedRamBank), pointer-ExternalRAMStart, n)
		}
		return m.ramBank.ReadBytes(0, pointer-ExternalRAMStart, n)
	}
	panic(fmt.Errorf("invalid address %X", pointer))
}

func (m *mbc1) Read(pointer uint16) byte {
	if pointer <= ROMBank0End {
		return m.romBank.Read(0, pointer)
	} else if pointer <= ROMBankNEnd {
		return m.romBank.Read(uint16(m.selectedRomBank), pointer)
	} else if ExternalRAMStart <= pointer && pointer <= ExternalRAMEnd {
		if !m.ramEnable {
			return 0
		}
		if m.ramBankingMode {
			return m.ramBank.Read(uint16(m.selectedRamBank), pointer-ExternalRAMStart)
		}
		return m.ramBank.Read(0, pointer-ExternalRAMStart)
	}
	panic(fmt.Errorf("invalid address %X", pointer))
}

func (m *mbc1) StoreBytes(pointer uint16, bytes []byte) {
	if pointer <= 0x1FFF {
		m.ramEnable = go_gb.FromBytes(bytes)&0xFF == 0x0A
	} else if pointer <= ROMBank0End {
		val := go_gb.FromBytes(bytes) & 0x1F
		if val == 0 {
			val = 1
		}
		m.selectedRomBank |= byte(val)
	} else if pointer <= 0x5FFF {
		if m.ramBankingMode {
			m.selectedRamBank = byte(go_gb.FromBytes(bytes) & 0x1F)
		} else {
			m.selectedRomBank |= byte((go_gb.FromBytes(bytes) & 0x3) << 5)
		}
	} else if pointer <= ROMBankNEnd {
		m.ramBankingMode = go_gb.FromBytes(bytes)&0x01 == 0x01
	} else {
		if m.ramBank != nil {
			m.ramBank.StoreBytes(uint16(m.selectedRamBank), pointer, bytes)
		}
	}
}

func (m *mbc1) Store(pointer uint16, val byte) {
	m.StoreBytes(pointer, []byte{val})
}

func (m *mbc1) LoadRom(bytes []byte) int {
	return m.romBank.LoadRom(bytes)
}
