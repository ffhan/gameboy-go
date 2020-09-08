package cpu

import "go-gb"

func jr(c *cpu) error {
	e := int16(c.readOpcode())
	pc := int16(c.pc) + e
	c.pc = uint16(pc)
	return nil
}

func jrcc(bit int) Instr {
	return func(c *cpu) error {
		if !c.getFlag(bit) {
			return jr(c)
		}
		return nil
	}
}

func jrc(bit int) Instr {
	return func(c *cpu) error {
		if c.getFlag(bit) {
			return jr(c)
		}
		return nil
	}
}

func ret(c *cpu) error {
	addrBytes := c.popStack(2)
	c.pc = (uint16(addrBytes[1]) << 8) | uint16(addrBytes[0])
	return nil
}

func retcc(bit int) Instr {
	return func(c *cpu) error {
		if !c.getFlag(bit) {
			return ret(c)
		}
		return nil
	}
}

func retc(bit int) Instr {
	return func(c *cpu) error {
		if c.getFlag(bit) {
			return ret(c)
		}
		return nil
	}
}

func jp(dst Ptr) Instr { // note: don't forget to check if it's a jump command (don't inc pc)
	return func(c *cpu) error {
		c.pc = go_gb.UnifyBytes(dst.Load(c))
		return nil
	}
}

// JP conditional complement
func jpcc(bit int, dst Ptr) Instr {
	instr := jp(dst)
	return func(c *cpu) error {
		if !c.getFlag(bit) {
			return instr(c)
		}
		return nil
	}
}

// JP conditional
func jpc(bit int, dst Ptr) Instr {
	instr := jp(dst)
	return func(c *cpu) error {
		if c.getFlag(bit) {
			return instr(c)
		}
		return nil
	}
}

func call(c *cpu) error {
	addr := c.memory.ReadBytes(c.pc, 2)
	c.pc += 2
	pcBytes := go_gb.SeparateUint16(c.pc)
	c.pushStack(pcBytes[1:2])
	c.pushStack(pcBytes[0:1])
	c.pc = go_gb.UnifyBytes(addr)
	return nil
}

func callc(bit int) Instr {
	return func(c *cpu) error {
		if c.getFlag(bit) {
			return call(c)
		}
		return nil
	}
}

func callcc(bit int) Instr {
	return func(c *cpu) error {
		if !c.getFlag(bit) {
			return call(c)
		}
		return nil
	}
}

func reti(c *cpu) error {
	c.ime = true
	return ret(c)
}

func rst(dst Ptr) Instr {
	return func(c *cpu) error {
		pcBytes := go_gb.SeparateUint16(c.pc)
		c.pushStack(pcBytes[1:2])
		c.pushStack(pcBytes[0:1])

		c.pc = go_gb.UnifyBytes(dst.Load(c))
		return nil
	}
}