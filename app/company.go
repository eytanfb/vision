package app

import "vision/config"

type Company struct {
	DisplayName    string   `json:"displayName"`
	FolderPathName string   `json:"folderPathName"`
	FullPath       string   `json:"fullPath"`
	SubFolders     []string `json:"subFolders"`
}

func CreateCompanyFromConfigCompany(company config.Company) Company {
	return Company{
		DisplayName:    company.DisplayName,
		FolderPathName: company.FolderPathName,
		FullPath:       company.FullPath,
		SubFolders:     company.SubFolders,
	}
}

func CompaniesFromConfig(companies []config.Company) []Company {
	var result []Company
	for _, company := range companies {
		result = append(result, CreateCompanyFromConfigCompany(company))
	}
	return result
}
