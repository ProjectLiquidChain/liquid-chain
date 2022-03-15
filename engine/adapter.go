package engine

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/vertexdlt/vertexvm/vm"
	"golang.org/x/crypto/blake2b"
)

const pointerSize = int(4)

func readAt(vm *vm.VM, ptr, size int) ([]byte, error) {
	data := make([]byte, size)
	_, err := vm.MemRead(data, ptr)
	return data, err
}

func (engine *Engine) chainStorageSet(vm *vm.VM, args ...uint64) (uint64, error) {
	keyPtr, keySize := int(args[0]), int(args[1])
	valuePtr, valueSize := int(args[2]), int(args[3])
	// Burn gas before actually execute
	cost := engine.gasPolicy.GetCostForStorage(valueSize)
	err := vm.BurnGas(cost)
	if err != nil {
		return 0, err
	}
	key, err := readAt(vm, keyPtr, keySize)
	if err != nil {
		return 0, err
	}
	value, err := readAt(vm, valuePtr, valueSize)
	if err != nil {
		return 0, err
	}
	err = engine.account.SetStorage(key, value)
	return uint64(len(value)), err
}

func (engine *Engine) chainStorageGet(vm *vm.VM, args ...uint64) (uint64, error) {
	keyPtr, keySize := int(args[0]), int(args[1])
	key, err := readAt(vm, keyPtr, keySize)
	if err != nil {
		return 0, err
	}
	valuePtr := int(uint32(args[2]))
	value, err := engine.account.GetStorage(key)
	if err != nil {
		return 0, err
	}
	byteSize, err := vm.MemWrite(value, valuePtr)
	return uint64(byteSize), err
}

func (engine *Engine) chainStorageSizeGet(vm *vm.VM, args ...uint64) (uint64, error) {
	keyPtr, keySize := int(args[0]), int(args[1])
	key, err := readAt(vm, keyPtr, keySize)
	if err != nil {
		return 0, err
	}
	value, err := engine.account.GetStorage(key)
	return uint64(len(value)), err
}

func (engine *Engine) chainGetCaller(vm *vm.VM, args ...uint64) (uint64, error) {
	_, err := vm.MemWrite(engine.caller[:], int(args[0]))
	return 0, err
}

func (engine *Engine) chainGetCreator(vm *vm.VM, args ...uint64) (uint64, error) {
	creator := engine.account.Creator
	_, err := vm.MemWrite(creator[:], int(args[0]))
	return 0, err
}

func (engine *Engine) chainPtrArgSizeGet(vm *vm.VM, args ...uint64) (uint64, error) {
	size, err := engine.ptrArgSizeGet(int(args[0]))
	return uint64(size), err
}

func (engine *Engine) chainPtrArgSizeSet(vm *vm.VM, args ...uint64) (uint64, error) {
	engine.ptrArgSizeSet(int(args[0]), int(args[1]))
	return 0, nil
}

func (engine *Engine) chainMethodBind(vm *vm.VM, args ...uint64) (uint64, error) {
	contractAddrBytes, err := readAt(vm, int(args[0]), crypto.AddressLength)
	if err != nil {
		return 0, err
	}
	contractAddr, err := crypto.AddressFromBytes(contractAddrBytes)
	if err != nil {
		return 0, err
	}

	invokedMethodBytes, err := readAt(vm, int(args[1]), int(args[2]))
	if err != nil {
		return 0, err
	}

	invokedMethod := string(invokedMethodBytes[:len(invokedMethodBytes)-1])
	aliasMethodBytes, err := readAt(vm, int(args[3]), int(args[4]))
	if err != nil {
		return 0, err
	}
	aliasMethod := string(aliasMethodBytes[:len(aliasMethodBytes)-1])
	engine.methodLookup[aliasMethod] = &foreignMethod{contractAddr, invokedMethod}
	return 0, nil
}

func (engine *Engine) chainBlockHeight(vm *vm.VM, args ...uint64) (uint64, error) {
	return engine.state.GetBlock().Height, nil
}

func (engine *Engine) chainBlockTime(vm *vm.VM, args ...uint64) (uint64, error) {
	return uint64(engine.state.GetBlock().Time), nil
}

func (engine *Engine) chainArgsWrite(vm *vm.VM, args ...uint64) (uint64, error) {
	bufferPtr, valuePtr, valueSize := int(args[0]), int(args[1]), int(args[2])
	bufferSize, _ := engine.ptrArgSizeGet(bufferPtr)
	memorySize := 4
	buf := make([]byte, memorySize)
	binary.LittleEndian.PutUint32(buf, uint32(valuePtr))
	_, err := vm.MemWrite(buf, bufferPtr+bufferSize)
	if err != nil {
		return 0, err
	}
	engine.ptrArgSizeSet(valuePtr, valueSize)
	engine.ptrArgSizeSet(bufferPtr, bufferSize+memorySize)
	return uint64(bufferPtr), nil
}

func (engine *Engine) chainArgsHash(vm *vm.VM, args ...uint64) (uint64, error) {
	bufferPtr, hashPtr := int(args[0]), int(args[1])
	bufferSize, err := engine.ptrArgSizeGet(bufferPtr)
	if err != nil {
		return 0, err
	}
	memorySize := 4
	argCnt := bufferSize / memorySize
	var values [][]byte
	for i := 0; i < argCnt; i++ {
		ptrMem, err := readAt(vm, bufferPtr+i*memorySize, memorySize)
		ptr := int(binary.LittleEndian.Uint32(ptrMem))
		ptrSize, err := engine.ptrArgSizeGet(ptr)
		if err != nil {
			return 0, err
		}
		value, err := readAt(vm, ptr, ptrSize)
		if err != nil {
			return 0, err
		}
		values = append(values, value)
	}
	result, err := rlp.EncodeToBytes(values)
	if err != nil {
		return 0, err
	}
	hash := blake2b.Sum256(result)
	vm.MemWrite(hash[:], hashPtr)
	return 0, nil
}

func (engine *Engine) chainEd25519Verify(vm *vm.VM, args ...uint64) (uint64, error) {
	addressPtr, hasherPtr, signaturePtr := int(args[0]), int(args[1]), int(args[2])
	addressBytes, err := readAt(vm, addressPtr, crypto.AddressLength)
	if err != nil {
		return 0, err
	}
	address, err := crypto.AddressFromBytes(addressBytes)
	if err != nil {
		return 0, err
	}
	hasher, err := readAt(vm, hasherPtr, 32)
	if err != nil {
		return 0, err
	}
	signature, err := readAt(vm, signaturePtr, 64)
	if err != nil {
		return 0, err
	}
	pubkey, err := address.PubKey()
	if !ed25519.Verify(pubkey, hasher, signature) {
		return 0, nil
	}
	return 1, nil
}

func (engine *Engine) chainGetContractAddress(vm *vm.VM, args ...uint64) (uint64, error) {
	addressPtr := int(args[0])
	contractAddr := engine.account.GetAddress()
	vm.MemWrite(contractAddr[:], addressPtr)
	return uint64(addressPtr), nil
}

func (engine *Engine) handleInvokeAlias(foreignMethod *foreignMethod, vm *vm.VM, args ...uint64) (uint64, error) {
	if engine.callDepth+1 > constant.MaxEngineCallDepth {
		return 0, errors.New("call depth limit reached")
	}

	foreignAccount, err := engine.state.LoadAccount(foreignMethod.contractAddress)
	if err != nil {
		return 0, err
	}
	contract, err := foreignAccount.GetContract()
	if err != nil {
		return 0, err
	}
	function, err := contract.Header.GetFunction(foreignMethod.name)
	if err != nil {
		return 0, err
	}
	if len(function.Parameters) != len(args) {
		return 0, errors.New("argument count mismatch")
	}
	var values [][]byte
	var bytes []byte
	for i, param := range function.Parameters {
		if param.IsArray {
			argPtr := int(args[i])
			size, _ := engine.ptrArgSizeGet(int(args[i]))
			bytes, err = readAt(vm, argPtr, size)
			if err != nil {
				return 0, err
			}
		} else {
			if param.Type.IsPointer() {
				argPtr := int(args[i])
				size := param.Type.GetMemorySize()
				bytes, err = readAt(vm, argPtr, size)
				if err != nil {
					return 0, err
				}
			} else {
				bytes = make([]byte, 8)
				binary.LittleEndian.PutUint64(bytes, args[i])
				size := param.Type.GetMemorySize()
				bytes = bytes[:size]
			}
		}
		values = append(values, bytes)

	}
	methodArgs, err := abi.EncodeFromBytes(function.Parameters, values)
	if err != nil {
		return 0, err
	}

	account, err := engine.state.LoadAccount(foreignMethod.contractAddress)
	if err != nil {
		return 0, err
	}
	childEngine := engine.newChildEngine(account)
	childEngine.setStats(engine.callDepth+1, engine.memAggr+vm.MemSize())
	return childEngine.Ignite(foreignMethod.name, methodArgs)
}

// GetFunction get host function for WebAssembly
func (engine *Engine) GetFunction(module, name string) vm.HostFunction {
	switch module {
	case "env":
		switch name {
		case "chain_storage_set":
			return engine.chainStorageSet
		case "chain_storage_get":
			return engine.chainStorageGet
		case "chain_storage_size_get":
			return engine.chainStorageSizeGet
		case "chain_get_caller":
			return engine.chainGetCaller
		case "chain_get_creator":
			return engine.chainGetCreator
		case "chain_method_bind":
			return engine.chainMethodBind
		case "chain_arg_size_get":
			return engine.chainPtrArgSizeGet
		case "chain_arg_size_set":
			return engine.chainPtrArgSizeSet
		case "chain_block_height":
			return engine.chainBlockHeight
		case "chain_block_time":
			return engine.chainBlockTime
		case "chain_args_write":
			return engine.chainArgsWrite
		case "chain_args_hash":
			return engine.chainArgsHash
		case "chain_ed25519_verify":
			return engine.chainEd25519Verify
		case "chain_get_contract_address":
			return engine.chainGetContractAddress
		default:
			contract, _ := engine.account.GetContract()
			if event, err := contract.Header.GetEvent(name); err == nil {
				return func(vm *vm.VM, args ...uint64) (uint64, error) {
					return engine.handleEmitEvent(event, vm, args...)
				}
			}

			if foreignMethod, ok := engine.methodLookup[name]; ok {
				return func(vm *vm.VM, args ...uint64) (uint64, error) {
					return engine.handleInvokeAlias(foreignMethod, vm, args...)
				}
			}
		}
	case "wasi_unstable":
		return wasiUnstableHandler(name)
	}
	return func(vm *vm.VM, args ...uint64) (uint64, error) {
		return 0, fmt.Errorf("unknown import %s for module %s", name, module)
	}
}
