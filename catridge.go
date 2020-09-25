package go_gb

import (
	"fmt"
	"strconv"
)

const (
	MiB = 1 << 20
	KiB = 1 << 10
)

type CartridgeType byte

func (c CartridgeType) String() string {
	switch c {
	case MbcROMOnly:
		return "ROM ONLY"
	case MbcMBC1:
		return "MBC1"
	case MbcMBC1RAM:
		return "MBC1+RAM"
	case MbcMBC1BATTERY:
		return "MBC1+RAM+BATTERY"
	case MbcMBC2:
		return "MBC2"
	case MbcMBC2BATTERY:
		return "MBC2+BATTERY"
	case MbcROMRAM:
		return "ROM+RAM"
	case MbcROMRAMBATTERY:
		return "ROM+RAM+BATTERY"
	case MbcMMM01:
		return "MMM01"
	case MbcMMM01RAM:
		return "MMM01+RAM"
	case MbcMMM01RAMBATTERY:
		return "MMM01+RAM+BATTERY"
	case MbcMBC3TIMERBATTERY:
		return "MBC3+TIMER+BATTERY"
	case MbcMBC3TIMERRAMBATTERY:
		return "MBC3+TIMER+RAM+BATTERY"
	case MbcMBC3:
		return "MBC3"
	case MbcMBC3RAM:
		return "MBC3+RAM"
	case MbcMBC3RAMBATTERY:
		return "MBC3+RAM+BATTERY"
	case MbcMBC5:
		return "MBC5"
	case MbcMBC5RAM:
		return "MBC5+RAM"
	case MbcMBC5RAMBATTERY:
		return "MBC5+RAM+BATTERY"
	case MbcMBC5RUMBLE:
		return "MBC5+RUMBLE"
	case MbcMBC5RUMBLERAM:
		return "MBC5+RUMBLE+RAM"
	case MbcMBC5RUMBLERAMBATTERY:
		return "MBC5+RUMBLE+RAM+BATTERY"
	case MbcMBC6:
		return "MBC6"
	case MbcMBC7SENSORRUMBLERAMBATTERY:
		return "MBC7+SENSOR+RUMBLE+RAM+BATTERY"
	case MbcPOCKETCAMERA:
		return "POCKET CAMERA"
	case MbcBANDAITAMA5:
		return "BANDAI TAMA5"
	case MbcHuC3:
		return "HuC3"
	case MbcHuC1RAMBATTERY:
		return "HuC1+RAM+BATTERY"
	}
	panic("invalid cartridge type " + strconv.Itoa(int(c)))
}

const (
	CartridgeTypeAddr    uint16 = 0x147
	CartridgeROMSizeAddr uint16 = 0x148
	CartridgeRAMSizeAddr uint16 = 0x149

	MbcROMOnly                    CartridgeType = 0x00
	MbcMBC1                       CartridgeType = 0x01
	MbcMBC1RAM                    CartridgeType = 0x02
	MbcMBC1BATTERY                CartridgeType = 0x03
	MbcMBC2                       CartridgeType = 0x05
	MbcMBC2BATTERY                CartridgeType = 0x06
	MbcROMRAM                     CartridgeType = 0x08
	MbcROMRAMBATTERY              CartridgeType = 0x09
	MbcMMM01                      CartridgeType = 0x0B
	MbcMMM01RAM                   CartridgeType = 0x0C
	MbcMMM01RAMBATTERY            CartridgeType = 0x0D
	MbcMBC3TIMERBATTERY           CartridgeType = 0x0F
	MbcMBC3TIMERRAMBATTERY        CartridgeType = 0x10
	MbcMBC3                       CartridgeType = 0x11
	MbcMBC3RAM                    CartridgeType = 0x12
	MbcMBC3RAMBATTERY             CartridgeType = 0x13
	MbcMBC5                       CartridgeType = 0x19
	MbcMBC5RAM                    CartridgeType = 0x1A
	MbcMBC5RAMBATTERY             CartridgeType = 0x1B
	MbcMBC5RUMBLE                 CartridgeType = 0x1C
	MbcMBC5RUMBLERAM              CartridgeType = 0x1D
	MbcMBC5RUMBLERAMBATTERY       CartridgeType = 0x1E
	MbcMBC6                       CartridgeType = 0x20
	MbcMBC7SENSORRUMBLERAMBATTERY CartridgeType = 0x22
	MbcPOCKETCAMERA               CartridgeType = 0xFC
	MbcBANDAITAMA5                CartridgeType = 0xFD
	MbcHuC3                       CartridgeType = 0xFE
	MbcHuC1RAMBATTERY             CartridgeType = 0xFF
)

type RomSize byte

func (r RomSize) GetSize() (uint, uint) {
	switch r {
	case 0x00:
		return 32 * KiB, 1
	case 0x01:
		return 64 * KiB, 4
	case 0x02:
		return 128 * KiB, 8
	case 0x03:
		return 256 * KiB, 16
	case 0x04:
		return 512 * KiB, 32
	case 0x05:
		return 1 * MiB, 64
	case 0x06:
		return 2 * MiB, 128
	case 0x07:
		return 4 * MiB, 256
	case 0x08:
		return 8 * MiB, 512
	case 0x52:
		return 1100 * KiB, 72
	case 0x53:
		return 1200 * KiB, 80
	case 0x54:
		return 1500 * KiB, 96
	}
	panic("invalid ROM size " + strconv.Itoa(int(r)))
}

func (r RomSize) String() string {
	size, banks := r.GetSize()
	var suffix string
	if banks > 1 || banks == 0 {
		suffix = "s"
	}
	if size >= MiB {
		return fmt.Sprintf("%d MiB in %d bank%s", size/MiB, banks, suffix)
	} else {
		return fmt.Sprintf("%d KiB in %d bank%s", size/KiB, banks, suffix)
	}
}

type RamSize byte

func (r RamSize) GetSize() (uint, uint) {
	switch r {
	case 0x00:
		return 0, 0 // None
	case 0x01:
		return 2 * KiB, 1 // 2 KBytes
	case 0x02:
		return 8 * KiB, 1 // 8 KBytes
	case 0x03:
		return 32 * KiB, 4 // 32 KBytes (4 banks of 8KBytes each)
	case 0x04:
		return 128 * KiB, 16 // 128 KBytes (16 banks of 8KBytes each)
	case 0x05:
		return 64 * KiB, 8 // 64 KBytes (8 banks of 8KBytes each)
	}
	panic("invalid RAM size " + strconv.Itoa(int(r)))
}

func (r RamSize) String() string {
	size, banks := r.GetSize()
	var suffix string
	if banks > 1 || banks == 0 {
		suffix = "s"
	}
	return fmt.Sprintf("%d KiB in %d bank%s", size/KiB, banks, suffix)
}
