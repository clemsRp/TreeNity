package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingGlobalScenarios = []sce.Scenario{
		// Global test
		{
			Name:    "Global - Invalid number of arguments (missing command)",
			Args:    "my-identifier",
			WantErr: true,
		},
	}
)
