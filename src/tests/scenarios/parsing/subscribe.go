package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingSubscribeScenarios = []sce.BasicScenario{
		// Subscribe tests
		{
			Name:    "Subscribe - Valid topic and subscriber",
			Args:    "my-identifier subscribe my-topic my-subscriber",
			WantErr: false,
		},
		{
			Name:    "Subscribe - Valid with single option",
			Args:    "my-identifier subscribe my-topic my-subscriber --raw",
			WantErr: false,
		},
		{
			Name:    "Subscribe - Valid with multiple options",
			Args:    "my-identifier subscribe my-topic my-subscriber --prefix my-prefix --raw",
			WantErr: false,
		},
		{
			Name:    "Subscribe - Invalid option value placement",
			Args:    "my-identifier subscribe my-topic my-subscriber --prefix --raw my-prefix",
			WantErr: true,
		},
		{
			Name:    "Subscribe - Missing required arguments",
			Args:    "my-identifier subscribe",
			WantErr: true,
		},
		{
			Name:    "Subscribe - Too many arguments",
			Args:    "my-identifier subscribe my-topic my-subscriber error",
			WantErr: true,
		},
		{
			Name:    "Subscribe - Invalid option presence",
			Args:    "my-identifier subscribe my-topic my-subscriber --error",
			WantErr: true,
		},

		// Subscribe prefix
		{
			Name:    "Subscribe - Invalid option value (float)",
			Args:    "my-identifier subscribe my-topic my-subscriber --prefix 5.5 --raw",
			WantErr: false,
		},
		{
			Name:    "Subscribe - Invalid option value (string)",
			Args:    "my-identifier subscribe my-topic my-subscriber --prefix string --raw",
			WantErr: false,
		},
		{
			Name:    "Subscribe - Valid option value (int)",
			Args:    "my-identifier subscribe my-topic my-subscriber --prefix 5 --raw",
			WantErr: false,
		},

		// Subscribe offset
		{
			Name:    "Subscribe - Invalid option value (float)",
			Args:    "my-identifier subscribe my-topic my-subscriber --offset 5.5 --raw",
			WantErr: true,
		},
		{
			Name:    "Subscribe - Invalid option value (string)",
			Args:    "my-identifier subscribe my-topic my-subscriber --offset string --raw",
			WantErr: true,
		},
		{
			Name:    "Subscribe - Valid option value (int)",
			Args:    "my-identifier subscribe my-topic my-subscriber --offset 5 --raw",
			WantErr: false,
		},
	}
)
