package cpu

import (
	"errors"
	go_gb "go-gb"
	"testing"
)

func TestOptable(t *testing.T) {
	if len(optable) != 256 {
		t.Fatal("invalid table size")
	}
	if len(cbOptable) != 256 {
		t.Fatal("invalid CB table size")
	}
}

func TestOptableCycles(t *testing.T) {
	table := [...]go_gb.MC{
		1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1,
		1, 3, 2, 2, 1, 1, 2, 1, 3, 2, 2, 2, 1, 1, 2, 1,

		3, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1,
		3, 3, 2, 2, 3, 3, 3, 1, 2, 2, 2, 2, 1, 1, 2, 1,

		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		2, 2, 2, 2, 2, 2, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1,

		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,

		5, 3, 4, 4, 6, 4, 2, 4, 2, 4, 3, 1, 3, 6, 2, 4,
		5, 3, 4, 0, 6, 4, 2, 4, 2, 4, 3, 0, 3, 0, 2, 4,

		3, 3, 2, 0, 0, 4, 2, 4, 4, 1, 4, 0, 0, 0, 2, 4,
		3, 3, 2, 1, 0, 4, 2, 4, 3, 2, 4, 1, 0, 0, 2, 4,
	}

	c := NewCpu()

	for i, op := range optable {
		failed := false
		var cyc go_gb.MC
		c.r[F] = 0
		func() {
			defer func() {
				err := recover()
				if err != nil {
					failed = true
					if e, ok := err.(error); ok && !errors.Is(e, InvalidOpErr) {
						t.Errorf("opcode %X panicked: %v\n", i, err)
					} else if !ok {
						t.Errorf("opcode %X panicked: %v\n", i, err)
					}
				}
			}()
			cyc = op(c)
		}()
		if !failed && table[i] != cyc+1 {
			t.Errorf("opcode %X expected %d cycles, got %d\n", i, table[i], cyc+1)
		}
	}
}
