package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	flightsorter "github.com/tadovas/go-flight-sort"
	"github.com/tadovas/go-flight-sort/helpers"
)

func run() error {
	http.HandleFunc("/calculate", helpers.HttpHandler(func(req *http.Request) (helpers.JsonValue, helpers.HttpCode, error) {
		if req.Method != http.MethodPost {
			return nil, http.StatusMethodNotAllowed, errors.New("only POST method allowed")
		}

		if req.Header.Get("Content-Type") != "application/json" {
			return nil, http.StatusUnsupportedMediaType, errors.New("only json payloads supported")
		}

		var input struct {
			Flights []flightsorter.Flight `json: "flights"`
		}
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("json parsing: %w", err)
		}

		flight, err := flightsorter.SortFlights(input.Flights)
		if err != nil {
			return nil, http.StatusUnprocessableEntity, fmt.Errorf("invalid input: %w", err)
		}
		var output = struct {
			Flight flightsorter.Flight `json: "flight"`
		}{
			Flight: flight,
		}

		return &output, http.StatusOK, nil
	}))

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("listener: %w", err)
	}
	fmt.Println("Will serve you now at :8080")
	return http.Serve(l, nil)
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Terminated")
}
