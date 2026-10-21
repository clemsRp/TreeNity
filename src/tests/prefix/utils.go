package prefix

import (
	"slices"
	"testing"
	pr "treenity/src/protocol"
	p "treenity/src/server/prefix"
	sce "treenity/src/tests/scenarios"
)

func UtilTestAction(t *testing.T, myTree *p.Tree, step sce.PrefixStep) {
	switch step.Action {
	case sce.PrefixActionInsert:
		err := myTree.Insert(step.Client)
		if (err != nil) != step.WantErr {
			t.Errorf("Insert() error = %v, WantErr %v", err, step.WantErr)
		}

	case sce.PrefixActionRemove:
		err := myTree.Remove(step.Client)
		if (err != nil) != step.WantErr {
			t.Errorf("Remove() error = %v, WantErr %v", err, step.WantErr)
		}

	case sce.PrefixActionFilter:
		clients := myTree.Filter(step.Key)
		if !SameClients(clients, step.WantClients) {
			t.Errorf("Filter() with '%s' key, got = %v, want %v", step.Key, clients, step.WantClients)
		}
	}
}

func SameClients(a, b []*pr.Client) bool {
	// Check lengths
	if len(a) != len(b) {
		return false
	}

	// Compare pointers
	for _, client_a := range a {
		if !slices.Contains(b, client_a) {
			return false
		}
	}

	return true
}
