package virtualmachine

var TestProgram = []int{
	LOAD, -3, //0
	ICONST, 2, //2
	ILT,     //4
	BRF, 10, //5
	ICONST, 1, //7
	RET,      //9
	LOAD, -3, //10
	LOAD, -3, //12
	ICONST, 1, //14
	ISUB,       //16
	CALL, 0, 1, //17
	IMUL,      //20
	RET,       //21
	ICONST, 1, //22 <--- entry point
	CALL, 0, 1, //24
	PRINT, //27
	HALT,  //28

}
