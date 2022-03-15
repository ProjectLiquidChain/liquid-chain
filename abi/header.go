package abi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

var (
	ErrDuplicatedFunctionsMethodID = errors.New("duplicated MethodID of functions")
	ErrDuplicatedEventsMethodID    = errors.New("duplicated MethodID of events")
)

// Event is emitting from engine
type Event struct {
	Name       string       `json:"name"`
	Parameters []*Parameter `json:"parameters"`
	id         crypto.MethodID
}

// Parameter describes a param of method
type Parameter struct {
	Name    string        `json:"name"`
	IsArray bool          `json:"-"`
	Type    PrimitiveType `json:"type"`
}

// Function describes a function in contract
type Function struct {
	Name       string       `json:"name"`
	Parameters []*Parameter `json:"parameters"`
	id         crypto.MethodID
}

// Header contains declaration for contract
type Header struct {
	Version   uint16
	Functions map[crypto.MethodID]*Function
	Events    map[crypto.MethodID]*Event
}

// GetFunctionByMethodID return Function by its id
func (h Header) GetFunctionByMethodID(id crypto.MethodID) (*Function, error) {
	if function, ok := h.Functions[id]; ok {
		return function, nil
	}
	return nil, fmt.Errorf("function with methodID %v not found", id)
}

// GetEvent return the event
func (h Header) GetEvent(name string) (*Event, error) {
	if event, ok := h.Events[crypto.GetMethodID(name)]; ok {
		return event, nil
	}
	return nil, fmt.Errorf("event %s not found", name)
}

// GetFunction returns function of a header from the func name
func (h Header) GetFunction(funcName string) (*Function, error) {
	if f, found := h.Functions[crypto.GetMethodID(funcName)]; found {
		return f, nil
	}
	return nil, fmt.Errorf("function %s not found", funcName)
}

// DecodeHeader decode byte array of header into header
func DecodeHeader(b []byte) (*Header, error) {
	var header struct {
		Version   uint16
		Functions []*Function
		Events    []*Event
	}
	if err := rlp.DecodeBytes(b, &header); err != nil {
		return nil, err
	}

	functions := make(map[crypto.MethodID]*Function)
	for _, function := range header.Functions {
		function.id = crypto.GetMethodID(function.Name)
		if _, duplicated := functions[function.id]; duplicated {
			return nil, ErrDuplicatedFunctionsMethodID
		}
		functions[function.id] = function
	}

	events := make(map[crypto.MethodID]*Event)
	for _, event := range header.Events {
		event.id = crypto.GetMethodID(event.Name)
		if _, duplicated := events[event.id]; duplicated {
			return nil, ErrDuplicatedEventsMethodID
		}
		events[event.id] = event
	}

	return &Header{header.Version, functions, events}, nil
}

// Encode encode a header struct into byte array
// encoding schema: version(2 bytes)|number of functions(1 byte)|function1|function2|...
func (h *Header) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(h)
}

func (h *Header) getEvents() []*Event {
	var ids []crypto.MethodID
	for id := range h.Events {
		ids = append(ids, id)
	}
	sortMethodIDs(ids)
	events := []*Event{}
	for _, id := range ids {
		events = append(events, h.Events[id])
	}
	return events
}

func sortMethodIDs(ids []crypto.MethodID) {
	sort.Slice(ids, func(i, j int) bool {
		return bytes.Compare(ids[i][:], ids[j][:]) == -1
	})
}

func (h *Header) getFunctions() []*Function {
	var ids []crypto.MethodID
	for id := range h.Functions {
		ids = append(ids, id)
	}
	sortMethodIDs(ids)
	functions := []*Function{}
	for _, id := range ids {
		functions = append(functions, h.Functions[id])
	}
	return functions
}

// EncodeRLP encodes a header to RLP format
func (h *Header) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, struct {
		Version   uint16
		Functions []*Function
		Events    []*Event
	}{
		Version:   h.Version,
		Functions: h.getFunctions(),
		Events:    h.getEvents(),
	})
}

// MarshalJSON returns json string of header
func (h *Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version   uint16      `json:"version"`
		Events    []*Event    `json:"events"`
		Functions []*Function `json:"functions"`
	}{
		Version:   h.Version,
		Events:    h.getEvents(),
		Functions: h.getFunctions(),
	})
}

// MarshalJSON returns json string of Parameter
func (p *Parameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name    string `json:"name"`
		IsArray bool   `json:"-"`
		Type    string `json:"type"`
		Size    uint   `json:"size,omitempty"`
	}{
		Name: p.Name,
		Type: p.Type.String(),
	})
}
