package prefix

import (
	"fmt"
	"slices"
	pr "treenity/src/protocol"
)

func (t *Tree) Insert(client *pr.Client) error {
	return t.InsertClient(client, []rune(client.Prefix))
}

func (t *Tree) InsertClient(client *pr.Client, prefix []rune) error {
	// Insert client
	if len(prefix) == 0 {
		// Handle duplication
		if slices.Contains(t.Clients, client) {
			return fmt.Errorf("'%s' client is already insert inside Tree", client.Id)
		}

		// Insert client inside Tree
		t.Clients = append(t.Clients, client)

		return nil
	}

	// Create tree if needed
	next_char := prefix[0]
	next_tree, ok := t.Children[next_char]
	if !ok {
		next_tree = CreateTree(next_char)
		t.Children[next_char] = next_tree
	}

	return next_tree.InsertClient(client, prefix[1:])
}

func (t *Tree) Remove(client *pr.Client) error {
	return t.RemoveClient(client, []rune(client.Prefix))
}

func (t *Tree) RemoveClient(client *pr.Client, prefix []rune) error {
	// Remove client
	if len(prefix) == 0 {
		// Handle missing client
		index := slices.Index(t.Clients, client)
		if index == -1 {
			return fmt.Errorf("'%s' client is not insert inside Tree", client.Id)
		}

		// Remove client from Tree
		t.Clients = append(t.Clients[:index], t.Clients[index+1:]...)

		return nil
	}

	// Check next tree exist
	next_char := prefix[0]
	next_tree, ok := t.Children[next_char]
	if !ok {
		return fmt.Errorf("'%s' client is not insert inside Tree, imcomplete path", client.Id)
	}

	// Remove recursively
	err := next_tree.RemoveClient(client, prefix[1:])
	if err != nil {
		return err
	}

	// Clean empty trees
	if len(next_tree.Clients) == 0 && len(next_tree.Children) == 0 {
		delete(t.Children, next_char)
	}

	return nil
}

func (t *Tree) Filter(key string) []*pr.Client {
	return t.FilterClient([]rune(key))
}

func (t *Tree) FilterClient(key []rune) []*pr.Client {
	// Get root clients
	filtered_clients := make([]*pr.Client, 0)
	filtered_clients = append(filtered_clients, t.Clients...)

	if len(key) > 0 {
		// Check sub path exist
		sub_tree, exist := t.Children[key[0]]

		// Get sub tree clients
		if exist {
			sub_filtered_clients := sub_tree.FilterClient(key[1:])
			filtered_clients = append(filtered_clients, sub_filtered_clients...)
		}
	}

	return filtered_clients
}
