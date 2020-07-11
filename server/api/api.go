package api

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/gorilla/schema"

	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/mattermost/mattermost-server/v5/mlog"

	"github.com/gorilla/mux"
)

const (
	// WSEventEncryptionStatusChange web socket broadcast event for status change
	WSEventEncryptionStatusChange = "encryption_status_change"
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
	apiRouter := h.Router.PathPrefix(config.APIPath).Subrouter()
	apiRouter.HandleFunc("/pub_key", h.extractUserIDMiddleware(h.handleGetPublicKey())).Methods(http.MethodGet)
	apiRouter.HandleFunc("/pub_key", h.extractUserIDMiddleware(h.handleSetPublicKey())).Methods(http.MethodPost)

	apiRouter.HandleFunc("/encryption_status", h.extractUserIDMiddleware(h.handleGetEncryptionStatus())).Methods(http.MethodGet)
	apiRouter.HandleFunc("/encryption_status", h.extractUserIDMiddleware(h.handleChangeEncryptionStatus())).Methods(http.MethodPost)
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
func (h *handler) handleGetPublicKeys() handlerWithUserID {
	type request struct {
		UserIDs []string `json:"user_ids"`
	}
	type response struct {
		PublicKeys []string `json:"public_keys"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ string) {

		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}

		userIDs := req.UserIDs

		pubKeys := make([]string, 0)

		if len(userIDs) == 0 {
			h.jsonError(w, Error{Message: "public key doesn't exists", StatusCode: http.StatusNoContent})
			return
		}

		for _, userID := range userIDs {
			pubKey, err := h.an.GetPublicKey(userID)
			if err != nil || pubKey == nil {
				h.jsonError(w, Error{Message: "public key doesn't exists", StatusCode: http.StatusNoContent})
				return
			}
			pubKeys = append(pubKeys, pubKey.String())
		}

		h.respondWithJSON(w, response{PublicKeys: pubKeys})
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

// handleSetPublicKey - returns handlerFunc which
// returns encryption status for channel and current user
func (h *handler) handleGetEncryptionStatus() handlerWithUserID {
	type request struct {
		ChannelID string `schema:"channel_id"`
	}

	type response struct {
		IsEncryptionEnabled bool `json:"is_encryption_enabled"`
	}

	return func(w http.ResponseWriter, r *http.Request, mattermostUserID string) {
		var req request
		err := schema.NewDecoder().Decode(&req, r.URL.Query())
		if err != nil {
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}

		isEnabled := h.an.IsEncryptionEnabledForChannel(req.ChannelID, mattermostUserID)

		h.respondWithJSON(w, response{
			IsEncryptionEnabled: isEnabled,
		})
	}
}

func (h *handler) handleChangeEncryptionStatus() handlerWithUserID {
	type request struct {
		ChannelID string `json:"channel_id"`
		Status    bool   `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request, mattermostUserID string) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.jsonError(w, Error{Message: "Bad Request", StatusCode: http.StatusBadRequest})
			return
		}
		if req.Status {
			unverifiedPlugins := h.an.UnverifiedPlugins()
			if len(unverifiedPlugins) > 0 {
				h.jsonError(w, Error{Message: "Unverified plugins detected", StatusCode: http.StatusForbidden})
				return
			}
		}

		err = h.an.SetEncryptionStatusForChannel(req.ChannelID, mattermostUserID, req.Status)
		if err != nil {
			h.jsonError(w, Error{
				Message:    "Error while changing encryption status",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		h.an.PublishWebSocketEvent(
			WSEventEncryptionStatusChange,

			map[string]interface{}{
				"status": req.Status,
			},

			&model.WebsocketBroadcast{
				UserId:    mattermostUserID,
				ChannelId: req.ChannelID,
			},
		)
		h.respondWithSuccess(w)
	}
}
