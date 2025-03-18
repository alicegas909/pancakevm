package virtualmachine

import (
	"fmt"
	"math"
)

type Value struct {
	IsFloat bool
	Int     int
	Float   float64
}

type VirtualMachine struct {
	Code  []int
	GMem  []Value
	Stack []Value

	Ip    int
	Sp    int
	Fp    int
	Asm   bool
	Trace bool
}

func (vm VirtualMachine) New(code []int, entry int, datasize int, asm bool, trace bool) VirtualMachine {
	vm.Code = code
	vm.GMem = make([]Value, datasize)
	vm.Stack = make([]Value, 100)
	vm.Ip = entry
	vm.Sp = -1
	vm.Asm = asm
	vm.Trace = trace

	return vm
}

func (vm *VirtualMachine) Processor() {
	TRUE := Value{IsFloat: false, Int: 1}
	FALSE := Value{IsFloat: false, Int: 0}
	for vm.Ip < len(vm.Code) {
		opcode := vm.Code[vm.Ip]

		if vm.Asm {
			vm.Disassemble()
		}

		if vm.Trace {
			vm.StackTrace()
		}
		switch opcode {
		// Integer Ops
		case ICONST:
			vm.Sp++
			vm.Ip++
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: vm.Code[vm.Ip]}
			vm.Ip++

		case IADD:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: a.Int + b.Int}
			vm.Ip++

		case ISUB:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: a.Int - b.Int}
			vm.Ip++

		case IMUL:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: a.Int * b.Int}
			vm.Ip++
		case ILT:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1.Int < v2.Int {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++

		case IGT:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1.Int > v2.Int {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++

		case IEQ:
			v1 := vm.Stack[vm.Sp]
			vm.Sp--
			v2 := vm.Stack[vm.Sp]
			if v1.Int == v2.Int {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++
		// Float Ops
		case FCONST:
			vm.Sp++
			vm.Ip++
			f := uint64(vm.Code[vm.Ip])
			fval := math.Float64frombits(f)
			vm.Stack[vm.Sp] = Value{IsFloat: true, Float: fval}
			vm.Ip++
		case FADD:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = Value{IsFloat: true, Float: a.Float + b.Float}
			vm.Ip++
		case FSUB:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = Value{IsFloat: true, Float: a.Float - b.Float}
			vm.Ip++
		case FMUL:
			b := vm.Stack[vm.Sp]
			vm.Sp--
			a := vm.Stack[vm.Sp]
			vm.Stack[vm.Sp] = Value{IsFloat: true, Float: a.Float * b.Float}
			vm.Ip++
		case FLT:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1.Float < v2.Float {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++
		case FGT:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1.Float > v2.Float {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++
		case FEQ:
			v2 := vm.Stack[vm.Sp]
			vm.Sp--
			v1 := vm.Stack[vm.Sp]
			if v1.Float == v2.Float {
				vm.Stack[vm.Sp] = TRUE
			} else {
				vm.Stack[vm.Sp] = FALSE
			}
			vm.Ip++
		case PRINT:
			v := vm.Stack[vm.Sp]
			vm.Sp--
			if v.IsFloat {
				fmt.Println(v.Float)
			} else {
				fmt.Println(v.Int)
			}
			vm.Ip++

		case CALL:
			vm.Ip++
			addr := vm.Code[vm.Ip]
			vm.Ip++
			numargs := vm.Code[vm.Ip]
			vm.Sp++
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: vm.Ip + 1}
			vm.Sp++
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: vm.Fp}
			vm.Sp++
			vm.Stack[vm.Sp] = Value{IsFloat: false, Int: numargs}
			vm.Fp = vm.Sp
			vm.Ip = addr

		case RET:

			rval := vm.Stack[vm.Sp]
			vm.Sp--
			vm.Sp = vm.Fp
			nargs := vm.Stack[vm.Sp].Int
			vm.Sp--
			oldFp := vm.Stack[vm.Sp].Int
			vm.Sp--
			retAddr := vm.Stack[vm.Sp].Int
			vm.Sp--
			vm.Sp -= nargs
			vm.Sp++
			vm.Stack[vm.Sp] = rval
			vm.Fp = oldFp
			vm.Ip = retAddr

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

// Disassemble prints the current instruction and the state of the stack
func (vm *VirtualMachine) Disassemble() {
	inst := vm.Code[vm.Ip]
	opname := "UNKNOWN"
	if name, ok := opnames[inst]; ok {
		opname = name
	}
	fmt.Printf("%04d: %-10s", vm.Ip, opname)

	// Print operand if instruction uses it
	switch inst {
	case ICONST, FCONST, GLOAD, GSTORE, BR, BRF, BRT, LOAD, STORE, FLOAD, FSTORE, CALL:
		if vm.Ip+1 < len(vm.Code) {
			if inst == FCONST {
				bits := uint64(vm.Code[vm.Ip+1])
				fmt.Printf(" %-8g", math.Float64frombits(bits))
			} else {
				fmt.Printf(" %-8d", vm.Code[vm.Ip+1])
			}
		} else {
			fmt.Printf(" %-8s", "INVALID")
		}
	default:
		fmt.Printf(" %-8s", "")
	}

	// Print stack state
	fmt.Printf(" [")
	for i := 0; i <= vm.Sp; i++ {
		if i > 0 {
			fmt.Printf(" ")
		}
		if vm.Stack[i].IsFloat {
			fmt.Printf("%g", vm.Stack[i].Float)
		} else {
			fmt.Printf("%d", vm.Stack[i].Int)
		}
	}
	fmt.Printf("]\n")
}

// StackTrace prints the entire program and marks the current instruction
func (vm *VirtualMachine) StackTrace() {
	i := 0
	for i < len(vm.Code) {
		// Store original position for comparison with IP
		originalPos := i

		fmt.Printf("%-6d ", i)
		inst := vm.Code[i]
		opname := "UNKNOWN"
		if name, ok := opnames[inst]; ok {
			opname = name
		}

		switch inst {
		case ICONST, FCONST, GLOAD, GSTORE, BR, BRF, BRT, LOAD, STORE, FLOAD, FSTORE, CALL:
			if i+1 < len(vm.Code) {
				if inst == FCONST {
					bits := uint64(vm.Code[i+1])
					fmt.Printf("%-12s %-8g", opname, math.Float64frombits(bits))
				} else {
					fmt.Printf("%-12s %-8d", opname, vm.Code[i+1])
				}
			} else {
				fmt.Printf("%-12s %-8s", opname, "INVALID")
			}
			i += 2 // Increment by 2 to account for operand
		default:
			fmt.Printf("%-12s %-8s", opname, "")
			i++ // Increment by 1 for instructions without operands
		}

		// Only print stack if this is the current instruction
		if originalPos == vm.Ip {
			fmt.Printf("   [")
			for j := 0; j <= vm.Sp; j++ {
				if j > 0 {
					fmt.Printf(" ")
				}
				if vm.Stack[j].IsFloat {
					fmt.Printf("%g", vm.Stack[j].Float)
				} else {
					fmt.Printf("%d", vm.Stack[j].Int)
				}
			}
			fmt.Printf("]")
			fmt.Printf(" <- current")
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------")
}
