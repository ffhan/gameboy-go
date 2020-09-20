package cpu

import (
	go_gb "go-gb"
	"go-gb/memory"
	"testing"
)

func initCpu(fill map[uint16]byte) *cpu {
	mmu := memory.NewMMU()
	c := NewCpu(mmu, nil)
	bytes := make([]byte, 0xFFFF+1)
	if fill != nil {
		for addr, val := range fill {
			bytes[addr] = val
		}
	}
	bytes[memory.CartridgeTypeAddr] = byte(memory.MbcROMRAM)
	bytes[memory.CartridgeROMSizeAddr] = 0x05
	bytes[memory.CartridgeRAMSizeAddr] = 0x03
	mmu.Init(bytes, go_gb.GB)
	return c
}

func TestCPU_doubleRegister(t *testing.T) {
	c := initCpu(nil)
	c.r[A] = 0xFA
	c.r[F] = 0x12

	if c.rMap[AF][1] != 0xFA || c.rMap[AF][0] != 0x12 {
		t.Fatal("invalid double registers")
	}
}

func TestCPU_doubleRegister_changeSingleRegister(t *testing.T) {
	c := initCpu(nil)
	c.r[A] = 0xFA
	c.r[F] = 0x12
	c.rMap[AF][1] = 0xAA

	if c.rMap[AF][1] != 0xAA || c.rMap[AF][0] != 0x12 {
		t.Fatal("invalid double registers")
	}
}

func TestCpu_readOpcode(t *testing.T) {
	c := initCpu(map[uint16]byte{0x0: 0x12})

	var cycles go_gb.MC

	if val := c.readOpcode(&cycles); val != 0x12 {
		t.Fatal("invalid opcode read")
	}
	if cycles != 1 {
		t.Errorf("expectd 1 cycle, got %d\n", cycles)
	}
	if c.pc != 1 {
		t.Fatal("PC invalid after reading opcode")
	}
}

func TestCpu_pushStack(t *testing.T) {
	c := initCpu(nil)
	input := []byte{1, 2, 3}
	var mc go_gb.MC
	c.pushStack(input, &mc)

	if mc != 3 {
		t.Errorf("expected 3 cycles, got %d\n", mc)
	}
	expected := uint16(0xFFFE - 3)
	if c.sp != expected {
		t.Errorf("expected SP to be on %X, got %X\n", expected, c.sp)
	}
	bytes := c.memory.ReadBytes(expected+1, 3)
	for i, val := range bytes {
		expected := input[2-i]
		if expected != val {
			t.Errorf("expected %d, got %d\n", expected, val)
		}
	}
}

func TestCpu_popStack(t *testing.T) {
	c := initCpu(nil)
	input := []byte{1, 2, 3}
	var mc go_gb.MC
	c.pushStack(input, &mc)

	if mc != 3 {
		t.Errorf("expected 3 cycles, got %d\n", mc)
	}
	initialSP := uint16(0xFFFE)
	expected := initialSP - 3
	if c.sp != expected {
		t.Errorf("expected SP to be on %X, got %X\n", expected, c.sp)
	}
	mc = go_gb.MC(0)
	stack := c.popStack(3, &mc)
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
	c := initCpu(nil)
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
	c := initCpu(nil)
	c.setFlag(BitZ, true)
	c.setFlag(BitN, true)
	c.setFlag(BitH, true)
	c.setFlag(BitC, true)

	if c.r[F] != 0xF0 {
		t.Errorf("expected %X, got %X\n", c.r[F], 0xF0)
	}
}

func TestCpu_Step_NOP(t *testing.T) {
	c := initCpu(nil)
	startPC := c.pc
	c.memory.Store(c.pc, 0)
	mc := c.Step()
	if mc != 1 {
		t.Errorf("expected 1 cycle, got %d\n", mc)
	}
	if c.pc != startPC+1 {
		t.Errorf("expected PC %X, got %X\n", startPC+1, c.pc)
	}
}
