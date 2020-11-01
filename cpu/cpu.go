package cpu

import (
	"go-gb"
	"reflect"
	"runtime"
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

// executes specific things on the cpu and returns the number of m cycles it took to execute
type Instr func(c *cpu) go_gb.MC

type timer interface {
	Step(mc go_gb.MC)
}

type cpu struct {
	pc, sp uint16

	r              [8]byte // registers (stored in little endian order)
	af, bc, de, hl []byte  // double registers (double registers store data in little endian)

	rMap [][]byte // register mappings

	bios   go_gb.Memory
	memory go_gb.MemoryBus
	hram   go_gb.Memory // used for direct stack access
	io     go_gb.Memory // used for direct IO access
	ier    go_gb.Memory // used for direct register access

	halt      bool
	stop      bool
	eiWaiting byte
	diWaiting byte
	ime       bool // Interrupt master enable
	dmaCycles go_gb.MC

	divTimer timer
	timer    timer

	serial go_gb.Serial

	ppu go_gb.PPU
}

func NewCpu(mmu go_gb.MemoryBus, ppu go_gb.PPU, timer timer, divTimer timer, serial go_gb.Serial) *cpu {
	c := &cpu{
		memory:   mmu,
		ppu:      ppu,
		hram:     mmu.HRAM(),
		io:       mmu.IO(),
		ier:      mmu.InterruptEnableRegister(),
		timer:    timer,
		divTimer: divTimer,
		serial:   serial,
	}
	c.init()
	return c
}

func (c *cpu) IME() bool {
	return c.ime
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
	c.memory.StoreBytes(pointer, b)
}

func (c *cpu) store(pointer uint16, val byte, mc *go_gb.MC) {
	*mc += go_gb.MC(1)
	c.memory.Store(pointer, val)
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
		v := c.memory.Read(c.sp)
		c.sp += 1
		*mc += 1
		bytes[i] = v
	}
	return bytes
}

func (c *cpu) pushStack(b []byte, mc *go_gb.MC) {
	for i := len(b) - 1; i >= 0; i-- {
		c.sp -= 1
		c.memory.Store(c.sp, b[i])
		*mc += 1
	}
}

func (c *cpu) GetRegister(r go_gb.RegisterName) []byte {
	return c.rMap[r]
}

func (c *cpu) init() {
	//c.pc = 0x0100
	//c.sp = 0xFFFE boot room fills this
	// todo: set r to init values
	// setting references to register arr
	//c.ime = true
	c.af = c.r[go_gb.F : go_gb.A+1]
	c.bc = c.r[go_gb.C : go_gb.B+1]
	c.de = c.r[go_gb.E : go_gb.D+1]
	c.hl = c.r[go_gb.L : go_gb.H+1]
	c.rMap = [][]byte{
		c.r[go_gb.F : go_gb.F+1], c.r[go_gb.A : go_gb.A+1], c.r[go_gb.C : go_gb.C+1], c.r[go_gb.B : go_gb.B+1],
		c.r[go_gb.E : go_gb.E+1], c.r[go_gb.D : go_gb.D+1], c.r[go_gb.L : go_gb.L+1], c.r[go_gb.H : go_gb.H+1],
		c.af, c.bc, c.de, c.hl,
	}
}

func (c *cpu) PC() uint16 {
	return c.pc
}

func (c *cpu) SP() uint16 {
	return c.sp
}

var instrs = map[string]bool{}
var opcodes = map[uint16]bool{}

//
//var cyc uint
//var every uint = 50
//
func (c *cpu) Step() go_gb.MC {
	var cycles go_gb.MC
	//if (c.pc == 0x1b05) && c.memory.Booted() {
	//	vramFile, err := os.Create("vram.txt")
	//	if err != nil {
	//		panic(err)
	//	}
	//	oamFile, err := os.Create("oam.txt")
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer vramFile.Close()
	//	defer oamFile.Close()
	//	c.memory.VRAM().Dump(vramFile)
	//	memory.DumpOam(c.memory.OAM(), c.memory.VRAM(), oamFile)
	//}
	//if c.memory.Booted() {
	//	cyc += 1
	//}
	//if cyc == every+1 {
	//	cyc = 0
	//}
	//if c.pc == 0x100 {
	//	print()
	//}
	if !c.halt && !c.stop {
		opcode := c.readOpcode(&cycles)
		var instr Instr
		if opcode == 0xCB {
			opcode = c.readOpcode(&cycles)
			opcodes[uint16(opcode)|0xCB00] = true
			instr = cbOptable[opcode]
		} else {
			instr = optable[opcode]
			opcodes[uint16(opcode)] = true
		}
		name := runtime.FuncForPC(reflect.ValueOf(instr).Pointer()).Name()
		if _, ok := instrs[name]; !ok {
			instrs[name] = true
		}

		/* Padamo na:
		; Set LCD control to Operation
		        ld      a,80h           ; 02b6 3e 80   >.
		        ldh     (40h),a         ; 02b8 e0 40   `@
		        ei                      ; 02ba fb   {
		        xor     a               ; 02bb af   /
		; Clear all interrupt flags
		        ldh     (0fh),a         ; 02bc e0 0f   `.
		Krene clearati inteerruptove na IF-u ali se VBLANK već izvede...
		*/
		cycles += instr(c)
	} else {
		cycles = 1
	}
	if c.memory.DMAInProgress() {
		c.dmaCycles += cycles
		if c.dmaCycles >= 40 {
			c.dmaCycles = 0
			c.memory.SetDMAInProgress(false)
		}
	}
	c.timer.Step(cycles)
	c.divTimer.Step(cycles)

	c.serial.Step(cycles)

	if c.ppu.Enabled() {
		c.ppu.Step(cycles)
	}
	c.handleEiDi()
	cycles += c.handleInterrupts()
	return cycles
}

func (c *cpu) handleInterrupts() go_gb.MC { // todo: should we count the cycles from the memory read?
	var cycles go_gb.MC
	if !c.ime { // interrupt master disabled
		return cycles
	}
	ifRegister := c.io.Read(go_gb.IF)
	ieRegister := c.ier.Read(go_gb.IE)

	//cycles += 2

	if ifRegister == 0 { // no interrupt flags set
		return cycles
	}

	for _, interrupt := range go_gb.Interrupts {
		if go_gb.ShouldServiceInterrupt(ieRegister, ifRegister, interrupt.Bit) {
			cycles += c.serviceInterrupt(ifRegister, interrupt)
			c.halt = false
			return cycles
		}
	}
	return cycles
}

func (c *cpu) serviceInterrupt(ifR byte, interrupt go_gb.Interrupt) go_gb.MC {
	var cycles go_gb.MC
	go_gb.Set(&ifR, int(interrupt.Bit), false)
	c.ime = false
	c.io.Store(go_gb.IF, ifR)
	callAddr(c, go_gb.ToBytes(interrupt.JpAddr, true), &cycles)
	if interrupt.Bit == go_gb.BitJoypad {
		c.stop = false // joypad interrupt removed stop
	}
	return 5
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
	register := &c.GetRegister(go_gb.F)[0]
	go_gb.Set(register, bit, val)
}

func (c *cpu) getFlag(bit int) bool {
	return go_gb.Bit(c.GetRegister(go_gb.F)[0], bit)
}
