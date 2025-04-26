package instructions

type WasmInstruction = byte

const (
	ConstI32 WasmInstruction = 0x41
	ConstI64 WasmInstruction = 0x42
	ConstF32 WasmInstruction = 0x43
	ConstF64 WasmInstruction = 0x44
	End      WasmInstruction = 0x0B
	AddI32   WasmInstruction = 0x6A
	SubI32   WasmInstruction = 0x6B
	MulI32   WasmInstruction = 0x6C
	DivI32   WasmInstruction = 0x6D
	GetLocal WasmInstruction = 0x20
	SetLocal WasmInstruction = 0x21
	TeeLocal WasmInstruction = 0x22
)
