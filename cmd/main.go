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
	"time"
)

func run() {
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

	sleepTime := time.Second

	debug := cpu.NewDebugger(c, os.Stdout)

	for {
		debug.Step()
		time.Sleep(sleepTime)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	js.Global().Set("run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go run()
		return nil
	}))
	wg.Wait()
}
