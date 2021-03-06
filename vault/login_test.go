package vault

import (
	"bytes"
	"encoding/json"
	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
	"testing"
)

func Test_loginBuffer(t *testing.T) {
	// Test with empty arg
	got, err := loginBuffer(nil)
	assert.NoErr(t, err)
	if got.Len() != 4 {
		t.Fatalf("result is %d bytes, should be 4 bytes", got.Len())
	}

	// Test contents of result
	js, err := json.Marshal(mock.Token)
	assert.NoErr(t, err)
	want := bytes.NewBuffer(js)
	got, err = loginBuffer(mock.Token)
	assert.NoErr(t, err)
	if got.String() != want.String() {
		t.Fatal("got incorrect bytestring")
	}
}
