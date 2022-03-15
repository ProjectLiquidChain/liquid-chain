package crypto

import (
	"crypto/ed25519"
	"encoding/base32"
	"encoding/json"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/crc16"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

const (
	versionByteAccountID byte = 11 << 3 // Base32-encodes to 'L...'
	// AddressLength size of a crypto address
	AddressLength = 35
)

// Address crypto address
type Address [AddressLength]byte

// EmptyAddress is address to deploy tx
var EmptyAddress = Address{}

func (address *Address) setBytes(b []byte) {
	if len(b) > len(address) {
		b = b[len(b)-AddressLength:]
	}
	copy(address[AddressLength-len(b):], b)
}

func (address Address) MarshalJSON() ([]byte, error) {
	if address == EmptyAddress {
		return json.Marshal(nil)
	}
	return json.Marshal(address.String())
}

// String Address string presentation
func (address *Address) String() string {
	return base32.StdEncoding.EncodeToString(address[:])
}

// PubKey retrieves public key of an address
func (address *Address) PubKey() (ed25519.PublicKey, error) {
	return decodeAddressBytes(address[:])
}

// AddressFromPubKey create an address from public key
func AddressFromPubKey(publicKey ed25519.PublicKey) Address {
	data := append([]byte{versionByteAccountID}, publicKey...)
	var a Address
	a.setBytes(append(data, crc16.Checksum(data)...))
	return a
}

// AddressFromString parse an address string to Address
func AddressFromString(address string) (Address, error) {
	pubKeyBytes, err := decodeAddress(address)
	if err != nil {
		return Address{}, err
	}
	pubkey := ed25519.PublicKey(pubKeyBytes)
	return AddressFromPubKey(pubkey), nil
}

// AddressFromBytes return an address given its bytes
func AddressFromBytes(b []byte) (Address, error) {
	if b == nil {
		return EmptyAddress, nil
	}

	var a Address
	_, err := decodeAddressBytes(b)
	if err != nil {
		return Address{}, err
	}
	a.setBytes(b)
	return a, nil
}

func decodeString(src string) ([]byte, error) {
	raw, err := base32.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, errors.Wrap(err, "base32 decode failed")
	}

	if len(raw) < 3 {
		return nil, errors.Errorf("encoded value is %d bytes; minimum valid length is 3", len(raw))
	}

	return raw, nil
}

func decodeAddress(src string) ([]byte, error) {
	raw, err := decodeString(src)
	if err != nil {
		return nil, err
	}
	return decodeAddressBytes(raw)
}

func decodeAddressBytes(raw []byte) ([]byte, error) {
	version := byte(raw[0])
	payload := raw[1 : len(raw)-2]
	checksum := raw[len(raw)-2:]
	original := raw[0 : len(raw)-2]

	if version != versionByteAccountID {
		return nil, errors.Errorf("Unexpected address version %x", version)
	}

	// checksum check
	if err := crc16.Validate(original, checksum); err != nil {
		return nil, err
	}
	return payload, nil
}

// NewDeploymentAddress returns new contract deployment address
func NewDeploymentAddress(senderAddress Address, senderNonce uint64) Address {
	senderBytes, _ := rlp.EncodeToBytes([]interface{}{senderAddress, senderNonce})
	res := blake2b.Sum256(senderBytes)
	address := AddressFromPubKey(res[:])
	return address
}
