package gowasmtk

import (
	"bytes"
	"os"

	"github.com/Orphoros/gowasmtk/instructions"
	"github.com/Orphoros/gowasmtk/types"
)

type WasmExportable interface {
	GetIndex() int
}

type WasmFunctionBuilder struct {
	paramTypes   []types.WasmType
	resultTypes  []types.WasmType
	code         []byte
	locals       [][]byte
	instructions []byte
	symbolTable  *wasmSymbolTable
	codeIndex    int
}

type WasmFunctionModule struct {
	sectionCode []byte
	typeIndex   int
	codeIndex   int
	funcType    wasmSectionFunctionType
}

func (m *WasmFunctionModule) GetIndex() int {
	return m.codeIndex
}

func NewWasmFunctionBuilder(symbolTable *wasmSymbolTable) *WasmFunctionBuilder {
	// reserve a slot for the function in the symbol table
	symbolTable.functions = append(symbolTable.functions, WasmFunctionModule{})

	return &WasmFunctionBuilder{
		paramTypes:   []types.WasmType{},
		resultTypes:  []types.WasmType{},
		code:         []byte{},
		locals:       [][]byte{},
		instructions: []byte{},
		symbolTable:  symbolTable,
		codeIndex:    len(symbolTable.functions) - 1,
	}
}

func (b *WasmFunctionBuilder) AddParam(paramType types.WasmType) *WasmFunctionBuilder {
	b.paramTypes = append(b.paramTypes, paramType)
	return b
}

func (b *WasmFunctionBuilder) AddReturn(resultType types.WasmType) *WasmFunctionBuilder {
	b.resultTypes = append(b.resultTypes, resultType)
	return b
}

func (b *WasmFunctionBuilder) AddLocal(n uint32, localType types.WasmType) *WasmFunctionBuilder {
	b.locals = append(b.locals, locals(n, localType))
	return b
}

func (b *WasmFunctionBuilder) AddInstrConstI32(n int32) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.ConstI32)
	b.instructions = append(b.instructions, leb128EncodeI(int64(n))...)
	return b
}

func (boolean *WasmFunctionBuilder) AddInstrConstI64(n int64) *WasmFunctionBuilder {
	boolean.instructions = append(boolean.instructions, instructions.ConstI64)
	boolean.instructions = append(boolean.instructions, leb128EncodeI(n)...)
	return boolean
}

func (b *WasmFunctionBuilder) AddInstrSetLocal(idx uint64) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.SetLocal)
	b.instructions = append(b.instructions, leb128EncodeU(idx)...)
	return b
}

func (b *WasmFunctionBuilder) AddInstrGetLocal(idx uint64) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.GetLocal)
	b.instructions = append(b.instructions, leb128EncodeU(idx)...)
	return b
}

func (b *WasmFunctionBuilder) AddInstrLocalTee(idx uint64) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.TeeLocal)
	b.instructions = append(b.instructions, leb128EncodeU(idx)...)
	return b
}

func (b *WasmFunctionBuilder) AddInstrAddI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.AddI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrSubI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.SubI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrMulI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.MulI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrDivI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.DivI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrIf(returnType types.PrimitiveType) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.If)
	b.instructions = append(b.instructions, returnType)
	return b
}

func (b *WasmFunctionBuilder) AddInstrElse() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.Else)
	return b
}

func (b *WasmFunctionBuilder) AddInstrEqI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.EqualI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrNotEqI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.NotEqualI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrLessThanI32S() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.LessThanSignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrLessThanI32U() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.LessThanUnsignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrGreaterThanI32S() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.GreaterThanSignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrGreaterThanI32U() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.GreaterThanUnsignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrLessThanEqI32S() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.LessThanEqualSignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrLessThanEqI32U() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.LessThanEqualUnsignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrGreaterThanEqI32S() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.GreaterThanEqualSignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrGreaterThanEqI32U() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.GreaterThanEqualUnsignedI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrEqzI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.EqzI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrAndI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.AndI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrOrI32() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.OrI32)
	return b
}

func (b *WasmFunctionBuilder) AddInstrCall(f *WasmFunctionModule) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.CallFunc)
	b.instructions = append(b.instructions, leb128EncodeU(uint64(f.GetIndex()))...)
	return b
}

func (b *WasmFunctionBuilder) AddInstrCallSelf() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.CallFunc)

	index := len(b.symbolTable.functions)

	b.instructions = append(b.instructions, leb128EncodeU(uint64(index))...)

	return b
}

func (b *WasmFunctionBuilder) AddInstrLoop(returnType types.PrimitiveType) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.Loop)
	b.instructions = append(b.instructions, returnType)
	return b
}

func (b *WasmFunctionBuilder) AddInstrBr(idx uint64) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.Br)
	b.instructions = append(b.instructions, leb128EncodeU(idx)...)
	return b
}

func (b *WasmFunctionBuilder) AddInstrBrIf(idx uint64) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.BrIf)
	b.instructions = append(b.instructions, leb128EncodeU(idx)...)
	return b
}

func (b *WasmFunctionBuilder) AddInstrBlock(returnType types.PrimitiveType) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.Block)
	b.instructions = append(b.instructions, returnType)
	return b
}

func (b *WasmFunctionBuilder) AddInstrEnd() *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instructions.End)
	return b
}

func (b *WasmFunctionBuilder) Build() WasmFunctionModule {
	funcType := funcType(b.paramTypes, b.resultTypes)
	typeIndex := -1
	for i, f := range b.symbolTable.functionTypes {
		if bytes.Equal(f, funcType) {
			typeIndex = i
			break
		}
	}
	if typeIndex == -1 {
		typeIndex = len(b.symbolTable.functionTypes)
		b.symbolTable.functionTypes = append(b.symbolTable.functionTypes, funcType)
	}

	m := WasmFunctionModule{
		sectionCode: b.buildFunctionCode(),
		typeIndex:   typeIndex,
		funcType:    funcType,
	}

	// Update the function in the symbol table
	b.symbolTable.functions[m.codeIndex] = m
	// Update the code index in the function
	m.codeIndex = len(b.symbolTable.functions) - 1

	return m
}

func (b *WasmFunctionBuilder) buildFunctionCode() []byte {
	return code(function(b.locals, b.instructions))
}

type WasmModuleBuilder struct {
	metaLanuages         []wasmMetadata
	metaTools            []wasmMetadata
	metaSdks             []wasmMetadata
	sectionFunctionTypes []wasmSectionFunctionType
	sectionFunction      []int // FIXME: Should be uint32
	sectionExports       []wasmSectionExportedModule
	sectionCode          [][]byte
	exportNames          []string
}

func NewWasmModuleBuilder(wasmSymbolTable *wasmSymbolTable) *WasmModuleBuilder {
	return &WasmModuleBuilder{
		metaLanuages:         []wasmMetadata{},
		metaTools:            []wasmMetadata{},
		metaSdks:             []wasmMetadata{},
		sectionFunctionTypes: wasmSymbolTable.functionTypes,
		sectionExports:       []wasmSectionExportedModule{},
		sectionCode:          [][]byte{},
		sectionFunction:      []int{},
		exportNames:          []string{},
	}
}

// Adds a source programming language to the module as metadata. This is an optional field. Examples of languages
// include "C" or "Rust". Multiple languages can be added to the module.
func (b *WasmModuleBuilder) AddMetaLanguage(name, version string) *WasmModuleBuilder {
	b.metaLanuages = append(b.metaLanuages, wasmMetadata{
		Name:    name,
		Version: version,
	})

	return b
}

// Adds an overall pipeline tool that produces and optimizes a given wasm module as metadata to the module.
// This is an optional field. Examples of tools include "LLVM" or "rustc". Multiple tools can be added to the module.
func (b *WasmModuleBuilder) AddMetaTool(name, version string) *WasmModuleBuilder {
	b.metaTools = append(b.metaTools, wasmMetadata{
		Name:    name,
		Version: version,
	})

	return b
}

// Adds SDK information to the module as metadata. This is an optional field.
// An SDK is a higher-level tool that can be installed to produce the wasm module.
// Examples of SDKs include "Emscripten" or "Webpack". Multiple SDKs can be added to the module.
func (b *WasmModuleBuilder) AddMetaSdk(name, version string) *WasmModuleBuilder {
	b.metaSdks = append(b.metaSdks, wasmMetadata{
		Name:    name,
		Version: version,
	})

	return b
}

// Register a function in the module. The function must be built using the WasmFunctionBuilder.
func (b *WasmModuleBuilder) AddFunction(function *WasmFunctionModule) *WasmModuleBuilder {
	b.sectionFunction = append(b.sectionFunction, function.typeIndex)
	b.sectionCode = append(b.sectionCode, function.sectionCode)

	return b
}

// Save the WASM module to a ".wasm" file. May return an error if the file cannot be created or written to.
func (b *WasmModuleBuilder) BuildWasmFile(fileName string) error {
	if len(fileName) < 5 || fileName[len(fileName)-5:] != ".wasm" {
		fileName += ".wasm"
	}

	return os.WriteFile(fileName, b.Build(), 0644)
}

// Export an item (function) from the module. The item must implement the WasmExportable interface.
// The name must be unique. If the name already exists, it will not be added again. The type of the item
// must be one of the WasmExportType constants. The item will be exported with the given name and type.
func (b *WasmModuleBuilder) Export(name string, exportType types.WasmExportType, item WasmExportable) *WasmModuleBuilder {
	found := false

	for _, n := range b.exportNames {
		if n == name {
			found = true
			break
		}
	}
	if found {
		return b
	}

	b.sectionExports = append(b.sectionExports, export(name, wasmExportDescription{
		Type:  exportType,
		Index: item.GetIndex(),
	}))
	b.exportNames = append(b.exportNames, name)

	return b
}

// Build the WASM bytecode. Returns the WASM bytecode as a byte slice.
func (b *WasmModuleBuilder) Build() []byte {
	sections := []wasmSection{}

	sections = append(sections, sectionType(b.sectionFunctionTypes...))
	funcIndices := make([]uint64, len(b.sectionFunction))
	for i, idx := range b.sectionFunction {
		funcIndices[i] = uint64(idx)
	}
	sections = append(sections, sectionFunc(funcIndices...))

	if len(b.sectionExports) > 0 {
		sections = append(sections, sectionExport(b.sectionExports...))
	}

	sections = append(sections, sectionCode(b.sectionCode...))

	if len(b.metaLanuages) > 0 || len(b.metaTools) > 0 || len(b.metaSdks) > 0 {
		sections = append(sections, sectionProducers(b.metaLanuages, b.metaTools, b.metaSdks))
	}

	return module(sections...)
}
