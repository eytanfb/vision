package app

type ViewManager struct {
	CurrentView       string
	Width             int
	Height            int
	Ready             bool
	TaskDetailsFocus  bool
	ItemDetailsFocus  bool
	HideSidebar       bool
	NavbarWidth       int
	SidebarWidth      int
	SidebarHeight     int
	DetailsViewWidth  int
	DetailsViewHeight int
	SummaryViewHeight int
	IsAddTaskView     bool
}

const (
	CompaniesView  = "companies"
	CategoriesView = "categories"
	DetailsView    = "details"
	AddTaskView    = "add_task"
)

func (vm *ViewManager) SetWidth(width int) {
	vm.Width = width
	vm.DetailsViewWidth = width - vm.SidebarWidth - 9
	vm.NavbarWidth = width - 5
}

func (vm *ViewManager) SetHeight(height int) {
	vm.Height = height
	vm.SidebarHeight = height - 10
	vm.SummaryViewHeight = height - 10
	vm.DetailsViewHeight = height - 20
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
	if vm.IsDetailsView() {
		vm.CurrentView = CategoriesView
	}
}

func (vm *ViewManager) GoToNextView(fm *FileManager, dm *DirectoryManager, tm *TaskManager) {
	if vm.IsCompanyView() {
		vm.CurrentView = CategoriesView
	} else if vm.IsCategoryView() {
		vm.CurrentView = DetailsView
		fm.FilesCursor = 0
		fm.Files = fm.FetchFiles(dm, tm)
	}
}

func (vm *ViewManager) Select(fm *FileManager, dm *DirectoryManager, tm *TaskManager) {
	if vm.IsCompanyView() {
		dm.AssignCompany()
	} else if vm.IsCategoryView() {
		dm.AssignCategory()
	}
	vm.GoToNextView(fm, dm, tm)
}

func (vm *ViewManager) ToggleHideSidebar() {
	vm.HideSidebar = !vm.HideSidebar
	if vm.HideSidebar {
		vm.DetailsViewWidth += vm.SidebarWidth
	} else {
		vm.DetailsViewWidth -= vm.SidebarWidth
	}
}
