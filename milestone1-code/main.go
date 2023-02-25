package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aleury/secretly/data"
	"github.com/aleury/secretly/handlers"
)

func main() {
	dataFilePath := os.Getenv("DATA_FILE_PATH")
	if dataFilePath == "" {
		log.Fatalln("expected DATA_FILE_PATH to be configured")
	}

	store, err := data.CreateSecretStore(dataFilePath)
	if err != nil {
		log.Fatalln("failed to initialize secret store")
	}

	mux := http.NewServeMux()
	handlers.Register(mux, store)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
