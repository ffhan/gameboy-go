package cpu

import (
	"go-gb"
)

func inc8bit(dst Ptr) Instr {
	return func(c *cpu) error {
		orig := dst.Load(c)[0]
		bytes := uint16(orig) + 1
		result := byte(bytes)
		c.setFlag(BitZ, bytes == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (bytes&0xF)+1 > 0xF) // TODO: probably wrong
		dst.Store(c, []byte{result})
		return nil
	}
}

func inc16bit(dst Ptr) Instr {
	return func(c *cpu) error {
		dst.Store(c, go_gb.MsbLsbBytes(go_gb.MsbLsb(dst.Load(c))+1, true))
		return nil
	}
}

func dec8bit(dst Ptr) Instr {
	return func(c *cpu) error {
		orig := dst.Load(c)[0]
		bytes := int16(orig) - 1
		result := byte(bytes)
		c.setFlag(BitZ, bytes == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (int16(orig)&0xF)-1 < 0) // TODO: probably wrong
		dst.Store(c, []byte{result})
		return nil
	}
}

func dec16bit(dst Ptr) Instr {
	return func(c *cpu) error {
		dst.Store(c, go_gb.MsbLsbBytes(uint16(int16(go_gb.MsbLsb(dst.Load(c)))-1), true))
		return nil
	}
}

func add8b(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := uint16(bytes[0])
		bytes = dst.Load(c)
		orig := uint16(bytes[0])
		dstVal := orig + srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (orig&0xF)+(srcVal&0xF) > 0xF)
		c.setFlag(BitC, (dstVal&0x100) != 0)
		dst.Store(c, []byte{result})
		return nil
	}
}

func add16b(dst, src Ptr) Instr {
	return func(c *cpu) error {
		dstVal := go_gb.MsbLsb(dst.Load(c))
		srcVal := go_gb.MsbLsb(src.Load(c))
		result := dstVal + srcVal
		c.setFlag(BitN, false)
		c.setFlag(BitH, ((dstVal&0xFFF)+(srcVal&0xFFF))&0x1000 != 0)
		c.setFlag(BitC, result > 0xFFFF)
		return nil
	}
}

func addSp(c *cpu) error {
	sp := sp()
	d := dx(8)
	orig := int16(go_gb.MsbLsb(sp.Load(c)))
	result := orig + int16(d.Load(c)[0])

	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < (orig&0xF))
	c.setFlag(BitC, (result&0xFF) < (orig&0xFF))
	return nil
}

func addHlSp(c *cpu) error {
	dst := rx(HL)
	orig := int16(c.sp)
	result := int16(c.sp) + int16(c.readOpcode())

	dst.Store(c, go_gb.MsbLsbBytes(uint16(result), true))

	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < (orig&0xF))
	c.setFlag(BitC, (result&0xFF) < (orig&0xFF))
	return nil
}

func adc8b(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := uint16(bytes[0])
		bytes = dst.Load(c)
		orig := uint16(bytes[0])
		carry := go_gb.BitToUint16(c.getFlag(BitC))
		dstVal := orig + srcVal + carry

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (orig&0xF)+(srcVal&0xF)+carry > 0xF)
		c.setFlag(BitC, (dstVal&0x100) != 0)
		dst.Store(c, []byte{result})
		return nil
	}
}

func sub(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := int16(bytes[0])
		bytes = dst.Load(c)
		orig := int16(bytes[0])
		dstVal := orig - srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF) < 0)
		c.setFlag(BitC, dstVal < 0)
		dst.Store(c, []byte{result})
		return nil
	}
}

func sbc(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := int16(bytes[0])
		bytes = dst.Load(c)
		orig := int16(bytes[0])
		carry := go_gb.BitToInt16(c.getFlag(BitC))
		dstVal := orig - srcVal - carry

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF)-carry < 0)
		c.setFlag(BitC, dstVal < 0)
		dst.Store(c, []byte{result})
		return nil
	}
}

func daa(c *cpu) error {
	panic("implement me")
}

func cp(dst, src Ptr) Instr {
	return func(c *cpu) error {
		bytes := src.Load(c)
		srcVal := int16(bytes[0])
		bytes = dst.Load(c)
		orig := int16(bytes[0])
		dstVal := orig - srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF) < 0)
		c.setFlag(BitC, dstVal < 0)
		dst.Store(c, []byte{result})
		return nil
	}
}
