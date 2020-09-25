package go_gb

import (
	"fmt"
)

type Cpu interface {
	Step() MC
	PC() uint16
	IME() bool
}

// picture processing unit
type PPU interface {
	Step(mc MC)
	IsVBlank() bool
	Enabled() bool
	Mode() byte
	CurrentLine() int
}

// sound processing unit
type SPU interface {
}

type CGBFlag byte

func (C CGBFlag) String() string {
	switch C {
	case CGBSupport:
		return "CGB & DMG supported"
	case OnlyCGB:
		return "CGB supported"
	default:
		return "DMG supported"
	}
}

type SGBFlag byte

func (S SGBFlag) String() string {
	switch S {
	case SGBSupport:
		return "SGB support"
	default:
		return "No SGB support"
	}
}

const (
	NoSGB      SGBFlag = 0x00
	SGBSupport SGBFlag = 0x03
)

const (
	NonCGB     CGBFlag = 0x00
	CGBSupport CGBFlag = 0x80
	OnlyCGB    CGBFlag = 0xC0
)

type Game struct {
	Rom           []byte
	Title         string
	CGBFlag       CGBFlag
	SGBFlag       SGBFlag
	CartridgeType CartridgeType
	RomSize       RomSize
	RamSize       RamSize
	NonJapanese   bool
}

func (g *Game) String() string {
	destination := "Non-Japanese"
	if !g.NonJapanese {
		destination = "Japanese"
	}
	return fmt.Sprintf("Title: %s\n%s\n%s\nCartridge type: %s\nROM size: %s\nRAM size: %s\nDestination: %s",
		g.Title, g.CGBFlag.String(), g.SGBFlag.String(),
		g.CartridgeType.String(), g.RomSize.String(), g.RamSize.String(), destination)
}

func NewGame(rom []byte) *Game {
	return &Game{
		Rom:           rom,
		Title:         string(rom[0x134:0x144]),
		CGBFlag:       CGBFlag(rom[0x143]),
		SGBFlag:       SGBFlag(rom[0x146]),
		CartridgeType: CartridgeType(rom[CartridgeTypeAddr]),
		RomSize:       RomSize(rom[CartridgeROMSizeAddr]),
		RamSize:       RamSize(rom[CartridgeRAMSizeAddr]),
		NonJapanese:   rom[0x14A] != 0,
	}
}

type GameBoy struct {
	cpu Cpu
	mmu Memory
	ppu PPU
	spu SPU
}

func (g *GameBoy) Run() {
	const (
		cpuFreq = 4_194_304 // Hz
		ppuFreq = 59.73     // Hz
	)
	for {
		mc := g.cpu.Step()
		g.ppu.Step(mc)
	}
}
