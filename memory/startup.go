package memory

import (
	"go-gb"
)

func getCartridge(memory []byte) go_gb.Cartridge {
	cartridgeType := memory[go_gb.CartridgeTypeAddr]
	switch cartridgeType {
	case 0x00:
		n := &noMBC{}
		n.LoadRom(memory)
		return n
	case 0x01:
		return NewMbc1(getRomBanks(memory), nil)
	case 0x02, 0x03: // todo: implement battery (save in files)
		return NewMbc1(getRomBanks(memory), getRamBanks(memory))
	case 0x05, 0x06:
		panic("implement MBC2")
	case 0x08, 0x09:
		mbc := &noMBC{ram: make([]byte, ExternalRAMEnd-ExternalRAMStart+1)}
		mbc.LoadRom(memory)
		return mbc
	}
	//panic(fmt.Errorf("implement cartridge type %X", cartridgeType))
	c := &noMBC{}
	c.LoadRom(memory)
	return c
}

func getRomBanks(memory []byte) *bank {
	banks := go_gb.RomSize(memory[go_gb.CartridgeROMSizeAddr])
	size, num := banks.GetSize()
	bank := newBank(uint16(num), uint16(size))
	bank.LoadRom(memory)
	return bank
}

func getRamBanks(memory []byte) *bank {
	val := go_gb.RamSize(memory[go_gb.CartridgeRAMSizeAddr])
	size, num := val.GetSize()
	return newBank(uint16(num), uint16(size))
}
