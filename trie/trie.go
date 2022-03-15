package trie

import (
	"bytes"
	"fmt"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/db"
)

// Trie is Merkle Patricia Trie
type Trie struct {
	db   db.Database
	root Node
}

var emptyRoot = common.HexToHash("45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0")

// New returns a Trie based
func New(rootHash common.Hash, db db.Database) (*Trie, error) {
	if db == nil {
		panic("Could not run trie.New without db.")
	}
	trie := &Trie{db: db}
	if rootHash != common.EmptyHash {
		rootNode, err := trie.loadNode(rootHash.Bytes())
		if err != nil {
			return nil, err
		}
		trie.root = rootNode
	}
	return trie, nil
}

// loadNode loads the node from database
func (tree *Trie) loadNode(node hashNode) (Node, error) {
	hash := common.BytesToHash(node)
	data := tree.db.Get(hash.Bytes())
	if data == nil {
		return nil, fmt.Errorf("Missing node data for node %s", string(hash.Bytes()))
	}
	return mustDecodeNode(hash.Bytes(), data), nil
}

func (tree *Trie) newFlag() nodeFlag { return nodeFlag{dirty: true} }

func (tree *Trie) insertToShortNode(node *shortNode, key []byte, value Node) (bool, Node, error) {
	matchedLength := computeCommonPrefixLength(key, node.Key)

	// Case 1: Match totally -> merged with current node
	if matchedLength == len(node.Key) {
		dirty, newNode, err := tree.insert(node.Value, key[matchedLength:], value)
		if !dirty { // If this update doesn't change anything
			return false, node, nil
		}
		if err != nil { // If this update failed
			return false, node, err
		}
		return true, &shortNode{node.Key, newNode, tree.newFlag()}, nil
	}

	// Case 2: Match partially -> branch out at different pos
	branch := &branchNode{flags: tree.newFlag()}
	var err error

	// Add the current node value into one branch
	_, branch.Children[node.Key[matchedLength]], err = tree.insert(nil, node.Key[matchedLength+1:], node.Value)
	if err != nil {
		return false, nil, err
	}

	// Insert new node into one other branch
	_, branch.Children[key[matchedLength]], err = tree.insert(nil, key[matchedLength+1:], value)
	if err != nil {
		return false, nil, err
	}

	// Replace this shortNode with the branch if diff occurs at index 0.
	if matchedLength == 0 {
		return true, branch, nil
	}

	// Otherwise, replace it with a short node leading up to the branch.
	return true, &shortNode{key[:matchedLength], branch, tree.newFlag()}, nil
}

func (tree *Trie) insertToBranchNode(node *branchNode, key []byte, value Node) (bool, Node, error) {
	dirty, newNode, err := tree.insert(node.Children[key[0]], key[1:], value)
	if !dirty || err != nil {
		return false, node, err
	}
	node = node.copy()
	node.flags = tree.newFlag()
	node.Children[key[0]] = newNode
	return true, node, nil
}

func (tree *Trie) insertToHashNode(node hashNode, key []byte, value Node) (bool, Node, error) {
	currentNode, err := tree.loadNode(node)
	if err != nil {
		return false, nil, err
	}
	dirty, newNode, err := tree.insert(currentNode, key, value)
	if !dirty {
		return false, currentNode, nil
	}
	if err != nil {
		return false, currentNode, err
	}
	return true, newNode, nil
}

func (tree *Trie) insert(node Node, key []byte, value Node) (bool, Node, error) {
	if len(key) == 0 {
		if v, ok := node.(valueNode); ok {
			return !bytes.Equal(v, value.(valueNode)), value, nil
		}
		return true, value, nil
	}
	switch node := node.(type) {
	case *shortNode:
		return tree.insertToShortNode(node, key, value)
	case *branchNode:
		return tree.insertToBranchNode(node, key, value)
	case hashNode:
		return tree.insertToHashNode(node, key, value)
	case nil:
		return true, &shortNode{
			Key:   key,
			Value: value,
			flags: tree.newFlag(),
		}, nil
	default:
		panic(fmt.Sprintf("%T: invalid node: %v", node, node))
	}
}

// Update will either insert or delete a key based on value
func (tree *Trie) Update(key, value []byte) error {
	hexKey := keybytesToHex(key)
	if len(value) > 0 {
		_, newRoot, err := tree.insert(tree.root, hexKey, append(valueNode{}, value...))
		if err != nil {
			return err
		}
		tree.root = newRoot
	} else {
		_, newRoot, err := tree.delete(tree.root, hexKey)
		if err != nil {
			return err
		}
		tree.root = newRoot
	}
	return nil
}

// Hash returns the root hash
func (tree *Trie) Hash() common.Hash {
	hash, cached, _ := tree.hashRoot(nil)
	tree.root = cached
	return common.BytesToHash(hash.(hashNode))
}

// Commit returns the root hash and write to disk db
func (tree *Trie) Commit() (common.Hash, error) {
	hash, cached, err := tree.hashRoot(tree.db)
	if err != nil {
		return common.Hash{}, err
	}
	tree.root = cached
	return common.BytesToHash(hash.(hashNode)), nil
}

func (tree *Trie) hashRoot(db db.Database) (Node, Node, error) {
	hasher := newHasher()
	defer returnHasherToPool(hasher)
	if tree.root == nil {
		// Could not cache nil root
		hash, _, err := hasher.hash(valueNode(nil), db, true)
		return hash, nil, err
	}
	return hasher.hash(tree.root, db, true)
}

// Get will retrieve the value of key in tree
// It will also update the root node if its path to leaf node requires loading from db
func (tree *Trie) Get(key []byte) ([]byte, error) {
	key = keybytesToHex(key)
	value, newRoot, reachedHashNode, err := tree.get(tree.root, key, 0)
	if err == nil && reachedHashNode {
		tree.root = newRoot
	}
	return value, err
}

func (tree *Trie) get(currentNode Node, key []byte, position int) (value []byte, newNode Node, reachedHashNode bool, err error) {
	switch node := currentNode.(type) {
	case *shortNode:
		if len(key)-position < len(node.Key) || !bytes.Equal(node.Key, key[position:position+len(node.Key)]) {
			// Given key mismatched node
			return nil, node, false, nil
		}

		value, newNode, reachedHashNode, err = tree.get(node.Value, key, position+len(node.Key))
		if err == nil && reachedHashNode {
			node = node.copy()
			node.Value = newNode
		}
		return value, node, reachedHashNode, err

	case *branchNode:
		value, newNode, reachedHashNode, err = tree.get(node.Children[key[position]], key, position+1)
		if err == nil && reachedHashNode {
			node = node.copy()
			node.Children[key[position]] = newNode
		}
		return value, node, reachedHashNode, err

	case nil:
		return nil, nil, false, nil

	case valueNode:
		return node, node, false, nil

	case hashNode:
		loadedNode, err := tree.loadNode(node)
		if err != nil {
			return nil, node, true, err
		}
		value, newNode, _, err := tree.get(loadedNode, key, position)
		return value, newNode, true, err

	default:
		panic(fmt.Sprintf("%T: invalid node: %v", currentNode, currentNode))
	}
}

func (tree *Trie) resolve(node Node) (Node, error) {
	if node, ok := node.(hashNode); ok {
		return tree.loadNode(node)
	}
	return node, nil
}

func (tree *Trie) delete(node Node, key []byte) (bool, Node, error) {
	switch node := node.(type) {
	case *shortNode:
		matchedLength := computeCommonPrefixLength(key, node.Key)
		if matchedLength < len(node.Key) {
			return false, node, nil
		}
		if matchedLength == len(key) {
			return true, nil, nil
		}
		if dirty, childNode, err := tree.delete(node.Value, key[len(node.Key):]); err == nil && dirty {
			switch childNode := childNode.(type) {
			case *shortNode:
				// Merge the child node with current node
				return true, &shortNode{
					Key:   append(node.Key, childNode.Key...),
					Value: childNode.Value,
					flags: tree.newFlag(),
				}, nil
			default:
				// Child node is value of current node
				return true, &shortNode{
					Key:   node.Key,
					Value: childNode,
					flags: tree.newFlag(),
				}, nil
			}
		} else {
			return false, node, err
		}

	case *branchNode:
		dirty, newNode, err := tree.delete(node.Children[key[0]], key[1:])
		if !dirty || err != nil {
			return false, node, err
		}
		node = node.copy()
		node.flags = tree.newFlag()
		node.Children[key[0]] = newNode
		pos := -1
		for i, child := range &node.Children {
			if child != nil {
				if pos == -1 {
					pos = i
				} else {
					pos = -2
					break
				}
			}
		}
		if pos >= 0 {
			if pos != 16 {
				childNode, err := tree.resolve(node.Children[pos])
				if err != nil {
					return false, nil, err
				}
				if childNode, ok := childNode.(*shortNode); ok {
					k := append([]byte{byte(pos)}, childNode.Key...)
					return true, &shortNode{k, childNode.Value, tree.newFlag()}, nil
				}
			}
			return true, &shortNode{[]byte{byte(pos)}, node.Children[pos], tree.newFlag()}, nil
		}
		return true, node, nil

	case valueNode:
		return true, nil, nil

	case nil:
		return false, nil, nil

	case hashNode:
		loadedNode, err := tree.loadNode(node)
		if err != nil {
			return false, nil, err
		}
		dirty, newNode, err := tree.delete(loadedNode, key)
		if !dirty || err != nil {
			return false, loadedNode, err
		}
		return true, newNode, nil

	default:
		panic(fmt.Sprintf("%T: invalid node: %v (%v)", node, node, key))
	}
}

// NodeIterator returns an iterator that returns nodes of the trie. Iteration starts at
// the key after the given start key.
func (tree *Trie) NodeIterator(start []byte) NodeIterator {
	return newNodeIterator(tree, start)
}
