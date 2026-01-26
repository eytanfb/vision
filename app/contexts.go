package app

// Context constants define the different view/state contexts
// that commands can be restricted to.
const (
	ContextKanban        = "kanban"
	ContextDetails       = "details"
	ContextItemDetails   = "item_details"
	ContextCategoryView  = "category_view"
	ContextWeeklyView    = "weekly_view"
	ContextRetrospective = "retrospective"
)
