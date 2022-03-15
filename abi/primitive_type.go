package abi

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// PrimitiveType PrimitiveType
type PrimitiveType uint

// enum for types
const (
	Uint8   PrimitiveType = 0x0
	Uint16  PrimitiveType = 0x1
	Uint32  PrimitiveType = 0x2
	Uint64  PrimitiveType = 0x3
	Int8    PrimitiveType = 0x4
	Int16   PrimitiveType = 0x5
	Int32   PrimitiveType = 0x6
	Int64   PrimitiveType = 0x7
	Float32 PrimitiveType = 0x8
	Float64 PrimitiveType = 0x9
	Address PrimitiveType = 0xa
)

// IsPointer return whether p is pointer or not
func (t PrimitiveType) IsPointer() bool {
	switch t {
	case Address:
		return true
	default:
		return false
	}
}

// IsAddress check if this type is an Address
func (t PrimitiveType) IsAddress() bool {
	return t == Address
}

func (t PrimitiveType) String() string {
	return map[PrimitiveType]string{
		Uint8:   "uint8",
		Uint16:  "uint16",
		Uint32:  "uint32",
		Uint64:  "uint64",
		Int8:    "int8",
		Int16:   "int16",
		Int32:   "int32",
		Int64:   "int64",
		Float32: "float32",
		Float64: "float64",
		Address: "address",
	}[t]
}

// GetMemorySize returns memory size for a primitive type
func (t PrimitiveType) GetMemorySize() int {
	switch t {
	case Address:
		return crypto.AddressLength
	case Uint8, Int8:
		return 1
	case Uint16, Int16:
		return 2
	case Uint32, Int32, Float32:
		return 4
	case Uint64, Int64, Float64:
		return 8
	default:
		panic("primitive type not found")
	}
}

// NewArgument returns a vm-compatible byte array from an interface
func (t PrimitiveType) NewArgument(value interface{}) ([]byte, error) {
	memorySize := t.GetMemorySize()
	buf := make([]byte, memorySize)
	switch t {
	case Address:
		address := value.(crypto.Address)
		copy(buf, address[:])
	case Uint8:
		buf[0] = byte(value.(uint8))
	case Uint16:
		binary.LittleEndian.PutUint16(buf, value.(uint16))
	case Uint32:
		binary.LittleEndian.PutUint32(buf, value.(uint32))
	case Uint64:
		binary.LittleEndian.PutUint64(buf, value.(uint64))
	case Int8:
		buf[0] = byte(value.(int8))
	case Int16:
		binary.LittleEndian.PutUint16(buf, uint16(value.(int16)))
	case Int32:
		binary.LittleEndian.PutUint32(buf, uint32(value.(int32)))
	case Int64:
		binary.LittleEndian.PutUint64(buf, uint64(value.(int64)))
	case Float32:
		binary.LittleEndian.PutUint32(buf, math.Float32bits(value.(float32)))
	case Float64:
		binary.LittleEndian.PutUint64(buf, math.Float64bits(value.(float64)))
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}
	return buf, nil
}

// NewArrayArgument returns a vm-compatible byte array from an interface of array
func (t PrimitiveType) NewArrayArgument(value interface{}) ([]byte, error) {
	var parsedArgs []byte
	switch t {
	case Address:
		parsed, ok := value.([]crypto.Address)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Uint8:
		parsed, ok := value.([]uint8)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Uint16:
		parsed, ok := value.([]uint16)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Uint32:
		parsed, ok := value.([]uint32)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Uint64:
		parsed, ok := value.([]uint64)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Int8:
		parsed, ok := value.([]int8)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Int16:
		parsed, ok := value.([]int16)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Int32:
		parsed, ok := value.([]int32)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Int64:
		parsed, ok := value.([]int64)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Float32:
		parsed, ok := value.([]float32)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	case Float64:
		parsed, ok := value.([]float64)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t.String())
		}
		for _, p := range parsed {
			arg, err := t.NewArgument(p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg...)
		}
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}

	return parsedArgs, nil
}
