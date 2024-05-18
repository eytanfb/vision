package app

type ViewManager struct {
	CurrentView              string
	Width                    int
	Height                   int
	Ready                    bool
	TaskDetailsFocus         bool
	ItemDetailsFocus         bool
	HideSidebar              bool
	NavbarWidth              int
	NavbarHeight             int
	SidebarWidth             int
	SidebarHeight            int
	DetailsViewWidth         int
	DetailsViewHeight        int
	SummaryViewHeight        int
	IsAddTaskView            bool
	IsAddSubTaskView         bool
	IsWeeklyView             bool
	IsFilterView             bool
	ShowCompanies            bool
	KanbanListCursor         int
	KanbanTaskCursor         int
	KanbanTasksCount         int
	KanbanViewLineDownFactor int
}

const (
	CompaniesView      = "companies"
	CategoriesView     = "categories"
	DetailsView        = "details"
	heightOffset       = 12
	detailsWidthOffset = 9
	navbarWidthOffset  = 5
)

func (vm *ViewManager) SetWidth(width int) {
	vm.Width = width
	vm.DetailsViewWidth = width - vm.SidebarWidth - detailsWidthOffset
	vm.NavbarWidth = width - navbarWidthOffset
}

func (vm *ViewManager) SetHeight(height int) {
	vm.Height = height
	vm.SidebarHeight = height - heightOffset
	vm.SummaryViewHeight = height - heightOffset
	vm.DetailsViewHeight = height - heightOffset
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
	if vm.IsDetailsView() && !vm.HideSidebar {
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
	if !vm.HideSidebar {
		if vm.IsCompanyView() {
			dm.AssignCompany()
		} else if vm.IsCategoryView() {
			dm.AssignCategory()
		}
		vm.GoToNextView(fm, dm, tm)
	}
}

func (vm *ViewManager) ToggleHideSidebar() {
	vm.HideSidebar = !vm.HideSidebar

	if vm.HideSidebar {
		vm.DetailsViewWidth += vm.SidebarWidth
		vm.ItemDetailsFocus = true
	} else {
		vm.DetailsViewWidth -= vm.SidebarWidth
		vm.ItemDetailsFocus = false
		vm.TaskDetailsFocus = false
	}
}

func (vm *ViewManager) ToggleWeeklyView() {
	vm.IsWeeklyView = !vm.IsWeeklyView
}

func (vm *ViewManager) KanbanLineDownAmount() int {
	return vm.KanbanTaskCursor / 1 * vm.KanbanViewLineDownFactor
}
