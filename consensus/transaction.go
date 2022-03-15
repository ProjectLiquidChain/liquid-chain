package consensus

import (
	"fmt"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

func (app *App) validateTx(tx *crypto.Transaction) error {
	if tx.Version != 1 {
		return fmt.Errorf("tx version %d not supported", tx.Version)
	}

	nonce := uint64(0)
	address := crypto.AddressFromPubKey(tx.Sender.PublicKey)
	account, err := app.State.LoadAccount(address)
	if err != nil {
		return err
	}
	if account != nil {
		nonce = account.Nonce
	}

	// Validate tx nonce
	if tx.Sender.Nonce != nonce {
		return fmt.Errorf("Invalid nonce. Expected %v, got %v", nonce, tx.Sender.Nonce)
	}

	// Validate tx signature
	signingHash := crypto.GetSigHash(tx)
	if valid := crypto.VerifySignature(tx.Sender.PublicKey, signingHash.Bytes(), tx.Signature); !valid {
		return fmt.Errorf("Invalid signature")
	}

	if tx.Payload.ID != (crypto.MethodID{}) {
		var contract *abi.Contract
		if tx.Receiver != crypto.EmptyAddress {
			account, err := app.State.LoadAccount(tx.Receiver)
			if err != nil {
				return err
			}
			if account == nil {
				return fmt.Errorf("Invoke nil contract")
			}
			if !account.IsContract() {
				return fmt.Errorf("Invoke a non-contract account")
			}
			contract, err = account.GetContract()
			if err != nil {
				return fmt.Errorf("Contract is missing, database might be corrupted")
			}
		} else {
			contract, err = abi.DecodeContract(tx.Payload.Contract)
			if err != nil {
				return err
			}
		}

		function, err := contract.Header.GetFunctionByMethodID(tx.Payload.ID)
		if err != nil {
			return err
		}

		_, err = abi.DecodeToBytes(function.Parameters, tx.Payload.Args)
		if err != nil {
			return err
		}
	}

	// Validate gas limit
	fee := uint64(tx.GasLimit) * uint64(tx.GasPrice)
	if !app.gasStation.Sufficient(address, fee) {
		return fmt.Errorf("Insufficient fee")
	}

	// Validate gas price
	if !app.gasStation.CheckGasPrice(tx.GasPrice) {
		return fmt.Errorf("Invalid gas price")
	}

	return nil
}
