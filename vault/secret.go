package vault

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type Secret interface {
	GetRequestID() string
	GetLeaseID() string
	IsRenewable() bool
	GetLeaseDuration() int
	GetData() map[string]interface{}
	GetMetadata() map[string]interface{}
}

// secret contains all data and metadata from a Vault secret
type secret struct {
	RequestID     string                            `json:"request_id"`
	LeaseID       string                            `json:"lease_id"`
	Renewable     bool                              `json:"renewable"`
	LeaseDuration int                               `json:"lease_duration"`
	Data          map[string]map[string]interface{} `json:"data"`
}

func (s *secret) GetRequestID() string {
	return s.RequestID
}

func (s *secret) GetLeaseID() string {
	return s.LeaseID
}

func (s *secret) IsRenewable() bool {
	return s.Renewable
}

func (s *secret) GetLeaseDuration() int {
	return s.LeaseDuration
}

func (s *secret) GetData() map[string]interface{} {
	return s.Data["data"]
}

func (s *secret) GetMetadata() map[string]interface{} {
	return s.Data["metadata"]
}

// GetSecret returns the secret from the provided path.
// In case of 403 response from server, the credentials will be renewed and the request retried once.
func (vault *Vault) GetSecret(path string) (Secret, error) {
	url := makeURL(vault.Config.Addr, path)

	req, err := secretsReq(url, vault.Token.Auth.ClientToken)
	if err != nil {
		return nil, err
	}

	secret := new(secret)
	if err = vault.do(req, &secret); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("while getting secret from Vault. path: %s url: %s", path, url))
	}

	return secret, nil
}

// secretsReq returns a http request for getting secrets from Vault
func secretsReq(url, auth string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "while building http request")
	}

	req.Header.Set("X-Vault-Token", auth)

	return req, nil
}

/**

// do performs a request to the vault and retries once if the request returns 403.
// Before retrying, the credentials will be renewed
// Results are unmarshaled to the destination, unless dst is an io.Writer whitch will get the raw response instead
func (vault Vault) do(req *http.Request, dst interface{}) (err error) {
	err = vault.Client.Do(req, &dst)
	if err != nil && strings.Contains(err.Error(), "status 403") {

		if err2 := vault.Authenticate(); err2 != nil {
			return errors.Wrapf(err2, "while renewing credentials, after first getting: %e", err)
		}

		err = vault.Client.Do(req, &dst)
	}

	return
}

// do2 performs a request to the vault and retries once if the request returns 403.
// Before retrying, the credentials will be renewed
// Results are unmarshaled to the destination, unless dst is an io.Writer whitch will get the raw response instead
func (vault Vault) do2(req *http.Request, dst interface{}) error {
	lease := time.Duration(vault.Token.LeaseDuration) * time.Second
	if time.Since(vault.Token.CreatedAt) > lease {
		if err := vault.Authenticate(); err != nil {
			return errors.Wrap(err, "while renewing expired credentials")
		}
	}

	return vault.Client.Do(req, &dst)
}

*/

func (vault Vault) do(req *http.Request, dst interface{}) error {
	return vault.Client.Do(req, &dst)
}

type rolesets struct {
	Data map[string][]string `json:"data"`
}

type gcpCredentials struct {
	LeaseDuration int                `json:"lease_duration"`
	Data          gcpCredentialsData `json:"data"`
}

type gcpCredentialsData struct {
	PrivateKeyData string `json:"private_key_data"`
}

// GetSecret returns the secret from the provided path.
// In case of 403 response from server, the credentials will be renewed and the request retried once.
func (vault *Vault) SetupAndRenewGcpCredentials(system string, roleset string, ttl int) error {
	err := validateRoleset(vault, system)
	if err != nil {
		return err
	}

	url := makeURL(vault.Config.Addr, fmt.Sprintf("%s/gcp/key/%s", system, roleset))
	fmt.Printf("url %s\n", url)

	req, err := secretsReq(url, vault.Token.Auth.ClientToken)
	if err != nil {
		return err
	}

	gcpCredentials := new(gcpCredentials)
	if err = vault.do(req, &gcpCredentials); err != nil {
		return errors.Wrap(err, fmt.Sprintf("while getting gcpCredentials from Vault. system: %s url: %s", system, url))
	}
	fmt.Printf("LeaseDuration %d\n", gcpCredentials.LeaseDuration)
	// fmt.Printf("PrivateKeyData %s\n", gcpCredentials.Data.PrivateKeyData)

	credsFilename := fmt.Sprintf("%s/go-creds.json", os.TempDir())
	err = ioutil.WriteFile(credsFilename, []byte(gcpCredentials.Data.PrivateKeyData), 0600)

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credsFilename)
	return nil
}

func validateRoleset(vault *Vault, system string) error {
	url := makeURL(vault.Config.Addr, fmt.Sprintf("%s/gcp/rolesets?list=true", system))

	req, err := secretsReq(url, vault.Token.Auth.ClientToken)
	if err != nil {
		return err
	}

	rolesets := new(rolesets)
	if err = vault.do(req, &rolesets); err != nil {
		return errors.Wrap(err, fmt.Sprintf("while getting rolesets from Vault. system: %s url: %s", system, url))
	}
	fmt.Printf("Rolesets %s\n", rolesets.Data["keys"])
	return nil
}

//{"request_id":"2bf836c8-ead6-cf57-b4c5-21644a0351a1","lease_id":"","renewable":false,"lease_duration":0,"data":{"keys":["storage_admin"]},"wrap_info":null,"warnings":null,"auth":null}
