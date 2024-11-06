package app

import "strings"

type DirectoryManager struct {
	Companies         []Company
	Categories        []string
	CompanyCategories map[string][]string
	SelectedCompany   Company
	SelectedCategory  string
	CompaniesCursor   int
	CategoriesCursor  int
}

func (dm *DirectoryManager) CurrentFolderPath() string {
	return dm.SelectedCompany.FolderPathName
}

func (dm *DirectoryManager) CompanyNames() []string {
	var names []string
	for _, company := range dm.Companies {
		names = append(names, company.DisplayName)
	}
	return names
}

func (dm *DirectoryManager) AssignCompany() {
	dm.SelectedCompany = dm.Companies[dm.CompaniesCursor]
}

func (dm *DirectoryManager) AssignCategory() {
	dm.SelectedCategory = dm.Categories[dm.CategoriesCursor]
}

func (dm *DirectoryManager) SelectCompany(companyName string) bool {
	for index, company := range dm.Companies {
		if strings.ToLower(company.DisplayName) == companyName {
			dm.SelectedCompany = company
			dm.CompaniesCursor = index
			dm.Categories = dm.CompanyCategories[company.FolderPathName]
			return true
		}
	}

	return false
}

func (dm *DirectoryManager) SelectCategory(categoryName string) bool {
	for _, category := range dm.Categories {
		if strings.ToLower(category) == categoryName {
			dm.SelectedCategory = categoryName
			return true
		}
	}

	return false
}

func (dm *DirectoryManager) CurrentCompanyName() string {
	return dm.SelectedCompany.DisplayName
}
