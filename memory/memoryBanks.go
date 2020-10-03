package memory

type bank struct {
	memory     []byte
	partSize   uint
	numOfParts uint
}

func newBank(numOfParts, partSize uint) *bank {
	return &bank{
		memory:     make([]byte, int(numOfParts)*int(partSize)),
		partSize:   partSize,
		numOfParts: numOfParts,
	}
}

func (b *bank) address(bank, pointer uint16) uint {
	return b.partSize*uint(bank) + uint(pointer)
}

func (b *bank) ReadBytes(bank, pointer, n uint16) []byte {
	address := b.address(bank, pointer)
	return b.memory[address : address+uint(n)]
}

func (b *bank) Read(bank, pointer uint16) byte {
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
