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

	Frequency time.Duration
	Throttle  bool
}

func NewScheduler(cpu go_gb.Cpu, ppu go_gb.PPU, lcd go_gb.Display) *scheduler {
	const PpuFrequency = 59.7
	ppuFreq := time.Duration(math.Round(float64(time.Second.Nanoseconds()) / PpuFrequency))
	return &scheduler{cpu: cpu, ppu: ppu, lcd: lcd, Frequency: ppuFreq, Throttle: true}
}

func (s *scheduler) Run() {
	fmt.Println(s.Frequency)

	var frames uint64
	var cycles uint64

	stopChan := make(chan bool)
	defer func() {
		stopChan <- true
	}()

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				fps := atomic.LoadUint64(&frames)
				inst := atomic.LoadUint64(&cycles)
				atomic.StoreUint64(&frames, 0)
				atomic.StoreUint64(&cycles, 0)
				fmt.Printf("FPS: %d\tCPU m cycles: %d\n", fps, inst)
			case <-stopChan:
				return
			}
		}
	}()

	start := time.Now()
	for {
		mc := s.cpu.Step()
		atomic.AddUint64(&cycles, uint64(mc))
		if !s.lcd.IsDrawing() {
			continue
		}
		if s.Throttle {
			start = start.Add(s.Frequency)
			time.Sleep(time.Until(start))
		}
		atomic.AddUint64(&frames, 1)
	}
}
