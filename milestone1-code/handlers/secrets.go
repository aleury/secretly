package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aleury/secretly/data"
	"github.com/aleury/secretly/types"
)

func secretHandler(store *data.SecretStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getSecretHandler(w, r, store)
		} else if r.Method == "POST" {
			createSecretHandler(w, r, store)
		} else {
			writeJson(w, http.StatusMethodNotAllowed, types.Error{Message: "method not allowed"})
		}
	}
}

func getSecretHandler(w http.ResponseWriter, r *http.Request, store *data.SecretStore) {
	secretId, err := parseSecretId(r)
	if err != nil {
		writeJson(w, http.StatusBadRequest, types.GetSecretResponse{})
		return
	}

	secret, err := store.Get(secretId)
	if err != nil {
		writeJson(w, http.StatusNotFound, types.GetSecretResponse{})
		return
	}

	writeJson(w, http.StatusOK, types.GetSecretResponse{Data: secret})
}

func createSecretHandler(w http.ResponseWriter, r *http.Request, store *data.SecretStore) {
	var req types.CreateSecretRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJson(w, http.StatusBadRequest, types.Error{Message: "invalid request body"})
		return
	}

	id, err := store.Put(req.PlainText)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, types.Error{Message: "failed to create secret"})
		return
	}

	writeJson(w, http.StatusCreated, types.CreateSecretResponse{ID: id})
}

func parseSecretId(r *http.Request) (string, error) {
	if r.URL.Path == "/" {
		return "", errors.New("secret id required")
	}
	return r.URL.Path[1:], nil
}

func writeJson(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
