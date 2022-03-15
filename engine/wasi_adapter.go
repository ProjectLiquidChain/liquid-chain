package engine

import (
	"fmt"

	"github.com/vertexdlt/vertexvm/vm"
)

func wasiUnstableHandler(name string) vm.HostFunction {
	switch name {
	case "proc_exit":
		return wasiProcExit
	case "proc_raise":
		return wasiProcRaise
	default:
		return func(vm *vm.VM, args ...uint64) (uint64, error) {
			return wasiDefaultHandler(name, vm, args...)
		}
	}
}

func wasiDefaultHandler(funcName string, vm *vm.VM, args ...uint64) (uint64, error) {
	return 0, fmt.Errorf("unsupported func call %s", funcName)
}

func wasiProcExit(vm *vm.VM, args ...uint64) (uint64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("invalid proc_exit argument")
	}
	return args[0], fmt.Errorf("process exit with code: %d", args[0])
}

func wasiProcRaise(vm *vm.VM, args ...uint64) (uint64, error) {
	return wasiProcExit(vm, args...)
}
