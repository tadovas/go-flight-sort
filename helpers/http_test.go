package helpers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type inputTest struct {
}

func (inputTest) Validate() error {
	return nil
}

type outputTest struct {
}

func TestJsonHandlerWritesCorrectJsonOutput(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/<irrelevant>", strings.NewReader("{}"))
	resp := httptest.NewRecorder()

	httpHandler := JsonHandler(http.MethodPost, "", func(ts inputTest) (map[string]interface{}, HttpCode, error) {
		return map[string]interface{}{
				"some_field": "abc",
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
	req := httptest.NewRequest(http.MethodPost, "/<irrelevant>", strings.NewReader("{}"))
	resp := httptest.NewRecorder()

	httpHandler := JsonHandler(http.MethodPost, "", func(inputTest) (outputTest, HttpCode, error) {
		return outputTest{}, http.StatusBadRequest, errors.New("big error")
	})

	httpHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-type"))
	assert.JSONEq(t, `{ "error" : "big error"}`, resp.Body.String())
}
