package gas

import (
	"github.com/vertexdlt/vertexvm/opcode"
)

// FreePolicy is a simple policy for first version
type FreePolicy struct {
	Policy
}

// GetCostForOp get cost from table
func (p *FreePolicy) GetCostForOp(op opcode.Opcode) uint64 {
	return 0
}

// GetCostForStorage size of data
func (p *FreePolicy) GetCostForStorage(size int) uint64 {
	return 0
}

// GetCostForContract creation
func (p *FreePolicy) GetCostForContract(size int) uint64 {
	return 0
}

// GetCostForEvent emission
func (p *FreePolicy) GetCostForEvent(size int) uint64 {
	return 0
}

// GetCostForMalloc returns cost for new memory allocation
func (p *FreePolicy) GetCostForMalloc(pages int) uint64 {
	return 0
}
