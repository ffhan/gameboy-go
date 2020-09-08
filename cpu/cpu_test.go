package cpu

import "testing"

func TestCPURegister(t *testing.T) {
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

	if c.readOpcode() != 0x12 {
		t.Fatal("invalid opcode read")
	}
	if c.pc != 0x101 {
		t.Fatal("PC invalid after reading opcode")
	}
}
