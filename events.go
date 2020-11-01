package go_gb

import (
	"strings"
	"sync"
)

var (
	Events = NewEventQueue(1000)
)

type event struct {
	data       string
	prev, next *event
}

func (e event) String() string {
	return e.data
}

type eventQueue struct {
	size, cap  int
	head, tail *event
	mutex      sync.Mutex
}

func (e *eventQueue) String() string {
	var sb strings.Builder
	for node := e.head; node != nil; node = node.next {
		sb.WriteString(node.String() + "\n")
	}
	return sb.String()
}

func NewEventQueue(cap int) *eventQueue {
	return &eventQueue{cap: cap}
}

func (e *eventQueue) Add(ev string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	newNode := &event{
		data: ev,
		prev: e.tail,
		next: nil,
	}
	if e.tail == nil {
		e.head = newNode
		e.tail = newNode
		e.size = 1
		return
	}

	if e.size == e.cap {
		_ = e.pop()
		e.tail.next = newNode
		e.tail = newNode
		return
	}
	e.tail.next = newNode
	e.tail = newNode
	e.size += 1
}

func (e *eventQueue) pop() *event {
	oldHead := e.head
	e.head = e.head.next
	e.head.prev = nil
	oldHead.next = nil
	return oldHead
}
