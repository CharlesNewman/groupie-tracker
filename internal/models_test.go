package internal

import (
	"reflect"
	"testing"
)

func TestArtistModel(t *testing.T) {
	artist := Artist{
		ID:           1,
		Image:        "image.jpg",
		Name:         "Queen",
		Members:      []string{"Freddie Mercury", "Brian May"},
		CreationDate: 1970,
		FirstAlbum:   "14-12-1973",
		Locations:    "/api/locations/1",
		ConcertDates: "/api/dates/1",
		Relations:    "/api/relation/1",
	}

	if artist.ID != 1 {
		t.Errorf("expected ID 1, got %d", artist.ID)
	}

	if artist.Name != "Queen" {
		t.Errorf("expected Queen, got %s", artist.Name)
	}

	if len(artist.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(artist.Members))
	}

	if artist.CreationDate != 1970 {
		t.Errorf("expected creation date 1970, got %d", artist.CreationDate)
	}
}

func TestRelationModel(t *testing.T) {
	expected := map[string][]string{
		"london-uk": {"01-01-2025", "02-01-2025"},
	}

	relation := Relation{
		ID:             1,
		DatesLocations: expected,
	}

	if relation.ID != 1 {
		t.Errorf("expected ID 1, got %d", relation.ID)
	}

	if !reflect.DeepEqual(relation.DatesLocations, expected) {
		t.Errorf(
			"expected %v, got %v",
			expected,
			relation.DatesLocations,
		)
	}
}

func TestRelationsResponseModel(t *testing.T) {
	response := RelationsResponse{
		Index: []Relation{
			{
				ID: 1,
				DatesLocations: map[string][]string{
					"athens-greece": {"10-08-2026"},
				},
			},
		},
	}

	if len(response.Index) != 1 {
		t.Errorf("expected 1 relation, got %d", len(response.Index))
	}

	if response.Index[0].ID != 1 {
		t.Errorf("expected relation ID 1, got %d", response.Index[0].ID)
	}
}

func TestLocationModel(t *testing.T) {
	location := Location{
		ID:        1,
		Locations: []string{"athens-greece", "london-uk"},
		Dates:     "/api/dates/1",
	}

	if location.ID != 1 {
		t.Errorf("expected ID 1, got %d", location.ID)
	}

	if len(location.Locations) != 2 {
		t.Errorf(
			"expected 2 locations, got %d",
			len(location.Locations),
		)
	}

	if location.Locations[0] != "athens-greece" {
		t.Errorf(
			"expected athens-greece, got %s",
			location.Locations[0],
		)
	}
}

func TestLocationResponseModel(t *testing.T) {
	response := LocationResponse{
		Index: []Location{
			{
				ID:        1,
				Locations: []string{"athens-greece"},
			},
		},
	}

	if len(response.Index) != 1 {
		t.Errorf("expected 1 location, got %d", len(response.Index))
	}

	if response.Index[0].ID != 1 {
		t.Errorf("expected location ID 1, got %d", response.Index[0].ID)
	}
}

func TestDateModel(t *testing.T) {
	date := Date{
		ID:    1,
		Dates: []string{"01-01-2025", "02-01-2025"},
	}

	if date.ID != 1 {
		t.Errorf("expected ID 1, got %d", date.ID)
	}

	if len(date.Dates) != 2 {
		t.Errorf("expected 2 dates, got %d", len(date.Dates))
	}

	if date.Dates[0] != "01-01-2025" {
		t.Errorf("expected 01-01-2025, got %s", date.Dates[0])
	}
}

func TestDateResponseModel(t *testing.T) {
	response := DateResponse{
		Index: []Date{
			{
				ID:    1,
				Dates: []string{"01-01-2025"},
			},
		},
	}

	if len(response.Index) != 1 {
		t.Errorf("expected 1 date, got %d", len(response.Index))
	}

	if response.Index[0].ID != 1 {
		t.Errorf("expected date ID 1, got %d", response.Index[0].ID)
	}
}

func TestPageDataModel(t *testing.T) {
	page := PageData{
		Artists: []Artist{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "SOJA"},
		},
		CurrentPage:  2,
		PreviousPage: 1,
		NextPage:     3,
		HasPrevious:  true,
		HasNext:      true,
	}

	if len(page.Artists) != 2 {
		t.Errorf("expected 2 artists, got %d", len(page.Artists))
	}

	if page.CurrentPage != 2 {
		t.Errorf(
			"expected current page 2, got %d",
			page.CurrentPage,
		)
	}

	if !page.HasPrevious {
		t.Error("expected HasPrevious to be true")
	}

	if !page.HasNext {
		t.Error("expected HasNext to be true")
	}
}

func TestFindModel(t *testing.T) {
	find := Find{
		Artist: Artist{
			ID:   1,
			Name: "Queen",
		},
		Relation: Relation{
			ID: 1,
			DatesLocations: map[string][]string{
				"athens-greece": {"10-08-2026"},
			},
		},
	}

	if find.Artist.ID != find.Relation.ID {
		t.Errorf(
			"expected matching IDs, got artist %d and relation %d",
			find.Artist.ID,
			find.Relation.ID,
		)
	}

	if find.Artist.Name != "Queen" {
		t.Errorf("expected Queen, got %s", find.Artist.Name)
	}
}
