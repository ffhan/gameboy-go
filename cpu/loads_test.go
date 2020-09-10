package cpu

import (
	go_gb "go-gb"
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

func checkReg(c *cpu, reg registerName) func() []byte {
	return func() []byte {
		return c.rMap[reg]
	}
}

func checkMr(c *cpu, reg registerName) func() []byte {
	return func() []byte {
		return c.memory.ReadBytes(go_gb.MsbLsb(c.rMap[reg]), 1)
	}
}

func checkSp(c *cpu) func() []byte {
	return func() []byte {
		return go_gb.MsbLsbBytes(c.sp, true)
	}
}

func checkMd(c *cpu, i int) func() []byte {
	offset := i / 8
	return func() []byte {
		addr := go_gb.MsbLsb(c.memory.ReadBytes(c.pc-2, uint16(offset)))
		return c.memory.ReadBytes(addr, 2)
	}
}

func TestLoad(t *testing.T) {
	c := NewCpu()
	table := []loadTest{
		{nil, checkReg(c, BC), load(rx(BC), dx(16)), []byte{0xFA, 0xCE}},
		{nil, checkReg(c, DE), load(rx(DE), dx(16)), []byte{0x12, 0x34}},
		{nil, checkReg(c, HL), load(rx(HL), dx(16)), []byte{0x56, 0x78}},
		{nil, checkSp(c), load(sp(), dx(16)), []byte{0x9A, 0xBC}},

		{func() { c.rMap[A][0] = 0xAE }, checkMr(c, BC), load(mr(BC), rx(A)), []byte{0xAE}},
		{func() { c.rMap[A][0] = 0xEA }, checkMr(c, DE), load(mr(DE), rx(A)), []byte{0xEA}},

		{nil, checkReg(c, B), load(rx(B), dx(8)), []byte{0x56}},
		{nil, checkReg(c, D), load(rx(D), dx(8)), []byte{0x57}},
		{nil, checkReg(c, H), load(rx(H), dx(8)), []byte{0x58}},
		{nil, checkMr(c, HL), load(mr(HL), dx(8)), []byte{0x59}},

		{func() { c.sp = 0xFFFE }, checkMd(c, 16), load(md(16), sp()), []byte{0xFF, 0xFE}},

		{func() { c.memory.StoreBytes(0, []byte{0x0A}); c.rMap[BC][0] = 0; c.rMap[BC][1] = 0 }, checkReg(c, A), load(rx(A), mr(BC)), []byte{0x0A}},
		{func() { c.memory.StoreBytes(0, []byte{0x1A}); c.rMap[DE][0] = 0; c.rMap[DE][1] = 0 }, checkReg(c, A), load(rx(A), mr(DE)), []byte{0x1A}},

		{nil, checkReg(c, C), load(rx(C), dx(8)), []byte{0x5A}},
		{nil, checkReg(c, E), load(rx(E), dx(8)), []byte{0x5B}},
		{nil, checkReg(c, L), load(rx(L), dx(8)), []byte{0x5C}},
		{nil, checkReg(c, A), load(rx(A), dx(8)), []byte{0x5D}},
	}
	c.memory.StoreBytes(c.pc, []byte{
		0xFA, 0xCE, 0x12, 0x34,
		0x56, 0x78, 0x9A, 0xBC,
		0x56, 0x57, 0x58, 0x59,
		0x00, 0x00, // load (nn), SP test moves PC
		0x5A, 0x5B, 0x5C, 0x5D,
	})
	for i, test := range table {
		if test.prepare != nil {
			test.prepare()
		}
		if err := test.in(c); err != nil {
			t.Error(err)
		}
		checkBytes(i+1, t, test.expected, test.results())
	}
}

func TestLoadHl(t *testing.T) {
	c := NewCpu()
	c.rMap[HL][0] = 0xAB
	c.rMap[HL][1] = 0xCD
	hlLSB := byte(0xCD)
	check := func(inc bool) func() []byte {
		return func() []byte {
			if inc {
				hlLSB += 1
			} else {
				hlLSB -= 1
			}
			if c.rMap[HL][0] != 0xAB || c.rMap[HL][1] != hlLSB {
				t.Errorf("expected HL to be %v, got %v\n", []byte{0xAB, hlLSB}, c.rMap[HL])
			}
			return c.rMap[A]
		}
	}
	c.rMap[A][0] = 0x69
	table := []loadTest{
		{nil, check(true), ldHl(nil, rx(A), true), []byte{0x69}},
		{nil, check(false), ldHl(nil, rx(A), false), []byte{0x69}},

		{func() { c.rMap[HL][1] = 0xFE; hlLSB = 0xFE }, check(true), ldHl(rx(A), nil, true), []byte{0x70}},
		{nil, check(false), ldHl(rx(A), nil, false), []byte{0x71}},
	}
	c.memory.StoreBytes(0xABFE, []byte{0x70, 0x71})
	for i, test := range table {
		if test.prepare != nil {
			test.prepare()
		}
		if err := test.in(c); err != nil {
			t.Error(err)
		}
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
	c := NewCpu()
	start := c.pc
	startSp := c.sp
	table := []spTest{
		{push(rx(BC)), func() { c.rMap[BC][0] = 0x12; c.rMap[BC][1] = 0x34 }, startSp - 2, []byte{0x34, 0x12}},
		{push(rx(DE)), func() { c.rMap[DE][0] = 0x12; c.rMap[DE][1] = 0x35 }, startSp - 2, []byte{0x35, 0x12}},
		{push(rx(HL)), func() { c.rMap[HL][0] = 0x12; c.rMap[HL][1] = 0x36 }, startSp - 2, []byte{0x36, 0x12}},
		{push(rx(AF)), func() { c.rMap[AF][0] = 0x12; c.rMap[AF][1] = 0x37 }, startSp - 2, []byte{0x37, 0x12}},
	}
	for i, test := range table {
		c.pc = start
		c.sp = startSp
		if test.prepare != nil {
			test.prepare()
		}
		if err := test.in(c); err != nil {
			t.Error(err)
		}
		if c.sp != test.expectedSp {
			t.Errorf("test %d expected SP %X, got %X\n", i+1, test.expectedSp, c.sp)
		}
		checkBytes(i+1, t, test.expected, c.popStack(2))
	}
}

func TestPop(t *testing.T) { // todo: test flags
	type poptest struct {
		register   registerName
		expectedSp uint16
		expected   []byte
	}
	c := NewCpu()
	c.pushStack([]byte{0xAB, 0xCD, 0xDE, 0xF0, 0xF1, 0xA1, 0xB2, 0xC4})
	startSp := c.sp
	table := []poptest{
		{BC, startSp + 2, []byte{0xC4, 0xB2}},
		{DE, startSp + 4, []byte{0xA1, 0xF1}},
		{HL, startSp + 6, []byte{0xF0, 0xDE}},
		{AF, startSp + 8, []byte{0xCD, 0xAB}},
	}
	for i, test := range table {
		if err := pop(rx(test.register))(c); err != nil {
			t.Error(err)
		}
		if c.sp != test.expectedSp {
			t.Errorf("test %d expected SP %X, got %X\n", i+1, test.expectedSp, c.sp)
		}
		checkBytes(i+1, t, test.expected, c.rMap[test.register])
	}
}

func TestLoadHlSp(t *testing.T) {
	c := NewCpu()
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
		if err := ldHlSp(c); err != nil {
			t.Error(err)
		}
		hl := go_gb.MsbLsb(c.rMap[HL])
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
