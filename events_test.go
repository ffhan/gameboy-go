package go_gb

import (
	"testing"
)

func testQueue(t *testing.T, queue *eventQueue, expectedSize int, expectedHeadData, expectedTailData string) {
	if queue.size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, queue.size)
	}
	if queue.tail.data != expectedTailData {
		t.Errorf("expected tail data to be \"%s\", got %s", expectedTailData, queue.tail.data)
	}
	if queue.head.data != expectedHeadData {
		t.Errorf("expected head data to be \"%s\", got %s", expectedHeadData, queue.head.data)
	}
}

func Test_eventQueue_Add(t *testing.T) {
	queue := NewEventQueue(3)
	queue.Add("1")
	if queue.head == nil || queue.head != queue.tail {
		t.Error("expected head and tail to be the same event")
	}
	testQueue(t, queue, 1, "1", "1")
	queue.Add("2")
	testQueue(t, queue, 2, "1", "2")
	queue.Add("3")
	testQueue(t, queue, 3, "1", "3")
	queue.Add("4")
	testQueue(t, queue, 3, "2", "4")
	queue.Add("5")
	testQueue(t, queue, 3, "3", "5")
}
