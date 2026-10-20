package hashmap

import "fmt"

type Node struct {
	Key   string
	Value any
	Next  *Node
}

func (n *Node) Append(key string, value any) int {
	if n.Key == key {
		n.Value = value
		return 0
	}

	if n.Next != nil {
		return n.Next.Append(key, value)
	}

	n.Next = &Node{
		Key:   key,
		Value: value,
		Next:  nil,
	}

	return 1
}

func (n *Node) Get(key string) (any, error) {
	if n.Key == key {
		return n.Value, nil

	} else if n.Next != nil {
		return n.Next.Get(key)
	}

	return nil, fmt.Errorf("'%s' key not found", key)
}

func (n *Node) Delete(key string) (int, error) {
	if n.Next == nil {
		return 0, fmt.Errorf("'%s' key not found", key)
	}

	if n.Next.Key == key {
		n.Next = n.Next.Next
		return 1, nil
	}

	return n.Next.Delete(key)
}

func (n *Node) Contains(key string) bool {
	if n.Key == key {
		return true
	}

	if n.Next == nil {
		return false
	}

	return n.Next.Contains(key)
}
