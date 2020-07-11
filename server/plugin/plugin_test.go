package plugin_test

import (
	"errors"
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/api"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/plugin"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	storeMock "github.com/bakurits/mattermost-plugin-anonymous/server/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_plugin_ServeHTTP_GetPublicKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)

	mockStore := storeMock.NewMockStore(ctrl)
	mockStore.EXPECT().LoadUser("key_in").Return(&store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1, 1},
	}, nil).AnyTimes()
	mockStore.EXPECT().LoadUser("key_in2").Return(&store.User{
		MattermostUserID: "key_in2",
		PublicKey:        []byte{2, 2},
	}, nil).AnyTimes()
	mockStore.EXPECT().LoadUser(gomock.Any()).Return(nil, errors.New("some error")).AnyTimes()

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
				URL:    fmt.Sprintf("%s/pub_keys", config.APIPath),
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
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_keys", config.APIPath),
				Body: struct {
					UserIDs []string `json:"user_ids"`
				}{UserIDs: []string{"abc"}},
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
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_keys", config.APIPath),
				Body: struct {
					UserIDs []string `json:"user_ids"`
				}{UserIDs: []string{"key_in"}},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					PublicKeys []string `json:"public_keys"`
				}{PublicKeys: []string{crypto.PublicKey([]byte{1, 1}).String()}},
			},
			userID: "abc",
		},
		{
			name: "test successs",
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/pub_keys", config.APIPath),
				Body: struct {
					UserIDs []string `json:"user_ids"`
				}{UserIDs: []string{"key_in", "key_in2"}},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					PublicKeys []string `json:"public_keys"`
				}{PublicKeys: []string{crypto.PublicKey([]byte{1, 1}).String(), crypto.PublicKey([]byte{2, 2}).String()}},
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
			p := plugin.NewWithStore(mockStore, nil)
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

	anonymousMock1 := mock.NewMockAnonymous(ctrl)
	anonymousMock1.EXPECT().SetEncryptionStatusForChannel("general", "storing_err_user", gomock.Any()).Return(errors.New("some error")).AnyTimes()
	anonymousMock1.EXPECT().SetEncryptionStatusForChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	anonymousMock1.EXPECT().PublishWebSocketEvent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	anonymousMock1.EXPECT().UnverifiedPlugins().Return([]anonymous.PluginIdentifier{}).AnyTimes()

	anonymousMock2 := mock.NewMockAnonymous(ctrl)
	anonymousMock2.EXPECT().SetEncryptionStatusForChannel("general", "storing_err_user", gomock.Any()).Return(errors.New("some error")).AnyTimes()
	anonymousMock2.EXPECT().SetEncryptionStatusForChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	anonymousMock2.EXPECT().PublishWebSocketEvent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	anonymousMock2.EXPECT().UnverifiedPlugins().Return([]anonymous.PluginIdentifier{{ID: "123", Version: "123"}}).AnyTimes()

	httpTest := test.HTTPTest{
		Assertions: is,
		Encoder:    test.EncodeJSON,
	}

	tests := []struct {
		name             string
		request          test.Request
		expectedResponse test.ExpectedResponse
		an               anonymous.Anonymous
		userID           string
	}{
		{
			name: "not authorized test",
			an:   anonymousMock1,
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
			an:   anonymousMock1,
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
			an:   anonymousMock1,
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
			name: "test unverified plugins",
			an:   anonymousMock2,
			request: test.Request{
				Method: "POST",
				URL:    fmt.Sprintf("%s/encryption_status", config.APIPath),
				Body: struct {
					ChannelID string `json:"channel_id"`
					Status    bool   `json:"status"`
				}{
					ChannelID: "general",
					Status:    true,
				},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusForbidden,
				ResponseType: test.ContentTypeJSON,
				Body: api.Error{
					Message:    "Unverified plugins detected",
					StatusCode: http.StatusForbidden,
				},
			},
			userID: "in_general",
		},
		{
			name: "test success",
			an:   anonymousMock1,
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
			p := plugin.NewWithAnonymous(tt.an, nil)
			req := httpTest.CreateHTTPRequest(tt.request)
			req.Header.Add("Mattermost-User-ID", tt.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(nil, rr, req)
			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}
