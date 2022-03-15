package token

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/storage"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
)

const contractAddressStr = "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
const ownerAddressStr = "LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"
const otherAddressStr = "LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY"
const nonExistentAddressStr = "LANXBHFABEPW5NDSIZUEIENR2LNQHYJ6464NYFVPLE6XKHTMCEZDCLM5"
const contractBalance = uint64(4319)
const ownerBalance = uint64(1000000000 - 10000 - 4319)
const otherBalance = uint64(10000)

func setup() *Token {
	state := storage.NewStateStorage(db.NewMemoryDB())
	if err := state.LoadState(&crypto.GenesisBlock); err != nil {
		panic(err)
	}

	header, err := abi.LoadHeaderFromFile("../test/testdata/gas-token-abi.json")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile("../test/testdata/gas-token.wasm")
	if err != nil {
		panic(err)
	}
	contract := &abi.Contract{
		Header: header,
		Code:   data,
	}
	contractBytes, err := rlp.EncodeToBytes(&contract)
	if err != nil {
		panic(err)
	}
	contractAddress, _ := crypto.AddressFromString(contractAddressStr)
	ownerAddress, _ := crypto.AddressFromString(ownerAddressStr)
	otherAddress, _ := crypto.AddressFromString(otherAddressStr)
	_, err = state.CreateAccount(ownerAddress, contractAddress, contractBytes)
	if err != nil {
		panic(err)
	}
	contractAccount, err := state.LoadAccount(contractAddress)
	if err != nil {
		panic(err)
	}
	token := NewToken(state, contractAccount)
	_, _, err = token.invokeContract(ownerAddress, "init", []string{strconv.FormatUint(1000000000, 10)})
	if err != nil {
		panic(err)
	}
	_, err = token.Transfer(ownerAddress, otherAddress, 10000, 0)
	if err != nil {
		panic(err)
	}
	_, err = token.Transfer(ownerAddress, contractAddress, 4319, 0)
	if err != nil {
		panic(err)
	}
	return token
}

func TestGetBalance(t *testing.T) {
	token := setup()
	contractAddress, _ := crypto.AddressFromString(contractAddressStr)
	ownerAddress, _ := crypto.AddressFromString(ownerAddressStr)
	otherAddress, _ := crypto.AddressFromString(otherAddressStr)
	nonExistentAddress, _ := crypto.AddressFromString(nonExistentAddressStr)

	ret, err := token.GetBalance(contractAddress)
	if err != nil {
		panic(err)
	}
	if ret != contractBalance {
		t.Errorf("Expect contract balance to be %v, got %v", contractBalance, ret)
	}
	ret, err = token.GetBalance(ownerAddress)
	if ret != ownerBalance {
		t.Errorf("Expect owner balance to be %v, got %v", ownerBalance, ret)
	}
	ret, err = token.GetBalance(otherAddress)
	if ret != otherBalance {
		t.Errorf("Expect other balance to be %v, got %v", otherBalance, ret)
	}
	ret, err = token.GetBalance(nonExistentAddress)
	if ret != 0 {
		t.Errorf("Expect non-existent balance to be %v, got %v", 0, ret)
	}
}

func TestTransferOK(t *testing.T) {
	token := setup()

	collector, _ := crypto.AddressFromString(contractAddressStr)
	caller, _ := crypto.AddressFromString(otherAddressStr)
	amount := uint64(100)

	events, err := token.Transfer(caller, collector, amount, 0)
	if err != nil {
		panic(err)
	}
	if len(events) != 1 {
		t.Errorf("Expect %v transfer event, got %v", 1, len(events))
	}
	ret, err := token.GetBalance(caller)
	if ret != otherBalance-amount {
		t.Errorf("Expect caller balance to be %v, got %v", otherBalance-amount, ret)
	}
	ret, err = token.GetBalance(collector)
	if ret != contractBalance+amount {
		t.Errorf("Expect collector balance to be %v, got %v", contractBalance+amount, ret)
	}
}

func TestTransferFail(t *testing.T) {
	token := setup()

	collector, _ := crypto.AddressFromString(contractAddressStr)
	caller, _ := crypto.AddressFromString(nonExistentAddressStr)
	_, err := token.Transfer(caller, collector, 100, 0)
	if err == nil || err.Error() != "process exit with code: 1" {
		t.Errorf("Expect token transfer failed")
	}
}

func TestInvokeUndefinedFunction(t *testing.T) {
	token := setup()
	ownerAddress, _ := crypto.AddressFromString(ownerAddressStr)
	_, _, err := token.invokeContract(ownerAddress, "undefined_function", []string{})
	if err == nil {
		t.Errorf("Expect contract invoke error")
	}
}

func TestInvokeWrongParameters(t *testing.T) {
	token := setup()
	ownerAddress, _ := crypto.AddressFromString(ownerAddressStr)
	_, _, err := token.invokeContract(ownerAddress, "mint", []string{})
	if err == nil {
		t.Errorf("Expect contract invoke error")
	}
}
