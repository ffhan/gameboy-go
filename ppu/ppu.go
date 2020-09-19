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

	frameBuffer [160 * 144 * 3]byte

	currentLine int
	currentMode byte
	modeClock   go_gb.MC
}

func (p *ppu) getBgTileMapAddr() uint16 {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 3) {
		return 0x9C00
	}
	return 0x9800
}

func (p *ppu) getWindowTileMapAddr() uint16 {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 6) {
		return 0x9C00
	}
	return 0x9800
}

func (p *ppu) getTileDataAddr() (uint16, bool) {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 4) {
		return 0x8000, true
	}
	return 0x8800, false
}

func (p *ppu) backgroundEnabled() bool {
	return go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 0)
}

func (p *ppu) windowEnabled() bool {
	return go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 5)
}

func (p *ppu) getScroll() (byte, byte) {
	scx := p.memory.Read(go_gb.LCDSCX)
	scy := p.memory.Read(go_gb.LCDSCY)
	return scx, scy
}

func (p *ppu) getLine() byte {
	return p.memory.Read(go_gb.LCDLY)
}

func (p *ppu) updateLine() {
	p.memory.Store(go_gb.LCDLY, byte(p.currentLine))
}

func (p *ppu) getWindow() (byte, byte) {
	wx := p.memory.Read(go_gb.LCDWX)
	wy := p.memory.Read(go_gb.LCDWY)
	return wx, wy
}

func (p *ppu) setMode(mode byte) {
	mode &= 0x3
	stat := p.memory.Read(go_gb.LCDStatusRegister) | mode
	p.memory.Store(go_gb.LCDStatusRegister, stat)
	p.modeClock = 0
	p.currentMode = mode
}

func (p *ppu) compareLyLyc() {
	val := p.memory.Read(go_gb.LCDLYC) == p.getLine()
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

func (p *ppu) getBgColor(colorNum byte) (byte, byte, byte) {
	col := (p.memory.Read(go_gb.LCDBGP) >> (colorNum * 2)) & 0x3
	switch col {
	case 0:
		return 0xFF, 0xFF, 0xFF
	case 1:
		return 0xCC, 0xCC, 0xCC
	case 2:
		return 0x77, 0x77, 0x77
	case 3:
		return 0, 0, 0
	}
	panic("invalid color number")
}

func (p *ppu) renderScanline() {
	p.renderBackgroundScanLine()
	p.renderSpriteOnScanLine()
}

func (p *ppu) renderBackgroundScanLine() {
	scx, scy := p.getScroll()
	wx, wy := p.getWindow()

	line := p.getLine()

	usingWindow := true
	if p.windowEnabled() {
		if wy <= line {
			usingWindow = false
		}
	}
	var yPos byte
	tileData, unsigned := p.getTileDataAddr()
	var mapAddr uint16
	if !usingWindow {
		mapAddr = p.getBgTileMapAddr()
		yPos = scy + line
	} else {
		mapAddr = p.getWindowTileMapAddr()
		yPos = line - wy
	}
	tileRow := uint16(yPos/8) * 32
	for pixel := byte(0); pixel < 160; pixel++ {
		xPos := pixel + scx
		if usingWindow && pixel >= wx {
			xPos = pixel - wx
		}
		tileCol := uint16(xPos) / 8

		tileLocation := tileData
		tileAddress := mapAddr + tileRow + tileCol
		tileId := uint16(p.memory.Read(tileAddress))
		if unsigned {
			tileLocation += tileId * 16
		} else {
			tileLocation = uint16(int(tileLocation) + (int(tileId)+128)*16)
		}

		lineNum := yPos % 8
		lineNum *= 2
		data1 := p.memory.Read(tileLocation + uint16(lineNum))
		data2 := p.memory.Read(tileLocation + uint16(lineNum) + 1)

		colorBit := int(xPos) % 8
		colorBit -= 7
		colorBit *= -1

		colorNum := (data2 >> colorBit) & 1
		colorNum <<= 1
		colorNum |= (data1 >> colorBit) & 1

		finally := int(line)
		if finally < 0 || finally > 143 || pixel < 0 || pixel > 159 {
			continue
		}

		red, green, blue := p.getBgColor(colorNum)

		p.frameBuffer[line*160+pixel] = red
		p.frameBuffer[line*160+pixel+1] = green
		p.frameBuffer[line*160+pixel+2] = blue
	}
}

func (p *ppu) renderSpriteOnScanLine() {

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
			// todo: render scanline to display
			p.setMode(0)
			p.renderScanline()
		}
	case 0:
		if p.modeClock >= 22 {
			p.currentLine += 1
			p.updateLine()
			if p.currentLine == 143 {
				p.setMode(1)
				// todo: VBLANK interrupt
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
			p.updateLine()
		}
	}
}
