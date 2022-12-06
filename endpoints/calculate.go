package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	flightsorter "github.com/tadovas/go-flight-sort"
	"github.com/tadovas/go-flight-sort/helpers"
)

// FlightsInput is top level struct used to deserialize calculate endpoints request body into list of flights
type FlightsInput struct {
	Flights []flightsorter.Flight `json:"flights"`
}

// Validate method does some sanity checking on FlightsInput and stops on first error detected
func (fi FlightsInput) Validate() error {
	if len(fi.Flights) < 1 {
		return errors.New("at least one flight expected")
	}
	for i, flight := range fi.Flights {
		switch {
		case flight.Source == "":
			return fmt.Errorf("flights[%v] source is empty", i)
		case flight.Dest == "":
			return fmt.Errorf("flights[%v] destination is empty", i)
		}
	}
	return nil
}

// FlightOutput is top level response structure for calculate endpoint
type FlightOutput struct {
	Flight flightsorter.Flight `json:"flight"`
}

// CalculateEndpoint represents /caluclate exposed by http. It takes request, does some checking, then calls flight sorter
// function and return appropriate results
func CalculateEndpoint(req *http.Request) (helpers.JsonValue, helpers.HttpCode, error) {
	if req.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("only POST method allowed")
	}

	if req.Header.Get("Content-Type") != "application/json" {
		return nil, http.StatusUnsupportedMediaType, errors.New("only json payloads supported")
	}

	var input FlightsInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("json parsing: %w", err)
	}
	if err := input.Validate(); err != nil {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf("input sanitization: %w", err)
	}

	flight, err := flightsorter.SortFlights(input.Flights)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf("invalid flight data: %w", err)
	}

	return &FlightOutput{Flight: flight}, http.StatusOK, nil
}
