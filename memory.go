package go_gb

import "io"

const (
	MemEntrypoint        uint16 = 0x0100
	MemNintendoLogoStart uint16 = 0x0104
	MemNintendoLogoEnd   uint16 = 0x0133
	MemTitleStart        uint16 = 0x0134
	MemTitleEnd          uint16 = 0x0143
	MemCGBFlag           uint16 = 0x0143
	MemRomSize           uint16 = 0x0148
	MemRamSize           uint16 = 0x0149
)

type Reader interface {
	Read(pointer uint16) byte
}

type Writer interface {
	Store(pointer uint16, val byte)
}

type Memory interface {
	ReadBytes(pointer, n uint16) []byte
	Reader
	StoreBytes(pointer uint16, bytes []byte)
	Writer
}

type Dumper interface {
	Dump(writer io.Writer)
}

func Update(memory Memory, address uint16, updateFunc func(b byte) byte) {
	val := memory.Read(address)
	val = updateFunc(val)
	memory.Store(address, val)
}

type RomLoader interface {
	LoadRom(rom []byte) int
}

type Cartridge interface {
	Memory
	RomLoader
}

type MemoryBus interface {
	Memory
	VRAM() Memory
	HRAM() Memory
	OAM() Memory
	IO() Memory
	InterruptEnableRegister() Memory
	Booted() bool
	DMAInProgress() bool
	SetDMAInProgress(val bool)
}
