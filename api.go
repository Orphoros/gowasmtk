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
}

type WasmFunctionModule struct {
	sectionCode []byte
	typeIndex   int
	codeIndex   int
	funcType    WasmSectionFunctionType
}

func (m *WasmFunctionModule) GetIndex() int {
	return m.codeIndex
}

func NewWasmFunctionBuilder(symbolTable *wasmSymbolTable) *WasmFunctionBuilder {
	return &WasmFunctionBuilder{
		paramTypes:   []types.WasmType{},
		resultTypes:  []types.WasmType{},
		code:         []byte{},
		locals:       [][]byte{},
		instructions: []byte{},
		symbolTable:  symbolTable,
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

func (b *WasmFunctionBuilder) AddInstruction(instruction instructions.WasmInstruction) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, instruction)
	return b
}

func (b *WasmFunctionBuilder) AddI32(n int32) *WasmFunctionBuilder {
	b.instructions = append(b.instructions, leb128EncodeI(int64(n))...)
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

	b.symbolTable.functions = append(b.symbolTable.functions, m)
	m.codeIndex = len(b.symbolTable.functions) - 1

	return m
}

func (b *WasmFunctionBuilder) buildFunctionCode() []byte {
	return code(function(b.locals, b.instructions))
}

type WasmModuleBuilder struct {
	metaLanuages         []WasmMetadata
	metaTools            []WasmMetadata
	metaSdks             []WasmMetadata
	sectionFunctionTypes []WasmSectionFunctionType
	sectionFunction      []int // FIXME: Should be uint32
	sectionExports       []WasmSectionExportedModule
	sectionCode          [][]byte
	exportNames          []string
}

func NewWasmModuleBuilder(wasmSymbolTable *wasmSymbolTable) *WasmModuleBuilder {
	return &WasmModuleBuilder{
		metaLanuages:         []WasmMetadata{},
		metaTools:            []WasmMetadata{},
		metaSdks:             []WasmMetadata{},
		sectionFunctionTypes: wasmSymbolTable.functionTypes,
		sectionExports:       []WasmSectionExportedModule{},
		sectionCode:          [][]byte{},
		sectionFunction:      []int{},
		exportNames:          []string{},
	}
}

// Adds a source programming language to the module as metadata. This is an optional field. Examples of languages
// include "C" or "Rust". Multiple languages can be added to the module.
func (b *WasmModuleBuilder) AddMetaLanguage(name, version string) *WasmModuleBuilder {
	b.metaLanuages = append(b.metaLanuages, WasmMetadata{
		Name:    name,
		Version: version,
	})

	return b
}

// Adds an overall pipeline tool that produces and optimizes a given wasm module as metadata to the module.
// This is an optional field. Examples of tools include "LLVM" or "rustc". Multiple tools can be added to the module.
func (b *WasmModuleBuilder) AddMetaTool(name, version string) *WasmModuleBuilder {
	b.metaTools = append(b.metaTools, WasmMetadata{
		Name:    name,
		Version: version,
	})

	return b
}

// Adds SDK information to the module as metadata. This is an optional field.
// An SDK is a higher-level tool that can be installed to produce the wasm module.
// Examples of SDKs include "Emscripten" or "Webpack". Multiple SDKs can be added to the module.
func (b *WasmModuleBuilder) AddMetaSdk(name, version string) *WasmModuleBuilder {
	b.metaSdks = append(b.metaSdks, WasmMetadata{
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

	b.sectionExports = append(b.sectionExports, export(name, WasmExportDescription{
		Type:  exportType,
		Index: item.GetIndex(),
	}))
	b.exportNames = append(b.exportNames, name)

	return b
}

// Build the WASM bytecode. Returns the WASM bytecode as a byte slice.
func (b *WasmModuleBuilder) Build() []byte {
	sections := []WasmSection{}

	sections = append(sections, sectionType(b.sectionFunctionTypes...))
	funcIndices := make([]uint64, len(b.sectionFunction))
	for i, idx := range b.sectionFunction {
		funcIndices[i] = uint64(idx)
	}
	sections = append(sections, sectionFunc(funcIndices...))

	if len(b.sectionExports) > 0 {
		sections = append(sections, sectionExport(b.sectionExports...))
	}

	sections = append(sections, SectionCode(b.sectionCode...))

	if len(b.metaLanuages) > 0 || len(b.metaTools) > 0 || len(b.metaSdks) > 0 {
		sections = append(sections, sectionProducers(b.metaLanuages, b.metaTools, b.metaSdks))
	}

	return module(sections...)
}
