package memory

import (
	"fmt"
	go_gb "go-gb"
)

type mbcType int

const (
	CartridgeTypeAddr    uint16 = 0x147
	CartridgeROMSizeAddr uint16 = 0x148
	CartridgeRAMSizeAddr uint16 = 0x149

	MiB = 1 << 20
	KiB = 1 << 10

	MbcROMOnly                    mbcType = 0x00
	MbcMBC1                       mbcType = 0x01
	MbcMBC1RAM                    mbcType = 0x02
	MbcMBC1BATTERY                mbcType = 0x03
	MbcMBC2                       mbcType = 0x05
	MbcMBC2BATTERY                mbcType = 0x06
	MbcROMRAM                     mbcType = 0x08
	MbcROMRAMBATTERY              mbcType = 0x09
	MbcMMM01                      mbcType = 0x0B
	MbcMMM01RAM                   mbcType = 0x0C
	MbcMMM01RAMBATTERY            mbcType = 0x0D
	MbcMBC3TIMERBATTERY           mbcType = 0x0F
	MbcMBC3TIMERRAMBATTERY        mbcType = 0x10
	MbcMBC3                       mbcType = 0x11
	MbcMBC3RAM                    mbcType = 0x12
	MbcMBC3RAMBATTERY             mbcType = 0x13
	MbcMBC5                       mbcType = 0x19
	MbcMBC5RAM                    mbcType = 0x1A
	MbcMBC5RAMBATTERY             mbcType = 0x1B
	MbcMBC5RUMBLE                 mbcType = 0x1C
	MbcMBC5RUMBLERAM              mbcType = 0x1D
	MbcMBC5RUMBLERAMBATTERY       mbcType = 0x1E
	MbcMBC6                       mbcType = 0x20
	MbcMBC7SENSORRUMBLERAMBATTERY mbcType = 0x22
	MbcPOCKETCAMERA               mbcType = 0xFC
	MbcBANDAITAMA5                mbcType = 0xFD
	MbcHuC3                       mbcType = 0xFE
	MbcHuC1RAMBATTERY             mbcType = 0xFF
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
	return m.ReadBytes(pointer, 1)[0]
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
	m.StoreBytes(pointer, []byte{val})
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
	return &mbc1{romBank: romBank, ramBank: ramBank, selectedRomBank: 1}
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
	return m.ReadBytes(pointer, 1)[0]
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
	}
}

func (m *mbc1) Store(pointer uint16, val byte) {
	m.StoreBytes(pointer, []byte{val})
}

func (m *mbc1) LoadRom(bytes []byte) int {
	return m.romBank.LoadRom(bytes)
}
