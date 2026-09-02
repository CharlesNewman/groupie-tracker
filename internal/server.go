package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var SyncInterval time.Duration = 10 // Interval in minutes

var (
	artists   []Artist
	relations []Relation
	locations []Location
	dates     []Date
)

func StartServer() error {
	// Initial FetchAPI
	err := FetchAPI()

	if err != nil {
		return fmt.Errorf("initial sync failed: %w", err)
	}

	// Background FetchAPI with ticker loop every 10 minutes
	go func() {
		for range time.Tick(SyncInterval * time.Minute) {

			err := FetchAPI()

			if err != nil {
				log.Printf("Background refresh error: %v", err)
				continue
			}
			log.Println("Data refreshed successfully in background!")
		}
	}()

	// Webapp Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := SetupRoutes()
	// Message on server start
	fmt.Println("Server running at http://localhost:" + port)

	server := &http.Server{
		Addr:              ":" + port,        // TCP address to listen on
		Handler:           mux,               // Handler to invoke (our router)
		ReadTimeout:       5 * time.Second,   // Max duration for reading the entire request
		ReadHeaderTimeout: 2 * time.Second,   // Max duration for reading request headers (Slowloris protection)
		WriteTimeout:      10 * time.Second,  // Max duration before timing out writes of the response
		IdleTimeout:       120 * time.Second, // Max amount of time to wait for the next request when keep-alive is enabled
	}
	// Serve webapp
	return server.ListenAndServe()
}

func FetchAPI() error {
	var artistsErr error
	var relationsErr error
	var locationsErr error
	var datesErr error

	// Channels
	artistChannel := make(chan ArtistResult)
	relationChannel := make(chan RelationResult)
	locationChannel := make(chan LocationResult)
	dateChannel := make(chan DateResult)

	// Fetch artists
	go func() {

		artists, err := FetchArtists()

		artistChannel <- ArtistResult{
			Artists: artists,
			Err:     err,
		}
	}()

	// Fetch relations
	go func() {

		relations, err := FetchRelations()

		relationChannel <- RelationResult{
			Relations: relations,
			Err:       err,
		}
	}()

	// Fetch locations
	go func() {

		locations, err := FetchLocations()

		locationChannel <- LocationResult{
			Locations: locations,
			Err:       err,
		}
	}()

	// Fetch dates
	go func() {

		dates, err := FetchDates()

		dateChannel <- DateResult{
			Dates: dates,
			Err:   err,
		}
	}()

	for i := 0; i < 4; i++ {
		select {
		case artistResult := <-artistChannel:
			artists = artistResult.Artists
			artistsErr = artistResult.Err

		case RelationResult := <-relationChannel:
			relations = RelationResult.Relations
			relationsErr = RelationResult.Err

		case LocationResult := <-locationChannel:
			locations = LocationResult.Locations
			locationsErr = LocationResult.Err

		case DateResult := <-dateChannel:
			dates = DateResult.Dates
			datesErr = DateResult.Err

		}
	}

	if artistsErr != nil {
		return fmt.Errorf("could not fetch artists: %w", artistsErr)
	}
	if len(artists) == 0 {
		return fmt.Errorf("no artists received")
	}
	if relationsErr != nil {
		return fmt.Errorf("could not fetch relations: %w", relationsErr)
	}
	if len(relations) == 0 {
		return fmt.Errorf("no relations received")
	}
	if locationsErr != nil {
		return fmt.Errorf("could not fetch locations: %w", locationsErr)
	}
	if len(locations) == 0 {
		return fmt.Errorf("no locations received")
	}
	if datesErr != nil {
		return fmt.Errorf("could not fetch dates: %w", datesErr)
	}
	if len(dates) == 0 {
		return fmt.Errorf("no dates received")
	}
	return nil
}

func SetupRoutes() *http.ServeMux {
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
