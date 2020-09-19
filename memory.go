package go_gb

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

type Memory interface {
	ReadBytes(pointer, n uint16) []byte
	Read(pointer uint16) byte
	StoreBytes(pointer uint16, bytes []byte)
	Store(pointer uint16, val byte)
}

type RomLoader interface {
	LoadRom(rom []byte) int
}

type Cartridge interface {
	Memory
	RomLoader
}
