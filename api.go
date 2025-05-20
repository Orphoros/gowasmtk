package gowasmtk

import (
	"bytes"
	"os"

	"github.com/Orphoros/gowasmtk/instructions"
	"github.com/Orphoros/gowasmtk/types"
)

type WasmFunctionBuilder struct {
	paramTypes   []types.WasmType
	resultTypes  []types.WasmType
	name         *string
	code         []byte
	locals       [][]byte
	instructions []byte
	symbolTable  *wasmSymbolTable
}

type WasmFunctionModule struct {
	sectionCode []byte
	typeIndex   int
	funcType    WasmSectionFunctionType
	exportName  *string
}

func (b *WasmFunctionModule) GetIndex() int {
	return b.typeIndex
}

func NewWasmFunctionBuilder(symbolTable *wasmSymbolTable) *WasmFunctionBuilder {
	return &WasmFunctionBuilder{
		paramTypes:   []types.WasmType{},
		resultTypes:  []types.WasmType{},
		name:         nil,
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

func (b *WasmFunctionBuilder) SetExported(name string) *WasmFunctionBuilder {
	b.name = &name
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

	return WasmFunctionModule{
		sectionCode: b.buildFunctionCode(),
		typeIndex:   typeIndex,
		funcType:    funcType,
		exportName:  b.name,
	}
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

	if function.exportName != nil {
		b.sectionExports = append(b.sectionExports, export(*function.exportName, WasmExportDescription{
			Type:  types.ExportFunctionType,
			Index: len(b.sectionFunction) - 1,
		}))
	}

	return b
}

// Save the WASM module to a ".wasm" file. May return an error if the file cannot be created or written to.
func (b *WasmModuleBuilder) BuildWasmFile(fileName string) error {
	if len(fileName) < 5 || fileName[len(fileName)-5:] != ".wasm" {
		fileName += ".wasm"
	}

	return os.WriteFile(fileName, b.Build(), 0644)
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
