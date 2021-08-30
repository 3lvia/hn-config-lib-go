package elvid

import (
	"log"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// Token exp
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	Profile      string `json:"profile"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	SessionState string `json:"session_state"`
	State        string `json:"state"`
	TokenType    string `json:"token_type"`
}

// GetToken provides the credentials of a user or service, and returns a token for sending with requests to a service.
func (elvid ElvID) GetToken(user, secret string) (token *Token, err error) {
	form := map[string][]string{
		"client_id":     []string{user},
		"client_secret": []string{secret},
		"grant_type":    []string{"client_credentials"},
	}

	err = elvid.client.PostForm(elvid.TokenEndPoint, form, &token)
	err = errors.Wrap(err, "while getting token from ElvID")

	return
}

func (elvid ElvID) IsValidAccessToken(accessToken string) (isTokenValid bool, err error) {
	isTokenValid = false

	// const jwksUrlEndpoint = `/.well-known/openid-configuration/jwks`
	// jwksUrl := authorityBaseUrl + jwksUrlEndpoint
	log.Println("JWKS URL: ", elvid.Configuration.JsonWebKeySetUri)

	// Create the keyfunc options. Refresh the JWKS every hour and log errors.
	refreshInterval := time.Hour
	options := keyfunc.Options{
		RefreshInterval: &refreshInterval,
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.KeyFunc\nError: %s", err.Error())
		},
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(elvid.Configuration.JsonWebKeySetUri, options)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}

	// Parse the JWT.
	token, err := jwt.Parse(accessToken, jwks.KeyFunc)
	if err != nil {
		log.Fatalf("Failed to parse the JWT.\nError: %s", err.Error())
	}
	log.Println()
	log.Println("Claims:\t\t", token.Claims)
	log.Println("Header:\t\t", token.Header)
	log.Println("Method.Alg():\t", token.Method.Alg())
	log.Println("Raw:\t\t", token.Raw)
	log.Println("Signature:\t\t", token.Signature)
	log.Println("Valid:\t\t", token.Valid)
	log.Println()

	// Check if the token is valid.
	if !token.Valid {
		log.Fatalf("The token is not valid.")
		isTokenValid = false
	} else {
		log.Println("The token is valid.")
		isTokenValid = true
	}

	return isTokenValid, err
}

// Append the raw token to the header of the provided request.
func (token Token) AppendToRequest(req *http.Request) {
	req.Header.Add("Authorization", token.AccessToken)
}
