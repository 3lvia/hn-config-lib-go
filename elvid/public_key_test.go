package elvid

import (
	"testing"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
)

func Test_ElvID_newPKS(t *testing.T) {
	tests := []struct {
		title              string
		expectError        bool
		expectErrorMessage string
	}{
		{
			title:       "Mock Public Key Set",
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// log.Println("------------------------------------")
			// log.Println(tt.title)
			elvid := &ElvID{}
			err := elvid.Configure(mock.Client)
			assert.NoErr(t, err)

			err = elvid.newPublicKeySet()
			assert.WantErr(t, tt.expectError, err, tt.expectErrorMessage)
			// if !tt.expectError && elvid.PublicKeySet.Keys == nil {
			// 	t.Error("Did not get any keys")
			// }
			// log.Println("------------------------------------")
		})
	}
}
