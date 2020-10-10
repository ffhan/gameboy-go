package wasm

import (
	"bytes"
	"fmt"
	go_gb "go-gb"
	cpu2 "go-gb/cpu"
	"go-gb/memory"
	"syscall/js"
)

type debugger struct {
	cpu go_gb.Cpu
	ppu go_gb.PPU

	memory go_gb.Memory
	oam    go_gb.Memory
	vram   go_gb.Memory

	stopped bool
	steps   int

	waitChan chan bool
}

func NewDebugger(c go_gb.Cpu, p go_gb.PPU, mem, io, oam, vram go_gb.Memory, joyPad Joypad) *debugger {
	d := &debugger{cpu: c, ppu: p, memory: mem, oam: oam, vram: vram, waitChan: make(chan bool)}
	joyPad.Subscribe(func(pressed bool) {
		if pressed {
			d.stopped = true

			var buf bytes.Buffer
			dumpCpu(d.cpu, d.ppu, buf)
			buf.Reset()
			dumpOam(oam, vram, buf)
			buf.Reset()
			dumpVram(io, vram, buf)
		}
	}, Pause)
	joyPad.Subscribe(func(pressed bool) {
		if pressed {
			d.stopped = false
			d.waitChan <- true
		}
	}, Continue)
	joyPad.Subscribe(func(pressed bool) {
		if pressed {
			d.steps = 1
		}
	}, Step)
	js.Global().Set("memoryRequest", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		start := uint16(args[0].Int())
		end := uint16(args[1].Int())

		dumpMemory(mem, start, end)
		return nil
	}))
	return d
}

func dumpCpu(cpu go_gb.Cpu, ppu go_gb.PPU, buf bytes.Buffer) {
	cpu2.DumpCpu(&buf, cpu, ppu)
	arr := js.Global().Get("Uint8Array").New(buf.Len())
	js.Global().Set("cpu", arr)
	js.CopyBytesToJS(js.Global().Get("cpu"), buf.Bytes())
}

func dumpMemory(mem go_gb.Memory, start uint16, end uint16) {
	var buf bytes.Buffer
	memory.DumpMemory(&buf, mem, start, end)
	arr := js.Global().Get("Uint8Array").New(buf.Len())
	js.Global().Set("mem", arr)
	js.CopyBytesToJS(js.Global().Get("mem"), buf.Bytes())
}

func dumpVram(io go_gb.Memory, vram go_gb.Memory, buf bytes.Buffer) {
	memory.DumpVram(io, vram, &buf)
	arr := js.Global().Get("Uint8Array").New(buf.Len())
	js.Global().Set("vram", arr)
	js.CopyBytesToJS(js.Global().Get("vram"), buf.Bytes())
}

func dumpOam(oam go_gb.Memory, vram go_gb.Memory, buf bytes.Buffer) {
	memory.DumpOam(oam, vram, &buf)
	arr := js.Global().Get("Uint8Array").New(buf.Len())
	js.Global().Set("oam", arr)
	js.CopyBytesToJS(js.Global().Get("oam"), buf.Bytes())
}

func (d *debugger) Wait() bool {
	if d.stopped {
		if d.steps > 0 {
			d.steps -= 1
			return false
		}
		fmt.Println("debugger stopped execution")
		defer fmt.Println("debugger continued execution")
		<-d.waitChan
		return true
	}
	return false
}
