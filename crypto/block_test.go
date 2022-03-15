package crypto

import (
	"bytes"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/google/go-cmp/cmp"
)

func TestBlock_Hash(t *testing.T) {
	type fields struct {
		block *Block
	}
	tests := []struct {
		name   string
		fields fields
		want   common.Hash
	}{{
		fields: fields{
			block: &Block{
				Time:            123,
				Height:          1,
				Parent:          common.HexToHash("2f636344b757343e13e7910eed1b832d769e1d113027424580a2faca232ce015"),
				StateRoot:       common.HexToHash("572343bcdac17dbae1aba2d1ccde3488adb169b18da8a4ecdffe11c8f1cc1f1f"),
				TransactionRoot: common.HexToHash("3e2e21d19f5c3491ea8d5416b44256c401596b184638e63d8ac34f073a686544"),
			},
		},
		want: common.HexToHash("f78a6b6423b6656ba57c777e35fc33dc728314afc189306577ccb47a6daeb551"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.block.Hash(); !cmp.Equal(got, tt.want) {
				t.Errorf("Block.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock(t *testing.T) {
	block := NewEmptyBlock(common.EmptyHash, 0, time.Unix(0, 0))
	block.SetStateRoot(common.BytesToHash([]byte{1, 2, 3}))
	block.SetTransactionRoot(common.BytesToHash([]byte{1, 2, 3}))
	encoded, _ := block.Encode()
	decodedBlock := MustDecodeBlock(encoded)
	if decodedBlock.Hash() != block.Hash() {
		t.Errorf("Got block hash after decoded = %v, want %v", decodedBlock.Hash(), block.Hash())
	}

	encodedNew, _ := decodedBlock.Encode()
	if !bytes.Equal(encoded, encodedNew) {
		t.Errorf("Encode not equal, got = %v, want %v", encodedNew, encoded)
	}
}

func TestMustDecodeBlock(t *testing.T) {
	// This decoding should panic
	defer func() { recover() }()
	MustDecodeBlock([]byte{1, 2, 3})
	t.Errorf("did not panic")
}
