package consensus

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/util"
)

type TestResource struct {
	app   *App
	dbDir string
}

func newTestResource() *TestResource {
	dbDir := "./tmp" + strconv.Itoa(rand.Intn(10000))
	err := os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	app := NewApp(dbDir, "")
	if err := app.State.LoadState(&crypto.GenesisBlock); err != nil {
		panic(err)
	}
	return &TestResource{
		app:   app,
		dbDir: dbDir,
	}
}

func (tr *TestResource) cleanData() {
	err := os.RemoveAll(tr.dbDir)
	if err != nil {
		panic(err)
	}
}

func TestApplyTx(t *testing.T) {
	tr := newTestResource()
	defer tr.cleanData()

	seed := make([]byte, 32)
	rand.Read(seed)

	seed2 := make([]byte, 32)
	rand.Read(seed2)

	// Setup deploy contract transaction
	sender := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)

	sender2 := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: ed25519.NewKeyFromSeed(seed2).Public().(ed25519.PublicKey),
	}

	data, err := util.BuildDeployTxPayload("../test/testdata/liquid-token.wasm", "../test/testdata/liquid-token-abi.json", "", []string{})
	if err != nil {
		t.Fatal(err)
	}
	deployTx := &crypto.Transaction{
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 0,
	}

	contractWithInitTxData, err := util.BuildDeployTxPayload("../test/testdata/liquid-token.wasm", "../test/testdata/liquid-token-abi.json", "init", []string{"1000"})
	if err != nil {
		panic(err)
	}
	deployWithInitTx := &crypto.Transaction{
		Sender:   &sender,
		Payload:  contractWithInitTxData,
		Receiver: crypto.EmptyAddress,
		GasLimit: 100000,
		GasPrice: 18,
	}

	// Setup result events after deploy contract
	contractAddress := crypto.NewDeploymentAddress(senderAddress, sender.Nonce)

	// Setup invoke contract transaction
	invokeData, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	invokeTx := &crypto.Transaction{
		Sender:   &sender,
		Receiver: contractAddress,
		Payload:  invokeData,
		GasLimit: 0,
		GasPrice: 0,
	}
	contractHeader, _ := abi.LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	mintEventHeader, _ := contractHeader.GetEvent("Mint")
	mintAmount := make([]byte, abi.Uint64.GetMemorySize())
	binary.LittleEndian.PutUint64(mintAmount, 1000)
	mintEventData, _ := abi.EncodeFromBytes(mintEventHeader.Parameters, [][]byte{senderAddress[:], mintAmount})

	// Setup invoke contract transaction
	invalidMethodInvokePayload, _ := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	invalidMethodInvokePayload.ID = crypto.GetMethodID("mint-invalid")
	invalidMethodInvokeTx := &crypto.Transaction{
		Sender:   &sender,
		Receiver: contractAddress,
		Payload:  invalidMethodInvokePayload,
		GasLimit: 0,
		GasPrice: 0,
	}

	invalidInitMethodInvokePayload, _ := util.BuildDeployTxPayload("../test/testdata/event-string.wasm", "../test/testdata/event-string-abi.json", "say", []string{"100"})
	invalidInitMethodInvokePayload.ID = initFunctionID
	invalidInitMethodInvokeTx := &crypto.Transaction{
		Sender:   &sender,
		Receiver: crypto.EmptyAddress,
		Payload:  invalidInitMethodInvokePayload,
		GasLimit: 0,
		GasPrice: 0,
	}

	deployInvalidRLPEncodingContractTx := &crypto.Transaction{
		Sender: &sender,
		Payload: &crypto.TxPayload{
			Contract: rlp.EmptyString,
		},
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 0,
	}

	corruptedContractPayloadWithInit, _ := util.BuildDeployTxPayload("../test/testdata/corrupted.wasm", "../test/testdata/liquid-token-abi.json", "init", []string{"123"})
	deployCorruptedContractWithInitTx := &crypto.Transaction{
		Sender:   &sender2,
		Payload:  corruptedContractPayloadWithInit,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 0,
	}

	corruptedContractPayload, _ := util.BuildDeployTxPayload("../test/testdata/corrupted.wasm", "../test/testdata/liquid-token-abi.json", "", []string{})
	deployCorruptedContractTx := &crypto.Transaction{
		Sender:   &sender2,
		Payload:  corruptedContractPayload,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 0,
	}
	corruptedContractAddress := crypto.NewDeploymentAddress(crypto.AddressFromPubKey(sender2.PublicKey), 0)

	igniteErrorPayload, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1"})
	igniteErrorTx := &crypto.Transaction{
		Sender:   &sender,
		Receiver: corruptedContractAddress,
		Payload:  igniteErrorPayload,
		GasLimit: 0,
		GasPrice: 0,
	}

	// Setup falsy tx to trigger reverse
	sender3 := crypto.TxSender{Nonce: uint64(2)}
	invalidInvokePayload, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	nonExistedPublicKey, _ := hex.DecodeString("1234567812345678")
	invalidContractAddress := crypto.AddressFromPubKey(nonExistedPublicKey)
	invalidInvokeTx := &crypto.Transaction{
		Sender:   &sender3,
		Receiver: invalidContractAddress,
		Payload:  invalidInvokePayload,
		GasLimit: 0,
		GasPrice: 0,
	}

	type args struct {
		app        *App
		tx         *crypto.Transaction
		gasStation gas.Station
	}
	tests := []struct {
		name       string
		args       args
		result     uint64
		code       crypto.ReceiptCode
		events     []*crypto.Event
		gasUsed    uint64
		wantErr    bool
		wantErrObj error
	}{{
		name:       "out of gas",
		args:       args{tr.app, deployTx, gas.NewLiquidStation(tr.app, crypto.Address{})},
		result:     0,
		code:       crypto.ReceiptCodeOutOfGas,
		events:     nil,
		gasUsed:    11100,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "valid deploy tx",
		args:       args{tr.app, deployTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeOK,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "deploy invalid rlp encoding contract",
		args:       args{tr.app, deployInvalidRLPEncodingContractTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeIgniteError,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    true,
		wantErrObj: errors.New("rlp: expected input list for struct { Header []uint8; Code []uint8 }"),
	}, {
		name:       "deploy corrupted contract with init, reverse",
		args:       args{tr.app, deployCorruptedContractWithInitTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeIgniteError,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "valid deploy corrupted contract tx",
		args:       args{tr.app, deployCorruptedContractTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeOK,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "ignite error",
		args:       args{tr.app, igniteErrorTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeIgniteError,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:   "valid deploy init contract tx",
		args:   args{tr.app, deployWithInitTx, gas.NewFreeStation(tr.app)},
		result: 0,
		code:   crypto.ReceiptCodeOK,
		events: []*crypto.Event{{
			Contract: contractAddress,
			Args:     mintEventData,
		}},
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:   "valid invoke tx",
		args:   args{tr.app, invokeTx, gas.NewFreeStation(tr.app)},
		result: 0,
		code:   crypto.ReceiptCodeOK,
		events: []*crypto.Event{{
			Contract: contractAddress,
			Args:     mintEventData,
		}},
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "invalid invoke tx, reverse",
		args:       args{tr.app, invalidInvokeTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeContractNotFound,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "invalid method",
		args:       args{tr.app, invalidMethodInvokeTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeMethodNotFound,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    false,
		wantErrObj: nil,
	}, {
		name:       "invalid method while init",
		args:       args{tr.app, invalidInitMethodInvokeTx, gas.NewFreeStation(tr.app)},
		result:     0,
		code:       crypto.ReceiptCodeMethodNotFound,
		events:     make([]*crypto.Event, 0),
		gasUsed:    0,
		wantErr:    true,
		wantErrObj: errors.New("function with methodID [68 214 68 31] not found"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.app.SetGasStation(tt.args.gasStation)
			receipt, err := tr.app.applyTransaction(tt.args.tx)
			if tt.wantErr && (err == nil) {
				t.Errorf("%s: applyTx() error = %v, wantErr %v", tt.name, err, tt.wantErrObj.Error())
			}
			if tt.wantErr && (err != nil) {
				if tt.wantErrObj.Error() != err.Error() {
					t.Errorf("%s: applyTx() error = %v, wantErr %v", tt.name, err, tt.wantErrObj.Error())
				}
				return
			}

			if receipt.Result != tt.result {
				t.Errorf("%s: applyTx() result = %v, want %v", tt.name, receipt.Result, tt.result)
			}

			if receipt.Code != tt.code {
				t.Errorf("%s: applyTx() receipt.Code = %v, want %v", tt.name, receipt.Code, tt.code)
			}

			if len(receipt.Events) == len(tt.events) {
				for i := range receipt.Events {
					if receipt.Events[i].Contract != tt.events[i].Contract {
						t.Errorf("%s: applyTx() event.contract = %s, want %s", tt.name, receipt.Events[i].Contract.String(), tt.events[i].Contract.String())
					}

					if !bytes.Equal(receipt.Events[i].Args, tt.events[i].Args) {
						t.Errorf("%s: applyTx() event.Params = %v, want %v", tt.name, receipt.Events[i].Args, tt.events[i].Args)
					}
				}
			} else {
				t.Errorf("%s: applyTx() events count = %v, want %v", tt.name, len(receipt.Events), len(tt.events))
			}

			if uint64(receipt.GasUsed) != tt.gasUsed {
				t.Errorf("%s: applyTx() gasUsed = %v, want %v", tt.name, receipt.GasUsed, tt.gasUsed)
			}
		})
	}
}
