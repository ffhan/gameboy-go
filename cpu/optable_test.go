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
	cbTable := [...]go_gb.MC{
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,

		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,

		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,

		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	}

	c := initCpu(nil)

	checkCycles(0, t, c, optable[:], table[:])
	checkCycles(0xCB00, t, c, cbOptable, cbTable[:])
}

func checkCycles(prefix int, t *testing.T, c *cpu, table []Instr, expected []go_gb.MC) {
	for i, op := range table {
		failed := false
		var cyc go_gb.MC
		c.r[F] = 0
		func() {
			defer func() {
				err := recover()
				if err != nil {
					failed = true
					if e, ok := err.(error); ok && !errors.Is(e, InvalidOpErr) {
						t.Errorf("opcode %X panicked: %v\n", i+prefix, err)
					} else if !ok {
						t.Errorf("opcode %X panicked: %v\n", i+prefix, err)
					}
				}
			}()
			cyc = op(c)
		}()
		if !failed && expected[i] != cyc+1 {
			t.Errorf("opcode %X expected %d cycles, got %d\n", i+prefix, expected[i], cyc+1)
		}
	}
}
