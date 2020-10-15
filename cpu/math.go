package cpu

import (
	"go-gb"
)

func inc8bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		origB := dst.Load(c, &cycles)
		orig := origB[0]
		bytes := uint16(orig) + 1
		result := byte(bytes)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (bytes&0xF)+1 > 0xF) // TODO: probably wrong
		dst.Store(c, []byte{result}, &cycles)
		return cycles
	}
}

func inc16bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := dst.Load(c, &cycles)
		dst.Store(c, go_gb.ToBytes(go_gb.FromBytes(bytes)+1, true), &cycles)
		return cycles + 1
	}
}

func dec8bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		origB := dst.Load(c, &cycles)
		orig := origB[0]
		bytes := uint16(orig) - 1
		result := byte(bytes)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-1 < 0) // TODO: probably wrong
		dst.Store(c, []byte{result}, &cycles)
		return cycles
	}
}

func dec16bit(dst Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := dst.Load(c, &cycles)
		dst.Store(c, go_gb.ToBytes(go_gb.FromBytes(bytes)-1, true), &cycles)
		return cycles + 1
	}
}

func add8b(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		srcVal := uint16(bytes[0])
		bytes = dst.Load(c, &cycles)
		orig := uint16(bytes[0])
		dstVal := orig + srcVal

		result := byte(dstVal)
		c.setFlag(BitH, (orig&0xF)+(srcVal&0xF) > 0xF)
		c.setFlag(BitC, dstVal > 0xFF)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		dst.Store(c, []byte{result}, &cycles)
		return cycles
	}
}

func add16b(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := dst.Load(c, &cycles)
		dstVal := uint(go_gb.FromBytes(bytes))
		bytes = src.Load(c, &cycles)
		srcVal := uint(go_gb.FromBytes(bytes))
		result := dstVal + srcVal
		c.setFlag(BitN, false)
		c.setFlag(BitH, ((dstVal&0xFFF)+(srcVal&0xFFF))&0x1000 != 0)
		c.setFlag(BitC, result > 0xFFFF)
		dst.Store(c, go_gb.ToBytes(uint16(result), true), &cycles)
		return cycles + 1
	}
}

func addSp(c *cpu) go_gb.MC {
	var cycles go_gb.MC
	sp := sp()
	d := dx(8)
	bytes := sp.Load(c, &cycles)
	orig := go_gb.FromBytes(bytes)
	val := int(int8(d.Load(c, &cycles)[0]))
	result := uint16(int(orig) + val)
	c.setFlag(BitZ, false)
	c.setFlag(BitN, false)
	c.setFlag(BitH, (result&0xF) < (orig&0xF))
	c.setFlag(BitC, (result&0xFF) < (orig&0xFF))
	sp.Store(c, go_gb.ToBytes(result, true), &cycles)
	return cycles + 2
}

func adc8b(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		srcVal := uint16(bytes[0])
		bytes = dst.Load(c, &cycles)
		orig := uint16(bytes[0])
		carry := go_gb.BitToUint16(c.getFlag(BitC))
		dstVal := orig + srcVal + carry

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, false)
		c.setFlag(BitH, (orig&0xF)+(srcVal&0xF)+carry > 0xF)
		c.setFlag(BitC, dstVal > 0xFF)
		dst.Store(c, []byte{result}, &cycles)
		return cycles
	}
}

func sub(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		srcVal := int16(bytes[0])
		bytes = dst.Load(c, &cycles)
		orig := int16(bytes[0])
		dstVal := orig - srcVal

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF) < 0)
		c.setFlag(BitC, dstVal < 0)
		dst.Store(c, []byte{result}, &cycles)
		return cycles
	}
}

func sbc(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		srcVal := int16(bytes[0])
		bytes = dst.Load(c, &cycles)
		orig := int16(bytes[0])
		carry := go_gb.BitToInt16(c.getFlag(BitC))
		dstVal := orig - srcVal - carry

		result := byte(dstVal)
		c.setFlag(BitZ, result == 0)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF)-carry < 0)
		c.setFlag(BitC, dstVal < 0)
		dst.Store(c, []byte{result}, &cycles)
		return cycles
	}
}

func daa(c *cpu) go_gb.MC {
	var mc go_gb.MC
	registerA := rx(go_gb.A)
	reg := registerA.Load(c, &mc)

	add := !c.getFlag(BitN)
	carry := c.getFlag(BitC)
	halfCarry := c.getFlag(BitH)

	result := uint16(reg[0])
	var correction uint16
	if carry {
		correction = 0x60
	}
	if halfCarry || add && ((result&0x0F) > 9) {
		correction |= 0x06
	}
	if carry || add && (result > 0x99) {
		correction |= 0x60
	}
	if add {
		result += correction
	} else {
		result -= correction
	}
	if ((correction << 2) & 0x100) != 0 {
		c.setFlag(BitC, true)
	}
	c.setFlag(BitH, false)
	storedResult := byte(result & 0xFF)
	registerA.Store(c, []byte{storedResult}, &mc)
	c.setFlag(BitZ, storedResult == 0)
	return mc
}

func cp(dst, src Ptr) Instr {
	return func(c *cpu) go_gb.MC {
		var cycles go_gb.MC
		bytes := src.Load(c, &cycles)
		srcVal := uint16(bytes[0])
		bytes = dst.Load(c, &cycles)
		orig := uint16(bytes[0])

		c.setFlag(BitZ, orig == srcVal)
		c.setFlag(BitN, true)
		c.setFlag(BitH, (orig&0xF)-(srcVal&0xF) < 0)
		c.setFlag(BitC, orig < srcVal)
		return cycles
	}
}
