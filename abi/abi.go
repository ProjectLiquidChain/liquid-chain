package abi

import (
	"fmt"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
)

// Encode return []byte from an inputted params and values pair
func Encode(params []*Parameter, values []interface{}) ([]byte, error) {
	if len(params) != len(values) {
		return []byte{0}, fmt.Errorf("Parameter count mismatch, expecting: %d, got: %d", len(params), len(values))
	}

	var rlpCompatibleArgs []interface{}

	for index, param := range params {
		if param.IsArray {
			arrayArg, err := param.Type.NewArrayArgument(values[index])
			if err != nil {
				return nil, err
			}
			rlpCompatibleArgs = append(rlpCompatibleArgs, arrayArg)
		} else {
			argument, err := param.Type.NewArgument(values[index])
			if err != nil {
				return nil, err
			}
			rlpCompatibleArgs = append(rlpCompatibleArgs, argument)
		}
	}
	result, err := rlp.EncodeToBytes(rlpCompatibleArgs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeToBytes returns uint64 array compatible with VM
func DecodeToBytes(params []*Parameter, bytes []byte) ([][]byte, error) {
	var decoded [][]byte
	err := rlp.DecodeBytes(bytes, &decoded)
	if err != nil {
		return nil, err
	}
	if len(params) != len(decoded) {
		return nil, fmt.Errorf("Argument count mismatch, expecting: %d, got: %d", len(params), len(decoded))
	}

	return decoded, nil
}

// EncodeFromBytes encodes arguments in byte format - an inverse of DecodeToBytes
func EncodeFromBytes(params []*Parameter, bytes [][]byte) ([]byte, error) {
	if len(params) != len(bytes) {
		return nil, fmt.Errorf("Argument count mismatch, expecting: %d, got: %d", len(params), len(bytes))
	}

	result, err := rlp.EncodeToBytes(bytes)
	if err != nil {
		return nil, err
	}

	return result, nil
}
