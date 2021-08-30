package main

import (
	"os"
	"testing"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
)

func Test_EnvironmentVariables(t *testing.T) {
	type arguments struct {
		environmentVariableKey string
	}
	tests := []struct {
		title              string
		arguments          arguments
		expectValue        bool
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "VAULT_ADDR",
			arguments:          arguments{os.Getenv("VAULT_ADDR")},
			expectValue:        false,
			expectError:        false,
			expectErrorMessage: "",
		},
		{
			title:              "GITHUB_TOKEN",
			arguments:          arguments{os.Getenv("GITHUB_TOKEN")},
			expectValue:        false,
			expectError:        false,
			expectErrorMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			_, variablePresent := os.LookupEnv(tt.arguments.environmentVariableKey)
			assert.Result(t, variablePresent, tt.expectValue)
		})
	}
}
