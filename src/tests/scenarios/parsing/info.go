package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingInfoScenarios = []sce.BasicScenario{
		// Info tests
		{
			Name:    "Info - Get valid info",
			Args:    "my-identifier info my-subscriber",
			WantErr: false,
		},
		{
			Name:    "Info - Missing subscriber argument",
			Args:    "my-identifier info",
			WantErr: true,
		},
		{
			Name:    "Info - Too many arguments",
			Args:    "my-identifier info my-subscriber error",
			WantErr: true,
		},
		{
			Name:    "Info - Invalid option presence",
			Args:    "my-identifier info my-subscriber --error",
			WantErr: true,
		},
	}
)
