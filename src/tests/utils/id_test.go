package parsing

import (
	"testing"
	"treenity/src/server/cmd"
	sce "treenity/src/tests/scenarios/utils"
)

func TestIdIsValid(t *testing.T) {
	for _, tt := range sce.ParsingIdScenarios {
		t.Run(tt.Name, func(t *testing.T) {
			got := cmd.IsValidId(tt.Id)
			if got != tt.WantErr {
				t.Errorf("IsValidId(%q) = %v, WantValid %v", tt.Id, got, tt.WantErr)
			}
		})
	}
}
