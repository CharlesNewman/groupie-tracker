package main

import (
	"errors"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	originalStartServer := startServer

	startServer = func() error {
		return nil
	}

	t.Cleanup(func() {
		startServer = originalStartServer
	})

	err := run()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRunReturnsServerError(t *testing.T) {
	originalStartServer := startServer

	expectedError := errors.New("test server error")

	startServer = func() error {
		return expectedError
	}

	t.Cleanup(func() {
		startServer = originalStartServer
	})

	err := run()

	if !errors.Is(err, expectedError) {
		t.Errorf(
			"expected error %v, got %v",
			expectedError,
			err,
		)
	}
}
