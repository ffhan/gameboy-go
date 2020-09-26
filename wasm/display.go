package wasm

import (
	"syscall/js"
)

type wasmDisplay struct {
	drawing bool
}

func NewWasmDisplay() *wasmDisplay {
	return &wasmDisplay{}
}

func (w *wasmDisplay) Draw(buffer []byte) {
	w.drawing = true
	js.CopyBytesToJS(js.Global().Get("buffer"), buffer)
	js.Global().Get("draw").Invoke()
}

func (w *wasmDisplay) IsDrawing() bool {
	defer func() {
		w.drawing = false
	}()
	return w.drawing
}
