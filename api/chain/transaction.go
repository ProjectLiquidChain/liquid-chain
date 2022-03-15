package chain

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/trie"
)

// GetTransactionParams contains query height
type GetTransactionParams struct {
	Hash string `json:"hash"`
}

// GetTransactionResult is response of GetTransactions
type GetTransactionResult struct {
	Transaction *transaction `json:"transaction"`
	Receipt     *receipt     `json:"receipt"`
}

// GetTransaction returns txs of given block
func (service *Service) GetTransaction(r *http.Request, params *GetTransactionParams, result *GetTransactionResult) error {
	service.syncLatestState()
	if _, err := hex.DecodeString(params.Hash); err != nil {
		return err
	}

	// Get block
	txHash := common.HexToHash(params.Hash)
	height, err := service.meta.TxHashToBlockHeight(txHash)
	if err != nil {
		return err
	}
	blockHash := service.meta.BlockHeightToBlockHash(height)
	if blockHash == common.EmptyHash {
		return fmt.Errorf("block %d not found", height)
	}

	block, err := service.block.GetBlock(blockHash)
	if err != nil {
		return err
	}

	// Get tx
	txTrie, err := trie.New(block.TransactionRoot, service.block)
	if err != nil {
		return err
	}
	rawTx, err := txTrie.Get(txHash.Bytes())
	if err != nil {
		return err
	}
	tx, err := crypto.DecodeTransaction(rawTx)
	if err != nil {
		return err
	}
	parsedTx, err := service.parseTransaction(tx, height)
	if err != nil {
		return err
	}

	result.Transaction = parsedTx

	// Get receipt
	receiptHash := service.meta.TxHashToReceiptHash(txHash)
	receiptTrie, err := trie.New(block.ReceiptRoot, service.block)
	if err != nil {
		return err
	}
	receiptBytes, err := receiptTrie.Get(receiptHash.Bytes())
	if err != nil {
		return err
	}
	receipt, err := crypto.DecodeReceipt(receiptBytes)
	if err != nil {
		return err
	}
	parsedReceipt, err := service.parseReceipt(receipt)
	if err != nil {
		return err
	}
	result.Receipt = parsedReceipt

	return nil
}
