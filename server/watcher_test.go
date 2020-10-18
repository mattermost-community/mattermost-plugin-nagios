package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	var plugin Plugin

	t.Run("invalid method", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		plugin.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)
		defer result.Body.Close()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)

		b, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)

		assert.Equal(t, "Bad Request\n", string(b))
	})

	t.Run("no token", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", nil)

		plugin.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)
		defer result.Body.Close()

		assert.Equal(t, http.StatusNotImplemented, result.StatusCode)

		b, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)

		assert.Equal(t, "This functionality is not configured.\n", string(b))
	})
}
