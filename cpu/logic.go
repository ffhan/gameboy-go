package cpu

import go_gb "go-gb"

func rlc(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		msb := a & 0x80
		a <<= 1
		if msb != 0 { // msb is set
			a |= 1
		} else {
			a &= 0xFE
		}

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, msb != 0)
		return 1 + mc
	}
}

func rlca(c *cpu) go_gb.MC {
	return rlc(rx(go_gb.A))(c)
}

func rl(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		msb := a & 0x80
		a <<= 1
		if c.getFlag(BitC) {
			a |= 1
		} else {
			a &= 0xFE
		}

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, msb != 0)
		return 1 + mc
	}
}

func rla(c *cpu) go_gb.MC {
	return rl(rx(go_gb.A))(c)
}

func rrc(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		lsb := a & 1
		a >>= 1
		if lsb != 0 { // msb is set
			a |= 0x80
		} else {
			a &= 0x7F
		}

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, lsb != 0)
		return 1 + mc
	}
}

func rrca(c *cpu) go_gb.MC {
	return rrc(rx(go_gb.A))(c)
}

func rr(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		lsb := a & 1
		a >>= 1
		if lsb != 0 {
			a |= 0x80
		} else {
			a &= 0x7F
		}
		if c.getFlag(BitC) { // msb is set
			a |= 0x80
		} else {
			a &= 0x7F
		}

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, lsb != 0)
		return 1 + mc
	}
}

func rra(c *cpu) go_gb.MC {
	return rr(rx(go_gb.A))(c)
}

func sla(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		msb := a & 0x80
		a <<= 1

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, msb != 0)
		return 1 + mc
	}
}

func sra(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		lsb := a & 1
		msb := a & 0x80
		a >>= 1
		a |= msb

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, lsb != 0)
		return 1 + mc
	}
}

func srl(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		lsb := a & 1
		a >>= 1

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, lsb != 0)
		return 1 + mc
	}
}

func swap(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		a := dst.Load(c, &mc)[0]
		msb := a & 0xF0
		lsb := a & 0x0F

		a = (lsb << 4) | (msb >> 4)

		dst.Store(c, []byte{a}, &mc)
		c.setFlag(BitZ, a == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		return 1 + mc
	}
}

func or(src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC

		dst := rx(go_gb.A)

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

func xor(src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC

		dst := rx(go_gb.A)

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

func cpl(c *cpu) go_gb.MC {
	val := c.GetRegister(go_gb.A)
	val[0] = ^val[0]
	c.setFlag(BitN, true)
	c.setFlag(BitH, true)
	return 0
}

func and(src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC

		dst := rx(go_gb.A)

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
