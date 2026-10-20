package hashmap

import (
	"fmt"
	"sync"
	pr "treenity/src/protocol"
)

type Hashmap struct {
	Len   int
	Nodes []*Node
	mu    sync.RWMutex
}

func CreateHashmap(size int) *Hashmap {
	var nodes []*Node
	for _ := range size {
		nodes = append(nodes, nil)
	}

	return &Hashmap{
		Nodes: nodes,
	}
}

func (h *Hashmap) Append(key string, value any) {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get hash
	hash := CalculateHash(key, len(h.Nodes))

	// Handle first element
	if h.Nodes[hash] == nil {
		h.Nodes[hash] = &Node{Key: key, Value: value, Next: nil}
		h.Len++
		return
	}

	// Append datas at hash's index
	h.Len += h.Nodes[hash].Append(key, value)
}

func (h *Hashmap) Get(key string) (any, error) {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get hash
	hash := CalculateHash(key, len(h.Nodes))

	// Get datas at hash's index
	return h.Nodes[hash].Get(key)
}

func (h *Hashmap) Delete(key string) error {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get hash
	hash := CalculateHash(key, len(h.Nodes))

	if h.Nodes[hash] == nil {
		return fmt.Errorf("'%s' key not found", key)
	}

	// Handle last element
	if h.Nodes[hash].Key == key {
		h.Nodes[hash] = h.Nodes[hash].Next
		h.Len--
		return nil
	}

	// Delete datas at hash's index
	diff, err := h.Nodes[hash].Delete(key)

	h.Len -= diff

	return err
}

func (h *Hashmap) Len() int {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.Len
}

func (h *Hashmap) Contains(key string) bool {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get hash
	hash := CalculateHash(key, len(h.Nodes))

	if h.Nodes[hash] == nil {
		return false
	}

	return h.Nodes[hash].Contains(key)
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
		hash ^= data_byte
		hash *= fnv_prime
	}

	// Limit hash to hashmap length
	hash %= length

	return hash
}
