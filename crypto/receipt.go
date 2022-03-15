package crypto

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"golang.org/x/crypto/blake2b"
)

// Event is emitted while executing transactions
type Event struct {
	ID       MethodID `json:"id"`
	Args     []byte   `json:"args"`
	Contract Address  `json:"contract"`
}

// Receipt reflects corresponding Transaction execution result
type Receipt struct {
	Transaction common.Hash
	Index       uint32      `json:"index"`
	Result      uint64      `json:"result"`
	GasUsed     uint32      `json:"gasUsed"`
	Code        ReceiptCode `json:"code"`
	Events      []*Event    `json:"events"`
	PostState   common.Hash
}

// Encode returns bytes representation of receipt
func (receipt Receipt) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(receipt)
}

// DecodeReceipt returns Receipt from bytes representation
func DecodeReceipt(raw []byte) (*Receipt, error) {
	var receipt Receipt
	if err := rlp.DecodeBytes(raw, &receipt); err != nil {
		return nil, err
	}
	return &receipt, nil
}

// Hash returns hash for storing receipt
func (receipt Receipt) Hash() common.Hash {
	hash, _ := receipt.Encode()
	return blake2b.Sum256(hash)
}
