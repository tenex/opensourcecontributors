package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestNullReplacer(t *testing.T) {
	badString := "Hello\x00World\x00"
	goodString := "Hello World "
	replacer := NewNullReplacer(strings.NewReader(badString), ' ')
	out, err := ioutil.ReadAll(replacer)
	if err != nil {
		t.Error(err)
	}
	if goodString != string(out) {
		t.Errorf("wanted [%s], got [%s]", goodString, out)
	}
}
