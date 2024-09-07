package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define the UI model
type UIModel struct {
	content string // Holds the content to display in the UI
}

func (m UIModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// Function to initialize the UI model
func InitialUIModel() UIModel {
	return UIModel{
		content: "Welcome to CLI Chat App! Type '/help' for commands.",
	}
}

// Update function for the UI
func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View function to render the UI
func (m UIModel) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(50).
		Render

	ui := border(m.content)

	return ui
}

// Start the Bubble Tea program
func StartUI() {
	p := tea.NewProgram(InitialUIModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
