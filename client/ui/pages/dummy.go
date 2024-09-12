package pages

import tea "github.com/charmbracelet/bubbletea"

type DummyModel struct{}

// Init function initializes the main menu model
func (m DummyModel) Init() tea.Cmd {
	return nil
}

func (m DummyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m DummyModel) View() string {
	return "DUMMY"
}

func NewDummyModel() DummyModel {
	return DummyModel{}
}
