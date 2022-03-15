package crypto

import (
	"crypto/ed25519"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"golang.org/x/crypto/blake2b"
)

// TxSender is sender of transaction
type TxSender struct {
	PublicKey ed25519.PublicKey `json:"publicKey"`
	Nonce     uint64            `json:"nonce"`
}

// TxPayload contains data to interact with smart contract
type TxPayload struct {
	ID       MethodID `json:"signature"`
	Args     []byte   `json:"args"`
	Contract []byte   `json:"contract"`
}

// Transaction is transaction of liquid-chain
type Transaction struct {
	Version   uint16     `json:"version"`
	Sender    *TxSender  `json:"sender"`
	Receiver  Address    `json:"receiver"`
	Payload   *TxPayload `json:"payload"`
	GasPrice  uint32     `json:"gasPrice"`
	GasLimit  uint32     `json:"gasLimit"`
	Signature []byte     `json:"signature"`
}

// Encode returns bytes representation of transaction
func (tx Transaction) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

// DecodeTransaction returns Transaction from bytes representation
func DecodeTransaction(raw []byte) (*Transaction, error) {
	var tx Transaction
	if err := rlp.DecodeBytes(raw, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

// Hash returns hash for storing transaction
func (tx Transaction) Hash() common.Hash {
	hash, _ := tx.Encode()
	return blake2b.Sum256(hash)
}
