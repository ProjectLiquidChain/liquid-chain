package abi

import (
	"errors"
	"testing"
)

func TestIsPointer(t *testing.T) {

	testTables := []struct {
		types     []string
		isPointer bool
	}{
		{
			types:     []string{"address"},
			isPointer: true,
		},
		{
			types:     []string{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"},
			isPointer: false,
		},
	}

	for index, table := range testTables {
		for _, testType := range table.types {
			p, _ := parsePrimitiveTypeFromString(testType)
			if p.IsPointer() != table.isPointer {
				t.Errorf("case %v: expecting type %v isPointer to be %v", index+1, p.String(), table.isPointer)
			}
		}
	}
}

func TestNewArrayArgument(t *testing.T) {
	var primitiveTypes []PrimitiveType
	types := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	for _, t := range types {
		pType, _ := parsePrimitiveTypeFromString(t)
		primitiveTypes = append(primitiveTypes, pType)
	}
	testErrorsTables := []struct {
		values interface{}
		err    error
	}{
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into address"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into uint8"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into uint16"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into uint32"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into uint64"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into int8"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into int16"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into int32"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into int64"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into float32"),
		},
		{
			values: []interface{}{"321432145321"},
			err:    errors.New("unable to convert array element into float64"),
		},
	}

	for index, table := range testErrorsTables {
		_, err := primitiveTypes[index].NewArrayArgument(table.values)
		if err == nil {
			t.Errorf("expecting error at case: %v", index)
		}
		if err.Error() != table.err.Error() {
			t.Errorf("Encoding case %v: error is incorrect, expecting error: %v, got: %v.", index+1, table.err, err)
		}
	}
}
