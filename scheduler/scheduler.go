package scheduler

import (
	"fmt"
	go_gb "go-gb"
	"math"
	"sync/atomic"
	"time"
)

type scheduler struct {
	cpu go_gb.Cpu
	ppu go_gb.PPU
}

func NewScheduler(cpu go_gb.Cpu, ppu go_gb.PPU) *scheduler {
	return &scheduler{cpu: cpu, ppu: ppu}
}

func (s *scheduler) Run() {
	const (
		CpuFrequency float64 = 4_194_304 / 4
		PpuFrequency         = 59.7
	)
	ppuFreq := time.Duration(math.Round(float64(time.Second.Nanoseconds()) / PpuFrequency))

	fmt.Println(ppuFreq)

	var frames uint64

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			fps := atomic.LoadUint64(&frames)
			atomic.StoreUint64(&frames, 0)
			fmt.Printf("FPS: %d\n", fps)
		}
	}()

	for {
		start := time.Now()
		s.cpu.Step()
		if s.cpu.IME() && s.ppu.IsVBlank() {
			time.Sleep(time.Until(start.Add(ppuFreq)))
			time.Sleep(100 * time.Millisecond)
			atomic.AddUint64(&frames, 1)
		}
	}
}
