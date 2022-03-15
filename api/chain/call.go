package chain

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/engine"
	"github.com/QuoineFinancial/liquid-chain/gas"
)

// CallParams is params to execute Call
type CallParams struct {
	Height  *uint64  `json:"height"`
	Address string   `json:"address"`
	Method  string   `json:"method"`
	Args    []string `json:"args"`
}

// CallResult is result of Call
type CallResult struct {
	Result string             `json:"result"`
	Code   crypto.ReceiptCode `json:"code"`
	Events []*call            `json:"events"`
}

// Call to execute function without tx creation in blockchain
func (service *Service) Call(r *http.Request, params *CallParams, result *CallResult) error {
	if params.Height == nil {
		service.syncLatestState()
	} else {
		service.syncStateAt(*params.Height)
	}

	address, err := crypto.AddressFromString(params.Address)
	if err != nil {
		return err
	}

	contractAccount, err := service.state.GetAccount(address)
	if err != nil {
		return err
	}

	if contractAccount == nil {
		return errors.New("contract with given address is missing")
	}

	senderAddress := crypto.EmptyAddress
	var app gas.App
	station := gas.NewFreeStation(app)
	execEngine := engine.NewEngine(service.state, contractAccount, senderAddress, station.GetPolicy(), 0)

	contract, err := contractAccount.GetContract()
	if err != nil {
		return err
	}

	function, err := contract.Header.GetFunction(params.Method)
	if err != nil {
		return err
	}

	args, err := abi.EncodeFromString(function.Parameters, params.Args)
	if err != nil {
		return err
	}

	igniteResult, err := execEngine.Ignite(params.Method, args)
	if err != nil {
		return err
	}

	result.Result = fmt.Sprintf("%x", igniteResult)
	result.Code = crypto.ReceiptCodeOK

	parsedEvents := []*call{}
	for _, event := range execEngine.GetEvents() {
		parsedEvent, err := service.parseEvent(event.ID, event.Args, event.Contract)
		if err != nil {
			return nil
		}
		parsedEvents = append(parsedEvents, parsedEvent)
	}
	result.Events = parsedEvents
	return nil
}
