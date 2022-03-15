package consensus

import (
	"bytes"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/engine"
)

// InitFunctionName is default init function name
const InitFunctionName = "init"

var (
	initFunctionID = crypto.GetMethodID(InitFunctionName)
)

func (app *App) applyTransaction(tx *crypto.Transaction) (*crypto.Receipt, error) {
	if tx.Receiver == crypto.EmptyAddress {
		return app.deployContract(tx)
	}
	return app.invokeContract(tx)
}

func (app *App) deployContract(tx *crypto.Transaction) (*crypto.Receipt, error) {
	receipt := crypto.Receipt{
		Transaction: tx.Hash(),
	}

	contractSize := len(tx.Payload.Contract)
	policy := app.gasStation.GetPolicy()
	receipt.GasUsed = uint32(policy.GetCostForContract(contractSize))
	if tx.GasLimit < receipt.GasUsed {
		receipt.Code = crypto.ReceiptCodeOutOfGas
		return &receipt, nil
	}

	contract, err := abi.DecodeContract(tx.Payload.Contract)
	if err != nil {
		return nil, err
	}

	// Create contract account
	senderAddress := crypto.AddressFromPubKey(tx.Sender.PublicKey)
	contractAddress := crypto.NewDeploymentAddress(senderAddress, tx.Sender.Nonce)
	contractAccount, err := app.State.CreateAccount(senderAddress, contractAddress, tx.Payload.Contract)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(tx.Payload.ID[:], initFunctionID[:]) {
		function, err := contract.Header.GetFunctionByMethodID(tx.Payload.ID)
		if err != nil {
			return nil, err
		}
		execEngine := engine.NewEngine(app.State, contractAccount, senderAddress, policy, uint64(tx.GasLimit-receipt.GasUsed))
		result, err := execEngine.Ignite(function.Name, tx.Payload.Args)
		receipt.GasUsed += uint32(execEngine.GetGasUsed())
		if err != nil {
			receipt.Code = crypto.ReceiptCodeIgniteError
			app.State.Revert()
		} else if !app.gasStation.Sufficient(senderAddress, uint64(receipt.GasUsed)*uint64(tx.GasPrice)) {
			receipt.Code = crypto.ReceiptCodeOutOfGas
			receipt.GasUsed = tx.GasLimit
			app.State.Revert()
		} else {
			receipt.Result = result
			receipt.Code = crypto.ReceiptCodeOK
			receipt.Events = append(receipt.Events, execEngine.GetEvents()...)
		}
	}

	// Create account for creator and increase nonce by 1
	if err := app.increaseNonce(senderAddress); err != nil {
		return nil, err
	}

	gasEvents := app.gasStation.Burn(senderAddress, uint64(receipt.GasUsed)*uint64(tx.GasPrice))
	receipt.Events = append(receipt.Events, gasEvents...)
	receipt.PostState = app.State.Hash()
	return &receipt, nil
}

func (app *App) invokeContract(tx *crypto.Transaction) (*crypto.Receipt, error) {
	receipt := crypto.Receipt{
		Transaction: tx.Hash(),
	}

	contractAccount, err := app.State.LoadAccount(tx.Receiver)
	if err != nil {
		panic(err)
	}

	if contractAccount == nil {
		receipt.Code = crypto.ReceiptCodeContractNotFound
		return &receipt, nil
	}

	contract, err := contractAccount.GetContract()
	if err != nil {
		panic(err)
	}
	function, err := contract.Header.GetFunctionByMethodID(tx.Payload.ID)
	if err != nil {
		receipt.Code = crypto.ReceiptCodeMethodNotFound
		return &receipt, nil
	}

	policy := app.gasStation.GetPolicy()
	senderAddress := crypto.AddressFromPubKey(tx.Sender.PublicKey)
	execEngine := engine.NewEngine(app.State, contractAccount, senderAddress, policy, uint64(tx.GasLimit))

	result, err := execEngine.Ignite(function.Name, tx.Payload.Args)
	receipt.GasUsed = uint32(execEngine.GetGasUsed())

	if err != nil {
		receipt.Code = crypto.ReceiptCodeIgniteError
		app.State.Revert()
	} else if !app.gasStation.Sufficient(senderAddress, uint64(receipt.GasUsed)*uint64(tx.GasPrice)) {
		receipt.Code = crypto.ReceiptCodeOutOfGas
		receipt.GasUsed = tx.GasLimit
		app.State.Revert()
	} else {
		receipt.Result = result
		receipt.Events = append(receipt.Events, execEngine.GetEvents()...)
	}

	// Create/get account for creator and increase nonce by 1
	if err := app.increaseNonce(senderAddress); err != nil {
		return nil, err
	}

	gasEvents := app.gasStation.Burn(senderAddress, uint64(receipt.GasUsed)*uint64(tx.GasPrice))
	receipt.Events = append(receipt.Events, gasEvents...)
	receipt.PostState = app.State.Hash()

	return &receipt, nil
}

func (app *App) increaseNonce(address crypto.Address) error {
	account, err := app.State.LoadAccount(address)
	if err != nil {
		return err
	}

	// Make sure account is created
	if account == nil {
		/* sender address is nil on first time submitting tx to blockchain
		therefore we need to create account with nonce 0 for it first */
		account, err = app.State.CreateAccount(address, address, nil)
		if err != nil {
			return err
		}
	}

	account.SetNonce(account.Nonce + 1)
	return nil
}
