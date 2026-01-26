package app

import "strings"

type DirectoryManager struct {
	Companies        []Company
	Categories       []string
	SelectedCompany  Company
	SelectedCategory string
	CompaniesCursor  int
	CategoriesCursor int
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
	if dm.CompaniesCursor >= 0 && dm.CompaniesCursor < len(dm.Companies) {
		dm.SelectedCompany = dm.Companies[dm.CompaniesCursor]
	}
}

func (dm *DirectoryManager) AssignCategory() {
	if dm.CategoriesCursor >= 0 && dm.CategoriesCursor < len(dm.Categories) {
		dm.SelectedCategory = dm.Categories[dm.CategoriesCursor]
	}
}

func (dm *DirectoryManager) SelectCompany(companyName string) bool {
	for index, company := range dm.Companies {
		if strings.ToLower(company.DisplayName) == companyName {
			dm.SelectedCompany = company
			dm.CompaniesCursor = index
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
