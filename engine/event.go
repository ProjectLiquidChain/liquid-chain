package engine

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/vertexdlt/vertexvm/vm"
)

func (engine *Engine) handleEmitEvent(eventHeader *abi.Event, vm *vm.VM, args ...uint64) (uint64, error) {
	var memBytes [][]byte
	for i, param := range eventHeader.Parameters {
		switch param.Type {
		case abi.Address:
			paramPtr := int(uint32(args[i]))
			size := param.Type.GetMemorySize()
			memValue, err := readAt(vm, paramPtr, size)
			if err != nil {
				return 0, err
			}
			if _, err := crypto.AddressFromBytes(memValue); err != nil {
				return 0, err
			}
			memBytes = append(memBytes, memValue)
		default:
			if param.IsArray {
				paramPtr := int(uint32(args[i]))
				size, err := engine.ptrArgSizeGet(paramPtr)
				if err != nil {
					return 0, err
				}
				array, err := readAt(vm, paramPtr, size)
				memBytes = append(memBytes, array)
			} else {
				size := abi.Uint64.GetMemorySize()
				value := make([]byte, size)
				binary.LittleEndian.PutUint64(value, args[i])
				memBytes = append(memBytes, value)
			}
		}
	}

	values, err := abi.EncodeFromBytes(eventHeader.Parameters, memBytes)
	if err != nil {
		return 0, err
	}

	cost := engine.gasPolicy.GetCostForEvent(len(values))
	if err := vm.BurnGas(cost); err != nil {
		return 0, err
	}

	engine.pushEvent(&crypto.Event{
		ID:       crypto.GetMethodID(eventHeader.Name),
		Contract: engine.account.GetAddress(),
		Args:     values,
	})
	return 0, nil
}
