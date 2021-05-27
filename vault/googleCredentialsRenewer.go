package vault

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

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

func (vault *Vault) SetupAndRenewGcpCredentials(system string, roleset string, ttl int) error {
	err := validateRoleset(vault, system, roleset)
	if err != nil {
		return err
	}

	// fmt.Printf("PrivateKeyData %s\n", gcpCredentials.Data.PrivateKeyData)
	leaseDuration, err := SetupGcpCredentials(vault, system, roleset, ttl)
	if err != nil {
		return err
	}

	renewTime := time.Now().Add(time.Duration(leaseDuration) * time.Second)
	// fmt.Printf("renewTime: %s", renewTime)

	go func() {
		for range time.Tick(time.Second * 1) {
			// fmt.Println("Ticking every 1 seconds")
			if time.Now().Add(60 * time.Second).After(renewTime) {
				leaseDuration, _ = SetupGcpCredentials(vault, system, roleset, ttl)
				renewTime = time.Now().Add(time.Duration(leaseDuration) * time.Second)
				// fmt.Printf("renewTime: %s", renewTime)
			}
		}
	}()
	return nil
}

func SetupGcpCredentials(vault *Vault, system string, roleset string, ttl int) (int, error) {
	url := makeURL(vault.Config.Addr, fmt.Sprintf("%s/gcp/key/%s", system, roleset))
	// fmt.Printf("url %s\n", url)

	var jsonStr = []byte(fmt.Sprintf(`{"ttl":"%d"}`, ttl))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vault.Token.Auth.ClientToken)

	if err != nil {
		return 0, err
	}

	gcpCredentials := new(gcpCredentials)
	if err = vault.do(req, &gcpCredentials); err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("while getting gcpCredentials from Vault. system: %s url: %s", system, url))
	}
	// fmt.Printf("LeaseDuration %d\n", gcpCredentials.LeaseDuration)

	credsFilename := fmt.Sprintf("%s/go-creds.json", os.TempDir())
	err = ioutil.WriteFile(credsFilename, []byte(gcpCredentials.Data.PrivateKeyData), 0600)

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credsFilename)
	fmt.Println("GOOGLE_APPLICATION_CREDENTIALS environment variable set")
	// fmt.Printf("PrivateKeyData %s\n", gcpCredentials.Data.PrivateKeyData[500:509])
	return gcpCredentials.LeaseDuration, nil
}

func validateRoleset(vault *Vault, system string, roleset string) error {
	url := makeURL(vault.Config.Addr, fmt.Sprintf("%s/gcp/rolesets?list=true", system))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", vault.Token.Auth.ClientToken)

	rolesets := new(rolesets)
	if err = vault.do(req, &rolesets); err != nil {
		return errors.Wrap(err, fmt.Sprintf("while getting rolesets from Vault. system: %s url: %s", system, url))
	}
	for _, a := range rolesets.Data["keys"] {
		if a == roleset {
			return nil
		}
	}
	return fmt.Errorf("Roleset '%s' does not exist for the system '%s'. Rolesets: %s", roleset, system, rolesets.Data["keys"])
}
