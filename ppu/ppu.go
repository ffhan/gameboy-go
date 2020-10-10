package ppu

import (
	go_gb "go-gb"
	"go-gb/memory"
	"sync"
)

const (
	frameBufASize = 160 * 144 * 4
)
const (
	White byte = iota
	LightGray
	DarkGray
	Black

	Transparent = White
)

type ppu struct {
	memory go_gb.Memory // used for usual memory access
	vram   go_gb.Memory // for skipping locks
	oam    go_gb.Memory // for skipping locks
	io     go_gb.Memory // optimized access to IO

	frameBuffer [160 * 144]byte // map colors in the display!

	currentLine int
	currentMode byte
	modeClock   go_gb.MC

	renderMutex sync.Mutex

	display go_gb.Display
}

func NewPpu(memory go_gb.Memory, vram go_gb.Memory, oam go_gb.Memory, io go_gb.Memory, display go_gb.Display) *ppu {
	return &ppu{memory: memory, vram: vram, oam: oam, io: io, currentMode: 2, display: display}
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

// returns tile data block 0 address and if addressing used unsigned offsets
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

func (p *ppu) oamInterrupt() {
	go_gb.Update(p.memory, go_gb.LCDSTAT, func(b byte) byte {
		go_gb.Set(&b, 5, true)
		return b
	})
	go_gb.Update(p.memory, go_gb.IF, func(b byte) byte {
		go_gb.Set(&b, int(go_gb.BitLCD), true)
		return b
	})
}

func (p *ppu) hblankInterrupt() {
	go_gb.Update(p.memory, go_gb.LCDSTAT, func(b byte) byte {
		go_gb.Set(&b, 3, true)
		return b
	})
	go_gb.Update(p.memory, go_gb.IF, func(b byte) byte {
		go_gb.Set(&b, int(go_gb.BitLCD), true)
		return b
	})
}

func (p *ppu) requestLCDSTATInterrupt() {
	go_gb.Update(p.io, go_gb.IF, func(b byte) byte {
		go_gb.Set(&b, int(go_gb.BitLCD), true)
		return b
	})
}

func (p *ppu) handleModeInterrupt(stat byte, modeBit int) {
	if go_gb.Bit(stat, modeBit) {
		p.requestLCDSTATInterrupt()
	}
}

func (p *ppu) vblankInterrupt() {
	go_gb.Update(p.memory, go_gb.IF, func(b byte) byte {
		go_gb.Set(&b, int(go_gb.BitVBlank), true)
		return b
	})
}

func (p *ppu) getScroll() (byte, byte) {
	scx := p.memory.Read(go_gb.LCDSCX)
	scy := p.memory.Read(go_gb.LCDSCY)
	return scx, scy
}

func (p *ppu) getLine() byte {
	return p.io.Read(go_gb.LCDLY)
}

func (p *ppu) updateLine() {
	p.io.Store(go_gb.LCDLY, byte(p.currentLine))
	p.compareLyLyc()
}

func (p *ppu) CurrentLine() int {
	return p.currentLine
}

func (p *ppu) getWindow() (byte, byte) {
	wx := p.memory.Read(go_gb.LCDWX)
	wy := p.memory.Read(go_gb.LCDWY)
	return wx, wy
}

func (p *ppu) setMode(mode byte, max go_gb.MC) {
	mode &= 0x3
	go_gb.Update(p.memory, go_gb.LCDSTAT, func(b byte) byte {
		return (b & 0xFC) | mode
	})
	p.currentMode = mode
	//fmt.Printf("mode %d clock %d line %d\n", p.currentMode, p.modeClock, p.currentLine)
	p.modeClock -= max
}

func (p *ppu) compareLyLyc() {
	val := p.io.Read(go_gb.LCDLYC) == p.getLine()
	go_gb.Update(p.io, go_gb.LCDSTAT, func(b byte) byte {
		go_gb.Set(&b, go_gb.LCDSTATCoincidenceFlag, val)
		return b
	})
	if val {
		go_gb.Update(p.memory, go_gb.IF, func(b byte) byte {
			go_gb.Set(&b, int(go_gb.BitLCD), true)
			return b
		})
	}
}

func (p *ppu) use8x16Sprites() bool {
	return go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 2)
}

func (p *ppu) getBgColor(colorNum byte) byte {
	return (p.io.Read(go_gb.LCDBGP) >> (colorNum * 2)) & 0x3
}
func (p *ppu) getSpriteColor(colorNum byte, address uint16) byte {
	return (p.io.Read(address) >> (colorNum * 2)) & 0x3
}

func (p *ppu) renderScanline() {
	p.renderBackgroundScanLine()
	if go_gb.Bit(p.io.Read(go_gb.LCDControlRegister), 1) {
		p.renderSpritesOnScanLine()
	}
}

func (p *ppu) Enabled() bool {
	return go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 7)
}

var seen = make(map[uint16]byte)

func (p *ppu) renderBackgroundScanLine() {
	scx, scy := p.getScroll()
	wx, wy := p.getWindow()

	line := p.getLine()
	usingWindow := p.windowEnabled() && wy <= line

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

	tileIds := p.vram.ReadBytes(mapAddr+tileRow+uint16(scx)/8, 20)

	lineNum := uint16(yPos % 8)
	lineNum *= 2

	var data1 [20]byte
	var data2 [20]byte
	for i := 0; i < 20; i++ {
		tileLocation := tileData
		//tileAddress := mapAddr + tileRow + tileCol
		tileId := uint16(tileIds[i])
		if unsigned {
			tileLocation += tileId * 16
		} else {
			if tileId < 128 {
				tileLocation += (tileId + 128) * 16
			} else {
				tileLocation += (tileId - 128) * 16
			}
		}

		data1[i] = p.vram.Read(tileLocation + lineNum)
		data2[i] = p.vram.Read(tileLocation + lineNum + 1)
	}

	var colors [4]byte
	for i := byte(0); i < 4; i++ {
		colors[i] = p.getBgColor(i)
	}

	//fmt.Printf("rendering background: line %d -> scx %d scy %d wx %d wy %d tiledata %X tileMap %X bg? %t tileRow %d\n",
	//	line, scx, scy, wx, wy, tileData, mapAddr, !usingWindow, tileRow)
	for pixel := byte(0); pixel < 160; pixel++ {
		pixel := pixel
		xPos := pixel + scx
		if usingWindow && pixel >= wx {
			xPos = pixel - wx
		}
		//tileCol := uint16(xPos) / 8

		data1 := data1[pixel/8]
		data2 := data2[pixel/8]

		colorBit := 7 - xPos%8

		colorNum := p.getColorNum(data1, data2, colorBit)

		// real color palettes will be done on the front end display
		colorId := colors[colorNum]
		bufferAddr := uint(line)*160 + uint(pixel)
		p.frameBuffer[bufferAddr] = colorId

		//fmt.Printf("pixel %d -> xPos %d tileCol %d tileLocation %X tileAddress %X tileId %d lineNum %d colorBit %d colorNum %d",
		//	pixel, xPos, tileCol, tileLocation, tileAddress, tileId, lineNum, colorBit, colorNum)
	}
}

func (p *ppu) getColorNum(data1 byte, data2 byte, colorBit byte) byte {
	colorNum := (data2 >> colorBit) & 1
	colorNum <<= 1
	colorNum |= (data1 >> colorBit) & 1
	return colorNum
}

func (p *ppu) renderSpritesOnScanLine() { // todo: handle 8x16 sprites
	use8x16 := p.use8x16Sprites()

	scanLine := p.getLine()
	ySize := byte(8)
	if use8x16 {
		ySize = 16
	}

	for sprite := 39; sprite >= 0; sprite-- {
		index := sprite * 4
		spriteData := p.oam.ReadBytes(memory.OAMStart+uint16(index), 4)
		yPos := spriteData[0] - 16
		xPos := spriteData[1] - 8
		if !(scanLine >= yPos && scanLine < (yPos+ySize)) {
			continue
		}
		tileLocation := spriteData[2]
		attributes := spriteData[3]

		var colorAddr uint16
		if go_gb.Bit(attributes, 4) {
			colorAddr = go_gb.LCDOBP1
		} else {
			colorAddr = go_gb.LCDOBP0
		}

		xFlip := go_gb.Bit(attributes, 5)
		yFlip := go_gb.Bit(attributes, 6)
		// If set to zero then sprite always rendered above bg
		// If set to 1, sprite is hidden behind the background and window
		// unless the color of the background or window is white, it's then rendered on top
		hidden := go_gb.Bit(attributes, 7)

		line := scanLine - yPos // line of the sprite

		if yFlip {
			line = ySize - 1 - line
		}

		line *= 2 // 2 bytes in a line

		dataAddress := (memory.VRAMStart + uint16(tileLocation)*16) + uint16(line)
		data := p.vram.ReadBytes(dataAddress, 2)
		low := data[0]
		high := data[1]

		for tilePixel := 7; tilePixel >= 0; tilePixel-- {
			colorBit := tilePixel
			if xFlip {
				colorBit = 7 - colorBit
			}

			pixel := xPos + byte(tilePixel)
			if pixel < 0 || pixel >= 160 {
				continue
			}

			colorNum := p.getColorNum(low, high, byte(colorBit))

			col := p.getSpriteColor(colorNum, colorAddr)
			//if col == Transparent {
			//	continue // don't update frame buffer
			//}
			frameBufferAddr := uint(scanLine)*160 + uint(pixel)
			if !hidden || p.frameBuffer[frameBufferAddr] == White {
				p.frameBuffer[frameBufferAddr] = col
			}
		}
	}
}

func (p *ppu) Mode() byte {
	return p.currentMode
}

func (p *ppu) Step(mc go_gb.MC) {
	p.modeClock += mc

	switch p.currentMode {
	case 2:
		if p.modeClock >= 20 {
			p.setMode(3, 20)
			p.renderScanline()
		}
	case 3:
		if p.modeClock >= 43 {
			p.setMode(0, 43)
		}
		p.handleModeInterrupt(p.io.Read(go_gb.LCDSTAT), 3)
	case 0:
		if p.modeClock >= 51 {
			p.currentLine += 1
			p.updateLine()
			if p.currentLine == 144 {
				p.vblankInterrupt()
				p.setMode(1, 51)
				p.display.Draw(p.frameBuffer[:])
			} else {
				p.setMode(2, 51)
			}
		}
		if p.currentLine == 144 {
			p.handleModeInterrupt(p.io.Read(go_gb.LCDSTAT), 4)
		} else {
			p.handleModeInterrupt(p.io.Read(go_gb.LCDSTAT), 5)
		}
	case 1:
		if p.modeClock >= 114 {
			p.currentLine += 1
			if p.currentLine > 153 {
				p.currentLine = 0
				p.setMode(2, 114)
			} else {
				p.modeClock -= 114
			}
			p.updateLine()
		}
		p.handleModeInterrupt(p.io.Read(go_gb.LCDSTAT), 5)
	}
}
