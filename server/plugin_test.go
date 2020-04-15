package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	tassert := assert.New(t)
	plugin := Plugin{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	tassert.NotNil(result)
	bodyBytes, err := ioutil.ReadAll(result.Body)
	tassert.Nil(err)
	bodyString := string(bodyBytes)

	tassert.Equal("Hello, world", bodyString)
}
