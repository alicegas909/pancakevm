package main

import (
	m "vm/virtualmachine"
)

func main() {

	var vm m.VirtualMachine
	vm = vm.New(m.TestProgram, 22, 5, false, true, false)
	vm.Processor()
}
