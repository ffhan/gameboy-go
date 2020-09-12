package cpu

import "go-gb"

func load(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		return mc + dst.Store(c, bytes)
	}
}

func ldHlSp(c *cpu) go_gb.MC {
	hl := rx(HL)
	bytes, mc := dx(8).Load(c)
	n := go_gb.FromBytes(bytes)
	spVal, mc2 := sp().Load(c)
	result := int16(go_gb.FromBytes(spVal)) + int16(n)
	hl.Store(c, go_gb.LsbMsbBytes(uint16(result), true))

	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < int16(c.sp&0xF))
	c.setFlag(BitC, (result&0xFF) < int16(c.sp&0xFF))
	return mc + mc2
}

// e.g. LD (HL+), A
func ldHl(dst, src Ptr, increment bool) Instr {
	var deferFunc func(c *cpu)
	var loadFunc Instr
	if (dst == nil && src == nil) || (dst != nil && src != nil) {
		panic("invalid ldHl call")
	}
	if dst == nil {
		dst = rx(HL)
		offset := 1
		if !increment {
			offset = -1
		}
		deferFunc = func(c *cpu) {
			bytes, _ := dst.Load(c) // ignoring load because it's always from HL
			dst.Store(c, go_gb.LsbMsbBytes(uint16(int(go_gb.FromBytes(bytes))+offset), true))
		}
		loadFunc = load(mPtr{dst}, src)
	}
	if src == nil {
		src = rx(HL)
		offset := 1
		if !increment {
			offset = -1
		}
		deferFunc = func(c *cpu) {
			bytes, _ := src.Load(c) // ignoring load because it's always from HL
			src.Store(c, go_gb.LsbMsbBytes(uint16(int(go_gb.FromBytes(bytes))+offset), true))
		}
		loadFunc = load(dst, mPtr{src})
	}
	return func(c *cpu) go_gb.MC {
		defer deferFunc(c)
		return loadFunc(c)
	}
}

func pop(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		val, mc := c.popStack(2)
		mc += dst.Store(c, val)
		return mc
	}
}

func push(src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		val, mc := src.Load(c)
		mc += c.pushStack(val)
		return mc + 1
	}
}
