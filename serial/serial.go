package serial

import (
	go_gb "go-gb"
	"io"
)

type serial struct {
	in io.Reader

	inWriter  io.ReadWriter
	outWriter io.ReadWriter

	inByte  *byte
	outByte byte
	counter byte

	cycles go_gb.MC

	memory go_gb.Memory
}

func NewSerial(in io.Reader, inWriter, outWriter io.ReadWriter, memory go_gb.Memory) *serial {
	return &serial{in: in, inWriter: inWriter, outWriter: outWriter, memory: memory}
}

func (s *serial) Stream() io.Reader {
	return s.inWriter
}

// this is as close to hardware spec as I can think of
func (s *serial) Step(mc go_gb.MC) { // todo: cgb clock speed for bit
	const (
		normalClockFreq = 8192 / 4 // 8192 Hz is 8192 T cycles or 8192 / 4 M cycles
	)
	sc := s.memory.Read(go_gb.SC)
	transferRequestedOrInProgress := go_gb.Bit(sc, 7)
	if !transferRequestedOrInProgress {
		return
	}
	sb := s.memory.Read(go_gb.SB)
	isMaster := go_gb.Bit(sc, 0)

	if s.inByte == nil {
		temp := []byte{0xFF}
		if s.in != nil {
			_, _ = s.in.Read(temp)
		} // if EOF temp will be 255 - that's okay
		s.inByte = &temp[0]
	}

	if isMaster {
		s.cycles += mc

		if s.cycles >= normalClockFreq {
			s.cycles -= normalClockFreq

			for i := 0; i < int(mc) && s.counter < 8; i++ {
				bit := (sb & 0x80) >> 7
				s.outByte = (s.outByte << 1) | bit
				sb = (sb << 1) | ((*s.inByte & 0x80) >> 7)
				*s.inByte <<= 1

				s.counter += 1
			}
			s.memory.Store(go_gb.SB, sb)
			if s.counter >= 8 { // transfer end
				defer func() {
					s.counter = 0
					s.inByte = nil
					s.outByte = 0
					s.cycles = 0
				}()

				// write input & output byte
				if s.inWriter != nil {
					_, _ = s.inWriter.Write([]byte{sb})
				}
				if s.outWriter != nil {
					_, _ = s.outWriter.Write([]byte{s.outByte})
				}
				// set SC transfer finish, enable serial interrupt
				s.memory.Store(go_gb.SC, 0)
				go_gb.Update(s.memory, go_gb.IF, func(b byte) byte {
					go_gb.Set(&b, int(go_gb.BitSerial), true)
					return b
				})
			}
		}
	} else {
		panic("unimplemented") // todo: figure out how to detect external clock
	}
}
