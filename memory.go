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

func ReadBytes(reader Reader, pointer uint16, n uint16) []byte {
	result := make([]byte, n)
	for i := uint16(0); i < n; i++ {
		result[i] = reader.Read(pointer + i)
	}
	return result
}

func WriteBytes(writer Writer, pointer uint16, data []byte) {
	for i, b := range data {
		i := uint16(i)
		writer.Store(pointer+i, b)
	}
}

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
