package cpu

import "testing"

type jptest struct {
	in                     Instr
	prepare                func()
	expectedSp, expectedPc uint16
}

func TestJp(t *testing.T) {
	c := NewCpu()
	start := c.pc
	table := []jptest{
		{jp(dx(16)), func() { c.memory.StoreBytes(c.pc, []byte{0xAB, 0xCD}) }, c.sp, 0xCDAB},

		{jrnc(BitZ), func() { c.memory.StoreBytes(c.pc, []byte{0xAB}) }, c.sp, start + 0xAB + 1},
		{jrnc(BitC), func() { c.memory.StoreBytes(c.pc, []byte{0xCD}) }, c.sp, start + 0xCD + 1},

		{jrc(BitZ), func() { c.memory.StoreBytes(c.pc, []byte{0xAB}) }, c.sp, start},
		{jrc(BitC), func() { c.memory.StoreBytes(c.pc, []byte{0xCD}) }, c.sp, start},

		{jpnc(BitZ, dx(16)), func() { c.memory.StoreBytes(c.pc, []byte{0xFB, 0xFA}) }, c.sp, 0xFAFB},
		{jpnc(BitC, dx(16)), func() { c.memory.StoreBytes(c.pc, []byte{0xFD, 0xFC}) }, c.sp, 0xFCFD},

		{jpnc(BitZ, dx(16)), func() { c.setFlag(BitZ, true); c.memory.StoreBytes(c.pc, []byte{0xFB, 0xFA}) }, c.sp, start},
		{jpnc(BitC, dx(16)), func() { c.setFlag(BitC, true); c.memory.StoreBytes(c.pc, []byte{0xFD, 0xFC}) }, c.sp, start},

		{jpc(BitZ, dx(16)), func() { c.setFlag(BitZ, true); c.memory.StoreBytes(c.pc, []byte{0xFB, 0xFA}) }, c.sp, 0xFAFB},
		{jpc(BitC, dx(16)), func() { c.setFlag(BitC, true); c.memory.StoreBytes(c.pc, []byte{0xFD, 0xFC}) }, c.sp, 0xFCFD},

		{jp(rx(HL)), func() { c.rMap[HL][0] = 0xFC; c.rMap[HL][1] = 0xFD }, c.sp, 0xFDFC},
	}
	for i, test := range table {
		c.pc = start
		if test.prepare != nil {
			test.prepare()
		}
		if err := test.in(c); err != nil {
			t.Error(err)
		}
		if c.pc != test.expectedPc {
			t.Errorf("test %d expected PC %X, got %X\n", i+1, test.expectedPc, c.pc)
		}
		if c.sp != test.expectedSp {
			t.Errorf("test %d expected SP %X, got %X\n", i+1, test.expectedSp, c.sp)
		}
	}
}
