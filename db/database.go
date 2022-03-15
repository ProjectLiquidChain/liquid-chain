package db

// Database generics inteface
type Database interface {
	Get(key []byte) []byte
	Put(key []byte, value []byte)
}
