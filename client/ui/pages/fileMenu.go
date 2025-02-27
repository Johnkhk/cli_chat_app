// filemenu.go
package pages

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

// FileMenuModel displays a table of file messages.
type FileMenuModel struct {
	files                  []ChatMessage
	fileTable              table.Model
	rpcClient              *app.RpcClient
	terminalWidth          int
	terminalHeight         int
	originalServerMessages []ChatMessage
	originalActiveUserID   int32
	originalActiveUsername string
	originalSelectedIdx    int
}

// NewFileMenuModel creates a new FileMenuModel.
// The 'files' slice should include all ChatMessage entries with non-text FileType.
func NewFileMenuModel(rpcClient *app.RpcClient, files []ChatMessage, terminalWidth int, terminalHeight int, originalServerMessages []ChatMessage, originalActiveUser int32, originalActiveUsername string, originalSelectedIdx int) FileMenuModel {
	columns := []table.Column{
		{Title: "File Name", Width: 30},
		{Title: "Type", Width: 10},
		{Title: "Size", Width: 10},
		{Title: "Timestamp", Width: 20}, // New timestamp column
	}

	// Build rows from file messages.
	var rows []table.Row
	for _, file := range files {
		sizeStr := strconv.FormatUint(file.FileSize, 10)
		timestampStr := file.Timestamp.Format("2006-01-02 15:04:05")
		rows = append(rows, table.Row{file.FileName, file.FileType, sizeStr, timestampStr})

	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Customize styles.
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	return FileMenuModel{
		files:                  files,
		fileTable:              t,
		rpcClient:              rpcClient,
		terminalWidth:          terminalWidth,
		terminalHeight:         terminalHeight,
		originalServerMessages: originalServerMessages,
		originalActiveUserID:   originalActiveUser,
		originalActiveUsername: originalActiveUsername,
		originalSelectedIdx:    originalSelectedIdx,
	}
}

// Init is part of the tea.Model interface.
func (m FileMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles key events for the file menu.
func (m FileMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "b":
			// Go back to the previous model.
			chatPanelModel := NewChatPanelModel(m.rpcClient)

			// Initialize the chat panel (e.g., fetch friends) and get the command for it
			chatCmd := chatPanelModel.Init() // Init returns the command for loading friends
			chatPanelModel.friendsModel.selected = m.originalSelectedIdx

			// Update the chat panel to handle the incoming message
			updatedChatPanelModel, updateCmd := chatPanelModel.Update(tea.WindowSizeMsg{Width: m.terminalWidth, Height: m.terminalHeight})
			castedChatPanelModel, ok := updatedChatPanelModel.(ChatPanelModel)
			if !ok {
				m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatPanelModel")
				return m, nil
			}
			castedChatPanelModel.chatModel.serverMessages = m.originalServerMessages
			return castedChatPanelModel, tea.Batch(chatCmd, updateCmd)
		case "enter":
			// When Enter is pressed, open the selected file.
			selectedRow := m.fileTable.Cursor()
			if selectedRow < 0 || selectedRow >= len(m.files) {
				return m, nil
			}
			file := m.files[selectedRow]
			// Write file data to a temporary file.
			tempDir := os.TempDir()
			tempPath := filepath.Join(tempDir, file.FileName)
			if err := os.WriteFile(tempPath, file.FileData, 0644); err != nil {
				m.rpcClient.Logger.Errorf("Failed to write file %s: %v", file.FileName, err)
				return m, nil
			}
			// Open the file with the OS default command.
			var openCmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin":
				openCmd = exec.Command("open", tempPath)
			case "linux":
				openCmd = exec.Command("xdg-open", tempPath)
			case "windows":
				openCmd = exec.Command("cmd", "/C", "start", tempPath)
			default:
				m.rpcClient.Logger.Error("Unsupported OS for opening files")
				return m, nil
			}
			if err := openCmd.Start(); err != nil {
				m.rpcClient.Logger.Errorf("Failed to open file: %v", err)
			}
			// Optionally, exit the file menu after opening.
			// return m, tea.Quit
			return m, nil
		default:
			m.fileTable, cmd = m.fileTable.Update(msg)
			return m, cmd
		}
	default:
		m.fileTable, cmd = m.fileTable.Update(msg)
		return m, cmd
	}
}

// View renders the file menu.
func (m FileMenuModel) View() string {
	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Render("Files")
	b.WriteString(title)
	b.WriteString("\n\n")
	b.WriteString(m.fileTable.View())
	b.WriteString("\n")
	b.WriteString("[ ↑/↓ to navigate | Enter to open | b to go back | q/esc to quit ]")
	return b.String()
}
