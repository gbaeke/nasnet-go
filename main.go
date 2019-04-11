package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gbaeke/nasnet-go/handlers"
	"github.com/mholt/certmagic"
)

func main() {
	//check environment variables to enable SSL
	port := getEnv("PORT", "9090")
	sslEnabled := getEnv("ssl", "false")
	hostName := getEnv("hostname", "")
	if hostName == "" && sslEnabled == "true" {
		log.Fatalln("Specify hostname environment variable when SSL is on")
	}
	stagingCA := getEnv("staging", "true")

	if sslEnabled == "true" {
		log.Println("SSL")

		// certmagic
		certmagic.Agreed = true
		certmagic.Email = "mail@mail.com"
		if stagingCA == "false" {
			log.Println("Using production CA")
			certmagic.CA = certmagic.LetsEncryptProductionCA
		} else {
			log.Println("Using staging CA")
			certmagic.CA = certmagic.LetsEncryptStagingCA
		}

		mux := handlers.RoutesMux()
		err := certmagic.HTTPS([]string{hostName}, mux)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		handlers.Routes()
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
