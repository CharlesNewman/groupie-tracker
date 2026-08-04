package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArtistDetailsBadRequest(t *testing.T) {
	handler := ArtistDetailsHandler(
		[]Artist{},
		[]Relation{},
		[]Location{},
		[]Date{},
	)

	request := httptest.NewRequest(http.MethodGet, "/artist?id=abc", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", response.Code)
	}
}

func TestArtistDetailsNotFound(t *testing.T) {
	handler := ArtistDetailsHandler(
		[]Artist{},
		[]Relation{},
		[]Location{},
		[]Date{},
	)

	request := httptest.NewRequest(http.MethodGet, "/artist?id=999", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", response.Code)
	}
}

func TestArtistDetailsMethodNotAllowed(t *testing.T) {
	handler := ArtistDetailsHandler(
		[]Artist{},
		[]Relation{},
		[]Location{},
		[]Date{},
	)

	request := httptest.NewRequest(http.MethodPost, "/artist?id=1", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", response.Code)
	}
}
