package hashmap

import sce "treenity/src/tests/scenarios"

var CollisionsScenarios = []sce.HashmapScenario{
	{
		Name: "Length and Contains with multiple collisions",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionAppend, Key: "hash_a", Value: 10},
			{Action: sce.ActionAppend, Key: "hash_b", Value: 20},
			{Action: sce.ActionLen, WantInt: 2},
			{Action: sce.ActionContains, Key: "hash_a", WantBool: true},
			{Action: sce.ActionContains, Key: "hash_b", WantBool: true},
			{Action: sce.ActionContains, Key: "hash_c", WantBool: false},
		},
	},
	{
		Name: "Length accuracy after deleting elements in collision chains",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionAppend, Key: "c1", Value: "val1"},
			{Action: sce.ActionAppend, Key: "c2", Value: "val2"},
			{Action: sce.ActionLen, WantInt: 2},
			// Supprime le premier de la collision
			{Action: sce.ActionDelete, Key: "c1", WantErr: false},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionContains, Key: "c1", WantBool: false},
			{Action: sce.ActionContains, Key: "c2", WantBool: true},
		},
	},
}
