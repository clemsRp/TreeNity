package hashmap

import (
	"testing"
	h "treenity/src/server/hashmap"
	sce "treenity/src/tests/scenarios/hashmap"
)

func TestCollisions(t *testing.T) {
	tests := sce.CollisionsScenarios

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			myMap := h.CreateHashmap(1)

			for _, step := range tt.Steps {
				UtilTestAction(t, myMap, step)
			}
		})
	}
}
