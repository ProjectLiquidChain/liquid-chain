package storage

import (
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/trie"
)

// ChainStorage is storage for block
type ChainStorage struct {
	db.Database
	txTrie       *trie.Trie
	receiptTrie  *trie.Trie
	CurrentBlock *crypto.Block
}

// NewChainStorage returns new instance of IndexStorage
func NewChainStorage(db db.Database) *ChainStorage {
	return &ChainStorage{db, nil, nil, nil}
}

// ComposeBlock compose currentBlock based on parent and proposed time
func (bs *ChainStorage) ComposeBlock(parent *crypto.Block, time time.Time) {
	bs.CurrentBlock = crypto.NewEmptyBlock(parent.Hash(), parent.Height+1, time)

	if txTrie, err := trie.New(common.EmptyHash, bs.Database); err != nil {
		panic(err)
	} else {
		bs.txTrie = txTrie
	}

	if receiptTrie, err := trie.New(common.EmptyHash, bs.Database); err != nil {
		panic(err)
	} else {
		bs.receiptTrie = receiptTrie
	}
}

// Commit puts currentBlock to storage
func (bs *ChainStorage) Commit(stateRoot common.Hash) common.Hash {
	if bs.CurrentBlock == nil {
		panic("ChainStorage.currentBlock is nil")
	}

	// Set state root
	bs.CurrentBlock.SetStateRoot(stateRoot)

	// Commit and set tx root
	txRootHash, err := bs.txTrie.Commit()
	if err != nil {
		panic(err)
	}
	bs.CurrentBlock.SetTransactionRoot(txRootHash)

	receiptRoot, err := bs.receiptTrie.Commit()
	if err != nil {
		panic(err)
	}
	bs.CurrentBlock.SetReceiptRoot(receiptRoot)

	// Store block
	hash := bs.CurrentBlock.Hash()
	rawBlock, err := bs.CurrentBlock.Encode()
	if err != nil {
		panic(err)
	}
	bs.Put(hash.Bytes(), rawBlock)

	return hash
}

// AddTransactionWithReceipt add tx and receipt to currentBlock
func (bs *ChainStorage) AddTransactionWithReceipt(tx *crypto.Transaction, receipt *crypto.Receipt) error {
	if bs.CurrentBlock == nil {
		panic("ChainStorage.currentBlock is nil")
	}

	rawTx, err := tx.Encode()
	if err != nil {
		return err
	}
	bs.txTrie.Update(tx.Hash().Bytes(), rawTx)
	bs.CurrentBlock.AddTransactions(tx)

	receipt.Index = uint32(len(bs.CurrentBlock.Receipts()))
	rawReceipt, err := receipt.Encode()
	if err != nil {
		return err
	}
	bs.receiptTrie.Update(receipt.Hash().Bytes(), rawReceipt)

	bs.CurrentBlock.AddReceipts(receipt)
	return nil
}

// GetBlock retrieves block by its hash
func (bs *ChainStorage) GetBlock(hash common.Hash) (*crypto.Block, error) {
	if hash == common.EmptyHash {
		return &crypto.GenesisBlock, nil
	}
	rawBlock := bs.Get(hash.Bytes())
	return crypto.DecodeBlock(rawBlock)
}

// MustGetBlock retrieves block by its hash, panic if failed
func (bs *ChainStorage) MustGetBlock(hash common.Hash) *crypto.Block {
	block, err := bs.GetBlock(hash)
	if err != nil {
		panic(err)
	}
	return block
}

// GetBlockTransactions returns transactions of given block
func (bs *ChainStorage) GetBlockTransactions(block *crypto.Block) ([]*crypto.Transaction, error) {
	txTrie, err := trie.New(block.TransactionRoot, bs.Database)
	if err != nil {
		return nil, err
	}
	iterator := trie.NewIterator(txTrie.NodeIterator(nil))
	txs := []*crypto.Transaction{}
	for iterator.Next() {
		tx, err := crypto.DecodeTransaction(iterator.Value)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// GetBlockReceipts returns receipts of given block
func (bs *ChainStorage) GetBlockReceipts(block *crypto.Block) ([]*crypto.Receipt, error) {
	receiptTrie, err := trie.New(block.ReceiptRoot, bs.Database)
	if err != nil {
		return nil, err
	}
	iterator := trie.NewIterator(receiptTrie.NodeIterator(nil))
	receipts := []*crypto.Receipt{}
	for iterator.Next() {
		receipt, err := crypto.DecodeReceipt(iterator.Value)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, nil
}
