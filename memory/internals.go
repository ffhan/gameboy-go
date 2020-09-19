package memory

import "fmt"

type mmap struct {
	start, end   uint16
	memory       []byte
	rLock, wLock bool
}

func (m *mmap) ReadBytes(pointer, n uint16) []byte {
	i := pointer - m.start
	return m.memory[i : i+n]
}

func (m *mmap) Read(pointer uint16) byte {
	i := pointer - m.start
	return m.memory[i]
}

func (m *mmap) StoreBytes(pointer uint16, bytes []byte) {
	i := pointer - m.start
	copy(m.memory[i:i+uint16(len(bytes))], bytes)
}

func (m *mmap) Store(pointer uint16, val byte) {
	m.memory[pointer-m.start] = val
}

func (m *mmap) LoadRom(bytes []byte) int {
	n := len(m.memory)
	copy(m.memory, bytes[:n])
	return n
}

func newMmap(start uint16, end uint16, memory []byte) *mmap {
	return &mmap{start: start, end: end, memory: memory}
}

type wram struct {
	bank         *bank
	selectedBank int
}

func (w *wram) ReadBytes(pointer, n uint16) []byte {
	if WRAMBank0Start <= pointer && pointer <= WRAMBank0End {
		return w.bank.ReadBytes(0, pointer-WRAMBank0Start, n)
	} else if WRAMBankNStart <= pointer && pointer <= WRAMBankNEnd {
		return w.bank.ReadBytes(uint16(w.selectedBank), pointer-WRAMBankNStart, n)
	}
	panic(fmt.Errorf("invalid address %X", pointer))
}

func (w *wram) Read(pointer uint16) byte {
	return w.ReadBytes(pointer, 1)[0]
}

func (w *wram) StoreBytes(pointer uint16, bytes []byte) {
	if WRAMBank0Start <= pointer && pointer <= WRAMBank0End {
		w.bank.StoreBytes(0, pointer-WRAMBank0Start, bytes)
		return
	} else if WRAMBankNStart <= pointer && pointer <= WRAMBankNEnd {
		w.bank.StoreBytes(uint16(w.selectedBank), pointer-WRAMBankNStart, bytes)
		return
	}
	panic(fmt.Errorf("invalid address %X", pointer))
}

func (w *wram) Store(pointer uint16, val byte) {
	w.StoreBytes(pointer, []byte{val})
}
