package demonstration // TODO: Rename to `package main` to running it locally.

import (
	"log"
	"net/http"
	"os"

	"github.com/3lvia/hn-config-lib-go/elvid"
	"github.com/3lvia/hn-config-lib-go/vault"
)

func demonstration() { // TODO: Rename to `func main()` before running it locally.
	vaultDemonstration()

	/*
		# Setup following environment variables.
		source_env ${HOME}
		export ELVID_BASE_URL="https://elvid.test-elvia.io"
		export ELVID_MACHINE_CLIENT_ID="00000000-0000-4000-8000-000000000000"
		export ELVID_MACHINE_CLIENT_SECRET="...."
		export ELVID_SCOPE="...."
		export GITHUB_TOKEN="...."
		export VAULT_ADDR="https://vault.dev-elvia.io"
		export VAULT_SECRET_PATH_VALUE="manual/kv/data/demonstration"
	*/
	elvidRequest := elvidMachineClientDemonstration()
	elvidApiDemonstration(elvidRequest)

	elvidUserClientDemonstration()
}

func vaultDemonstration() {
	vault, err := vault.New()
	if err != nil {
		log.Fatal(err)
	}

	secret, err := vault.GetSecret(os.Getenv("VAULT_SECRET_PATH_VALUE"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("${VAULT_SECRET_PATH_VALUE}", secret)

	log.Println("VaultSecret.GetData:\t\t", secret.GetData())
	log.Println("VaultSecret.GetLeaseDuration:\t", secret.GetLeaseDuration())
	log.Println("VaultSecret.GetLeaseID:\t\t", secret.GetLeaseID())
	log.Println("VaultSecret.GetMetadata:\t\t", secret.GetMetadata())
	log.Println("VaultSecret.GetCreatedTime:\t\t", secret.GetMetadata()["created_time"])
	log.Println("VaultSecret.GetDeletionTime:\t\t", secret.GetMetadata()["deletion_time"])
	log.Println("VaultSecret.GetDestroyed:\t\t", secret.GetMetadata()["destroyed"])
	log.Println("VaultSecret.GetVersion:\t\t", secret.GetMetadata()["version"])
	log.Println("VaultSecret.GetRequestID:\t\t", secret.GetRequestID())
}

func elvidUserClientDemonstration() {
	elvid, err := elvid.New()
	if err != nil {
		log.Fatal(err)
	}

	isValidAccessToken, err := elvid.HasValidUserClientAccessToken(os.Getenv("ELVID_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("IsValidAccessToken:", isValidAccessToken)
}

func elvidMachineClientDemonstration() *http.Request {
	elvid, err := elvid.New()
	if err != nil {
		log.Fatal(err)
	}

	elvidBearerToken, err := elvid.GetToken(os.Getenv("ELVID_MACHINE_CLIENT_ID"), os.Getenv("ELVID_MACHINE_CLIENT_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	elvidClientRequest, err := http.NewRequest("", "api.url", nil)
	if err != nil {
		log.Fatal(err)
	}

	elvidBearerToken.AppendToRequest(elvidClientRequest)
	return elvidClientRequest
}

func elvidApiDemonstration(elvidClientRequest *http.Request) {
	elvid, err := elvid.New()
	if err != nil {
		log.Fatal(err)
	}

	err = elvid.AuthorizeRequest(elvidClientRequest, os.Getenv("ELVID_SCOPE"))
	if err != nil {
		log.Fatal("Invalid access token!")
	}
	log.Println("End of successful request demonstration.")
}
