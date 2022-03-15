package trie

import (
	"bytes"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/db"
)

func newEmpty() *Trie {
	db := db.NewMemoryDB()
	trie, _ := New(common.Hash{}, db)
	return trie
}

func updateString(trie *Trie, k, v string) {
	err := trie.Update([]byte(k), []byte(v))
	if err != nil {
		panic(err)
	}
}

func getString(trie *Trie, k string) []byte {
	value, _ := trie.Get([]byte(k))
	return value
}

func deleteString(trie *Trie, k string) {
	err := trie.Update([]byte(k), nil)
	if err != nil {
		panic(err)
	}
}

func TestNull(t *testing.T) {
	trie := newEmpty()
	key := make([]byte, 32)
	value := []byte("test")
	err := trie.Update(key, value)
	if err != nil {
		panic(err)
	}
	storedValue, _ := trie.Get(key)
	if !bytes.Equal(storedValue, value) {
		t.Error("wrong value")
	}
}

func TestMutable(t *testing.T) {
	trie := newEmpty()
	key := []byte{1, 2}
	value := []byte{1, 2}
	err := trie.Update(key, value)
	if err != nil {
		panic(err)
	}

	// Mutable key and value
	key[0] = 2
	value[0] = 2

	v, err := trie.Get([]byte{1, 2})
	if err != nil {
		t.Error(err.Error())
	}
	expectedValue := []byte{1, 2}
	if !bytes.Equal(v, expectedValue) {
		t.Errorf("Expected value: %v, actual value: %v", expectedValue, v)
	}
}

func TestInsert(t *testing.T) {
	trie := newEmpty()

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	exp := common.HexToHash("6ca394ff9b13d6690a51dea30b1b5c43108e52944d30b9095227c49bae03ff8b")
	root := trie.Hash()
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}

	trie = newEmpty()
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	exp = common.HexToHash("e9d7f23f40cd82fe35f5a7a6778c3503f775f3623ba7a71fb335f0eee29dac8a")
	root, err := trie.Commit()
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}
	if err != nil {
		t.Errorf("expected nil got %x", err)
	}
}

func TestGet(t *testing.T) {
	trie := newEmpty()
	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	for i := 0; i < 2; i++ {

		if res := getString(trie, "dog"); !bytes.Equal(res, []byte("puppy")) {
			t.Errorf("expected puppy got %x", res)
		}

		if res := getString(trie, "doe"); !bytes.Equal(res, []byte("reindeer")) {
			t.Errorf("expected reindeer got %x", res)
		}

		if res := getString(trie, "dogglesworth"); !bytes.Equal(res, []byte("cat")) {
			t.Errorf("expected cat got %x", res)
		}

		unknown := getString(trie, "unknown")
		if unknown != nil {
			t.Errorf("expected nil got %x", unknown)
		}

		if i == 1 {
			return
		}
		trie.Commit()
	}
}

func TestDelete(t *testing.T) {
	trie := newEmpty()
	vals := []struct{ k, v string }{
		{"doge", "coin"},
		{"dog", "puppy"},
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"ether", ""},
		{"shaman", ""},
		{"doge", ""},
		{"dog", ""},
		{"do", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			updateString(trie, val.k, val.v)
			val.v = "123"
		} else {
			deleteString(trie, val.k)
		}
	}

	hash := trie.Hash()
	exp := common.HexToHash("c75904529cfd1d0d709c8c059858951bdc93663df88f6385dba042c011087111")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestEmptyValues(t *testing.T) {
	trie := newEmpty()

	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		updateString(trie, val.k, val.v)
	}

	hash := trie.Hash()
	exp := common.HexToHash("79a9b42da0e261b9f3ca9e78560ac8d486bcce2da8a5ddb2df8721d4c0dc2d0a")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}
