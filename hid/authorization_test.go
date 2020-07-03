package hid

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"

	"github.com/dgrijalva/jwt-go"
)

// at overrides time value for tests and restores default value after
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
} //time.Unix(0, 0)

func Test_HID_AuthorizeRequest(t *testing.T) {
	type args struct {
		audience string
		scope    string
	}
	tests := []struct {
		name      string
		time      time.Time
		args      args
		wantErr   bool
		errWanted string
	}{
		{
			name:      "invalid audience",
			time:      time.Now().Add(time.Minute * 30),
			args:      args{mock.ID, mock.ID},
			wantErr:   true,
			errWanted: "Invalid audience",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at(tt.time, func() {
				hid, err := New()
				assert.NoErr(t, err)

				token, err := hid.GetToken(os.Getenv("TEST_HID_ID"), os.Getenv("TEST_HID_SECRET"))
				assert.NoErr(t, err)

				req, err := http.NewRequest("GET", mock.Addr, nil)
				assert.NoErr(t, err)
				req.Header.Add("Authorization", token.Raw)

				err = hid.AuthorizeRequest(req, tt.args.audience, tt.args.scope)
				assert.WantErr(t, tt.wantErr, err, tt.errWanted)
			})
		})
	}
}
