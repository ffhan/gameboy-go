package cpu

import (
	"errors"
	go_gb "go-gb"
)

func NOP(c *cpu) go_gb.MC {
	return 0
}

func STOP(c *cpu) go_gb.MC { // todo: halt until button pressed (joypad interrupt?)
	c.stop = true
	return 0
}

func halt(c *cpu) go_gb.MC {
	c.halt = true
	return 0
}

func scf(c *cpu) go_gb.MC {
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, true)
	return 0
}

func ccf(c *cpu) go_gb.MC {
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, !c.getFlag(BitC))
	return 0
}

func prefix(c *cpu) go_gb.MC {
	c.cbLookup = true
	return 0
}

var InvalidOpErr = errors.New("non-mapped operation called")

func invalid(c *cpu) go_gb.MC {
	panic(InvalidOpErr)
}

func di(c *cpu) go_gb.MC {
	c.diWaiting = 2
	return 0
}

func ei(c *cpu) go_gb.MC {
	c.eiWaiting = 2
	return 0
}
