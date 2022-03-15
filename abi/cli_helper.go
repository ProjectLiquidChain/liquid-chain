package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// ParameterFile ParameterFile
type ParameterFile struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// HeaderFile representation of Header file
type HeaderFile struct {
	Version uint16 `json:"version"`
	Events  []struct {
		Name       string          `json:"name"`
		Parameters []ParameterFile `json:"parameters"`
	} `json:"events"`
	Functions []struct {
		Name       string          `json:"name"`
		Parameters []ParameterFile `json:"parameters"`
	}
}

// IsArray returns true if ParameterFile type ends in [] (an array)
func (p ParameterFile) IsArray() bool {
	return strings.HasSuffix(p.Type, "[]")
}

func parsePrimitiveTypeFromString(t string) (PrimitiveType, error) {
	var primitiveType PrimitiveType
	switch t {
	case "address":
		primitiveType = Address
	case "uint8":
		primitiveType = Uint8
	case "uint16":
		primitiveType = Uint16
	case "uint32":
		primitiveType = Uint32
	case "uint64":
		primitiveType = Uint64
	case "int8":
		primitiveType = Int8
	case "int16":
		primitiveType = Int16
	case "int32":
		primitiveType = Int32
	case "int64":
		primitiveType = Int64
	case "float32":
		primitiveType = Float32
	case "float64":
		primitiveType = Float64
	default:
		return primitiveType, fmt.Errorf("not supported type: %s for parsePrimitiveTypeFromString", t)
	}
	return primitiveType, nil
}

func parseArrayArgsFromString(t PrimitiveType, value string) (interface{}, error) {
	if !(value[0] == '[' && value[len(value)-1] == ']') {
		return nil, fmt.Errorf("wrong array value format, expected [value], got: %s", value)
	}

	args := strings.Split(value[1:len(value)-1], ",")

	switch t {
	case Address:
		slices := []crypto.Address{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(crypto.Address))
		}
		return slices, nil
	case Uint8:
		slices := []uint8{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(uint8))
		}
		return slices, nil
	case Uint16:
		slices := []uint16{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(uint16))
		}
		return slices, nil
	case Uint32:
		slices := []uint32{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(uint32))
		}
		return slices, nil
	case Uint64:
		slices := []uint64{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(uint64))
		}
		return slices, nil
	case Int8:
		slices := []int8{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(int8))
		}
		return slices, nil
	case Int16:
		slices := []int16{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(int16))
		}
		return slices, nil
	case Int32:
		slices := []int32{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(int32))
		}
		return slices, nil
	case Int64:
		slices := []int64{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(int64))
		}
		return slices, nil
	case Float32:
		slices := []float32{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(float32))
		}
		return slices, nil
	case Float64:
		slices := []float64{}
		for _, arg := range args {
			result, err := parseArgFromString(t, arg)
			if err != nil {
				return nil, err
			}
			slices = append(slices, result.(float64))
		}
		return slices, nil
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}
}

func parseArgFromString(t PrimitiveType, value string) (interface{}, error) {
	var result interface{}
	value = strings.TrimSpace(value)
	switch t {
	case Address:
		return crypto.AddressFromString(value)
	case Uint8:
		param, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, err
		}
		result = uint8(param)
	case Uint16:
		param, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, err
		}
		result = uint16(param)
	case Uint32:
		param, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}
		result = uint32(param)
	case Uint64:
		param, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		result = uint64(param)
	case Int8:
		param, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return nil, err
		}
		result = int8(param)
	case Int16:
		param, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return nil, err
		}
		result = int16(param)
	case Int32:
		param, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		result = int32(param)
	case Int64:
		param, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		result = int64(param)
	case Float32:
		param, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		result = float32(param)
	case Float64:
		param, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		result = float64(param)
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}
	return result, nil
}

// EncodeFromString return []byte from an inputted types and values type of string slices
func EncodeFromString(params []*Parameter, values []string) ([]byte, error) {
	var interfaces []interface{}
	if len(params) != len(values) {
		return []byte{0}, fmt.Errorf("Argument count mismatch, expecting: %d, got: %d", len(params), len(values))
	}
	for index, param := range params {
		if param.IsArray {
			arg, err := parseArrayArgsFromString(param.Type, values[index])
			if err != nil {
				return []byte{0}, err
			}
			interfaces = append(interfaces, arg)
		} else {
			arg, err := parseArgFromString(param.Type, values[index])
			if err != nil {
				return []byte{0}, err
			}
			interfaces = append(interfaces, arg)
		}
	}

	encoded, err := Encode(params, interfaces)
	if err != nil {
		return []byte{0}, err
	}
	return encoded, nil
}

// EncodeHeaderToBytes encode a header file into byte array
func EncodeHeaderToBytes(path string) ([]byte, error) {
	var headerFile HeaderFile
	headerFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(headerFileContent, &headerFile); err != nil {
		return nil, err
	}
	header := struct {
		Version   uint16
		Functions []*Function
		Events    []*Event
	}{
		Version:   headerFile.Version,
		Functions: []*Function{},
		Events:    []*Event{},
	}

	for _, hFunction := range headerFile.Functions {
		function := Function{
			Name:       hFunction.Name,
			Parameters: []*Parameter{},
			id:         crypto.GetMethodID(hFunction.Name),
		}
		for _, hParam := range hFunction.Parameters {
			var parameter Parameter
			var paramType PrimitiveType
			var pType string
			parameter.IsArray = hParam.IsArray()
			if parameter.IsArray {
				pType = hParam.Type[:len(hParam.Type)-2]
			} else {
				pType = hParam.Type
			}
			paramType, err = parsePrimitiveTypeFromString(pType)
			if err != nil {
				return nil, err
			}
			parameter.Type = paramType
			parameter.Name = hParam.Name
			function.Parameters = append(function.Parameters, &parameter)
		}
		header.Functions = append(header.Functions, &function)
	}

	for _, hEvent := range headerFile.Events {
		event := Event{
			Name:       hEvent.Name,
			Parameters: []*Parameter{},
		}
		for _, hParam := range hEvent.Parameters {
			var parameter Parameter
			var paramType PrimitiveType
			var pType string
			parameter.IsArray = hParam.IsArray()
			if parameter.IsArray {
				pType = hParam.Type[:len(hParam.Type)-2]
			} else {
				pType = hParam.Type
			}
			paramType, err = parsePrimitiveTypeFromString(pType)
			if err != nil {
				return nil, err
			}
			parameter.Type = paramType
			parameter.Name = hParam.Name
			event.Parameters = append(event.Parameters, &parameter)
		}
		header.Events = append(header.Events, &event)
	}

	sort.Slice(header.Functions, func(i, j int) bool {
		return bytes.Compare(header.Functions[i].id[:], header.Functions[j].id[:]) == -1
	})

	sort.Slice(header.Events, func(i, j int) bool {
		return bytes.Compare(header.Events[i].id[:], header.Events[j].id[:]) == -1
	})

	return rlp.EncodeToBytes(header)
}

// LoadHeaderFromFile load a header file into Header
func LoadHeaderFromFile(path string) (*Header, error) {
	bytes, err := EncodeHeaderToBytes(path)
	if err != nil {
		return nil, err
	}
	return DecodeHeader(bytes)
}
