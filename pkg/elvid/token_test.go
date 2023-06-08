package elvid

import (
	"os"
	"testing"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
)

func Test_ElvID_GetToken(t *testing.T) {
	type arguments struct {
		user   string
		secret string
	}
	tests := []struct {
		title              string
		arguments          arguments
		expectValue        int
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "Mock Credentials",
			arguments:          arguments{mock.ID, mock.Secret},
			expectError:        true,
			expectErrorMessage: "400 Bad Request",
		},
		{
			title:              "Actual Machine Client Credentials",
			arguments:          arguments{os.Getenv("ELVID_MACHINE_CLIENT_ID"), os.Getenv("ELVID_MACHINE_CLIENT_SECRET")},
			expectValue:        3600,
			expectError:        false,
			expectErrorMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// log.Println("------------------------------------")
			// log.Println(tt.title)
			elvid, err := New()
			assert.NoErr(t, err)

			gotToken, err := elvid.GetToken(tt.arguments.user, tt.arguments.secret)
			// log.Println("gotToken:", gotToken)
			// log.Println("expectError:", tt.expectError)
			// log.Println("expectErrorMessage:", tt.expectErrorMessage)
			// log.Println("err:", err)
			assert.WantErr(t, tt.expectError, err, tt.expectErrorMessage)

			if !tt.expectError {
				got := gotToken.ExpiresIn
				assert.Result(t, got, tt.expectValue)
				// log.Println("AccessToken:\t", gotToken.AccessToken)
				// log.Println("ExpiresAt:\t", gotToken.ExpiresAt)
				// log.Println("ExpiresIn:\t", gotToken.ExpiresIn)
				// log.Println("IdToken:\t", gotToken.IdToken)
				// log.Println("Profile:\t", gotToken.Profile)
				// log.Println("RefreshToken:\t", gotToken.RefreshToken)
				// log.Println("Scope:\t", gotToken.Scope)
				// log.Println("SessionState:\t", gotToken.SessionState)
				// log.Println("State:\t", gotToken.State)
				// log.Println("TokenType:\t", gotToken.TokenType)
			}
			// log.Println("------------------------------------")
		})
	}
}
