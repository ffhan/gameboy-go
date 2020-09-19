package go_gb

type Cpu interface {
	Step() MC
}

// picture processing unit
type PPU interface {
	Paint()
}

// sound processing unit
type SPU interface {
}

type GameBoy struct {
	cpu Cpu
	mmu Memory
	ppu PPU
	spu SPU
}
