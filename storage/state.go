package storage

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/trie"
)

// StateStorage is the global account state consisting of many address->state mapping
type StateStorage struct {
	db.Database
	block             *crypto.Block
	stateTrie         *trie.Trie
	accounts          map[crypto.Address]*Account
	accountCheckpoint common.Hash
}

// NewStateStorage returns a state storage
func NewStateStorage(db db.Database) *StateStorage {
	return &StateStorage{Database: db}
}

// MustLoadState do LoadState, but panic if error
func (state *StateStorage) MustLoadState(block *crypto.Block) {
	if err := state.LoadState(block); err != nil {
		panic(err)
	}
}

// LoadState load state root of block into trie
func (state *StateStorage) LoadState(block *crypto.Block) error {
	stateTrie, err := trie.New(block.StateRoot, state.Database)
	if err != nil {
		return err
	}

	state.block = block
	state.stateTrie = stateTrie
	state.accountCheckpoint = block.StateRoot
	state.accounts = make(map[crypto.Address]*Account)

	return nil
}

// GetBlock returns block that inits current state
func (state *StateStorage) GetBlock() *crypto.Block {
	return state.block
}

// Hash retrives hash of entire state
func (state *StateStorage) Hash() common.Hash {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		// Update account storage
		account.StorageHash = account.storage.Hash()

		// Update account
		raw, _ := rlp.EncodeToBytes(account)
		if err = state.stateTrie.Update(account.address[:], raw); err != nil {
			panic(err)
		}
	}
	return state.stateTrie.Hash()
}

// Commit stores all dirty Accounts to storage.trie
func (state *StateStorage) Commit() common.Hash {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		if account.IsContract() {
			// Update contract
			state.Put(account.ContractHash.Bytes(), account.contract)
		}

		// Update account storage
		if account.StorageHash, err = account.storage.Commit(); err != nil {
			panic(err)
		}

		// Update account
		raw, err := rlp.EncodeToBytes(account)
		if err != nil {
			panic(err)
		}

		if err := state.stateTrie.Update(account.address[:], raw); err != nil {
			panic(err)
		}

		account.dirty = false
	}

	stateRootHash, err := state.stateTrie.Commit()
	if err != nil {
		panic(err)
	}

	state.accountCheckpoint = stateRootHash
	return stateRootHash
}

// Revert state to last checkpoint
func (state *StateStorage) Revert() {
	t, err := trie.New(state.accountCheckpoint, state.Database)
	if err != nil {
		panic(err)
	}
	state.stateTrie = t
	state.accounts = make(map[crypto.Address]*Account)
}
