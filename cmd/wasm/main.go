package main

import (
	"bytes"
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/scheduler"
	"go-gb/wasm"
	"io/ioutil"
	"sync"
	"syscall/js"
)

func run() (go_gb.Cpu, go_gb.Memory, go_gb.PPU, go_gb.Display) {
	mmu := memory.NewMMU()
	rom := make([]byte, 2*1<<20)
	n := js.CopyBytesToGo(rom, js.Global().Get("rom"))
	mmu.Init(rom[:n], go_gb.GB)

	game, err := go_gb.LoadGame(ioutil.NopCloser(bytes.NewBuffer(rom[:n])))
	if err != nil {
		panic(err)
	}
	js.Global().Set("title", game.Title)
	js.Global().Set("cartridgeType", game.CartridgeType.String())
	js.Global().Set("sgb", game.SGBFlag.String())
	js.Global().Set("cgb", game.CGBFlag.String())
	js.Global().Set("romSize", game.RomSize.String())
	js.Global().Set("ramSize", game.RamSize.String())
	js.Global().Set("nonJapanese", game.NonJapanese)

	fmt.Println("initialized mmu")

	lcd := wasm.NewWasmDisplay()

	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), mmu.IO(), lcd)
	c := cpu.NewCpu(mmu, ppu)

	return c, mmu, ppu, lcd
}

type Runner interface {
	Run()
	AddStopper(addr uint16)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var cpu go_gb.Cpu
	var ppu go_gb.PPU
	var mmu go_gb.Memory
	var lcd go_gb.Display

	var sched Runner

	js.Global().Set("run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cpu, mmu, ppu, lcd = run()
		sched = scheduler.NewScheduler(cpu, ppu, lcd)
		sched.AddStopper(0x100)
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
