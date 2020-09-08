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

type Instr func(c *cpu) error

type cpu struct {
	pc, sp uint16

	r              [8]byte // registers
	af, bc, de, hl []byte  // double registers

	rMap [][]byte // register mappings

	memory go_gb.MemoryBus

	instrQueue instrQueue
	ime        bool // Interrupt master enable
}

func (c *cpu) readOpcode() byte {
	val := c.memory.Read(c.pc)
	c.pc += 1
	return val
}

func (c *cpu) popStack(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		c.sp += 1
		bytes[i] = c.memory.Read(c.sp)
	}
	return bytes
}

func (c *cpu) pushStack(b []byte) {
	for _, val := range b {
		c.memory.Store(c.sp, val)
		c.sp -= 1
	}
}

func (c *cpu) getRegister(r registerName) []byte {
	return c.rMap[r]
}

func (c *cpu) init() {
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
