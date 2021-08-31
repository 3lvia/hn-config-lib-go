package elvid

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
	"github.com/golang-jwt/jwt/v4"
)

// at overrides time value for tests and restores default value after
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
} //time.Unix(0, 0)

func Test_ElvID_AuthorizeRequest(t *testing.T) {
	type arguments struct {
		scope string
	}
	tests := []struct {
		title              string
		time               time.Time
		arguments          arguments
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:              "Invalid Scope",
			time:               time.Now().Add(time.Minute * 30),
			arguments:          arguments{mock.ID},
			expectError:        true,
			expectErrorMessage: "Did not find valid scope",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			at(tt.time, func() {
				// log.Println("------------------------------------")
				// log.Println(tt.title)
				elvid, err := New()
				assert.NoErr(t, err)

				token, err := elvid.GetToken(os.Getenv("ELVID_MACHINE_CLIENT_ID"), os.Getenv("ELVID_MACHINE_CLIENT_SECRET"))
				assert.NoErr(t, err)

				req, err := http.NewRequest("GET", mock.Addr, nil)
				assert.NoErr(t, err)
				req.Header.Add("Authorization", token.AccessToken)
				// req.Header.Add("Authorization", os.Getenv("ELVID_ACCESS_TOKEN"))

				err = elvid.AuthorizeRequest(req, tt.arguments.scope)
				assert.WantErr(t, tt.expectError, err, tt.expectErrorMessage)
				// log.Println("------------------------------------")
			})
		})
	}
}
