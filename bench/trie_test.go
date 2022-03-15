package bench

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"github.com/google/uuid"
)

type Node struct {
	key   []byte
	value []byte
}

const nodeCount = 1000000
const keyLength = 32
const valueLength = 128

var nodes []Node

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}

func init() {
	for i := 0; i < nodeCount; i++ {
		key := randomBytes(keyLength)
		value := randomBytes(valueLength)
		nodes = append(nodes, Node{key, value})
	}
}

func benchmarkInsert(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		db := db.NewMemoryDB()
		root := common.HexToHash("")
		b.ReportAllocs()
		tree, _ := trie.New(root, db)
		for j := 0; j < n; j++ {
			tree.Update(nodes[j].key, nodes[j].value)
		}
		tree.Commit()
	}
}

func benchmarkInsertDisk(n int, b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		id, _ := uuid.NewUUID()
		path := fmt.Sprintf("./data-"+id.String(), n, i)
		database := db.NewRocksDB(path)
		root := common.HexToHash("")
		tree, _ := trie.New(root, database)
		for j := 0; j < n; j++ {
			if err := tree.Update(nodes[j].key, nodes[j].value); err != nil {
				panic(err)
			}
		}
		_, err := tree.Commit()
		if err != nil {
			panic(err)
		}
		os.RemoveAll(path)
	}
}

func benchmarkGetDisk(n int, b *testing.B) {
	id, _ := uuid.NewUUID()
	path := fmt.Sprintf("./data-" + id.String())
	database := db.NewRocksDB(path)
	root := common.HexToHash("")
	tree, _ := trie.New(root, database)
	for i := 0; i < n; i++ {
		if err := tree.Update(nodes[i].key, nodes[i].value); err != nil {
			panic(err)
		}
	}
	hash, _ := tree.Commit()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		newTree, _ := trie.New(hash, database)
		// for j := 0; j < n; j++ {
		v, _ := newTree.Get(nodes[i%n].key)
		if !bytes.Equal(v, nodes[i%n].value) {
			b.Error("Wrong data")
		}
		// }
	}
	os.RemoveAll(path)
}

// Memory
func BenchmarkInsert1(b *testing.B)       { benchmarkInsert(1, b) }
func BenchmarkInsert100(b *testing.B)     { benchmarkInsert(100, b) }
func BenchmarkInsert10000(b *testing.B)   { benchmarkInsert(10000, b) }
func BenchmarkInsert1000000(b *testing.B) { benchmarkInsert(1000000, b) }

// Disk
func BenchmarkInsertDisk1(b *testing.B)      { benchmarkInsertDisk(1, b) }
func BenchmarkInsertDisk100(b *testing.B)    { benchmarkInsertDisk(100, b) }
func BenchmarkInsertDisk10000(b *testing.B)  { benchmarkInsertDisk(10000, b) }
func BenchmarkInsertDisk100000(b *testing.B) { benchmarkInsertDisk(10000, b) }

func BenchmarkGetDisk1(b *testing.B)      { benchmarkGetDisk(1, b) }
func BenchmarkGetDisk100(b *testing.B)    { benchmarkGetDisk(100, b) }
func BenchmarkGetDisk10000(b *testing.B)  { benchmarkGetDisk(10000, b) }
func BenchmarkGetDisk100000(b *testing.B) { benchmarkGetDisk(100000, b) }
