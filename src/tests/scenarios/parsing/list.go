package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingListScenarios = []sce.Scenario{
		// List tests
		{
			Name:    "List - Valid list",
			Args:    "my-identifier list",
			WantErr: false,
		},
		{
			Name:    "List - Too many arguments",
			Args:    "my-identifier list error",
			WantErr: true,
		},
		{
			Name:    "List - Invalid option presence",
			Args:    "my-identifier list --error",
			WantErr: true,
		},
	}
)
