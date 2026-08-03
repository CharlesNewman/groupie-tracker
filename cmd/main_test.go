package main

import (
	"errors"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	originalStartServer := startServer
	defer func() {
		startServer = originalStartServer
	}()

	startServer = func() error {
		return nil
	}

	err := run()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRunServerError(t *testing.T) {
	originalStartServer := startServer
	defer func() {
		startServer = originalStartServer
	}()

	expectedError := errors.New("test server error")

	startServer = func() error {
		return expectedError
	}

	err := run()

	if !errors.Is(err, expectedError) {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}
