package db

import (
	"encoding/hex"
)

// MemoryDB simple memory database
type MemoryDB struct {
	cache map[string][]byte
}

// NewMemoryDB return new in-memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		cache: make(map[string][]byte),
	}
}

// Get returns the value based on key
func (db *MemoryDB) Get(key []byte) []byte {
	return db.cache[hex.EncodeToString(key)]
}

// Put inserts an key-value pair to database
func (db *MemoryDB) Put(key []byte, value []byte) {
	db.cache[hex.EncodeToString(key)] = value
}
