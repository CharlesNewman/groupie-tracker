package internal

import "testing"

func TestFetchArtists(t *testing.T) {
	artists, err := FetchArtists()
	if err != nil {
		t.Fatalf("FetchArtists returned an error: %v", err)
	}

	if len(artists) == 0 {
		t.Fatal("expected artists, got an empty list")
	}

	if artists[0].ID == 0 {
		t.Error("expected the first artist to have an ID")
	}

	if artists[0].Name == "" {
		t.Error("expected the first artist to have a name")
	}
}

func TestFetchRelations(t *testing.T) {
	relations, err := FetchRelations()
	if err != nil {
		t.Fatalf("FetchRelations returned an error: %v", err)
	}

	if len(relations) == 0 {
		t.Fatal("expected relations, got an empty list")
	}

	if relations[0].ID == 0 {
		t.Error("expected the first relation to have an ID")
	}

	if relations[0].DatesLocations == nil {
		t.Error("expected the first relation to contain dates and locations")
	}
}

func TestFetchLocations(t *testing.T) {
	locations, err := FetchLocations()
	if err != nil {
		t.Fatalf("FetchLocations returned an error: %v", err)
	}

	if len(locations) == 0 {
		t.Fatal("expected locations, got an empty list")
	}

	if locations[0].ID == 0 {
		t.Error("expected the first location to have an ID")
	}

	if len(locations[0].Locations) == 0 {
		t.Error("expected the first location to contain locations")
	}
}

func TestFetchDates(t *testing.T) {
	dates, err := FetchDates()
	if err != nil {
		t.Fatalf("FetchDates returned an error: %v", err)
	}

	if len(dates) == 0 {
		t.Fatal("expected dates, got an empty list")
	}

	if dates[0].ID == 0 {
		t.Error("expected the first date entry to have an ID")
	}

	if len(dates[0].Dates) == 0 {
		t.Error("expected the first date entry to contain dates")
	}
}

func TestAPIDataIDsMatch(t *testing.T) {
	artists, err := FetchArtists()
	if err != nil {
		t.Fatalf("could not fetch artists: %v", err)
	}

	relations, err := FetchRelations()
	if err != nil {
		t.Fatalf("could not fetch relations: %v", err)
	}

	locations, err := FetchLocations()
	if err != nil {
		t.Fatalf("could not fetch locations: %v", err)
	}

	dates, err := FetchDates()
	if err != nil {
		t.Fatalf("could not fetch dates: %v", err)
	}

	if len(artists) != len(relations) {
		t.Errorf(
			"artists and relations have different lengths: %d and %d",
			len(artists),
			len(relations),
		)
	}

	if len(artists) != len(locations) {
		t.Errorf(
			"artists and locations have different lengths: %d and %d",
			len(artists),
			len(locations),
		)
	}

	if len(artists) != len(dates) {
		t.Errorf(
			"artists and dates have different lengths: %d and %d",
			len(artists),
			len(dates),
		)
	}

	for i := 0; i < len(artists); i++ {
		if artists[i].ID != relations[i].ID {
			t.Errorf(
				"artist ID %d does not match relation ID %d",
				artists[i].ID,
				relations[i].ID,
			)
		}

		if artists[i].ID != locations[i].ID {
			t.Errorf(
				"artist ID %d does not match location ID %d",
				artists[i].ID,
				locations[i].ID,
			)
		}

		if artists[i].ID != dates[i].ID {
			t.Errorf(
				"artist ID %d does not match date ID %d",
				artists[i].ID,
				dates[i].ID,
			)
		}
	}
}
