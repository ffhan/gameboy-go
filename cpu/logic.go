package cpu

import go_gb "go-gb"

func rlca(c *cpu) go_gb.MC {
	return rc(c, rx(A), true, true)
}

func rrca(c *cpu) go_gb.MC {
	return rc(c, rx(A), false, true)
}

func r(c *cpu, dst Ptr, left bool) go_gb.MC {
	bytes, mc := dst.Load(c)
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
	mc += dst.Store(c, []byte{b})
	c.setFlag(BitZ, b == 0)
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, old == 1)
	return mc
}

func rla(c *cpu) go_gb.MC {
	return r(c, rx(A), true)
}

func rra(c *cpu) go_gb.MC {
	return r(c, rx(A), false)
}

func and(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := bytes[0]
		bytes, mc2 := dst.Load(c)
		orig := bytes[0]
		dstVal := orig & srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, true)
		c.setFlag(BitC, false)
		return mc + mc2 + dst.Store(c, []byte{dstVal})
	}
}

func xor(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := bytes[0]
		bytes, mc2 := dst.Load(c)
		orig := bytes[0]
		dstVal := orig ^ srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		return mc + mc2 + dst.Store(c, []byte{dstVal})
	}
}

func or(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := bytes[0]
		bytes, mc2 := dst.Load(c)
		orig := bytes[0]
		dstVal := orig | srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		return mc + mc2 + dst.Store(c, []byte{dstVal})
	}
}

func cpl(c *cpu) go_gb.MC {
	val := c.getRegister(A)
	val[0] = ^val[0]
	c.setFlag(BitN, true)
	c.setFlag(BitH, true)
	return 0
}

func rc(c *cpu, dst Ptr, left, resetZ bool) go_gb.MC {
	bytes, mc := dst.Load(c)
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
	mc += dst.Store(c, []byte{b})
	if resetZ {
		c.setFlag(BitZ, false)
	} else {
		c.setFlag(BitZ, b == 0)
	}
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, old == 1)
	return mc
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
	bytes, mc := dst.Load(c)
	b := bytes[0]
	var old byte
	if left {
		old = b >> 7
		b <<= 1
	} else {
		old = b & 1
		b >>= 1
	}
	mc += dst.Store(c, []byte{b})
	c.setFlag(BitZ, b == 0)
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, old == 1)
	return mc
}

func sla(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		return sa(c, dst, true)
	}
}

func sra(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := dst.Load(c)
		b := bytes[0]
		var old byte
		msb := b & 0xF0
		old = b & 1
		b >>= 1
		b |= msb
		mc += dst.Store(c, []byte{b})
		c.setFlag(BitZ, b == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, old == 1)
		return mc
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
		bytes, mc := dst.Load(c)
		val := bytes[0]
		msn := val & 0xF0
		lsn := val & 0x0F
		result := (msn >> 4) | (lsn << 4)
		mc += dst.Store(c, []byte{result})

		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		return mc
	}
}
