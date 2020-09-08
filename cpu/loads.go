package cpu

import "go-gb"

func load(dst, src Ptr) Instr {
	return func(c *cpu) error {
		dst.Store(c, src.Load(c))
		return nil
	}
}

func ldHlSp(c *cpu) error {
	hl := rx(HL)
	n := go_gb.UnifyBytes(dx(8).Load(c))
	result := int16(go_gb.UnifyBytes(sp().Load(c))) + int16(n)
	hl.Store(c, go_gb.SeparateUint16(uint16(result)))

	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < int16(c.sp&0xF))
	c.setFlag(BitC, (result&0xFF) < int16(c.sp&0xFF))
	return nil
}

func loadIncrement(dst, src Ptr, increment bool) Instr {
	return func(c *cpu) error {
		if dst == nil {
			dst = rx(HL)
			if increment {
				defer dst.Store(c, go_gb.SeparateUint16(go_gb.UnifyBytes(dst.Load(c))+1))
			} else {
				defer dst.Store(c, go_gb.SeparateUint16(go_gb.UnifyBytes(dst.Load(c))-1))
			}
			return load(mPtr{dst}, src)(c)
		}
		if src == nil {
			src = rx(HL)
			if increment {
				defer src.Store(c, go_gb.SeparateUint16(go_gb.UnifyBytes(src.Load(c))+1))
			} else {
				defer src.Store(c, go_gb.SeparateUint16(go_gb.UnifyBytes(src.Load(c))-1))
			}
			return load(dst, mPtr{src})(c)
		}
		return nil
	}
}

func loadIo(dst, src Ptr) Instr {
	return func(c *cpu) error {
		n := uint16(src.Load(c)[0])
		dst.Store(c, c.memory.ReadBytes(go_gb.IOPortsStart+n, 1))
		return nil
	}
}

func pop(dst Ptr) Instr {
	return func(c *cpu) error {
		val := c.popStack(2)
		dst.Store(c, val)
		return nil
	}
}

func push(src Ptr) Instr {
	return func(c *cpu) error {
		val := src.Load(c)
		c.pushStack(val)
		return nil
	}
}