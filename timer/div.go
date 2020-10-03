package timer

import go_gb "go-gb"

const (
	divFreq = 16384 / 4 // Hz (T cycles) to M cycles
)

type divTimer struct {
	io            go_gb.Memory
	currentCycles go_gb.MC
}

func NewDivTimer(io go_gb.Memory) *divTimer {
	return &divTimer{io: io}
}

func (t *divTimer) Step(cycles go_gb.MC) {
	t.currentCycles += cycles

	if t.currentCycles >= divFreq {
		go_gb.Update(t.io, go_gb.DIV, func(b byte) byte {
			return b + 1
		})
		t.currentCycles -= divFreq
	}
}
