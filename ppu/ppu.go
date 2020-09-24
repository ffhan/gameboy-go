package ppu

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/memory"
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

	frameBuffer [160 * 144]byte // map colors in the display!

	currentLine int
	currentMode byte
	modeClock   go_gb.MC

	display go_gb.Display
}

func NewPpu(memory go_gb.Memory, vram go_gb.Memory, oam go_gb.Memory, display go_gb.Display) *ppu {
	return &ppu{memory: memory, vram: vram, oam: oam, currentMode: 2, display: display}
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

func (p *ppu) vblankInterrupt() {
	go_gb.Update(p.memory, go_gb.IF, func(b byte) byte {
		go_gb.Set(&b, int(go_gb.BitVBlank), true)
		return b
	})
}

func (p *ppu) IsVBlank() bool {
	return go_gb.Bit(p.memory.Read(go_gb.IF), int(go_gb.BitVBlank))
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
	go_gb.Update(p.memory, go_gb.LCDSTAT, func(b byte) byte {
		return (b & 0xFC) | mode
	})
	p.modeClock = 0
	p.currentMode = mode
	fmt.Printf("mode %d clock %d line %d\n", p.currentMode, p.modeClock, p.currentLine)
}

func (p *ppu) compareLyLyc() {
	val := p.memory.Read(go_gb.LCDLYC) == p.getLine()
	stat := p.memory.Read(go_gb.LCDSTAT)
	go_gb.Set(&stat, go_gb.LCDSTATCoincidenceFlag, val)
	go_gb.Set(&stat, go_gb.LCDSTATCoincidenceInterrupt, val)
	p.memory.Store(go_gb.LCDSTAT, stat)

	interrupt := p.memory.Read(go_gb.IF)
	if val {
		go_gb.Set(&interrupt, 1, true)
		p.memory.Store(go_gb.IF, interrupt)
	}
}

func (p *ppu) use8x16Sprites() bool {
	return go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 2)
}

func (p *ppu) getColor(colorNum byte, address uint16) byte {
	return (p.memory.Read(address) >> (colorNum * 2)) & 0x3
}

func (p *ppu) renderScanline() {
	p.renderBackgroundScanLine()
	p.renderSpritesOnScanLine()
}

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

	//fmt.Printf("rendering background: line %d -> scx %d scy %d wx %d wy %d tiledata %X tileMap %X bg? %t tileRow %d\n",
	//	line, scx, scy, wx, wy, tileData, mapAddr, !usingWindow, tileRow)
	for pixel := byte(0); pixel < 160; pixel++ {
		xPos := pixel + scx
		if usingWindow && pixel >= wx {
			xPos = pixel - wx
		}
		tileCol := uint16(xPos) / 8

		tileLocation := tileData
		tileAddress := mapAddr + tileRow + tileCol
		tileId := uint16(p.vram.Read(tileAddress))
		if unsigned {
			tileLocation += tileId * 16
		} else {
			if tileId < 128 {
				tileLocation += (tileId + 128) * 16
			} else {
				tileLocation -= (tileId - 128) * 16
			}
		}

		lineNum := yPos % 8
		lineNum *= 2
		data1 := p.vram.Read(tileLocation + uint16(lineNum))
		data2 := p.vram.Read(tileLocation + uint16(lineNum) + 1)

		colorBit := 7 - xPos%8

		colorNum := p.getColorNum(data1, data2, colorBit)

		finally := int(line)
		if finally < 0 || finally > 143 || pixel < 0 || pixel > 159 {
			continue
		}

		// real color palettes will be done on the front end display
		colorId := p.getColor(colorNum, go_gb.LCDBGP)
		p.frameBuffer[line*160+pixel] = colorId

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

func (p *ppu) renderSpritesOnScanLine() {
	use8x16 := p.use8x16Sprites()

	scanLine := p.getLine()
	ySize := byte(8)
	if use8x16 {
		ySize = 16
	}

	for sprite := byte(0); sprite < 40; sprite++ {
		index := sprite * 4
		spriteData := p.oam.ReadBytes(memory.OAMStart+uint16(index), 4)
		yPos := spriteData[0] - 16
		xPos := spriteData[1] - 8
		tileLocation := spriteData[2]
		attributes := spriteData[3]

		xFlip := go_gb.Bit(attributes, 5)
		yFlip := go_gb.Bit(attributes, 6)

		if scanLine >= yPos && scanLine < (yPos+ySize) {
			line := scanLine - yPos // line of the sprite

			if yFlip {
				line = ySize - 1 - line
			}

			line *= 2 // 2 bytes in a line

			dataAddress := (0x8000 + uint16(tileLocation)*16) + uint16(line)
			data := p.vram.ReadBytes(dataAddress, 2)
			data1 := data[0]
			data2 := data[1]

			for tilePixel := 7; tilePixel >= 0; tilePixel-- {
				colorBit := tilePixel
				if xFlip {
					colorBit = 7 - colorBit
				}

				colorNum := p.getColorNum(data1, data2, byte(colorBit))
				var colorAddr uint16
				if go_gb.Bit(attributes, 4) {
					colorAddr = go_gb.LCDOBP1
				} else {
					colorAddr = go_gb.LCDOBP0
				}

				col := p.getColor(colorNum, colorAddr)
				if col == Transparent {
					continue // don't update frame buffer
				}

				pixel := xPos + byte(tilePixel)
				p.frameBuffer[scanLine*160+pixel] = col
			}
		}
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
			// todo: render scanline to display
			p.setMode(0)
			p.renderScanline()
		}
	case 0:
		if p.modeClock >= 51 {
			p.currentLine += 1
			p.updateLine()
			if p.currentLine == 143 {
				p.setMode(1)
				p.vblankInterrupt()
				p.display.Draw(p.frameBuffer[:])
				fmt.Printf("drawing %v\n", p.frameBuffer)
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
