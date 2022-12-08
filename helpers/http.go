package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// HttpCode wraps int used for http codes in http package
// adds more readability to operate on HttpCode instead of plain int
type HttpCode int

// Validator is a constraint interface for any input value to be validated before calling actual endpoint
type Validator interface {
	Validate() error
}

type errorResponse struct {
	Error string `json:"error"`
}

// JsonHandler takes EndpointFunc and converts it into http.HandlerFunc by taking care of method and media matching,
// input deserialization into json, validation and proper response return.
func JsonHandler[I Validator, R any](method string, mediaType string, endpoint func(input I) (R, HttpCode, error)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		response, httpCode := func() (interface{}, HttpCode) {
			input, httpCode, err := doRequestParsing[I](method, mediaType, request)
			if err != nil {
				return errorResponse{Error: err.Error()}, httpCode
			}
			var response interface{}
			response, httpCode, err = endpoint(input)
			if err != nil {
				response = errorResponse{Error: err.Error()}
			}
			return response, httpCode
		}()

		if err := writeJsonResponse(writer, response, httpCode); err != nil {
			fmt.Println("Error writing response: ", err)
		}
	}
}

// doRequestParsing function takes request and does some basic validation against method, media type, tries to parse
// json into given input type I, also calls Validate method. It either returns validated input or http code and error
func doRequestParsing[I Validator](method, mediaType string, req *http.Request) (I, HttpCode, error) {
	var input I
	if req.Method != method {
		return input, http.StatusMethodNotAllowed, errors.New("unexpected http method")
	}

	if req.Header.Get("Content-Type") != mediaType {
		return input, http.StatusUnsupportedMediaType, fmt.Errorf("unexpected media type")
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return input, http.StatusBadRequest, fmt.Errorf("json parsing: %w", err)
	}
	if err := input.Validate(); err != nil {
		return input, http.StatusUnprocessableEntity, fmt.Errorf("input sanitization: %w", err)
	}
	return input, http.StatusOK, nil
}

// writeJsonResponse serializes passed JsonValue into json, adds json headers to response and writes passed
// http code and serialized json into http.ResponseWriter
func writeJsonResponse(writer http.ResponseWriter, val interface{}, code HttpCode) error {
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
