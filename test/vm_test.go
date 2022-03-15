package test

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"

	"github.com/QuoineFinancial/liquid-chain/engine"
)

func TestVM(t *testing.T) {
	contract := loadContract("testdata/liquid-token-abi.json", "testdata/liquid-token.wasm")
	header := contract.Header
	contractBytes, _ := rlp.EncodeToBytes(&contract)
	caller, _ := crypto.AddressFromString("LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R")
	contractAddress, _ := crypto.AddressFromString("LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53")

	state := storage.NewStateStorage(db.NewMemoryDB())
	if err := state.LoadState(&crypto.GenesisBlock); err != nil {
		t.Fatal(err)
	}

	accountState, _ := state.CreateAccount(caller, contractAddress, contractBytes)
	execEngine := engine.NewEngine(state, accountState, caller, &gas.FreePolicy{}, 0)
	toAddress := "LB3EICIUKOUYCY4D7T2O6RKL7ISEPISNKUXNILDTJ76V2PDZVT5ZDP3U"
	var mint = 10000000
	mintAmount := strconv.Itoa(mint)
	var transfer = 30
	transferAmount := strconv.Itoa(transfer)

	mintFunction, err := header.GetFunction("mint")
	if err != nil {
		panic(err)
	}
	mintArgs, err := abi.EncodeFromString(mintFunction.Parameters, []string{mintAmount})
	if err != nil {
		panic(err)
	}
	_, err = execEngine.Ignite("mint", mintArgs)
	if err != nil {
		panic(err)
	}

	transferFunction, err := header.GetFunction("transfer")
	if err != nil {
		panic(err)
	}
	transferArgs, err := abi.EncodeFromString(transferFunction.Parameters, []string{toAddress, transferAmount})
	if err != nil {
		panic(err)
	}
	_, err = execEngine.Ignite("transfer", transferArgs)
	if err != nil {
		panic(err)
	}

	getBalanceFunction, err := header.GetFunction("get_balance")
	if err != nil {
		panic(err)
	}
	getBalanceMint, _ := abi.EncodeFromString(getBalanceFunction.Parameters, []string{caller.String()})
	if err != nil {
		panic(err)
	}

	getBalanceTo, err := abi.EncodeFromString(getBalanceFunction.Parameters, []string{toAddress})
	if err != nil {
		panic(err)
	}
	ret, err := execEngine.Ignite("get_balance", getBalanceTo)
	if err != nil {
		t.Error(err)
	}
	if int(ret) != transfer {
		t.Errorf("Expect return value to be %v, got %v", transfer, ret)
	}
	ret, err = execEngine.Ignite("get_balance", getBalanceMint)
	if err != nil {
		t.Error(err)
	}
	if int(ret) != mint-transfer {
		t.Errorf("Expect return value to be %v, got %v", mint-transfer, ret)
	}
}

func loadContract(abiPath, wasmPath string) *abi.Contract {
	header, err := abi.LoadHeaderFromFile(abiPath)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		panic(err)
	}

	return &abi.Contract{
		Header: header,
		Code:   data,
	}
}
