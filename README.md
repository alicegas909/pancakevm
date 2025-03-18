# pancakevm
Extremely naive stack based virtual machine in go. Currently supports 64 bit integer and floating point operations.

If you want to mess around with it, as of now you can only write programs in the "TestProgram" array located in "debug_vm.go". A simple program that adds two integers and prints the result would look like:

```go
var TestProgram = []int{
	ICONST, 4,
	ICONST, 6,
	IADD,
	PRINT,
	HALT,
}
```
A slightly more complicated program that stores two floats into global memory, calls a function that multiplies the two floats, and prints the result would look like:
```go
var TestProgram = []int{
	FCONST, f(3.5),
	GSTORE, 0,
	FCONST, f(2.5),
	GSTORE, 1,
	CALL, 13, 0, // call fuction at address 13 with zero arguments
	PRINT, //function returns to here
	HALT,
	GLOAD, 0, //address 13 is here
	GLOAD, 1,
	FMUL,
	RET,
}

```
To run the virtual machine, simply build "main.go" and run the executable. You can additionally use the command line argument "-trace" view the disassembly.
