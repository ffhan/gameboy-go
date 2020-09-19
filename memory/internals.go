package memory

type mmap struct {
	start, end uint16
	memory     []byte
}

func newMmap(start uint16, end uint16, memory []byte) *mmap {
	return &mmap{start: start, end: end, memory: memory}
}

func (m *mmap) ReadBytes(pointer, n uint16) []byte {
	i := pointer - m.start
	return m.memory[i : i+n]
}

func (m *mmap) Read(pointer uint16) byte {
	i := pointer - m.start
	return m.memory[i]
}

func (m *mmap) StoreBytes(pointer uint16, bytes []byte) {
	i := pointer - m.start
	copy(m.memory[i:i+uint16(len(bytes))], bytes)
}

func (m *mmap) Store(pointer uint16, val byte) {
	m.memory[pointer-m.start] = val
}

func (m *mmap) LoadRom(bytes []byte) int {
	n := len(m.memory)
	copy(m.memory, bytes[:n])
	return n
}

type lockedMemory struct {
}

func (l *lockedMemory) ReadBytes(pointer, n uint16) []byte {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = 0xFF
	}
	return bytes
}

func (l *lockedMemory) Read(pointer uint16) byte {
	return l.ReadBytes(pointer, 1)[0]
}

func (l *lockedMemory) StoreBytes(pointer uint16, bytes []byte) {
}

func (l *lockedMemory) Store(pointer uint16, val byte) {
}
