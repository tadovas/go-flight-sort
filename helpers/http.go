package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpCode int

type JsonValue interface{}

type EndpointFunc func(req *http.Request) (JsonValue, HttpCode, error)

func HttpHandler(endpoint EndpointFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		responseValue, httpStatus, err := endpoint(request)
		if err != nil {
			responseValue = struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			}
		}
		if err := writeJsonResponse(writer, responseValue, httpStatus); err != nil {
			fmt.Println("Error writing response: ", err)
		}
	}
}

func writeJsonResponse(writer http.ResponseWriter, val JsonValue, code HttpCode) error {
	jsonBytes, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return fmt.Errorf("json error: %w", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(int(code))
	if _, err := writer.Write(jsonBytes); err != nil {
		return fmt.Errorf("content write error: %w", err)
	}
	return nil
}
