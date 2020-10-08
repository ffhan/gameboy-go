package memory

import (
	"fmt"
	"io"
)

type mmap struct {
	start, end uint16
	memory     []byte
}

func newMmap(start uint16, end uint16, memory []byte) *mmap {
	return &mmap{start: start, end: end + 1, memory: memory}
}

func (m *mmap) Dump(writer io.Writer) {
	for i := uint16(0x9800); i < m.end; i++ {
		o := i - 0x9800
		if o > 0 && o%32 == 0 {
			fmt.Fprintln(writer)
		}
		fmt.Fprintf(writer, "%02X ", m.memory[i-m.start])
	}
	fmt.Fprintln(writer)
	fmt.Fprintln(writer)

	for i := uint16(0x8000); i <= 0x97FF; {
		pxl1 := m.memory[i-m.start]
		pxl2 := m.memory[i-m.start+1]
		i += 2
		o := i - 0x8000
		for p := 7; p >= 0; p-- {
			color := (((pxl2 >> p) & 1) << 1) | ((pxl1 >> p) & 1)
			var char rune
			switch color {
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
		if o > 0 && o%2 == 0 {
			fmt.Fprintln(writer)
		}
		if o > 0 && o%16 == 0 {
			fmt.Fprintln(writer)
			fmt.Fprintf(writer, "ID: %X\n", o/0x10)
		}
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
