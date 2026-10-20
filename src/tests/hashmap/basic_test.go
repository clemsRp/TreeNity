package hashmap

import (
	"testing"
	h "treenity/src/server/hashmap"
	sce "treenity/src/tests/scenarios/hashmap"
)

func TestBasic(t *testing.T) {
	tests := sce.BasicScenarios

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			myMap := h.CreateHashmap(5)

			for _, step := range tt.Steps {
				UtilTestAction(t, myMap, step)
			}
		})
	}
}
