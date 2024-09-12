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

	// Logo Style
	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")). // Example color for the title
			Bold(true).                       // Makes the title bold
			Align(lipgloss.Center)            // Aligns the title to the center

	// Title Style
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")). // Example color for the logo
			Bold(true).                        // Makes the logo bold
			Underline(true).                   // Optionally underline the logo
			Align(lipgloss.Center)

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
			BorderForeground(lipgloss.Color("240")). // Gray color for inactive borders
			Padding(0, 2)                            // Optional padding to avoid cutting off content

	blueBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")). // Blue color for active borders
			Padding(0, 2)                           // Optional padding to avoid cutting off content

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	// windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
