package chain

import (
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

type argument struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type call struct {
	Contract string     `json:"contract,omitempty"`
	Name     string     `json:"name,omitempty"`
	Args     []argument `json:"args,omitempty"`
}

type receipt struct {
	Index       uint32             `json:"index"`
	Transaction common.Hash        `json:"transaction"`
	Result      string             `json:"result"`
	GasUsed     uint32             `json:"gasUsed"`
	Code        crypto.ReceiptCode `json:"code"`
	Events      []call             `json:"events"`
	PostState   common.Hash        `json:"postState"`
}

type transactionType string

const (
	transactionTypeDeploy transactionType = "deploy"
	transactionTypeInvoke transactionType = "invoke"
)

type transaction struct {
	Hash        common.Hash     `json:"hash"`
	Type        transactionType `json:"type"`
	BlockHeight uint64          `json:"height"`
	Version     uint16          `json:"version"`
	Sender      crypto.Address  `json:"sender"`
	Nonce       uint64          `json:"nonce"`
	Receiver    crypto.Address  `json:"receiver"`
	Payload     call            `json:"payload"`
	GasPrice    uint32          `json:"gasPrice"`
	GasLimit    uint32          `json:"gasLimit"`
	Signature   []byte          `json:"signature"`
}

type block struct {
	Hash            common.Hash   `json:"hash"`
	Transactions    []transaction `json:"transactions"`
	Receipts        []receipt     `json:"receipts"`
	Height          uint64        `json:"height"`
	Time            uint64        `json:"time"`
	Parent          common.Hash   `json:"parent"`
	StateRoot       common.Hash   `json:"stateRoot"`
	TransactionRoot common.Hash   `json:"transactionRoot"`
	ReceiptRoot     common.Hash   `json:"receiptRoot"`
}
