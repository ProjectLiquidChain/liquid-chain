package storage

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/common"
)

const (
	blockHeightToBlockHashPrefix byte = 0x0
	txHashToBlockHeightPrefix    byte = 0x1
	latestBlockHeightPrefix      byte = 0x2
	txHashToReceiptHashPrefix    byte = 0x3
)

func (index *MetaStorage) encodeTxHashToReceiptHashKey(hash common.Hash) []byte {
	return index.encodeKey(txHashToReceiptHashPrefix, hash.Bytes())
}

func (index *MetaStorage) encodeTxHashToBlockHeightKey(hash common.Hash) []byte {
	return index.encodeKey(txHashToBlockHeightPrefix, hash.Bytes())
}

func (index *MetaStorage) encodeBlockHeightToBlockHashKey(height uint64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, height)
	return index.encodeKey(blockHeightToBlockHashPrefix, key)
}

func (index *MetaStorage) encodeLatestBlockHeightKey() []byte {
	return index.encodeKey(latestBlockHeightPrefix, []byte{})
}

func (index *MetaStorage) encodeKey(prefix byte, key []byte) []byte {
	return append([]byte{byte(prefix)}, key...)
}
