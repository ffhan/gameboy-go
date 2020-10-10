package cpu

import (
	"errors"
	"fmt"
	"go-gb"
)

type Store interface {
	Store(c *cpu, b []byte, mc *go_gb.MC)
}

type Loader interface {
	Load(c *cpu, mc *go_gb.MC) []byte
}

type Ptr interface {
	Store
	Loader
}

type reg struct {
	addr go_gb.RegisterName
}

// loads a value from a register
//
// examples: A, HL, BC
func rx(addr go_gb.RegisterName) reg {
	return reg{addr}
}

// loads from memory addressed in the pointer
//
// examples: (A), (HL), (BC)...
func mr(addr go_gb.RegisterName) mPtr {
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

func (r reg) Store(c *cpu, b []byte, mc *go_gb.MC) {
	copy(c.GetRegister(r.addr), b)
}

func (r reg) Load(c *cpu, mc *go_gb.MC) []byte {
	return c.GetRegister(r.addr)
}

type offset struct {
	dst Ptr
}

func (o offset) Store(c *cpu, b []byte, mc *go_gb.MC) {
	o.dst.Store(c, []byte{b[0], 0xFF}, mc)
}

func (o offset) Load(c *cpu, mc *go_gb.MC) []byte {
	bytes := o.dst.Load(c, mc)
	return []byte{bytes[0], 0xFF}
}

func off(dst Ptr) offset {
	return offset{dst}
}

type mPtr struct {
	addr Ptr
}

func (m mPtr) Store(c *cpu, b []byte, mc *go_gb.MC) {
	bytes := m.addr.Load(c, mc)
	c.storeBytes(go_gb.FromBytes(bytes), b, mc)
}

func (m mPtr) Load(c *cpu, mc *go_gb.MC) []byte {
	pointerAddr := m.addr.Load(c, mc)
	result := c.readBytes(go_gb.FromBytes(pointerAddr), 1, mc)
	return result
}

type data struct {
	size int
}

var InvalidStoreErr = errors.New("invalid store call")

func (d data) Store(c *cpu, b []byte, mc *go_gb.MC) {
	panic(InvalidStoreErr)
}

func (d data) Load(c *cpu, mc *go_gb.MC) []byte {
	n := d.size / 8
	return c.readFromPc(uint16(n), mc)
}

type stackPtr struct {
}

func (s stackPtr) Store(c *cpu, b []byte, mc *go_gb.MC) {
	if len(b) != 2 {
		panic(fmt.Errorf("invalid SP store %v", b))
	}
	c.sp = go_gb.FromBytes(b)
}

func (s stackPtr) Load(c *cpu, mc *go_gb.MC) []byte {
	return go_gb.ToBytes(c.sp, true)
}

type pc struct {
}

func (p pc) Store(c *cpu, b []byte, mc *go_gb.MC) {
	c.pc = go_gb.FromBytes(b)
}

func (p pc) Load(c *cpu, mc *go_gb.MC) []byte {
	return go_gb.ToBytes(c.pc, true)
}

type hardcoded struct {
	val byte
}

func (h hardcoded) Store(c *cpu, b []byte, mc *go_gb.MC) {
	panic(InvalidStoreErr)
}

func (h hardcoded) Load(c *cpu, mc *go_gb.MC) []byte {
	return []byte{h.val}
}

func hc(b byte) hardcoded {
	return hardcoded{b}
}

// handles SP loading and storing
func sp() stackPtr {
	return stackPtr{}
}
