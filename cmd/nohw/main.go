package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"io/ioutil"
	"os"
)

func main() {
	mmu := memory.NewMMU()
	rom := make([]byte, 2*1<<20)

	rom, err := ioutil.ReadFile("roms/Tetris (World) (Rev A).gb")
	if err != nil {
		panic(err)
	}

	fmt.Println(go_gb.NewGame(rom))

	mmu.Init(rom, go_gb.GB)

	fmt.Println("initialized mmu")

	lcd := go_gb.NewNopDisplay()

	//mmuD := memory.NewDebugger(mmu, os.Stdout)
	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), lcd)
	c := cpu.NewDebugger(cpu.NewCpu(mmu, ppu), os.Stdout)

	//sysD := debugger.NewSystemDebugger(c, mmuD)
	c.Debug(true)

	for {
		c.Step()
		c.PC()
		if ppu.IsVBlank() {
			print()
		}
		if c.PC() == 0x64 {
			//sysD.Debug(true)
			//lcd.Debug(true)
			print()
		}
	}
}
