package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

// ChatMessage represents a message in the chat.
type ChatMessage struct {
	Sender  string // "self" or "other"
	Message string
}

// Define a custom message type for received messages.
type ReceivedMessage struct {
	Sender  string
	Message string
}
type ChatModel struct {
	viewport       viewport.Model
	messages       []ChatMessage
	textarea       textarea.Model
	selfStyle      lipgloss.Style
	otherStyle     lipgloss.Style
	err            error
	rpcClient      *app.RpcClient
	ctx            context.Context
	cancel         context.CancelFunc
	activeUserID   int32 // Add this field to track the active user ID
	activeUsername string
}

// NewChatModel initializes a new ChatModel.
func NewChatModel(rpcClient *app.RpcClient) ChatModel {
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

	ctx, cancel := context.WithCancel(context.Background())

	return ChatModel{
		textarea:   ta,
		messages:   []ChatMessage{},
		viewport:   vp,
		selfStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")).PaddingLeft(2),                        // Purple for self messages with left padding
		otherStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Align(lipgloss.Right).PaddingRight(2), // Green for other messages and right-aligned
		err:        nil,
		rpcClient:  rpcClient,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Init initializes the model with a command to listen to the message channel.
func (m ChatModel) Init() tea.Cmd {
	// Create a command to start listening to the message channel.
	m.rpcClient.Logger.Info("Initializing the chat model and starting to listen to the message channel.")
	return tea.Batch(
		textarea.Blink,
		m.listenToMessageChannel(), // Start listening to the message channel.
	)
}

// listenToMessageChannel listens to the message channel and sends messages to the Tea program.
func (m ChatModel) listenToMessageChannel() tea.Cmd {

	return func() tea.Msg {

		for {
			select {
			case <-m.ctx.Done():
				m.rpcClient.Logger.Info("Context cancelled, stopping message channel listener.")
				return nil

			case msg, ok := <-m.rpcClient.ChatClient.MessageChannel:
				if !ok {
					// If the channel is closed, log it and return.
					m.rpcClient.Logger.Warn("Message channel closed. Exiting listener.")
					return nil
				}

				if msg == nil {
					m.rpcClient.Logger.Warn("Received a nil message from the channel, ignoring.")
					continue
				}

				// Log the received message and return it as a ReceivedMessage.
				m.rpcClient.Logger.Infof("Received message from channel: Sender=%s, Message=%s, Status=%s", msg.SenderUsername, string(msg.EncryptedMessage), msg.Status)
				return ReceivedMessage{
					Sender:  msg.SenderUsername,
					Message: string(msg.EncryptedMessage),
				}
			}
		}
	}
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
			m.cancel() // Cancel the context to stop listening to the stream.
			return m, tea.Quit
		case tea.KeyEnter:
			m.rpcClient.Logger.Info("Enter key pressed.")

			// Check if the textarea has content before sending the message.
			userMessage := m.textarea.Value()
			if strings.TrimSpace(userMessage) == "" {
				// If the message is empty, don't send it.
				m.rpcClient.Logger.Warn("Attempted to send an empty message.")
				return m, nil
			}
			m.rpcClient.Logger.Infof("Sending message: %s", userMessage)

			// Add the user's message with "self" style and send it to the server.
			m.messages = append(m.messages, ChatMessage{
				Sender:  "self",
				Message: userMessage,
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// Send the message to the server.
			err := m.rpcClient.ChatClient.SendUnencryptedMessage(m.ctx, uint32(m.activeUserID), userMessage) // Assuming recipient ID is 1.
			if err != nil {
				m.rpcClient.Logger.Errorf("Failed to send message: %v", err)
				return m, tea.Quit
			}

		case tea.KeyCtrlA:
			// Simulate receiving an "other" message on pressing Ctrl+A (for testing purposes).
			m.messages = append(m.messages, ChatMessage{
				Sender:  "other",
				Message: "This is a reply from the other person.",
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case ReceivedMessage:
		// Handle incoming messages from the channel.
		m.rpcClient.Logger.Infof("Received message from sender %s: %s", msg.Sender, msg.Message)
		m.messages = append(m.messages, ChatMessage{
			Sender:  msg.Sender,
			Message: msg.Message,
		})
		m.viewport.SetContent(m.renderMessages())
		m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, m.listenToMessageChannel() // Continue listening to the message channel.

	case errMsg:
		// Handle errors from the channel.
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

// View renders the chat view.
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
			spacesNeeded := viewportWidth - lipgloss.Width(msgContent) - 2 // Adjust right padding.
			if spacesNeeded > 0 {
				styledMessage = lipgloss.NewStyle().MarginLeft(spacesNeeded).Render(styledMessage)
			}
		case "server":
			continue // Skip rendering server welcome messages.
		default:
			// Render messages from any other sender with a default style
			defaultStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")). // Use a unique color for other usernames.
				Align(lipgloss.Left).
				PaddingLeft(1) // Left padding to differentiate from "self"
			styledMessage = defaultStyle.Render(fmt.Sprintf("%s: %s", msg.Sender, msg.Message))
		}

		// Append the rendered message to the list
		renderedMessages = append(renderedMessages, styledMessage)
	}

	return strings.Join(renderedMessages, "\n")
}

// SetActiveUser sets the active user for the chat and clears the message history.
// func (m *ChatModel) SetActiveUser(userID int32, username string) {
// 	m.rpcClient.Logger.Infof("Setting active user for chat: ID=%d, Username=%s", userID, username)
// 	m.activeUserID = userID
// 	m.activeUsername = username
// 	m.messages = []ChatMessage{} // Clear existing messages when switching users. or load history eventually
// 	content := fmt.Sprintf("No messages yet. Start the conversation with %s", username)
// 	m.viewport.SetContent(content) // Update viewport content.
// }

func (m *ChatModel) SetActiveUser(userID int32, username string) {
	m.rpcClient.Logger.Infof("Setting active user for chat: ID=%d, Username=%s", userID, username)
	m.activeUserID = userID
	m.activeUsername = username
	m.messages = []ChatMessage{} // Clear existing messages when switching users.

	// Fetch chat history between the current user and the active user.
	m.rpcClient.Logger.Infof("Fetching chat history between users: CurrentUserID=%d, ActiveUserID=%d", m.rpcClient.CurrentUserID, userID)
	chatHistory, err := m.rpcClient.ChatClient.Store.GetChatHistoryBetweenUsers(m.rpcClient.CurrentUserID, uint32(userID))
	// chatHistory, err := m.rpcClient.ChatClient.Store.GetAllChatHistory()
	if err != nil {
		m.rpcClient.Logger.Errorf("Failed to get chat history: %v", err)
		m.viewport.SetContent(fmt.Sprintf("Failed to load chat history with %s", username))
		return
	}

	// Load chat history into the model.
	for _, msg := range chatHistory {
		// m.rpcClient.Logger.Info("msg: ", msg)
		sender := "self"
		if msg.SenderID != m.rpcClient.CurrentUserID {
			sender = username
		}
		m.messages = append(m.messages, ChatMessage{
			Sender:  sender,
			Message: msg.Message,
		})
	}

	if len(m.messages) == 0 {
		m.viewport.SetContent(fmt.Sprintf("No messages yet. Start the conversation with %s", username))
	} else {
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
	}
}
