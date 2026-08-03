package main

import (
	"fmt"

	"groupie-tracker/internal"
)

var startServer = internal.StartServer

func run() error {
	return startServer()
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
