package demonstration

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/vault"
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

func Test_VaultGetSecretPath(t *testing.T) {
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
			title:              "VAULT_SECRET_PATH_VALUE",
			arguments:          arguments{os.Getenv("VAULT_SECRET_PATH_VALUE")},
			expectValue:        true,
			expectError:        false,
			expectErrorMessage: "",
		},
		{
			title:              "INVALID_VAULT_SECRET_PATH_VALUE",
			arguments:          arguments{os.Getenv("INVALID_VAULT_SECRET_PATH_VALUE")},
			expectValue:        false,
			expectError:        true,
			expectErrorMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			vault, err := vault.New()
			assert.Result(t, err, nil)

			_, err = vault.GetSecret(tt.arguments.environmentVariableKey)
			hasError := err != nil
			assert.Result(t, hasError, !tt.expectValue)
		})
	}
}

func Test_VaultGetSecretData(t *testing.T) {
	type arguments struct {
		environmentVariableKey string
	}
	tests := []struct {
		title              string
		arguments          arguments
		expectValue        string
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "VAULT_SECRET_PATH_VALUE",
			arguments:          arguments{os.Getenv("VAULT_SECRET_PATH_VALUE")}, // "manual/kv/data/demonstration"
			expectValue:        "This secret is used for demonstration purposes only.",
			expectError:        false,
			expectErrorMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			vault, err := vault.New()
			assert.Result(t, err, nil)

			secret, err := vault.GetSecret(tt.arguments.environmentVariableKey)
			hasExpectValueAsSubstring := strings.Contains(fmt.Sprintf("%v", secret.GetData()["description"]), tt.expectValue)
			assert.Result(t, !tt.expectError, hasExpectValueAsSubstring)
		})
	}
}
