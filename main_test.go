package main

import "testing"

func TestHash(t *testing.T) {
	result := Hash("hello")
	expected := "5d41402abc4b2a76b9719d911017c592"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
