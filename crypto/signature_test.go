package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/google/go-cmp/cmp"
)

func TestSign(t *testing.T) {
	type args struct {
		privateKey ed25519.PrivateKey
		message    []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{{
		args: args{
			privateKey: ed25519.NewKeyFromSeed([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2}),
			message:    []byte("Hello"),
		},
		want: []byte{255, 97, 130, 204, 8, 94, 179, 251, 116, 129, 109, 108, 78, 211, 46, 8, 230, 12, 195, 77, 254, 118, 133, 255, 251, 211, 102, 218, 216, 146, 66, 77, 194, 215, 69, 77, 24, 168, 24, 186, 114, 237, 181, 78, 177, 57, 17, 75, 186, 210, 204, 37, 55, 219, 237, 197, 176, 105, 188, 209, 82, 183, 208, 7},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sign(tt.args.privateKey, tt.args.message); !cmp.Equal(got, tt.want) {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSigHash(t *testing.T) {
	type args struct {
		tx *Transaction
	}
	tests := []struct {
		name string
		args args
		want common.Hash
	}{{
		args: args{
			tx: &Transaction{
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
		},
		want: common.HexToHash("1c3346e14d541923b557a4ca8bb28f5dd26ca2f2cfe6e3d2558eb5a66864d45e"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSigHash(tt.args.tx); !cmp.Equal(got, tt.want) {
				t.Errorf("GetSigHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifySignature(t *testing.T) {
	seedFirst := make([]byte, 32)
	rand.Read(seedFirst)

	seedSecond := make([]byte, 32)
	rand.Read(seedSecond)

	type args struct {
		publicKey ed25519.PublicKey
		message   []byte
		signature []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "valid",
		args: args{
			publicKey: ed25519.NewKeyFromSeed(seedFirst).Public().(ed25519.PublicKey),
			message:   []byte("Hello"),
			signature: ed25519.Sign(ed25519.NewKeyFromSeed(seedFirst), []byte("Hello")),
		},
		want: true,
	}, {
		name: "invalid",
		args: args{
			publicKey: ed25519.NewKeyFromSeed(seedFirst).Public().(ed25519.PublicKey),
			message:   []byte("Hello"),
			signature: ed25519.Sign(ed25519.NewKeyFromSeed(seedSecond), []byte("Hello")),
		},
		want: false,
	}, {
		name: "invalid",
		args: args{
			publicKey: ed25519.NewKeyFromSeed(seedFirst).Public().(ed25519.PublicKey),
			message:   []byte("Hello"),
			signature: []byte{1, 2, 3},
		},
		want: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifySignature(tt.args.publicKey, tt.args.message, tt.args.signature); got != tt.want {
				t.Errorf("VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
