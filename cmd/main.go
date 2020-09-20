package main

import (
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	const (
		CpuFrequency = 4_194_304
	)
	mmu := memory.NewMMU()
	rom, err := ioutil.ReadFile("roms/Tetris (World) (Rev A).gb")
	if err != nil {
		panic(err)
	}
	mmu.Init(rom, go_gb.GB)
	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM())
	c := cpu.NewCpu(mmu, ppu)

	sleepTime := time.Second / CpuFrequency

	debug := cpu.NewDebugger(c, os.Stdout)

	for {
		debug.Step()
		time.Sleep(sleepTime)
	}
}
