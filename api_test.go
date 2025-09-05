package gowasmtk

import (
	"log"
	"testing"

	"github.com/Orphoros/gowasmtk/instructions"
	"github.com/Orphoros/gowasmtk/types"

	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

type apiTestCase struct {
	input      *WasmModuleBuilder
	nameOfMain string
	args       []interface{}
	expected   interface{}
}

func TestModule(t *testing.T) {
	t.Run("should build a void wasm module", func(t *testing.T) {
		functionType := funcType([]types.WasmType{}, []types.WasmType{})

		typeSection := sectionType(functionType)

		funcSection := sectionFunc(0)

		mainFunc := function([][]byte{}, []byte{instructions.End})

		mainFuncCode := code(mainFunc)

		codeSection := sectionCode(
			mainFuncCode,
		)
		exportedMainFunc := export("main", wasmExportDescription{
			Type:  types.ExportFunctionType,
			Index: 0,
		}, 0)

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

	t.Run("should build a wasm module with function call arithmetic", func(t *testing.T) {
		wasmSymbolTable := NewSymbolTable(nil)

		adder := NewWasmFunctionBuilder(wasmSymbolTable).
			AddParam(types.I32).
			AddParam(types.I32).
			AddReturn(types.I32).
			AddInstrGetLocal(0).
			AddInstrGetLocal(1).
			AddInstrAddI32().
			AddInstrEnd().
			Build()

		main := NewWasmFunctionBuilder(wasmSymbolTable).
			AddParam(types.I32).
			AddReturn(types.I32).
			AddLocal(1, types.I32).
			AddInstrConstI32(10).
			AddInstrSetLocal(1).
			AddInstrGetLocal(0).
			AddInstrGetLocal(1).
			AddInstrCall(&adder).
			AddInstrEnd().
			Build()

		mod := NewWasmModuleBuilder(wasmSymbolTable).
			AddFunction(&adder).
			AddFunction(&main).
			Export("main", types.ExportFunctionType, &main).
			AddMetaSdk("Orp", "0.0.1").
			AddMetaLanguage("Shark", "0.0.1").
			AddMetaTool("GoWasmTK", "0.0.1")

		runModValueTest(t, apiTestCase{
			input:      mod,
			nameOfMain: "main",
			args:       []interface{}{45},
			expected:   int32(55), // 45 + 10
		})

	})

	t.Run("should build a wasm module with conditional", func(t *testing.T) {
		wasmSymbolTable := NewSymbolTable(nil)
		main := NewWasmFunctionBuilder(wasmSymbolTable).
			// get a i32 from parameter, if it is 0, return 0, else return 1
			AddParam(types.I32).
			AddReturn(types.I32).
			AddInstrGetLocal(0).
			AddInstrEqzI32().
			AddInstrIf(types.I32).
			AddInstrConstI32(0).
			AddInstrElse().
			AddInstrConstI32(1).
			AddInstrEnd().
			AddInstrEnd().
			Build()

		mod := NewWasmModuleBuilder(wasmSymbolTable).
			AddFunction(&main).
			Export("main", types.ExportFunctionType, &main).
			AddMetaSdk("Orp", "0.0.1").
			AddMetaLanguage("Shark", "0.0.1").
			AddMetaTool("GoWasmTK", "0.0.1")

		runModValueTest(t, apiTestCase{
			input:      mod,
			nameOfMain: "main",
			args:       []interface{}{0},
			expected:   int32(0),
		})

		runModValueTest(t, apiTestCase{
			input:      mod,
			nameOfMain: "main",
			args:       []interface{}{42},
			expected:   int32(1),
		})
	})

	t.Run("should build a wasm module with loop", func(t *testing.T) {
		wasmSymbolTable := NewSymbolTable(nil)
		// get a i32 from parameter, calculate fibonacci
		fib := NewWasmFunctionBuilder(wasmSymbolTable).
			AddParam(types.I32).
			AddReturn(types.I32).
			// if (n < 2)
			AddInstrGetLocal(0).
			AddInstrConstI32(2).
			AddInstrLessThanI32S().
			AddInstrIf(types.I32).
			// then: return n
			AddInstrGetLocal(0).
			AddInstrElse().
			// else: return fib(n-2) + fib(n-1)
			AddInstrGetLocal(0).
			AddInstrConstI32(2).
			AddInstrSubI32().
			AddInstrCallSelf().
			AddInstrGetLocal(0).
			AddInstrConstI32(1).
			AddInstrSubI32().
			AddInstrCallSelf().
			AddInstrAddI32().
			AddInstrEnd(). // close if block
			AddInstrEnd(). // close function body
			Build()
		main := NewWasmFunctionBuilder(wasmSymbolTable).
			// get a i32 from parameter, calculate fibonacci
			AddParam(types.I32).
			AddReturn(types.I32).
			AddInstrGetLocal(0).
			AddInstrCall(&fib).
			AddInstrEnd().
			Build()

		mod := NewWasmModuleBuilder(wasmSymbolTable).
			AddFunction(&main).
			AddFunction(&fib).
			Export("main", types.ExportFunctionType, &main).
			AddMetaSdk("Orp", "0.0.1").
			AddMetaLanguage("Shark", "0.0.1").
			AddMetaTool("GoWasmTK", "0.0.1")

		runModValueTest(t, apiTestCase{
			input:      mod,
			nameOfMain: "main",
			args:       []interface{}{10},
			expected:   int32(55), // fib of 10 is 55
		})

		runModValueTest(t, apiTestCase{
			input:      mod,
			nameOfMain: "main",
			args:       []interface{}{20},
			expected:   int32(6765), // fib of 20 is 6765
		})
	})

	t.Run("should build a wasm module with imports", func(t *testing.T) {
		imports := []WasmImportDeclaration{
			{
				ModuleName:   "basicMath",
				FunctionName: "addOne",
				ParamTypes:   []types.WasmType{types.I32},
				ResultTypes:  []types.WasmType{types.I32},
			},
		}
		wasmSymbolTable := NewSymbolTable(&imports)
		adder := NewWasmFunctionBuilder(wasmSymbolTable).
			AddParam(types.I32).
			AddParam(types.I32).
			AddReturn(types.I32).
			AddInstrGetLocal(0).
			AddInstrGetLocal(1).
			AddInstrAddI32().
			AddInstrEnd().
			Build()

		main := NewWasmFunctionBuilder(wasmSymbolTable).
			// get a i32 from parameter, call adder with the constant of 10, then call imported addOne, return the result
			AddParam(types.I32).
			AddReturn(types.I32).
			AddLocal(1, types.I32).
			AddInstrConstI32(10).
			AddInstrSetLocal(1).
			AddInstrGetLocal(0).
			AddInstrGetLocal(1).
			AddInstrCall(&adder).
			AddInstrCallImport(&imports[0]). // call the first import
			AddInstrEnd().
			Build()

		mod := NewWasmModuleBuilder(wasmSymbolTable).
			AddFunction(&adder).
			AddFunction(&main).
			Export("main", types.ExportFunctionType, &main).
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

func runModValueTest(t *testing.T, test apiTestCase) {
	t.Helper()
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, err := wasmer.NewModule(store, test.input.Build())

	if err != nil {
		log.Fatal("module compilation error: %w\n", err)
		return
	}

	importObject := wasmer.NewImportObject()
	instance, err := wasmer.NewInstance(module, importObject)

	if err != nil {
		log.Fatal("instance error: %w\n", err)
		return
	}

	main, err := instance.Exports.GetFunction(test.nameOfMain)

	if err != nil {
		log.Fatal("function retrieval error: %w\n", err)
		return
	}

	result, err := main(test.args...)

	if err != nil {
		log.Fatal("function call error: %w\n", err)
		return
	}

	if result != test.expected {
		t.Fatalf("expected %v, got %v", test.expected, result)
	}

}
