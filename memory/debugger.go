package memory

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

type debugger struct {
	go_gb.MemoryBus
	output  io.Writer
	debugOn bool
}

func NewDebugger(memory go_gb.MemoryBus, output io.Writer) *debugger {
	return &debugger{MemoryBus: memory, output: output}
}

func (d *debugger) Debug(val bool) {
	d.debugOn = val
}

func (d *debugger) printf(format string, args ...interface{}) {
	if d.debugOn {
		fmt.Fprintf(d.output, format, args...)
	}
}

func (d *debugger) ReadBytes(pointer, n uint16) []byte {
	bytes := d.MemoryBus.ReadBytes(pointer, n)
	d.printf("read %d bytes from %X: %v\n", n, pointer, bytes)
	return bytes
}

func (d *debugger) Read(pointer uint16) byte {
	b := d.MemoryBus.Read(pointer)
	d.printf("read byte from %X: %X\n", pointer, b)
	return b
}

func (d *debugger) StoreBytes(pointer uint16, bytes []byte) {
	d.MemoryBus.StoreBytes(pointer, bytes)
	d.printf("stored bytes to %X: %v\n", pointer, bytes)
}

func (d *debugger) Store(pointer uint16, val byte) {
	d.MemoryBus.Store(pointer, val)
	d.printf("stored byte to %X: %v\n", pointer, val)
}

func DumpMemory(writer io.Writer, memory go_gb.Memory, start, end uint16) {
	result := memory.ReadBytes(start, end-start+1)
	for i, b := range result {
		i := uint16(i)
		fmt.Fprintf(writer, "%X: %X\t%d\t%b\n", start+i, b, b, b)
	}
}
