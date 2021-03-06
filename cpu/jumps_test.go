package cpu

import (
	go_gb "go-gb"
	"testing"
)

type jptest struct {
	in                     Instr
	prepare                func()
	expectedSp, expectedPc uint16
	expectedMC             go_gb.MC
}

func TestJp(t *testing.T) {
	c := initCpu(nil)
	c.pc = 0xC000
	start := c.pc
	table := []jptest{
		{jp(dx(16)), func() { c.memory.StoreBytes(c.pc, []byte{0xAB, 0xCD}) }, c.sp, 0xCDAB, 3},

		{jrnc(BitZ), func() { c.memory.StoreBytes(c.pc, []byte{0xAB}) }, c.sp, start + 0xAB + 1, 2},
		{jrnc(BitC), func() { c.memory.StoreBytes(c.pc, []byte{0xCD}) }, c.sp, start + 0xCD + 1, 2},

		{jrc(BitZ), func() { c.memory.StoreBytes(c.pc, []byte{0xAB}) }, c.sp, start, 1},
		{jrc(BitC), func() { c.memory.StoreBytes(c.pc, []byte{0xCD}) }, c.sp, start, 1},

		{jpnc(BitZ, dx(16)), func() { c.memory.StoreBytes(c.pc, []byte{0xFB, 0xFA}) }, c.sp, 0xFAFB, 3},
		{jpnc(BitC, dx(16)), func() { c.memory.StoreBytes(c.pc, []byte{0xFD, 0xFC}) }, c.sp, 0xFCFD, 3},

		{jpnc(BitZ, dx(16)), func() { c.setFlag(BitZ, true); c.memory.StoreBytes(c.pc, []byte{0xFB, 0xFA}) }, c.sp, start, 2},
		{jpnc(BitC, dx(16)), func() { c.setFlag(BitC, true); c.memory.StoreBytes(c.pc, []byte{0xFD, 0xFC}) }, c.sp, start, 2},

		{jpc(BitZ, dx(16)), func() { c.setFlag(BitZ, true); c.memory.StoreBytes(c.pc, []byte{0xFB, 0xFA}) }, c.sp, 0xFAFB, 3},
		{jpc(BitC, dx(16)), func() { c.setFlag(BitC, true); c.memory.StoreBytes(c.pc, []byte{0xFD, 0xFC}) }, c.sp, 0xFCFD, 3},

		{jpHl, func() { c.rMap[go_gb.HL][1] = 0xFC; c.rMap[go_gb.HL][0] = 0xFD }, c.sp, 0xFCFD, 0},

		{call, func() { c.memory.StoreBytes(c.pc, []byte{0x00, 0xDF}) }, c.sp - 2, 0xDF00, 5},
		{ret, nil, c.sp, 0xC002, 3},
	}
	for i, test := range table {
		c.pc = start
		if test.prepare != nil {
			test.prepare()
		}
		if mc := test.in(c); mc != test.expectedMC {
			t.Errorf("test %d expected MC %d, got %d\n", i+1, test.expectedMC, mc)
		}
		if c.pc != test.expectedPc {
			t.Errorf("test %d expected PC %X, got %X\n", i+1, test.expectedPc, c.pc)
		}
		if c.sp != test.expectedSp {
			t.Errorf("test %d expected SP %X, got %X\n", i+1, test.expectedSp, c.sp)
		}
	}
}
