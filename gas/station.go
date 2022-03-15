package gas

import (
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// Station interface for check and burn gas
type Station interface {
	Sufficient(addr crypto.Address, gas uint64) bool
	Burn(addr crypto.Address, gas uint64) []*crypto.Event

	// CheckGasPrice checks whether the given price is suitable for current gas station or not
	CheckGasPrice(price uint32) bool

	// Switch decides whether we should replace current gas station for the new one
	// If yes, it will replace the gasStation of the app and return true
	// Otherwise, it returns false
	Switch() bool

	GetPolicy() Policy
}

// Token interface
type Token interface {
	GetBalance(addr crypto.Address) (uint64, error)
	Transfer(caller crypto.Address, addr crypto.Address, amount uint64, memo uint64) ([]*crypto.Event, error)
	GetAccount() *storage.Account
}

// App interface
type App interface {
	SetGasStation(gasStation Station)
	GetGasContractToken() Token
}
