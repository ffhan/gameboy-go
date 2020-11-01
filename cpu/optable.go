package cpu

import "go-gb"

var ( // todo check table (especially loads)
	optable = [...]Instr{
		NOP, load(rx(go_gb.BC), dx(16)), load(mr(go_gb.BC), rx(go_gb.A)), inc16bit(rx(go_gb.BC)), inc8bit(rx(go_gb.B)), dec8bit(rx(go_gb.B)), load(rx(go_gb.B), dx(8)), rlca, load(md(16), sp()), add16b(rx(go_gb.HL), rx(go_gb.BC)), load(rx(go_gb.A), mr(go_gb.BC)), dec16bit(rx(go_gb.BC)), inc8bit(rx(go_gb.C)), dec8bit(rx(go_gb.C)), load(rx(go_gb.C), dx(8)), rrca,
		STOP, load(rx(go_gb.DE), dx(16)), load(mr(go_gb.DE), rx(go_gb.A)), inc16bit(rx(go_gb.DE)), inc8bit(rx(go_gb.D)), dec8bit(rx(go_gb.D)), load(rx(go_gb.D), dx(8)), rla, jr, add16b(rx(go_gb.HL), rx(go_gb.DE)), load(rx(go_gb.A), mr(go_gb.DE)), dec16bit(rx(go_gb.DE)), inc8bit(rx(go_gb.E)), dec8bit(rx(go_gb.E)), load(rx(go_gb.E), dx(8)), rra,
		jrnc(BitZ), load(rx(go_gb.HL), dx(16)), ldHl(nil, rx(go_gb.A), true), inc16bit(rx(go_gb.HL)), inc8bit(rx(go_gb.H)), dec8bit(rx(go_gb.H)), load(rx(go_gb.H), dx(8)), daa, jrc(BitZ), add16b(rx(go_gb.HL), rx(go_gb.HL)), ldHl(rx(go_gb.A), nil, true), dec16bit(rx(go_gb.HL)), inc8bit(rx(go_gb.L)), dec8bit(rx(go_gb.L)), load(rx(go_gb.L), dx(8)), cpl,
		jrnc(BitC), load(sp(), dx(16)), ldHl(nil, rx(go_gb.A), false), inc16bit(sp()), inc8bit(mr(go_gb.HL)), dec8bit(mr(go_gb.HL)), load(mr(go_gb.HL), dx(8)), scf, jrc(BitC), add16b(rx(go_gb.HL), sp()), ldHl(rx(go_gb.A), nil, false), dec16bit(sp()), inc8bit(rx(go_gb.A)), dec8bit(rx(go_gb.A)), load(rx(go_gb.A), dx(8)), ccf,
		load(rx(go_gb.B), rx(go_gb.B)), load(rx(go_gb.B), rx(go_gb.C)), load(rx(go_gb.B), rx(go_gb.D)), load(rx(go_gb.B), rx(go_gb.E)), load(rx(go_gb.B), rx(go_gb.H)), load(rx(go_gb.B), rx(go_gb.L)), load(rx(go_gb.B), mr(go_gb.HL)), load(rx(go_gb.B), rx(go_gb.A)), load(rx(go_gb.C), rx(go_gb.B)), load(rx(go_gb.C), rx(go_gb.C)), load(rx(go_gb.C), rx(go_gb.D)), load(rx(go_gb.C), rx(go_gb.E)), load(rx(go_gb.C), rx(go_gb.H)), load(rx(go_gb.C), rx(go_gb.L)), load(rx(go_gb.C), mr(go_gb.HL)), load(rx(go_gb.C), rx(go_gb.A)),
		load(rx(go_gb.D), rx(go_gb.B)), load(rx(go_gb.D), rx(go_gb.C)), load(rx(go_gb.D), rx(go_gb.D)), load(rx(go_gb.D), rx(go_gb.E)), load(rx(go_gb.D), rx(go_gb.H)), load(rx(go_gb.D), rx(go_gb.L)), load(rx(go_gb.D), mr(go_gb.HL)), load(rx(go_gb.D), rx(go_gb.A)), load(rx(go_gb.E), rx(go_gb.B)), load(rx(go_gb.E), rx(go_gb.C)), load(rx(go_gb.E), rx(go_gb.D)), load(rx(go_gb.E), rx(go_gb.E)), load(rx(go_gb.E), rx(go_gb.H)), load(rx(go_gb.E), rx(go_gb.L)), load(rx(go_gb.E), mr(go_gb.HL)), load(rx(go_gb.E), rx(go_gb.A)),
		load(rx(go_gb.H), rx(go_gb.B)), load(rx(go_gb.H), rx(go_gb.C)), load(rx(go_gb.H), rx(go_gb.D)), load(rx(go_gb.H), rx(go_gb.E)), load(rx(go_gb.H), rx(go_gb.H)), load(rx(go_gb.H), rx(go_gb.L)), load(rx(go_gb.H), mr(go_gb.HL)), load(rx(go_gb.H), rx(go_gb.A)), load(rx(go_gb.L), rx(go_gb.B)), load(rx(go_gb.L), rx(go_gb.C)), load(rx(go_gb.L), rx(go_gb.D)), load(rx(go_gb.L), rx(go_gb.E)), load(rx(go_gb.L), rx(go_gb.H)), load(rx(go_gb.L), rx(go_gb.L)), load(rx(go_gb.L), mr(go_gb.HL)), load(rx(go_gb.L), rx(go_gb.A)),
		load(mr(go_gb.HL), rx(go_gb.B)), load(mr(go_gb.HL), rx(go_gb.C)), load(mr(go_gb.HL), rx(go_gb.D)), load(mr(go_gb.HL), rx(go_gb.E)), load(mr(go_gb.HL), rx(go_gb.H)), load(mr(go_gb.HL), rx(go_gb.L)), halt, load(mr(go_gb.HL), rx(go_gb.A)), load(rx(go_gb.A), rx(go_gb.B)), load(rx(go_gb.A), rx(go_gb.C)), load(rx(go_gb.A), rx(go_gb.D)), load(rx(go_gb.A), rx(go_gb.E)), load(rx(go_gb.A), rx(go_gb.H)), load(rx(go_gb.A), rx(go_gb.L)), load(rx(go_gb.A), mr(go_gb.HL)), load(rx(go_gb.A), rx(go_gb.A)),

		add8b(rx(go_gb.B)), add8b(rx(go_gb.C)), add8b(rx(go_gb.D)), add8b(rx(go_gb.E)), add8b(rx(go_gb.H)), add8b(rx(go_gb.L)), add8b(mr(go_gb.HL)), add8b(rx(go_gb.A)), adc8b(rx(go_gb.B)), adc8b(rx(go_gb.C)), adc8b(rx(go_gb.D)), adc8b(rx(go_gb.E)), adc8b(rx(go_gb.H)), adc8b(rx(go_gb.L)), adc8b(mr(go_gb.HL)), adc8b(rx(go_gb.A)),
		sub(rx(go_gb.B)), sub(rx(go_gb.C)), sub(rx(go_gb.D)), sub(rx(go_gb.E)), sub(rx(go_gb.H)), sub(rx(go_gb.L)), sub(mr(go_gb.HL)), sub(rx(go_gb.A)), sbc(rx(go_gb.B)), sbc(rx(go_gb.C)), sbc(rx(go_gb.D)), sbc(rx(go_gb.E)), sbc(rx(go_gb.H)), sbc(rx(go_gb.L)), sbc(mr(go_gb.HL)), sbc(rx(go_gb.A)),
		and(rx(go_gb.B)), and(rx(go_gb.C)), and(rx(go_gb.D)), and(rx(go_gb.E)), and(rx(go_gb.H)), and(rx(go_gb.L)), and(mr(go_gb.HL)), and(rx(go_gb.A)), xor(rx(go_gb.B)), xor(rx(go_gb.C)), xor(rx(go_gb.D)), xor(rx(go_gb.E)), xor(rx(go_gb.H)), xor(rx(go_gb.L)), xor(mr(go_gb.HL)), xor(rx(go_gb.A)),
		or(rx(go_gb.B)), or(rx(go_gb.C)), or(rx(go_gb.D)), or(rx(go_gb.E)), or(rx(go_gb.H)), or(rx(go_gb.L)), or(mr(go_gb.HL)), or(rx(go_gb.A)), cp(rx(go_gb.A), rx(go_gb.B)), cp(rx(go_gb.A), rx(go_gb.C)), cp(rx(go_gb.A), rx(go_gb.D)), cp(rx(go_gb.A), rx(go_gb.E)), cp(rx(go_gb.A), rx(go_gb.H)), cp(rx(go_gb.A), rx(go_gb.L)), cp(rx(go_gb.A), mr(go_gb.HL)), cp(rx(go_gb.A), rx(go_gb.A)),

		retnc(BitZ), pop(rx(go_gb.BC)), jpnc(BitZ, dx(16)), jp(dx(16)), callcc(BitZ), push(rx(go_gb.BC)), add8b(dx(8)), rst(hc(0x00)), retc(BitZ), ret, jpc(BitZ, dx(16)), prefix, callc(BitZ), call, adc8b(dx(8)), rst(hc(0x08)),
		retnc(BitC), pop(rx(go_gb.DE)), jpnc(BitC, dx(16)), invalid, callcc(BitC), push(rx(go_gb.DE)), sub(dx(8)), rst(hc(0x10)), retc(BitC), reti, jpc(BitC, dx(16)), invalid, callc(BitC), invalid, sbc(dx(8)), rst(hc(0x18)),
		load(mem(off(dx(8))), rx(go_gb.A)), pop(rx(go_gb.HL)), load(mem(off(rx(go_gb.C))), rx(go_gb.A)), invalid, invalid, push(rx(go_gb.HL)), and(dx(8)), rst(hc(0x20)), addSp, jpHl, load(md(16), rx(go_gb.A)), invalid, invalid, invalid, xor(dx(8)), rst(hc(0x28)),
		load(rx(go_gb.A), mem(off(dx(8)))), pop(rx(go_gb.AF)), load(rx(go_gb.A), mem(off(rx(go_gb.C)))), di, invalid, push(rx(go_gb.AF)), or(dx(8)), rst(hc(0x30)), ldHlSp, ldSpHl, load(rx(go_gb.A), md(16)), ei, invalid, invalid, cp(rx(go_gb.A), dx(8)), rst(hc(0x38)),
	}
	cbOptable = createCbOptable()
)

func createCbOptable() []Instr {
	registers := [...]Ptr{rx(go_gb.B), rx(go_gb.C), rx(go_gb.D), rx(go_gb.E), rx(go_gb.H), rx(go_gb.L), mr(go_gb.HL), rx(go_gb.A)}
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
			in = bit(hc(byte((i-4*16)/8)), reg)
		case 16, 17, 18, 19, 20, 21, 22, 23:
			in = res(hc(byte((i-8*16)/8)), reg)
		default:
			in = set(hc(byte((i-12*16)/8)), reg)
		}
		table = append(table, in)
	}
	if len(table) != size {
		panic("invalid cb table")
	}
	return table
}
