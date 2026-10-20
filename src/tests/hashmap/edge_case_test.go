package hashmap

import (
	"testing"
	h "treenity/src/server/hashmap"
	sce "treenity/src/tests/scenarios/hashmap"
)

func TestEdgeCase(t *testing.T) {
	tests := sce.EdgeCaseScenarios

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			myMap := h.CreateHashmap(10)

			for _, step := range tt.Steps {
				UtilTestAction(t, myMap, step)
			}
		})
	}
}
