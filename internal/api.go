package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchArtists() ([]Artist, error) {
	apidata, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}

	// Close the API response body when FetchArtists finishes
	// defer means: Run this line later, when the current function ends
	defer apidata.Body.Close()

	var artists []Artist

	if apidata.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", apidata.Status)
	}

	err = json.NewDecoder(apidata.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func FetchRelations() ([]Relation, error) {
	apidata, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}

	defer apidata.Body.Close()

	var relations RelationsResponse

	if apidata.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", apidata.Status)
	}

	err = json.NewDecoder(apidata.Body).Decode(&relations)
	if err != nil {
		return nil, err
	}

	return relations.Index, nil
}

func FetchLocations() ([]Location, error) {
	apidata, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return nil, err
	}
	defer apidata.Body.Close()

	var locations LocationResponse

	if apidata.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", apidata.Status)
	}

	err = json.NewDecoder(apidata.Body).Decode(&locations)
	if err != nil {
		return nil, err
	}

	return locations.Index, nil
}

func FetchDates() ([]Date, error) {
	apidata, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		return nil, err
	}
	defer apidata.Body.Close()

	var dates DateResponse
	if apidata.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", apidata.Status)
	}

	err = json.NewDecoder(apidata.Body).Decode(&dates)
	if err != nil {
		return nil, err
	}
	return dates.Index, nil
}
