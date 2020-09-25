package cpu

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

type debugger struct {
	cpu     *cpu
	output  io.Writer
	debugOn bool
}

func NewDebugger(cpu *cpu, output io.Writer) *debugger {
	return &debugger{cpu: cpu, output: output}
}

func (d *debugger) Debug(val bool) {
	d.debugOn = val
}

func (d *debugger) PC() uint16 {
	return d.cpu.pc
}

func (d *debugger) Step() go_gb.MC {
	pc := d.cpu.pc
	op := uint16(d.cpu.memory.Read(d.cpu.pc))
	if op == 0xCB {
		op = (op << 8) | uint16(d.cpu.memory.Read(d.cpu.pc+1))
	}
	mc := d.cpu.Step()
	d.print(op, pc)
	return mc
}

func (d *debugger) IME() bool {
	return d.cpu.ime
}

func (d *debugger) print(opcode uint16, pc uint16) {
	if d.debugOn {
		fmt.Fprintf(d.output, "OP: %X\tPC: %X\tSP: %X\ta: %X\tf: %X\tb: %X\tc: %X\td: %X\te: %X\th: %X\tl: %X\tZNHC: %b\n",
			opcode, pc, d.cpu.sp,
			d.cpu.r[A], d.cpu.r[F],
			d.cpu.r[B], d.cpu.r[C],
			d.cpu.r[D], d.cpu.r[E],
			d.cpu.r[H], d.cpu.r[L], d.cpu.r[F]>>4)
	}
}
