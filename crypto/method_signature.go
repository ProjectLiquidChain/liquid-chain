package crypto

import "golang.org/x/crypto/blake2b"

const methodIDLength int = 4

// MethodID is first 4 bytes taken from hash of method name
type MethodID [methodIDLength]byte

// GetMethodID return MethodID of method name
func GetMethodID(name string) MethodID {
	var id MethodID
	hash := blake2b.Sum256([]byte(name))
	copy(id[:methodIDLength], hash[:methodIDLength])
	return id
}
