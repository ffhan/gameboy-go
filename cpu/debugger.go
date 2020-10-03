package cpu

import (
	"fmt"
	go_gb "go-gb"
	"io"
)

type debugger struct {
	cpu           *cpu
	output        io.Writer
	debugOn       bool
	useInstrNames bool
}

func NewDebugger(cpu *cpu, output io.Writer) *debugger {
	return &debugger{cpu: cpu, output: output}
}

func (d *debugger) Debug(val bool) {
	d.debugOn = val
}

func (d *debugger) PrintInstructionNames(val bool) {
	if val {
		go_gb.InitInstructions()
	}
	d.useInstrNames = val
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
	var instruction string
	if d.useInstrNames {
		if opcode > 0xFF {
			instruction = go_gb.Prefixed[byte(opcode)].String()
		} else {
			instruction = go_gb.Unprefixed[byte(opcode)].String()
		}
	}
	if d.debugOn {
		fmt.Fprintf(d.output, "OP: %04X\tPC: %04X\tSP: %04X\ta: %02X\tf: %02X\tb: %02X\tc: %02X\td: %02X\te: %02X\th: %02X\tl: %02X\tZNHC: %04b Instruction: '%s' PPU mode: %d line: %d\n",
			opcode, pc, d.cpu.sp,
			d.cpu.r[A], d.cpu.r[F],
			d.cpu.r[B], d.cpu.r[C],
			d.cpu.r[D], d.cpu.r[E],
			d.cpu.r[H], d.cpu.r[L], d.cpu.r[F]>>4,
			instruction,
			d.cpu.ppu.Mode(), d.cpu.ppu.CurrentLine())
	}
}
