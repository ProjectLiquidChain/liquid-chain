package trie

import (
	"hash"
	"sync"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/db"
	"golang.org/x/crypto/blake2b"
)

const (
	hashLength = 32
)

type hasher struct {
	buf   buffer
	sha   state
	blake hash.Hash
}

// state wraps sha3.state
type state interface {
	hash.Hash
	Read([]byte) (int, error)
}

type buffer []byte

var emptyBuffer = make(buffer, 0, 17*32) // cap is as large as a full fullNode.

func (b *buffer) Write(data []byte) (n int, err error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *buffer) Reset() { *b = emptyBuffer }

// hashers live in a global db.
var hasherPool = sync.Pool{
	New: func() interface{} {
		b, _ := blake2b.New256([]byte{})
		return &hasher{
			buf:   emptyBuffer,
			blake: b,
		}
	},
}

func newHasher() *hasher { return hasherPool.Get().(*hasher) }

func returnHasherToPool(h *hasher) { hasherPool.Put(h) }

func (h *hasher) hash(node Node, db db.Database, force bool) (Node, Node, error) {
	if hash, dirty := node.cache(); hash != nil {
		if db == nil {
			return hash, node, nil
		}
		if !dirty {
			switch node.(type) {
			case *branchNode, *shortNode:
				return hash, hash, nil
			default:
				return hash, node, nil
			}
		}
	}

	// Trie not processed yet or needs storage, walk the children
	collapsed, cached, err := h.hashChildren(node, db)
	if err != nil {
		return hashNode{}, node, err
	}
	hashed, err := h.store(collapsed, db, force)
	if err != nil {
		return hashNode{}, node, err
	}

	// Cache the hash of the node for later reuse and remove
	// the dirty flag in commit mode. It's fine to assign these values directly
	// without copying the node first because hashChildren copies it.
	cachedHash, _ := hashed.(hashNode)
	switch cachedNode := cached.(type) {
	case *shortNode:
		cachedNode.flags.hash = cachedHash
		if db != nil {
			cachedNode.flags.dirty = false
		}
	case *branchNode:
		cachedNode.flags.hash = cachedHash
		if db != nil {
			cachedNode.flags.dirty = false
		}
	}
	return hashed, cached, nil
}

// hashChildren replaces the children of a node with their hashes if the encoded
// size of the child is larger than a hash, returning the collapsed node as well
// as a replacement for the original node with the child hashes cached in.
func (h *hasher) hashChildren(original Node, db db.Database) (Node, Node, error) {
	var err error

	switch node := original.(type) {
	case *shortNode:
		// Hash the short node's child, caching the newly hashed subtree
		collapsed, cached := node.copy(), node.copy()
		collapsed.Key = hexToCompact(node.Key)
		cached.Key = node.Key
		if _, ok := node.Value.(valueNode); !ok {
			collapsed.Value, cached.Value, err = h.hash(node.Value, db, false)
			if err != nil {
				return original, original, err
			}
		}
		return collapsed, cached, nil

	case *branchNode:
		// Hash the full node's children, caching the newly hashed subtrees
		collapsed, cached := node.copy(), node.copy()
		for i := 0; i < 16; i++ {
			if node.Children[i] == nil {
				continue
			}
			collapsed.Children[i], cached.Children[i], err = h.hash(node.Children[i], db, false)
			if err != nil {
				return original, original, err
			}
		}
		cached.Children[16] = node.Children[16]
		return collapsed, cached, nil

	default:
		// Value and hash nodes don't have children so they're left as were
		return node, original, nil
	}
}

// store hashes the node n and if we have a storage layer specified, it writes
// the key/value pair to it and tracks any node->child references as well as any
// node->external trie references.
func (h *hasher) store(node Node, db db.Database, force bool) (Node, error) {

	// Don't store hashes or empty nodes.
	if _, isHash := node.(hashNode); node == nil || isHash {
		return node, nil
	}

	// Generate the RLP encoding of the node
	h.buf.Reset()
	if err := rlp.Encode(&h.buf, node); err != nil {
		panic("encode error: " + err.Error())
	}
	if len(h.buf) < hashLength && !force {
		return node, nil // Nodes smaller than 32 bytes are stored inside their parent
	}

	// Larger nodes are replaced by their hash and stored in the database.
	hash, _ := node.cache()
	if hash == nil {
		var err error
		hash, err = h.makeHashNode()
		if err != nil {
			return nil, err
		}
	}

	// Store node
	if db != nil {
		db.Put(hash, h.buf)
	}

	return hash, nil
}

func (h *hasher) makeHashNode() (hashNode, error) {
	h.blake.Reset()
	_, err := h.blake.Write(h.buf)
	if err != nil {
		return nil, err
	}
	return h.blake.Sum([]byte{}), err
}
