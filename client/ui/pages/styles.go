package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	// Title Style
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")). // Example color for the title
			Bold(true).                       // Makes the title bold
			Underline(true).                  // Optionally underline the title
			Align(lipgloss.Center)            // Aligns the title to the center

	// Buttons
	focusedSubmitButton = focusedStyle.Render("[ Submit ]")
	blurredSubmitButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedBackButton   = focusedStyle.Render("[ Back ]")
	blurredBackButton   = fmt.Sprintf("[ %s ]", blurredStyle.Render("Back"))

	// main menu
	leftPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Align(lipgloss.Center)
	rightPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("205")).
			Align(lipgloss.Center)
	grayBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")) // Gray color for inactive borders

	blueBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")) // Blue color for active borders

)
