package wasm

import (
	go_gb "go-gb"
	"syscall/js"
)

type wasmDisplay struct {
	drawing          bool
	buffer, drawFunc js.Value
}

func NewWasmDisplay() *wasmDisplay {
	return &wasmDisplay{
		drawing:  false,
		buffer:   js.Global().Get("buffer"),
		drawFunc: js.Global().Get("draw"),
	}
}

func (w *wasmDisplay) Draw(buffer []byte) {
	w.drawing = true
	js.CopyBytesToJS(w.buffer, buffer)
	w.drawFunc.Invoke()
	go_gb.Events.Add("drawn to display")
}

func (w *wasmDisplay) IsDrawing() bool {
	defer func() {
		w.drawing = false
	}()
	return w.drawing
}
