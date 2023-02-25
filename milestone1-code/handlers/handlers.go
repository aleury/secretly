package handlers

import (
	"net/http"

	"github.com/aleury/secretly/data"
)

func Register(mux *http.ServeMux, store *data.SecretStore) {
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler(store))
}
