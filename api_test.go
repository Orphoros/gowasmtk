package gowasmtk

import (
	"log"
	"testing"

	"github.com/Orphoros/gowasmtk/instructions"
	"github.com/Orphoros/gowasmtk/types"
)

func TestModule(t *testing.T) {
	t.Run("should build a void wasm module", func(t *testing.T) {
		functionType := funcType([]types.WasmType{}, []types.WasmType{})

		typeSection := sectionType(functionType)

		funcSection := sectionFunc(0)

		mainFunc := function([][]byte{}, []byte{instructions.End})

		mainFuncCode := code(mainFunc)

		codeSection := SectionCode(
			mainFuncCode,
		)
		exportedMainFunc := export("main", WasmExportDescription{
			Type:  types.ExportFunctionType,
			Index: 0,
		})

		exportSection := sectionExport(exportedMainFunc)

		mod := module(typeSection, funcSection, exportSection, codeSection)

		expected := []byte{
			0x00, 0x61, 0x73, 0x6D, // magic
			0x01, 0x00, 0x00, 0x00, // version
			0x01, 0x04, 0x01, 0x60, 0x00, 0x00, // type section
			0x03, 0x02, 0x01, 0x00, // function section
			0x07, 0x08, 0x01, 0x04, 0x6d, 0x61, 0x69, 0x6E, 0x00, 0x00, // export section
			0x0A, 0x04, 0x01, 0x02, 0x00, 0x0B, // code section
		}

		if len(mod) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(mod))
		}

		for i := 0; i < len(mod); i++ {
			if mod[i] != expected[i] {
				t.Fatalf("expected %x, got %x", expected[i], mod[i])
			}
		}
	})

	t.Run("should build a wasm module with arithmetic", func(t *testing.T) {
		wasmSymbolTable := NewSymbolTable()

		f1 := NewWasmFunctionBuilder(wasmSymbolTable).
			AddParam(types.I32).
			AddReturn(types.I32).
			AddLocal(1, types.I32).
			AddInstruction(instructions.ConstI32).
			AddI32(42).
			AddInstruction(instructions.SetLocal).
			AddI32(1).
			AddInstruction(instructions.GetLocal).
			AddI32(0).
			AddInstruction(instructions.GetLocal).
			AddI32(1).
			AddInstruction(instructions.AddI32).
			AddInstruction(instructions.End).
			Build()

		f2 := NewWasmFunctionBuilder(wasmSymbolTable).
			AddParam(types.I32).
			AddReturn(types.I32).
			AddLocal(1, types.I32).
			AddInstruction(instructions.ConstI32).
			AddI32(10).
			AddInstruction(instructions.SetLocal).
			AddI32(1).
			AddInstruction(instructions.GetLocal).
			AddI32(0).
			AddInstruction(instructions.GetLocal).
			AddI32(1).
			AddInstruction(instructions.AddI32).
			AddInstruction(instructions.End).
			Build()

		mod := NewWasmModuleBuilder(wasmSymbolTable).
			AddFunction(&f1).
			AddFunction(&f2).
			Export("f1", types.ExportFunctionType, &f1).
			Export("f2", types.ExportFunctionType, &f2).
			AddMetaSdk("Orp", "0.0.1").
			AddMetaLanguage("Shark", "0.0.1").
			AddMetaTool("GoWasmTK", "0.0.1")

		err := mod.BuildWasmFile("mod.wasm")
		if err != nil {
			log.Fatal("error: %w\n", err)
			return
		}
	})
}
