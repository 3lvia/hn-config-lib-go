package elvid

import (
	"fmt"

	"github.com/pkg/errors"
)

// PublicKeySet (Public Key Set) stores a slice of public keys and their metadata
type PublicKeySet struct {
	Keys []struct {
		KeyID     string   `json:"kid"`
		Algorithm string   `json:"alg"`
		X5C       []string `json:"x5c"`
	} `json:"keys"`
}

// Renews the stored public key set for the external ElvID server.
func (elvid *ElvID) newPublicKeySet() error {
	// log.Println("elvid.JsonWebKeySetUri:\t", elvid.JsonWebKeySetUri)
	// log.Println("elvid.PublicKeySet:\t\t", elvid.PublicKeySet)
	err := elvid.client.Get(elvid.JsonWebKeySetUri, &elvid.PublicKeySet)
	if err != nil {
		return errors.Wrap(err, "while renewing ElvID public key set")
	}

	return nil
}

// keyFunc converts a jwt token to a RSA public key
func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	cert, err := getPemCert(token)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
}

// getPemCert extracts the pem certificate from a jwt token
func getPemCert(token *jwt.Token) (cert string, err error) {
	for _, k := range tmpPublicKeySet.Keys {
		if kid, ok := token.Header["kid"].(string); ok {
			if kid == k.KeyID {
				cert = "-----BEGIN CERTIFICATE-----\n" + k.X5C[0] + "\n-----END CERTIFICATE-----"
				return
			}
		} else {
			return "", errors.New("expecting JWT header to have string kid")
		}
	}

	return "", errors.New("Unable to find corresponding kid")
}

var tmpPublicKeySet PublicKeySet // needs to be globally accessible because of how dgrijalva/jwt-go works. Not for caching; set before each use, nil after use.

func provideKeys(PublicKeySet PublicKeySet) {
	tmpPublicKeySet = PublicKeySet
}

func revokeKeys() {
	tmpPublicKeySet = PublicKeySet{nil}
}
