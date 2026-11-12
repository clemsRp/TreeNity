package parsing

import (
	sce "treenity/src/tests/scenarios"
)

var (
	ParsingIdScenarios = []sce.IdScenario{
		{
			Name:    "Id - valid id",
			Id:      "leo",
			WantErr: true,
		},
		{
			Name:    "Id - valid UPPERCASE id",
			Id:      "LEO",
			WantErr: true,
		},
		{
			Name:    "Id - invalid character id",
			Id:      "leo!",
			WantErr: false,
		},
		{
			Name:    "Id - missing id",
			Id:      "",
			WantErr: false,
		},
		{
			Name:    "Id - all character authorised id",
			Id:      "A-z_0.9",
			WantErr: true,
		},
		{
			Name:    "Id - length limit id",
			Id:      "asfdkjnbvdtrjiopklhjnfdasvgbeoqa",
			WantErr: true,
		},
		{
			Name:    "Id - too long id",
			Id:      "asfdkjnbvdtrjiopklhjnfdasvgbeoqaasdfasdfasdfasdfasdfasdfas",
			WantErr: false,
		},
		{
			Name:    "Id - 33 character id",
			Id:      "123456789012345678901234567890123",
			WantErr: false,
		},
		{
			Name:    "Id - min length id",
			Id:      "a",
			WantErr: true,
		},
		{
			Name:    "Id - ' character in id",
			Id:      "it'sme",
			WantErr: false,
		},
		{
			Name:    "Id - _ character in id",
			Id:      "it_s_me",
			WantErr: true,
		},
		{
			Name:    "Id - - character in id",
			Id:      "it-s-me",
			WantErr: true,
		},
		{
			Name:    "Id - . character in id",
			Id:      "it.s.me",
			WantErr: true,
		},
		{
			Name:    "Id - space in middle id",
			Id:      "its me",
			WantErr: false,
		},
		{
			Name:    "Id - space at end id",
			Id:      "its_me ",
			WantErr: false,
		},
		{
			Name:    "Id - accented character id",
			Id:      "léo",
			WantErr: false,
		},
		{
			Name:    "Id - new line id",
			Id:      "test\n",
			WantErr: false,
		},
		{
			Name:    "Id - slash character id",
			Id:      "test/here",
			WantErr: false,
		},
	}
)
