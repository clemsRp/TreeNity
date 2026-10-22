package hashmap

import sce "treenity/src/tests/scenarios"

var CollisionsScenarios = []sce.HashmapScenario{
	{
		Name: "Length and Contains with multiple collisions",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionAppend, Key: "hash_a", Value: "10"},
			{Action: sce.HashmapActionAppend, Key: "hash_b", Value: "20"},
			{Action: sce.HashmapActionLen, WantInt: 2},
			{Action: sce.HashmapActionContains, Key: "hash_a", WantBool: true},
			{Action: sce.HashmapActionContains, Key: "hash_b", WantBool: true},
			{Action: sce.HashmapActionContains, Key: "hash_c", WantBool: false},
		},
	},
	{
		Name: "Length accuracy after deleting elements in collision chains",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionAppend, Key: "c1", Value: "val1"},
			{Action: sce.HashmapActionAppend, Key: "c2", Value: "val2"},
			{Action: sce.HashmapActionLen, WantInt: 2},
			// Supprime le premier de la collision
			{Action: sce.HashmapActionDelete, Key: "c1", WantErr: false},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionContains, Key: "c1", WantBool: false},
			{Action: sce.HashmapActionContains, Key: "c2", WantBool: true},
		},
	},
}
