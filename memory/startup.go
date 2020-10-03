package memory

import (
	"fmt"
	"go-gb"
)

func getCartridge(memory []byte) go_gb.Cartridge {
	cartridgeType := memory[go_gb.CartridgeTypeAddr]
	var mbc go_gb.Cartridge
	switch cartridgeType {
	case 0x00:
		mbc = &noMBC{}
	case 0x01:
		mbc = NewMbc1(getRomBanks(memory), nil)
	case 0x02, 0x03: // todo: implement battery (save in files)
		mbc = NewMbc1(getRomBanks(memory), getRamBanks(memory))
	case 0x05, 0x06:
		panic("implement MBC2")
	case 0x08, 0x09:
		mbc = &noMBC{ram: make([]byte, ExternalRAMEnd-ExternalRAMStart+1)}
	default:
		panic(fmt.Errorf("implement cartridge type %X", cartridgeType))
	}
	mbc.LoadRom(memory)
	return mbc
}

func getRomBanks(memory []byte) *bank {
	banks := go_gb.RomSize(memory[go_gb.CartridgeROMSizeAddr])
	size, num := banks.GetSize()
	return newBank(num, size)
}

func getRamBanks(memory []byte) *bank {
	val := go_gb.RamSize(memory[go_gb.CartridgeRAMSizeAddr])
	size, num := val.GetSize()
	return newBank(num, size)
}
