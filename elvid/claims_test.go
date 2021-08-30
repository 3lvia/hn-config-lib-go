package elvid

import (
	"testing"
	"time"

	"github.com/3lvia/hn-config-lib-go/testing/mock"
)

func Test_VerifyClaims(t *testing.T) {
	type arguments struct {
		issuer string
		scope  string
	}
	tests := []struct {
		title              string
		time               time.Time
		arguments          arguments
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "invalid audience",
			time:               time.Now().Add(time.Minute * 30),
			arguments:          arguments{mock.Addr, mock.ID},
			expectError:        true,
			expectErrorMessage: "Invalid audience",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			at(tt.time, func() {
				/**
				elvid, err := New()
				assert.NoErr(t, err)

				token, err := elvid.GetToken(os.Getenv("TEST_ELVID_ID"), os.Getenv("TEST_ELVID_SECRET"))
				assert.NoErr(t, err)

				jwt, err := elvid.authenticate(token.Raw)
				assert.NoErr(t, err)

				err = VerifyClaims(jwt, tt.arguments.issuer, tt.arguments.audience, tt.arguments.scope)
				assert.WantErr(t, tt.expectError, err, tt.expectErrorMessage)
				*/
			})
		})
	}
}
