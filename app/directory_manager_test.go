package app

import (
	"testing"
)

func TestDirectoryManager_CurrentFolderPathReturnsSelectedCompanyFolderPathName(t *testing.T) {
	// Arrange
	companyName := "TestCompany"
	company := Company{DisplayName: companyName, FolderPathName: companyName}
	dm := DirectoryManager{SelectedCompany: company}

	// Act
	result := dm.CurrentFolderPath()

	// Assert
	if result != companyName {
		t.Errorf("Expected %s, got %s", companyName, result)
	}
}
