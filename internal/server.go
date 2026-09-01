package internal

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

func StartServer() error {

	var wg sync.WaitGroup

	var artists []Artist
	var relations []Relation
	var locations []Location
	var dates []Date
	var artistsErr error
	var relationsErr error
	var locationsErr error
	var datesErr error

	wg.Add(4)

	// Fetch artists
	go func() {
		defer wg.Done()

		artists, artistsErr = FetchArtists()
	}()

	// Fetch relations

	go func() {
		defer wg.Done()
		relations, relationsErr = FetchRelations()
	}()

	// Fetch locations

	go func() {
		defer wg.Done()

		locations, locationsErr = FetchLocations()

	}()

	// Fetch dates

	go func() {
		defer wg.Done()
		dates, datesErr = FetchDates()
	}()

	wg.Wait()

	// Artists Error
	if artistsErr != nil {
		return fmt.Errorf("could not fetch artists: %w", artistsErr)
	}

	if len(artists) == 0 {
		return fmt.Errorf("no artists received")
	}

	// Relation Error
	if relationsErr != nil {
		return fmt.Errorf("could not fetch relations: %w", relationsErr)
	}

	if len(relations) == 0 {
		return fmt.Errorf("no relations received")
	}

	// Location Error

	if locationsErr != nil {
		return fmt.Errorf("could not fetch locations: %w", locationsErr)
	}

	if len(locations) == 0 {
		return fmt.Errorf("no locations received")
	}

	// Date Error
	if datesErr != nil {
		return fmt.Errorf("could not fetch dates: %w", datesErr)
	}

	if len(dates) == 0 {
		return fmt.Errorf("no dates received")
	}

	mux := SetupRoutes(artists, relations, locations, dates)

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running at http://localhost:" + port)

	return http.ListenAndServe(":"+port, mux)
}

func SetupRoutes(artists []Artist, relations []Relation, locations []Location, dates []Date) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Handler(artists, locations))
	mux.HandleFunc("/artist", ArtistDetailsHandler(artists, relations, locations, dates))
	mux.HandleFunc("/search", SearchHandler(artists))
	mux.HandleFunc("/suggestions", SuggestionHandler(artists))
	mux.HandleFunc("/artist-details", DetailsHandler(artists, relations))
	mux.HandleFunc("/filter", FilterHandler(artists, locations))

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", fileServer),
	)

	return mux
}
