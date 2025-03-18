package virtualmachine

import "math"

func f(fl float64) int {
	return int(math.Float64bits(fl))
}

var TestProgram = []int{
	FCONST, f(3.5),
	GSTORE, 0,
	FCONST, f(2.5),
	GSTORE, 1,
	CALL, 13, 0,
	PRINT,
	HALT,
	GLOAD, 0,
	GLOAD, 1,
	FMUL,
	RET,
}
