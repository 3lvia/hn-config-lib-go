package elvid

import (
	"net/http"
	"os"

	"github.com/3lvia/hn-config-lib-go/service"
)

// IDManager represents a service that is able to provide clients with authorization tokens with the GetToken function,
// and is capable of authorizing these incoming tokens for the server with the AuthorizeRequest function.
type IDManager interface {
	GetToken(user, secret string) (token *Token, err error)
	AuthorizeRequest(r *http.Request, scope string) error
}

// ElvID holds the configurations and keys necessary to communicate with the ElvID service.
type ElvID struct {
	Configuration
	PublicKeySet
}

// New creates a new ElvID, performs necessary setup, and returns it as an IDManager
func New() (IDManager, error) {
	elvid := new(ElvID)
	cert := os.Getenv(envars["cert"])
	err := service.Setup(elvid, cert)
	return elvid, err
}

// ConnectToServer performs necessary setup for connections to the external ElvID service
func (elvid *ElvID) ConnectToServer() error {
	return elvid.newPublicKeySet()
}
