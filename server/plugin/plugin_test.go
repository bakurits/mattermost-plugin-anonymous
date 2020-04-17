package plugin_test

import (
	"errors"
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/api"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/plugin"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	storeMock "github.com/bakurits/mattermost-plugin-anonymous/server/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_plugin_ServeHTTP_GetPublicKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)

	mockStore := storeMock.NewMockStore(ctrl)
	mockStore.EXPECT().LoadUser("key_in").Return(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}, nil)
	mockStore.EXPECT().LoadUser(gomock.Any()).Return(nil, errors.New("some error"))

	mockStore.EXPECT().StoreUser(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}).Return(nil)
	mockStore.EXPECT().StoreUser(gomock.Any()).Return(errors.New("some error"))

	httpTest := test.HTTPTest{
		Assertions: tassert,
		Encoder:    test.EncodeJSON,
	}

	tests := []struct {
		name             string
		request          test.Request
		expectedResponse test.ExpectedResponse
		config           *config.Config
		userID           string
	}{
		{
			name: "test bad request",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/pub_key", config.PathAPI),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Bad Request",
					StatusCode: http.StatusBadRequest,
				},
			},
			userID: "abc",
		},
		{
			name: "test not registered user",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/pub_key?user_id=%s", config.PathAPI, "asd"),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusNoContent,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "public key doesn't exists",
					StatusCode: http.StatusNoContent,
				},
			},
			userID: "abc",
		},
		{
			name: "test success",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/pub_key?user_id=%s", config.PathAPI, "key_in"),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: crypto.PublicKey([]byte{1, 1}).String()},
			},
			userID: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := plugin.NewWithStore(mockStore, nil)
			req := httpTest.CreateHTTPRequest(tt.request)
			req.Header.Add("Mattermost-User-ID", tt.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(nil, rr, req)
			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}

func Test_plugin_ServeHTTP_SetPublicKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)

	mockStore := storeMock.NewMockStore(ctrl)
	mockStore.EXPECT().StoreUser(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}).Return(nil)
	mockStore.EXPECT().StoreUser(gomock.Any()).Return(errors.New("some error"))

	httpTest := test.HTTPTest{
		Assertions: tassert,
		Encoder:    test.EncodeJSON,
	}

	tests := []struct {
		name             string
		request          test.Request
		expectedResponse test.ExpectedResponse
		config           *config.Config
		userID           string
	}{
		{
			name: "test bad request",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_key", config.PathAPI),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Bad Request",
					StatusCode: http.StatusBadRequest,
				},
			},
			userID: "key_in",
		},
		{
			name: "test bad public key",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_key", config.PathAPI),
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: "~~1"},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Public key format is incorrect",
					StatusCode: http.StatusBadRequest,
				},
			},
			userID: "key_in",
		},
		{
			name: "test not authorized",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_key", config.PathAPI),
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: crypto.PublicKey([]byte{1, 1, 1}).String()},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusUnauthorized,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Not Authorized",
					StatusCode: http.StatusUnauthorized,
				},
			},
			userID: "",
		},
		{
			name: "test success",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_key", config.PathAPI),
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: crypto.PublicKey([]byte{1, 1}).String()},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					Status string `json:"status"`
				}{Status: "OK"},
			},
			userID: "key_in",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := plugin.NewWithStore(mockStore, nil)
			req := httpTest.CreateHTTPRequest(tt.request)
			req.Header.Add("Mattermost-User-ID", tt.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(nil, rr, req)
			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}
