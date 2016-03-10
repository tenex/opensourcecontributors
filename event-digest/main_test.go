package main

import (
	"os"
	"testing"
)

func TestStreamReaderWithCrazyNullJSON(t *testing.T) {
	eventFile, err := os.Open("../fixtures/2990966171.json")
	users := UsernameSet{}
	evtCount, err := digestStream(eventFile, users)
	if err != nil {
		t.Error(err)
	}
	if len(users) != 2 {
		t.Errorf("got %v users, expected 2", len(users))
	}
	if evtCount != 2 {
		t.Errorf("got %v events, expected 2", evtCount)
	}
}
