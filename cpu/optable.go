package cpu

import (
	go_gb "go-gb"
)

var ( // todo check table (especially loads)
	optable = [...]Instr{
		NOP, load(rx(BC), dx(16)), load(mr(BC), rx(A)), inc16bit(rx(BC)), inc8bit(rx(B)), dec8bit(rx(B)), load(rx(B), dx(8)), rlca, load(md(16), sp()), add16b(rx(HL), rx(BC)), load(rx(A), mr(BC)), dec16bit(rx(BC)), inc8bit(rx(C)), dec8bit(rx(C)), load(rx(C), dx(8)), rrca,
		STOP, load(rx(DE), dx(16)), load(mr(DE), rx(A)), inc16bit(rx(DE)), inc8bit(rx(D)), dec8bit(rx(D)), load(rx(D), dx(8)), rla, jr, add16b(rx(HL), rx(DE)), load(rx(A), mr(DE)), dec16bit(rx(DE)), inc8bit(rx(E)), dec8bit(rx(E)), load(rx(E), dx(8)), rra,
		jrnc(BitZ), load(rx(HL), dx(16)), ldHl(nil, rx(A), true), inc16bit(rx(HL)), inc8bit(rx(H)), dec8bit(rx(H)), load(rx(H), dx(8)), daa, jrc(BitZ), add16b(rx(HL), rx(HL)), ldHl(rx(A), nil, true), dec16bit(rx(HL)), inc8bit(rx(L)), dec8bit(rx(L)), load(rx(L), dx(8)), cpl,
		jrnc(BitC), load(sp(), dx(16)), ldHl(nil, rx(A), false), inc16bit(sp()), inc8bit(mr(HL)), dec8bit(mr(HL)), load(mr(HL), dx(8)), scf, jrc(BitC), addHlSp, ldHl(rx(A), nil, false), dec16bit(sp()), inc8bit(rx(A)), dec8bit(rx(A)), load(rx(A), dx(8)), ccf,
		load(rx(B), rx(B)), load(rx(B), rx(C)), load(rx(B), rx(D)), load(rx(B), rx(E)), load(rx(B), rx(H)), load(rx(B), rx(L)), load(rx(B), mr(HL)), load(rx(B), rx(A)), load(rx(C), rx(B)), load(rx(C), rx(C)), load(rx(C), rx(D)), load(rx(C), rx(E)), load(rx(C), rx(H)), load(rx(C), rx(L)), load(rx(C), mr(HL)), load(rx(C), rx(A)),
		load(rx(D), rx(B)), load(rx(D), rx(C)), load(rx(D), rx(D)), load(rx(D), rx(E)), load(rx(D), rx(H)), load(rx(D), rx(L)), load(rx(D), mr(HL)), load(rx(D), rx(A)), load(rx(E), rx(B)), load(rx(E), rx(C)), load(rx(E), rx(D)), load(rx(E), rx(E)), load(rx(E), rx(H)), load(rx(E), rx(L)), load(rx(E), mr(HL)), load(rx(E), rx(A)),
		load(rx(H), rx(B)), load(rx(H), rx(C)), load(rx(H), rx(D)), load(rx(H), rx(E)), load(rx(H), rx(H)), load(rx(H), rx(L)), load(rx(H), mr(HL)), load(rx(H), rx(A)), load(rx(L), rx(B)), load(rx(L), rx(C)), load(rx(L), rx(D)), load(rx(L), rx(E)), load(rx(L), rx(H)), load(rx(L), rx(L)), load(rx(L), mr(HL)), load(rx(L), rx(A)),
		load(mr(HL), rx(B)), load(mr(HL), rx(C)), load(mr(HL), rx(D)), load(mr(HL), rx(E)), load(mr(HL), rx(H)), load(mr(HL), rx(L)), halt, load(mr(HL), rx(A)), load(rx(A), rx(B)), load(rx(A), rx(C)), load(rx(A), rx(D)), load(rx(A), rx(E)), load(rx(A), rx(H)), load(rx(A), rx(L)), load(rx(A), mr(HL)), load(rx(A), rx(A)),

		add8b(rx(A), rx(B)), add8b(rx(A), rx(C)), add8b(rx(A), rx(D)), add8b(rx(A), rx(E)), add8b(rx(A), rx(H)), add8b(rx(A), rx(L)), add8b(rx(A), mr(HL)), add8b(rx(A), rx(A)), adc8b(rx(A), rx(B)), adc8b(rx(A), rx(C)), adc8b(rx(A), rx(D)), adc8b(rx(A), rx(E)), adc8b(rx(A), rx(H)), adc8b(rx(A), rx(L)), adc8b(rx(A), mr(HL)), adc8b(rx(A), rx(A)),
		sub(rx(A), rx(B)), sub(rx(A), rx(C)), sub(rx(A), rx(D)), sub(rx(A), rx(E)), sub(rx(A), rx(H)), sub(rx(A), rx(L)), sub(rx(A), mr(HL)), sub(rx(A), rx(A)), sbc(rx(A), rx(B)), sbc(rx(A), rx(C)), sbc(rx(A), rx(D)), sbc(rx(A), rx(E)), sbc(rx(A), rx(H)), sbc(rx(A), rx(L)), sbc(rx(A), mr(HL)), sbc(rx(A), rx(A)),
		and(rx(A), rx(B)), and(rx(A), rx(C)), and(rx(A), rx(D)), and(rx(A), rx(E)), and(rx(A), rx(H)), and(rx(A), rx(L)), and(rx(A), mr(HL)), and(rx(A), rx(A)), xor(rx(A), rx(B)), xor(rx(A), rx(C)), xor(rx(A), rx(D)), xor(rx(A), rx(E)), xor(rx(A), rx(H)), xor(rx(A), rx(L)), xor(rx(A), mr(HL)), xor(rx(A), rx(A)),
		or(rx(A), rx(B)), or(rx(A), rx(C)), or(rx(A), rx(D)), or(rx(A), rx(E)), or(rx(A), rx(H)), or(rx(A), rx(L)), or(rx(A), mr(HL)), or(rx(A), rx(A)), cp(rx(A), rx(B)), cp(rx(A), rx(C)), cp(rx(A), rx(D)), cp(rx(A), rx(E)), cp(rx(A), rx(H)), cp(rx(A), rx(L)), cp(rx(A), mr(HL)), cp(rx(A), rx(A)),

		retcc(BitZ), pop(rx(BC)), jpnc(BitZ, md(16)), jp(dx(16)), callcc(BitZ), push(rx(BC)), add8b(rx(A), dx(8)), rst(hc(0x00)), retc(BitZ), ret, jpc(BitZ, dx(16)), prefix, callc(BitZ), call, adc8b(rx(A), dx(8)), rst(hc(0x08)),
		retcc(BitC), pop(rx(DE)), jpnc(BitC, md(16)), invalid, callcc(BitC), push(rx(DE)), sub(rx(A), dx(8)), rst(hc(0x10)), retc(BitC), reti, jpc(BitC, dx(16)), invalid, callc(BitC), invalid, sbc(rx(A), dx(8)), rst(hc(0x18)),
		loadIo(md(8), rx(A)), pop(rx(HL)), load(mr(C), rx(A)), invalid, invalid, push(rx(HL)), and(rx(A), dx(8)), rst(hc(0x20)), addSp, jp(mr(HL)), load(md(16), rx(A)), invalid, invalid, invalid, xor(rx(A), dx(8)), rst(hc(0x28)),
		loadIo(rx(A), md(8)), pop(rx(AF)), load(rx(A), mr(C)), di, invalid, push(rx(AF)), or(rx(A), dx(8)), rst(hc(0x30)), ldHlSp, load(sp(), rx(HL)), load(rx(A), md(16)), ei, invalid, invalid, cp(rx(A), dx(8)), rst(hc(0x38)),
	}
	cbOptable = createCbOptable()
)

func createCbOptable() []Instr {
	registers := [...]Ptr{rx(B), rx(C), rx(D), rx(E), rx(H), rx(L), mr(HL), rx(A)}
	size := 256
	table := make([]Instr, 0, size)
	for i := 0; i < size; i++ {
		offset := i / 8
		reg := registers[i%8]
		var in Instr
		switch offset {
		case 0:
			in = rlc(reg)
		case 1:
			in = rrc(reg)
		case 2:
			in = rl(reg)
		case 3:
			in = rr(reg)
		case 4:
			in = sla(reg)
		case 5:
			in = sra(reg)
		case 6:
			in = swap(reg)
		case 7:
			in = srl(reg)
		case 8, 9, 10, 11, 12, 13, 14, 15:
			in = bit(hc(byte(offset-8)), reg)
		case 16, 17, 18, 19, 20, 21, 22, 23:
			in = res(hc(byte(offset-8)), reg)
		default:
			in = set(hc(byte(offset-8)), reg)
		}
		table = append(table, in)
	}
	if len(table) != size {
		panic("invalid cb table")
	}
	return table
}

func (c *cpu) Step() error {
	instr := optable[c.readOpcode()]
	err := instr(c) // instructions might push different opcodes before
	if c.eiWaiting {
		c.ime = true
		c.eiWaiting = false
	} else if c.diWaiting {
		c.ime = false
		c.diWaiting = false
	}
	return err
}

func (c *cpu) setFlag(bit int, val bool) {
	register := &c.getRegister(F)[0]
	go_gb.Set(register, bit, val)
}

func (c *cpu) getFlag(bit int) bool {
	return go_gb.Bit(c.getRegister(F)[0], bit)
}
