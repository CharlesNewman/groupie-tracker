package internal

import "testing"

func TestArtistFields(t *testing.T) {
	artist := Artist{
		ID:           1,
		Name:         "Test Artist",
		Members:      []string{"Alice", "Bob"},
		CreationDate: 2000,
		FirstAlbum:   "01-01-2005",
	}

	if artist.ID != 1 {
		t.Errorf("expected ID 1, got %d", artist.ID)
	}

	if artist.Name != "Test Artist" {
		t.Errorf("expected Test Artist, got %s", artist.Name)
	}

	if len(artist.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(artist.Members))
	}

	if artist.CreationDate != 2000 {
		t.Errorf(
			"expected creation year 2000, got %d",
			artist.CreationDate,
		)
	}

	if artist.FirstAlbum != "01-01-2005" {
		t.Errorf(
			"expected first album 01-01-2005, got %s",
			artist.FirstAlbum,
		)
	}
}

func TestRelationFields(t *testing.T) {
	relation := Relation{
		ID: 1,
		DatesLocations: map[string][]string{
			"athens-greece": {"01-01-2026"},
		},
	}

	if relation.ID != 1 {
		t.Errorf("expected ID 1, got %d", relation.ID)
	}

	dates := relation.DatesLocations["athens-greece"]

	if len(dates) != 1 {
		t.Errorf("expected 1 date, got %d", len(dates))
	}
}

func TestLocationFields(t *testing.T) {
	location := Location{
		ID: 1,
		Locations: []string{
			"athens-greece",
			"london-uk",
		},
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
}

func TestDateFields(t *testing.T) {
	date := Date{
		ID: 1,
		Dates: []string{
			"01-01-2026",
			"02-02-2026",
		},
	}

	if date.ID != 1 {
		t.Errorf("expected ID 1, got %d", date.ID)
	}

	if len(date.Dates) != 2 {
		t.Errorf("expected 2 dates, got %d", len(date.Dates))
	}
}

func TestPageDataFields(t *testing.T) {
	pageData := PageData{
		Artists: []Artist{
			{
				ID:   1,
				Name: "Test Artist",
			},
		},
		CurrentPage:  1,
		PreviousPage: 0,
		NextPage:     2,
		HasPrevious:  false,
		HasNext:      true,
		LocationList: []string{
			"athens-greece",
			"london-uk",
		},
	}

	if len(pageData.Artists) != 1 {
		t.Errorf(
			"expected 1 artist, got %d",
			len(pageData.Artists),
		)
	}

	if pageData.CurrentPage != 1 {
		t.Errorf(
			"expected current page 1, got %d",
			pageData.CurrentPage,
		)
	}

	if pageData.HasPrevious {
		t.Error("expected HasPrevious to be false")
	}

	if !pageData.HasNext {
		t.Error("expected HasNext to be true")
	}

	if len(pageData.LocationList) != 2 {
		t.Errorf(
			"expected 2 locations, got %d",
			len(pageData.LocationList),
		)
	}
}

func TestFindFields(t *testing.T) {
	result := Find{
		Artist: Artist{
			ID:   1,
			Name: "Test Artist",
		},
		Relation: Relation{
			ID: 1,
		},
		LocationCount: 2,
		ConcertCount:  3,
	}

	if result.Artist.ID != result.Relation.ID {
		t.Error("artist ID and relation ID should match")
	}

	if result.LocationCount != 2 {
		t.Errorf(
			"expected 2 locations, got %d",
			result.LocationCount,
		)
	}

	if result.ConcertCount != 3 {
		t.Errorf(
			"expected 3 concerts, got %d",
			result.ConcertCount,
		)
	}
}
