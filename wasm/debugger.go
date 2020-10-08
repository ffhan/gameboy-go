package wasm

import (
	"bytes"
	"fmt"
	go_gb "go-gb"
	"go-gb/memory"
	"syscall/js"
)

type debugger struct {
	oam  go_gb.Memory
	vram go_gb.DumpableMemory

	stopped bool
	steps   int

	waitChan chan bool
}

func NewDebugger(oam go_gb.Memory, vram go_gb.DumpableMemory, joyPad Joypad) *debugger {
	d := &debugger{oam: oam, vram: vram, waitChan: make(chan bool)}
	joyPad.Subscribe(func(pressed bool) {
		if pressed {
			d.stopped = true
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
	joyPad.Subscribe(func(pressed bool) {
		if pressed {
			var buf bytes.Buffer
			memory.DumpOam(oam, vram, &buf)
			arr := js.Global().Get("Uint8Array").New(buf.Len())
			js.Global().Set("oam", arr)
			js.CopyBytesToJS(js.Global().Get("oam"), buf.Bytes())
		}
	}, Oam)
	joyPad.Subscribe(func(pressed bool) {
		if pressed {
			var buf bytes.Buffer
			vram.Dump(&buf)
			arr := js.Global().Get("Uint8Array").New(buf.Len())
			js.Global().Set("vram", arr)
			js.CopyBytesToJS(js.Global().Get("vram"), buf.Bytes())
		}
	}, Vram)
	return d
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
