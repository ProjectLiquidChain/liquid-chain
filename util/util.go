package util

import (
	"io/ioutil"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// BuildInvokeTxPayload builds data for invoke transaction
func BuildInvokeTxPayload(headerPath string, methodName string, params []string) (*crypto.TxPayload, error) {
	header, err := abi.LoadHeaderFromFile(headerPath)
	if err != nil {
		return nil, err
	}

	function, err := header.GetFunction(methodName)
	if err != nil {
		return nil, err
	}

	encodedArgs, err := abi.EncodeFromString(function.Parameters, params)
	if err != nil {
		return nil, err
	}

	return &crypto.TxPayload{
		ID:   crypto.GetMethodID(methodName),
		Args: encodedArgs,
	}, nil
}

// BuildDeployTxPayload builds data for deploy transaction
func BuildDeployTxPayload(codePath string, headerPath string, initFuncName string, params []string) (*crypto.TxPayload, error) {
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		return nil, err
	}

	encodedHeader, err := abi.EncodeHeaderToBytes(headerPath)
	if err != nil {
		return nil, err
	}

	header, err := abi.DecodeHeader(encodedHeader)
	if err != nil {
		return nil, err
	}

	contractCode, err := rlp.EncodeToBytes(&abi.Contract{Header: header, Code: code})
	if err != nil {
		return nil, err
	}

	payload := crypto.TxPayload{
		Contract: contractCode,
	}

	function, err := header.GetFunction(initFuncName)
	if err == nil {
		encodedArgs, err := abi.EncodeFromString(function.Parameters, params)
		if err != nil {
			return nil, err
		}
		payload.ID = crypto.GetMethodID(initFuncName)
		payload.Args = encodedArgs
	}

	return &payload, nil
}
