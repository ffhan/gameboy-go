package debugger

import go_gb "go-gb"

type Debugger interface {
	Debug(val bool)
}

type CpuDebugger interface {
	go_gb.Cpu
	Debugger
}

type PpuDebugger interface {
	go_gb.PPU
	Debugger
}

type MemoryDebugger interface {
	go_gb.Memory
	Debugger
}

type debugger struct {
	cpu     CpuDebugger
	mmu     MemoryDebugger
	debugOn bool
}

func NewSystemDebugger(cpu CpuDebugger, mmu MemoryDebugger) *debugger {
	return &debugger{cpu: cpu, mmu: mmu}
}

func (d *debugger) Debug(val bool) {
	d.debugOn = val
	d.cpu.Debug(val)
	d.mmu.Debug(val)
}
