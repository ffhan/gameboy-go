package cpu

import (
	go_gb "go-gb"
)

func rlca(c *cpu) go_gb.MC {
	return rc(c, rx(go_gb.A), true, true) - 1
}

func rrca(c *cpu) go_gb.MC {
	return rc(c, rx(go_gb.A), false, true) - 1
}

func r(c *cpu, dst Ptr, left bool) go_gb.MC {
	var cycles go_gb.MC
	bytes := dst.Load(c, &cycles)
	b := bytes[0]
	var old byte
	carry := c.getFlag(BitC)
	if left {
		old = b >> 7
		b <<= 1
		if carry {
			b |= 1
		}
	} else {
		old = b & 1
		b >>= 1
		if carry {
			b |= 0x80
		}
	}
	dst.Store(c, []byte{b}, &cycles)
	c.setFlag(BitZ, b == 0)
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, old == 1)
	return cycles + 1
}

func rla(c *cpu) go_gb.MC {
	return r(c, rx(go_gb.A), true) - 1
}

func rra(c *cpu) go_gb.MC {
	return r(c, rx(go_gb.A), false) - 1
}

func and(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC

		bytes := src.Load(c, &cycles)
		srcVal := bytes[0]
		bytes = dst.Load(c, &cycles)
		orig := bytes[0]
		dstVal := orig & srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, true)
		c.setFlag(BitC, false)
		dst.Store(c, []byte{dstVal}, &cycles)
		return cycles
	}
}

func xor(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC

		bytes := src.Load(c, &cycles)
		srcVal := bytes[0]
		bytes = dst.Load(c, &cycles)
		orig := bytes[0]
		dstVal := orig ^ srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		dst.Store(c, []byte{dstVal}, &cycles)
		return cycles
	}
}

func or(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		srcVal := bytes[0]
		bytes = dst.Load(c, &cycles)
		orig := bytes[0]
		dstVal := orig | srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		dst.Store(c, []byte{dstVal}, &cycles)
		return cycles
	}
}

func cpl(c *cpu) go_gb.MC {
	val := c.GetRegister(go_gb.A)
	val[0] = ^val[0]
	c.setFlag(BitN, true)
	c.setFlag(BitH, true)
	return 0
}

func rc(c *cpu, dst Ptr, left, resetZ bool) go_gb.MC {
	var cycles go_gb.MC
	bytes := dst.Load(c, &cycles)
	b := bytes[0]
	var old byte
	if left {
		old = b >> 7
		b <<= 1
		b |= old
	} else {
		old = b & 1
		b >>= 1
		b |= old << 7
	}
	dst.Store(c, []byte{b}, &cycles)
	if resetZ {
		c.setFlag(BitZ, false)
	} else {
		c.setFlag(BitZ, b == 0)
	}
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, old == 1)
	return cycles + 1
}

func rlc(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return rc(c, dst, true, false)
	}
}

func rrc(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return rc(c, dst, false, false)
	}
}

func rl(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return r(c, dst, true)
	}
}

func rr(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return r(c, dst, false)
	}
}

func sa(c *cpu, dst Ptr, left bool) go_gb.MC {
	var cycles go_gb.MC
	bytes := dst.Load(c, &cycles)
	b := bytes[0]
	var old byte
	if left {
		old = b >> 7
		b <<= 1
	} else {
		old = b & 1
		b >>= 1
	}
	dst.Store(c, []byte{b}, &cycles)
	c.setFlag(BitZ, b == 0)
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, old == 1)
	return cycles + 1
}

func sla(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return sa(c, dst, true)
	}
}

func sra(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := dst.Load(c, &cycles)
		b := bytes[0]
		var old byte
		msb := b & 0xF0
		old = b & 1
		b >>= 1
		b |= msb
		dst.Store(c, []byte{b}, &cycles)
		c.setFlag(BitZ, b == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, old == 1)
		return cycles + 1
	}
}

// note: srl is analogous to sla, not sra!
func srl(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return sa(c, dst, false)
	}
}

func swap(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := dst.Load(c, &cycles)
		val := bytes[0]
		msn := val & 0xF0
		lsn := val & 0x0F
		result := (msn >> 4) | (lsn << 4)
		dst.Store(c, []byte{result}, &cycles)

		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		return cycles + 1
	}
}
