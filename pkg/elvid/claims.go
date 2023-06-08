package elvid

import (
	"github.com/pkg/errors"
)

type claims struct {
	ClientID string   `json:"client_id,omitempty"`
	Scope    []string `json:"scope,omitempty"`
	jwt.StandardClaims
}

func verifyClaims(token *jwt.Token, issuer, scope string) error {
	if !token.Claims.(*claims).VerifyIssuer(issuer, false) {
		return errors.New("Invalid issuer")
	}

	return verifyScope(token, scope)
}

func verifyScope(token *jwt.Token, scope string) error {
	if claims, ok := token.Claims.(*claims); ok {
		for _, s := range claims.Scope {
			if s == scope {
				return nil // scope found, verified scope
			}
		}
		return errors.New("Did not find valid scope")
	}
	return errors.New("Invalid token")
}
