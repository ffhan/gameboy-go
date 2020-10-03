package memory

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"io"
)

type debugger struct {
	memory  cpu.MemoryBus
	output  io.Writer
	debugOn bool
}

func NewDebugger(memory cpu.MemoryBus, output io.Writer) *debugger {
	return &debugger{memory: memory, output: output}
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
	bytes := d.memory.ReadBytes(pointer, n)
	d.printf("read %d bytes from %X: %v\n", n, pointer, bytes)
	return bytes
}

func (d *debugger) Read(pointer uint16) byte {
	b := d.memory.Read(pointer)
	d.printf("read byte from %X: %X\n", pointer, b)
	return b
}

func (d *debugger) StoreBytes(pointer uint16, bytes []byte) {
	d.memory.StoreBytes(pointer, bytes)
	d.printf("stored bytes to %X: %v\n", pointer, bytes)
}

func (d *debugger) Store(pointer uint16, val byte) {
	d.memory.Store(pointer, val)
	d.printf("stored byte to %X: %v\n", pointer, val)
}

func (d *debugger) HRAM() go_gb.Memory {
	return d.memory.HRAM()
}

func (d *debugger) IO() go_gb.Memory {
	return d.memory.IO()
}

func (d *debugger) InterruptEnableRegister() go_gb.Memory {
	return d.memory.InterruptEnableRegister()
}
