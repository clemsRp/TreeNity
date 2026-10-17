package parsing

import (
	"testing"
	"treenity/src/client"
	sce "treenity/src/tests/scenarios/parsing"
)

func TestProduceParsing(t *testing.T) {
	// Init client
	c := client.NewClient()

	// Define all tests
	tests := sce.ParsingProduceScenarios

	UtilTestScenario(t, c, tests)
}
