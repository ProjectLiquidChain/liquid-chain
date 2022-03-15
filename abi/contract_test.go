package abi

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestDecodeContract(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	contract := Contract{
		Header: h,
		Code:   []byte{1},
	}
	encodedContract, err := rlp.EncodeToBytes(&contract)
	if err != nil {
		t.Error(err)
	}

	decodedContract, err := DecodeContract(encodedContract)
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})

	if diff := cmp.Diff(*decodedContract, contract, opts); diff != "" {
		t.Errorf("Decode contract %v is incorrect, expected: %v, got: %v, diff: %v", contract, contract, decodedContract, diff)
	}
}

func TestMarshalJSON(t *testing.T) {
	abiFile := "../test/testdata/liquid-token-abi.json"
	want, _ := LoadHeaderFromFile(abiFile)
	code, _ := hex.DecodeString("1")
	contract := Contract{
		Header: want,
		Code:   code,
	}
	jsonBytes, _ := contract.Header.MarshalJSON()

	tmpPath, _ := uuid.NewUUID()
	defer os.RemoveAll(tmpPath.String())
	if err := ioutil.WriteFile(tmpPath.String(), jsonBytes, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	got, _ := LoadHeaderFromFile(tmpPath.String())
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Got header loaded after JSON marshal = %v, want %v", got, want)
	}
}
