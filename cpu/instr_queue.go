package cpu

import (
	"fmt"
	"strings"
)

type state struct {
	OP                     uint16
	PC                     uint16
	SP                     uint16
	A, F, B, C, D, E, H, L byte
	Flags                  byte
	Instruction            string
	PpuMode                byte
	PpuLine                byte
}

func (s state) String() string {
	return fmt.Sprintf("OP: %04X\tPC: %04X\tSP: %04X\ta: %02X\tf: %02X\tb: %02X\tc: %02X\td: %02X\te: %02X\th: %02X\tl: %02X\tZNHC: %04b Instruction: '%s' PPU mode: %d line: %d\n",
		s.OP, s.PC, s.SP,
		s.A, s.F, s.B, s.C,
		s.D, s.E, s.H, s.L,
		s.Flags, s.Instruction,
		s.PpuMode, s.PpuLine,
	)
}

type instructionNode struct {
	next  *instructionNode
	value state
}

type instructionQueue struct {
	head     *instructionNode
	tail     *instructionNode
	size     int
	capacity int
}

func (i *instructionQueue) Tail() state {
	return i.tail.value
}

func NewInstructionQueue(capacity int) *instructionQueue {
	return &instructionQueue{capacity: capacity}
}

func (i *instructionQueue) Push(op, pc, sp uint16, a, f, b, c, d, e, h, l, flags byte, instruction string, ppuMode, ppuLine byte) {
	newNode := &instructionNode{
		next: nil,
		value: state{
			OP:          op,
			PC:          pc,
			SP:          sp,
			A:           a,
			F:           f,
			B:           b,
			C:           c,
			D:           d,
			E:           e,
			H:           h,
			L:           l,
			Flags:       flags,
			Instruction: instruction,
			PpuMode:     ppuMode,
			PpuLine:     ppuLine,
		},
	}
	if i.head == nil {
		i.head = newNode
		i.tail = newNode
		i.size = 1
		return
	}
	if i.size == i.capacity {
		if i.head == i.tail {
			i.head = newNode
			i.tail = newNode
			return
		}
		i.head = i.head.next
		i.tail.next = newNode
		i.tail = newNode
		return
	}
	i.tail.next = newNode
	i.tail = newNode
	i.size += 1
}

func (i *instructionQueue) String() string {
	var sb strings.Builder
	for node := i.head; node != nil; node = node.next {
		sb.WriteString(node.value.String())
	}
	return sb.String()
}

func (i *instructionQueue) Size() int {
	return i.size
}
