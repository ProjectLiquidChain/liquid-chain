package engine

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/vertexdlt/vertexvm/vm"
	vertex "github.com/vertexdlt/vertexvm/vm"
)

type testcase struct {
	name          string
	params        []uint64
	expected      uint64
	expectedError string
	entry         string
}

// Read code from wasm file and init a VM with basic configuration
func getVM(filename string) *vm.VM {
	wasm := fmt.Sprintf("testdata/%s.wasm", filename)
	code, err := ioutil.ReadFile(wasm)
	if err != nil {
		panic(err)
	}

	engine := &Engine{}
	gasLimit := uint64(0)
	gasPolicy := &gas.FreePolicy{}
	vm, err := vertex.NewVM(code, gasPolicy, &vm.Gas{Limit: gasLimit}, engine)
	if err != nil {
		panic(err)
	}
	return vm
}

func TestAdapterGetFunction(t *testing.T) {
	tests := []testcase{
		{name: "exit_invalid", entry: "calc", params: []uint64{1}, expectedError: "invalid proc_exit argument"},
		{name: "exit", entry: "calc", params: []uint64{1}, expectedError: "process exit with code: 1"},
	}

	for _, test := range tests {
		vm := getVM(test.name)
		fnID, ok := vm.GetFunctionIndex(test.entry)
		if !ok {
			t.Error("cannot get function export")
		}
		ret, err := vm.Invoke(fnID, test.params...)
		if err.Error() != test.expectedError {
			t.Errorf("Test %s: Expect return value to be %s, got %s", test.name, test.expectedError, err)
		}
		if ret != test.expected {
			t.Errorf("Test %s: Expect return value to be %d, got %d", test.name, test.expected, ret)
		}
	}
}
