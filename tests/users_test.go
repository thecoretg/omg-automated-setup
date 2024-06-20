package main

import (
	"testing"

	"github.com/thecoretg/omg-user-automation/internal/mac"
)

func TestGetUser(t *testing.T) {
	// run checkroot
	err := mac.CheckRoot()
	if err != nil {
		t.Errorf("Error checking root: %v", err)
	}

}

func TestLookupUser(t *testing.T) {
	exists, err := mac.CheckUserExists("danny")
	if err != nil {
		t.Errorf("Error checking user: %v", err)
	}

	if !exists {
		t.Errorf("User does not exist")
	}

}
