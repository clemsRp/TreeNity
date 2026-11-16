package prefix

import (
	"sync"
	pr "treenity/src/protocol"
)

type Tree struct {
	Char       byte
	Children   map[byte]*Tree
	Clients    []*pr.Client
	AllClients []*pr.Client
	mu         sync.RWMutex
}

func CreateTree(char byte) *Tree {
	return &Tree{
		Char:       char,
		Children:   make(map[byte]*Tree),
		Clients:    make([]*pr.Client, 0),
		AllClients: make([]*pr.Client, 0),
	}
}
