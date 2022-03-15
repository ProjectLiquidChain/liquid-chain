package gas

import (
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// DummyStation provide a dummy gas station only for testing purpose
type DummyStation struct {
	app    App
	policy Policy
}

// Sufficient gas of an address is enough for burn
func (station *DummyStation) Sufficient(addr crypto.Address, gas uint64) bool {
	return gas != 0
}

// Burn gas, do nothing
func (station *DummyStation) Burn(addr crypto.Address, gas uint64) []*crypto.Event {
	return nil
}

// Switch on fee
func (station *DummyStation) Switch() bool {
	return false
}

// GetPolicy free
func (station *DummyStation) GetPolicy() Policy {
	return station.policy
}

// CheckGasPrice of transaction
func (station *DummyStation) CheckGasPrice(price uint32) bool {
	return false
}

// NewDummyStation constructor
func NewDummyStation(app App) Station {
	return &DummyStation{
		app:    app,
		policy: &FreePolicy{},
	}
}
