package instructions

type WasmInstruction = byte

const (
	ConstI32                    WasmInstruction = 0x41 // (1)
	ConstI64                    WasmInstruction = 0x42 // (1)
	ConstF32                    WasmInstruction = 0x43 // (0.0)
	ConstF64                    WasmInstruction = 0x44 // (0.0)
	End                         WasmInstruction = 0x0B
	AddI32                      WasmInstruction = 0x6A // (+)
	SubI32                      WasmInstruction = 0x6B // (-)
	MulI32                      WasmInstruction = 0x6C // (*)
	DivI32                      WasmInstruction = 0x6D // (/)
	GetLocal                    WasmInstruction = 0x20
	SetLocal                    WasmInstruction = 0x21
	TeeLocal                    WasmInstruction = 0x22
	CallFunc                    WasmInstruction = 0x10
	If                          WasmInstruction = 0x04
	Else                        WasmInstruction = 0x05
	EqualI32                    WasmInstruction = 0x46 // (1 == 1)
	NotEqualI32                 WasmInstruction = 0x47 // (1 != 1)
	LessThanSignedI32           WasmInstruction = 0x48 // (-1 < -1)
	LessThanUnsignedI32         WasmInstruction = 0x49 // (1 < 1)
	GreaterThanSignedI32        WasmInstruction = 0x4A // (-1 > -1)
	GreaterThanUnsignedI32      WasmInstruction = 0x4B // (1 > 1)
	LessThanEqualSignedI32      WasmInstruction = 0x4C // (-1 <= -1)
	LessThanEqualUnsignedI32    WasmInstruction = 0x4D // (1 <= 1)
	GreaterThanEqualSignedI32   WasmInstruction = 0x4E // (-1 >= -1)
	GreaterThanEqualUnsignedI32 WasmInstruction = 0x4F // (1 >= 1)
	EqzI32                      WasmInstruction = 0x45 // (a == 0)
	AndI32                      WasmInstruction = 0x71 // (a & b)
	OrI32                       WasmInstruction = 0x72 // (a | b)
	Block                       WasmInstruction = 0x02
	Loop                        WasmInstruction = 0x03
	Br                          WasmInstruction = 0x0C
	BrIf                        WasmInstruction = 0x0D
)
