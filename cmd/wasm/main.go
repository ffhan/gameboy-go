package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/debugger"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/scheduler"
	"go-gb/wasm"
	"os"
	"sync"
	"syscall/js"
)

func run() (debugger.CpuDebugger, debugger.MemoryDebugger, go_gb.PPU, go_gb.Display) {
	mmu := memory.NewMMU()
	rom := make([]byte, 2*1<<20)
	n := js.CopyBytesToGo(rom, js.Global().Get("document").Get("rom"))
	mmu.Init(rom[:n], go_gb.GB)

	fmt.Println("initialized mmu")

	lcd := wasm.NewWasmDisplay()

	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), lcd)
	mmuD := memory.NewDebugger(mmu, os.Stdout)
	c := cpu.NewCpu(mmuD, ppu)

	return cpu.NewDebugger(c, os.Stdout), mmuD, ppu, lcd
}

type Runner interface {
	Run()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var cpu debugger.CpuDebugger
	var ppu go_gb.PPU
	var mmu debugger.MemoryDebugger
	var lcd go_gb.Display

	var systemDebugger debugger.Debugger
	var sched Runner

	js.Global().Set("run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cpu, mmu, ppu, lcd = run()
		systemDebugger = debugger.NewSystemDebugger(cpu, mmu)
		systemDebugger.Debug(false)
		sched = scheduler.NewScheduler(cpu, ppu, lcd)
		return nil
	}))
	js.Global().Set("step", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cpu.Step()
		return nil
	}))
	js.Global().Set("start", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go sched.Run()
		return nil
	}))
	wg.Wait()
}
