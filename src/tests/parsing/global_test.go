package parsing

import (
	"testing"
	"treenity/src/client"
	sce "treenity/src/tests/scenarios/parsing"
)

func TestGlobalParsing(t *testing.T) {
	// Init client
	c := client.NewClient()

	// Define all tests
	tests := sce.ParsingGlobalScenarios

	UtilTestScenario(t, c, tests)
}
