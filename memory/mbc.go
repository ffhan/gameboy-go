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

func getCartridge(memory []byte) Memory {
	cartridgeType := memory[CartridgeTypeAddr]
	switch cartridgeType {
	case 0x00:
		return &noMBC{}
	case 0x01:
		return NewMbc1(getRomBanks(memory), nil)
	case 0x02, 0x03: // todo: implement battery (save in files)
		return NewMbc1(getRomBanks(memory), getRamBanks(memory))
	case 0x05, 0x06:
		panic("implement MBC2")
	case 0x08, 0x09:
		return &noMBC{ram: make([]byte, ExternalRAMEnd-ExternalRAMStart+1)}
	}
	panic(fmt.Errorf("implement cartridge type %X", cartridgeType))
}

func getRomBanks(memory []byte) *bank {
	banks := memory[CartridgeROMSizeAddr]
	switch banks {
	case 0x00:
		return newBank(1, 32*KiB)
	case 0x01:
		return newBank(4, 16*KiB)
	case 0x02:
		return newBank(8, 16*KiB)
	case 0x03:
		return newBank(16, 16*KiB)
	case 0x04:
		return newBank(32, 16*KiB)
	case 0x05:
		return newBank(64, 16*KiB)
	case 0x06:
		return newBank(128, 16*KiB)
	case 0x07:
		return newBank(256, 16*KiB)
	case 0x08:
		return newBank(512, 16*KiB)
	case 0x52:
		return newBank(72, 16*KiB)
	case 0x53:
		return newBank(80, 16*KiB)
	case 0x54:
		return newBank(96, 16*KiB)
	}
	panic("invalid number of banks")
}

func getRamBanks(memory []byte) *bank {
	val := memory[CartridgeRAMSizeAddr]
	switch val {
	case 0x00:
		return nil
	case 0x01:
		return newBank(1, 2*KiB)
	case 0x02:
		return newBank(1, 8*KiB)
	case 0x03:
		return newBank(4, 8*KiB)
	case 0x04:
		return newBank(16, 8*KiB)
	case 0x05:
		return newBank(8, 8*KiB)
	}
	panic("invalid RAM size")
}

type bank struct {
	memory   []byte
	partSize uint16
}

func newBank(numOfParts, partSize uint16) *bank {
	return &bank{
		memory:   make([]byte, numOfParts*partSize),
		partSize: partSize,
	}
}

func (b *bank) address(bank, pointer uint16) uint16 {
	return b.partSize*bank + pointer
}

func (b *bank) ReadBytes(bank, pointer, n uint16) []byte {
	address := b.address(bank, pointer)
	return b.memory[address : address+n]
}

func (b *bank) Read(bank, pointer uint16) byte {
	address := b.address(bank, pointer)
	return b.memory[address]
}

func (b *bank) StoreBytes(bank, pointer uint16, bytes []byte) {
	address := b.address(bank, pointer)
	copy(b.memory[address:address+uint16(len(bytes))], bytes)
}

func (b *bank) Store(bank, pointer uint16, val byte) {
	b.memory[b.address(bank, pointer)] = val
}

type noMBC struct {
	rom [ROMBankNEnd + 1]byte
	ram []byte
}

func (n2 *noMBC) ReadBytes(pointer, n uint16) []byte {
	if pointer <= ROMBankNEnd {
		return n2.rom[pointer : pointer+n]
	} else if ExternalRAMStart <= pointer && pointer <= ExternalRAMEnd {
		if n2.ram == nil {
			return make([]byte, n)
		}
		address := pointer - ExternalRAMStart
		return n2.ram[address : address+n]
	}
	panic(fmt.Errorf("invalid memory access at address %X", pointer))
}

func (n2 *noMBC) Read(pointer uint16) byte {
	return n2.ReadBytes(pointer, 1)[0]
}

func (n2 *noMBC) StoreBytes(pointer uint16, bytes []byte) {
	if n2.ram == nil {
		return
	}
	address := pointer - ExternalRAMStart
	copy(n2.ram[address:address+uint16(len(bytes))], bytes)
}

func (n2 *noMBC) Store(pointer uint16, val byte) {
	n2.StoreBytes(pointer, []byte{val})
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
