package ipc

import (
	sce "treenity/src/tests/scenarios"
)

var ActionsScenarios = []sce.IpcScenario{
	{
		Name: "Basic FIFO management test",
		Steps: []sce.IpcStep{
			// Init FIFO's
			{Action: sce.IpcActionCreateReader, Pid: 0, WantErr: false},
			{Action: sce.IpcActionCreateWriter, Pid: 0, WantErr: false},

			// Communicate
			{Action: sce.IpcActionWrite, Pid: 0, Text: "user.create:topix", WantErr: false},
			{Action: sce.IpcActionRead, Pid: 0, Text: "user.create:topix", WantErr: false},
		},
	},
	{
		Name: "Basic FIFO error test",
		Steps: []sce.IpcStep{
			// Communicate on invalid workers
			{Action: sce.IpcActionWrite, Pid: -1, WantErr: true},
			{Action: sce.IpcActionRead, Pid: -1, WantErr: true},
		},
	},
	{
		Name: "Multiple sequential messages test",
		Steps: []sce.IpcStep{
			{Action: sce.IpcActionCreateReader, Pid: 1, WantErr: false},
			{Action: sce.IpcActionCreateWriter, Pid: 1, WantErr: false},

			// Send multiple message
			{Action: sce.IpcActionWrite, Pid: 1, Text: "message.one", WantErr: false},
			{Action: sce.IpcActionWrite, Pid: 1, Text: "message.two", WantErr: false},
			{Action: sce.IpcActionWrite, Pid: 1, Text: "message.three", WantErr: false},

			// Check messages order
			{Action: sce.IpcActionRead, Pid: 1, Text: "message.one", WantErr: false},
			{Action: sce.IpcActionRead, Pid: 1, Text: "message.two", WantErr: false},
			{Action: sce.IpcActionRead, Pid: 1, Text: "message.three", WantErr: false},
		},
	},
}
