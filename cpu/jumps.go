package cpu

import "go-gb"

func jr(c *cpu) go_gb.MC {
	var cycles go_gb.MC
	opcode := c.readOpcode(&cycles)
	e := int(opcode)
	pc := int(c.pc) + e
	c.setPc(uint16(pc), &cycles)
	return cycles
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
	var cycles go_gb.MC
	addrBytes := c.popStack(2, &cycles)
	c.setPc((uint16(addrBytes[1])<<8)|uint16(addrBytes[0]), &cycles)
	return cycles
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
	var cycles go_gb.MC
	reg := rx(HL).Load(c, &cycles)
	c.pc = go_gb.FromBytes(reg)
	return cycles
}

func jp(dst Ptr) Instr { // note: don't forget to check if it's a jump command (don't inc pc)
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := dst.Load(c, &cycles)
		c.setPc(go_gb.FromBytes(bytes), &cycles)
		return cycles
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
	var cycles go_gb.MC
	addr := c.readFromPc(2, &cycles)
	callAddr(c, addr, &cycles)
	return cycles
}

func callAddr(c *cpu, addr []byte, mc *go_gb.MC) {
	pcBytes := go_gb.ToBytesReverse(c.pc, true)
	c.pushStack(pcBytes, mc)
	c.pc = go_gb.FromBytes(addr)
	if mc != nil {
		*mc += 1
	}
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
		var cycles go_gb.MC
		pcBytes := go_gb.ToBytesReverse(c.pc, true)
		c.pushStack(pcBytes, &cycles)

		bytes := dst.Load(c, &cycles)
		c.setPc(go_gb.FromBytes(bytes), &cycles)
		return cycles
	}
}
