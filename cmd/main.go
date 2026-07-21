package main

import (
	"fmt"
	"net/http"
	"os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Groupie Tracker is running")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running at http://localhost:" + port)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
