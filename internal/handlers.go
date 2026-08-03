package internal

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Handler(artists []Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			ErrorHandler(w, "Invalid path", http.StatusNotFound)
			return
		}

		if r.Method != http.MethodGet {
			ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			ErrorHandler(w, "Template Error", http.StatusInternalServerError)
			return
		}
		readPage := r.URL.Query().Get("page")
		currentPage := 1
		pageSize := 12
		maxPage := (len(artists) + pageSize - 1) / pageSize
		if readPage != "" {
			pageNumber, err := strconv.Atoi(readPage)
			if err != nil {
				ErrorHandler(w, "Invalid page number", http.StatusBadRequest)
				return
			}
			if pageNumber < 1 {
				ErrorHandler(w, "Invalid page number", http.StatusBadRequest)
				return
			} else if pageNumber > maxPage {
				ErrorHandler(w, "Invalid page number", http.StatusNotFound)
				return
			} else {
				currentPage = pageNumber
			}
		}

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
			ErrorHandler(w, "Could not display page", http.StatusInternalServerError)
			return
		}
	}
}

func SearchHandler(artists []Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		SearchRequest := r.URL.Query().Get("query")
		lowerSearchRequest := strings.ToLower(SearchRequest)
		var matches []Artist
		for i := 0; i < len(artists); i++ {
			artist := artists[i]
			lowerArtist := strings.ToLower(artist.Name)
			if strings.Contains(lowerArtist, lowerSearchRequest) {
				matches = append(matches, artist)
			}
		}
		pageData := PageData{
			Artists:     matches,
			CurrentPage: 1,
		}
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			ErrorHandler(w, "Could not display request", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, pageData)
		if err != nil {
			ErrorHandler(w, "Could not display page", http.StatusInternalServerError)
			return
		}
	}
}

func SuggestionHandler(artists []Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			//Use ErrorHandler for full pages if you do here it will break response.json() in javascript
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		SuggestionRequest := r.URL.Query().Get("query")
		lowerSuggestionRequest := strings.ToLower(SuggestionRequest)
		var matches []Artist
		for i := 0; i < len(artists); i++ {
			artist := artists[i]
			lowerArtist := strings.ToLower(artist.Name)
			if strings.Contains(lowerArtist, lowerSuggestionRequest) {
				matches = append(matches, artist)
			}
		}
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(matches)
		if err != nil {
			//Use ErrorHandler for full pages if you do here it will break response.json() in javascript
			http.Error(w, "Could not encode suggestions", http.StatusInternalServerError)
			return
		}
	}
}

func ArtistDetailsHandler(artists []Artist, relations []Relation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//The errors work with javascript here as well thats why we use http.Error instead the ErrorHandler

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idText := r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(idText)
		if err != nil {
			http.Error(w, "Could not find artist ID", http.StatusBadRequest)
			return
		}
		var FoundArtist Artist
		var FoundRelation Relation

		type Find struct {
			Artist   Artist   `json:"artist"`
			Relation Relation `json:"relation"`
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

func DetailsHandler(artists []Artist, relations []Relation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idText := r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(idText)
		if err != nil {
			ErrorHandler(w, "Could not find artist ID", http.StatusBadRequest)
			return
		}
		type Find struct {
			Artist   Artist   `json:"artist"`
			Relation Relation `json:"relation"`
		}

		var FoundArtist Artist
		var FoundRelation Relation

		for i := 0; i < len(artists); i++ {
			artist := artists[i]
			if artist.ID == idInt {
				FoundArtist = artist
				break
			}
		}
		for j := 0; j < len(relations); j++ {
			relation := relations[j]
			if relation.ID == idInt {
				FoundRelation = relation
				break
			}
		}
		if FoundArtist.ID == 0 || FoundRelation.ID == 0 {
			ErrorHandler(w, "Mismatch Artist ID", http.StatusNotFound)
			return
		}
		result := Find{
			Artist:   FoundArtist,
			Relation: FoundRelation,
		}

		funcMap := template.FuncMap{
			"formatLocation": func(location string) string {
				location = strings.ReplaceAll(location, "_", " ")
				location = strings.ReplaceAll(location, "-", " ")
				location = strings.Title(location)
				location = strings.ReplaceAll(location, "Usa", "USA")
				location = strings.ReplaceAll(location, "Uk", "UK")
				return location
			},
		}

		tmpl, err := template.New("artist-details.html").
			Funcs(funcMap).
			ParseFiles("templates/artist-details.html")

		if err != nil {
			ErrorHandler(w, "Could not load details page", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, result)
		if err != nil {
			ErrorHandler(w, "Could not display details page", http.StatusInternalServerError)
			return
		}
	}
}

func ErrorHandler(w http.ResponseWriter, message string, statusCode int) {

	type ErrorData struct {
		Code    int
		Message string
	}

	data := ErrorData{
		Code:    statusCode,
		Message: message,
	}
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not display error page", http.StatusInternalServerError)
		return
	}
}
