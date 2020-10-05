package cpu

import (
	"errors"
	go_gb "go-gb"
	"testing"
)

func TestReg_Load(t *testing.T) {
	c := initCpu(nil)
	input := byte(0xFA)
	c.r[A] = input
	a := rx(A)
	var mc go_gb.MC
	result := a.Load(c, &mc)
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
	var mc go_gb.MC
	c := initCpu(nil)
	input := byte(0xFA)
	rx(A).Store(c, []byte{input}, &mc)
	if input != c.r[A] {
		t.Errorf("expected %X, got %X\n", input, c.r[A])
	}
	if mc != 0 {
		t.Error("MC should be 0")
	}
}

func TestReg_Load16bit(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	input := []byte{0xFA, 0xCE}
	c.r[A] = input[1]
	c.r[F] = input[0]
	a := rx(AF)
	result := a.Load(c, &mc)
	if len(result) != 2 {
		t.Errorf("expected len 2, got %d\n", len(result))
	}
	for i, res := range result {
		if res != input[i] {
			t.Errorf("expected result %X, got %X\n", input[i], res)
		}
	}
	if mc != 0 {
		t.Error("MC should be 0")
	}
}

func TestReg_Store16bit(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	input := []byte{0xFA, 0xCE}
	rx(AF).Store(c, input, &mc)
	for i, res := range c.rMap[AF] {
		if res != input[i] {
			t.Errorf("expected result %X, got %X\n", input[i], res)
		}
	}
	if mc != 0 {
		t.Error("MC should be 0")
	}
}

func TestMPtr_Load(t *testing.T) {
	var mc go_gb.MC
	c := initCpu(nil)
	c.pc = 0xA000
	input := []byte{0x02, 0xA0, 0xAB}
	addr := c.pc
	c.memory.StoreBytes(addr, input)
	expected := []byte{0xAB}
	bytes := md(16).Load(c, &mc)
	for i := range bytes {
		if bytes[i] != expected[i] {
			t.Errorf("expected %X, got %X\n", expected[i], bytes[i])
		}
	}
	expectedMC := go_gb.MC(3)
	if mc != expectedMC {
		t.Errorf("expected MC %d, got %d\n", expectedMC, mc)
	}
}

func TestMPtr_Store(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	c.pc = 0xA000
	c.memory.StoreBytes(c.pc, []byte{0x00, 0xA0})
	input := []byte{0xFE, 0xBD}
	md(16).Store(c, input, &mc)
	result := c.memory.ReadBytes(0xA000, 2)
	for i := range result {
		if result[i] != input[i] {
			t.Errorf("expected %X, got %X\n", input[i], result[i])
		}
	}
	expected := go_gb.MC(4)
	if mc != expected {
		t.Errorf("expected MC %d, got %d\n", expected, mc)
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
	dx(16).Store(initCpu(nil), []byte{}, nil)
}

func TestData_Load(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	input := []byte{1, 2}
	c.pc = 0xA000
	c.memory.StoreBytes(c.pc, input)
	bytes := dx(16).Load(c, &mc)
	for i := range bytes {
		if bytes[i] != input[i] {
			t.Errorf("expected %X, got %X\n", input[i], bytes[i])
		}
	}
	expected := go_gb.MC(2)
	if mc != expected {
		t.Errorf("expected MC %d, got %d\n", expected, mc)
	}
}

func TestHardcoded_Load(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	bytes := hc(4).Load(c, &mc)
	if len(bytes) != 1 {
		t.Fatal("invalid bytes length")
	}
	if bytes[0] != 4 {
		t.Errorf("expected 4, got %d\n", bytes[0])
	}
	expected := go_gb.MC(0)
	if mc != expected {
		t.Errorf("expected MC %d, got %d\n", expected, mc)
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
	hc(6).Store(initCpu(nil), []byte{}, nil)
}

func TestMr(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	expected := byte(0xFE)
	c.r[B] = 0xFF
	c.r[C] = 0x00
	c.memory.Store(0xFF00, expected)
	bytes := mr(BC).Load(c, &mc)
	if len(bytes) != 1 {
		t.Errorf("expected len 1, got %d\n", len(bytes))
	}
	if bytes[0] != expected {
		t.Errorf("expected %X, got %X\n", expected, bytes[0])
	}
	expectedMC := go_gb.MC(1)
	if mc != expectedMC {
		t.Errorf("expected MC %d, got %d\n", expectedMC, mc)
	}
}

func TestOffset_Load_Reg(t *testing.T) { // offset doesn't work
	var mc go_gb.MC

	c := initCpu(nil)
	c.rMap[F][0] = 0xAB
	o := off(rx(AF))
	expected := uint16(0xFFAB)
	result := go_gb.FromBytes(o.Load(c, &mc))
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
	expectedMc := go_gb.MC(0)
	if mc != expectedMc {
		t.Errorf("expected MC %d, got %d\n", expectedMc, mc)
	}
}

func TestOffset_Store_Reg(t *testing.T) { // offset doesn't work
	var mc go_gb.MC

	c := initCpu(nil)
	o := off(rx(A))
	expected := byte(0xAB)
	o.Store(c, []byte{0x0B}, &mc)
	result := c.rMap[A][0]
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
	expectedMc := go_gb.MC(0)
	if mc != expectedMc {
		t.Errorf("expected MC %d, got %d\n", expectedMc, mc)
	}
}

func TestOffset_Load_Mptr8bit(t *testing.T) {
	var mc go_gb.MC
	c := initCpu(nil)
	c.pc = 0xA000
	c.r[A] = 0xAA
	expected := byte(0xF1)
	c.memory.Store(0xFFAA, expected)
	c.memory.Store(c.pc, 0xAA)
	o := mem(off(dx(8)))
	result := o.Load(c, &mc)[0]
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
	expectedMc := go_gb.MC(2)
	if mc != expectedMc {
		t.Errorf("expected MC %d, got %d\n", expectedMc, mc)
	}
}

func TestOffset_Load_Mptr16bit(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	c.pc = 0xA000
	c.r[A] = 0xAA
	expected := byte(0xF1)
	c.memory.Store(0xFFAA, expected)
	c.memory.StoreBytes(c.pc, []byte{0xAA})
	o := mem(off(dx(16)))
	result := o.Load(c, &mc)[0]
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
	expectedMc := go_gb.MC(3)
	if mc != expectedMc {
		t.Errorf("expected MC %d, got %d\n", expectedMc, mc)
	}
}

func TestOffset_Store_Mptr8bit(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	c.pc = 0xA000
	expected := byte(0xAB)
	c.memory.Store(c.pc, 0xAB)
	o := mem(off(dx(8)))
	o.Store(c, []byte{expected}, &mc)
	result := c.memory.Read(0xFFAB)
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
	expectedMc := go_gb.MC(2)
	if mc != expectedMc {
		t.Errorf("expected MC %d, got %d\n", expectedMc, mc)
	}
}

func TestOffset_Store_Mptr16bit(t *testing.T) {
	var mc go_gb.MC

	c := initCpu(nil)
	c.pc = 0xA000
	expected := byte(0xAB)
	c.memory.StoreBytes(c.pc, []byte{0x34, 0x02})
	o := mem(off(dx(16)))
	o.Store(c, []byte{expected}, &mc)
	result := c.memory.Read(0xACDE)
	if result != expected {
		t.Errorf("expected %X, got %X\n", expected, result)
	}
	expectedMc := go_gb.MC(3)
	if mc != expectedMc {
		t.Errorf("expected MC %d, got %d\n", expectedMc, mc)
	}
}
