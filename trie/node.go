package trie

import (
	"errors"
	"fmt"
	"io"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
)

// Node is the unit of trie
type Node interface {
	cache() (node hashNode, isDirty bool)
}

type (
	branchNode struct {
		Children [17]Node
		flags    nodeFlag
	}
	shortNode struct {
		Key   []byte
		Value Node
		flags nodeFlag
	}
	valueNode []byte
	hashNode  []byte
)

func (n *branchNode) copy() *branchNode { copy := *n; return &copy }
func (n *shortNode) copy() *shortNode   { copy := *n; return &copy }

// EncodeRLP encodes a full node into the consensus RLP format.
func (n *branchNode) EncodeRLP(w io.Writer) error {
	var nodes [17]Node
	for i, child := range &n.Children {
		if child != nil {
			nodes[i] = child
		} else {
			nodes[i] = valueNode(nil)
		}
	}
	return rlp.Encode(w, nodes)
}

// flags contains caching-related metadata about a node.
type nodeFlag struct {
	hash  []byte // cached hash of the node (may be nil)
	dirty bool
}

func (n *branchNode) cache() (hashNode, bool) { return n.flags.hash, n.flags.dirty }
func (n *shortNode) cache() (hashNode, bool)  { return n.flags.hash, n.flags.dirty }
func (n hashNode) cache() (hashNode, bool)    { return nil, true }
func (n valueNode) cache() (hashNode, bool)   { return nil, true }

func mustDecodeNode(hash, buf []byte) Node {
	if s, _, err := rlp.SplitString(buf); err == nil && len(s) == 0 {
		// Handle nil node
		return nil
	}
	node, err := decodeNode(hash, buf)
	if err != nil {
		panic(fmt.Sprintf("node %x: %v", hash, err))
	}
	return node
}

// decodeNode parses the RLP encoding of a trie node.
func decodeNode(hash, buf []byte) (Node, error) {
	if len(buf) == 0 {
		return nil, errors.New("Unexpected end of buffer")
	}
	elements, _, err := rlp.SplitList(buf)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	switch count, _ := rlp.CountValues(elements); count {
	case 0:
		return valueNode(nil), nil
	case 2:
		node, err := decodeShortNode(hash, elements)
		return node, err
	case 17:
		node, err := decodeBranchNode(hash, elements)
		return node, err
	default:
		return nil, fmt.Errorf("Node elements count invalid: %v", count)
	}
}

func decodeShortNode(hash, elements []byte) (*shortNode, error) {
	keyByte, rest, err := rlp.SplitString(elements)
	if err != nil {
		return nil, err
	}
	key := compactToHex(keyByte)
	flag := nodeFlag{hash: hash}

	// Leaf node
	if hasTerm(key) {
		value, _, err := rlp.SplitString(rest)
		if err != nil {
			return nil, fmt.Errorf("invalid value node: %v", err)
		}
		return &shortNode{
			Key:   key,
			Value: append(valueNode{}, value...),
			flags: flag,
		}, nil
	}

	// Extension node
	node, _, err := decodeRef(rest)
	if err != nil {
		return nil, err
	}
	return &shortNode{
		Key:   key,
		Value: node,
		flags: flag,
	}, nil
}

func decodeBranchNode(hash, elements []byte) (*branchNode, error) {
	node := &branchNode{flags: nodeFlag{hash: hash}}
	for i := 0; i < 16; i++ {
		child, rest, err := decodeRef(elements)
		if err != nil {
			return node, err
		}
		node.Children[i] = child
		elements = rest
	}
	value, _, err := rlp.SplitString(elements)
	if err != nil {
		return node, err
	}
	if len(value) > 0 {
		node.Children[16] = append(valueNode{}, value...)
	}
	return node, nil
}

func decodeRef(buf []byte) (Node, []byte, error) {
	kind, value, rest, err := rlp.Split(buf)
	if err != nil {
		return nil, buf, err
	}
	switch {
	case kind == rlp.List:
		// 'embedded' node reference. The encoding must be smaller
		// than a hash in order to be valid.
		size := len(buf) - len(rest)

		if size > common.HashLength {
			err := fmt.Errorf("oversized embedded node (size is %d bytes, want size < %d)", size, common.HashLength)
			return nil, buf, err
		}
		n, err := decodeNode(nil, buf)
		return n, rest, err
	case kind == rlp.String && len(value) == 0:
		// empty node
		return nil, rest, nil
	case kind == rlp.String && len(value) == 32:
		return append(hashNode{}, value...), rest, nil
	default:
		return nil, nil, fmt.Errorf("invalid RLP string size %d (want 0 or 32)", len(value))
	}
}
