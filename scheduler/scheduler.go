package scheduler

import (
	"fmt"
	go_gb "go-gb"
	"math"
	"sync/atomic"
	"time"
)

type scheduler struct {
	cpu      go_gb.Cpu
	ppu      go_gb.PPU
	lcd      go_gb.Display
	stoppers map[uint16]bool
}

func NewScheduler(cpu go_gb.Cpu, ppu go_gb.PPU, lcd go_gb.Display) *scheduler {
	return &scheduler{cpu: cpu, ppu: ppu, lcd: lcd, stoppers: make(map[uint16]bool)}
}

func (s *scheduler) AddStopper(addr uint16) {
	s.stoppers[addr] = true
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

	for {
		if _, ok := s.stoppers[s.cpu.PC()]; ok {
			return
		}
		start := time.Now()
		mc := s.cpu.Step()
		atomic.AddUint64(&cycles, uint64(mc))
		if s.lcd.IsDrawing() {
			time.Sleep(time.Until(start.Add(ppuFreq)))
			atomic.AddUint64(&frames, 1)
			//fmt.Println(time.Now().Sub(start))
		}
	}
}
