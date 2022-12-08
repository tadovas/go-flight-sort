package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/tadovas/go-flight-sort/endpoints"
	"github.com/tadovas/go-flight-sort/helpers"
)

func run() error {
	http.HandleFunc("/calculate", helpers.JsonHandler(http.MethodPost, "application/json", endpoints.CalculateEndpoint))

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
