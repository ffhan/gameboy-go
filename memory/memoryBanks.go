package memory

import "fmt"

type bank struct {
	memory     []byte
	partSize   uint
	numOfParts uint
}

func newBank(numOfParts, totalSize uint) *bank {
	return &bank{
		memory:     make([]byte, totalSize),
		partSize:   totalSize / numOfParts,
		numOfParts: numOfParts,
	}
}

func (b *bank) address(bank, pointer uint16) uint {
	result := b.partSize*uint(bank) + uint(pointer)
	if result >= uint(len(b.memory)) {
		panic(fmt.Sprintf("want bank %d for pointer %X on bank (%d, %d, %d)\n", bank, pointer, len(b.memory), b.numOfParts, b.partSize))
	}
	return result
}

func (b *bank) ReadBytes(bank, pointer, n uint16) []byte {
	address := b.address(bank, pointer)
	return b.memory[address : address+uint(n)]
}

func (b *bank) Read(bank, pointer uint16) byte {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("panic at the disco", len(b.memory), cap(b.memory), b.partSize, b.numOfParts)
			panic(err)
		}
	}()
	address := b.address(bank, pointer)
	return b.memory[address]
}

func (b *bank) StoreBytes(bank, pointer uint16, bytes []byte) {
	address := b.address(bank, pointer)
	copy(b.memory[address:address+uint(len(bytes))], bytes)
}

func (b *bank) Store(bank, pointer uint16, val byte) {
	b.memory[b.address(bank, pointer)] = val
}

func (b *bank) LoadRom(bytes []byte) int {
	n := len(b.memory)
	if n > len(bytes) {
		n = len(bytes)
	}
	copy(b.memory, bytes[:n])
	return n
}
