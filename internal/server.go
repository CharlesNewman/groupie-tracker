package internal

import (
	"fmt"
	"net/http"
	"os"
)

func StartServer() error {
	// Fetch artists
	artists, err := FetchArtists()
	if err != nil {
		return fmt.Errorf("could not fetch artists: %w", err)
	}

	if len(artists) == 0 {
		return fmt.Errorf("no artists received")
	}

	// Fetch relations
	relations, err := FetchRelations()
	if err != nil {
		return fmt.Errorf("could not fetch relations: %w", err)
	}

	if len(relations) == 0 {
		return fmt.Errorf("no relations received")
	}

	// Fetch locations
	locations, err := FetchLocations()
	if err != nil {
		return fmt.Errorf("could not fetch locations: %w", err)
	}

	if len(locations) == 0 {
		return fmt.Errorf("no locations received")
	}

	// Fetch dates
	dates, err := FetchDates()
	if err != nil {
		return fmt.Errorf("could not fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return fmt.Errorf("no dates received")
	}

	mux := SetupRoutes(artists, relations)

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running at http://localhost:" + port)

	return http.ListenAndServe(":"+port, mux)
}

func SetupRoutes(artists []Artist, relations []Relation) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Handler(artists))
	mux.HandleFunc("/artist", ArtistDetailsHandler(artists, relations))
	mux.HandleFunc("/search", SearchHandler(artists))
	mux.HandleFunc("/suggestions", SuggestionHandler(artists))
	mux.HandleFunc("/artist-details", DetailsHandler(artists, relations))

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", fileServer),
	)

	return mux
}
