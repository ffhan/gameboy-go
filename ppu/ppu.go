package ppu

import (
	go_gb "go-gb"
)

const (
	frameBufASize = 160 * 144 * 4
)

type ppu struct {
	memory go_gb.Memory // used for usual memory access
	vram   go_gb.Memory // for skipping locks
	oam    go_gb.Memory // for skipping locks

	currentLine int
	currentMode byte
	modeClock   go_gb.MC
}

func (p *ppu) getTileMapAddr() uint16 {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 3) {
		return 0x9C00
	}
	return 0x9800
}

func (p *ppu) getScXY() (byte, byte) {
	scx := p.memory.Read(go_gb.LCDSCX)
	scy := p.memory.Read(go_gb.LCDSCY)
	return scx, scy
}

func (p *ppu) getLy() byte {
	return p.memory.Read(go_gb.LCDLY)
}

func (p *ppu) updateLy() {
	p.memory.Store(go_gb.LCDLY, byte(p.currentLine))
}

func (p *ppu) getWindowPosition() (byte, byte) {
	wx := p.memory.Read(go_gb.LCDWX)
	wy := p.memory.Read(go_gb.LCDWY)
	return wx, wy
}

func (p *ppu) getTileDataAddr() uint16 {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 4) {
		return 0x8000
	}
	return 0x8800
}

func (p *ppu) setMode(mode byte) {
	mode &= 0x3
	stat := p.memory.Read(go_gb.LCDStatusRegister) | mode
	p.memory.Store(go_gb.LCDStatusRegister, stat)
	p.modeClock = 0
	p.currentMode = mode
}

func (p *ppu) compareLyLyc() {
	val := p.memory.Read(go_gb.LCDLYC) == p.getLy()
	stat := p.memory.Read(go_gb.LCDStatusRegister)
	go_gb.Set(&stat, go_gb.LCDSTATCoincidenceFlag, val)
	go_gb.Set(&stat, go_gb.LCDSTATCoincidenceInterrupt, val)
	p.memory.Store(go_gb.LCDStatusRegister, stat)

	interrupt := p.memory.Read(go_gb.IF)
	if val {
		go_gb.Set(&interrupt, 1, true)
		p.memory.Store(go_gb.IF, interrupt)
	}
}

func (p *ppu) Step(mc go_gb.MC) {
	defer p.compareLyLyc()
	p.modeClock += mc

	switch p.currentMode {
	case 2:
		if p.modeClock >= 20 {
			p.setMode(3)
		}
	case 3:
		if p.modeClock >= 43 {
			// todo: HBLANK interrupt
			p.setMode(0)
			p.renderScanline()
		}
	case 0:
		if p.modeClock >= 22 {
			p.currentLine += 1
			p.updateLy()
			if p.currentLine == 143 {
				p.setMode(1)
				// todo: VBLANK interrupt
				// todo: render to display
			} else {
				p.setMode(2)
			}
		}
	case 1:
		if p.modeClock >= 114 {
			p.currentLine += 1
			if p.currentLine > 153 {
				p.currentLine = 0
				p.setMode(2)
			}
			p.updateLy()
		}
	}
}

func (p *ppu) renderScanline() {

}
