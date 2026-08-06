package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestMainPrintsServerError(t *testing.T) {
	originalTransport := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("test server error")
	})
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not create stdout pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = originalStdout
		reader.Close()
		writer.Close()
	})

	main()

	if err := writer.Close(); err != nil {
		t.Fatalf("could not close stdout writer: %v", err)
	}
	os.Stdout = originalStdout

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("could not read stdout: %v", err)
	}

	expected := "Server error: could not fetch artists: " +
		"Get \"https://groupietrackers.herokuapp.com/api/artists\": test server error\n"
	if string(output) != expected {
		t.Errorf("expected output %q, got %q", expected, output)
	}
}
