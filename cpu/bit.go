package cpu

import go_gb "go-gb"

func bit(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		bitAddr, mc2 := bit.Load(c)
		c.setFlag(BitZ, go_gb.Bit(bytes[0], int(bitAddr[0])))
		c.setFlag(BitN, false)
		c.setFlag(BitH, true)
		return mc + mc2
	}
}

func set(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		b, mc := src.Load(c)
		bytes, mc2 := bit.Load(c)
		go_gb.Set(&b[0], int(bytes[0]), true)
		src.Store(c, []byte{b[0]})
		return mc + mc2
	}
}

func res(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		b, mc := src.Load(c)
		bytes, mc2 := bit.Load(c)
		go_gb.Set(&b[0], int(bytes[0]), false)
		src.Store(c, []byte{b[0]})
		return mc + mc2
	}
}
