package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createHttpClientWithDefaultHeaders(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		client := createHTTPClientWithDefaultHeaders(nil)
		assert.NotNil(t, client.Transport)
		proxyTransport, casted := client.Transport.(*proxyTransport)
		assert.True(t, casted)
		assert.Nil(t, proxyTransport.defaultHeaders)
	})

	t.Run("Some headers", func(t *testing.T) {
		m := map[string][]string{
			"Header1": {"val1", "val2"},
		}
		client := createHTTPClientWithDefaultHeaders(m)
		assert.NotNil(t, client.Transport)
		proxyTransport, casted := client.Transport.(*proxyTransport)
		assert.True(t, casted)
		assert.Equal(t, m, proxyTransport.defaultHeaders)
	})
}

func Test_applyDefaultHeaders(t *testing.T) {
	input := map[string][]string{
		"Header1": {"val1", "val2"},
	}

	r := &http.Request{}
	applyDefaultHeaders(r, input)

	assert.Equal(t, http.Header(input), r.Header)
}
