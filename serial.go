package go_gb

import "io"

const (
	SB uint16 = 0xFF01 // Serial transfer data (R/W)
	SC uint16 = 0xFF02 // Serial transfer control (R/W)
)

type Serial interface {
	Stream() io.Reader
	Step(mc MC)
}
