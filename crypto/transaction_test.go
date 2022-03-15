package crypto

import (
	"crypto/ed25519"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/google/go-cmp/cmp"
)

func TestTransaction_Serialize(t *testing.T) {
	type fields struct {
		Version   uint16
		Sender    *TxSender
		Receiver  Address
		Payload   *TxPayload
		GasPrice  uint32
		GasLimit  uint32
		Signature []byte
		Receipt   *Receipt
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{{
		fields: fields{
			Version: 1,
			Sender: &TxSender{
				Nonce:     uint64(0),
				PublicKey: ed25519.NewKeyFromSeed(make([]byte, 32)).Public().(ed25519.PublicKey),
			},
			Receiver: Address{},
			Payload: &TxPayload{
				Contract: []byte{1, 2, 3},
				ID:       GetMethodID("Transfer"),
				Args:     []byte{4, 5, 6},
			},
			GasPrice:  1,
			GasLimit:  2,
			Signature: []byte{7, 8, 9},
			Receipt: &Receipt{
				Result:  1,
				GasUsed: 2,
				Code:    ReceiptCodeOK,
				Events: []*Event{{
					Contract: Address{},
					Args:     []byte{10, 11, 12},
				}},
			},
		},
		want:    []byte{248, 92, 1, 226, 160, 59, 106, 39, 188, 206, 182, 164, 45, 98, 163, 168, 208, 42, 111, 13, 115, 101, 50, 21, 119, 29, 226, 67, 166, 58, 192, 72, 161, 139, 89, 218, 41, 128, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 205, 132, 34, 60, 57, 226, 131, 4, 5, 6, 131, 1, 2, 3, 1, 2, 131, 7, 8, 9},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Transaction{
				Version:   tt.fields.Version,
				Sender:    tt.fields.Sender,
				Receiver:  tt.fields.Receiver,
				Payload:   tt.fields.Payload,
				GasPrice:  tt.fields.GasPrice,
				GasLimit:  tt.fields.GasLimit,
				Signature: tt.fields.Signature,
			}
			got, err := tx.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Transaction.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_Deserialize(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Transaction
		wantErr bool
	}{{
		name: "invalid",
		args: args{
			raw: []byte{1, 2, 3},
		},
		wantErr: true,
	}, {
		name: "valid",
		args: args{
			raw: []byte{248, 92, 1, 226, 160, 59, 106, 39, 188, 206, 182, 164, 45, 98, 163, 168, 208, 42, 111, 13, 115, 101, 50, 21, 119, 29, 226, 67, 166, 58, 192, 72, 161, 139, 89, 218, 41, 128, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 205, 132, 34, 60, 57, 226, 131, 4, 5, 6, 131, 1, 2, 3, 1, 2, 131, 7, 8, 9},
		},
		want: Transaction{
			Version: 1,
			Sender: &TxSender{
				Nonce:     uint64(0),
				PublicKey: ed25519.NewKeyFromSeed(make([]byte, 32)).Public().(ed25519.PublicKey),
			},
			Receiver: Address{},
			Payload: &TxPayload{
				Contract: []byte{1, 2, 3},
				ID:       GetMethodID("Transfer"),
				Args:     []byte{4, 5, 6},
			},
			GasPrice:  1,
			GasLimit:  2,
			Signature: []byte{7, 8, 9},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := DecodeTransaction(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if equal := cmp.Equal((*tx), tt.want); !equal {
					t.Errorf("Transaction.Deserialize() %v, want %v", tx, tt.wantErr)
				}
			}
		})
	}
}

func TestTransaction_Hash(t *testing.T) {
	type fields struct {
		Version   uint16
		Sender    *TxSender
		Receiver  Address
		Payload   *TxPayload
		GasPrice  uint32
		GasLimit  uint32
		Signature []byte
		Receipt   *Receipt
	}
	tests := []struct {
		name   string
		fields fields
		want   common.Hash
	}{{
		fields: fields{
			Version: 1,
			Sender: &TxSender{
				Nonce:     uint64(0),
				PublicKey: ed25519.NewKeyFromSeed(make([]byte, 32)).Public().(ed25519.PublicKey),
			},
			Receiver: Address{},
			Payload: &TxPayload{
				Contract: []byte{1, 2, 3},
				ID:       GetMethodID("Transfer"),
				Args:     []byte{4, 5, 6},
			},
			GasPrice:  1,
			GasLimit:  2,
			Signature: []byte{7, 8, 9},
			Receipt: &Receipt{
				Result:  1,
				GasUsed: 2,
				Code:    ReceiptCodeOK,
				Events: []*Event{{
					Contract: Address{},
					Args:     []byte{10, 11, 12},
				}},
			},
		},
		want: common.HexToHash("146dc33a1f41bc435a323ce1bd0c3ceaf49fe14dc869b4da7ba482fb52735f41"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Transaction{
				Version:   tt.fields.Version,
				Sender:    tt.fields.Sender,
				Receiver:  tt.fields.Receiver,
				Payload:   tt.fields.Payload,
				GasPrice:  tt.fields.GasPrice,
				GasLimit:  tt.fields.GasLimit,
				Signature: tt.fields.Signature,
			}
			if got := tx.Hash(); !cmp.Equal(got, tt.want) {
				t.Errorf("Transaction.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
