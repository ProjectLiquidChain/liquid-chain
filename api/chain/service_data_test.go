package chain

import (
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/QuoineFinancial/liquid-chain/consensus"
	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
	"github.com/tendermint/tendermint/abci/types"
)

type testResource struct {
	service *Service
	app     *consensus.App
	dbDir   string
}

func (resource testResource) tearDown() {
	os.RemoveAll(resource.dbDir)
}

func (resource testResource) seed() {
	app := resource.app

	type txRequest struct {
		tx                        *crypto.Transaction
		expectedResponseCheckTx   types.ResponseCheckTx
		expectedResponseDeliverTx types.ResponseDeliverTx
	}

	type round struct {
		height int64
		time   time.Time

		txRequests []txRequest
	}

	rounds := []round{{
		height: 1,
		time:   time.Unix(1, 0),
		txRequests: []txRequest{{
			tx:                        resource.getDeployTx(0, 0),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}, {
			tx:                        resource.getDeployTx(1, 0),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}, {
			tx:                        resource.getDeployTx(2, 0),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}},
	}, {
		height: 2,
		time:   time.Unix(2, 0),
		txRequests: []txRequest{{
			tx:                        resource.getInvokeTx(1, 1),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}, {
			tx:                        resource.getInvokeTx(0, 1),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}, {
			tx:                        resource.getInvokeTx(2, 1),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}},
	}, {
		height: 3,
		time:   time.Unix(3, 0),
		txRequests: []txRequest{{
			tx:                        resource.getDeployEventStringTx(0, 2),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: consensus.ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: consensus.ResponseCodeOK},
		}},
	}, {
		height: 4,
		time:   time.Unix(4, 0),
		txRequests: []txRequest{{
			tx:                      resource.getInvokeNilContractTx(0, 3),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: consensus.ResponseCodeNotOK, Log: "Invoke nil contract"},
		}, {
			tx:                      resource.getInvalidMaxSizeTx(0, 3),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: consensus.ResponseCodeNotOK, Log: fmt.Sprintf("Transaction size exceed %vB", constant.MaxTransactionSize)},
		}, {
			tx:                      resource.getInvalidSignatureTx(0, 3),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: consensus.ResponseCodeNotOK, Log: "Invalid signature"},
		}, {
			tx:                      resource.getInvalidNonceTx(0, 123),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: consensus.ResponseCodeNotOK, Log: "Invalid nonce. Expected 3, got 123"},
		}, {
			tx:                      resource.getInvalidGasPriceTx(0, 3),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: consensus.ResponseCodeNotOK, Log: "Invalid gas price"},
		}, {
			tx:                      resource.getInvokeNonContractTx(0, 3),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: consensus.ResponseCodeNotOK, Log: "Invoke a non-contract account"},
		}},
	}}

	appHash := []byte{}
	for _, round := range rounds {
		app.BeginBlock(types.RequestBeginBlock{
			Header: types.Header{
				Height:  round.height,
				Time:    round.time,
				AppHash: appHash,
			},
		})

		for _, txRequest := range round.txRequests {
			rawTx, _ := txRequest.tx.Encode()
			responseCheckTx := app.CheckTx(types.RequestCheckTx{Tx: rawTx})
			if responseCheckTx.Code == consensus.ResponseCodeOK {
				app.DeliverTx(types.RequestDeliverTx{Tx: rawTx})
			}
		}

		responseCommit := app.Commit()
		appHash = responseCommit.Data
	}
}

func newTestResource() *testResource {
	rand.Seed(time.Now().UTC().UnixNano())
	dbDir := "./tmp" + strconv.Itoa(rand.Intn(10000)) + "/"

	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		panic(err)
	}

	app := consensus.NewApp(filepath.Join(dbDir, "liquid"), "")
	if err := app.State.LoadState(&crypto.GenesisBlock); err != nil {
		panic(err)
	}

	service := NewService(nil, app.Meta, app.State, app.Chain)
	return &testResource{service, app, dbDir}
}

func (testResource) getSenderWithNonce(senderIndex byte, nonce int) (crypto.TxSender, ed25519.PrivateKey) {
	seed := append(make([]byte, 31), senderIndex)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(nonce),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	return sender, privateKey
}

func (resource testResource) getDeployEventStringTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	data, err := util.BuildDeployTxPayload("../../test/testdata/event-string.wasm", "../../test/testdata/event-string-abi.json", "", []string{})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 1,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getDeployTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	data, err := util.BuildDeployTxPayload("../../test/testdata/liquid-token.wasm", "../../test/testdata/liquid-token-abi.json", "", []string{})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 1,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getInvokeTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getInvalidMaxSizeTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, _ := resource.getSenderWithNonce(senderIndex, nonce)
	type maxSizeContart [constant.MaxTransactionSize]byte
	var contract maxSizeContart
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  &crypto.TxPayload{Contract: contract[:]},
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 1,
	}
	return tx
}

func (resource testResource) getInvalidSignatureTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, _ := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
	}
	tx.Signature = []byte{1, 2, 3}
	return tx
}

func (resource testResource) getInvalidNonceTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getInvalidGasPriceTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 0,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getInvokeNilContractTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 123),
		GasLimit: 0,
		GasPrice: 0,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getInvokeNonContractTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: senderAddress,
		GasLimit: 0,
		GasPrice: 0,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (resource testResource) getInvalidSerializedTx(senderIndex byte, nonce int) *crypto.Transaction {
	sender, privateKey := resource.getSenderWithNonce(senderIndex, nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: senderAddress,
		GasLimit: 0,
		GasPrice: 0,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}
