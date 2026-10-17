package parsing

import (
	"testing"
	"treenity/src/client"
	sce "treenity/src/tests/scenarios/parsing"
)

func TestCreateParsing(t *testing.T) {
	// Init client
	c := client.NewClient()

	// Define all tests
	tests := sce.ParsingCreateScenarios

	UtilTestScenario(t, c, tests)
}
