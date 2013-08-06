package main

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonServerReturnsJsonDocumentWithRightHeaders(t *testing.T) {

	req := new(http.Request)

	builder := func(req *http.Request) interface{} {
		return map[string]interface{}{
			"test": "test_value",
		}
	}

	jsonServer := JsonServer(builder)
	responseWriter := httptest.NewRecorder()

	// call to test
	jsonServer(responseWriter, req)

	contentType := responseWriter.Header().Get("Content-Type")
	assert.Equal(t, contentType, "application/json")

	assert.Equal(t, "{\"test\":\"test_value\"}\n", responseWriter.Body.String())

}
