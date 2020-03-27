// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/mattermost/mattermost-server/v5/mlog"

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

// NewHTTPHandler initializes the router.
func NewHTTPHandler() *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	apiRouter := h.Router.PathPrefix(config.PathAPI).Subrouter()
	apiRouter.HandleFunc("/pub_key", h.handleGetPublicKey()).Methods("GET")
	apiRouter.HandleFunc("/pub_key", h.handleSetPublicKey()).Methods("POST")
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

// handleGetPublicKey handle get public key request
func (h *Handler) handleGetPublicKey() http.HandlerFunc {

	type response struct {
		PublicKey []byte `json:"public_key"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		anonymousAPI := anonymous.FromContext(r.Context())
		pubKey, err := anonymousAPI.GetPublicKey()
		if err != nil {
			h.jsonError(w, Error{
				Message:    "public key doesn't exists",
				StatusCode: http.StatusNoContent,
			})
		}

		h.respondWithJSON(w, response{PublicKey: pubKey})
	}
}

// handleSetPublicKey - handle set public key request
func (h *Handler) handleSetPublicKey() http.HandlerFunc {

	type request struct {
		PublicKey []byte `json:"public_key"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		anonymousAPI := anonymous.FromContext(r.Context())

		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.jsonError(w, Error{Message: "Bad Request.", StatusCode: http.StatusBadRequest})
		}

		err = anonymousAPI.StorePublicKey(req.PublicKey)
		if err != nil {
			h.jsonError(w, Error{Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
		}

		h.respondWithSuccess(w)
	}
}
