package parsing

import (
	"testing"
	"treenity/src/client"
	sce "treenity/src/tests/scenarios/parsing"
)

func TestInfoParsing(t *testing.T) {
	// Init client
	c := client.NewClient()

	// Define all tests
	tests := sce.ParsingInfoScenarios

	UtilTestScenario(t, c, tests)
}
