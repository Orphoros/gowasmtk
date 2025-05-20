package gowasmtk

type wasmSymbolTable struct {
	functionTypes []wasmSectionFunctionType
	functions     []WasmFunctionModule
}

func NewSymbolTable() *wasmSymbolTable {
	return &wasmSymbolTable{
		functionTypes: []wasmSectionFunctionType{},
		functions:     []WasmFunctionModule{},
	}
}
