package ppu

import (
	go_gb "go-gb"
)

const (
	frameBufASize = 160 * 144 * 4
)

type ppu struct {
	memory go_gb.Memory

	frameBuffer  [160 * 144 * 3]byte
	frameBufferA [frameBufASize]byte

	colorMapping [4 * 3]byte

	bgMapA [256 * 256 * 4]byte
}

func (p *ppu) getTileMapAddr() uint16 {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 3) {
		return 0x9C00
	}
	return 0x9800
}

func (p *ppu) getTileDataAddr() uint16 {
	if go_gb.Bit(p.memory.Read(go_gb.LCDControlRegister), 4) {
		return 0x8000
	}
	return 0x8800
}

func (p *ppu) GetFrameBuffer() [frameBufASize]byte {
	return p.frameBufferA
}

func (p *ppu) setBgColor(row, col int, pVal, colorval byte) {
	colorFromPalette := (pVal >> (2 * colorval)) & 3
	p.bgMapA[(row*256*4)+(col*4)] = p.colorMapping[colorFromPalette*3]
	p.bgMapA[(row*256*4)+(col*4)+1] = p.colorMapping[colorFromPalette*3+1]
	p.bgMapA[(row*256*4)+(col*4)+2] = p.colorMapping[colorFromPalette*3+2]
	p.bgMapA[(row*256*4)+(col*4)+3] = 0xFF
}

func (p *ppu) drawFrame() {
	for r := 0; r < 144; r++ {
		for col := 0; col < 160; col++ {
			yOffA := r * 256 * 4
			xOffA := col * 4
			p.frameBufferA[(r*160*4)+(col*4)] = p.bgMapA[yOffA+xOffA]
			p.frameBufferA[(r*160*4)+(col*4)+1] = p.bgMapA[yOffA+xOffA+1]
			p.frameBufferA[(r*160*4)+(col*4)+2] = p.bgMapA[yOffA+xOffA+2]
			p.frameBufferA[(r*160*4)+(col*4)+3] = p.bgMapA[yOffA+xOffA+3]
		}
	}
}

func (p *ppu) calcBg(row byte) {
	scx := p.memory.Read(go_gb.LCDSCX)
	scy := p.memory.Read(go_gb.LCDSCY)

	tileMap := p.getTileMapAddr()
	tileData := p.getTileDataAddr()

	pVal := p.memory.Read(go_gb.LCDBGP)
	for j := 0; j < 256; j++ {
		offY := uint16(row + scy)
		offX := uint16(byte(j) + scx)

		tileId := p.memory.Read(tileMap + ((offY / 8 * 32) + (offX / 8)))

		var colorval byte
		if tileData == 0x8800 {
			colorval = (p.memory.Read(tileData+0x800+uint16(int8(tileId)*0x10)+(offY%8*2)) >> (7 - (offX % 8)) & 0x1) + ((p.memory.Read(tileData+0x800+uint16(int8(tileId)*0x10)+(offY%8*2)+1) >> (7 - (offX % 8)) & 0x1) * 2)
		} else {
			colorval = (p.memory.Read(tileData+(uint16(tileId)*2)+(offY%8*2)) >> (7 - (offX % 8)) & 0x1) + (p.memory.Read(tileData+(uint16(tileId)*2)+(offY%8*2)+1)>>(7-(offX%8))&0x1)*2
		}
		p.setBgColor(int(row), j, pVal, colorval)
	}
}
