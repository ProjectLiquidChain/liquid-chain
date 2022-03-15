package db

import (
	"bytes"
	"os"
	"testing"
)

var testVector = []struct {
	key   string
	value string
}{
	{"hello", "world"},
	{"dang", "nguyen"},
	{"block", "chain"},
	{"merkle", "tree"},
}

func TestDB(t *testing.T) {
	// Setup
	path := "./test-db"
	db := NewRocksDB(path)

	// Put
	for _, item := range testVector {
		db.Put([]byte(item.key), []byte(item.value))
	}

	// Get
	for _, item := range testVector {
		actual := db.Get([]byte(item.key))
		if !bytes.Equal(actual, []byte(item.value)) {
			t.Errorf("Value getting from db is different from expected. Expected: %v. Actual: %v", item.value, actual)
		}
	}

	// Tear down
	os.RemoveAll(path)
}

func TestMemoryDB(t *testing.T) {
	// Setup
	db := NewMemoryDB()

	// Put
	for _, item := range testVector {
		db.Put([]byte(item.key), []byte(item.value))
	}

	// Get
	for _, item := range testVector {
		actual := db.Get([]byte(item.key))
		if !bytes.Equal(actual, []byte(item.value)) {
			t.Errorf("Value getting from db is different from expected. Expected: %v. Actual: %v", item.value, actual)
		}
	}
}
