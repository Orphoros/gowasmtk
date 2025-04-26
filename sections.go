package gowasmtk

import (
	"github.com/Orphoros/gowasmtk/types"
)

type SectionId = byte

type WasmVector []byte
type WasmSection []byte
type WasmSecionFunction = WasmSection
type WasmSectionType = WasmSection
type WasmSectionExport = WasmSection

type WasmMetadata struct {
	Name    string
	Version string
}

type WasmSectionExportedModule = []byte

type WasmSectionFunctionType = []byte

type WasmExportDescription = struct {
	Type  types.WasmExportType
	Index int
}

const (
	SectionIdCustom   SectionId = 0x00
	SectionIdType     SectionId = 0x01
	SectionIdFunction SectionId = 0x03
	SectionIdCode     SectionId = 0x0A
	SectionIdExport   SectionId = 0x07
)

func same(s string) WasmVector {
	return vec([]byte(s))
}

func export(name string, exportdescs WasmExportDescription) WasmSectionExportedModule {
	var descs []byte

	desc := append([]byte{exportdescs.Type}, leb128EncodeU(uint64(exportdescs.Index))...)
	descs = append(descs, desc...)

	data := append(
		same(name),
		descs...,
	)

	return data
}

func sectionExport(exports ...WasmSectionExportedModule) WasmSectionExport {
	vector := vecNested(exports)

	return section(SectionIdExport, vector)
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

func SectionCode(codes ...[]byte) []byte {
	return section(SectionIdCode, vecNested(codes))
}

func funcType(paramTypes []types.WasmType, resultTypes []types.WasmType) WasmSectionFunctionType {
	// FIXME: Result cannnot be an array.
	return append(
		WasmSectionFunctionType{types.FunctionType},
		append(
			vec(paramTypes),
			vec(resultTypes)...,
		)...,
	)
}

func sectionType(functypes ...WasmSectionFunctionType) WasmSectionType {
	sectionVec := WasmVector{}

	sectionVec = append(sectionVec, vecNested(functypes)...)

	return section(SectionIdType, sectionVec)

}

func sectionFunc(typeidxs ...uint64) WasmSecionFunction {
	var typeidxsBytes []byte
	for _, idx := range typeidxs {
		typeidxsBytes = append(typeidxsBytes, leb128EncodeU(idx)...)
	}

	return section(SectionIdFunction, vec(typeidxsBytes))
}

func section(id SectionId, contents WasmVector) WasmSection {
	wasmSection := WasmSection{}

	wasmSection = append(wasmSection, id)
	wasmSection = append(wasmSection, leb128EncodeU(uint64(len(contents)))...)
	wasmSection = append(wasmSection, contents...)

	return wasmSection
}

func encodeString(s string) []byte {
	return append(leb128EncodeU(uint64(len(s))), []byte(s)...)
}

func sectionCustom(name string, payload []byte) WasmSection {
	// Build custom section data.
	var section []byte
	section = append(section, SectionIdCustom)

	// Custom section name is encoded as a string.
	customName := encodeString(name)
	fullContent := append(customName, payload...)

	// Append section length.
	section = append(section, leb128EncodeU(uint64(len(fullContent)))...)
	// Append the custom section payload.
	section = append(section, fullContent...)

	return section
}

func sectionProducers(languages []WasmMetadata, tools []WasmMetadata, sdk []WasmMetadata) WasmSection {
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
