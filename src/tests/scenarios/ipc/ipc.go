package ipc

import sce "treenity/src/tests/scenarios"

var (
	IpcTestScenario = []sce.IpcScenario {
		{
			Name: "Ipc - valid payload",
			Input: "abc",
			Want: "abc",
		},
	}
)

