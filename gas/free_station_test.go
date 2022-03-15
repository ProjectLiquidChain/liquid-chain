package gas

import (
	"errors"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

const contractAddressStr = "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
const otherAddressStr = "LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY"

type MockFreeToken struct {
	Token
}

func (token *MockFreeToken) GetBalance(addr crypto.Address) (uint64, error) {
	return 100, nil
}

func (token *MockFreeToken) Transfer(caller crypto.Address, addr crypto.Address, amount uint64, memo uint64) ([]*crypto.Event, error) {
	if addr.String() != contractAddressStr {
		panic("Expected collector is gas contract address")
	}
	if amount == 10000 {
		return nil, errors.New("Token transfer failed")
	}
	return []*crypto.Event{}, nil
}

func (token *MockFreeToken) GetAccount() *storage.Account {
	contractAddress, _ := crypto.AddressFromString(contractAddressStr)
	return &storage.Account{
		Nonce:   0,
		Creator: contractAddress,
	}
}

type MockFreeAppNoToken struct {
	App
}

func (app *MockFreeAppNoToken) GetGasContractToken() Token {
	return nil
}

type MockFreeApp struct {
	App
}

func (app *MockFreeApp) SetGasStation(station Station) {

}

func (app *MockFreeApp) GetGasContractToken() Token {
	return &MockFreeToken{}
}

func TestFreeNoSwitch(t *testing.T) {
	app := &MockFreeAppNoToken{}
	station := NewFreeStation(app)
	ret := station.Switch()
	if ret {
		t.Error("Expected return false")
	}
}

func TestFreeSwitch(t *testing.T) {
	app := &MockFreeApp{}
	station := NewFreeStation(app)
	ret := station.Switch()
	if !ret {
		t.Error("Expected return true")
	}
}

func TestFreeSufficient(t *testing.T) {
	app := &MockFreeApp{}
	station := NewFreeStation(app)
	otherAddress, _ := crypto.AddressFromString(otherAddressStr)
	ret := station.Sufficient(otherAddress, 10)

	if !ret {
		t.Error("Expected return true")
	}
}

func TestFreeBurn(t *testing.T) {
	app := &MockFreeApp{}
	station := NewFreeStation(app)
	otherAddress, _ := crypto.AddressFromString(otherAddressStr)
	station.Burn(otherAddress, 10)

	ret := station.Burn(otherAddress, 0)
	if ret != nil {
		t.Error("Expected return nil")
	}
}

func TestFreeCheckGasPrice(t *testing.T) {
	app := &MockFreeAppNoToken{}
	station := NewFreeStation(app)

	if station.CheckGasPrice(0) {
		t.Error("Expected return false")
	}

	if !station.CheckGasPrice(10) {
		t.Error("Expected return true")
	}
}
