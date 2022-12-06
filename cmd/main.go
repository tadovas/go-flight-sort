package main

import (
	"fmt"
	"os"
)

func run() error {
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Terminated")
}
