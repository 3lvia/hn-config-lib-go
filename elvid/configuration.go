package elvid

import (
	"os"

	"github.com/3lvia/hn-config-lib-go/libhttp"
	"github.com/pkg/errors"
)

const (
	DiscoveryEndpoint           = "/.well-known/openid-configuration"
	JsonWebKeySetEndpoint       = "/.well-known/openid-configuration/jwks"
	AuthorizationEndpoint       = "/connect/authorize"
	TokenEndpoint               = "/connect/token"
	UserInfoEndpoint            = "/connect/userinfo"
	EndSessionEndpoint          = "/connect/endsession"
	CheckSessionEndpoint        = "/connect/checksession"
	RevocationEndpoint          = "/connect/revocation"
	IntrospectionEndpoint       = "/connect/introspect"
	DeviceAuthorizationEndpoint = "/connect/deviceauthorization"
)

var envars = map[string]string{
	"addr":      "ELVID_BASE_URL",
	"cert":      "ELVID_CACERT",
	"discovery": "ELVID_DISCOVERY",
}

// Configuration contains the configuration information needed to do the initial setup and renewal of the ElvID service
type Configuration struct {
	Issuer           string `json:"issuer"`
	JsonWebKeySetUri string `json:"jwks_uri"`
	// AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndPoint string `json:"token_endpoint"`
	// UserInfoEndpoint            string `json:"userinfo_endpoint"`
	// EndSessionEndpoint          string `json:"end_session_endpoint"`
	// CheckSessionEndpoint        string `json:"check_session_iframe"`
	// RevocationEndpoint          string `json:"revocation_endpoint"`
	// IntrospectionEndpoint       string `json:"introspection_endpoint"`
	// DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`

	client libhttp.Client
}

func (elvid *ElvID) Configure(client libhttp.Client) error {
	addr := os.Getenv(envars["addr"])
	if addr == "" {
		return errors.New("missing env var " + envars["addr"])
	}

	discovery := os.Getenv(envars["discovery"])
	if discovery == "" {
		discovery = DiscoveryEndpoint
	}

	err := client.Get(addr+discovery, &elvid)
	if err != nil {
		elvid.JsonWebKeySetUri = addr + JsonWebKeySetEndpoint
		elvid.TokenEndPoint = addr + TokenEndpoint
	}

	elvid.Issuer = addr
	elvid.client = client

	return nil
}
