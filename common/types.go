// implementation copied from ethereum/go-ethereum

package common

import (
	"encoding/hex"
	"encoding/json"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// EmptyHash is initial hash
var EmptyHash = Hash{}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return BytesToHash(b)
}

// HashToHex encodes b as a hex string
func HashToHex(b []byte) string {
	enc := make([]byte, len(b)*2)
	hex.Encode(enc, b)
	return string(enc)
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return HashToHex(h[:])
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte {
	if h == EmptyHash {
		return nil
	}
	return h[:]
}
