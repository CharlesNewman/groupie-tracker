package main

import (
	"fmt"
	"net/http"
	"os"

	"groupie-tracker/internal"
)

func main() {
	// Artists
	artists, err := internal.FetchArtists()
	if err != nil {
		fmt.Println("API error:", err)
		return
	}

	if len(artists) == 0 {
		fmt.Println("No artists received")
		return
	}
	fmt.Println("Artists:", len(artists))
	// Relations
	relations, err := internal.FetchRelations()
	if err != nil {
		fmt.Println("API error:", err)
		return
	}
	if len(relations) == 0 {
		fmt.Println("No location and dates received")
		return
	}
	// Locations
	locations, err := internal.FetchLocations()
	if err != nil {
		fmt.Println("API error:", err)
		return
	}
	if len(locations) == 0 {
		fmt.Println("No locations received")
		return
	}
	// Dates
	dates, err := internal.FetchDates()
	if err != nil {
		fmt.Println("API error:", err)
		return
	}
	if len(dates) == 0 {
		fmt.Println("No dates received")
		return
	}

	// this is how i send data to the html
	mux := http.NewServeMux()
	mux.HandleFunc("/", internal.Handler(artists))
	mux.HandleFunc("/artist", internal.ArtistDetailsHandler(artists, relations))

	// the StripPrefix is used to prevent the html to repeat a path
	// and ending up somewere it doesnt exists
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

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
