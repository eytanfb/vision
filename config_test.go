package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	file := "config.json"

	config, err := LoadConfig(file)
	if err != nil {
		t.Errorf("Error loading config file: %s", err)
	}

	if len(config.Companies) == 0 {
		t.Errorf("No companies found in config file")
	}

	if len(config.Companies) != 3 {
		t.Errorf("Expected 3 companies, got %d", len(config.Companies))
	}

	if len(config.Categories) == 0 {
		t.Errorf("No categories found in config file")
	}

	if len(config.Categories) != 9 {
		t.Errorf("Expected 9 categories, got %d", len(config.Categories))
	}
}
