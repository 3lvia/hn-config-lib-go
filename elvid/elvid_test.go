package elvid

import (
	"testing"

	"github.com/3lvia/hn-config-lib-go/env"
	"github.com/3lvia/hn-config-lib-go/testing/assert"
)

// replaceEnv compacts environment variables handling to increase readability of tests.
func replaceEnv(t *testing.T, vars []string) {
	t.Helper()

	err := env.Clear(testenv...)
	assert.NoErr(t, err)

	if vars != nil {
		err = env.Set(vars...)
		assert.NoErr(t, err)
	}
}

func Test_New(t *testing.T) {
	err := env.Save(testenv...)
	assert.NoErr(t, err)

	tests := []struct {
		title              string
		expectValue        string
		environmentSlice   []string
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "No Environment Variables",
			expectError:        true, // Fails successfully if ElvID is not running locally
			expectErrorMessage: "missing env var ELVID_BASE_URL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// log.Println("------------------------------------")
			// log.Println(tt.title)
			replaceEnv(t, tt.environmentSlice)

			_, err := New()
			assert.WantErr(t, tt.expectError, err, tt.expectErrorMessage)

			err = env.Reset()
			assert.NoErr(t, err)
			// log.Println("------------------------------------")
		})
	}

	err = env.Reset()
	assert.NoErr(t, err)
}
