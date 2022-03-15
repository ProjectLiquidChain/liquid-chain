package gas

import (
	"testing"

	"github.com/vertexdlt/vertexvm/opcode"
)

func TestAlphaPolicy(t *testing.T) {
	policy := AlphaPolicy{}
	cost := policy.GetCostForOp(opcode.Select)
	if cost != 5 {
		t.Errorf("Expect cost %v, got %v", 5, cost)
	}
	cost = policy.GetCostForStorage(100)
	if cost != 100 {
		t.Errorf("Expect cost %v, got %v", 100, cost)
	}
	cost = policy.GetCostForContract(100)
	if cost != 100 {
		t.Errorf("Expect cost %v, got %v", 100, cost)
	}
	cost = policy.GetCostForEvent(100)
	if cost != 100 {
		t.Errorf("Expect cost %v, got %v", 100, cost)
	}
	cost = policy.GetCostForMalloc(1)
	if cost != GasMemoryPage {
		t.Errorf("Expect cost %v, got %v", GasMemoryPage, cost)
	}
}
