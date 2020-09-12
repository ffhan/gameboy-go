package cpu

import "go-gb"

func jr(c *cpu) go_gb.MC {
	opcode, mc := c.readOpcode()
	e := int(opcode)
	pc := int(c.pc) + e
	return mc + c.setPc(uint16(pc))
}

func jrnc(bit int) Instr {
	return func(c *cpu) go_gb.MC {
		if !c.getFlag(bit) {
			return jr(c)
		}
		return 1
	}
}

func jrc(bit int) Instr {
	return func(c *cpu) go_gb.MC {
		if c.getFlag(bit) {
			return jr(c)
		}
		return 1
	}
}

func ret(c *cpu) go_gb.MC {
	addrBytes, mc := c.popStack(2)
	return mc + c.setPc((uint16(addrBytes[1])<<8)|uint16(addrBytes[0]))
}

func retnc(bit int) Instr {
	return func(c *cpu) go_gb.MC {
		if !c.getFlag(bit) {
			return ret(c) + 1
		}
		return 1
	}
}

func retc(bit int) Instr {
	return func(c *cpu) go_gb.MC {
		if c.getFlag(bit) {
			return ret(c) + 1
		}
		return 1
	}
}

func jpHl(c *cpu) go_gb.MC {
	reg, _ := rx(HL).Load(c)
	c.pc = go_gb.FromBytes(reg)
	return 0
}

func jp(dst Ptr) Instr { // note: don't forget to check if it's a jump command (don't inc pc)
	return func(c *cpu) go_gb.MC {
		bytes, mc := dst.Load(c)
		return mc + c.setPc(go_gb.FromBytes(bytes))
	}
}

// JP NOT conditional
func jpnc(bit int, dst Ptr) Instr {
	instr := jp(dst)
	return func(c *cpu) go_gb.MC {
		if !c.getFlag(bit) {
			return instr(c)
		}
		return 2
	}
}

// JP conditional
func jpc(bit int, dst Ptr) Instr {
	instr := jp(dst)
	return func(c *cpu) go_gb.MC {
		if c.getFlag(bit) {
			return instr(c)
		}
		return 2
	}
}

func call(c *cpu) go_gb.MC {
	addr, m := c.readFromPc(2)
	mc := callAddr(c, addr)
	return mc + m
}

func callAddr(c *cpu, addr []byte) go_gb.MC {
	pcBytes := go_gb.MsbLsbBytes(c.pc, true)
	cycles := c.pushStack(pcBytes)
	c.pc = go_gb.FromBytes(addr)
	return cycles + 1 // for reading SP
}

func callc(bit int) Instr {
	return func(c *cpu) go_gb.MC {
		if c.getFlag(bit) {
			return call(c)
		}
		return 2
	}
}

func callcc(bit int) Instr {
	return func(c *cpu) go_gb.MC {
		if !c.getFlag(bit) {
			return call(c)
		}
		return 2
	}
}

func reti(c *cpu) go_gb.MC {
	c.ime = true
	return ret(c)
}

func rst(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		pcBytes := go_gb.MsbLsbBytes(c.pc, true)
		cycles := c.pushStack(pcBytes)

		bytes, mc := dst.Load(c)
		return cycles + mc + c.setPc(go_gb.FromBytes(bytes))
	}
}
