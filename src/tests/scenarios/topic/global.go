package topic

import sce "treenity/src/tests/scenarios"

var ActionsScenarios = []sce.TopicScenario{
	{
		Name: "Empty topics list",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionList, WantNames: nil},
			{Action: sce.TopicActionGet, TopicName: "ghost", WantExist: false},
		},
	},
	{
		Name: "Create keeps creation order",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "zeta"},
			{Action: sce.TopicActionCreate, TopicName: "alpha"},
			{Action: sce.TopicActionCreate, TopicName: "mid"},
			{Action: sce.TopicActionList, WantNames: []string{"zeta", "alpha", "mid"}},
			{Action: sce.TopicActionGet, TopicName: "alpha", WantExist: true},
		},
	},
	{
		Name: "Create twice returns the same topic without duplicate",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "t"},
			{Action: sce.TopicActionCreate, TopicName: "t", WantSame: true},
			{Action: sce.TopicActionList, WantNames: []string{"t"}},
		},
	},
	{
		Name: "Delete removes from map and list",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "a"},
			{Action: sce.TopicActionCreate, TopicName: "b"},
			{Action: sce.TopicActionCreate, TopicName: "c"},
			{Action: sce.TopicActionDelete, TopicName: "b"},
			{Action: sce.TopicActionGet, TopicName: "b", WantExist: false},
			{Action: sce.TopicActionList, WantNames: []string{"a", "c"}},
		},
	},
	{
		Name: "Delete stops the topic",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "t"},
			{Action: sce.TopicActionIsStopped, TopicName: "t", WantStopped: false},
			{Action: sce.TopicActionDelete, TopicName: "t"},
			{Action: sce.TopicActionIsStopped, TopicName: "t", WantStopped: true},
		},
	},
	{
		Name: "Recreate after delete gives a fresh running topic",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "t"},
			{Action: sce.TopicActionDelete, TopicName: "t"},
			{Action: sce.TopicActionCreate, TopicName: "t"},
			{Action: sce.TopicActionIsStopped, TopicName: "t", WantStopped: false},
			{Action: sce.TopicActionList, WantNames: []string{"t"}},
		},
	},
	{
		Name: "Delete unknown topic fails with code 2",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionDelete, TopicName: "ghost", WantErr: true, WantCode: 2},
		},
	},
	{
		Name: "Delete twice fails the second time",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "t"},
			{Action: sce.TopicActionDelete, TopicName: "t"},
			{Action: sce.TopicActionDelete, TopicName: "t", WantErr: true, WantCode: 2},
		},
	},
	{
		Name: "Failed delete does not alter the list",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "keep"},
			{Action: sce.TopicActionDelete, TopicName: "ghost", WantErr: true},
			{Action: sce.TopicActionList, WantNames: []string{"keep"}},
		},
	},
	{
		Name: "StopAll stops every topic but keeps them registered",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "a"},
			{Action: sce.TopicActionCreate, TopicName: "b"},
			{Action: sce.TopicActionStopAll},
			{Action: sce.TopicActionIsStopped, TopicName: "a", WantStopped: true},
			{Action: sce.TopicActionIsStopped, TopicName: "b", WantStopped: true},
			{Action: sce.TopicActionList, WantNames: []string{"a", "b"}},
		},
	},
	{
		Name: "StopAll twice then Delete does not panic",
		Steps: []sce.TopicStep{
			{Action: sce.TopicActionCreate, TopicName: "t"},
			{Action: sce.TopicActionStopAll},
			{Action: sce.TopicActionStopAll},
			{Action: sce.TopicActionDelete, TopicName: "t"},
		},
	},
}
