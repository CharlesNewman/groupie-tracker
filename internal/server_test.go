package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// setTestData populates the package-level artists/relations/locations/dates
// vars that SetupRoutes reads from, and restores their previous values
// after the test finishes.
func setTestData(t *testing.T) {
	t.Helper()

	originalArtists := artists
	originalRelations := relations
	originalLocations := locations
	originalDates := dates

	artists = []Artist{
		{
			ID:           1,
			Name:         "Test Band",
			Members:      []string{"Alice", "Bob", "Charlie"},
			CreationDate: 2000,
			FirstAlbum:   "01-01-2005",
		},
		{
			ID:           2,
			Name:         "Solo Artist",
			Members:      []string{"David"},
			CreationDate: 2015,
			FirstAlbum:   "01-01-2016",
		},
	}

	relations = []Relation{
		{ID: 1},
		{ID: 2},
	}

	locations = []Location{
		{
			ID:        1,
			Locations: []string{"athens-greece"},
		},
		{
			ID:        2,
			Locations: []string{"london-uk"},
		},
	}

	dates = []Date{
		{ID: 1},
		{ID: 2},
	}

	t.Cleanup(func() {
		artists = originalArtists
		relations = originalRelations
		locations = originalLocations
		dates = originalDates
	})
}

// Tests run from the internal folder.
// This temporarily moves them to the project root
// so templates/ and static/ can be found.
func moveToProjectRoot(t *testing.T) {
	t.Helper()

	currentDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not read current directory: %v", err)
	}

	err = os.Chdir("..")
	if err != nil {
		t.Fatalf("could not move to project root: %v", err)
	}

	t.Cleanup(func() {
		err := os.Chdir(currentDirectory)
		if err != nil {
			t.Errorf("could not restore directory: %v", err)
		}
	})
}

func TestSetupRoutes(t *testing.T) {
	setTestData(t)

	mux := SetupRoutes()

	if mux == nil {
		t.Fatal("SetupRoutes returned nil")
	}
}

func TestHomeRoute(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if !strings.Contains(response.Body.String(), "Test Band") {
		t.Error("home page did not display Test Band")
	}
}

func TestInvalidRoute(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/invalid-route",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.Code,
		)
	}
}

func TestHomeRejectsPost(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			response.Code,
		)
	}
}

func TestSuggestionsRoute(t *testing.T) {
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/suggestions?query=test",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	var result []Artist

	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("could not decode suggestions: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(result))
	}

	if result[0].Name != "Test Band" {
		t.Errorf(
			"expected Test Band, got %s",
			result[0].Name,
		)
	}
}

func TestFilterByCreationYear(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?creationMin=1999&creationMax=2001",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Test Band") {
		t.Error("expected Test Band in filtered results")
	}

	if strings.Contains(body, "Solo Artist") {
		t.Error("Solo Artist should not appear in filtered results")
	}
}

func TestFilterByMembers(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?members=3",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Test Band") {
		t.Error("expected Test Band in filtered results")
	}

	if strings.Contains(body, "Solo Artist") {
		t.Error("Solo Artist should not appear in filtered results")
	}
}

func TestFilterByLocation(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?location=london-uk",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Solo Artist") {
		t.Error("expected Solo Artist in filtered results")
	}

	if strings.Contains(body, "Test Band") {
		t.Error("Test Band should not appear in filtered results")
	}
}

func TestStaticRoute(t *testing.T) {
	moveToProjectRoot(t)
	setTestData(t)

	mux := SetupRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/static/style.css",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}
