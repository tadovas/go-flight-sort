package helpers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonHandlerWritesCorrectJsonOutput(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/<irrelevant>", nil)
	resp := httptest.NewRecorder()

	httpHandler := JsonHandler(func(req *http.Request) (JsonValue, HttpCode, error) {
		return struct {
				SomeField string `json:"some_field"`
			}{
				SomeField: "abc",
			},
			http.StatusOK,
			nil
	})

	httpHandler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-type"))
	assert.JSONEq(t, `{ "some_field" : "abc"}`, resp.Body.String())
}

func TestJsonHandlerWritesCorrectError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/<irrelevant>", nil)
	resp := httptest.NewRecorder()

	httpHandler := JsonHandler(func(req *http.Request) (JsonValue, HttpCode, error) {
		return nil, http.StatusBadRequest, errors.New("big error")
	})

	httpHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-type"))
	assert.JSONEq(t, `{ "error" : "big error"}`, resp.Body.String())
}
