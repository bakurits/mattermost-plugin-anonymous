// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"encoding/json"
	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/mattermost/mattermost-server/v5/mlog"
	"net/http"

	"github.com/gorilla/mux"
)

// Error - returned error message for api errors
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// Handler is an http.Handler for all plugin HTTP endpoints
type Handler struct {
	*mux.Router
}

// InitRouter initializes the router.
func NewHTTPHandler() *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	apiRouter := h.Router.PathPrefix(config.PathAPI).Subrouter()
	apiRouter.HandleFunc("/pub_key", h.HandleGetPublicKey).Methods("GET")
	apiRouter.HandleFunc("/pub_key", h.HandleSetPublicKey).Methods("POST")
	return h
}

func (h *Handler) jsonError(w http.ResponseWriter, err Error) {
	w.WriteHeader(err.StatusCode)
	h.respondWithJSON(w, err)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		mlog.Error(err.Error())
	}
}

func (h *Handler) respondWithSuccess(w http.ResponseWriter) {
	_, err := w.Write([]byte("{\"status\": \"OK\"}"))
	if err != nil {
		mlog.Error(err.Error())
	}
}

// HandleGetPublicKey handle get public key request
func (h *Handler) HandleGetPublicKey(w http.ResponseWriter, r *http.Request) {
	anonymousApi := anonymous.FromContext(r.Context())
	pubKey, err := anonymousApi.GetPublicKey()
	if err != nil {
		h.jsonError(w, Error{
			Message:    "public key doesn't exists",
			StatusCode: http.StatusNoContent,
		})
	}

	h.respondWithJSON(w, struct {
		PublicKey []byte `json:"public_key"`
	}{PublicKey: pubKey})
}

// SetPublicKeyRequest - struct for parsing setPublicKey request body
type SetPublicKeyRequest struct {
	PublicKey []byte `json:"public_key"`
}

// HandleSetPublicKey - handle set public key request
func (h *Handler) HandleSetPublicKey(w http.ResponseWriter, r *http.Request) {
	anonymousApi := anonymous.FromContext(r.Context())

	var request SetPublicKeyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.jsonError(w, Error{Message: "Bad Request.", StatusCode: http.StatusBadRequest})
	}

	err = anonymousApi.StorePublicKey(request.PublicKey)
	if err != nil {
		h.jsonError(w, Error{Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
	}

	h.respondWithSuccess(w)
}
