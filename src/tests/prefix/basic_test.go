package prefix

import (
	"testing"
	p "treenity/src/server/prefix"
	sce "treenity/src/tests/scenarios/prefix"
)

func TestBasic(t *testing.T) {
	tests := sce.ActionsScenarios

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			myTree := p.CreateTree(0)

			for _, step := range tt.Steps {
				UtilTestAction(t, myTree, step)
			}
		})
	}
}
