package memory

import "fmt"

type wram struct {
	bank         *bank
	selectedBank int
}

func (w *wram) Memory() []byte {
	return w.bank.memory
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
	if WRAMBank0Start <= pointer && pointer <= WRAMBank0End {
		return w.bank.Read(0, pointer-WRAMBank0Start)
	} else if WRAMBankNStart <= pointer && pointer <= WRAMBankNEnd {
		return w.bank.Read(uint16(w.selectedBank), pointer-WRAMBankNStart)
	}
	panic(fmt.Errorf("invalid address %X", pointer))
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
	if WRAMBank0Start <= pointer && pointer <= WRAMBank0End {
		w.bank.Store(0, pointer-WRAMBank0Start, val)
		return
	} else if WRAMBankNStart <= pointer && pointer <= WRAMBankNEnd {
		w.bank.Store(uint16(w.selectedBank), pointer-WRAMBankNStart, val)
		return
	}
	panic(fmt.Errorf("invalid address %X", pointer))
}
