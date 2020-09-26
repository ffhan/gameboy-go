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
	lcd go_gb.Display
}

func NewScheduler(cpu go_gb.Cpu, ppu go_gb.PPU, lcd go_gb.Display) *scheduler {
	return &scheduler{cpu: cpu, ppu: ppu, lcd: lcd}
}

func (s *scheduler) Run() {
	fmt.Println("started sched")
	const (
		CpuFrequency float64 = 4_194_304 / 4
		PpuFrequency         = 59.7
	)
	ppuFreq := time.Duration(math.Round(float64(time.Second.Nanoseconds()) / PpuFrequency))

	fmt.Println(ppuFreq)

	var frames uint64
	var cycles uint64

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			fps := atomic.LoadUint64(&frames)
			inst := atomic.LoadUint64(&cycles)
			atomic.StoreUint64(&frames, 0)
			atomic.StoreUint64(&cycles, 0)
			fmt.Printf("FPS: %d\tCPU m cycles: %d\n", fps, inst)
		}
	}()

	for {
		start := time.Now()
		mc := s.cpu.Step()
		atomic.AddUint64(&cycles, uint64(mc))
		if s.lcd.IsDrawing() {
			time.Sleep(time.Until(start.Add(ppuFreq)))
			atomic.AddUint64(&frames, 1)
			fmt.Println(time.Now().Sub(start))
		}
	}
}
