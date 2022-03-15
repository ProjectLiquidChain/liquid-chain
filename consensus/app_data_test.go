package consensus

import (
	"crypto/ed25519"

	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
)

func (tr TestResource) getSenderWithNonce(nonce int) (crypto.TxSender, ed25519.PrivateKey) {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(nonce),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	return sender, privateKey
}

func (tr TestResource) getDeployTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	data, err := util.BuildDeployTxPayload("../test/testdata/liquid-token.wasm", "../test/testdata/liquid-token-abi.json", "", []string{})
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

func (tr TestResource) getInvokeTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
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

func (tr TestResource) getInvalidMaxSizeTx(nonce int) *crypto.Transaction {
	sender, _ := tr.getSenderWithNonce(nonce)
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

func (tr TestResource) getInvalidSignatureTx(nonce int) *crypto.Transaction {
	sender, _ := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
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

func (tr TestResource) getInvalidNonceTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
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

func (tr TestResource) getInvalidGasPriceTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
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

func (tr TestResource) getInvokeNilContractTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 123),
		GasLimit: 0,
		GasPrice: 1,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (tr TestResource) getInvokeNonContractTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: senderAddress,
		GasLimit: 0,
		GasPrice: 1,
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (tr TestResource) getInvalidSerializedTx(nonce int) *crypto.Transaction {
	sender, privateKey := tr.getSenderWithNonce(nonce)
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxPayload("../test/testdata/liquid-token-abi.json", "mint", []string{"1000"})
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
