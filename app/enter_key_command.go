package app

type EnterKeyCommand struct{}

func (j EnterKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "companies" {
		for _, company := range m.Companies {
			if company.DisplayName == m.Companies[m.CompaniesCursor].DisplayName {
				m.SelectedCompany = company
				break
			}
		}
		m.CurrentView = "categories"
	} else if m.CurrentView == "categories" {
		m.SelectedCategory = m.Categories[m.CategoriesCursor]
		m.CurrentView = "details"
		m.FilesCursor = 0
		m.Files = m.FetchFiles()
	}

	return nil
}

func (j EnterKeyCommand) HelpText() string {
	return "EnterKeyCommand help text"
}

func (j EnterKeyCommand) AllowedStates() []string {
	return []string{}
}
