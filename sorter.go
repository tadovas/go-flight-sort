package flight_sorter

import (
	"errors"
)

var (
	ErrNoFlights  = errors.New("no flights provided")
	ErrFlightLoop = errors.New("flight loop detected")
	ErrBrokenPath = errors.New("broken flight path")
)

type Flight struct {
	Start Airport `json:"start"`
	Dest  Airport `json:"dest"`
}

type Airport string

type airportSet map[Airport]struct{}

func SortFlights(flights []Flight) (Flight, error) {
	if len(flights) == 0 {
		return Flight{}, ErrNoFlights
	} else if len(flights) == 1 {
		return flights[0], nil
	}

	startAirports := make(airportSet)
	destinationAirports := make(airportSet)
	for _, flight := range flights {
		startAirports[flight.Start] = struct{}{}
		destinationAirports[flight.Dest] = struct{}{}
	}

	start, err := findAirportWithNoReference(startAirports, destinationAirports)
	if err != nil {
		return Flight{}, err
	}

	destination, err := findAirportWithNoReference(destinationAirports, startAirports)
	if err != nil {
		return Flight{}, err
	}

	return Flight{Start: start, Dest: destination}, nil
}

func findAirportWithNoReference(set airportSet, other airportSet) (Airport, error) {
	var res []Airport
	for key, _ := range set {
		if _, found := other[key]; !found {
			res = append(res, key)
		}
	}
	switch {
	case len(res) < 1:
		return "", ErrFlightLoop
	case len(res) > 1:
		return "", ErrBrokenPath
	}
	return res[0], nil
}
