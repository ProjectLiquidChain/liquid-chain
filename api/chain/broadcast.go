package chain

import (
	"encoding/base64"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// BroadcastParams is params to broadcast transaction
type BroadcastParams struct {
	RawTransaction string `json:"rawTx"`
}

// BroadcastResult is result of broadcast
type BroadcastResult struct {
	Code            uint32 `json:"code"`
	Log             string `json:"log"`
	TransactionHash string `json:"hash"`
}

// Broadcast delivers transaction to blockchain
func (service *Service) Broadcast(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}

	tx, err := crypto.DecodeTransaction(bytes)
	if err != nil {
		return err
	}

	tmResult, err := service.tmAPI.BroadcastTxSync(bytes)
	if err != nil {
		return err
	}

	result.TransactionHash = tx.Hash().String()
	result.Log = tmResult.Log
	result.Code = tmResult.Code

	return nil
}

// BroadcastAsync broadcast but wont wait for transaction's checkTx result
func (service *Service) BroadcastAsync(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}

	tx, err := crypto.DecodeTransaction(bytes)
	if err != nil {
		return err
	}

	tmResult, err := service.tmAPI.BroadcastTxAsync(bytes)
	if err != nil {
		return err
	}

	result.TransactionHash = tx.Hash().String()
	result.Log = tmResult.Log
	result.Code = tmResult.Code

	return nil
}

// BroadcastCommit broadcast and wait until the transaction is committed in a block or fail to pass checkTx
func (service *Service) BroadcastCommit(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}

	tx, err := crypto.DecodeTransaction(bytes)
	if err != nil {
		return err
	}

	tmResult, err := service.tmAPI.BroadcastTxCommit(bytes)
	if err != nil {
		return err
	}

	result.TransactionHash = tx.Hash().String()
	result.Log = tmResult.CheckTx.Log
	result.Code = tmResult.CheckTx.Code

	return nil
}
