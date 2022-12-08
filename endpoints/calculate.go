package endpoints

import (
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

// CalculateEndpoint represents /calculate exposed by http. It takes already parsed input struct, calls flight sorting function
// and returns FlightOutput with single flight or appropriate error
func CalculateEndpoint(input FlightsInput) (FlightOutput, helpers.HttpCode, error) {
	flight, err := flightsorter.SortFlights(input.Flights)
	if err != nil {
		return FlightOutput{}, http.StatusUnprocessableEntity, fmt.Errorf("invalid flight data: %w", err)
	}

	return FlightOutput{Flight: flight}, http.StatusOK, nil
}
