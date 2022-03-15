package test

import (
	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/engine"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

func TestVMEvent(t *testing.T) {
	contract := loadContract("testdata/event-string-abi.json", "testdata/event-string.wasm")
	contractBytes, _ := rlp.EncodeToBytes(&contract)
	caller, _ := crypto.AddressFromString("LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R")
	contractAddress, _ := crypto.AddressFromString("LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53")

	state := storage.NewStateStorage(db.NewMemoryDB())
	if err := state.LoadState(&crypto.GenesisBlock); err != nil {
		t.Fatal(err)
	}

	accountState, _ := state.CreateAccount(caller, contractAddress, contractBytes)
	execEngine := engine.NewEngine(state, accountState, caller, &gas.FreePolicy{}, 0)

	sayFunction, err := contract.Header.GetFunction("say")
	if err != nil {
		panic(err)
	}
	sayArgs, err := abi.EncodeFromString(sayFunction.Parameters, []string{"1"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := execEngine.Ignite("say", sayArgs); err != nil {
		t.Fatal(err)
	}

	events := execEngine.GetEvents()
	if countEvent := len(events); countEvent != 1 {
		t.Errorf("Expected only one event is emitting, got %v", countEvent)
	}

	expectedMessage := "Checking"
	event, _ := contract.Header.GetEvent("Say")
	messages, _ := abi.DecodeToBytes(event.Parameters, events[0].Args)

	if message := messages[0]; string(message) != expectedMessage {
		t.Errorf("Expected message %v, got %v", expectedMessage, string(message))
	}
}
