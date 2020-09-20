package memory

import (
	"go-gb"
)

func getCartridge(memory []byte) go_gb.Cartridge {
	cartridgeType := memory[CartridgeTypeAddr]
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
	banks := memory[CartridgeROMSizeAddr]
	var bank *bank
	switch banks {
	case 0x00:
		bank = newBank(1, 32*KiB)
	case 0x01:
		bank = newBank(4, 16*KiB)
	case 0x02:
		bank = newBank(8, 16*KiB)
	case 0x03:
		bank = newBank(16, 16*KiB)
	case 0x04:
		bank = newBank(32, 16*KiB)
	case 0x05:
		bank = newBank(64, 16*KiB)
	case 0x06:
		bank = newBank(128, 16*KiB)
	case 0x07:
		bank = newBank(256, 16*KiB)
	case 0x08:
		bank = newBank(512, 16*KiB)
	case 0x52:
		bank = newBank(72, 16*KiB)
	case 0x53:
		bank = newBank(80, 16*KiB)
	case 0x54:
		bank = newBank(96, 16*KiB)
	}
	if bank == nil {
		panic("invalid number of banks")
	}
	bank.LoadRom(memory)
	return bank
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
