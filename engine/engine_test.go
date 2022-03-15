package engine

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"golang.org/x/crypto/blake2b"
)

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

func TestEngineIgnite(t *testing.T) {
	contractCreator, _ := crypto.AddressFromString("LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R")
	mathAddress, _ := crypto.AddressFromString("LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53")
	utilAddress, _ := crypto.AddressFromString("LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY")
	state := storage.NewStateStorage(db.NewMemoryDB())
	if err := state.LoadState(&crypto.Block{
		Height: 1,
		Time:   uint64(time.Unix(1578905663, 0).UTC().Unix()),
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		callee        *abi.Contract
		calleeAddress crypto.Address
		caller        *abi.Contract
		callerAddress crypto.Address
		funcName      string
		args          []string
		want          uint64
		wantErr       string
	}{
		{
			name:          "chained ignite",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "hypotenuse",
			args:          []string{"3", "4"},
			want:          math.Float64bits(5),
		},
		{
			name:          "chained ignite with array param",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "variance",
			args:          []string{"[1,2,3,4,5]"},
			want:          2,
		},
		{
			name:          "chained ignite with events",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "xor_checksum",
			args:          []string{"LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"},
			want:          149,
		},
		{
			name:          "ignite unknown imported function",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "average",
			args:          []string{"[1,2,3,4,5]"},
			wantErr:       "unknown import get_average for module env",
		},
		{
			name:          "chained ignite overflow",
			callee:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			calleeAddress: utilAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "mean",
			args:          []string{"[1,2,3,4,5]"},
			wantErr:       "call depth limit reached",
		},
		{
			name:          "ignite block time",
			caller:        loadContract("testdata/blockinfo-abi.json", "testdata/blockinfo.wasm"),
			callerAddress: utilAddress,
			funcName:      "block_time",
			args:          []string{},
			want:          1578905663,
		},
		{
			name:          "ignite block height",
			caller:        loadContract("testdata/blockinfo-abi.json", "testdata/blockinfo.wasm"),
			callerAddress: utilAddress,
			funcName:      "block_height",
			args:          []string{},
			want:          1,
		},
		{
			name:          "chained ignite with invoke address validation",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "mod_invoke",
			args:          []string{"LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"},
			wantErr:       "Unexpected address version 0",
		},
		{
			name:          "chained ignite with event address validation",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "mod_emit",
			args:          []string{"LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"},
			wantErr:       "Unexpected address version 0",
		},
		{
			name:          "chained ignite with event plarray",
			callee:        nil,
			calleeAddress: crypto.EmptyAddress,
			caller:        loadContract("testdata/event-string-abi.json", "testdata/event-string.wasm"),
			callerAddress: utilAddress,
			funcName:      "say",
			args:          []string{"0"},
			want:          0,
		},
		{
			name:          "chained ignite with mismatched params",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "parity",
			args:          []string{"1", "2"},
			wantErr:       "argument count mismatch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contractBytes, _ := rlp.EncodeToBytes(&tt.caller)
			callerAccount, _ := state.CreateAccount(contractCreator, tt.callerAddress, contractBytes)
			execEngine := NewEngine(state, callerAccount, contractCreator, &gas.FreePolicy{}, 0)
			if tt.callee != nil {
				if tt.calleeAddress.String() != tt.callerAddress.String() {
					contractBytes, _ := rlp.EncodeToBytes(&tt.callee)
					if _, err := state.CreateAccount(contractCreator, tt.calleeAddress, contractBytes); err != nil {
						panic(err)
					}

				}
				// contract init
				initFunc := "init"
				function, err := tt.caller.Header.GetFunction(initFunc)
				if err != nil {
					panic(err)
				}
				args, err := abi.EncodeFromString(function.Parameters, []string{tt.calleeAddress.String()})
				if err != nil {
					panic(err)
				}
				_, err = execEngine.Ignite(initFunc, args)
				if err != nil {
					panic(err)
				}
			}

			// exec
			function, err := tt.caller.Header.GetFunction(tt.funcName)
			if err != nil {
				panic(err)
			}
			args, err := abi.EncodeFromString(function.Parameters, tt.args)
			if err != nil {
				panic(err)
			}
			got, err := execEngine.Ignite(tt.funcName, args)
			errString := ""
			if err != nil {
				errString = err.Error()
			}
			if tt.wantErr != errString {
				t.Errorf("Engine.Ignite() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && got != tt.want {
				t.Errorf("Engine.Ignite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelegatedCall(t *testing.T) {
	seedStr := "38621fc10a1e56089192360454da6277ce6289bc8f6ef29f34c62a0267940e50"
	trigger, _ := crypto.AddressFromString("LCVYYRXAH2OJY7ONOO5QNRWOAVQXPXIA6T2C6PE2BFBEZX56UBZOJ5GX")
	holder, _ := crypto.AddressFromString("LCFDZPMMHQTNPX64NQR2D5GBJVITHC2M7VFVFLP24YVRKOSI5CSDV4SS")
	contract, _ := crypto.AddressFromString("LBHYGTPXRBYOVQ74XAWIZPJL3ZM4XN27QP53ERQQTKF4BQE4RMIVYBWT")
	amount := uint64(50)
	nonce := uint32(0)
	values := []interface{}{trigger, contract, amount, nonce}
	parameters := []*abi.Parameter{
		&abi.Parameter{Type: abi.Address},
		&abi.Parameter{Type: abi.Address},
		&abi.Parameter{Type: abi.Uint64},
		&abi.Parameter{Type: abi.Uint32},
	}
	result, err := abi.Encode(parameters, values)
	if err != nil {
		panic(err)
	}
	seed, err := hex.DecodeString(seedStr)
	if err != nil {
		panic(err)
	}
	hash := blake2b.Sum256(result)
	privKey := ed25519.NewKeyFromSeed(seed)
	signature := ed25519.Sign(privKey, hash[:])

	pubkey, _ := holder.PubKey()
	if !ed25519.Verify(pubkey, hash[:], signature) {
		panic("invalid signature")
	}

	abiContract := loadContract("testdata/delegated-token-abi.json", "testdata/delegated-token.wasm")
	contractBytes, _ := rlp.EncodeToBytes(abiContract)
	state := storage.NewStateStorage(db.NewMemoryDB())
	if err = state.LoadState(&crypto.Block{
		Height: 1,
		Time:   1578905663,
	}); err != nil {
		t.Fatal(err)
	}

	contractAccount, _ := state.CreateAccount(holder, contract, contractBytes)
	execEngine := NewEngine(state, contractAccount, holder, &gas.FreePolicy{}, 0)

	function, err := abiContract.Header.GetFunction("mint")
	if err != nil {
		panic(err)
	}
	args, err := abi.EncodeFromString(function.Parameters, []string{fmt.Sprintf("%v", amount)})
	if err != nil {
		panic(err)
	}
	got, err := execEngine.Ignite("mint", args)
	if err != nil {
		panic(err)
	}

	function, err = abiContract.Header.GetFunction("delegated_transfer")
	if err != nil {
		panic(err)
	}
	args, err = abi.EncodeFromString(function.Parameters, []string{
		trigger.String(),
		fmt.Sprintf("%v", amount),
		holder.String(),
		fmt.Sprintf("%v", nonce),
		strings.Join(strings.Split(fmt.Sprint(signature), " "), ","),
	})
	if err != nil {
		panic(err)
	}
	got, err = execEngine.Ignite("delegated_transfer", args)
	if got != 0 {
		t.Errorf("Engine.Ignite() = %v, want %v", got, 0)
	}
}
