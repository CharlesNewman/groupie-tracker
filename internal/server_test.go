package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func testArtists() []Artist {
	return []Artist{
		{
			ID:           1,
			Name:         "Queen",
			Image:        "queen.jpg",
			Members:      []string{"Freddie Mercury", "Brian May"},
			CreationDate: 1970,
			FirstAlbum:   "14-12-1973",
		},
	}
}

func testRelations() []Relation {
	return []Relation{
		{
			ID: 1,
			DatesLocations: map[string][]string{
				"london-uk": {"01-01-2026", "02-01-2026"},
			},
		},
	}
}

func testLocations() []Location {
	return []Location{
		{
			ID: 1,
			Locations: []string{
				"london-uk",
				"athens-greece",
			},
		},
	}
}

func testDates() []Date {
	return []Date{
		{
			ID: 1,
			Dates: []string{
				"01-01-2026",
				"02-01-2026",
			},
		},
	}
}

func setupTestRoutes() *http.ServeMux {
	return SetupRoutes(
		testArtists(),
		testRelations(),
		testLocations(),
		testDates(),
	)
}

func TestHomeRoute(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
}

func TestInvalidRoute(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", response.Code)
	}
}

func TestHomeWrongMethod(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", response.Code)
	}
}

func TestSearchRoute(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?query=Queen",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "Queen") {
		t.Error("expected response to contain Queen")
	}
}

func TestSuggestionsRoute(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/suggestions?query=Que",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	contentType := response.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected JSON, got %s", contentType)
	}

	if !strings.Contains(response.Body.String(), "Queen") {
		t.Error("expected JSON response to contain Queen")
	}
}

func TestArtistJSONRoute(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/artist?id=1",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Queen") {
		t.Error("expected JSON response to contain Queen")
	}

	if !strings.Contains(body, `"locationCount":2`) {
		t.Error("expected locationCount to be 2")
	}

	if !strings.Contains(body, `"concertCount":2`) {
		t.Error("expected concertCount to be 2")
	}
}

func TestArtistInvalidID(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/artist?id=abc",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", response.Code)
	}
}

func TestArtistNotFound(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/artist?id=999",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", response.Code)
	}
}

func TestArtistDetailsPage(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/artist-details?id=1",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "Queen") {
		t.Error("expected details page to contain Queen")
	}
}

func TestInvalidPageNumber(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/?page=abc",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", response.Code)
	}
}

func TestPageNotFound(t *testing.T) {
	mux := setupTestRoutes()

	request := httptest.NewRequest(
		http.MethodGet,
		"/?page=999",
		nil,
	)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", response.Code)
	}
}
