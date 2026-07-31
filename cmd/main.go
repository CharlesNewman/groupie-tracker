package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"groupie-tracker/internal"
)

func Handler(artists []internal.Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, artists)
		if err != nil {
			http.Error(w, "Could not displace page", http.StatusInternalServerError)
		}
	}
}

func ArtistDetailsHandler(artists []internal.Artist, relations []internal.Relation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idText := r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(idText)
		if err != nil {
			http.Error(w, "Could not find artist ID", http.StatusBadRequest)
			return
		}
		var FoundArtist internal.Artist
		var FoundRelation internal.Relation

		type Find struct {
			Artist   internal.Artist   `json:"artist"`
			Relation internal.Relation `json:"relation"`
		}

		for i := 0; i < len(artists); i++ {
			artist := artists[i]
			if artist.ID == idInt {
				FoundArtist = artist
				break
			}
		}
		for i := 0; i < len(relations); i++ {
			relation := relations[i]
			if relation.ID == idInt {
				FoundRelation = relation
				break
			}
		}
		if FoundArtist.ID == 0 || FoundRelation.ID == 0 {
			http.Error(w, "Mismatch Artist ID", http.StatusNotFound)
			return
		}
		result := Find{
			Artist:   FoundArtist,
			Relation: FoundRelation,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, "Could not encode result", http.StatusInternalServerError)
			return
		}
	}
}

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
	fmt.Println("Dates:", len(dates))

	// this is how i send data to the html
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler(artists))
	mux.HandleFunc("/artist", ArtistDetailsHandler(artists, relations))

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
