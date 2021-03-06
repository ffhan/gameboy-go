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
)

func main() {
	logs, err := os.Create("output.log")
	if err != nil {
		panic(err)
	}
	defer logs.Close()

	mmu := memory.NewMMU()
	file, err := os.Open("roms/gb-test-roms-master/cpu_instrs/cpu_instrs.gb")
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

	debugger := cpu.NewDebugger(realCpu, logs, cpu.NewInstructionQueue(100000))
	debugger.PrintEveryCycle = false
	debugger.Debug(true)
	debugger.PrintInstructionNames(true)

	defer func() {
		err := recover()
		switch err.(type) {
		case error:
			fmt.Printf("PC: %X -> err: %v\n", realCpu.PC(), err)
		case string:
			fmt.Printf("PC: %X -> err: %s\n", realCpu.PC(), err)
		}
		_, _ = debugger.Dump()
		panic(err)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	sched := scheduler.NewScheduler(debugger, ppu, lcd)
	sched.Throttle = false
	//sched.AddStopper(0x100)

	go func() {
		fmt.Println("waiting for the signal", os.Getpid())
		signal := <-sig
		fmt.Println("received a signal", signal.String())
		debugger.Dump()
		fmt.Println("dumped instr queue")
		os.Exit(0)
	}()

	sched.Run()
}
