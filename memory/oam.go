package memory

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

func DumpOam(io, oam, vram go_gb.Memory, writer io.Writer) {
	const (
		rows             = 4
		tilesPerRow      = 10
		tileSize         = 8
		bytesPerTileLine = 2
		bytesPerSprite   = 4
		bytesPerTile     = 16
	)
	read := io.Read(go_gb.LCDControlRegister)
	fmt.Fprintf(writer, "sprites enabled: %t 8x16 sprites: %t\n",
		go_gb.Bit(read, 1),
		go_gb.Bit(read, 2))
	for sprite := uint16(0); sprite < rows*tilesPerRow; sprite++ {
		addr := OAMStart + sprite*4
		y := oam.Read(addr)
		x := oam.Read(addr + 1)
		tileId := oam.Read(addr + 2)
		attributes := oam.Read(addr + 3)

		var colorAddr uint16
		if go_gb.Bit(attributes, 4) {
			colorAddr = go_gb.LCDOBP1
		} else {
			colorAddr = go_gb.LCDOBP0
		}
		var colors [4]byte
		for i := byte(1); i < 4; i++ {
			colors[i] = (io.Read(colorAddr) >> (i * 2)) & 0x3
		}

		fmt.Fprintf(writer, "sprite %d: x %d y %d tileId %02X attributes %08b palette: %v\n",
			sprite, x, y, tileId, attributes, colors)
	}
	for line := 0; line < rows*tileSize; line++ {
		tileRow := line / tileSize
		rowStart := tileRow * tilesPerRow
		for tile := 0; tile < tilesPerRow; tile++ {
			tileNum := rowStart + tile
			oamAddr := OAMStart + uint16(bytesPerSprite*tileNum)
			attrs := oam.Read(oamAddr + 3)

			var colorAddr uint16
			if go_gb.Bit(attrs, 4) {
				colorAddr = go_gb.LCDOBP1
			} else {
				colorAddr = go_gb.LCDOBP0
			}
			tileId := oam.Read(oamAddr + 2)

			xFlip := go_gb.Bit(attrs, 5)
			yFlip := go_gb.Bit(attrs, 6)

			if yFlip {
				// todo: implement yflip
			}

			addr := VRAMStart + uint16(tileId)*uint16(bytesPerTile) + (uint16(line) % tileSize * bytesPerTileLine)
			low := vram.Read(addr)
			high := vram.Read(addr + 1)

			var colors [4]byte
			for i := byte(1); i < 4; i++ {
				colors[i] = (io.Read(colorAddr) >> (i * 2)) & 0x3
			}

			for pxl := 7; pxl >= 0; pxl-- {
				colorBit := pxl
				if xFlip {
					colorBit = 7 - colorBit
				}

				colorNum := (high >> colorBit) & 1
				colorNum <<= 1
				colorNum |= (low >> colorBit) & 1

				colorId := colors[colorNum]
				var char rune
				switch colorId {
				case 0:
					char = '▁'
				case 1:
					char = '░'
				case 2:
					char = '▒'
				case 3:
					char = '▓'
				}
				fmt.Fprint(writer, string(char))
			}
		}
		fmt.Fprintln(writer)
	}
}
