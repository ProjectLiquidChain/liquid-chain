package crypto

import (
	"crypto/ed25519"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"golang.org/x/crypto/blake2b"
)

// Sign return signature of message when signing using privateKey
func Sign(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}

// GetSigHash returns hash for signing transaction
func GetSigHash(tx *Transaction) common.Hash {
	encoded, _ := rlp.EncodeToBytes([]interface{}{
		tx.Version,
		tx.Sender,
		tx.Receiver,
		tx.Payload,
		tx.GasPrice,
		tx.GasLimit,
	})
	return blake2b.Sum256(encoded)
}

// VerifySignature verify whether signature valid or not
func VerifySignature(publicKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}
