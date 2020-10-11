package cpu

import (
	"fmt"
	go_gb "go-gb"
	"io"
	"os"
)

type debuggerQueue interface {
	Push(op, pc, sp uint16, a, f, b, c, d, e, h, l, flags byte, instruction string, ppuMode, ppuLine byte)
	fmt.Stringer
}

type debugger struct {
	cpu           *cpu
	output        io.Writer
	debugOn       bool
	useInstrNames bool

	PrintEveryCycle bool

	instructionQueue debuggerQueue
}

func NewDebugger(cpu *cpu, output io.Writer, instructionQueue debuggerQueue) *debugger {
	return &debugger{cpu: cpu, output: output, instructionQueue: instructionQueue}
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

var setD = false

func (d *debugger) Step() go_gb.MC {
	pc := d.cpu.pc
	op := uint16(d.cpu.memory.Read(d.cpu.pc))
	if op == 0xCB {
		op = (op << 8) | uint16(d.cpu.memory.Read(d.cpu.pc+1))
	}
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(d.output, "failed OP %X at %X, printing stack trace...\n", op, pc)
			fmt.Fprintln(d.output, d.instructionQueue.String())
			panic(err)
		}
	}()
	mc := d.cpu.Step()
	d.print(op, pc)
	return mc
}

func (d *debugger) IME() bool {
	return d.cpu.ime
}

func (d *debugger) print(opcode uint16, pc uint16) {
	if d.debugOn {
		var instruction string
		if d.useInstrNames {
			if opcode > 0xFF {
				instruction = go_gb.Prefixed[byte(opcode)].String()
			} else {
				instruction = go_gb.Unprefixed[byte(opcode)].String()
			}
		}
		queue := d.instructionQueue
		sp, a, f, b, c, d, e, h, l, flags, ppuMode, ppuLine := d.cpu.sp, d.cpu.r[go_gb.A], d.cpu.r[go_gb.F], d.cpu.r[go_gb.B], d.cpu.r[go_gb.C], d.cpu.r[go_gb.D], d.cpu.r[go_gb.E], d.cpu.r[go_gb.H], d.cpu.r[go_gb.L], d.cpu.r[go_gb.F]>>4, d.cpu.ppu.Mode(), d.cpu.ppu.CurrentLine()
		queue.Push(opcode, pc, sp, a, f, b, c, d, e, h, l, flags, instruction, ppuMode, byte(ppuLine))
	}
	if d.PrintEveryCycle {
		DumpCpu(os.Stdout, d.cpu, d.cpu.ppu)
	}
}

func (d *debugger) SP() uint16 {
	return d.cpu.sp
}

func (d *debugger) GetRegister(name go_gb.RegisterName) []byte {
	return d.cpu.GetRegister(name)
}

func DumpCpu(writer io.Writer, c go_gb.Cpu, p go_gb.PPU) {
	fmt.Fprintf(writer, "PC: %04X\tSP: %04X\ta: %02X\tf: %02X\tb: %02X\tc: %02X\td: %02X\te: %02X\th: %02X\tl: %02X\tZNHC: %04b PPU mode: %d line: %d\n",
		c.PC(), c.SP(),
		c.GetRegister(go_gb.A)[0], c.GetRegister(go_gb.F)[0],
		c.GetRegister(go_gb.B)[0], c.GetRegister(go_gb.C)[0],
		c.GetRegister(go_gb.D)[0], c.GetRegister(go_gb.E)[0],
		c.GetRegister(go_gb.H)[0], c.GetRegister(go_gb.L)[0],
		c.GetRegister(go_gb.F)[0]>>4,
		p.Mode(), p.CurrentLine())
}
