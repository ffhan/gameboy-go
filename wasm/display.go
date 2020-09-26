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

func (w *wasmDisplay) mapColor(col byte) (r, g, b byte) {
	switch col {
	case 0:
		return 255, 255, 255
	case 1:
		return 0xCC, 0xCC, 0xCC
	case 2:
		return 0x77, 0x77, 0x77
	case 3:
		return 0, 0, 0
	}
	panic("invalid color")
}

func (w *wasmDisplay) Draw(buffer []byte) {
	w.drawing = true
	js.CopyBytesToJS(js.Global().Get("document").Get("buffer"), buffer)
	js.Global().Get("draw").Invoke()
}

func (w *wasmDisplay) IsDrawing() bool {
	defer func() {
		w.drawing = false
	}()
	return w.drawing
}
