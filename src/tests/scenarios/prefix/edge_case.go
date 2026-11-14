package prefix

import (
	pr "treenity/src/protocol"
	sce "treenity/src/tests/scenarios"
)

var (
	client3 = &pr.Client{
		Id:     "client3",
		Prefix: "user.",
	}

	client4 = &pr.Client{
		Id:     "client4",
		Prefix: "",
	}

	client5 = &pr.Client{
		Id:     "client5",
		Prefix: "",
	}
)

var EdgeCaseScenarios = []sce.PrefixScenario{
	{
		Name: "Empty tree test",
		Steps: []sce.PrefixStep{
			// Test Remove
			{Action: sce.PrefixActionRemove, Client: client5, WantErr: true},

			// Test Filter
			{Action: sce.PrefixActionFilter, Key: "", WantClients: []*pr.Client{}},
		},
	},
	{
		Name: "Empty tree management test",
		Steps: []sce.PrefixStep{
			// Insert clients
			{Action: sce.PrefixActionInsert, Client: client3, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client4, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client5, WantErr: false},

			// Test management actions
			{Action: sce.PrefixActionFilter, Key: "user.", WantClients: []*pr.Client{client3, client4, client5}},
			{Action: sce.PrefixActionRemove, Client: client3, WantErr: false},
			{Action: sce.PrefixActionFilter, Key: "user.", WantClients: []*pr.Client{client4, client5}},
		},
	},
}
