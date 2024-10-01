package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChatMessage represents a message in the chat.
type ChatMessage struct {
	Sender  string // "self" or "other"
	Message string
}

type ChatModel struct {
	viewport   viewport.Model
	messages   []ChatMessage
	textarea   textarea.Model
	selfStyle  lipgloss.Style
	otherStyle lipgloss.Style
	err        error
}

func NewChatModel() ChatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(50, 10) // Increased viewport width to show more messages.
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChatModel{
		textarea:   ta,
		messages:   []ChatMessage{},
		viewport:   vp,
		selfStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")).PaddingLeft(2),                        // Purple for self messages with left padding
		otherStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Align(lipgloss.Right).PaddingRight(2), // Green for other messages and right-aligned
		err:        nil,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			// Add the user's message with "self" style
			m.messages = append(m.messages, ChatMessage{
				Sender:  "self",
				Message: m.textarea.Value(),
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()

		// Simulate receiving an "other" message on pressing Ctrl+A (for testing purposes)
		case tea.KeyCtrlA:
			m.messages = append(m.messages, ChatMessage{
				Sender:  "other",
				Message: "This is a reply from the other person.",
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

// renderMessages iterates over the chat messages and applies styles based on sender.
func (m ChatModel) renderMessages() string {
	var renderedMessages []string
	viewportWidth := m.viewport.Width // Use viewport width to determine right alignment.
	for _, msg := range m.messages {
		var styledMessage string
		switch msg.Sender {
		case "self":
			// Render self messages with the left-aligned style.
			styledMessage = m.selfStyle.Render(fmt.Sprintf("You: %s", msg.Message))
		case "other":
			// Render other messages right-aligned with padding to push to the right side.
			msgContent := fmt.Sprintf("Other: %s", msg.Message)
			styledMessage = m.otherStyle.Render(msgContent)
			spacesNeeded := viewportWidth - lipgloss.Width(msgContent) - 2 // Adjust right padding
			if spacesNeeded > 0 {
				styledMessage = lipgloss.NewStyle().MarginLeft(spacesNeeded).Render(styledMessage)
			}
		}
		renderedMessages = append(renderedMessages, styledMessage)
	}
	return strings.Join(renderedMessages, "\n")
}
