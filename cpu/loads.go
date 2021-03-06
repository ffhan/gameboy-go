package cpu

import "go-gb"

func load(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		dst.Store(c, bytes, &cycles)
		return cycles
	}
}

func ldSpHl(c *cpu) go_gb.MC {
	return load(sp(), rx(go_gb.HL))(c) + 1
}

func ldHlSp(c *cpu) go_gb.MC {
	var cycles go_gb.MC

	hl := rx(go_gb.HL)
	n := int8(go_gb.FromBytes(dx(8).Load(c, &cycles)))
	var result uint16
	spVal := go_gb.FromBytes(sp().Load(c, &cycles))
	if n < 0 {
		result = spVal - uint16(0xFF^byte(n))
	} else {
		result = spVal + uint16(n)
	}

	hl.Store(c, go_gb.ToBytes(result, true), &cycles)

	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < (c.sp&0xF))
	c.setFlag(BitC, (result&0xFF) < (c.sp&0xFF))
	return cycles + 1
}

// e.g. LD (HL+), A
func ldHl(dst, src Ptr, increment bool) Instr {
	var deferFunc func(c *cpu)
	var loadFunc Instr
	if (dst == nil && src == nil) || (dst != nil && src != nil) {
		panic("invalid ldHl call")
	}
	if dst == nil {
		dst = rx(go_gb.HL)
		offset := 1
		if !increment {
			offset = -1
		}
		deferFunc = func(c *cpu) {
			var mc go_gb.MC
			bytes := dst.Load(c, &mc) // ignoring load because it's always from HL
			dst.Store(c, go_gb.ToBytes(uint16(int(go_gb.FromBytes(bytes))+offset), true), &mc)
		}
		loadFunc = load(mPtr{dst}, src)
	}
	if src == nil {
		src = rx(go_gb.HL)
		offset := 1
		if !increment {
			offset = -1
		}
		deferFunc = func(c *cpu) {
			var mc go_gb.MC
			bytes := src.Load(c, &mc) // ignoring load because it's always from HL
			src.Store(c, go_gb.ToBytes(uint16(int(go_gb.FromBytes(bytes))+offset), true), &mc)
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
		var cycles go_gb.MC
		val := c.popStack(2, &cycles)
		dst.Store(c, val, &cycles)
		return cycles
	}
}

func push(src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var mc go_gb.MC
		val := src.Load(c, &mc)
		c.pushStack(val, &mc)
		return mc + 1
	}
}
