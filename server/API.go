package main

import (
	"encoding/json"
	"errors"
	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"net/http"
)

// UserIDHeaderKey - key of an http header where the user id is stored
const UserIDHeaderKey = "Mattermost-User-ID"

// APIError - returned error message for api errors
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func writerErrorWrapper(n int, err error) {
	if err != nil {
		mlog.Error(err.Error())
	}
}

// writeAPIError writes api error as json in response
func writeAPIError(w http.ResponseWriter, err *APIError) {
	b, parse_err := json.Marshal(err)
	if parse_err != nil {
		mlog.Error(parse_err.Error())
		return
	}
	w.WriteHeader(err.StatusCode)
	writerErrorWrapper(w.Write(b))
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	resp, parse_err := json.Marshal(data)
	if parse_err != nil {
		mlog.Error(parse_err.Error())
	}
	writerErrorWrapper(w.Write(resp))
}

func writeSuccess(w http.ResponseWriter) {
	writerErrorWrapper(w.Write([]byte("{\"status\": \"OK\"}")))
}

// getUserIDFromRequest reads mattermost user ID from request
func getUserIDFromRequest(r *http.Request) (string, error) {
	userID := r.Header.Get(UserIDHeaderKey)
	if userID == "" {
		return "", errors.New("not authorized")
	}
	return userID, nil
}

// APICallHandler api call handler interface
type APICallHandler func(p *Plugin, c *plugin.Context, w http.ResponseWriter, r *http.Request)

// HandleGetPublicKey handle get public key request
func HandleGetPublicKey(p *Plugin, _ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAPIError(w, &APIError{
			Message:    "Wrong Http method",
			StatusCode: http.StatusMethodNotAllowed,
		})
	}

	userID, err := getUserIDFromRequest(r)
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

	respondWithJSON(w, struct {
		PublicKey []byte `json:"public_key"`
	}{PublicKey: pubKey})
}

// SetPublicKeyRequest - struct for parsing setPublicKey request body
type SetPublicKeyRequest struct {
	PublicKey []byte `json:"public_key"`
}

// HandleSetPublicKey - handle set public key request
func HandleSetPublicKey(p *Plugin, _ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, &APIError{
			Message:    "Wrong Http method",
			StatusCode: http.StatusMethodNotAllowed,
		})
	}

	userID, err := getUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	var request SetPublicKeyRequest
	errno := json.NewDecoder(r.Body).Decode(&request)
	if errno != nil {
		writeAPIError(w, &APIError{Message: "Bad Request.", StatusCode: http.StatusBadRequest})
	}

	errno = p.storePublicKey(request.PublicKey, userID)
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
