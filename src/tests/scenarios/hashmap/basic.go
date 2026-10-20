package hashmap

import sce "treenity/src/tests/scenarios"

var BasicScenarios = []sce.HashmapScenario{
	{
		Name: "Check initial length and empty contains",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionLen, WantInt: 0},
			{Action: sce.ActionContains, Key: "missing", WantBool: false},
		},
	},
	{
		Name: "Length increments on insert and checks contains",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionLen, WantInt: 0},
			{Action: sce.ActionAppend, Key: "item1", Value: "val1"},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionContains, Key: "item1", WantBool: true},
			{Action: sce.ActionAppend, Key: "item2", Value: "val2"},
			{Action: sce.ActionLen, WantInt: 2},
			{Action: sce.ActionContains, Key: "item2", WantBool: true},
		},
	},
	{
		Name: "Updating existing key does not change length",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionAppend, Key: "k", Value: "v1"},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionAppend, Key: "k", Value: "v2"},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionGet, Key: "k", Want: "v2", WantErr: false},
		},
	},
	{
		Name: "Length decrements and contains returns false on delete",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionAppend, Key: "to_delete", Value: "data"},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionDelete, Key: "to_delete", WantErr: false},
			{Action: sce.ActionLen, WantInt: 0},
			{Action: sce.ActionContains, Key: "to_delete", WantBool: false},
		},
	},
}
