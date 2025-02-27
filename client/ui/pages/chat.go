package pages

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/lib"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
)

// ChatMessage represents a message in the chat.
type ChatMessage struct {
	Sender    string // "self" or "other"
	Message   string
	FileType  string
	FileSize  uint64
	FileName  string
	FileData  []byte
	Timestamp time.Time // Changed from uint64 to time.Time
}

// Define a custom message type for received messages.
type ReceivedMessage struct {
	SenderID  uint32
	Sender    string
	Message   string
	FileType  string
	FileSize  uint64
	FileName  string
	FileData  []byte
	Timestamp time.Time
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
	serverMessages []ChatMessage
}

const gap = "\n\n"

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

				// Decrypt the message and handle any errors.

				// Check if the message is encrypted.
				if msg.EncryptionType == chat.EncryptionType_SIGNAL || msg.EncryptionType == chat.EncryptionType_PREKEY {
					// Decrypt the message using the Signal protocol.
					decryptedBytes, err := m.rpcClient.ChatClient.DecryptMessage(m.ctx, msg)
					decrypted := string(decryptedBytes)
					if err != nil {
						m.rpcClient.Logger.Errorf("Failed to decrypt message: %v", err)
						return errMsg{err}
					}
					// Log the decrypted message and return it as a ReceivedMessage.
					m.rpcClient.Logger.Infof("Received decrypted message from channel: Sender=%s, Message=%s, Status=%s", msg.SenderUsername, decrypted, msg.Status)
					return ReceivedMessage{
						SenderID:  msg.SenderId,
						Sender:    msg.SenderUsername,
						Message:   decrypted,
						FileType:  msg.FileType,
						FileSize:  msg.FileSize,
						FileName:  msg.FileName,
						Timestamp: time.Now().UTC(),
					}
				} else {

					// Log the received message and return it as a ReceivedMessage.
					m.rpcClient.Logger.Infof("Received message from channel: Sender=%s, Message=%s, Status=%s", msg.SenderUsername, string(msg.EncryptedMessage), msg.Status)
					return ReceivedMessage{
						SenderID:  msg.SenderId,
						Sender:    msg.SenderUsername,
						Message:   string(msg.EncryptedMessage),
						FileType:  msg.FileType,
						FileSize:  msg.FileSize,
						FileName:  msg.FileName,
						Timestamp: time.Now().UTC(),
					}
				}
			}
		}
	}
}

func (m ChatModel) OpenFileMenu(originalServerMessages []ChatMessage, originalActiveUser int32, originalActiveUsername string) tea.Cmd {
	return func() tea.Msg {
		return OpenFileMenuMsg{originalServerMessages, originalActiveUser, originalActiveUsername}
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

			userMessage := m.textarea.Value()
			if strings.TrimSpace(userMessage) == "" {
				// If the message is empty, don't send it.
				m.rpcClient.Logger.Warn("Attempted to send an empty message.")
				return m, nil
			}

			if userMessage == "/open" {
				return m, m.OpenFileMenu(m.serverMessages, m.activeUserID, m.activeUsername)

			}

			// Check for file sending command. For example: "/file /path/to/file.jpg"
			if strings.HasPrefix(userMessage, "/file ") {
				filePath := strings.TrimSpace(strings.TrimPrefix(userMessage, "/file "))
				m.rpcClient.Logger.Infof("Sending file: %s", filePath)

				// Read the file.
				fileData, err := os.ReadFile(filePath)
				if err != nil {
					m.rpcClient.Logger.Errorf("Failed to read file %s: %v", filePath, err)
					// Optionally show an error in the chat.
					m.messages = append(m.messages, ChatMessage{
						Sender:   "self",
						Message:  fmt.Sprintf("Error reading file: %s", filePath),
						FileType: "text",
						FileSize: 0,
						FileName: filePath,
					})
					m.viewport.SetContent(m.renderMessages())
					m.textarea.Reset()
					m.viewport.GotoBottom()
					return m, nil
				}

				// Determine file type (simple logic based on file extension).
				fileName := filepath.Base(filePath)
				ext := strings.ToLower(filepath.Ext(fileName))
				fileType := "file"
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" || ext == ".tiff" || ext == ".ico" || ext == ".webp" {
					fileType = "image"
				} else if ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".mkv" || ext == ".flv" || ext == ".wmv" {
					fileType = "video"
				} else {
					fileType = "file"
				}

				// Send the file message. For images, your store logic will treat it differently.
				err = m.rpcClient.ChatClient.SendMessage(m.ctx, uint32(m.activeUserID), m.rpcClient.CurrentDeviceID, fileData, &lib.SendMessageOptions{
					FileType: fileType,
					FileSize: uint64(len(fileData)),
					FileName: fileName,
				})
				if err != nil {
					m.rpcClient.Logger.Errorf("Failed to send file: %v", err)
					m.messages = append(m.messages, ChatMessage{
						Sender:   "self",
						Message:  fmt.Sprintf("Error sending file: %s of type %s", fileName, fileType),
						FileType: "text",
						FileSize: 0,
						FileName: fileName,
					})
				} else {
					// Append a representation of the sent file to the chat history.
					m.messages = append(m.messages, ChatMessage{
						Sender:   "self",
						Message:  fmt.Sprintf("[sent file] %s", fileName),
						FileType: fileType,
						FileSize: uint64(len(fileData)),
						FileName: fileName,
						FileData: fileData,
					})
				}
				m.viewport.SetContent(m.renderMessages())
				m.textarea.Reset()
				m.viewport.GotoBottom()
				return m, nil
			}

			// Prevent sending messages if activeUserID is not set.
			if m.activeUserID == 0 {
				m.serverMessages = append(m.serverMessages, ChatMessage{
					Sender:   "server",
					Message:  "Please select a user (that is not me) to chat with. Add friends to your friends list to chat with them.",
					FileType: "text",
					FileSize: 0,
					FileName: "",
				})
				m.messages = m.serverMessages
				m.viewport.SetContent(m.renderMessages())
				m.textarea.Reset()
				m.viewport.GotoBottom()
				return m, nil
			}

			m.rpcClient.Logger.Infof("Sending message: %s", userMessage)

			// Add the user's message with "self" style and send it to the server.
			m.messages = append(m.messages, ChatMessage{
				Sender:    "self",
				Message:   userMessage,
				FileType:  "text",
				FileSize:  uint64(len([]byte(userMessage))),
				FileName:  "",
				Timestamp: time.Now().UTC(),
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// Send the message to the server as a text message.
			err := m.rpcClient.ChatClient.SendMessage(m.ctx, uint32(m.activeUserID), m.rpcClient.CurrentDeviceID, []byte(userMessage), &lib.SendMessageOptions{
				FileType: "text",
				FileSize: uint64(len([]byte(userMessage))),
				FileName: "",
			})
			if err != nil {
				m.rpcClient.Logger.Errorf("Failed to send message: %v", err)
				return m, tea.Quit
			}

		case tea.KeyCtrlA:
			// Simulate receiving an "other" message on pressing Ctrl+A (for testing purposes).
			m.messages = append(m.messages, ChatMessage{
				Sender:   "other",
				Message:  "This is a reply from the other person.",
				FileType: "text",
				FileSize: 0,
				FileName: "",
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case ReceivedMessage:
		// If the message is from the server, add it to the server messages.
		if m.activeUserID == 0 && msg.SenderID == 0 {
			m.serverMessages = append(m.serverMessages, ChatMessage{
				Sender:   msg.Sender,
				Message:  msg.Message,
				FileType: "text",
			})
			m.messages = m.serverMessages
			m.viewport.SetContent(m.renderMessages())
			return m, m.listenToMessageChannel()
		}

		// Check if the message is from the active user
		if msg.SenderID == uint32(m.activeUserID) {
			m.rpcClient.Logger.Infof("Processing message from sender %s: %s", msg.Sender, msg.Message)
			m.messages = append(m.messages, ChatMessage{
				Sender:    msg.Sender,
				Message:   msg.Message,
				FileType:  msg.FileType,
				FileSize:  msg.FileSize,
				FileName:  msg.FileName,
				FileData:  msg.FileData,
				Timestamp: msg.Timestamp,
			})
			m.viewport.SetContent(m.renderMessages())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
		return m, m.listenToMessageChannel()

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
		// "%s\n\n%s",
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}

// renderMessages iterates over the chat messages and applies styles based on sender.
func (m ChatModel) renderMessages() string {
	var renderedMessages []string

	for _, msg := range m.messages {
		var styledMessage string

		// Format timestamp
		timeStr := ""
		if !msg.Timestamp.IsZero() {
			timeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")). // Subtle gray color for timestamp
				Italic(true)
			timeStr = timeStyle.Render(fmt.Sprintf("[%s] ", msg.Timestamp.Format("15:04")))
		}

		// Determine sender prefix styling
		var senderPrefix string
		if msg.Sender == "self" {
			senderPrefix = m.selfStyle.Render("You: ")
		} else {
			defaultStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")). // Unique color for other usernames
				Align(lipgloss.Left).
				PaddingLeft(1)
			senderPrefix = defaultStyle.Render(fmt.Sprintf("%s: ", msg.Sender))
		}

		// Check if this message represents a file (non-text)
		if msg.FileType != "text" {
			fileStyle := lipgloss.NewStyle().
				Bold(true).
				Background(lipgloss.Color("21")). // blue background
				Foreground(lipgloss.Color("15"))  // white text

			fileDetails := fmt.Sprintf("%s (%s, %d bytes)", msg.FileName, msg.FileType, msg.FileSize)
			styledMessage = timeStr + senderPrefix + fileStyle.Render(fileDetails)
		} else {
			// Normal text message
			styledMessage = timeStr + senderPrefix + msg.Message
		}

		renderedMessages = append(renderedMessages, styledMessage)
	}

	// Wrap the content to fit the viewport
	return lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(renderedMessages, "\n"))
}

func (m *ChatModel) SetActiveUser(userID int32, username string) {
	m.rpcClient.Logger.Infof("Setting active user for chat: ID=%d, Username=%s", userID, username)
	m.activeUserID = userID
	m.activeUsername = username
	m.messages = []ChatMessage{} // Clear existing messages when switching users.

	// Special case for server
	if userID == 0 {
		m.messages = m.serverMessages
		m.viewport.SetContent(m.renderMessages())
		return
	}

	// Fetch chat history between the current user and the active user.
	m.rpcClient.Logger.Infof("Fetching chat history between users: CurrentUserID=%d, ActiveUserID=%d", m.rpcClient.CurrentUserID, userID)
	chatHistory, err := m.rpcClient.ChatClient.Store.GetChatHistory(m.rpcClient.CurrentUserID, uint32(userID))
	if err != nil {
		m.rpcClient.Logger.Errorf("Failed to get chat history: %v", err)
		m.viewport.SetContent(fmt.Sprintf("Failed to load chat history with %s", username))
		return
	}

	// Load chat history into the model.
	for _, msg := range chatHistory {
		sender := "self"
		if msg.SenderID != m.rpcClient.CurrentUserID {
			sender = username
		}
		m.messages = append(m.messages, ChatMessage{
			Sender:    sender,
			Message:   msg.Message,
			FileType:  msg.FileType,
			FileSize:  msg.FileSize,
			FileName:  msg.FileName,
			FileData:  msg.Media,
			Timestamp: msg.Timestamp,
		})
	}

	if len(m.messages) == 0 {
		m.viewport.SetContent(fmt.Sprintf("No messages yet. Start the conversation with %s", username))
	} else {
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
	}
}
