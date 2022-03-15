package gas

import (
	"github.com/vertexdlt/vertexvm/opcode"
)

type gasTable [256]uint64

// Cost for generics operations
const (
	GasStack      uint64 = 1
	GasFrame      uint64 = 1
	GasJump       uint64 = 3
	GasBlock      uint64 = 5
	GasNumUnary   uint64 = 1
	GasNumBinary  uint64 = 2
	GasMemory     uint64 = 1
	GasMemoryPage uint64 = 1024
)

func newGasTable() gasTable {
	return gasTable{
		opcode.Block:             GasFrame + GasBlock,
		opcode.Loop:              GasFrame + GasBlock,
		opcode.If:                GasFrame + GasBlock + GasJump,
		opcode.Else:              GasBlock + GasJump,
		opcode.End:               GasFrame + GasStack + GasStack,
		opcode.Br:                GasFrame + GasJump,
		opcode.BrIf:              GasFrame + GasStack + GasNumUnary + GasJump,
		opcode.BrTable:           GasStack + GasFrame + GasJump,
		opcode.Return:            GasJump,
		opcode.Call:              GasFrame + GasStack,
		opcode.CallIndirect:      GasFrame + GasFrame + GasStack,
		opcode.Drop:              GasStack,
		opcode.Select:            GasStack + GasStack + GasStack + GasNumUnary + GasStack,
		opcode.GetLocal:          GasFrame + GasStack + GasStack,
		opcode.SetLocal:          GasFrame + GasStack + GasStack,
		opcode.TeeLocal:          GasFrame + GasStack + GasStack,
		opcode.GetGlobal:         GasFrame + GasStack + GasStack,
		opcode.SetGlobal:         GasFrame + GasStack + GasStack,
		opcode.I32Load:           GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load:           GasFrame + GasFrame + GasStack + GasMemory,
		opcode.F32Load:           GasFrame + GasFrame + GasStack + GasMemory,
		opcode.F64Load:           GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Load8S:         GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Load8U:         GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Load16S:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Load16U:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load8S:         GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load8U:         GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load16S:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load16U:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load32S:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Load32U:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Store:          GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Store:          GasFrame + GasFrame + GasStack + GasMemory,
		opcode.F32Store:          GasFrame + GasFrame + GasStack + GasMemory,
		opcode.F64Store:          GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Store8:         GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I32Store16:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Store8:         GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Store16:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.I64Store32:        GasFrame + GasFrame + GasStack + GasMemory,
		opcode.MemorySize:        GasFrame + GasStack,
		opcode.MemoryGrow:        GasFrame + GasStack + GasStack,
		opcode.I32Const:          GasFrame + GasStack,
		opcode.I64Const:          GasFrame + GasStack,
		opcode.F32Const:          GasFrame + GasStack,
		opcode.F64Const:          GasFrame + GasStack,
		opcode.I32Eqz:            GasStack + GasNumUnary + GasStack,
		opcode.I32Eq:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Ne:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32LtS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32LtU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32GtS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32GtU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32LeS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32LeU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32GeS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32GeU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Eqz:            GasStack + GasNumUnary + GasStack,
		opcode.I64Eq:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Ne:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64LtS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64LtU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64GtS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64GtU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64LeS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64LeU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64GeS:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64GeU:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Eq:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Ne:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Lt:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Gt:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Le:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Ge:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Eq:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Ne:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Lt:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Gt:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Le:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Ge:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Clz:            GasStack + GasNumUnary + GasStack,
		opcode.I32Ctz:            GasStack + GasNumUnary + GasStack,
		opcode.I32Popcnt:         GasStack + GasNumUnary + GasStack,
		opcode.I32Add:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Sub:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Mul:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32DivS:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32DivU:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32RemS:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32RemU:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32And:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Or:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Xor:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Shl:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32ShrS:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32ShrU:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Rotl:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32Rotr:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Clz:            GasStack + GasNumUnary + GasStack,
		opcode.I64Ctz:            GasStack + GasNumUnary + GasStack,
		opcode.I64Popcnt:         GasStack + GasNumUnary + GasStack,
		opcode.I64Add:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Sub:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Mul:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64DivS:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64DivU:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64RemS:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64RemU:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64And:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Or:             GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Xor:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Shl:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64ShrS:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64ShrU:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Rotl:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I64Rotr:           GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Abs:            GasStack + GasNumUnary + GasStack,
		opcode.F32Neg:            GasStack + GasNumUnary + GasStack,
		opcode.F32Ceil:           GasStack + GasNumUnary + GasStack,
		opcode.F32Floor:          GasStack + GasNumUnary + GasStack,
		opcode.F32Trunc:          GasStack + GasNumUnary + GasStack,
		opcode.F32Nearest:        GasStack + GasNumUnary + GasStack,
		opcode.F32Sqrt:           GasStack + GasNumUnary + GasStack,
		opcode.F32Add:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Sub:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Mul:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Div:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Min:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Max:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F32Copysign:       GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Abs:            GasStack + GasNumUnary + GasStack,
		opcode.F64Neg:            GasStack + GasNumUnary + GasStack,
		opcode.F64Ceil:           GasStack + GasNumUnary + GasStack,
		opcode.F64Floor:          GasStack + GasNumUnary + GasStack,
		opcode.F64Trunc:          GasStack + GasNumUnary + GasStack,
		opcode.F64Nearest:        GasStack + GasNumUnary + GasStack,
		opcode.F64Sqrt:           GasStack + GasNumUnary + GasStack,
		opcode.F64Add:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Sub:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Mul:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Div:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Min:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Max:            GasStack + GasStack + GasNumBinary + GasStack,
		opcode.F64Copysign:       GasStack + GasStack + GasNumBinary + GasStack,
		opcode.I32WrapI64:        GasStack + GasNumUnary + GasStack,
		opcode.I32TruncSF32:      GasStack + GasNumUnary + GasStack,
		opcode.I32TruncUF32:      GasStack + GasNumUnary + GasStack,
		opcode.I32TruncSF64:      GasStack + GasNumUnary + GasStack,
		opcode.I32TruncUF64:      GasStack + GasNumUnary + GasStack,
		opcode.I64ExtendSI32:     GasStack + GasNumUnary + GasStack,
		opcode.I64ExtendUI32:     GasStack + GasNumUnary + GasStack,
		opcode.I64TruncSF32:      GasStack + GasNumUnary + GasStack,
		opcode.I64TruncUF32:      GasStack + GasNumUnary + GasStack,
		opcode.I64TruncSF64:      GasStack + GasNumUnary + GasStack,
		opcode.I64TruncUF64:      GasStack + GasNumUnary + GasStack,
		opcode.F32ConvertSI32:    GasStack + GasNumUnary + GasStack,
		opcode.F32ConvertUI32:    GasStack + GasNumUnary + GasStack,
		opcode.F32ConvertSI64:    GasStack + GasNumUnary + GasStack,
		opcode.F32ConvertUI64:    GasStack + GasNumUnary + GasStack,
		opcode.F32DemoteF64:      GasStack + GasNumUnary + GasStack,
		opcode.F64ConvertSI32:    GasStack + GasNumUnary + GasStack,
		opcode.F64ConvertUI32:    GasStack + GasNumUnary + GasStack,
		opcode.F64ConvertSI64:    GasStack + GasNumUnary + GasStack,
		opcode.F64ConvertUI64:    GasStack + GasNumUnary + GasStack,
		opcode.F64PromoteF32:     GasStack + GasNumUnary + GasStack,
		opcode.I32ReinterpretF32: GasStack + GasNumUnary + GasStack,
		opcode.I64ReinterpretF64: GasStack + GasNumUnary + GasStack,
		opcode.F32ReinterpretI32: GasStack + GasNumUnary + GasStack,
		opcode.F64ReinterpretI64: GasStack + GasNumUnary + GasStack,
	}
}

var gasAlphaTable = newGasTable()

// AlphaPolicy is a simple policy for first version
type AlphaPolicy struct {
	Policy
}

// GetCostForOp get cost from table
func (p *AlphaPolicy) GetCostForOp(op opcode.Opcode) uint64 {
	return gasAlphaTable[op]
}

// GetCostForStorage size of data
func (p *AlphaPolicy) GetCostForStorage(size int) uint64 {
	return uint64(size)
}

// GetCostForContract creation
func (p *AlphaPolicy) GetCostForContract(size int) uint64 {
	return uint64(size)
}

// GetCostForEvent emission
func (p *AlphaPolicy) GetCostForEvent(size int) uint64 {
	return uint64(size)
}

// GetCostForMalloc returns cost for new memory allocation
func (p *AlphaPolicy) GetCostForMalloc(pages int) uint64 {
	return GasMemoryPage * uint64(pages)
}
