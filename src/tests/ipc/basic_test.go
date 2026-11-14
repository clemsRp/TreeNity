package ipc

import (
	"testing"
	sce "treenity/src/tests/scenarios/ipc"
)

func TestBasic(t *testing.T) {
	tests := sce.ActionsScenarios

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Init ipc worker map
			ipcMap := make(map[int]*IpcWorker)
			ipcMap[-1] = &IpcWorker{
				Reader: nil,
				Writer: nil,
			}

			// Execute scenarios
			for _, step := range tt.Steps {
				UtilTestAction(t, ipcMap, step)
			}

			// Clean FIFO paths
			for _, worker := range ipcMap {
				worker.Reader.Close()
				worker.Writer.Close()
			}
		})
	}
}
