package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingProduceScenarios = []sce.Scenario{
		// Produce tests
		{
			Name:    "Produce - Valid topic",
			Args:    "my-identifier produce my-topic",
			WantErr: false,
		},
		{
			Name:    "Produce - Valid topic with option",
			Args:    "my-identifier produce my-topic --raw",
			WantErr: false,
		},
		{
			Name:    "Produce - Missing topic argument",
			Args:    "my-identifier produce",
			WantErr: true,
		},
		{
			Name:    "Produce - Too many arguments",
			Args:    "my-identifier produce my-topic error",
			WantErr: true,
		},
		{
			Name:    "Produce - Invalid option presence",
			Args:    "my-identifier produce my-topic --error",
			WantErr: true,
		},
	}
)
