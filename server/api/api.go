// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"

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

// handler is an http.handler for all plugin HTTP endpoints
type handler struct {
	*mux.Router
	an anonymous.Anonymous
}

type handlerWithUserID func(w http.ResponseWriter, r *http.Request, userID string)

// NewHTTPHandler initializes the router.
func NewHTTPHandler(an anonymous.Anonymous) http.Handler {
	return newHandler(an)
}

func newHandler(an anonymous.Anonymous) *handler {
	h := &handler{
		Router: mux.NewRouter(),
		an:     an,
	}
	apiRouter := h.Router.PathPrefix(config.PathAPI).Subrouter()
	apiRouter.HandleFunc("/pub_key", h.extractUserIDMiddleware(h.handleGetPublicKey())).Methods("GET")
	apiRouter.HandleFunc("/pub_key", h.extractUserIDMiddleware(h.handleSetPublicKey())).Methods("POST")
	return h
}

func (h *handler) jsonError(w http.ResponseWriter, err Error) {
	w.WriteHeader(err.StatusCode)
	h.respondWithJSON(w, err)
}

func (h *handler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		mlog.Error(err.Error())
	}
}

func (h *handler) respondWithSuccess(w http.ResponseWriter) {
	h.respondWithJSON(w, struct {
		Status string `json:"status"`
	}{Status: "OK"})
}

func (h *handler) extractUserIDMiddleware(handler handlerWithUserID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mattermostUserID := r.Header.Get("Mattermost-User-ID")
		if mattermostUserID == "" {
			h.jsonError(w, Error{Message: "Not Authorized", StatusCode: http.StatusUnauthorized})
			return
		}
		handler(w, r, mattermostUserID)
	}
}

// handleGetPublicKey handle get public key request
func (h *handler) handleGetPublicKey() handlerWithUserID {
	type request struct {
		UserID string `schema:"user_id"`
	}
	type response struct {
		PublicKey string `json:"public_key"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ string) {
		var req request
		err := schema.NewDecoder().Decode(&req, r.URL.Query())
		if err != nil || req.UserID == "" {
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}

		pubKey, err := h.an.GetPublicKey(req.UserID)
		if err != nil || pubKey == nil {
			h.jsonError(w, Error{Message: "public key doesn't exists", StatusCode: http.StatusNoContent})
			return
		}

		h.respondWithJSON(w, response{PublicKey: pubKey.String()})
	}
}

// handleSetPublicKey - handle set public key request
func (h *handler) handleSetPublicKey() handlerWithUserID {
	type request struct {
		PublicKey string `json:"public_key"`
	}

	return func(w http.ResponseWriter, r *http.Request, mattermostUserID string) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}

		pubKey, err := crypto.PublicKeyFromString(req.PublicKey)
		if err != nil {
			h.jsonError(w, Error{Message: "Public key format is incorrect", StatusCode: http.StatusBadRequest})
			return
		}

		err = h.an.StorePublicKey(mattermostUserID, pubKey)
		if err != nil {
			h.jsonError(w, Error{Message: "Not Authorized", StatusCode: http.StatusUnauthorized})
			return
		}

		h.respondWithSuccess(w)
	}
}
