package cpu

import (
	"go-gb"
	memory2 "go-gb/memory"
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

	memory go_gb.Memory

	halt      bool
	eiWaiting byte
	diWaiting byte
	ime       bool // Interrupt master enable
	cbLookup  bool

	storeFunctions map[uint16]func(*cpu, []byte)
}

func NewCpu() *cpu {
	return &cpu{}
}

func (c *cpu) readOpcode(mc *go_gb.MC) byte {
	val := c.memory.Read(c.pc)
	*mc += 1 // we purposefully don't check for nil in mc because it should always be cycle counted
	// discard the result if you want not to count cycles.
	c.pc += 1
	return val
}

func (c *cpu) readBytes(pointer, n uint16, mc *go_gb.MC) []byte {
	*mc += go_gb.MC(n)
	return c.memory.ReadBytes(pointer, n)
}

func (c *cpu) read(pointer uint16, mc *go_gb.MC) byte {
	*mc += 1
	return c.memory.Read(pointer)
}

func (c *cpu) storeBytes(pointer uint16, b []byte, mc *go_gb.MC) {
	*mc += go_gb.MC(len(b))
	if f, ok := c.storeFunctions[pointer]; ok {
		f(c, b)
	}
	c.memory.StoreBytes(pointer, b)
}

func (c *cpu) store(pointer uint16, val byte, mc *go_gb.MC) {
	c.storeBytes(pointer, []byte{val}, mc)
}

func (c *cpu) readFromPc(size uint16, mc *go_gb.MC) []byte {
	val := c.memory.ReadBytes(c.pc, size)
	*mc += go_gb.MC(size)
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
		v := c.memory.Read(c.sp)
		*mc += 1
		bytes[i] = v
	}
	return bytes
}

func (c *cpu) pushStack(b []byte, mc *go_gb.MC) {
	for _, val := range b {
		c.memory.Store(c.sp, val)
		*mc += 1
		c.sp -= 1
	}
}

func (c *cpu) getRegister(r registerName) []byte {
	return c.rMap[r]
}

func (c *cpu) Init(rom []byte, gbType go_gb.GameboyType) {
	mmu := memory2.NewMMU()
	mmu.Init(rom, gbType)
	c.memory = mmu

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
	c.storeFunctions = map[uint16]func(c *cpu, bytes []byte){
		go_gb.LCDDMA: dma,
	}
}

func dma(c *cpu, bytes []byte) {
	source := go_gb.FromBytes(bytes)
	result := c.memory.ReadBytes(source, 0x9F+1)
	c.memory.StoreBytes(memory2.OAMStart, result)
}

func (c *cpu) Step() go_gb.MC {
	var cycles go_gb.MC
	if !c.halt {
		opcode := c.readOpcode(&cycles)
		var instr Instr
		if c.cbLookup {
			instr = cbOptable[opcode]
			c.cbLookup = false
		} else {
			instr = optable[opcode]
		}
		cycles += instr(c)
	} else {
		cycles = 1
	}
	// todo: handle PPU BEFORE interrupts (hblank/vblank/stat interrupts)
	c.handleEiDi()
	cycles += c.handleInterrupts()
	return cycles
}

func (c *cpu) handleInterrupts() go_gb.MC { // todo: should we count the cycles from the memory read?
	var cycles go_gb.MC
	if !c.ime { // interrupt master disabled
		return cycles
	}
	ifRegister := c.memory.Read(go_gb.IF)
	ieRegister := c.memory.Read(go_gb.IE)

	cycles += 2

	if ifRegister == 0 { // no interrupt flags set
		return cycles
	}

	for _, interrupt := range go_gb.Interrupts {
		if go_gb.ShouldServiceInterrupt(ieRegister, ifRegister, interrupt.Bit) {
			cycles += c.serviceInterrupt(ifRegister, interrupt)
			return cycles
		}
	}
	return cycles
}

func (c *cpu) serviceInterrupt(ifR byte, interrupt go_gb.Interrupt) go_gb.MC {
	var cycles go_gb.MC
	go_gb.Set(&ifR, int(interrupt.Bit), false)
	c.ime = false
	c.memory.Store(go_gb.IF, ifR) // todo: should we update cycles during interrupts?
	cycles += 1
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
