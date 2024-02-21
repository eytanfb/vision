package app

type KKeyCommand struct{}

func (j KKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "details" {
		if m.ItemDetailsFocus {
			m.Viewport.LineUp(10)
		} else {
			m.FilesCursor--
			if m.FilesCursor < 0 {
				m.FilesCursor = 0
			}
		}
	} else if m.CurrentView == "categories" {
		m.CategoriesCursor--
		if m.CategoriesCursor < 0 {
			m.CategoriesCursor = 0
		}
	} else if m.CurrentView == "companies" {
		m.CompaniesCursor--
		if m.CompaniesCursor < 0 {
			m.CompaniesCursor = 0
		}
	}

	return nil
}

func (j KKeyCommand) HelpText() string {
	return "KKeyCommand help text"
}

func (j KKeyCommand) AllowedStates() []string {
	return []string{}
}
