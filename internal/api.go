package internal

import (
	"encoding/json"
	"net/http"
)

func FetchArtists() ([]Artist, error) {
	apidata, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer apidata.Body.Close()

	var artists []Artist

	err = json.NewDecoder(apidata.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
