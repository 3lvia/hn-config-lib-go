package main

import (
	"fmt"
	"log"
	"time"

	"github.com/3lvia/hn-config-lib-go/vault"
)

func main() {
	myVault, err := vault.New()
	if err != nil {
		log.Fatal(err)
	}

	// Get a secret from the vault
	err = myVault.SetupAndRenewGcpCredentials("monitoring", "super", 300)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println(time.Now())
		time.Sleep(60 * time.Second)
		// for i := 1; i < 60; i++ {
		// 	time.Sleep(60 * time.Second)
		// 	fmt.Print(".")
		// }
	}
}
