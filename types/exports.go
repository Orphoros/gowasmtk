package types

type WasmExportType = WasmType

const (
	ExportFunctionType WasmExportType = 0x00
	ExportTableType    WasmExportType = 0x01
	ExportMemoryType   WasmExportType = 0x02
)
