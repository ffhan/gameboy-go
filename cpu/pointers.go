package cpu

import (
	"fmt"
	"go-gb"
)

type Ptr interface {
	Store(c *cpu, b []byte)
	Load(c *cpu) []byte
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

func (r reg) Store(c *cpu, b []byte) {
	copy(c.getRegister(r.addr), b)
}

func (r reg) Load(c *cpu) []byte {
	return c.getRegister(r.addr)
}

type mPtr struct {
	addr Ptr
}

func (m mPtr) Store(c *cpu, b []byte) {
	c.memory.StoreBytes(go_gb.UnifyBytes(m.addr.Load(c)), b)
}

func (m mPtr) Load(c *cpu) []byte {
	return c.memory.ReadBytes(go_gb.UnifyBytes(m.addr.Load(c)), 1)
}

type data struct {
	size int
}

func (d data) Store(c *cpu, b []byte) {
	panic("cannot store in data")
}

func (d data) Load(c *cpu) []byte {
	n := d.size / 8
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[n-1-i] = c.readOpcode()
	}
	return bytes
}

type stackPtr struct {
}

func (s stackPtr) Store(c *cpu, b []byte) {
	if len(b) != 2 {
		panic(fmt.Errorf("invalid SP store %v", b))
	}
	c.sp = go_gb.UnifyBytes(b)
}

func (s stackPtr) Load(c *cpu) []byte {
	return go_gb.SeparateUint16(c.sp)
}

type pc struct {
}

func (p pc) Store(c *cpu, b []byte) {
	c.pc = go_gb.UnifyBytes(b)
}

func (p pc) Load(c *cpu) []byte {
	return go_gb.SeparateUint16(c.pc)
}

type hardcoded struct {
	val byte
}

func (h hardcoded) Store(c *cpu, b []byte) {
	panic("cannot store in hardcoded values")
}

func (h hardcoded) Load(c *cpu) []byte {
	return []byte{h.val}
}

func hc(b byte) hardcoded {
	return hardcoded{b}
}

// handles SP loading and storing
func sp() stackPtr {
	return stackPtr{}
}
