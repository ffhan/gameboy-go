package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/wasm"
	"os"
	"sync"
	"syscall/js"
)

func run() go_gb.Cpu {
	const (
		CpuFrequency = 4_194_304
	)
	mmu := memory.NewMMU()
	rom := make([]byte, 2*1<<20)
	n := js.CopyBytesToGo(rom, js.Global().Get("document").Get("rom"))
	mmu.Init(rom[:n], go_gb.GB)

	fmt.Println("initialized mmu")

	lcd := wasm.NewWasmDisplay()

	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), lcd)
	c := cpu.NewCpu(mmu, ppu)

	return cpu.NewDebugger(c, os.Stdout)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var c go_gb.Cpu

	js.Global().Set("run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c = run()
		return nil
	}))
	js.Global().Set("step", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.Step()
		return nil
	}))
	wg.Wait()
}
