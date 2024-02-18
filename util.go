package main

import "strings"

// IsValidCategory checks if the provided argument is a valid category.
func IsValidCategory(arg string) bool {
	categories := []string{"tasks", "meetings", "standups", "people", "projects"}
	for _, c := range categories {
		if arg == c {
			return true
		}
	}
	return false
}

// FindCompany checks if the company exists in the config and returns it if found.
func FindCompany(arg string, cfg *Config) (Company, bool) {
	// Assuming your config has a way to list companies. Adjust accordingly.
	for _, company := range cfg.Companies {
		if strings.ToLower(arg) == strings.ToLower(company.DisplayName) {
			return company, true
		}
	}
	return Company{}, false
}
