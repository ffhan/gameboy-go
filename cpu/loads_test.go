package cpu

import (
	go_gb "go-gb"
	"go-gb/memory"
	"testing"
)

type loadTest struct {
	prepare  func()
	results  func() []byte
	in       Instr
	expected []byte
}

func checkBytes(num int, t *testing.T, expected, results []byte) {
	if len(expected) != len(results) {
		t.Errorf("test %d expected and results sizes dont match\n", num)
	}
	for i := range expected {
		if expected[i] != results[i] {
			t.Errorf("test %d expected %X, got %X\n", num, expected[i], results[i])
		}
	}
}

func checkReg(c *cpu, reg go_gb.RegisterName) func() []byte {
	return func() []byte {
		return c.rMap[reg]
	}
}

func checkMr(c *cpu, reg go_gb.RegisterName) func() []byte {
	return func() []byte {
		return c.memory.ReadBytes(go_gb.FromBytes(c.rMap[reg]), 1)
	}
}

func checkSp(c *cpu) func() []byte {
	return func() []byte {
		return go_gb.ToBytes(c.sp, true)
	}
}

func checkMd(c *cpu, i int) func() []byte {
	offset := i / 8
	return func() []byte {
		addr := go_gb.FromBytes(c.memory.ReadBytes(c.pc-2, uint16(offset)))
		return c.memory.ReadBytes(addr, 2)
	}
}

func TestLoad(t *testing.T) {
	c := initCpu(nil)
	c.pc = memory.WRAMBank0Start
	table := []loadTest{
		{nil, checkReg(c, go_gb.BC), load(rx(go_gb.BC), dx(16)), []byte{0xFA, 0xCE}},
		{nil, checkReg(c, go_gb.DE), load(rx(go_gb.DE), dx(16)), []byte{0xFE, 0xCE}},
		{nil, checkReg(c, go_gb.HL), load(rx(go_gb.HL), dx(16)), []byte{0xCE, 0xFF}},
		{nil, checkSp(c), load(sp(), dx(16)), []byte{0x9A, 0xBC}},

		{func() { c.rMap[go_gb.A][0] = 0xAE }, checkMr(c, go_gb.BC), load(mr(go_gb.BC), rx(go_gb.A)), []byte{0xAE}},
		{func() { c.rMap[go_gb.A][0] = 0xEA }, checkMr(c, go_gb.DE), load(mr(go_gb.DE), rx(go_gb.A)), []byte{0xEA}},

		{nil, checkReg(c, go_gb.B), load(rx(go_gb.B), dx(8)), []byte{0x56}},
		{nil, checkReg(c, go_gb.D), load(rx(go_gb.D), dx(8)), []byte{0x57}},
		{nil, checkReg(c, go_gb.H), load(rx(go_gb.H), dx(8)), []byte{0xAA}},
		{func() { c.sp = 0xFFFE }, checkMd(c, 16), load(md(16), sp()), []byte{0xFE, 0xFF}},
		{nil, checkMr(c, go_gb.HL), load(mr(go_gb.HL), dx(8)), []byte{0xAB}},

		{func() {
			c.memory.StoreBytes(0xFBFA, []byte{0x0A})
			c.rMap[go_gb.BC][0] = 0xFA
			c.rMap[go_gb.BC][1] = 0xFB
		}, checkReg(c, go_gb.A), load(rx(go_gb.A), mr(go_gb.BC)), []byte{0x0A}},
		{func() {
			c.memory.StoreBytes(0xFCFB, []byte{0x1A})
			c.rMap[go_gb.DE][0] = 0xFB
			c.rMap[go_gb.DE][1] = 0xFC
		}, checkReg(c, go_gb.A), load(rx(go_gb.A), mr(go_gb.DE)), []byte{0x1A}},

		{nil, checkReg(c, go_gb.C), load(rx(go_gb.C), dx(8)), []byte{0x5A}},
		{nil, checkReg(c, go_gb.E), load(rx(go_gb.E), dx(8)), []byte{0x5B}},
		{nil, checkReg(c, go_gb.L), load(rx(go_gb.L), dx(8)), []byte{0x5C}},
		{nil, checkReg(c, go_gb.A), load(rx(go_gb.A), dx(8)), []byte{0x5D}},
		{func() { c.sp = 0x1200 }, checkReg(c, go_gb.HL), ldHlSp, []byte{0x34, 0x12}},
		{func() {
			c.r[go_gb.A] = 0xF1
			c.r[go_gb.H] = 0xC0
			c.r[go_gb.L] = 0x00
		}, checkReg(c, go_gb.HL), ldHl(nil, rx(go_gb.A), true), []byte{0x01, 0xC0}},
		{func() {
			c.r[go_gb.A] = 0xF1
			c.r[go_gb.H] = 0xC0
			c.r[go_gb.L] = 0x01
		}, checkReg(c, go_gb.HL), ldHl(nil, rx(go_gb.A), false), []byte{0x00, 0xC0}},
		{func() {
			c.r[go_gb.A] = 0xF1
			c.r[go_gb.H] = 0xC0
			c.r[go_gb.L] = 0x00
		}, checkReg(c, go_gb.HL), ldHl(rx(go_gb.A), nil, true), []byte{0x01, 0xC0}},
		{func() {
			c.r[go_gb.A] = 0xF1
			c.r[go_gb.H] = 0xC0
			c.r[go_gb.L] = 0x01
		}, checkReg(c, go_gb.HL), ldHl(rx(go_gb.A), nil, false), []byte{0x00, 0xC0}},
	}
	c.memory.StoreBytes(c.pc, []byte{
		0xFA, 0xCE, 0xFE, 0xCE,
		0xCE, 0xFF, 0x9A, 0xBC,
		0x56, 0x57, 0xAA, 0xAB,
		0xCD, 0xAB, // load (nn), SP test moves PC
		0x5A, 0x5B, 0x5C, 0x5D,
		0x34,
	})
	for i, test := range table {
		if test.prepare != nil {
			test.prepare()
		}
		test.in(c)
		checkBytes(i+1, t, test.expected, test.results())
		t.Logf("test %d completed", i+1)
	}
}

func TestLoadHl(t *testing.T) {
	c := initCpu(nil)
	c.rMap[go_gb.HL][1] = 0xAB
	c.rMap[go_gb.HL][0] = 0xCD
	hlLSB := byte(0xCD)
	check := func(inc bool) func() []byte {
		return func() []byte {
			if inc {
				hlLSB += 1
			} else {
				hlLSB -= 1
			}
			if c.rMap[go_gb.HL][1] != 0xAB || c.rMap[go_gb.HL][0] != hlLSB {
				t.Errorf("expected HL to be %v, got %v\n", []byte{0xAB, hlLSB}, c.rMap[go_gb.HL])
			}
			return c.rMap[go_gb.A]
		}
	}
	c.rMap[go_gb.A][0] = 0x69
	table := []loadTest{
		{nil, check(true), ldHl(nil, rx(go_gb.A), true), []byte{0x69}},
		{nil, check(false), ldHl(nil, rx(go_gb.A), false), []byte{0x69}},

		{func() { c.rMap[go_gb.HL][0] = 0xFE; hlLSB = 0xFE }, check(true), ldHl(rx(go_gb.A), nil, true), []byte{0x70}},
		{nil, check(false), ldHl(rx(go_gb.A), nil, false), []byte{0x71}},
	}
	c.memory.StoreBytes(0xABFE, []byte{0x70, 0x71})
	for i, test := range table {
		if test.prepare != nil {
			test.prepare()
		}
		test.in(c)
		checkBytes(i+1, t, test.expected, test.results())
	}
}

type spTest struct {
	in         Instr
	prepare    func()
	expectedSp uint16
	expected   []byte
}

func TestPush(t *testing.T) {
	c := initCpu(nil)
	start := c.pc
	startSp := c.sp
	table := []spTest{
		{push(rx(go_gb.BC)), func() { c.rMap[go_gb.BC][0] = 0x12; c.rMap[go_gb.BC][1] = 0x34 }, startSp - 2, []byte{0x34, 0x12}},
		{push(rx(go_gb.DE)), func() { c.rMap[go_gb.DE][0] = 0x12; c.rMap[go_gb.DE][1] = 0x35 }, startSp - 2, []byte{0x35, 0x12}},
		{push(rx(go_gb.HL)), func() { c.rMap[go_gb.HL][0] = 0x12; c.rMap[go_gb.HL][1] = 0x36 }, startSp - 2, []byte{0x36, 0x12}},
		{push(rx(go_gb.AF)), func() { c.rMap[go_gb.AF][0] = 0x12; c.rMap[go_gb.AF][1] = 0x37 }, startSp - 2, []byte{0x37, 0x12}},
	}
	for i, test := range table {
		c.pc = start
		c.sp = startSp
		if test.prepare != nil {
			test.prepare()
		}
		test.in(c)
		if c.sp != test.expectedSp {
			t.Errorf("test %d expected SP %X, got %X\n", i+1, test.expectedSp, c.sp)
		}
		var mc go_gb.MC
		checkBytes(i+1, t, test.expected, c.popStack(2, &mc))
		if mc != 2 {
			t.Errorf("expected 2 cycles, got %d\n", mc)
		}
	}
}

func TestPop(t *testing.T) { // todo: test flags
	type poptest struct {
		register   go_gb.RegisterName
		expectedSp uint16
		expected   []byte
	}
	c := initCpu(nil)
	var mc go_gb.MC
	in := []byte{0xAB, 0xCD, 0xDE, 0xF0, 0xF1, 0xA1, 0xB2, 0xC4}
	c.pushStack(in, &mc)
	if mc != go_gb.MC(len(in)) {
		t.Errorf("expected %d cycles, got %d\n", len(in), mc)
	}
	startSp := c.sp
	table := []poptest{
		{go_gb.BC, startSp + 2, []byte{0xC4, 0xB2}},
		{go_gb.DE, startSp + 4, []byte{0xA1, 0xF1}},
		{go_gb.HL, startSp + 6, []byte{0xF0, 0xDE}},
		{go_gb.AF, startSp + 8, []byte{0xCD, 0xAB}},
	}
	for i, test := range table {
		pop(rx(test.register))(c)
		if c.sp != test.expectedSp {
			t.Errorf("test %d expected SP %X, got %X\n", i+1, test.expectedSp, c.sp)
		}
		checkBytes(i+1, t, test.expected, c.rMap[test.register])
	}
}

func TestLoadHlSp(t *testing.T) {
	c := initCpu(nil)
	c.pc = 0xA000
	type hlsptest struct {
		prepare    func()
		z, n, h, c bool
		expected   uint16
	}
	table := []hlsptest{
		{func() { c.memory.Store(c.pc, 0xAB); c.sp = 0xFF00 }, false, false, false, false, 0xFFAB},
		{func() { c.memory.Store(c.pc, 0x01); c.sp = 0xFF0F }, false, false, true, false, 0xFF10},
		{func() { c.memory.Store(c.pc, 0x1B); c.sp = 0x00F0 }, false, false, false, true, 0x010B},
		{func() { c.memory.Store(c.pc, 0xFF); c.sp = 0x00FF }, false, false, true, true, 0x01FE},
	}
	for i, test := range table {
		test.prepare()
		ldHlSp(c)
		hl := go_gb.FromBytes(c.rMap[go_gb.HL])
		if hl != test.expected {
			t.Errorf("test %d expected %X, got %X\n", i+1, test.expected, hl)
		}
		if c.getFlag(BitZ) != test.z {
			t.Errorf("test %d Flag Z is wrong", i+1)
		}
		if c.getFlag(BitN) != test.n {
			t.Errorf("test %d Flag N is wrong", i+1)
		}
		if c.getFlag(BitH) != test.h {
			t.Errorf("test %d Flag H is wrong", i+1)
		}
		if c.getFlag(BitC) != test.c {
			t.Errorf("test %d Flag C is wrong", i+1)
		}
	}
}
