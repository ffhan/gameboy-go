package main

import (
	"bytes"
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"go-gb/scheduler"
	"go-gb/serial"
	"go-gb/timer"
	"go-gb/wasm"
	"io/ioutil"
	"sync"
	"syscall/js"
)

func run() (go_gb.Cpu, go_gb.MemoryBus, go_gb.PPU, go_gb.Display, wasm.Joypad) {
	mmu := memory.NewMMU()
	rom := make([]byte, 2*1<<20)
	n := js.CopyBytesToGo(rom, js.Global().Get("rom"))

	joypad := wasm.NewJoypad() // todo: fix this relationship

	mmu.Init(rom[:n], go_gb.GB, joypad)
	joypad.Init(mmu.IO())

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

	divTimer := timer.NewDivTimer(mmu.IO())
	timer := timer.NewTimer(mmu.IO())

	serialPort := serial.NewSerial(nil, nil, nil, mmu.IO())

	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), mmu.IO(), lcd)
	c := cpu.NewCpu(mmu, ppu, timer, divTimer, serialPort)

	return c, mmu, ppu, lcd, joypad
}

type Runner interface {
	Run()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var cpu go_gb.Cpu
	var ppu go_gb.PPU
	var mmu go_gb.MemoryBus
	var lcd go_gb.Display
	var joypad wasm.Joypad

	var sched Runner

	js.Global().Set("run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cpu, mmu, ppu, lcd, joypad = run()
		s := scheduler.NewScheduler(cpu, ppu, lcd)
		s.Controller = wasm.NewDebugger(cpu, ppu, mmu, mmu.IO(), mmu.OAM(), mmu.VRAM(), joypad)
		sched = s
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
