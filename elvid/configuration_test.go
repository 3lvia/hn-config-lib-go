package elvid

import (
	"testing"

	"github.com/3lvia/hn-config-lib-go/env"
	"github.com/3lvia/hn-config-lib-go/libhttp"
	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
)

var testenv = []string{
	envars["addr"],
	envars["cert"],
	envars["discovery"],
}

func Test_ElvID_Configuration(t *testing.T) {
	err := env.Save(testenv...)
	assert.NoErr(t, err)

	tests := []struct {
		title              string
		environmentSlice   []string
		client             libhttp.Client
		expectValue        Configuration
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "no environment variables",
			client:             mock.Client,
			expectValue:        Configuration{},
			expectError:        true,
			expectErrorMessage: "missing env var " + envars["addr"],
		}, {
			title:            "defaulting values",
			environmentSlice: []string{envars["addr"], mock.Addr},
			client:           mock.ClientForbidden,
			expectValue: Configuration{
				mock.Addr,
				mock.Addr + JsonWebKeySetEndpoint,
				// mock.Addr + AuthorizationEndpoint,
				mock.Addr + TokenEndpoint,
				mock.ClientForbidden,
				// mock.Addr + UserInfoEndpoint,
				// mock.Addr + EndSessionEndpoint,
				// mock.Addr + CheckSessionEndpoint,
				// mock.Addr + RevocationEndpoint,
				// mock.Addr + IntrospectionEndpoint,
				// mock.Addr + DeviceAuthorizationEndpoint,
			},
			expectError: false,
		}, {
			title:            "with all environment variables",
			environmentSlice: []string{envars["addr"], mock.Addr, envars["discovery"], mock.Path},
			client:           mock.Client,
			expectValue:      Configuration{Issuer: mock.Addr, client: mock.Client},
			expectError:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			replaceEnv(t, tt.environmentSlice)

			elvid := new(ElvID)
			err = elvid.Configure(tt.client)

			assert.WantErr(t, tt.expectError, err, tt.expectErrorMessage)
			assert.Result(t, elvid.Configuration, tt.expectValue)
		})
	}

	err = env.Reset()
	assert.NoErr(t, err)
}
