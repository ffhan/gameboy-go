package cpu

import (
	"errors"
	go_gb "go-gb"
	"testing"
)

func TestReg_Load(t *testing.T) {
	c := NewCpu()
	input := byte(0xFA)
	c.r[A] = input
	a := rx(A)
	result, mc := a.Load(c)
	if len(result) != 1 {
		t.Errorf("expected len 1, got %d\n", len(result))
	}
	if result[0] != input {
		t.Errorf("expected result %X, got %X\n", input, result[0])
	}
	if mc != 0 {
		t.Error("MC should be 0")
	}
}

func TestReg_Store(t *testing.T) {
	c := NewCpu()
	input := byte(0xFA)
	rx(A).Store(c, []byte{input})
	if input != c.r[A] {
		t.Errorf("expected %X, got %X\n", input, c.r[A])
	}
}

func TestReg_Load16bit(t *testing.T) {
	c := NewCpu()
	input := []byte{0xFA, 0xCE}
	c.r[A] = input[0]
	c.r[F] = input[1]
	a := rx(AF)
	result := a.Load(c)
	if len(result) != 2 {
		t.Errorf("expected len 2, got %d\n", len(result))
	}
	for i, res := range result {
		if res != input[i] {
			t.Errorf("expected result %X, got %X\n", input[i], res)
		}
	}
}

func TestReg_Store16bit(t *testing.T) {
	c := NewCpu()
	input := []byte{0xFA, 0xCE}
	rx(AF).Store(c, input)
	for i, res := range c.rMap[AF] {
		if res != input[i] {
			t.Errorf("expected result %X, got %X\n", input[i], res)
		}
	}
}

func TestMPtr_Load(t *testing.T) {
	c := NewCpu()
	input := []byte{0x01, 0x02, 0xAB}
	addr := c.pc
	c.memory.StoreBytes(addr, input)
	expected := []byte{0xAB}
	bytes := md(16).Load(c)
	for i := range bytes {
		if bytes[i] != expected[i] {
			t.Errorf("expected %X, got %X\n", expected[i], bytes[i])
		}
	}
}

func TestMPtr_Store(t *testing.T) {
	c := NewCpu()
	input := []byte{0x01, 0x02, 0xBD}
	input = []byte{0xBD, 0xFE}
	md(16).Store(c, input)
	result := c.memory.ReadBytes(0, 2)
	for i := range result {
		if result[i] != input[i] {
			t.Errorf("expected %X, got %X\n", input[i], result[i])
		}
	}
}

func TestData_Store(t *testing.T) {
	defer func() {
		i := recover()
		if !errors.Is(i.(error), InvalidStoreErr) {
			panic(i.(error))
		}
		t.Log("invalid store has been called")
	}()
	dx(16).Store(NewCpu(), []byte{})
}

func TestData_Load(t *testing.T) {
	c := NewCpu()
	input := []byte{1, 2}
	c.memory.StoreBytes(c.pc, input)
	bytes := dx(16).Load(c)
	for i := range bytes {
		if bytes[i] != input[i] {
			t.Errorf("expected %X, got %X\n", input[i], bytes[i])
		}
	}
}

func TestHardcoded_Load(t *testing.T) {
	c := NewCpu()
	bytes := hc(4).Load(c)
	if len(bytes) != 1 {
		t.Fatal("invalid bytes length")
	}
	if bytes[0] != 4 {
		t.Errorf("expected 4, got %d\n", bytes[0])
	}
}

func TestHardcoded_Store(t *testing.T) {
	defer func() {
		err := recover().(error)
		if !errors.Is(err, InvalidStoreErr) {
			panic(err)
		}
		t.Log("invalid store has been called")
	}()
	hc(6).Store(NewCpu(), []byte{})
}

func TestMr(t *testing.T) {
	c := NewCpu()
	expected := byte(0xFE)
	c.r[B] = 0x01
	c.r[C] = 0x02
	c.memory.Store(0x0102, expected)
	bytes := mr(BC).Load(c)
	if len(bytes) != 1 {
		t.Errorf("expected len 1, got %d\n", len(bytes))
	}
	if bytes[0] != expected {
		t.Errorf("expected %X, got %X\n", expected, bytes[0])
	}
}

func TestOffset_Load_Reg(t *testing.T) { // offset doesn't work
	c := NewCpu()
	c.rMap[F][0] = 0xAB
	o := off(rx(AF), 0xFF00)
	expected := uint16(0xFFAB)
	result := go_gb.FromBytes(o.Load(c))
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
}

func TestOffset_Store_Reg(t *testing.T) { // offset doesn't work
	c := NewCpu()
	o := off(rx(A), 0xA0)
	expected := byte(0xAB)
	o.Store(c, []byte{0x0B})
	result := c.rMap[A][0]
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
}

func TestOffset_Load_Mptr8bit(t *testing.T) {
	c := NewCpu()
	c.rMap[A][0] = 0xAA
	expected := byte(0xF1)
	c.memory.Store(0xFFAA, expected)
	c.memory.Store(c.pc, 0xAA)
	o := mem(off(dx(8), 0xFF00))
	result := o.Load(c)[0]
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
}

func TestOffset_Load_Mptr16bit(t *testing.T) {
	c := NewCpu()
	c.rMap[A][0] = 0xAA
	expected := byte(0xF1)
	c.memory.Store(0xFFAA, expected)
	c.memory.StoreBytes(c.pc, []byte{0x00, 0xAA})
	o := mem(off(dx(16), 0xFF00))
	result := o.Load(c)[0]
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
}

func TestOffset_Store_Mptr8bit(t *testing.T) {
	c := NewCpu()
	expected := byte(0xAB)
	c.memory.Store(c.pc, 0xAB)
	o := mem(off(dx(8), 0xFF00))
	o.Store(c, []byte{expected})
	result := c.memory.Read(0xFFAB)
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
}

func TestOffset_Store_Mptr16bit(t *testing.T) {
	c := NewCpu()
	expected := byte(0xAB)
	c.memory.StoreBytes(c.pc, []byte{0x02, 0x34})
	o := mem(off(dx(16), 0xAAAA))
	o.Store(c, []byte{expected})
	result := c.memory.Read(0xACDE)
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
}
