package main

import (
	"fmt"
	"groupie-tracker/internal"
	"net/http"
	"os"
)

func homeHandler(artists []internal.Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, artist := range artists {
			fmt.Fprintln(w, artist.Name)
		}
	}
}

func main() {
	artists, err := internal.FetchArtists()
	if err != nil {
		fmt.Println("API error:", err)
		return
	}

	if len(artists) == 0 {
		fmt.Println("No artists received")
		return
	}

	fmt.Println("Number of artists:", len(artists))
	fmt.Println("First artist:", artists[0].Name)

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler(artists))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running at http://localhost:" + port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
