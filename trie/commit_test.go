package trie

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/google/uuid"
)

type TestNode struct {
	key   []byte
	value []byte
}

const nodeCount = 10000
const keyLength = 64
const valueLength = 128

var nodes []TestNode

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}

func init() {
	for i := 0; i < nodeCount; i++ {
		key := randomBytes(i%keyLength + 10)
		value := randomBytes(valueLength)
		nodes = append(nodes, TestNode{key, value})
	}
	nodes = append(nodes, TestNode{[]byte("doe"), []byte("reindeer")})
	nodes = append(nodes, TestNode{[]byte("dogglesworth"), []byte("cat")})
	nodes = append(nodes, TestNode{[]byte("dog"), []byte("puppy")})
	nodes = append(nodes, TestNode{[]byte("doge"), []byte("coin")})
}

func TestTrieWithDiskStorage(t *testing.T) {
	id, _ := uuid.NewUUID()
	path := fmt.Sprintf("./data-" + id.String())
	database := db.NewRocksDB(path)
	root := common.HexToHash("")
	tree, _ := New(root, database)
	tree.Hash()
	for i := 0; i < nodeCount; i++ {
		if err := tree.Update(nodes[i].key, nodes[i].value); err != nil {
			panic(err)
		}
	}
	tree.Hash()
	hash, _ := tree.Commit()

	newTree, _ := New(hash, database)
	for i := 0; i < nodeCount; i++ {
		if err := newTree.Update(nodes[i].key, nodes[i].value); err != nil {
			panic(err)
		}
	}
	for i := 0; i < nodeCount; i++ {
		v, _ := newTree.Get(nodes[i].key)
		if !bytes.Equal(v, nodes[i].value) {
			t.Error("Wrong data")
		}
	}

	newDeleteTree, _ := New(hash, database)
	for i := 0; i < nodeCount; i++ {
		if err := newDeleteTree.Update(nodes[i].key, []byte("")); err != nil {
			panic(err)
		}
	}

	os.RemoveAll(path)
}
