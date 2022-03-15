package storage

import (
	"encoding/binary"
	"errors"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
)

var (
	// ErrTransactionNotFound used when tx hash not found in meta
	ErrTransactionNotFound = errors.New("transaction not found")
)

// MetaStorage is storage of indexes
type MetaStorage struct {
	db.Database
}

// NewMetaStorage returns new instance of IndexStorage
func NewMetaStorage(db db.Database) *MetaStorage {
	return &MetaStorage{db}
}

// StoreBlockMetas extracts all indexes and store it
func (ms *MetaStorage) StoreBlockMetas(block *crypto.Block) error {
	ms.Put(
		ms.encodeBlockHeightToBlockHashKey(block.Height),
		block.Hash().Bytes(),
	)

	blockHeightByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(blockHeightByte, block.Height)
	for _, tx := range block.Transactions() {
		ms.Put(
			ms.encodeTxHashToBlockHeightKey(tx.Hash()),
			blockHeightByte,
		)
	}

	for _, receipt := range block.Receipts() {
		ms.Put(
			ms.encodeTxHashToReceiptHashKey(receipt.Transaction),
			receipt.Hash().Bytes(),
		)
	}

	if block.Height > ms.LatestBlockHeight() {
		ms.Put(
			ms.encodeLatestBlockHeightKey(),
			blockHeightByte,
		)
	}

	return nil
}

// LatestBlockHeight retrieves latest block height
func (ms *MetaStorage) LatestBlockHeight() uint64 {
	blockHeightByte := ms.Get(ms.encodeLatestBlockHeightKey())
	if len(blockHeightByte) == 0 {
		return crypto.GenesisBlock.Height
	}
	return binary.LittleEndian.Uint64(blockHeightByte)
}

// BlockHeightToBlockHash retrieves block hash by its height
func (ms *MetaStorage) BlockHeightToBlockHash(height uint64) common.Hash {
	hash := ms.Get(ms.encodeBlockHeightToBlockHashKey(height))
	return common.BytesToHash(hash)
}

// TxHashToBlockHeight retrieves height of block which contains tx
func (ms *MetaStorage) TxHashToBlockHeight(txHash common.Hash) (uint64, error) {
	blockHeightByte := ms.Get(ms.encodeTxHashToBlockHeightKey(txHash))
	if len(blockHeightByte) == 0 {
		return 0, ErrTransactionNotFound
	}
	return binary.LittleEndian.Uint64(blockHeightByte), nil
}

// TxHashToReceiptHash retrieves height of block which contains tx
func (ms *MetaStorage) TxHashToReceiptHash(txHash common.Hash) common.Hash {
	receiptHashBytes := ms.Get(ms.encodeTxHashToReceiptHashKey(txHash))
	return common.BytesToHash(receiptHashBytes)
}
