package cpu

import "testing"

func TestCPU_doubleRegister(t *testing.T) {
	c := NewCpu()
	c.r[A] = 0xFA
	c.r[F] = 0x12

	if c.rMap[AF][0] != 0xFA || c.rMap[AF][1] != 0x12 {
		t.Fatal("invalid double registers")
	}
}

func TestCpu_readOpcode(t *testing.T) {
	c := NewCpu()
	c.memory.StoreBytes(0x100, []byte{0x12})

	if val, mc := c.readOpcode(); val != 0x12 || mc != 1 {
		t.Fatal("invalid opcode read or MC is not 1")
	}
	if c.pc != 0x101 {
		t.Fatal("PC invalid after reading opcode")
	}
}

func TestCpu_pushStack(t *testing.T) {
	c := NewCpu()
	input := []byte{1, 2, 3}
	c.pushStack(input)
	expected := uint16(0xFFFE - 3)
	if c.sp != expected {
		t.Errorf("expected SP to be on %X, got %X\n", expected, c.sp)
	}
	bytes, mc := c.memory.ReadBytes(expected+1, 3)
	if mc != 3 {
		t.Errorf("expected MC %d, got %d\n", 3, mc)
	}
	for i, val := range bytes {
		expected := input[2-i]
		if expected != val {
			t.Errorf("expected %d, got %d\n", expected, val)
		}
	}
}

func TestCpu_popStack(t *testing.T) {
	c := NewCpu()
	input := []byte{1, 2, 3}
	c.pushStack(input)
	initialSP := uint16(0xFFFE)
	expected := initialSP - 3
	if c.sp != expected {
		t.Errorf("expected SP to be on %X, got %X\n", expected, c.sp)
	}
	stack, mc := c.popStack(3)
	if mc != 3 {
		t.Errorf("expected MC %d, got %d\n", 3, mc)
	}
	for i, val := range stack {
		expected := input[2-i]
		if expected != val {
			t.Errorf("expected %d, got %d\n", expected, val)
		}
	}
	if c.sp != initialSP {
		t.Errorf("expected SP to be on %X, got %X\n", initialSP, c.sp)
	}
}

func TestCpu_getFlag(t *testing.T) {
	c := NewCpu()
	c.r[F] |= 0xA0 // 1010
	if !c.getFlag(BitZ) {
		t.Error("Bit Z should be set")
	}
	if c.getFlag(BitN) {
		t.Error("Bit N should not be set")
	}
	if !c.getFlag(BitH) {
		t.Error("Bit H should be set")
	}
	if c.getFlag(BitC) {
		t.Error("Bit C should not be set")
	}
}

func TestCpu_setFlag(t *testing.T) {
	c := NewCpu()
	c.setFlag(BitZ, true)
	c.setFlag(BitN, true)
	c.setFlag(BitH, true)
	c.setFlag(BitC, true)

	if c.r[F] != 0xF0 {
		t.Errorf("expected %X, got %X\n", c.r[F], 0xF0)
	}
}

func TestCpu_Step_NOP(t *testing.T) {
	c := NewCpu()
	startPC := c.pc
	c.memory.Store(c.pc, 0)
	c.Step()
	if c.pc != startPC+1 {
		t.Errorf("expected PC %X, got %X\n", startPC+1, c.pc)
	}
}
