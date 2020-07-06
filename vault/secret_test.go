package vault

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
)

func TestVault_GetSecret(t *testing.T) {
	tests := []struct {
		name      string
		vault     Vault
		want      Secret
		wantErr   bool
		errWanted string
	}{
		{
			name:      "forbidden access",
			vault:     Vault{Config: Config{Client: mock.ClientForbidden}},
			want:      nil,
			wantErr:   true,
			errWanted: "403: forbidden",
		}, {
			name:    "access granted",
			vault:   Vault{Config: Config{Client: mock.Client}},
			want:    &secret{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.vault.GetSecret("")
			assert.WantErr(t, tt.wantErr, err, tt.errWanted)
			assert.DeepResult(t, got, tt.want)
		})
	}
}

func Test_secretsReq(t *testing.T) {
	type args struct {
		url  string
		auth string
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "build request",
			args: args{mock.Addr, mock.Token},
			want: mock.Request(t, "GET", mock.Addr, "", "X-Vault-Token", mock.Token),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := secretsReq(tt.args.url, tt.args.auth)

			assert.NoErr(t, err)

			assert.Result(t, got.Method, tt.want.Method)
			assert.DeepResult(t, got.URL, tt.want.URL)
			assert.DeepResult(t, got.Header, tt.want.Header)
		})
	}
}

func TestVault_do(t *testing.T) {
	tests := []struct {
		name      string
		vault     Vault
		wantErr   bool
		errWanted string
	}{
		{
			name:      "forbidden access, without authentication",
			vault:     Vault{Config: Config{Client: mock.ClientForbidden}},
			wantErr:   true,
			errWanted: "403: forbidden",
		}, {
			name:    "access granted",
			vault:   Vault{Config: Config{Client: mock.Client}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vault.do(nil, nil)
			assert.WantErr(t, tt.wantErr, err, tt.errWanted)
		})
	}
}

func Test_deserializeFromJSON(t *testing.T) {
	// Arrange
	var secret secret

	// Act
	err := json.Unmarshal([]byte(secretResponse), &secret)
	if err != nil {
		t.Errorf("unexpected error while deserializing %v", err)
	}
	serviceAccountKey := secret.Data["data"]["service-account-key"]

	// Assert
	if serviceAccountKey != "my-secret-key" {
		t.Errorf("unexpected service account key %s", secret.Data["data"]["servic-account-key"])
	}
}

const secretResponse string  = `{"request_id":"464ac0ff-fa12-13ca-9e6d-ca3be05ac802","lease_id":"","renewable":false,"lease_duration":0,"data":{"data":{"service-account-key":"my-secret-key"},"metadata":{"created_time":"2020-06-02T08:26:38.487373863Z","deletion_time":"","destroyed":false,"version":1}},"wrap_info":null,"warnings":null,"auth":null}`
