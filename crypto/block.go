package crypto

import (
	"time"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"golang.org/x/crypto/blake2b"
)

// GenesisBlock is the first block of liquid chain
var GenesisBlock = Block{
	Height:          0,
	Time:            0,
	Parent:          common.EmptyHash,
	StateRoot:       common.EmptyHash,
	TransactionRoot: common.EmptyHash,
}

// Block is unit of Liquid chain
type Block struct {
	hash         common.Hash
	transactions []*Transaction
	receipts     []*Receipt
	txTrie       *trie.Trie
	receiptTrie  *trie.Trie

	Height          uint64      `json:"height"`
	Time            uint64      `json:"time"`
	Parent          common.Hash `json:"parent"`
	StateRoot       common.Hash `json:"stateRoot"`
	TransactionRoot common.Hash `json:"transactionRoot"`
	ReceiptRoot     common.Hash `json:"receiptRoot"`
}

// Transactions returns transactions of block
func (block *Block) Transactions() []*Transaction {
	return block.transactions
}

// Receipts returns receipts of block
func (block *Block) Receipts() []*Receipt {
	return block.receipts
}

// AddTransactions adds transactions to block
func (block *Block) AddTransactions(txs ...*Transaction) {
	block.transactions = append(block.transactions, txs...)
}

// AddReceipts adds receipts to block
func (block *Block) AddReceipts(receipts ...*Receipt) {
	block.receipts = append(block.receipts, receipts...)
}

// SetStateRoot sets StateRoot of block
func (block *Block) SetStateRoot(hash common.Hash) {
	block.StateRoot.SetBytes(hash.Bytes())
}

// SetTransactionRoot sets TransactionRoot of block
func (block *Block) SetTransactionRoot(hash common.Hash) {
	block.TransactionRoot.SetBytes(hash.Bytes())
}

// SetReceiptRoot sets ReceiptRoot of block
func (block *Block) SetReceiptRoot(hash common.Hash) {
	block.ReceiptRoot.SetBytes(hash.Bytes())
}

// Hash returns blake2b hash of rlp encoding of block
func (block *Block) Hash() common.Hash {
	if block.hash == common.EmptyHash {
		encoded, _ := rlp.EncodeToBytes(block)
		blockChecksum := blake2b.Sum256(encoded)
		blockHash := common.BytesToHash(blockChecksum[:])
		block.hash = blockHash
	}
	return block.hash
}

// NewEmptyBlock creates empty block
func NewEmptyBlock(parent common.Hash, height uint64, blockTime time.Time) *Block {
	return &Block{
		Parent: parent,
		Height: height,
		Time:   uint64(blockTime.UTC().Unix()),
	}
}

// Encode returns bytes array of block
func (block *Block) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(block)
}

// DecodeBlock returns block from encoded byte array
func DecodeBlock(rawBlock []byte) (*Block, error) {
	var block Block
	if err := rlp.DecodeBytes(rawBlock, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// MustDecodeBlock returns block from encoded byte array. It will panic in case of error
func MustDecodeBlock(rawBlock []byte) *Block {
	block, err := DecodeBlock(rawBlock)
	if err != nil {
		panic(err)
	}
	return block
}
