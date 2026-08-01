package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

func Handler(artists []Artist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
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

func ArtistDetailsHandler(artists []Artist, relations []Relation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
