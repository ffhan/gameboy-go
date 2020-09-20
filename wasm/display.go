package wasm

import (
	"fmt"
	"syscall/js"
)

type wasmDisplay struct {
}

func NewWasmDisplay() *wasmDisplay {
	return &wasmDisplay{}
}

func (w *wasmDisplay) Draw(scanLine int, bufferLine []byte) {
	fmt.Println("buffer", bufferLine)
	js.CopyBytesToJS(js.Global().Get("document").Get("bufferLine"),
		bufferLine)
	js.Global().Get("document").Set("scanLine", scanLine)
	js.Global().Get("draw").Invoke()
}
