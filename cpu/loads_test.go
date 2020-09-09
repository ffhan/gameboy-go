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
		return go_gb.MsbLsbBytes(c.sp)
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
