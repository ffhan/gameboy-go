package memory

import (
	go_gb "go-gb"
	"io"
)

type debugger struct {
	memory go_gb.Memory
	output io.Writer
}

func NewDebugger(memory go_gb.Memory, output io.Writer) *debugger {
	return &debugger{memory: memory, output: output}
}

func (d *debugger) ReadBytes(pointer, n uint16) []byte {
	bytes := d.memory.ReadBytes(pointer, n)
	//fmt.Fprintf(d.output, "read %d bytes from %X: %v\n", n, pointer, bytes)
	return bytes
}

func (d *debugger) Read(pointer uint16) byte {
	b := d.memory.Read(pointer)
	//fmt.Fprintf(d.output, "read byte from %X: %X\n", pointer, b)
	return b
}

func (d *debugger) StoreBytes(pointer uint16, bytes []byte) {
	d.memory.StoreBytes(pointer, bytes)
	//fmt.Fprintf(d.output, "stored bytes to %X: %v\n", pointer, bytes)
}

func (d *debugger) Store(pointer uint16, val byte) {
	d.memory.Store(pointer, val)
	//fmt.Fprintf(d.output, "stored byte to %X: %v\n", pointer, val)
}
