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

func run() (debugger.CpuDebugger, debugger.MemoryDebugger, go_gb.PPU) {
	mmu := memory.NewMMU()
	rom := make([]byte, 2*1<<20)
	n := js.CopyBytesToGo(rom, js.Global().Get("document").Get("rom"))
	mmu.Init(rom[:n], go_gb.GB)

	fmt.Println("initialized mmu")

	lcd := wasm.NewWasmDisplay()

	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), lcd)
	mmuD := memory.NewDebugger(mmu, os.Stdout)
	c := cpu.NewCpu(mmuD, ppu)

	return cpu.NewDebugger(c, os.Stdout), mmuD, ppu
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var c debugger.CpuDebugger
	var p go_gb.PPU
	var m debugger.MemoryDebugger

	var systemDebugger debugger.Debugger

	js.Global().Set("run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c, m, p = run()
		systemDebugger = debugger.NewSystemDebugger(c, m)
		systemDebugger.Debug(false)
		return nil
	}))
	js.Global().Set("step", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			for {
				c.Step()
				if p.IsVBlank() || p.Mode() == 1 {
					return
				}
			}
		}()
		return nil
	}))
	js.Global().Set("start", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go scheduler.NewScheduler(c, p).Run()
		return nil
	}))
	wg.Wait()
}
