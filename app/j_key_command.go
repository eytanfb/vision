package app

type JKeyCommand struct{}

func (j JKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "companies" {
		m.CompaniesCursor++
		if m.CompaniesCursor >= len(m.Companies) {
			m.CompaniesCursor = len(m.Companies) - 1
		}
	} else if m.CurrentView == "categories" {
		m.CategoriesCursor++
		if m.CategoriesCursor >= len(m.Categories) {
			m.CategoriesCursor = len(m.Categories) - 1
		}
	} else if m.CurrentView == "details" {
		if m.ItemDetailsFocus {
			if m.TaskDetailsFocus {
				m.TasksCursor++
				if m.TasksCursor >= len(m.Tasks) {
					m.TasksCursor = len(m.Tasks) - 1
				}
			} else {
				m.Viewport.LineDown(10)
			}
		} else {
			m.FilesCursor++
			if m.FilesCursor >= len(m.Files) {
				m.FilesCursor = len(m.Files) - 1
			}
		}
	}

	return nil
}

func (j JKeyCommand) HelpText() string {
	return "JKeyCommand help text"
}

func (j JKeyCommand) AllowedStates() []string {
	return []string{}
}
