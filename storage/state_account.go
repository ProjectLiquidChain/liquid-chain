package storage

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"golang.org/x/crypto/blake2b"
)

// Account stores information related to the account
type Account struct {
	Nonce        uint64         `json:"nonce"`
	ContractHash common.Hash    `json:"contractHash"`
	StorageHash  common.Hash    `json:"storageHash"`
	Creator      crypto.Address `json:"creator"`

	dirty    bool
	address  crypto.Address
	storage  *trie.Trie
	contract []byte
}

// loadAccount load the account from disk
func (state *StateStorage) loadAccount(address crypto.Address) (*Account, error) {
	raw, err := state.stateTrie.Get(address[:])
	if err != nil {
		return nil, err
	}
	var account Account
	if len(raw) == 0 {
		return nil, nil
	}
	if err := rlp.DecodeBytes(raw, &account); err != nil {
		return nil, err
	}
	account.address = address
	account.contract = state.Get(account.ContractHash.Bytes())
	if account.storage, err = trie.New(account.StorageHash, state.Database); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccount retrieve the account state at addr
func (state *StateStorage) GetAccount(address crypto.Address) (*Account, error) {
	loadedAccount, err := state.loadAccount(address)
	if err != nil {
		return nil, err
	}
	return loadedAccount, nil
}

// LoadAccount retrieve the account state at addr
func (state *StateStorage) LoadAccount(address crypto.Address) (*Account, error) {
	if state.accounts[address] == nil {
		loadedAccount, err := state.loadAccount(address)
		if err != nil {
			return nil, err
		}
		state.accounts[address] = loadedAccount
	}
	return state.accounts[address], nil
}

// CreateAccount create a new account state for addr
func (state *StateStorage) CreateAccount(creator crypto.Address, address crypto.Address, contract []byte) (*Account, error) {
	storage, err := trie.New(common.Hash{}, state.Database)
	if err != nil {
		return nil, err
	}
	account := &Account{
		Nonce:    0,
		Creator:  creator,
		address:  address,
		storage:  storage,
		contract: contract,
		dirty:    true,
	}

	account.setContract(contract)
	state.accounts[address] = account
	return account, nil
}

// GetStorage get the value at key of storage
func (account *Account) GetStorage(key []byte) ([]byte, error) {
	return account.storage.Get(key)
}

// SetStorage set the account storage
func (account *Account) SetStorage(key, value []byte) error {
	account.dirty = true
	return account.storage.Update(key, value)
}

// GetAddress returns state address
func (account *Account) GetAddress() crypto.Address {
	return account.address
}

// IsContract check whether this is an contract account or a normal account
func (account *Account) IsContract() bool {
	return account.ContractHash != common.EmptyHash
}

// GetContract retrieves contract code for account state
func (account *Account) GetContract() (*abi.Contract, error) {
	return abi.DecodeContract(account.contract)
}

// SetNonce stores the latest nonce to account state
func (account *Account) SetNonce(nonce uint64) {
	account.dirty = true
	account.Nonce = nonce
}

// GetCreator contract creator
func (account *Account) GetCreator() crypto.Address {
	return account.Creator
}

func (account *Account) setContract(contract []byte) {
	account.dirty = true
	account.contract = contract
	if len(account.contract) > 0 {
		account.ContractHash = blake2b.Sum256(contract)
	}
}
