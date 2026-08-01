package main

import (
	"fmt"

	"groupie-tracker/internal"
)

func main() {
	err := internal.StartServer()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
