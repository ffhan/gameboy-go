package memory

import (
	"fmt"
	"go-gb"
)

const (
	ROMBank0Start           uint16 = 0x0000
	ROMBank0End             uint16 = 0x3FFF
	ROMBankNStart           uint16 = 0x4000
	ROMBankNEnd             uint16 = 0x7FFF
	VRAMStart               uint16 = 0x8000
	VRAMEnd                 uint16 = 0x9FFF
	ExternalRAMStart        uint16 = 0xA000
	ExternalRAMEnd          uint16 = 0xBFFF
	WRAMBank0Start          uint16 = 0xC000
	WRAMBank0End            uint16 = 0xCFFF
	WRAMBankNStart          uint16 = 0xD000
	WRAMBankNEnd            uint16 = 0xDFFF
	ECHORAMStart            uint16 = 0xE000
	ECHORAMEnd              uint16 = 0xFDFF
	OAMStart                uint16 = 0xFE00
	OAMEnd                  uint16 = 0xFE9F
	UnusableStart           uint16 = 0xFEA0
	UnusableEnd             uint16 = 0xFEFF
	IOPortsStart            uint16 = 0xFF00
	IOPortsEnd              uint16 = 0xFF7F
	HRAMStart               uint16 = 0xFF80
	HRAMEnd                 uint16 = 0xFFFE
	InterruptEnableRegister uint16 = 0xFFFF
)

type mmu struct {
	internalMemory          [0xFFFF + 1]byte
	cartridge               go_gb.Cartridge
	vram                    go_gb.Memory
	wram                    go_gb.Memory
	echo                    go_gb.Memory
	oam                     go_gb.Memory
	unusable                go_gb.Memory
	io                      go_gb.Memory
	hram                    go_gb.Memory
	interruptEnableRegister go_gb.Memory

	locked *lockedMemory
}

func NewMMU() *mmu {
	m := &mmu{}
	return m
}

// returns true if x in [start, end], false otherwise
func inInterval(pointer, start, end uint16) bool {
	return start <= pointer && pointer <= end
}

func (m *mmu) createMmapWithRedirection(start, end, redirectStart, redirectEnd uint16) *mmap {
	if (end - start) != (redirectEnd - redirectStart) {
		panic(fmt.Errorf("invalid redirection (%X, %X) and (%X, %X)", start, end, redirectStart, redirectEnd))
	}
	return newMmap(start, end, m.internalMemory[redirectStart:int(redirectEnd)+1])
}

func (m *mmu) createMmap(start, end uint16) *mmap {
	return m.createMmapWithRedirection(start, end, start, end)
}

func (m *mmu) Init(rom []byte, gbType go_gb.GameboyType) {
	var wramMemory go_gb.Memory
	if gbType == go_gb.CGB {
		wramMemory = &wram{bank: newBank(8, 1<<12), selectedBank: 1}
	} else {
		wramMemory = &wram{bank: newBank(2, 1<<12), selectedBank: 1}
	}
	m.cartridge = getCartridge(rom)
	m.vram = m.createMmap(VRAMStart, VRAMEnd)
	m.wram = wramMemory
	m.echo = m.createMmapWithRedirection(ECHORAMStart, ECHORAMEnd, WRAMBank0Start, 0xDDFF)
	m.oam = m.createMmap(OAMStart, OAMEnd)
	m.unusable = m.createMmap(UnusableStart, UnusableEnd)
	m.io = m.createMmap(IOPortsStart, IOPortsEnd)
	m.hram = m.createMmap(HRAMStart, HRAMEnd)
	m.interruptEnableRegister = m.createMmap(InterruptEnableRegister, InterruptEnableRegister)

	m.locked = &lockedMemory{}
}

// takes a pointer and returns a whole portion of the memory responsible
func (m *mmu) Route(pointer uint16) go_gb.Memory {
	if inInterval(pointer, ROMBank0Start, ROMBankNEnd) {
		return m.cartridge
	} else if inInterval(pointer, VRAMStart, VRAMEnd) {
		locked := m.Read(go_gb.LCDStatusRegister)&0x3 == 3
		if !locked {
			return m.locked
		}
		return m.vram
	} else if inInterval(pointer, ExternalRAMStart, ExternalRAMEnd) {
		return m.cartridge
	} else if inInterval(pointer, WRAMBank0Start, WRAMBankNEnd) {
		return m.wram
	} else if inInterval(pointer, ECHORAMStart, ECHORAMEnd) {
		return m.echo
	} else if inInterval(pointer, OAMStart, OAMEnd) {
		locked := m.Read(go_gb.LCDStatusRegister)&0x3 > 1
		if locked {
			return m.locked
		}
		return m.oam
	} else if inInterval(pointer, UnusableStart, UnusableEnd) {
		return m.unusable
	} else if inInterval(pointer, IOPortsStart, IOPortsEnd) {
		return m.io
	} else if inInterval(pointer, HRAMStart, HRAMEnd) {
		return m.hram
	} else if pointer == InterruptEnableRegister {
		return m.interruptEnableRegister
	}
	panic(fmt.Errorf("invalid pointer %X", pointer))
}

func (m *mmu) ReadBytes(pointer, n uint16) []byte {
	return m.Route(pointer).ReadBytes(pointer, n)
}

func (m *mmu) Read(pointer uint16) byte {
	return m.Route(pointer).Read(pointer)
}

func (m *mmu) StoreBytes(pointer uint16, bytes []byte) {
	m.Route(pointer).StoreBytes(pointer, bytes)
}

func (m *mmu) Store(pointer uint16, val byte) {
	m.Route(pointer).Store(pointer, val)
}

func (m *mmu) LoadRom(rom []byte) int {
	return m.cartridge.LoadRom(rom)
}
