package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	structure = `package memory

type bios struct {
	rom [0x100]byte
}

func NewBios(rom [0x100]byte) *bios {
	b := &bios{rom: rom}
	b.init()
	return b
}

func (b *bios) init() {
	b.rom = [0x100]byte{%s}
}

func (b *bios) ReadBytes(pointer, n uint16) []byte {
	return b.rom[pointer : pointer+n]
}

func (b *bios) Read(pointer uint16) byte {
	return b.rom[pointer]
}

func (b *bios) StoreBytes(pointer uint16, bytes []byte) {
	panic("invalid op")
}

func (b *bios) Store(pointer uint16, val byte) {
	panic("invalid op")
}
`
)

func main() {
	rom, err := ioutil.ReadFile("cmd/boot.gb")
	if err != nil {
		panic(err)
	}

	var sb strings.Builder
	for i, val := range rom {
		sb.Write([]byte(strconv.Itoa(int(val))))
		if i != len(rom)-1 {
			sb.Write([]byte{','})
			if i%16 == 15 {
				sb.Write([]byte{'\n'})
			}
		}
	}

	file, err := os.Create("memory/bios.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, structure, sb.String())
}
