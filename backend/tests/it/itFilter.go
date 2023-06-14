package it

import (
	"os"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping integration test")
	}
}

func MainFilter() {
	if os.Getenv("INTEGRATION") == "" {
		os.Exit(0)
	}
}
