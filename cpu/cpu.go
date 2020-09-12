package cpu

import (
	"go-gb"
)

const (
	BitZ = 7
	BitN = 6
	BitH = 5
	BitC = 4

	PrefixCB byte = 0xCB
	PrefixDD byte = 0xDD
	PrefixED byte = 0xED
	PrefixFD byte = 0xFD
)

type registerName uint16

const (
	A registerName = iota
	F
	B
	C
	D
	E
	H
	L
	AF
	BC
	DE
	HL
)

// executes specific things on the cpu and returns the number of m cycles it took to execute
type Instr func(c *cpu) go_gb.MC

type cpu struct {
	pc, sp uint16

	r              [8]byte // registers
	af, bc, de, hl []byte  // double registers

	rMap [][]byte // register mappings

	memory *go_gb.MemoryBus

	halt      bool
	eiWaiting byte
	diWaiting byte
	ime       bool // Interrupt master enable
}

func NewCpu() *cpu {
	c := &cpu{}
	c.init()
	return c
}

func (c *cpu) readOpcode() (byte, go_gb.MC) {
	val, m := c.memory.Read(c.pc)
	c.pc += 1
	return val, m
}
func (c *cpu) readFromPc(size uint16) ([]byte, go_gb.MC) {
	val, m := c.memory.ReadBytes(c.pc, size)
	c.pc += size
	return val, m
}

func (c *cpu) setPc(val uint16) go_gb.MC {
	c.pc = val
	return 1
}

func (c *cpu) popStack(size int) ([]byte, go_gb.MC) {
	bytes := make([]byte, size)
	m := go_gb.MC(0)
	for i := 0; i < size; i++ {
		c.sp += 1
		v, mc := c.memory.Read(c.sp)
		bytes[i] = v
		m += mc
	}
	return bytes, m
}

func (c *cpu) pushStack(b []byte) go_gb.MC {
	m := go_gb.MC(0)
	for _, val := range b {
		m += c.memory.Store(c.sp, val)
		c.sp -= 1
	}
	return m
}

func (c *cpu) getRegister(r registerName) []byte {
	return c.rMap[r]
}

func (c *cpu) init() {
	c.memory = go_gb.NewMemoryBus()
	c.pc = 0x0100
	c.sp = 0xFFFE
	// todo: set r to init values
	// setting references to register arr
	c.af = c.r[A : F+1]
	c.bc = c.r[B : C+1]
	c.de = c.r[D : E+1]
	c.hl = c.r[H : L+1]
	c.rMap = [][]byte{
		c.r[A : A+1], c.r[F : F+1], c.r[B : B+1], c.r[C : C+1],
		c.r[D : D+1], c.r[E : E+1], c.r[H : H+1], c.r[L : L+1],
		c.af, c.bc, c.de, c.hl,
	}
}

func (c *cpu) Step() {
	var cycles go_gb.MC
	if !c.halt {
		opcode, mc := c.readOpcode()
		instr := optable[opcode]
		cycles += instr(c)
		cycles += mc
	} else {
		cycles = 1
	}
	c.handleEiDi()
	c.handleInterrupts()
}

func (c *cpu) handleInterrupts() {
	if !c.ime { // interrupt master disabled
		return
	}
	var cycles go_gb.MC
	ifRegister, m := c.memory.Read(go_gb.IF)
	ieRegister, m2 := c.memory.Read(go_gb.IE)
	cycles += m + m2

	if ifRegister == 0 { // no interrupt flags set
		return
	}

	for _, interrupt := range go_gb.Interrupts {
		if go_gb.ShouldServiceInterrupt(ieRegister, ifRegister, interrupt.Bit) {
			c.serviceInterrupt(ifRegister, interrupt)
			return
		}
	}
}

func (c *cpu) serviceInterrupt(ifR byte, interrupt go_gb.Interrupt) go_gb.MC {
	var cycles go_gb.MC
	go_gb.Set(&ifR, int(interrupt.Bit), false)
	c.ime = false
	cycles += c.memory.Store(go_gb.IF, ifR)
	cycles += callAddr(c, go_gb.MsbLsbBytes(interrupt.JpAddr, true))
	return cycles
}

func (c *cpu) handleEiDi() {
	if c.eiWaiting != 0 {
		c.eiWaiting -= 1
		if c.eiWaiting == 0 {
			c.ime = true
		}
	}
	if c.diWaiting != 0 {
		c.diWaiting -= 1
		if c.diWaiting == 0 {
			c.ime = false
		}
	}
}

func (c *cpu) setFlag(bit int, val bool) {
	register := &c.getRegister(F)[0]
	go_gb.Set(register, bit, val)
}

func (c *cpu) getFlag(bit int) bool {
	return go_gb.Bit(c.getRegister(F)[0], bit)
}
