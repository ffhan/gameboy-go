package cpu

import (
	"errors"
)

func NOP(c *cpu) error {
	return nil
}

func STOP(c *cpu) error { // todo: halt until button pressed (joypad interrupt?)
	panic("implement me")
}

func halt(c *cpu) error {
	c.halt = true
	return nil
}

func scf(c *cpu) error {
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, true)
	return nil
}

func ccf(c *cpu) error {
	c.setFlag(BitN, false)
	c.setFlag(BitH, false)
	c.setFlag(BitC, !c.getFlag(BitC))
	return nil
}

func prefix(c *cpu) error {
	return cbOptable[c.readOpcode()](c)
}

func invalid(c *cpu) error {
	return errors.New("non-mapped operation called")
}

func di(c *cpu) error {
	c.diWaiting = 2
	return nil
}

func ei(c *cpu) error {
	c.eiWaiting = 2
	return nil
}
