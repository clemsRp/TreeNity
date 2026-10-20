package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingCreateScenarios = []sce.BasicScenario{
		// Create tests
		{
			Name:    "Create - Valid topic",
			Args:    "my-identifier create my-topic",
			WantErr: false,
		},
		{
			Name:    "Create - Too many arguments",
			Args:    "my-identifier create my-topic error",
			WantErr: true,
		},
		{
			Name:    "Create - Missing topic argument",
			Args:    "my-identifier create",
			WantErr: true,
		},
		{
			Name:    "Create - Invalid option presence",
			Args:    "my-identifier create my-topic --error",
			WantErr: true,
		},
	}
)
