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
		rpcClient: rpcClient,
		tabs:      []string{"Chat", "Friends"},
		activeTab: 0, // Default to the first tab (Chat)
		// tabContent: []tea.Model{chatModel, friendsModel},
		tabContent: []tea.Model{&chatModel, &friendsModel}, // Store friendsModel as a pointer

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
	// windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Top).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	// windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	// windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder())
)

// Update function for main menu to handle key inputs and window resizing
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resizing
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		// Update the size of all child models
		for i, content := range m.tabContent {
			if friendsContent, ok := content.(*friendsModel); ok {
				m.rpcClient.Logger.Infof("Updating tab %d with width %d and height %d", i, m.terminalWidth, m.terminalHeight)
				// friendsContent.list.SetSize(m.terminalWidth, m.terminalHeight)
				w := int(0.8 * float64(m.terminalWidth))
				h := int(0.5 * float64(m.terminalHeight))
				friendsContent.list.SetSize(w, h)

				m.tabContent[i] = friendsContent
			}
		}
		// if friendsContent, ok := m.tabContent[m.activeTab].(*friendsModel); ok {

		// 	m.rpcClient.Logger.Infof("Updating active tab %d with width %d and height %d", m.activeTab, m.terminalWidth, m.terminalHeight)
		// 	w := int(0.8 * float64(m.terminalWidth))
		// 	h := int(0.5 * float64(m.terminalHeight))
		// 	friendsContent.list.SetSize(w, h)
		// 	m.tabContent[m.activeTab] = friendsContent
		// }

	case tea.KeyMsg:
		// Handle key inputs
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.rpcClient.Logger.Info("Exiting the application from main menu")
			return m, tea.Quit

		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)

		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
		}
	}

	// Update the currently active tab's content
	updatedModel, subCmd := m.tabContent[m.activeTab].Update(msg)
	if updatedContent, ok := updatedModel.(tea.Model); ok {
		m.tabContent[m.activeTab] = updatedContent
	}
	cmds = append(cmds, subCmd)

	// Combine all commands using tea.Batch
	cmd = tea.Batch(cmds...)

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

	// Calculate available width and height for child components
	availableWidth := (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())
	availableHeight := m.terminalHeight - 5 // Adjust height as needed

	// Render the content of the active tab within the calculated dimensions
	doc.WriteString(
		windowStyle.
			Width(availableWidth).
			Height(availableHeight).
			Render(m.tabContent[m.activeTab].View()),
	)

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
