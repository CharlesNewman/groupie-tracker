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
		if strings.TrimSpace(SearchRequest) == "" {
			ErrorHandler(w, "Search cannot be empty", http.StatusBadRequest)
			return
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

func ArtistDetailsHandler(artists []Artist, relations []Relation, locations []Location, dates []Date) http.HandlerFunc {
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
		var locationCount int
		var concertCount int

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

		for i := 0; i < len(locations); i++ {
			location := locations[i]

			if location.ID == idInt {
				locationCount = len(location.Locations)
				break
			}
		}

		for i := 0; i < len(dates); i++ {
			date := dates[i]

			if date.ID == idInt {
				concertCount = len(date.Dates)
				break
			}
		}

		if FoundArtist.ID == 0 || FoundRelation.ID == 0 {
			http.Error(w, "Mismatch Artist ID", http.StatusNotFound)
			return
		}
		result := Find{
			Artist:        FoundArtist,
			Relation:      FoundRelation,
			ConcertCount:  concertCount,
			LocationCount: locationCount,
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

func FilterHandler(artists []Artist, locations []Location) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		creationMinText := r.URL.Query().Get("creationMin")
		creationMaxText := r.URL.Query().Get("creationMax")
		firstAlbumMinText := r.URL.Query().Get("firstAlbumMin")
		firstAlbumMaxText := r.URL.Query().Get("firstAlbumMax")
		numberOfMembersMaxText := r.URL.Query().Get("membersMax")
		numberOfMembersMinText := r.URL.Query().Get("membersMin")
		location := r.URL.Query().Get("location")

		if creationMinText == "" {
			creationMinText = "1900"
		}

		if creationMaxText == "" {
			creationMaxText = "2026"
		}

		if firstAlbumMinText == "" {
			firstAlbumMinText = "1900"
		}

		if firstAlbumMaxText == "" {
			firstAlbumMaxText = "2026"
		}

		if numberOfMembersMinText == "" {
			numberOfMembersMinText = "1"
		}

		if numberOfMembersMaxText == "" {
			numberOfMembersMaxText = "20"
		}

		minYear, err := strconv.Atoi(creationMinText)
		if err != nil {
			ErrorHandler(w, "Invalid minimum year", http.StatusBadRequest)
			return
		}

		maxYear, err := strconv.Atoi(creationMaxText)
		if err != nil {
			ErrorHandler(w, "Invalid maximum year", http.StatusBadRequest)
			return
		}

		albumMinYear, err := strconv.Atoi(firstAlbumMinText)
		if err != nil {
			ErrorHandler(w, "Invalid minimum album year", http.StatusBadRequest)
			return
		}

		albumMaxYear, err := strconv.Atoi(firstAlbumMaxText)
		if err != nil {
			ErrorHandler(w, "Invalid maximum year", http.StatusBadRequest)
			return
		}

		membersMaxNum, err := strconv.Atoi(numberOfMembersMaxText)
		if err != nil {
			ErrorHandler(w, "Invalid number of Members", http.StatusBadRequest)
			return
		}

		membersMinNum, err := strconv.Atoi(numberOfMembersMinText)
		if err != nil {
			ErrorHandler(w, "Invalid number of Members", http.StatusBadRequest)
			return
		}

		var filteredArtists []Artist

		for _, artist := range artists {

			matchesLocation := location == ""

			for _, locationData := range locations {
				if artist.ID != locationData.ID {
					continue
				}

				for _, artistLocation := range locationData.Locations {
					if strings.EqualFold(strings.TrimSpace(artistLocation), strings.TrimSpace(location)) {
						matchesLocation = true
						break
					}
				}
				break
			}

			firstAlbumDateText := artist.FirstAlbum
			result := ""
			for i := 6; i < len(firstAlbumDateText); i++ {
				result += string(firstAlbumDateText[i])
			}

			intResult, err := strconv.Atoi(result)
			if err != nil {
				ErrorHandler(w, "Invalid first album year", http.StatusBadRequest)
				return
			}

			if artist.CreationDate >= minYear &&
				artist.CreationDate <= maxYear &&
				intResult >= albumMinYear &&
				intResult <= albumMaxYear &&
				len(artist.Members) >= membersMinNum &&
				len(artist.Members) <= membersMaxNum &&
				matchesLocation {
				filteredArtists = append(filteredArtists, artist)
			}
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			ErrorHandler(w, "Could not display request", http.StatusInternalServerError)
			return
		}
		pageData := PageData{
			Artists:     filteredArtists,
			CurrentPage: 1,
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			ErrorHandler(w, "Could not display page", http.StatusInternalServerError)
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
