package main

import (
	"flag"
	m "vm/virtualmachine"
)

func main() {
	trace := flag.Bool("trace", false, "Enable instruction tracing")
	asm := flag.Bool("asm", false, "Enable assembly output")
	flag.Parse()
	var vm m.VirtualMachine
	vm = vm.New(m.TestProgram, 0, 10, *asm, *trace)
	vm.Processor()
}
