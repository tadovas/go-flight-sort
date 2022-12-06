package flight_sorter

import (
	"errors"
)

var (
	ErrNoFlights  = errors.New("no flights provided")
	ErrFlightLoop = errors.New("flight loop detected")
	ErrBrokenPath = errors.New("broken flight path")
)

// Airport is a wrapper around string to make code more readable
type Airport string

// Flight represent a single flight for user from Source Airport to Dest Airport
type Flight struct {
	Source Airport `json:"source"`
	Dest   Airport `json:"dest"`
}

// airportSet represents unique set of either start airports or destination airports
type airportSet map[Airport]struct{}

// SortFlights takes a list of direct users flights in any order
// and returns a single flight with initial source and final destination if possible.
// Otherwise it might return:
// ErrNoFlights - if input contains no flights
// ErrFlightLoop - if input contains flight loops (i.e. input or part of input is closed flight loop)
// ErrBrokenPath - if input contains flights which diverge (same source for different destinations) or
// or merge - flights with different sources has same destination. Also if there are more than one initial source
// or final destination
func SortFlights(flights []Flight) (Flight, error) {
	if len(flights) == 0 {
		return Flight{}, ErrNoFlights
	} else if len(flights) == 1 {
		return flights[0], nil
	}

	sourceAirports := make(airportSet)
	destinationAirports := make(airportSet)
	for _, flight := range flights {
		sourceAirports[flight.Source] = struct{}{}
		destinationAirports[flight.Dest] = struct{}{}
	}

	start, err := findAirportWithNoReference(sourceAirports, destinationAirports)
	if err != nil {
		return Flight{}, err
	}

	destination, err := findAirportWithNoReference(destinationAirports, sourceAirports)
	if err != nil {
		return Flight{}, err
	}

	return Flight{Source: start, Dest: destination}, nil
}

// findAirportWithNoReference takes set of airports and returns only the airport which is not in other set
// i.e. if first set is sources and other is destinations, that means if input is valid an exactly one airport will
// be in source set but not in destination set. The same applies vice versa - given destinations there should be exactly
// one airport which is in destinations but not in source airport set. Any other results either end up in flight loops or
// in broken path.
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
