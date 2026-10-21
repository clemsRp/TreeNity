package hashmap

import sce "treenity/src/tests/scenarios"

var EdgeCaseScenarios = []sce.HashmapScenario{
	{
		Name: "Contains on completely empty map",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionContains, Key: "", WantBool: false},
			{Action: sce.HashmapActionContains, Key: "random", WantBool: false},
			{Action: sce.HashmapActionLen, WantInt: 0},
		},
	},
	{
		Name: "Failed delete does not alter length",
		Steps: []sce.HashmapStep{
			{Action: sce.HashmapActionAppend, Key: "keep", Value: "safe"},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionDelete, Key: "ghost", WantErr: true},
			{Action: sce.HashmapActionLen, WantInt: 1},
			{Action: sce.HashmapActionContains, Key: "keep", WantBool: true},
		},
	},
}
