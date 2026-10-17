package parsing

import (
	"strings"
	"testing"
	"treenity/src/client"
)

func TestParsing(t *testing.T) {
	// Init client
	c := client.NewClient()

	// Define all tests
	tests := ParsingScenarios

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Split command
			args := strings.Fields(tt.args)
			ctx, err := c.GetContext(args)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && ctx == nil {
				t.Errorf("GetContext() returned a nil context without error")
			}
		})
	}
}
