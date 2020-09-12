package cpu

import (
	"go-gb"
)

func inc8bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		origB, mc := dst.Load(c)
		orig := origB[0]
		bytes := uint16(orig) + 1
		result := byte(bytes)
		c.setFlag(BitZ, bytes == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (bytes&0xF)+1 > 0xF) // TODO: probably wrong
		return mc + dst.Store(c, []byte{result})
	}
}

func inc16bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := dst.Load(c)
		return mc + dst.Store(c, go_gb.LsbMsbBytes(go_gb.FromBytes(bytes)+1, true))
	}
}

func dec8bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		origB, mc := dst.Load(c)
		orig := origB[0]
		bytes := int16(orig) - 1
		result := byte(bytes)
		c.setFlag(BitZ, bytes == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (int16(orig)&0xF)-1 < 0) // TODO: probably wrong
		return mc + dst.Store(c, []byte{result})
	}
}

func dec16bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := dst.Load(c)
		return mc + dst.Store(c, go_gb.LsbMsbBytes(uint16(int16(go_gb.FromBytes(bytes))-1), true))
	}
}

func add8b(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := uint16(bytes[0])
		bytes, mc2 := dst.Load(c)
		orig := uint16(bytes[0])
		dstVal := orig + srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (orig&0xF)+(srcVal&0xF) > 0xF)
		c.setFlag(BitC, (dstVal&0x100) != 0)
		return mc + mc2 + dst.Store(c, []byte{result})
	}
}

func add16b(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := dst.Load(c)
		dstVal := go_gb.FromBytes(bytes)
		bytes, mc2 := src.Load(c)
		srcVal := go_gb.FromBytes(bytes)
		result := dstVal + srcVal
		c.setFlag(BitN, false)
		c.setFlag(BitH, ((dstVal&0xFFF)+(srcVal&0xFFF))&0x1000 != 0)
		c.setFlag(BitC, result > 0xFFFF)
		return mc + mc2 + dst.Store(c, go_gb.LsbMsbBytes(result, true)) + 1
	}
}

func addSp(c *cpu) go_gb.MC {
	sp := sp()
	d := dx(8)
	bytes, mc := sp.Load(c)
	orig := int16(go_gb.FromBytes(bytes))
	val, mc2 := d.Load(c)
	result := orig + int16(val[0])

	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < (orig&0xF))
	c.setFlag(BitC, (result&0xFF) < (orig&0xFF))
	return mc + mc2 + sp.Store(c, go_gb.LsbMsbBytes(uint16(result), true))
}

func adc8b(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := uint16(bytes[0])
		bytes, mc2 := dst.Load(c)
		orig := uint16(bytes[0])
		carry := go_gb.BitToUint16(c.getFlag(BitC))
		dstVal := orig + srcVal + carry

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (orig&0xF)+(srcVal&0xF)+carry > 0xF)
		c.setFlag(BitC, (dstVal&0x100) != 0)
		return mc + mc2 + dst.Store(c, []byte{result})
	}
}

func sub(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := int16(bytes[0])
		bytes, mc2 := dst.Load(c)
		orig := int16(bytes[0])
		dstVal := orig - srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF) < 0)
		c.setFlag(BitC, dstVal < 0)
		return mc + mc2 + dst.Store(c, []byte{result})
	}
}

func sbc(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := int16(bytes[0])
		bytes, mc2 := dst.Load(c)
		orig := int16(bytes[0])
		carry := go_gb.BitToInt16(c.getFlag(BitC))
		dstVal := orig - srcVal - carry

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF)-carry < 0)
		c.setFlag(BitC, dstVal < 0)
		return mc + mc2 + dst.Store(c, []byte{result})
	}
}

func daa(c *cpu) go_gb.MC {
	panic("implement me")
}

func cp(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		bytes, mc := src.Load(c)
		srcVal := int16(bytes[0])
		bytes, mc2 := dst.Load(c)
		orig := int16(bytes[0])
		dstVal := orig - srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF) < 0)
		c.setFlag(BitC, dstVal < 0)
		return mc + mc2 + dst.Store(c, []byte{result})
	}
}
