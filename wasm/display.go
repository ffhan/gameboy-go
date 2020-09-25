package wasm

import (
	"syscall/js"
)

type wasmDisplay struct {
	buffer []byte
}

func NewWasmDisplay() *wasmDisplay {
	buffer := make([]byte, 160*144*4)
	for i := range buffer {
		if i%4 == 3 {
			buffer[i] = 255
		}
	}
	return &wasmDisplay{buffer: buffer}
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
	for i, pixel := range buffer {
		r, g, b := w.mapColor(pixel)
		w.buffer[i*4] = r
		w.buffer[i*4+1] = g
		w.buffer[i*4+2] = b
	}
	js.CopyBytesToJS(js.Global().Get("document").Get("buffer"), w.buffer)
	js.Global().Get("draw").Invoke()
}
