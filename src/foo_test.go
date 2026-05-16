package main

import (
	"testing"
)

func TestFoo (t *testing.T) {
	want := "foo"
	if want != "foo" {
		t.Error("Ahoy")
	}
}
