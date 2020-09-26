package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/scheduler"
	"os"
	"testing"
)

func TestRunning(t *testing.T) {
	mmu := memory.NewMMU()
	file, err := os.Open("roms/Tetris (World) (Rev A).gb")
	if err != nil {
		panic(err)
	}

	game, err := go_gb.LoadGame(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(game)

	mmu.Init(game.Rom, go_gb.GB)

	lcd := go_gb.NewNopDisplay()

	//mmuD := memory.NewDebugger(mmu, os.Stdout)
	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), mmu.IO(), lcd)
	realCpu := cpu.NewCpu(mmu, ppu)
	//c := cpu.NewDebugger(realCpu, os.Stdout)

	//sysD := debugger.NewSystemDebugger(c, mmuD)
	//c.Debug(false)

	sched := scheduler.NewScheduler(realCpu, ppu, lcd)
	sched.AddStopper(0x100)
	sched.Run()
}
