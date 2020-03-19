package main

import (
	"encoding/json"
	"errors"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"net/http"
)

const USER_ID_HEADER_KEY = "Mattermost-User-ID"

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// writeAPIError writes api error as json in response
func writeAPIError(w http.ResponseWriter, err *APIError) {
	b, _ := json.Marshal(err)
	w.WriteHeader(err.StatusCode)
	_, _ = w.Write(b)
}

func respondWithJson(w http.ResponseWriter, data interface{}) {
	resp, _ := json.Marshal(data)
	_, _ = w.Write(resp)
}

func writeSuccess(w http.ResponseWriter) {
	_, _ = w.Write([]byte("{\"status\": \"OK\"}"))
}

// getUserIdFromRequest reads mattermost user ID from request
func getUserIdFromRequest(r *http.Request) (string, error) {
	userID := r.Header.Get(USER_ID_HEADER_KEY)
	if userID == "" {
		return "", errors.New("not authorized")
	}
	return userID, nil
}

// APICallHandler api call handler interface
type APICallHandler func(p *Plugin, c *plugin.Context, w http.ResponseWriter, r *http.Request)

func HandleGetPublicKey(p *Plugin, _ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	pubKey, err := p.getPublicKey(userID)
	if err != nil {
		writeAPIError(w, &APIError{
			Message:    "public key doesn't exists",
			StatusCode: http.StatusNoContent,
		})
	}

	respondWithJson(w, struct {
		PublicKey string `json:"public_key"`
	}{PublicKey: pubKey})
}

func HandleSetPublicKey(p *Plugin, _ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	pubKey := r.FormValue("term")

	errno := p.storePublicKey(pubKey, userID)
	if errno != nil {
		writeAPIError(w, &APIError{Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
	}

	writeSuccess(w)
}

var (
	apiHandlers = map[string]APICallHandler{
		"/api/pub_key/get": HandleGetPublicKey,
		"/api/pub_key/set": HandleSetPublicKey,
	}
)

// ServeHTTP serves API calls
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if f, ok := apiHandlers[r.URL.Path]; ok {
		f(p, c, w, r)
	} else {
		http.NotFound(w, r)
	}
}
