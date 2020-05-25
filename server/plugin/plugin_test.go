package plugin_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous/mock"

	"github.com/bakurits/mattermost-plugin-anonymous/server/api"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/plugin"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	mockStore "github.com/bakurits/mattermost-plugin-anonymous/server/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_plugin_ServeHTTP_GetPublicKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	is := assert.New(t)

	storeMock := mockStore.NewMockStore(ctrl)
	storeMock.EXPECT().LoadUser("key_in").Return(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}, nil)
	storeMock.EXPECT().LoadUser(gomock.Any()).Return(nil, errors.New("some error"))

	storeMock.EXPECT().StoreUser(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}).Return(nil)
	storeMock.EXPECT().StoreUser(gomock.Any()).Return(errors.New("some error"))

	httpTest := test.HTTPTest{
		Assertions: is,
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
				URL:    fmt.Sprintf("%s/pub_key", config.APIPath),
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
				URL:    fmt.Sprintf("%s/pub_key?user_id=%s", config.APIPath, "asd"),
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
				URL:    fmt.Sprintf("%s/pub_key?user_id=%s", config.APIPath, "key_in"),
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
			p := plugin.NewWithStore(storeMock, nil)
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

	storeMock := mockStore.NewMockStore(ctrl)
	storeMock.EXPECT().StoreUser(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}).Return(nil)
	storeMock.EXPECT().StoreUser(gomock.Any()).Return(errors.New("some error"))

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
				URL:    fmt.Sprintf("%s/pub_key", config.APIPath),
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
				URL:    fmt.Sprintf("%s/pub_key", config.APIPath),
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
				URL:    fmt.Sprintf("%s/pub_key", config.APIPath),
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
				URL:    fmt.Sprintf("%s/pub_key", config.APIPath),
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
			p := plugin.NewWithStore(storeMock, nil)
			req := httpTest.CreateHTTPRequest(tt.request)
			req.Header.Add("Mattermost-User-ID", tt.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(nil, rr, req)
			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}

func Test_plugin_ServeHTTP_GetEncryptionStatus(t *testing.T) {

	ctrl := gomock.NewController(t)
	is := assert.New(t)

	anonymousMock := mock.NewMockAnonymous(ctrl)
	anonymousMock.EXPECT().IsEncryptionEnabledForChannel("general", "in_general").Return(true).AnyTimes()
	anonymousMock.EXPECT().IsEncryptionEnabledForChannel(gomock.Any(), gomock.Any()).Return(false).AnyTimes()

	httpTest := test.HTTPTest{
		Assertions: is,
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
			name: "not authorized test",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/encryption_status?channel_id=%s", config.APIPath, "general"),
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
			name: "test bad request",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/encryption_status?channel_d=%s", config.APIPath, "general"),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Bad Request",
					StatusCode: http.StatusBadRequest,
				},
			},
			userID: "storing_err_user",
		},
		{
			name: "test disabled",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/encryption_status?channel_id=%s", config.APIPath, "general"),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					IsEncryptionEnabled bool `json:"is_encryption_enabled"`
				}{IsEncryptionEnabled: false},
			},
			userID: "storing_err_user",
		},
		{
			name: "test enabled",
			request: test.Request{
				Method: "GET",
				URL:    fmt.Sprintf("%s/encryption_status?channel_id=%s", config.APIPath, "general"),
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					IsEncryptionEnabled bool `json:"is_encryption_enabled"`
				}{IsEncryptionEnabled: true},
			},
			userID: "in_general",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := plugin.NewWithAnonymous(anonymousMock, nil)
			req := httpTest.CreateHTTPRequest(tt.request)
			req.Header.Add("Mattermost-User-ID", tt.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(nil, rr, req)
			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}

func Test_plugin_ServeHTTP_ChangeEncryptionStatus(t *testing.T) {

	ctrl := gomock.NewController(t)
	is := assert.New(t)

	anonymousMock := mock.NewMockAnonymous(ctrl)
	anonymousMock.EXPECT().SetEncryptionStatusForChannel("general", "storing_err_user", gomock.Any()).Return(errors.New("some error")).AnyTimes()
	anonymousMock.EXPECT().SetEncryptionStatusForChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	anonymousMock.EXPECT().PublishWebSocketEvent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	httpTest := test.HTTPTest{
		Assertions: is,
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
			name: "not authorized test",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/encryption_status", config.APIPath),
				Body: struct {
					ChannelID string `json:"channel_id"`
					Status    bool   `json:"status"`
				}{
					ChannelID: "general",
					Status:    false,
				},
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
			name: "test bad request",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/encryption_status", config.APIPath),
				Body: struct {
					ChannelID int  `json:"channel_id"`
					Status    bool `json:"status"`
				}{
					ChannelID: 1,
					Status:    false,
				},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Bad Request",
					StatusCode: http.StatusBadRequest,
				},
			},
			userID: "storing_err_user",
		},
		{
			name: "test can't change status",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/encryption_status", config.APIPath),
				Body: struct {
					ChannelID string `json:"channel_id"`
					Status    bool   `json:"status"`
				}{
					ChannelID: "general",
					Status:    false,
				},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Error while changing encryption status",
					StatusCode: http.StatusBadRequest,
				},
			},
			userID: "storing_err_user",
		},
		{
			name: "test success",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/encryption_status", config.APIPath),
				Body: struct {
					ChannelID string `json:"channel_id"`
					Status    bool   `json:"status"`
				}{
					ChannelID: "general",
					Status:    false,
				},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					Status string `json:"status"`
				}{Status: "OK"},
			},
			userID: "in_general",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := plugin.NewWithAnonymous(anonymousMock, nil)
			req := httpTest.CreateHTTPRequest(tt.request)
			req.Header.Add("Mattermost-User-ID", tt.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(nil, rr, req)
			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}
