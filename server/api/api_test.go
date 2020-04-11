package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	mockAnonymous "github.com/bakurits/mattermost-plugin-anonymous/server/anonymous/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_handler_handleGetPublicKey(t *testing.T) {

	tassert := assert.New(t)

	ctrl := gomock.NewController(t)

	anonymousMock := mockAnonymous.NewMockAnonymous(ctrl)
	anonymousMock.EXPECT().GetPublicKey("key_in").Return([]byte{1, 1}, nil)
	anonymousMock.EXPECT().GetPublicKey(gomock.Any()).Return(nil, errors.New("some error"))

	handler := newHandler().handleGetPublicKey()

	httpTest := test.HTTPTest{
		Assertions: tassert,
		Encoder:    test.EncodeJSON,
	}

	tests := []struct {
		name             string
		request          test.Request
		expectedResponse test.ExpectedResponse
	}{
		{
			name: "test bad request",
			request: test.Request{
				Method: "GET",
				URL:    "/api/v1/pub_key",
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: Error{
					Message:    "Bad Request",
					StatusCode: http.StatusBadRequest,
				},
			},
		},
		{
			name: "test not registered user",
			request: test.Request{
				Method: "GET",
				URL:    "/api/v1/pub_key?user_id=asd",
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusNoContent,
				ResponseType: test.ContentTypeJSON,
				Body: Error{
					Message:    "public key doesn't exists",
					StatusCode: http.StatusNoContent,
				},
			},
		},
		{
			name: "test success",
			request: test.Request{
				Method: "GET",
				URL:    "/api/v1/pub_key?user_id=key_in",
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusOK,
				ResponseType: test.ContentTypeJSON,
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: crypto.PublicKey([]byte{1, 1}).String()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httpTest.CreateHTTPRequest(tt.request)

			ctx := req.Context()
			ctx = anonymous.Context(ctx, anonymousMock)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}

func Test_handler_handleSetPublicKey(t *testing.T) {
	tassert := assert.New(t)

	ctrl := gomock.NewController(t)

	anonymousMock := mockAnonymous.NewMockAnonymous(ctrl)
	anonymousMock.EXPECT().StorePublicKey(crypto.PublicKey([]byte{1, 1})).Return(nil)
	anonymousMock.EXPECT().StorePublicKey(gomock.Any()).Return(errors.New("some error"))

	handler := newHandler().handleSetPublicKey()

	httpTest := test.HTTPTest{
		Assertions: tassert,
		Encoder:    test.EncodeJSON,
	}

	tests := []struct {
		name             string
		request          test.Request
		expectedResponse test.ExpectedResponse
	}{
		{
			name: "test bad request",
			request: test.Request{
				Method: "POST",
				URL:    "/api/v1/pub_key",
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: Error{
					Message:    "Bad Request",
					StatusCode: http.StatusBadRequest,
				},
			},
		},
		{
			name: "test bad public key",
			request: test.Request{
				Method: "GET",
				URL:    "/api/v1/pub_key",
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: "~~1"},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusBadRequest,
				ResponseType: test.ContentTypeJSON,
				Body: Error{
					Message:    "Public key format is incorrect",
					StatusCode: http.StatusBadRequest,
				},
			},
		},
		{
			name: "test not authorized",
			request: test.Request{
				Method: "GET",
				URL:    "/api/v1/pub_key",
				Body: struct {
					PublicKey string `json:"public_key"`
				}{PublicKey: crypto.PublicKey([]byte{1, 1, 1}).String()},
			},
			expectedResponse: test.ExpectedResponse{
				StatusCode:   http.StatusUnauthorized,
				ResponseType: test.ContentTypeJSON,
				Body: Error{
					Message:    "Not Authorized",
					StatusCode: http.StatusUnauthorized,
				},
			},
		},
		{
			name: "test success",
			request: test.Request{
				Method: "GET",
				URL:    "/api/v1/pub_key",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httpTest.CreateHTTPRequest(tt.request)

			ctx := req.Context()
			ctx = anonymous.Context(ctx, anonymousMock)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			httpTest.CompareHTTPResponse(rr, tt.expectedResponse)
		})
	}
}
