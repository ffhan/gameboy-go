package scheduler

import (
	"fmt"
	go_gb "go-gb"
	"math"
	"os"
	"sync/atomic"
	"time"
)

type Controller interface {
	Wait() bool
}

type scheduler struct {
	cpu go_gb.Cpu
	ppu go_gb.PPU
	lcd go_gb.Display

	Frequency time.Duration
	Throttle  bool

	Controller Controller
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

	dump := false
	start := time.Now()

	go func() {
		const seconds = 1
		t := time.NewTicker(seconds * time.Second)
		for {
			select {
			case <-t.C:
				fps := atomic.LoadUint64(&frames)
				inst := atomic.LoadUint64(&cycles)
				atomic.StoreUint64(&frames, 0)
				atomic.StoreUint64(&cycles, 0)
				fmt.Printf("%s: FPS: %f\tCPU m cycles: %d, start: %s\n", time.Now().String(), float64(fps)/seconds, inst, start.String())
			case <-stopChan:
				return
			}
		}
	}()

	for {
		if s.Controller != nil && s.Controller.Wait() {
			start = time.Now()
		} // optionally wait (e.g. user debugging)
		mc := s.cpu.Step()
		atomic.AddUint64(&cycles, uint64(mc))
		if !s.lcd.IsDrawing() {
			if dump {
				go_gb.DumpDisplay(os.Stdout, s.lcd.(*go_gb.NopDisplay))
				dump = false
			}
			continue
		}
		if s.Throttle {
			start = start.Add(s.Frequency)
			time.Sleep(time.Until(start))
		}
		atomic.AddUint64(&frames, 1)
	}
}
