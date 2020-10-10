package memory

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

func DumpOam(oam, vram go_gb.Memory, writer io.Writer) {
	for spriteAddr := uint16(0xFE00); spriteAddr < 0xFE9F; spriteAddr += 4 {
		y := oam.Read(spriteAddr)
		x := oam.Read(spriteAddr + 1)
		tileNum := oam.Read(spriteAddr + 2)
		attrs := oam.Read(spriteAddr + 3)

		tileAddr := VRAMStart + uint16(tileNum)
		tile := vram.ReadBytes(tileAddr, 16)

		for b := 0; b < 16; b += 2 {
			b1 := tile[b]
			b2 := tile[b+1]
			for i := 0; i < 8; i++ {
				c1 := (b1 & 0x80) >> 7
				c2 := (b2 & 0x80) >> 7
				col := (c2 << 1) | c1
				b1 <<= 1
				b2 <<= 1
				var char rune
				switch col {
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
			fmt.Fprintln(writer)
		}
		fmt.Fprintf(writer, "x: %d y: %d attrs: %b\n", x, y, attrs)
	}
}
