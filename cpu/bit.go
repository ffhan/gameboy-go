package cpu

import go_gb "go-gb"

func bit(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) error {
		c.setFlag(BitZ, go_gb.Bit(src.Load(c)[0], int(bit.Load(c)[0])))
		c.setFlag(BitN, false)
		c.setFlag(BitH, true)
		return nil
	}
}

func set(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) error {
		b := src.Load(c)[0]
		go_gb.Set(&b, int(bit.Load(c)[0]), true)
		src.Store(c, []byte{b})
		return nil
	}
}

func res(bit hardcoded, src Ptr) Instr {
	return func(c *cpu) error {
		b := src.Load(c)[0]
		go_gb.Set(&b, int(bit.Load(c)[0]), false)
		src.Store(c, []byte{b})
		return nil
	}
}
