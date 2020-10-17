package serial

import (
	"errors"
	go_gb "go-gb"
	"io"
)

type ExternalSerial interface {
	io.Reader
	Ready() bool // determines if the next bit should be read.
}

var NopSerial = &noopSerial{}

type noopSerial struct {
}

func (s *noopSerial) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (s *noopSerial) Ready() bool {
	return true
}

type serial struct {
	in ExternalSerial

	inWriter  io.ReadWriter
	outWriter io.ReadWriter

	inByte  *byte
	outByte byte
	counter byte

	cycles go_gb.MC

	memory go_gb.Memory

	readyForSending bool
}

func NewSerial(in ExternalSerial, inWriter, outWriter io.ReadWriter, memory go_gb.Memory) *serial {
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

			s.readyForSending = true // this seems wrong
			for i := 0; i < int(mc) && s.counter < 8; i++ {
				sb = s.process(sb)
			}
			s.afterProcess(sb)
		}
	} else if s.in != nil {
		if s.in.Ready() {
			sb = s.process(sb)
			s.afterProcess(sb)
		}
	}
}

func (s *serial) process(sb byte) byte {
	bit := (sb & 0x80) >> 7
	s.outByte = (s.outByte << 1) | bit
	sb = (sb << 1) | ((*s.inByte & 0x80) >> 7)
	*s.inByte <<= 1

	s.counter += 1
	return sb
}

func (s *serial) afterProcess(sb byte) {
	s.memory.Store(go_gb.SB, sb)
	if s.counter >= 8 { // transfer end
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

		s.counter = 0
		s.inByte = nil
		s.outByte = 0
		s.cycles = 0
		s.readyForSending = false
	}
}

func (s *serial) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, errors.New("empty byte slice")
	}
	p[0] = s.memory.Read(go_gb.SB)
	return 1, nil
}

func (s *serial) Ready() bool {
	return s.readyForSending
}
