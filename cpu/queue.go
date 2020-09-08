package cpu

type instructionNode struct {
	exec       Instr
	prev, next *instructionNode
}

type instrQueue struct {
	head, tail *instructionNode
	size       int
}

func (q *instrQueue) Push(in Instr) {
	node := &instructionNode{
		exec: in,
		prev: nil,
		next: nil,
	}
	if q.head == nil && q.tail == nil {
		q.head = node
	}
	q.tail.next = node
	q.tail = node
	q.size += 1
}

func (q *instrQueue) Size() int {
	return q.size
}

// a -> b -> c
// a -> c
// a
func (q *instrQueue) Pop() Instr {
	instr := q.head
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.size -= 1
	return instr.exec
}
