package parsing

import (
	"strings"
	"testing"
	"treenity/src/client"
	sce "treenity/src/tests/scenarios"
)

func UtilTestScenario(t *testing.T, c *client.Client, tests []sce.Scenario) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Split command
			args := strings.Fields(tt.Args)
			ctx, err := c.GetContext(args)

			if (err != nil) != tt.WantErr {
				t.Errorf("GetContext() error = %v, WantErr %v", err, tt.WantErr)
				return
			}

			if !tt.WantErr && ctx == nil {
				t.Errorf("GetContext() returned a nil context without error")
			}
		})
	}
}
