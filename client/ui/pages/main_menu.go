package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

// mainMenuModel structure with terminal size fields and state
type mainMenuModel struct {
	rpcClient      *app.RpcClient
	terminalWidth  int
	terminalHeight int
	activeTab      int
	tabs           []string
	tabContent     []tea.Model
}

// Initialize the main menu model with tabs
func NewMainMenuModel(rpcClient *app.RpcClient) mainMenuModel {
	// Create dummy models for the tab contents
	chatModel := NewDummyModel() // Replace with your actual Chat tab model
	friendsModel := NewFriendsModel(rpcClient)

	return mainMenuModel{
		rpcClient:  rpcClient,
		tabs:       []string{"Chat", "Friends"},
		activeTab:  0, // Default to the first tab (Chat)
		tabContent: []tea.Model{chatModel, friendsModel},
		// tabContent: []string{chatModel, friendsModel},
	}
}

// Tab border styling
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	// windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder())
)

// Update function for main menu to handle key inputs and window resizing
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.rpcClient.Logger.Info("Exiting the application from main menu")
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	// Update the currently active tab's content
	updatedModel, cmd := m.tabContent[m.activeTab].Update(msg)
	m.tabContent[m.activeTab] = updatedModel

	// Return the updated mainMenuModel and any command to be executed
	return m, cmd
}

// View function renders the Main Menu UI with tabs
func (m mainMenuModel) View() string {
	doc := strings.Builder{}

	// Render the tabs
	var renderedTabs []string
	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		// Adjust borders for the tabs
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
			// border.BottomLeft = "│"

		} else if isLast && isActive {
			// border.BottomRight = "┐" // Adjust to match content window's top corner
			border.BottomRight = "│" // Adjust to match content window's top corner

		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Width(m.terminalWidth/len(m.tabs)-5).Render(t))
		// renderedTabs = append(renderedTabs, style.Width(m.terminalWidth/len(m.tabs)).Margin(4, 4).Render(t))
	}

	// Combine the rendered tabs into a single row
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Height(m.terminalHeight - 5).Render(m.tabContent[m.activeTab].View()))

	return docStyle.Align(lipgloss.Center).Width(m.terminalWidth).Height(m.terminalHeight).
		Render(doc.String())
}

// Init function initializes the main menu model
func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

// Helper functions for min and max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
