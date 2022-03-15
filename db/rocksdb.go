package db

import (
	"github.com/linxGnu/grocksdb"
)

// RocksDB use map to store and retrieve value
type RocksDB struct {
	instance *grocksdb.DB
}

// NewRocksDB returns a new instance of the RocksDB
func NewRocksDB(path string) *RocksDB {
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))
	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	instance, err := grocksdb.OpenDb(opts, path)
	if err != nil {
		panic(err)
	}
	return &RocksDB{instance}
}

// Get returns the value based on key
func (db *RocksDB) Get(key []byte) []byte {
	ro := grocksdb.NewDefaultReadOptions()
	ro.SetFillCache(true)
	value, err := db.instance.Get(ro, key)
	if err != nil {
		panic(err)
	}
	return value.Data()
}

// Put inserts an key-value pair to database
func (db *RocksDB) Put(key []byte, value []byte) {
	wo := grocksdb.NewDefaultWriteOptions()
	wo.SetSync(false)
	if err := db.instance.Put(wo, key, value); err != nil {
		panic(err)
	}
}

func (db *RocksDB) Close() {
	db.instance.Close()
}
