package gowasmtk

type wasmSymbolTable struct {
	functionTypes []WasmSectionFunctionType
	functions     []WasmFunctionModule
}

func NewSymbolTable() *wasmSymbolTable {
	return &wasmSymbolTable{
		functionTypes: []WasmSectionFunctionType{},
		functions:     []WasmFunctionModule{},
	}
}
