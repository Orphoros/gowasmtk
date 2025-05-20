package gowasmtk

import (
	"github.com/Orphoros/gowasmtk/types"
	"github.com/jcalabro/leb128"
)

func vec(elements []byte) wasmVector {
	var result []byte

	result = append(result, leb128EncodeU(uint64(len(elements)))...)
	result = append(result, elements...)

	return result
}

func vecNested(elements [][]byte) wasmVector {
	var result []byte

	result = append(result, leb128EncodeU(uint64(len(elements)))...)

	for _, element := range elements {
		result = append(result, element...)
	}

	return result
}

func leb128EncodeU(n uint64) wasmVector {
	return leb128.EncodeU64(n)
}

func leb128EncodeI(n int64) wasmVector {
	return leb128.EncodeS64(n)
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

func version() []byte {
	return []byte{0x01, 0x00, 0x00, 0x00}
}

func magic() []byte {
	return stringToBytes("\000asm")
}

func module(sections ...wasmSection) []byte {
	mod := []byte{}

	mod = append(mod, magic()...)
	mod = append(mod, version()...)
	for _, section := range sections {
		mod = append(mod, section...)
	}

	return mod
}

func locals(amount uint32, localType types.WasmType) []byte {
	return append(
		leb128EncodeU(uint64(amount)),
		localType,
	)
}
