package cpu

import (
	"errors"
	"fmt"
	"go-gb"
)

type Store interface {
	Store(c *cpu, b []byte) go_gb.MC
}

type Loader interface {
	Load(c *cpu) ([]byte, go_gb.MC)
}

type Ptr interface {
	Store
	Loader
}

type reg struct {
	addr registerName
}

// loads a value from a register
//
// examples: A, HL, BC
func rx(addr registerName) reg {
	return reg{addr}
}

// loads from memory addressed in the pointer
//
// examples: (A), (HL), (BC)...
func mr(addr registerName) mPtr {
	return mPtr{rx(addr)}
}

// loads from memory from the current PC (size 8/16 bits)
//
// nn
func dx(size int) data {
	return data{size: size}
}

// loads from memory addressed in the next (size) bytes (read from memory from the current PC, see dx)
//
// (nn)
func md(size int) mPtr {
	return mPtr{dx(size)}
}

func mem(ptr Ptr) mPtr {
	return mPtr{ptr}
}

func (r reg) Store(c *cpu, b []byte) go_gb.MC {
	copy(c.getRegister(r.addr), b)
	return 0
}

func (r reg) Load(c *cpu) ([]byte, go_gb.MC) {
	return c.getRegister(r.addr), 0
}

type offset struct {
	dst    Ptr
	offset int
}

func (o offset) Store(c *cpu, b []byte) go_gb.MC {
	return o.dst.Store(c, go_gb.LsbMsbBytes(uint16(int(go_gb.FromBytes(b))+o.offset), len(b) == 2 || o.offset > 0xFF))
}

func (o offset) Load(c *cpu) ([]byte, go_gb.MC) {
	bytes, mc := o.dst.Load(c)
	return go_gb.LsbMsbBytes(uint16(int(go_gb.FromBytes(bytes))+o.offset), len(bytes) == 2 || o.offset > 0xFF), mc
}

func off(dst Ptr, o int) offset {
	return offset{dst, o}
}

type mPtr struct {
	addr Ptr
}

func (m mPtr) Store(c *cpu, b []byte) go_gb.MC {
	bytes, mc := m.addr.Load(c)
	mc += c.memory.StoreBytes(go_gb.FromBytes(bytes), b)
	return mc
}

func (m mPtr) Load(c *cpu) ([]byte, go_gb.MC) {
	pointerAddr, mc := m.addr.Load(c)
	result, m2 := c.memory.ReadBytes(go_gb.FromBytes(pointerAddr), 1)
	return result, mc + m2
}

type data struct {
	size int
}

var InvalidStoreErr = errors.New("invalid store call")

func (d data) Store(c *cpu, b []byte) go_gb.MC {
	panic(InvalidStoreErr)
}

func (d data) Load(c *cpu) ([]byte, go_gb.MC) {
	var cycles go_gb.MC
	n := d.size / 8
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		b, mc := c.readOpcode()
		bytes[i] = b
		cycles += mc
	}
	return bytes, cycles
}

type stackPtr struct {
}

func (s stackPtr) Store(c *cpu, b []byte) go_gb.MC {
	if len(b) != 2 {
		panic(fmt.Errorf("invalid SP store %v", b))
	}
	c.sp = go_gb.FromBytes(b)
	return 1
}

func (s stackPtr) Load(c *cpu) ([]byte, go_gb.MC) {
	return go_gb.LsbMsbBytes(c.sp, true), 1
}

type pc struct {
}

func (p pc) Store(c *cpu, b []byte) go_gb.MC {
	c.pc = go_gb.FromBytes(b)
	return 0
}

func (p pc) Load(c *cpu) ([]byte, go_gb.MC) {
	return go_gb.LsbMsbBytes(c.pc, true), 0
}

type hardcoded struct {
	val byte
}

func (h hardcoded) Store(c *cpu, b []byte) go_gb.MC {
	panic(InvalidStoreErr)
}

func (h hardcoded) Load(c *cpu) ([]byte, go_gb.MC) {
	return []byte{h.val}, 0
}

func hc(b byte) hardcoded {
	return hardcoded{b}
}

// handles SP loading and storing
func sp() stackPtr {
	return stackPtr{}
}
