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

// For the pages
type PageData struct {
	Artists      []internal.Artist
	CurrentPage  int
	PreviousPage int
	NextPage     int
	HasPrevious  bool
	HasNext      bool
}

// ////
func Handler(artists []internal.Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}

		//For the pages
		readPage := r.URL.Query().Get("page")
		currentPage := 1
		if readPage != "" {
			pageNumber, err := strconv.Atoi(readPage)
			if err != nil {
				http.Error(w, "Invalid page number", http.StatusBadRequest)
				return
			}
			if pageNumber < 1 {
				http.Error(w, "Invalid page number", http.StatusBadRequest)
				return
			} else {
				currentPage = pageNumber
			}
		}
		pageSize := 12
		start := (currentPage - 1) * pageSize
		end := start + pageSize

		if end > len(artists) {
			end = len(artists)
		}
		PreviousPage := currentPage - 1
		NextPage := currentPage + 1
		HasPrevious := currentPage > 1
		HasNext := end < len(artists)

		pageData := PageData{
			Artists:      artists[start:end],
			CurrentPage:  currentPage,
			PreviousPage: PreviousPage,
			NextPage:     NextPage,
			HasPrevious:  HasPrevious,
			HasNext:      HasNext,
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			http.Error(w, "Could not display page", http.StatusInternalServerError)
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
