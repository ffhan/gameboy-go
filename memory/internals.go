package memory

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

type mmap struct {
	start, end uint16
	memory     []byte
}

func newMmap(start uint16, end uint16, memory []byte) *mmap {
	return &mmap{start: start, end: end + 1, memory: memory}
}

func DumpVram(io, vram go_gb.Memory, writer io.Writer) {
	wx := io.Read(go_gb.LCDWX)
	wy := io.Read(go_gb.LCDWY)
	scx := io.Read(go_gb.LCDSCX)
	scy := io.Read(go_gb.LCDSCY)
	fmt.Fprintf(writer, "LY: %d, scx: %d, scy: %d, wx: %d, wy: %d\n",
		io.Read(go_gb.LCDLY), scx, scy, wx, wy)

	for line := 0; line < 256; line++ {
		line := byte(line)
		usingWindow := go_gb.Bit(io.Read(go_gb.LCDControlRegister), 5) && wy <= line

		var tileData uint16
		var unsigned bool
		if go_gb.Bit(io.Read(go_gb.LCDControlRegister), 4) {
			tileData, unsigned = 0x8000, true
		} else {
			tileData, unsigned = 0x8800, false
		}

		yPos := line
		var mapAddr uint16
		lookBit := 3
		if usingWindow {
			lookBit = 6
			yPos = line - wy
		}
		if go_gb.Bit(io.Read(go_gb.LCDControlRegister), lookBit) {
			mapAddr = 0x9C00
		} else {
			mapAddr = 0x9800
		}
		tileRow := uint16(yPos/8) * 32

		tileIds := vram.ReadBytes(mapAddr+tileRow, 32)

		lineNum := uint16(yPos%8) * 2

		var data1 [32]byte
		var data2 [32]byte

		for i := 0; i < 32; i++ {
			tileLocation := tileData
			tileId := uint16(tileIds[i])

			if unsigned {
				tileLocation += tileId * 16
			} else {
				if tileId < 128 {
					tileLocation += (tileId + 128) * 16
				} else {
					tileLocation -= (tileId - 128) * 16
				}
			}

			data1[i] = vram.Read(tileLocation + lineNum)
			data2[i] = vram.Read(tileLocation + lineNum + 1)
		}

		var colors [4]byte
		for i := byte(0); i < 4; i++ {
			colors[i] = (io.Read(go_gb.LCDBGP) >> (i * 2)) & 0x3
		}

		for pixel := 0; pixel < 256; pixel++ {
			pixel := byte(pixel)

			xPos := pixel
			if usingWindow && pixel >= wx {
				xPos = pixel - wx
			}

			data1 := data1[pixel/8]
			data2 := data2[pixel/8]

			colorBit := 7 - xPos%8
			colorNum := (data2 >> colorBit) & 1
			colorNum <<= 1
			colorNum |= (data1 >> colorBit) & 1

			colorId := colors[colorNum]
			var char rune
			switch colorId {
			case 0:
				char = '▓'
			case 1:
				char = '▒'
			case 2:
				char = '░'
			case 3:
				char = '▁'
			}
			fmt.Fprint(writer, string(char))
		}
		fmt.Fprintln(writer)
	}
}

func (m *mmap) ReadBytes(pointer, n uint16) []byte {
	i := pointer - m.start
	return m.memory[i : i+n]
}

func (m *mmap) Read(pointer uint16) byte {
	i := pointer - m.start
	return m.memory[i]
}

func (m *mmap) StoreBytes(pointer uint16, bytes []byte) {
	i := pointer - m.start
	copy(m.memory[i:i+uint16(len(bytes))], bytes)
}

func (m *mmap) Store(pointer uint16, val byte) {
	m.memory[pointer-m.start] = val
}

func (m *mmap) LoadRom(bytes []byte) int {
	n := len(m.memory)
	copy(m.memory, bytes[:n])
	return n
}

type lockedMemory struct {
}

func (l *lockedMemory) ReadBytes(pointer, n uint16) []byte {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = 0xFF
	}
	return bytes
}

func (l *lockedMemory) Read(pointer uint16) byte {
	return 0xFF
}

func (l *lockedMemory) StoreBytes(pointer uint16, bytes []byte) {
	fmt.Printf("storing locked bytes at %X\n", pointer)
}

func (l *lockedMemory) Store(pointer uint16, val byte) {
	fmt.Printf("storing a locked byte at %X\n", pointer)
}
