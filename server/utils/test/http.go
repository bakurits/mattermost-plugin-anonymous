package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

// ContentType type for http Content-Type
type ContentType string

const (
	//ContentTypeJSON string associated with content-type json
	ContentTypeJSON = "application/json"
	//ContentTypeHTML string associated with content-type html
	ContentTypeHTML = "text/html"
)

// Request stores http request basic data
type Request struct {
	Method string
	URL    string
	Body   interface{}
}

// CreateHTTPRequest creates http request with basic data
func (test *HTTPTest) CreateHTTPRequest(request Request) *http.Request {
	var body io.Reader
	data, err := test.Encoder(request.Body)
	test.NoError(err)
	body = bytes.NewBuffer(data)

	req, err := http.NewRequest(request.Method, request.URL, body)
	test.NoError(err, "Error while creating request")
	return req
}

// ExpectedResponse stores expected responce basic data
type ExpectedResponse struct {
	StatusCode   int
	ResponseType ContentType
	Body         interface{}
}

// HTTPTest encapsulates data for testing needs
type HTTPTest struct {
	*assert.Assertions
	Encoder func(interface{}) ([]byte, error)
}

// CompareHTTPResponse compares expected responce with real one
func (test *HTTPTest) CompareHTTPResponse(rr *httptest.ResponseRecorder, expected ExpectedResponse) {

	test.Equal(expected.StatusCode, rr.Code, "Http status codes are different")

	expectedBody, err := test.Encoder(expected.Body)
	test.NoError(err)

	test.Equal(string(expected.ResponseType), rr.Header().Get("Content-Type"))

	gotBody := rr.Body.Bytes()

	test.Equal(expectedBody, gotBody)
}
