package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func handlerTestData() ([]Artist, []Relation, []Location, []Date) {
	artists := []Artist{
		{
			ID:           1,
			Name:         "Test Band",
			Members:      []string{"A", "B", "C"},
			CreationDate: 2000,
			FirstAlbum:   "01-01-2005",
		},
		{
			ID:           2,
			Name:         "Solo Artist",
			Members:      []string{"D"},
			CreationDate: 2015,
			FirstAlbum:   "01-01-2016",
		},
		{
			ID:           3,
			Name:         "Large Band",
			Members:      []string{"A", "B", "C", "D", "E"},
			CreationDate: 1995,
			FirstAlbum:   "01-01-1998",
		},
	}

	relations := []Relation{
		{
			ID: 1,
			DatesLocations: map[string][]string{
				"athens-greece": {"01-01-2026"},
			},
		},
		{
			ID: 2,
			DatesLocations: map[string][]string{
				"london-uk": {"02-02-2026"},
			},
		},
		{
			ID: 3,
			DatesLocations: map[string][]string{
				"texas-usa": {"03-03-2026"},
			},
		},
	}

	locations := []Location{
		{
			ID:        1,
			Locations: []string{"athens-greece"},
		},
		{
			ID:        2,
			Locations: []string{"london-uk"},
		},
		{
			ID:        3,
			Locations: []string{"texas-usa"},
		},
	}

	dates := []Date{
		{ID: 1, Dates: []string{"01-01-2026"}},
		{ID: 2, Dates: []string{"02-02-2026"}},
		{ID: 3, Dates: []string{"03-03-2026"}},
	}

	return artists, relations, locations, dates
}

func moveHandlersTestsToRoot(t *testing.T) {
	t.Helper()

	currentDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current directory: %v", err)
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

func TestHandlerHomePage(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	Handler(artists, locations).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Test Band") {
		t.Error("expected Test Band on home page")
	}

	if !strings.Contains(body, "Athens Greece") {
		t.Error("expected formatted location in location list")
	}
}

func TestHandlerInvalidPath(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/not-found",
		nil,
	)
	response := httptest.NewRecorder()

	Handler(artists, locations).ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.Code,
		)
	}
}

func TestHandlerRejectsPost(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	Handler(artists, locations).ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			response.Code,
		)
	}
}

func TestSearchHandler(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, _, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?query=solo",
		nil,
	)
	response := httptest.NewRecorder()

	SearchHandler(artists).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Solo Artist") {
		t.Error("expected Solo Artist in search results")
	}

	if strings.Contains(body, "Test Band") {
		t.Error("Test Band should not appear in search results")
	}
}

func TestSearchHandlerEmptyQuery(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, _, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?query=",
		nil,
	)
	response := httptest.NewRecorder()

	SearchHandler(artists).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}
}

func TestSuggestionHandler(t *testing.T) {
	artists, _, _, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/suggestions?query=large",
		nil,
	)
	response := httptest.NewRecorder()

	SuggestionHandler(artists).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var results []Artist

	err := json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Name != "Large Band" {
		t.Errorf("expected Large Band, got %s", results[0].Name)
	}
}

func TestArtistDetailsHandler(t *testing.T) {
	artists, relations, locations, dates := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/artist?id=1",
		nil,
	)
	response := httptest.NewRecorder()

	ArtistDetailsHandler(
		artists,
		relations,
		locations,
		dates,
	).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var result Find

	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("could not decode artist details: %v", err)
	}

	if result.Artist.Name != "Test Band" {
		t.Errorf(
			"expected Test Band, got %s",
			result.Artist.Name,
		)
	}

	if result.LocationCount != 1 {
		t.Errorf(
			"expected 1 location, got %d",
			result.LocationCount,
		)
	}

	if result.ConcertCount != 1 {
		t.Errorf(
			"expected 1 concert, got %d",
			result.ConcertCount,
		)
	}
}

func TestArtistDetailsInvalidID(t *testing.T) {
	artists, relations, locations, dates := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/artist?id=wrong",
		nil,
	)
	response := httptest.NewRecorder()

	ArtistDetailsHandler(
		artists,
		relations,
		locations,
		dates,
	).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}
}

func TestHandlerFilterByCreationYear(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?creationMin=1999&creationMax=2001",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Test Band") {
		t.Error("expected Test Band in filtered results")
	}

	if strings.Contains(body, "Solo Artist") {
		t.Error("Solo Artist should not appear")
	}
}

func TestFilterByFirstAlbumYear(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?firstAlbumMin=1997&firstAlbumMax=1999",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	body := response.Body.String()

	if !strings.Contains(body, "Large Band") {
		t.Error("expected Large Band in filtered results")
	}

	if strings.Contains(body, "Test Band") {
		t.Error("Test Band should not appear")
	}
}

func TestHandlerFilterByMembers(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?members=3",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	body := response.Body.String()

	if !strings.Contains(body, "Test Band") {
		t.Error("expected Test Band with 3 members")
	}

	if strings.Contains(body, "Large Band") {
		t.Error("Large Band should not appear")
	}
}

func TestFilterAtLeastMembersValues(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?members=3&members=4&members=5&members=6&members=7&members=8",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	body := response.Body.String()

	if !strings.Contains(body, "Test Band") {
		t.Error("expected Test Band with at least 3 members")
	}

	if !strings.Contains(body, "Large Band") {
		t.Error("expected Large Band with at least 3 members")
	}

	if strings.Contains(body, "Solo Artist") {
		t.Error("Solo Artist should not appear")
	}
}

func TestHandlerFilterByLocation(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?location=london-uk",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	body := response.Body.String()

	if !strings.Contains(body, "Solo Artist") {
		t.Error("expected Solo Artist for London")
	}

	if strings.Contains(body, "Test Band") {
		t.Error("Test Band should not appear")
	}
}

func TestCombinedFilters(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?creationMin=1990&creationMax=2000&members=5&location=texas-usa",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	body := response.Body.String()

	if !strings.Contains(body, "Large Band") {
		t.Error("expected Large Band in combined results")
	}

	if strings.Contains(body, "Test Band") {
		t.Error("Test Band should not appear")
	}
}

func TestFilterRejectsInvalidYear(t *testing.T) {
	moveHandlersTestsToRoot(t)

	artists, _, locations, _ := handlerTestData()

	request := httptest.NewRequest(
		http.MethodGet,
		"/filter?creationMin=wrong",
		nil,
	)
	response := httptest.NewRecorder()

	FilterHandler(artists, locations).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}
}

func TestFormatLocation(t *testing.T) {
	functions := indexFuncMap()

	formatLocation, ok := functions["formatLocation"].(func(string) string)
	if !ok {
		t.Fatal("formatLocation was not found")
	}

	result := formatLocation("north_carolina-usa")

	if result != "North Carolina USA" {
		t.Errorf(
			"expected North Carolina USA, got %s",
			result,
		)
	}
}
