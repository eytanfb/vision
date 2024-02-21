package app

import (
	"testing"
)

func TestRemoveDatesFromText(t *testing.T) {
	text := "- [ ] Task 1  🛫  2021-07-01"
	expected := "- [ ] Task 1"
	result := RemoveDatesFromText(text)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
