package prefix

import (
	pr "treenity/src/protocol"
	sce "treenity/src/tests/scenarios"
)

var (
	client0 = &pr.Client{
		Id:     "client0",
		Prefix: "user.",
	}

	client1 = &pr.Client{
		Id:     "client1",
		Prefix: "us",
	}

	client2 = &pr.Client{
		Id:     "client2",
		Prefix: "",
	}
)

var ActionsScenarios = []sce.PrefixScenario{
	{
		Name: "Basic remove management test",
		Steps: []sce.PrefixStep{
			// Insert clients
			{Action: sce.PrefixActionInsert, Client: client0, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client1, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client2, WantErr: false},

			// Test Remove
			{Action: sce.PrefixActionRemove, Client: client2, WantErr: false},
			{Action: sce.PrefixActionRemove, Client: client2, WantErr: true},
			{Action: sce.PrefixActionRemove, Client: client1, WantErr: false},
			{Action: sce.PrefixActionRemove, Client: nil, WantErr: true},
		},
	},
	{
		Name: "Basic filter management test",
		Steps: []sce.PrefixStep{
			// Insert clients
			{Action: sce.PrefixActionInsert, Client: client0, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client1, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client2, WantErr: false},

			// Test Filter
			{Action: sce.PrefixActionFilter, Key: "user.", WantClients: []*pr.Client{client0, client1, client2}},
			{Action: sce.PrefixActionFilter, Key: "user.create", WantClients: []*pr.Client{client0, client1, client2}},
			{Action: sce.PrefixActionFilter, Key: "ustencil", WantClients: []*pr.Client{client1, client2}},
			{Action: sce.PrefixActionFilter, Key: "us", WantClients: []*pr.Client{client1, client2}},
			{Action: sce.PrefixActionFilter, Key: "u", WantClients: []*pr.Client{client2}},
			{Action: sce.PrefixActionFilter, Key: "", WantClients: []*pr.Client{client2}},
			{Action: sce.PrefixActionFilter, Key: "nothing", WantClients: []*pr.Client{client2}},
		},
	},
	{
		Name: "Basic insert management test",
		Steps: []sce.PrefixStep{
			// Test Insertion
			{Action: sce.PrefixActionInsert, Client: client0, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client1, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client2, WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client2, CustomPrefix: "salut", WantErr: false},
			{Action: sce.PrefixActionInsert, Client: client2, WantErr: true},
			{Action: sce.PrefixActionInsert, Client: nil, WantErr: true},
		},
	},
}
