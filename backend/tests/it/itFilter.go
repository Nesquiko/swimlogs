package it

import (
	"fmt"
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
		fmt.Println("Skipping integration test")
		os.Exit(0)
	}
}
