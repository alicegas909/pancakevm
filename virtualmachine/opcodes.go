package virtualmachine

type Opcode = int

const (
	IADD Opcode = iota + 1
	ISUB
	IMUL
	ILT
	IGT
	IEQ
	BR
	BRT
	BRF
	ICONST
	LOAD
	GLOAD
	STORE
	GSTORE
	PRINT
	POP
	CALL
	RET
	HALT
)

var opnames = map[int]string{
	IADD:   "IADD",
	ISUB:   "ISUB",
	IMUL:   "IMUL",
	ILT:    "ILT",
	IGT:    "IGT",
	IEQ:    "IEQ",
	BR:     "BR",
	BRT:    "BRT",
	BRF:    "BRF",
	ICONST: "ICONST",
	LOAD:   "LOAD",
	GLOAD:  "GLOAD",
	STORE:  "STORE",
	GSTORE: "GSTORE",
	PRINT:  "PRINT",
	POP:    "POP",
	CALL:   "CALL",
	RET:    "RET",
	HALT:   "HALT",
}
