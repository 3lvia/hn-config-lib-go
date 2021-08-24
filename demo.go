package main

import (
	"log"
	"net/http"
	"os"

	"github.com/3lvia/hn-config-lib-go/hid"
	"github.com/3lvia/hn-config-lib-go/vault"
)

// This is a runnable version of the `example_test.go` file.
// It needs appropriate environment variables before executing.
// See README.md for more information.
func main() {
	vaultDemonstation()
	// requestInstance := elvidClientDemonstation()
	// elvidApiDemonstration(requestInstance)
}

func vaultDemonstation() {
	vaultInstance, err := vault.New()
	if err != nil {
		log.Fatal(err)
	}

	// Example:
	// export VAULT_SECRET_PATH_VALUE="onetime/kv/data/azurerm-redis-cache/onetime"
	log.Println("Vault Secret Path: ", os.Getenv("VAULT_SECRET_PATH_VALUE"))
	secretInstance, err := vaultInstance.GetSecret(os.Getenv("VAULT_SECRET_PATH_VALUE"))
	if err != nil {
		log.Fatal(err)
	}

	dataInstance := secretInstance.GetData()
	log.Println("\tVault Secret Data: ", dataInstance)
	log.Println("\t\tHostName: ", dataInstance["hostname"])
	log.Println("\t\tPrimary Access Key: ", dataInstance["primary-access-key"])
	log.Println("\t\tPrimary Connection String: ", dataInstance["primary-connection-string"])
	log.Println("\t\tSecondary Access Key: ", dataInstance["secondary-access-key"])
	log.Println("\t\tSecondary Connection String: ", dataInstance["secondary-connection-string"])

	log.Println("GetLeaseDuration: ", secretInstance.GetLeaseDuration())
	log.Println("GetLeaseID: ", secretInstance.GetLeaseID())
	log.Println("GetMetadata: ", secretInstance.GetMetadata())
	log.Println("\tCreatedTime: ", secretInstance.GetMetadata()["created_time"])
	log.Println("\tDeletionTime: ", secretInstance.GetMetadata()["deletion_time"])
	log.Println("\tDestroyed? ", secretInstance.GetMetadata()["destroyed"])
	log.Println("\tVersion: ", secretInstance.GetMetadata()["version"])
	log.Println("GetRequestID: ", secretInstance.GetRequestID())
	log.Println("IsRenewable: ", secretInstance.IsRenewable())
}

func elvidClientDemonstation() *http.Request {
	// Make reusable HID item
	myHID, err := hid.New()
	if err != nil {
		log.Fatal(err)
	}

	// Get a bearer token from HID
	myToken, err := myHID.GetToken(os.Getenv("TEST_HID_ID"), os.Getenv("TEST_HID_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	// Make http.Request as usual
	myRequest, err := http.NewRequest("", "api.url", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add bearer token to http request header
	myToken.AppendToRequest(myRequest)

	// Send token to API with requests
	return myRequest
}

func elvidApiDemonstration(myRequest *http.Request) {
	// Make reusable HID item
	myHID, err := hid.New()
	if err != nil {
		log.Fatal(err)
	}

	// Verify if token is valid. Invalid token throws an error
	err = myHID.AuthorizeRequest(myRequest, os.Getenv("TEST_AUDIENCE"), os.Getenv("TEST_SCOPE"))
	if err != nil {
		log.Fatal("Token is invalid")
	}

	// Handle the request
	log.Println("The request was successfull")
}
