package gas

import (
	"github.com/vertexdlt/vertexvm/vm"
)

// Policy for gas cost
type Policy interface {
	vm.GasPolicy
	GetCostForStorage(size int) uint64
	GetCostForContract(size int) uint64
	GetCostForEvent(size int) uint64
}
