package hashmap

import (
	"sync"
	pr "treenity/src/protocol"
)

type Hashmap struct {
	Size   int
	Length int
	Nodes  []*Node
	mu     sync.RWMutex
}

func CreateHashmap(size int) *Hashmap {
	nodes := make([]*Node, size)

	return &Hashmap{
		Nodes:  nodes,
		Size:   size,
		Length: 0,
	}
}

func CalculateHash(key string, length int) int {
	// Cast FNV-1a variables
	fnv_prime := uint64(pr.FnvPrime)
	fnv_offset_basis := uint64(pr.FnvOffsetBasis)

	// Get starting datas
	hash := fnv_offset_basis
	data_bytes := []byte(key)

	// Execute hash algorythm
	for _, data_byte := range data_bytes {
		hash ^= uint64(data_byte)
		hash *= fnv_prime
	}

	// Limit hash to hashmap length
	hash %= uint64(length)

	return int(hash)
}
