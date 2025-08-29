package gowasmtk

type wasmSymbolTable struct {
	functionTypes []wasmSectionFunctionType
	functions     []WasmFunctionModule
	imports       *[]WasmImportDeclaration
}

func NewSymbolTable(imports *[]WasmImportDeclaration) *wasmSymbolTable {
	return &wasmSymbolTable{
		functionTypes: []wasmSectionFunctionType{},
		functions:     []WasmFunctionModule{},
		imports:       imports,
	}
}
