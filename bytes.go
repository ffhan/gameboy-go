package go_gb

import "fmt"

func Bit(b byte, i int) bool {
	return (b>>i)&1 == 1
}

func Set(b *byte, i int, val bool) {
	if val {
		*b |= 1 << i
	} else {
		*b &= ^(1 << i)
	}
}

func Toggle(b *byte, i int) bool {
	*b ^= 1 << i
	return (*b>>i)&1 == 1
}

func FromBytes(b []byte) uint16 {
	if len(b) > 2 {
		panic(fmt.Errorf("%v cannot be transformed to uint16", b))
	}
	if len(b) == 2 {
		return (uint16(b[1]) << 8) | uint16(b[0])
	}
	return uint16(b[0])
}

func ToBytes(val uint16, word bool) []byte {
	if word {
		return []byte{byte(val & 0xFF), byte((val & 0xFF00) >> 8)}
	}
	return []byte{byte(val & 0xFF)}
}

func FromBytesReverse(b []byte) uint16 {
	if len(b) > 2 {
		panic(fmt.Errorf("%v cannot be transformed to uint16", b))
	}
	if len(b) == 2 {
		return (uint16(b[1]) << 8) | uint16(b[0])
	}
	return uint16(b[0])
}

func ToBytesReverse(val uint16, word bool) []byte {
	if word {
		return []byte{byte((val & 0xFF00) >> 8), byte(val & 0xFF)}
	}
	return []byte{byte(val & 0xFF)}
}

func BitToByte(val bool) byte {
	if val {
		return 1
	}
	return 0
}

func BitToUint16(val bool) uint16 {
	if val {
		return 1
	}
	return 0
}

func BitToInt16(val bool) int16 {
	if val {
		return 1
	}
	return 0
}

func Reverse(b []byte) {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - 1 - i
		tmp := b[i]
		b[i] = b[j]
		b[j] = tmp
	}
}
