package abi

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEncodeHeaderFromFile(t *testing.T) {
	encoded, err := EncodeHeaderToBytes("../test/testdata/liquid-token-abi.json")
	if err != nil {
		t.Errorf("error: %s", err)
	}
	result := []byte{0xf8, 0x9d, 0x1, 0xf8, 0x56, 0xd0, 0x84, 0x69, 0x6e, 0x69, 0x74, 0xca, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xda, 0x88, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0xd0, 0xc5, 0x82, 0x74, 0x6f, 0x80, 0xa, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xd0, 0x84, 0x6d, 0x69, 0x6e, 0x74, 0xca, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xd8, 0x8b, 0x67, 0x65, 0x74, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0xcb, 0xca, 0x87, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x80, 0xa, 0xf8, 0x42, 0xd6, 0x84, 0x4d, 0x69, 0x6e, 0x74, 0xd0, 0xc5, 0x82, 0x74, 0x6f, 0x80, 0xa, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xea, 0x88, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0xe0, 0xc7, 0x84, 0x66, 0x72, 0x6f, 0x6d, 0x80, 0xa, 0xc5, 0x82, 0x74, 0x6f, 0x80, 0xa, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xc7, 0x84, 0x6d, 0x65, 0x6d, 0x6f, 0x80, 0x3}
	if !bytes.Equal(encoded, result) {
		t.Errorf("Encoding is incorrect,\nexpected:\t%#v\nreality:\t%#v.", result, encoded)
	}
}

func TestDecodeHeader(t *testing.T) {
	tests := []struct {
		name        string
		inputFile   string
		result      []byte
		resultError error
	}{{
		name:        "normal",
		inputFile:   "../test/testdata/liquid-token-abi.json",
		result:      []byte{0xf8, 0x9d, 0x1, 0xf8, 0x56, 0xd0, 0x84, 0x69, 0x6e, 0x69, 0x74, 0xca, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xda, 0x88, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0xd0, 0xc5, 0x82, 0x74, 0x6f, 0x80, 0xa, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xd0, 0x84, 0x6d, 0x69, 0x6e, 0x74, 0xca, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xd8, 0x8b, 0x67, 0x65, 0x74, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0xcb, 0xca, 0x87, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x80, 0xa, 0xf8, 0x42, 0xd6, 0x84, 0x4d, 0x69, 0x6e, 0x74, 0xd0, 0xc5, 0x82, 0x74, 0x6f, 0x80, 0xa, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xea, 0x88, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0xe0, 0xc7, 0x84, 0x66, 0x72, 0x6f, 0x6d, 0x80, 0xa, 0xc5, 0x82, 0x74, 0x6f, 0x80, 0xa, 0xc9, 0x86, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x80, 0x3, 0xc7, 0x84, 0x6d, 0x65, 0x6d, 0x6f, 0x80, 0x3},
		resultError: nil,
	}, {
		name:        "duplicated MethodID of events",
		inputFile:   "../test/testdata/events-duplicated-method-id.json",
		result:      nil,
		resultError: ErrDuplicatedEventsMethodID,
	}, {
		name:        "duplicated MethodID of functions",
		inputFile:   "../test/testdata/functions-duplicated-method-id.json",
		result:      nil,
		resultError: ErrDuplicatedFunctionsMethodID,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			bytes, err := EncodeHeaderToBytes(test.inputFile)
			if err != nil {
				t.Error(err)
			}

			_, err = DecodeHeader(bytes)
			if test.resultError != nil {
				if err != test.resultError {
					t.Errorf("want err %v, got %v", test.resultError, err)

				}
			} else {
				if err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff(bytes, test.result); diff != "" {
					t.Errorf("DecodeHeader error, expected: %#v, got: %#v", test.result, bytes)
				}
			}

		})
	}

}

func TestGetEvent(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	event, err := h.GetEvent("Transfer")
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})
	if diff := cmp.Diff(event, h.Events[crypto.GetMethodID("Transfer")], opts); diff != "" {
		t.Errorf("GetEvent of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Events[crypto.GetMethodID("Transfer")], event, diff)
	}

	notFoundEvent, err := h.GetEvent("nil")
	if err == nil {
		t.Error("expecting error is nil for getting not found event")
	}
	if notFoundEvent != nil || err.Error() != "event nil not found" {
		t.Errorf("Error of GetEvent of %v is incorrect, expected: %v, got: %v", h, "event nil not found", err.Error())
	}
}

func TestGetFunction(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	event, err := h.GetFunction("transfer")
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})
	if diff := cmp.Diff(event, h.Functions[crypto.GetMethodID("transfer")], opts); diff != "" {
		t.Errorf("GetFunction of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Functions[crypto.GetMethodID("transfer")], event, diff)
	}

	notFoundFunction, err := h.GetFunction("nil")
	if err == nil {
		t.Error("expecting error is nil for getting not found function")
	}
	if notFoundFunction != nil || err.Error() != "function nil not found" {
		t.Errorf("Error of GetFunction of %v is incorrect, expected: %v, got: %v", h, "function nil not found", err.Error())
	}
}

func TestGetFunctionByMethodID(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	function, err := h.GetFunctionByMethodID(crypto.GetMethodID("transfer"))
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})
	if diff := cmp.Diff(function, h.Functions[crypto.GetMethodID("transfer")], opts); diff != "" {
		t.Errorf("GetFunction of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Functions[crypto.GetMethodID("transfer")], function, diff)
	}

	notFoundFunction, err := h.GetFunctionByMethodID(crypto.MethodID{})
	if err == nil {
		t.Error("expecting error is nil for getting not found function")
	}
	expectedErr := fmt.Sprintf("function with methodID %v not found", crypto.MethodID{})
	if notFoundFunction != nil || err.Error() != expectedErr {
		t.Errorf("Error of GetFunction of %v is incorrect, expected: %v, got: %v", h, expectedErr, err.Error())
	}
}

func TestSortMethodIDs(t *testing.T) {
	tests := []struct {
		inputs []crypto.MethodID
		result []crypto.MethodID
	}{{
		inputs: []crypto.MethodID{
			{0, 1, 4, 0},
			{0, 1, 2, 0},
			{0, 1, 3, 0},
		},
		result: []crypto.MethodID{
			{0, 1, 2, 0},
			{0, 1, 3, 0},
			{0, 1, 4, 0},
		},
	}}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			sortMethodIDs(tt.inputs)
			for i := range tt.inputs {
				if !bytes.Equal(tt.inputs[i][:], tt.result[i][:]) {
					t.Errorf("sortMethodIDs failed\n- expected\n%#v\n- got:\n%#v", tt.result, tt.inputs)
				}
			}
		})
	}
}
