package cpu

import go_gb "go-gb"

func bit(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		c.setFlag(BitZ, !go_gb.Bit(src.Load(c, &cycles)[0], int(bit.Load(c, &cycles)[0])))
		c.setFlag(BitN, false)
		c.setFlag(BitH, true)
		return cycles + 1
	}
}

func set(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		b := src.Load(c, &cycles)
		bytes := bit.Load(c, &cycles)

		go_gb.Set(&b[0], int(bytes[0]), true)
		src.Store(c, []byte{b[0]}, &cycles)
		return cycles + 1
	}
}

func res(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		b := src.Load(c, &cycles)
		bytes := bit.Load(c, &cycles)
		go_gb.Set(&b[0], int(bytes[0]), false)
		src.Store(c, []byte{b[0]}, &cycles)
		return cycles + 1
	}
}
