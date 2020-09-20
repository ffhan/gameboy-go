package cpu

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

type debugger struct {
	cpu    *cpu
	output io.Writer
}

func NewDebugger(cpu *cpu, output io.Writer) *debugger {
	return &debugger{cpu: cpu, output: output}
}

func (d *debugger) Step() go_gb.MC {
	d.print(d.cpu.memory.Read(d.cpu.pc))
	return d.cpu.Step()
}

func (d *debugger) print(opcode byte) {
	fmt.Fprintf(d.output, "OP: %X\tPC: %X\tSP: %X\ta: %X\tf: %X\tb: %X\tc: %X\td: %X\te: %X\th: %X\tl: %X\n", opcode, d.cpu.pc, d.cpu.sp, d.cpu.r[A], d.cpu.r[F], d.cpu.r[B], d.cpu.r[C], d.cpu.r[D], d.cpu.r[E], d.cpu.r[H], d.cpu.r[L])
}
