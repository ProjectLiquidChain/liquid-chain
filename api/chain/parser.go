package chain

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

func parseParam(param *abi.Parameter, value []byte) (string, error) {
	if param.IsArray {
		return base64.StdEncoding.EncodeToString(value), nil
	}
	switch param.Type {
	case abi.Address:
		address, err := crypto.AddressFromBytes(value)
		if err != nil {
			return "", err
		}
		return address.String(), nil
	case abi.Uint8:
		return fmt.Sprintf("%d", uint8(value[0])), nil
	case abi.Uint16:
		return fmt.Sprintf("%d", binary.LittleEndian.Uint16(value)), nil
	case abi.Uint32:
		return fmt.Sprintf("%d", binary.LittleEndian.Uint32(value)), nil
	case abi.Uint64:
		return fmt.Sprintf("%d", binary.LittleEndian.Uint64(value)), nil
	case abi.Int8:
		return fmt.Sprintf("%d", int8(value[0])), nil
	case abi.Int16:
		return fmt.Sprintf("%d", int16(binary.LittleEndian.Uint16(value))), nil
	case abi.Int32:
		return fmt.Sprintf("%d", int32(binary.LittleEndian.Uint32(value))), nil
	case abi.Int64:
		return fmt.Sprintf("%d", int64(binary.LittleEndian.Uint32(value))), nil
	case abi.Float32:
		return fmt.Sprintf("%f", math.Float32frombits(binary.LittleEndian.Uint32(value))), nil
	case abi.Float64:
		return fmt.Sprintf("%f", math.Float64frombits(binary.LittleEndian.Uint64(value))), nil
	}

	return "", errors.New("unsupported type")
}

func parseFunction(methodID crypto.MethodID, args []byte, contract *abi.Contract) (*call, error) {
	if contract == nil {
		return nil, nil
	}

	function := contract.Header.Functions[methodID]
	parsedArgs, err := abi.DecodeToBytes(function.Parameters, args)
	if err != nil {
		return nil, err
	}

	call := call{
		Name: function.Name,
	}
	for i, arg := range parsedArgs {
		param := function.Parameters[i]
		value, err := parseParam(param, arg)
		if err != nil {
			return nil, err
		}
		typeName := param.Type.String()
		if param.IsArray {
			typeName = typeName + "[]"
		}
		call.Args = append(call.Args, argument{
			Type:  typeName,
			Name:  param.Name,
			Value: value,
		})
	}
	return &call, nil
}

func (service *Service) parseEvent(methodID crypto.MethodID, args []byte, address crypto.Address) (*call, error) {
	account, err := service.state.GetAccount(address)
	if err != nil {
		return nil, err
	}

	contract, err := account.GetContract()
	if err != nil {
		return nil, err
	}

	event := contract.Header.Events[methodID]

	parsedArgs, err := abi.DecodeToBytes(event.Parameters, args)
	if err != nil {
		return nil, err
	}

	call := call{
		Contract: address.String(),
		Name:     event.Name,
	}
	for i, arg := range parsedArgs {
		param := event.Parameters[i]
		value, err := parseParam(param, arg)
		if err != nil {
			return nil, err
		}
		typeName := param.Type.String()
		if param.IsArray {
			typeName = typeName + "[]"
		}
		call.Args = append(call.Args, argument{
			Type:  typeName,
			Name:  param.Name,
			Value: value,
		})
	}
	return &call, nil
}

func (service *Service) parseTransaction(tx *crypto.Transaction, blockHeight uint64) (*transaction, error) {
	parsedTx := transaction{
		BlockHeight: blockHeight,
		Hash:        tx.Hash(),
		Version:     tx.Version,
		Sender:      crypto.AddressFromPubKey(tx.Sender.PublicKey),
		Nonce:       tx.Sender.Nonce,
		Receiver:    tx.Receiver,
		GasPrice:    tx.GasPrice,
		GasLimit:    tx.GasLimit,
		Signature:   tx.Signature,
	}

	var contract *abi.Contract
	if tx.Receiver != crypto.EmptyAddress {
		parsedTx.Type = transactionTypeInvoke
		account, err := service.state.GetAccount(tx.Receiver)
		if err != nil {
			return nil, err
		}
		c, err := account.GetContract()
		if err != nil {
			return nil, err
		}
		contract = c
	} else {
		parsedTx.Type = transactionTypeDeploy
		parsedTx.Receiver = crypto.NewDeploymentAddress(
			crypto.AddressFromPubKey(tx.Sender.PublicKey),
			tx.Sender.Nonce,
		)
		c, err := abi.DecodeContract(tx.Payload.Contract)
		if err == nil {
			contract = c
		}
	}

	if tx.Payload.ID != (crypto.MethodID{}) {
		parsedPayload, err := parseFunction(tx.Payload.ID, tx.Payload.Args, contract)
		if err != nil {
			return nil, err
		}
		if parsedPayload != nil {
			parsedTx.Payload = *parsedPayload
		}
	}

	return &parsedTx, nil
}

func (service *Service) parseReceipt(r *crypto.Receipt) (*receipt, error) {
	parsedReceipt := receipt{
		Index:       r.Index,
		Transaction: r.Transaction,
		Result:      fmt.Sprintf("%x", r.Result),
		Code:        r.Code,
		GasUsed:     r.GasUsed,
		Events:      make([]call, 0),
		PostState:   r.PostState,
	}
	for _, event := range r.Events {
		parsedEvent, err := service.parseEvent(event.ID, event.Args, event.Contract)
		if err != nil {
			return nil, err
		}
		parsedReceipt.Events = append(parsedReceipt.Events, *parsedEvent)
	}
	return &parsedReceipt, nil
}

func (service *Service) parseBlock(rawBlock *crypto.Block) (*block, error) {
	parsedBlock := block{
		Hash:            rawBlock.Hash(),
		Height:          rawBlock.Height,
		Time:            rawBlock.Time,
		Parent:          rawBlock.Parent,
		StateRoot:       rawBlock.StateRoot,
		TransactionRoot: rawBlock.TransactionRoot,
		ReceiptRoot:     rawBlock.ReceiptRoot,
		Transactions:    []transaction{},
		Receipts:        []receipt{},
	}

	for _, tx := range rawBlock.Transactions() {
		parsedTx, err := service.parseTransaction(tx, rawBlock.Height)
		if err != nil {
			return nil, err
		}
		parsedBlock.Transactions = append(parsedBlock.Transactions, *parsedTx)
	}

	txHashToReceipt := make(map[common.Hash]*receipt)
	for _, receipt := range rawBlock.Receipts() {
		parsedReceipt, err := service.parseReceipt(receipt)
		txHashToReceipt[parsedReceipt.Transaction] = parsedReceipt
		if err != nil {
			return nil, err
		}
		parsedBlock.Receipts = append(parsedBlock.Receipts, *parsedReceipt)
	}

	sort.Slice(parsedBlock.Receipts, func(i, j int) bool {
		return parsedBlock.Receipts[i].Index < parsedBlock.Receipts[j].Index
	})

	sort.Slice(parsedBlock.Transactions, func(i, j int) bool {
		return txHashToReceipt[parsedBlock.Transactions[i].Hash].Index < txHashToReceipt[parsedBlock.Transactions[j].Hash].Index
	})

	return &parsedBlock, nil
}
