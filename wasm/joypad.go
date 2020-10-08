package wasm

import (
	"fmt"
	go_gb "go-gb"
	"syscall/js"
)

type Key byte

const (
	ButtonRight Key = iota
	ButtonLeft
	ButtonUp
	ButtonDown
	ButtonA
	ButtonB
	Select
	Start
)

type joypad struct {
	io            go_gb.Memory
	currentlyHeld map[Key]bool
}

func NewJoypad() *joypad {
	j := &joypad{currentlyHeld: make(map[Key]bool)}
	js.Global().Set("keyDown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := Key(args[0].Int())
		j.KeyDown(key)
		fmt.Printf("go key %d down\n", key)
		return nil
	}))
	js.Global().Set("keyUp", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := Key(args[0].Int())
		j.KeyUp(key)
		fmt.Printf("go key %d up\n", key)
		return nil
	}))
	return j
}

func (j *joypad) Init(io go_gb.Memory) {
	j.io = io
}

func (j *joypad) KeyDown(key Key) {
	j.currentlyHeld[key] = true
}

func (j *joypad) KeyUp(key Key) {
	j.currentlyHeld[key] = false
}

func (j *joypad) Read(pointer uint16) byte {
	if pointer != go_gb.JOYP {
		panic("invalid read from JOYP")
	}
	buttons := j.io.Read(go_gb.JOYP)&0x20 != 0
	if buttons {
		result := byte(0x20)
		for i := 0; i < 4; i++ {
			currentlyHeld, pressedBefore := j.currentlyHeld[Key(i)]
			if !pressedBefore || !currentlyHeld {
				result |= 1 << i
			}
		}
		return result
	}
	result := byte(0x10)
	for i := 0; i < 4; i++ {
		currentlyHeld, pressedBefore := j.currentlyHeld[Key(i+4)]
		if !pressedBefore || !currentlyHeld {
			result |= 1 << i
		}
	}
	return result
}
