package abi

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/google/go-cmp/cmp"
)

func parseParameterFromString(s string) (Parameter, error) {
	var p Parameter
	if strings.HasSuffix(s, "[]") {
		p.IsArray = true
		t, err := parsePrimitiveTypeFromString(s[:strings.Index(s, "[")])
		if err != nil {
			return Parameter{}, err
		}
		p.Type = t
	} else {
		p.IsArray = false
		t, err := parsePrimitiveTypeFromString(s)
		if err != nil {
			return Parameter{}, err
		}
		p.Type = t
	}
	return p, nil
}

func TestEncode(t *testing.T) {
	address, _ := crypto.AddressFromString("LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O")
	address2, _ := crypto.AddressFromString("LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O")
	addresses := []crypto.Address{address, address2}
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types  []*Parameter
		values []interface{}
		result []byte
	}{
		{
			types:  parameters1,
			values: []interface{}{address, uint8(88), uint16(43221), uint32(3333324342), uint64(3213214325432656666), int8(88), int16(4321), int32(-34325), int64(-321452), float32(8321.38), float64(-4321452.1188)},
			result: []byte{248, 86, 163, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 130, 213, 168, 132, 54, 126, 174, 198, 136, 26, 35, 156, 150, 103, 161, 151, 44, 88, 130, 225, 16, 132, 235, 121, 255, 255, 136, 84, 24, 251, 255, 255, 255, 255, 255, 132, 133, 5, 2, 70, 136, 81, 107, 154, 7, 43, 124, 80, 193},
		},
		{
			types:  parameters2,
			values: []interface{}{addresses, []uint8{uint8(88), uint8(255)}, []uint16{uint16(555), uint16(12333)}, []uint32{uint32(3333324342), uint32(3333324342), uint32(33324342)}, []uint64{uint64(3213214325432656666), uint64(32145467)}, []int8{int8(88), int8(-88)}, []int16{int16(333), int16(-542)}, []int32{int32(43298), int32(-321432)}, []int64{int64(-23425254), int64(10875498375)}, []float32{float32(-1341.233), float32(50492.235)}, []float64{float64(-132341.233), float64(50454392.235)}},
			result: []byte{248, 170, 184, 70, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 130, 88, 255, 132, 43, 2, 45, 48, 140, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 144, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 130, 88, 168, 132, 77, 1, 226, 253, 136, 34, 169, 0, 0, 104, 24, 251, 255, 144, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 136, 117, 167, 167, 196, 60, 60, 69, 71, 144, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
		},
	}
	for index, table := range testTables {
		encoded, err := Encode(table.types, table.values)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if !bytes.Equal(encoded, table.result) {
			t.Errorf("Encoding case %v: encode of %v is incorrect, expected: %v, got: %v.", index+1, table.values, table.result, encoded)
		}
	}

	testErrorsTables := []struct {
		types  []*Parameter
		values []interface{}
		err    error
	}{
		{
			types:  parameters1,
			values: []interface{}{float64(-4321452.1188)},
			err:    errors.New("Parameter count mismatch, expecting: 11, got: 1"),
		},
	}

	for index, table := range testErrorsTables {
		_, err := Encode(table.types, table.values)
		if err == nil {
			t.Errorf("expecting error at case: %v", index)
		}
		if err.Error() != table.err.Error() {
			t.Errorf("Encoding case %v: error is incorrect, expecting error: %v, got: %v.", index+1, table.err, err)
		}
	}
}

func TestEncodeFromString(t *testing.T) {
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types  []*Parameter
		values []string
		result []byte
	}{
		{
			types:  parameters1,
			values: []string{"LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O", "88", "43221", "3333324342", "3213214325432656666", "88", "4321", "-34325", "-321452", "8321.38", "-4321452.1188"},
			result: []byte{248, 86, 163, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 130, 213, 168, 132, 54, 126, 174, 198, 136, 26, 35, 156, 150, 103, 161, 151, 44, 88, 130, 225, 16, 132, 235, 121, 255, 255, 136, 84, 24, 251, 255, 255, 255, 255, 255, 132, 133, 5, 2, 70, 136, 81, 107, 154, 7, 43, 124, 80, 193},
		},
		{
			types:  parameters2,
			values: []string{"[LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O, LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O]", "[88,255]", "[555,12333]", "[3333324342,3333324342,33324342]", "[3213214325432656666,32145467]", "[88,-88]", "[333,-542]", "[43298,-321432]", "[-23425254,10875498375]", "[-1341.233,50492.235]", "[-132341.233,50454392.235]"},
			result: []byte{248, 170, 184, 70, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 130, 88, 255, 132, 43, 2, 45, 48, 140, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 144, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 130, 88, 168, 132, 77, 1, 226, 253, 136, 34, 169, 0, 0, 104, 24, 251, 255, 144, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 136, 117, 167, 167, 196, 60, 60, 69, 71, 144, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
		},
	}

	for index, table := range testTables {
		encoded, err := EncodeFromString(table.types, table.values)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if !bytes.Equal(encoded, table.result) {
			t.Errorf("Encoding case %v: encode of %v is incorrect, expected: %v, got: %v.", index+1, table.values, table.result, encoded)
		}
	}

	testErrorsTables := []struct {
		types  []*Parameter
		values []string
		err    error
	}{
		{
			types:  parameters1,
			values: []string{"-4321452.1188"},
			err:    errors.New("Argument count mismatch, expecting: 11, got: 1"),
		},
	}

	for index, table := range testErrorsTables {
		_, err := EncodeFromString(table.types, table.values)
		if err == nil {
			t.Errorf("expecting error at case: %v", index)
		}
		if err.Error() != table.err.Error() {
			t.Errorf("Encoding case %v: error is incorrect, expecting error: %v, got: %v.", index+1, table.err, err)
		}
	}
}

func TestBytesEncoding(t *testing.T) {
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types   []*Parameter
		encoded []byte
		decoded [][]byte
	}{
		{
			types:   parameters1,
			encoded: []byte{248, 86, 163, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 130, 213, 168, 132, 54, 126, 174, 198, 136, 26, 35, 156, 150, 103, 161, 151, 44, 88, 130, 225, 16, 132, 235, 121, 255, 255, 136, 84, 24, 251, 255, 255, 255, 255, 255, 132, 133, 5, 2, 70, 136, 81, 107, 154, 7, 43, 124, 80, 193},
			decoded: [][]byte{[]byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238}, []byte{88}, []byte{213, 168}, []byte{54, 126, 174, 198}, []byte{26, 35, 156, 150, 103, 161, 151, 44}, []byte{88}, []byte{225, 16}, []byte{235, 121, 255, 255}, []byte{84, 24, 251, 255, 255, 255, 255, 255}, []byte{133, 5, 2, 70}, []byte{81, 107, 154, 7, 43, 124, 80, 193}},
		},
		{
			types:   parameters2,
			encoded: []byte{248, 170, 184, 70, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 130, 88, 255, 132, 43, 2, 45, 48, 140, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 144, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 130, 88, 168, 132, 77, 1, 226, 253, 136, 34, 169, 0, 0, 104, 24, 251, 255, 144, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 136, 117, 167, 167, 196, 60, 60, 69, 71, 144, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
			decoded: [][]byte{[]byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238}, []byte{88, 255}, []byte{43, 2, 45, 48}, []byte{54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1}, []byte{26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0}, []byte{88, 168}, []byte{77, 1, 226, 253}, []byte{34, 169, 0, 0, 104, 24, 251, 255}, []byte{26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0}, []byte{117, 167, 167, 196, 60, 60, 69, 71}, []byte{160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65}},
		},
	}

	for index, table := range testTables {
		decoded, err := DecodeToBytes(table.types, table.encoded)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if diff := cmp.Diff(decoded, table.decoded); diff != "" {
			t.Errorf("Decoding case %v: decode of %v is incorrect, expected: %v, got: %v, diff: %v", index+1, table.encoded, table.decoded, decoded, diff)
		}

		encoded, err := EncodeFromBytes(table.types, decoded)
		if err != nil {
			panic(err)
		}
		if diff := cmp.Diff(encoded, table.encoded); diff != "" {
			t.Errorf("Encoding case %v: encode of %v is incorrect, expected: %v, got: %v, diff: %v", index+1, decoded, table.encoded, encoded, diff)
		}
	}

	testEncodeFromBytesErrorsTables := []struct {
		types   []*Parameter
		decoded [][]byte
		err     error
	}{
		{
			types:   parameters2,
			decoded: [][]byte{[]byte{1}},
			err:     errors.New("Argument count mismatch, expecting: 11, got: 1"),
		},
	}

	for index, table := range testEncodeFromBytesErrorsTables {
		_, err := EncodeFromBytes(table.types, table.decoded)
		if err == nil {
			t.Errorf("expecting error at case: %v", index)
		}
		if err.Error() != table.err.Error() {
			t.Errorf("Encoding case %v: error is incorrect, expecting error: %v, got: %v.", index+1, table.err, err)
		}
	}

	testDecodeToBytesErrorsTables := []struct {
		types   []*Parameter
		encoded []byte
		err     error
	}{
		{
			types:   parameters1,
			encoded: []byte{242, 88, 130, 213, 168, 132, 54, 126, 174, 198, 136, 26, 35, 156, 150, 103, 161, 151, 44, 88, 130, 225, 16, 132, 235, 121, 255, 255, 136, 84, 24, 251, 255, 255, 255, 255, 255, 132, 133, 5, 2, 70, 136, 81, 107, 154, 7, 43, 124, 80, 193},
			err:     errors.New("Argument count mismatch, expecting: 11, got: 10"),
		},
	}

	for index, table := range testDecodeToBytesErrorsTables {
		_, err := DecodeToBytes(table.types, table.encoded)
		if err == nil {
			t.Errorf("expecting error at case: %v", index)
		}
		if err.Error() != table.err.Error() {
			t.Errorf("Encoding case %v: error is incorrect, expecting error: %v, got: %v.", index+1, table.err, err)
		}
	}
}
