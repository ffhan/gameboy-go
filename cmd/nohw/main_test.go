package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/scheduler"
	"go-gb/serial"
	"go-gb/timer"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestRunning(t *testing.T) {
	logs, err := os.Create("output.log")
	if err != nil {
		panic(err)
	}
	defer logs.Close()

	mmu := memory.NewMMU()
	file, err := os.Open("roms/Street Fighter II (UE) [S][!].gb")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	game, err := go_gb.LoadGame(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(game)

	mmu.Init(game.Rom, go_gb.GB, go_gb.NOPJoypad)

	lcd := go_gb.NewNopDisplay()

	divTimer := timer.NewDivTimer(mmu.IO())
	timer := timer.NewTimer(mmu.IO())

	//mmuD := memory.NewDebugger(mmu, os.Stdout)
	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), mmu.IO(), lcd)
	mmuD := memory.NewDebugger(mmu, logs)
	mmuD.Debug(false)

	serialFile, err := os.Create("serial.txt")
	if err != nil {
		panic(err)
	}

	serialPort := serial.NewSerial(serial.NopSerial, nil, serialFile, mmu.IO())

	realCpu := cpu.NewCpu(mmuD, ppu, timer, divTimer, serialPort)

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

	sig := make(chan os.Signal, 10)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)

	debugger := cpu.NewDebugger(realCpu, logs, cpu.NewInstructionQueue(100000))
	debugger.PrintEveryCycle = false
	debugger.Debug(true)
	debugger.PrintInstructionNames(true)
	sched := scheduler.NewScheduler(debugger, ppu, lcd)
	sched.Throttle = false
	//sched.AddStopper(0x100)

	go func() {
		<-sig
		debugger.Dump()
		fmt.Println("dumped instr queue")
	}()

	sched.Run()
}
