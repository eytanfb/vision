package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Company struct {
	DisplayName    string   `json:"displayName"`
	FolderPathName string   `json:"folderPathName"`
	FullPath       string   `json:"fullPath"`
	SubFolders     []string `json:"subFolders"`
	Color          string   `json:"color"`
}

type Config struct {
	Companies      []Company `json:"companies"`
	Categories     []string
	DefaultCompany string
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file at %s: %w", path, err)
	}

	if file == nil {
		return nil, fmt.Errorf("could not open file at %s", path)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode file at %s: %w", path, err)
	}

	config.Categories = LoadCategories(&config)

	defaultCompany := os.Getenv("VISION_DEFAULT_COMPANY")
	if defaultCompany == "" {
		defaultCompany = "clerky"
	}

	config.DefaultCompany = defaultCompany

	return &config, nil
}

func LoadCategories(config *Config) []string {
	categoryList := []string{}

	for _, company := range config.Companies {
		for _, subFolder := range company.SubFolders {
			categoryList = append(categoryList, subFolder)
		}
	}

	return uniqueStrings(categoryList)
}

func uniqueStrings(input []string) []string {
	unique := make(map[string]bool)
	var result []string

	for _, value := range input {
		if _, ok := unique[value]; !ok {
			unique[value] = true
			result = append(result, value)
		}
	}

	return result
}
