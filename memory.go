package go_gb

import (
	"fmt"
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

type mmap struct {
	start, end uint16
	memory     []byte
}

func (m *mmap) ReadBytes(pointer, n uint16) ([]byte, MC) {
	i := pointer - m.start
	return m.memory[i : i+n], MC(n)
}

func (m *mmap) Read(pointer uint16) (byte, MC) {
	i := pointer - m.start
	return m.memory[i], 1
}

func (m *mmap) StoreBytes(pointer uint16, bytes []byte) MC {
	i := pointer - m.start
	copy(m.memory[i:i+uint16(len(bytes))+1], bytes)
	return MC(len(bytes))
}

func (m *mmap) Store(pointer uint16, val byte) MC {
	m.memory[pointer-m.start] = val
	return 1
}

func newMmap(start uint16, end uint16, memory []byte) *mmap {
	return &mmap{start: start, end: end, memory: memory}
}

type MemoryBus struct {
	completeMem             [0xFFFF + 1]byte
	romBank0                *mmap
	romBankN                *mmap
	vram                    *mmap
	externalRam             *mmap
	wramBank0               *mmap
	wramBankN               *mmap
	echo                    *mmap
	oam                     *mmap
	unusable                *mmap
	io                      *mmap
	hram                    *mmap
	interruptEnableRegister *mmap
}

func NewMemoryBus() *MemoryBus {
	m := &MemoryBus{}
	m.init()
	return m
}

// returns true if x in [start, end], false otherwise
func inInterval(pointer, start, end uint16) bool {
	return start <= pointer && pointer <= end
}

func (m *MemoryBus) createMmapWithRedirection(start, end, redirectStart, redirectEnd uint16) *mmap {
	if (end - start) != (redirectEnd - redirectStart) {
		panic(fmt.Errorf("invalid redirection (%X, %X) and (%X, %X)", start, end, redirectStart, redirectEnd))
	}
	return newMmap(start, end, m.completeMem[redirectStart:int(redirectEnd)+1])
}

func (m *MemoryBus) createMmap(start, end uint16) *mmap {
	return m.createMmapWithRedirection(start, end, start, end)
}

func (m *MemoryBus) init() {
	m.romBank0 = m.createMmap(ROMBank0Start, ROMBank0End)
	m.romBankN = m.createMmap(ROMBankNStart, ROMBankNEnd)
	m.vram = m.createMmap(VRAMStart, VRAMEnd)
	m.externalRam = m.createMmap(ExternalRAMStart, ExternalRAMEnd)
	m.wramBank0 = m.createMmap(WRAMBank0Start, WRAMBank0End)
	m.wramBankN = m.createMmap(WRAMBankNStart, WRAMBankNEnd)
	m.echo = m.createMmapWithRedirection(ECHORAMStart, ECHORAMEnd, WRAMBank0Start, 0xDDFF)
	m.oam = m.createMmap(OAMStart, OAMEnd)
	m.unusable = m.createMmap(UnusableStart, UnusableEnd)
	m.io = m.createMmap(IOPortsStart, IOPortsEnd)
	m.hram = m.createMmap(HRAMStart, HRAMEnd)
	m.interruptEnableRegister = m.createMmap(InterruptEnableRegister, InterruptEnableRegister)
}

// takes a pointer and returns a whole portion of the memory responsible
func (m *MemoryBus) Route(pointer uint16) *mmap {
	if inInterval(pointer, ROMBank0Start, ROMBank0End) {
		return m.romBank0
	} else if inInterval(pointer, ROMBankNStart, ROMBankNEnd) {
		return m.romBankN
	} else if inInterval(pointer, VRAMStart, VRAMEnd) {
		return m.vram
	} else if inInterval(pointer, ExternalRAMStart, ExternalRAMEnd) {
		return m.externalRam
	} else if inInterval(pointer, WRAMBank0Start, WRAMBank0End) {
		return m.wramBank0
	} else if inInterval(pointer, WRAMBankNStart, WRAMBankNEnd) {
		return m.wramBankN
	} else if inInterval(pointer, ECHORAMStart, ECHORAMEnd) {
		return m.echo
	} else if inInterval(pointer, OAMStart, OAMEnd) {
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

func (m *MemoryBus) ReadBytes(pointer, n uint16) ([]byte, MC) {
	return m.Route(pointer).ReadBytes(pointer, n)
}

func (m *MemoryBus) Read(pointer uint16) (byte, MC) {
	return m.Route(pointer).Read(pointer)
}

func (m *MemoryBus) StoreBytes(pointer uint16, bytes []byte) MC {
	m.Route(pointer).StoreBytes(pointer, bytes)
	return MC(len(bytes))
}

func (m *MemoryBus) Store(pointer uint16, val byte) MC {
	m.Route(pointer).Store(pointer, val)
	return MC(1)
}
