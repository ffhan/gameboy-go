package go_gb

import (
	"fmt"
	"io"
)

const (
	LCDControlRegister uint16 = 0xFF40
	LCDSTAT            uint16 = 0xFF41

	LCDSCY uint16 = 0xFF42 // Scroll Y (R/W)
	LCDSCX uint16 = 0xFF43 // Scroll X (R/W)
	LCDLY  uint16 = 0xFF44 // LCDC Y-Coordinate (R)
	LCDLYC uint16 = 0xFF45 // LY Compare (R/W)
	LCDWY  uint16 = 0xFF4A // LCD Window Y position (R/W)
	LCDWX  uint16 = 0xFF4B // LCD Window X position minus 7 (R/W)

	/*
		 LCD BG Palette data (R/W)

		This register assigns gray shades to the color numbers of the BG and Window tiles.

		  Bit 7-6 - Shade for Color Number 3
		  Bit 5-4 - Shade for Color Number 2
		  Bit 3-2 - Shade for Color Number 1
		  Bit 1-0 - Shade for Color Number 0

		The four possible gray shades are:

		  0  White
		  1  Light gray
		  2  Dark gray
		  3  Black

	*/
	LCDBGP  uint16 = 0xFF47
	LCDOBP0 uint16 = 0xFF48 // Object Palette 0 data (R/W)
	LCDOBP1 uint16 = 0xFF49 // Object Palette 1 data (R/W)

	LCDBCPS uint16 = 0xFF68 // Background palette index
	LCDBCPD uint16 = 0xFF69 // Background palette data
	LCDOCPS uint16 = 0xFF6A // Sprite palette index
	LCDOCPD uint16 = 0xFF6B // Sprite palette data

	LCDVBK uint16 = 0xFF4F // LCD VRAM Bank

	LCDDMA uint16 = 0xFF46 // LCD OAM DMA transfer and start address (W)

	LCDHDMA1 uint16 = 0xFF51 // LCD CGB Mode Only - New DMA Source, High
	LCDHDMA2 uint16 = 0xFF52 // LCD CGB Mode Only - New DMA Source, Low
	LCDHDMA3 uint16 = 0xFF53 // LCD CGB Mode Only - New DMA Destination, High
	LCDHDMA4 uint16 = 0xFF54 // LCD CGB Mode Only - New DMA Destination, Low
	LCDHDMA5 uint16 = 0xFF55 // LCD CGB Mode Only - New DMA Length/Mode/Start
)

const (
	LCDSTATCoincidenceFlag = iota + 2
	LCDSTATHBlankInterrupt
	LCDSTATVBlankInterrupt
	LCDSTATOAMInterruptFlag
	LCDSTATCoincidenceInterrupt
)

type Display interface {
	Draw(bufferLine []byte)
	// calling this method returns if the display is drawing, and sets it to false after the method call
	IsDrawing() bool
}

type NopDisplay struct {
	debugOn   bool
	isDrawing bool

	buffer []byte
}

func NewNopDisplay() *NopDisplay {
	return &NopDisplay{}
}

func (n *NopDisplay) Debug(val bool) {
	n.debugOn = val
}

func (n *NopDisplay) IsDrawing() bool {
	defer func() {
		n.isDrawing = false
	}()
	return n.isDrawing
}

func (n *NopDisplay) Draw(buffer []byte) {
	n.buffer = buffer
	n.isDrawing = true
	if n.debugOn {
		fmt.Printf("screen buffer: %v\n", buffer)
	}
}

func DumpDisplay(writer io.Writer, display *NopDisplay) {
	for i, val := range display.buffer {
		var char rune
		switch val {
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
		if (i % 160) == 159 {
			fmt.Fprintln(writer)
		}
	}
}
