package main

import (
	"fmt"
	go_gb "go-gb"
	"go-gb/cpu"
	"go-gb/memory"
	"go-gb/ppu"
	"os"
)

func main() {
	mmu := memory.NewMMU()
	file, err := os.Open("roms/Tetris (World) (Rev A).gb")
	if err != nil {
		panic(err)
	}

	game, err := go_gb.LoadGame(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(game)

	mmu.Init(game.Rom, go_gb.GB)

	lcd := go_gb.NewNopDisplay()

	//mmuD := memory.NewDebugger(mmu, os.Stdout)
	ppu := ppu.NewPpu(mmu, mmu.VRAM(), mmu.OAM(), lcd)
	c := cpu.NewDebugger(cpu.NewCpu(mmu, ppu), os.Stdout)

	//sysD := debugger.NewSystemDebugger(c, mmuD)
	c.Debug(false)

	for {
		c.Step()
		c.PC()
		if ppu.IsVBlank() {
			print()
		}
		if c.PC() == 0x8c {
			//sysD.Debug(true)
			//lcd.Debug(true)
			print()
		}
	}
}
