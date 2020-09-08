package cpu

func rca(left bool) Instr {
	return func(c *cpu) error {
		b := c.getRegister(A)[0]
		var old byte
		if left {
			old = b >> 7
			b <<= 1
			b ^= old
		} else {
			old = b & 1
			b >>= 1
			b ^= old << 7
		}
		c.getRegister(A)[0] = b
		c.setFlag(BitZ, b == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, old == 1)
		return nil
	}
}

func ra(left bool) Instr {
	return func(c *cpu) error {
		b := c.getRegister(A)[0]
		var old byte
		carry := c.getFlag(BitC)
		if left {
			old = b >> 7
			b <<= 1
			if carry {
				b ^= 1
			}
		} else {
			old = b & 1
			b >>= 1
			if carry {
				b ^= 0x80
			}
		}
		c.getRegister(A)[0] = b
		c.setFlag(BitZ, b == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, old == 1)
		return nil
	}
}

func and(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := bytes[0]
		bytes = dst.Load(c)
		orig := bytes[0]
		dstVal := orig & srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, true)
		c.setFlag(BitC, false)
		dst.Store(c, []byte{dstVal})
		return nil
	}
}

func xor(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := bytes[0]
		bytes = dst.Load(c)
		orig := bytes[0]
		dstVal := orig ^ srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		dst.Store(c, []byte{dstVal})
		return nil
	}
}

func or(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := bytes[0]
		bytes = dst.Load(c)
		orig := bytes[0]
		dstVal := orig | srcVal

		c.setFlag(BitZ, dstVal == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, false)
		c.setFlag(BitC, false)
		dst.Store(c, []byte{dstVal})
		return nil
	}
}

func cpl(c *cpu) error {
	val := c.getRegister(A)
	val[0] = ^val[0]
	c.setFlag(BitN, true)
	c.setFlag(BitH, true)
	return nil
}
