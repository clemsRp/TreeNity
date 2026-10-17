package parsing

import (
	"testing"
	"treenity/src/client"
	sce "treenity/src/tests/scenarios/parsing"
)

func TestListParsing(t *testing.T) {
	// Init client
	c := client.NewClient()

	// Define all tests
	tests := sce.ParsingListScenarios

	UtilTestScenario(t, c, tests)
}
