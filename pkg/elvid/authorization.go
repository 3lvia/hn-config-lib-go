package elvid

import (
	"net/http"

	"github.com/pkg/errors"
)

// AuthorizeRequest takes an incoming request on behalf of the service and extracts the token from the "Authorization" header.
// The token is then checked for authenticity, and then the claims of that token is verified against the provided scope.
func (elvid *ElvID) AuthorizeRequest(r *http.Request, scope string) error {
	rawToken := r.Header.Get("Authorization")

	token, err := elvid.authenticate(rawToken)
	if err != nil {
		return errors.Wrap(err, "while authenticating")
	}

	err = verifyClaims(token, elvid.Issuer, scope)
	if err != nil {
		return errors.Wrap(err, "while verifying claims")
	}

	return nil
}

// authenticate verifies the authenticity of a provided raw token
func (elvid *ElvID) authenticate(rawToken string) (*jwt.Token, error) {
	provideKeys(elvid.PublicKeySet)
	defer revokeKeys()

	token, err := parseToken(rawToken)
	if err != nil {

		// might fail if stored public keys are outdated. Renew keys and retry once
		err = elvid.newPublicKeySet()
		if err != nil {
			return nil, err
		}

		token, err = parseToken(rawToken)
		if err != nil {
			return nil, err
		}
	}

	if !token.Valid {
		return token, errors.New("Invalid token")
	}

	return token, nil
}

func parseToken(rawToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(rawToken, &claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}
