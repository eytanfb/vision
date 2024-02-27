package app

type ViewManager struct {
	CurrentView      string
	Width            int
	Height           int
	Ready            bool
	TaskDetailsFocus bool
	ItemDetailsFocus bool
}

func (vm ViewManager) IsCompanyView() bool {
	return vm.CurrentView == CompaniesView
}

func (vm ViewManager) IsCategoryView() bool {
	return vm.CurrentView == CategoriesView
}

func (vm ViewManager) IsDetailsView() bool {
	return vm.CurrentView == DetailsView
}

func (vm ViewManager) IsTaskDetailsFocus() bool {
	return vm.IsDetailsView() && vm.TaskDetailsFocus
}

func (vm ViewManager) IsItemDetailsFocus() bool {
	return vm.IsDetailsView() && vm.ItemDetailsFocus
}

func (vm *ViewManager) GoToPreviousView() {
	if vm.IsCategoryView() {
		vm.CurrentView = CompaniesView
	} else if vm.IsDetailsView() {
		vm.CurrentView = CategoriesView
	}
}

func (vm *ViewManager) GoToNextView(fm *FileManager, dm *DirectoryManager) {
	if vm.IsCompanyView() {
		vm.CurrentView = CategoriesView
	} else if vm.IsCategoryView() {
		vm.CurrentView = DetailsView
		fm.FilesCursor = 0
		fm.Files = fm.FetchFiles(dm)
	}
}

func (vm *ViewManager) Select(fm *FileManager, dm *DirectoryManager) {
	if vm.IsCompanyView() {
		dm.AssignCompany()
	} else if vm.IsCategoryView() {
		dm.AssignCategory()
	}
	vm.GoToNextView(fm, dm)
}
