package chain

import (
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/common"
)

// LatestBlockParams is params for latest block request
type LatestBlockParams struct{}

// BlockByHeightParams contains query height
type BlockByHeightParams struct {
	Height uint64 `json:"height"`
}

// BlockResult is response of GetBlock
type BlockResult struct {
	Block *block `json:"block"`
}

// GetLatestBlock return the block by height
func (service *Service) GetLatestBlock(r *http.Request, params *LatestBlockParams, result *BlockResult) error {
	service.syncLatestState()
	blockHash := service.meta.BlockHeightToBlockHash(service.meta.LatestBlockHeight())
	block, err := service.block.GetBlock(blockHash)
	if err != nil {
		return err
	}

	txs, err := service.block.GetBlockTransactions(block)
	if err != nil {
		return err
	}
	block.AddTransactions(txs...)

	receipts, err := service.block.GetBlockReceipts(block)
	if err != nil {
		return err
	}
	block.AddReceipts(receipts...)

	parsedBlock, err := service.parseBlock(block)
	if err != nil {
		return err
	}
	result.Block = parsedBlock
	return nil
}

// GetBlockByHeight return block by its height
func (service *Service) GetBlockByHeight(r *http.Request, params *BlockByHeightParams, result *BlockResult) error {
	service.syncLatestState()

	blockHash := service.meta.BlockHeightToBlockHash(params.Height)
	if blockHash == common.EmptyHash {
		return fmt.Errorf("block %d not found", params.Height)
	}

	block, err := service.block.GetBlock(blockHash)
	if err != nil {
		return err
	}

	txs, err := service.block.GetBlockTransactions(block)
	if err != nil {
		return err
	}
	block.AddTransactions(txs...)

	receipts, err := service.block.GetBlockReceipts(block)
	if err != nil {
		return err
	}
	block.AddReceipts(receipts...)

	parsedBlock, err := service.parseBlock(block)
	if err != nil {
		return err
	}
	result.Block = parsedBlock
	return nil
}
