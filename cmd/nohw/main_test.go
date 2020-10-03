package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/scheduler"
	"go-gb/timer"
	"os"
	"testing"
)

func TestRunning(t *testing.T) {
	logs, err := os.Create("output.log")
	if err != nil {
		panic(err)
	}
	defer logs.Close()

	mmu := memory.NewMMU()
	file, err := os.Open("roms/Tetris (World) (Rev A).gb")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	game, err := go_gb.LoadGame(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(game)

	mmu.Init(game.Rom, go_gb.GB)

	lcd := go_gb.NewNopDisplay()

	divTimer := timer.NewDivTimer(mmu.IO())
	timer := timer.NewTimer(mmu.IO())

	//mmuD := memory.NewDebugger(mmu, os.Stdout)
	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), mmu.IO(), lcd)
	mmuD := memory.NewDebugger(mmu, logs)
	mmuD.Debug(false)

	realCpu := cpu.NewCpu(mmuD, ppu, timer, divTimer)

	defer func() {
		err := recover()
		switch err.(type) {
		case error:
			fmt.Printf("PC: %X -> err: %v\n", realCpu.PC(), err)
		case string:
			fmt.Printf("PC: %X -> err: %s\n", realCpu.PC(), err)
		}
		panic(err)
	}()

	debugger := cpu.NewDebugger(realCpu, logs)
	debugger.Debug(true)
	debugger.PrintInstructionNames(true)
	sched := scheduler.NewScheduler(debugger, ppu, lcd)
	//sched.AddStopper(0x100)
	sched.Run()
}
