package prefix

import (
	pr "treenity/src/protocol"
)

type Tree struct {
	Char     rune
	Children map[rune]*Tree
	Clients  []*pr.Client
}

func CreateTree(char rune) *Tree {
	return &Tree{
		Char:     char,
		Children: make(map[rune]*Tree),
		Clients:  make([]*pr.Client, 0),
	}
}
