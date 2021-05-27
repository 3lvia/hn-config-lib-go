package vault

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/3lvia/hn-config-lib-go/env"
	"github.com/3lvia/hn-config-lib-go/testing/assert"
	"github.com/3lvia/hn-config-lib-go/testing/mock"
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
		name      string
		envslice  []string
		wantErr   bool
		errWanted string
	}{
		{
			name:      "no environment variables",
			envslice:  []string{},
			wantErr:   true,
			errWanted: "missing env var VAULT_ADDR",
		}, {
			name:      "broken authentification",
			envslice:  []string{envars["addr"], mock.Addr, envars["github"], mock.Token},
			wantErr:   true,
			errWanted: "while do-ing http request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replaceEnv(t, tt.envslice)

			_, err := New()

			assert.WantErr(t, tt.wantErr, err, tt.errWanted)
		})
	}

	err = env.Reset()
	assert.NoErr(t, err)
}
func Test_GetSecret(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	myVault, err := New()
	if err != nil {
		log.Fatal(err)
	}

	// Get a secret from the vault
	mySecret, err := myVault.GetSecret("libraries/kv/data/manual/integrationtest")
	if err != nil {
		log.Fatal(err)
	}
	data := mySecret.GetData()["mykey"].(string)
	if len(data) == 0 {
		t.Errorf("len(creds) == 0")
	}
}

func Test_SetupAndRenewGcpCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	myVault, err := New()
	if err != nil {
		log.Fatal(err)
	}

	// Get a secret from the vault
	err = myVault.SetupAndRenewGcpCredentials("monitoring", "storage_admin", 5)
	if err != nil {
		log.Fatal(err)
	}

	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(creds) == 0 {
		t.Errorf("len(creds) == 0")
	}
}

func Test_SetupAndRenewGcpCredentials_VerifyRenewing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	myVault, err := New()
	if err != nil {
		log.Fatal(err)
	}

	err = myVault.SetupAndRenewGcpCredentials("monitoring", "storage_admin", 1)
	if err != nil {
		log.Fatal(err)
	}

	creds1, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	// fmt.Printf("creds1: %s\n", creds1[500:509])
	if err != nil {
		log.Fatal(err)
	}
	if len(creds1) == 0 {
		t.Errorf("len(creds1) == 0")
	}

	time.Sleep(5 * time.Second)

	creds2, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	// fmt.Printf("creds2: %s\n", creds2[500:509])
	if err != nil {
		log.Fatal(err)
	}
	if len(creds2) == 0 {
		t.Errorf("len(creds1) == 0")
	}
	if string(creds1) == string(creds2) {
		t.Errorf("creds1 == creds2 after renewing")
	}
}
