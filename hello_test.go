package main

import "testing"

func TestHi(t *testing.T) {
	hi := Hi()
	if "Hello"!= hi {
		t.Errorf("Can't even say Hi... got %q", hi)
	}
}