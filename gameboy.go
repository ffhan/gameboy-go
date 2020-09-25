package go_gb

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
