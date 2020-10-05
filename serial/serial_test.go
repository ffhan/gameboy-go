package serial

import (
	"bytes"
	"errors"
	"fmt"
	go_gb "go-gb"
	"io"
	"os"
	"testing"
)

func TestEmptyRead(t *testing.T) {
	b := []byte{0xFF}
	n, err := os.Stdin.Read(b)
	if err != nil {
		t.Log(err)
	}
	if errors.Is(err, io.EOF) {
		t.Log("found EOF")
	}
	fmt.Println(n, b)
}

type mockMemory struct {
	SB, SC byte
}

func (m mockMemory) ReadBytes(pointer, n uint16) []byte {
	panic("implement me")
}

func (m mockMemory) Read(pointer uint16) byte {
	switch pointer {
	case go_gb.SB:
		return m.SB
	}
	return m.SC
}

func (m *mockMemory) StoreBytes(pointer uint16, bytes []byte) {
	panic("implement me")
}

func (m *mockMemory) Store(pointer uint16, val byte) {
	switch pointer {
	case go_gb.SB:
		m.SB = val
	default:
		m.SC = val
	}
}

func TestSerial_Step(t *testing.T) {
	in := bytes.NewBufferString("abc")
	m := &mockMemory{SB: '1', SC: 0x81}
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	s := NewSerial(in, &buf, &buf2, m)

	const freq = 8192 / 4

	s.cycles = freq
	s.Step(8)
	m.SC = 0x81
	m.SB = '2'
	s.cycles = freq
	s.Step(4)
	s.cycles = freq
	s.Step(4)
	m.SC = 0x81
	m.SB = '3'
	s.cycles = freq
	s.Step(3)
	s.cycles = freq
	s.Step(4)
	s.cycles = freq
	s.Step(2)

	c := buf.String()
	t.Log(c)
	if c != "abc" {
		t.Error("invalid input")
	}
	c = buf2.String()
	t.Log(c)
	if c != "123" {
		t.Error("invalid output")
	}
}
