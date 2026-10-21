package hashmap

import sce "treenity/src/tests/scenarios"

var BasicScenarios = []sce.HashmapScenario{
	{
		Name: "Check initial length and empty contains",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionLen, WantInt: 0},
			{Action: sce.HashmapActionContains, Key: "missing", WantBool: false},
		},
	},
	{
		Name: "Length increments on insert and checks contains",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionLen, WantInt: 0},
			{Action: sce.HashmapActionAppend, Key: "item1", Value: "val1"},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionContains, Key: "item1", WantBool: true},
			{Action: sce.HashmapActionAppend, Key: "item2", Value: "val2"},
			{Action: sce.HashmapActionLen, WantInt: 2},
			{Action: sce.HashmapActionContains, Key: "item2", WantBool: true},
		},
	},
	{
		Name: "Updating existing key does not change length",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionAppend, Key: "k", Value: "v1"},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionAppend, Key: "k", Value: "v2"},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionGet, Key: "k", Want: "v2", WantErr: false},
		},
	},
	{
		Name: "Length decrements and contains returns false on delete",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionAppend, Key: "to_delete", Value: "data"},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionDelete, Key: "to_delete", WantErr: false},
			{Action: sce.HashmapActionLen, WantInt: 0},
			{Action: sce.HashmapActionContains, Key: "to_delete", WantBool: false},
		},
	},
}
