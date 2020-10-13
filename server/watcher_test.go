package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	plugin := Plugin{
		configuration: &configuration{
			Token: "test",
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.NotNil(t, result)
	defer result.Body.Close()

	bodyBytes, err := ioutil.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Unauthorized\n", string(bodyBytes))
}
