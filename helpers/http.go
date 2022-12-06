package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HttpCode wraps int used for http codes in http package
// adds more readability to operate on HttpCode instead of plain int
type HttpCode int

// JsonValue represents pointer to real struct we will serialize when endpoint function does its job
// same as HttpCode it serves for code readability
type JsonValue interface{}

// EndpointFunc takes request as a parameter, does specific endpont logic and is expected to return
// either pointer to struct for json serialization or error, http code should be always present.
type EndpointFunc func(req *http.Request) (JsonValue, HttpCode, error)

type errorResponse struct {
	Error string `json:"error"`
}

// JsonHandler takes EndpointFunc and converts it into http.HandlerFunc by taking care of result serialization
// and proper reponse writes
func JsonHandler(endpoint EndpointFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		responseValue, httpStatus, err := endpoint(request)
		if err != nil {
			responseValue = errorResponse{Error: err.Error()}
		}
		if err := writeJsonResponse(writer, responseValue, httpStatus); err != nil {
			fmt.Println("Error writing response: ", err)
		}
	}
}

// writeJsonResponse serializes passed JsonValue into json, adds json headers to response and writes passed
// http code and serialized json into http.ResponseWriter
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
