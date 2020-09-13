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
	F registerName = iota
	A
	C
	B
	E
	D
	L
	H
	AF
	BC
	DE
	HL
)

// executes specific things on the cpu and returns the number of m cycles it took to execute
type Instr func(c *cpu) go_gb.MC

type cpu struct {
	pc, sp uint16

	r              [8]byte // registers (stored in little endian order)
	af, bc, de, hl []byte  // double registers (double registers store data in little endian)

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

func (c *cpu) readOpcode(mc *go_gb.MC) byte {
	val := c.memory.Read(c.pc, mc)
	c.pc += 1
	return val
}
func (c *cpu) readFromPc(size uint16, mc *go_gb.MC) []byte {
	val := c.memory.ReadBytes(c.pc, size, mc)
	c.pc += size
	return val
}

func (c *cpu) setPc(val uint16, mc *go_gb.MC) {
	c.pc = val
	if mc != nil {
		*mc += 1
	}
}

func (c *cpu) popStack(size int, mc *go_gb.MC) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		c.sp += 1
		v := c.memory.Read(c.sp, mc)
		bytes[i] = v
	}
	return bytes
}

func (c *cpu) pushStack(b []byte, mc *go_gb.MC) {
	for _, val := range b {
		c.memory.Store(c.sp, val, mc)
		c.sp -= 1
	}
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
	c.af = c.r[F : A+1]
	c.bc = c.r[C : B+1]
	c.de = c.r[E : D+1]
	c.hl = c.r[L : H+1]
	c.rMap = [][]byte{
		c.r[F : F+1], c.r[A : A+1], c.r[C : C+1], c.r[B : B+1],
		c.r[E : E+1], c.r[D : D+1], c.r[L : L+1], c.r[H : H+1],
		c.af, c.bc, c.de, c.hl,
	}
}

func (c *cpu) Step() {
	var cycles go_gb.MC
	if !c.halt {
		opcode := c.readOpcode(&cycles)
		instr := optable[opcode]
		cycles += instr(c)
	} else {
		cycles = 1
	}
	c.handleEiDi()
	c.handleInterrupts()
}

func (c *cpu) handleInterrupts() { // todo: should we count the cycles from the memory read?
	if !c.ime { // interrupt master disabled
		return
	}
	var cycles go_gb.MC
	ifRegister := c.memory.Read(go_gb.IF, &cycles)
	ieRegister := c.memory.Read(go_gb.IE, &cycles)

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
	c.memory.Store(go_gb.IF, ifR, &cycles)
	callAddr(c, go_gb.ToBytesReverse(interrupt.JpAddr, true), &cycles)
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
