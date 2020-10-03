package timer

import go_gb "go-gb"

type timer struct {
	io            go_gb.Memory
	enabled       bool
	frequency     go_gb.MC
	currentCycles go_gb.MC
}

func NewTimer(io go_gb.Memory) *timer {
	return &timer{io: io}
}

func mapFrequency(clockSelect byte) go_gb.MC {
	switch clockSelect & 0x3 {
	case 0b00:
		return 4096 / 4
	case 0b01:
		return 262144 / 4
	case 0b10:
		return 65536 / 4
	case 0b11:
		return 16384 / 4
	}
	panic("invalid frequency")
}

func (t *timer) loadConfig() {
	timerControl := t.io.Read(go_gb.TAC)
	t.enabled = go_gb.Bit(timerControl, 2)
	t.frequency = mapFrequency(timerControl)
}

func (t *timer) Step(cycles go_gb.MC) {
	t.loadConfig()
	if !t.enabled {
		return
	}
	t.currentCycles += cycles

	if t.currentCycles >= t.frequency {
		go_gb.Update(t.io, go_gb.TIMA, func(b byte) byte {
			b += 1
			if b == 0 {
				go_gb.Update(t.io, go_gb.IF, func(b byte) byte {
					go_gb.Set(&b, int(go_gb.BitTimer), true)
					return b
				})
				return t.io.Read(go_gb.TMA)
			}
			return b
		})
		t.currentCycles -= t.frequency
	}
}
