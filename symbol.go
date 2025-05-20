package gowasmtk

type wasmSymbolTable struct {
	functionTypes []WasmSectionFunctionType
}

func NewSymbolTable() *wasmSymbolTable {
	return &wasmSymbolTable{
		functionTypes: []WasmSectionFunctionType{},
	}
}
