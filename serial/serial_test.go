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

type mockExternalSerial struct {
	io.Reader
	IsReady bool
}

func (m *mockExternalSerial) Ready() bool {
	return m.IsReady
}

func TestSerial_StepMaster(t *testing.T) {
	in := bytes.NewBufferString("abc")
	m := &mockMemory{SB: '1', SC: 0x81}
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	s := NewSerial(&mockExternalSerial{in, false}, &buf, &buf2, m)

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

func TestSerial_StepPassive(t *testing.T) {
	in := bytes.NewBufferString("abc")
	m := &mockMemory{SB: '1', SC: 0x80}
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	externalSerial := &mockExternalSerial{in, false}
	s := NewSerial(externalSerial, &buf, &buf2, m)

	for i := 0; i < 23; i++ { // not ready - no bytes exchanged
		s.Step(go_gb.MC(i))
	}
	for i := 0; i <= 16; i++ { // one single byte exchanged
		externalSerial.IsReady = i%2 == 0
		s.Step(go_gb.MC(i))
	}
	m.SB = '2'
	m.SC = 0x80
	externalSerial.IsReady = true
	for i := 0; i < 16; i++ {
		if i == 8 {
			m.SC = 0x80
			m.SB = '3'
		}
		s.Step(go_gb.MC(i))
	}
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
