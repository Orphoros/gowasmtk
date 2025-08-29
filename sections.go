package gowasmtk

import (
	"github.com/Orphoros/gowasmtk/types"
)

type sectionId = byte

type wasmVector []byte
type wasmSection []byte
type WasmSecionFunction = wasmSection
type WasmSectionType = wasmSection
type WasmSectionExport = wasmSection
type WasmSectionImport = wasmSection

type wasmMetadata struct {
	Name    string
	Version string
}

type wasmSectionExportedModule = []byte
type wasmSectionFunctionType = []byte
type wasmSectionImportedModule = []byte

type wasmExportDescription = struct {
	Type  types.WasmExportType
	Index int
}

var importdesc = struct {
	function func(index uint32) wasmSectionImportedModule
}{
	func(index uint32) wasmSectionImportedModule {
		ve := wasmSectionImportedModule{types.ImportFunctionType}
		ve = append(ve, leb128EncodeU(uint64(index))...)
		return ve
	},
}

const (
	sectionIdCustom   sectionId = 0x00
	sectionIdType     sectionId = 0x01
	sectionIdImport   sectionId = 0x02
	sectionIdFunction sectionId = 0x03
	sectionIdCode     sectionId = 0x0A
	sectionIdExport   sectionId = 0x07
)

func name(s string) wasmVector {
	return vec([]byte(s))
}

func imports(modName, funcName string, importdesc wasmSectionImportedModule) wasmSectionImportedModule {
	return append(name(modName), append(name(funcName), importdesc...)...)
}

func sectionImports(imports ...wasmSectionImportedModule) WasmSectionImport {
	sectionVec := wasmVector{}

	sectionVec = append(sectionVec, vecNested(imports)...)

	return section(sectionIdImport, sectionVec)
}

func export(exportName string, exportdescs wasmExportDescription, numImportDeclarations int) wasmSectionExportedModule {
	var descs []byte

	desc := append([]byte{exportdescs.Type}, leb128EncodeU(uint64(exportdescs.Index+numImportDeclarations))...)
	descs = append(descs, desc...)

	data := append(
		name(exportName),
		descs...,
	)

	return data
}

func sectionExport(exports ...wasmSectionExportedModule) WasmSectionExport {
	vector := vecNested(exports)

	return section(sectionIdExport, vector)
}

func code(f []byte) []byte {
	return append(
		leb128EncodeU(uint64(len(f))),
		f...,
	)
}

func function(locals [][]byte, body []byte) []byte {
	return append(
		vecNested(locals),
		body...,
	)
}

func sectionCode(codes ...[]byte) []byte {
	return section(sectionIdCode, vecNested(codes))
}

func funcType(paramTypes []types.WasmType, resultTypes []types.WasmType) wasmSectionFunctionType {
	// FIXME: Result cannnot be an array.
	return append(
		wasmSectionFunctionType{types.FunctionType},
		append(
			vec(paramTypes),
			vec(resultTypes)...,
		)...,
	)
}

func sectionType(functypes ...wasmSectionFunctionType) WasmSectionType {
	sectionVec := wasmVector{}

	sectionVec = append(sectionVec, vecNested(functypes)...)

	return section(sectionIdType, sectionVec)

}

func sectionFunc(typeidxs ...uint64) WasmSecionFunction {
	var typeidxsBytes []byte
	for _, idx := range typeidxs {
		typeidxsBytes = append(typeidxsBytes, leb128EncodeU(idx)...)
	}

	return section(sectionIdFunction, vec(typeidxsBytes))
}

func section(id sectionId, contents wasmVector) wasmSection {
	wasmSection := wasmSection{}

	wasmSection = append(wasmSection, id)
	wasmSection = append(wasmSection, leb128EncodeU(uint64(len(contents)))...)
	wasmSection = append(wasmSection, contents...)

	return wasmSection
}

func encodeString(s string) []byte {
	return append(leb128EncodeU(uint64(len(s))), []byte(s)...)
}

func sectionCustom(name string, payload []byte) wasmSection {
	// Build custom section data.
	var section []byte
	section = append(section, sectionIdCustom)

	// Custom section name is encoded as a string.
	customName := encodeString(name)
	fullContent := append(customName, payload...)

	// Append section length.
	section = append(section, leb128EncodeU(uint64(len(fullContent)))...)
	// Append the custom section payload.
	section = append(section, fullContent...)

	return section
}

func sectionProducers(languages []wasmMetadata, tools []wasmMetadata, sdk []wasmMetadata) wasmSection {
	var payload []byte
	var length uint64 = 0

	if len(languages) > 0 {
		length += 1
		payload = append(payload, encodeString("language")...)
		payload = append(payload, leb128EncodeU(uint64(len(languages)))...)
		for _, lang := range languages {
			payload = append(payload, encodeString(lang.Name)...)
			payload = append(payload, encodeString(lang.Version)...)
		}
	}
	if len(tools) > 0 {
		length += 1
		payload = append(payload, encodeString("processed-by")...)
		payload = append(payload, leb128EncodeU(uint64(len(tools)))...)
		for _, tool := range tools {
			payload = append(payload, encodeString(tool.Name)...)
			payload = append(payload, encodeString(tool.Version)...)
		}
	}
	if len(sdk) > 0 {
		length += 1
		payload = append(payload, encodeString("sdk")...)
		payload = append(payload, leb128EncodeU(uint64(len(sdk)))...)
		for _, sdkItem := range sdk {
			payload = append(payload, encodeString(sdkItem.Name)...)
			payload = append(payload, encodeString(sdkItem.Version)...)
		}
	}

	return sectionCustom("producers", append(leb128EncodeU(length), payload...))
}
