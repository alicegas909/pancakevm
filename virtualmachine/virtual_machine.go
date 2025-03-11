package virtualmachine

import (
	"fmt"
)

type VirtualMachine struct {
	Code  []int
	GMem  []int
	Stack []int

	Ip    int
	Sp    int
	Fp    int
	Asm   bool
	Trace bool
	Dbg   bool
}

func (vm VirtualMachine) New(code []int, entry int, datasize int, asm bool, trace bool, dbg bool) VirtualMachine {
	vm.Code = code
	vm.GMem = make([]int, datasize)
	vm.Stack = make([]int, 100)
	vm.Ip = entry
	vm.Sp = -1
	vm.Asm = asm
	vm.Trace = trace
	vm.Dbg = dbg
	return vm
}

func (vm *VirtualMachine) Processor() {
	TRUE := 1
	FALSE := 0
	for vm.Ip < len(vm.Code) {
		opcode := vm.Code[vm.Ip]

		if vm.Dbg {
			op := opnames[opcode]
			fmt.Print(vm.Ip, " ", op, " ")
		}

		if vm.Trace {
			vm.StackTrace()
		}
		switch opcode {
		case ICONST:
			vm.Sp++
			vm.Ip++
			vm.Stack[vm.Sp] = vm.Code[vm.Ip]
			vm.Ip++

		case PRINT:
			v := vm.Stack[vm.Sp]
			vm.Sp--
			fmt.Println(v)
			vm.Ip++

		case CALL:
			vm.Ip++
			addr := vm.Code[vm.Ip]
			vm.Ip++
			numargs := vm.Code[vm.Ip]
			vm.Sp++
			vm.Stack[vm.Sp] = numargs
			vm.Sp++
			vm.Stack[vm.Sp] = vm.Fp
			vm.Sp++
			vm.Stack[vm.Sp] = vm.Ip
			vm.Fp = vm.Sp
			vm.Ip = addr

		case RET:
			rval := vm.Stack[vm.Sp]
			vm.Sp--
			vm.Sp = vm.Fp
			vm.Ip = vm.Stack[vm.Sp]
			vm.Sp--
			nargs := vm.Stack[vm.Sp]
			vm.Sp--
			vm.Sp -= nargs
			vm.Sp++
			vm.Stack[vm.Sp] = rval
			vm.Ip++

		case GLOAD:
			vm.Ip++
			addr := vm.Code[vm.Ip]
			v := vm.GMem[addr]
			vm.Sp++
			vm.Stack[vm.Sp] = v
			vm.Ip++

		case LOAD:
			vm.Ip++
			offset := vm.Code[vm.Ip]
			vm.Sp++
			addr := vm.Fp + offset
			vm.Stack[vm.Sp] = vm.Stack[addr]
			vm.Ip++

		case GSTORE:
			vm.Ip++
			addr := vm.Code[vm.Ip]
			v := vm.Stack[vm.Sp]
			vm.Sp--
			vm.GMem[addr] = v
			vm.Ip++

		case STORE:
			vm.Ip++
			offset := vm.Code[vm.Ip]
			addr := vm.Fp + offset
			vm.Stack[addr] = vm.Stack[vm.Sp]
			vm.Sp--
			vm.Ip++

		case IADD:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = a + b
			vm.Ip++

		case ISUB:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = a - b
			vm.Ip++

		case IMUL:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = a * b
			vm.Ip++
			if vm.Dbg {
				print("mul\n")
			}
		case BR:
			vm.Ip++
			vm.Ip = vm.Code[vm.Ip]

			vm.Ip++
			addr := vm.Code[vm.Ip]
			v := vm.Stack[vm.Sp]
			vm.Sp--
			if v == TRUE {
				vm.Ip = addr
			} else {
				vm.Ip++
			}

		case BRF:
			vm.Ip++
			addr := vm.Code[vm.Ip]
			v := vm.Stack[vm.Sp]
			vm.Sp--
			if v == FALSE {
				vm.Ip = addr
			} else {
				vm.Ip++
			}

		case ILT:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1 < v2 {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++

		case IGT:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1 > v2 {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++

		case IEQ:
			v1 := vm.Stack[vm.Sp]
			vm.Sp--
			v2 := vm.Stack[vm.Sp]
			if v1 == v2 {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++
		case POP:
			vm.Sp--
			vm.Ip++
		case HALT:
			return
		default:
			break
		}
	}
}

func (vm *VirtualMachine) Disassemble() {
	inst := vm.Code[vm.Ip]
	fmt.Printf("%04d: %-10s", vm.Ip, opnames[inst])

	// Print operand if instruction uses it
	switch inst {
	case ICONST, GLOAD, GSTORE, BR, BRF, BRT, LOAD, STORE, CALL:
		fmt.Printf(" %-8d", vm.Code[vm.Ip+1])
	default:
		fmt.Printf(" %-8s", "")
	}

	// Print stack state
	fmt.Printf(" [")
	for i := 0; i <= vm.Sp; i++ {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%d", vm.Stack[i])
	}
	fmt.Printf("]\n")
}

func (vm *VirtualMachine) StackTrace() {
	i := 0
	for i < len(vm.Code) {
		fmt.Printf("%-6d ", i)
		inst := vm.Code[i]
		switch inst {
		case ICONST, GLOAD, GSTORE, BR, BRF, BRT, LOAD, STORE, CALL:
			fmt.Printf("%-12s %-8d", opnames[inst], vm.Code[i+1])
			i++
		default:
			fmt.Printf("%-10s %-8s", opnames[inst], "")
		}
		if i == vm.Ip {
			fmt.Printf("   [")
			for j := 0; j <= vm.Sp; j++ {
				if j > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%d", vm.Stack[j])
			}
			fmt.Printf("]")
			fmt.Printf(" <- current")
		}
		fmt.Println()
		i++
	}
	fmt.Println("-------------------------------------------")
}
