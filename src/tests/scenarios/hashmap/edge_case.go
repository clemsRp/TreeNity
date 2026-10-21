package hashmap

import sce "treenity/src/tests/scenarios"

var EdgeCaseScenarios = []sce.HashmapScenario{
	{
		Name: "Contains on completely empty map",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionContains, Key: "", WantBool: false},
			{Action: sce.ActionContains, Key: "random", WantBool: false},
			{Action: sce.ActionLen, WantInt: 0},
		},
	},
	{
		Name: "Failed delete does not alter length",
		Steps: []sce.HashmapStep{
			{Action: sce.ActionAppend, Key: "keep", Value: "safe"},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionDelete, Key: "ghost", WantErr: true},
			{Action: sce.ActionLen, WantInt: 1},
			{Action: sce.ActionContains, Key: "keep", WantBool: true},
		},
	},
}
