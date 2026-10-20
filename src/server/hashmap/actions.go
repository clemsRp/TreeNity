package hashmap

import (
	"fmt"
	"os"
)

func (h *Hashmap) Append(key string, value any) {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get hash
	hash := CalculateHash(key, len(h.Nodes))

	// Handle first element
	if h.Nodes[hash] == nil {
		h.Nodes[hash] = &Node{Key: key, Value: value, Next: nil}
		h.Length++
		return
	}

	// Append datas at hash's index
	h.Length += h.Nodes[hash].Append(key, value)
}

func (h *Hashmap) Get(key string) (any, error) {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get hash
	hash := CalculateHash(key, len(h.Nodes))

	if h.Nodes[hash] == nil {
		return nil, fmt.Errorf("'%s' key not found", key)
	}

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
		h.Length--
		return nil
	}

	// Delete datas at hash's index
	diff, err := h.Nodes[hash].Delete(key)

	h.Length -= diff

	return err
}

func (h *Hashmap) Len() int {
	// Handle Thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.Length
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

func (h *Hashmap) Print() {
	fmt.Fprintf(os.Stderr, "\n\nSize: %d\n", h.Size)
	fmt.Fprintf(os.Stderr, "Length: %d\n", h.Length)

	for deep := range h.Size {
		fmt.Fprintf(os.Stderr, "\nCouche num %d\n", deep)
		curr := h.Nodes[deep]
		for curr != nil {
			fmt.Fprintf(os.Stderr, "Key: %s, Value: %s\n", curr.Key, curr.Value)
			curr = curr.Next
		}
	}
}
